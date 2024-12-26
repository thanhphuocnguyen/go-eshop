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
    id = $1 AND
    archived = COALESCE(sqlc.narg('archived'), archived);

-- name: GetProductDetail :many
SELECT
    sqlc.embed(products),
    img.image_id AS image_id,
    img.image_url AS image_url,
    img.is_primary AS image_is_primary
FROM
    products
LEFT JOIN images AS img ON products.id = img.product_id
WHERE
    products.id = $1 AND
    archived = COALESCE(sqlc.narg('archived'), archived)
ORDER BY
    img.is_primary DESC;

-- name: ListProducts :many
SELECT
    products.*,
    img.image_id AS image_id,
    img.image_url AS image_url,
    img.is_primary AS image_is_primary
FROM
    products
LEFT JOIN images AS img ON products.id = img.product_id AND img.is_primary = TRUE
WHERE
    archived = COALESCE(sqlc.narg('archived'), archived) AND
    name ILIKE COALESCE(sqlc.narg('name'), name) AND
    sku ILIKE COALESCE(sqlc.narg('sku'), sku)
ORDER BY
    id
LIMIT $1
OFFSET $2;


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
    id = sqlc.arg('id')
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM
    products
WHERE
    id = $1;

-- name: ArchiveProduct :exec
UPDATE
    products
SET
    archived = true
WHERE
    id = $1
RETURNING *;

-- name: UpdateProductStock :exec
UPDATE
    products
SET
    stock = stock + $2
WHERE
    id = $1
RETURNING *;