package postgres

import (
	"context"
	"fmt"
	"time"

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
}

type PaymentInfo struct {
	ID        int32              `json:"id"`
	Amount    float64            `json:"amount"`
	Method    sqlc.PaymentMethod `json:"method"`
	Status    sqlc.PaymentStatus `json:"status"`
	CreatedAt time.Time          `json:"created_at"`
}

type CheckoutCartTxResult struct {
	sqlc.Order
	Items   []sqlc.OrderItem `json:"items"`
	Payment PaymentInfo      `json:"payment"`
}

func (s *Postgres) CheckoutCartTx(ctx context.Context, arg CheckoutCartTxParams) (CheckoutCartTxResult, error) {
	var result CheckoutCartTxResult
	err := s.execTx(ctx, func(q *sqlc.Queries) error {
		var err error

		cartItems, err := s.GetCartItems(ctx, arg.CartID)
		if err != nil {
			log.Error().Err(err).Msg("GetCartItems")
			return err
		}

		if len(cartItems) == 0 {
			log.Error().Msg("Cart is empty")
			return fmt.Errorf("cart is empty")
		}

		// create order
		totalPrice := float64(0)
		for _, item := range cartItems {
			if item.ProductStock < int32(item.Quantity) {
				log.Error().Msg("Product out of stock")
				return fmt.Errorf("product out of stock")
			}
			price, _ := item.ProductPrice.Float64Value()
			totalPrice += price.Float64 * float64(item.Quantity)
		}

		result.Order, err = s.CreateOrder(ctx, sqlc.CreateOrderParams{
			UserID:        arg.UserID,
			UserAddressID: arg.AddressID,
			TotalPrice:    util.GetPgNumericInt(totalPrice),
		})

		if err != nil {
			log.Error().Err(err).Msg("CreateOrder")
			return err
		}

		// create order items
		result.Items = make([]sqlc.OrderItem, len(cartItems))
		for i, item := range cartItems {
			orderItem, err := s.CreateOrderItem(ctx, sqlc.CreateOrderItemParams{
				ProductID: item.ProductID,
				OrderID:   result.Order.OrderID,
				Quantity:  int32(item.Quantity),
				Price:     item.ProductPrice,
			})
			if err != nil {
				log.Error().Err(err).Msg("CreateOrderItem")
				return err
			}
			result.Items[i] = orderItem

			// update product stock
			_, err = s.UpdateProduct(ctx, sqlc.UpdateProductParams{
				ProductID: item.ProductID,
				Stock: pgtype.Int4{
					Int32: item.ProductStock - int32(item.Quantity),
					Valid: true,
				},
			})
			if err != nil {
				log.Error().Err(err).Msg("DecreaseProductStock")
				return err
			}
		}

		// create payment transaction
		paymentRow, err := s.CreatePaymentTransaction(ctx, sqlc.CreatePaymentTransactionParams{
			OrderID: result.Order.OrderID,
			Amount:  util.GetPgNumericInt(totalPrice),
			Method:  arg.PaymentMethod,
		})
		if err != nil {
			log.Error().Err(err).Msg("CreatePaymentTransaction")
			return err
		}

		result.Payment = PaymentInfo{
			ID:        paymentRow.ID,
			Amount:    totalPrice,
			Method:    arg.PaymentMethod,
			Status:    sqlc.PaymentStatusPending,
			CreatedAt: paymentRow.CreatedAt.Time,
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
