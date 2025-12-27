package elasticsearch

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	elasticSearchV8 "github.com/elastic/go-elasticsearch/v8"
)

type Client struct {
	es *elasticSearchV8.Client
}

func NewClient(url string) (*Client, error) {
	cfg := elasticSearchV8.Config{
		Addresses: []string{
			url,
		},
		RetryOnStatus: []int{http.StatusTooManyRequests},
		Transport:     &http.Transport{TLSHandshakeTimeout: 10 * time.Second},
	}

	es, err := elasticSearchV8.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating the client: %s", err)
	}

	return &Client{es: es}, nil
}

func (c *Client) CreateIndex(index string, mapping interface{}) error {
	mappingReader, err := ToReader(mapping)
	if err != nil {
		return fmt.Errorf("error converting mapping to reader: %s", err)
	}

	res, err := c.es.Indices.Create(index, c.es.Indices.Create.WithBody(mappingReader))
	if err != nil {
		return fmt.Errorf("error creating index: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error creating index: %s", res.String())
	}

	return nil
}

func (c *Client) IndexDocument(index string, documentID string, document interface{}) error {
	documentReader, err := ToReader(document)
	if err != nil {
		return fmt.Errorf("error converting document to reader: %s", err)
	}

	res, err := c.es.Index(index, documentReader, c.es.Index.WithDocumentID(documentID))
	if err != nil {
		return fmt.Errorf("error indexing document: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.String())
	}

	return nil
}

func (c *Client) QueryDocuments(index string, query interface{}) ([]interface{}, error) {
	queryReader, err := ToReader(query)
	if err != nil {
		return nil, fmt.Errorf("error converting query to reader: %s", err)
	}

	res, err := c.es.Search(
		c.es.Search.WithIndex(index),
		c.es.Search.WithBody(queryReader),
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

func (c *Client) UpdateDocument(index string, documentID string, update interface{}) error {
	updateReader, err := ToReader(update)
	if err != nil {
		return fmt.Errorf("error converting update to reader: %s", err)
	}

	res, err := c.es.Update(index, documentID, updateReader)
	if err != nil {
		return fmt.Errorf("error updating document: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error updating document: %s", res.String())
	}

	return nil
}

func (c *Client) DeleteDocument(index string, documentID string) error {
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
