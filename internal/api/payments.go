package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
	"github.com/thanhphuocnguyen/go-eshop/pkg/paymentsrv"
)

func (sv *Server) getStripeConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"public_key": sv.config.StripePublishableKey})
}

// @Summary Initiate payment
// @Description Initiate payment
// @Tags payment
// @Accept json
// @Produce json
// @Param request body PaymentRequest true "Payment request"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[PaymentResponse]
// @Failure 400 {object} ApiResponse[gin.H]
// @Failure 401 {object} ApiResponse[gin.H]
// @Failure 403 {object} ApiResponse[gin.H]
// @Failure 404 {object} ApiResponse[gin.H]
// @Failure 500 {object} ApiResponse[gin.H]
// @Router /payment [post]
func (sv *Server) createPaymentIntentHandler(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusUnauthorized, createErrorResponse[PaymentResponse](UnauthorizedCode, "", errors.New("authorization payload is not provided")))
		return
	}
	user, err := sv.repo.GetUserByID(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[PaymentResponse](NotFoundCode, "", err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[PaymentResponse](InternalServerErrorCode, "", err))
		return
	}
	var req PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[PaymentResponse](InvalidBodyCode, "", err))
		return
	}

	ord, err := sv.repo.GetOrder(c, uuid.MustParse(req.OrderID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[PaymentResponse](NotFoundCode, "", err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[PaymentResponse](InternalServerErrorCode, "", errors.New("order not found")))
		return
	}

	if ord.CustomerID != user.ID {
		c.JSON(http.StatusForbidden, createErrorResponse[PaymentResponse](PermissionDeniedCode, "", errors.New("permission denied")))
		return
	}

	payment, err := sv.repo.GetPaymentByOrderID(c, ord.ID)
	if err != nil && !errors.Is(err, repository.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, createErrorResponse[PaymentResponse](InternalServerErrorCode, "", errors.New("order not found")))
		return
	}

	if payment.ID != uuid.Nil && payment.Status != repository.PaymentStatusCancelled {
		c.JSON(http.StatusBadRequest, createErrorResponse[PaymentResponse](InvalidPaymentCode, "", errors.New("payment already exists")))
		return
	}

	total, _ := ord.TotalPrice.Float64Value()
	// create new payment
	createPaymentParams := repository.CreatePaymentParams{
		OrderID: ord.ID,
		Amount:  utils.GetPgNumericFromFloat(total.Float64),
		Method:  repository.PaymentMethod(req.PaymentMethod),
	}
	var resp CreatePaymentIntentResponse
	switch req.PaymentMethod {
	case string(repository.PaymentMethodStripe):
		stripeInstance, err := paymentsrv.NewStripePayment(sv.config.StripeSecretKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse[PaymentResponse](InternalServerErrorCode, "", err))
			return
		}
		sv.paymentCtx.SetStrategy(stripeInstance)
		createPaymentIntentResult, err := sv.paymentCtx.CreatePaymentIntent(total.Float64, user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse[PaymentResponse](InternalServerErrorCode, "", err))
			return
		}
		paymentIntent := createPaymentIntentResult.(*stripe.PaymentIntent)

		createPaymentParams.PaymentIntentID = &paymentIntent.ID
		resp.ClientSecret = &paymentIntent.ClientSecret
	}

	payment, err = sv.repo.CreatePayment(c, createPaymentParams)

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[PaymentResponse](InternalServerErrorCode, "", err))
		return
	}

	resp.PaymentID = payment.ID.String()
	c.JSON(http.StatusOK, createSuccessResponse(c, resp, "", nil, nil))
	return
}

