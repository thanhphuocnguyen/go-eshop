package repository

import "github.com/google/uuid"

const (
	// ProductImageType is the type of image for product
	ProductImageType = "product"
	// ProductVariantImageType is the type of image for product variant
	ProductVariantImageType = "product_variant"
	// ProductVariantImageType is the type of image for product variant
)

type ProductAttributesTxParam struct {
	AttributeID int32   `json:"attribute_id"`
	ValueIDs    []int32 `json:"value_ids"`
}

type CreateProductVariantTxParams struct {
	Price      float64                    `json:"price" binding:"required,gt=0"`
	Stock      int32                      `json:"stock" binding:"required,gt=0"`
	Sku        string                     `json:"sku" binding:"required,alphanum"`
	Weight     *float64                   `json:"weight" binding:"omitempty,gt=0"`
	IsActive   *bool                      `json:"is_active" binding:"omitempty"`
	ImageUrl   *string                    `json:"image_url" binding:"omitempty"`
	Attributes []ProductAttributesTxParam `json:"attributes" binding:"omitempty,dive"`
}

type UpdateProductVariantTxParams struct {
	ID         *string                    `json:"id" binding:"omitempty,uuid"`
	Price      *float64                   `json:"price" binding:"required,gt=0"`
	Stock      *int32                     `json:"stock" binding:"required,gt=0"`
	Sku        *string                    `json:"sku" binding:"required,alphanum"`
	Weight     *float64                   `json:"weight" binding:"omitempty,gt=0"`
	IsActive   *bool                      `json:"is_active" binding:"omitempty"`
	ImageUrl   *string                    `json:"image_url" binding:"omitempty"`
	Attributes []ProductAttributesTxParam `json:"attributes" binding:"omitempty,dive"`
}

type ProductTxResult struct {
	ID       uuid.UUID   `json:"id"`
	Variants []uuid.UUID `json:"variants"`
}
