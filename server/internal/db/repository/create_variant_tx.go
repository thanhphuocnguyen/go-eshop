package repository

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/util"
)

type CreateVariantTxParam struct {
	ProductID    int64
	VariantName  string
	VariantPrice float64
	VariantStock int32
	Attributes   []int32
}

type CreateVariantTxResult struct {
	Variant ProductVariant
}

func (s *pgRepo) CreateVariantTx(ctx context.Context, arg CreateVariantTxParam) (CreateVariantTxResult, error) {
	var result CreateVariantTxResult
	err := s.execTx(ctx, func(q *Queries) error {
		createParam := CreateVariantParams{
			ProductID:    arg.ProductID,
			VariantName:  arg.VariantName,
			VariantPrice: util.GetPgNumericFromFloat(arg.VariantPrice),
			VariantStock: arg.VariantStock,
		}
		variant, err := q.CreateVariant(ctx, createParam)
		if err != nil {
			log.Error().Err(err).Msg("CreateVariant")
			return err
		}
		result.Variant = variant
		if len(arg.Attributes) > 0 {
			addVariantAttributeParam := make([]AddVariantAttributesParams, len(arg.Attributes))
			for i, attrID := range arg.Attributes {
				addVariantAttributeParam[i] = AddVariantAttributesParams{
					VariantID:        util.GetPgTypeInt8(variant.VariantID),
					AttributeValueID: attrID,
				}
			}
			_, err = q.AddVariantAttributes(ctx, addVariantAttributeParam)
			if err != nil {
				log.Error().Err(err).Msg("AddVariantAttributes")
				return err
			}
		}
		return nil
	})
	return result, err
}
