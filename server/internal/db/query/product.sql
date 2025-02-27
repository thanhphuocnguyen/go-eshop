-- name: CreateProduct :one
INSERT INTO products (
    product_id,
    name,
    description,
    brand_id,
    collection_id,
    category_id) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6)
RETURNING *;

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

-- name: GetProductVariantsByID :one
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
    pv.variant_id, pv.sku, pv.price, pv.stock_quantity, pv.discount,
    a.attribute_id AS attribute_id, a.name as attribute_name,
    va.variant_attribute_id AS variant_attribute_id, va.value as variant_attribute_value,
    img.image_id, img.image_url, img.product_id as img_product_id, img.variant_id as img_variant_id
FROM
    products p
LEFT JOIN product_variants AS pv ON p.product_id = pv.product_id
LEFT JOIN variant_attributes AS va ON pv.variant_id = va.variant_id
LEFT JOIN attributes AS a ON va.attribute_id = a.attribute_id
LEFT JOIN images AS img ON p.product_id = img.product_id
WHERE
    p.product_id = $1 AND
    archived = COALESCE(sqlc.narg('archived'), false)
ORDER BY
    pv.variant_id, a.attribute_id, va.variant_attribute_id DESC;

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
LEFT JOIN images AS img ON p.product_id = img.product_id
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
LEFT JOIN images AS img ON products.product_id = img.product_id
LEFT JOIN product_variants AS pv ON products.product_id = pv.product_id
WHERE
    products.product_id = $1 AND
    archived = COALESCE(sqlc.narg('archived'), false) AND
    pv.variant_id = COALESCE(sqlc.narg('variant_id'), pv.variant_id);

-- name: GetProductsByCategory :many
SELECT
    p.*,
    img.image_id AS image_id,
    img.image_url AS image_url,
    MIN(pv.price)::DECIMAL AS min_price,
    MAX(pv.price)::DECIMAL AS max_price,
    MAX(pv.price)::SMALLINT AS discount,
    COUNT(pv.variant_id) AS variant_count
FROM
    products AS p
JOIN product_variants AS pv ON p.product_id = pv.product_id
LEFT JOIN images AS img ON p.product_id = img.product_id
WHERE
    p.archived = COALESCE(sqlc.narg('archived'), archived) AND
    p.name ILIKE COALESCE(sqlc.narg('name'), name) AND
    pv.sku ILIKE COALESCE(sqlc.narg('sku'), sku) AND
    p.category_id = $1
GROUP BY
    p.product_id, img.image_id
ORDER BY
    p.product_id
LIMIT
    $2
OFFSET
    $3;

-- name: GetProductsByCollection :many
SELECT
    p.*,
    img.image_id AS image_id,
    img.image_url AS image_url,
    MIN(pv.price)::DECIMAL AS min_price,
    MAX(pv.price)::DECIMAL AS max_price,
    MAX(pv.price)::SMALLINT AS discount,
    COUNT(pv.variant_id) AS variant_count
FROM
    products AS p
JOIN product_variants AS pv ON p.product_id = pv.product_id
LEFT JOIN images AS img ON p.product_id = img.product_id
WHERE
    p.archived = COALESCE(sqlc.narg('archived'), archived) AND
    p.name ILIKE COALESCE(sqlc.narg('name'), name) AND
    pv.sku ILIKE COALESCE(sqlc.narg('sku'), sku) AND
    p.collection_id = $1
GROUP BY
    p.product_id, img.image_id
ORDER BY
    p.product_id
LIMIT
    $2
OFFSET
    $3;

-- name: GetProductsByBrand :many
SELECT
    p.*,
    img.image_id AS image_id,
    img.image_url AS image_url,
    MIN(pv.price)::DECIMAL AS min_price,
    MAX(pv.price)::DECIMAL AS max_price,
    MAX(pv.price)::SMALLINT AS discount,
    COUNT(pv.variant_id) AS variant_count
FROM
    products AS p
JOIN product_variants AS pv ON p.product_id = pv.product_id
LEFT JOIN images AS img ON p.product_id = img.product_id
WHERE
    p.archived = COALESCE(sqlc.narg('archived'), archived) AND
    p.name ILIKE COALESCE(sqlc.narg('name'), name) AND
    pv.sku ILIKE COALESCE(sqlc.narg('sku'), sku) AND
    p.brand_id = $1
GROUP BY
    p.product_id, img.image_id
ORDER BY
    p.product_id
LIMIT
    $2
OFFSET
    $3;

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
    brand_id = coalesce(sqlc.narg('brand_id'), brand_id),
    collection_id = coalesce(sqlc.narg('collection_id'), collection_id),
    category_id = coalesce(sqlc.narg('category_id'), category_id),
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
INSERT INTO products (product_id, category_id, collection_id, brand_id, name, description) VALUES ($1, $2, $3, $4, $5, $6);