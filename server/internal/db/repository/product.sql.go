// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: product.sql

package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AddBulkProductsParams struct {
	ProductID    uuid.UUID   `json:"product_id"`
	CategoryID   pgtype.Int4 `json:"category_id"`
	CollectionID pgtype.Int4 `json:"collection_id"`
	BrandID      pgtype.Int4 `json:"brand_id"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
}

const archiveProduct = `-- name: ArchiveProduct :exec
UPDATE
    products
SET
    archived = true,
    updated_at = NOW()
WHERE
    product_id = $1
`

func (q *Queries) ArchiveProduct(ctx context.Context, productID uuid.UUID) error {
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
    name ILIKE COALESCE($2, name)
`

type CountProductsParams struct {
	Archived pgtype.Bool `json:"archived"`
	Name     pgtype.Text `json:"name"`
}

func (q *Queries) CountProducts(ctx context.Context, arg CountProductsParams) (int64, error) {
	row := q.db.QueryRow(ctx, countProducts, arg.Archived, arg.Name)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createProduct = `-- name: CreateProduct :one
INSERT INTO products (
    product_id,
    name,
    description,
    brand_id,
    collection_id,
    category_id) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6)
RETURNING product_id, category_id, collection_id, brand_id, name, description, archived, created_at, updated_at
`

type CreateProductParams struct {
	ProductID    uuid.UUID   `json:"product_id"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	BrandID      pgtype.Int4 `json:"brand_id"`
	CollectionID pgtype.Int4 `json:"collection_id"`
	CategoryID   pgtype.Int4 `json:"category_id"`
}

func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error) {
	row := q.db.QueryRow(ctx, createProduct,
		arg.ProductID,
		arg.Name,
		arg.Description,
		arg.BrandID,
		arg.CollectionID,
		arg.CategoryID,
	)
	var i Product
	err := row.Scan(
		&i.ProductID,
		&i.CategoryID,
		&i.CollectionID,
		&i.BrandID,
		&i.Name,
		&i.Description,
		&i.Archived,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteProduct = `-- name: DeleteProduct :exec
DELETE FROM
    products
WHERE
    product_id = $1
`

func (q *Queries) DeleteProduct(ctx context.Context, productID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteProduct, productID)
	return err
}

const getProductByID = `-- name: GetProductByID :one
SELECT
    products.product_id, products.category_id, products.collection_id, products.brand_id, products.name, products.description, products.archived, products.created_at, products.updated_at,
    COUNT(pv.variant_id) AS variant_count
FROM
    products
JOIN product_variants AS pv ON products.product_id = pv.product_id
WHERE
    products.product_id = $1 AND
    archived = COALESCE($2, false)
GROUP BY
    products.product_id
`

type GetProductByIDParams struct {
	ProductID uuid.UUID   `json:"product_id"`
	Archived  pgtype.Bool `json:"archived"`
}

type GetProductByIDRow struct {
	ProductID    uuid.UUID   `json:"product_id"`
	CategoryID   pgtype.Int4 `json:"category_id"`
	CollectionID pgtype.Int4 `json:"collection_id"`
	BrandID      pgtype.Int4 `json:"brand_id"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Archived     bool        `json:"archived"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	VariantCount int64       `json:"variant_count"`
}

func (q *Queries) GetProductByID(ctx context.Context, arg GetProductByIDParams) (GetProductByIDRow, error) {
	row := q.db.QueryRow(ctx, getProductByID, arg.ProductID, arg.Archived)
	var i GetProductByIDRow
	err := row.Scan(
		&i.ProductID,
		&i.CategoryID,
		&i.CollectionID,
		&i.BrandID,
		&i.Name,
		&i.Description,
		&i.Archived,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.VariantCount,
	)
	return i, err
}

const getProductDetail = `-- name: GetProductDetail :many
SELECT
    p.product_id, p.category_id, p.collection_id, p.brand_id, p.name, p.description, p.archived, p.created_at, p.updated_at,
    pv.variant_id, pv.sku, pv.price, pv.stock_quantity, pv.discount,
    a.attribute_id AS attribute_id, a.name as attribute_name,
    va.variant_attribute_id AS variant_attribute_id, va.value as variant_attribute_value,
    img.image_id, img.image_url, img.product_id as img_product_id, img.variant_id as img_variant_id
FROM
    products p
LEFT JOIN product_variants AS pv ON p.product_id = pv.product_id
LEFT JOIN variant_attributes AS va ON pv.variant_id = va.variant_id
LEFT JOIN attributes AS a ON va.attribute_id = a.attribute_id
LEFT JOIN images AS img ON p.product_id = img.product_id
WHERE
    p.product_id = $1 AND
    archived = COALESCE($2, false)
ORDER BY
    pv.variant_id, a.attribute_id, va.variant_attribute_id DESC
`

type GetProductDetailParams struct {
	ProductID uuid.UUID   `json:"product_id"`
	Archived  pgtype.Bool `json:"archived"`
}

type GetProductDetailRow struct {
	Product               Product        `json:"product"`
	VariantID             pgtype.UUID    `json:"variant_id"`
	Sku                   pgtype.Text    `json:"sku"`
	Price                 pgtype.Numeric `json:"price"`
	StockQuantity         pgtype.Int4    `json:"stock_quantity"`
	Discount              pgtype.Int2    `json:"discount"`
	AttributeID           pgtype.Int4    `json:"attribute_id"`
	AttributeName         pgtype.Text    `json:"attribute_name"`
	VariantAttributeID    pgtype.Int4    `json:"variant_attribute_id"`
	VariantAttributeValue pgtype.Text    `json:"variant_attribute_value"`
	ImageID               pgtype.Int4    `json:"image_id"`
	ImageUrl              pgtype.Text    `json:"image_url"`
	ImgProductID          pgtype.UUID    `json:"img_product_id"`
	ImgVariantID          pgtype.UUID    `json:"img_variant_id"`
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
			&i.Product.CategoryID,
			&i.Product.CollectionID,
			&i.Product.BrandID,
			&i.Product.Name,
			&i.Product.Description,
			&i.Product.Archived,
			&i.Product.CreatedAt,
			&i.Product.UpdatedAt,
			&i.VariantID,
			&i.Sku,
			&i.Price,
			&i.StockQuantity,
			&i.Discount,
			&i.AttributeID,
			&i.AttributeName,
			&i.VariantAttributeID,
			&i.VariantAttributeValue,
			&i.ImageID,
			&i.ImageUrl,
			&i.ImgProductID,
			&i.ImgVariantID,
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

const getProductVariantsByID = `-- name: GetProductVariantsByID :one
SELECT
    products.product_id, category_id, collection_id, brand_id, name, description, archived, products.created_at, products.updated_at, variant_id, pv.product_id, price, discount, stock_quantity, sku, pv.created_at, pv.updated_at
FROM
    products
JOIN
    product_variants pv ON products.product_id = pv.product_id
WHERE
    products.product_id = $1 AND
    archived = COALESCE($2, false) AND
    pv.variant_id = COALESCE($3, pv.variant_id)
`

type GetProductVariantsByIDParams struct {
	ProductID uuid.UUID   `json:"product_id"`
	Archived  pgtype.Bool `json:"archived"`
	VariantID pgtype.UUID `json:"variant_id"`
}

type GetProductVariantsByIDRow struct {
	ProductID     uuid.UUID      `json:"product_id"`
	CategoryID    pgtype.Int4    `json:"category_id"`
	CollectionID  pgtype.Int4    `json:"collection_id"`
	BrandID       pgtype.Int4    `json:"brand_id"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	Archived      bool           `json:"archived"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	VariantID     uuid.UUID      `json:"variant_id"`
	ProductID_2   uuid.UUID      `json:"product_id_2"`
	Price         pgtype.Numeric `json:"price"`
	Discount      int16          `json:"discount"`
	StockQuantity int32          `json:"stock_quantity"`
	Sku           pgtype.Text    `json:"sku"`
	CreatedAt_2   time.Time      `json:"created_at_2"`
	UpdatedAt_2   time.Time      `json:"updated_at_2"`
}

func (q *Queries) GetProductVariantsByID(ctx context.Context, arg GetProductVariantsByIDParams) (GetProductVariantsByIDRow, error) {
	row := q.db.QueryRow(ctx, getProductVariantsByID, arg.ProductID, arg.Archived, arg.VariantID)
	var i GetProductVariantsByIDRow
	err := row.Scan(
		&i.ProductID,
		&i.CategoryID,
		&i.CollectionID,
		&i.BrandID,
		&i.Name,
		&i.Description,
		&i.Archived,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.VariantID,
		&i.ProductID_2,
		&i.Price,
		&i.Discount,
		&i.StockQuantity,
		&i.Sku,
		&i.CreatedAt_2,
		&i.UpdatedAt_2,
	)
	return i, err
}

const getProductWithImage = `-- name: GetProductWithImage :one
SELECT
    products.product_id, products.category_id, products.collection_id, products.brand_id, products.name, products.description, products.archived, products.created_at, products.updated_at,
    img.image_id AS image_id,
    img.image_url AS image_url
FROM
    products
LEFT JOIN images AS img ON products.product_id = img.product_id
LEFT JOIN product_variants AS pv ON products.product_id = pv.product_id
WHERE
    products.product_id = $1 AND
    archived = COALESCE($2, false) AND
    pv.variant_id = COALESCE($3, pv.variant_id)
`

type GetProductWithImageParams struct {
	ProductID uuid.UUID   `json:"product_id"`
	Archived  pgtype.Bool `json:"archived"`
	VariantID pgtype.UUID `json:"variant_id"`
}

type GetProductWithImageRow struct {
	ProductID    uuid.UUID   `json:"product_id"`
	CategoryID   pgtype.Int4 `json:"category_id"`
	CollectionID pgtype.Int4 `json:"collection_id"`
	BrandID      pgtype.Int4 `json:"brand_id"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Archived     bool        `json:"archived"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	ImageID      pgtype.Int4 `json:"image_id"`
	ImageUrl     pgtype.Text `json:"image_url"`
}

func (q *Queries) GetProductWithImage(ctx context.Context, arg GetProductWithImageParams) (GetProductWithImageRow, error) {
	row := q.db.QueryRow(ctx, getProductWithImage, arg.ProductID, arg.Archived, arg.VariantID)
	var i GetProductWithImageRow
	err := row.Scan(
		&i.ProductID,
		&i.CategoryID,
		&i.CollectionID,
		&i.BrandID,
		&i.Name,
		&i.Description,
		&i.Archived,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ImageID,
		&i.ImageUrl,
	)
	return i, err
}

const getProducts = `-- name: GetProducts :many
SELECT
    p.product_id, p.category_id, p.collection_id, p.brand_id, p.name, p.description, p.archived, p.created_at, p.updated_at,
    img.image_id AS image_id, img.image_url AS image_url,
    MIN(pv.price)::DECIMAL AS min_price,
    MAX(pv.price)::DECIMAL AS max_price,
    COUNT(pv.variant_id) AS variant_count
FROM
    products as p
JOIN product_variants AS pv ON p.product_id = pv.product_id
LEFT JOIN images AS img ON p.product_id = img.product_id
WHERE
    archived = COALESCE($3, archived) AND
    name ILIKE COALESCE($4, name) AND
    sku ILIKE COALESCE($5, sku)
GROUP BY
    p.product_id, img.image_id
ORDER BY
    p.product_id
LIMIT 
    $1
OFFSET
    $2
`

type GetProductsParams struct {
	Limit    int32       `json:"limit"`
	Offset   int32       `json:"offset"`
	Archived pgtype.Bool `json:"archived"`
	Name     pgtype.Text `json:"name"`
	Sku      pgtype.Text `json:"sku"`
}

type GetProductsRow struct {
	ProductID    uuid.UUID      `json:"product_id"`
	CategoryID   pgtype.Int4    `json:"category_id"`
	CollectionID pgtype.Int4    `json:"collection_id"`
	BrandID      pgtype.Int4    `json:"brand_id"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Archived     bool           `json:"archived"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	ImageID      pgtype.Int4    `json:"image_id"`
	ImageUrl     pgtype.Text    `json:"image_url"`
	MinPrice     pgtype.Numeric `json:"min_price"`
	MaxPrice     pgtype.Numeric `json:"max_price"`
	VariantCount int64          `json:"variant_count"`
}

func (q *Queries) GetProducts(ctx context.Context, arg GetProductsParams) ([]GetProductsRow, error) {
	rows, err := q.db.Query(ctx, getProducts,
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
	var items []GetProductsRow
	for rows.Next() {
		var i GetProductsRow
		if err := rows.Scan(
			&i.ProductID,
			&i.CategoryID,
			&i.CollectionID,
			&i.BrandID,
			&i.Name,
			&i.Description,
			&i.Archived,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ImageID,
			&i.ImageUrl,
			&i.MinPrice,
			&i.MaxPrice,
			&i.VariantCount,
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

const getProductsByBrand = `-- name: GetProductsByBrand :many
SELECT
    p.product_id, p.category_id, p.collection_id, p.brand_id, p.name, p.description, p.archived, p.created_at, p.updated_at,
    img.image_id AS image_id,
    img.image_url AS image_url,
    MIN(pv.price)::DECIMAL AS min_price,
    MAX(pv.price)::DECIMAL AS max_price,
    MAX(pv.price)::SMALLINT AS discount,
    COUNT(pv.variant_id) AS variant_count
FROM
    products AS p
JOIN product_variants AS pv ON p.product_id = pv.product_id
LEFT JOIN images AS img ON p.product_id = img.product_id
WHERE
    p.archived = COALESCE($4, archived) AND
    p.name ILIKE COALESCE($5, name) AND
    pv.sku ILIKE COALESCE($6, sku) AND
    p.brand_id = $1
GROUP BY
    p.product_id, img.image_id
ORDER BY
    p.product_id
LIMIT
    $2
OFFSET
    $3
`

type GetProductsByBrandParams struct {
	BrandID  pgtype.Int4 `json:"brand_id"`
	Limit    int32       `json:"limit"`
	Offset   int32       `json:"offset"`
	Archived pgtype.Bool `json:"archived"`
	Name     pgtype.Text `json:"name"`
	Sku      pgtype.Text `json:"sku"`
}

type GetProductsByBrandRow struct {
	ProductID    uuid.UUID      `json:"product_id"`
	CategoryID   pgtype.Int4    `json:"category_id"`
	CollectionID pgtype.Int4    `json:"collection_id"`
	BrandID      pgtype.Int4    `json:"brand_id"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Archived     bool           `json:"archived"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	ImageID      pgtype.Int4    `json:"image_id"`
	ImageUrl     pgtype.Text    `json:"image_url"`
	MinPrice     pgtype.Numeric `json:"min_price"`
	MaxPrice     pgtype.Numeric `json:"max_price"`
	Discount     int16          `json:"discount"`
	VariantCount int64          `json:"variant_count"`
}

func (q *Queries) GetProductsByBrand(ctx context.Context, arg GetProductsByBrandParams) ([]GetProductsByBrandRow, error) {
	rows, err := q.db.Query(ctx, getProductsByBrand,
		arg.BrandID,
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
	var items []GetProductsByBrandRow
	for rows.Next() {
		var i GetProductsByBrandRow
		if err := rows.Scan(
			&i.ProductID,
			&i.CategoryID,
			&i.CollectionID,
			&i.BrandID,
			&i.Name,
			&i.Description,
			&i.Archived,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ImageID,
			&i.ImageUrl,
			&i.MinPrice,
			&i.MaxPrice,
			&i.Discount,
			&i.VariantCount,
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

const getProductsByCategory = `-- name: GetProductsByCategory :many
SELECT
    p.product_id, p.category_id, p.collection_id, p.brand_id, p.name, p.description, p.archived, p.created_at, p.updated_at,
    img.image_id AS image_id,
    img.image_url AS image_url,
    MIN(pv.price)::DECIMAL AS min_price,
    MAX(pv.price)::DECIMAL AS max_price,
    MAX(pv.price)::SMALLINT AS discount,
    COUNT(pv.variant_id) AS variant_count
FROM
    products AS p
JOIN product_variants AS pv ON p.product_id = pv.product_id
LEFT JOIN images AS img ON p.product_id = img.product_id
WHERE
    p.archived = COALESCE($4, archived) AND
    p.name ILIKE COALESCE($5, name) AND
    pv.sku ILIKE COALESCE($6, sku) AND
    p.category_id = $1
GROUP BY
    p.product_id, img.image_id
ORDER BY
    p.product_id
LIMIT
    $2
OFFSET
    $3
`

type GetProductsByCategoryParams struct {
	CategoryID pgtype.Int4 `json:"category_id"`
	Limit      int32       `json:"limit"`
	Offset     int32       `json:"offset"`
	Archived   pgtype.Bool `json:"archived"`
	Name       pgtype.Text `json:"name"`
	Sku        pgtype.Text `json:"sku"`
}

type GetProductsByCategoryRow struct {
	ProductID    uuid.UUID      `json:"product_id"`
	CategoryID   pgtype.Int4    `json:"category_id"`
	CollectionID pgtype.Int4    `json:"collection_id"`
	BrandID      pgtype.Int4    `json:"brand_id"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Archived     bool           `json:"archived"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	ImageID      pgtype.Int4    `json:"image_id"`
	ImageUrl     pgtype.Text    `json:"image_url"`
	MinPrice     pgtype.Numeric `json:"min_price"`
	MaxPrice     pgtype.Numeric `json:"max_price"`
	Discount     int16          `json:"discount"`
	VariantCount int64          `json:"variant_count"`
}

func (q *Queries) GetProductsByCategory(ctx context.Context, arg GetProductsByCategoryParams) ([]GetProductsByCategoryRow, error) {
	rows, err := q.db.Query(ctx, getProductsByCategory,
		arg.CategoryID,
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
	var items []GetProductsByCategoryRow
	for rows.Next() {
		var i GetProductsByCategoryRow
		if err := rows.Scan(
			&i.ProductID,
			&i.CategoryID,
			&i.CollectionID,
			&i.BrandID,
			&i.Name,
			&i.Description,
			&i.Archived,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ImageID,
			&i.ImageUrl,
			&i.MinPrice,
			&i.MaxPrice,
			&i.Discount,
			&i.VariantCount,
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

const getProductsByCollection = `-- name: GetProductsByCollection :many
SELECT
    p.product_id, p.category_id, p.collection_id, p.brand_id, p.name, p.description, p.archived, p.created_at, p.updated_at,
    img.image_id AS image_id,
    img.image_url AS image_url,
    MIN(pv.price)::DECIMAL AS min_price,
    MAX(pv.price)::DECIMAL AS max_price,
    MAX(pv.price)::SMALLINT AS discount,
    COUNT(pv.variant_id) AS variant_count
FROM
    products AS p
JOIN product_variants AS pv ON p.product_id = pv.product_id
LEFT JOIN images AS img ON p.product_id = img.product_id
WHERE
    p.archived = COALESCE($4, archived) AND
    p.name ILIKE COALESCE($5, name) AND
    pv.sku ILIKE COALESCE($6, sku) AND
    p.collection_id = $1
GROUP BY
    p.product_id, img.image_id
ORDER BY
    p.product_id
LIMIT
    $2
OFFSET
    $3
`

type GetProductsByCollectionParams struct {
	CollectionID pgtype.Int4 `json:"collection_id"`
	Limit        int32       `json:"limit"`
	Offset       int32       `json:"offset"`
	Archived     pgtype.Bool `json:"archived"`
	Name         pgtype.Text `json:"name"`
	Sku          pgtype.Text `json:"sku"`
}

type GetProductsByCollectionRow struct {
	ProductID    uuid.UUID      `json:"product_id"`
	CategoryID   pgtype.Int4    `json:"category_id"`
	CollectionID pgtype.Int4    `json:"collection_id"`
	BrandID      pgtype.Int4    `json:"brand_id"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Archived     bool           `json:"archived"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	ImageID      pgtype.Int4    `json:"image_id"`
	ImageUrl     pgtype.Text    `json:"image_url"`
	MinPrice     pgtype.Numeric `json:"min_price"`
	MaxPrice     pgtype.Numeric `json:"max_price"`
	Discount     int16          `json:"discount"`
	VariantCount int64          `json:"variant_count"`
}

func (q *Queries) GetProductsByCollection(ctx context.Context, arg GetProductsByCollectionParams) ([]GetProductsByCollectionRow, error) {
	rows, err := q.db.Query(ctx, getProductsByCollection,
		arg.CollectionID,
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
	var items []GetProductsByCollectionRow
	for rows.Next() {
		var i GetProductsByCollectionRow
		if err := rows.Scan(
			&i.ProductID,
			&i.CategoryID,
			&i.CollectionID,
			&i.BrandID,
			&i.Name,
			&i.Description,
			&i.Archived,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ImageID,
			&i.ImageUrl,
			&i.MinPrice,
			&i.MaxPrice,
			&i.Discount,
			&i.VariantCount,
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
    brand_id = coalesce($3, brand_id),
    collection_id = coalesce($4, collection_id),
    category_id = coalesce($5, category_id),
    updated_at = NOW()
WHERE
    product_id = $6
RETURNING product_id, category_id, collection_id, brand_id, name, description, archived, created_at, updated_at
`

type UpdateProductParams struct {
	Name         pgtype.Text `json:"name"`
	Description  pgtype.Text `json:"description"`
	BrandID      pgtype.Int4 `json:"brand_id"`
	CollectionID pgtype.Int4 `json:"collection_id"`
	CategoryID   pgtype.Int4 `json:"category_id"`
	ProductID    uuid.UUID   `json:"product_id"`
}

func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error) {
	row := q.db.QueryRow(ctx, updateProduct,
		arg.Name,
		arg.Description,
		arg.BrandID,
		arg.CollectionID,
		arg.CategoryID,
		arg.ProductID,
	)
	var i Product
	err := row.Scan(
		&i.ProductID,
		&i.CategoryID,
		&i.CollectionID,
		&i.BrandID,
		&i.Name,
		&i.Description,
		&i.Archived,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
