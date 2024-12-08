// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: cart.sql

package sqlc

import (
	"context"
)

const addProductToCart = `-- name: AddProductToCart :one
INSERT INTO cart_items(cart_id, product_id, quantity) VALUES ($1, $2, $3) RETURNING id, product_id, cart_id, quantity, created_at
`

type AddProductToCartParams struct {
	CartID    int64 `json:"cart_id"`
	ProductID int64 `json:"product_id"`
	Quantity  int16 `json:"quantity"`
}

func (q *Queries) AddProductToCart(ctx context.Context, arg AddProductToCartParams) (CartItem, error) {
	row := q.db.QueryRow(ctx, addProductToCart, arg.CartID, arg.ProductID, arg.Quantity)
	var i CartItem
	err := row.Scan(
		&i.ID,
		&i.ProductID,
		&i.CartID,
		&i.Quantity,
		&i.CreatedAt,
	)
	return i, err
}

const createCart = `-- name: CreateCart :one
INSERT INTO 
    carts (user_id) 
VALUES 
    ($1) 
RETURNING id, checked_out_at, user_id, updated_at, created_at
`

func (q *Queries) CreateCart(ctx context.Context, userID int64) (Cart, error) {
	row := q.db.QueryRow(ctx, createCart, userID)
	var i Cart
	err := row.Scan(
		&i.ID,
		&i.CheckedOutAt,
		&i.UserID,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getCartByUserID = `-- name: GetCartByUserID :one
SELECT id, checked_out_at, user_id, updated_at, created_at FROM carts WHERE user_id = $1
`

func (q *Queries) GetCartByUserID(ctx context.Context, userID int64) (Cart, error) {
	row := q.db.QueryRow(ctx, getCartByUserID, userID)
	var i Cart
	err := row.Scan(
		&i.ID,
		&i.CheckedOutAt,
		&i.UserID,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getCartDetailByUserID = `-- name: GetCartDetailByUserID :many
SELECT carts.id, carts.checked_out_at, carts.user_id, carts.updated_at, carts.created_at, cart_items.id, cart_items.product_id, cart_items.cart_id, cart_items.quantity, cart_items.created_at, products.id, products.name, products.description, products.sku, products.image_url, products.stock, products.archived, products.price, products.updated_at, products.created_at
FROM carts
JOIN cart_items ON carts.id = cart_items.cart_id
JOIN products ON cart_items.product_id = products.id
WHERE carts.user_id = $1
`

type GetCartDetailByUserIDRow struct {
	Cart     Cart     `json:"cart"`
	CartItem CartItem `json:"cart_item"`
	Product  Product  `json:"product"`
}

func (q *Queries) GetCartDetailByUserID(ctx context.Context, userID int64) ([]GetCartDetailByUserIDRow, error) {
	rows, err := q.db.Query(ctx, getCartDetailByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetCartDetailByUserIDRow{}
	for rows.Next() {
		var i GetCartDetailByUserIDRow
		if err := rows.Scan(
			&i.Cart.ID,
			&i.Cart.CheckedOutAt,
			&i.Cart.UserID,
			&i.Cart.UpdatedAt,
			&i.Cart.CreatedAt,
			&i.CartItem.ID,
			&i.CartItem.ProductID,
			&i.CartItem.CartID,
			&i.CartItem.Quantity,
			&i.CartItem.CreatedAt,
			&i.Product.ID,
			&i.Product.Name,
			&i.Product.Description,
			&i.Product.Sku,
			&i.Product.ImageUrl,
			&i.Product.Stock,
			&i.Product.Archived,
			&i.Product.Price,
			&i.Product.UpdatedAt,
			&i.Product.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const removeProductFromCart = `-- name: RemoveProductFromCart :exec
DELETE FROM cart_items WHERE cart_id = $1 AND product_id = $2
`

type RemoveProductFromCartParams struct {
	CartID    int64 `json:"cart_id"`
	ProductID int64 `json:"product_id"`
}

func (q *Queries) RemoveProductFromCart(ctx context.Context, arg RemoveProductFromCartParams) error {
	_, err := q.db.Exec(ctx, removeProductFromCart, arg.CartID, arg.ProductID)
	return err
}

const updateProductQuantity = `-- name: UpdateProductQuantity :exec
UPDATE cart_items SET quantity = $1 WHERE cart_id = $2 AND product_id = $3
`

type UpdateProductQuantityParams struct {
	Quantity  int16 `json:"quantity"`
	CartID    int64 `json:"cart_id"`
	ProductID int64 `json:"product_id"`
}

func (q *Queries) UpdateProductQuantity(ctx context.Context, arg UpdateProductQuantityParams) error {
	_, err := q.db.Exec(ctx, updateProductQuantity, arg.Quantity, arg.CartID, arg.ProductID)
	return err
}
