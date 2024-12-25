package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/sqlc"
	"github.com/thanhphuocnguyen/go-eshop/internal/util"
)

type CheckoutCartTxParams struct {
	UserID      int64              `json:"user_id"`
	CartID      int64              `json:"cart_id"`
	AddressID   int64              `json:"address_id"`
	PaymentType sqlc.PaymentMethod `json:"payment_type"`
	IsCod       bool               `json:"is_cod"`
	//TODO: MakPaymentTransaction func() error
}
type CheckoutCartTxResult struct {
	sqlc.Order
	Items []sqlc.OrderItem
}

func (s *Postgres) CheckoutCartTx(ctx context.Context, arg CheckoutCartTxParams) (CheckoutCartTxResult, error) {
	var result CheckoutCartTxResult
	err := s.execTx(ctx, func(q *sqlc.Queries) error {
		var err error
		// get cart details
		cartDetails, err := s.GetCartDetail(ctx, arg.CartID)
		if err != nil {
			log.Error().Err(err).Msg("GetCartDetail")
			return err
		}
		if len(cartDetails) == 0 {
			log.Error().Msg("cart is empty")
			return fmt.Errorf("cart is empty")
		}
		// create order items concurrently with goroutine and channel to handle error
		var orderItems []sqlc.OrderItem
		totalPrice := float64(0)
		for _, item := range cartDetails {
			// create order item for each cart item
			item, err := s.CreateOrderItem(ctx, sqlc.CreateOrderItemParams{
				ProductID: item.Product.ID,
				OrderID:   result.Order.ID,
				Quantity:  int32(item.CartItem.Quantity),
				Price:     item.Product.Price,
			})
			price, _ := item.Price.Float64Value()

			totalPrice += price.Float64 * float64(item.Quantity)
			if err != nil {
				log.Error().Err(err).Msg("CreateOrderItem")
				return err
			}
			orderItems = append(orderItems, item)
		}
		result.Items = orderItems
		// set cart checkout at time
		err = s.SetCartCheckoutAt(ctx, sqlc.SetCartCheckoutAtParams{
			ID:         arg.CartID,
			CheckoutAt: util.GetPgTypeTimestamp(time.Now()),
		})
		if err != nil {
			log.Error().Err(err).Msg("SetCartCheckoutAt")
			return err
		}
		parsedPrice, _ := util.ParsePgNumeric(totalPrice)

		result.Order, err = s.CreateOrder(ctx, sqlc.CreateOrderParams{
			UserID:        arg.UserID,
			UserAddressID: arg.AddressID,
			IsCod:         arg.IsCod,
			CartID:        arg.CartID,
			TotalPrice:    parsedPrice,
		})

		if err != nil {
			log.Error().Err(err).Msg("CreateOrder")
			return err
		}

		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("execTx")
	}

	return result, err
}
