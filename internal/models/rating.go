package models

import (
	"mime/multipart"
)

type PostHelpfulRatingModel struct {
	Helpful bool `json:"helpful"`
}

type PostReplyRatingModel struct {
	RatingID string `json:"ratingId" validate:"required"`
	Content  string `json:"content" validate:"required"`
}

type PostRatingFormData struct {
	OrderItemID string                  `form:"orderItemId" validate:"required"`
	Rating      float64                 `form:"rating" validate:"required,min=1,max=5"`
	Title       string                  `form:"title" validate:"required"`
	Content     string                  `form:"content" validate:"required"`
	Files       []*multipart.FileHeader `form:"files" validate:"omitempty"`
}
