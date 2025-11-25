package models

import "time"

type DiscountLinkObject struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Price *float64 `json:"price,omitempty"`
}

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
	RuleType      string `json:"ruleType" binding:"required,oneof=condition action"`
	ConditionType string `json:"conditionType" binding:"required,oneof=product_collection product_category product_brand product_attribute cart_total cart_item_quantity"`
}

type UpdateDiscountRuleModel struct {
	RuleType      *string `json:"ruleType" binding:"omitempty,oneof=condition action"`
	ConditionType *string `json:"conditionType" binding:"omitempty,oneof=product_collection product_category product_brand product_attribute cart_total cart_item_quantity"`
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
