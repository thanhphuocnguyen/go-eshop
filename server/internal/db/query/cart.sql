-- name: CreateCart :one
INSERT INTO 
    carts (cart_id, user_id) 
VALUES 
    ($1, $2) 
RETURNING *;

-- name: GetCart :one
SELECT * FROM carts WHERE user_id = $1 LIMIT 1;

-- name: RemoveProductFromCart :exec
DELETE FROM cart_items WHERE cart_id = $1 AND cart_item_id = $2;

-- name: UpdateCart :exec
UPDATE carts SET updated_at = NOW() WHERE cart_id = $1 RETURNING *;

-- Cart Item Section

-- name: CreateCartItem :one
INSERT INTO cart_items 
    (cart_id, product_id, variant_id, quantity) 
VALUES 
    ($1, $2, $3, $4) 
RETURNING *;

-- name: UpdateCartItemQuantity :exec
UPDATE cart_items SET quantity = $1 WHERE cart_item_id = $2 RETURNING *;

-- name: CountCartItem :one
SELECT COUNT(*) FROM cart_items WHERE cart_id = $1;

-- name: GetCartItem :one
SELECT * FROM cart_items WHERE cart_item_id = $1;

-- name: GetCartItemWithProduct :one
SELECT ci.*, 
    p.name AS product_name,
    pv.price, pv.stock_quantity, pv.sku,
    img.image_url
FROM cart_items ci
JOIN products AS p ON ci.product_id = p.product_id
JOIN product_variants AS pv ON ci.variant_id = p.variant_id
JOIN variant_attributes AS va ON pv.variant_id = va.variant_id
LEFT JOIN images as img ON p.product_id = img.product_id
WHERE ci.cart_item_id = $1
ORDER BY ci.created_at DESC, ci.cart_item_id DESC, p.product_id, pv.variant_id, va.attribute_id;

-- name: GetCartItemByVariantID :one
SELECT * 
FROM 
    cart_items 
WHERE 
    variant_id = $1
LIMIT 1;

-- name: GetCartItemsByID :many
SELECT cart_items.*, 
    p.name AS product_name,
    img.image_url AS image_url,
    pv.price, pv.stock_quantity, pv.discount, pv.sku,
    va.variant_attribute_id, va.value AS attribute_value,
    a.name AS attribute_name
FROM cart_items
JOIN products AS p ON cart_items.product_id = p.product_id
JOIN product_variants AS pv ON cart_items.variant_id = pv.variant_id
JOIN variant_attributes AS va ON pv.variant_id = va.variant_id
JOIN attributes AS a ON va.attribute_id = a.attribute_id
LEFT JOIN images as img ON p.product_id = img.product_id
WHERE cart_id = $1
ORDER BY cart_items.created_at DESC, cart_items.cart_item_id DESC, p.product_id, pv.variant_id, va.attribute_id;

-- name: ClearCart :exec
DELETE FROM cart_items WHERE cart_id = $1;