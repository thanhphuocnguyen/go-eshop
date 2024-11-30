-- name: CreateCart :one
INSERT INTO 
    carts (user_id) 
VALUES 
    ($1) 
RETURNING *;

-- name: GetCartByID :one
SELECT c.id, c.user_id, ci.quantity, p.id as product_id, p.name, p.description, p.price, p.image_url
FROM carts AS c
INNER JOIN cart_items AS ci ON carts.id = cart_items.cart_id
INNER JOIN products AS p ON cart_items.product_id = products.id
WHERE c.id = $1
LIMIT 1;

-- name: AddProductToCart :one
INSERT INTO cart_items(cart_id, product_id, quantity) VALUES ($1, $2, $3) RETURNING *;

-- name: RemoveProductFromCart :exec
DELETE FROM cart_items WHERE cart_id = $1 AND product_id = $2;