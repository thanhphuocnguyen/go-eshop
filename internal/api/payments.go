package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/constants"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
	"github.com/thanhphuocnguyen/go-eshop/pkg/payment"
)

// Setup payment-related routes
func (sv *Server) addPaymentRoutes(rg *gin.RouterGroup) {
	payments := rg.Group("/payments").Use(authenticateMiddleware(sv.tokenGenerator))
	{
		payments.GET(":id", sv.getPayment)
		payments.GET("stripe-config", sv.getStripeConfig)
		payments.POST("", sv.CreatePaymentIntent)
		payments.PUT(":orderId", sv.changePaymentStatus)
	}
}

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
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 403 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /payment [post]
func (sv *Server) CreatePaymentIntent(c *gin.Context) {
	authPayload, ok := c.MustGet(constants.AuthPayLoad).(*auth.TokenPayload)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.CreateErr(UnauthorizedCode, errors.New("authorization payload is not provided")))
		return
	}
	user, err := sv.repo.GetUserByID(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, err))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	var req models.PaymentModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	ord, err := sv.repo.GetOrder(c, uuid.MustParse(req.OrderID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, err))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, errors.New("order not found")))
		return
	}

	if ord.UserID != user.ID {
		c.JSON(http.StatusForbidden, dto.CreateErr(PermissionDeniedCode, errors.New("permission denied")))
		return
	}

	pmRow, err := sv.repo.GetPaymentByOrderID(c, ord.ID)
	if err != nil && !errors.Is(err, repository.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, errors.New("order not found")))
		return
	}

	if pmRow.ID != uuid.Nil && pmRow.Status != repository.PaymentStatusCancelled {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidPaymentCode, errors.New("payment already exists")))
		return
	}

	total, _ := ord.TotalPrice.Float64Value()
	paymentMethodId := uuid.MustParse(req.PaymentMethodID)
	// create new payment
	createPaymentParams := repository.CreatePaymentParams{
		OrderID:         ord.ID,
		Amount:          utils.GetPgNumericFromFloat(total.Float64),
		PaymentMethodID: paymentMethodId,
	}
	var resp dto.PaymentIntentSecret
	paymentMethod, err := sv.repo.GetPaymentMethodByID(c, paymentMethodId)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, err))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	intent, err := sv.paymentSrv.CreatePaymentIntent(c, paymentMethod.Code, payment.PaymentRequest{
		Amount:      int64(total.Float64 * 100), // convert to smallest currency unit
		Currency:    payment.USD,
		Email:       user.Email,
		Description: "Payment for order " + ord.ID.String(),
		Metadata:    map[string]string{"order_id": ord.ID.String()},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	createPaymentParams.PaymentIntentID = &intent.ID
	resp.ClientSecret = &intent.ClientSecret

	pmRow, err = sv.repo.CreatePayment(c, createPaymentParams)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	resp.PaymentID = pmRow.ID.String()
	c.JSON(http.StatusOK, dto.CreateDataResp(c, resp, nil, nil))
}

// @Summary Get payment  by order ID
// @Description Get payment  by order ID
// @Tags payment
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[PaymentResponse]
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 403 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /payment/{id} [get]
func (sv *Server) getPayment(c *gin.Context) {
	var param models.UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	payment, err := sv.repo.GetPaymentByID(c, uuid.MustParse(param.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	var details interface{}
	paymentMethod, err := sv.repo.GetPaymentMethodByID(c, payment.PaymentMethodID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	details, err = sv.paymentSrv.GetPayment(c, *payment.PaymentIntentID, paymentMethod.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	resp := dto.PaymentDetail{
		ID:      payment.ID.String(),
		Gateway: payment.Gateway,
		Status:  payment.Status,
		Details: details,
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, resp, nil, nil))
}

// @Summary Change payment status
// @Description Change payment status
// @Tags payment
// @Accept json
// @Produce json
// @Param paymentId path string true "Payment ID"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[PaymentResponse]
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 403 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /payment/{paymentId} [get]
func (sv *Server) changePaymentStatus(c *gin.Context) {
	var param models.UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, errors.New("order not found")))
		return
	}
	var req models.UpdatePaymentStatusModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, errors.New("order not found")))
		return
	}
	payment, err := sv.repo.GetPaymentByID(c, uuid.MustParse(param.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, errors.New("order not found")))
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, errors.New("order not found")))
		return
	}
	order, err := sv.repo.GetOrder(c, payment.OrderID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, err))
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	if order.Status != repository.OrderStatusPending {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidPaymentCode, errors.New("order is not pending")))
		return
	}

	if payment.Gateway != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidPaymentCode, errors.New("cannot change payment status for stripe payment")))
		return
	}

	if req.Status == repository.PaymentStatusCancelled {
		if payment.Status == repository.PaymentStatusSuccess {
			c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidPaymentCode, errors.New("cannot cancel payment that is already success")))
			return
		}
		if payment.Status == repository.PaymentStatusCancelled {
			c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidPaymentCode, errors.New("payment is already cancelled")))
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
			c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
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
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	resp := dto.PaymentDetail{
		ID:     payment.ID.String(),
		Status: req.Status,
	}
	c.JSON(http.StatusOK, dto.CreateDataResp(c, resp, nil, nil))
}
