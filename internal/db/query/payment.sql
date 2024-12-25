-- name: CreatePayment :one
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