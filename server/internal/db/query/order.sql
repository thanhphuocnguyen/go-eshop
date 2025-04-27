-- name: CreateOrder :one
INSERT INTO orders (
    id,
    customer_id,
    user_address_id,
    customer_email,
    customer_name,
    customer_phone,
    total_price
)
VALUES 
    ($1,$2,$3,$4, $5, $6, $7)
RETURNING *;

-- name: GetOrder :one
SELECT
    *
FROM
    orders ord
WHERE
    id = $1
LIMIT 1;

-- name: GetOrderDetails :many
SELECT
    ord.*, 
    oi.quantity, oi.price_per_unit_snapshot, oi.variant_sku_snapshot, oi.product_name_snapshot, oi.line_total_snapshot,
    oi.attributes_snapshot, oi.id as order_item_id,
    p.name as product_name, p.id as variant_id,
    u_addr.street, u_addr.ward, u_addr.district, u_addr.city, 
    img.url as image_url,
    pm.status as payment_status, pm.id as payment_id, pm.amount as payment_amount, pm.payment_method, pm.payment_gateway, pm.refund_id
FROM
    orders ord
JOIN
    order_items oi ON oi.id = ord.id
JOIN
    products p ON oi.variant_id = p.id
JOIN
    user_addresses u_addr ON ord.user_address_id = u_addr.user_address_id
LEFT JOIN
    payments pm ON ord.id = pm.id
LEFT JOIN 
    image_assignments ia ON p.id = ia.entity_id AND ia.entity_type = 'product'
LEFT JOIN
    images img ON img.id = ia.image_id
WHERE
    ord.id = $1;

-- name: ListOrders :many
SELECT
    ord.*, pm.status as payment_status, COUNT(oit.id) as total_items
FROM
    orders ord
JOIN order_items oit ON ord.id = oit.id
LEFT JOIN payments pm ON ord.id = pm.id
WHERE
    customer_id = COALESCE(sqlc.narg('customer_id'), customer_id) AND
    ord.status = COALESCE(sqlc.narg('status'), ord.status) AND
    ord.created_at >= COALESCE(sqlc.narg('start_date'), ord.created_at) AND
    ord.created_at <= COALESCE(sqlc.narg('end_date'), ord.created_at)
GROUP BY ord.id, pm.status
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
    id = sqlc.arg('id')
RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM
    orders
WHERE
    id = $1;

-- name: CreateOrderItem :one
INSERT INTO
    order_items (id, order_id, variant_id, quantity, price_per_unit_snapshot, variant_sku_snapshot, product_name_snapshot, line_total_snapshot, attributes_snapshot)
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: CreateBulkOrderItems :copyfrom
INSERT INTO
    order_items (id, order_id, variant_id, quantity, price_per_unit_snapshot, variant_sku_snapshot, product_name_snapshot, line_total_snapshot, attributes_snapshot)
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8, $9);

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
    customer_id = $1 AND
    status = COALESCE(sqlc.narg('status'), status) AND
    created_at >= COALESCE(sqlc.narg('start_date'), created_at) AND
    created_at <= COALESCE(sqlc.narg('end_date'), created_at);