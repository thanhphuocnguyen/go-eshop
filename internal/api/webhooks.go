package api

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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
func (server *Server) stripeEventHandler(c *gin.Context) {
	var evt stripe.Event
	if err := c.ShouldBindJSON(&evt); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[bool](InvalidEventCode, "", err))
		return
	}

	// unmarshal the event data
	evtData := evt.Data.Object
	var id *string
	if evtData["id"] != nil {
		id = utils.StringPtr(evtData["id"].(string))
	}
	log.Info().Interface("evt type", evt.Type).Msg("Received stripe event")

	var failureCode *string
	var failureMsg *string
	if evtData["failure_code"] != nil {
		failureCode = utils.StringPtr(evtData["failure_code"].(string))
	}
	if evtData["failure_message"] != nil {
		failureMsg = utils.StringPtr(evtData["failure_message"].(string))
	}

	if strings.HasPrefix(string(evt.Type), "payment_intent.") {
		piAmount := evtData["amount"].(float64)
		payment, err := server.repo.GetPaymentByPaymentIntentID(c, id)
		if err != nil {
			if errors.Is(err, repository.ErrRecordNotFound) {
				// If the payment is not found, we can ignore the event
				log.Info().Msgf("Payment not found for payment intent ID: %s", *id)
				c.JSON(http.StatusOK, createSuccessResponse(c, true, "", nil, nil))
				return
			}
			c.JSON(http.StatusBadRequest, createErrorResponse[bool](InternalServerErrorCode, "", err))
			return
		}
		updateTransactionStatus := repository.UpdatePaymentParams{
			ID:              payment.ID,
			GatewayChargeID: id,
			PaymentMethod: repository.NullPaymentMethod{
				PaymentMethod: payment.PaymentMethod,
				Valid:         true,
			},
			PaymentGateway:         payment.PaymentGateway,
			Amount:                 utils.GetPgNumericFromFloat((float64(piAmount) / 100)),
			GatewayPaymentIntentID: id,
			ErrorCode:              failureCode,
			ErrorMessage:           failureMsg,
		}

		switch evt.Type {
		case stripe.EventTypePaymentIntentSucceeded:
			updateTransactionStatus.Status = repository.NullPaymentStatus{
				PaymentStatus: repository.PaymentStatusSuccess,
				Valid:         true,
			}
			server.taskDistributor.SendOrderCreatedEmailTask(c,
				&worker.PayloadSendOrderCreatedEmailTask{
					PaymentID: payment.ID,
				},
				asynq.MaxRetry(10),
				asynq.ProcessIn(time.Second*3),
				asynq.Queue(worker.QueueDefault))
		case stripe.EventTypePaymentIntentCanceled:
			updateTransactionStatus.Status = repository.NullPaymentStatus{
				PaymentStatus: repository.PaymentStatusCancelled,
				Valid:         true,
			}
		case stripe.EventTypePaymentIntentPaymentFailed:
			updateTransactionStatus.Status = repository.NullPaymentStatus{
				PaymentStatus: repository.PaymentStatusFailed,
				Valid:         true,
			}

		default:
			log.Info().Msgf("Unhandled event type: %s", evt.Type)
			c.JSON(http.StatusOK, createSuccessResponse(c, true, "", nil, nil))
			return
		}
		err = server.repo.UpdatePayment(c, updateTransactionStatus)
	}
	if strings.HasPrefix(string(evt.Type), "charge.") {
		piAmount := evtData["amount"].(float64)
		var paymentIntentID *string
		if evtData["payment_intent"] != nil {
			paymentIntentID = utils.StringPtr(evtData["payment_intent"].(string))
		}
		payment, err := server.repo.GetPaymentByPaymentIntentID(c, paymentIntentID)
		if err != nil {
			if errors.Is(err, repository.ErrRecordNotFound) {
				// If the payment is not found, we can ignore the event
				log.Info().Msgf("Payment not found for payment intent ID: %s", *paymentIntentID)
				c.JSON(http.StatusOK, createSuccessResponse(c, true, "", nil, nil))
				return
			}
			c.JSON(http.StatusBadRequest, createErrorResponse[bool](InternalServerErrorCode, "", err))
			return
		}
		var createPaymentTransactionArg = repository.CreatePaymentTransactionParams{
			PaymentID:              payment.ID,
			Amount:                 utils.GetPgNumericFromFloat((float64(piAmount) / 100)),
			GatewayTransactionID:   id,
			Status:                 repository.PaymentStatusPending,
			GatewayResponseCode:    failureCode,
			GatewayResponseMessage: failureMsg,
		}

		switch evt.Type {
		case stripe.EventTypeChargeSucceeded:
			createPaymentTransactionArg.Status = repository.PaymentStatusSuccess
			updateTransactionStatus := repository.UpdatePaymentParams{
				ID:              payment.ID,
				GatewayChargeID: id,
			}
			err = server.repo.UpdatePayment(c, updateTransactionStatus)
			if err != nil {
				c.JSON(http.StatusBadRequest, createErrorResponse[bool](InternalServerErrorCode, "", err))
				return
			}
		case stripe.EventTypeChargeFailed:
			createPaymentTransactionArg.Status = repository.PaymentStatusFailed
		case stripe.EventTypeChargeRefunded:
			createPaymentTransactionArg.Status = repository.PaymentStatusRefunded
		case stripe.EventTypeChargeCaptured:
			createPaymentTransactionArg.Status = repository.PaymentStatusProcessing
		case stripe.EventTypeChargePending:
			createPaymentTransactionArg.Status = repository.PaymentStatusPending
		case stripe.EventTypeChargeUpdated:
			if evtData["status"] != nil {
				status := evtData["status"].(string)
				switch status {
				case "succeeded":
					createPaymentTransactionArg.Status = repository.PaymentStatusSuccess
				case "pending":
					createPaymentTransactionArg.Status = repository.PaymentStatusPending
				case "failed":
					createPaymentTransactionArg.Status = repository.PaymentStatusFailed
				case "canceled":
					createPaymentTransactionArg.Status = repository.PaymentStatusCancelled
				case "refunded":
					createPaymentTransactionArg.Status = repository.PaymentStatusRefunded
				case "processing":
					createPaymentTransactionArg.Status = repository.PaymentStatusProcessing
				}
			}
		}
		_, err = server.repo.CreatePaymentTransaction(c, createPaymentTransactionArg)
		if err != nil {
			c.JSON(http.StatusBadRequest, createErrorResponse[bool](InternalServerErrorCode, "", err))
			return
		}
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, true, "", nil, nil))
}
