package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhphuocnguyen/go-eshop/internal/auth"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	paymentService "github.com/thanhphuocnguyen/go-eshop/pkg/payment"
)

type PaymentRequest struct {
	CartID  int64  `json:"cart_id" binding:"required,min=1"`
	Gateway string `json:"gateway" binding:"required"`
}
type GetPaymentByOrderIDParam struct {
	OrderID int64 `uri:"order_id" binding:"required"`
}

type GetPaymentParam struct {
	PaymentID string `uri:"payment_id" binding:"required"`
}

type PaymentResponse struct {
	TransactionID string                    `json:"transaction_id"`
	Gateway       repository.PaymentGateway `json:"gateway,omitempty"`
	Status        repository.PaymentStatus  `json:"status,omitempty"`
	Details       interface{}               `json:"details"`
}

func (sv *Server) initiatePayment(c *gin.Context) {

}

func (sv *Server) getPayment(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusUnauthorized, mapErrResp(fmt.Errorf("authorization payload is not provided")))
		return
	}
	var param GetPaymentByOrderIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	order, err := sv.repo.GetOrder(c, param.OrderID)
	if err != nil {
		if err == repository.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, mapErrResp(err))
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	if order.UserID != authPayload.UserID {
		c.JSON(http.StatusForbidden, mapErrResp(fmt.Errorf("forbidden")))
		return
	}
	payment, err := sv.repo.GetPaymentTransactionByOrderID(c, param.OrderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	var details interface{}
	if payment.PaymentGateway.Valid {
		stripeInstance, err := paymentService.NewStripePayment(sv.config.StripeSecretKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, mapErrResp(err))
			return
		}
		sv.paymentCtx.SetStrategy(stripeInstance)

		details, err = sv.paymentCtx.GetPaymentObject(payment.PaymentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, mapErrResp(err))
			return
		}
	}

	c.JSON(http.StatusOK, GenericResponse[PaymentResponse]{&PaymentResponse{
		TransactionID: payment.PaymentID,
		Gateway:       payment.PaymentGateway.PaymentGateway,
		Status:        payment.Status,
		Details:       details,
	}, nil, nil})
}
