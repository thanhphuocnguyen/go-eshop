package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/getmapping"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/getsettings"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type ESClient struct {
	es *elasticsearch.TypedClient
}

func NewClient(url string, username string, password string) (*ESClient, error) {
	// customTransport := &http.Transport{
	// 	TLSClientConfig: &tls.Config{
	// 		InsecureSkipVerify: true, // Bypass certificate verification
	// 	},
	// }

	// 2. Create a custom http.Client using the custom Transport
	// httpClient := &http.Client{
	// 	Transport: customTransport,
	// 	Timeout:   time.Second * 30, // Set an appropriate timeout
	// }
	cfg := elasticsearch.Config{
		Addresses: []string{
			url,
		},
		// RetryOnStatus: []int{http.StatusTooManyRequests},
		// Transport:     httpClient.Transport,
	}

	es, err := elasticsearch.NewTypedClient(cfg)

	if err != nil {
		return nil, fmt.Errorf("error creating the client: %s", err)
	}

	return &ESClient{es: es}, nil
}

func (c *ESClient) Ping(ctx context.Context) error {
	res, err := c.es.Ping().Do(ctx)
	if err != nil {
		return err
	}

	if !res {
		return fmt.Errorf("error pinging elasticsearch: no response")
	}

	return nil
}

func (c *ESClient) CreateIndex(ctx context.Context, index string, mapping *types.TypeMapping, setting *types.IndexSettings) error {
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

func (c *ESClient) IndexExists(ctx context.Context, index string) (bool, error) {
	res, err := c.es.Indices.Exists(index).Do(ctx)
	if err != nil {
		return false, fmt.Errorf("error checking if index exists: %s", err)
	}

	return res, nil
}

func (c *ESClient) IndexDocument(ctx context.Context, index string, documentID string, document interface{}) error {

	_, err := c.es.Index(index).Request(document).Id(documentID).Do(ctx)
	if err != nil {
		return fmt.Errorf("error indexing document: %s", err)
	}

	return nil
}

func (c *ESClient) QueryDocuments(ctx context.Context, index string, query *search.Request) ([]json.RawMessage, error) {

	res, err := c.es.Search().Index(index).Request(query).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("error searching documents: %s", err)
	}

	var documents []json.RawMessage
	for _, hit := range res.Hits.Hits {
		documents = append(documents, hit.Source_)
	}

	return documents, nil
}

func (c *ESClient) UpdateDocument(ctx context.Context, index string, documentID string, updateDoc json.RawMessage) error {

	_, err := c.es.Update(index, documentID).Request(&update.Request{
		Upsert: updateDoc,
		Doc:    updateDoc,
	}).Do(ctx)
	if err != nil {
		return fmt.Errorf("error updating document: %s", err)
	}

	return nil
}

func (c *ESClient) DeleteDocument(ctx context.Context, index string, documentID string) error {
	_, err := c.es.Delete(index, documentID).Do(ctx)
	if err != nil {
		return fmt.Errorf("error deleting document: %s", err)
	}

	return nil
}

func (c *ESClient) BulkIndexDocuments(ctx context.Context, index string, request []esutil.BulkIndexerItem) error {
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

func (c *ESClient) DeleteIndex(ctx context.Context, index string) error {
	res, err := c.es.Indices.Exists(index).Do(ctx)
	if !res {
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

func (c *ESClient) CreateMapping(ctx context.Context, index string, mapping map[string]types.Property) error {
	_, err := c.es.Indices.PutMapping(index).Properties(mapping).Do(ctx)
	if err != nil {
		return fmt.Errorf("error creating mapping: %s", err)
	}

	return nil
}

func (c *ESClient) GetMapping(ctx context.Context, index string) (getmapping.Response, error) {
	res, err := c.es.Indices.GetMapping().Index(index).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting mapping: %s", err)
	}
	return res, nil
}

func (c *ESClient) GetSetting(ctx context.Context, index string) (getsettings.Response, error) {
	res, err := c.es.Indices.GetSettings().Index(index).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting mapping: %s", err)
	}
	return res, nil
}

func (c *ESClient) Close(ctx context.Context) error {
	// The elasticsearch-go client does not have a Close method,
	// but if you had any resources to clean up, you would do it here.
	c.es.Close(ctx)
	return nil
}
