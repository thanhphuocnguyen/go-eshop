package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/sqlc"
)

type CheckoutCartTxParams struct {
	UserID int64 `json:"user_id"`
	CartID int64 `json:"cart_id"`
	sqlc.CreateOrderParams
	//TODO: MakPaymentTransaction func() error
}

type CheckoutCartTxResult struct {
	Order sqlc.Order `json:"order"`
}

func (s *Postgres) CheckoutCartTx(ctx context.Context, arg CheckoutCartTxParams) (CheckoutCartTxResult, error) {
	var result CheckoutCartTxResult
	err := s.execTx(ctx, func(q *sqlc.Queries) error {
		var err error
		errChan := make(chan error)
		// create order
		result.Order, err = s.CreateOrder(ctx, sqlc.CreateOrderParams{
			UserID:      arg.UserID,
			PaymentType: arg.PaymentType,
			IsCod:       arg.IsCod,
		})

		if err != nil {
			log.Error().Err(err).Msg("CreateOrder")
			return err
		}
		// get cart details
		cartDetails, err := s.GetCartDetail(ctx, arg.CartID)
		if err != nil {
			log.Error().Err(err).Msg("GetCartDetail")
			return err
		}
		// create order items concurrently with goroutine and channel to handle error
		go func() {
			for _, item := range cartDetails {
				// create order item for each cart item
				_, err = s.CreateOrderItem(ctx, sqlc.CreateOrderItemParams{
					ProductID: item.Product.ID,
					OrderID:   result.Order.ID,
					Quantity:  int32(item.CartItem.Quantity),
					Price:     item.Product.Price,
				})

				if err != nil {
					errChan <- err
					return
				}
			}
			// close the channel when all order items are created
			close(errChan)
		}()

		// wait for all order items to be created
		for err := range errChan {
			if err != nil {
				log.Error().Err(err).Msg("CreateOrderItem")
				return err
			}
		}

		// set cart checkout at time
		err = s.SetCartCheckoutAt(ctx, sqlc.SetCartCheckoutAtParams{
			ID: arg.CartID,
			CheckedOutAt: pgtype.Timestamptz{
				Time:  time.Now(),
				Valid: true,
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("SetCartCheckoutAt")
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("execTx")
	}

	return result, err
}
