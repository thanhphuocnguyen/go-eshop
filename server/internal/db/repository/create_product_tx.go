package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

// CreateProductTx creates a new product in a transaction
type CreateProductTxParam struct {
	Name         string                         `json:"name" binding:"required,min=3,max=100"`
	Description  string                         `json:"description" binding:"omitempty,min=6,max=1000"`
	Price        float64                        `json:"price" binding:"required,gt=0"`
	Sku          string                         `json:"sku" binding:"required"`
	Slug         string                         `json:"slug" binding:"omitempty"`
	CategoryID   string                         `json:"category_id,omitempty" binding:"omitempty,uuid"`
	BrandID      string                         `json:"brand_id,omitempty" binding:"omitempty,uuid"`
	CollectionID *string                        `json:"collection_id,omitempty" binding:"omitempty,uuid"`
	Variants     []CreateProductVariantTxParams `json:"variants,omitempty"`
}

func (s *pgRepo) CreateProductTx(ctx context.Context, arg CreateProductTxParam) (ProductTxResult, error) {
	var result ProductTxResult
	err := s.execTx(ctx, func(q *Queries) error {
		createProductParam := CreateProductParams{
			ID:   uuid.New(),
			Name: arg.Name,
		}

		createProductParam.BasePrice = utils.GetPgNumericFromFloat(arg.Price * 100)
		createProductParam.Description = utils.GetPgTypeText(arg.Description)
		createProductParam.Slug = arg.Slug
		createProductParam.BaseSku = utils.GetPgTypeText(arg.Sku)

		createProductParam.CategoryID = utils.GetPgTypeUUIDFromString(arg.CategoryID)

		createProductParam.BrandID = utils.GetPgTypeUUIDFromString(arg.BrandID)
		if arg.CollectionID != nil {
			createProductParam.CollectionID = utils.GetPgTypeUUIDFromString(*arg.CollectionID)
		}

		product, err := q.CreateProduct(ctx, createProductParam)
		if err != nil {
			log.Error().Err(err).Timestamp()
			return err
		}
		result.Variants = make([]uuid.UUID, 0)
		for _, variantReq := range arg.Variants {
			createVariantParam := CreateProductVariantParams{
				ProductID: product.ID,
				ID:        uuid.New(),
				Sku:       variantReq.Sku,
				Price:     utils.GetPgNumericFromFloat(variantReq.Price),
				Stock:     variantReq.Stock,
			}

			if variantReq.Weight != nil {
				createVariantParam.Weight = utils.GetPgNumericFromFloat(*variantReq.Weight)
			}

			createdVariant, err := q.CreateProductVariant(ctx, createVariantParam)

			if err != nil {
				log.Error().Err(err).Msg("CreateVariantProduct")
				return err
			}

			for _, attr := range variantReq.Attributes {
				for i, valueID := range attr.ValueIDs {
					createProductAttributesParam := CreateProductVariantAttributeParams{
						VariantID:        createdVariant.ID,
						AttributeValueID: valueID,
					}

					_, err := q.CreateProductVariantAttribute(ctx, createProductAttributesParam)
					if err != nil {
						log.Error().Err(err).Msg(fmt.Sprintf("index: %d, id: %d", i, valueID))
						return err
					}
				}
			}
			result.Variants = append(result.Variants, createdVariant.ID)
		}

		result.ID = product.ID
		return nil
	})

	return result, err
}
