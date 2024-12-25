
-- name: UpdateCartItemQuantity :exec
UPDATE cart_items SET quantity = $1 WHERE id = $2 RETURNING *;

-- name: CountCartItem :one
SELECT COUNT(*) FROM cart_items WHERE cart_id = $1;

-- name: GetCartItem :one
SELECT * FROM cart_items WHERE id = $1;

-- name: ClearCart :exec
DELETE FROM cart_items WHERE cart_id = $1;