package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v81"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/worker"
)

// @Summary Stripe webhook
// @Description Stripe webhook
// @Tags webhook
// @Accept json
// @Produce json
// @Success 200 {object} ApiResponse[bool]
// @Failure 400 {object} ApiResponse[bool]
// @Failure 500 {object} ApiResponse[bool]
// @Router /webhook/stripe [post]
func (server *Server) stripeWebhook(c *gin.Context) {
	var evt stripe.Event
	if err := c.ShouldBindJSON(&evt); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[bool](InvalidEventCode, "", err))
		return
	}

	switch evt.Type {
	case stripe.EventTypePaymentIntentSucceeded:
	case stripe.EventTypePaymentIntentCanceled:
	case stripe.EventTypePaymentIntentPaymentFailed:
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(evt.Data.Raw, &paymentIntent)
		if err != nil {
			c.JSON(http.StatusBadRequest, createErrorResponse[bool]("", "", err))
			return
		}
		log.Info().Interface("evt type", evt.Type).Msg("Received stripe event")

		payment, err := server.repo.GetPaymentTransactionByID(c, paymentIntent.ID)
		if err != nil {
			if errors.Is(err, repository.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, createErrorResponse[bool]("payment_not_found", "", err))
				return
			}
			c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
			return
		}
		updateTransactionStatus := repository.UpdatePaymentTransactionParams{
			ID: payment.ID,
			Status: repository.NullPaymentStatus{
				Valid: true,
			},
		}
		switch evt.Type {
		case stripe.EventTypePaymentIntentSucceeded:
			server.taskDistributor.SendOrderCreatedEmailTask(c,
				&worker.PayloadSendOrderCreatedEmailTask{
					PaymentID: payment.ID,
					OrderID:   payment.OrderID,
				},
				asynq.MaxRetry(3),
				asynq.Queue("email"),
				asynq.ProcessIn(time.Second*5))
			updateTransactionStatus.Status.PaymentStatus = repository.PaymentStatusSuccess
		case stripe.EventTypePaymentIntentCanceled:
			updateTransactionStatus.Status.PaymentStatus = repository.PaymentStatusCancelled
		case stripe.EventTypePaymentIntentPaymentFailed:
			updateTransactionStatus.Status.PaymentStatus = repository.PaymentStatusFailed
		}
		err = server.repo.UpdatePaymentTransaction(c, updateTransactionStatus)
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
			return
		}
	default:
		log.Info().Msgf("Unhandled event type: %s", evt.Type)
		c.JSON(http.StatusOK, createSuccessResponse(c, true, "", nil, nil))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, true, "", nil, nil))
}
