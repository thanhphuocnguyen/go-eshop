package repository

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

// CreateProductTx creates a new product in a transaction
type UpdateProductTxParam struct {
	ProductID   int64
	Name        *string
	Description *string
	CategoryID  *int32
	Variants    []UpdateVariantTxParams
}

type UpdateProductTxResult struct {
	Product  Product `json:"product"`
	Variants []UpdateVariantTxResult
}

func (s *pgRepo) UpdateProductTx(ctx context.Context, arg UpdateProductTxParam) (result UpdateProductTxResult, err error) {
	err = s.execTx(ctx, func(q *Queries) error {
		updateProductParam := UpdateProductParams{}
		if arg.Name != nil {
			updateProductParam.Name = utils.GetPgTypeText(*arg.Name)
		}
		if arg.Description != nil {
			updateProductParam.Description = utils.GetPgTypeText(*arg.Description)
		}

		product, err := q.UpdateProduct(ctx, updateProductParam)
		if err != nil {
			log.Error().Err(err).Msg("CreateProduct")
			return err
		}

		if arg.CategoryID != nil {
			q.AddProductToCollection(ctx, AddProductToCollectionParams{
				CategoryID: *arg.CategoryID,
				ProductID:  product.ProductID,
			})
		}

		result.Product = product
		result.Variants = make([]UpdateVariantTxResult, 0)
		if len(arg.Variants) > 0 {
			for _, params := range arg.Variants {
				variantUpdated, attributesUpdated, err := updateVariantUtil(ctx, q, params)
				if err != nil {
					return err
				}
				result.Variants = append(result.Variants, UpdateVariantTxResult{
					Variant:    variantUpdated,
					Attributes: attributesUpdated,
				})
			}
		}
		return nil
	})
	return result, err
}
