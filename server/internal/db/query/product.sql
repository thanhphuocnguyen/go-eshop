-- name: CreateProduct :one
INSERT INTO products (name, description) VALUES ($1, $2) RETURNING *;

-- name: GetProductByID :one
SELECT
    products.*,
    COUNT(pv.variant_id) AS variant_count
FROM
    products
JOIN product_variants AS pv ON products.product_id = pv.product_id
WHERE
    products.product_id = $1 AND
    archived = COALESCE(sqlc.narg('archived'), false)
GROUP BY
    products.product_id;

-- name: GetProductWithVariantByID :one
SELECT
    *
FROM
    products
JOIN
    product_variants pv ON products.product_id = pv.product_id
WHERE
    products.product_id = $1 AND
    archived = COALESCE(sqlc.narg('archived'), false) AND
    pv.variant_id = COALESCE(sqlc.narg('variant_id'), pv.variant_id);

-- name: GetProductDetail :many
SELECT
    sqlc.embed(p),
    pv.variant_id, pv.sku, pv.price, pv.stock_quantity,
    a.attribute_id AS attribute_id, a.name as attribute_name,
    va.variant_attribute_id AS variant_attribute_id, va.value as variant_attribute_value,
    img.image_id, img.product_id as img_product_id, img.variant_id as img_variant_id, img.image_url, img.primary AS image_primary
FROM
    products p
JOIN product_variants AS pv ON p.product_id = pv.product_id
JOIN variant_attributes AS va ON pv.variant_id = va.variant_id
JOIN attributes AS a ON va.attribute_id = a.attribute_id
LEFT JOIN images AS img ON p.product_id = img.product_id
WHERE
    p.product_id = $1 AND
    archived = COALESCE(sqlc.narg('archived'), false)
ORDER BY
    pv.variant_id, a.attribute_id, va.variant_attribute_id, img.primary DESC;

-- name: GetProducts :many
SELECT
    p.*,
    img.image_id AS image_id, img.image_url AS image_url,
    MIN(pv.price)::DECIMAL AS min_price,
    MAX(pv.price)::DECIMAL AS max_price,
    COUNT(pv.variant_id) AS variant_count
FROM
    products as p
JOIN product_variants AS pv ON p.product_id = pv.product_id
LEFT JOIN images AS img ON p.product_id = img.product_id AND img.primary = TRUE
WHERE
    archived = COALESCE(sqlc.narg('archived'), archived) AND
    name ILIKE COALESCE(sqlc.narg('name'), name) AND
    sku ILIKE COALESCE(sqlc.narg('sku'), sku)
GROUP BY
    p.product_id, img.image_id
ORDER BY
    p.product_id
LIMIT 
    $1
OFFSET
    $2;

-- name: GetProductWithImage :one
SELECT
    products.*,
    img.image_id AS image_id,
    img.image_url AS image_url
FROM
    products
LEFT JOIN images AS img ON products.product_id = img.product_id AND img.primary = TRUE
LEFT JOIN product_variants AS pv ON products.product_id = pv.product_id
WHERE
    products.product_id = $1 AND
    archived = COALESCE(sqlc.narg('archived'), false) AND
    pv.variant_id = COALESCE(sqlc.narg('variant_id'), pv.variant_id);

-- name: CountProducts :one
SELECT
    COUNT(*)
FROM
    products
WHERE
    archived = COALESCE(sqlc.narg('archived'), archived) AND
    name ILIKE COALESCE(sqlc.narg('name'), name);

-- name: UpdateProduct :one
UPDATE
    products
SET
    name = coalesce(sqlc.narg('name'), name),
    description = coalesce(sqlc.narg('description'), description),
    updated_at = NOW()
WHERE
    product_id = sqlc.arg('product_id')
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM
    products
WHERE
    product_id = $1;

-- name: ArchiveProduct :exec
UPDATE
    products
SET
    archived = true,
    updated_at = NOW()
WHERE
    product_id = $1;

-- name: AddBulkProducts :copyfrom
INSERT INTO
    products (
        name,
        description
    )
VALUES
    (
        $1,
        $2
    );