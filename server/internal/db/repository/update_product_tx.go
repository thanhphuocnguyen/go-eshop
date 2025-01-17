package repository

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/util"
)

// CreateProductTx creates a new product in a transaction
type UpdateProductTxParam struct {
	ProductID   int64
	Name        *string
	Description *string
	Price       *float64
	Discount    *int32
	Stock       *int32
	BrandID     *int64
	CategoryID  *int32
	Sku         *string
	Variants    []UpdateVariantTxParams
}

type UpdateProductTxResult struct {
	Product Product          `json:"product"`
	Variant []ProductVariant `json:"variant"`
}

func (s *pgRepo) UpdateProductTx(ctx context.Context, arg UpdateProductTxParam) (result UpdateProductTxResult, err error) {
	err = s.execTx(ctx, func(q *Queries) error {
		updateProductParam := UpdateProductParams{}
		if arg.Name != nil {
			updateProductParam.Name = util.GetPgTypeText(*arg.Name)
		}
		if arg.Description != nil {
			updateProductParam.Description = util.GetPgTypeText(*arg.Description)
		}
		if arg.Price != nil {
			updateProductParam.Price = util.GetPgNumericFromFloat(*arg.Price)
		}
		if arg.Stock != nil {
			updateProductParam.Stock = util.GetPgTypeInt4(*arg.Stock)
		}

		if arg.Discount != nil {
			updateProductParam.Discount = util.GetPgTypeInt4(*arg.Discount)
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
		result.Variant = make([]ProductVariant, 0)
		if len(arg.Variants) > 0 {
			for _, params := range arg.Variants {
				rs, err := updateVariantUtil(ctx, q, params)
				if err != nil {
					return err
				}
				result.Variant = append(result.Variant, *rs)
			}
		}
		return nil
	})
	return result, err
}
