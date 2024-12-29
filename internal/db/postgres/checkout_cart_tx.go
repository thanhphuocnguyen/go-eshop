package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/sqlc"
	"github.com/thanhphuocnguyen/go-eshop/internal/util"
)

type CheckoutCartTxParams struct {
	UserID        int64              `json:"user_id"`
	CartID        int32              `json:"cart_id"`
	AddressID     int64              `json:"address_id"`
	PaymentMethod sqlc.PaymentMethod `json:"payment_method"`
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
		cart, err := s.GetCart(ctx, arg.UserID)
		if err != nil {
			log.Error().Err(err).Msg("GetCartDetail")
			return err
		}

		cartItems, err := s.GetCartItems(ctx, cart.ID)
		if err != nil {
			log.Error().Err(err).Msg("GetCartItems")
			return err
		}
		// create order items concurrently with goroutine and channel to handle error
		result.Items = make([]sqlc.OrderItem, len(cartItems))

		totalPrice := float64(0)
		for i, item := range cartItems {
			// create order orderItem for each cart orderItem
			if item.ProductStock < int32(item.Quantity) {
				log.Error().Msg("Product out of stock")
				return fmt.Errorf("product out of stock")
			}

			_, err = s.UpdateProduct(ctx, sqlc.UpdateProductParams{
				ID: item.ProductID,
				Stock: pgtype.Int4{
					Int32: item.ProductStock - int32(item.Quantity),
					Valid: true,
				},
			})

			if err != nil {
				log.Error().Err(err).Msg("DecreaseProductStock")
				return err
			}

			orderItem, err := s.CreateOrderItem(ctx, sqlc.CreateOrderItemParams{
				ProductID: item.ProductID,
				OrderID:   result.Order.ID,
				Quantity:  int32(item.Quantity),
				Price:     item.ProductPrice,
			})

			if err != nil {
				log.Error().Err(err).Msg("CreateOrderItem")
				return err
			}

			price, _ := orderItem.Price.Float64Value()
			totalPrice += price.Float64 * float64(orderItem.Quantity)
			result.Items[i] = orderItem
		}

		parsedPrice, _ := util.ParsePgTypeNumber(totalPrice)
		// create order
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
		result.Payment, err = s.CreatePaymentTransaction(ctx, sqlc.CreatePaymentTransactionParams{
			OrderID: result.Order.ID,
			Amount:  parsedPrice,
			Method:  arg.PaymentMethod,
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
