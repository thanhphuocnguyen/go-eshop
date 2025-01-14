-- name: AddProductToCart :one
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
SELECT cart_items.*, p.name AS product_name, 
    p.price AS product_price, p.stock AS product_stock,
    pv.variant_price, pv.variant_stock,
    img.image_url
FROM cart_items
JOIN products AS p ON cart_items.product_id = p.product_id
JOIN product_variants AS pv ON cart_items.variant_id = p.variant_id
LEFT JOIN images as img ON p.product_id = img.product_id AND img.primary = true
WHERE cart_items.cart_item_id = $1;

-- name: GetCartItemByProductID :one
SELECT * FROM cart_items WHERE product_id = $1;

-- name: GetCartItems :many
SELECT cart_items.*, p.name AS product_name, p.price AS product_price, p.stock AS product_stock, img.image_url AS image_url
FROM cart_items
JOIN products AS p ON cart_items.product_id = p.product_id
LEFT JOIN images as img ON p.product_id = img.product_id AND img.primary = true
WHERE cart_id = $1;

-- name: ClearCart :exec
DELETE FROM cart_items WHERE cart_id = $1;