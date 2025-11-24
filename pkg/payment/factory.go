package payment

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
)

type PaymentIntentRequest struct {
	OrderID      uuid.UUID         `json:"order_id"`
	Amount       float64           `json:"amount"`
	Currency     string            `json:"currency"`
	Email        string            `json:"email"`
	Description  string            `json:"description"`
	CustomerInfo CustomerInfo      `json:"customer_info"`
	Metadata     map[string]string `json:"metadata"`
}

type CustomerInfo struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type PaymentIntentResult struct {
	PaymentIntentID string                 `json:"payment_intent_id"`
	ClientSecret    *string                `json:"client_secret,omitempty"`
	RedirectURL     *string                `json:"redirect_url,omitempty"`
	QRCode          *string                `json:"qr_code,omitempty"`
	ExpiresAt       *int64                 `json:"expires_at,omitempty"`
	Instructions    *string                `json:"instructions,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

type PaymentIntentFactory interface {
	CreatePaymentIntent(ctx context.Context, req PaymentIntentRequest, config repository.PaymentMethodConfig) (*PaymentIntentResult, error)
	GetPaymentMethodType() string
}

type PaymentIntentFactoryRegistry struct {
	factories map[string]PaymentIntentFactory
}

func NewPaymentIntentFactoryRegistry() *PaymentIntentFactoryRegistry {
	return &PaymentIntentFactoryRegistry{
		factories: make(map[string]PaymentIntentFactory),
	}
}

func (r *PaymentIntentFactoryRegistry) Register(factory PaymentIntentFactory) {
	r.factories[factory.GetPaymentMethodType()] = factory
}

func (r *PaymentIntentFactoryRegistry) CreatePaymentIntent(
	ctx context.Context,
	methodType string,
	req PaymentIntentRequest,
	config repository.PaymentMethodConfig,
) (*PaymentIntentResult, error) {
	factory, exists := r.factories[methodType]
	if !exists {
		return nil, fmt.Errorf("payment method type %s not supported", methodType)
	}

	return factory.CreatePaymentIntent(ctx, req, config)
}

func (r *PaymentIntentFactoryRegistry) GetSupportedMethods() []string {
	methods := make([]string, 0, len(r.factories))
	for methodType := range r.factories {
		methods = append(methods, methodType)
	}
	return methods
}
