package api

import (
	"time"

	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
)

type URISlugParam struct {
	ID string `uri:"id" binding:"required"`
}

type UriIDParam struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type SlugParam struct {
	Slug string `uri:"slug" binding:"required"`
}

type VerifyEmailQuery struct {
	VerifyCode string `form:"verify_code" binding:"required,min=1"`
}

type ProductVariantParam struct {
	ID string `uri:"variant_id" binding:"required,uuid"`
}

// Meta info about the request
type MetaInfo struct {
	Timestamp string `json:"timestamp"`
	RequestID string `json:"requestId"`
	Path      string `json:"path"`
	Method    string `json:"method"`
}

// Response types - unchanged
type ApiResponse[T any] struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       *T          `json:"data,omitempty,omitzero"`
	Error      *ApiError   `json:"error,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Meta       *MetaInfo   `json:"meta"`
}

// Error structure for detailed errors
type ApiError struct {
	Code    string `json:"code"`
	Details string `json:"details"`
	Stack   string `json:"stack,omitempty"` // Hide in production
}

// Pagination info (for paginated endpoints)
type Pagination struct {
	Total           int64 `json:"total"`
	Page            int64 `json:"page"`
	PageSize        int64 `json:"pageSize"`
	TotalPages      int64 `json:"totalPages"`
	HasNextPage     bool  `json:"hasNextPage"`
	HasPreviousPage bool  `json:"hasPreviousPage"`
}

type PaginationQueryParams struct {
	Page     int64 `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize int64 `form:"page_size,default=20" binding:"omitempty,min=1,max=100"`
}

type ProductListItemResponse struct {
	ID               string                   `json:"id"`
	Name             string                   `json:"name"`
	Description      string                   `json:"description"`
	ShortDescription *string                  `json:"short_description"`
	Attributes       []string                 `json:"attributes"`
	BasePrice        float64                  `json:"price,omitzero"`
	BaseSku          string                   `json:"sku"`
	UpdatedAt        string                   `json:"updated_at"`
	IsActive         bool                     `json:"is_active"`
	Slug             string                   `json:"slug"`
	CreatedAt        string                   `json:"created_at"`
	RatingCount      int32                    `json:"rating_count"`
	OneStarCount     int32                    `json:"one_star_count"`
	TwoStarCount     int32                    `json:"two_star_count"`
	ThreeStarCount   int32                    `json:"three_star_count"`
	FourStarCount    int32                    `json:"four_star_count"`
	FiveStarCount    int32                    `json:"five_star_count"`
	Variants         []ProductVariantModel    `json:"variants"`
	ProductImages    []ProductImageModel      `json:"product_images"`
	Collection       *GeneralCategoryResponse `json:"collection,omitempty"`
	Brand            *GeneralCategoryResponse `json:"brand,omitempty"`
	Category         *GeneralCategoryResponse `json:"category,omitempty"`
}

type GeneralCategoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CategoryListResponse struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description *string            `json:"description,omitempty"`
	Slug        string             `json:"slug"`
	Published   bool               `json:"published,omitempty"`
	Remarkable  bool               `json:"remarkable,omitempty"`
	CreatedAt   string             `json:"created_at,omitempty"`
	UpdatedAt   string             `json:"updated_at,omitempty"`
	ImageUrl    *string            `json:"image_url,omitempty"`
	Products    []ProductListModel `json:"products,omitempty"`
}

type CategoryResponse struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	Description *string                 `json:"description,omitempty"`
	Slug        string                  `json:"slug"`
	Published   bool                    `json:"published,omitempty"`
	Remarkable  bool                    `json:"remarkable,omitempty"`
	CreatedAt   string                  `json:"created_at,omitempty"`
	UpdatedAt   string                  `json:"updated_at,omitempty"`
	ImageUrl    *string                 `json:"image_url,omitempty"`
	Products    []CategoryLinkedProduct `json:"products"`
}

type PublicIDParam struct {
	PublicID string `uri:"public_id" binding:"required"`
}

type RemoveImageParams struct {
	ID        string  `uri:"id" binding:"required,uuid"`
	VariantID *string `uri:"id" binding:"omitempty,uuid"`
	ImageID   string  `uri:"image_id" binding:"required,uuid"`
}

type ImageResponse struct {
	ID          string   `json:"id"`
	ExternalID  string   `json:"external_id"`
	Url         string   `json:"url"`
	MimeType    string   `json:"mime_type,omitempty"`
	FileSize    int64    `json:"file_size,omitzero"`
	Assignments []string `json:"assignments,omitempty"`
}

type CollectionsQueryParams struct {
	PaginationQueryParams
	Collections *[]int32 `form:"collection_ids,omitempty"`
}

type CollectionDetailResponse struct {
	Collection CategoryResponse          `json:"collection"`
	Categories []FiltersModel            `json:"categories"`
	Brands     []FiltersModel            `json:"brands"`
	Attributes map[string][]FiltersModel `json:"attributes"`
}

type RatingsQueryParams struct {
	PaginationQueryParams
	Status *string `form:"status" binding:"omitempty,oneof=approved rejected pending"`
}

type CreatePaymentIntentResponse struct {
	PaymentID    string  `json:"payment_id"`
	ClientSecret *string `json:"client_secret"`
}

type PaymentResponse struct {
	ID      string                    `json:"id"`
	Gateway repository.PaymentGateway `json:"gateway,omitempty"`
	Status  repository.PaymentStatus  `json:"status,omitempty"`
	Details interface{}               `json:"details"`
}

