package api

import (
	"mime/multipart"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
)

type URISlugParam struct {
	ID string `uri:"id" binding:"required"`
}

type UriIDParam struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// Response types - unchanged
type ApiResponse[T any] struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       *T          `json:"data,omitempty"`
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

// Meta info about the request
type MetaInfo struct {
	Timestamp string `json:"timestamp"`
	RequestID string `json:"requestId"`
	Path      string `json:"path"`
	Method    string `json:"method"`
}

type PaginationQueryParams struct {
	Page     int64 `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize int64 `form:"page_size,default=20" binding:"omitempty,min=1,max=100"`
}

func createErrorResponse[T any](code string, msg string, err error) ApiResponse[T] {
	return ApiResponse[T]{
		Success: false,
		Data:    nil,
		Error: &ApiError{
			Code:    code,
			Details: msg,
			Stack:   err.Error(),
		},
	}
}

func createSuccessResponse[T any](c *gin.Context, data T, message string, pagination *Pagination, err *ApiError) ApiResponse[T] {
	resp := ApiResponse[T]{
		Success:    true,
		Message:    message,
		Data:       &data,
		Pagination: pagination,

		Meta: &MetaInfo{
			Timestamp: time.Now().Format(time.RFC3339),
			RequestID: c.GetString("RequestID"),
			Path:      c.FullPath(),
			Method:    c.Request.Method,
		},
	}
	if err != nil {
		resp.Error = err
	}
	return resp
}

func createPagination(page, pageSize, total int64) *Pagination {
	return &Pagination{
		Page:            page,
		PageSize:        pageSize,
		Total:           total,
		TotalPages:      total / int64(pageSize),
		HasNextPage:     total > int64(page*pageSize),
		HasPreviousPage: page > 1,
	}
}

type ProductQueries struct {
	Page         int64    `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize     int64    `form:"page_size,default=20" binding:"omitempty,min=1,max=100"`
	Search       *string  `form:"search" binding:"omitempty,max=1000"`
	CategoryIDs  []string `form:"category_ids" binding:"omitempty,uuidslice"`
	BrandID      *string  `form:"brand_id" binding:"omitempty,uuid4"`
	CollectionID *string  `form:"collection_id" binding:"omitempty,uuid4"`
}

type ProductAttributeDetail struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	ValueObject AttributeValue `json:"value_object"`
}

type ImageAssignment struct {
	ID           string `json:"id"`
	EntityID     string `json:"entity_id"`
	EntityType   string `json:"entity_type"`
	DisplayOrder int16  `json:"display_order"`
	Role         string `json:"role"`
}

type ProductImageModel struct {
	ID                 string            `json:"id"`
	Url                string            `json:"url"`
	ExternalID         string            `json:"external_id"`
	Role               string            `json:"role"`
	VariantAssignments []ImageAssignment `json:"assignments"`
}

type ProductVariantModel struct {
	ID         string                   `json:"id"`
	Price      float64                  `json:"price"`
	StockQty   int32                    `json:"stock_qty"`
	IsActive   bool                     `json:"is_active"`
	Sku        *string                  `json:"sku,omitempty"`
	Attributes []ProductAttributeDetail `json:"attributes"`
}

type UpdateProductImageAssignments struct {
	ID           int32  `json:"id"`
	EntityID     string `json:"entity_id"`
	DisplayOrder int16  `json:"display_order"`
}

type GeneralCategoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProductModel struct {
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

type ProductListModel struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	VariantCount int64    `json:"variant_count,omitzero"`
	MinPrice     float64  `json:"min_price,omitzero"`
	MaxPrice     float64  `json:"max_price,omitzero"`
	Slug         string   `json:"slug,omitempty"`
	Sku          string   `json:"sku"`
	ImgUrl       *string  `json:"image_url,omitempty"`
	AvgRating    *float64 `json:"avg_rating,omitempty"`
	ReviewCount  *int32   `json:"review_count,omitempty"`
	ImgID        *string  `json:"image_id,omitempty"`
	CreatedAt    string   `json:"created_at,omitempty"`
	UpdatedAt    string   `json:"updated_at,omitempty"`
}

