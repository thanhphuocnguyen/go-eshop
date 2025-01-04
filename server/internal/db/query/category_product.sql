-- name: AddProductToCollection :one
INSERT INTO
    category_products (category_id, product_id, sort_order)
VALUES
    ($1, $2, $3)
RETURNING *;

-- name: RemoveProductFromCollection :one
DELETE FROM
    category_products
WHERE
    category_id = $1
    AND product_id = $2
RETURNING *;

-- name: GetCollectionProducts :many
SELECT
    p.*
FROM
    products p
    JOIN category_products cp ON p.product_id = cp.product_id
WHERE
    cp.category_id = $1;
