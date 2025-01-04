package api

import (
	"encoding/json"
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
	var paymentIntent stripe.PaymentIntent
	err := json.Unmarshal(evt.Data.Raw, &paymentIntent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	payment, err := server.repo.GetPaymentTransactionByID(c, paymentIntent.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	updateTransactionStatus := repository.UpdatePaymentTransactionParams{
		PaymentID: payment.PaymentID,
	}
	switch evt.Type {
	case stripe.EventTypePaymentIntentSucceeded:
		updateTransactionStatus.Status = repository.NullPaymentStatus{
			PaymentStatus: repository.PaymentStatusSuccess,
			Valid:         true,
		}
		err = server.repo.UpdatePaymentTransaction(c, updateTransactionStatus)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Handle successful payment
	case stripe.EventTypePaymentIntentCanceled:
		updateTransactionStatus.Status = repository.NullPaymentStatus{
			PaymentStatus: repository.PaymentStatusCancelled,
			Valid:         true,
		}
		err = server.repo.UpdatePaymentTransaction(c, updateTransactionStatus)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Handle canceled payment
	case stripe.EventTypePaymentIntentPaymentFailed:
		updateTransactionStatus.Status = repository.NullPaymentStatus{
			PaymentStatus: repository.PaymentStatusFailed,
			Valid:         true,
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
