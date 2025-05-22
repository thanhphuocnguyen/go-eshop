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
	DisplayOrder *int16                `form:"displayOrder" binding:"omitempty"`
	Remarkable   *bool                 `form:"remarkable" binding:"omitempty"`
	Image        *multipart.FileHeader `form:"image" binding:"omitempty"`
}

type CategoryProductRequest struct {
	SortOrder int16 `json:"sortOrder,omitempty"`
}

type PostHelpfulRatingRequest struct {
	Helpful bool `json:"helpful"`
}

type PostReplyRatingRequest struct {
	RatingID string `json:"ratingId" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

type PaymentRequest struct {
	OrderID       string `json:"orderId" binding:"required,uuid"`
	PaymentMethod string `json:"paymentMethod" binding:"required,oneof=cod stripe"`
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
	IsDefault bool    `json:"isDefault,omitempty" binding:"omitempty"`
}

type UpdateAddressRequest struct {
	Phone     *string `json:"phone" binding:"omitempty"`
	Address   *string `json:"address1" binding:"omitempty"`
	Ward      *string `json:"ward" binding:"omitempty"`
	District  *string `json:"district" binding:"omitempty"`
	City      *string `json:"city" binding:"omitempty"`
	IsDefault *bool   `json:"isDefault" binding:"omitempty"`
}

type UpdateCartItemQtyRequest struct {
	Quantity int16 `json:"quantity" binding:"required,gt=0"`
}

type CheckoutRequest struct {
	PaymentMethod      string   `json:"paymentMethod" binding:"required,oneof=code stripe"`
	PaymentGateway     *string  `json:"paymentGateway" binding:"omitempty,oneof=stripe"`
	AddressID          *string  `json:"addressId" binding:"omitempty,uuid"`
	Email              *string  `json:"email" binding:"omitempty,email"`
	FullName           *string  `json:"fullName" binding:"omitempty"`
	Address            *Address `json:"address" binding:"omitempty"`
	PaymentRecipeEmail *string  `json:"paymentReceiptEmail" binding:"omitempty,email"`
}

type AssignmentRequest struct {
	VariantIDs []string `json:"variantIds" binding:"required"`
}

type UpdateUserRequest struct {
	UserID   uuid.UUID `json:"userId" binding:"required,uuid"`
	FullName *string   `json:"fullname,omitempty" binding:"omitempty,min=3,max=32"`
	Email    *string   `json:"email" binding:"email,max=255,min=6"`
	Phone    *string   `json:"phone" binding:"omitempty,min=8,max=15"`
}

type CollectionProductRequest struct {
	SortOrder int16 `json:"sortOrder,omitempty"`
}

type UpdateCategoryRequest struct {
	Name         *string               `form:"name" binding:"omitempty,min=3,max=255"`
	Description  *string               `form:"description" binding:"omitempty,max=1000"`
	Slug         *string               `form:"slug" binding:"omitempty,min=3,max=255"`
	Published    *bool                 `form:"published" binding:"omitempty"`
	Remarkable   *bool                 `form:"remarkable" binding:"omitempty"`
	DisplayOrder *int16                `form:"displayOrder" binding:"omitempty"`
	Image        *multipart.FileHeader `form:"image" binding:"omitempty"`
}

type ProductQueries struct {
	Page         int64    `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize     int64    `form:"pageSize,default=20" binding:"omitempty,min=1,max=100"`
	Search       *string  `form:"search" binding:"omitempty,max=1000"`
	CategoryIDs  []string `form:"categoryIds" binding:"omitempty,uuidslice"`
	BrandID      *string  `form:"brandId" binding:"omitempty,uuid4"`
	CollectionID *string  `form:"collectionId" binding:"omitempty,uuid4"`
}

type OrderListQuery struct {
	PaginationQueryParams
	Status        *string `form:"status,omitempty" binding:"omitempty,oneof=pending confirmed delivering delivered completed cancelled refunded"`
	PaymentStatus *string `form:"paymentStatus,omitempty" binding:"omitempty,oneof=pending succeeded failed cancelled refunded"`
}

