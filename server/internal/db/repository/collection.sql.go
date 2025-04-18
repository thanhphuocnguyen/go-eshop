// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: collection.sql

package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const countCollections = `-- name: CountCollections :one
SELECT count(*)
FROM collections
`

func (q *Queries) CountCollections(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countCollections)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createCollection = `-- name: CreateCollection :one
INSERT INTO collections 
    (id, name, slug, description, image_url, image_id)
VALUES 
    ($1, $2, $3, $4, $5, $6)
RETURNING id, name, image_url, image_id, description, slug, remarkable, display_order, published, created_at, updated_at
`

type CreateCollectionParams struct {
	ID          uuid.UUID   `json:"id"`
	Name        string      `json:"name"`
	Slug        string      `json:"slug"`
	Description pgtype.Text `json:"description"`
	ImageUrl    pgtype.Text `json:"image_url"`
	ImageID     pgtype.Text `json:"image_id"`
}

func (q *Queries) CreateCollection(ctx context.Context, arg CreateCollectionParams) (Collection, error) {
	row := q.db.QueryRow(ctx, createCollection,
		arg.ID,
		arg.Name,
		arg.Slug,
		arg.Description,
		arg.ImageUrl,
		arg.ImageID,
	)
	var i Collection
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ImageUrl,
		&i.ImageID,
		&i.Description,
		&i.Slug,
		&i.Remarkable,
		&i.DisplayOrder,
		&i.Published,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteCollection = `-- name: DeleteCollection :exec
DELETE FROM collections WHERE id = $1
`

func (q *Queries) DeleteCollection(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteCollection, id)
	return err
}

const getCollectionByID = `-- name: GetCollectionByID :one
SELECT c.id, c.name, c.image_url, c.image_id, c.description, c.slug, c.remarkable, c.display_order, c.published, c.created_at, c.updated_at
FROM collections c 
WHERE c.id = $1
LIMIT 1
`

