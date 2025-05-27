package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"golang.org/x/sync/errgroup"
)

type UpdateProductImages struct {
	ID          string   `json:"id"`
	Role        *string  `json:"role"`
	IsRemoved   *bool    `json:"omitempty,is_removed"`
	Assignments []string `json:"assignments,omitempty"`
}

type UpdateProductTxParams struct {
	Name             *string                        `json:"name" binding:"omitempty,min=3,max=100"`
	Description      *string                        `json:"description" binding:"omitempty,min=6,max=5000"`
	ShortDescription *string                        `json:"short_description" binding:"omitempty,max=2000"`
	Price            *float64                       `json:"price" binding:"omitempty,gt=0"`
	Sku              *string                        `json:"sku" binding:"omitempty"`
	Slug             *string                        `json:"slug" binding:"omitempty"`
	Stock            *int32                         `json:"stock" binding:"omitempty,gt=0"`
	CategoryID       *string                        `json:"categoryId,omitempty" binding:"omitempty,uuid"`
	BrandID          *string                        `json:"brandId,omitempty" binding:"omitempty,uuid"`
	CollectionID     *string                        `json:"collectionId,omitempty" binding:"omitempty,uuid"`
	Attributes       []string                       `json:"attributes" binding:"omitempty"`
	Images           []UpdateProductImages          `json:"images" binding:"omitempty,dive"`
	Variants         []UpdateProductVariantTxParams `json:"variants" binding:"omitempty,dive"`
	RemoveImageFn    func(ctx context.Context, externalID string) (string, error)
}

