// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package repository

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CardType string

const (
	CardTypeDebit  CardType = "debit"
	CardTypeCredit CardType = "credit"
)

func (e *CardType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = CardType(s)
	case string:
		*e = CardType(s)
	default:
		return fmt.Errorf("unsupported scan type for CardType: %T", src)
	}
	return nil
}

type NullCardType struct {
	CardType CardType `json:"card_type"`
	Valid    bool     `json:"valid"` // Valid is true if CardType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullCardType) Scan(value interface{}) error {
	if value == nil {
		ns.CardType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.CardType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullCardType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.CardType), nil
}

type CartStatus string

const (
	CartStatusActive     CartStatus = "active"
	CartStatusCheckedOut CartStatus = "checked_out"
)

func (e *CartStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = CartStatus(s)
	case string:
		*e = CartStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for CartStatus: %T", src)
	}
	return nil
}

type NullCartStatus struct {
	CartStatus CartStatus `json:"cart_status"`
	Valid      bool       `json:"valid"` // Valid is true if CartStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullCartStatus) Scan(value interface{}) error {
	if value == nil {
		ns.CartStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.CartStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullCartStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.CartStatus), nil
}

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusDelivering OrderStatus = "delivering"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefunded   OrderStatus = "refunded"
	OrderStatusCompleted  OrderStatus = "completed"
)

func (e *OrderStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = OrderStatus(s)
	case string:
		*e = OrderStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for OrderStatus: %T", src)
	}
	return nil
}

type NullOrderStatus struct {
	OrderStatus OrderStatus `json:"order_status"`
	Valid       bool        `json:"valid"` // Valid is true if OrderStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullOrderStatus) Scan(value interface{}) error {
	if value == nil {
		ns.OrderStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.OrderStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullOrderStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.OrderStatus), nil
}

type PaymentGateway string

const (
	PaymentGatewayStripe     PaymentGateway = "stripe"
	PaymentGatewayPaypal     PaymentGateway = "paypal"
	PaymentGatewayVisa       PaymentGateway = "visa"
	PaymentGatewayMastercard PaymentGateway = "mastercard"
	PaymentGatewayApplePay   PaymentGateway = "apple_pay"
	PaymentGatewayGooglePay  PaymentGateway = "google_pay"
	PaymentGatewayPostpaid   PaymentGateway = "postpaid"
	PaymentGatewayMomo       PaymentGateway = "momo"
	PaymentGatewayZaloPay    PaymentGateway = "zalo_pay"
	PaymentGatewayVnPay      PaymentGateway = "vn_pay"
)

func (e *PaymentGateway) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PaymentGateway(s)
	case string:
		*e = PaymentGateway(s)
	default:
		return fmt.Errorf("unsupported scan type for PaymentGateway: %T", src)
	}
	return nil
}

type NullPaymentGateway struct {
	PaymentGateway PaymentGateway `json:"payment_gateway"`
	Valid          bool           `json:"valid"` // Valid is true if PaymentGateway is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPaymentGateway) Scan(value interface{}) error {
	if value == nil {
		ns.PaymentGateway, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PaymentGateway.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPaymentGateway) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PaymentGateway), nil
}

type PaymentMethod string

const (
	PaymentMethodCard         PaymentMethod = "card"
	PaymentMethodCod          PaymentMethod = "cod"
	PaymentMethodWallet       PaymentMethod = "wallet"
	PaymentMethodPostpaid     PaymentMethod = "postpaid"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
)

func (e *PaymentMethod) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PaymentMethod(s)
	case string:
		*e = PaymentMethod(s)
	default:
		return fmt.Errorf("unsupported scan type for PaymentMethod: %T", src)
	}
	return nil
}

type NullPaymentMethod struct {
	PaymentMethod PaymentMethod `json:"payment_method"`
	Valid         bool          `json:"valid"` // Valid is true if PaymentMethod is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPaymentMethod) Scan(value interface{}) error {
	if value == nil {
		ns.PaymentMethod, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PaymentMethod.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPaymentMethod) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PaymentMethod), nil
}

type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusSuccess    PaymentStatus = "success"
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusCancelled  PaymentStatus = "cancelled"
	PaymentStatusRefunded   PaymentStatus = "refunded"
	PaymentStatusProcessing PaymentStatus = "processing"
)

func (e *PaymentStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PaymentStatus(s)
	case string:
		*e = PaymentStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for PaymentStatus: %T", src)
	}
	return nil
}

type NullPaymentStatus struct {
	PaymentStatus PaymentStatus `json:"payment_status"`
	Valid         bool          `json:"valid"` // Valid is true if PaymentStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPaymentStatus) Scan(value interface{}) error {
	if value == nil {
		ns.PaymentStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PaymentStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPaymentStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PaymentStatus), nil
}

type UserRole string

const (
	UserRoleAdmin     UserRole = "admin"
	UserRoleUser      UserRole = "user"
	UserRoleModerator UserRole = "moderator"
)

func (e *UserRole) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = UserRole(s)
	case string:
		*e = UserRole(s)
	default:
		return fmt.Errorf("unsupported scan type for UserRole: %T", src)
	}
	return nil
}

type NullUserRole struct {
	UserRole UserRole `json:"user_role"`
	Valid    bool     `json:"valid"` // Valid is true if UserRole is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullUserRole) Scan(value interface{}) error {
	if value == nil {
		ns.UserRole, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.UserRole.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullUserRole) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.UserRole), nil
}

