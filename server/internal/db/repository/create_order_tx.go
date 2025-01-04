package repository

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/util"
)

type CreateOrderTxParams struct {
	UserID                int64
	CartID                int32
	PaymentMethod         PaymentMethod
	PaymentGateway        PaymentGateway
	AddressID             int64
	PaymentID             string
	CreateOrderItemParams []CreateOrderItemParams
	TotalPrice            float64
}

func (s *pgRepo) CreateOrderTx(ctx context.Context, arg CreateOrderTxParams) (int64, error) {
	var result int64
	err := s.execTx(ctx, func(q *Queries) (err error) {

		order, err := s.CreateOrder(ctx, CreateOrderParams{
			UserID:        arg.UserID,
			UserAddressID: arg.AddressID,
			TotalPrice:    util.GetPgNumericFromFloat(arg.TotalPrice),
		})
		result = order.OrderID
		if err != nil {
			log.Error().Err(err).Msg("CreateOrder")
			return err
		}

		// create payment transaction
		_, err = q.CreatePaymentTransaction(ctx, CreatePaymentTransactionParams{
			PaymentID:     arg.PaymentID,
			OrderID:       order.OrderID,
			Amount:        util.GetPgNumericFromFloat(arg.TotalPrice),
			PaymentMethod: arg.PaymentMethod,
			PaymentGateway: NullPaymentGateway{
				PaymentGateway: arg.PaymentGateway,
				Valid:          true,
			},
		})

		if err != nil {
			log.Error().Err(err).Msg("CreatePaymentTransaction")
			return err
		}

		// clear cart
		err = s.ClearCart(ctx, arg.CartID)
		if err != nil {
			log.Error().Err(err).Msg("ClearCart")
			return err
		}

		return nil
	})

	if err != nil {
		log.Error().Err(err).Msg("execTx")
	}

	return result, err
}
