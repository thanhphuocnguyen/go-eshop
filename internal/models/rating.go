package models

import (
	"mime/multipart"
)

type PostHelpfulRatingModel struct {
	Helpful bool `json:"helpful"`
}

type PostReplyRatingModel struct {
	RatingID string `json:"ratingId" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

type PostRatingFormData struct {
	OrderItemID string                  `form:"orderItemId" binding:"required"`
	Rating      float64                 `form:"rating" binding:"required,min=1,max=5"`
	Title       string                  `form:"title" binding:"required"`
	Content     string                  `form:"content" binding:"required"`
	Files       []*multipart.FileHeader `form:"files" binding:"omitempty"`
}
type RatingsQueryParams struct {
	PaginationQuery
	Status *string `form:"status" binding:"omitempty,oneof=approved rejected pending"`
}
