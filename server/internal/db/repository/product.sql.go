// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: product.sql

package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

const archiveProduct = `-- name: ArchiveProduct :exec
UPDATE
    products
SET
    archived = true,
    updated_at = NOW()
WHERE
    product_id = $1
`

func (q *Queries) ArchiveProduct(ctx context.Context, productID int64) error {
	_, err := q.db.Exec(ctx, archiveProduct, productID)
	return err
}

const countProducts = `-- name: CountProducts :one
SELECT
    COUNT(*)
FROM
    products
WHERE
    archived = COALESCE($1, archived) AND
    name ILIKE COALESCE($2, name) AND
    sku ILIKE COALESCE($3, sku)
`

type CountProductsParams struct {
	Archived pgtype.Bool `json:"archived"`
	Name     pgtype.Text `json:"name"`
	Sku      pgtype.Text `json:"sku"`
}

func (q *Queries) CountProducts(ctx context.Context, arg CountProductsParams) (int64, error) {
	row := q.db.QueryRow(ctx, countProducts, arg.Archived, arg.Name, arg.Sku)
	var count int64
	err := row.Scan(&count)
	return count, err
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
RETURNING product_id, name, description, sku, stock, archived, price, updated_at, created_at
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
		&i.ProductID,
		&i.Name,
		&i.Description,
		&i.Sku,
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
    product_id = $1
`

func (q *Queries) DeleteProduct(ctx context.Context, productID int64) error {
	_, err := q.db.Exec(ctx, deleteProduct, productID)
	return err
}

const getProduct = `-- name: GetProduct :one
SELECT
    product_id, name, description, sku, stock, archived, price, updated_at, created_at
FROM
    products
WHERE
    product_id = $1 AND
    archived = COALESCE($2, FALSE)
`

type GetProductParams struct {
	ProductID int64       `json:"product_id"`
	Archived  pgtype.Bool `json:"archived"`
}

func (q *Queries) GetProduct(ctx context.Context, arg GetProductParams) (Product, error) {
	row := q.db.QueryRow(ctx, getProduct, arg.ProductID, arg.Archived)
	var i Product
	err := row.Scan(
		&i.ProductID,
		&i.Name,
		&i.Description,
		&i.Sku,
		&i.Stock,
		&i.Archived,
		&i.Price,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getProductDetail = `-- name: GetProductDetail :many
SELECT
    products.product_id, products.name, products.description, products.sku, products.stock, products.archived, products.price, products.updated_at, products.created_at,
    img.image_id AS image_id,
    img.image_url AS image_url,
    img.primary AS image_primary
FROM
    products
LEFT JOIN images AS img ON products.product_id = img.product_id
WHERE
    products.product_id = $1 AND
    archived = COALESCE($2, false)
ORDER BY
    img.primary DESC
`

type GetProductDetailParams struct {
	ProductID int64       `json:"product_id"`
	Archived  pgtype.Bool `json:"archived"`
}

type GetProductDetailRow struct {
	Product      Product     `json:"product"`
	ImageID      pgtype.Int4 `json:"image_id"`
	ImageUrl     pgtype.Text `json:"image_url"`
	ImagePrimary pgtype.Bool `json:"image_primary"`
}

func (q *Queries) GetProductDetail(ctx context.Context, arg GetProductDetailParams) ([]GetProductDetailRow, error) {
	rows, err := q.db.Query(ctx, getProductDetail, arg.ProductID, arg.Archived)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProductDetailRow
	for rows.Next() {
		var i GetProductDetailRow
		if err := rows.Scan(
			&i.Product.ProductID,
			&i.Product.Name,
			&i.Product.Description,
			&i.Product.Sku,
			&i.Product.Stock,
			&i.Product.Archived,
			&i.Product.Price,
			&i.Product.UpdatedAt,
			&i.Product.CreatedAt,
			&i.ImageID,
			&i.ImageUrl,
			&i.ImagePrimary,
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

const getProductWithImage = `-- name: GetProductWithImage :one
SELECT
    products.product_id, products.name, products.description, products.sku, products.stock, products.archived, products.price, products.updated_at, products.created_at,
    img.image_id AS image_id,
    img.image_url AS image_url
FROM
    products
LEFT JOIN images AS img ON products.product_id = img.product_id AND img.primary = TRUE
WHERE
    products.product_id = $1 AND
    archived = COALESCE($2, false)
`

type GetProductWithImageParams struct {
	ProductID int64       `json:"product_id"`
	Archived  pgtype.Bool `json:"archived"`
}

type GetProductWithImageRow struct {
	ProductID   int64          `json:"product_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Sku         string         `json:"sku"`
	Stock       int32          `json:"stock"`
	Archived    bool           `json:"archived"`
	Price       pgtype.Numeric `json:"price"`
	UpdatedAt   time.Time      `json:"updated_at"`
	CreatedAt   time.Time      `json:"created_at"`
	ImageID     pgtype.Int4    `json:"image_id"`
	ImageUrl    pgtype.Text    `json:"image_url"`
}

