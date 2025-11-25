package models

type ProductVariantParam struct {
	ID string `uri:"variantId" binding:"required,uuid"`
}
type URIVariantParam struct {
	ProductID string `uri:"id" binding:"required,uuid"`
	VariantID string `uri:"variantId" binding:"required,uuid"`
}
type ProductQuery struct {
	Page         int64     `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize     int64     `form:"pageSize,default=20" binding:"omitempty,min=1,max=100"`
	Search       *string   `form:"search" binding:"omitempty,max=1000"`
	CategoryIDs  *[]string `form:"categoryIds" binding:"omitempty,uuidslice"`
	BrandIDs     *[]string `form:"brandIds" binding:"omitempty,uuid4"`
	CollectionID *[]string `form:"collectionIds" binding:"omitempty,uuid4"`
}

type CreateProductModel struct {
	BasePrice          float64 `json:"price" binding:"required,gt=0"`
	DiscountPercentage *int16  `json:"discountPercentage" binding:"omitempty,gte=0,lte=100"`
	ShortDescription   *string `json:"shortDescription" binding:"omitempty,max=1000"`
	Name               string  `json:"name" binding:"required,min=3,max=255"`
	Description        string  `json:"description" binding:"required"`
	BaseSku            string  `json:"sku" binding:"required,min=3,max=100"`
	Slug               string  `json:"slug" binding:"required,min=3,max=255"`
	BrandID            string  `json:"brandId" binding:"required,uuid"`

	Attributes    []int32  `json:"attributes" binding:"omitempty"`
	CategoryIDs   []string `json:"categoryIds" binding:"omitempty,uuidslice"`
	CollectionIDs []string `json:"collectionIds" binding:"omitempty,uuidslice"`
}

type UpdateProductModel struct {
	IsActive           *bool    `json:"isActive" binding:"omitempty"`
	BasePrice          *float64 `json:"price" binding:"omitempty,gt=0"`
	DiscountPercentage *int16   `json:"discountPercentage" binding:"omitempty,gte=0,lte=100"`
	Name               *string  `json:"name" binding:"omitempty,min=3,max=255"`
	Description        *string  `json:"description" binding:"omitempty"`
	ShortDescription   *string  `json:"shortDescription" binding:"omitempty,max=1000"`
	BaseSku            *string  `json:"sku" binding:"omitempty,min=3,max=100"`
	Slug               *string  `json:"slug" binding:"omitempty,min=3,max=255"`
	BrandID            *string  `json:"brandId" binding:"omitempty,uuid"`

	Attributes    *[]int32  `json:"attributes" binding:"omitempty"`
	CategoryIDs   *[]string `json:"categoryIds" binding:"omitempty,uuidslice"`
	CollectionIDs *[]string `json:"collectionIds" binding:"omitempty,uuidslice"`
}

type CreateProdVariantModel struct {
	Price           float64  `json:"price" binding:"required,gt=0"`
	StockQty        int32    `json:"stockQty" binding:"required,gte=0"`
	IsActive        bool     `json:"isActive" binding:"required"`
	AttributeValues []int64  `json:"attributeValues" binding:"required"`
	Description     *string  `json:"description" binding:"omitempty"`
	Weight          *float64 `json:"weight" binding:"omitempty,gt=0"`
}

type UpdateProdVariantModel struct {
	Price           *float64 `json:"price" binding:"omitempty,gt=0"`
	StockQty        *int32   `json:"stockQty" binding:"omitempty,gte=0"`
	IsActive        *bool    `json:"isActive" binding:"omitempty"`
	Description     *string  `json:"description" binding:"omitempty"`
	AttributeValues *[]int64 `json:"attributeValues" binding:"required,omitempty"`
	Weight          *float64 `json:"weight" binding:"omitempty,gt=0"`
}
