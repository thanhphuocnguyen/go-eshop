-- name: CreatePayment :one
INSERT INTO payments (order_id, amount, payment_method_id, gateway, status, payment_intent_id, charge_id, net_amount) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: GetPaymentByID :one
SELECT * FROM payments WHERE id = $1 LIMIT 1;

-- name: GetPaymentByOrderID :one
SELECT * FROM payments WHERE order_id = $1 LIMIT 1;

-- name: GetPaymentByPaymentIntentID :one
SELECT * FROM payments WHERE payment_intent_id = $1 LIMIT 1;

-- name: UpdatePayment :exec
UPDATE payments
SET
    amount = COALESCE(sqlc.narg(amount), amount),
    refund_id = COALESCE(sqlc.narg(refund_id), refund_id),
    status = COALESCE(sqlc.narg(status), status),
    gateway = COALESCE(sqlc.narg(gateway), gateway),
    payment_intent_id = COALESCE(sqlc.narg(payment_intent_id), payment_intent_id),
    charge_id = COALESCE(sqlc.narg(charge_id), charge_id),
    error_code = COALESCE(sqlc.narg(error_code), error_code),
    error_message = COALESCE(sqlc.narg(error_message), error_message),
    updated_at = NOW()
WHERE id = $1;

-- name: DeletePayment :exec
DELETE FROM payments WHERE id = $1;

-- Payment Transactions --
-- name: CreatePaymentTransaction :one
INSERT INTO payment_transactions (payment_id,amount,status,gateway_transaction_id,gateway_response_code,gateway_response_message) VALUES ($1,$2,$3,$4,$5,$6) RETURNING *;

-- name: GetPaymentTransactionByID :one
SELECT * FROM payment_transactions WHERE id = $1 LIMIT 1;

-- name: GetPaymentTransactionByPaymentID :one
SELECT * FROM payment_transactions WHERE payment_id = $1 LIMIT 1;

-- name: UpdatePaymentTransaction :exec
UPDATE payment_transactions
SET
    amount = COALESCE(sqlc.narg(amount), amount),
    status = COALESCE(sqlc.narg(status), status),
    gateway_transaction_id = COALESCE(sqlc.narg(gateway_transaction_id), gateway_transaction_id),
    gateway_response_code = COALESCE(sqlc.narg(gateway_response_code), gateway_response_code),
    gateway_response_message = COALESCE(sqlc.narg(gateway_response_message), gateway_response_message)
WHERE id = $1;

-- name: DeletePaymentTransaction :exec
DELETE FROM payment_transactions WHERE id = $1;

-- Payment Methods --
-- name: GetPaymentMethodByID :one
SELECT * FROM payment_methods WHERE id = $1 LIMIT 1;

-- name: GetPaymentMethods :many
SELECT * FROM payment_methods ORDER BY name ASC;

-- name: ListPaymentMethods :many
SELECT * FROM payment_methods WHERE is_active = TRUE ORDER BY name ASC;