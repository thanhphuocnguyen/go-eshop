package payment

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrGatewayNotFound   = errors.New("payment gateway not found")
	ErrGatewayNotEnabled = errors.New("payment gateway not enabled")
	ErrInvalidAmount     = errors.New("invalid payment amount")
	ErrInvalidCurrency   = errors.New("invalid currency")
	ErrInvalidEmail      = errors.New("invalid email address")
)

// GatewayFactory creates payment gateway instances
type GatewayFactory func(config GatewayConfig) (PaymentGateway, error)

// PaymentService orchestrates payment operations across multiple gateways
type PaymentService struct {
	gateways  map[string]PaymentGateway
	factories map[string]GatewayFactory
	mu        sync.RWMutex
}

// NewPaymentService creates a new payment service instance
func NewPaymentService() *PaymentService {
	return &PaymentService{
		gateways:  make(map[string]PaymentGateway),
		factories: make(map[string]GatewayFactory),
	}
}

// RegisterGateway registers a gateway factory
func (ps *PaymentService) RegisterGateway(name string, factory GatewayFactory) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.factories[name] = factory
}

// AddGateway adds a configured gateway
func (ps *PaymentService) AddGateway(config GatewayConfig) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	factory, exists := ps.factories[config.Name]
	if !exists {
		return fmt.Errorf("gateway factory %s not registered", config.Name)
	}

	gateway, err := factory(config)
	if err != nil {
		return fmt.Errorf("failed to create gateway %s: %w", config.Name, err)
	}

	ps.gateways[config.Name] = gateway
	return nil
}

// CreatePaymentIntent creates a payment intent using the primary gateway
func (ps *PaymentService) CreatePaymentIntent(ctx context.Context, gatewayName string, req PaymentRequest) (*PaymentIntent, error) {
	if err := ps.validatePaymentRequest(req); err != nil {
		return nil, err
	}
	gateway, err := ps.getGateway(gatewayName)
	if err != nil {
		return nil, err
	}

	return ps.createPaymentIntentWithGateway(ctx, req, gateway)
}

// CreatePaymentIntentWithGateway creates a payment intent using a specific gateway
func (ps *PaymentService) CreatePaymentIntentWithGateway(ctx context.Context, req PaymentRequest, gatewayName string) (*PaymentIntent, error) {
	if err := ps.validatePaymentRequest(req); err != nil {
		return nil, err
	}

	gateway, err := ps.getGateway(gatewayName)
	if err != nil {
		return nil, err
	}

	return ps.createPaymentIntentWithGateway(ctx, req, gateway)
}

// ConfirmPayment confirms a payment with fallback support
func (ps *PaymentService) ConfirmPayment(ctx context.Context, intentID string, gatewayName string) (*PaymentResult, error) {
	gateway, err := ps.getGateway(gatewayName)
	if err != nil {
		return nil, err
	}

	result, err := gateway.ConfirmPayment(ctx, intentID)

	return result, nil
}

// RefundPayment processes a refund
func (ps *PaymentService) RefundPayment(ctx context.Context, req RefundRequest, gatewayName string) (*RefundResult, error) {
	gateway, err := ps.getGateway(gatewayName)
	if err != nil {
		return nil, err
	}

	return gateway.RefundPayment(ctx, req)
}

// GetPayment retrieves payment details
func (ps *PaymentService) GetPayment(ctx context.Context, transactionID string, gatewayName string) (*PaymentIntent, error) {
	gateway, err := ps.getGateway(gatewayName)
	if err != nil {
		return nil, err
	}

	return gateway.GetPayment(ctx, transactionID)
}

// Helper methods

func (ps *PaymentService) validatePaymentRequest(req PaymentRequest) error {
	if req.Amount <= 0 {
		return ErrInvalidAmount
	}

	if req.Currency == "" {
		return ErrInvalidCurrency
	}

	if req.Email == "" {
		return ErrInvalidEmail
	}

	return nil
}

func (ps *PaymentService) getGateway(name string) (PaymentGateway, error) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	gateway, exists := ps.gateways[name]
	if !exists {
		return nil, ErrGatewayNotFound
	}

	return gateway, nil
}

func (ps *PaymentService) createPaymentIntentWithGateway(ctx context.Context, req PaymentRequest, gateway PaymentGateway) (*PaymentIntent, error) {
	return gateway.CreatePaymentIntent(ctx, req)
}
