package models

import (
	"mime/multipart"
)

type CreateCategoryModel struct {
	Name         string                `form:"name" validate:"required,min=3,max=255"`
	Slug         string                `form:"slug" validate:"required,min=3,max=255"`
	DisplayOrder *int16                `form:"displayOrder" validate:"omitnil,omitempty"`
	Description  *string               `form:"description" validate:"omitnil,omitempty,max=1000"`
	Image        *multipart.FileHeader `form:"image" validate:"omitnil,omitempty"`
}

type UpdateCategoryModel struct {
	Name         *string               `form:"name" validate:"omitnil,omitempty,min=3,max=255"`
	Description  *string               `form:"description" validate:"omitnil,omitempty,max=1000"`
	Slug         *string               `form:"slug" validate:"omitnil,omitempty,min=3,max=255"`
	Published    *bool                 `form:"published" validate:"omitnil,omitempty"`
	DisplayOrder *int16                `form:"displayOrder" validate:"omitnil,omitempty"`
	Image        *multipart.FileHeader `form:"image" validate:"omitnil,omitempty"`
}
