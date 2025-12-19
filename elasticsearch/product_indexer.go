package elasticsearch

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type ProductIndexer struct {
	esStore ESStore
}

func NewProductIndexer(store ESStore) *ProductIndexer {
	return &ProductIndexer{esStore: store}
}

func (pi *ProductIndexer) IndexProduct(product Product) error {
	return pi.esStore.IndexDocument("products", product.ID, product)
}

func (pi *ProductIndexer) UpdateProduct(productID string, updatedProduct Product) error {
	return pi.esStore.UpdateDocument("products", productID, updatedProduct)
}

func (pi *ProductIndexer) DeleteProduct(productID string) error {
	return pi.esStore.DeleteDocument("products", productID)
}
