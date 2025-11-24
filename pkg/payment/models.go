package payment

import (
	"time"
)

// PaymentStatus represents the current state of a payment
type PaymentStatus string

const (
	StatusPending   PaymentStatus = "pending"
	StatusCompleted PaymentStatus = "completed"
	StatusFailed    PaymentStatus = "failed"
	StatusRefunded  PaymentStatus = "refunded"
	StatusCancelled PaymentStatus = "cancelled"
)

// Currency represents supported currencies
type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
)

// PaymentRequest encapsulates all payment creation data
type PaymentRequest struct {
	Amount      int64             `json:"amount"` // Amount in smallest currency unit
	Currency    Currency          `json:"currency"`
	Email       string            `json:"email"`
	Description string            `json:"description"`
	Metadata    map[string]string `json:"metadata"`
}

// PaymentIntent represents a payment intention
type PaymentIntent struct {
	ID           string            `json:"id"`
	Amount       int64             `json:"amount"`
	Currency     Currency          `json:"currency"`
	Status       PaymentStatus     `json:"status"`
	ClientSecret string            `json:"client_secret,omitempty"`
	Email        string            `json:"email"`
	Description  string            `json:"description"`
	Metadata     map[string]string `json:"metadata"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// PaymentResult represents the outcome of a payment operation
type PaymentResult struct {
	Success       bool              `json:"success"`
	TransactionID string            `json:"transaction_id"`
	Status        PaymentStatus     `json:"status"`
	Message       string            `json:"message"`
	Metadata      map[string]string `json:"metadata"`
	ProcessedAt   time.Time         `json:"processed_at"`
}

// RefundRequest encapsulates refund data
type RefundRequest struct {
	TransactionID string `json:"transaction_id"`
	Amount        int64  `json:"amount"` // Amount to refund, nil for full refund
	Reason        string `json:"reason"`
}

// RefundResult represents the outcome of a refund
type RefundResult struct {
	RefundID      string        `json:"refund_id"`
	TransactionID string        `json:"transaction_id"`
	Amount        int64         `json:"amount"`
	Status        PaymentStatus `json:"status"`
	Reason        string        `json:"reason"`
	CreatedAt     time.Time     `json:"created_at"`
}
