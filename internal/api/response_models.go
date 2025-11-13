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
	Page     int64   `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize int64   `form:"pageSize,default=20" binding:"omitempty,min=1,max=100"`
	Search   *string `form:"search" binding:"omitempty,omitzero,max=1000"`
}

type ProductDetailItemResponse struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	ShortDescription *string  `json:"shortDescription"`
	Attributes       []string `json:"attributes"`
	BasePrice        float64  `json:"price,omitzero"`
	BaseSku          string   `json:"sku"`
	IsActive         bool     `json:"isActive"`
	Slug             string   `json:"slug"`

	RatingCount    int32 `json:"ratingCount"`
	OneStarCount   int32 `json:"oneStarCount"`
	TwoStarCount   int32 `json:"twoStarCount"`
	ThreeStarCount int32 `json:"threeStarCount"`
	FourStarCount  int32 `json:"fourStarCount"`
	FiveStarCount  int32 `json:"fiveStarCount"`

	UpdatedAt string `json:"updatedAt"`
	CreatedAt string `json:"createdAt"`

	MaxDiscountValue float64 `json:"maxDiscountValue,omitzero"`
	DiscountType     *string `json:"discountType,omitempty"`

	Variants      []ProductVariantModel    `json:"variants"`
	ProductImages []ProductImageModel      `json:"productImages"`
	Collection    *GeneralCategoryResponse `json:"collection,omitempty"`
	Brand         *GeneralCategoryResponse `json:"brand,omitempty"`
	Category      *GeneralCategoryResponse `json:"category,omitempty"`
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
	CreatedAt   string             `json:"createdAt,omitempty"`
	UpdatedAt   string             `json:"updatedAt,omitempty"`
	ImageUrl    *string            `json:"imageUrl,omitempty"`
	Products    []ProductListModel `json:"products,omitempty"`
}

type CategoryResponse struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	Description *string                 `json:"description,omitempty"`
	Slug        string                  `json:"slug"`
	Published   bool                    `json:"published,omitempty"`
	Remarkable  bool                    `json:"remarkable,omitempty"`
	CreatedAt   string                  `json:"createdAt,omitempty"`
	UpdatedAt   string                  `json:"updatedAt,omitempty"`
	ImageUrl    *string                 `json:"imageUrl,omitempty"`
	Products    []CategoryLinkedProduct `json:"products"`
}

type PublicIDParam struct {
	PublicID string `uri:"public_id" binding:"required"`
}

type RemoveImageParams struct {
	ID        string  `uri:"id" binding:"required,uuid"`
	VariantID *string `uri:"id" binding:"omitempty,uuid"`
	ImageID   string  `uri:"imageId" binding:"required,uuid"`
}

type ImageResponse struct {
	ID          string   `json:"id"`
	ExternalID  string   `json:"externalId"`
	Url         string   `json:"url"`
	MimeType    string   `json:"mimeType,omitempty"`
	FileSize    int64    `json:"fileSize,omitzero"`
	Assignments []string `json:"assignments,omitempty"`
}

type CollectionsQueryParams struct {
	PaginationQueryParams
	Collections *[]int32 `form:"collectionIds,omitempty"`
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
	PaymentID    string  `json:"paymentId"`
	ClientSecret *string `json:"clientSecret"`
}

type PaymentResponse struct {
	ID      string                   `json:"id"`
	Gateway *string                  `json:"gateway,omitempty"`
	Status  repository.PaymentStatus `json:"status,omitempty"`
	Details interface{}              `json:"details"`
}

type OrderItemResponse struct {
	ID                 string                             `json:"id"`
	VariantID          string                             `json:"variantId"`
	Name               string                             `json:"name"`
	ImageUrl           *string                            `json:"imageUrl"`
	LineTotal          float64                            `json:"lineTotal"`
	Quantity           int16                              `json:"quantity"`
	AttributesSnapshot []repository.AttributeDataSnapshot `json:"attributesSnapshot"`
	Rating             *RatingModel                       `json:"rating,omitempty"`
}

type OrderDetailResponse struct {
	ID            uuid.UUID                          `json:"id"`
	Total         float64                            `json:"total"`
	Status        repository.OrderStatus             `json:"status"`
	CustomerName  string                             `json:"customerName"`
	CustomerEmail string                             `json:"customerEmail"`
	PaymentInfo   PaymentInfoModel                   `json:"paymentInfo,omitempty"`
	ShippingInfo  repository.ShippingAddressSnapshot `json:"shippingInfo"`
	Products      []OrderItemResponse                `json:"products"`
	CreatedAt     time.Time                          `json:"createdAt"`
}

type OrderListResponse struct {
	ID            uuid.UUID                `json:"id"`
	Total         float64                  `json:"total"`
	TotalItems    int32                    `json:"totalItems"`
	Status        repository.OrderStatus   `json:"status"`
	PaymentStatus repository.PaymentStatus `json:"paymentStatus"`
	CustomerName  string                   `json:"customerName"`
	CustomerEmail string                   `json:"customerEmail"`
	CreatedAt     time.Time                `json:"createdAt"`
	UpdatedAt     time.Time                `json:"updatedAt"`
}

type AddressResponse struct {
	ID        string    `json:"id"`
	Default   bool      `json:"default"`
	CreatedAt time.Time `json:"createdAt"`
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
	VerifiedEmail     bool                `json:"verifiedEmail,omitempty"`
	VerifiedPhone     bool                `json:"verifiedPhone,omitempty"`
	PasswordChangedAt string              `json:"passwordChangedAt,omitempty"`
	Addresses         []AddressResponse   `json:"addresses"`
	CreatedAt         string              `json:"createdAt,omitempty"`
	UpdatedAt         string              `json:"updatedAt,omitempty"`
}

type CartItemResponse struct {
	ID         string                             `json:"id" binding:"required,uuid"`
	ProductID  string                             `json:"productId" binding:"required,uuid"`
	VariantID  string                             `json:"variantId" binding:"required,uuid"`
	Name       string                             `json:"name"`
	Quantity   int16                              `json:"quantity"`
	Price      float64                            `json:"price"`
	StockQty   int32                              `json:"stock"`
	Sku        *string                            `json:"sku,omitempty"`
	ImageURL   *string                            `json:"imageUrl,omitempty"`
	CategoryID *string                            `json:"categoryId,omitempty"`
	Attributes []repository.AttributeDataSnapshot `json:"attributes"`
}

type CartDetailResponse struct {
	ID         uuid.UUID          `json:"id"`
	TotalPrice float64            `json:"totalPrice"`
	CartItems  []CartItemResponse `json:"cartItems"`
	UpdatedAt  time.Time          `json:"updatedAt,omitempty"`
	CreatedAt  time.Time          `json:"createdAt"`
}

type DiscountListItemResponseModel struct {
	ID            string  `json:"id"`
	Code          string  `json:"code"`
	Description   *string `json:"description"`
	DiscountType  string  `json:"discountType"`
	DiscountValue float64 `json:"discountValue"`
	ExpiresAt     string  `json:"expiresAt"`
	StartsAt      string  `json:"startsAt"`
	UsedCount     int32   `json:"usedCount"`
	UsageLimit    *int32  `json:"usageLimit"`
	IsActive      bool    `json:"isActive"`
	CreatedAt     string  `json:"createdAt"`
	UpdatedAt     string  `json:"updatedAt"`
	MinPurchase   float64 `json:"minPurchase"`
	MaxDiscount   float64 `json:"maxDiscount"`
}

type DiscountLinkObject struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Price *float64 `json:"price,omitempty"`
}

type DiscountUsageHistory struct {
	ID             string    `json:"id"`
	OrderID        string    `json:"orderId"`
	CustomerName   string    `json:"customerName"`
	Amount         float64   `json:"amount"`
	DiscountAmount float64   `json:"discountAmount"`
	Date           time.Time `json:"date" format:"2006-01-02"`
}

type DiscountDetailResponseModel struct {
	ID            string                 `json:"id"`
	Code          string                 `json:"code"`
	Description   *string                `json:"description"`
	DiscountType  string                 `json:"discountType"`
	DiscountValue float64                `json:"discountValue"`
	ExpiresAt     string                 `json:"expiresAt"`
	StartsAt      string                 `json:"startsAt"`
	UsedCount     int32                  `json:"usedCount"`
	UsageLimit    *int32                 `json:"usageLimit"`
	IsActive      bool                   `json:"isActive"`
	CreatedAt     string                 `json:"createdAt"`
	UpdatedAt     string                 `json:"updatedAt"`
	MinPurchase   float64                `json:"minPurchase"`
	MaxDiscount   float64                `json:"maxDiscount"`
	UsageHistory  []DiscountUsageHistory `json:"usageHistory"`
}

type LoginResponse struct {
	ID                    uuid.UUID `json:"sessionId"`
	AccessToken           string    `json:"accessToken"`
	AccessTokenExpiresAt  time.Time `json:"accessTokenExpiresIn"`
	RefreshToken          string    `json:"refreshToken"`
	RefreshTokenExpiresAt time.Time `json:"refreshTokenExpiresAt"`
}

type RefreshTokenResponse struct {
	AccessToken          string    `json:"accessToken"`
	AccessTokenExpiresAt time.Time `json:"accessTokenExpiresAt"`
}

type ProductRatingModel struct {
	ID               uuid.UUID          `json:"id"`
	FirstName        string             `json:"firstName"`
	LastName         string             `json:"lastName"`
	ProductName      string             `json:"productName,omitempty"`
	UserID           uuid.UUID          `json:"userId"`
	Rating           float64            `json:"rating"`
	ReviewTitle      string             `json:"reviewTitle"`
	IsVisible        bool               `json:"isVisible"`
	IsApproved       bool               `json:"isApproved"`
	ReviewContent    string             `json:"reviewContent"`
	VerifiedPurchase bool               `json:"verifiedPurchase"`
	Count            int64              `json:"count"`
	Images           []RatingImageModel `json:"images"`
}

type ProductListModel struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	VariantCount int64    `json:"variantCount,omitzero"`
	BasePrice    float64  `json:"basePrice,omitzero"`
	MinPrice     float64  `json:"minPrice,omitzero"`
	MaxPrice     float64  `json:"maxPrice,omitzero"`
	Slug         string   `json:"slug,omitempty"`
	Sku          string   `json:"sku"`
	ImgUrl       *string  `json:"imageUrl,omitempty"`
	AvgRating    *float64 `json:"avgRating,omitempty"`
	ReviewCount  *int32   `json:"reviewCount,omitempty"`
	ImgID        *string  `json:"imageId,omitempty"`
	CreatedAt    string   `json:"createdAt,omitempty"`
	UpdatedAt    string   `json:"updatedAt,omitempty"`
}

type FiltersModel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProductVariantImageModel struct {
	ID        string `json:"id"`
	VariantID string `json:"variantId,omitempty"`
	ImageUrl  string `json:"imageUrl"`
	ImageID   string `json:"imageId"`
}

type CategoryLinkedProduct struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	VariantCount int32   `json:"variantCount"`
	ImageUrl     *string `json:"imageUrl,omitempty"`
}

type ProductVariantModel struct {
	ID         string                  `json:"id"`
	Price      float64                 `json:"price"`
	StockQty   int32                   `json:"stockQty"`
	IsActive   bool                    `json:"isActive"`
	Sku        *string                 `json:"sku,omitempty"`
	Attributes []ProductAttributeModel `json:"attributes"`
}

type ImageAssignmentModel struct {
	ID           string `json:"id"`
	EntityID     string `json:"entityId"`
	EntityType   string `json:"entityType"`
	DisplayOrder int16  `json:"displayOrder"`
	Role         string `json:"role"`
}

type ProductImageModel struct {
	ID                 string                 `json:"id"`
	Url                string                 `json:"url"`
	ExternalID         string                 `json:"externalId"`
	Role               string                 `json:"role"`
	VariantAssignments []ImageAssignmentModel `json:"assignments"`
}

type OrderItemAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type PaymentInfoModel struct {
	ID           string  `json:"id"`
	RefundID     *string `json:"refundId"`
	Amount       float64 `json:"amount"`
	IntendID     *string `json:"intentId"`
	ClientSecret *string `json:"clientSecret"`
	GateWay      *string `json:"gateway"`
	Method       string  `json:"method"`
	Status       string  `json:"status"`
}

type RatingModel struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Rating    float64   `json:"rating"`
	CreatedAt time.Time `json:"createdAt"`
}

type RatingImageModel struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type ProductAttributeModel struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	ValueObject AttributeValue `json:"valueObject"`
}

type AttributeValue struct {
	ID           uuid.UUID `json:"id"`
	Code         string    `json:"code"`
	Name         *string   `json:"name"`
	IsActive     *bool     `json:"isActive"`
	DisplayOrder *int16    `json:"displayOrder"`
}

type AttributeResponse struct {
	ID        uuid.UUID        `json:"id"`
	Name      string           `json:"name"`
	Values    []AttributeValue `json:"values,omitempty"`
	CreatedAt string           `json:"createdAt"`
	UpdatedAt string           `json:"updatedAt"`
}
