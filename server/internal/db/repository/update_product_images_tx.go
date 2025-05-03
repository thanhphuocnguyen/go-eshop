package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type UpdateProdImagesTxArgs struct {
	ImageID    int32       `json:"image_id"`
	EntityID   uuid.UUID   `json:"entity_id"`
	EntityType string      `json:"entity_type"`
	Role       *string     `json:"role"`
	VariantIDs []uuid.UUID `json:"variant_ids"`
}

func (s *pgRepo) UpdateProductImagesTx(ctx context.Context, arg []UpdateProdImagesTxArgs) (err error) {
	err = s.execTx(ctx, func(q *Queries) (err error) {
		for _, image := range arg {
			if image.Role != nil {
				var defaultDisplayOrder int16 = 1
				err = q.UpdateProductImageAssignment(ctx, UpdateProductImageAssignmentParams{
					ImageID:      image.ImageID,
					EntityID:     image.EntityID,
					EntityType:   image.EntityType,
					DisplayOrder: &defaultDisplayOrder,
					Role:         image.Role,
				})
				if err != nil {
					log.Error().Err(err).Msg("Failed to update product image assignment")
					return
				}
			}
			// Remove all old image assignments
			err = q.DeleteImageAssignments(ctx, DeleteImageAssignmentsParams{
				ImageID:    image.ImageID,
				EntityType: VariantEntityType,
			})

			if err != nil {
				log.Error().Err(err).Msg("Failed to delete image assignments")
				return
			}

			// If there are no variant IDs, we can return early
			if len(image.VariantIDs) > 0 {
				// Create new image assignments
				createBulkImgAssignmentParams := make([]CreateBulkImageAssignmentsParams, 0)
				for _, variantID := range image.VariantIDs {
					createImgAssignmentParams := CreateBulkImageAssignmentsParams{
						ImageID:      image.ImageID,
						EntityID:     variantID,
						EntityType:   VariantEntityType,
						DisplayOrder: 1,
					}
					if image.Role != nil {
						createImgAssignmentParams.Role = *image.Role
					}
					createBulkImgAssignmentParams = append(createBulkImgAssignmentParams, createImgAssignmentParams)
				}
				_, err = q.CreateBulkImageAssignments(ctx, createBulkImgAssignmentParams)
				if err != nil {
					log.Error().Err(err).Msg("Failed to create bulk image assignments")
					return
				}
			}

		}

		return
	})
	return
}