func (q *Queries) GetCollectionByID(ctx context.Context, id uuid.UUID) (Collection, error) {
	row := q.db.QueryRow(ctx, getCollectionByID, id)
	var i Collection
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ImageUrl,
		&i.ImageID,
		&i.Description,
		&i.Slug,
		&i.Remarkable,
		&i.DisplayOrder,
		&i.Published,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getCollectionByIDWithProducts = `-- name: GetCollectionByIDWithProducts :many
SELECT 
    c.id, c.name, c.image_url, c.image_id, c.description, c.slug, c.remarkable, c.display_order, c.published, c.created_at, c.updated_at,
    p.name as product_name, p.id as product_id, p.description as product_description, 
    p.base_price as product_price, 
    p.base_sku as product_sku, p.slug as product_slug,
    img.id as image_id, img.url as product_image_url
FROM 
    collections c
LEFT JOIN
    products p ON c.id = p.collection_id
LEFT JOIN
    image_assignments ia ON p.id = ia.entity_id AND ia.entity_type = 'product'
LEFT JOIN
    images img ON img.id = ia.image_id
WHERE c.id = $1
`

type GetCollectionByIDWithProductsRow struct {
	ID                 uuid.UUID      `json:"id"`
	Name               string         `json:"name"`
	ImageUrl           pgtype.Text    `json:"image_url"`
	ImageID            pgtype.Text    `json:"image_id"`
	Description        pgtype.Text    `json:"description"`
	Slug               string         `json:"slug"`
	Remarkable         pgtype.Bool    `json:"remarkable"`
	DisplayOrder       pgtype.Int4    `json:"display_order"`
	Published          bool           `json:"published"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	ProductName        pgtype.Text    `json:"product_name"`
	ProductID          pgtype.UUID    `json:"product_id"`
	ProductDescription pgtype.Text    `json:"product_description"`
	ProductPrice       pgtype.Numeric `json:"product_price"`
	ProductSku         pgtype.Text    `json:"product_sku"`
	ProductSlug        pgtype.Text    `json:"product_slug"`
	ImageID_2          pgtype.Int4    `json:"image_id_2"`
	ProductImageUrl    pgtype.Text    `json:"product_image_url"`
}

func (q *Queries) GetCollectionByIDWithProducts(ctx context.Context, id uuid.UUID) ([]GetCollectionByIDWithProductsRow, error) {
	rows, err := q.db.Query(ctx, getCollectionByIDWithProducts, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetCollectionByIDWithProductsRow{}
	for rows.Next() {
		var i GetCollectionByIDWithProductsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.ImageUrl,
			&i.ImageID,
			&i.Description,
			&i.Slug,
			&i.Remarkable,
			&i.DisplayOrder,
			&i.Published,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ProductName,
			&i.ProductID,
			&i.ProductDescription,
			&i.ProductPrice,
			&i.ProductSku,
			&i.ProductSlug,
			&i.ImageID_2,
			&i.ProductImageUrl,
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

const getCollections = `-- name: GetCollections :many
SELECT 
    c.id, c.name, c.image_url, c.image_id, c.description, c.slug, c.remarkable, c.display_order, c.published, c.created_at, c.updated_at
FROM collections AS c
LIMIT $1 OFFSET $2
`

type GetCollectionsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetCollections(ctx context.Context, arg GetCollectionsParams) ([]Collection, error) {
	rows, err := q.db.Query(ctx, getCollections, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Collection{}
	for rows.Next() {
		var i Collection
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.ImageUrl,
			&i.ImageID,
			&i.Description,
			&i.Slug,
			&i.Remarkable,
			&i.DisplayOrder,
			&i.Published,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getCollectionsByIDs = `-- name: GetCollectionsByIDs :many
SELECT 
    c.id, c.name, c.image_url, c.image_id, c.description, c.slug, c.remarkable, c.display_order, c.published, c.created_at, c.updated_at, 
    p.name as product_name, p.id, p.description,
    p.base_price as product_price, 
    p.base_sku as product_sku, p.slug as product_slug,
    img.id as image_id, img.url
FROM collections AS c
LEFT JOIN products AS p ON c.id = p.collection_id
LEFT JOIN image_assignments AS ia ON p.id = ia.entity_id AND ia.entity_type = 'product'
LEFT JOIN images AS img ON img.id = ia.image_id
WHERE c.id = ANY($3::int[])
GROUP BY c.id, p.id, img.id, img.url
LIMIT $1 OFFSET $2
`

type GetCollectionsByIDsParams struct {
	Limit  int32   `json:"limit"`
	Offset int32   `json:"offset"`
	Ids    []int32 `json:"ids"`
}

type GetCollectionsByIDsRow struct {
	ID            uuid.UUID      `json:"id"`
	Name          string         `json:"name"`
	ImageUrl      pgtype.Text    `json:"image_url"`
	ImageID       pgtype.Text    `json:"image_id"`
	Description   pgtype.Text    `json:"description"`
	Slug          string         `json:"slug"`
	Remarkable    pgtype.Bool    `json:"remarkable"`
	DisplayOrder  pgtype.Int4    `json:"display_order"`
	Published     bool           `json:"published"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	ProductName   pgtype.Text    `json:"product_name"`
	ID_2          pgtype.UUID    `json:"id_2"`
	Description_2 pgtype.Text    `json:"description_2"`
	ProductPrice  pgtype.Numeric `json:"product_price"`
	ProductSku    pgtype.Text    `json:"product_sku"`
	ProductSlug   pgtype.Text    `json:"product_slug"`
	ImageID_2     pgtype.Int4    `json:"image_id_2"`
	Url           pgtype.Text    `json:"url"`
}

func (q *Queries) GetCollectionsByIDs(ctx context.Context, arg GetCollectionsByIDsParams) ([]GetCollectionsByIDsRow, error) {
	rows, err := q.db.Query(ctx, getCollectionsByIDs, arg.Limit, arg.Offset, arg.Ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetCollectionsByIDsRow{}
	for rows.Next() {
		var i GetCollectionsByIDsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.ImageUrl,
			&i.ImageID,
			&i.Description,
			&i.Slug,
			&i.Remarkable,
			&i.DisplayOrder,
			&i.Published,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ProductName,
			&i.ID_2,
			&i.Description_2,
			&i.ProductPrice,
			&i.ProductSku,
			&i.ProductSlug,
			&i.ImageID_2,
			&i.Url,
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

type SeedCollectionsParams struct {
	Name        string      `json:"name"`
	Description pgtype.Text `json:"description"`
	ImageUrl    pgtype.Text `json:"image_url"`
}

const updateCollectionWith = `-- name: UpdateCollectionWith :one
UPDATE collections
SET 
    name = COALESCE($2, name), 
    image_url = COALESCE($3, image_url), 
    image_id = COALESCE($4, image_id),
    description = COALESCE($5, description),
    remarkable = COALESCE($6, remarkable),
    slug = COALESCE($7, slug),
    published = COALESCE($8, published),
    updated_at = now()
WHERE id = $1
RETURNING id, name, image_url, image_id, description, slug, remarkable, display_order, published, created_at, updated_at
`

type UpdateCollectionWithParams struct {
	ID          uuid.UUID   `json:"id"`
	Name        pgtype.Text `json:"name"`
	ImageUrl    pgtype.Text `json:"image_url"`
	ImageID     pgtype.Text `json:"image_id"`
	Description pgtype.Text `json:"description"`
	Remarkable  pgtype.Bool `json:"remarkable"`
	Slug        pgtype.Text `json:"slug"`
	Published   pgtype.Bool `json:"published"`
}

func (q *Queries) UpdateCollectionWith(ctx context.Context, arg UpdateCollectionWithParams) (Collection, error) {
	row := q.db.QueryRow(ctx, updateCollectionWith,
		arg.ID,
		arg.Name,
		arg.ImageUrl,
		arg.ImageID,
		arg.Description,
		arg.Remarkable,
		arg.Slug,
		arg.Published,
	)
	var i Collection
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ImageUrl,
		&i.ImageID,
		&i.Description,
		&i.Slug,
		&i.Remarkable,
		&i.DisplayOrder,
		&i.Published,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
