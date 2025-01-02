-- name: CreateProduct :one
INSERT INTO
    products (
        name,
        description,
        sku,
        stock,
        price
    )
VALUES
    (
        $1,
        $2,
        $3,
        $4,
        $5
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
    img.image_id AS image_id,
    img.image_url AS image_url,
    img.primary AS image_primary
FROM
    products
LEFT JOIN images AS img ON products.product_id = img.product_id
WHERE
    products.product_id = $1 AND
    archived = COALESCE(sqlc.narg('archived'), false)
ORDER BY
    img.primary DESC;

-- name: ListProducts :many
SELECT
    products.*,
    img.image_id AS image_id,
    img.image_url AS image_url
FROM
    products
LEFT JOIN images AS img ON products.product_id = img.product_id AND img.primary = TRUE
WHERE
    archived = COALESCE(sqlc.narg('archived'), archived) AND
    name ILIKE COALESCE(sqlc.narg('name'), name) AND
    sku ILIKE COALESCE(sqlc.narg('sku'), sku)
ORDER BY
    products.product_id
LIMIT $1
OFFSET $2;

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
WHERE
    products.product_id = $1 AND
    archived = COALESCE(sqlc.narg('archived'), false);
