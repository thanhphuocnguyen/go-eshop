package dto

type ProductAttribute struct {
	ID     int32                  `json:"attributeId"`
	Name   string                 `json:"attributeName"`
	Values []AttributeValueDetail `json:"attributeValues,omitempty"`
}
type ProductListItem struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Description        string   `json:"description"`
	BasePrice          float64  `json:"basePrice,omitzero"`
	Slug               string   `json:"slug,omitempty"`
	Sku                string   `json:"sku"`
	ImageUrl           *string  `json:"imageUrl,omitempty"`
	AvgRating          *float64 `json:"avgRating,omitempty"`
	DiscountPercentage *int16   `json:"discountPercentage,omitzero"`
	ReviewCount        *int32   `json:"reviewCount,omitempty"`
	ImgID              *string  `json:"imageId,omitempty"`
	CreatedAt          string   `json:"createdAt,omitempty"`
	UpdatedAt          string   `json:"updatedAt,omitempty"`
}

type ProductSummary struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Price        float64  `json:"price,omitzero"`
	Slug         string   `json:"slug,omitempty"`
	ImageUrl     *string  `json:"imageUrl,omitempty"`
	AvgRating    *float64 `json:"avgRating,omitempty"`
	VariantCount int16    `json:"variantCount,omitzero"`
	ReviewCount  *int32   `json:"reviewCount,omitempty"`
	ImageID      *string  `json:"imageId,omitempty"`
	CreatedAt    string   `json:"createdAt,omitempty"`
	UpdatedAt    string   `json:"updatedAt,omitempty"`
}
type VariantDetail struct {
	ID         string                 `json:"id"`
	Price      float64                `json:"price"`
	Stock      int32                  `json:"stock"`
	IsActive   bool                   `json:"isActive"`
	Sku        string                 `json:"sku,omitempty"`
	Weight     *float64               `json:"weight,omitempty"`
	ImageUrl   *string                `json:"imageUrl,omitempty"`
	ImageID    *string                `json:"imageId,omitempty"`
	Attributes []AttributeValueDetail `json:"attributeValues,omitempty"`
	CreatedAt  string                 `json:"createdAt,omitempty"`
	UpdatedAt  string                 `json:"updatedAt,omitempty"`
}
type ProductDetail struct {
	ID                 string  `json:"id"`
	Name               string  `json:"name"`
	Description        string  `json:"description"`
	ShortDescription   *string `json:"shortDescription"`
	BasePrice          float64 `json:"price,omitzero"`
	BaseSku            string  `json:"sku"`
	IsActive           bool    `json:"isActive"`
	Slug               string  `json:"slug"`
	ImageUrl           *string `json:"imageUrl,omitempty"`
	ImageId            *string `json:"imageId,omitempty"`
	DiscountPercentage *int16  `json:"discountPercentage,omitempty"`

	RatingCount    int32 `json:"ratingCount"`
	OneStarCount   int32 `json:"oneStarCount"`
	TwoStarCount   int32 `json:"twoStarCount"`
	ThreeStarCount int32 `json:"threeStarCount"`
	FourStarCount  int32 `json:"fourStarCount"`
	FiveStarCount  int32 `json:"fiveStarCount"`

	UpdatedAt string `json:"updatedAt"`
	CreatedAt string `json:"createdAt"`

	Brand       GeneralCategory    `json:"brand,omitempty"`
	Attributes  []ProductAttribute `json:"attributes,omitempty"`
	Collections []GeneralCategory  `json:"collections,omitempty"`
	Categories  []GeneralCategory  `json:"categories,omitempty"`
	Variations  []VariantDetail    `json:"variants,omitempty"`
}

type CartItemDetail struct {
	ID             string            `json:"id" binding:"required,uuid"`
	ProductID      string            `json:"productId" binding:"required,uuid"`
	VariantID      string            `json:"variantId" binding:"required,uuid"`
	Name           string            `json:"name"`
	Quantity       int16             `json:"quantity"`
	DiscountAmount float64           `json:"discountAmount"`
	Price          float64           `json:"price"`
	StockQty       int32             `json:"stock"`
	Sku            *string           `json:"sku,omitempty"`
	ImageURL       *string           `json:"imageUrl,omitempty"`
	Attributes     []AttributeDetail `json:"attributes"`
}
