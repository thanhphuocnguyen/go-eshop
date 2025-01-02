-- name: CreatePaymentTransaction :one
INSERT INTO
    payments (
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
        $4
    )
RETURNING *;

-- name: GetPaymentTransactionByID :one
SELECT
    *
FROM
    payments
WHERE
    payment_id = $1
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
    payment_method = COALESCE(sqlc.narg(payment_method), payment_method)
WHERE
    payment_id = $1;

-- name: DeletePaymentTransaction :exec
DELETE FROM
    payments
WHERE
    payment_id = $1;