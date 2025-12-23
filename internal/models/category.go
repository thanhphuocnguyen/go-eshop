package models

import (
	"mime/multipart"
)

type CreateCategoryModel struct {
	DisplayOrder *int16                `form:"displayOrder" validate:"omitempty"`
	Description  *string               `form:"description" validate:"omitempty,max=1000"`
	Name         string                `form:"name" validate:"required,min=3,max=255"`
	Slug         string                `form:"slug" validate:"required,min=3,max=255"`
	Image        *multipart.FileHeader `form:"image" validate:"omitempty"`
}

type UpdateCategoryModel struct {
	Name         *string               `form:"name" validate:"omitempty,min=3,max=255"`
	Description  *string               `form:"description" validate:"omitempty,max=1000"`
	Slug         *string               `form:"slug" validate:"omitempty,min=3,max=255"`
	Published    *bool                 `form:"published" validate:"omitempty"`
	DisplayOrder *int16                `form:"displayOrder" validate:"omitempty"`
	Image        *multipart.FileHeader `form:"image" validate:"omitempty"`
}
