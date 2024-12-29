// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: payment.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createPaymentTransaction = `-- name: CreatePaymentTransaction :one
INSERT INTO
    payments (
        order_id,
        amount,
        method
    )
VALUES
    (
        $1,
        $2,
        $3
    )
RETURNING id, order_id, amount, method, status, gateway, transaction_id, created_at, updated_at
`

type CreatePaymentTransactionParams struct {
	OrderID int64          `json:"order_id"`
	Amount  pgtype.Numeric `json:"amount"`
	Method  PaymentMethod  `json:"method"`
}

func (q *Queries) CreatePaymentTransaction(ctx context.Context, arg CreatePaymentTransactionParams) (Payment, error) {
	row := q.db.QueryRow(ctx, createPaymentTransaction, arg.OrderID, arg.Amount, arg.Method)
	var i Payment
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.Amount,
		&i.Method,
		&i.Status,
		&i.Gateway,
		&i.TransactionID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deletePaymentTransaction = `-- name: DeletePaymentTransaction :exec
DELETE FROM
    payments
WHERE
    id = $1
`

func (q *Queries) DeletePaymentTransaction(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deletePaymentTransaction, id)
	return err
}

const getPaymentTransactionByID = `-- name: GetPaymentTransactionByID :one
SELECT
    id, order_id, amount, method, status, gateway, transaction_id, created_at, updated_at
FROM
    payments
WHERE
    id = $1
LIMIT 1
`

func (q *Queries) GetPaymentTransactionByID(ctx context.Context, id int32) (Payment, error) {
	row := q.db.QueryRow(ctx, getPaymentTransactionByID, id)
	var i Payment
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.Amount,
		&i.Method,
		&i.Status,
		&i.Gateway,
		&i.TransactionID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPaymentTransactionByOrderID = `-- name: GetPaymentTransactionByOrderID :one
SELECT
    id, order_id, amount, method, status, gateway, transaction_id, created_at, updated_at
FROM
    payments
WHERE
    order_id = $1
LIMIT 1
`

func (q *Queries) GetPaymentTransactionByOrderID(ctx context.Context, orderID int64) (Payment, error) {
	row := q.db.QueryRow(ctx, getPaymentTransactionByOrderID, orderID)
	var i Payment
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.Amount,
		&i.Method,
		&i.Status,
		&i.Gateway,
		&i.TransactionID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updatePaymentTransaction = `-- name: UpdatePaymentTransaction :exec
UPDATE
    payments
SET
    order_id = COALESCE($2, order_id),
    amount = COALESCE($3, amount),
    method = COALESCE($4, method)
WHERE
    id = $1
`

type UpdatePaymentTransactionParams struct {
	ID      int32             `json:"id"`
	OrderID pgtype.Int8       `json:"order_id"`
	Amount  pgtype.Numeric    `json:"amount"`
	Method  NullPaymentMethod `json:"method"`
}

func (q *Queries) UpdatePaymentTransaction(ctx context.Context, arg UpdatePaymentTransactionParams) error {
	_, err := q.db.Exec(ctx, updatePaymentTransaction,
		arg.ID,
		arg.OrderID,
		arg.Amount,
		arg.Method,
	)
	return err
}
