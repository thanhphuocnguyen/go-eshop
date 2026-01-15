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
)

type ESStore struct {
	es *elasticsearch.Client
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

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating the client: %s", err)
	}

	return &ESStore{es: es}, nil
}

func (c *ESStore) Ping() error {
	res, err := c.es.Ping()
	if err != nil {
		return fmt.Errorf("error pinging elasticsearch: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error pinging elasticsearch: %s", res.String())
	}

	return nil
}

func (c *ESStore) CreateIndex(index string, mapping interface{}) error {
	res, err := c.es.Indices.Create(index, c.es.Indices.Create.WithBody(esutil.NewJSONReader(&mapping)))
	if err != nil {
		return fmt.Errorf("error creating index: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error creating index: %s", res.String())
	}

	return nil
}

func (c *ESStore) IndexExists(index string) (bool, error) {
	res, err := c.es.Indices.Exists([]string{index})
	if err != nil {
		return false, fmt.Errorf("error checking if index exists: %s", err)
	}
	defer res.Body.Close()

	return res.StatusCode == 200, nil
}

func (c *ESStore) IndexDocument(index string, documentID string, document interface{}) error {

	res, err := c.es.Index(index, esutil.NewJSONReader(&document), c.es.Index.WithDocumentID(documentID))
	if err != nil {
		return fmt.Errorf("error indexing document: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.String())
	}

	return nil
}

func (c *ESStore) QueryDocuments(index string, query interface{}) ([]interface{}, error) {

	res, err := c.es.Search(
		c.es.Search.WithIndex(index),
		c.es.Search.WithBody(esutil.NewJSONReader(&query)),
	)
	if err != nil {
		return nil, fmt.Errorf("error searching documents: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error searching documents: %s", res.String())
	}

	var result struct {
		Hits struct {
			Hits []struct {
				Source json.RawMessage `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %s", err)
	}

	var documents []interface{}
	for _, hit := range result.Hits.Hits {
		var doc interface{}
		if err := json.Unmarshal(hit.Source, &doc); err != nil {
			log.Printf("error unmarshaling document: %s", err)
			continue
		}
		documents = append(documents, doc)
	}

	return documents, nil
}

func (c *ESStore) UpdateDocument(index string, documentID string, update interface{}) error {

	res, err := c.es.Update(index, documentID, esutil.NewJSONReader(&update))
	if err != nil {
		return fmt.Errorf("error updating document: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error updating document: %s", res.String())
	}

	return nil
}

func (c *ESStore) DeleteDocument(index string, documentID string) error {
	res, err := c.es.Delete(index, documentID)
	if err != nil {
		return fmt.Errorf("error deleting document: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error deleting document: %s", res.String())
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

func (c *ESStore) DeleteIndex(index string) error {
	res, err := c.es.Indices.Delete([]string{index})
	if err != nil {
		return fmt.Errorf("error deleting index: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error deleting index: %s", res.String())
	}

	return nil
}

func (c *ESStore) CreateMapping(index string, mapping interface{}) error {
	res, err := c.es.Indices.PutMapping([]string{index}, esutil.NewJSONReader(&mapping))
	if err != nil {
		return fmt.Errorf("error creating mapping: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error creating mapping: %s", res.String())
	}
	return nil
}
