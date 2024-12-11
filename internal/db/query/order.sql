-- name: CreateOrder :one
INSERT INTO
    orders (
        user_id,
        payment_type,
        is_cod
    )
VALUES
    (
        $1,
        $2,
        $3
    )
RETURNING *;

-- name: GetOrder :one
SELECT
    *
FROM
    orders
WHERE
    id = $1
LIMIT 1;

-- name: GetOrderDetails :many
SELECT
    sqlc.embed(orders), sqlc.embed(order_items), sqlc.embed(products)
FROM
    orders
JOIN
    order_items ON order_items.order_id = orders.id
JOIN
    products ON order_items.product_id = products.id
WHERE
    orders.id = $1;

-- name: ListOrders :many
SELECT
    *
FROM
    orders
WHERE
    user_id = $1
ORDER BY
    id
LIMIT $2
OFFSET $3;

-- name: UpdateOrder :one
UPDATE
    orders
SET
    status = coalesce(sqlc.narg('status'), status),
    shipping_id = coalesce(sqlc.narg('shipping_id'), shipping_id),
    payment_status = coalesce(sqlc.narg('payment_status'), payment_status),
    confirmed_at = coalesce(sqlc.narg('confirmed_at'), confirmed_at),
    cancelled_at = coalesce(sqlc.narg('cancelled_at'), cancelled_at),
    delivered_at = coalesce(sqlc.narg('delivered_at'), delivered_at),
    updated_at = sqlc.arg('updated_at')
WHERE
    id = sqlc.arg('id')
RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM
    orders
WHERE
    id = $1;

-- name: CreateOrderItem :one
INSERT INTO
    order_items (
        product_id,
        order_id,
        quantity,
        price
    )
VALUES
    (
        $1,
        $2,
        $3,
        $4
    )
RETURNING *;

-- name: ListOrderItems :many
SELECT
    *
FROM
    order_items
WHERE
    order_id = $1
ORDER BY
    id
LIMIT $2
OFFSET $3;