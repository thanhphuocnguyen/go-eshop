package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/sqlc"
)

type CheckoutCartParams struct {
	UserID int64 `json:"user_id"`
	sqlc.SetCartCheckoutAtParams
	sqlc.CreateOrderParams
	//TODO: MakPaymentTransaction func() error
}

type CheckoutCartTxResult struct {
	Order sqlc.Order `json:"order"`
}

func (s *Postgres) CheckoutCartTx(ctx context.Context, arg CheckoutCartParams) (CheckoutCartTxResult, error) {
	var result CheckoutCartTxResult
	err := s.execTx(ctx, func(q *sqlc.Queries) error {
		var err error
		errChan := make(chan error)
		result.Order, err = s.CreateOrder(ctx, sqlc.CreateOrderParams{
			UserID:      arg.UserID,
			PaymentType: arg.PaymentType,
			IsCod:       arg.IsCod,
		})

		if err != nil {
			log.Error().Err(err).Msg("CreateOrder")
			return err
		}

		cartDetails, err := s.GetCartDetail(ctx, arg.ID)
		if err != nil {
			log.Error().Err(err).Msg("GetCartDetail")
			return err
		}

		go func() {
			for _, item := range cartDetails {
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
		}()

		for err := range errChan {
			if err != nil {
				log.Error().Err(err).Msg("CreateOrderItem")
				return err
			}
		}

		s.SetCartCheckoutAt(ctx, sqlc.SetCartCheckoutAtParams{
			ID: arg.ID,
			CheckedOutAt: pgtype.Timestamptz{
				Time:  time.Now(),
				Valid: true,
			},
		})
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("execTx")
	}

	return result, err
}