func (q *Queries) GetProductWithImage(ctx context.Context, arg GetProductWithImageParams) (GetProductWithImageRow, error) {
	row := q.db.QueryRow(ctx, getProductWithImage, arg.ProductID, arg.Archived)
	var i GetProductWithImageRow
	err := row.Scan(
		&i.ProductID,
		&i.Name,
		&i.Description,
		&i.Sku,
		&i.Stock,
		&i.Archived,
		&i.Price,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.ImageID,
		&i.ImageUrl,
	)
	return i, err
}

const listProducts = `-- name: ListProducts :many
SELECT
    products.product_id, products.name, products.description, products.sku, products.stock, products.archived, products.price, products.updated_at, products.created_at,
    img.image_id AS image_id,
    img.image_url AS image_url
FROM
    products
LEFT JOIN images AS img ON products.product_id = img.product_id AND img.primary = TRUE
WHERE
    archived = COALESCE($3, archived) AND
    name ILIKE COALESCE($4, name) AND
    sku ILIKE COALESCE($5, sku)
ORDER BY
    products.product_id
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

type ListProductsRow struct {
	ProductID   int64          `json:"product_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Sku         string         `json:"sku"`
	Stock       int32          `json:"stock"`
	Archived    bool           `json:"archived"`
	Price       pgtype.Numeric `json:"price"`
	UpdatedAt   time.Time      `json:"updated_at"`
	CreatedAt   time.Time      `json:"created_at"`
	ImageID     pgtype.Int4    `json:"image_id"`
	ImageUrl    pgtype.Text    `json:"image_url"`
}

func (q *Queries) ListProducts(ctx context.Context, arg ListProductsParams) ([]ListProductsRow, error) {
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
	var items []ListProductsRow
	for rows.Next() {
		var i ListProductsRow
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
			&i.ImageID,
			&i.ImageUrl,
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
    stock = coalesce($4, stock),
    price = coalesce($5, price),
    updated_at = NOW()
WHERE
    product_id = $6
RETURNING product_id, name, description, sku, stock, archived, price, updated_at, created_at
`

type UpdateProductParams struct {
	Name        pgtype.Text    `json:"name"`
	Description pgtype.Text    `json:"description"`
	Sku         pgtype.Text    `json:"sku"`
	Stock       pgtype.Int4    `json:"stock"`
	Price       pgtype.Numeric `json:"price"`
	ProductID   int64          `json:"product_id"`
}

func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error) {
	row := q.db.QueryRow(ctx, updateProduct,
		arg.Name,
		arg.Description,
		arg.Sku,
		arg.Stock,
		arg.Price,
		arg.ProductID,
	)
	var i Product
	err := row.Scan(
		&i.ProductID,
		&i.Name,
		&i.Description,
		&i.Sku,
		&i.Stock,
		&i.Archived,
		&i.Price,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const updateProductStock = `-- name: UpdateProductStock :exec
UPDATE
    products
SET
    stock = stock + $2
WHERE
    product_id = $1
RETURNING product_id, name, description, sku, stock, archived, price, updated_at, created_at
`

type UpdateProductStockParams struct {
	ProductID int64 `json:"product_id"`
	Stock     int32 `json:"stock"`
}

func (q *Queries) UpdateProductStock(ctx context.Context, arg UpdateProductStockParams) error {
	_, err := q.db.Exec(ctx, updateProductStock, arg.ProductID, arg.Stock)
	return err
}