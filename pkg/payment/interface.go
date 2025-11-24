package payment

import "context"

// PaymentGateway defines the contract for all payment providers
type PaymentGateway interface {
	// Name returns the gateway identifier
	Name() string

	// CreatePaymentIntent creates a payment intent
	CreatePaymentIntent(ctx context.Context, req PaymentRequest) (*PaymentIntent, error)

	// ConfirmPayment processes and confirms a payment
	ConfirmPayment(ctx context.Context, intentID string) (*PaymentResult, error)

	// GetPayment retrieves payment details
	GetPayment(ctx context.Context, transactionID string) (*PaymentIntent, error)

	// RefundPayment processes a refund
	RefundPayment(ctx context.Context, req RefundRequest) (*RefundResult, error)

	// CancelPayment cancels a pending payment
	CancelPayment(ctx context.Context, intentID string) error

	// VerifyWebhook verifies incoming webhook signatures
	VerifyWebhook(payload []byte, signature string) error

	// Health checks if the gateway is operational
	Health(ctx context.Context) error
}

// GatewayConfig holds configuration for payment gateways
type GatewayConfig struct {
	Name          string            `json:"name"`
	APIKey        string            `json:"api_key"`
	SecretKey     string            `json:"secret_key"`
	WebhookSecret string            `json:"webhook_secret"`
	Environment   string            `json:"environment"` // sandbox, production
	Extra         map[string]string `json:"extra"`
}
