package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v81"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
	"github.com/thanhphuocnguyen/go-eshop/pkg/payment"
)

// ---------------------------------------------- API Models ----------------------------------------------
type OrderListQuery struct {
	PaginationQueryParams
	Status        *string `form:"status,omitempty" binding:"omitempty,oneof=pending confirmed delivering delivered completed cancelled refunded"`
	PaymentStatus *string `form:"payment_status,omitempty" binding:"omitempty,oneof=pending succeeded failed cancelled refunded"`
}

type OrderIDParams struct {
	ID string `uri:"id" binding:"required,min=1"`
}

type OrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending confirmed delivering delivered completed"`
	Reason string `json:"reason,omitempty"`
}

type OrderItemResponse struct {
	ID                 string                             `json:"id"`
	Name               string                             `json:"name"`
	ImageUrl           *string                            `json:"image_url"`
	AttributesSnapshot []repository.AttributeDataSnapshot `json:"attributes_snapshot"`
	LineTotal          float64                            `json:"line_total"`
	Quantity           int16                              `json:"quantity"`
}
type PaymentInfo struct {
	ID           string  `json:"id"`
	RefundID     *string `json:"refund_id"`
	Amount       float64 `json:"amount"`
	IntendID     *string `json:"intent_id"`
	ClientSecret *string `json:"client_secret"`
	GateWay      *string `json:"gateway"`
	Method       string  `json:"method"`
	Status       string  `json:"status"`
}

type OrderDetailResponse struct {
	ID            uuid.UUID                          `json:"id"`
	Total         float64                            `json:"total"`
	Status        repository.OrderStatus             `json:"status"`
	CustomerName  string                             `json:"customer_name"`
	CustomerEmail string                             `json:"customer_email"`
	PaymentInfo   *PaymentInfo                       `json:"payment_info,omitempty"`
	ShippingInfo  repository.ShippingAddressSnapshot `json:"shipping_info"`
	Products      []OrderItemResponse                `json:"products"`
	CreatedAt     time.Time                          `json:"created_at"`
}
type OrderListResponse struct {
	ID            uuid.UUID                `json:"id"`
	Total         float64                  `json:"total"`
	TotalItems    int32                    `json:"total_items"`
	Status        repository.OrderStatus   `json:"status"`
	PaymentStatus repository.PaymentStatus `json:"payment_status"`
	CustomerName  string                   `json:"customer_name"`
	CustomerEmail string                   `json:"customer_email"`
	CreatedAt     time.Time                `json:"created_at"`
	UpdatedAt     time.Time                `json:"updated_at"`
}

type RefundOrderRequest struct {
	Reason string `json:"reason" binding:"required,oneof=defective damaged fraudulent requested_by_customer"`
}
type CancelOrderRequest struct {
	Reason string `json:"reason" binding:"required,oneof=duplicate fraudulent requested_by_customer abandoned"`
}

//---------------------------------------------- API Handlers ----------------------------------------------

// @Summary List orders
// @Description List orders of the current user
// @Tags orders
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param status query string false "Filter by status"
// @Param payment_status query string false "Filter by payment status"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[[]OrderListResponse]
// @Failure 400 {object} ApiResponse[[]OrderListResponse]
// @Failure 401 {object} ApiResponse[[]OrderListResponse]
// @Failure 500 {object} ApiResponse[[]OrderListResponse]
// @Router /order/list [get]
func (sv *Server) getOrdersHandler(c *gin.Context) {
	tokenPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusUnauthorized, createErrorResponse[[]OrderListResponse](UnauthorizedCode, "", errors.New("authorization payload is not provided")))
		return
	}

	var orderListQuery OrderListQuery
	if err := c.ShouldBindQuery(&orderListQuery); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[[]OrderListResponse](InvalidBodyCode, "", err))
		return
	}

	dbParams := repository.GetOrdersParams{
		// CustomerID: utils.GetPgTypeUUID(tokenPayload.UserID),
		Limit:  orderListQuery.PageSize,
		Offset: (orderListQuery.Page - 1) * orderListQuery.PageSize,
	}

	fetchedOrderRows, err := sv.repo.GetOrders(c, dbParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[[]OrderListResponse](InternalServerErrorCode, "", err))
		return
	}

	count, err := sv.repo.CountOrders(c, repository.CountOrdersParams{CustomerID: utils.GetPgTypeUUID(tokenPayload.UserID)})

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[[]OrderListResponse](InternalServerErrorCode, "", err))
		return
	}

	var orderResponses []OrderListResponse
	for _, aggregated := range fetchedOrderRows {
		total, _ := aggregated.TotalPrice.Float64Value()
		orderResponses = append(orderResponses, OrderListResponse{
			ID:            aggregated.ID,
			Total:         total.Float64,
			TotalItems:    int32(aggregated.TotalItems),
			Status:        aggregated.Status,
			PaymentStatus: aggregated.PaymentStatus.PaymentStatus,
			CreatedAt:     aggregated.CreatedAt.UTC(),
			UpdatedAt:     aggregated.UpdatedAt.UTC(),
		})
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, orderResponses, "success", &Pagination{
		Total: count, Page: orderListQuery.Page,
		PageSize:        orderListQuery.PageSize,
		TotalPages:      count / int64(orderListQuery.PageSize),
		HasNextPage:     count > int64(orderListQuery.PageSize),
		HasPreviousPage: orderListQuery.Page > 1,
	}, nil))
}

