package dto

import "github.com/thanhphuocnguyen/go-eshop/internal/db/repository"

type PaymentIntentSecret struct {
	PaymentID    string  `json:"paymentId"`
	ClientSecret *string `json:"clientSecret"`
}

type PaymentDetail struct {
	ID      string                   `json:"id"`
	Gateway *string                  `json:"gateway,omitempty"`
	Status  repository.PaymentStatus `json:"status,omitempty"`
	// Details contains gateway-specific payment information.
	// The type may vary depending on the payment provider, and is typically a struct or map[string]interface{}.
	Details any `json:"details"`
}

type PaymentInfo struct {
	ID           string  `json:"id"`
	RefundID     *string `json:"refundId"`
	Amount       float64 `json:"amount"`
	IntendID     *string `json:"intentId"`
	ClientSecret *string `json:"clientSecret"`
	GateWay      *string `json:"gateway"`
	Method       string  `json:"method"`
	Status       string  `json:"status"`
}
