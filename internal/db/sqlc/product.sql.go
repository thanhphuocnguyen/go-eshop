// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: product.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const archiveProduct = `-- name: ArchiveProduct :exec
UPDATE
    products
SET
    archived = true
WHERE
    id = $1
RETURNING id, name, description, sku, image_url, stock, archived, price, updated_at, created_at
`

func (q *Queries) ArchiveProduct(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, archiveProduct, id)
	return err
}

const createProduct = `-- name: CreateProduct :one
INSERT INTO
    products (
        name,
        description,
        sku,
        stock,
        price
    )
VALUES
    (
        $1,
        $2,
        $3,
        $4,
        $5
    )
RETURNING id, name, description, sku, image_url, stock, archived, price, updated_at, created_at
`

type CreateProductParams struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Sku         string         `json:"sku"`
	Stock       int32          `json:"stock"`
	Price       pgtype.Numeric `json:"price"`
}

func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error) {
	row := q.db.QueryRow(ctx, createProduct,
		arg.Name,
		arg.Description,
		arg.Sku,
		arg.Stock,
		arg.Price,
	)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Sku,
		&i.ImageUrl,
		&i.Stock,
		&i.Archived,
		&i.Price,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const deleteProduct = `-- name: DeleteProduct :exec
DELETE FROM
    products
WHERE
    id = $1
`

func (q *Queries) DeleteProduct(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteProduct, id)
	return err
}

const getProduct = `-- name: GetProduct :one
SELECT
    id, name, description, sku, image_url, stock, archived, price, updated_at, created_at
FROM
    products
WHERE
    id = $1
LIMIT 1
`

func (q *Queries) GetProduct(ctx context.Context, id int64) (Product, error) {
	row := q.db.QueryRow(ctx, getProduct, id)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Sku,
		&i.ImageUrl,
		&i.Stock,
		&i.Archived,
		&i.Price,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const listProducts = `-- name: ListProducts :many
SELECT
    id, name, description, sku, image_url, stock, archived, price, updated_at, created_at
FROM
    products
WHERE
    archived = COALESCE($3, archived) AND
    name ILIKE COALESCE($4, name) AND
    sku ILIKE COALESCE($5, sku)
ORDER BY
    id
LIMIT $1
OFFSET $2
`

type ListProductsParams struct {
	Limit    int32       `json:"limit"`
	Offset   int32       `json:"offset"`
	Archived pgtype.Bool `json:"archived"`
	Name     pgtype.Text `json:"name"`
	Sku      pgtype.Text `json:"sku"`
}

func (q *Queries) ListProducts(ctx context.Context, arg ListProductsParams) ([]Product, error) {
	rows, err := q.db.Query(ctx, listProducts,
		arg.Limit,
		arg.Offset,
		arg.Archived,
		arg.Name,
		arg.Sku,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Product{}
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Sku,
			&i.ImageUrl,
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

const updateProduct = `-- name: UpdateProduct :one
UPDATE
    products
SET
    name = coalesce($1, name),
    description = coalesce($2, description),
    sku = coalesce($3, sku),
    image_url = coalesce($4, image_url),
    stock = coalesce($5, stock),
    price = coalesce($6, price),
    updated_at = NOW()
WHERE
    id = $7
RETURNING id, name, description, sku, image_url, stock, archived, price, updated_at, created_at
`

type UpdateProductParams struct {
	Name        pgtype.Text    `json:"name"`
	Description pgtype.Text    `json:"description"`
	Sku         pgtype.Text    `json:"sku"`
	ImageUrl    pgtype.Text    `json:"image_url"`
	Stock       pgtype.Int4    `json:"stock"`
	Price       pgtype.Numeric `json:"price"`
	ID          int64          `json:"id"`
}

func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error) {
	row := q.db.QueryRow(ctx, updateProduct,
		arg.Name,
		arg.Description,
		arg.Sku,
		arg.ImageUrl,
		arg.Stock,
		arg.Price,
		arg.ID,
	)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Sku,
		&i.ImageUrl,
		&i.Stock,
		&i.Archived,
		&i.Price,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const updateProductImage = `-- name: UpdateProductImage :exec
UPDATE
    products
SET
    image_url = $2
WHERE
    id = $1
RETURNING id, name, description, sku, image_url, stock, archived, price, updated_at, created_at
`

type UpdateProductImageParams struct {
	ID       int64       `json:"id"`
	ImageUrl pgtype.Text `json:"image_url"`
}

func (q *Queries) UpdateProductImage(ctx context.Context, arg UpdateProductImageParams) error {
	_, err := q.db.Exec(ctx, updateProductImage, arg.ID, arg.ImageUrl)
	return err
}

const updateProductStock = `-- name: UpdateProductStock :exec
UPDATE
    products
SET
    stock = stock + $2
WHERE
    id = $1
RETURNING id, name, description, sku, image_url, stock, archived, price, updated_at, created_at
`

type UpdateProductStockParams struct {
	ID    int64 `json:"id"`
	Stock int32 `json:"stock"`
}

func (q *Queries) UpdateProductStock(ctx context.Context, arg UpdateProductStockParams) error {
	_, err := q.db.Exec(ctx, updateProductStock, arg.ID, arg.Stock)
	return err
}
