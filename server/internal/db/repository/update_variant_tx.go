package repository

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/util"
)

type UpdateVariantTxParams struct {
	VariantID  int64    `json:"variant_id" binding:"required,min=1"`
	Price      *float64 `json:"price" binding:"omitempty,min=0"`
	Name       *string  `json:"name" binding:"omitempty,min=1"`
	SKU        *string  `json:"sku" binding:"omitempty,min=1"`
	Stock      *int32   `json:"stock" binding:"omitempty,min=0"`
	Attributes []int32  `json:"attributes" binding:"omitempty,dive,min=1"`
}

type UpdateVariantTxResult struct {
	Variant ProductVariant `json:"variant"`
}

func (repo *pgRepo) UpdateVariantTx(ctx context.Context, arg UpdateVariantTxParams) (UpdateVariantTxResult, error) {
	var result UpdateVariantTxResult
	err := repo.execTx(ctx, func(q *Queries) error {
		updated, err := updateVariantUtil(ctx, q, arg)
		if err != nil {
			return err
		}
		result.Variant = *updated
		return nil
	})

	return result, err
}

func updateVariantUtil(c context.Context, q *Queries, params UpdateVariantTxParams) (*ProductVariant, error) {
	var variant ProductVariant
	var updateParams UpdateVariantParams = UpdateVariantParams{
		VariantID: params.VariantID,
	}

	if params.Name != nil {
		updateParams.VariantName = util.GetPgTypeText(*params.Name)
	}
	if params.SKU != nil {
		updateParams.VariantSku = util.GetPgTypeText(*params.SKU)
	}
	if params.Price != nil {
		updateParams.VariantPrice = util.GetPgNumericFromFloat(*params.Price)
	}
	if params.Stock != nil {
		updateParams.VariantStock = util.GetPgTypeInt4(*params.Stock)
	}
	if len(params.Attributes) > 0 {
		attributes, err := q.GetVariantAttributes(c, params.VariantID)
		if err != nil {
			log.Error().Err(err).Msg("GetVariantAttributes")
			return nil, err
		}
		existingAttributes := make(map[int32]bool)
		newAttributes := make([]int32, 0)
		for _, attr := range attributes {
			existingAttributes[attr.AttributeValueID] = false
		}

		for _, attr := range params.Attributes {
			if _, ok := existingAttributes[attr]; ok {
				existingAttributes[attr] = true
			} else {
				newAttributes = append(newAttributes, attr)
			}
		}

		for attr, exists := range existingAttributes {
			if !exists {
				err = q.DeleteVariantAttribute(c, DeleteVariantAttributeParams{
					AttributeValueID: attr,
					VariantID:        params.VariantID,
				})
				if err != nil {
					log.Error().Err(err).Msg("DeleteVariantAttribute")
					return nil, err
				}
			}
		}

		for _, attr := range newAttributes {
			_, err := q.CreateVariantAttribute(c, CreateVariantAttributeParams{
				VariantID:        variant.VariantID,
				AttributeValueID: attr,
			})
			if err != nil {
				log.Error().Err(err).Msg("CreateVariantAttribute")
				return nil, err
			}
		}
	}

	updated, err := q.UpdateVariant(c, updateParams)
	if err != nil {
		log.Error().Err(err).Msg("UpdateVariant")
		return nil, err
	}
	return &updated, nil
}
