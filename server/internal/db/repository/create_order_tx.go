package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

type CustomerInfoTxArgs struct {
	FullName string
	Email    string
	Phone    string
}

type CreateOrderTxArgs struct {
	UserID                uuid.UUID
	CartID                uuid.UUID
	AddressID             int64
	CustomerInfo          CustomerInfoTxArgs
	CreateOrderItemParams []CreateBulkOrderItemsParams
	TotalPrice            float64
}

func (s *pgRepo) CreateOrderTx(ctx context.Context, arg CreateOrderTxArgs) (uuid.UUID, error) {
	var result uuid.UUID
	err := s.execTx(ctx, func(q *Queries) (err error) {
		params := CreateOrderParams{
			ID:            uuid.New(),
			CustomerID:    arg.UserID,
			UserAddressID: arg.AddressID,
			TotalPrice:    utils.GetPgNumericFromFloat(arg.TotalPrice),
			CustomerEmail: arg.CustomerInfo.Email,
			CustomerName:  arg.CustomerInfo.FullName,
			CustomerPhone: utils.GetPgTypeText(arg.CustomerInfo.Phone),
		}

		order, err := s.CreateOrder(ctx, params)

		result = order.ID
		if err != nil {
			log.Error().Err(err).Msg("CreateOrder")
			return err
		}

		for i := range arg.CreateOrderItemParams {
			arg.CreateOrderItemParams[i].OrderID = order.ID

		}
		_, err = q.CreateBulkOrderItems(ctx, arg.CreateOrderItemParams)

		if err != nil {
			log.Error().Err(err).Msg("CreateOrderItem")
			return err
		}

		// clear cart
		err = s.CheckoutCart(ctx, CheckoutCartParams{
			OrderID: utils.GetPgTypeUUID(order.ID),
			ID:      arg.CartID,
		})

		if err != nil {
			log.Error().Err(err).Msg("CheckoutCart")
			return err
		}

		return nil
	})

	if err != nil {
		log.Error().Err(err).Msg("execTx")
	}

	return result, err
}
