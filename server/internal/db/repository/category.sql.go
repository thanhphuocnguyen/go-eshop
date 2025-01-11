// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: category.sql

package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

const countCollections = `-- name: CountCollections :one
SELECT count(*)
FROM categories
WHERE category_id = COALESCE($1, category_id)
`

func (q *Queries) CountCollections(ctx context.Context, categoryID pgtype.Int4) (int64, error) {
	row := q.db.QueryRow(ctx, countCollections, categoryID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createCollection = `-- name: CreateCollection :one
INSERT INTO categories (name, description, sort_order, published)
VALUES ($1, $2, $3, $4)
RETURNING category_id, name, description, sort_order, published, created_at, updated_at
`

type CreateCollectionParams struct {
	Name        string      `json:"name"`
	Description pgtype.Text `json:"description"`
	SortOrder   int16       `json:"sort_order"`
	Published   bool        `json:"published"`
}

func (q *Queries) CreateCollection(ctx context.Context, arg CreateCollectionParams) (Category, error) {
	row := q.db.QueryRow(ctx, createCollection,
		arg.Name,
		arg.Description,
		arg.SortOrder,
		arg.Published,
	)
	var i Category
	err := row.Scan(
		&i.CategoryID,
		&i.Name,
		&i.Description,
		&i.SortOrder,
		&i.Published,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getCollection = `-- name: GetCollection :many
SELECT 
    c.category_id, c.name, c.description, c.sort_order, c.published, c.created_at, c.updated_at, 
    p.name, p.description, p.price, p.discount, 
    cp.product_id, 
    i.image_id, i.image_url
FROM categories c
JOIN category_products cp ON cp.category_id = c.category_id
JOIN products p ON cp.product_id = p.product_id AND p.published = TRUE
LEFT JOIN images i ON p.product_id = i.product_id AND i.primary = TRUE
WHERE categories.category_id = $1
`

type GetCollectionRow struct {
	CategoryID    int32          `json:"category_id"`
	Name          string         `json:"name"`
	Description   pgtype.Text    `json:"description"`
	SortOrder     int16          `json:"sort_order"`
	Published     bool           `json:"published"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	Name_2        string         `json:"name_2"`
	Description_2 string         `json:"description_2"`
	Price         pgtype.Numeric `json:"price"`
	Discount      int32          `json:"discount"`
	ProductID     int64          `json:"product_id"`
	ImageID       pgtype.Int4    `json:"image_id"`
	ImageUrl      pgtype.Text    `json:"image_url"`
}

func (q *Queries) GetCollection(ctx context.Context, categoryID int32) ([]GetCollectionRow, error) {
	rows, err := q.db.Query(ctx, getCollection, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCollectionRow
	for rows.Next() {
		var i GetCollectionRow
		if err := rows.Scan(
			&i.CategoryID,
			&i.Name,
			&i.Description,
			&i.SortOrder,
			&i.Published,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name_2,
			&i.Description_2,
			&i.Price,
			&i.Discount,
			&i.ProductID,
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

const getCollectionByName = `-- name: GetCollectionByName :one
SELECT c.category_id, c.name, c.description, c.sort_order, published, c.created_at, c.updated_at, cp.category_id, cp.product_id, cp.sort_order, p.product_id, p.name, p.description, sku, stock, discount, archived, price, p.updated_at, p.created_at
FROM categories c
JOIN category_products cp ON c.category_id = cp.category_id
JOIN products p ON cp.product_id = p.product_id AND p.published = TRUE
WHERE c.name = $1 AND c.published = TRUE
LIMIT 1
`

type GetCollectionByNameRow struct {
	CategoryID    int32          `json:"category_id"`
	Name          string         `json:"name"`
	Description   pgtype.Text    `json:"description"`
	SortOrder     int16          `json:"sort_order"`
	Published     bool           `json:"published"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	CategoryID_2  int32          `json:"category_id_2"`
	ProductID     int64          `json:"product_id"`
	SortOrder_2   int16          `json:"sort_order_2"`
	ProductID_2   int64          `json:"product_id_2"`
	Name_2        string         `json:"name_2"`
	Description_2 string         `json:"description_2"`
	Sku           pgtype.Text    `json:"sku"`
	Stock         int32          `json:"stock"`
	Discount      int32          `json:"discount"`
	Archived      bool           `json:"archived"`
	Price         pgtype.Numeric `json:"price"`
	UpdatedAt_2   time.Time      `json:"updated_at_2"`
	CreatedAt_2   time.Time      `json:"created_at_2"`
}

func (q *Queries) GetCollectionByName(ctx context.Context, name string) (GetCollectionByNameRow, error) {
	row := q.db.QueryRow(ctx, getCollectionByName, name)
	var i GetCollectionByNameRow
	err := row.Scan(
		&i.CategoryID,
		&i.Name,
		&i.Description,
		&i.SortOrder,
		&i.Published,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.CategoryID_2,
		&i.ProductID,
		&i.SortOrder_2,
		&i.ProductID_2,
		&i.Name_2,
		&i.Description_2,
		&i.Sku,
		&i.Stock,
		&i.Discount,
		&i.Archived,
		&i.Price,
		&i.UpdatedAt_2,
		&i.CreatedAt_2,
	)
	return i, err
}

const getCollections = `-- name: GetCollections :many
SELECT 
    c.category_id, c.name, c.description, c.sort_order, c.published,
    p.name, p.description, p.price, p.discount, 
    cp.product_id, 
    images.image_id, images.image_url
FROM categories c
JOIN category_products cp ON c.category_id = cp.category_id
JOIN products p ON category_products.product_id = p.product_id AND p.published = TRUE
LEFT JOIN images ON p.product_id = images.product_id AND images.primary = TRUE
WHERE categories.category_id = ANY($1::int[]) AND published = TRUE
ORDER BY c.sort_order, cp.sort_order
`

type GetCollectionsRow struct {
	CategoryID    int32          `json:"category_id"`
	Name          string         `json:"name"`
	Description   pgtype.Text    `json:"description"`
	SortOrder     int16          `json:"sort_order"`
	Published     bool           `json:"published"`
	Name_2        string         `json:"name_2"`
	Description_2 string         `json:"description_2"`
	Price         pgtype.Numeric `json:"price"`
	Discount      int32          `json:"discount"`
	ProductID     int64          `json:"product_id"`
	ImageID       pgtype.Int4    `json:"image_id"`
	ImageUrl      pgtype.Text    `json:"image_url"`
}

func (q *Queries) GetCollections(ctx context.Context, categoryIds []int32) ([]GetCollectionsRow, error) {
	rows, err := q.db.Query(ctx, getCollections, categoryIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCollectionsRow
	for rows.Next() {
		var i GetCollectionsRow
		if err := rows.Scan(
			&i.CategoryID,
			&i.Name,
			&i.Description,
			&i.SortOrder,
			&i.Published,
			&i.Name_2,
			&i.Description_2,
			&i.Price,
			&i.Discount,
			&i.ProductID,
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

const removeCollection = `-- name: RemoveCollection :exec
DELETE FROM categories
WHERE category_id = $1
`

func (q *Queries) RemoveCollection(ctx context.Context, categoryID int32) error {
	_, err := q.db.Exec(ctx, removeCollection, categoryID)
	return err
}

type SeedCollectionsParams struct {
	Name        string      `json:"name"`
	Description pgtype.Text `json:"description"`
	SortOrder   int16       `json:"sort_order"`
	Published   bool        `json:"published"`
}

const updateCollection = `-- name: UpdateCollection :one
UPDATE categories
SET 
    name = COALESCE($2, name), 
    description = COALESCE($3, description), 
    sort_order = COALESCE($4, sort_order), 
    published = COALESCE($5, published),
    updated_at = now()
WHERE category_id = $1
RETURNING category_id, name, description, sort_order, published, created_at, updated_at
`

type UpdateCollectionParams struct {
	CategoryID  int32       `json:"category_id"`
	Name        pgtype.Text `json:"name"`
	Description pgtype.Text `json:"description"`
	SortOrder   pgtype.Int2 `json:"sort_order"`
	Published   pgtype.Bool `json:"published"`
}

func (q *Queries) UpdateCollection(ctx context.Context, arg UpdateCollectionParams) (Category, error) {
	row := q.db.QueryRow(ctx, updateCollection,
		arg.CategoryID,
		arg.Name,
		arg.Description,
		arg.SortOrder,
		arg.Published,
	)
	var i Category
	err := row.Scan(
		&i.CategoryID,
		&i.Name,
		&i.Description,
		&i.SortOrder,
		&i.Published,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
