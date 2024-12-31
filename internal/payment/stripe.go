package payment

import (
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/client"
	"github.com/thanhphuocnguyen/go-eshop/config"
)

type StripeClient struct {
	clientAPI *client.API
}

func NewStripeClient(cfg config.Config) PaymentService {
	config := &stripe.BackendConfig{
		MaxNetworkRetries: stripe.Int64(1), // Zero retries
	}
	sc := &client.API{}
	sc.Init(cfg.StripeSecretKey, &stripe.Backends{
		API:     stripe.GetBackendWithConfig(stripe.APIBackend, config),
		Uploads: stripe.GetBackendWithConfig(stripe.UploadsBackend, config),
	})
	stripeClient := &StripeClient{
		clientAPI: sc,
	}
	return stripeClient
}

func (client *StripeClient) MakePayment(param PaymentParam) (PaymentResult, error) {
	params := &stripe.PaymentIntentParams{
		Amount:        stripe.Int64(int64(param.Amount)),
		Currency:      stripe.String(string(stripe.CurrencyUSD)),
		PaymentMethod: stripe.String(param.Source),
		ReceiptEmail:  stripe.String(param.Email),
	}
	pi, err := client.clientAPI.PaymentIntents.New(params)
	if err != nil {
		return PaymentResult{}, err
	}
	return PaymentResult{
		TransactionID: pi.ID,
		Amount:        float64(pi.Amount),
		Currency:      string(pi.Currency),
		Source:        pi.PaymentMethod.ID,
		Email:         pi.ReceiptEmail,
	}, nil
}

func (client *StripeClient) CapturePaymentIntent(paymentIntentID string) (*stripe.PaymentIntent, error) {
	pi, err := client.clientAPI.PaymentIntents.Capture(paymentIntentID, nil)
	if err != nil {
		return nil, err
	}
	return pi, nil
}

func (client *StripeClient) RetrievePaymentIntent(paymentIntentID string) (*stripe.PaymentIntent, error) {
	pi, err := client.clientAPI.PaymentIntents.Get(paymentIntentID, nil)
	if err != nil {
		return nil, err
	}
	return pi, nil
}

func (client *StripeClient) CancelPayment(transactionID string) error {
	_, err := client.clientAPI.PaymentIntents.Cancel(transactionID, nil)
	if err != nil {
		return err
	}
	return nil
}
