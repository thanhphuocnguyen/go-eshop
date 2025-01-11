-- name: GetVariantByID :one
SELECT
    *
FROM
    product_variants pv
WHERE
    pv.variant_id = $1
LIMIT 1;

-- name: GetVariantDetails :many
SELECT
    pv.*,
    a.attribute_name,
    av.attribute_value,
    va.variant_attribute_id,
    av.attribute_value_id,
    a.attribute_id
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


-- name: UpdateVariant :one
UPDATE
    product_variants
SET
    variant_name = COALESCE(sqlc.narg('variant_name'), variant_name),
    variant_sku = COALESCE(sqlc.narg('variant_sku'), variant_sku),
    variant_price = COALESCE(sqlc.narg('variant_price'), variant_price),
    variant_stock = COALESCE(sqlc.narg('variant_stock'), variant_stock)
WHERE
    variant_id = $1
RETURNING *;


-- name: DeleteVariant :exec
DELETE FROM
    product_variants
WHERE
    variant_id = $1;

-- name: ListVariants :many
SELECT
    p.*,
    pv.variant_sku, pv.variant_price, pv.variant_stock
FROM
    product_variants pv
JOIN
    products p ON pv.product_id = p.product_id
ORDER BY
    pv.variant_id
LIMIT $1
OFFSET $2;


-- name: CreateVariantAttribute :one
INSERT INTO variant_attributes (
    variant_id,
    attribute_value_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: AddVariantAttributes :copyfrom
INSERT INTO variant_attributes (
    variant_id,
    attribute_value_id
) VALUES (
    $1, $2
);

-- name: DeleteVariantAttribute :exec
DELETE FROM
    variant_attributes
WHERE
    variant_attribute_id = $1;

-- name: UpdateVariantAttribute :one
UPDATE
    variant_attributes
SET
    variant_id = $2,
    attribute_value_id = $3
WHERE
    variant_attribute_id = $1
RETURNING *;


-- name: GetVariantAttributeByID :one
SELECT
    *
FROM
    variant_attributes
WHERE
    variant_attribute_id = $1
LIMIT 1;
