package api

import (
	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

func mapToUserResponse(user repository.User) UserResponse {
	return UserResponse{
		ID:                user.ID,
		Addresses:         []AddressResponse{},
		Email:             user.Email,
		FullName:          user.FirstName,
		Role:              user.Role,
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
		Default:  address.Default,
		ID:       address.ID.String(),
	}
}

func mapToProductResponse(productRows repository.GetProductDetailRow) ProductDetailItemResponse {

	basePrice, _ := productRows.BasePrice.Float64Value()
	attributes := make([]string, len(productRows.Attributes))
	for i, attr := range productRows.Attributes {
		attributes[i] = attr.String()
	}
	discountValue, _ := productRows.MaxDiscountValue.Float64Value()

	resp := ProductDetailItemResponse{
		ID:               productRows.ProductID.String(),
		Name:             productRows.Name,
		BasePrice:        basePrice.Float64,
		ShortDescription: productRows.ShortDescription,
		Attributes:       attributes,
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
				if a.ID == row.AttrID.String() {
					attrIdx = j
					break
				}
			}

			if attrIdx != -1 {
				// If the attribute already exists, do nothing
				continue
			}

			variants[variantIdx].Attributes = append(variants[variantIdx].Attributes, ProductAttributeModel{
				ID:   row.AttrID.String(),
				Name: row.AttrName,
				ValueObject: AttributeValue{
					ID:           row.AttrValID,
					Code:         row.AttrValCode,
					Name:         &row.AttrValName,
					IsActive:     row.IsActive,
					DisplayOrder: &row.AttrDisplayOrder,
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
						ID:   row.AttrID.String(),
						Name: row.AttrName,
						ValueObject: AttributeValue{
							ID:           row.AttrValID,
							Code:         row.AttrValCode,
							Name:         &row.AttrValName,
							IsActive:     row.IsActive,
							DisplayOrder: &row.AttrDisplayOrder,
						},
					},
				},
			}
			variants = append(variants, variant)
		}

	}
	return variants
}

func mapToProductImages(productID uuid.UUID, imageRows []repository.GetProductImagesAssignedRow) []ProductImageModel {
	// log.Debug().Msgf("mapToProductImages: %v", imageRows)
	images := make([]ProductImageModel, 0)
	for _, row := range imageRows {
		existingImageIdx := -1
		for i, image := range images {
			if image.ID == row.ID.String() {
				existingImageIdx = i
				break
			}
		}
		if existingImageIdx != -1 {
			image := ImageAssignmentModel{
				ID:           row.ID.String(),
				EntityID:     row.EntityID.String(),
				EntityType:   row.EntityType,
				Role:         row.Role,
				DisplayOrder: row.DisplayOrder,
			}
			if row.EntityID != productID {
				// If the image already exists, append the assignment to the existing image
				images[existingImageIdx].VariantAssignments = append(images[existingImageIdx].VariantAssignments, image)
			}
		} else {
			// If the image does not exist, add it to the list of images
			image := ProductImageModel{
				ID:                 row.ID.String(),
				Url:                row.Url,
				ExternalID:         row.ExternalID,
				Role:               row.Role,
				VariantAssignments: make([]ImageAssignmentModel, 0),
			}

			if row.EntityID != productID {
				image.VariantAssignments = append(image.VariantAssignments, ImageAssignmentModel{
					ID:           row.ID.String(),
					EntityID:     row.EntityID.String(),
					EntityType:   row.EntityType,
					Role:         row.Role,
					DisplayOrder: row.DisplayOrder,
				})
			}
			images = append(images, image)
		}
	}
	return images
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
		ImgUrl:       productRow.ImgUrl,
		AvgRating:    &avgRating,
		ReviewCount:  &productRow.RatingCount,
		VariantCount: productRow.VariantCount,
		CreatedAt:    productRow.CreatedAt.String(),
		UpdatedAt:    productRow.UpdatedAt.String(),
	}
	if productRow.ImgID.Valid {
		id, _ := uuid.FromBytes(productRow.ImgID.Bytes[:])
		product.ImgID = utils.StringPtr(id.String())
	}

	return product
}

func mapAddressResponse(address repository.UserAddress) AddressResponse {
	return AddressResponse{
		ID:        address.ID.String(),
		Default:   address.Default,
		CreatedAt: address.CreatedAt,
		Phone:     address.PhoneNumber,
		Street:    address.Street,
		Ward:      address.Ward,
		District:  address.District,
		City:      address.City,
	}
}
