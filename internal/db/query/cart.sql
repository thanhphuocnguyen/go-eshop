-- name: CreateCart :one
INSERT INTO 
    carts (user_id) 
VALUES 
    ($1) 
RETURNING *;

-- name: GetCart :one
SELECT * FROM carts WHERE checked_out_at IS NULL AND user_id = $1 LIMIT 1;

-- name: GetCartDetail :many
SELECT sqlc.embed(carts), sqlc.embed(cart_items), sqlc.embed(products)
FROM carts
JOIN cart_items ON carts.id = cart_items.cart_id
JOIN products ON cart_items.product_id = products.id
WHERE carts.user_id = $1 AND carts.checked_out_at IS NULL;

-- name: AddProductToCart :one
INSERT INTO cart_items(cart_id, product_id, quantity) VALUES ($1, $2, $3) RETURNING *;

-- name: RemoveProductFromCart :exec
DELETE FROM cart_items WHERE cart_id = $1 AND id = $2;

-- name: ClearCart :exec
DELETE FROM cart_items WHERE cart_id = $1;

-- name: SetCartCheckoutAt :exec
UPDATE carts SET checked_out_at = $1, updated_at = NOW() WHERE id = $2;

-- name: UpdateCart :exec
UPDATE carts SET updated_at = NOW() WHERE id = $1 RETURNING *;
