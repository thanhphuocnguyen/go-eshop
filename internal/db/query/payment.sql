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