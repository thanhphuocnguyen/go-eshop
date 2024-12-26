package postgres

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/sqlc"
	"github.com/thanhphuocnguyen/go-eshop/internal/util"
)

type CheckoutCartTxParams struct {
	UserID        int64              `json:"user_id"`
	CartID        int32              `json:"cart_id"`
	AddressID     int64              `json:"address_id"`
	PaymentMethod sqlc.PaymentMethod `json:"payment_method"`
	IsCod         bool               `json:"is_cod"`
	//TODO: MakPaymentTransaction func() error
}
type CheckoutCartTxResult struct {
	sqlc.Order
	Items   []sqlc.OrderItem
	Payment sqlc.Payment
}

func (s *Postgres) CheckoutCartTx(ctx context.Context, arg CheckoutCartTxParams) (CheckoutCartTxResult, error) {
	var result CheckoutCartTxResult
	err := s.execTx(ctx, func(q *sqlc.Queries) error {
		var err error
		// get cart details
		_, err = s.GetCart(ctx, arg.UserID)
		if err != nil {
			log.Error().Err(err).Msg("GetCartDetail")
			return err
		}

		cartItems, err := s.GetCartItems(ctx, arg.CartID)
		if err != nil {
			log.Error().Err(err).Msg("GetCartItems")
			return err
		}
		// create order items concurrently with goroutine and channel to handle error
		var orderItems []sqlc.OrderItem
		totalPrice := float64(0)
		for _, item := range cartItems {
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

		parsedPrice, _ := util.ParsePgNumeric(totalPrice)

		result.Order, err = s.CreateOrder(ctx, sqlc.CreateOrderParams{
			UserID:        arg.UserID,
			UserAddressID: arg.AddressID,
			TotalPrice:    parsedPrice,
		})

		if err != nil {
			log.Error().Err(err).Msg("CreateOrder")
			return err
		}

		// create payment transaction
		payment, err := s.CreatePaymentTransaction(ctx, sqlc.CreatePaymentTransactionParams{
			OrderID: result.Order.ID,
			Amount:  parsedPrice, Method: arg.PaymentMethod,
		})
		if err != nil {
			log.Error().Err(err).Msg("CreatePaymentTransaction")
			return err
		}
		result.Payment = payment

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
