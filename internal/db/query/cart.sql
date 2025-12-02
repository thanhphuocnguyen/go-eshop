-- name: CreateCart :one
INSERT INTO carts (user_id, session_id) VALUES ($1, $2) RETURNING *;

-- name: GetCart :one
SELECT carts.*, COUNT(cart_items.id) AS item_count FROM carts
LEFT JOIN cart_items ON cart_items.cart_id = carts.id
WHERE ((carts.user_id IS NOT NULL AND carts.user_id = $1) OR (carts.session_id IS NOT NULL AND carts.session_id = $2)) AND carts.order_id IS NULL
GROUP BY carts.id
ORDER BY carts.updated_at DESC
LIMIT 1;

-- name: GetCartDetails :one
SELECT sqlc.embed(carts), COUNT(cart_items.id) as item_count, SUM(cart_items.quantity * product_variants.price) AS total_price, SUM(cart_items.quantity) AS total_quantity FROM carts
LEFT JOIN cart_items ON cart_items.cart_id = carts.id
JOIN product_variants ON product_variants.id = cart_items.variant_id
WHERE carts.id = $1
GROUP BY carts.id;

-- name: RemoveProductFromCart :exec
DELETE FROM cart_items WHERE cart_id = $1 AND id = $2;

-- name: UpdateCartTimestamp :exec
UPDATE carts SET updated_at = NOW() WHERE id = $1 RETURNING *;

-- name: CheckoutCart :exec
UPDATE carts SET order_id = $1 WHERE id = $2 RETURNING *;

-- Cart Item Section
-- name: AddCartItem :one
INSERT INTO cart_items (id, cart_id, variant_id, quantity) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateCartItemQuantity :exec
UPDATE cart_items SET quantity = $1 WHERE id = $2 RETURNING *;

-- name: CountCartItems :one
SELECT COUNT(*) FROM cart_items WHERE id = $1;

-- name: GetCartItem :one
SELECT * FROM cart_items WHERE id = $1 AND cart_id = $2;

-- name: GetCartItemByProductVariantID :one
SELECT * FROM cart_items WHERE variant_id = $1 AND cart_id = $2;

-- name: GetCartItems :many
SELECT
    sqlc.embed(ci),
    pv.price AS variant_price, pv.sku AS variant_sku, pv.stock AS variant_stock, pv.image_url AS variant_image_url,
    p.name AS product_name, p.id AS product_id, p.discount_percentage AS product_discount_percentage, p.brand_id AS product_brand_id,
    JSONB_AGG(
    DISTINCT JSONB_BUILD_OBJECT(
            'id', av.id,
            'name', a.name,
            'value', av.value
        )
    ) AS attributes,
    ARRAY_AGG(DISTINCT c.id)::uuid[] AS category_ids,
    ARRAY_AGG(DISTINCT col.id)::uuid[] AS collection_ids
FROM cart_items AS ci
JOIN product_variants AS pv ON pv.id = ci.variant_id
JOIN products AS p ON p.id = pv.product_id
JOIN variant_attribute_values AS vav ON vav.variant_id = pv.id
JOIN attribute_values AS av ON vav.attribute_value_id = av.id
JOIN attributes AS a ON av.attribute_id = a.id
LEFT JOIN category_products AS pc ON pc.product_id = p.id
LEFT JOIN categories AS c ON c.id = pc.category_id
LEFT JOIN collection_products AS pcol ON pcol.product_id = p.id
LEFT JOIN collections AS col ON col.id = pcol.collection_id
WHERE ci.cart_id = $1
GROUP BY ci.id, pv.id, p.id
ORDER BY ci.added_at, ci.id, pv.id DESC;

-- name: ClearCart :exec
DELETE FROM cart_items WHERE id = $1;
