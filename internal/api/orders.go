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
	"github.com/thanhphuocnguyen/go-eshop/pkg/pmgateway"
)

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
// @Router /order/list [get]
func (sv *Server) getOrdersHandler(c *gin.Context) {
	tokenPayload, ok := c.MustGet(AuthPayLoad).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusUnauthorized, createErr(UnauthorizedCode, errors.New("authorization payload is not provided")))
		return
	}

	var orderListQuery OrderListQuery
	if err := c.ShouldBindQuery(&orderListQuery); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	dbParams := repository.GetOrdersParams{
		Limit:  orderListQuery.PageSize,
		Offset: (orderListQuery.Page - 1) * orderListQuery.PageSize,
	}

	if tokenPayload.RoleCode != "admin" {
		dbParams.CustomerID = utils.GetPgTypeUUID(tokenPayload.UserID)
	}

	fetchedOrderRows, err := sv.repo.GetOrders(c, dbParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	count, err := sv.repo.CountOrders(c, repository.CountOrdersParams{CustomerID: utils.GetPgTypeUUID(tokenPayload.UserID)})

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	var orderResponses []OrderListResponse
	for _, aggregated := range fetchedOrderRows {
		// Convert PaymentStatus interface{} to PaymentStatus type
		paymentStatus := repository.PaymentStatusPending
		if aggregated.PaymentStatus.Valid {
			paymentStatus = aggregated.PaymentStatus.PaymentStatus
		}

		total, _ := aggregated.TotalPrice.Float64Value()
		orderResponses = append(orderResponses, OrderListResponse{
			ID:            aggregated.ID,
			Total:         total.Float64,
			TotalItems:    int32(aggregated.TotalItems),
			Status:        aggregated.Status,
			PaymentStatus: paymentStatus,
			CreatedAt:     aggregated.CreatedAt.UTC(),
			UpdatedAt:     aggregated.UpdatedAt.UTC(),
		})
	}

	c.JSON(http.StatusOK, createDataResp(c, orderResponses, createPagination(orderListQuery.Page, orderListQuery.PageSize, count), nil))
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
func (sv *Server) getOrderDetailHandler(c *gin.Context) {
	var params UriIDParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	var resp *OrderDetailResponse

	if err := sv.cacheSrv.Get(c, "order_detail:"+params.ID, &resp); err == nil {
		if resp != nil {
			c.JSON(http.StatusOK, createDataResp(c, resp, nil, nil))
			return
		}
	}

	order, err := sv.repo.GetOrder(c, uuid.MustParse(params.ID))
	if err != nil {
		if err == repository.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, createErr(NotFoundCode, err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	total, _ := order.TotalPrice.Float64Value()
	paymentAmount, _ := order.PaymentAmount.Float64Value()
	resp = &OrderDetailResponse{
		ID:            order.ID,
		Total:         total.Float64,
		CustomerName:  order.CustomerName,
		CustomerEmail: order.CustomerEmail,
		Status:        order.Status,
		ShippingInfo:  order.ShippingAddress,
		PaymentInfo: PaymentInfoModel{
			ID:       order.PaymentID.String(),
			Amount:   paymentAmount.Float64,
			IntendID: order.PaymentIntentID,
			GateWay:  order.Gateway,
			Method:   string(order.PaymentMethod),
			Status:   string(order.PaymentStatus),
		},
		CreatedAt: order.CreatedAt.UTC(),
		Products:  []OrderItemResponse{},
	}

	orderItemRows, err := sv.repo.GetOrderItems(c, order.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	paymentInfo, err := sv.repo.GetPaymentByOrderID(c, order.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}
	var apiErr *ApiError = nil

	if paymentInfo.Status == repository.PaymentStatusPending &&
		paymentInfo.PaymentIntentID != nil &&
		order.PaymentMethod == "Stripe" {
		resp.PaymentInfo.IntendID = paymentInfo.PaymentIntentID
		stripeInstance, err := pmgateway.NewStripePayment(sv.config.StripeSecretKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
			return
		}
		sv.paymentSrv.SetStrategy(stripeInstance)
		paymentDetail, err := sv.paymentSrv.GetPaymentObject(*paymentInfo.PaymentIntentID)
		if err != nil {
			log.Err(err).Msg("failed to get payment intent")
			apiErr = &ApiError{
				Code:    InternalServerErrorCode,
				Details: "failed to get payment intent",
				Stack:   err,
			}
		} else {
			stripeObject, _ := paymentDetail.(*stripe.PaymentIntent)
			resp.PaymentInfo.ClientSecret = &stripeObject.ClientSecret
		}
	}

	orderItems := make([]OrderItemResponse, 0, len(orderItemRows))

	for _, item := range orderItemRows {
		lineTotal, _ := item.LineTotalSnapshot.Float64Value()
		itemResp := OrderItemResponse{
			ID:                 item.ID.String(),
			VariantID:          item.VariantID.String(),
			Name:               item.ProductName,
			ImageUrl:           item.ImageUrl,
			LineTotal:          lineTotal.Float64,
			AttributesSnapshot: item.AttributesSnapshot,
			Quantity:           item.Quantity,
		}
		if item.RatingID.Valid {
			id, _ := uuid.FromBytes(item.RatingID.Bytes[:])
			rating, _ := item.Rating.Float64Value()
			itemResp.Rating = &RatingModel{
				ID:        id.String(),
				Title:     *item.ReviewTitle,
				Content:   *item.ReviewContent,
				Rating:    rating.Float64,
				CreatedAt: item.RatingCreatedAt.Time.UTC(),
			}
		}
		orderItems = append(orderItems, itemResp)
	}

	resp.Products = orderItems

	if err := sv.cacheSrv.Set(c, "order_detail:"+params.ID, resp, utils.TimeDurationPtr(5*time.Minute)); err != nil {
		log.Err(err).Msg("failed to cache order detail")
		apiErr = &ApiError{
			Code:    InternalServerErrorCode,
			Details: "failed to cache order detail",
			Stack:   err,
		}
	}

	c.JSON(http.StatusOK, createDataResp(c, resp, nil, apiErr))
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
	tokenPayload, _ := c.MustGet(AuthPayLoad).(*auth.Payload)
	var params UriIDParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}
	order, err := sv.repo.GetOrder(c, uuid.MustParse(params.ID))
	if err != nil {
		if err == repository.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, createErr(NotFoundCode, err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}
	if order.CustomerID != tokenPayload.UserID {
		c.JSON(http.StatusForbidden, createErr(PermissionDeniedCode, errors.New("you do not have permission to access this order")))
		return
	}
	if order.Status != repository.OrderStatusDelivered {
		c.JSON(http.StatusBadRequest, createErr(InvalidPaymentCode, errors.New("order cannot be confirmed")))
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
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}
	var apiErr *ApiError
	if err := sv.cacheSrv.Delete(c, "order_detail:"+params.ID); err != nil {
		log.Err(err).Msg("failed to delete order detail cache")
		apiErr = &ApiError{
			Code:    InternalServerErrorCode,
			Details: "failed to delete order detail cache",
			Stack:   err,
		}
	}

	c.JSON(http.StatusOK, createDataResp(c, true, nil, apiErr))
}

// @Summary Cancel order
// @Description Cancel order by order ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[OrderListResponse]
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /order/{orderId}/cancel [put]
func (sv *Server) cancelOrder(c *gin.Context) {
	tokenPayload, _ := c.MustGet(AuthPayLoad).(*auth.Payload)

	var params UriIDParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}
	var req CancelOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	order, err := sv.repo.GetOrder(c, uuid.MustParse(params.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}
	userRole := c.GetString(UserRole)
	if order.CustomerID != tokenPayload.UserID && userRole != "admin" {
		c.JSON(http.StatusForbidden, createErr(PermissionDeniedCode, errors.New("you do not have permission to access this order")))
		return
	}

	paymentRow, err := sv.repo.GetPaymentByOrderID(c, order.ID)
	if err != nil && !errors.Is(err, repository.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	// if order status is not pending or user is not admin
	if order.Status != repository.OrderStatusPending || (paymentRow.Status != repository.PaymentStatusPending) {
		c.JSON(http.StatusBadRequest, createErr(PermissionDeniedCode, errors.New("order cannot be cancelled")))
		return
	}

	// if order
	cancelOrderTxParams := repository.CancelOrderTxArgs{
		OrderID: uuid.MustParse(params.ID),
		CancelPaymentFromMethod: func(paymentID string, method string) error {
			switch method {
			case "Stripe":
				stripeInstance, err := pmgateway.NewStripePayment(sv.config.StripeSecretKey)
				if err != nil {
					c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
					return err
				}
				sv.paymentSrv.SetStrategy(stripeInstance)
				_, err = sv.paymentSrv.CancelPayment(paymentID, req.Reason)
				return err
			default:
				return nil
			}
		},
	}
	ordId, err := sv.repo.CancelOrderTx(c, cancelOrderTxParams)

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(repository.ErrDeadlockDetected.InternalQuery, err))
		return
	}
	sv.cacheSrv.Delete(c, "order_detail:"+params.ID)
	c.JSON(http.StatusOK, createDataResp(c, ordId, nil, nil))
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
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /order/{orderId}/status [put]
func (sv *Server) changeOrderStatus(c *gin.Context) {
	var params UriIDParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}
	var req OrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}
	order, err := sv.repo.GetOrder(c, uuid.MustParse(params.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErr(NotFoundCode, err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}
	if order.Status == repository.OrderStatusDelivered || order.Status == repository.OrderStatusCancelled || order.Status == repository.OrderStatusRefunded {
		c.JSON(http.StatusBadRequest, createErr(InvalidPaymentCode, errors.New("order cannot be changed")))
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
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	if err := sv.cacheSrv.Delete(c, "order_detail:"+params.ID); err != nil {
		log.Err(err).Msg("failed to delete order detail cache")
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}
	sv.cacheSrv.Delete(c, "order_detail:"+params.ID)

	c.JSON(http.StatusOK, createDataResp(c, rs, nil, nil))
}

// @Summary Refund order
// @Description Refund order by order ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[OrderListResponse]
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /order/{orderId}/refund [put]
func (sv *Server) refundOrder(c *gin.Context) {
	var params UriIDParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}
	var req RefundOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}
	order, err := sv.repo.GetOrder(c, uuid.MustParse(params.ID))
	if err != nil {
		if err == repository.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, createErr(NotFoundCode, err))
		}
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	if order.Status != repository.OrderStatusDelivered {
		c.JSON(http.StatusBadRequest, createErr(InvalidPaymentCode, errors.New("order cannot be refunded")))
		return
	}
	var reason pmgateway.RefundReason
	var amountRefund int64
	switch req.Reason {
	case "defective":
		reason = pmgateway.RefundReasonByDefectiveOrDamaged
		amountRefund = order.TotalPrice.Int.Int64()
	case "damaged":
		reason = pmgateway.RefundReasonByDefectiveOrDamaged
		amountRefund = order.TotalPrice.Int.Int64()
	case "fraudulent":
		reason = pmgateway.RefundReasonByFraudulent
		amountRefund = order.TotalPrice.Int.Int64()
	case "requested_by_customer":
		amountRefund = order.TotalPrice.Int.Int64() * 90 / 100
		reason = pmgateway.RefundReasonRequestedByCustomer
	}
	err = sv.repo.RefundOrderTx(c, repository.RefundOrderTxArgs{
		OrderID: uuid.MustParse(params.ID),
		RefundPaymentFromMethod: func(paymentID string, method string) (string, error) {
			switch method {
			case "Stripe":
				stripeInstance, err := pmgateway.NewStripePayment(sv.config.StripeSecretKey)
				if err != nil {
					return "", err
				}
				sv.paymentSrv.SetStrategy(stripeInstance)

				return sv.paymentSrv.RefundPayment(paymentID, amountRefund, reason)
			default:
				return "", nil
			}
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}
	sv.cacheSrv.Delete(c, "order_detail:"+params.ID)

	c.JSON(http.StatusOK, createDataResp(c, order, nil, nil))
}

// --- Admin API ---

// @Summary Get all orders (Admin endpoint)
// @Description Get all orders with pagination and filtering
// @Tags Admin
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Param status query string false "Filter by status"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[[]OrderListResponse]
// @Failure 401 {object} ErrorResp
// @Failure 403 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/orders [get]
func (sv *Server) getAdminOrdersHandler(c *gin.Context) {
	var orderListQuery OrderListQuery
	if err := c.ShouldBindQuery(&orderListQuery); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
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
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
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
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	var orderResponses []OrderListResponse
	for _, aggregated := range fetchedOrderRows {
		// Convert PaymentStatus interface{} to PaymentStatus type
		paymentStatus := repository.PaymentStatusPending
		if aggregated.PaymentStatus.Valid {
			paymentStatus = aggregated.PaymentStatus.PaymentStatus
		}

		total, _ := aggregated.TotalPrice.Float64Value()
		orderResponses = append(orderResponses, OrderListResponse{
			ID:            aggregated.ID,
			Total:         total.Float64,
			TotalItems:    int32(aggregated.TotalItems),
			Status:        aggregated.Status,
			CustomerName:  aggregated.CustomerName,
			CustomerEmail: aggregated.CustomerEmail,
			PaymentStatus: paymentStatus,
			CreatedAt:     aggregated.CreatedAt.UTC(),
			UpdatedAt:     aggregated.UpdatedAt.UTC(),
		})
	}

	c.JSON(http.StatusOK, createDataResp(c, orderResponses, createPagination(orderListQuery.Page, orderListQuery.PageSize, count), nil))
}

// @Summary Get order details by ID (Admin endpoint)
// @Description Get detailed information about an order by its ID
// @Tags Admin
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[OrderDetailResponse]
// @Failure 401 {object} ErrorResp
// @Failure 403 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/orders/{id} [get]
func (sv *Server) getAdminOrderDetailHandler(c *gin.Context) {
	// Reuse the existing order detail handler since admin has access to all orders
	sv.getOrderDetailHandler(c)
}
