package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/auth"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	paymentService "github.com/thanhphuocnguyen/go-eshop/pkg/payment"
)

type PaymentRequest struct {
	CartID  int64  `json:"cart_id" binding:"required,min=1"`
	Gateway string `json:"gateway" binding:"required"`
}
type GetPaymentByOrderIDParam struct {
	OrderID uuid.UUID `uri:"order_id" binding:"required,uuid"`
}

type GetPaymentParam struct {
	ID string `uri:"payment_id" binding:"required"`
}

type PaymentResponse struct {
	TransactionID string                    `json:"transaction_id"`
	Gateway       repository.PaymentGateway `json:"gateway,omitempty"`
	Status        repository.PaymentStatus  `json:"status,omitempty"`
	Details       interface{}               `json:"details"`
}

type ChangePaymentStatusReq struct {
	Status repository.PaymentStatus `json:"status" binding:"required"`
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
// @Success 200 {object} PaymentResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /payment [post]
func (sv *Server) initiatePayment(c *gin.Context) {
	var param GetPaymentByOrderIDParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	ord, err := sv.repo.GetOrder(c, param.OrderID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse(http.StatusNotFound, "", errors.New("order not found")))
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusNotFound, "", errors.New("order not found")))
		return
	}

	_, err = sv.repo.GetPaymentTransactionByOrderID(c, param.OrderID)
	if err == nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", errors.New("order already paid")))
	}

	if err != nil && err != repository.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusNotFound, "", errors.New("order not found")))
		return
	}
	// Only for cod order
	payment, err := sv.repo.CreatePaymentTransaction(c, repository.CreatePaymentTransactionParams{
		ID:            uuid.New().String(),
		OrderID:       ord.ID,
		Amount:        ord.TotalPrice,
		PaymentMethod: repository.PaymentMethodCod,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusNotFound, "", errors.New("order not found")))
		return
	}
	resp := PaymentResponse{
		TransactionID: payment.ID,
		Gateway:       payment.PaymentGateway.PaymentGateway,
		Status:        payment.Status,
		Details:       nil,
	}
	c.JSON(http.StatusOK, createSuccessResponse(c, resp, "", nil, nil))
}

// @Summary Get payment transaction by order ID
// @Description Get payment transaction by order ID
// @Tags payment
// @Accept json
// @Produce json
// @Param order_id path int true "Order ID"
// @Security BearerAuth
// @Success 200 {object} PaymentResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /payment/{order_id} [get]
func (sv *Server) getPayment(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusUnauthorized, createErrorResponse(http.StatusUnauthorized, "", errors.New("authorization payload is not provided")))
		return
	}
	var param GetPaymentByOrderIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusNotFound, "", errors.New("order not found")))
		return
	}

	order, err := sv.repo.GetOrder(c, param.OrderID)
	if err != nil {
		if err == repository.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, createErrorResponse(http.StatusNotFound, "", errors.New("order not found")))
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusNotFound, "", errors.New("order not found")))
		return
	}
	if order.CustomerID != authPayload.UserID {
		c.JSON(http.StatusForbidden, createErrorResponse(http.StatusForbidden, "", errors.New("forbidden")))
		return
	}
	payment, err := sv.repo.GetPaymentTransactionByOrderID(c, param.OrderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusNotFound, "", errors.New("order not found")))
		return
	}
	var details interface{}
	if payment.PaymentGateway.Valid {
		stripeInstance, err := paymentService.NewStripePayment(sv.config.StripeSecretKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusNotFound, "", errors.New("order not found")))
			return
		}
		sv.paymentCtx.SetStrategy(stripeInstance)

		details, err = sv.paymentCtx.GetPaymentObject(payment.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusNotFound, "", errors.New("order not found")))
			return
		}
	}
	resp := PaymentResponse{
		TransactionID: payment.ID,
		Gateway:       payment.PaymentGateway.PaymentGateway,
		Status:        payment.Status,
		Details:       details,
	}
	c.JSON(http.StatusOK, createSuccessResponse(c, resp, "", nil, nil))
}

// @Summary Change payment status
// @Description Change payment status
// @Tags payment
// @Accept json
// @Produce json
// @Param payment_id path string true "Payment ID"
// @Security BearerAuth
// @Success 200 {object} PaymentResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /payment/{payment_id} [get]
func (sv *Server) changePaymentStatus(c *gin.Context) {
	var param GetPaymentParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusNotFound, "", errors.New("order not found")))
		return
	}
	var req ChangePaymentStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusNotFound, "", errors.New("order not found")))
		return
	}
	payment, err := sv.repo.GetPaymentTransactionByID(c, param.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse(http.StatusNotFound, "", errors.New("order not found")))
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusNotFound, "", errors.New("order not found")))
		return
	}
	order, err := sv.repo.GetOrder(c, payment.OrderID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse(http.StatusNotFound, "", errors.New("order not found")))
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusNotFound, "", errors.New("order not found")))
		return
	}
	if order.Status != repository.OrderStatusPending {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", errors.New("order is not pending")))
		return
	}

	if payment.PaymentGateway.Valid {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", errors.New("cannot change payment status for stripe payment")))
		return
	}

	if req.Status == repository.PaymentStatusCancelled {
		if payment.Status == repository.PaymentStatusSuccess {
			c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", errors.New("cannot cancel payment that is already success")))
			return
		}
		if payment.Status == repository.PaymentStatusCancelled {
			c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", errors.New("payment is already cancelled")))
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
			c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusNotFound, "", errors.New("order not found")))
			return
		}

	}

	err = sv.repo.UpdatePaymentTransaction(c, repository.UpdatePaymentTransactionParams{
		ID: payment.ID,
		Status: repository.NullPaymentStatus{
			PaymentStatus: req.Status,
			Valid:         true,
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusNotFound, "", errors.New("order not found")))
		return
	}
	resp := PaymentResponse{
		TransactionID: payment.ID,
		Status:        req.Status,
	}
	c.JSON(http.StatusOK, createSuccessResponse(c, resp, "", nil, nil))
}
