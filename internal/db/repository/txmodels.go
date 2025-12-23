package repository

import "github.com/google/uuid"

type ProductAttributesTxParam struct {
	ID      string `json:"id" validate:"required,uuid"`
	ValueID int64  `json:"value_id" validate:"required"`
}

type CreateProductVariantTxParams struct {
	Price      float64                    `json:"price" validate:"required,gt=0"`
	Stock      int32                      `json:"stock_qty" validate:"required,gte=0"`
	Weight     *float64                   `json:"weight" validate:"omitempty,gte=0"`
	IsActive   *bool                      `json:"is_active" validate:"omitempty"`
	Attributes []ProductAttributesTxParam `json:"attributes" validate:"min=1,dive"`
}

type UpdateProductVariantTxParams struct {
	ID         *string                    `json:"id" validate:"omitempty,uuid"`
	Price      *float64                   `json:"price" validate:"required,gt=0"`
	Stock      *int32                     `json:"stock_qty" validate:"required,gte=0"`
	Weight     *float64                   `json:"weight" validate:"omitempty,gte=0"`
	IsActive   *bool                      `json:"is_active" validate:"omitempty"`
	Attributes []ProductAttributesTxParam `json:"attributes" validate:"min=1,dive"`
}

type ProductTxResult struct {
	ID uuid.UUID `json:"id"`
}

type CreateProductTxArgs struct {
	Product       CreateProductParams
	Attributes    []int32
	CategoryIDs   []string
	CollectionIDs []string
}

type UpdateProductTxArgs struct {
	Product       UpdateProductParams
	Attributes    *[]int32
	CategoryIDs   *[]string
	CollectionIDs *[]string
}

type AttributeDataSnapshot struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ShippingAddressSnapshot struct {
	Street   string `json:"street" validate:"required"`
	Ward     string `json:"ward" validate:"required"`
	District string `json:"district" validate:"required"`
	City     string `json:"city" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
}