// @Summary Get order detail
// @Description Get order detail by order ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Security BearerAuth
// @Success 200 {object} OrderDetailResponse
// @Failure 400 {object} OrderDetailResponse
// @Failure 401 {object} OrderDetailResponse
// @Failure 500 {object} OrderDetailResponse
// @Router /order/{order_id} [get]
func (sv *Server) getOrderDetailHandler(c *gin.Context) {
	var params OrderIDParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[OrderDetailResponse](InvalidBodyCode, "", err))
		return
	}

	order, err := sv.repo.GetOrder(c, uuid.MustParse(params.ID))
	if err != nil {
		if err == repository.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, createErrorResponse[OrderDetailResponse](NotFoundCode, "", err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[OrderDetailResponse](InternalServerErrorCode, "", err))
		return
	}

	total, _ := order.TotalPrice.Float64Value()

	resp := OrderDetailResponse{
		ID:            order.ID,
		Total:         total.Float64,
		CustomerName:  order.CustomerName,
		CustomerEmail: order.CustomerEmail,
		PaymentInfo:   nil,
		Status:        order.Status,
		ShippingInfo:  order.ShippingAddress,
		CreatedAt:     order.CreatedAt.UTC(),
		Products:      []OrderItemResponse{},
	}

	orderItemRows, err := sv.repo.GetOrderProducts(c, order.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[OrderDetailResponse](InternalServerErrorCode, "", err))
		return
	}

	paymentInfo, err := sv.repo.GetPaymentByOrderID(c, order.ID)

	var apiErr *ApiError = nil

	if err != nil {
		log.Err(err).Msg("Order does not have payment info")
		apiErr = &ApiError{
			Code:    InternalServerErrorCode,
			Details: "failed to get payment info",
			Stack:   err.Error(),
		}
	}
	if err == nil {
		amount, _ := paymentInfo.Amount.Float64Value()
		resp.PaymentInfo = &PaymentInfo{
			ID:       paymentInfo.ID.String(),
			RefundID: paymentInfo.RefundID,
			GateWay:  (*string)(&paymentInfo.PaymentGateway.PaymentGateway),
			Method:   string(paymentInfo.PaymentMethod),
			Status:   string(paymentInfo.Status),
			Amount:   amount.Float64,
		}
		if paymentInfo.Status == repository.PaymentStatusPending && paymentInfo.GatewayPaymentIntentID != nil && paymentInfo.PaymentMethod == repository.PaymentMethodStripe {
			resp.PaymentInfo.IntendID = paymentInfo.GatewayPaymentIntentID
			stripeInstance, err := payment.NewStripePayment(sv.config.StripeSecretKey)
			if err != nil {
				c.JSON(http.StatusInternalServerError, createErrorResponse[PaymentResponse](InternalServerErrorCode, "", err))
				return
			}
			sv.paymentCtx.SetStrategy(stripeInstance)
			paymentDetail, err := sv.paymentCtx.GetPaymentObject(*paymentInfo.GatewayPaymentIntentID)
			if err != nil {
				log.Err(err).Msg("failed to get payment intent")
				apiErr = &ApiError{
					Code:    InternalServerErrorCode,
					Details: "failed to get payment intent",
					Stack:   err.Error(),
				}
			} else {
				stripeObject, _ := paymentDetail.(*stripe.PaymentIntent)
				resp.PaymentInfo.ClientSecret = &stripeObject.ClientSecret
			}
		}
	}

	orderItems := make([]OrderItemResponse, 0, len(orderItemRows))

	for _, item := range orderItemRows {
		lineTotal, _ := item.LineTotalSnapshot.Float64Value()
		orderItems = append(orderItems, OrderItemResponse{
			ID:                 item.VariantID.String(),
			Name:               item.ProductName,
			ImageUrl:           item.ImageUrl,
			LineTotal:          lineTotal.Float64,
			AttributesSnapshot: item.AttributesSnapshot,
			Quantity:           item.Quantity,
		})
	}

	resp.Products = orderItems

	c.JSON(http.StatusOK, createSuccessResponse(c, resp, "success", nil, apiErr))
}

