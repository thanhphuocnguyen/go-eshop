package repository

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/util"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

// CreateProductTx creates a new product in a transaction
type CreateProductTxParam struct {
	Name        string
	Description string
	Price       float64
	Discount    *int32
	Stock       int32
	BrandID     *int64
	CategoryID  *int32
	Variants    []CreateVariantTxParam
}

type CreateProductTxResult struct {
	Product Product          `json:"product"`
	Variant []ProductVariant `json:"variant"`
}

func (s *pgRepo) CreateProductTx(ctx context.Context, arg CreateProductTxParam) (CreateProductTxResult, error) {
	var result CreateProductTxResult
	err := s.execTx(ctx, func(q *Queries) error {
		createProductParam := CreateProductParams{
			Name:        arg.Name,
			Description: arg.Description,
			Sku:         util.GetPgTypeText(utils.GenerateSKU()),
			Stock:       arg.Stock,
			Price:       util.GetPgNumericFromFloat(arg.Price),
		}
		if arg.Discount != nil {
			createProductParam.Discount = *arg.Discount
		}

		product, err := q.CreateProduct(ctx, createProductParam)
		if err != nil {
			log.Error().Err(err).Msg("CreateProduct")
			return err
		}

		result.Product = product
		result.Variant = make([]ProductVariant, 0)
		if len(arg.Variants) > 0 {
			for _, params := range arg.Variants {
				params.ProductID = product.ProductID
				variantCreated, err := createVariantUtil(ctx, q, params)
				if err != nil {
					log.Error().Err(err).Msg("CreateVariant")
					return err
				}
				result.Variant = append(result.Variant, *variantCreated)
			}
		}
		if arg.CategoryID != nil {
			maxSortOrder, err := q.GetMaxSortOrderInCollection(ctx, *arg.CategoryID)
			if err != nil {
				log.Error().Err(err).Msg("GetMaxSortOrderInCollection")
				return err
			}

			_, err = q.AddProductToCollection(ctx, AddProductToCollectionParams{
				CategoryID: *arg.CategoryID,
				ProductID:  product.ProductID,
				SortOrder:  maxSortOrder + 1,
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
