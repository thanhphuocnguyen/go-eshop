package elasticsearch

import (
	"fmt"
	"io"
)

// ExampleUsage demonstrates how to use the ToReader utility functions
func ExampleUsage() {
	// Example 1: Converting a search query struct to io.Reader
	searchQuery := struct {
		Query struct {
			Match struct {
				Title string `json:"title"`
			} `json:"match"`
		} `json:"query"`
		Size int `json:"size"`
	}{
		Size: 10,
	}
	searchQuery.Query.Match.Title = "golang"

	reader, err := ToReader(searchQuery)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Search query as io.Reader created successfully\n")
	_ = reader

	// Example 2: Converting a map to io.Reader
	bulkData := map[string]interface{}{
		"index": map[string]interface{}{
			"_index": "products",
			"_id":    "1",
		},
	}

	bulkReader, err := ToReader(bulkData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Bulk data as io.Reader created successfully\n")
	_ = bulkReader

	// Example 3: Using ToReaderWithValidation for safe conversions
	emptyData := map[string]interface{}{}
	_, err = ToReaderWithValidation(emptyData)
	if err != nil {
		fmt.Printf("Validation caught empty data: %v\n", err)
	}

	// Example 4: Converting JSON string to io.Reader
	jsonString := `{"query": {"match_all": {}}}`
	stringReader, err := ToReader(jsonString)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("JSON string as io.Reader created successfully\n")
	_ = stringReader

	// Example 5: Using MustToReader when you're certain conversion will succeed
	safeData := map[string]string{
		"name":  "Product Name",
		"brand": "Brand Name",
	}
	mustReader := MustToReader(safeData)
	fmt.Printf("MustToReader conversion successful\n")
	_ = mustReader

	// Example 6: Pretty-printing for debugging
	debugData := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"title": map[string]interface{}{
					"type": "text",
				},
				"price": map[string]interface{}{
					"type": "float",
				},
			},
		},
	}

	prettyReader, err := ToReaderPretty(debugData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Pretty JSON reader created for debugging\n")
	_ = prettyReader
}

// ExampleWithElasticsearchClient shows how to use ToReader with actual ES operations
func ExampleWithElasticsearchClient(client *Client) {
	// Creating an index with mapping
	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"title": map[string]interface{}{
					"type": "text",
				},
				"description": map[string]interface{}{
					"type": "text",
				},
				"price": map[string]interface{}{
					"type": "float",
				},
				"created_at": map[string]interface{}{
					"type": "date",
				},
			},
		},
	}

	if err := client.CreateIndex("products", mapping); err != nil {
		fmt.Printf("Error creating index: %v\n", err)
		return
	}

	// Indexing a document
	product := map[string]interface{}{
		"title":       "Golang Programming Book",
		"description": "Learn Go programming language",
		"price":       29.99,
		"created_at":  "2023-12-18T00:00:00Z",
	}

	if err := client.IndexDocument("products", "1", product); err != nil {
		fmt.Printf("Error indexing document: %v\n", err)
		return
	}

	// Searching with a complex query
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"title": "golang",
						},
					},
				},
				"filter": []map[string]interface{}{
					{
						"range": map[string]interface{}{
							"price": map[string]interface{}{
								"gte": 20.0,
								"lte": 50.0,
							},
						},
					},
				},
			},
		},
		"size": 10,
		"sort": []map[string]interface{}{
			{
				"created_at": map[string]interface{}{
					"order": "desc",
				},
			},
		},
	}

	results, err := client.QueryDocuments("products", searchQuery)
	if err != nil {
		fmt.Printf("Error querying documents: %v\n", err)
		return
	}

	fmt.Printf("Found %d products\n", len(results))

	// Updating a document
	updateDoc := map[string]interface{}{
		"doc": map[string]interface{}{
			"price": 24.99, // Update price
		},
	}

	if err := client.UpdateDocument("products", "1", updateDoc); err != nil {
		fmt.Printf("Error updating document: %v\n", err)
		return
	}

	fmt.Println("All operations completed successfully!")
}

// AdvancedExamples shows more complex usage patterns
func AdvancedExamples() {
	// Example 1: Bulk operations
	bulkOperations := []interface{}{
		map[string]interface{}{
			"index": map[string]interface{}{
				"_index": "products",
				"_id":    "1",
			},
		},
		map[string]interface{}{
			"title": "Product 1",
			"price": 10.99,
		},
		map[string]interface{}{
			"index": map[string]interface{}{
				"_index": "products",
				"_id":    "2",
			},
		},
		map[string]interface{}{
			"title": "Product 2",
			"price": 15.99,
		},
	}

	// Convert each operation to string and combine
	var bulkBody []byte
	for _, op := range bulkOperations {
		reader, err := ToReader(op)
		if err != nil {
			fmt.Printf("Error converting bulk operation: %v\n", err)
			continue
		}

		// In real implementation, you'd collect these and join with newlines
		_ = reader
	}

	// Example 2: Aggregation queries
	aggregationQuery := map[string]interface{}{
		"size": 0, // Don't return documents, just aggregations
		"aggs": map[string]interface{}{
			"price_ranges": map[string]interface{}{
				"range": map[string]interface{}{
					"field": "price",
					"ranges": []map[string]interface{}{
						{"to": 10.0},
						{"from": 10.0, "to": 50.0},
						{"from": 50.0},
					},
				},
			},
			"avg_price": map[string]interface{}{
				"avg": map[string]interface{}{
					"field": "price",
				},
			},
		},
	}

	aggReader, err := ToReader(aggregationQuery)
	if err != nil {
		fmt.Printf("Error creating aggregation reader: %v\n", err)
		return
	}
	_ = aggReader

	fmt.Println("Advanced examples completed!")
}

// Helper function to read from io.Reader for demonstration
func readFromReader(r io.Reader) (string, error) {
	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}
	return string(buf[:n]), nil
}
