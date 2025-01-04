package payment

import (
	"errors"

	"github.com/thanhphuocnguyen/go-eshop/config"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
)

type PaymentStrategy interface {
	InitiatePayment(amount float64, email string) (string, error)
	ProcessPayment(transactionID string) (string, error)
	GetPaymentObject(transactionID string) (interface{}, error)
	RefundPayment(transactionID, reason string) (string, error)
	CancelPayment(transactionID, reason string) (string, error)
}

var (
	ErrPaymentStrategyNotSet = errors.New("payment strategy not set")
	ErrEmptySecretKey        = errors.New("empty secret key")
)

type PaymentContext struct {
	strategy PaymentStrategy
}

func (p *PaymentContext) SetStrategy(strategy PaymentStrategy) error {
	if p.strategy == nil {
		return ErrPaymentStrategyNotSet
	}
	p.strategy = strategy
	return nil
}

func (p *PaymentContext) InitiatePayment(amount float64, email string) (string, error) {
	return p.strategy.InitiatePayment(amount, email)
}

func (p *PaymentContext) ProcessPayment(transactionID string) (string, error) {

	return p.strategy.ProcessPayment(transactionID)
}

func (p *PaymentContext) RefundPayment(transactionID string, reason string) (string, error) {
	return p.strategy.RefundPayment(transactionID, reason)
}

func (p *PaymentContext) CancelPayment(transactionID string, reason string) (string, error) {
	return p.strategy.CancelPayment(transactionID, reason)
}

func (p *PaymentContext) GetPaymentObject(transactionID string) (interface{}, error) {
	return p.strategy.GetPaymentObject(transactionID)
}

func (p *PaymentContext) GetPaymentGatewayInstanceFromName(name repository.PaymentGateway, cfg config.Config) PaymentStrategy {
	switch name {
	case repository.PaymentGatewayStripe:
		instance, err := NewStripePayment(cfg.StripeSecretKey)
		if err != nil {
			return nil
		}
		return instance
	default:
		return nil
	}
}
