
-- name: UpdateCartItemQuantity :exec
UPDATE cart_items SET quantity = $1 WHERE id = $2 RETURNING *;

-- name: CountCartItem :one
SELECT COUNT(*) FROM cart_items WHERE cart_id = $1;

-- name: GetCartItem :one
SELECT * FROM cart_items WHERE id = $1;

-- name: GetCartItems :many
SELECT sqlc.embed(cart_items), sqlc.embed(products), sqlc.embed(images) FROM cart_items
JOIN products ON cart_items.product_id = products.id
JOIN images ON products.id = images.product_id AND images.is_primary = true
WHERE cart_id = $1;

-- name: ClearCart :exec
DELETE FROM cart_items WHERE cart_id = $1;