package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
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
// @Success 200 {object} GenericResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /webhook/stripe [post]
func (server *Server) stripeWebhook(c *gin.Context) {
	var evt stripe.Event
	if err := c.ShouldBindJSON(&evt); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	switch evt.Type {
	case stripe.EventTypePaymentIntentSucceeded: // TODO: handle success payment send email order created to customer
	case stripe.EventTypePaymentIntentCanceled:
	case stripe.EventTypePaymentIntentPaymentFailed:
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(evt.Data.Raw, &paymentIntent)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Info().Interface("evt type", evt.Type).Msg("Received stripe event")

		payment, err := server.repo.GetPaymentTransactionByID(c, paymentIntent.ID)
		if err != nil {
			if errors.Is(err, repository.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		updateTransactionStatus := repository.UpdatePaymentTransactionParams{
			PaymentID: payment.PaymentID,
			Status: repository.NullPaymentStatus{
				Valid: true,
			},
		}
		switch evt.Type {
		case stripe.EventTypePaymentIntentSucceeded:
			server.taskDistributor.SendOrderCreatedEmailTask(c, &worker.PayloadSendOrderCreatedEmailTask{})
			updateTransactionStatus.Status.PaymentStatus = repository.PaymentStatusSuccess
		case stripe.EventTypePaymentIntentCanceled:
			updateTransactionStatus.Status.PaymentStatus = repository.PaymentStatusCancelled
		case stripe.EventTypePaymentIntentPaymentFailed:
			updateTransactionStatus.Status.PaymentStatus = repository.PaymentStatusFailed
		}
		err = server.repo.UpdatePaymentTransaction(c, updateTransactionStatus)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	default:
		log.Info().Msgf("Unhandled event type: %s", evt.Type)
		c.JSON(http.StatusOK, gin.H{"message": "Unhandled event type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
