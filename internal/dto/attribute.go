package dto

type AttributeValueDetail struct {
	ID    int64   `json:"id"`
	Value string  `json:"value"`
	Name  *string `json:"name,omitempty"`
}

type AttributeDetail struct {
	ID     int32                  `json:"id"`
	Name   string                 `json:"name"`
	Values []AttributeValueDetail `json:"values"`
}
