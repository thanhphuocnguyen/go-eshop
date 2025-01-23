package payment

import (
	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
	"github.com/stripe/stripe-go/v81/refund"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

type stripePayment struct {
}

// InitiatePayment implements PaymentStrategy.
func (s *stripePayment) InitiatePayment(amount float64, email string) (string, error) {
	log.Info().Msg("InitiatePayment")
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amount * utils.MUL)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
	}
	if email != "" {
		params.ReceiptEmail = stripe.String(email)
	}

	rs, err := paymentintent.New(params)

	if err != nil {
		return "", err
	}

	return rs.ID, nil
}

// ProcessPayment implements PaymentStrategy.
func (s *stripePayment) ProcessPayment(transactionID string) (string, error) {
	rs, err := paymentintent.Confirm(transactionID, &stripe.PaymentIntentConfirmParams{})
	return rs.ID, err
}

// RefundPayment implements PaymentStrategy.
func (s *stripePayment) RefundPayment(transactionID string, amount int64, reason RefundReason) (string, error) {
	rs, err := refund.New(&stripe.RefundParams{
		PaymentIntent: stripe.String(transactionID),
		Amount:        stripe.Int64(amount),
		Reason:        stripe.String("requested_by_customer"),
	})
	return rs.ID, err
}

func (s *stripePayment) CancelPayment(transactionID string, reason CancelReason) (string, error) {
	rs, err := paymentintent.Cancel(transactionID, &stripe.PaymentIntentCancelParams{
		CancellationReason: stripe.String(string(reason)),
	})
	if err != nil {
		return "", err
	}
	return rs.ID, nil
}

func (s *stripePayment) GetPaymentObject(transactionID string) (interface{}, error) {
	rs, err := paymentintent.Get(transactionID, nil)
	return rs, err
}

func NewStripePayment(secretKey string) (PaymentStrategy, error) {
	if secretKey == "" {
		return nil, ErrEmptySecretKey
	}
	stripe.Key = secretKey
	return &stripePayment{}, nil
}
