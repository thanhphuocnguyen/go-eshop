// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: image.sql

package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createImage = `-- name: CreateImage :one
INSERT INTO images (product_id, variant_id, image_url, external_id) VALUES ($1, $2, $3, $4) RETURNING image_id, product_id, variant_id, image_url, external_id, "primary", created_at, updated_at
`

type CreateImageParams struct {
	ProductID  pgtype.Int8 `json:"product_id"`
	VariantID  pgtype.Int8 `json:"variant_id"`
	ImageUrl   string      `json:"image_url"`
	ExternalID pgtype.Text `json:"external_id"`
}

func (q *Queries) CreateImage(ctx context.Context, arg CreateImageParams) (Image, error) {
	row := q.db.QueryRow(ctx, createImage,
		arg.ProductID,
		arg.VariantID,
		arg.ImageUrl,
		arg.ExternalID,
	)
	var i Image
	err := row.Scan(
		&i.ImageID,
		&i.ProductID,
		&i.VariantID,
		&i.ImageUrl,
		&i.ExternalID,
		&i.Primary,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteImage = `-- name: DeleteImage :exec
DELETE FROM images WHERE image_id = $1
`

func (q *Queries) DeleteImage(ctx context.Context, imageID int32) error {
	_, err := q.db.Exec(ctx, deleteImage, imageID)
	return err
}

const getImageByExternalID = `-- name: GetImageByExternalID :one
SELECT image_id, product_id, variant_id, image_url, external_id, "primary", created_at, updated_at FROM images WHERE external_id = $1 LIMIT 1
`

func (q *Queries) GetImageByExternalID(ctx context.Context, externalID pgtype.Text) (Image, error) {
	row := q.db.QueryRow(ctx, getImageByExternalID, externalID)
	var i Image
	err := row.Scan(
		&i.ImageID,
		&i.ProductID,
		&i.VariantID,
		&i.ImageUrl,
		&i.ExternalID,
		&i.Primary,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getImageByID = `-- name: GetImageByID :one
SELECT image_id, product_id, variant_id, image_url, external_id, "primary", created_at, updated_at FROM images WHERE image_id = $1 LIMIT 1
`

func (q *Queries) GetImageByID(ctx context.Context, imageID int32) (Image, error) {
	row := q.db.QueryRow(ctx, getImageByID, imageID)
	var i Image
	err := row.Scan(
		&i.ImageID,
		&i.ProductID,
		&i.VariantID,
		&i.ImageUrl,
		&i.ExternalID,
		&i.Primary,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getImagesByProductID = `-- name: GetImagesByProductID :many
SELECT image_id, product_id, variant_id, image_url, external_id, "primary", created_at, updated_at FROM images WHERE product_id = $1
`

func (q *Queries) GetImagesByProductID(ctx context.Context, productID pgtype.Int8) ([]Image, error) {
	rows, err := q.db.Query(ctx, getImagesByProductID, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Image
	for rows.Next() {
		var i Image
		if err := rows.Scan(
			&i.ImageID,
			&i.ProductID,
			&i.VariantID,
			&i.ImageUrl,
			&i.ExternalID,
			&i.Primary,
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

const getImagesByVariantID = `-- name: GetImagesByVariantID :many
SELECT image_id, product_id, variant_id, image_url, external_id, "primary", created_at, updated_at FROM images WHERE variant_id = $1
`

func (q *Queries) GetImagesByVariantID(ctx context.Context, variantID pgtype.Int8) ([]Image, error) {
	rows, err := q.db.Query(ctx, getImagesByVariantID, variantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Image
	for rows.Next() {
		var i Image
		if err := rows.Scan(
			&i.ImageID,
			&i.ProductID,
			&i.VariantID,
			&i.ImageUrl,
			&i.ExternalID,
			&i.Primary,
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

const getPrimaryImageByProductID = `-- name: GetPrimaryImageByProductID :one
SELECT image_id, product_id, variant_id, image_url, external_id, "primary", created_at, updated_at FROM images WHERE product_id = $1 AND "primary" = TRUE LIMIT 1
`

func (q *Queries) GetPrimaryImageByProductID(ctx context.Context, productID pgtype.Int8) (Image, error) {
	row := q.db.QueryRow(ctx, getPrimaryImageByProductID, productID)
	var i Image
	err := row.Scan(
		&i.ImageID,
		&i.ProductID,
		&i.VariantID,
		&i.ImageUrl,
		&i.ExternalID,
		&i.Primary,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPrimaryImageByVariantID = `-- name: GetPrimaryImageByVariantID :one
SELECT image_id, product_id, variant_id, image_url, external_id, "primary", created_at, updated_at FROM images WHERE variant_id = $1 AND "primary" = TRUE LIMIT 1
`

func (q *Queries) GetPrimaryImageByVariantID(ctx context.Context, variantID pgtype.Int8) (Image, error) {
	row := q.db.QueryRow(ctx, getPrimaryImageByVariantID, variantID)
	var i Image
	err := row.Scan(
		&i.ImageID,
		&i.ProductID,
		&i.VariantID,
		&i.ImageUrl,
		&i.ExternalID,
		&i.Primary,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const setPrimaryImage = `-- name: SetPrimaryImage :exec
UPDATE images SET "primary" = TRUE WHERE image_id = $1
`

func (q *Queries) SetPrimaryImage(ctx context.Context, imageID int32) error {
	_, err := q.db.Exec(ctx, setPrimaryImage, imageID)
	return err
}

const unsetPrimaryImage = `-- name: UnsetPrimaryImage :exec
UPDATE
    images
SET
    "primary" = false
WHERE
    product_id = $1 AND "primary" = true
`

func (q *Queries) UnsetPrimaryImage(ctx context.Context, productID pgtype.Int8) error {
	_, err := q.db.Exec(ctx, unsetPrimaryImage, productID)
	return err
}

const updateImage = `-- name: UpdateImage :exec
UPDATE images 
SET 
    product_id = COALESCE($2, product_id),
    variant_id = COALESCE($3, variant_id), 
    image_url = COALESCE($4, image_url), 
    external_id = COALESCE($5, external_id) 
WHERE image_id = $1
`

type UpdateImageParams struct {
	ImageID    int32       `json:"image_id"`
	ProductID  pgtype.Int8 `json:"product_id"`
	VariantID  pgtype.Int8 `json:"variant_id"`
	ImageUrl   pgtype.Text `json:"image_url"`
	ExternalID pgtype.Text `json:"external_id"`
}

func (q *Queries) UpdateImage(ctx context.Context, arg UpdateImageParams) error {
	_, err := q.db.Exec(ctx, updateImage,
		arg.ImageID,
		arg.ProductID,
		arg.VariantID,
		arg.ImageUrl,
		arg.ExternalID,
	)
	return err
}
