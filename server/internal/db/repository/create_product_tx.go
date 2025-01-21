package repository

import (
	"context"

	"github.com/rs/zerolog/log"
)

// CreateProductTx creates a new product in a transaction
type CreateProductTxParam struct {
	Name        string                 `json:"name" binding:"required,min=3,max=100"`
	Description string                 `json:"description" binding:"required,min=10,max=1000"`
	CategoryID  *int32                 `json:"category_id,omitempty"`
	Variants    []CreateVariantTxParam `json:"variants" binding:"omitempty,dive"`
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
			Description: arg.Description,
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
		if arg.CategoryID != nil {
			maxSortOrder, err := q.GetMaxSortOrderInCollection(ctx, *arg.CategoryID)
			if err != nil {
				log.Error().Err(err).Msg("GetMaxSortOrderInCollection")
				return err
			}
			if maxSortOrder == nil {
				maxSortOrder = 0
			}

			_, err = q.AddProductToCollection(ctx, AddProductToCollectionParams{
				CategoryID: *arg.CategoryID,
				ProductID:  product.ProductID,
				SortOrder:  int16(maxSortOrder.(int) + 1),
			})
			if err != nil {
				log.Error().Err(err).Msg("AddProductToCollection")
				return err
			}
		}

		return nil
	})

	return result, err
}
