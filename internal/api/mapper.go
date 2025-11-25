package api

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

func mapToUserResponse(user repository.User, roleCode string) dto.UserDetail {
	return dto.UserDetail{
		ID:                user.ID,
		Addresses:         []dto.AddressDetail{},
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

func mapAddressToAddressResponse(address repository.UserAddress) dto.AddressDetail {
	return dto.AddressDetail{
		Phone:    address.PhoneNumber,
		Street:   address.Street,
		Ward:     address.Ward,
		District: address.District,
		City:     address.City,
		Default:  address.IsDefault,
		ID:       address.ID.String(),
	}
}

func mapToProductDetailResponse(row repository.GetProductDetailRow) dto.ProductDetail {
	basePrice, _ := row.BasePrice.Float64Value()

	resp := dto.ProductDetail{
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
		Categories:  []dto.GeneralCategory{},
		Collections: []dto.GeneralCategory{},
		Attributes:  []dto.ProductAttribute{},
		Brand:       dto.GeneralCategory{},
		Variations:  []dto.VariantDetail{},
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

func mapToAdminProductResponse(productRow repository.Product) dto.ProductListItem {
	basePrice, _ := productRow.BasePrice.Float64Value()

	avgRating := utils.GetAvgRating(productRow.RatingCount, productRow.OneStarCount, productRow.TwoStarCount, productRow.ThreeStarCount, productRow.FourStarCount, productRow.FiveStarCount)
	product := dto.ProductListItem{
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

func mapToShopProductResponse(productRow repository.GetProductListRow) dto.ProductSummary {
	price, _ := productRow.MinPrice.Float64Value()
	avgRating := utils.GetAvgRating(productRow.RatingCount, productRow.OneStarCount, productRow.TwoStarCount, productRow.ThreeStarCount, productRow.FourStarCount, productRow.FiveStarCount)
	product := dto.ProductSummary{
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

func mapToVariantListModelDto(row repository.GetProductVariantListRow) dto.VariantDetail {
	price, _ := row.Price.Float64Value()
	variant := dto.VariantDetail{
		ID:       row.ID.String(),
		Price:    price.Float64,
		Stock:    row.Stock,
		IsActive: *row.IsActive,
		Sku:      row.Sku,
		ImageUrl: row.ImageUrl,
	}
	variant.Attributes = []dto.AttributeValueDetail{}
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

func mapAddressResponse(address repository.UserAddress) dto.AddressDetail {
	return dto.AddressDetail{
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
