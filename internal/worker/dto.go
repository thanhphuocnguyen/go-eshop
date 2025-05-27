package worker

import "github.com/google/uuid"

type PayloadVerifyEmail struct {
	UserID uuid.UUID `json:"user_id"`
}

type VerifyEmailData struct {
	UserID     uuid.UUID
	Email      string
	FullName   string
	VerifyLink string
}

type PayloadSendOrderCreatedEmailTask struct {
	PaymentID uuid.UUID `json:"paymentId"`
}

type OrderCreatedItems struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Qty   int     `json:"qty"`
}
type OrderCreatedEmailData struct {
	OrderID  uuid.UUID           `json:"orderId"`
	Total    float64             `json:"total"`
	FullName string              `json:"fullName"`
	Items    []OrderCreatedItems `json:"items"`
}
