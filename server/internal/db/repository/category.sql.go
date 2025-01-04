// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: category.sql

package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createCollection = `-- name: CreateCollection :one
INSERT INTO categories (name, sort_order, image_url, published)
VALUES ($1, $2, $3, $4)
RETURNING category_id, name, sort_order, image_url, published, created_at, updated_at
`

type CreateCollectionParams struct {
	Name      string      `json:"name"`
	SortOrder int16       `json:"sort_order"`
	ImageUrl  pgtype.Text `json:"image_url"`
	Published bool        `json:"published"`
}

func (q *Queries) CreateCollection(ctx context.Context, arg CreateCollectionParams) (Category, error) {
	row := q.db.QueryRow(ctx, createCollection,
		arg.Name,
		arg.SortOrder,
		arg.ImageUrl,
		arg.Published,
	)
	var i Category
	err := row.Scan(
		&i.CategoryID,
		&i.Name,
		&i.SortOrder,
		&i.ImageUrl,
		&i.Published,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getCollection = `-- name: GetCollection :one
SELECT categories.category_id, categories.name, categories.sort_order, categories.image_url, categories.published, categories.created_at, categories.updated_at, images.image_id, images.product_id, images.variant_id, images.image_url, images.external_id, images."primary", images.created_at, images.updated_at, products.product_id, products.name, products.description, products.sku, products.stock, products.archived, products.price, products.updated_at, products.created_at, category_products.category_id, category_products.product_id, category_products.sort_order
FROM categories
JOIN category_products ON categories.category_id = category
JOIN products ON category_products.product_id = products.product_id
JOIN images ON products.product_id = images.product_id AND images.primary = TRUE
WHERE categories.category_id = $1
LIMIT 1
`

type GetCollectionRow struct {
	Category        Category        `json:"category"`
	Image           Image           `json:"image"`
	Product         Product         `json:"product"`
	CategoryProduct CategoryProduct `json:"category_product"`
}

func (q *Queries) GetCollection(ctx context.Context, categoryID int32) (GetCollectionRow, error) {
	row := q.db.QueryRow(ctx, getCollection, categoryID)
	var i GetCollectionRow
	err := row.Scan(
		&i.Category.CategoryID,
		&i.Category.Name,
		&i.Category.SortOrder,
		&i.Category.ImageUrl,
		&i.Category.Published,
		&i.Category.CreatedAt,
		&i.Category.UpdatedAt,
		&i.Image.ImageID,
		&i.Image.ProductID,
		&i.Image.VariantID,
		&i.Image.ImageUrl,
		&i.Image.ExternalID,
		&i.Image.Primary,
		&i.Image.CreatedAt,
		&i.Image.UpdatedAt,
		&i.Product.ProductID,
		&i.Product.Name,
		&i.Product.Description,
		&i.Product.Sku,
		&i.Product.Stock,
		&i.Product.Archived,
		&i.Product.Price,
		&i.Product.UpdatedAt,
		&i.Product.CreatedAt,
		&i.CategoryProduct.CategoryID,
		&i.CategoryProduct.ProductID,
		&i.CategoryProduct.SortOrder,
	)
	return i, err
}

const getCollectionByName = `-- name: GetCollectionByName :one
SELECT category_id, name, sort_order, image_url, published, created_at, updated_at
FROM categories
WHERE name = $1
LIMIT 1
`

func (q *Queries) GetCollectionByName(ctx context.Context, name string) (Category, error) {
	row := q.db.QueryRow(ctx, getCollectionByName, name)
	var i Category
	err := row.Scan(
		&i.CategoryID,
		&i.Name,
		&i.SortOrder,
		&i.ImageUrl,
		&i.Published,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
