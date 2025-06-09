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

type CreatePaymentResult struct {
	PaymentID       string        `json:"paymentId"`
	ClientSecret    *string       `json:"clientSecret,omitempty"`
	PaymentIntentID string        `json:"paymentIntentId"`
	TotalPrice      float64       `json:"totalPrice"`
	OrderID         uuid.UUID     `json:"orderId"`
	Status          PaymentStatus `json:"status"`
}

type CreatePaymentArgs struct {
	Amount  float64
	Email   string
	OrderID uuid.UUID
}

type CheckoutCartTxArgs struct {
	UserID                uuid.UUID
	CartID                uuid.UUID
	CustomerInfo          CustomerInfoTxArgs
	CreateOrderItemParams []CreateBulkOrderItemsParams
	TotalPrice            float64
	DiscountPrice         float64
	DiscountID            *uuid.UUID
	ShippingAddress       ShippingAddressSnapshot
	PaymentMethod         PaymentMethod
	PaymentGateway        *string
	CreatePaymentFn       func(orderID uuid.UUID, paymentMethod PaymentMethod) (paymentIntentID string, clientSecret *string, err error)
}

func (s *pgRepo) CheckoutCartTx(ctx context.Context, arg CheckoutCartTxArgs) (CreatePaymentResult, error) {
	var result CreatePaymentResult
	err := s.execTx(ctx, func(q *Queries) (err error) {
		params := CreateOrderParams{
			CustomerID:      arg.UserID,
			ShippingAddress: arg.ShippingAddress,
			TotalPrice:      utils.GetPgNumericFromFloat(arg.TotalPrice),
			CustomerEmail:   arg.CustomerInfo.Email,
			CustomerName:    arg.CustomerInfo.FullName,
			CustomerPhone:   arg.CustomerInfo.Phone,
		}

		order, err := s.CreateOrder(ctx, params)
		if err != nil {
			log.Error().Err(err).Msg("CreateOrder")
			return err
		}
		paymentAmount := arg.TotalPrice - arg.DiscountPrice
		if paymentAmount < 0 {
			paymentAmount = 0
		}
		result.TotalPrice = paymentAmount

		createPaymentArgs := CreatePaymentParams{
			OrderID: order.ID,
			Amount:  utils.GetPgNumericFromFloat(paymentAmount),
		}

		createPaymentArgs.Method = arg.PaymentMethod
		createPaymentArgs.Gateway = arg.PaymentGateway

		if arg.DiscountPrice != 0 {
			if arg.DiscountPrice > arg.TotalPrice {
				arg.DiscountPrice = arg.TotalPrice
			}
			_, err := q.InsertOrderDiscount(ctx, InsertOrderDiscountParams{
				OrderID:        order.ID,
				DiscountID:     *arg.DiscountID,
				DiscountAmount: utils.GetPgNumericFromFloat(arg.DiscountPrice),
			})
			if err != nil {
				log.Error().Err(err).Msg("GetDiscountByCode")
				return err
			}
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
		if paymentAmount > 0 {
			paymentIntentID, clientSecret, err := arg.CreatePaymentFn(order.ID, createPaymentArgs.Method)
			createPaymentArgs.PaymentIntentID = &paymentIntentID
			if err != nil {
				log.Error().Err(err).Msg("CreatePaymentFn")
				return err
			}
			result.ClientSecret = clientSecret
			result.PaymentIntentID = paymentIntentID
		} else {
			createPaymentArgs.PaymentIntentID = nil
			createPaymentArgs.Status = PaymentStatusSuccess

		}

		// create payment transaction
		payment, err := q.CreatePayment(ctx, createPaymentArgs)
		if err != nil {
			log.Error().Err(err).Msg("CreatePayment")
			return err
		}

		result.PaymentID = payment.ID.String()

		return nil
	})

	if err != nil {
		log.Error().Err(err).Msg("execTx")
	}

	return result, err
}
