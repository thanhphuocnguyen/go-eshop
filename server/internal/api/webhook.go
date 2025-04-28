package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v81"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
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

		payment, err := server.repo.GetPaymentByPaymentIntentID(c, utils.GetPgTypeText(paymentIntent.ID))
		if err != nil {
			if errors.Is(err, repository.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, createErrorResponse[bool]("payment_not_found", "", err))
				return
			}
			c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
			return
		}
		updateTransactionStatus := repository.UpdatePaymentParams{
			ID: payment.ID,
			Status: repository.NullPaymentStatus{
				Valid: true,
			},
		}
		switch evt.Type {
		case stripe.EventTypePaymentIntentSucceeded:
			server.taskDistributor.SendOrderCreatedEmailTask(c,
				&worker.PayloadSendOrderCreatedEmailTask{
					PaymentID: payment.ID.String(),
					OrderID:   payment.OrderID,
				},
				asynq.MaxRetry(3),
				asynq.Queue("email"),
				asynq.ProcessIn(time.Second*5))
			server.repo.CreatePaymentTransaction(c, repository.CreatePaymentTransactionParams{
				ID:                     uuid.New(),
				PaymentID:              payment.ID,
				Amount:                 payment.Amount,
				Status:                 repository.PaymentStatusSuccess,
				GatewayTransactionID:   payment.GatewayPaymentIntentID,
				GatewayResponseCode:    utils.GetPgTypeText(evt.LastResponse.Status),
				GatewayResponseMessage: utils.GetPgTypeText(string(evt.LastResponse.RawJSON)),
			})
			updateTransactionStatus.Status.PaymentStatus = repository.PaymentStatusSuccess
		case stripe.EventTypePaymentIntentCanceled:
			server.repo.CreatePaymentTransaction(c, repository.CreatePaymentTransactionParams{
				ID:                     uuid.New(),
				PaymentID:              payment.ID,
				Amount:                 payment.Amount,
				Status:                 repository.PaymentStatusCancelled,
				GatewayTransactionID:   payment.GatewayPaymentIntentID,
				GatewayResponseCode:    utils.GetPgTypeText(evt.LastResponse.Status),
				GatewayResponseMessage: utils.GetPgTypeText(string(evt.LastResponse.RawJSON)),
			})
			updateTransactionStatus.Status.PaymentStatus = repository.PaymentStatusCancelled
		case stripe.EventTypePaymentIntentPaymentFailed:
			updateTransactionStatus.Status.PaymentStatus = repository.PaymentStatusFailed
			server.repo.CreatePaymentTransaction(c, repository.CreatePaymentTransactionParams{
				ID:                     uuid.New(),
				PaymentID:              payment.ID,
				Amount:                 payment.Amount,
				Status:                 repository.PaymentStatusFailed,
				GatewayTransactionID:   payment.GatewayPaymentIntentID,
				GatewayResponseCode:    utils.GetPgTypeText(evt.LastResponse.Status),
				GatewayResponseMessage: utils.GetPgTypeText(string(evt.LastResponse.RawJSON)),
			})
			updateTransactionStatus.Status.PaymentStatus = repository.PaymentStatusFailed
			updateTransactionStatus.ErrorMessage = utils.GetPgTypeText(string(evt.LastResponse.RawJSON))
			updateTransactionStatus.ErrorCode = utils.GetPgTypeText(evt.LastResponse.Status)
		}
		updateTransactionStatus.UpdatedAt = utils.GetPgTypeTimestamp(time.Now())
		err = server.repo.UpdatePayment(c, updateTransactionStatus)
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
