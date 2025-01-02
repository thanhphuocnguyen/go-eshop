package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/util"
)

type CheckoutCartTxParams struct {
	UserID         int64
	CartID         int32
	AddressID      int64
	CartItems      []GetCartItemsRow
	PaymentGateway PaymentGateway
	PaymentMethod  PaymentMethod
}

type PaymentInfo struct {
	ID             int32           `json:"id"`
	Amount         float64         `json:"amount"`
	PaymentMethod  PaymentMethod   `json:"payment_method"`
	PaymentGateWay *PaymentGateway `json:"payment_gateway"`
	PaymentStatus  PaymentStatus   `json:"payment_status"`
	CreatedAt      time.Time       `json:"created_at"`
}

type OrderItemDetail struct {
	ID          int64          `json:"id"`
	ProductID   int64          `json:"product_id"`
	ProductName string         `json:"product_name"`
	Quantity    int32          `json:"quantity"`
	ImageUrl    *string        `json:"image_url"`
	Price       pgtype.Numeric `json:"price"`
	SubTotal    float64        `json:"sub_total"`
}

type CheckoutCartTxResult struct {
	Order
	Items []OrderItemDetail `json:"items"`
}

func (s *pgRepo) CheckoutCartTx(ctx context.Context, arg CheckoutCartTxParams) (CheckoutCartTxResult, error) {
	var result CheckoutCartTxResult
	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		// create order
		totalPrice := float64(0)
		for _, item := range arg.CartItems {
			if item.ProductStock < int32(item.Quantity) {
				log.Error().Msg("Product out of stock")
				return fmt.Errorf("product out of stock")
			}
			price, _ := item.ProductPrice.Float64Value()
			totalPrice += price.Float64 * float64(item.Quantity)
		}

		result.Order, err = s.CreateOrder(ctx, CreateOrderParams{
			UserID:        arg.UserID,
			UserAddressID: arg.AddressID,
			TotalPrice:    util.GetPgNumericInt(totalPrice),
		})

		if err != nil {
			log.Error().Err(err).Msg("CreateOrder")
			return err
		}

		// create order items
		result.Items = make([]OrderItemDetail, len(arg.CartItems))
		for i, item := range arg.CartItems {
			itemCreated, err := s.CreateOrderItem(ctx, CreateOrderItemParams{
				ProductID: item.ProductID,
				OrderID:   result.Order.OrderID,
				Quantity:  int32(item.Quantity),
				Price:     item.ProductPrice,
			})
			if err != nil {
				log.Error().Err(err).Msg("CreateOrderItem")
				return err
			}
			subTotal := util.PgNumericToFloat64(item.ProductPrice)

			result.Items[i] = OrderItemDetail{
				ID:          itemCreated.OrderItemID,
				ProductID:   item.ProductID,
				ProductName: item.ProductName,
				ImageUrl:    &item.ImageUrl.String,
				Quantity:    int32(item.Quantity),
				Price:       item.ProductPrice,
				SubTotal:    subTotal * float64(item.Quantity),
			}

			// update product stock
			_, err = s.UpdateProduct(ctx, UpdateProductParams{
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
