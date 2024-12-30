-- name: CreateCart :one
INSERT INTO 
    carts (user_id) 
VALUES 
    ($1) 
RETURNING *;

-- name: GetCart :one
SELECT * FROM carts WHERE user_id = $1 LIMIT 1;

-- name: RemoveProductFromCart :exec
DELETE FROM cart_items WHERE cart_id = $1 AND id = $2;

-- name: UpdateCart :exec
UPDATE carts SET updated_at = NOW() WHERE id = $1 RETURNING *;
