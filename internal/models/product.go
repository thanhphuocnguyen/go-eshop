package models

type ProductVariantParam struct {
	VariantID string `uri:"variantId" validate:"required,uuid"`
}
type URIVariantParam struct {
	ProductID string `uri:"id" validate:"required,uuid"`
	VariantID string `uri:"variantId" validate:"required,uuid"`
}

type ProductQuery struct {
	Page         int64     `form:"page,default=1" validate:"omitempty,min=1"`
	PageSize     int64     `form:"pageSize,default=20" validate:"omitempty,min=1,max=100"`
	Search       *string   `form:"search" validate:"omitempty,max=1000"`
	CategoryIDs  *[]string `form:"categoryIds" validate:"omitnil,omitempty,uuidslice"`
	BrandIDs     *[]string `form:"brandIds" validate:"omitnil,omitempty,uuidslice"`
	CollectionID *[]string `form:"collectionIds" validate:"omitnil,omitempty,uuidslice"`
}

type CreateProductModel struct {
	BasePrice          float64 `json:"price" validate:"required,gt=0"`
	DiscountPercentage *int16  `json:"discountPercentage" validate:"omitempty,gte=0,lte=100"`
	ShortDescription   *string `json:"shortDescription" validate:"omitempty,max=1000"`
	Name               string  `json:"name" validate:"required,min=3,max=255"`
	Description        string  `json:"description" validate:"required"`
	BaseSku            string  `json:"sku" validate:"required,min=3,max=100"`
	Slug               string  `json:"slug" validate:"required,min=3,max=255"`
	BrandID            string  `json:"brandId" validate:"required,uuid"`

	Attributes    []int32  `json:"attributes" validate:"omitempty"`
	CategoryIDs   []string `json:"categoryIds" validate:"omitempty,uuidslice"`
	CollectionIDs []string `json:"collectionIds" validate:"omitempty,uuidslice"`

	Variants []CreateProdVariantModel `json:"variants" validate:"omitempty,dive,required"`
}

type UpdateProductModel struct {
	IsActive           *bool    `json:"isActive,omitempty" validate:"omitempty"`
	BasePrice          *float64 `json:"price,omitempty" validate:"omitempty,gt=0"`
	DiscountPercentage *int16   `json:"discountPercentage,omitempty" validate:"omitempty,gte=0,lte=100"`
	Name               *string  `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	Description        *string  `json:"description,omitempty" validate:"omitempty"`
	ShortDescription   *string  `json:"shortDescription,omitempty" validate:"omitempty,max=1000"`
	BaseSku            *string  `json:"sku,omitempty" validate:"omitempty,min=3,max=100"`
	Slug               *string  `json:"slug,omitempty" validate:"omitempty,min=3,max=255"`
	BrandID            *string  `json:"brandId,omitempty" validate:"omitempty,uuid"`

	Attributes    *[]int32  `json:"attributes,omitempty" validate:"omitempty"`
	CategoryIDs   *[]string `json:"categoryIds,omitempty" validate:"omitnil,omitempty,uuidslice"`
	CollectionIDs *[]string `json:"collectionIds,omitempty" validate:"omitnil,omitempty,uuidslice"`
}

type CreateProdVariantModel struct {
	Price           float64  `json:"price" validate:"required,gt=0"`
	StockQty        int32    `json:"stockQty" validate:"required,gte=0"`
	IsActive        bool     `json:"isActive" validate:"required"`
	AttributeValues []int64  `json:"attributeValues" validate:"required"`
	Description     string   `json:"description" validate:"omitempty"`
	Weight          *float64 `json:"weight" validate:"omitnil,omitempty,gt=0"`
}

type UpdateProdVariantModel struct {
	Price           *float64 `json:"price" validate:"omitempty,gt=0"`
	StockQty        *int32   `json:"stockQty" validate:"omitempty,gte=0"`
	IsActive        *bool    `json:"isActive" validate:"omitempty"`
	Description     *string  `json:"description" validate:"omitempty"`
	AttributeValues *[]int64 `json:"attributeValues" validate:"omitnil,omitempty"`
	Weight          *float64 `json:"weight" validate:"omitempty,gt=0"`
}
