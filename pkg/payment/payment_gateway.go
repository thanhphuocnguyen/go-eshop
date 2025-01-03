package payment

import "errors"

type PaymentStrategy interface {
	InitiatePayment(amount float64, email string) (string, error)
	ProcessPayment(transactionID string) (string, error)
	GetPaymentObject(transactionID string) (interface{}, error)
	RefundPayment(transactionID string) (string, error)
	CancelPayment(transactionID string, reason string) (string, error)
}

var (
	ErrPaymentStrategyNotSet = errors.New("payment strategy not set")
	ErrEmptySecretKey        = errors.New("empty secret key")
)

type PaymentContext struct {
	strategy PaymentStrategy
}

func (p *PaymentContext) SetStrategy(strategy PaymentStrategy) {
	p.strategy = strategy
}

func (p *PaymentContext) InitiatePayment(amount float64, email string) (string, error) {
	if p.strategy == nil {
		return "", ErrPaymentStrategyNotSet
	}
	return p.strategy.InitiatePayment(amount, email)
}

func (p *PaymentContext) ProcessPayment(transactionID string) (string, error) {
	if p.strategy == nil {
		return "", ErrPaymentStrategyNotSet
	}
	return p.strategy.ProcessPayment(transactionID)
}

func (p *PaymentContext) RefundPayment(transactionID string) (string, error) {
	if p.strategy == nil {
		return "", ErrPaymentStrategyNotSet
	}
	return p.strategy.RefundPayment(transactionID)
}

func (p *PaymentContext) CancelPayment(transactionID string, reason string) (string, error) {
	if p.strategy == nil {
		return "", ErrPaymentStrategyNotSet
	}
	return p.strategy.CancelPayment(transactionID, reason)
}

func (p *PaymentContext) GetPaymentObject(transactionID string) (interface{}, error) {
	if p.strategy == nil {
		return nil, ErrPaymentStrategyNotSet
	}
	return p.strategy.GetPaymentObject(transactionID)
}
