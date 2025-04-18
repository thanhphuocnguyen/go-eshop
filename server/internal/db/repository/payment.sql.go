// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: payment.sql

package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createPaymentTransaction = `-- name: CreatePaymentTransaction :one
INSERT INTO
    payments (
        id,
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
RETURNING id, order_id, amount, payment_method, status, payment_gateway, refund_id, created_at, updated_at
`

type CreatePaymentTransactionParams struct {
	ID             string             `json:"id"`
	OrderID        uuid.UUID          `json:"order_id"`
	Amount         pgtype.Numeric     `json:"amount"`
	PaymentMethod  PaymentMethod      `json:"payment_method"`
	PaymentGateway NullPaymentGateway `json:"payment_gateway"`
}

func (q *Queries) CreatePaymentTransaction(ctx context.Context, arg CreatePaymentTransactionParams) (Payment, error) {
	row := q.db.QueryRow(ctx, createPaymentTransaction,
		arg.ID,
		arg.OrderID,
		arg.Amount,
		arg.PaymentMethod,
		arg.PaymentGateway,
	)
	var i Payment
	err := row.Scan(
		&i.ID,
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
    id = $1
`

func (q *Queries) DeletePaymentTransaction(ctx context.Context, id string) error {
	_, err := q.db.Exec(ctx, deletePaymentTransaction, id)
	return err
}

const getPaymentTransactionByID = `-- name: GetPaymentTransactionByID :one
SELECT
    id, order_id, amount, payment_method, status, payment_gateway, refund_id, created_at, updated_at
FROM
    payments
WHERE
    id = $1
LIMIT 1
`

func (q *Queries) GetPaymentTransactionByID(ctx context.Context, id string) (Payment, error) {
	row := q.db.QueryRow(ctx, getPaymentTransactionByID, id)
	var i Payment
	err := row.Scan(
		&i.ID,
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
    id, order_id, amount, payment_method, status, payment_gateway, refund_id, created_at, updated_at
FROM
    payments
WHERE
    order_id = $1
LIMIT 1
`

func (q *Queries) GetPaymentTransactionByOrderID(ctx context.Context, orderID uuid.UUID) (Payment, error) {
	row := q.db.QueryRow(ctx, getPaymentTransactionByOrderID, orderID)
	var i Payment
	err := row.Scan(
		&i.ID,
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
    id = $1
`

type UpdatePaymentTransactionParams struct {
	ID            string            `json:"id"`
	Amount        pgtype.Numeric    `json:"amount"`
	PaymentMethod NullPaymentMethod `json:"payment_method"`
	RefundID      pgtype.Text       `json:"refund_id"`
	Status        NullPaymentStatus `json:"status"`
}

func (q *Queries) UpdatePaymentTransaction(ctx context.Context, arg UpdatePaymentTransactionParams) error {
	_, err := q.db.Exec(ctx, updatePaymentTransaction,
		arg.ID,
		arg.Amount,
		arg.PaymentMethod,
		arg.RefundID,
		arg.Status,
	)
	return err
}
