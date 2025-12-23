package dto

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

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
	ID             string            `json:"id" validate:"required,uuid"`
	ProductID      string            `json:"productId" validate:"required,uuid"`
	VariantID      string            `json:"variantId" validate:"required,uuid"`
	Name           string            `json:"name"`
	Quantity       int16             `json:"quantity"`
	DiscountAmount float64           `json:"discountAmount"`
	Price          float64           `json:"price"`
	StockQty       int32             `json:"stock"`
	Sku            *string           `json:"sku,omitempty"`
	ImageURL       *string           `json:"imageUrl,omitempty"`
	Attributes     []AttributeDetail `json:"attributes"`
}

func MapToProductDetailResponse(row repository.GetProductDetailRow) ProductDetail {
	basePrice, _ := row.BasePrice.Float64Value()

	resp := ProductDetail{
		ID:                 row.ID.String(),
		Name:               row.Name,
		BasePrice:          basePrice.Float64,
		ShortDescription:   row.ShortDescription,
		Description:        row.Description,
		BaseSku:            row.BaseSku,
		Slug:               row.Slug,
		RatingCount:        row.RatingCount,
		OneStarCount:       row.OneStarCount,
		TwoStarCount:       row.TwoStarCount,
		ThreeStarCount:     row.ThreeStarCount,
		FourStarCount:      row.FourStarCount,
		FiveStarCount:      row.FiveStarCount,
		DiscountPercentage: row.DiscountPercentage,
		UpdatedAt:          row.UpdatedAt.String(),
		CreatedAt:          row.CreatedAt.String(),

		IsActive: *row.IsActive,
		ImageUrl: row.ImageUrl,
		ImageId:  row.ImageID,

		// Initialize slices
		Categories:  []GeneralCategory{},
		Collections: []GeneralCategory{},
		Attributes:  []ProductAttribute{},
		Brand:       GeneralCategory{},
		Variations:  []VariantDetail{},
	}

	// Unmarshal JSON data
	if err := json.Unmarshal(row.Attributes, &resp.Attributes); err != nil {
		log.Error().Err(err).Msg("Unmarshal attributes")
	}
	if err := json.Unmarshal(row.Categories, &resp.Categories); err != nil {
		log.Error().Err(err).Msg("Unmarshal categories")
	}
	if err := json.Unmarshal(row.Collections, &resp.Collections); err != nil {
		log.Error().Err(err).Msg("Unmarshal collections")
	}
	if err := json.Unmarshal(row.Brand, &resp.Brand); err != nil {
		log.Error().Err(err).Msg("Unmarshal brand")
	}
	if err := json.Unmarshal(row.Variants, &resp.Variations); err != nil {
		log.Error().Err(err).Msg("Unmarshal variants")
	}

	return resp
}

func MapToAdminProductResponse(productRow repository.Product) ProductListItem {
	basePrice, _ := productRow.BasePrice.Float64Value()

	avgRating := utils.GetAvgRating(
		productRow.RatingCount,
		productRow.OneStarCount,
		productRow.TwoStarCount,
		productRow.ThreeStarCount,
		productRow.FourStarCount,
		productRow.FiveStarCount,
	)

	product := ProductListItem{
		ID:                 productRow.ID.String(),
		Name:               productRow.Name,
		Description:        productRow.Description,
		BasePrice:          basePrice.Float64,
		Sku:                productRow.BaseSku,
		Slug:               productRow.Slug,
		AvgRating:          &avgRating,
		ImageUrl:           productRow.ImageUrl,
		ImgID:              productRow.ImageID,
		ReviewCount:        &productRow.RatingCount,
		CreatedAt:          productRow.CreatedAt.String(),
		UpdatedAt:          productRow.UpdatedAt.String(),
		DiscountPercentage: productRow.DiscountPercentage,
	}

	return product
}

func MapToShopProductResponse(productRow repository.GetProductListRow) ProductSummary {
	price, _ := productRow.MinPrice.Float64Value()
	avgRating := utils.GetAvgRating(
		productRow.RatingCount,
		productRow.OneStarCount,
		productRow.TwoStarCount,
		productRow.ThreeStarCount,
		productRow.FourStarCount,
		productRow.FiveStarCount,
	)

	product := ProductSummary{
		ID:           productRow.ID.String(),
		Name:         productRow.Name,
		Price:        price.Float64,
		VariantCount: int16(productRow.VariantCount),
		Slug:         productRow.Slug,
		AvgRating:    &avgRating,
		ImageUrl:     productRow.ImageUrl,
		ImageID:      productRow.ImageID,
		ReviewCount:  &productRow.RatingCount,
		CreatedAt:    productRow.CreatedAt.String(),
		UpdatedAt:    productRow.UpdatedAt.String(),
	}

	return product
}

func MapToVariantListModelDto(row repository.GetProductVariantListRow) VariantDetail {
	price, _ := row.Price.Float64Value()
	variant := VariantDetail{
		ID:       row.ID.String(),
		Price:    price.Float64,
		Stock:    row.Stock,
		IsActive: *row.IsActive,
		Sku:      row.Sku,
		ImageUrl: row.ImageUrl,
	}
	variant.Attributes = []AttributeValueDetail{}
	err := json.Unmarshal(row.AttributeValues, &variant.Attributes)
	if err != nil {
		log.Error().Err(err).Msg("Unmarshal variant attribute values")
	}
	if row.Weight.Valid {
		weight, _ := row.Weight.Float64Value()
		variant.Weight = &weight.Float64
	}

	return variant
}
