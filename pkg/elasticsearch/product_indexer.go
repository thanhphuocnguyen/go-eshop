package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/getmapping"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/getsettings"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/dynamicmapping"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
)

// ProductIndexDocument represents the document structure for Elasticsearch indexing
// It matches the product_mapping.json schema
type ProductIndexDocument struct {
	ID               string    `json:"id"`
	Name             string    `json:"name,omitempty"`
	Description      string    `json:"description,omitempty"`
	ShortDescription *string   `json:"shortDescription,omitempty"`
	IsActive         bool      `json:"isActive"`
	Slug             string    `json:"slug,omitempty"`
	Sku              string    `json:"sku,omitempty"`
	Price            float64   `json:"price,omitzero"`
	BasePrice        float64   `json:"basePrice,omitzero"`
	Brand            string    `json:"brand,omitempty"`
	ImageUrl         *string   `json:"imageUrl,omitempty"`
	AvgRating        *float64  `json:"avgRating,omitempty,omitzero"`
	InStock          bool      `json:"inStock"`
	RatingCount      int       `json:"ratingCount,omitempty,omitzero"`
	CreatedAt        time.Time `json:"createdAt,omitempty"`
	Attributes       []string  `json:"attributes,omitempty"`
	Categories       []string  `json:"categories,omitempty"`
	Collections      []string  `json:"collections,omitempty"`
}

type ProductIndexer struct {
	esStore *ESClient
}

func NewProductIndexer(store *ESClient) *ProductIndexer {
	return &ProductIndexer{esStore: store}
}

func (pi *ProductIndexer) CreateIndex(ctx context.Context) error {
	numberOfShard := "1"
	analyzer := types.CustomAnalyzer{
		Type:      "custom",
		Tokenizer: "standard",
		Filter:    []string{"lowercase", "asciifolding"},
	}
	settings := &types.IndexSettings{
		NumberOfShards: &numberOfShard,
		Analysis: &types.IndexSettingsAnalysis{
			Analyzer: map[string]types.Analyzer{
				"custom_analyzer": &analyzer,
			},
		},
	}
	nameAnalyzer := "custom_analyzer"
	mappings := &types.TypeMapping{
		Dynamic: &dynamicmapping.Strict,
		Properties: map[string]types.Property{
			"id": types.NewKeywordProperty(),
			"name": types.TextProperty{
				Analyzer: &nameAnalyzer,
				Fields: map[string]types.Property{
					"keyword": types.NewKeywordProperty(),
				},
			},
			"description":      types.NewTextProperty(),
			"shortDescription": types.NewTextProperty(),
			"isActive":         types.NewBooleanProperty(),
			"slug":             types.NewKeywordProperty(),
			"sku":              types.NewKeywordProperty(),
			"price":            types.NewDoubleNumberProperty(),
			"imageUrl":         types.NewKeywordProperty(),
			"ratingCount":      types.NewIntegerNumberProperty(),

			"avgRating": types.NewFloatNumberProperty(),
			"basePrice": types.NewDoubleNumberProperty(),
			"inStock":   types.NewBooleanProperty(),
			"createdAt": types.NewDateProperty(),
			"categories": types.TextProperty{
				Analyzer: &nameAnalyzer,
				Fields: map[string]types.Property{
					"keyword": types.NewKeywordProperty(),
				},
			},
			"collections": types.TextProperty{
				Analyzer: &nameAnalyzer,
				Fields: map[string]types.Property{
					"keyword": types.NewKeywordProperty(),
				},
			},
			"brand": types.TextProperty{
				Analyzer: &nameAnalyzer,
				Fields: map[string]types.Property{
					"keyword": types.NewKeywordProperty(),
				},
			},
			"attributes": types.TextProperty{
				Analyzer: &nameAnalyzer,
				Fields: map[string]types.Property{
					"keyword": types.NewKeywordProperty(),
				},
			},
		},
	}
	// res, err := pi.esStore.CreateIndex()
	return pi.esStore.CreateIndex(ctx, PRODUCT_INDEX, mappings, settings)
}

// func readMappingFromFile(path string) (string, error) {
// 	absPath, err := filepath.Abs("pkg/elasticsearch/" + path)
// 	if err != nil {
// 		return "", err
// 	}

// 	fmt.Println(absPath)
// 	data, err := os.ReadFile(absPath)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(data), nil
// }

func (pi *ProductIndexer) IndexProduct(ctx context.Context, product repository.GetProductsForIndexingRow) error {
	indexDoc := pi.convertToIndexDocument(product)
	return pi.esStore.IndexDocument(ctx, PRODUCT_INDEX, product.ID.String(), indexDoc)
}

