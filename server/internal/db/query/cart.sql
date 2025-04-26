-- name: CreateCart :one
INSERT INTO carts (id, user_id, session_id) VALUES ($1, $2, $3) RETURNING *;

-- name: GetCart :one
SELECT *
FROM carts
WHERE (
    (carts.user_id IS NOT NULL AND carts.user_id = $1) OR 
    (carts.session_id IS NOT NULL AND carts.session_id = $2)
) AND carts.order_id IS NULL
ORDER BY carts.updated_at DESC
LIMIT 1;

-- name: RemoveProductFromCart :exec
DELETE FROM cart_items WHERE cart_id = $1 AND id = $2;

-- name: UpdateCartTimestamp :exec
UPDATE carts SET updated_at = NOW() WHERE id = $1 RETURNING *;

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
    p.id AS product_id, p.name AS product_name,
    ci.id as cart_item_id, ci.quantity,
    av.id as attr_val_id, av.code AS attr_val_code, av.name as attr_val_name, a.name AS attr_name, a.id AS attr_id,
    i.id AS image_id, i.url AS image_url
FROM cart_items AS ci
JOIN product_variants AS pv ON pv.id = ci.variant_id
JOIN products AS p ON p.id = pv.product_id
JOIN variant_attribute_values AS vav ON vav.variant_id = pv.id
JOIN attribute_values AS av ON vav.attribute_value_id = av.id
JOIN attributes AS a ON av.attribute_id = a.id
LEFT JOIN image_assignments AS ia ON ia.entity_id = pv.id AND ia.entity_type = 'variant'
LEFT JOIN images AS i ON i.id = ia.image_id
WHERE ci.cart_id = $1
ORDER BY ci.added_at, ci.id, pv.id DESC;

-- name: ClearCart :exec
DELETE FROM cart_items WHERE id = $1;