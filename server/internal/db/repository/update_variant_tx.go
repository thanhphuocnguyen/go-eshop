package repository

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/util"
)

type UpdateVariantTxParams struct {
	VariantID  int64                          `json:"variant_id" binding:"required,min=1"`
	Price      *float64                       `json:"price" binding:"omitempty,min=0"`
	Name       *string                        `json:"name" binding:"omitempty,min=1"`
	Sku        *string                        `json:"sku" binding:"omitempty,min=1"`
	Stock      *int32                         `json:"stock" binding:"omitempty,min=0"`
	Attributes []UpdateVariantAttributeParams `json:"attributes" binding:"omitempty,dive,min=1"`
}

type UpdateVariantTxResult struct {
	Variant    ProductVariant     `json:"variant"`
	Attributes []VariantAttribute `json:"attributes"`
}

func (repo *pgRepo) UpdateVariantTx(ctx context.Context, arg UpdateVariantTxParams) (UpdateVariantTxResult, error) {
	var result UpdateVariantTxResult
	err := repo.execTx(ctx, func(q *Queries) error {
		updated, variantsUpdated, err := updateVariantUtil(ctx, q, arg)
		if err != nil {
			return err
		}
		result.Variant = updated
		result.Attributes = variantsUpdated
		return nil
	})

	return result, err
}

func updateVariantUtil(c context.Context, q *Queries, params UpdateVariantTxParams) (ProductVariant, []VariantAttribute, error) {
	var attributes []VariantAttribute = make([]VariantAttribute, 0)
	var updateParams UpdateVariantParams = UpdateVariantParams{
		VariantID: params.VariantID,
	}

	if params.Sku != nil {
		updateParams.Sku = util.GetPgTypeText(*params.Sku)
	}
	if params.Price != nil {
		updateParams.Price = util.GetPgNumericFromFloat(*params.Price)
	}
	if params.Stock != nil {
		updateParams.StockQuantity = util.GetPgTypeInt4(*params.Stock)
	}

	if len(params.Attributes) > 0 {
		for _, attr := range params.Attributes {
			updatedVariantAttribute, err := q.UpdateVariantAttribute(c, UpdateVariantAttributeParams{
				VariantAttributeID: attr.VariantAttributeID,
				Value:              attr.Value,
			})
			if err != nil {
				log.Error().Err(err).Msg("UpdateVariantAttribute")
				return ProductVariant{}, nil, err
			}
			attributes = append(attributes, updatedVariantAttribute)
		}

	}

	updated, err := q.UpdateVariant(c, updateParams)
	if err != nil {
		log.Error().Err(err).Msg("UpdateVariant")
		return updated, nil, err
	}
	return updated, nil, nil
}
