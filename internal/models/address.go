package models

type CreateAddress struct {
	Phone     string  `json:"phone" validate:"required,min=10,max=15"`
	Street    string  `json:"street" validate:"required"`
	District  string  `json:"district" validate:"required"`
	City      string  `json:"city" validate:"required"`
	Ward      *string `json:"ward,omitempty" validate:"omitempty,max=100"`
	IsDefault bool    `json:"isDefault,omitempty" validate:"omitempty"`
}

type UpdateAddress struct {
	Phone     *string `json:"phone" validate:"omitempty"`
	Address   *string `json:"address1" validate:"omitempty"`
	Ward      *string `json:"ward" validate:"omitempty"`
	District  *string `json:"district" validate:"omitempty"`
	City      *string `json:"city" validate:"omitempty"`
	IsDefault *bool   `json:"isDefault" validate:"omitempty"`
}
