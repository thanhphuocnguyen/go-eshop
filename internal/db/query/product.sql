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
    id = $1
LIMIT 1;

-- name: ListProducts :many
SELECT
    *
FROM
    products
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
    image_url = coalesce(sqlc.narg('image_url'), image_url),
    stock = coalesce(sqlc.narg('stock'), stock),
    price = coalesce(sqlc.narg('price'), price),
    updated_at = sqlc.arg('updated_at')
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

-- name: UpdateProductImage :exec
UPDATE
    products
SET
    image_url = $2
WHERE
    id = $1
RETURNING *;