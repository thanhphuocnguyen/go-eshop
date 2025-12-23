package models

type CheckoutModel struct {
	PaymentMethodId string   `json:"paymentMethodId" validate:"required,uuid"`
	DiscountCodes   []string `json:"discountCodes" validate:"omitempty"`
}

type UpdateCartItemQtyModel struct {
	Quantity int16 `json:"quantity" validate:"required,gt=0"`
}
