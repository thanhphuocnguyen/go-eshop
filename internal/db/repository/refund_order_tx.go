package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type RefundOrderTxArgs struct {
	OrderID                 uuid.UUID
	RefundPaymentFromMethod func(paymentID string, method PaymentMethod) (string, error)
}

func (pg *pgRepo) RefundOrderTx(ctx context.Context, args RefundOrderTxArgs) (err error) {
	pg.execTx(ctx, func(q *Queries) error {
		// refund payment
		payment, err := q.GetPaymentByOrderID(ctx, args.OrderID)
		if err != nil {
			log.Error().Err(err).Msg("GetPaymentTransactionByOrderID")
			return err
		}
		// if payment is not found, we don't need to refund it
		if payment.Status != PaymentStatusSuccess {
			return errors.New("payment is not successful, can't refund")
		}

		// refund payment from gateway if it's not refunded yet
		if args.RefundPaymentFromMethod != nil {
			refundID, err := args.RefundPaymentFromMethod(*payment.PaymentIntentID, payment.Method)
			if err != nil {
				log.Error().Err(err).Msg("RefundPaymentFromGateway")
				return err
			}
			updateParams := UpdatePaymentParams{
				ID: payment.ID,
				Status: NullPaymentStatus{
					PaymentStatus: PaymentStatusRefunded,
					Valid:         true,
				},
			}
			if refundID != "" {
				updateParams.RefundID = &refundID
			}
			err = q.UpdatePayment(ctx, updateParams)
			if err != nil {
				log.Error().Err(err).Msg("UpdatePaymentTransaction")
				return err
			}
		}

		// refund order
		_, err = q.UpdateOrder(ctx, UpdateOrderParams{
			ID: args.OrderID,
			Status: NullOrderStatus{
				OrderStatus: OrderStatusRefunded,
				Valid:       true,
			},
		})

		if err != nil {
			log.Error().Err(err).Msg("UpdateOrder")
			return err
		}
		return nil
	})
	return
}
