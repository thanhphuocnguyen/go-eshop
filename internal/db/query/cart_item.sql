
-- name: UpdateCartItemQuantity :exec
UPDATE cart_items SET quantity = $1 WHERE id = $2 RETURNING *;

-- name: CountCartItem :one
SELECT COUNT(*) FROM cart_items WHERE cart_id = $1;

-- name: GetCartItem :one
SELECT * FROM cart_items WHERE id = $1;


-- name: GetCartItemByProductID :one
SELECT * FROM cart_items WHERE product_id = $1;

-- name: GetCartItems :many
SELECT cart_items.*, p.name AS product_name, p.price AS product_price, p.stock AS product_stock, img.image_url AS image_url
FROM cart_items
JOIN products AS p ON cart_items.product_id = p.id
LEFT JOIN images as img ON p.id = img.product_id AND img.is_primary = true
WHERE cart_id = $1;

-- name: ClearCart :exec
DELETE FROM cart_items WHERE cart_id = $1;