package repository

import (
	"time"

	"github.com/google/uuid"
)

const (
	ThumbnailRole = "thumbnail"
	SmallRole     = "small"
	GalleryRole   = "gallery"
	SliderRole    = "slider"
	IconRole      = "icon"
	LogoRole      = "logo"
	AvatarRole    = "avatar"
)

type DiscountType string

const (
	PercentageDiscount  DiscountType = "percentage"
	FixedAmountDiscount DiscountType = "fixed_amount"
)

// Enum types for database enums
type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusSuccess    PaymentStatus = "success"
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusCancelled  PaymentStatus = "cancelled"
	PaymentStatusRefunded   PaymentStatus = "refunded"
	PaymentStatusProcessing PaymentStatus = "processing"
)

type PaymentMethod string

const (
	PaymentMethodCreditCard   PaymentMethod = "credit_card"
	PaymentMethodDebitCard    PaymentMethod = "debit_card"
	PaymentMethodPaypal       PaymentMethod = "paypal"
	PaymentMethodStripe       PaymentMethod = "stripe"
	PaymentMethodApplePay     PaymentMethod = "apple_pay"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodCOD          PaymentMethod = "cod"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefunded   OrderStatus = "refunded"
)

type UserRole string

const (
	UserRoleAdmin     UserRole = "admin"
	UserRoleUser      UserRole = "user"
	UserRoleModerator UserRole = "moderator"
)

type CartStatus string

const (
	CartStatusActive     CartStatus = "active"
	CartStatusCheckedOut CartStatus = "checked_out"
)

type CardType string

const (
	CardTypeDebit  CardType = "debit"
	CardTypeCredit CardType = "credit"
)

type EntityType string

const (
	EntityTypeProduct    EntityType = "product"
	EntityTypeCategory   EntityType = "category"
	EntityTypeCollection EntityType = "collection"
	EntityTypeBrand      EntityType = "brand"
	EntityTypeUser       EntityType = "user"
)

// Nullable types
type NullPaymentStatus struct {
	PaymentStatus PaymentStatus
	Valid         bool
}

type NullOrderStatus struct {
	OrderStatus OrderStatus
	Valid       bool
}

type NullUserRole struct {
	UserRole UserRole
	Valid    bool
}

type NullPaymentMethod struct {
	PaymentMethod PaymentMethod
	Valid         bool
}

// VerifyEmail type for email verification
type VerifyEmail struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Email      string
	VerifyCode string
	IsUsed     bool
	CreatedAt  time.Time
	ExpiredAt  time.Time
}
