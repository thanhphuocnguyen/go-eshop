package dto

import (
	"encoding/json"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

func MapToUserResponse(user repository.User, roleCode string) UserDetail {
	return UserDetail{
		ID:                user.ID,
		Addresses:         []AddressDetail{},
		Email:             user.Email,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		RoleID:            user.RoleID.String(),
		RoleCode:          roleCode,
		Phone:             user.PhoneNumber,
		Username:          user.Username,
		VerifiedEmail:     user.VerifiedEmail,
		VerifiedPhone:     user.VerifiedPhone,
		CreatedAt:         user.CreatedAt.String(),
		UpdatedAt:         user.UpdatedAt.String(),
		PasswordChangedAt: user.PasswordChangedAt.String(),
	}
}

func MapAddressToAddressResponse(address repository.UserAddress) AddressDetail {
	return AddressDetail{
		Phone:    address.PhoneNumber,
		Street:   address.Street,
		Ward:     address.Ward,
		District: address.District,
		City:     address.City,
		Default:  address.IsDefault,
		ID:       address.ID.String(),
	}
}

func MapToProductDetailResponse(row repository.GetProductDetailRow) ProductDetail {
	basePrice, _ := row.BasePrice.Float64Value()

	resp := ProductDetail{
		ID:                 row.ID.String(),
		Name:               row.Name,
		BasePrice:          basePrice.Float64,
		ShortDescription:   row.ShortDescription,
		Description:        row.Description,
		BaseSku:            row.BaseSku,
		Slug:               row.Slug,
		RatingCount:        row.RatingCount,
		OneStarCount:       row.OneStarCount,
		TwoStarCount:       row.TwoStarCount,
		ThreeStarCount:     row.ThreeStarCount,
		FourStarCount:      row.FourStarCount,
		FiveStarCount:      row.FiveStarCount,
		DiscountPercentage: row.DiscountPercentage,
		UpdatedAt:          row.UpdatedAt.String(),
		CreatedAt:          row.CreatedAt.String(),

		IsActive: *row.IsActive,
		ImageUrl: row.ImageUrl,
		ImageId:  row.ImageID,

		// Initialize slices
		Categories:  []GeneralCategory{},
		Collections: []GeneralCategory{},
		Attributes:  []ProductAttribute{},
		Brand:       GeneralCategory{},
		Variations:  []VariantDetail{},
	}

	// Unmarshal JSON data
	if err := json.Unmarshal(row.Attributes, &resp.Attributes); err != nil {
		log.Error().Err(err).Msg("Unmarshal attributes")
	}
	if err := json.Unmarshal(row.Categories, &resp.Categories); err != nil {
		log.Error().Err(err).Msg("Unmarshal categories")
	}
	if err := json.Unmarshal(row.Collections, &resp.Collections); err != nil {
		log.Error().Err(err).Msg("Unmarshal collections")
	}
	if err := json.Unmarshal(row.Brand, &resp.Brand); err != nil {
		log.Error().Err(err).Msg("Unmarshal brand")
	}
	if err := json.Unmarshal(row.Variants, &resp.Variations); err != nil {
		log.Error().Err(err).Msg("Unmarshal variants")
	}

	return resp
}

func MapToAdminProductResponse(productRow repository.Product) ProductListItem {
	basePrice, _ := productRow.BasePrice.Float64Value()

	avgRating := utils.GetAvgRating(
		productRow.RatingCount,
		productRow.OneStarCount,
		productRow.TwoStarCount,
		productRow.ThreeStarCount,
		productRow.FourStarCount,
		productRow.FiveStarCount,
	)

	product := ProductListItem{
		ID:                 productRow.ID.String(),
		Name:               productRow.Name,
		Description:        productRow.Description,
		BasePrice:          basePrice.Float64,
		Sku:                productRow.BaseSku,
		Slug:               productRow.Slug,
		AvgRating:          &avgRating,
		ImageUrl:           productRow.ImageUrl,
		ImgID:              productRow.ImageID,
		ReviewCount:        &productRow.RatingCount,
		CreatedAt:          productRow.CreatedAt.String(),
		UpdatedAt:          productRow.UpdatedAt.String(),
		DiscountPercentage: productRow.DiscountPercentage,
	}

	return product
}

func MapToShopProductResponse(productRow repository.GetProductListRow) ProductSummary {
	price, _ := productRow.MinPrice.Float64Value()
	avgRating := utils.GetAvgRating(
		productRow.RatingCount,
		productRow.OneStarCount,
		productRow.TwoStarCount,
		productRow.ThreeStarCount,
		productRow.FourStarCount,
		productRow.FiveStarCount,
	)

	product := ProductSummary{
		ID:           productRow.ID.String(),
		Name:         productRow.Name,
		Price:        price.Float64,
		VariantCount: int16(productRow.VariantCount),
		Slug:         productRow.Slug,
		AvgRating:    &avgRating,
		ImageUrl:     productRow.ImageUrl,
		ImageID:      productRow.ImageID,
		ReviewCount:  &productRow.RatingCount,
		CreatedAt:    productRow.CreatedAt.String(),
		UpdatedAt:    productRow.UpdatedAt.String(),
	}

	return product
}

func MapToVariantListModelDto(row repository.GetProductVariantListRow) VariantDetail {
	price, _ := row.Price.Float64Value()
	variant := VariantDetail{
		ID:       row.ID.String(),
		Price:    price.Float64,
		Stock:    row.Stock,
		IsActive: *row.IsActive,
		Sku:      row.Sku,
		ImageUrl: row.ImageUrl,
	}
	variant.Attributes = []AttributeValueDetail{}
	err := json.Unmarshal(row.AttributeValues, &variant.Attributes)
	if err != nil {
		log.Error().Err(err).Msg("Unmarshal variant attribute values")
	}
	if row.Weight.Valid {
		weight, _ := row.Weight.Float64Value()
		variant.Weight = &weight.Float64
	}

	return variant
}

func MapAddressResponse(address repository.UserAddress) AddressDetail {
	return AddressDetail{
		ID:        address.ID.String(),
		Default:   address.IsDefault,
		CreatedAt: address.CreatedAt,
		Phone:     address.PhoneNumber,
		Street:    address.Street,
		Ward:      address.Ward,
		District:  address.District,
		City:      address.City,
	}
}

type ErrorResp struct {
	Error ApiError `json:"error"`
}

func CreateErr(code string, err error) ErrorResp {
	return ErrorResp{
		Error: ApiError{
			Code:    code,
			Details: err.Error(),
			Stack:   err,
		},
	}
}

func CreateDataResp[T any](c *gin.Context, data T, pagination *Pagination, err *ApiError) ApiResponse[T] {
	resp := ApiResponse[T]{
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

func CreatePagination(page, pageSize, total int64) *Pagination {
	return &Pagination{
		Page:            page,
		PageSize:        pageSize,
		Total:           total,
		TotalPages:      total / int64(pageSize),
		HasNextPage:     total > int64(page*pageSize),
		HasPreviousPage: page > 1,
	}
}

func IsStructEmpty(s interface{}) bool {
	return unsafe.Sizeof(s) == 0
}
