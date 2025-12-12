package models

type CheckoutModel struct {
	PaymentMethodId string   `json:"paymentMethodId" binding:"required,uuid"`
	DiscountCodes   []string `json:"discountCodes" binding:"omitempty"`
}

type UpdateCartItemQtyModel struct {
	Quantity int16 `json:"quantity" binding:"required,gt=0"`
}