func (pi *ProductIndexer) UpdateProduct(ctx context.Context, productID string, updatedProduct ProductIndexDocument) error {
	jsonData, err := json.Marshal(updatedProduct)
	if err != nil {
		return err
	}
	return pi.esStore.UpdateDocument(ctx, PRODUCT_INDEX, productID, jsonData)
}

func (pi *ProductIndexer) DeleteProduct(ctx context.Context, productID string) error {
	return pi.esStore.DeleteDocument(ctx, PRODUCT_INDEX, productID)
}

// convertToIndexDocument converts database row to Elasticsearch index document
func (pi *ProductIndexer) convertToIndexDocument(row repository.GetProductsForIndexingRow) ProductIndexDocument {
	doc := ProductIndexDocument{
		ID:               row.ID.String(),
		Name:             row.Name,
		Description:      row.Description,
		ShortDescription: row.ShortDescription,
		IsActive:         row.IsActive != nil && *row.IsActive,
		Slug:             row.Slug,
		Sku:              row.BaseSku,
		Categories:       row.Categories,
		Brand:            row.Brand,
		InStock:          row.TotalStock != nil && *row.TotalStock > 0,
		CreatedAt:        row.CreatedAt,
		RatingCount:      int(row.RatingCount),
		ImageUrl:         row.ImageUrl,
		Collections:      row.Collections,
		Attributes:       row.AttributeValues,
	}

	// Convert price from pgtype.Numeric to float64
	if row.BasePrice.Valid {
		price, _ := row.BasePrice.Float64Value()
		doc.Price = price.Float64
	}
	if row.MinPrice.Valid {
		price, _ := row.MinPrice.Float64Value()
		doc.Price = price.Float64
	}
	// Convert rating from pgtype.Numeric to float64
	if row.AvgRating.Valid {
		rating, _ := row.AvgRating.Float64Value()
		doc.AvgRating = &rating.Float64
	}

	// Parse attributes from JSON bytes

	return doc
}

func (pi *ProductIndexer) BulkIndexProducts(ctx context.Context, products []repository.GetProductsForIndexingRow) (uint64, error) {
	bulkRequests := make([]esutil.BulkIndexerItem, 0, len(products))

	countSuccessful := uint64(0)
	for _, product := range products {
		// Convert to index document format
		indexDoc := pi.convertToIndexDocument(product)
		// Marshal the product index document to JSON
		body, err := json.Marshal(indexDoc)
		if err != nil {
			log.Printf("ERROR: Failed to marshal product %s: %v", product.ID.String(), err)
			continue
		}

		req := esutil.BulkIndexerItem{
			Action:     "index",
			Index:      PRODUCT_INDEX,
			DocumentID: product.ID.String(),
			Body:       bytes.NewReader(body),
			// OnSuccess is called for each successful operation
			OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
				atomic.AddUint64(&countSuccessful, 1)
				fmt.Printf("[%d] %s test/%s \n", res.Status, res.Result, item.DocumentID)
			},
			// OnFailure is called for each failed operation
			OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
				if err != nil {
					log.Printf("ERROR: Failed to index product %s: %v", item.DocumentID, err)
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

func (pi *ProductIndexer) BulkDeleteProducts(ctx context.Context) error {
	return pi.esStore.DeleteIndex(ctx, PRODUCT_INDEX)
}

func (pi *ProductIndexer) GetMapping(ctx context.Context) (getmapping.Response, error) {
	return pi.esStore.GetMapping(ctx, PRODUCT_INDEX)
}
func (pi *ProductIndexer) GetSettings(ctx context.Context) (getsettings.Response, error) {
	return pi.esStore.GetSetting(ctx, PRODUCT_INDEX)
}

func (pi *ProductIndexer) IndexExists(ctx context.Context) (bool, error) {
	return pi.esStore.IndexExists(ctx, PRODUCT_INDEX)
}

func (pi *ProductIndexer) GetProducts(ctx context.Context) ([]interface{}, error) {
	sz := 1000
	query := &search.Request{
		Size: &sz,
		Query: &types.Query{
			MatchAll: types.NewMatchAllQuery(),
		},
	}
	res, err := pi.esStore.QueryDocuments(ctx, PRODUCT_INDEX, query)
	if err != nil {
		return nil, err
	}
	fmt.Println(len(res))
	return res, nil
}

func (pi *ProductIndexer) SearchProducts(ctx context.Context, query *search.Request) ([]interface{}, error) {

	res, err := pi.esStore.QueryDocuments(ctx, PRODUCT_INDEX, query)
	if err != nil {
		return nil, err
	}
	fmt.Println(len(res))
	return res, nil
}