// @Summary confirm received order payment info
// @Description confirm received order payment info
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[bool]
// @Failure 400 {object} ApiResponse[gin.H]
// @Failure 401 {object} ApiResponse[gin.H]
// @Failure 500 {object} ApiResponse[gin.H]
// @Router /order/{order_id}/confirm-received [put]
func (sv *Server) confirmOrderPayment(c *gin.Context) {
	tokenPayload, _ := c.MustGet(authorizationPayload).(*auth.Payload)
	var params OrderIDParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[gin.H](InvalidBodyCode, "", err))
		return
	}
	order, err := sv.repo.GetOrder(c, uuid.MustParse(params.ID))
	if err != nil {
		if err == repository.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, createErrorResponse[gin.H](NotFoundCode, "", err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "", err))
		return
	}
	if order.CustomerID != tokenPayload.UserID {
		c.JSON(http.StatusForbidden, createErrorResponse[gin.H](PermissionDeniedCode, "", errors.New("You do not have permission to access this order")))
		return
	}
	if order.Status != repository.OrderStatusDelivered {
		c.JSON(http.StatusBadRequest, createErrorResponse[gin.H](InvalidPaymentCode, "", errors.New("order cannot be confirmed")))
		return
	}

	orderUpdateParams := repository.UpdateOrderParams{
		ID: order.ID,
		Status: repository.NullOrderStatus{
			OrderStatus: repository.OrderStatusCompleted,
			Valid:       true,
		},
		ConfirmedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	}
	_, err = sv.repo.UpdateOrder(c, orderUpdateParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
		return
	}
	c.JSON(http.StatusOK, createSuccessResponse(c, true, "success", nil, nil))
}

// @Router /order/{order_id}/confirm_payment [put]
// @Summary Cancel order
// @Description Cancel order by order ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[OrderListResponse]
// @Failure 400 {object} ApiResponse[OrderListResponse]
// @Failure 401 {object} ApiResponse[OrderListResponse]
// @Failure 500 {object} ApiResponse[OrderListResponse]
// @Router /order/{order_id}/cancel [put]
func (sv *Server) cancelOrder(c *gin.Context) {
	tokenPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusUnauthorized, createErrorResponse[OrderListResponse](UnauthorizedCode, "", errors.New("authorization payload is not provided")))
		return
	}
	user, err := sv.repo.GetUserByID(c, tokenPayload.UserID)
	if err != nil {
		if err == repository.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, createErrorResponse[OrderListResponse](NotFoundCode, "", err))
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[OrderListResponse](InternalServerErrorCode, "", err))
		return
	}
	var params OrderIDParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[OrderListResponse](InvalidBodyCode, "", err))
		return
	}
	var req CancelOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[OrderListResponse](InvalidBodyCode, "", err))
		return
	}
	order, err := sv.repo.GetOrder(c, uuid.MustParse(params.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[OrderListResponse](InternalServerErrorCode, "", err))
		return
	}

	if order.CustomerID != tokenPayload.UserID && user.Role != repository.UserRoleAdmin {
		c.JSON(http.StatusForbidden, createErrorResponse[OrderListResponse](PermissionDeniedCode, "", errors.New("You do not have permission to access this order")))
		return
	}

	// if order status is not pending or user is not admin
	if order.Status != repository.OrderStatusPending || user.Role != repository.UserRoleAdmin {
		c.JSON(http.StatusBadRequest, createErrorResponse[OrderListResponse](PermissionDeniedCode, "", errors.New("order cannot be cancelled")))
		return
	}

	var reason payment.CancelReason
	switch req.Reason {
	case "duplicate":
		reason = payment.CancelReasonDuplicate
	case "fraudulent":
		reason = payment.CancelReasonFraudulent
	case "abandoned":
		reason = payment.CancelReasonAbandoned
	case "requested_by_customer":
		reason = payment.CancelReasonRequestedByCustomer
	default:
		reason = payment.CancelReasonRequestedByCustomer
	}

	// if order
	order, err = sv.repo.CancelOrderTx(c, repository.CancelOrderTxArgs{
		OrderID: uuid.MustParse(params.ID),
		CancelPaymentFromGateway: func(paymentID string, gateway repository.PaymentGateway) error {
			switch gateway {
			case repository.PaymentGatewayStripe:
				stripeInstance, err := payment.NewStripePayment(sv.config.StripeSecretKey)
				if err != nil {
					c.JSON(http.StatusInternalServerError, createErrorResponse[OrderListResponse](InternalServerErrorCode, "", err))
					return err
				}
				sv.paymentCtx.SetStrategy(stripeInstance)
				_, err = sv.paymentCtx.CancelPayment(paymentID, reason)
				return err
			default:
				return nil
			}
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[OrderListResponse](repository.ErrDeadlockDetected.InternalQuery, "", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, order, "success", nil, nil))
}

// @Summary Change order status
// @Description Change order status by order ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param status body string true "Status"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[OrderListResponse]
// @Failure 400 {object} ApiResponse[OrderListResponse]
// @Failure 401 {object} ApiResponse[OrderListResponse]
// @Failure 500 {object} ApiResponse[OrderListResponse]
// @Router /order/{order_id}/status [put]
func (sv *Server) changeOrderStatus(c *gin.Context) {
	var params OrderIDParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[OrderListResponse](InvalidBodyCode, "", err))
		return
	}
	var req OrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[OrderListResponse](InvalidBodyCode, "", err))
		return
	}
	order, err := sv.repo.GetOrder(c, uuid.MustParse(params.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[OrderListResponse](NotFoundCode, "", err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[OrderListResponse](InternalServerErrorCode, "", err))
		return
	}
	if order.Status == repository.OrderStatusCompleted || order.Status == repository.OrderStatusCancelled || order.Status == repository.OrderStatusRefunded {
		c.JSON(http.StatusBadRequest, createErrorResponse[OrderListResponse](InvalidPaymentCode, "", errors.New("order cannot be changed")))
		return
	}

	status := repository.OrderStatus(req.Status)

	updateParams := repository.UpdateOrderParams{
		ID: order.ID,
		Status: repository.NullOrderStatus{
			OrderStatus: status,
			Valid:       true,
		},
	}
	if status == repository.OrderStatusConfirmed {
		updateParams.ConfirmedAt = pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		}
	}
	if status == repository.OrderStatusDelivering {
		updateParams.DeliveredAt = pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		}
	}

	rs, err := sv.repo.UpdateOrder(c, updateParams)

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[OrderListResponse](InternalServerErrorCode, "", err))
		return
	}
	c.JSON(http.StatusOK, createSuccessResponse(c, rs, "success", nil, nil))
}

