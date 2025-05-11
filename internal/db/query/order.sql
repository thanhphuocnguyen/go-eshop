-- name: CreateOrder :one
INSERT INTO orders (
    customer_id,
    customer_email,
    customer_name,
    customer_phone,
    total_price,
    shipping_address
)
VALUES 
    ($1, $2, $3, $4,  $5, $6)
RETURNING *;

-- name: GetOrder :one
SELECT
    *
FROM
    orders
WHERE
    id = $1
LIMIT 1;


-- name: GetOrderProducts :many
SELECT
    oi.*,
    p.name as product_name,
    i.url as image_url
FROM
    order_items oi
JOIN
    product_variants pv ON oi.variant_id = pv.id
JOIN
    products p ON pv.product_id = p.id
LEFT JOIN image_assignments AS ia ON ia.entity_id = pv.id AND ia.entity_type = 'variant'
LEFT JOIN images AS i ON i.id = ia.image_id
WHERE
    oi.order_id = $1;

-- name: GetOrders :many
SELECT
    ord.*, pm.status as payment_status, COUNT(oi.id) as total_items
FROM
    orders ord
LEFT JOIN order_items oi ON ord.id = oi.id
LEFT JOIN payments pm ON ord.id = pm.order_id
WHERE
    ord.customer_id = COALESCE(sqlc.narg('customer_id'), ord.customer_id) AND
    ord.status = COALESCE(sqlc.narg('status'), ord.status) AND
    ord.created_at >= COALESCE(sqlc.narg('start_date'), ord.created_at) AND
    ord.created_at <= COALESCE(sqlc.narg('end_date'), ord.created_at) AND
    pm.status = COALESCE(sqlc.narg('payment_status'), pm.status)
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
    order_items (order_id, variant_id, quantity, price_per_unit_snapshot, variant_sku_snapshot, product_name_snapshot, line_total_snapshot, attributes_snapshot)
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: CreateBulkOrderItems :copyfrom
INSERT INTO
    order_items (order_id, variant_id, quantity, price_per_unit_snapshot, variant_sku_snapshot, product_name_snapshot, line_total_snapshot, attributes_snapshot)
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetOrderItemByID :one
SELECT
    oi.id as order_item_id,
    o.id as order_id,
    p.id as product_id,
    pv.id as variant_id,
    o.customer_id
FROM
    order_items oi
JOIN
    product_variants pv ON oi.variant_id = pv.id
JOIN
    products p ON pv.product_id = p.id
JOIN
    orders o ON oi.order_id = o.id
WHERE
    oi.id = $1
LIMIT 1;

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
    orders ord
LEFT JOIN payments p ON ord.id = p.order_id
WHERE
    ord.status = COALESCE(sqlc.narg('status'), ord.status) AND
    customer_id = COALESCE(sqlc.narg('customer_id'), customer_id) AND
    p.status = COALESCE(sqlc.narg('payment_status'), p.status) AND
    ord.created_at >= COALESCE(sqlc.narg('start_date'), ord.created_at) AND
    ord.created_at <= COALESCE(sqlc.narg('end_date'), ord.created_at);