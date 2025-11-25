package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
)

type OrderItemDetail struct {
	ID                 string                             `json:"id"`
	VariantID          string                             `json:"variantId"`
	Name               string                             `json:"name"`
	ImageUrl           *string                            `json:"imageUrl"`
	LineTotal          float64                            `json:"lineTotal"`
	Quantity           int16                              `json:"quantity"`
	AttributesSnapshot []repository.AttributeDataSnapshot `json:"attributesSnapshot"`
	Rating             *RatingDetail                      `json:"rating,omitempty"`
}

type OrderDetail struct {
	ID            uuid.UUID                          `json:"id"`
	Total         float64                            `json:"total"`
	Status        repository.OrderStatus             `json:"status"`
	CustomerName  string                             `json:"customerName"`
	CustomerEmail string                             `json:"customerEmail"`
	PaymentInfo   PaymentInfo                        `json:"paymentInfo"`
	ShippingInfo  repository.ShippingAddressSnapshot `json:"shippingInfo"`
	Products      []OrderItemDetail                  `json:"products"`
	CreatedAt     time.Time                          `json:"createdAt"`
}

type OrderListItem struct {
	ID            uuid.UUID                `json:"id"`
	Total         float64                  `json:"total"`
	TotalItems    int32                    `json:"totalItems"`
	Status        repository.OrderStatus   `json:"status"`
	PaymentStatus repository.PaymentStatus `json:"paymentStatus"`
	CustomerName  string                   `json:"customerName"`
	CustomerEmail string                   `json:"customerEmail"`
	CreatedAt     time.Time                `json:"createdAt"`
	UpdatedAt     time.Time                `json:"updatedAt"`
}

type OrderItemAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
