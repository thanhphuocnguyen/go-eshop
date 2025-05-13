package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

type CreateProductTxArgs struct {
	Name             string                         `json:"name" binding:"required,min=3,max=100"`
	Description      string                         `json:"description" binding:"omitempty,min=6,max=5000"`
	ShortDescription *string                        `json:"short_description" binding:"omitempty,max=2000"`
	Price            float64                        `json:"price" binding:"required,gt=0"`
	Sku              string                         `json:"sku" binding:"required"`
	Slug             string                         `json:"slug" binding:"omitempty"`
	CategoryID       string                         `json:"category_id,omitempty" binding:"omitempty,uuid"`
	BrandID          string                         `json:"brand_id,omitempty" binding:"omitempty,uuid"`
	CollectionID     *string                        `json:"collection_id,omitempty" binding:"omitempty,uuid"`
	Attributes       []string                       `json:"attributes" binding:"min=1"`
	Variants         []CreateProductVariantTxParams `json:"variants,min=1"`
}

func (s *pgRepo) CreateProductTx(ctx context.Context, arg CreateProductTxArgs) (productID uuid.UUID, err error) {
	err = s.execTx(ctx, func(q *Queries) (err error) {
		attributes := make([]uuid.UUID, len(arg.Attributes))
		for i, attr := range arg.Attributes {
			attributes[i] = uuid.MustParse(attr)
		}
		createParams := CreateProductParams{
			Name:        arg.Name,
			Attributes:  attributes,
			Description: arg.Description,
		}

		createParams.BasePrice = utils.GetPgNumericFromFloat(arg.Price)
		createParams.ShortDescription = arg.ShortDescription
		createParams.Slug = arg.Slug
		createParams.BaseSku = arg.Sku

		createParams.CategoryID = utils.GetPgTypeUUIDFromString(arg.CategoryID)

		createParams.BrandID = utils.GetPgTypeUUIDFromString(arg.BrandID)
		if arg.CollectionID != nil {
			createParams.CollectionID = utils.GetPgTypeUUIDFromString(*arg.CollectionID)
		}

		product, err := q.CreateProduct(ctx, createParams)
		if err != nil {
			log.Error().Err(err).Timestamp().Msg("CreateProduct")

			return err
		}

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

		for _, variant := range arg.Variants {
			sku := ""
			attrValues := make([]AttributeValue, 0)
			for _, attr := range variant.Attributes {
				attrValues = append(attrValues, attributeValueMp[attr.ValueID])
			}

			variantSku := GetVariantSKUWithAttributeNames(product.BaseSku, attrValues)

			sku = variantSku

			createVariantParams := CreateProductVariantParams{
				ProductID: product.ID,
				Sku:       sku,
			}

			createVariantParams.Stock = *&variant.Stock
			createVariantParams.Weight = utils.GetPgNumericFromFloat(*variant.Weight)
			createVariantParams.Price = utils.GetPgNumericFromFloat(variant.Price)

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

		return nil
	})

	return productID, err
}
