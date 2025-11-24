package dto

import (
	"encoding/json"

	"github.com/google/uuid"
)

type DiscountRule struct {
	ID         uuid.UUID       `json:"id"`
	DiscountID uuid.UUID       `json:"discount_id"`
	RuleType   string          `json:"rule_type"`
	RuleValue  json.RawMessage `json:"rule_value"`
}

// Example rule value structures
type ProductRule struct {
	ProductIDs []uuid.UUID `json:"product_ids"`
}

type CategoryRule struct {
	CategoryIDs      []uuid.UUID `json:"category_ids"`
	ExcludeSaleItems bool        `json:"exclude_sale_items"`
}

type CartTotalRule struct {
	MinTotal float64 `json:"min_total"`
	MaxTotal float64 `json:"max_total"`
}

type PurchaseQuantityRule struct {
	MinQuantity int `json:"min_quantity"`
	MaxQuantity int `json:"max_quantity"`
}

type FirstTimeBuyerRule struct {
	IsFirstTimeBuyer bool `json:"is_first_time_buyer"`
}

type BrandRule struct {
	BrandIDs []uuid.UUID `json:"brand_ids"`
}

type CustomerSegmentRule struct {
	MinTotalSpent     float64 `json:"min_total_spent"`
	IsNewCustomer     bool    `json:"is_new_customer"`
	MaxPreviousOrders int     `json:"max_previous_orders"`
	CustomerType      string  `json:"customer_type"`
}
