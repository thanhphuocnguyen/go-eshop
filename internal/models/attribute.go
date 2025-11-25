package models

type AttributeParam struct {
	ID      int32  `uri:"id" binding:"required"`
	ValueID *int64 `uri:"valueId" binding:"omitempty"`
}
type ProductAttributeValue struct {
	Value   string `json:"value"`
	ValueId int64  `json:"valueId"`
}

type AttributeModel struct {
	Name string `json:"name" binding:"required"`
}

type AttributeValueModel struct {
	Value string `json:"value" binding:"required"`
}

type AttributesQuery struct {
	IDs []int32 `form:"ids" binding:"omitempty"`
}
