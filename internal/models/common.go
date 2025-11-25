package models

type URISlugParam struct {
	Slug string `uri:"slug" binding:"required"`
}

type UriIDParam struct {
	ID string `uri:"id" binding:"required,uuid"`
}
type PublicIDParam struct {
	PublicID string `uri:"publicId" binding:"required"`
}

type PaginationQuery struct {
	Page     int64   `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize int64   `form:"pageSize,default=20" binding:"omitempty,min=1,max=100"`
	Search   *string `form:"search" binding:"omitempty,omitzero,max=1000"`
}
