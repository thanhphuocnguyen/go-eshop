package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
)

// Setup order-related routes
func (sv *Server) addOrderRoutes(r chi.Router) {
	r.Route("/orders", func(r chi.Router) {
		// Apply authentication middleware
		r.Use(func(h http.Handler) http.Handler {
			return authenticateMiddleware(h, sv.tokenGenerator)
		})

		r.Get("/", sv.getOrders)
		r.Get("/{id}", sv.getOrderDetail)
		r.Put("/{id}/confirm-received", sv.confirmOrderPayment)
		r.Post("/{id}/cancel", sv.adminCancelOrder)
	})
}

// @Summary List orders
// @Description List orders of the current user
// @Tags orders
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Param status query string false "Filter by status"
// @Param payment_status query string false "Filter by payment status"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[[]OrderListResponse]
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /orders [get]
func (sv *Server) getOrders(w http.ResponseWriter, r *http.Request) {
	tokenPayload, ok := r.Context().Value("auth").(*auth.TokenPayload)
	if !ok {
		RespondUnauthorized(w, UnauthorizedCode, errors.New("authorization payload is not provided"))
		return
	}

	paginationQuery := ParsePaginationQuery(r)

	// Parse additional query parameters
	queryParams := r.URL.Query()
	var orderListQuery models.OrderListQuery
	orderListQuery.Page = paginationQuery.Page
	orderListQuery.PageSize = paginationQuery.PageSize

	if status := queryParams.Get("status"); status != "" {
		orderListQuery.Status = &status
	}
	if paymentStatus := queryParams.Get("payment_status"); paymentStatus != "" {
		orderListQuery.PaymentStatus = &paymentStatus
	}

	validate := validator.New()
	if err := validate.Struct(&orderListQuery); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	dbParams := repository.GetOrdersParams{
		Limit:  orderListQuery.PageSize,
		Offset: (orderListQuery.Page - 1) * orderListQuery.PageSize,
	}

	if tokenPayload.RoleCode != "admin" {
		dbParams.UserID = utils.GetPgTypeUUID(tokenPayload.UserID)
	}

	fetchedOrderRows, err := sv.repo.GetOrders(r.Context(), dbParams)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	count, err := sv.repo.CountOrders(r.Context(), repository.CountOrdersParams{UserID: utils.GetPgTypeUUID(tokenPayload.UserID)})

	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	var orderResponses []dto.OrderListItem
	for _, aggregated := range fetchedOrderRows {
		// Convert PaymentStatus interface{} to PaymentStatus type
		paymentStatus := repository.PaymentStatusPending
		if aggregated.PaymentStatus.Valid {
			paymentStatus = aggregated.PaymentStatus.PaymentStatus
		}

		total, _ := aggregated.TotalPrice.Float64Value()
		orderResponses = append(orderResponses, dto.OrderListItem{
			ID:            aggregated.ID,
			Total:         total.Float64,
			TotalItems:    int32(aggregated.TotalItems),
			Status:        aggregated.Status,
			PaymentStatus: paymentStatus,
			CreatedAt:     aggregated.CreatedAt.UTC(),
			UpdatedAt:     aggregated.UpdatedAt.UTC(),
		})
	}

	RespondSuccessWithPagination(w, r, orderResponses, dto.CreatePagination(orderListQuery.Page, orderListQuery.PageSize, count))
}

