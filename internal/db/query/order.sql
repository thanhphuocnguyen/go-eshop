-- name: CreateOrder :one
INSERT INTO
    orders (
        user_id,
        user_address_id,
        total_price
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
    orders.*, 
    order_items.quantity, order_items.price as item_price, order_items.id as order_item_id,
    products.name as product_name, products.id as product_id,
    user_addresses.street, user_addresses.ward, user_addresses.district, user_addresses.city, 
    images.image_url,
    payments.status as payment_status
FROM
    orders
LEFT JOIN
    payments ON orders.id = payments.order_id
LEFT JOIN
    order_items ON order_items.order_id = orders.id
LEFT JOIN
    products ON order_items.product_id = products.id
LEFT JOIN 
    images ON products.id = images.product_id AND images.is_primary = true
LEFT JOIN
    user_addresses ON orders.user_address_id = user_addresses.id
WHERE
    orders.id = $1;

-- name: ListOrders :many
SELECT
    orders.*, payments.status as payment_status, COUNT(order_items.id) as total_items
FROM
    orders
LEFT JOIN payments ON orders.id = payments.order_id
LEFT JOIN order_items ON orders.id = order_items.order_id
WHERE
    user_id = COALESCE(sqlc.arg('user_id'), user_id) AND
    orders.status = COALESCE(sqlc.narg('status'), orders.status) AND
    orders.created_at >= COALESCE(sqlc.narg('start_date'), orders.created_at) AND
    orders.created_at <= COALESCE(sqlc.narg('end_date'), orders.created_at)
GROUP BY orders.id, payments.status
ORDER BY
    orders.created_at DESC
LIMIT $1
OFFSET $2;

-- name: UpdateOrder :one
UPDATE
    orders
SET
    status = coalesce(sqlc.narg('status'), status),
    confirmed_at = coalesce(sqlc.narg('confirmed_at'), confirmed_at),
    cancelled_at = coalesce(sqlc.narg('cancelled_at'), cancelled_at),
    delivered_at = coalesce(sqlc.narg('delivered_at'), delivered_at),
    user_address_id = coalesce(sqlc.narg('user_address_id'), user_address_id),
    updated_at = now()
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

-- name: CountOrders :one
SELECT
    COUNT(*)
FROM
    orders
WHERE
    user_id = $1 AND
    status = COALESCE(sqlc.narg('status'), status) AND
    created_at >= COALESCE(sqlc.narg('start_date'), created_at) AND
    created_at <= COALESCE(sqlc.narg('end_date'), created_at);