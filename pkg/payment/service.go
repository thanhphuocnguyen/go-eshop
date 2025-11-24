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
	primary   string
	fallbacks []string
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

// SetPrimaryGateway sets the primary payment gateway
func (ps *PaymentService) SetPrimaryGateway(name string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, exists := ps.gateways[name]; !exists {
		return ErrGatewayNotFound
	}

	ps.primary = name
	return nil
}

// SetFallbackGateways sets fallback gateways in order of preference
func (ps *PaymentService) SetFallbackGateways(names ...string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for _, name := range names {
		if _, exists := ps.gateways[name]; !exists {
			return fmt.Errorf("fallback gateway %s not found", name)
		}
	}

	ps.fallbacks = names
	return nil
}

// CreatePaymentIntent creates a payment intent using the primary gateway
func (ps *PaymentService) CreatePaymentIntent(ctx context.Context, req PaymentRequest) (*PaymentIntent, error) {
	if err := ps.validatePaymentRequest(req); err != nil {
		return nil, err
	}

	gateway, err := ps.getPrimaryGateway()
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
	if err != nil {
		// Try fallback gateways if primary fails
		return ps.confirmPaymentWithFallback(ctx, intentID)
	}

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

func (ps *PaymentService) getPrimaryGateway() (PaymentGateway, error) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if ps.primary == "" {
		return nil, errors.New("no primary gateway configured")
	}

	gateway, exists := ps.gateways[ps.primary]
	if !exists {
		return nil, ErrGatewayNotFound
	}

	return gateway, nil
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

func (ps *PaymentService) confirmPaymentWithFallback(ctx context.Context, intentID string) (*PaymentResult, error) {
	ps.mu.RLock()
	fallbacks := ps.fallbacks
	ps.mu.RUnlock()

	for _, gatewayName := range fallbacks {
		gateway, err := ps.getGateway(gatewayName)
		if err != nil {
			continue
		}

		result, err := gateway.ConfirmPayment(ctx, intentID)
		if err == nil {
			return result, nil
		}
	}

	return nil, errors.New("all payment gateways failed")
}
