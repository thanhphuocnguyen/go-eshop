package models

type OrderStatusModel struct {
	Status string `json:"status" validate:"required,oneof=pending confirmed delivering delivered completed"`
	Reason string `json:"reason,omitempty"`
}

type RefundOrderModel struct {
	Reason string `json:"reason" validate:"required"`
}

type CancelOrderModel struct {
	Reason string `json:"reason" validate:"required"`
}
type OrderListQuery struct {
	PaginationQuery
	Status        *string `form:"status,omitempty" validate:"omitempty,oneof=pending confirmed delivering delivered completed cancelled refunded"`
	PaymentStatus *string `form:"paymentStatus,omitempty" validate:"omitempty,oneof=pending succeeded failed cancelled refunded"`
}
