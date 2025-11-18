package api

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
)

type PaginationQueryParams struct {
	Page     int64   `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize int64   `form:"pageSize,default=20" binding:"omitempty,min=1,max=100"`
	Search   *string `form:"search" binding:"omitempty,omitzero,max=1000"`
}

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
	OrderID         string `json:"orderId" binding:"required,uuid"`
	PaymentMethodID string `json:"paymentMethodId" binding:"required,uuid"`
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
	Reason string `json:"reason" binding:"required"`
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

type CheckoutAddress struct {
	Phone    string  `json:"phone"`
	Street   string  `json:"street"`
	Ward     *string `json:"ward,omitempty"`
	District string  `json:"district"`
	City     string  `json:"city"`
}

type CheckoutRequest struct {
	PaymentMethodId string           `json:"paymentMethod" binding:"required,oneof=cod stripe"`
	PaymentGateway  *string          `json:"paymentGateway" binding:"omitempty,oneof=stripe"`
	AddressID       *string          `json:"addressId" binding:"omitempty,uuid"`
	Email           *string          `json:"email" binding:"omitempty,email"`
	FirstName       *string          `json:"firstName,omitempty" binding:"omitempty,min=3,max=32"`
	LastName        *string          `json:"lastName,omitempty" binding:"omitempty,min=3,max=32"`
	Address         *CheckoutAddress `json:"address" binding:"omitempty"`
	DiscountCode    *string          `json:"discountCode" binding:"omitempty,min=5,max=32,alphanum"`
}

type AssignmentRequest struct {
	VariantIDs []string `json:"variantIds" binding:"required"`
}

type UpdateUserRequest struct {
	UserID    uuid.UUID `json:"userId" binding:"required,uuid"`
	FirstName *string   `json:"firstName,omitempty" binding:"omitempty,min=3,max=32"`
	LastName  *string   `json:"lastName,omitempty" binding:"omitempty,min=3,max=32"`
	Email     *string   `json:"email" binding:"email,max=255,min=6"`
	Phone     *string   `json:"phone" binding:"omitempty,min=8,max=15"`
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
	StartsAt          time.Time `json:"startsAt" binding:"required" time_format:"RFC3339"`
	ExpiresAt         time.Time `json:"expiresAt" binding:"omitempty" time_format:"RFC3339"`
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
	FromDate      *time.Time `form:"fromDate, default=" binding:"omitempty"`
	ToDate        *time.Time `form:"toDate" binding:"omitempty"`
}

type RegisterRequestBody struct {
	Username  string                `json:"username" binding:"required,min=3,max=32,lowercase"`
	Password  string                `json:"password" binding:"required,min=6,max=32"`
	FirstName string                `json:"firstName,omitempty" binding:"omitempty,min=3,max=32"`
	LastName  string                `json:"lastName,omitempty" binding:"omitempty,min=3,max=32"`
	Phone     string                `json:"phone" binding:"required,min=10,max=15"`
	Email     string                `json:"email" binding:"required,email,max=255,min=6"`
	Address   *CreateAddressRequest `json:"address" binding:"omitempty,required"`
}

type LoginRequest struct {
	Username *string `form:"username" binding:"omitempty,max=32"`
	Email    *string `form:"email" binding:"omitempty,email,max=255"`
	Password string  `form:"password" binding:"required,min=6,max=32"`
}

type BrandsQueries struct {
	PaginationQueryParams
}

type BrandProductRequest struct {
	SortOrder int16 `json:"sortOrder,omitempty"`
}

type AttributeRequest struct {
	Name string `json:"name" binding:"required"`
}

type AttributeValuesReq struct {
	Value string `json:"value" binding:"required"`
}

type AttributeParam struct {
	ID      int32  `uri:"id" binding:"required"`
	ValueID *int64 `uri:"valueId" binding:"omitempty"`
}

type GetAttributesQuery struct {
	IDs []int32 `form:"ids" binding:"omitempty"`
}

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
	VerifyCode string `form:"verifyCode" binding:"required,min=1"`
}

type ProductVariantParam struct {
	ID string `uri:"variantId" binding:"required,uuid"`
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
	Stack   error  `json:"stack,omitempty"` // Hide in production
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

type ManageProductDetailResp struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	ShortDescription *string `json:"shortDescription"`
	BasePrice        float64 `json:"price,omitzero"`
	BaseSku          string  `json:"sku"`
	IsActive         bool    `json:"isActive"`
	Slug             string  `json:"slug"`
	ImageUrl         string  `json:"imageUrl"`
	ImageId          string  `json:"imageId,omitempty"`

	RatingCount    int32 `json:"ratingCount"`
	OneStarCount   int32 `json:"oneStarCount"`
	TwoStarCount   int32 `json:"twoStarCount"`
	ThreeStarCount int32 `json:"threeStarCount"`
	FourStarCount  int32 `json:"fourStarCount"`
	FiveStarCount  int32 `json:"fiveStarCount"`

	UpdatedAt string `json:"updatedAt"`
	CreatedAt string `json:"createdAt"`

	Attributes []int32                  `json:"attributes,omitempty"`
	Collection *GeneralCategoryResponse `json:"collection,omitempty"`
	Brand      *GeneralCategoryResponse `json:"brand,omitempty"`
	Category   *GeneralCategoryResponse `json:"category,omitempty"`
}

type GeneralCategoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CategoryListResponse struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	Description *string                  `json:"description,omitempty"`
	Slug        string                   `json:"slug"`
	Published   bool                     `json:"published,omitempty"`
	Remarkable  bool                     `json:"remarkable,omitempty"`
	CreatedAt   string                   `json:"createdAt,omitempty"`
	UpdatedAt   string                   `json:"updatedAt,omitempty"`
	ImageUrl    *string                  `json:"imageUrl,omitempty"`
	Products    []ManageProductListModel `json:"products,omitempty"`
}

type CategoryResponse struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	Description *string                 `json:"description,omitempty"`
	Slug        string                  `json:"slug"`
	Published   bool                    `json:"published,omitempty"`
	CreatedAt   string                  `json:"createdAt,omitempty"`
	UpdatedAt   string                  `json:"updatedAt,omitempty"`
	ImageUrl    *string                 `json:"imageUrl,omitempty"`
	Products    []CategoryLinkedProduct `json:"products"`
}

type PublicIDParam struct {
	PublicID string `uri:"publicId" binding:"required"`
}

type RemoveImageParams struct {
	ID        string  `uri:"id" binding:"required,uuid"`
	VariantID *string `uri:"variantId" binding:"omitempty,uuid"`
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
	Collection CategoryResponse              `json:"collection"`
	Categories []FiltersModel                `json:"categories"`
	Brands     []FiltersModel                `json:"brands"`
	Attributes map[string][]FilterIntIDModel `json:"attributes"`
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
	// Details contains gateway-specific payment information.
	// The type may vary depending on the payment provider, and is typically a struct or map[string]interface{}.
	Details any `json:"details"`
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
	PaymentInfo   PaymentInfoModel                   `json:"paymentInfo"`
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

type UserDetail struct {
	ID                uuid.UUID         `json:"id"`
	RoleID            string            `json:"roleId"`
	RoleCode          repository.Role   `json:"roleCode"`
	Username          string            `json:"username"`
	FirstName         string            `json:"firstName"`
	LastName          string            `json:"lastName"`
	Email             string            `json:"email,omitempty"`
	Phone             string            `json:"phone,omitempty"`
	AvatarURL         *string           `json:"avatarUrl,omitempty"`
	AvatarID          *string           `json:"avatarId,omitempty"`
	Locked            bool              `json:"locked,omitempty"`
	VerifiedEmail     bool              `json:"verifiedEmail,omitempty"`
	VerifiedPhone     bool              `json:"verifiedPhone,omitempty"`
	PasswordChangedAt string            `json:"passwordChangedAt,omitempty"`
	Addresses         []AddressResponse `json:"addresses"`
	CreatedAt         string            `json:"createdAt,omitempty"`
	LastLoginAt       string            `json:"lastLoginAt,omitempty"`
	UpdatedAt         string            `json:"updatedAt,omitempty"`
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
	UpdatedAt  time.Time          `json:"updatedAt"`
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

type ManageProductListModel struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	VariantCount int64    `json:"variantCount,omitzero"`
	BasePrice    float64  `json:"basePrice,omitzero"`
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

type FilterIntIDModel struct {
	ID   int64  `json:"id"`
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

type VariantModel struct {
	ID         string           `json:"id"`
	Price      float64          `json:"price"`
	StockQty   int32            `json:"stockQty"`
	IsActive   bool             `json:"isActive"`
	Weight     *float64         `json:"weight,omitempty"`
	Sku        *string          `json:"sku,omitempty"`
	Attributes []AttributeValue `json:"attributeValues,omitempty"`
	CreatedAt  string           `json:"createdAt"`
	UpdatedAt  string           `json:"updatedAt"`
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
	ID  int64  `json:"id"`
	URL string `json:"url"`
}

type AttributeValue struct {
	ID    int64   `json:"id"`
	Value string  `json:"value"`
	Name  *string `json:"name,omitempty"`
}

type AttributeRespModel struct {
	ID     int32            `json:"id"`
	Name   string           `json:"name"`
	Values []AttributeValue `json:"values"`
}

type CreateProductReq struct {
	Name             string  `json:"name" binding:"required,min=3,max=255"`
	Description      string  `json:"description" binding:"required"`
	ShortDescription *string `json:"shortDescription" binding:"omitempty,max=1000"`
	Attributes       []int32 `json:"attributes" binding:"required,required"`
	BasePrice        float64 `json:"price" binding:"required,gt=0"`
	BaseSku          string  `json:"sku" binding:"required,min=3,max=100"`
	IsActive         bool    `json:"isActive" binding:"required"`
	Slug             string  `json:"slug" binding:"required,min=3,max=255"`

	CollectionID *string `json:"collectionId" binding:"omitempty,uuid4"`
	BrandID      string  `json:"brandId" binding:"omitempty,uuid4"`
	CategoryID   string  `json:"categoryId" binding:"required,uuid4"`
}
type UpdateProductReq struct {
	Name             *string  `json:"name" binding:"min=3,max=255"`
	Description      *string  `json:"description"`
	ShortDescription *string  `json:"shortDescription" binding:"omitempty,max=1000"`
	Attributes       *[]int32 `json:"attributes" binding:""`
	BasePrice        *float64 `json:"price" binding:"gt=0"`
	BaseSku          *string  `json:"sku" binding:"min=3,max=100"`
	IsActive         *bool    `json:"isActive" binding:""`
	Slug             *string  `json:"slug" binding:"min=3,max=255"`

	CollectionID *string `json:"collectionId" binding:"omitempty,uuid4"`
	BrandID      *string `json:"brandId" binding:"omitempty,uuid4"`
	CategoryID   *string `json:"categoryId" binding:"uuid4"`
}

type CreateProdVariantReq struct {
	Price           float64  `json:"price" binding:"required,gt=0"`
	StockQty        int32    `json:"stockQty" binding:"required,gte=0"`
	IsActive        bool     `json:"isActive" binding:"required"`
	AttributeValues []int64  `json:"attributeValues" binding:"required"`
	Description     *string  `json:"description" binding:"omitempty"`
	Weight          *float64 `json:"weight" binding:"omitempty,gt=0"`
}

type UpdateProdVariantReq struct {
	Price           *float64 `json:"price" binding:"omitempty,gt=0"`
	StockQty        *int32   `json:"stockQty" binding:"omitempty,gte=0"`
	IsActive        *bool    `json:"isActive" binding:"omitempty"`
	Description     *string  `json:"description" binding:"omitempty"`
	AttributeValues *[]int64 `json:"attributeValues" binding:"required,omitempty"`
	Weight          *float64 `json:"weight" binding:"omitempty,gt=0"`
}

type URIVariantParam struct {
	ProductID string `uri:"productId" binding:"required,uuid"`
	VariantID string `uri:"variantId" binding:"required,uuid"`
}
