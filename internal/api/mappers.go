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

func mapToProductResponse(productRows repository.GetProductDetailRow) ManageProductDetailResp {
	basePrice, _ := productRows.BasePrice.Float64Value()

	resp := ManageProductDetailResp{
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

func mapToListProductResponse(productRow repository.Product) ManageProductListModel {
	basePrice, _ := productRow.BasePrice.Float64Value()

	avgRating := utils.GetAvgRating(productRow.RatingCount, productRow.OneStarCount, productRow.TwoStarCount, productRow.ThreeStarCount, productRow.FourStarCount, productRow.FiveStarCount)
	product := ManageProductListModel{
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
