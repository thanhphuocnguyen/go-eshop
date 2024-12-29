// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user_address.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createAddress = `-- name: CreateAddress :one
INSERT INTO
    user_addresses (
        user_id,
        phone,
        address_1,
        address_2,
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
        $6,
        $7
    ) RETURNING id, user_id, phone, address_1, address_2, ward, district, city, is_primary, is_deleted, created_at, updated_at, deleted_at
`

type CreateAddressParams struct {
	UserID   int64       `json:"user_id"`
	Phone    string      `json:"phone"`
	Address1 string      `json:"address_1"`
	Address2 pgtype.Text `json:"address_2"`
	Ward     pgtype.Text `json:"ward"`
	District string      `json:"district"`
	City     string      `json:"city"`
}

func (q *Queries) CreateAddress(ctx context.Context, arg CreateAddressParams) (UserAddress, error) {
	row := q.db.QueryRow(ctx, createAddress,
		arg.UserID,
		arg.Phone,
		arg.Address1,
		arg.Address2,
		arg.Ward,
		arg.District,
		arg.City,
	)
	var i UserAddress
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Phone,
		&i.Address1,
		&i.Address2,
		&i.Ward,
		&i.District,
		&i.City,
		&i.IsPrimary,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const deleteAddress = `-- name: DeleteAddress :exec
UPDATE
    user_addresses
SET
    is_deleted = true
WHERE
    id = $1 AND user_id = $2
`

type DeleteAddressParams struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
}

func (q *Queries) DeleteAddress(ctx context.Context, arg DeleteAddressParams) error {
	_, err := q.db.Exec(ctx, deleteAddress, arg.ID, arg.UserID)
	return err
}

const getAddress = `-- name: GetAddress :one
SELECT
    id, user_id, phone, address_1, address_2, ward, district, city, is_primary, is_deleted, created_at, updated_at, deleted_at
FROM
    user_addresses
WHERE
    id = $1 AND user_id = $2 AND is_deleted = false
LIMIT 1
`

type GetAddressParams struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
}

func (q *Queries) GetAddress(ctx context.Context, arg GetAddressParams) (UserAddress, error) {
	row := q.db.QueryRow(ctx, getAddress, arg.ID, arg.UserID)
	var i UserAddress
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Phone,
		&i.Address1,
		&i.Address2,
		&i.Ward,
		&i.District,
		&i.City,
		&i.IsPrimary,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getAddresses = `-- name: GetAddresses :many
SELECT
    id, user_id, phone, address_1, address_2, ward, district, city, is_primary, is_deleted, created_at, updated_at, deleted_at
FROM
    user_addresses
WHERE
    user_id = $1 AND is_deleted = false
ORDER BY
    is_primary DESC, id ASC
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
			&i.ID,
			&i.UserID,
			&i.Phone,
			&i.Address1,
			&i.Address2,
			&i.Ward,
			&i.District,
			&i.City,
			&i.IsPrimary,
			&i.IsDeleted,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
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
    id, user_id, phone, address_1, address_2, ward, district, city, is_primary, is_deleted, created_at, updated_at, deleted_at
FROM
    user_addresses
WHERE
    user_id = $1 AND is_primary = true AND is_deleted = false
LIMIT 1
`

func (q *Queries) GetPrimaryAddress(ctx context.Context, userID int64) (UserAddress, error) {
	row := q.db.QueryRow(ctx, getPrimaryAddress, userID)
	var i UserAddress
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Phone,
		&i.Address1,
		&i.Address2,
		&i.Ward,
		&i.District,
		&i.City,
		&i.IsPrimary,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const resetPrimaryAddress = `-- name: ResetPrimaryAddress :exec
UPDATE
    user_addresses
SET
    is_primary = false
WHERE
    user_id = $1 AND is_primary = true
`

func (q *Queries) ResetPrimaryAddress(ctx context.Context, userID int64) error {
	_, err := q.db.Exec(ctx, resetPrimaryAddress, userID)
	return err
}

const setPrimaryAddress = `-- name: SetPrimaryAddress :exec
UPDATE
    user_addresses
SET
    is_primary = $1
WHERE
    user_id = $2 AND id = $3 AND is_deleted = false
`

type SetPrimaryAddressParams struct {
	IsPrimary bool  `json:"is_primary"`
	UserID    int64 `json:"user_id"`
	ID        int64 `json:"id"`
}

func (q *Queries) SetPrimaryAddress(ctx context.Context, arg SetPrimaryAddressParams) error {
	_, err := q.db.Exec(ctx, setPrimaryAddress, arg.IsPrimary, arg.UserID, arg.ID)
	return err
}

const updateAddress = `-- name: UpdateAddress :one
UPDATE
    user_addresses
SET
    phone = coalesce($1, phone),
    address_1 = coalesce($2, address_1),
    address_2 = coalesce($3, address_2),
    ward = coalesce($4, ward),
    district = coalesce($5, district),
    city = coalesce($6, city),
    is_primary = coalesce($7, is_primary)
WHERE
    id = $8 AND user_id = $9 AND is_deleted = false
RETURNING id, user_id, phone, address_1, address_2, ward, district, city, is_primary, is_deleted, created_at, updated_at, deleted_at
`

type UpdateAddressParams struct {
	Phone     pgtype.Text `json:"phone"`
	Address1  pgtype.Text `json:"address_1"`
	Address2  pgtype.Text `json:"address_2"`
	Ward      pgtype.Text `json:"ward"`
	District  pgtype.Text `json:"district"`
	City      pgtype.Text `json:"city"`
	IsPrimary pgtype.Bool `json:"is_primary"`
	ID        int64       `json:"id"`
	UserID    int64       `json:"user_id"`
}

func (q *Queries) UpdateAddress(ctx context.Context, arg UpdateAddressParams) (UserAddress, error) {
	row := q.db.QueryRow(ctx, updateAddress,
		arg.Phone,
		arg.Address1,
		arg.Address2,
		arg.Ward,
		arg.District,
		arg.City,
		arg.IsPrimary,
		arg.ID,
		arg.UserID,
	)
	var i UserAddress
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Phone,
		&i.Address1,
		&i.Address2,
		&i.Ward,
		&i.District,
		&i.City,
		&i.IsPrimary,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}
