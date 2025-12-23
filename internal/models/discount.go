package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AddDiscount struct {
	IsStackable       bool      `json:"isStackable" validate:"omitempty"`
	IsActive          bool      `json:"isActive" validate:"required"`
	DiscountValue     float64   `json:"discountValue" validate:"required,gt=0"`
	Code              string    `json:"code" validate:"required,min=5,max=32,alphanum"`
	Name              string    `json:"name" validate:"required,min=3,max=100"`
	DiscountType      string    `json:"discountType" validate:"required"`
	ValidFrom         time.Time `json:"validFrom" validate:"required" time_format:"2006-01-02T15:04:05Z07:00"`
	ValidUntil        time.Time `json:"validUntil" validate:"omitempty" time_format:"2006-01-02T15:04:05Z07:00"`
	Priority          *int32    `json:"priority" validate:"omitempty,gte=0"`
	Description       *string   `json:"description" validate:"omitempty,max=1000"`
	MinOrderValue     *float64  `json:"minOrderValue" validate:"omitempty,gt=0"`
	MaxDiscountAmount *float64  `json:"maxDiscountAmount" validate:"omitempty,gt=0"`
	UsageLimit        *int32    `json:"usageLimit" validate:"omitempty,gte=0"`
	UsagePerUser      *int32    `json:"usagePerUser" validate:"omitempty,gte=0"`
}

type AddDiscountRule struct {
	RuleType  string                 `json:"ruleType" validate:"required,oneof=product category customer_segment brand collection first_time_buyer purchase_quantity"`
	RuleValue map[string]interface{} `json:"ruleValue" validate:"required,min=1"`
}

type UpdateDiscountRule struct {
	RuleType  *string                `json:"ruleType" validate:"omitempty,oneof=condition action"`
	RuleValue map[string]interface{} `json:"ruleValue" validate:"omitempty,min=1"`
}

type UriRuleIDParam struct {
	ID     string `uri:"id" validate:"required,uuid"`
	RuleID string `uri:"ruleId" validate:"required,uuid"`
}

type UpdateDiscount struct {
	IsStackable       *bool      `json:"isStackable" validate:"omitempty"`
	IsActive          *bool      `json:"isActive" validate:"omitempty"`
	DiscountValue     *float64   `json:"discountValue" validate:"omitempty,gt=0"`
	Code              *string    `json:"code" validate:"omitempty,min=5,max=32,alphanum"`
	Name              *string    `json:"name" validate:"omitempty,min=3,max=100"`
	DiscountType      *string    `json:"discountType" validate:"omitempty"`
	ValidFrom         *time.Time `json:"validFrom" validate:"omitempty" time_format:"2006-01-02T15:04:05Z07:00"`
	ValidUntil        *time.Time `json:"validUntil" validate:"omitempty" time_format:"2006-01-02T15:04:05Z07:00"`
	Priority          *int32     `json:"priority" validate:"omitempty,gte=0"`
	Description       *string    `json:"description" validate:"omitempty,max=1000"`
	MinOrderValue     *float64   `json:"minOrderValue" validate:"omitempty,gt=0"`
	MaxDiscountAmount *float64   `json:"maxDiscountAmount" validate:"omitempty,gt=0"`
	UsageLimit        *int32     `json:"usageLimit" validate:"omitempty,gte=0"`
	UsagePerUser      *int32     `json:"usagePerUser" validate:"omitempty,gte=0"`
}

type DiscountListQuery struct {
	DiscountType  *string    `form:"discountType" validate:"omitempty,oneof=percentage fixed_amount"`
	IsActive      *bool      `from:"isActive" validate:"omitempty"`
	DiscountValue *float64   `form:"discountValue" validate:"omitempty,gt=0"`
	FromDate      *time.Time `form:"fromDate" validate:"omitempty" time_format:"2006-01-02T15:04:05Z07:00"`
	ToDate        *time.Time `form:"toDate" validate:"omitempty" time_format:"2006-01-02T15:04:05Z07:00"`
}

type DiscountRule struct {
	ID         uuid.UUID       `json:"id"`
	DiscountID uuid.UUID       `json:"discount_id"`
	RuleType   string          `json:"ruleType"`
	RuleValue  json.RawMessage `json:"ruleValue"`
}

// Example rule value structures
type ProductRule struct {
	ProductIDs []uuid.UUID `json:"product_ids"`
}

type CategoryRule struct {
	CategoryIDs      []string `json:"categoryIds"`
	ExcludeSaleItems bool     `json:"excludeSaleItems"`
}

type CollectionRule struct {
	CollectionIDs    []uuid.UUID `json:"collectionIds"`
	ExcludeSaleItems bool        `json:"excludeSaleItems"`
}

type PurchaseQuantityRule struct {
	MinQuantity int `json:"minQuantity"`
	MaxQuantity int `json:"maxQuantity"`
}

type FirstTimeBuyerRule struct {
	IsFirstTimeBuyer bool `json:"isFirstTimeBuyer"`
}

type BrandRule struct {
	BrandIDs []uuid.UUID `json:"brandIds"`
}

type CustomerSegmentRule struct {
	MinTotalSpent     *float64 `json:"minTotalSpent"`
	IsNewCustomer     bool     `json:"isNewCustomer"`
	MaxPreviousOrders *int     `json:"maxPreviousOrders"`
	CustomerType      *string  `json:"customerType"`
}

type CheckDiscountApplicabilityRequest struct {
	DiscountCodes []string `json:"discountCodes" validate:"required"`
	CartID        string   `json:"cartId" validate:"required,uuid"`
}

type CheckDiscountApplicabilityResponse struct {
	IsApplicable bool    `json:"isApplicable"`
	Message      string  `json:"message,omitempty"`
	DiscountAmt  float64 `json:"discountAmount,omitempty"`
}
