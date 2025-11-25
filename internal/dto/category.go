package dto

type AdminCategoryDetail struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description *string           `json:"description,omitempty"`
	Slug        string            `json:"slug"`
	Published   bool              `json:"published,omitempty"`
	CreatedAt   string            `json:"createdAt,omitempty"`
	UpdatedAt   string            `json:"updatedAt,omitempty"`
	ImageUrl    *string           `json:"imageUrl,omitempty"`
	Products    []ProductListItem `json:"products,omitempty"`
}

type CategoryDetail struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description *string          `json:"description,omitempty"`
	Slug        string           `json:"slug"`
	Published   bool             `json:"published,omitempty"`
	CreatedAt   string           `json:"createdAt,omitempty"`
	ImageUrl    *string          `json:"imageUrl,omitempty"`
	Products    []ProductSummary `json:"products"`
}

type GeneralCategory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
