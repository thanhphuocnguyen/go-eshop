// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user_address.sql

package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countAddresses = `-- name: CountAddresses :one
SELECT
    COUNT(*)
FROM
    user_addresses
`

func (q *Queries) CountAddresses(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countAddresses)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createAddress = `-- name: CreateAddress :one
INSERT INTO
    user_addresses (
        user_id,
        phone,
        street,
        ward,
        district,
        city
    )
VALUES
    (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6
    ) RETURNING user_address_id, user_id, phone, street, ward, district, city, "default", deleted, created_at, updated_at
`

type CreateAddressParams struct {
	UserID   int64       `json:"user_id"`
	Phone    string      `json:"phone"`
	Street   string      `json:"street"`
	Ward     pgtype.Text `json:"ward"`
	District string      `json:"district"`
	City     string      `json:"city"`
}

func (q *Queries) CreateAddress(ctx context.Context, arg CreateAddressParams) (UserAddress, error) {
	row := q.db.QueryRow(ctx, createAddress,
		arg.UserID,
		arg.Phone,
		arg.Street,
		arg.Ward,
		arg.District,
		arg.City,
	)
	var i UserAddress
	err := row.Scan(
		&i.UserAddressID,
		&i.UserID,
		&i.Phone,
		&i.Street,
		&i.Ward,
		&i.District,
		&i.City,
		&i.Default,
		&i.Deleted,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteAddress = `-- name: DeleteAddress :exec
UPDATE
    user_addresses
SET
    deleted = true,
    updated_at = now()
WHERE
    user_address_id = $1 AND user_id = $2
`

type DeleteAddressParams struct {
	UserAddressID int64 `json:"user_address_id"`
	UserID        int64 `json:"user_id"`
}

func (q *Queries) DeleteAddress(ctx context.Context, arg DeleteAddressParams) error {
	_, err := q.db.Exec(ctx, deleteAddress, arg.UserAddressID, arg.UserID)
	return err
}

const getAddress = `-- name: GetAddress :one
SELECT
    user_address_id, user_id, phone, street, ward, district, city, "default", deleted, created_at, updated_at
FROM
    user_addresses
WHERE
    user_address_id = $1 AND user_id = $2 AND is_deleted = false
LIMIT 1
`

type GetAddressParams struct {
	UserAddressID int64 `json:"user_address_id"`
	UserID        int64 `json:"user_id"`
}

func (q *Queries) GetAddress(ctx context.Context, arg GetAddressParams) (UserAddress, error) {
	row := q.db.QueryRow(ctx, getAddress, arg.UserAddressID, arg.UserID)
	var i UserAddress
	err := row.Scan(
		&i.UserAddressID,
		&i.UserID,
		&i.Phone,
		&i.Street,
		&i.Ward,
		&i.District,
		&i.City,
		&i.Default,
		&i.Deleted,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAddresses = `-- name: GetAddresses :many
SELECT
    user_address_id, user_id, phone, street, ward, district, city, "default", deleted, created_at, updated_at
FROM
    user_addresses
WHERE
    user_id = $1 AND deleted = false
ORDER BY
    "default" DESC, user_address_id ASC
`

func (q *Queries) GetAddresses(ctx context.Context, userID int64) ([]UserAddress, error) {
	rows, err := q.db.Query(ctx, getAddresses, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UserAddress
	for rows.Next() {
		var i UserAddress
		if err := rows.Scan(
			&i.UserAddressID,
			&i.UserID,
			&i.Phone,
			&i.Street,
			&i.Ward,
			&i.District,
			&i.City,
			&i.Default,
			&i.Deleted,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getPrimaryAddress = `-- name: GetPrimaryAddress :one
SELECT
    user_address_id, user_id, phone, street, ward, district, city, "default", deleted, created_at, updated_at
FROM
    user_addresses
WHERE
    user_id = $1 AND "default" = true AND deleted = false
LIMIT 1
`

func (q *Queries) GetPrimaryAddress(ctx context.Context, userID int64) (UserAddress, error) {
	row := q.db.QueryRow(ctx, getPrimaryAddress, userID)
	var i UserAddress
	err := row.Scan(
		&i.UserAddressID,
		&i.UserID,
		&i.Phone,
		&i.Street,
		&i.Ward,
		&i.District,
		&i.City,
		&i.Default,
		&i.Deleted,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const resetPrimaryAddress = `-- name: ResetPrimaryAddress :exec
UPDATE
    user_addresses
SET
    "default" = false
WHERE
    user_id = $1 AND "default" = true
`

func (q *Queries) ResetPrimaryAddress(ctx context.Context, userID int64) error {
	_, err := q.db.Exec(ctx, resetPrimaryAddress, userID)
	return err
}

const setPrimaryAddress = `-- name: SetPrimaryAddress :exec
UPDATE
    user_addresses
SET
    "default" = $1
WHERE
    user_id = $2 AND user_address_id = $3 AND deleted = false
`

type SetPrimaryAddressParams struct {
	Default       bool  `json:"default"`
	UserID        int64 `json:"user_id"`
	UserAddressID int64 `json:"user_address_id"`
}

func (q *Queries) SetPrimaryAddress(ctx context.Context, arg SetPrimaryAddressParams) error {
	_, err := q.db.Exec(ctx, setPrimaryAddress, arg.Default, arg.UserID, arg.UserAddressID)
	return err
}

const updateAddress = `-- name: UpdateAddress :one
UPDATE
    user_addresses
SET
    phone = coalesce($1, phone),
    street = coalesce($2, street),
    ward = coalesce($3, ward),
    district = coalesce($4, district),
    city = coalesce($5, city),
    "default" = coalesce($6, "default")
WHERE
    user_address_id = $7 AND user_id = $8 AND deleted = false
RETURNING user_address_id, user_id, phone, street, ward, district, city, "default", deleted, created_at, updated_at
`

type UpdateAddressParams struct {
	Phone         pgtype.Text `json:"phone"`
	Street        pgtype.Text `json:"street"`
	Ward          pgtype.Text `json:"ward"`
	District      pgtype.Text `json:"district"`
	City          pgtype.Text `json:"city"`
	Default       pgtype.Bool `json:"default"`
	UserAddressID int64       `json:"user_address_id"`
	UserID        int64       `json:"user_id"`
}

func (q *Queries) UpdateAddress(ctx context.Context, arg UpdateAddressParams) (UserAddress, error) {
	row := q.db.QueryRow(ctx, updateAddress,
		arg.Phone,
		arg.Street,
		arg.Ward,
		arg.District,
		arg.City,
		arg.Default,
		arg.UserAddressID,
		arg.UserID,
	)
	var i UserAddress
	err := row.Scan(
		&i.UserAddressID,
		&i.UserID,
		&i.Phone,
		&i.Street,
		&i.Ward,
		&i.District,
		&i.City,
		&i.Default,
		&i.Deleted,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}