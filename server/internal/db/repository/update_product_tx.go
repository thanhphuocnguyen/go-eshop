package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

// CreateProductTx creates a new product in a transaction
type UpdateProductTxParam struct {
	Name          *string                        `json:"name" binding:"omitempty,min=3,max=100"`
	Description   *string                        `json:"description" binding:"omitempty,min=6,max=1000"`
	Price         *float64                       `json:"price" binding:"omitempty,gt=0"`
	Sku           *string                        `json:"sku" binding:"omitempty,alphanum"`
	Slug          *string                        `json:"slug" binding:"omitempty"`
	Stock         *int32                         `json:"stock" binding:"omitempty,gt=0"`
	CategoryID    *string                        `json:"category_id,omitempty" binding:"omitempty,uuid"`
	BrandID       *string                        `json:"brand_id,omitempty" binding:"omitempty,uuid"`
	CollectionID  *string                        `json:"collection_id,omitempty" binding:"omitempty,uuid"`
	RemovedImages []int32                        `json:"removed_images" binding:"omitempty,dive"`
	Variants      []UpdateProductVariantTxParams `json:"variants" binding:"omitempty,dive"`
}

func (s *pgRepo) UpdateProductTx(ctx context.Context, productID uuid.UUID, arg UpdateProductTxParam) (rs ProductTxResult, err error) {
	err = s.execTx(ctx, func(q *Queries) error {
		log.Info().Interface("variants", arg.Variants).Msg("Log variants")
		updateProductParam := UpdateProductParams{
			ID: productID,
		}
		if arg.Name != nil {
			updateProductParam.Name = utils.GetPgTypeText(*arg.Name)
		}
		if arg.Description != nil {
			updateProductParam.Description = utils.GetPgTypeText(*arg.Description)
		}
		if arg.Slug != nil {
			updateProductParam.Slug = utils.GetPgTypeText(*arg.Slug)
		}
		if arg.Sku != nil {
			updateProductParam.BaseSku = utils.GetPgTypeText(*arg.Sku)
		}
		if arg.CategoryID != nil {
			updateProductParam.CategoryID = utils.GetPgTypeUUIDFromString(*arg.CategoryID)
		}
		if arg.CollectionID != nil {
			updateProductParam.CollectionID = utils.GetPgTypeUUIDFromString(*arg.CollectionID)
		}
		if arg.BrandID != nil {
			updateProductParam.BrandID = utils.GetPgTypeUUIDFromString(*arg.BrandID)
		}
		if arg.Price != nil {
			updateProductParam.BasePrice = utils.GetPgNumericFromFloat(*arg.Price)
		}

		product, err := q.UpdateProduct(ctx, updateProductParam)
		if err != nil {
			log.Error().Err(err).Msg("UpdateProduct")
			return err
		}

		if len(arg.Variants) > 0 {
			rs.Variants = make([]uuid.UUID, 0)
			for _, variant := range arg.Variants {
				if variant.ID != nil {
					updateVariantParams := UpdateProductVariantParams{
						ID: uuid.MustParse(*variant.ID),
					}
					if variant.Stock != nil {
						updateVariantParams.Stock = utils.GetPgTypeInt4(*variant.Stock)
					}
					if variant.Weight != nil {
						updateVariantParams.Weight = utils.GetPgNumericFromFloat(*variant.Weight)
					}
					if variant.Price != nil {
						updateVariantParams.Price = utils.GetPgNumericFromFloat(*variant.Price)
					}
					if variant.Sku != nil {
						updateVariantParams.Sku = utils.GetPgTypeText(*variant.Sku)
					}
					if variant.IsActive != nil {
						updateVariantParams.IsActive = utils.GetPgTypeBool(*variant.IsActive)
					}

					updated, err := q.UpdateProductVariant(ctx, updateVariantParams)
					if err != nil {
						log.Error().Err(err).Msg("UpdateProductVariant")
						return err
					}
					rs.Variants = append(rs.Variants, updated.ID)
					log.Info().Interface("attributes ", variant.Attributes).Msg("Log attributes")
					if len(variant.Attributes) > 0 {
						// delete old attributes
						err := q.DeleteProductVariantAttributes(ctx, updated.ID)
						if err != nil {
							log.Error().Err(err).Msg("DeleteProductVariantAttributes")
							return err
						}
						// create new attributes
						createBulkProductVariantAttributesParam := make([]CreateBulkProductVariantAttributeParams, 0)
						for _, attr := range variant.Attributes {
							if len(attr.ValueIDs) == 0 {
								continue
							}
							for _, valueID := range attr.ValueIDs {
								createBulkProductVariantAttributesParam = append(createBulkProductVariantAttributesParam, CreateBulkProductVariantAttributeParams{
									VariantID:        updated.ID,
									AttributeValueID: valueID,
								})
							}
						}
						_, err = q.CreateBulkProductVariantAttribute(ctx, createBulkProductVariantAttributesParam)
						if err != nil {
							log.Error().Err(err).Msg("CreateProductVariantAttribute")
							return err
						}
					}
				} else {
					createVariantParams := CreateProductVariantParams{
						ID:        uuid.New(),
						ProductID: product.ID,
					}
					if variant.Stock != nil {
						createVariantParams.Stock = *variant.Stock
					}
					if variant.Weight != nil {
						createVariantParams.Weight = utils.GetPgNumericFromFloat(*variant.Weight)
					}
					if variant.Price != nil {
						createVariantParams.Price = utils.GetPgNumericFromFloat(*variant.Price)
					}
					if variant.Sku != nil {
						createVariantParams.Sku = *variant.Sku
					}

					created, err := q.CreateProductVariant(ctx, createVariantParams)
					if err != nil {
						log.Error().Err(err).Msg("CreateProductVariant")
						return err
					}

					rs.Variants = append(rs.Variants, created.ID)

					if len(variant.Attributes) > 0 {
						createBulkProductVariantAttributesParam := make([]CreateBulkProductVariantAttributeParams, 0)
						for _, attr := range variant.Attributes {
							if len(attr.ValueIDs) == 0 {
								continue
							}
							for _, valueID := range attr.ValueIDs {
								createBulkProductVariantAttributesParam = append(createBulkProductVariantAttributesParam, CreateBulkProductVariantAttributeParams{
									VariantID:        created.ID,
									AttributeValueID: valueID,
								})
							}
						}
						_, err = q.CreateBulkProductVariantAttribute(ctx, createBulkProductVariantAttributesParam)
						if err != nil {
							log.Error().Err(err).Msg("CreateProductVariantAttribute")
							return err
						}
					}
				}
			}
		}

		rs.ID = product.ID

		return nil
	})
	return
}
