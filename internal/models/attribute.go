package models

type AttributeParam struct {
	ID      int32  `uri:"id" validate:"required"`
	ValueID *int64 `uri:"valueId" validate:"omitempty"`
}
type ProductAttributeValue struct {
	Value   string `json:"value"`
	ValueId int64  `json:"valueId"`
}

type AttributeModel struct {
	Name string `json:"name" validate:"required"`
}

type AttributeValueModel struct {
	Value string `json:"value" validate:"required"`
}

type AttributesQuery struct {
	IDs []int32 `form:"ids" validate:"omitempty"`
}