type Attribute struct {
	AttributeID int32  `json:"attribute_id"`
	Name        string `json:"name"`
}

type AttributeValue struct {
	AttributeValueID int32  `json:"attribute_value_id"`
	AttributeID      int32  `json:"attribute_id"`
	Value            string `json:"value"`
}

type Cart struct {
	CartID    int32     `json:"cart_id"`
	UserID    int64     `json:"user_id"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type CartItem struct {
	CartItemID int32     `json:"cart_item_id"`
	ProductID  int64     `json:"product_id"`
	CartID     int32     `json:"cart_id"`
	Quantity   int16     `json:"quantity"`
	CreatedAt  time.Time `json:"created_at"`
}

type Category struct {
	CategoryID int32       `json:"category_id"`
	Name       string      `json:"name"`
	SortOrder  int16       `json:"sort_order"`
	ImageUrl   pgtype.Text `json:"image_url"`
	Published  bool        `json:"published"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

type CategoryProduct struct {
	CategoryID int32 `json:"category_id"`
	ProductID  int64 `json:"product_id"`
	SortOrder  int16 `json:"sort_order"`
}

type Image struct {
	ImageID    int32            `json:"image_id"`
	ProductID  pgtype.Int8      `json:"product_id"`
	VariantID  pgtype.Int8      `json:"variant_id"`
	ImageUrl   string           `json:"image_url"`
	ExternalID pgtype.Text      `json:"external_id"`
	Primary    pgtype.Bool      `json:"primary"`
	CreatedAt  pgtype.Timestamp `json:"created_at"`
	UpdatedAt  pgtype.Timestamp `json:"updated_at"`
}

type Order struct {
	OrderID       int64              `json:"order_id"`
	UserID        int64              `json:"user_id"`
	UserAddressID int64              `json:"user_address_id"`
	TotalPrice    pgtype.Numeric     `json:"total_price"`
	Status        OrderStatus        `json:"status"`
	ConfirmedAt   pgtype.Timestamptz `json:"confirmed_at"`
	DeliveredAt   pgtype.Timestamptz `json:"delivered_at"`
	CancelledAt   pgtype.Timestamptz `json:"cancelled_at"`
	RefundedAt    pgtype.Timestamptz `json:"refunded_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
	CreatedAt     time.Time          `json:"created_at"`
}

type OrderItem struct {
	OrderItemID int64          `json:"order_item_id"`
	ProductID   int64          `json:"product_id"`
	OrderID     int64          `json:"order_id"`
	Quantity    int32          `json:"quantity"`
	Price       pgtype.Numeric `json:"price"`
}

type Payment struct {
	PaymentID      string             `json:"payment_id"`
	OrderID        int64              `json:"order_id"`
	Amount         pgtype.Numeric     `json:"amount"`
	PaymentMethod  PaymentMethod      `json:"payment_method"`
	Status         PaymentStatus      `json:"status"`
	PaymentGateway NullPaymentGateway `json:"payment_gateway"`
	RefundID       pgtype.Text        `json:"refund_id"`
	CreatedAt      pgtype.Timestamptz `json:"created_at"`
	UpdatedAt      pgtype.Timestamptz `json:"updated_at"`
}

type Product struct {
	ProductID   int64          `json:"product_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Sku         string         `json:"sku"`
	Stock       int32          `json:"stock"`
	Archived    bool           `json:"archived"`
	Price       pgtype.Numeric `json:"price"`
	UpdatedAt   time.Time      `json:"updated_at"`
	CreatedAt   time.Time      `json:"created_at"`
}

type ProductVariant struct {
	VariantID int64            `json:"variant_id"`
	ProductID int64            `json:"product_id"`
	Sku       string           `json:"sku"`
	Price     pgtype.Numeric   `json:"price"`
	Stock     pgtype.Int8      `json:"stock"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}

type Session struct {
	SessionID    uuid.UUID `json:"session_id"`
	UserID       int64     `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	Blocked      bool      `json:"blocked"`
	ExpiredAt    time.Time `json:"expired_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type User struct {
	UserID            int64     `json:"user_id"`
	Role              UserRole  `json:"role"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	Phone             string    `json:"phone"`
	Fullname          string    `json:"fullname"`
	HashedPassword    string    `json:"hashed_password"`
	VerifiedEmail     bool      `json:"verified_email"`
	VerifiedPhone     bool      `json:"verified_phone"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type UserAddress struct {
	UserAddressID int64       `json:"user_address_id"`
	UserID        int64       `json:"user_id"`
	Phone         string      `json:"phone"`
	Street        string      `json:"street"`
	Ward          pgtype.Text `json:"ward"`
	District      string      `json:"district"`
	City          string      `json:"city"`
	Default       bool        `json:"default"`
	Deleted       bool        `json:"deleted"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

type UserPaymentInfo struct {
	PaymentMethodID int32              `json:"payment_method_id"`
	UserID          pgtype.Int8        `json:"user_id"`
	CardNumber      string             `json:"card_number"`
	CardholderName  string             `json:"cardholder_name"`
	ExpirationDate  pgtype.Date        `json:"expiration_date"`
	BillingAddress  string             `json:"billing_address"`
	Default         pgtype.Bool        `json:"default"`
	CreatedAt       pgtype.Timestamptz `json:"created_at"`
	UpdatedAt       pgtype.Timestamptz `json:"updated_at"`
}

type VariantAttribute struct {
	VariantAttributeID int32 `json:"variant_attribute_id"`
	VariantID          int32 `json:"variant_id"`
	ValueID            int32 `json:"value_id"`
}
