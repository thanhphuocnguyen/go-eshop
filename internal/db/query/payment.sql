-- name: CreatePaymentTransaction :one
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
    method = COALESCE(sqlc.narg(method), method)
WHERE
    id = $1;

-- name: DeletePaymentTransaction :exec
DELETE FROM
    payments
WHERE
    id = $1;