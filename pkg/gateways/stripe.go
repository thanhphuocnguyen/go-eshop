package gateways

import (
	"context"
	"fmt"
	"time"

	"github.com/stripe/stripe-go/v84"
	"github.com/thanhphuocnguyen/go-eshop/pkg/payment"
)

type StripeGateway struct {
	apiKey        string
	webhookSecret string
	environment   string
	client        *stripe.Client
}

func NewStripeGateway(config payment.GatewayConfig) (payment.PaymentGateway, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("stripe API key is required")
	}
	client := stripe.NewClient(config.APIKey)

	return &StripeGateway{
		apiKey:        config.APIKey,
		webhookSecret: config.WebhookSecret,
		environment:   config.Environment,
		client:        client,
	}, nil
}

func (s *StripeGateway) Name() string {
	return "stripe"
}

func (s *StripeGateway) CreatePaymentIntent(ctx context.Context, req payment.PaymentRequest) (*payment.PaymentIntent, error) {
	// Implement Stripe PaymentIntent creation
	// This would use the actual Stripe SDK
	intent, err := s.client.V1PaymentIntents.Create(ctx, &stripe.PaymentIntentCreateParams{
		Amount:       &req.Amount,
		Description:  &req.Description,
		ReceiptEmail: &req.Email,
		Metadata:     req.Metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Stripe payment intent: %w", err)
	}

	return &payment.PaymentIntent{
		ID:           intent.ID,
		Amount:       intent.Amount,
		Currency:     payment.Currency(intent.Currency),
		Status:       payment.PaymentStatus(intent.Status),
		Email:        intent.Customer.Email,
		Description:  intent.Description,
		ClientSecret: intent.ClientSecret,
		CreatedAt:    time.Unix(intent.Created, 0),
		Metadata:     intent.Metadata,
	}, nil
}

func (s *StripeGateway) ConfirmPayment(ctx context.Context, intentID string) (*payment.PaymentResult, error) {
	// Implement Stripe payment confirmation

	rs, err := s.client.V1PaymentIntents.Confirm(ctx, intentID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm Stripe payment intent: %w", err)
	}

	result := &payment.PaymentResult{
		Success:       rs.Status == stripe.PaymentIntentStatusSucceeded,
		TransactionID: rs.ID,
		Status:        payment.PaymentStatus(rs.Status),
		Message:       "Payment confirmed",
		Metadata:      rs.Metadata,
		ProcessedAt:   time.Now(),
	}

	return result, nil
}

func (s *StripeGateway) GetPayment(ctx context.Context, transactionID string) (*payment.PaymentIntent, error) {
	intent, err := s.client.V1PaymentIntents.Retrieve(ctx, transactionID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve Stripe payment intent: %w", err)
	}
	return &payment.PaymentIntent{
		ID:           intent.ID,
		Amount:       intent.Amount,
		Currency:     payment.Currency(intent.Currency),
		Status:       payment.PaymentStatus(intent.Status),
		Email:        intent.Customer.Email,
		Description:  intent.Description,
		ClientSecret: intent.ClientSecret,
		CreatedAt:    time.Unix(intent.Created, 0),
		Metadata:     intent.Metadata,
	}, nil
}

func (s *StripeGateway) RefundPayment(ctx context.Context, req payment.RefundRequest) (*payment.RefundResult, error) {
	// Implement refund logic
	refundRs, err := s.client.V1Refunds.Create(ctx, &stripe.RefundCreateParams{
		PaymentIntent: &req.TransactionID,
		Amount:        &req.Amount,
		Reason:        &req.Reason,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Stripe refund: %w", err)
	}
	return &payment.RefundResult{
		RefundID:      refundRs.ID,
		TransactionID: req.TransactionID,
		Amount:        req.Amount,
		Status:        payment.StatusCompleted,
		Reason:        req.Reason,
		CreatedAt:     time.Now(),
	}, nil
}

func (s *StripeGateway) CancelPayment(ctx context.Context, intentID string) error {
	// Implement payment cancellation
	_, err := s.client.V1PaymentIntents.Cancel(ctx, intentID, nil)
	if err != nil {
		return fmt.Errorf("failed to cancel Stripe payment intent: %w", err)
	}
	return nil
}

func (s *StripeGateway) VerifyWebhook(payload []byte, signature string) error {
	// Implement webhook verification
	return nil
}

func (s *StripeGateway) Health(ctx context.Context) error {
	// Implement health check
	_, err := s.client.V1Balance.Retrieve(ctx, nil)
	if err != nil {
		return fmt.Errorf("stripe health check failed: %w", err)
	}
	return nil
}
