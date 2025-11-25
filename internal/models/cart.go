package models

type CheckoutModel struct {
	PaymentMethodId string   `json:"paymentMethod" binding:"required,oneof=cod stripe"`
	AddressId       string   `json:"addressId" binding:"required,uuid"`
	DiscountCodes   []string `json:"discountCodes" binding:"omitempty"`
}

type UpdateCartItemQtyModel struct {
	Quantity int16 `json:"quantity" binding:"required,gt=0"`
}
