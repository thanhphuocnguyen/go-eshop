-- name: CreateProduct :one
INSERT INTO
    products (
        name,
        description,
        sku,
        stock,
        price,
        discount
    )
VALUES
    (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6
    )
RETURNING *;

-- name: GetProduct :one
SELECT
    *
FROM
    products
WHERE
    product_id = $1 AND
    archived = COALESCE(sqlc.narg('archived'), FALSE);

-- name: GetProductDetail :many
SELECT
    sqlc.embed(products),
    img.image_id AS image_id, img.image_url AS image_url, img.primary AS image_primary,
    pv.variant_id AS variant_id, pv.variant_name, pv.variant_sku, pv.variant_price, pv.variant_stock,
    a.attribute_id AS attribute_id, a.attribute_name,
    va.variant_attribute_id AS variant_attribute_id, va.attribute_value_id AS attribute_value_id,
    av.attribute_value_id AS attribute_value_id, av.attribute_value
FROM
    products
LEFT JOIN images AS img ON products.product_id = img.product_id
LEFT JOIN product_variants AS pv ON products.product_id = pv.product_id
LEFT JOIN variant_attributes AS va ON pv.variant_id = va.variant_id
LEFT JOIN attribute_values AS av ON va.attribute_value_id = av.attribute_value_id
LEFT JOIN attributes AS a ON av.attribute_id = a.attribute_id
WHERE
    products.product_id = $1 AND
    archived = COALESCE(sqlc.narg('archived'), false)
ORDER BY
    img.primary DESC;

-- name: ListProducts :many
SELECT
    p.*,
    img.image_id AS image_id, img.image_url AS image_url,
    COUNT(pv.variant_id) AS variant_count
FROM
    products as p
LEFT JOIN images AS img ON p.product_id = img.product_id AND img.primary = TRUE
LEFT JOIN product_variants AS pv ON p.product_id = pv.product_id
WHERE
    archived = COALESCE(sqlc.narg('archived'), archived) AND
    name ILIKE COALESCE(sqlc.narg('name'), name) AND
    sku ILIKE COALESCE(sqlc.narg('sku'), sku)
ORDER BY
    p.product_id
LIMIT 
    $1
OFFSET
    $2;

-- name: CountProducts :one
SELECT
    COUNT(*)
FROM
    products
WHERE
    archived = COALESCE(sqlc.narg('archived'), archived) AND
    name ILIKE COALESCE(sqlc.narg('name'), name) AND
    sku ILIKE COALESCE(sqlc.narg('sku'), sku);

-- name: UpdateProduct :one
UPDATE
    products
SET
    name = coalesce(sqlc.narg('name'), name),
    description = coalesce(sqlc.narg('description'), description),
    sku = coalesce(sqlc.narg('sku'), sku),
    stock = coalesce(sqlc.narg('stock'), stock),
    price = coalesce(sqlc.narg('price'), price),
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

-- name: UpdateProductStock :exec
UPDATE
    products
SET
    stock = stock + $2
WHERE
    product_id = $1
RETURNING *;

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

-- name: SeedProducts :copyfrom
INSERT INTO
    products (
        name,
        description,
        sku,
        stock,
        price,
        discount
    )
VALUES
    (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6
    );
