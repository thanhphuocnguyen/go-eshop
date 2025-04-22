package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

type UpdateProductVariantsTxParam struct {
	Variants []UpdateProductVariantTxParams `json:"variants" binding:"required,dive"`
}

type UpdateProductVariantsTxResult struct {
	UpdatedIDs []uuid.UUID `json:"updated_ids"`
	CreatedIDs []uuid.UUID `json:"created_ids"`
}

func (s *pgRepo) UpdateProductVariantsTx(ctx context.Context, productID uuid.UUID, arg UpdateProductVariantsTxParam) (rs UpdateProductVariantsTxResult, err error) {
	err = s.execTx(ctx, func(q *Queries) (err error) {
		product, err := q.GetProductByID(ctx, GetProductByIDParams{
			ID: productID,
		})
		if err != nil {
			log.Error().Err(err).Msg("GetProductByID")
			return
		}

		if len(arg.Variants) > 0 {
			rs.UpdatedIDs = make([]uuid.UUID, 0)
			rs.CreatedIDs = make([]uuid.UUID, 0)
			for _, variant := range arg.Variants {
				sku := ""
				attrValueIds := make([]int32, 0)
				for _, attr := range variant.Attributes {
					attrValueIds = append(attrValueIds, attr.ValueID)
				}
				variantSku, err := GetVariantSKUWithAttributeNames(q, ctx, product.BaseSku, attrValueIds)
				if err != nil {
					log.Error().Err(err).Msg("GetVariantSKUWithAttributeNames Failed")
					return err
				}
				sku = variantSku

				if variant.ID != nil {
					updateVariantParams := UpdateProductVariantParams{
						ID:  uuid.MustParse(*variant.ID),
						Sku: utils.GetPgTypeText(sku),
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

					if variant.IsActive != nil {
						updateVariantParams.IsActive = utils.GetPgTypeBool(*variant.IsActive)
					}

					updated, err := q.UpdateProductVariant(ctx, updateVariantParams)
					if err != nil {
						log.Error().Err(err).Msg("UpdateProductVariant")
						return err
					}
					rs.UpdatedIDs = append(rs.UpdatedIDs, updated.ID)
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
							createBulkProductVariantAttributesParam = append(createBulkProductVariantAttributesParam, CreateBulkProductVariantAttributeParams{
								VariantID:        updated.ID,
								AttributeValueID: attr.ValueID,
							})
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
						Sku:       sku,
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

					created, err := q.CreateProductVariant(ctx, createVariantParams)
					if err != nil {
						log.Error().Err(err).Msg("CreateProductVariant")
						return err
					}

					rs.CreatedIDs = append(rs.CreatedIDs, created.ID)

					if len(variant.Attributes) > 0 {
						createBulkProductVariantAttributesParam := make([]CreateBulkProductVariantAttributeParams, len(variant.Attributes))
						for i, attr := range variant.Attributes {
							createBulkProductVariantAttributesParam[i] = CreateBulkProductVariantAttributeParams{
								VariantID:        created.ID,
								AttributeValueID: attr.ValueID,
							}
						}

						rs, err := q.CreateBulkProductVariantAttribute(ctx, createBulkProductVariantAttributesParam)
						if err != nil {
							log.Error().Err(err).Msg("CreateProductVariantAttribute")
							return err
						}

						log.Debug().Msgf("CreateProductVariantAttribute %v", rs)
					}
				}
			}
		}

		return
	})
	return
}
