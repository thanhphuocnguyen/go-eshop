package payment

import (
	"errors"
)

type RefundReason string

const (
	RefundReasonRequestedByCustomer  RefundReason = "Requested by customer"
	RefundReasonByDefectiveOrDamaged RefundReason = "Defective or damaged"
	RefundReasonByFraudulent         RefundReason = "Fraudulent"
)

type PaymentStrategy interface {
	InitiatePayment(amount float64, email string) (string, error)
	ProcessPayment(transactionID string) (string, error)
	GetPaymentObject(transactionID string) (interface{}, error)
	RefundPayment(transactionID string, amount int64, reason RefundReason) (string, error)
	CancelPayment(transactionID, reason string) (string, error)
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

	return p.strategy.ProcessPayment(transactionID)
}

func (p *PaymentContext) RefundPayment(transactionID string, amount int64, reason RefundReason) (string, error) {
	return p.strategy.RefundPayment(transactionID, amount, reason)
}

func (p *PaymentContext) CancelPayment(transactionID string, reason string) (string, error) {
	return p.strategy.CancelPayment(transactionID, reason)
}

func (p *PaymentContext) GetPaymentObject(transactionID string) (interface{}, error) {
	return p.strategy.GetPaymentObject(transactionID)
}
