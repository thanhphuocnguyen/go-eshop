// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

import (
	"database/sql/driver"
	"fmt"
	"time"

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

type OrderStatus string

const (
	OrderStatusWaitForConfirming OrderStatus = "wait_for_confirming"
	OrderStatusConfirmed         OrderStatus = "confirmed"
	OrderStatusDelivering        OrderStatus = "delivering"
	OrderStatusDelivered         OrderStatus = "delivered"
	OrderStatusCancelled         OrderStatus = "cancelled"
	OrderStatusRefunded          OrderStatus = "refunded"
	OrderStatusCompleted         OrderStatus = "completed"
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

type PaymentStatus string

const (
	PaymentStatusNotPaid PaymentStatus = "not_paid"
	PaymentStatusPaid    PaymentStatus = "paid"
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

type PaymentType string

const (
	PaymentTypeCash     PaymentType = "cash"
	PaymentTypeTransfer PaymentType = "transfer"
)

func (e *PaymentType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PaymentType(s)
	case string:
		*e = PaymentType(s)
	default:
		return fmt.Errorf("unsupported scan type for PaymentType: %T", src)
	}
	return nil
}

type NullPaymentType struct {
	PaymentType PaymentType `json:"payment_type"`
	Valid       bool        `json:"valid"` // Valid is true if PaymentType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPaymentType) Scan(value interface{}) error {
	if value == nil {
		ns.PaymentType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PaymentType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPaymentType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PaymentType), nil
}

type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
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
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AttributeValeu struct {
	ID          int64       `json:"id"`
	AttributeID int64       `json:"attribute_id"`
	Value       string      `json:"value"`
	Color       pgtype.Text `json:"color"`
	CreatedAt   time.Time   `json:"created_at"`
}

type Cart struct {
	ID           int64              `json:"id"`
	CheckedOutAt pgtype.Timestamptz `json:"checked_out_at"`
	UserID       int64              `json:"user_id"`
	UpdatedAt    time.Time          `json:"updated_at"`
	CreatedAt    time.Time          `json:"created_at"`
}

type CartItem struct {
	ID        int64     `json:"id"`
	ProductID int64     `json:"product_id"`
	CartID    int64     `json:"cart_id"`
	Quantity  int32     `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}

type Category struct {
	ID        int64       `json:"id"`
	Name      string      `json:"name"`
	ImageUrl  pgtype.Text `json:"image_url"`
	Published bool        `json:"published"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type CategoryProduct struct {
	CategoryID int64 `json:"category_id"`
	ProductID  int64 `json:"product_id"`
}

type Order struct {
	ID            int64              `json:"id"`
	UserID        int64              `json:"user_id"`
	Status        OrderStatus        `json:"status"`
	ShippingID    pgtype.Int8        `json:"shipping_id"`
	PaymentType   PaymentType        `json:"payment_type"`
	PaymentStatus PaymentStatus      `json:"payment_status"`
	IsCod         bool               `json:"is_cod"`
	ConfirmedAt   pgtype.Timestamptz `json:"confirmed_at"`
	CancelledAt   pgtype.Timestamptz `json:"cancelled_at"`
	DeliveredAt   pgtype.Timestamptz `json:"delivered_at"`
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

type PaymentInfo struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	CardNumber  string    `json:"card_number"`
	ExpiredDate time.Time `json:"expired_date"`
	VccCode     string    `json:"vcc_code"`
	CardType    CardType  `json:"card_type"`
	IsVerified  bool      `json:"is_verified"`
	CreatedAt   time.Time `json:"created_at"`
}

type Product struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Sku         string         `json:"sku"`
	ImageUrl    string         `json:"image_url"`
	Stock       int32          `json:"stock"`
	Archived    bool           `json:"archived"`
	Price       pgtype.Numeric `json:"price"`
	UpdatedAt   time.Time      `json:"updated_at"`
	CreatedAt   time.Time      `json:"created_at"`
}

type Shipping struct {
	ID            int64          `json:"id"`
	Vendor        string         `json:"vendor"`
	OrderID       int64          `json:"order_id"`
	Fee           pgtype.Numeric `json:"fee"`
	Phone         string         `json:"phone"`
	EstimatedDays int32          `json:"estimated_days"`
	CreatedAt     time.Time      `json:"created_at"`
}

type User struct {
	ID                int64     `json:"id"`
	Role              UserRole  `json:"role"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	FullName          string    `json:"full_name"`
	HashedPassword    string    `json:"hashed_password"`
	VerifiedEmail     bool      `json:"verified_email"`
	VerifiedPhone     bool      `json:"verified_phone"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	CreatedAt         time.Time `json:"created_at"`
}