// @Summary Refund order
// @Description Refund order by order ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[OrderListResponse]
// @Failure 400 {object} ApiResponse[OrderListResponse]
// @Failure 401 {object} ApiResponse[OrderListResponse]
// @Failure 500 {object} ApiResponse[OrderListResponse]
// @Router /order/{order_id}/refund [put]
func (sv *Server) refundOrder(c *gin.Context) {
	var params OrderIDParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[OrderListResponse](InvalidBodyCode, "", err))
		return
	}
	var req RefundOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[OrderListResponse](InvalidBodyCode, "", err))
		return
	}
	order, err := sv.repo.GetOrder(c, uuid.MustParse(params.ID))
	if err != nil {
		if err == repository.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, createErrorResponse[OrderListResponse](NotFoundCode, "", err))
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[OrderListResponse](InternalServerErrorCode, "", err))
		return
	}

	if order.Status != repository.OrderStatusDelivered {
		c.JSON(http.StatusBadRequest, createErrorResponse[OrderListResponse](InvalidPaymentCode, "", errors.New("order cannot be refunded")))
		return
	}
	var reason payment.RefundReason
	var amountRefund int64
	switch req.Reason {
	case "defective":
		reason = payment.RefundReasonByDefectiveOrDamaged
		amountRefund = order.TotalPrice.Int.Int64()
	case "damaged":
		reason = payment.RefundReasonByDefectiveOrDamaged
		amountRefund = order.TotalPrice.Int.Int64()
	case "fraudulent":
		reason = payment.RefundReasonByFraudulent
		amountRefund = order.TotalPrice.Int.Int64()
	case "requested_by_customer":
		amountRefund = order.TotalPrice.Int.Int64() * 90 / 100
		reason = payment.RefundReasonRequestedByCustomer
	}
	err = sv.repo.RefundOrderTx(c, repository.RefundOrderTxArgs{
		OrderID: uuid.MustParse(params.ID),
		RefundPaymentFromGateway: func(paymentID string, gateway repository.PaymentGateway) (string, error) {
			switch gateway {
			case repository.PaymentGatewayStripe:
				stripeInstance, err := payment.NewStripePayment(sv.config.StripeSecretKey)
				if err != nil {
					c.JSON(http.StatusInternalServerError, createErrorResponse[OrderListResponse](InternalServerErrorCode, "", err))
					return "", err
				}
				sv.paymentCtx.SetStrategy(stripeInstance)

				return sv.paymentCtx.RefundPayment(paymentID, amountRefund, reason)
			default:
				return "", nil
			}
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[OrderListResponse](InternalServerErrorCode, "", err))
		return
	}
	c.JSON(http.StatusOK, createSuccessResponse(c, order, "success", nil, nil))
}

