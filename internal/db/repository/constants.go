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

// Payment method constants for backward compatibility
const (
	PaymentMethodCodeCreditCard   = "credit_card"
	PaymentMethodCodeDebitCard    = "debit_card"
	PaymentMethodCodePaypal       = "paypal"
	PaymentMethodCodeStripe       = "stripe"
	PaymentMethodCodeApplePay     = "apple_pay"
	PaymentMethodCodeBankTransfer = "bank_transfer"
	PaymentMethodCodeCOD          = "cod"
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

// User role constants for backward compatibility
type Role string

const (
	UserRoleCodeAdmin     Role = "admin"
	UserRoleCodeUser      Role = "user"
	UserRoleCodeModerator Role = "moderator"
)

type CartStatus string

const (
	CartStatusActive     CartStatus = "active"
	CartStatusCheckedOut CartStatus = "checked_out"
)

// EntityType reference table struct
type EntityType struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	IsActive    bool      `json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Entity type constants for backward compatibility
const (
	EntityTypeCodeProduct       = "product"
	EntityTypeCodeVariant       = "variant"
	EntityTypeCodeProductRating = "product_rating"
	EntityTypeCodeCategory      = "category"
	EntityTypeCodeCollection    = "collection"
	EntityTypeCodeBrand         = "brand"
	EntityTypeCodeUser          = "user"
)

// Nullable types - updated for new reference table structs
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

type NullEntityType struct {
	EntityType EntityType
	Valid      bool
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
