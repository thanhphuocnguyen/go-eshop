package dto

import "time"

type RatingDetail struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Rating    float64   `json:"rating"`
	CreatedAt time.Time `json:"createdAt"`
}

type RatingImage struct {
	ID  int64  `json:"id"`
	URL string `json:"url"`
}
type ProductRatingDetail struct {
	ID               string        `json:"id"`
	FirstName        string        `json:"firstName"`
	LastName         string        `json:"lastName"`
	ProductName      string        `json:"productName,omitempty"`
	UserID           string        `json:"userId"`
	Rating           float64       `json:"rating"`
	ReviewTitle      string        `json:"reviewTitle"`
	IsVisible        bool          `json:"isVisible"`
	IsApproved       bool          `json:"isApproved"`
	ReviewContent    string        `json:"reviewContent"`
	VerifiedPurchase bool          `json:"verifiedPurchase"`
	Count            int64         `json:"count"`
	Images           []RatingImage `json:"images"`
}
