package repository

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

type CreateVariantTxParam struct {
	ProductID    int64
	VariantPrice float64
	Discount     int16
	VariantSku   *string
	VariantStock int32
	Attributes   []struct {
		AttributeID int32
		Value       string
	}
}

func (s *pgRepo) CreateVariantTx(ctx context.Context, arg CreateVariantTxParam) (ProductVariant, error) {
	var variant ProductVariant
	err := s.execTx(ctx, func(q *Queries) error {
		result, txErr := createVariantUtil(ctx, q, arg)
		if txErr != nil {
			return txErr
		}
		variant = *result
		return nil
	})

	return variant, err
}

func createVariantUtil(ctx context.Context, q *Queries, arg CreateVariantTxParam) (*ProductVariant, error) {
	createParam := CreateVariantParams{
		ProductID:     arg.ProductID,
		Price:         utils.GetPgNumericFromFloat(arg.VariantPrice),
		StockQuantity: arg.VariantStock,
		Discount:      arg.Discount,
	}

	if arg.VariantSku != nil {
		createParam.Sku = utils.GetPgTypeText(*arg.VariantSku)
	}

	variant, err := q.CreateVariant(ctx, createParam)
	if err != nil {
		log.Error().Err(err).Msg("CreateVariant")
		return nil, err
	}
	if len(arg.Attributes) > 0 {
		addVariantAttributeParam := make([]CreateBulkVariantAttributesParams, len(arg.Attributes))
		for i, param := range arg.Attributes {
			addVariantAttributeParam[i] = CreateBulkVariantAttributesParams{
				VariantID:   variant.VariantID,
				AttributeID: param.AttributeID,
				Value:       param.Value,
			}
		}
		_, err = q.CreateBulkVariantAttributes(ctx, addVariantAttributeParam)

		if err != nil {
			log.Error().Err(err).Msg("AddVariantAttributes")
			return nil, err
		}
	}

	return &variant, nil
}
