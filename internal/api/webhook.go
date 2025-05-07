package api

import (
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
	// unmarshal the event data

	evtData := evt.Data.Object
	var chargeID *string
	if evtData["id"] != nil {
		chargeID = StringPtr(evtData["id"].(string))
	}
	var pID *string
	if evtData["payment_intent"] != nil {
		pID = StringPtr(evtData["payment_intent"].(string))
	}
	log.Info().Interface("evt type", evt.Type).Msg("Received stripe event")
	payment, err := server.repo.GetPaymentByPaymentIntentID(c, pID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			// If the payment is not found, we can ignore the event
			log.Info().Msgf("Payment not found for payment intent ID: %s", *pID)
			c.JSON(http.StatusOK, createSuccessResponse(c, true, "", nil, nil))
			return
		}
		c.JSON(http.StatusBadRequest, createErrorResponse[bool](InternalServerErrorCode, "", err))
		return
	}
	piAmount := evtData["amount"].(float64)
	var failureCode *string
	var failureMsg *string
	if evtData["failure_code"] != nil {
		failureCode = StringPtr(evtData["failure_code"].(string))
	}
	if evtData["failure_message"] != nil {
		failureMsg = StringPtr(evtData["failure_message"].(string))
	}
	var createPaymentTransactionArg = repository.CreatePaymentTransactionParams{
		ID:                     uuid.New(),
		PaymentID:              payment.ID,
		Amount:                 utils.GetPgNumericFromFloat((float64(piAmount) / 100)),
		GatewayTransactionID:   chargeID,
		Status:                 repository.PaymentStatusPending,
		GatewayResponseCode:    failureCode,
		GatewayResponseMessage: failureMsg,
	}

	updateTransactionStatus := repository.UpdatePaymentParams{
		ID:              payment.ID,
		GatewayChargeID: chargeID,
		PaymentMethod: repository.NullPaymentMethod{
			PaymentMethod: payment.PaymentMethod,
			Valid:         true,
		},
		PaymentGateway:         payment.PaymentGateway,
		Amount:                 utils.GetPgNumericFromFloat((float64(piAmount) / 100)),
		GatewayPaymentIntentID: pID,
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
				PaymentID: payment.ID.String(),
				OrderID:   payment.OrderID,
			},
			asynq.MaxRetry(3),
			asynq.Queue("email"),
			asynq.ProcessIn(time.Second*5))
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
	case stripe.EventTypeChargeFailed:
		createPaymentTransactionArg.Status = repository.PaymentStatusFailed
	case stripe.EventTypeChargeSucceeded:
		createPaymentTransactionArg.Status = repository.PaymentStatusSuccess
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
				updateTransactionStatus.Status = repository.NullPaymentStatus{
					PaymentStatus: repository.PaymentStatusSuccess,
					Valid:         true,
				}
			case "pending":
				createPaymentTransactionArg.Status = repository.PaymentStatusPending
				updateTransactionStatus.Status = repository.NullPaymentStatus{
					PaymentStatus: repository.PaymentStatusPending,
					Valid:         true,
				}
			case "failed":
				createPaymentTransactionArg.Status = repository.PaymentStatusFailed
				updateTransactionStatus.Status = repository.NullPaymentStatus{
					PaymentStatus: repository.PaymentStatusFailed,
					Valid:         true,
				}
			case "canceled":
				createPaymentTransactionArg.Status = repository.PaymentStatusCancelled
				updateTransactionStatus.Status = repository.NullPaymentStatus{
					PaymentStatus: repository.PaymentStatusCancelled,
					Valid:         true,
				}
			case "refunded":
				createPaymentTransactionArg.Status = repository.PaymentStatusRefunded
				updateTransactionStatus.Status = repository.NullPaymentStatus{
					PaymentStatus: repository.PaymentStatusRefunded,
					Valid:         true,
				}
			case "processing":
				createPaymentTransactionArg.Status = repository.PaymentStatusProcessing
				updateTransactionStatus.Status = repository.NullPaymentStatus{
					PaymentStatus: repository.PaymentStatusProcessing,
					Valid:         true,
				}
			}
		}
	default:
		log.Info().Msgf("Unhandled event type: %s", evt.Type)
		c.JSON(http.StatusBadGateway, createSuccessResponse(c, true, "", nil, nil))
		return
	}
	_, err = server.repo.CreatePaymentTransaction(c, createPaymentTransactionArg)
	if err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[bool]("duplicate_transaction", "", err))
		return
	}
	err = server.repo.UpdatePayment(c, updateTransactionStatus)
	if err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[bool](InternalServerErrorCode, "", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, true, "", nil, nil))
}
