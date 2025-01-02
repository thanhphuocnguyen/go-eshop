package repository

import (
	"context"

	"github.com/rs/zerolog/log"
)

type CancelOrderTxParams struct {
	OrderID int64
}

func (pg *pgRepo) CancelOrderTx(ctx context.Context, params CancelOrderTxParams) (err error) {
	pg.execTx(ctx, func(q *Queries) error {
		_, err := q.UpdateOrder(ctx, UpdateOrderParams{
			OrderID: params.OrderID,
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
			OrderID: params.OrderID,
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
