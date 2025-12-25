package models

import (
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
)

type PaymentModel struct {
	OrderID         string `json:"orderId" validate:"required,uuid"`
	PaymentMethodID string `json:"paymentMethodId" validate:"required,uuid"`
}

type UpdatePaymentStatusModel struct {
	Status repository.PaymentStatus `json:"status" validate:"required"`
}
