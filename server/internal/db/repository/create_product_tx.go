package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

// CreateProductTx creates a new product in a transaction
type CreateProductTxParam struct {
	Name         string                         `json:"name" binding:"required,min=3,max=100"`
	Description  string                         `json:"description" binding:"omitempty,min=6,max=5000"`
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
		createParams := CreateProductParams{
			ID:   uuid.New(),
			Name: arg.Name,
		}

		createParams.BasePrice = utils.GetPgNumericFromFloat(arg.Price)
		createParams.Description = utils.GetPgTypeText(arg.Description)
		createParams.Slug = arg.Slug
		createParams.BaseSku = utils.GetPgTypeText(arg.Sku)

		createParams.CategoryID = utils.GetPgTypeUUIDFromString(arg.CategoryID)

		createParams.BrandID = utils.GetPgTypeUUIDFromString(arg.BrandID)
		if arg.CollectionID != nil {
			createParams.CollectionID = utils.GetPgTypeUUIDFromString(*arg.CollectionID)
		}

		product, err := q.CreateProduct(ctx, createParams)
		if err != nil {
			log.Error().Err(err).Timestamp()
			return err
		}
		result.Variants = make([]uuid.UUID, 0)
		createVariantParams := make([]CreateBulkProductVariantsParams, len(arg.Variants))
		createVariantAttributeParams := make([]CreateBulkProductVariantAttributeParams, 0)
		variantIDs := make([]uuid.UUID, 0)
		for i, variantReq := range arg.Variants {
			attributeIDs := make([]int32, 0)
			for _, attr := range variantReq.Attributes {
				attributeIDs = append(attributeIDs, attr.ValueID)
			}
			variantSku, err := GetVariantSKUWithAttributeNames(q, ctx, arg.Sku, attributeIDs)
			if err != nil {
				log.Error().Err(err).Msg("GetVariantSKUWithAttributeNames Failed")
				return err
			}
			newVariantID := uuid.New()
			createVariantParam := CreateBulkProductVariantsParams{
				ProductID: product.ID,
				ID:        newVariantID,
				Sku:       variantSku,
				Price:     utils.GetPgNumericFromFloat(variantReq.Price),
				Stock:     variantReq.Stock,
			}

			variantIDs = append(variantIDs, newVariantID)
			if variantReq.Weight != nil {
				createVariantParam.Weight = utils.GetPgNumericFromFloat(*variantReq.Weight)
			}

			createVariantParams[i] = createVariantParam

			for _, attr := range variantReq.Attributes {
				createProductAttributesParam := CreateBulkProductVariantAttributeParams{
					VariantID:        newVariantID,
					AttributeValueID: attr.ValueID,
				}
				createVariantAttributeParams = append(createVariantAttributeParams, createProductAttributesParam)
			}
		}

		rs, err := q.CreateBulkProductVariants(ctx, createVariantParams)
		if err != nil {
			log.Error().Err(err).Msg("CreateVariantProduct Failed")
			return err
		}
		log.Debug().Msgf("CreateVariantProduct %v", rs)

		rs, err = q.CreateBulkProductVariantAttribute(ctx, createVariantAttributeParams)
		if err != nil {
			log.Error().Err(err).Msg("CreateVariantProductAttribute Failed")
			return err
		}
		log.Debug().Msgf("CreateVariantProductAttribute %v", rs)
		result.Variants = variantIDs

		result.ID = product.ID
		return nil
	})

	return result, err
}
