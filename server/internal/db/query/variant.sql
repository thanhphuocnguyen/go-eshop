-- name: CreateVariant :one
INSERT INTO product_variants (
    product_id,
    sku,
    price,
    stock_quantity,
    discount
) VALUES (
    $1, $2, $3, $4, $5
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
    pv.variant_id = $1
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
    updated_at = NOW()
WHERE
    variant_id = $1
RETURNING *;


-- name: DeleteVariant :exec
DELETE FROM
    product_variants
WHERE
    variant_id = $1;