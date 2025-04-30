package repository

import "github.com/google/uuid"

type ProductAttributesTxParam struct {
	// AttributeID int32 `json:"attribute_id" binding:"required"`
	ID      int32 `json:"id" binding:"required"`
	ValueID int32 `json:"value_id" binding:"required"`
}

type CreateProductVariantTxParams struct {
	Price      float64                    `json:"price" binding:"required,gt=0"`
	Stock      int32                      `json:"stock_qty" binding:"required,gte=0"`
	Weight     *float64                   `json:"weight" binding:"omitempty,gte=0"`
	IsActive   *bool                      `json:"is_active" binding:"omitempty"`
	Attributes []ProductAttributesTxParam `json:"attributes" binding:"min=1,dive"`
}

type UpdateProductVariantTxParams struct {
	ID         *string                    `json:"id" binding:"omitempty,uuid"`
	Price      *float64                   `json:"price" binding:"required,gt=0"`
	Stock      *int32                     `json:"stock_qty" binding:"required,gte=0"`
	Weight     *float64                   `json:"weight" binding:"omitempty,gte=0"`
	IsActive   *bool                      `json:"is_active" binding:"omitempty"`
	Attributes []ProductAttributesTxParam `json:"attributes" binding:"min=1,dive"`
}

type ProductTxResult struct {
	ID uuid.UUID `json:"id"`
}

type AttributeDataSnapshot struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ShippingAddressSnapshot struct {
	Street   string `json:"street" binding:"required"`
	Ward     string `json:"ward" binding:"required"`
	District string `json:"district" binding:"required"`
	City     string `json:"city" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
}
