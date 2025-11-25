package models

type OrderStatusModel struct {
	Status string `json:"status" binding:"required,oneof=pending confirmed delivering delivered completed"`
	Reason string `json:"reason,omitempty"`
}

type RefundOrderModel struct {
	Reason string `json:"reason" binding:"required"`
}

type CancelOrderModel struct {
	Reason string `json:"reason" binding:"required"`
}
type OrderListQuery struct {
	PaginationQuery
	Status        *string `form:"status,omitempty" binding:"omitempty,oneof=pending confirmed delivering delivered completed cancelled refunded"`
	PaymentStatus *string `form:"paymentStatus,omitempty" binding:"omitempty,oneof=pending succeeded failed cancelled refunded"`
}
