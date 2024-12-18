// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: order.sql

package sqlc

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

const createOrder = `-- name: CreateOrder :one
INSERT INTO
    orders (
        user_id,
        payment_type,
        is_cod
    )
VALUES
    (
        $1,
        $2,
        $3
    )
RETURNING id, user_id, status, shipping_id, payment_type, payment_status, is_cod, confirmed_at, cancelled_at, delivered_at, updated_at, created_at
`

type CreateOrderParams struct {
	UserID      int64       `json:"user_id"`
	PaymentType PaymentType `json:"payment_type"`
	IsCod       bool        `json:"is_cod"`
}

func (q *Queries) CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error) {
	row := q.db.QueryRow(ctx, createOrder, arg.UserID, arg.PaymentType, arg.IsCod)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Status,
		&i.ShippingID,
		&i.PaymentType,
		&i.PaymentStatus,
		&i.IsCod,
		&i.ConfirmedAt,
		&i.CancelledAt,
		&i.DeliveredAt,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const createOrderItem = `-- name: CreateOrderItem :one
INSERT INTO
    order_items (
        product_id,
        order_id,
        quantity,
        price
    )
VALUES
    (
        $1,
        $2,
        $3,
        $4
    )
RETURNING id, product_id, order_id, quantity, price, created_at
`

type CreateOrderItemParams struct {
	ProductID int64          `json:"product_id"`
	OrderID   int64          `json:"order_id"`
	Quantity  int32          `json:"quantity"`
	Price     pgtype.Numeric `json:"price"`
}

func (q *Queries) CreateOrderItem(ctx context.Context, arg CreateOrderItemParams) (OrderItem, error) {
	row := q.db.QueryRow(ctx, createOrderItem,
		arg.ProductID,
		arg.OrderID,
		arg.Quantity,
		arg.Price,
	)
	var i OrderItem
	err := row.Scan(
		&i.ID,
		&i.ProductID,
		&i.OrderID,
		&i.Quantity,
		&i.Price,
		&i.CreatedAt,
	)
	return i, err
}

const deleteOrder = `-- name: DeleteOrder :exec
DELETE FROM
    orders
WHERE
    id = $1
`

func (q *Queries) DeleteOrder(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteOrder, id)
	return err
}

const getOrder = `-- name: GetOrder :one
SELECT
    id, user_id, status, shipping_id, payment_type, payment_status, is_cod, confirmed_at, cancelled_at, delivered_at, updated_at, created_at
FROM
    orders
WHERE
    id = $1
LIMIT 1
`

