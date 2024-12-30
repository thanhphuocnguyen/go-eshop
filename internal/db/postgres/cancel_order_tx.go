package postgres

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/sqlc"
)

type CancelOrderTxParams struct {
	OrderID int64
}

func (pg *Postgres) CancelOrderTx(ctx context.Context, params CancelOrderTxParams) (err error) {
	pg.execTx(ctx, func(q *sqlc.Queries) error {
		_, err := q.UpdateOrder(ctx, sqlc.UpdateOrderParams{
			ID: params.OrderID,
			Status: sqlc.NullOrderStatus{
				OrderStatus: sqlc.OrderStatusCancelled,
				Valid:       true,
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("UpdateOrder")
			return err
		}

		// refilling stock
		orderItems, err := q.ListOrderItems(ctx, sqlc.ListOrderItemsParams{
			OrderID: params.OrderID,
		})
		if err != nil {
			log.Error().Err(err).Msg("ListOrderItems")
			return err
		}
		for _, item := range orderItems {
			err = q.UpdateProductStock(ctx, sqlc.UpdateProductStockParams{
				Stock: item.Quantity,
				ID:    item.ProductID,
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
