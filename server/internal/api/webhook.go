package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v81"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
)

func (server *Server) stripeWebhook(c *gin.Context) {
	var evt stripe.Event
	if err := c.ShouldBindJSON(&evt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch evt.Type {
	case stripe.EventTypePaymentIntentSucceeded:
	case stripe.EventTypePaymentIntentCanceled:
	case stripe.EventTypePaymentIntentPaymentFailed:
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(evt.Data.Raw, &paymentIntent)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Info().Msgf("Received event of type: %v", evt)

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
