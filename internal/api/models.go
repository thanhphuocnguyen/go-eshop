package api

import (
	"time"

	"github.com/google/uuid"
)

type ProductRatingModel struct {
	ID               uuid.UUID          `json:"id"`
	Name             string             `json:"name"`
	ProductName      string             `json:"product_name,omitempty"`
	UserID           uuid.UUID          `json:"user_id"`
	Rating           float64            `json:"rating"`
	ReviewTitle      string             `json:"review_title"`
	IsVisible        bool               `json:"is_visible"`
	IsApproved       bool               `json:"is_approved"`
	ReviewContent    string             `json:"review_content"`
	VerifiedPurchase bool               `json:"verified_purchase"`
	HelpfulVotes     int32              `json:"helpful_votes"`
	UnhelpfulVotes   int32              `json:"unhelpful_votes"`
	Images           []RatingImageModel `json:"images"`
}

type ProductListModel struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	VariantCount int64    `json:"variant_count,omitzero"`
	MinPrice     float64  `json:"min_price,omitzero"`
	MaxPrice     float64  `json:"max_price,omitzero"`
	Slug         string   `json:"slug,omitempty"`
	Sku          string   `json:"sku"`
	ImgUrl       *string  `json:"image_url,omitempty"`
	AvgRating    *float64 `json:"avg_rating,omitempty"`
	ReviewCount  *int32   `json:"review_count,omitempty"`
	ImgID        *string  `json:"image_id,omitempty"`
	CreatedAt    string   `json:"created_at,omitempty"`
	UpdatedAt    string   `json:"updated_at,omitempty"`
}

type FiltersModel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProductVariantImageModel struct {
	ID        string `json:"id"`
	VariantID string `json:"variant_id,omitempty"`
	ImageUrl  string `json:"image_url"`
	ImageID   string `json:"image_id"`
}

type CategoryLinkedProduct struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	VariantCount int32   `json:"variant_count"`
	ImageUrl     *string `json:"image_url,omitempty"`
}

type ProductVariantModel struct {
	ID         string                  `json:"id"`
	Price      float64                 `json:"price"`
	StockQty   int32                   `json:"stock_qty"`
	IsActive   bool                    `json:"is_active"`
	Sku        *string                 `json:"sku,omitempty"`
	Attributes []ProductAttributeModel `json:"attributes"`
}

type ImageAssignmentModel struct {
	ID           string `json:"id"`
	EntityID     string `json:"entity_id"`
	EntityType   string `json:"entity_type"`
	DisplayOrder int16  `json:"display_order"`
	Role         string `json:"role"`
}

type ProductImageModel struct {
	ID                 string                 `json:"id"`
	Url                string                 `json:"url"`
	ExternalID         string                 `json:"external_id"`
	Role               string                 `json:"role"`
	VariantAssignments []ImageAssignmentModel `json:"assignments"`
}

type Address struct {
	Phone    string  `json:"phone"`
	Street   string  `json:"street"`
	Ward     *string `json:"ward,omitempty"`
	District string  `json:"district"`
	City     string  `json:"city"`
}

type OrderItemAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type PaymentInfoModel struct {
	ID           string  `json:"id"`
	RefundID     *string `json:"refund_id"`
	Amount       float64 `json:"amount"`
	IntendID     *string `json:"intent_id"`
	ClientSecret *string `json:"client_secret"`
	GateWay      *string `json:"gateway"`
	Method       string  `json:"method"`
	Status       string  `json:"status"`
}

type RatingModel struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Rating    float64   `json:"rating"`
	CreatedAt time.Time `json:"created_at"`
}

type RatingImageModel struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type ProductAttributeModel struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	ValueObject AttributeValue `json:"value_object"`
}
