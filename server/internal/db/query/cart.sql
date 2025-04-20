-- name: CreateCart :one
INSERT INTO carts (id, user_id, session_id) VALUES ($1, $2, $3) RETURNING *;

-- name: GetCart :one
SELECT * FROM carts WHERE user_id = $1 LIMIT 1;

-- name: RemoveProductFromCart :exec
DELETE FROM cart_items WHERE cart_id = $1 AND id = $2;

-- name: UpdateCart :exec
UPDATE carts SET updated_at = NOW() WHERE id = $1 RETURNING *;

-- Cart Item Section
-- name: GetCartItemByProductID :one
SELECT * FROM cart_items WHERE id = $1 AND variant_id = $2;

-- name: CreateCartItem :one
INSERT INTO cart_items (id, cart_id, variant_id, quantity) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateCartItemQuantity :exec
UPDATE cart_items SET quantity = $1 WHERE id = $2 RETURNING *;

-- name: CountCartItems :one
SELECT COUNT(*) FROM cart_items WHERE id = $1;

-- name: GetCartItem :one
SELECT * FROM cart_items WHERE id = $1;

-- name: GetCartItemWithProduct :one
SELECT 
    ci.*, 
    p.name AS product_name, p.brand_id, p.collection_id, p.category_id,
    pv.price, pv.stock, pv.sku
FROM cart_items ci
JOIN product_variants AS pv ON ci.variant_id = p.id
JOIN products AS p ON p.id = ci.product_id
WHERE ci.id = $1
ORDER BY ci.added_at, ci.id, pv.id DESC;

-- name: GetCartItemsByID :many
SELECT sqlc.embed(ci), sqlc.embed(p),
    av.value AS attribute_value, a.name AS attribute_name
FROM cart_items as ci
JOIN product_variants AS pv ON pv.id = ci.variant_id
JOIN products AS p ON p.id = pv.product_id
LEFT JOIN variant_attribute_values AS vav ON vav.id = pa.variant_id
LEFT JOIN attribute_values AS av ON vav.attribute_value_id = av.id
LEFT JOIN attributes AS a ON av.attribute_id = a.id
WHERE cart_id = $1
ORDER BY ci.id, ci.added_at DESC;

-- name: ClearCart :exec
DELETE FROM cart_items WHERE id = $1;