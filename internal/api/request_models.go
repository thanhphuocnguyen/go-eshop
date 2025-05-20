package api

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
)

type CreateCategoryRequest struct {
	Name         string                `form:"name" binding:"required,min=3,max=255"`
	Slug         string                `form:"slug" binding:"required,min=3,max=255"`
	Description  *string               `form:"description" binding:"omitempty,max=1000"`
	DisplayOrder *int16                `form:"display_order" binding:"omitempty"`
	Remarkable   *bool                 `form:"remarkable" binding:"omitempty"`
	Image        *multipart.FileHeader `form:"image" binding:"omitempty"`
}

type CategoryProductRequest struct {
	SortOrder int16 `json:"sort_order,omitempty"`
}

type PostHelpfulRatingRequest struct {
	Helpful bool `json:"helpful"`
}

type PostReplyRatingRequest struct {
	RatingID string `json:"rating_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

type PaymentRequest struct {
	OrderID       string `json:"order_id" binding:"required,uuid"`
	PaymentMethod string `json:"payment_method" binding:"required,oneof=cod stripe"`
}

type UpdatePaymentStatusRequest struct {
	Status repository.PaymentStatus `json:"status" binding:"required"`
}

type OrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending confirmed delivering delivered completed"`
	Reason string `json:"reason,omitempty"`
}

type RefundOrderRequest struct {
	Reason string `json:"reason" binding:"required,oneof=defective damaged fraudulent requested_by_customer"`
}

type CancelOrderRequest struct {
	Reason string `json:"reason" binding:"required,oneof=duplicate fraudulent requested_by_customer abandoned"`
}

type CreateAddressRequest struct {
	Phone     string  `json:"phone" binding:"required,min=10,max=15"`
	Street    string  `json:"street" binding:"required"`
	District  string  `json:"district" binding:"required"`
	City      string  `json:"city" binding:"required"`
	Ward      *string `json:"ward,omitempty" binding:"omitempty,max=100"`
	IsDefault bool    `json:"is_default,omitempty" binding:"omitempty"`
}

type UpdateAddressRequest struct {
	Phone     *string `json:"phone" binding:"omitempty"`
	Address   *string `json:"address_1" binding:"omitempty"`
	Ward      *string `json:"ward" binding:"omitempty"`
	District  *string `json:"district" binding:"omitempty"`
	City      *string `json:"city" binding:"omitempty"`
	IsDefault *bool   `json:"is_default" binding:"omitempty"`
}

type UpdateCartItemQtyRequest struct {
	Quantity int16 `json:"quantity" binding:"required,gt=0"`
}

type CheckoutRequest struct {
	PaymentMethod      string   `json:"payment_method" binding:"required,oneof=code stripe"`
	PaymentGateway     *string  `json:"payment_gateway" binding:"omitempty,oneof=stripe"`
	AddressID          *string  `json:"address_id" binding:"omitempty,uuid"`
	Email              *string  `json:"email" binding:"omitempty,email"`
	FullName           *string  `json:"full_name" binding:"omitempty"`
	Address            *Address `json:"address" binding:"omitempty"`
	PaymentRecipeEmail *string  `json:"payment_receipt_email" binding:"omitempty,email"`
}

type AssignmentRequest struct {
	VariantIDs []string `json:"variant_ids" binding:"required"`
}

type UpdateUserRequest struct {
	UserID   uuid.UUID `json:"user_id" binding:"required,uuid"`
	FullName *string   `json:"fullname,omitempty" binding:"omitempty,min=3,max=32"`
	Email    *string   `json:"email" binding:"email,max=255,min=6"`
	Phone    *string   `json:"phone" binding:"omitempty,min=8,max=15"`
}

type CollectionProductRequest struct {
	SortOrder int16 `json:"sort_order,omitempty"`
}

type UpdateCategoryRequest struct {
	Name         *string               `form:"name" binding:"omitempty,min=3,max=255"`
	Description  *string               `form:"description" binding:"omitempty,max=1000"`
	Slug         *string               `form:"slug" binding:"omitempty,min=3,max=255"`
	Published    *bool                 `form:"published" binding:"omitempty"`
	Remarkable   *bool                 `form:"remarkable" binding:"omitempty"`
	DisplayOrder *int16                `form:"display_order" binding:"omitempty"`
	Image        *multipart.FileHeader `form:"image" binding:"omitempty"`
}

type ProductQueries struct {
	Page         int64    `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize     int64    `form:"page_size,default=20" binding:"omitempty,min=1,max=100"`
	Search       *string  `form:"search" binding:"omitempty,max=1000"`
	CategoryIDs  []string `form:"category_ids" binding:"omitempty,uuidslice"`
	BrandID      *string  `form:"brand_id" binding:"omitempty,uuid4"`
	CollectionID *string  `form:"collection_id" binding:"omitempty,uuid4"`
}

type OrderListQuery struct {
	PaginationQueryParams
	Status        *string `form:"status,omitempty" binding:"omitempty,oneof=pending confirmed delivering delivered completed cancelled refunded"`
	PaymentStatus *string `form:"payment_status,omitempty" binding:"omitempty,oneof=pending succeeded failed cancelled refunded"`
}

type PostRatingFormData struct {
	OrderItemID string                  `form:"order_item_id" binding:"required"`
	Rating      float64                 `form:"rating" binding:"required,min=1,max=5"`
	Title       string                  `form:"title" binding:"required"`
	Content     string                  `form:"content" binding:"required"`
	Files       []*multipart.FileHeader `form:"files" binding:"omitempty"`
}

type CreateDiscountRequest struct {
	ID                string    `json:"id" binding:"required,uuid"`
	Code              string    `json:"code" binding:"required,min=5,max=32"`
	Description       string    `json:"description" binding:"required"`
	DiscountType      string    `json:"discount_type" binding:"required,oneof=percentage fixed_amount"`
	DiscountValue     float64   `json:"discount_value" binding:"required,gt=0"`
	MinPurchaseAmount *float64  `json:"min_purchase_amount" binding:"omitempty,gt=0"`
	MaxDiscountAmount *float64  `json:"max_discount_amount" binding:"omitempty,gt=0"`
	UsageLimit        *int32    `json:"usage_limit" binding:"omitempty,gte=0"`
	IsActive          bool      `json:"is_active" binding:"required"`
	StartsAt          time.Time `json:"starts_at" binding:"required"`
	ExpiresAt         time.Time `json:"expires_at" binding:"required"`

	// Related entities
	Products   []string `json:"products,omitempty" binding:"omitempty,uuidslice"`
	Categories []string `json:"categories,omitempty" binding:"omitempty,uuidslice"`
	Users      []string `json:"users,omitempty" binding:"omitempty,uuidslice"`
}

type DiscountListQuery struct {
	PaginationQueryParams
	Search        *string    `form:"search" binding:"omitempty,max=1000"`
	DiscountType  *string    `form:"discount_type" binding:"omitempty,oneof=percentage fixed_amount"`
	IsActive      *bool      `from:"is_active" binding:"omitempty"`
	DiscountValue *float64   `form:"discount_value" binding:"omitempty,gt=0"`
	FromDate      *time.Time `form:"from_date" binding:"omitempty"`
	ToDate        *time.Time `form:"to_date" binding:"omitempty"`
}
