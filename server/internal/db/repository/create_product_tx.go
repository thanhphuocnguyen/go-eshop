package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

// CreateProductTx creates a new product in a transaction
type CreateProductTxParam struct {
	CategoryID   *int32                 `json:"category_id,omitempty"`
	CollectionID *int32                 `json:"collection_id,omitempty"`
	BrandID      *int32                 `json:"brand_id,omitempty"`
	Name         string                 `json:"name" binding:"required,min=3,max=100"`
	Description  string                 `json:"description" binding:"required,min=10,max=1000"`
	Variants     []CreateVariantTxParam `json:"variants" binding:"omitempty,dive"`
}

type CreateProductTxResult struct {
	Product  Product          `json:"product"`
	Variants []ProductVariant `json:"variants"`
}

func (s *pgRepo) CreateProductTx(ctx context.Context, arg CreateProductTxParam) (CreateProductTxResult, error) {
	var result CreateProductTxResult
	err := s.execTx(ctx, func(q *Queries) error {
		createProductParam := CreateProductParams{
			Name:        arg.Name,
			ProductID:   uuid.New(),
			Description: arg.Description,
		}

		if arg.CategoryID != nil {
			createProductParam.CategoryID = utils.GetPgTypeInt4(*arg.CategoryID)
		}
		if arg.CollectionID != nil {
			createProductParam.CollectionID = utils.GetPgTypeInt4(*arg.CollectionID)
		}
		if arg.BrandID != nil {
			createProductParam.BrandID = utils.GetPgTypeInt4(*arg.BrandID)
		}
		product, err := q.CreateProduct(ctx, createProductParam)
		if err != nil {
			log.Error().Err(err).Msg("CreateProduct")
			return err
		}

		result.Product = product
		result.Variants = make([]ProductVariant, 0)
		if len(arg.Variants) > 0 {
			for _, params := range arg.Variants {
				params.ProductID = product.ProductID
				variantCreated, err := createVariantUtil(ctx, q, params)
				if err != nil {
					log.Error().Err(err).Msg("CreateVariant")
					return err
				}
				result.Variants = append(result.Variants, *variantCreated)
			}
		}

		return nil
	})

	return result, err
}
