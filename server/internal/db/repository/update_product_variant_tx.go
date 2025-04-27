package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

type UpdateProdVariantsTxArgs struct {
	Variants []UpdateProductVariantTxParams `json:"variants" binding:"required,dive"`
}

type UpdateProductVariantsTxResult struct {
	UpdatedIDs []uuid.UUID `json:"updated_ids"`
	CreatedIDs []uuid.UUID `json:"created_ids"`
}

func (s *pgRepo) UpdateProductVariantsTx(ctx context.Context, productID uuid.UUID, arg UpdateProdVariantsTxArgs) (rs UpdateProductVariantsTxResult, err error) {
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
			// get all attribute values
			attrValMp := make(map[int32]bool)
			attrValueIds := make([]int32, 0)
			for _, variant := range arg.Variants {
				for _, attr := range variant.Attributes {
					if attrValMp[attr.ValueID] {
						continue
					}
					attrValueIds = append(attrValueIds, attr.ValueID)
					attrValMp[attr.ValueID] = true
				}
			}

			attributeValueMp := make(map[int32]AttributeValue)
			attributeValueRows, err := q.GetAttributeValuesByIDs(ctx, attrValueIds)

			if err != nil {
				log.Error().Err(err).Msg("GetAttributeValuesByIDs")
				return err
			}
			for _, attributeValue := range attributeValueRows {
				if _, ok := attributeValueMp[attributeValue.ID]; !ok {
					attributeValueMp[attributeValue.ID] = attributeValue
				}
			}
			for i, variant := range arg.Variants {
				sku := ""
				attrValues := make([]AttributeValue, 0)
				for _, attr := range variant.Attributes {
					attrValues = append(attrValues, attributeValueMp[attr.ValueID])
				}

				variantSku := GetVariantSKUWithAttributeNames(product.BaseSku, attrValues)

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
						log.Error().Err(err).Msgf("UpdateProductVariant at index %d", i)
						return err
					}

					rs.UpdatedIDs = append(rs.UpdatedIDs, updated.ID)
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

						_, err := q.CreateBulkProductVariantAttribute(ctx, createBulkProductVariantAttributesParam)
						if err != nil {
							log.Error().Err(err).Msg("CreateProductVariantAttribute")
							return err
						}

					}
				}
			}
		}

		return
	})
	return
}
