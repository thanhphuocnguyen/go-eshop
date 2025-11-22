package api

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

func mapToUserResponse(user repository.User, roleCode repository.Role) UserDetail {
	return UserDetail{
		ID:                user.ID,
		Addresses:         []AddressResponse{},
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

func mapAddressToAddressResponse(address repository.UserAddress) AddressResponse {
	return AddressResponse{
		Phone:    address.PhoneNumber,
		Street:   address.Street,
		Ward:     address.Ward,
		District: address.District,
		City:     address.City,
		Default:  address.IsDefault,
		ID:       address.ID.String(),
	}
}

func mapToProductDetailResponse(productRows repository.GetProductDetailRow) ProductDetailDto {
	basePrice, _ := productRows.BasePrice.Float64Value()

	resp := ProductDetailDto{
		ID:               productRows.ID.String(),
		Name:             productRows.Name,
		BasePrice:        basePrice.Float64,
		ShortDescription: productRows.ShortDescription,
		Description:      productRows.Description,
		BaseSku:          productRows.BaseSku,
		Slug:             productRows.Slug,
		RatingCount:      productRows.RatingCount,
		OneStarCount:     productRows.OneStarCount,
		TwoStarCount:     productRows.TwoStarCount,
		ThreeStarCount:   productRows.ThreeStarCount,
		FourStarCount:    productRows.FourStarCount,
		FiveStarCount:    productRows.FiveStarCount,

		UpdatedAt: productRows.UpdatedAt.String(),
		CreatedAt: productRows.CreatedAt.String(),

		IsActive: *productRows.IsActive,
		ImageUrl: productRows.ImageUrl,
		ImageId:  productRows.ImageID,

		// Initialize slices
		Categories:  []GeneralCategoryResponse{},
		Collections: []GeneralCategoryResponse{},
		Attributes:  []ProductAttribute{},
		Brand:       GeneralCategoryResponse{},
	}

	// Unmarshal JSON data
	if err := json.Unmarshal(productRows.Attributes, &resp.Attributes); err != nil {
		log.Error().Err(err).Msg("Unmarshal attributes")
	}
	if err := json.Unmarshal(productRows.Categories, &resp.Categories); err != nil {
		log.Error().Err(err).Msg("Unmarshal categories")
	}
	if err := json.Unmarshal(productRows.Collections, &resp.Collections); err != nil {
		log.Error().Err(err).Msg("Unmarshal collections")
	}
	if err := json.Unmarshal(productRows.Brand, &resp.Brand); err != nil {
		log.Error().Err(err).Msg("Unmarshal brand")
	}

	return resp
}

func mapToAdminProductResponse(productRow repository.Product) ProductListDTO {
	basePrice, _ := productRow.BasePrice.Float64Value()

	avgRating := utils.GetAvgRating(productRow.RatingCount, productRow.OneStarCount, productRow.TwoStarCount, productRow.ThreeStarCount, productRow.FourStarCount, productRow.FiveStarCount)
	product := ProductListDTO{
		ID:          productRow.ID.String(),
		Name:        productRow.Name,
		Description: productRow.Description,
		BasePrice:   basePrice.Float64,
		Sku:         productRow.BaseSku,
		Slug:        productRow.Slug,
		AvgRating:   &avgRating,
		ImageUrl:    productRow.ImageUrl,
		ImgID:       productRow.ImageID,
		ReviewCount: &productRow.RatingCount,
		CreatedAt:   productRow.CreatedAt.String(),
		UpdatedAt:   productRow.UpdatedAt.String(),
	}

	return product
}

func mapToShopProductResponse(productRow repository.GetProductListRow) ProductSummary {
	price, _ := productRow.MinPrice.Float64Value()
	avgRating := utils.GetAvgRating(productRow.RatingCount, productRow.OneStarCount, productRow.TwoStarCount, productRow.ThreeStarCount, productRow.FourStarCount, productRow.FiveStarCount)
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

func mapToVariantListModelDto(row repository.GetProductVariantListRow) VariantModelDto {
	price, _ := row.Price.Float64Value()
	variant := VariantModelDto{
		ID:       row.ID.String(),
		Price:    price.Float64,
		Stock:    row.Stock,
		IsActive: *row.IsActive,
		Sku:      row.Sku,
		ImageUrl: row.ImageUrl,
	}
	variant.Attributes = []AttributeValue{}
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

func mapAddressResponse(address repository.UserAddress) AddressResponse {
	return AddressResponse{
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
