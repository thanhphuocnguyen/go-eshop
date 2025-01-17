package repository

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/util"
)

type CreateVariantTxParam struct {
	ProductID    int64   `json:"product_id" validate:"required"`
	VariantName  string  `json:"variant_name" validate:"required"`
	VariantPrice float64 `json:"variant_price" validate:"required"`
	VariantSku   *string `json:"variant_sku" validate:"omitempty"`
	VariantStock int32   `json:"variant_stock" validate:"required"`
	Attributes   []int32 `json:"attributes" validate:"omitempty"`
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
		ProductID:    arg.ProductID,
		VariantName:  arg.VariantName,
		VariantPrice: util.GetPgNumericFromFloat(arg.VariantPrice),
		VariantStock: arg.VariantStock,
	}

	if arg.VariantSku != nil {
		createParam.VariantSku = util.GetPgTypeText(*arg.VariantSku)
	}

	variant, err := q.CreateVariant(ctx, createParam)
	if err != nil {
		log.Error().Err(err).Msg("CreateVariant")
		return nil, err
	}
	if len(arg.Attributes) > 0 {
		addVariantAttributeParam := make([]CreateBulkVariantAttributeParams, len(arg.Attributes))
		for i, attrID := range arg.Attributes {
			addVariantAttributeParam[i] = CreateBulkVariantAttributeParams{
				VariantID:        variant.VariantID,
				AttributeValueID: attrID,
			}
		}
		_, err = q.CreateBulkVariantAttribute(ctx, addVariantAttributeParam)
		if err != nil {
			log.Error().Err(err).Msg("AddVariantAttributes")
			return nil, err
		}
	}
	return &variant, nil
}
