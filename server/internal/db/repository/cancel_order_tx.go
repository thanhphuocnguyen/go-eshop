package repository

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
)

type CancelOrderTxArgs struct {
	OrderID                  int64
	CancelPaymentFromGateway func(paymentID string, gateway PaymentGateway) error
}

func (pg *pgRepo) CancelOrderTx(ctx context.Context, args CancelOrderTxArgs) (err error) {
	pg.execTx(ctx, func(q *Queries) error {
		// cancel payment
		payment, err := q.GetPaymentTransactionByOrderID(ctx, args.OrderID)
		if !errors.Is(err, ErrRecordNotFound) {
			log.Error().Err(err).Msg("GetPaymentTransactionByOrderID")
		}
		// if payment is not found, we don't need to cancel it
		if err == nil && payment.PaymentGateway.Valid {
			if payment.Status == PaymentStatusSuccess {
				return errors.New("payment is already successful, need to refund")
			}

			// cancel payment from gateway if it's not cancelled yet
			if args.CancelPaymentFromGateway != nil {
				err = args.CancelPaymentFromGateway(payment.PaymentID, payment.PaymentGateway.PaymentGateway)
				if err != nil {
					log.Error().Err(err).Msg("CancelPaymentFromGateway")
					return err
				}
				err := q.UpdatePaymentTransaction(ctx, UpdatePaymentTransactionParams{
					PaymentID: payment.PaymentID,
					Status: NullPaymentStatus{
						PaymentStatus: PaymentStatusCancelled,
						Valid:         true,
					},
				})
				if err != nil {
					log.Error().Err(err).Msg("UpdatePaymentTransaction")
					return err
				}
			}
		}

		// cancel order
		_, err = q.UpdateOrder(ctx, UpdateOrderParams{
			OrderID: args.OrderID,
			Status: NullOrderStatus{
				OrderStatus: OrderStatusCancelled,
				Valid:       true,
			},
		})

		if err != nil {
			log.Error().Err(err).Msg("UpdateOrder")
			return err
		}

		// refilling stock
		orderItems, err := q.ListOrderItems(ctx, ListOrderItemsParams{
			OrderID: args.OrderID,
		})

		if err != nil {
			log.Error().Err(err).Msg("ListOrderItems")
			return err
		}
		for _, item := range orderItems {
			err = q.UpdateProductStock(ctx, UpdateProductStockParams{
				Stock:     item.Quantity,
				ProductID: item.ProductID,
			})
			if err != nil {
				log.Error().Err(err).Msg("UpdateProductStock")
				return err
			}
		}
		return nil
	})

	return err
}
