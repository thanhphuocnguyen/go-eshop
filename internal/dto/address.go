package dto

import (
	"time"

	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
)

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

func MapAddressResponse(address repository.UserAddress) AddressDetail {
	return AddressDetail{
		ID:        address.ID.String(),
		Default:   address.IsDefault,
		CreatedAt: address.CreatedAt,
		Phone:     address.PhoneNumber,
		Street:    address.Street,
		Ward:      address.Ward,
		District:  address.District,
		City:      address.City,
	}
}
