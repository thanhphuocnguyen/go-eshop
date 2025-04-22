package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type UpdateImageAssignmentsTxParam struct {
	ImageID    int32       `json:"image_id"`
	VariantIDs []uuid.UUID `json:"variant_ids"`
}

func (s *pgRepo) UpdateImageAssignmentsTx(ctx context.Context, arg []UpdateImageAssignmentsTxParam) (err error) {
	err = s.execTx(ctx, func(q *Queries) (err error) {
		if len(arg) > 0 {
			for _, imageAssignment := range arg {
				// Remove all old image assignments
				err = q.DeleteImageAssignments(ctx, DeleteImageAssignmentsParams{
					ImageID:    imageAssignment.ImageID,
					EntityType: VariantEntityType,
				})

				if err != nil {
					log.Error().Err(err).Msg("Failed to delete image assignments")
					return
				}

				// If there are no variant IDs, we can return early
				if len(imageAssignment.VariantIDs) == 0 {
					log.Debug().Msg("No variant IDs provided, skipping image assignment creation")
					return
				}
				// Create new image assignments
				createBulkImgAssignmentParams := make([]CreateBulkImageAssignmentsParams, 0)
				for _, variantID := range imageAssignment.VariantIDs {
					createBulkImgAssignmentParams = append(createBulkImgAssignmentParams, CreateBulkImageAssignmentsParams{
						ImageID:      imageAssignment.ImageID,
						EntityID:     variantID,
						EntityType:   VariantEntityType,
						DisplayOrder: 1,
					})
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
