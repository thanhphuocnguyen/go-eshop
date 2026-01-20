package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

// CreateProductTx creates a product with its attributes, categories, and collections
// within a single database transaction to ensure consistency
func (repo *pgRepo) CreateProductTx(ctx context.Context, arg CreateProductTxArgs) (Product, error) {
	var result Product

	err := repo.execTx(ctx, func(q *Queries) error {
		var err error

		// Create the product first
		result, err = q.CreateProduct(ctx, arg.Product)
		if err != nil {
			log.Error().Err(err).Msg("CreateProduct failed in transaction")
			return err
		}

		// Create product attributes
		for _, attrID := range arg.Attributes {
			_, err = q.CreateProductAttribute(ctx, CreateProductAttributeParams{
				ProductID:   result.ID,
				AttributeID: attrID,
			})
			if err != nil {
				log.Error().Err(err).Msg("CreateProductAttribute failed in transaction")
				return err
			}
		}

		// Add product to categories
		if len(arg.CategoryIDs) > 0 {
			categoryParams := make([]AddProductsToCategoryParams, len(arg.CategoryIDs))
			for i, catID := range arg.CategoryIDs {
				categoryParams[i] = AddProductsToCategoryParams{
					CategoryID: uuid.MustParse(catID),
					ProductID:  result.ID,
				}
			}
			_, err = q.AddProductsToCategory(ctx, categoryParams)
			if err != nil {
				log.Error().Err(err).Msg("AddProductsToCategory failed in transaction")
				return err
			}
		}

		// Add product to collections
		if len(arg.CollectionIDs) > 0 {
			collectionParams := make([]AddProductsToCollectionParams, len(arg.CollectionIDs))
			for i, collID := range arg.CollectionIDs {
				collectionParams[i] = AddProductsToCollectionParams{
					CollectionID: uuid.MustParse(collID),
					ProductID:    result.ID,
				}
			}
			_, err = q.AddProductsToCollection(ctx, collectionParams)
			if err != nil {
				log.Error().Err(err).Msg("AddProductsToCollection failed in transaction")
				return err
			}
		}

		// Create product variants if provided
		if len(arg.Variants) > 0 {
			// Get attribute values to generate SKU
			for _, variantParam := range arg.Variants {
				attributeValues, err := q.GetAttributeValuesByIDs(ctx, variantParam.AttributeValues)
				if err != nil {
					log.Error().Err(err).Msg("GetAttributeValuesByIDs failed in transaction")
					return err
				}

				variantSku := GetVariantSKUWithAttributeNames(result.BaseSku, attributeValues)

				createVariantParams := CreateProductVariantParams{
					ProductID:   result.ID,
					Description: &variantParam.Description,
					Sku:         variantSku,
					Price:       utils.GetPgNumericFromFloat(variantParam.Price),
					Stock:       variantParam.Stock,
				}
				if variantParam.Weight != nil {
					createVariantParams.Weight = utils.GetPgNumericFromFloat(*variantParam.Weight)
				}

				variant, err := q.CreateProductVariant(ctx, createVariantParams)
				if err != nil {
					log.Error().Err(err).Msg("CreateProductVariant failed in transaction")
					return err
				}

				// Create variant attribute associations
				if len(variantParam.AttributeValues) > 0 {
					variantAttrParams := make([]CreateBulkProductVariantAttributeParams, len(variantParam.AttributeValues))
					for i, attrValID := range variantParam.AttributeValues {
						variantAttrParams[i] = CreateBulkProductVariantAttributeParams{
							VariantID:        variant.ID,
							AttributeValueID: attrValID,
						}
					}
					_, err = q.CreateBulkProductVariantAttribute(ctx, variantAttrParams)
					if err != nil {
						log.Error().Err(err).Msg("CreateBulkProductVariantAttribute failed in transaction")
						return err
					}
				}
			}
		}

		return nil
	})

	return result, err
}