func (s *pgRepo) UpdateProductTx(ctx context.Context, productID uuid.UUID, arg UpdateProductTxParams) (err error) {
	err = s.execTx(ctx, func(q *Queries) (err error) {
		product, err := q.GetProductByID(ctx, GetProductByIDParams{
			ID: productID,
		})
		updateProductParam := UpdateProductParams{
			ID: productID,
		}

		if arg.Name != nil {
			updateProductParam.Name = arg.Name
		}
		if arg.Description != nil {
			updateProductParam.Description = arg.Description
		}
		if arg.ShortDescription != nil {
			updateProductParam.ShortDescription = arg.ShortDescription
		}
		if arg.Slug != nil {
			updateProductParam.Slug = arg.Slug
		}
		if arg.Sku != nil {
			updateProductParam.BaseSku = arg.Sku
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
		if len(arg.Attributes) > 0 {
			attributes := make([]uuid.UUID, len(arg.Attributes))
			for i, attr := range arg.Attributes {
				attributes[i] = uuid.MustParse(attr)
			}
			updateProductParam.Attributes = attributes
		}

		product, err = q.UpdateProduct(ctx, updateProductParam)

		if err != nil {
			log.Error().Err(err).Msg("UpdateProduct")
			return err
		}

		if len(arg.Images) > 0 {
			errGroup, _ := errgroup.WithContext(ctx)
			imgAssignmentArgs := make([]UpdateProdImagesTxArgs, 0)
			for _, image := range arg.Images {
				if image.IsRemoved != nil && *image.IsRemoved {
					errGroup.Go(func() (err error) {
						img, err := q.GetImageFromID(ctx, GetImageFromIDParams{
							ID:         uuid.MustParse(image.ID),
							EntityType: string(EntityTypeProduct),
						})

						if err != nil {

							return err
						}

						msg, err := arg.RemoveImageFn(ctx, img.ExternalID)
						if err != nil {
							return fmt.Errorf("failed to remove image: %w, reason: %s", err, msg)
						}

						// Remove image from product
						err = q.DeleteProductImage(ctx, uuid.MustParse(image.ID))
						if err != nil {
							return err
						}
						return err
					})
				} else {
					// parse the assignments to UUIDs
					assignmentIds := make([]uuid.UUID, 0)
					for _, assignment := range image.Assignments {
						if assignment == productID.String() {
							continue
						}
						assignmentIds = append(assignmentIds, uuid.MustParse(assignment))
					}
					// append the image assignment to the list
					args := UpdateProdImagesTxArgs{
						ImageID:    image.ID,
						Role:       image.Role,
						EntityID:   productID,
						EntityType: string(EntityTypeProduct),
						VariantIDs: assignmentIds,
					}

					imgAssignmentArgs = append(imgAssignmentArgs, args)
				}

			}

			err = errGroup.Wait()
			if err != nil {
				return err
			}

			if len(imgAssignmentArgs) != 0 {
				// update the image assignments
				err := updateProductImages(ctx, q, imgAssignmentArgs)
				if err != nil {
					return err
				}
			}

		}
		if err != nil {
			log.Error().Err(err).Msg("GetProductByID")
			return err
		}

		if len(arg.Variants) > 0 {
			// get all attribute values
			attrValMp := make(map[string]bool)
			attrValueIds := make([]uuid.UUID, 0)
			for _, variant := range arg.Variants {
				for _, attr := range variant.Attributes {
					if attrValMp[attr.ValueID] {
						continue
					}
					id, parseErr := uuid.Parse(attr.ValueID)
					if parseErr != nil {
						log.Error().Err(parseErr).Msgf("Parse UUID from %s", attr.ValueID)
						return parseErr
					}
					attrValueIds = append(attrValueIds, id)
					attrValMp[attr.ValueID] = true
				}
			}

			attributeValueMp := make(map[string]AttributeValue)
			attributeValueRows, err := q.GetAttributeValuesByIDs(ctx, attrValueIds)

			if err != nil {
				log.Error().Err(err).Msg("GetAttributeValuesByIDs")
				return err
			}
			for _, attributeValue := range attributeValueRows {
				if _, ok := attributeValueMp[attributeValue.ID.String()]; !ok {
					attributeValueMp[attributeValue.ID.String()] = attributeValue
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
						Sku: &sku,
					}
					if variant.Stock != nil {
						updateVariantParams.Stock = variant.Stock
					}
					if variant.Weight != nil {
						updateVariantParams.Weight = utils.GetPgNumericFromFloat(*variant.Weight)
					}
					if variant.Price != nil {
						updateVariantParams.Price = utils.GetPgNumericFromFloat(*variant.Price)
					}

					if variant.IsActive != nil {
						updateVariantParams.IsActive = variant.IsActive
					}

					updated, err := q.UpdateProductVariant(ctx, updateVariantParams)
					if err != nil {
						log.Error().Err(err).Msgf("UpdateProductVariant at index %d", i)
						return err
					}

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
								AttributeValueID: uuid.MustParse(attr.ValueID),
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

					if len(variant.Attributes) > 0 {
						createBulkProductVariantAttributesParam := make([]CreateBulkProductVariantAttributeParams, len(variant.Attributes))
						for i, attr := range variant.Attributes {
							createBulkProductVariantAttributesParam[i] = CreateBulkProductVariantAttributeParams{
								VariantID:        created.ID,
								AttributeValueID: uuid.MustParse(attr.ValueID),
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

type UpdateProdImagesTxArgs struct {
	ImageID    string      `json:"imageId"`
	EntityID   uuid.UUID   `json:"entityId"`
	EntityType string      `json:"entityType"`
	Role       *string     `json:"role"`
	VariantIDs []uuid.UUID `json:"variantIds"`
}

func updateProductImages(ctx context.Context, q *Queries, arg []UpdateProdImagesTxArgs) error {
	for _, image := range arg {
		if image.Role != nil {
			var defaultDisplayOrder int16 = 1
			err := q.UpdateProductImageAssignment(ctx, UpdateProductImageAssignmentParams{
				ImageID:      uuid.MustParse(image.ImageID),
				EntityID:     image.EntityID,
				EntityType:   image.EntityType,
				DisplayOrder: &defaultDisplayOrder,
				Role:         image.Role,
			})
			if err != nil {
				log.Error().Err(err).Msg("Failed to update product image assignment")
				return err
			}
		}
		// Remove all old image assignments
		err := q.DeleteImageAssignments(ctx, DeleteImageAssignmentsParams{
			ImageID:    uuid.MustParse(image.ImageID),
			EntityType: string(EntityTypeProductVariant),
		})

		if err != nil {
			log.Error().Err(err).Msg("Failed to delete image assignments")
			return err
		}

		// If there are no variant IDs, we can return early
		if len(image.VariantIDs) > 0 {
			// Create new image assignments
			createBulkImgAssignmentParams := make([]InsertBulkImageAssignmentsParams, 0)
			for _, variantID := range image.VariantIDs {
				createImgAssignmentParams := InsertBulkImageAssignmentsParams{
					ImageID:      uuid.MustParse(image.ImageID),
					EntityID:     variantID,
					Role:         "gallery",
					EntityType:   string(EntityTypeProductVariant),
					DisplayOrder: 1,
				}
				if image.Role != nil {
					createImgAssignmentParams.Role = *image.Role
				}
				createBulkImgAssignmentParams = append(createBulkImgAssignmentParams, createImgAssignmentParams)
			}
			_, err = q.InsertBulkImageAssignments(ctx, createBulkImgAssignmentParams)
			if err != nil {
				log.Error().Err(err).Msg("Failed to create bulk image assignments")
				return err
			}
		}
	}
	return nil
}
