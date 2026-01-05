package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
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

func (pi *ProductIndexer) IndexProduct(product repository.Product) error {
	return pi.esStore.IndexDocument(PRODUCT_INDEX, product.ID.String(), product)
}

func (pi *ProductIndexer) UpdateProduct(productID string, updatedProduct repository.Product) error {
	return pi.esStore.UpdateDocument(PRODUCT_INDEX, productID, updatedProduct)
}

func (pi *ProductIndexer) DeleteProduct(productID string) error {
	return pi.esStore.DeleteDocument(PRODUCT_INDEX, productID)
}

func (pi *ProductIndexer) BulkIndexProducts(ctx context.Context, products []repository.Product) (uint64, error) {
	bulkRequests := make([]esutil.BulkIndexerItem, 0, len(products))
	countSuccessful := uint64(0)
	for _, product := range products {
		body, _ := json.Marshal(product)
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
