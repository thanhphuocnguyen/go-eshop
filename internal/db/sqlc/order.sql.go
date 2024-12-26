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
        user_address_id,
        total_price
    )
VALUES
    (
        $1,
        $2,
        $3
    )
RETURNING id, user_id, user_address_id, total_price, status, confirmed_at, delivered_at, cancelled_at, refunded_at, updated_at, created_at
`

type CreateOrderParams struct {
	UserID        int64          `json:"user_id"`
	UserAddressID int64          `json:"user_address_id"`
	TotalPrice    pgtype.Numeric `json:"total_price"`
}

func (q *Queries) CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error) {
	row := q.db.QueryRow(ctx, createOrder, arg.UserID, arg.UserAddressID, arg.TotalPrice)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.UserAddressID,
		&i.TotalPrice,
		&i.Status,
		&i.ConfirmedAt,
		&i.DeliveredAt,
		&i.CancelledAt,
		&i.RefundedAt,
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
    id, user_id, user_address_id, total_price, status, confirmed_at, delivered_at, cancelled_at, refunded_at, updated_at, created_at
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
		&i.UserAddressID,
		&i.TotalPrice,
		&i.Status,
		&i.ConfirmedAt,
		&i.DeliveredAt,
		&i.CancelledAt,
		&i.RefundedAt,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getOrderDetails = `-- name: GetOrderDetails :many
SELECT
    orders.id, orders.user_id, orders.user_address_id, orders.total_price, orders.status, orders.confirmed_at, orders.delivered_at, orders.cancelled_at, orders.refunded_at, orders.updated_at, orders.created_at, order_items.id, order_items.product_id, order_items.order_id, order_items.quantity, order_items.price, order_items.created_at, products.id, products.name, products.description, products.sku, products.stock, products.archived, products.price, products.updated_at, products.created_at, user_addresses.id, user_addresses.user_id, user_addresses.phone, user_addresses.address_1, user_addresses.address_2, user_addresses.ward, user_addresses.district, user_addresses.city, user_addresses.is_primary, user_addresses.is_deleted, user_addresses.created_at, user_addresses.updated_at, user_addresses.deleted_at
FROM
    orders
LEFT JOIN
    order_items ON order_items.order_id = orders.id
LEFT JOIN
    products ON order_items.product_id = products.id
LEFT JOIN
    user_addresses ON orders.user_address_id = user_addresses.id
WHERE
    orders.id = $1
`

type GetOrderDetailsRow struct {
	Order       Order       `json:"order"`
	OrderItem   OrderItem   `json:"order_item"`
	Product     Product     `json:"product"`
	UserAddress UserAddress `json:"user_address"`
}

func (q *Queries) GetOrderDetails(ctx context.Context, id int64) ([]GetOrderDetailsRow, error) {
	rows, err := q.db.Query(ctx, getOrderDetails, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetOrderDetailsRow
	for rows.Next() {
		var i GetOrderDetailsRow
		if err := rows.Scan(
			&i.Order.ID,
			&i.Order.UserID,
			&i.Order.UserAddressID,
			&i.Order.TotalPrice,
			&i.Order.Status,
			&i.Order.ConfirmedAt,
			&i.Order.DeliveredAt,
			&i.Order.CancelledAt,
			&i.Order.RefundedAt,
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
			&i.Product.Stock,
			&i.Product.Archived,
			&i.Product.Price,
			&i.Product.UpdatedAt,
			&i.Product.CreatedAt,
			&i.UserAddress.ID,
			&i.UserAddress.UserID,
			&i.UserAddress.Phone,
			&i.UserAddress.Address1,
			&i.UserAddress.Address2,
			&i.UserAddress.Ward,
			&i.UserAddress.District,
			&i.UserAddress.City,
			&i.UserAddress.IsPrimary,
			&i.UserAddress.IsDeleted,
			&i.UserAddress.CreatedAt,
			&i.UserAddress.UpdatedAt,
			&i.UserAddress.DeletedAt,
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
	var items []OrderItem
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
    orders.id, orders.user_id, orders.user_address_id, orders.total_price, orders.status, orders.confirmed_at, orders.delivered_at, orders.cancelled_at, orders.refunded_at, orders.updated_at, orders.created_at, count(*) as total_items, sum(order_items.price) as total_price
FROM
    orders
JOIN order_items ON order_items.order_id = orders.id
WHERE
    user_id = $1
GROUP BY
    orders.id
ORDER BY
    orders.id
LIMIT $2
OFFSET $3
`

type ListOrdersParams struct {
	UserID int64 `json:"user_id"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type ListOrdersRow struct {
	Order      Order `json:"order"`
	TotalItems int64 `json:"total_items"`
	TotalPrice int64 `json:"total_price"`
}

func (q *Queries) ListOrders(ctx context.Context, arg ListOrdersParams) ([]ListOrdersRow, error) {
	rows, err := q.db.Query(ctx, listOrders, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListOrdersRow
	for rows.Next() {
		var i ListOrdersRow
		if err := rows.Scan(
			&i.Order.ID,
			&i.Order.UserID,
			&i.Order.UserAddressID,
			&i.Order.TotalPrice,
			&i.Order.Status,
			&i.Order.ConfirmedAt,
			&i.Order.DeliveredAt,
			&i.Order.CancelledAt,
			&i.Order.RefundedAt,
			&i.Order.UpdatedAt,
			&i.Order.CreatedAt,
			&i.TotalItems,
			&i.TotalPrice,
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
    confirmed_at = coalesce($2, confirmed_at),
    cancelled_at = coalesce($3, cancelled_at),
    delivered_at = coalesce($4, delivered_at),
    user_address_id = coalesce($5, user_address_id),
    updated_at = $6
WHERE
    id = $7
RETURNING id, user_id, user_address_id, total_price, status, confirmed_at, delivered_at, cancelled_at, refunded_at, updated_at, created_at
`

type UpdateOrderParams struct {
	Status        NullOrderStatus    `json:"status"`
	ConfirmedAt   pgtype.Timestamptz `json:"confirmed_at"`
	CancelledAt   pgtype.Timestamptz `json:"cancelled_at"`
	DeliveredAt   pgtype.Timestamptz `json:"delivered_at"`
	UserAddressID pgtype.Int8        `json:"user_address_id"`
	UpdatedAt     time.Time          `json:"updated_at"`
	ID            int64              `json:"id"`
}

func (q *Queries) UpdateOrder(ctx context.Context, arg UpdateOrderParams) (Order, error) {
	row := q.db.QueryRow(ctx, updateOrder,
		arg.Status,
		arg.ConfirmedAt,
		arg.CancelledAt,
		arg.DeliveredAt,
		arg.UserAddressID,
		arg.UpdatedAt,
		arg.ID,
	)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.UserAddressID,
		&i.TotalPrice,
		&i.Status,
		&i.ConfirmedAt,
		&i.DeliveredAt,
		&i.CancelledAt,
		&i.RefundedAt,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}
