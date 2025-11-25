package models

type CreateAddress struct {
	Phone     string  `json:"phone" binding:"required,min=10,max=15"`
	Street    string  `json:"street" binding:"required"`
	District  string  `json:"district" binding:"required"`
	City      string  `json:"city" binding:"required"`
	Ward      *string `json:"ward,omitempty" binding:"omitempty,max=100"`
	IsDefault bool    `json:"isDefault,omitempty" binding:"omitempty"`
}

type UpdateAddress struct {
	Phone     *string `json:"phone" binding:"omitempty"`
	Address   *string `json:"address1" binding:"omitempty"`
	Ward      *string `json:"ward" binding:"omitempty"`
	District  *string `json:"district" binding:"omitempty"`
	City      *string `json:"city" binding:"omitempty"`
	IsDefault *bool   `json:"isDefault" binding:"omitempty"`
}
