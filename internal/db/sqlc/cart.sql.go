// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: cart.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
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

const clearCart = `-- name: ClearCart :exec
DELETE FROM cart_items WHERE cart_id = $1
`

func (q *Queries) ClearCart(ctx context.Context, cartID int64) error {
	_, err := q.db.Exec(ctx, clearCart, cartID)
	return err
}

const createCart = `-- name: CreateCart :one
INSERT INTO 
    carts (user_id) 
VALUES 
    ($1) 
RETURNING id, checkout_at, user_id, updated_at, created_at
`

func (q *Queries) CreateCart(ctx context.Context, userID int64) (Cart, error) {
	row := q.db.QueryRow(ctx, createCart, userID)
	var i Cart
	err := row.Scan(
		&i.ID,
		&i.CheckoutAt,
		&i.UserID,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getCart = `-- name: GetCart :one
SELECT id, checkout_at, user_id, updated_at, created_at FROM carts WHERE checkout_at IS NULL AND user_id = $1 LIMIT 1
`

func (q *Queries) GetCart(ctx context.Context, userID int64) (Cart, error) {
	row := q.db.QueryRow(ctx, getCart, userID)
	var i Cart
	err := row.Scan(
		&i.ID,
		&i.CheckoutAt,
		&i.UserID,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getCartDetail = `-- name: GetCartDetail :many
SELECT carts.id, carts.checkout_at, carts.user_id, carts.updated_at, carts.created_at, cart_items.id, cart_items.product_id, cart_items.cart_id, cart_items.quantity, cart_items.created_at, products.id, products.name, products.description, products.sku, products.image_url, products.stock, products.archived, products.price, products.updated_at, products.created_at
FROM carts
JOIN cart_items ON carts.id = cart_items.cart_id
JOIN products ON cart_items.product_id = products.id
WHERE carts.user_id = $1 AND carts.checkout_at IS NULL
`

type GetCartDetailRow struct {
	Cart     Cart     `json:"cart"`
	CartItem CartItem `json:"cart_item"`
	Product  Product  `json:"product"`
}

func (q *Queries) GetCartDetail(ctx context.Context, userID int64) ([]GetCartDetailRow, error) {
	rows, err := q.db.Query(ctx, getCartDetail, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetCartDetailRow{}
	for rows.Next() {
		var i GetCartDetailRow
		if err := rows.Scan(
			&i.Cart.ID,
			&i.Cart.CheckoutAt,
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
DELETE FROM cart_items WHERE cart_id = $1 AND id = $2
`

type RemoveProductFromCartParams struct {
	CartID int64 `json:"cart_id"`
	ID     int64 `json:"id"`
}

func (q *Queries) RemoveProductFromCart(ctx context.Context, arg RemoveProductFromCartParams) error {
	_, err := q.db.Exec(ctx, removeProductFromCart, arg.CartID, arg.ID)
	return err
}

const setCartCheckoutAt = `-- name: SetCartCheckoutAt :exec
UPDATE carts SET checkout_at = $1, updated_at = NOW() WHERE id = $2
`

type SetCartCheckoutAtParams struct {
	CheckoutAt pgtype.Timestamptz `json:"checkout_at"`
	ID         int64              `json:"id"`
}

func (q *Queries) SetCartCheckoutAt(ctx context.Context, arg SetCartCheckoutAtParams) error {
	_, err := q.db.Exec(ctx, setCartCheckoutAt, arg.CheckoutAt, arg.ID)
	return err
}

const updateCart = `-- name: UpdateCart :exec
UPDATE carts SET updated_at = NOW() WHERE id = $1 RETURNING id, checkout_at, user_id, updated_at, created_at
`

func (q *Queries) UpdateCart(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, updateCart, id)
	return err
}
