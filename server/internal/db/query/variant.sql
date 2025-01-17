-- name: CreateVariant :one
INSERT INTO product_variants (
    product_id,
    variant_name,
    variant_sku,
    variant_price,
    variant_stock
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetVariantByID :one
SELECT
    *
FROM
    product_variants pv
WHERE
    pv.variant_id = $1
LIMIT 1;

-- name: GetVariantByProductID :many
SELECT
    pv.*,
    a.attribute_name, a.attribute_id,
    av.attribute_value_id, av.attribute_value, av.color,
    va.variant_attribute_id
FROM
    product_variants pv
JOIN
    variant_attributes va ON pv.variant_id = va.variant_id
JOIN
    attribute_values av ON va.attribute_value_id = av.attribute_value_id
JOIN
    attributes a ON av.attribute_id = a.attribute_id
WHERE
    pv.product_id = $1;

-- name: GetVariantDetails :many
SELECT
    pv.*,
    a.attribute_name, a.attribute_id,
    av.attribute_value_id, av.attribute_value, av.color,
    va.variant_attribute_id
FROM
    product_variants pv
JOIN
    variant_attributes va ON pv.variant_id = va.variant_id
JOIN
    attribute_values av ON va.attribute_value_id = av.attribute_value_id
JOIN
    attributes a ON av.attribute_id = a.attribute_id
WHERE
    pv.variant_id = $1;

-- name: UpdateVariant :one
UPDATE
    product_variants
SET
    variant_name = COALESCE(sqlc.narg('variant_name'), variant_name),
    variant_sku = COALESCE(sqlc.narg('variant_sku'), variant_sku),
    variant_price = COALESCE(sqlc.narg('variant_price'), variant_price),
    variant_stock = COALESCE(sqlc.narg('variant_stock'), variant_stock),
    updated_at = NOW()
WHERE
    variant_id = $1
RETURNING *;


-- name: DeleteVariant :exec
DELETE FROM
    product_variants
WHERE
    variant_id = $1;


---- Variant Attributes ----

-- name: CreateVariantAttribute :one
INSERT INTO variant_attributes (
    variant_id,
    attribute_value_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: CreateBulkVariantAttribute :copyfrom
INSERT INTO variant_attributes (
    variant_id,
    attribute_value_id
) VALUES (
    $1, $2
);

-- name: GetVariantAttributeByID :one
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
    variant_id = $1
ORDER BY
    variant_attribute_id;

-- name: UpdateVariantAttribute :one
UPDATE
    variant_attributes
SET
    variant_id = $2,
    attribute_value_id = $3,
    updated_at = NOW()
WHERE
    variant_attribute_id = $1
RETURNING *;


-- name: DeleteVariantAttribute :exec
DELETE FROM
    variant_attributes
WHERE
    attribute_value_id = $1 AND variant_id = $2;