// --- Admin API ---

// @Summary Get all orders (Admin endpoint)
// @Description Get all orders with pagination and filtering
// @Tags admin
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param status query string false "Filter by status"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[[]OrderListResponse]
// @Failure 401 {object} ApiResponse[[]OrderListResponse]
// @Failure 403 {object} ApiResponse[[]OrderListResponse]
// @Failure 500 {object} ApiResponse[[]OrderListResponse]
// @Router /admin/orders [get]
func (sv *Server) getAdminOrdersHandler(c *gin.Context) {
	var orderListQuery OrderListQuery
	if err := c.ShouldBindQuery(&orderListQuery); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[[]OrderListResponse](InvalidBodyCode, "", err))
		return
	}

	dbParams := repository.GetOrdersParams{
		Limit:  orderListQuery.PageSize,
		Offset: (orderListQuery.Page - 1) * orderListQuery.PageSize,
	}

	if orderListQuery.Status != nil {
		dbParams.Status = repository.NullOrderStatus{
			OrderStatus: repository.OrderStatus(*orderListQuery.Status),
			Valid:       true,
		}
	}

	if orderListQuery.PaymentStatus != nil {
		dbParams.PaymentStatus = repository.NullPaymentStatus{
			PaymentStatus: repository.PaymentStatus(*orderListQuery.PaymentStatus),
			Valid:         true,
		}
	}

	fetchedOrderRows, err := sv.repo.GetOrders(c, dbParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[[]OrderListResponse](InternalServerErrorCode, "", err))
		return
	}

	countParams := repository.CountOrdersParams{}
	if orderListQuery.Status != nil {
		countParams.Status = repository.NullOrderStatus{
			OrderStatus: repository.OrderStatus(*orderListQuery.Status),
			Valid:       true,
		}
	}

	if orderListQuery.PaymentStatus != nil {
		countParams.PaymentStatus = repository.NullPaymentStatus{
			PaymentStatus: repository.PaymentStatus(*orderListQuery.PaymentStatus),
			Valid:         true,
		}
	}

	count, err := sv.repo.CountOrders(c, countParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[[]OrderListResponse](InternalServerErrorCode, "", err))
		return
	}

	var orderResponses []OrderListResponse
	for _, aggregated := range fetchedOrderRows {
		total, _ := aggregated.TotalPrice.Float64Value()
		orderResponses = append(orderResponses, OrderListResponse{
			ID:            aggregated.ID,
			Total:         total.Float64,
			TotalItems:    int32(aggregated.TotalItems),
			Status:        aggregated.Status,
			CustomerName:  aggregated.CustomerName,
			CustomerEmail: aggregated.CustomerEmail,
			PaymentStatus: aggregated.PaymentStatus.PaymentStatus,
			CreatedAt:     aggregated.CreatedAt.UTC(),
			UpdatedAt:     aggregated.UpdatedAt.UTC(),
		})
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, orderResponses, "success", &Pagination{
		Total:           count,
		Page:            orderListQuery.Page,
		PageSize:        orderListQuery.PageSize,
		TotalPages:      (count + orderListQuery.PageSize - 1) / orderListQuery.PageSize,
		HasNextPage:     count > int64(orderListQuery.Page*orderListQuery.PageSize),
		HasPreviousPage: orderListQuery.Page > 1,
	}, nil))
}

// @Summary Get order details by ID (Admin endpoint)
// @Description Get detailed information about an order by its ID
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[OrderDetailResponse]
// @Failure 401 {object} ApiResponse[OrderDetailResponse]
// @Failure 403 {object} ApiResponse[OrderDetailResponse]
// @Failure 404 {object} ApiResponse[OrderDetailResponse]
// @Failure 500 {object} ApiResponse[OrderDetailResponse]
// @Router /admin/orders/{id} [get]
func (sv *Server) getAdminOrderDetailHandler(c *gin.Context) {
	// Reuse the existing order detail handler since admin has access to all orders
	sv.getOrderDetailHandler(c)
}
