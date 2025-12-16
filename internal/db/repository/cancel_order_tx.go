package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
)

type CancelOrderTxArgs struct {
	OrderID                 uuid.UUID
	CancelPaymentFromMethod func(ID string, method string) error
}

func (repo *pgRepo) CancelOrderTx(ctx context.Context, args CancelOrderTxArgs) (uuid.UUID, error) {
	var orderID = args.OrderID
	err := repo.execTx(ctx, func(q *Queries) error {
		// cancel payment
		payment, err := q.GetPaymentByOrderID(ctx, args.OrderID)
		if !errors.Is(err, ErrRecordNotFound) {
			log.Error().Err(err).Msg("GetPaymentTransactionByOrderID")
		}
		paymentMethod, methodErr := q.GetPaymentMethodByID(ctx, payment.PaymentMethodID)
		if err == nil && methodErr != nil {
			log.Error().Err(err).Msg("GetPaymentMethodByID")
			return err
		}
		// if payment is not found, we don't need to cancel it
		if !errors.Is(err, ErrRecordNotFound) && payment.Gateway != nil {
			if payment.Status == PaymentStatusSuccess {
				return errors.New("payment is already successful, need to refund")
			}

			// cancel payment from gateway if it's not cancelled yet
			if args.CancelPaymentFromMethod != nil && payment.PaymentIntentID != nil {
				err = args.CancelPaymentFromMethod(*payment.PaymentIntentID, paymentMethod.Code)
				if err != nil {
					log.Error().Err(err).Msg("CancelPaymentFromGateway")
					return err
				}
				err := q.UpdatePaymentTransaction(ctx, UpdatePaymentTransactionParams{
					ID: payment.ID,
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
			ID: args.OrderID,
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
		// TODO: #1 implement UpdateProductStock
		for _, item := range orderItems {
			_, err := q.UpdateProductStock(ctx, UpdateProductStockParams{
				ID:    item.VariantID,
				Stock: -int32(item.Quantity),
			})

			if err != nil {
				log.Error().Err(err).Msg("UpdateProductStock")
				return err
			}
		}
		return nil
	})

	return orderID, err
}
