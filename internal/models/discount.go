package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AddDiscountModel struct {
	IsStackable       bool      `json:"isStackable" binding:"omitempty"`
	IsActive          bool      `json:"isActive" binding:"required"`
	DiscountValue     float64   `json:"discountValue" binding:"required,gt=0"`
	Code              string    `json:"code" binding:"required,min=5,max=32,alphanum"`
	Name              string    `json:"name" binding:"required,min=3,max=100"`
	DiscountType      string    `json:"discountType" binding:"required"`
	ValidFrom         time.Time `json:"validFrom" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
	ValidUntil        time.Time `json:"validUntil" binding:"omitempty" time_format:"2006-01-02T15:04:05Z07:00"`
	Priority          *int32    `json:"priority" binding:"omitempty,gte=0"`
	Description       *string   `json:"description" binding:"omitempty,max=1000"`
	MinOrderValue     *float64  `json:"minOrderValue" binding:"omitempty,gt=0"`
	MaxDiscountAmount *float64  `json:"maxDiscountAmount" binding:"omitempty,gt=0"`
	UsageLimit        *int32    `json:"usageLimit" binding:"omitempty,gte=0"`
}

type AddDiscountRuleModel struct {
	RuleType  string                 `json:"ruleType" binding:"required,oneof=product category customer_segment"`
	RuleValue map[string]interface{} `json:"ruleValue" binding:"required,min=1"`
}

type UpdateDiscountRuleModel struct {
	RuleType  *string                `json:"ruleType" binding:"omitempty,oneof=condition action"`
	RuleValue map[string]interface{} `json:"ruleValue" binding:"omitempty,min=1"`
}

type UriRuleIDParam struct {
	ID     string `uri:"id" binding:"required,uuid"`
	RuleID string `uri:"ruleId" binding:"required,uuid"`
}

type UpdateDiscountModel struct {
	IsStackable       *bool      `json:"isStackable" binding:"omitempty"`
	IsActive          *bool      `json:"isActive" binding:"omitempty"`
	DiscountValue     *float64   `json:"discountValue" binding:"omitempty,gt=0"`
	Code              *string    `json:"code" binding:"omitempty,min=5,max=32,alphanum"`
	Name              *string    `json:"name" binding:"omitempty,min=3,max=100"`
	DiscountType      *string    `json:"discountType" binding:"omitempty"`
	ValidFrom         *time.Time `json:"validFrom" binding:"omitempty" time_format:"2006-01-02T15:04:05Z07:00"`
	ValidUntil        *time.Time `json:"validUntil" binding:"omitempty" time_format:"2006-01-02T15:04:05Z07:00"`
	Priority          *int32     `json:"priority" binding:"omitempty,gte=0"`
	Description       *string    `json:"description" binding:"omitempty,max=1000"`
	MinOrderValue     *float64   `json:"minOrderValue" binding:"omitempty,gt=0"`
	MaxDiscountAmount *float64   `json:"maxDiscountAmount" binding:"omitempty,gt=0"`
	UsageLimit        *int32     `json:"usageLimit" binding:"omitempty,gte=0"`
}

type DiscountListQuery struct {
	PaginationQuery
	DiscountType  *string    `form:"discountType" binding:"omitempty,oneof=percentage fixed_amount"`
	IsActive      *bool      `from:"isActive" binding:"omitempty"`
	DiscountValue *float64   `form:"discountValue" binding:"omitempty,gt=0"`
	FromDate      *time.Time `form:"fromDate" binding:"omitempty" time_format:"2006-01-02T15:04:05Z07:00"`
	ToDate        *time.Time `form:"toDate" binding:"omitempty" time_format:"2006-01-02T15:04:05Z07:00"`
}

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
