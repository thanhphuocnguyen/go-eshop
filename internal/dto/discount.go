package dto

import "time"

type DiscountListItem struct {
	ID            string   `json:"id"`
	Code          string   `json:"code"`
	Description   *string  `json:"description"`
	DiscountType  string   `json:"discountType"`
	DiscountValue float64  `json:"discountValue"`
	ValidUntil    string   `json:"validUntil"`
	ValidFrom     string   `json:"validFrom"`
	TimeUsed      int32    `json:"usedCount"`
	IsActive      bool     `json:"isActive"`
	CreatedAt     string   `json:"createdAt"`
	UpdatedAt     string   `json:"updatedAt"`
	UsageLimit    *int32   `json:"usageLimit,omitempty"`
	MinPurchase   *float64 `json:"minPurchase,omitzero,omitempty"`
	MaxDiscount   *float64 `json:"maxDiscount,omitempty"`
}

type DiscountUsageHistory struct {
	ID             string    `json:"id"`
	OrderID        string    `json:"orderId"`
	CustomerName   string    `json:"customerName"`
	Amount         float64   `json:"amount"`
	DiscountAmount float64   `json:"discountAmount"`
	Date           time.Time `json:"date" format:"2006-01-02"`
}

type DiscountDetail struct {
	ID            string                 `json:"id"`
	Code          string                 `json:"code"`
	Description   *string                `json:"description"`
	DiscountType  string                 `json:"discountType"`
	DiscountValue float64                `json:"discountValue"`
	ValidUntil    string                 `json:"validUntil"`
	ValidFrom     string                 `json:"validFrom"`
	TimesUsed     int32                  `json:"timesUsed"`
	UsageLimit    *int32                 `json:"usageLimit"`
	IsActive      bool                   `json:"isActive"`
	CreatedAt     string                 `json:"createdAt"`
	UpdatedAt     string                 `json:"updatedAt"`
	MinPurchase   float64                `json:"minPurchase"`
	MaxDiscount   float64                `json:"maxDiscount"`
	UsageHistory  []DiscountUsageHistory `json:"usageHistory"`
}