// UpdateProductTx updates a product with its attributes, categories, and collections
// within a single database transaction to ensure consistency
func (repo *pgRepo) UpdateProductTx(ctx context.Context, arg UpdateProductTxArgs) (Product, error) {
	var result Product

	err := repo.execTx(ctx, func(q *Queries) error {
		var err error

		// Update the product first
		result, err = q.UpdateProduct(ctx, arg.Product)
		if err != nil {
			log.Error().Err(err).Msg("UpdateProduct failed in transaction")
			return err
		}

		// Update product attributes if provided
		if arg.Attributes != nil {
			// Delete existing attributes first
			err = q.DeleteProductAttributesByProductID(ctx, arg.Product.ID)
			if err != nil {
				log.Error().Err(err).Msg("DeleteProductAttributesByProductID failed in transaction")
				return err
			}

			// Create new attributes
			if len(*arg.Attributes) > 0 {
				prodAttrParams := make([]CreateBulkProductAttributesParams, len(*arg.Attributes))
				for i, attrID := range *arg.Attributes {
					prodAttrParams[i] = CreateBulkProductAttributesParams{
						ProductID:   arg.Product.ID,
						AttributeID: attrID,
					}
				}
				_, err = q.CreateBulkProductAttributes(ctx, prodAttrParams)
				if err != nil {
					log.Error().Err(err).Msg("CreateBulkProductAttributes failed in transaction")
					return err
				}
			}
		}

		// Update categories if provided
		if arg.CategoryIDs != nil {
			// Remove existing category associations
			err = q.RemoveProductsFromCategory(ctx, arg.Product.ID)
			if err != nil {
				log.Error().Err(err).Msg("RemoveProductsFromCategory failed in transaction")
				return err
			}

			// Add new category associations
			if len(*arg.CategoryIDs) > 0 {
				categoryParams := make([]AddProductsToCategoryParams, len(*arg.CategoryIDs))
				for i, catID := range *arg.CategoryIDs {
					categoryParams[i] = AddProductsToCategoryParams{
						CategoryID: uuid.MustParse(catID),
						ProductID:  arg.Product.ID,
					}
				}
				_, err = q.AddProductsToCategory(ctx, categoryParams)
				if err != nil {
					log.Error().Err(err).Msg("AddProductsToCategory failed in transaction")
					return err
				}
			}
		}

		// Update collections if provided
		if arg.CollectionIDs != nil {
			// Remove existing collection associations
			err = q.RemoveProductsFromCollection(ctx, arg.Product.ID)
			if err != nil {
				log.Error().Err(err).Msg("RemoveProductsFromCollection failed in transaction")
				return err
			}

			// Add new collection associations
			if len(*arg.CollectionIDs) > 0 {
				collectionParams := make([]AddProductsToCollectionParams, len(*arg.CollectionIDs))
				for i, collID := range *arg.CollectionIDs {
					collectionParams[i] = AddProductsToCollectionParams{
						CollectionID: uuid.MustParse(collID),
						ProductID:    arg.Product.ID,
					}
				}
				_, err = q.AddProductsToCollection(ctx, collectionParams)
				if err != nil {
					log.Error().Err(err).Msg("AddProductsToCollection failed in transaction")
					return err
				}
			}
		}

		return nil
	})

	return result, err
}

// CreateProductVariantTx creates a product variant with its attribute values
// within a single database transaction to ensure consistency
func (repo *pgRepo) CreateProductVariantTx(ctx context.Context, arg CreateProductVariantTxArgs) (ProductVariant, error) {
	var result ProductVariant

	err := repo.execTx(ctx, func(q *Queries) error {
		var err error

		// Create the product variant first
		result, err = q.CreateProductVariant(ctx, arg.Variant)
		if err != nil {
			log.Error().Err(err).Msg("CreateProductVariant failed in transaction")
			return err
		}

		// Create variant attribute associations
		if len(arg.AttributeValues) > 0 {
			variantAttrParams := make([]CreateBulkProductVariantAttributeParams, len(arg.AttributeValues))
			for i, attrValID := range arg.AttributeValues {
				variantAttrParams[i] = CreateBulkProductVariantAttributeParams{
					VariantID:        result.ID,
					AttributeValueID: attrValID,
				}
			}
			_, err = q.CreateBulkProductVariantAttribute(ctx, variantAttrParams)
			if err != nil {
				log.Error().Err(err).Msg("CreateBulkProductVariantAttribute failed in transaction")
				return err
			}
		}

		return nil
	})

	return result, err
}
