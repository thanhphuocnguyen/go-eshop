-- name: CreateCart :one
INSERT INTO carts (user_id, session_id) VALUES ($1, $2) RETURNING *;

-- name: GetCart :one
SELECT * FROM carts
WHERE ((carts.user_id IS NOT NULL AND carts.user_id = $1) OR (carts.session_id IS NOT NULL AND carts.session_id = $2)) AND carts.order_id IS NULLORDER BY carts.updated_at DESC
ORDER BY carts.updated_at DESC
LIMIT 1;

-- name: RemoveProductFromCart :exec
DELETE FROM cart_items WHERE cart_id = $1 AND id = $2;

-- name: UpdateCartTimestamp :exec
UPDATE carts SET updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: CheckoutCart :exec
UPDATE carts SET order_id = $1 WHERE id = $2 RETURNING *;

-- Cart Item Section
-- name: CreateCartItem :one
INSERT INTO cart_items (id, cart_id, variant_id, quantity) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateCartItemQuantity :exec
UPDATE cart_items SET quantity = $1 WHERE id = $2 RETURNING *;

-- name: CountCartItems :one
SELECT COUNT(*) FROM cart_items WHERE id = $1;

-- name: GetCartItem :one
SELECT * FROM cart_items WHERE id = $1 AND cart_id = $2;

-- name: GetCartItemByProductVariantID :one
SELECT * FROM cart_items WHERE variant_id = $1 AND cart_id = $2;

-- name: GetCartItems :many
SELECT
    sqlc.embed(ci), 
    pv.id AS variant_id, pv.price, pv.stock, pv.sku, pv.stock as stock_qty,
    p.id AS product_id, p.name AS product_name, p.category_id,
    ci.id as cart_item_id, ci.quantity,
    av.id as attr_val_id, av.code AS attr_val_code, av.name as attr_val_name, a.name AS attr_name, a.id AS attr_id,
    pi.id AS image_id, pi.image_url AS image_url
FROM cart_items AS ci
JOIN product_variants AS pv ON pv.id = ci.variant_id
JOIN products AS p ON p.id = pv.product_id
JOIN variant_attribute_values AS vav ON vav.variant_id = pv.id
JOIN attribute_values AS av ON vav.attribute_value_id = av.id
JOIN attributes AS a ON av.attribute_id = a.id
LEFT JOIN product_images AS pi ON pi.product_id = p.id AND pi.is_primary = true
WHERE ci.cart_id = $1
ORDER BY ci.added_at, ci.id, pv.id DESC;

-- name: GetCartItemsForOrder :many
SELECT 
    sqlc.embed(ci), 
    pv.id AS variant_id, pv.price, pv.stock, pv.sku, pv.stock as stock_qty,
    p.name AS product_name, p.id AS product_id, p.category_id,
    av.code AS attr_val_code, av.name as attr_val_name, a.name AS attr_name
FROM cart_items AS ci
JOIN product_variants AS pv ON pv.id = ci.variant_id
JOIN products AS p ON p.id = pv.product_id
JOIN variant_attribute_values AS vav ON vav.variant_id = pv.id
JOIN attribute_values AS av ON vav.attribute_value_id = av.id
JOIN attributes AS a ON av.attribute_id = a.id
WHERE ci.cart_id = $1
ORDER BY ci.added_at, ci.id, pv.id DESC;

-- name: ClearCart :exec
DELETE FROM cart_items WHERE id = $1;

-- name: GetAvailableDiscountsForCart :many
SELECT d.id, d.code, d.description, d.discount_type, d.discount_value, d.starts_at, d.expires_at, dc.category_id, dp.product_id,
       SUM(CASE WHEN dp.product_id IS NOT NULL THEN 1 ELSE 0 END) AS product_count,
       SUM(CASE WHEN dc.category_id IS NOT NULL THEN 1 ELSE 0 END) AS category_count
FROM cart_items ci
JOIN product_variants pv ON ci.variant_id = pv.id
JOIN products p ON pv.product_id = p.id
LEFT JOIN discount_products dp ON p.id = dp.product_id
LEFT JOIN discount_categories dc ON p.category_id = dc.category_id
LEFT JOIN discounts d ON dp.discount_id = d.id OR dc.discount_id = d.id
WHERE ci.cart_id = $1
AND d.deleted_at IS NULL
AND (d.starts_at IS NULL OR d.starts_at <= NOW())
AND (d.expires_at IS NULL OR d.expires_at >= NOW())
GROUP BY d.id, p.id, dc.category_id, dp.product_id
HAVING COUNT(DISTINCT dp.product_id) > 0 OR COUNT(DISTINCT dc.category_id) > 0
ORDER BY d.starts_at DESC, d.expires_at DESC;