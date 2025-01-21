-- name: CreateOrder :one
INSERT INTO orders (order_id,user_id,user_address_id,total_price) VALUES ($1,$2,$3,$4) RETURNING *;

-- name: GetOrder :one
SELECT
    *
FROM
    orders ord
WHERE
    order_id = $1
LIMIT 1;

-- name: GetOrderDetails :many
SELECT
    ord.*, 
    oit.quantity, oit.price as item_price, oit.order_item_id as order_item_id,
    p.name as product_name, p.product_id as product_id,
    u_addr.street, u_addr.ward, u_addr.district, u_addr.city, 
    images.image_url,
    pm.status as payment_status, pm.payment_id as payment_id, pm.amount as payment_amount, pm.payment_method as payment_method, pm.payment_gateway as payment_gateway, pm.refund_id as refund_id,
    pv.variant_id
FROM
    orders ord
LEFT JOIN
    payments pm ON ord.order_id = pm.order_id
LEFT JOIN
    order_items oit ON oit.order_id = ord.order_id
LEFT JOIN
    products p ON oit.product_id = p.product_id
LEFT JOIN 
    product_variants AS pv ON oit.variant_id = p.variant_id
LEFT JOIN 
    images ON p.product_id = images.product_id AND images.primary = true
LEFT JOIN
    user_addresses u_addr ON ord.user_address_id = u_addr.user_address_id
WHERE
    ord.order_id = $1;

-- name: ListOrders :many
SELECT
    ord.*, pm.status as payment_status, COUNT(oit.order_item_id) as total_items
FROM
    orders ord
LEFT JOIN payments pm ON ord.order_id = pm.order_id
LEFT JOIN order_items oit ON ord.order_id = oit.order_id
WHERE
    user_id = COALESCE(sqlc.narg('user_id'), user_id) AND
    ord.status = COALESCE(sqlc.narg('status'), ord.status) AND
    ord.created_at >= COALESCE(sqlc.narg('start_date'), ord.created_at) AND
    ord.created_at <= COALESCE(sqlc.narg('end_date'), ord.created_at)
GROUP BY ord.order_id, pm.status
ORDER BY
    ord.created_at DESC
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
    order_id = sqlc.arg('order_id')
RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM
    orders
WHERE
    order_id = $1;

-- name: CreateOrderItem :one
INSERT INTO
    order_items (
        product_id,
        variant_id,
        order_id,
        quantity,
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

-- name: ListOrderItems :many
SELECT
    *
FROM
    order_items
WHERE
    order_id = $1
ORDER BY
    order_item_id
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