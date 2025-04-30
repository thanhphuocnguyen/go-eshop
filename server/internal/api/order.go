package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
	"github.com/thanhphuocnguyen/go-eshop/pkg/payment"
)

// ---------------------------------------------- API Models ----------------------------------------------
type OrderListQuery struct {
	PaginationQueryParams
	Status        string `form:"status"`
	PaymentStatus string `form:"payment_status"`
}

type OrderIDParams struct {
	ID string `uri:"id" binding:"required,min=1"`
}

type OrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending confirmed delivering delivered completed"`
}

type OrderItemResponse struct {
	ID                 string                             `json:"id"`
	Name               string                             `json:"name"`
	ImageUrl           *string                            `json:"image_url"`
	AttributesSnapshot []repository.AttributeDataSnapshot `json:"attribute_snapshot"`
	LineTotal          float64                            `json:"line_total"`
	Quantity           int16                              `json:"quantity"`
}
type PaymentInfo struct {
	ID             string  `json:"id"`
	RefundID       *string `json:"refund_id"`
	PaymentAmount  float64 `json:"payment_amount"`
	PaymentGateWay *string `json:"payment_gateway"`
	PaymentMethod  string  `json:"payment_method"`
	PaymentStatus  string  `json:"payment_status"`
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
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
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
	user, err := sv.repo.GetUserByID(c, tokenPayload.UserID)
	if err != nil {
		if err == repository.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, createErrorResponse[[]OrderListResponse](NotFoundCode, "", err))
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[[]OrderListResponse](InternalServerErrorCode, "", err))
		return
	}

	dbParams := repository.ListOrdersParams{
		Limit:  20,
		Offset: 1,
	}

	dbParams.Limit = int64(orderListQuery.PageSize)
	dbParams.Offset = int64(orderListQuery.Page-1) * int64(orderListQuery.PageSize)

	if user.Role != repository.UserRoleAdmin {
		dbParams.CustomerID = utils.GetPgTypeUUID(tokenPayload.UserID)
	}

	listOrderRows, err := sv.repo.ListOrders(c, dbParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[[]OrderListResponse](InternalServerErrorCode, "", err))
		return
	}

	count, err := sv.repo.CountOrders(c, repository.CountOrdersParams{CustomerID: tokenPayload.UserID})

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[[]OrderListResponse](InternalServerErrorCode, "", err))
		return
	}

	var orderResponses []OrderListResponse
	for _, aggregated := range listOrderRows {
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
		TotalPages:      int(count / int64(orderListQuery.PageSize)),
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
	var shippingInfo repository.ShippingAddressSnapshot

	resp := OrderDetailResponse{
		ID:            order.ID,
		Total:         total.Float64,
		CustomerName:  order.CustomerName,
		CustomerEmail: order.CustomerEmail,
		PaymentInfo:   nil,
		Status:        order.Status,
		ShippingInfo:  shippingInfo,
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
		resp.PaymentInfo = &PaymentInfo{
			ID:             paymentInfo.ID.String(),
			RefundID:       paymentInfo.RefundID,
			PaymentGateWay: (*string)(&paymentInfo.PaymentGateway.PaymentGateway),
			PaymentMethod:  string(paymentInfo.PaymentMethod),
			PaymentStatus:  string(paymentInfo.Status),
		}
	}

	orderItems := make([]OrderItemResponse, 0, len(orderItemRows))

	for _, item := range orderItemRows {
		var attrSnapshot []repository.AttributeDataSnapshot
		lineTotal, _ := item.LineTotalSnapshot.Float64Value()
		orderItems = append(orderItems, OrderItemResponse{
			ID:                 item.VariantID.String(),
			Name:               item.ProductName,
			ImageUrl:           item.ImageUrl,
			LineTotal:          lineTotal.Float64,
			AttributesSnapshot: attrSnapshot,
			Quantity:           item.Quantity,
		})
	}

	resp.Products = orderItems

	c.JSON(http.StatusOK, createSuccessResponse(c, resp, "success", nil, apiErr))
}

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
		reason = payment.RefundReasonRequestedByCustomer
		amountRefund = order.TotalPrice.Int.Int64() * 90 / 100
	default:
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
