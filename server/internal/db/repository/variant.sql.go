// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: variant.sql

package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type CreateBulkVariantAttributesParams struct {
	VariantID   int64  `json:"variant_id"`
	AttributeID int32  `json:"attribute_id"`
	Value       string `json:"value"`
}

const createVariant = `-- name: CreateVariant :one
INSERT INTO product_variants (
    product_id,
    sku,
    price,
    stock_quantity,
    discount
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING variant_id, product_id, price, discount, stock_quantity, sku, created_at, updated_at
`

type CreateVariantParams struct {
	ProductID     int64          `json:"product_id"`
	Sku           pgtype.Text    `json:"sku"`
	Price         pgtype.Numeric `json:"price"`
	StockQuantity int32          `json:"stock_quantity"`
	Discount      int16          `json:"discount"`
}

func (q *Queries) CreateVariant(ctx context.Context, arg CreateVariantParams) (ProductVariant, error) {
	row := q.db.QueryRow(ctx, createVariant,
		arg.ProductID,
		arg.Sku,
		arg.Price,
		arg.StockQuantity,
		arg.Discount,
	)
	var i ProductVariant
	err := row.Scan(
		&i.VariantID,
		&i.ProductID,
		&i.Price,
		&i.Discount,
		&i.StockQuantity,
		&i.Sku,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createVariantAttribute = `-- name: CreateVariantAttribute :one

INSERT INTO variant_attributes (
    variant_id,
    attribute_id,
    value
) VALUES (
    $1, $2, $3
) RETURNING variant_attribute_id, variant_id, attribute_id, value, created_at, updated_at
`

type CreateVariantAttributeParams struct {
	VariantID   int64  `json:"variant_id"`
	AttributeID int32  `json:"attribute_id"`
	Value       string `json:"value"`
}

// Variant Attributes
func (q *Queries) CreateVariantAttribute(ctx context.Context, arg CreateVariantAttributeParams) (VariantAttribute, error) {
	row := q.db.QueryRow(ctx, createVariantAttribute, arg.VariantID, arg.AttributeID, arg.Value)
	var i VariantAttribute
	err := row.Scan(
		&i.VariantAttributeID,
		&i.VariantID,
		&i.AttributeID,
		&i.Value,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteVariant = `-- name: DeleteVariant :exec
DELETE FROM
    product_variants
WHERE
    variant_id = $1
`

func (q *Queries) DeleteVariant(ctx context.Context, variantID int64) error {
	_, err := q.db.Exec(ctx, deleteVariant, variantID)
	return err
}

const getVariantAttribute = `-- name: GetVariantAttribute :one
SELECT
    variant_attribute_id, variant_id, attribute_id, value, created_at, updated_at
FROM
    variant_attributes
WHERE
    variant_attribute_id = $1
LIMIT 1
`

func (q *Queries) GetVariantAttribute(ctx context.Context, variantAttributeID int32) (VariantAttribute, error) {
	row := q.db.QueryRow(ctx, getVariantAttribute, variantAttributeID)
	var i VariantAttribute
	err := row.Scan(
		&i.VariantAttributeID,
		&i.VariantID,
		&i.AttributeID,
		&i.Value,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getVariantAttributes = `-- name: GetVariantAttributes :many
SELECT
    variant_attribute_id, variant_id, attribute_id, value, created_at, updated_at
FROM
    variant_attributes
WHERE
    variant_id = $1
`

func (q *Queries) GetVariantAttributes(ctx context.Context, variantID int64) ([]VariantAttribute, error) {
	rows, err := q.db.Query(ctx, getVariantAttributes, variantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []VariantAttribute
	for rows.Next() {
		var i VariantAttribute
		if err := rows.Scan(
			&i.VariantAttributeID,
			&i.VariantID,
			&i.AttributeID,
			&i.Value,
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

const getVariantByID = `-- name: GetVariantByID :one
SELECT
    pv.variant_id, pv.product_id, pv.price, pv.discount, pv.stock_quantity, pv.sku, pv.created_at, pv.updated_at,
    p.name as product_name, p.product_id
FROM
    product_variants pv
JOIN
    products p ON pv.product_id = p.product_id
WHERE
    pv.variant_id = $1 AND pv.product_id = $2
LIMIT 1
`

type GetVariantByIDParams struct {
	VariantID int64 `json:"variant_id"`
	ProductID int64 `json:"product_id"`
}

type GetVariantByIDRow struct {
	VariantID     int64          `json:"variant_id"`
	ProductID     int64          `json:"product_id"`
	Price         pgtype.Numeric `json:"price"`
	Discount      int16          `json:"discount"`
	StockQuantity int32          `json:"stock_quantity"`
	Sku           pgtype.Text    `json:"sku"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	ProductName   string         `json:"product_name"`
	ProductID_2   int64          `json:"product_id_2"`
}

func (q *Queries) GetVariantByID(ctx context.Context, arg GetVariantByIDParams) (GetVariantByIDRow, error) {
	row := q.db.QueryRow(ctx, getVariantByID, arg.VariantID, arg.ProductID)
	var i GetVariantByIDRow
	err := row.Scan(
		&i.VariantID,
		&i.ProductID,
		&i.Price,
		&i.Discount,
		&i.StockQuantity,
		&i.Sku,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ProductName,
		&i.ProductID_2,
	)
	return i, err
}

const getVariantByProductID = `-- name: GetVariantByProductID :many
SELECT
    pv.variant_id, pv.product_id, pv.price, pv.discount, pv.stock_quantity, pv.sku, pv.created_at, pv.updated_at,
    a.name as attribute_name, a.attribute_id,
    va.variant_attribute_id, va.value as attribute_value
FROM
    product_variants pv
JOIN
    variant_attributes va ON pv.variant_id = va.variant_id
JOIN
    attributes a ON va.attribute_id = a.attribute_id
WHERE
    pv.product_id = $1
`

type GetVariantByProductIDRow struct {
	VariantID          int64          `json:"variant_id"`
	ProductID          int64          `json:"product_id"`
	Price              pgtype.Numeric `json:"price"`
	Discount           int16          `json:"discount"`
	StockQuantity      int32          `json:"stock_quantity"`
	Sku                pgtype.Text    `json:"sku"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	AttributeName      string         `json:"attribute_name"`
	AttributeID        int32          `json:"attribute_id"`
	VariantAttributeID int32          `json:"variant_attribute_id"`
	AttributeValue     string         `json:"attribute_value"`
}

func (q *Queries) GetVariantByProductID(ctx context.Context, productID int64) ([]GetVariantByProductIDRow, error) {
	rows, err := q.db.Query(ctx, getVariantByProductID, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetVariantByProductIDRow
	for rows.Next() {
		var i GetVariantByProductIDRow
		if err := rows.Scan(
			&i.VariantID,
			&i.ProductID,
			&i.Price,
			&i.Discount,
			&i.StockQuantity,
			&i.Sku,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.AttributeName,
			&i.AttributeID,
			&i.VariantAttributeID,
			&i.AttributeValue,
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

const getVariantDetails = `-- name: GetVariantDetails :many
SELECT
    pv.variant_id, pv.product_id, pv.price, pv.discount, pv.stock_quantity, pv.sku, pv.created_at, pv.updated_at,
    a.name as attribute_name, a.attribute_id,
    va.variant_attribute_id, va.value
FROM
    product_variants pv
JOIN
    variant_attributes va ON pv.variant_id = va.variant_id
JOIN
    attributes a ON va.attribute_id = a.attribute_id
WHERE
    pv.variant_id = $1
`

type GetVariantDetailsRow struct {
	VariantID          int64          `json:"variant_id"`
	ProductID          int64          `json:"product_id"`
	Price              pgtype.Numeric `json:"price"`
	Discount           int16          `json:"discount"`
	StockQuantity      int32          `json:"stock_quantity"`
	Sku                pgtype.Text    `json:"sku"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	AttributeName      string         `json:"attribute_name"`
	AttributeID        int32          `json:"attribute_id"`
	VariantAttributeID int32          `json:"variant_attribute_id"`
	Value              string         `json:"value"`
}

func (q *Queries) GetVariantDetails(ctx context.Context, variantID int64) ([]GetVariantDetailsRow, error) {
	rows, err := q.db.Query(ctx, getVariantDetails, variantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetVariantDetailsRow
	for rows.Next() {
		var i GetVariantDetailsRow
		if err := rows.Scan(
			&i.VariantID,
			&i.ProductID,
			&i.Price,
			&i.Discount,
			&i.StockQuantity,
			&i.Sku,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.AttributeName,
			&i.AttributeID,
			&i.VariantAttributeID,
			&i.Value,
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

const removeVariantAttribute = `-- name: RemoveVariantAttribute :exec
DELETE FROM
    variant_attributes
WHERE
    variant_attribute_id = $1
`

func (q *Queries) RemoveVariantAttribute(ctx context.Context, variantAttributeID int32) error {
	_, err := q.db.Exec(ctx, removeVariantAttribute, variantAttributeID)
	return err
}

const updateVariant = `-- name: UpdateVariant :one
UPDATE
    product_variants
SET
    sku = COALESCE($2, sku),
    price = COALESCE($3, price),
    stock_quantity = COALESCE($4, stock_quantity),
    discount = COALESCE($5, discount),
    updated_at = NOW()
WHERE
    variant_id = $1
RETURNING variant_id, product_id, price, discount, stock_quantity, sku, created_at, updated_at
`

type UpdateVariantParams struct {
	VariantID     int64          `json:"variant_id"`
	Sku           pgtype.Text    `json:"sku"`
	Price         pgtype.Numeric `json:"price"`
	StockQuantity pgtype.Int4    `json:"stock_quantity"`
	Discount      pgtype.Int2    `json:"discount"`
}

func (q *Queries) UpdateVariant(ctx context.Context, arg UpdateVariantParams) (ProductVariant, error) {
	row := q.db.QueryRow(ctx, updateVariant,
		arg.VariantID,
		arg.Sku,
		arg.Price,
		arg.StockQuantity,
		arg.Discount,
	)
	var i ProductVariant
	err := row.Scan(
		&i.VariantID,
		&i.ProductID,
		&i.Price,
		&i.Discount,
		&i.StockQuantity,
		&i.Sku,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateVariantAttribute = `-- name: UpdateVariantAttribute :one
UPDATE
    variant_attributes
SET
    value = COALESCE($2, value)
WHERE
    variant_attribute_id = $1
RETURNING variant_attribute_id, variant_id, attribute_id, value, created_at, updated_at
`

type UpdateVariantAttributeParams struct {
	VariantAttributeID int32       `json:"variant_attribute_id"`
	Value              pgtype.Text `json:"value"`
}

func (q *Queries) UpdateVariantAttribute(ctx context.Context, arg UpdateVariantAttributeParams) (VariantAttribute, error) {
	row := q.db.QueryRow(ctx, updateVariantAttribute, arg.VariantAttributeID, arg.Value)
	var i VariantAttribute
	err := row.Scan(
		&i.VariantAttributeID,
		&i.VariantID,
		&i.AttributeID,
		&i.Value,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
