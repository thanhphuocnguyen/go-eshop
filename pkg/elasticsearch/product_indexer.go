package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/elastic/go-elasticsearch/v9/esutil"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
)

type ProductIndexer struct {
	esStore *ESStore
}

func NewProductIndexer(store *ESStore) *ProductIndexer {
	return &ProductIndexer{esStore: store}
}

func (pi *ProductIndexer) CreateIndex() error {
	// remove index first
	err := pi.esStore.DeleteIndex(PRODUCT_INDEX)
	if err != nil {
		return err
	}
	// read mapping from file
	mapping, err := readMappingFromFile("mappings/product_mapping.json")
	if err != nil {
		return err
	}
	// res, err := pi.esStore.CreateIndex()
	return pi.esStore.CreateIndex(PRODUCT_INDEX, mapping)
}

func readMappingFromFile(path string) (string, error) {
	absPath, err := filepath.Abs("pkg/elasticsearch/" + path)
	if err != nil {
		return "", err
	}

	fmt.Println(absPath)
	data, err := os.ReadFile(absPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (pi *ProductIndexer) IndexProduct(product repository.GetProductsForIndexingRow) error {
	return pi.esStore.IndexDocument(PRODUCT_INDEX, product.ID.String(), product)
}

func (pi *ProductIndexer) UpdateProduct(productID string, updatedProduct repository.GetProductsForIndexingRow) error {
	return pi.esStore.UpdateDocument(PRODUCT_INDEX, productID, updatedProduct)
}

func (pi *ProductIndexer) DeleteProduct(productID string) error {
	return pi.esStore.DeleteDocument(PRODUCT_INDEX, productID)
}

func (pi *ProductIndexer) BulkIndexProducts(ctx context.Context, products []repository.GetProductsForIndexingRow) (uint64, error) {
	bulkRequests := make([]esutil.BulkIndexerItem, 0, len(products))

	countSuccessful := uint64(0)
	for _, product := range products {
		productIndexBody := map[string]interface{}{
			"id":                product.ID,
			"name":              product.Name,
			"description":       product.Description,
			"short_description": product.ShortDescription,
			"is_active":         product.IsActive,
			"slug":              product.Slug,
			"sku":               product.BaseSku,
			"price":             product.BasePrice,
			"categories":        product.Categories,
			"brand":             product.Brand,
			"rating":            product.AvgRating,
			"in_stock":          true,
			"created_at":        product.CreatedAt,
			"attributes":        product.Attributes,
		}

		// Marshal the product index body to JSON
		body, _ := json.Marshal(productIndexBody)
		req := esutil.BulkIndexerItem{
			Action:     "index",
			Index:      PRODUCT_INDEX,
			DocumentID: product.ID.String(),
			Body:       bytes.NewReader(body),
			// OnSuccess is called for each successful operation
			OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
				atomic.AddUint64(&countSuccessful, 1)
			},

			// OnFailure is called for each failed operation
			OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
				if err != nil {
					log.Printf("ERROR: %s", err)
				} else {
					log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
				}
			},
		}
		bulkRequests = append(bulkRequests, req)
	}
	err := pi.esStore.BulkIndexDocuments(ctx, PRODUCT_INDEX, bulkRequests)

	return countSuccessful, err
}

func (pi *ProductIndexer) BulkDeleteProducts() error {
	return pi.esStore.DeleteIndex(PRODUCT_INDEX)
}

func (pi *ProductIndexer) GetMapping() (map[string]interface{}, error) {
	return pi.esStore.GetMapping(PRODUCT_INDEX)
}
