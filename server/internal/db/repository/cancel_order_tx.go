package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

type CancelOrderTxArgs struct {
	OrderID                  uuid.UUID
	CancelPaymentFromGateway func(paymentID string, gateway PaymentGateway) error
}

func (pg *pgRepo) CancelOrderTx(ctx context.Context, args CancelOrderTxArgs) (order Order, err error) {
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
		order, err = q.UpdateOrder(ctx, UpdateOrderParams{
			OrderID: args.OrderID,
			Status: NullOrderStatus{
				OrderStatus: OrderStatusCancelled,
				Valid:       true,
			},
			CancelledAt: pgtype.Timestamptz{
				Time:  time.Now(),
				Valid: true,
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
		// TODO: implement UpdateProductStock
		for _, item := range orderItems {
			_, err = q.UpdateVariant(ctx, UpdateVariantParams{
				VariantID:     item.VariantID,
				StockQuantity: utils.GetPgTypeInt4(int32(item.Quantity)),
			})

			if err != nil {
				log.Error().Err(err).Msg("UpdateProductStock")
				return err
			}
		}
		return nil
	})

	return order, err
}
