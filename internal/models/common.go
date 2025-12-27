package models

type URISlugParam struct {
	Slug string `uri:"slug" validate:"required"`
}

type UriIDParam struct {
	ID string `uri:"id" validate:"required,uuid"`
}
type PublicIDParam struct {
	PublicID string `uri:"publicId" validate:"required"`
}

type PaginationQuery struct {
	Page     int64   `form:"page,default=1" validate:"omitempty,min=1"`
	PageSize int64   `form:"pageSize,default=20" validate:"omitempty,min=1,max=100"`
	Search   *string `form:"search" validate:"omitempty,omitzero,max=1000"`
}
