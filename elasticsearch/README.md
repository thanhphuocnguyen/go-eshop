# Elasticsearch Package for Product Indexing

This package provides functionality for indexing product objects in Elasticsearch. It includes methods for creating, updating, deleting, and querying product documents.

## Package Structure

- **client.go**: Implements the Elasticsearch client, handling connections and requests.
- **product_indexer.go**: Exports the `ProductIndexer` struct with methods for indexing product objects.
- **types.go**: Defines data structures for product objects and related interfaces.
- **queries.go**: Contains functions for constructing and executing queries against the product index.
- **mappings/product_mapping.json**: Defines the mapping for the product index, specifying data types and properties.
- **config.go**: Manages configuration settings for connecting to Elasticsearch.

## Usage Instructions

1. **Setup Elasticsearch**: Ensure you have an Elasticsearch server running and accessible.
2. **Configure Connection**: Update the configuration settings in `config.go` with your server details.
3. **Indexing Products**:
   - Create an instance of `ProductIndexer`.
   - Use `IndexProduct`, `UpdateProduct`, or `DeleteProduct` methods to manage product documents.
4. **Querying Products**: Utilize the functions in `queries.go` to search for products based on your criteria.

## Examples

```go
// Example of indexing a product
productIndexer := NewProductIndexer()
err := productIndexer.IndexProduct(product)
if err != nil {
    log.Fatalf("Error indexing product: %v", err)
}

// Example of querying products
results, err := productIndexer.QueryProducts(query)
if err != nil {
    log.Fatalf("Error querying products: %v", err)
}
```

## Additional Information

Refer to the official Elasticsearch documentation for more details on index management and query syntax.