// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: category_product.sql

package repository

import (
	"context"
)

const addProductToCollection = `-- name: AddProductToCollection :one
INSERT INTO
    category_products (category_id, product_id, sort_order)
VALUES
    ($1, $2, $3)
RETURNING category_id, product_id, sort_order
`

type AddProductToCollectionParams struct {
	CategoryID int32 `json:"category_id"`
	ProductID  int64 `json:"product_id"`
	SortOrder  int16 `json:"sort_order"`
}

func (q *Queries) AddProductToCollection(ctx context.Context, arg AddProductToCollectionParams) (CategoryProduct, error) {
	row := q.db.QueryRow(ctx, addProductToCollection, arg.CategoryID, arg.ProductID, arg.SortOrder)
	var i CategoryProduct
	err := row.Scan(&i.CategoryID, &i.ProductID, &i.SortOrder)
	return i, err
}

const getCollectionProducts = `-- name: GetCollectionProducts :many
SELECT
    p.product_id, p.name, p.description, p.sku, p.stock, p.archived, p.price, p.updated_at, p.created_at
FROM
    products p
    JOIN category_products cp ON p.product_id = cp.product_id
WHERE
    cp.category_id = $1
`

func (q *Queries) GetCollectionProducts(ctx context.Context, categoryID int32) ([]Product, error) {
	rows, err := q.db.Query(ctx, getCollectionProducts, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Product
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ProductID,
			&i.Name,
			&i.Description,
			&i.Sku,
			&i.Stock,
			&i.Archived,
			&i.Price,
			&i.UpdatedAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const removeProductFromCollection = `-- name: RemoveProductFromCollection :one
DELETE FROM
    category_products
WHERE
    category_id = $1
    AND product_id = $2
RETURNING category_id, product_id, sort_order
`

type RemoveProductFromCollectionParams struct {
	CategoryID int32 `json:"category_id"`
	ProductID  int64 `json:"product_id"`
}

func (q *Queries) RemoveProductFromCollection(ctx context.Context, arg RemoveProductFromCollectionParams) (CategoryProduct, error) {
	row := q.db.QueryRow(ctx, removeProductFromCollection, arg.CategoryID, arg.ProductID)
	var i CategoryProduct
	err := row.Scan(&i.CategoryID, &i.ProductID, &i.SortOrder)
	return i, err
}
