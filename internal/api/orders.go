package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/constants"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
)

// Setup order-related routes
func (sv *Server) addOrderRoutes(rg *gin.RouterGroup) {
	orders := rg.Group("/orders", authenticateMiddleware(sv.tokenGenerator))
	{
		orders.GET("", sv.getOrders)
		orders.GET(":id", sv.getOrderDetail)
		orders.PUT(":id/confirm-received", sv.confirmOrderPayment)
		orders.POST(":id/cancel", sv.adminCancelOrder)
	}
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
func (sv *Server) getOrders(c *gin.Context) {
	tokenPayload, ok := c.MustGet(constants.AuthPayLoad).(*auth.TokenPayload)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.CreateErr(UnauthorizedCode, errors.New("authorization payload is not provided")))
		return
	}

	var orderListQuery models.OrderListQuery
	if err := c.ShouldBindQuery(&orderListQuery); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	dbParams := repository.GetOrdersParams{
		Limit:  orderListQuery.PageSize,
		Offset: (orderListQuery.Page - 1) * orderListQuery.PageSize,
	}

	if tokenPayload.RoleCode != "admin" {
		dbParams.UserID = utils.GetPgTypeUUID(tokenPayload.UserID)
	}

	fetchedOrderRows, err := sv.repo.GetOrders(c, dbParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	count, err := sv.repo.CountOrders(c, repository.CountOrdersParams{UserID: utils.GetPgTypeUUID(tokenPayload.UserID)})

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
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

	c.JSON(http.StatusOK, dto.CreateDataResp(c, orderResponses, dto.CreatePagination(orderListQuery.Page, orderListQuery.PageSize, count), nil))
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
func (sv *Server) getOrderDetail(c *gin.Context) {
	var params models.UriIDParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	var resp *dto.OrderDetail = nil

	if err := sv.cacheSrv.Get(c, "order_detail:"+params.ID, &resp); err == nil {
		if resp != nil {
			c.JSON(http.StatusOK, dto.CreateDataResp(c, resp, nil, nil))
			return
		}
	}

	order, err := sv.repo.GetOrder(c, uuid.MustParse(params.ID))
	if err != nil {
		if err == repository.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, err))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
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

	paymentInfo, err := sv.repo.GetPaymentByOrderID(c, order.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	var apiErr *dto.ApiError = nil

	if paymentInfo.Status == repository.PaymentStatusPending &&
		paymentInfo.PaymentIntentID != nil {
		paymentDetail, err := sv.paymentSrv.GetPayment(c, *paymentInfo.PaymentIntentID, *paymentInfo.Gateway)
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

	discountRows, err := sv.repo.GetOrderDiscounts(c, order.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
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

	orderItemRows, err := sv.repo.GetOrderItems(c, order.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
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

	c.JSON(http.StatusOK, dto.CreateDataResp(c, resp, nil, apiErr))
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
func (sv *Server) confirmOrderPayment(c *gin.Context) {
	tokenPayload, _ := c.MustGet(constants.AuthPayLoad).(*auth.TokenPayload)
	var params models.UriIDParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	order, err := sv.repo.GetOrder(c, uuid.MustParse(params.ID))
	if err != nil {
		if err == repository.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, err))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	if order.UserID != tokenPayload.UserID {
		c.JSON(http.StatusForbidden, dto.CreateErr(PermissionDeniedCode, errors.New("you do not have permission to access this order")))
		return
	}
	if order.Status != repository.OrderStatusDelivered {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidPaymentCode, errors.New("order cannot be confirmed")))
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
	_, err = sv.repo.UpdateOrder(c, orderUpdateParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	var apiErr *dto.ApiError
	if err := sv.cacheSrv.Delete(c, "order_detail:"+params.ID); err != nil {
		log.Err(err).Msg("failed to delete order detail cache")
		apiErr = &dto.ApiError{
			Code:    InternalServerErrorCode,
			Details: "failed to delete order detail cache",
			Stack:   err,
		}
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, true, nil, apiErr))
}
