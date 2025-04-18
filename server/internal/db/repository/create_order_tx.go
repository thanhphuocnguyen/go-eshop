package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

type CreateOrderTxParams struct {
	UserID                uuid.UUID
	CartID                uuid.UUID
	PaymentMethod         PaymentMethod
	PaymentGateway        PaymentGateway
	AddressID             int64
	PaymentID             string
	CreateOrderItemParams []CreateOrderItemParams
	TotalPrice            float64
}

func (s *pgRepo) CreateOrderTx(ctx context.Context, arg CreateOrderTxParams) (uuid.UUID, error) {
	var result uuid.UUID
	err := s.execTx(ctx, func(q *Queries) (err error) {
		order, err := s.CreateOrder(ctx, CreateOrderParams{
			ID:            uuid.New(),
			CustomerID:    arg.UserID,
			UserAddressID: arg.AddressID,
			TotalPrice:    utils.GetPgNumericFromFloat(arg.TotalPrice),
		})

		result = order.ID
		if err != nil {
			log.Error().Err(err).Msg("CreateOrder")
			return err
		}

		for _, createOrderItemParam := range arg.CreateOrderItemParams {
			createOrderItemParam.OrderID = order.ID
			_, err = q.CreateOrderItem(ctx, createOrderItemParam)
			if err != nil {
				log.Error().Err(err).Msg("CreateOrderItem")
				return err
			}
		}

		// create payment transaction
		_, err = q.CreatePaymentTransaction(ctx, CreatePaymentTransactionParams{
			ID:            arg.PaymentID,
			OrderID:       order.ID,
			Amount:        utils.GetPgNumericFromFloat(arg.TotalPrice),
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