type ProductCreateResp struct {
	ID string `json:"id"`
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

type CreateCategoryRequest struct {
	Name         string                `form:"name" binding:"required,min=3,max=255"`
	Slug         string                `form:"slug" binding:"required,min=3,max=255"`
	Description  *string               `form:"description" binding:"omitempty,max=1000"`
	DisplayOrder *int16                `form:"display_order" binding:"omitempty"`
	Remarkable   *bool                 `form:"remarkable" binding:"omitempty"`
	Image        *multipart.FileHeader `form:"image" binding:"omitempty"`
}

type CategoryLinkedProduct struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	VariantCount int32   `json:"variant_count"`
	ImageUrl     *string `json:"image_url,omitempty"`
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

type SlugParam struct {
	Slug string `uri:"slug" binding:"required"`
}

type CategoryProductRequest struct {
	SortOrder int16 `json:"sort_order,omitempty"`
}

type PublicIDParam struct {
	PublicID string `uri:"public_id" binding:"required"`
}
type RemoveImageParams struct {
	ID        string  `uri:"id" binding:"required,uuid"`
	VariantID *string `uri:"id" binding:"omitempty,uuid"`
	ImageID   string  `uri:"image_id" binding:"required,uuid"`
}

type AssignmentRequest struct {
	VariantIDs []string `json:"variant_ids" binding:"required"`
}

type ProductVariantImageModel struct {
	ID        string `json:"id"`
	VariantID string `json:"variant_id,omitempty"`
	ImageUrl  string `json:"image_url"`
	ImageID   string `json:"image_id"`
}

type ImageResponse struct {
	ID          string   `json:"id"`
	ExternalID  string   `json:"external_id"`
	Url         string   `json:"url"`
	MimeType    string   `json:"mime_type,omitempty"`
	FileSize    int64    `json:"file_size,omitzero"`
	Assignments []string `json:"assignments,omitempty"`
}

type getCollectionsQueries struct {
	PaginationQueryParams
	Collections *[]int32 `form:"collection_ids,omitempty"`
}

type CollectionProductRequest struct {
	SortOrder int16 `json:"sort_order,omitempty"`
}
type FiltersModel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type CollectionDetailResponse struct {
	Collection CategoryResponse          `json:"collection"`
	Categories []FiltersModel            `json:"categories"`
	Brands     []FiltersModel            `json:"brands"`
	Attributes map[string][]FiltersModel `json:"attributes"`
}

type PostRatingFormData struct {
	OrderItemID string                  `form:"order_item_id" binding:"required"`
	Rating      float64                 `form:"rating" binding:"required,min=1,max=5"`
	Title       string                  `form:"title" binding:"required"`
	Content     string                  `form:"content" binding:"required"`
	Files       []*multipart.FileHeader `form:"files" binding:"omitempty"`
}
type RatingsQueryParams struct {
	PaginationQueryParams
	Status *string `form:"status" binding:"omitempty,oneof=approved rejected pending"`
}
type PostHelpfulRatingRequest struct {
	Helpful bool `json:"helpful"`
}

type PostReplyRatingRequest struct {
	RatingID string `json:"rating_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

type RatingImageModel struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type ProductRatingModel struct {
	ID               uuid.UUID          `json:"id"`
	Name             string             `json:"name"`
	ProductName      string             `json:"product_name,omitempty"`
	UserID           uuid.UUID          `json:"user_id"`
	Rating           float64            `json:"rating"`
	ReviewTitle      string             `json:"review_title"`
	IsVisible        bool               `json:"is_visible"`
	IsApproved       bool               `json:"is_approved"`
	ReviewContent    string             `json:"review_content"`
	VerifiedPurchase bool               `json:"verified_purchase"`
	HelpfulVotes     int32              `json:"helpful_votes"`
	UnhelpfulVotes   int32              `json:"unhelpful_votes"`
	Images           []RatingImageModel `json:"images"`
}

type PaymentRequest struct {
	OrderID       string `json:"order_id" binding:"required,uuid"`
	PaymentMethod string `json:"payment_method" binding:"required,oneof=cod stripe"`
}

type CreatePaymentIntentResult struct {
	PaymentID    string  `json:"payment_id"`
	ClientSecret *string `json:"client_secret"`
}

type PaymentResponse struct {
	ID      string                    `json:"id"`
	Gateway repository.PaymentGateway `json:"gateway,omitempty"`
	Status  repository.PaymentStatus  `json:"status,omitempty"`
	Details interface{}               `json:"details"`
}

type ChangePaymentStatusReq struct {
	Status repository.PaymentStatus `json:"status" binding:"required"`
}

type OrderListQuery struct {
	PaginationQueryParams
	Status        *string `form:"status,omitempty" binding:"omitempty,oneof=pending confirmed delivering delivered completed cancelled refunded"`
	PaymentStatus *string `form:"payment_status,omitempty" binding:"omitempty,oneof=pending succeeded failed cancelled refunded"`
}

type OrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending confirmed delivering delivered completed"`
	Reason string `json:"reason,omitempty"`
}

type Rating struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Rating    float64   `json:"rating"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderItemResponse struct {
	ID                 string                             `json:"id"`
	VariantID          string                             `json:"variant_id"`
	Name               string                             `json:"name"`
	ImageUrl           *string                            `json:"image_url"`
	LineTotal          float64                            `json:"line_total"`
	Quantity           int16                              `json:"quantity"`
	AttributesSnapshot []repository.AttributeDataSnapshot `json:"attributes_snapshot"`
	Rating             *Rating                            `json:"rating,omitempty"`
}
type PaymentInfo struct {
	ID           string  `json:"id"`
	RefundID     *string `json:"refund_id"`
	Amount       float64 `json:"amount"`
	IntendID     *string `json:"intent_id"`
	ClientSecret *string `json:"client_secret"`
	GateWay      *string `json:"gateway"`
	Method       string  `json:"method"`
	Status       string  `json:"status"`
}

type OrderDetailResponse struct {
	ID            uuid.UUID                          `json:"id"`
	Total         float64                            `json:"total"`
	Status        repository.OrderStatus             `json:"status"`
	CustomerName  string                             `json:"customer_name"`
	CustomerEmail string                             `json:"customer_email"`
	PaymentInfo   *PaymentInfo                       `json:"payment_info,omitempty"`
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

type RefundOrderRequest struct {
	Reason string `json:"reason" binding:"required,oneof=defective damaged fraudulent requested_by_customer"`
}
type CancelOrderRequest struct {
	Reason string `json:"reason" binding:"required,oneof=duplicate fraudulent requested_by_customer abandoned"`
}

type CreateAddressReq struct {
	Phone     string  `json:"phone" binding:"required,min=10,max=15"`
	Street    string  `json:"street" binding:"required"`
	District  string  `json:"district" binding:"required"`
	City      string  `json:"city" binding:"required"`
	Ward      *string `json:"ward,omitempty" binding:"omitempty,max=100"`
	IsDefault bool    `json:"is_default,omitempty" binding:"omitempty"`
}

type UpdateAddressReq struct {
	Phone     *string `json:"phone" binding:"omitempty"`
	Address   *string `json:"address_1" binding:"omitempty"`
	Ward      *string `json:"ward" binding:"omitempty"`
	District  *string `json:"district" binding:"omitempty"`
	City      *string `json:"city" binding:"omitempty"`
	IsDefault *bool   `json:"is_default" binding:"omitempty"`
}

// ------------------------------ API Models ------------------------------
type Address struct {
	Phone    string  `json:"phone"`
	Street   string  `json:"street"`
	Ward     *string `json:"ward,omitempty"`
	District string  `json:"district"`
	City     string  `json:"city"`
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

type UpdateUserRequest struct {
	UserID   uuid.UUID `json:"user_id" binding:"required,uuid"`
	FullName *string   `json:"fullname,omitempty" binding:"omitempty,min=3,max=32"`
	Email    *string   `json:"email" binding:"email,max=255,min=6"`
	Phone    *string   `json:"phone" binding:"omitempty,min=8,max=15"`
}

type VerifyEmailQuery struct {
	VerifyCode string `form:"verify_code" binding:"required,min=1"`
}

type ProductVariantParam struct {
	ID string `uri:"variant_id" binding:"required,uuid"`
}

type OrderItemAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
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

type UpdateCartItemQtyReq struct {
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

type CheckoutResponse struct {
	OrderID         uuid.UUID `json:"order_id"`
	PaymentID       string    `json:"payment_id"`
	PaymentIntentID *string   `json:"payment_intent_id,omitempty"`
	ClientSecret    *string   `json:"client_secret,omitempty"`
	TotalPrice      float64   `json:"total_price"`
}
