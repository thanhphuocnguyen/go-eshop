-- name: CreateVariant :one
INSERT INTO product_variants (
    variant_id,
    product_id,
    sku,
    price,
    stock_quantity,
    discount
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetVariantByID :one
SELECT
    pv.*,
    p.name as product_name, p.product_id
FROM
    product_variants pv
JOIN
    products p ON pv.product_id = p.product_id
WHERE
    pv.variant_id = $1 AND pv.product_id = $2
LIMIT 1;

-- name: GetVariantByProductID :many
SELECT
    pv.*,
    a.name as attribute_name, a.attribute_id,
    va.variant_attribute_id, va.value as attribute_value
FROM
    product_variants pv
JOIN
    variant_attributes va ON pv.variant_id = va.variant_id
JOIN
    attributes a ON va.attribute_id = a.attribute_id
WHERE
    pv.product_id = $1;

-- name: GetVariantDetails :many
SELECT
    pv.*,
    a.name as attribute_name, a.attribute_id,
    va.variant_attribute_id, va.value
FROM
    product_variants pv
JOIN
    variant_attributes va ON pv.variant_id = va.variant_id
JOIN
    attributes a ON va.attribute_id = a.attribute_id
WHERE
    pv.variant_id = $1;

-- name: UpdateVariant :one
UPDATE
    product_variants
SET
    sku = COALESCE(sqlc.narg('sku'), sku),
    price = COALESCE(sqlc.narg('price'), price),
    stock_quantity = COALESCE(sqlc.narg('stock_quantity'), stock_quantity),
    discount = COALESCE(sqlc.narg('discount'), discount),
    updated_at = NOW()
WHERE
    variant_id = $1
RETURNING *;


-- name: DeleteVariant :exec
DELETE FROM
    product_variants
WHERE
    variant_id = $1;

-- Variant Attributes

-- name: CreateVariantAttribute :one
INSERT INTO variant_attributes (
    variant_id,
    attribute_id,
    value
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetVariantAttribute :one
SELECT
    *
FROM
    variant_attributes
WHERE
    variant_attribute_id = $1
LIMIT 1;

-- name: GetVariantAttributes :many
SELECT
    *
FROM
    variant_attributes
WHERE
    variant_id = $1;

-- name: UpdateVariantAttribute :one
UPDATE
    variant_attributes
SET
    value = COALESCE(sqlc.narg('value'), value)
WHERE
    variant_attribute_id = $1
RETURNING *;

-- name: RemoveVariantAttribute :exec
DELETE FROM
    variant_attributes
WHERE
    variant_attribute_id = $1;

-- name: CreateBulkVariantAttributes :copyfrom
INSERT INTO variant_attributes (
    variant_id,
    attribute_id,
    value
) VALUES
    ($1, $2, $3);