type OrderItemResponse struct {
	ID                 string                             `json:"id"`
	VariantID          string                             `json:"variant_id"`
	Name               string                             `json:"name"`
	ImageUrl           *string                            `json:"image_url"`
	LineTotal          float64                            `json:"line_total"`
	Quantity           int16                              `json:"quantity"`
	AttributesSnapshot []repository.AttributeDataSnapshot `json:"attributes_snapshot"`
	Rating             *RatingModel                       `json:"rating,omitempty"`
}

type OrderDetailResponse struct {
	ID            uuid.UUID                          `json:"id"`
	Total         float64                            `json:"total"`
	Status        repository.OrderStatus             `json:"status"`
	CustomerName  string                             `json:"customer_name"`
	CustomerEmail string                             `json:"customer_email"`
	PaymentInfo   *PaymentInfoModel                  `json:"payment_info,omitempty"`
	ShippingInfo  repository.ShippingAddressSnapshot `json:"shipping_info"`
	Products      []OrderItemResponse                `json:"products"`
	CreatedAt     time.Time                          `json:"created_at"`
}

type OrderListResponse struct {
	ID            uuid.UUID                `json:"id"`
	Total         float64                  `json:"total"`
	TotalItems    int32                    `json:"total_items"`
	Status        repository.OrderStatus   `json:"status"`
	PaymentStatus repository.PaymentStatus `json:"payment_status"`
	CustomerName  string                   `json:"customer_name"`
	CustomerEmail string                   `json:"customer_email"`
	CreatedAt     time.Time                `json:"created_at"`
	UpdatedAt     time.Time                `json:"updated_at"`
}

type AddressResponse struct {
	ID        string    `json:"id"`
	Default   bool      `json:"default"`
	CreatedAt time.Time `json:"created_at"`
	Phone     string    `json:"phone"`
	Street    string    `json:"street"`
	Ward      *string   `json:"ward,omitempty"`
	District  string    `json:"district"`
	City      string    `json:"city"`
}

type UserResponse struct {
	ID                uuid.UUID           `json:"id"`
	Role              repository.UserRole `json:"role"`
	Username          string              `json:"username"`
	FullName          string              `json:"fullname"`
	Email             string              `json:"email,omitempty"`
	Phone             string              `json:"phone,omitempty"`
	VerifiedEmail     bool                `json:"verified_email,omitempty"`
	VerifiedPhone     bool                `json:"verified_phone,omitempty"`
	PasswordChangedAt string              `json:"password_changed_at,omitempty"`
	Addresses         []AddressResponse   `json:"addresses"`
	CreatedAt         string              `json:"created_at,omitempty"`
	UpdatedAt         string              `json:"updated_at,omitempty"`
}

type CartItemResponse struct {
	ID         string                             `json:"id" binding:"required,uuid"`
	ProductID  string                             `json:"product_id" binding:"required,uuid"`
	VariantID  string                             `json:"variant_id" binding:"required,uuid"`
	Name       string                             `json:"name"`
	Quantity   int16                              `json:"quantity"`
	Price      float64                            `json:"price"`
	Discount   int16                              `json:"discount"`
	StockQty   int32                              `json:"stock"`
	Sku        *string                            `json:"sku,omitempty"`
	ImageURL   *string                            `json:"image_url,omitempty"`
	Attributes []repository.AttributeDataSnapshot `json:"attributes"`
}

type CartDetailResponse struct {
	ID         uuid.UUID          `json:"id"`
	TotalPrice float64            `json:"total_price"`
	CartItems  []CartItemResponse `json:"cart_items"`
	UpdatedAt  time.Time          `json:"updated_at,omitempty"`
	CreatedAt  time.Time          `json:"created_at"`
}

type CheckoutResponse struct {
	OrderID         uuid.UUID `json:"order_id"`
	PaymentID       string    `json:"payment_id"`
	PaymentIntentID *string   `json:"payment_intent_id,omitempty"`
	ClientSecret    *string   `json:"client_secret,omitempty"`
	TotalPrice      float64   `json:"total_price"`
}

type DiscountListItemResponseModel struct {
	ID            string  `json:"id"`
	Code          string  `json:"code"`
	Description   string  `json:"description"`
	DiscountType  string  `json:"discount_type"`
	DiscountValue float64 `json:"discount_value"`
	ExpiredAt     string  `json:"expired_at"`
	StartsAt      string  `json:"starts_at"`
	UsedCount     int32   `json:"used_count"`
	UsageLimit    int32   `json:"usage_limit"`
	IsActive      bool    `json:"is_active"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	MinPurchase   int32   `json:"min_purchase"`
	MaxDiscount   int32   `json:"max_discount"`
}

type DiscountLinkObject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DiscountDetailResponseModel struct {
	ID            string               `json:"id"`
	Code          string               `json:"code"`
	Description   string               `json:"description"`
	DiscountType  string               `json:"discount_type"`
	DiscountValue float64              `json:"discount_value"`
	ExpiredAt     string               `json:"expired_at"`
	StartsAt      string               `json:"starts_at"`
	UsedCount     int32                `json:"used_count"`
	UsageLimit    int32                `json:"usage_limit"`
	IsActive      bool                 `json:"is_active"`
	CreatedAt     string               `json:"created_at"`
	UpdatedAt     string               `json:"updated_at"`
	MinPurchase   int32                `json:"min_purchase"`
	MaxDiscount   int32                `json:"max_discount"`
	Products      []DiscountLinkObject `json:"products"`
	Categories    []DiscountLinkObject `json:"categories"`
	Users         []DiscountLinkObject `json:"variants"`
}