type PostRatingFormData struct {
	OrderItemID string                  `form:"orderItemId" binding:"required"`
	Rating      float64                 `form:"rating" binding:"required,min=1,max=5"`
	Title       string                  `form:"title" binding:"required"`
	Content     string                  `form:"content" binding:"required"`
	Files       []*multipart.FileHeader `form:"files" binding:"omitempty"`
}

type CreateDiscountRequest struct {
	Code              string    `json:"code" binding:"required,min=5,max=32,alphanum"`
	DiscountType      string    `json:"discountType" binding:"required,oneof=percentage fixed_amount"`
	DiscountValue     float64   `json:"discountValue" binding:"required,gt=0"`
	IsActive          bool      `json:"isActive" binding:"required"`
	StartsAt          time.Time `json:"startsAt" binding:"required" time_format:"2006-01-02T15:04:05+07:00"`
	ExpiresAt         time.Time `json:"expiresAt" binding:"omitempty" time_format:"2006-01-02T15:04:05+07:00"`
	Description       *string   `json:"description" binding:"omitempty,max=1000"`
	MinPurchaseAmount *float64  `json:"minPurchaseAmount" binding:"omitempty,gt=0"`
	MaxDiscountAmount *float64  `json:"maxDiscountAmount" binding:"omitempty,gt=0"`
	UsageLimit        *int32    `json:"usageLimit" binding:"omitempty,gte=0"`

	// Related entities
	Products   []string `json:"products,omitempty" binding:"omitempty,uuidslice"`
	Categories []string `json:"categories,omitempty" binding:"omitempty,uuidslice"`
	Users      []string `json:"users,omitempty" binding:"omitempty,uuidslice"`
}

type DiscountListQuery struct {
	PaginationQueryParams
	DiscountType  *string    `form:"discountType" binding:"omitempty,oneof=percentage fixed_amount"`
	IsActive      *bool      `from:"isActive" binding:"omitempty"`
	DiscountValue *float64   `form:"discountValue" binding:"omitempty,gt=0"`
	FromDate      *time.Time `form:"fromDate" binding:"omitempty"`
	ToDate        *time.Time `form:"toDate" binding:"omitempty"`
}

type RegisterRequestBody struct {
	Username string                `json:"username" binding:"required,min=3,max=32,lowercase"`
	Password string                `json:"password" binding:"required,min=6,max=32"`
	FullName string                `json:"fullname" binding:"required,min=3,max=32"`
	Phone    string                `json:"phone" binding:"required,min=10,max=15"`
	Email    string                `json:"email" binding:"required,email,max=255,min=6"`
	Address  *CreateAddressRequest `json:"address" binding:"omitempty,required"`
}

type LoginRequest struct {
	Username *string `form:"username" binding:"omitempty,min=3,max=32"`
	Email    *string `form:"email" binding:"omitempty,email,max=255,min=6"`
	Password string  `json:"password" binding:"required,min=6,max=32"`
}

type BrandsQueries struct {
	PaginationQueryParams
	ids []int32 `form:"ids,omitempty"`
}

type BrandProductRequest struct {
	SortOrder int16 `json:"sortOrder,omitempty"`
}

type AttributeValueRequest struct {
	Code         string  `json:"code" binding:"required"`
	Name         *string `json:"name" binding:"omitempty"`
	DisplayOrder *int16  `json:"displayOrder" binding:"omitempty,min=0"`
	IsActive     *bool   `json:"isActive" binding:"omitempty"`
}

type UpdateAttributeValueRequest struct {
	ID                    *string `json:"id" binding:"omitempty,uuid"`
	AttributeValueRequest `json:",inline"`
}

type CreateAttributeRequest struct {
	Name   string                  `json:"name" binding:"required"`
	Values []AttributeValueRequest `json:"values,omitempty"`
}

type UpdateAttributeRequest struct {
	Name   string                        `json:"name" binding:"required"`
	Values []UpdateAttributeValueRequest `json:"values,omitempty"`
}

type AttributeParam struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type GetAttributesQuery struct {
	IDs []uuid.UUID `form:"ids" binding:"omitempty"`
}
