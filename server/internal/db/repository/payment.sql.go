// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: payment.sql

package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createPaymentTransaction = `-- name: CreatePaymentTransaction :one
INSERT INTO
    payments (
        payment_id,
        order_id,
        amount,
        payment_method,
        payment_gateway
    )
VALUES
    (
        $1,
        $2,
        $3,
        $4,
        $5
    )
RETURNING payment_id, order_id, amount, payment_method, status, payment_gateway, refund_id, created_at, updated_at
`

type CreatePaymentTransactionParams struct {
	PaymentID      string             `json:"payment_id"`
	OrderID        int64              `json:"order_id"`
	Amount         pgtype.Numeric     `json:"amount"`
	PaymentMethod  PaymentMethod      `json:"payment_method"`
	PaymentGateway NullPaymentGateway `json:"payment_gateway"`
}

func (q *Queries) CreatePaymentTransaction(ctx context.Context, arg CreatePaymentTransactionParams) (Payment, error) {
	row := q.db.QueryRow(ctx, createPaymentTransaction,
		arg.PaymentID,
		arg.OrderID,
		arg.Amount,
		arg.PaymentMethod,
		arg.PaymentGateway,
	)
	var i Payment
	err := row.Scan(
		&i.PaymentID,
		&i.OrderID,
		&i.Amount,
		&i.PaymentMethod,
		&i.Status,
		&i.PaymentGateway,
		&i.RefundID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deletePaymentTransaction = `-- name: DeletePaymentTransaction :exec
DELETE FROM
    payments
WHERE
    payment_id = $1
`

func (q *Queries) DeletePaymentTransaction(ctx context.Context, paymentID string) error {
	_, err := q.db.Exec(ctx, deletePaymentTransaction, paymentID)
	return err
}

const getPaymentTransactionByID = `-- name: GetPaymentTransactionByID :one
SELECT
    payment_id, order_id, amount, payment_method, status, payment_gateway, refund_id, created_at, updated_at
FROM
    payments
WHERE
    payment_id = $1
LIMIT 1
`

func (q *Queries) GetPaymentTransactionByID(ctx context.Context, paymentID string) (Payment, error) {
	row := q.db.QueryRow(ctx, getPaymentTransactionByID, paymentID)
	var i Payment
	err := row.Scan(
		&i.PaymentID,
		&i.OrderID,
		&i.Amount,
		&i.PaymentMethod,
		&i.Status,
		&i.PaymentGateway,
		&i.RefundID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPaymentTransactionByOrderID = `-- name: GetPaymentTransactionByOrderID :one
SELECT
    payment_id, order_id, amount, payment_method, status, payment_gateway, refund_id, created_at, updated_at
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
		&i.PaymentID,
		&i.OrderID,
		&i.Amount,
		&i.PaymentMethod,
		&i.Status,
		&i.PaymentGateway,
		&i.RefundID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updatePaymentTransaction = `-- name: UpdatePaymentTransaction :exec
UPDATE
    payments
SET
    amount = COALESCE($2, amount),
    payment_method = COALESCE($3, payment_method),
    refund_id = COALESCE($4, refund_id),
    status = COALESCE($5, status)
WHERE
    payment_id = $1
`

type UpdatePaymentTransactionParams struct {
	PaymentID     string            `json:"payment_id"`
	Amount        pgtype.Numeric    `json:"amount"`
	PaymentMethod NullPaymentMethod `json:"payment_method"`
	RefundID      pgtype.Text       `json:"refund_id"`
	Status        NullPaymentStatus `json:"status"`
}

func (q *Queries) UpdatePaymentTransaction(ctx context.Context, arg UpdatePaymentTransactionParams) error {
	_, err := q.db.Exec(ctx, updatePaymentTransaction,
		arg.PaymentID,
		arg.Amount,
		arg.PaymentMethod,
		arg.RefundID,
		arg.Status,
	)
	return err
}