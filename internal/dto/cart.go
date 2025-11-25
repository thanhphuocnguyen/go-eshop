package dto

import (
	"time"

	"github.com/google/uuid"
)

type CartDetail struct {
	ID             uuid.UUID        `json:"id"`
	TotalPrice     float64          `json:"totalPrice"`
	DiscountAmount float64          `json:"discountAmount"`
	UpdatedAt      *time.Time       `json:"updatedAt"`
	CreatedAt      time.Time        `json:"createdAt"`
	CartItems      []CartItemDetail `json:"cartItems"`
}
