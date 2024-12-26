// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

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
	PaymentGatewayRazorpay   PaymentGateway = "razorpay"
	PaymentGatewayVisa       PaymentGateway = "visa"
	PaymentGatewayMastercard PaymentGateway = "mastercard"
	PaymentGatewayAmex       PaymentGateway = "amex"
	PaymentGatewayApplePay   PaymentGateway = "apple_pay"
	PaymentGatewayGooglePay  PaymentGateway = "google_pay"
	PaymentGatewayAmazonPay  PaymentGateway = "amazon_pay"
	PaymentGatewayPhonePe    PaymentGateway = "phone_pe"
	PaymentGatewayPaytm      PaymentGateway = "paytm"
	PaymentGatewayUpi        PaymentGateway = "upi"
	PaymentGatewayWallet     PaymentGateway = "wallet"
	PaymentGatewayCod        PaymentGateway = "cod"
	PaymentGatewayPostpaid   PaymentGateway = "postpaid"
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
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodPaypal     PaymentMethod = "paypal"
	PaymentMethodCod        PaymentMethod = "cod"
	PaymentMethodDebitCard  PaymentMethod = "debit_card"
	PaymentMethodApplePay   PaymentMethod = "apple_pay"
	PaymentMethodWallet     PaymentMethod = "wallet"
	PaymentMethodPostpaid   PaymentMethod = "postpaid"
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
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusSuccess PaymentStatus = "success"
	PaymentStatusFailed  PaymentStatus = "failed"
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
	ID          int32  `json:"id"`
	AttributeID int32  `json:"attribute_id"`
	Value       string `json:"value"`
}

type Cart struct {
	ID        int32     `json:"id"`
	UserID    int64     `json:"user_id"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type CartItem struct {
	ID        int32     `json:"id"`
	ProductID int64     `json:"product_id"`
	CartID    int32     `json:"cart_id"`
	Quantity  int16     `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}

type Category struct {
	ID        int32       `json:"id"`
	Name      string      `json:"name"`
	SortOrder int16       `json:"sort_order"`
	ImageUrl  pgtype.Text `json:"image_url"`
	Published bool        `json:"published"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type CategoryProduct struct {
	CategoryID int32 `json:"category_id"`
	ProductID  int64 `json:"product_id"`
}

type Image struct {
	ImageID      int32            `json:"image_id"`
	ProductID    pgtype.Int8      `json:"product_id"`
	VariantID    pgtype.Int8      `json:"variant_id"`
	ImageUrl     string           `json:"image_url"`
	CloudinaryID pgtype.Text      `json:"cloudinary_id"`
	IsPrimary    pgtype.Bool      `json:"is_primary"`
	CreatedAt    pgtype.Timestamp `json:"created_at"`
	UpdatedAt    pgtype.Timestamp `json:"updated_at"`
}

type Order struct {
	ID            int64              `json:"id"`
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
	ID        int64          `json:"id"`
	ProductID int64          `json:"product_id"`
	OrderID   int64          `json:"order_id"`
	Quantity  int32          `json:"quantity"`
	Price     pgtype.Numeric `json:"price"`
	CreatedAt time.Time      `json:"created_at"`
}

type Payment struct {
	ID            int32              `json:"id"`
	OrderID       int64              `json:"order_id"`
	Amount        pgtype.Numeric     `json:"amount"`
	Method        PaymentMethod      `json:"method"`
	Status        PaymentStatus      `json:"status"`
	Gateway       NullPaymentGateway `json:"gateway"`
	TransactionID pgtype.Text        `json:"transaction_id"`
	CreatedAt     pgtype.Timestamp   `json:"created_at"`
	UpdatedAt     pgtype.Timestamp   `json:"updated_at"`
}

type Product struct {
	ID          int64          `json:"id"`
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
	ID        int64            `json:"id"`
	ProductID int64            `json:"product_id"`
	Sku       string           `json:"sku"`
	Price     pgtype.Numeric   `json:"price"`
	Stock     pgtype.Int8      `json:"stock"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}

type Session struct {
	ID           uuid.UUID `json:"id"`
	UserID       int64     `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiredAt    time.Time `json:"expired_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type User struct {
	ID                int64     `json:"id"`
	Role              UserRole  `json:"role"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	Phone             string    `json:"phone"`
	FullName          string    `json:"full_name"`
	HashedPassword    string    `json:"hashed_password"`
	VerifiedEmail     bool      `json:"verified_email"`
	VerifiedPhone     bool      `json:"verified_phone"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type UserAddress struct {
	ID        int64              `json:"id"`
	UserID    int64              `json:"user_id"`
	Phone     string             `json:"phone"`
	Address1  string             `json:"address_1"`
	Address2  pgtype.Text        `json:"address_2"`
	Ward      pgtype.Text        `json:"ward"`
	District  string             `json:"district"`
	City      string             `json:"city"`
	IsPrimary bool               `json:"is_primary"`
	IsDeleted bool               `json:"is_deleted"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	DeletedAt pgtype.Timestamptz `json:"deleted_at"`
}

type VariantAttribute struct {
	ID        int32 `json:"id"`
	VariantID int32 `json:"variant_id"`
	ValueID   int32 `json:"value_id"`
}
