-- name: CreateOrder :one
INSERT INTO
    orders (
        user_id,
        payment_type,
        is_cod,
        user_address_id,
        cart_id,
        total_price
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
    sqlc.embed(orders), sqlc.embed(order_items), sqlc.embed(products), sqlc.embed(user_addresses)
FROM
    orders
JOIN
    order_items ON order_items.order_id = orders.id
JOIN
    products ON order_items.product_id = products.id
JOIN
    user_addresses ON orders.user_address_id = user_addresses.id
WHERE
    orders.id = $1;

-- name: ListOrders :many
SELECT
    sqlc.embed(orders), count(*) as total_items, sum(order_items.price) as total_price
FROM
    orders
JOIN order_items ON order_items.order_id = orders.id
WHERE
    user_id = $1
GROUP BY
    orders.id
ORDER BY
    orders.id
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
    user_address_id = coalesce(sqlc.narg('user_address_id'), user_address_id),
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