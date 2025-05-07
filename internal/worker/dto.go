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
	PaymentID string    `json:"payment_id"`
	OrderID   uuid.UUID `json:"order_id"`
}

type OrderCreatedItems struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Qty   int     `json:"qty"`
}
type OrderCreatedEmailData struct {
	OrderID  uuid.UUID           `json:"order_id"`
	Total    float64             `json:"total"`
	FullName string              `json:"full_name"`
	Items    []OrderCreatedItems `json:"items"`
}
