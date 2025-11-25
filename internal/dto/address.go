package dto

import "time"

type AddressDetail struct {
	ID        string    `json:"id"`
	Default   bool      `json:"default"`
	CreatedAt time.Time `json:"createdAt"`
	Phone     string    `json:"phone"`
	Street    string    `json:"street"`
	Ward      *string   `json:"ward,omitempty"`
	District  string    `json:"district"`
	City      string    `json:"city"`
}
