// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: brand.sql

package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const countBrands = `-- name: CountBrands :one
SELECT count(*)
FROM brands
`

func (q *Queries) CountBrands(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countBrands)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createBrand = `-- name: CreateBrand :one
INSERT INTO brands 
    (id, name, slug, description, image_url, image_id)
VALUES 
    ($1, $2, $3, $4, $5, $6)
RETURNING id, name, image_url, image_id, description, slug, remarkable, display_order, published, created_at, updated_at
`

type CreateBrandParams struct {
	ID          uuid.UUID   `json:"id"`
	Name        string      `json:"name"`
	Slug        string      `json:"slug"`
	Description pgtype.Text `json:"description"`
	ImageUrl    pgtype.Text `json:"image_url"`
	ImageID     pgtype.Text `json:"image_id"`
}

func (q *Queries) CreateBrand(ctx context.Context, arg CreateBrandParams) (Brand, error) {
	row := q.db.QueryRow(ctx, createBrand,
		arg.ID,
		arg.Name,
		arg.Slug,
		arg.Description,
		arg.ImageUrl,
		arg.ImageID,
	)
	var i Brand
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

const deleteBrand = `-- name: DeleteBrand :exec
DELETE FROM brands WHERE id = $1
`

func (q *Queries) DeleteBrand(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteBrand, id)
	return err
}

const getBrandByID = `-- name: GetBrandByID :one
SELECT c.id, c.name, c.image_url, c.image_id, c.description, c.slug, c.remarkable, c.display_order, c.published, c.created_at, c.updated_at 
FROM brands c 
WHERE c.id = $1 LIMIT 1
`

func (q *Queries) GetBrandByID(ctx context.Context, id uuid.UUID) (Brand, error) {
	row := q.db.QueryRow(ctx, getBrandByID, id)
	var i Brand
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

const getBrands = `-- name: GetBrands :many
SELECT
    c.id, c.name, c.image_url, c.image_id, c.description, c.slug, c.remarkable, c.display_order, c.published, c.created_at, c.updated_at
FROM
    brands AS c
LIMIT $1
OFFSET $2
`

type GetBrandsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetBrands(ctx context.Context, arg GetBrandsParams) ([]Brand, error) {
	rows, err := q.db.Query(ctx, getBrands, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Brand{}
	for rows.Next() {
		var i Brand
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

const getBrandsByIDs = `-- name: GetBrandsByIDs :many
SELECT 
    c.id, c.name, c.image_url, c.image_id, c.description, c.slug, c.remarkable, c.display_order, c.published, c.created_at, c.updated_at, 
    p.name as product_name, p.id as product_id, p.description, p.base_price as product_price, p.base_sku as product_sku, p.slug as product_slug
FROM brands AS c
LEFT JOIN products AS p ON p.brand_id = c.id
WHERE c.id = ANY($3::int[])
GROUP BY c.id, p.id
LIMIT $1 OFFSET $2
`

type GetBrandsByIDsParams struct {
	Limit  int32   `json:"limit"`
	Offset int32   `json:"offset"`
	Ids    []int32 `json:"ids"`
}

type GetBrandsByIDsRow struct {
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
	ProductID     pgtype.UUID    `json:"product_id"`
	Description_2 pgtype.Text    `json:"description_2"`
	ProductPrice  pgtype.Numeric `json:"product_price"`
	ProductSku    pgtype.Text    `json:"product_sku"`
	ProductSlug   pgtype.Text    `json:"product_slug"`
}

func (q *Queries) GetBrandsByIDs(ctx context.Context, arg GetBrandsByIDsParams) ([]GetBrandsByIDsRow, error) {
	rows, err := q.db.Query(ctx, getBrandsByIDs, arg.Limit, arg.Offset, arg.Ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetBrandsByIDsRow{}
	for rows.Next() {
		var i GetBrandsByIDsRow
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
			&i.Description_2,
			&i.ProductPrice,
			&i.ProductSku,
			&i.ProductSlug,
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

const updateBrandWith = `-- name: UpdateBrandWith :one
UPDATE brands
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

type UpdateBrandWithParams struct {
	ID          uuid.UUID   `json:"id"`
	Name        pgtype.Text `json:"name"`
	ImageUrl    pgtype.Text `json:"image_url"`
	ImageID     pgtype.Text `json:"image_id"`
	Description pgtype.Text `json:"description"`
	Remarkable  pgtype.Bool `json:"remarkable"`
	Slug        pgtype.Text `json:"slug"`
	Published   pgtype.Bool `json:"published"`
}

func (q *Queries) UpdateBrandWith(ctx context.Context, arg UpdateBrandWithParams) (Brand, error) {
	row := q.db.QueryRow(ctx, updateBrandWith,
		arg.ID,
		arg.Name,
		arg.ImageUrl,
		arg.ImageID,
		arg.Description,
		arg.Remarkable,
		arg.Slug,
		arg.Published,
	)
	var i Brand
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
