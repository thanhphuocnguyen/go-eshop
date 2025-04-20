package repository

import "github.com/google/uuid"

const (
	// ProductImageType is the type of image for product
	ProductImageType        = "product"
	ProductVariantImageType = "product_variant"
)

type ProductAttributesTxParam struct {
	// AttributeID int32 `json:"attribute_id" binding:"required"`
	ValueID int32 `json:"value_id" binding:"required"`
}

type CreateProductVariantTxParams struct {
	Price      float64                    `json:"price" binding:"required,gt=0"`
	Stock      int32                      `json:"stock" binding:"required,gt=0"`
	Weight     *float64                   `json:"weight" binding:"omitempty,gt=0"`
	IsActive   *bool                      `json:"is_active" binding:"omitempty"`
	Attributes []ProductAttributesTxParam `json:"attributes" binding:"min=1,dive"`
}

type UpdateProductVariantTxParams struct {
	ID         *string                     `json:"id" binding:"omitempty,uuid"`
	Price      *float64                    `json:"price" binding:"required,gt=0"`
	Stock      *int32                      `json:"stock" binding:"required,gt=0"`
	Weight     *float64                    `json:"weight" binding:"omitempty,gt=0"`
	IsActive   *bool                       `json:"is_active" binding:"omitempty"`
	Attributes *[]ProductAttributesTxParam `json:"attributes" binding:"omitempty,min=1,dive"`
}

type ProductTxResult struct {
	ID       uuid.UUID   `json:"id"`
	Variants []uuid.UUID `json:"variants"`
}