// @Summary Get payment  by order ID
// @Description Get payment  by order ID
// @Tags payment
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[PaymentResponse]
// @Failure 400 {object} ApiResponse[gin.H]
// @Failure 401 {object} ApiResponse[gin.H]
// @Failure 403 {object} ApiResponse[gin.H]
// @Failure 404 {object} ApiResponse[gin.H]
// @Failure 500 {object} ApiResponse[gin.H]
// @Router /payment/{id} [get]
func (sv *Server) getPaymentHandler(c *gin.Context) {
	var param UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[PaymentResponse](InvalidBodyCode, "", err))
		return
	}

	payment, err := sv.repo.GetPaymentByID(c, uuid.MustParse(param.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[PaymentResponse](InternalServerErrorCode, "", err))
		return
	}

	var details interface{}
	if payment.Method == repository.PaymentMethodStripe {
		stripeInstance, err := paymentsrv.NewStripePayment(sv.config.StripeSecretKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse[PaymentResponse](InternalServerErrorCode, "", err))
			return
		}
		sv.paymentCtx.SetStrategy(stripeInstance)

		details, err = sv.paymentCtx.GetPaymentObject(*payment.PaymentIntentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse[PaymentResponse](InternalServerErrorCode, "", err))
			return
		}
	}

	resp := PaymentResponse{
		ID:      payment.ID.String(),
		Gateway: payment.Gateway,
		Status:  payment.Status,
		Details: details,
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, resp, "", nil, nil))
}

// @Summary Change payment status
// @Description Change payment status
// @Tags payment
// @Accept json
// @Produce json
// @Param paymentId path string true "Payment ID"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[PaymentResponse]
// @Failure 400 {object} ApiResponse[gin.H]
// @Failure 401 {object} ApiResponse[gin.H]
// @Failure 403 {object} ApiResponse[gin.H]
// @Failure 404 {object} ApiResponse[gin.H]
// @Failure 500 {object} ApiResponse[gin.H]
// @Router /payment/{paymentId} [get]
func (sv *Server) changePaymentStatusHandler(c *gin.Context) {
	var param UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[PaymentResponse](InvalidBodyCode, "", errors.New("order not found")))
		return
	}
	var req UpdatePaymentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[PaymentResponse](InvalidBodyCode, "", errors.New("order not found")))
		return
	}
	payment, err := sv.repo.GetPaymentByID(c, uuid.MustParse(param.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[PaymentResponse](NotFoundCode, "", errors.New("order not found")))
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[PaymentResponse](InternalServerErrorCode, "", errors.New("order not found")))
		return
	}
	order, err := sv.repo.GetOrder(c, payment.OrderID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[PaymentResponse](NotFoundCode, "", err))
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[PaymentResponse](InternalServerErrorCode, "", err))
		return
	}
	if order.Status != repository.OrderStatusPending {
		c.JSON(http.StatusBadRequest, createErrorResponse[PaymentResponse](InvalidPaymentCode, "", errors.New("order is not pending")))
		return
	}

	if payment.Gateway != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[PaymentResponse](InvalidPaymentCode, "", errors.New("cannot change payment status for stripe payment")))
		return
	}

	if req.Status == repository.PaymentStatusCancelled {
		if payment.Status == repository.PaymentStatusSuccess {
			c.JSON(http.StatusBadRequest, createErrorResponse[PaymentResponse](InvalidPaymentCode, "", errors.New("cannot cancel payment that is already success")))
			return
		}
		if payment.Status == repository.PaymentStatusCancelled {
			c.JSON(http.StatusBadRequest, createErrorResponse[PaymentResponse](InvalidPaymentCode, "", errors.New("payment is already cancelled")))
			return
		}
	}

	if req.Status == repository.PaymentStatusSuccess {
		_, err := sv.repo.UpdateOrder(c, repository.UpdateOrderParams{
			ID: order.ID,
			Status: repository.NullOrderStatus{
				OrderStatus: repository.OrderStatusDelivered,
				Valid:       true,
			},
			DeliveredAt: utils.GetPgTypeTimestamp(time.Now()),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse[PaymentResponse](InternalServerErrorCode, "", err))
			return
		}

	}

	err = sv.repo.UpdatePayment(c, repository.UpdatePaymentParams{
		ID: payment.ID,
		Status: repository.NullPaymentStatus{
			PaymentStatus: req.Status,
			Valid:         true,
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[PaymentResponse](InternalServerErrorCode, "", err))
		return
	}
	resp := PaymentResponse{
		ID:     payment.ID.String(),
		Status: req.Status,
	}
	c.JSON(http.StatusOK, createSuccessResponse(c, resp, "", nil, nil))
}
