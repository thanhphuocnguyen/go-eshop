package models

import (
	"mime/multipart"
)

type CategoryProductRequest struct {
	SortOrder int16 `json:"sortOrder,omitempty"`
}

type CreateCategoryModel struct {
	DisplayOrder *int16                `form:"displayOrder" binding:"omitempty"`
	Description  *string               `form:"description" binding:"omitempty,max=1000"`
	Name         string                `form:"name" binding:"required,min=3,max=255"`
	Slug         string                `form:"slug" binding:"required,min=3,max=255"`
	Image        *multipart.FileHeader `form:"image" binding:"omitempty"`
}

type UpdateCategoryModel struct {
	Name         *string               `form:"name" binding:"omitempty,min=3,max=255"`
	Description  *string               `form:"description" binding:"omitempty,max=1000"`
	Slug         *string               `form:"slug" binding:"omitempty,min=3,max=255"`
	Published    *bool                 `form:"published" binding:"omitempty"`
	DisplayOrder *int16                `form:"displayOrder" binding:"omitempty"`
	Image        *multipart.FileHeader `form:"image" binding:"omitempty"`
}

type CollectionsQueryParams struct {
	PaginationQuery
	Collections *[]int32 `form:"collectionIds,omitempty"`
}
