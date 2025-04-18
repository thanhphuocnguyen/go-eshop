-- name: CreatePaymentTransaction :one
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
RETURNING *;

-- name: GetPaymentTransactionByID :one
SELECT
    *
FROM
    payments
WHERE
    id = $1
LIMIT 1;

-- name: GetPaymentTransactionByOrderID :one
SELECT
    *
FROM
    payments
WHERE
    order_id = $1
LIMIT 1;

-- name: UpdatePaymentTransaction :exec
UPDATE
    payments
SET
    amount = COALESCE(sqlc.narg(amount), amount),
    payment_method = COALESCE(sqlc.narg(payment_method), payment_method),
    refund_id = COALESCE(sqlc.narg(refund_id), refund_id),
    status = COALESCE(sqlc.narg(status), status)
WHERE
    id = $1;

-- name: DeletePaymentTransaction :exec
DELETE FROM
    payments
WHERE
    id = $1;