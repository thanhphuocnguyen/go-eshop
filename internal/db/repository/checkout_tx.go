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
	DiscountIDs           []uuid.UUID
	ShippingAddress       ShippingAddressSnapshot
	PaymentMethodID       uuid.UUID
	PaymentGateway        *string
	CreatePaymentFn       func(orderID uuid.UUID, method string) (paymentIntentID string, clientSecret *string, err error)
}

func (s *pgRepo) CheckoutCartTx(ctx context.Context, arg CheckoutCartTxArgs) (CreatePaymentResult, error) {
	var result CreatePaymentResult
	err := s.execTx(ctx, func(q *Queries) (err error) {
		params := CreateOrderParams{
			UserID:          arg.UserID,
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
			OrderID:         order.ID,
			PaymentMethodID: arg.PaymentMethodID,
			Amount:          utils.GetPgNumericFromFloat(paymentAmount),
		}

		createPaymentArgs.Gateway = arg.PaymentGateway

		if arg.DiscountPrice != 0 {
			if arg.DiscountPrice > arg.TotalPrice {
				arg.DiscountPrice = arg.TotalPrice
			}
			for _, id := range arg.DiscountIDs {
				_, err := q.AddDiscountUsage(ctx, AddDiscountUsageParams{
					OrderID:        order.ID,
					DiscountID:     id,
					DiscountAmount: utils.GetPgNumericFromFloat(arg.DiscountPrice),
					UserID:         arg.UserID,
				})
				if err != nil {
					log.Error().Err(err).Msg("AddDiscountUsage")
					return err
				}
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
			method, err := s.GetPaymentMethodByID(ctx, arg.PaymentMethodID)
			if err != nil {
				log.Error().Err(err).Msg("GetPaymentMethodByID")
				return err
			}
			paymentIntentID, clientSecret, err := arg.CreatePaymentFn(order.ID, method.Code)
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
