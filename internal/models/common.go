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
