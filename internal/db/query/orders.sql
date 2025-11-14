-- name: CreateOrder :one
INSERT INTO orders (customer_id,customer_email,customer_name,customer_phone,total_price,shipping_address) VALUES ($1, $2, $3, $4,  $5, $6) RETURNING *;

-- name: GetOrder :one
SELECT
    orders.*,
    pm.id as payment_id,
    pm.status as payment_status,
    pm.amount as payment_amount,
    pmt.code as payment_method,
    pmt.id as payment_method_id,
    pm.gateway,
    pm.payment_intent_id,
    pm.created_at as payment_created_at,
    d.code,
    od.discount_amount
FROM orders
JOIN payments pm ON orders.id = pm.order_id
JOIN payment_methods pmt ON pm.payment_method_id = pmt.id
LEFT JOIN order_discounts od ON orders.id = od.order_id
LEFT JOIN discounts d ON od.discount_id = d.id
WHERE orders.id = $1
LIMIT 1;


-- name: GetOrderItems :many
SELECT
    oi.*,
    p.name as product_name, pi.image_url as image_url,
    rv.id as rating_id, rv.rating, rv.review_title, rv.review_content, rv.created_at as rating_created_at
FROM order_items oi
JOIN product_variants pv ON oi.variant_id = pv.id
JOIN products p ON pv.product_id = p.id
LEFT JOIN product_images AS pi ON pi.product_id = p.id AND pi.is_primary = true
LEFT JOIN product_ratings rv ON rv.order_item_id = oi.id
WHERE oi.order_id = $1;

-- name: GetOrders :many
SELECT ord.*, pm.status as payment_status, COUNT(oi.id) as total_items
FROM orders ord
LEFT JOIN order_items oi ON ord.id = oi.id
LEFT JOIN payments pm ON ord.id = pm.order_id
WHERE
    ord.customer_id = COALESCE(sqlc.narg('customer_id'), ord.customer_id) AND
    ord.status = COALESCE(sqlc.narg('status'), ord.status) AND
    ord.created_at >= COALESCE(sqlc.narg('start_date'), ord.created_at) AND
    ord.created_at <= COALESCE(sqlc.narg('end_date'), ord.created_at) AND
    pm.status = COALESCE(sqlc.narg('payment_status'), pm.status)
GROUP BY ord.id, pm.status
ORDER BY ord.created_at DESC
LIMIT $1
OFFSET $2;

-- name: UpdateOrder :one
UPDATE orders
SET
    status = coalesce(sqlc.narg('status'), status),
    confirmed_at = coalesce(sqlc.narg('confirmed_at'), confirmed_at),
    cancelled_at = coalesce(sqlc.narg('cancelled_at'), cancelled_at),
    delivered_at = coalesce(sqlc.narg('delivered_at'), delivered_at),
    updated_at = now()
WHERE id = sqlc.arg('id') RETURNING orders.id;

-- name: DeleteOrder :exec
DELETE FROM orders WHERE id = $1;

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, variant_id, quantity, price_per_unit_snapshot, variant_sku_snapshot, product_name_snapshot, line_total_snapshot, attributes_snapshot) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: CreateBulkOrderItems :copyfrom
INSERT INTO order_items (order_id, variant_id, quantity, price_per_unit_snapshot, variant_sku_snapshot, product_name_snapshot, line_total_snapshot, attributes_snapshot) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetOrderItemByID :one
SELECT oi.id as order_item_id, o.id as order_id, p.id as product_id, pv.id as variant_id, o.customer_id
FROM order_items oi
JOIN product_variants pv ON oi.variant_id = pv.id
JOIN products p ON pv.product_id = p.id
JOIN orders o ON oi.order_id = o.id
WHERE oi.id = $1
LIMIT 1;

-- name: GetOrderItemsByOrderID :many
SELECT oi.id as order_item_id, o.id as order_id, p.id as product_id, pv.id as variant_id, o.customer_id
FROM order_items oi
JOIN product_variants pv ON oi.variant_id = pv.id
JOIN products p ON pv.product_id = p.id
JOIN orders o ON oi.order_id = o.id
WHERE oi.order_id = $1
ORDER BY oi.id;

-- name: ListOrderItems :many
SELECT * FROM order_items WHERE order_id = $1 ORDER BY id LIMIT $2 OFFSET $3;

-- name: CountOrders :one
SELECT COUNT(*)
FROM orders ord
LEFT JOIN payments p ON ord.id = p.order_id
WHERE
    ord.status = COALESCE(sqlc.narg('status'), ord.status) AND
    customer_id = COALESCE(sqlc.narg('customer_id'), customer_id) AND
    p.status = COALESCE(sqlc.narg('payment_status'), p.status) AND
    ord.created_at >= COALESCE(sqlc.narg('start_date'), ord.created_at) AND
    ord.created_at <= COALESCE(sqlc.narg('end_date'), ord.created_at);