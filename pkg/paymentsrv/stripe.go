package paymentsrv

import (
	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
	"github.com/stripe/stripe-go/v81/refund"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

type stripePayment struct {
}

// CreatePaymentIntent implements PaymentStrategy.
func (s *stripePayment) CreatePaymentIntent(amount float64, email string) (result interface{}, err error) {
	log.Info().Msg("CreatePaymentIntent For")
	params := &stripe.PaymentIntentParams{
		Amount:        stripe.Int64(int64(amount * utils.MUL)),
		Currency:      stripe.String(string(stripe.CurrencyUSD)),
		ReceiptEmail:  stripe.String(email),
		PaymentMethod: stripe.String("pm_card_visa"),
	}
	if email != "" {
		params.ReceiptEmail = stripe.String(email)
	}

	rs, err := paymentintent.New(params)

	if err != nil {
		return nil, err
	}

	return rs, nil
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

func (s *stripePayment) CancelPayment(transactionID string, reason string) (string, error) {
	rs, err := paymentintent.Cancel(transactionID, &stripe.PaymentIntentCancelParams{
		CancellationReason: stripe.String(reason),
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
