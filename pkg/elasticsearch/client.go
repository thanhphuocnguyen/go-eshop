package elasticsearch

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esutil"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/update"
	"github.com/elastic/go-elasticsearch/v9/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v9/typedapi/indices/getmapping"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
)

type ESStore struct {
	es *elasticsearch.TypedClient
}

func NewClient(url string) (*ESStore, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			url,
		},
		RetryOnStatus: []int{http.StatusTooManyRequests},

		Transport: &http.Transport{
			TLSHandshakeTimeout:   10 * time.Second,
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: 30 * time.Second,
			DialContext:           (&net.Dialer{Timeout: 10 * time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12},
		},
	}

	es, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating the client: %s", err)
	}

	return &ESStore{es: es}, nil
}

func (c *ESStore) Ping(ctx context.Context) error {
	res, err := c.es.Ping().Do(ctx)
	if err != nil {
		return fmt.Errorf("error pinging elasticsearch: %s", err)
	}

	if !res {
		return fmt.Errorf("error pinging elasticsearch: no response")
	}

	return nil
}

func (c *ESStore) CreateIndex(ctx context.Context, index string, mapping *types.TypeMapping, setting *types.IndexSettings) error {
	res, err := c.es.Indices.Create(index).Request(&create.Request{
		Mappings: mapping,
		Settings: setting,
	}).Do(ctx)
	if err != nil {
		return fmt.Errorf("error creating index: %s", err)
	}

	if !res.Acknowledged {
		return fmt.Errorf("error creating index: not acknowledged")
	}

	return nil
}

func (c *ESStore) IndexExists(ctx context.Context, index string) (bool, error) {
	res, err := c.es.Indices.Exists(index).Do(ctx)
	if err != nil {
		return false, fmt.Errorf("error checking if index exists: %s", err)
	}

	return res, nil
}

func (c *ESStore) IndexDocument(ctx context.Context, index string, documentID string, document interface{}) error {

	_, err := c.es.Index(index).Request(document).Id(documentID).Do(ctx)
	if err != nil {
		return fmt.Errorf("error indexing document: %s", err)
	}

	return nil
}

func (c *ESStore) QueryDocuments(ctx context.Context, index string, query *search.Request) ([]interface{}, error) {

	res, err := c.es.Search().Index(index).Request(query).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("error searching documents: %s", err)
	}

	var documents []interface{}
	for _, hit := range res.Hits.Hits {
		var doc interface{}
		if err := json.Unmarshal(hit.Source_, &doc); err != nil {
			log.Printf("error unmarshaling document: %s", err)
			continue
		}
		documents = append(documents, doc)
	}

	return documents, nil
}

func (c *ESStore) UpdateDocument(ctx context.Context, index string, documentID string, updateDoc json.RawMessage) error {

	_, err := c.es.Update(index, documentID).Request(&update.Request{
		Upsert: updateDoc,
		Doc:    updateDoc,
	}).Do(ctx)
	if err != nil {
		return fmt.Errorf("error updating document: %s", err)
	}

	return nil
}

func (c *ESStore) DeleteDocument(ctx context.Context, index string, documentID string) error {
	_, err := c.es.Delete(index, documentID).Do(ctx)
	if err != nil {
		return fmt.Errorf("error deleting document: %s", err)
	}

	return nil
}

func (c *ESStore) BulkIndexDocuments(ctx context.Context, index string, request []esutil.BulkIndexerItem) error {
	bulkIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:     c.es,
		Index:      index,
		NumWorkers: 4,
	})

	if err != nil {
		return fmt.Errorf("error creating the bulk indexer: %s", err)
	}

	for _, doc := range request {
		err = bulkIndexer.Add(ctx, doc)
		if err != nil {
			return fmt.Errorf("error adding document to the bulk indexer: %s", err)
		}
	}

	if err := bulkIndexer.Close(ctx); err != nil {
		return fmt.Errorf("error closing the bulk indexer: %s", err)
	}
	return nil
}

func (c *ESStore) DeleteIndex(ctx context.Context, index string) error {
	res, err := c.es.Indices.Exists(index).Do(ctx)
	if res {
		// Index does not exist, nothing to delete
		return nil
	}
	delRs, err := c.es.Indices.Delete(index).Do(ctx)
	if err != nil {
		return fmt.Errorf("error deleting index: %s", err)
	}

	if !delRs.Acknowledged {
		return fmt.Errorf("error deleting index: %s", "not acknowledged")
	}

	return nil
}

func (c *ESStore) CreateMapping(ctx context.Context, index string, mapping map[string]types.Property) error {
	_, err := c.es.Indices.PutMapping(index).Properties(mapping).Do(ctx)
	if err != nil {
		return fmt.Errorf("error creating mapping: %s", err)
	}

	return nil
}

func (c *ESStore) GetMapping(ctx context.Context, index string) (getmapping.Response, error) {
	res, err := c.es.Indices.GetMapping().Index(index).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting mapping: %s", err)
	}
	return res, nil
}
