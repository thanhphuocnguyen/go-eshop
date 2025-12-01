package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// CreateProductTx creates a product with its attributes, categories, and collections
// within a single database transaction to ensure consistency
func (store *pgRepo) CreateProductTx(ctx context.Context, arg CreateProductTxArgs) (Product, error) {
	var result Product

	err := store.execTx(ctx, func(q *Queries) error {
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

		return nil
	})

	return result, err
}

// UpdateProductTx updates a product with its attributes, categories, and collections
// within a single database transaction to ensure consistency
func (store *pgRepo) UpdateProductTx(ctx context.Context, arg UpdateProductTxArgs) (Product, error) {
	var result Product

	err := store.execTx(ctx, func(q *Queries) error {
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