// @Summary Get order detail
// @Description Get order detail by order ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Security BearerAuth
// @Success 200 {object} OrderDetailResponse
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /order/{orderId} [get]
func (sv *Server) getOrderDetail(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		RespondBadRequest(w, InvalidBodyCode, errors.New("id parameter is required"))
		return
	}

	var resp *dto.OrderDetail = nil

	if err := sv.cacheSrv.Get(r.Context(), "order_detail:"+idParam, &resp); err == nil {
		if resp != nil {
			RespondSuccess(w, r, resp)
			return
		}
	}

	order, err := sv.repo.GetOrder(r.Context(), uuid.MustParse(idParam))
	if err != nil {
		if err == repository.ErrRecordNotFound {
			RespondNotFound(w, NotFoundCode, err)
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	total, _ := order.TotalPrice.Float64Value()
	paymentAmount, _ := order.PaymentAmount.Float64Value()
	resp = &dto.OrderDetail{
		ID:            order.ID,
		Total:         total.Float64,
		CustomerName:  order.CustomerName,
		CustomerEmail: order.CustomerEmail,
		Status:        order.Status,
		ShippingInfo:  order.ShippingAddress,
		CreatedAt:     order.CreatedAt.UTC(),
		UpdatedAt:     order.UpdatedAt.UTC(),
		LineItems:     []dto.LineItem{},
	}

	if order.PaymentID.Valid {
		pmId, _ := uuid.FromBytes(order.PaymentID.Bytes[:])
		resp.PaymentInfo = dto.PaymentInfo{
			ID:       pmId.String(),
			Amount:   paymentAmount.Float64,
			IntendID: order.PaymentIntentID,
			GateWay:  order.Gateway,
			Method:   *order.PaymentMethod,
			Status:   string(order.PaymentStatus.PaymentStatus),
		}
	}

	paymentInfo, err := sv.repo.GetPaymentByOrderID(r.Context(), order.ID)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	var apiErr *dto.ApiError = nil

	if paymentInfo.Status == repository.PaymentStatusPending &&
		paymentInfo.PaymentIntentID != nil {
		paymentDetail, err := sv.paymentSrv.GetPayment(r.Context(), *paymentInfo.PaymentIntentID, *paymentInfo.Gateway)
		if err != nil {
			log.Err(err).Msg("failed to get payment intent")
			apiErr = &dto.ApiError{
				Code:    InternalServerErrorCode,
				Details: "failed to get payment intent",
				Stack:   err,
			}
		} else {
			resp.PaymentInfo.ClientSecret = &paymentDetail.ClientSecret
		}
	}

	discountRows, err := sv.repo.GetOrderDiscounts(r.Context(), order.ID)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	var discounts []dto.OrderDiscount
	for _, discount := range discountRows {
		discountValue, _ := discount.DiscountValue.Float64Value()
		discounts = append(discounts, dto.OrderDiscount{
			ID:            discount.ID.String(),
			Code:          discount.Code,
			Description:   discount.Description,
			DiscountType:  string(discount.DiscountType),
			DiscountValue: discountValue.Float64,
		})
	}
	resp.Discounts = discounts

	orderItemRows, err := sv.repo.GetOrderItems(r.Context(), order.ID)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	lineItems := make([]dto.LineItem, 0, len(orderItemRows))
	for _, item := range orderItemRows {
		lineTotal, _ := item.LineTotalSnapshot.Float64Value()
		var discountAmount float64
		if item.DiscountedPrice.Valid {
			discount, _ := item.DiscountedPrice.Float64Value()
			discountAmount = discount.Float64
		}
		var price float64
		if item.PricePerUnitSnapshot.Valid {
			p, _ := item.PricePerUnitSnapshot.Float64Value()
			price = p.Float64
		}
		itemResp := dto.LineItem{
			ID:                 item.ID.String(),
			VariantID:          item.VariantID.String(),
			Name:               item.ProductName,
			ImageUrl:           item.ImageUrl,
			LineTotal:          lineTotal.Float64,
			DiscountAmount:     discountAmount,
			Price:              price,
			AttributesSnapshot: item.AttributesSnapshot,
			Quantity:           item.Quantity,
			CreatedAt:          item.CreatedAt.UTC(),
			UpdatedAt:          item.UpdatedAt.UTC(),
		}
		lineItems = append(lineItems, itemResp)
	}

	resp.LineItems = lineItems

	response := dto.CreateDataResp(r.Context(), resp, nil, apiErr)
	RespondJSON(w, http.StatusOK, response)
}

// @Summary confirm received order payment info
// @Description confirm received order payment info
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[bool]
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /order/{orderId}/confirm-received [put]
func (sv *Server) confirmOrderPayment(w http.ResponseWriter, r *http.Request) {
	tokenPayload, ok := r.Context().Value("auth").(*auth.TokenPayload)
	if !ok {
		RespondUnauthorized(w, UnauthorizedCode, errors.New("authorization payload is not provided"))
		return
	}

	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		RespondBadRequest(w, InvalidBodyCode, errors.New("id parameter is required"))
		return
	}

	order, err := sv.repo.GetOrder(r.Context(), uuid.MustParse(idParam))
	if err != nil {
		if err == repository.ErrRecordNotFound {
			RespondNotFound(w, NotFoundCode, err)
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	if order.UserID != tokenPayload.UserID {
		RespondForbidden(w, PermissionDeniedCode, errors.New("you do not have permission to access this order"))
		return
	}
	if order.Status != repository.OrderStatusDelivered {
		RespondBadRequest(w, InvalidPaymentCode, errors.New("order cannot be confirmed"))
		return
	}

	orderUpdateParams := repository.UpdateOrderParams{
		ID: order.ID,
		Status: repository.NullOrderStatus{
			OrderStatus: repository.OrderStatusDelivered,
			Valid:       true,
		},
		ConfirmedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	}
	_, err = sv.repo.UpdateOrder(r.Context(), orderUpdateParams)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	var apiErr *dto.ApiError
	if err := sv.cacheSrv.Delete(r.Context(), "order_detail:"+idParam); err != nil {
		log.Err(err).Msg("failed to delete order detail cache")
		apiErr = &dto.ApiError{
			Code:    InternalServerErrorCode,
			Details: "failed to delete order detail cache",
			Stack:   err,
		}
	}

	response := dto.CreateDataResp(r.Context(), true, nil, apiErr)
	RespondJSON(w, http.StatusOK, response)
}