func (q *Queries) GetOrder(ctx context.Context, id int64) (Order, error) {
	row := q.db.QueryRow(ctx, getOrder, id)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Status,
		&i.ShippingID,
		&i.PaymentType,
		&i.PaymentStatus,
		&i.IsCod,
		&i.ConfirmedAt,
		&i.CancelledAt,
		&i.DeliveredAt,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getOrderDetails = `-- name: GetOrderDetails :many
SELECT
    orders.id, orders.user_id, orders.status, orders.shipping_id, orders.payment_type, orders.payment_status, orders.is_cod, orders.confirmed_at, orders.cancelled_at, orders.delivered_at, orders.updated_at, orders.created_at, order_items.id, order_items.product_id, order_items.order_id, order_items.quantity, order_items.price, order_items.created_at, products.id, products.name, products.description, products.sku, products.image_url, products.stock, products.archived, products.price, products.updated_at, products.created_at
FROM
    orders
JOIN
    order_items ON order_items.order_id = orders.id
JOIN
    products ON order_items.product_id = products.id
WHERE
    orders.id = $1
`

type GetOrderDetailsRow struct {
	Order     Order     `json:"order"`
	OrderItem OrderItem `json:"order_item"`
	Product   Product   `json:"product"`
}

func (q *Queries) GetOrderDetails(ctx context.Context, id int64) ([]GetOrderDetailsRow, error) {
	rows, err := q.db.Query(ctx, getOrderDetails, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetOrderDetailsRow{}
	for rows.Next() {
		var i GetOrderDetailsRow
		if err := rows.Scan(
			&i.Order.ID,
			&i.Order.UserID,
			&i.Order.Status,
			&i.Order.ShippingID,
			&i.Order.PaymentType,
			&i.Order.PaymentStatus,
			&i.Order.IsCod,
			&i.Order.ConfirmedAt,
			&i.Order.CancelledAt,
			&i.Order.DeliveredAt,
			&i.Order.UpdatedAt,
			&i.Order.CreatedAt,
			&i.OrderItem.ID,
			&i.OrderItem.ProductID,
			&i.OrderItem.OrderID,
			&i.OrderItem.Quantity,
			&i.OrderItem.Price,
			&i.OrderItem.CreatedAt,
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

const listOrderItems = `-- name: ListOrderItems :many
SELECT
    id, product_id, order_id, quantity, price, created_at
FROM
    order_items
WHERE
    order_id = $1
ORDER BY
    id
LIMIT $2
OFFSET $3
`

type ListOrderItemsParams struct {
	OrderID int64 `json:"order_id"`
	Limit   int32 `json:"limit"`
	Offset  int32 `json:"offset"`
}

func (q *Queries) ListOrderItems(ctx context.Context, arg ListOrderItemsParams) ([]OrderItem, error) {
	rows, err := q.db.Query(ctx, listOrderItems, arg.OrderID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []OrderItem{}
	for rows.Next() {
		var i OrderItem
		if err := rows.Scan(
			&i.ID,
			&i.ProductID,
			&i.OrderID,
			&i.Quantity,
			&i.Price,
			&i.CreatedAt,
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

const listOrders = `-- name: ListOrders :many
SELECT
    id, user_id, status, shipping_id, payment_type, payment_status, is_cod, confirmed_at, cancelled_at, delivered_at, updated_at, created_at
FROM
    orders
WHERE
    user_id = $1
ORDER BY
    id
LIMIT $2
OFFSET $3
`

type ListOrdersParams struct {
	UserID int64 `json:"user_id"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListOrders(ctx context.Context, arg ListOrdersParams) ([]Order, error) {
	rows, err := q.db.Query(ctx, listOrders, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Order{}
	for rows.Next() {
		var i Order
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Status,
			&i.ShippingID,
			&i.PaymentType,
			&i.PaymentStatus,
			&i.IsCod,
			&i.ConfirmedAt,
			&i.CancelledAt,
			&i.DeliveredAt,
			&i.UpdatedAt,
			&i.CreatedAt,
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

const updateOrder = `-- name: UpdateOrder :one
UPDATE
    orders
SET
    status = coalesce($1, status),
    shipping_id = coalesce($2, shipping_id),
    payment_status = coalesce($3, payment_status),
    confirmed_at = coalesce($4, confirmed_at),
    cancelled_at = coalesce($5, cancelled_at),
    delivered_at = coalesce($6, delivered_at),
    updated_at = $7
WHERE
    id = $8
RETURNING id, user_id, status, shipping_id, payment_type, payment_status, is_cod, confirmed_at, cancelled_at, delivered_at, updated_at, created_at
`

type UpdateOrderParams struct {
	Status        NullOrderStatus    `json:"status"`
	ShippingID    pgtype.Int8        `json:"shipping_id"`
	PaymentStatus NullPaymentStatus  `json:"payment_status"`
	ConfirmedAt   pgtype.Timestamptz `json:"confirmed_at"`
	CancelledAt   pgtype.Timestamptz `json:"cancelled_at"`
	DeliveredAt   pgtype.Timestamptz `json:"delivered_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
	ID            int64              `json:"id"`
}

func (q *Queries) UpdateOrder(ctx context.Context, arg UpdateOrderParams) (Order, error) {
	row := q.db.QueryRow(ctx, updateOrder,
		arg.Status,
		arg.ShippingID,
		arg.PaymentStatus,
		arg.ConfirmedAt,
		arg.CancelledAt,
		arg.DeliveredAt,
		arg.UpdatedAt,
		arg.ID,
	)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Status,
		&i.ShippingID,
		&i.PaymentType,
		&i.PaymentStatus,
		&i.IsCod,
		&i.ConfirmedAt,
		&i.CancelledAt,
		&i.DeliveredAt,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}
