package api

import (
	"github.com/google/uuid"
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

func mapToProductResponse(productRows repository.GetProductDetailRow) ProductDetailItemResponse {

	basePrice, _ := productRows.BasePrice.Float64Value()

	discountValue, _ := productRows.MaxDiscountValue.Float64Value()

	resp := ProductDetailItemResponse{
		ID:               productRows.ProductID.String(),
		Name:             productRows.Name,
		BasePrice:        basePrice.Float64,
		ShortDescription: productRows.ShortDescription,
		Attributes:       productRows.Attributes,
		Description:      productRows.Description,
		BaseSku:          productRows.BaseSku,
		Slug:             productRows.Slug,
		RatingCount:      productRows.RatingCount,
		OneStarCount:     productRows.OneStarCount,
		TwoStarCount:     productRows.TwoStarCount,
		ThreeStarCount:   productRows.ThreeStarCount,
		FourStarCount:    productRows.FourStarCount,
		FiveStarCount:    productRows.FiveStarCount,

		MaxDiscountValue: discountValue.Float64,
		DiscountType:     productRows.DiscountType,

		UpdatedAt: productRows.UpdatedAt.String(),
		CreatedAt: productRows.CreatedAt.String(),

		IsActive:      *productRows.IsActive,
		Variants:      make([]ProductVariantModel, 0),
		ProductImages: make([]ProductImageModel, 0),
	}

	if productRows.BrandID.Valid {
		id, _ := uuid.FromBytes(productRows.BrandID.Bytes[:])
		resp.Brand = &GeneralCategoryResponse{
			ID:   id.String(),
			Name: *productRows.BrandName,
		}
	}
	if productRows.CategoryID.Valid {
		id, _ := uuid.FromBytes(productRows.CategoryID.Bytes[:])
		resp.Category = &GeneralCategoryResponse{
			ID:   id.String(),
			Name: *productRows.CategoryName,
		}
	}
	if productRows.CollectionID.Valid {
		collectionID, _ := uuid.FromBytes(productRows.CollectionID.Bytes[:])
		resp.Collection = &GeneralCategoryResponse{
			ID:   collectionID.String(),
			Name: *productRows.CollectionName,
		}
	}

	return resp
}

func mapToVariantResp(variantRows []repository.GetProductVariantsRow) []ProductVariantModel {
	variants := make([]ProductVariantModel, 0)
	for _, row := range variantRows {
		variantIdx := -1
		for i, v := range variants {
			if v.ID == row.ID.String() {
				variantIdx = i
				break
			}
		}
		if variantIdx != -1 {
			// If the variant already exists, append the attribute to the existing variant
			attrIdx := -1
			for j, a := range variants[variantIdx].Attributes {
				if a.ID == row.AttrID {
					attrIdx = j
					break
				}
			}

			if attrIdx != -1 {
				// If the attribute already exists, do nothing
				continue
			}

			variants[variantIdx].Attributes = append(variants[variantIdx].Attributes, ProductAttributeModel{
				ID:   row.AttrID,
				Name: row.AttrName,
				ValueObject: AttributeValue{
					ID:    row.AttrValID,
					Value: row.AttrValue,
				},
			})
		} else {
			// If the variant does not exist, add it to the list of variants
			price, _ := row.Price.Float64Value()
			variant := ProductVariantModel{
				ID:       row.ID.String(),
				Price:    price.Float64,
				StockQty: row.Stock,
				IsActive: *row.IsActive,
				Sku:      &row.Sku,
				Attributes: []ProductAttributeModel{
					{
						ID:   row.AttrID,
						Name: row.AttrName,
						ValueObject: AttributeValue{
							ID:    row.AttrValID,
							Value: row.AttrValue,
						},
					},
				},
			}
			variants = append(variants, variant)
		}

	}
	return variants
}

func mapToListProductResponse(productRow repository.GetProductsRow) ProductListModel {
	minPrice, _ := productRow.MinPrice.Float64Value()
	maxPrice, _ := productRow.MaxPrice.Float64Value()
	basePrice, _ := productRow.BasePrice.Float64Value()
	if minPrice.Float64 == 0 {
		minPrice = basePrice
	}
	if maxPrice.Float64 == 0 {
		maxPrice = basePrice
	}
	avgRating := utils.GetAvgRating(productRow.RatingCount, productRow.OneStarCount, productRow.TwoStarCount, productRow.ThreeStarCount, productRow.FourStarCount, productRow.FiveStarCount)
	product := ProductListModel{
		ID:           productRow.ID.String(),
		Name:         productRow.Name,
		Description:  productRow.Description,
		BasePrice:    basePrice.Float64,
		MinPrice:     minPrice.Float64,
		MaxPrice:     maxPrice.Float64,
		Sku:          productRow.BaseSku,
		Slug:         productRow.Slug,
		ImgUrl:       &productRow.ImgUrl,
		AvgRating:    &avgRating,
		ReviewCount:  &productRow.RatingCount,
		VariantCount: productRow.VariantCount,
		CreatedAt:    productRow.CreatedAt.String(),
		UpdatedAt:    productRow.UpdatedAt.String(),
	}
	product.ImgID = &productRow.ImgID

	return product
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
