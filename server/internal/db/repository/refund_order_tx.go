package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
)

type RefundOrderTxArgs struct {
	OrderID                  uuid.UUID
	RefundPaymentFromGateway func(paymentID string, gateway PaymentGateway) (string, error)
}

func (pg *pgRepo) RefundOrderTx(ctx context.Context, args RefundOrderTxArgs) (err error) {
	pg.execTx(ctx, func(q *Queries) error {
		// refund payment
		payment, err := q.GetPaymentTransactionByOrderID(ctx, args.OrderID)
		if err != nil {
			log.Error().Err(err).Msg("GetPaymentTransactionByOrderID")
			return err
		}
		// if payment is not found, we don't need to refund it
		if payment.PaymentGateway.Valid {
			if payment.Status != PaymentStatusSuccess {
				return errors.New("payment is not successful, can't refund")
			}

			// refund payment from gateway if it's not refunded yet
			if args.RefundPaymentFromGateway != nil {
				refundID, err := args.RefundPaymentFromGateway(payment.ID, payment.PaymentGateway.PaymentGateway)
				if err != nil {
					log.Error().Err(err).Msg("RefundPaymentFromGateway")
					return err
				}
				updateParams := UpdatePaymentTransactionParams{
					ID: payment.ID,
					Status: NullPaymentStatus{
						PaymentStatus: PaymentStatusRefunded,
						Valid:         true,
					},
				}
				if refundID != "" {
					updateParams.RefundID = pgtype.Text{
						String: refundID,
						Valid:  true,
					}
				}
				err = q.UpdatePaymentTransaction(ctx, updateParams)
				if err != nil {
					log.Error().Err(err).Msg("UpdatePaymentTransaction")
					return err
				}
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
