package models

type AddCartItemModel struct {
	VariantID string `json:"variantId" binding:"required,uuid"`
	Quantity  int16  `json:"quantity" binding:"required,gt=0"`
}

type CheckoutModel struct {
	PaymentMethodId string   `json:"paymentMethod" binding:"required,oneof=cod stripe"`
	AddressId       string   `json:"addressId" binding:"required,uuid"`
	DiscountCodes   []string `json:"discountCodes" binding:"omitempty"`
}
