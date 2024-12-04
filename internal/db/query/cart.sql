-- name: CreateCart :one
INSERT INTO 
    carts (user_id) 
VALUES 
    ($1) 
RETURNING *;

-- name: GetCartByUserID :many
SELECT sqlc.embed(carts), sqlc.embed(cart_items), sqlc.embed(products)
FROM carts
JOIN cart_items ON carts.id = cart_items.cart_id
JOIN products ON cart_items.product_id = products.id
WHERE carts.user_id = $1;

-- name: AddProductToCart :one
INSERT INTO cart_items(cart_id, product_id, quantity) VALUES ($1, $2, $3) RETURNING *;

-- name: RemoveProductFromCart :exec
DELETE FROM cart_items WHERE cart_id = $1 AND product_id = $2;

-- name: UpdateProductQuantity :exec
UPDATE cart_items SET quantity = $1 WHERE cart_id = $2 AND product_id = $3;