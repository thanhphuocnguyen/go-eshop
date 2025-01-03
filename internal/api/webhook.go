package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
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
		err = server.repo.UpdatePaymentTransaction(c, repository.UpdatePaymentTransactionParams{
			PaymentID: payment.PaymentID,
			Status: repository.NullPaymentStatus{
				PaymentStatus: repository.PaymentStatusSuccess,
				Valid:         true,
			},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Handle successful payment
	case stripe.EventTypePaymentIntentCanceled:
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
		err = server.repo.UpdatePaymentTransaction(c, repository.UpdatePaymentTransactionParams{
			PaymentID: payment.PaymentID,
			Status: repository.NullPaymentStatus{
				PaymentStatus: repository.PaymentStatus(repository.PaymentStatusCancelled),
				Valid:         true,
			},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Handle canceled payment
	case stripe.EventTypePaymentIntentPaymentFailed:
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
		err = server.repo.UpdatePaymentTransaction(c, repository.UpdatePaymentTransactionParams{
			PaymentID: payment.PaymentID,
			Status: repository.NullPaymentStatus{
				PaymentStatus: repository.PaymentStatus(repository.PaymentStatusFailed),
				Valid:         true,
			},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	case stripe.EventTypePaymentIntentProcessing:
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
		err = server.repo.UpdatePaymentTransaction(c, repository.UpdatePaymentTransactionParams{
			PaymentID: payment.PaymentID,
			Status: repository.NullPaymentStatus{
				PaymentStatus: repository.PaymentStatus(repository.PaymentStatusProcessing),
				Valid:         true,
			},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unhandled event type"})
		return
	}
}
