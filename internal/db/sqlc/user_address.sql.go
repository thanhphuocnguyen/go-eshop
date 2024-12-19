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
    ) RETURNING id, user_id, phone, address_1, address_2, ward, district, city, is_primary
`

type CreateAddressParams struct {
	UserID   int64       `json:"user_id"`
	Phone    string      `json:"phone"`
	Address1 string      `json:"address_1"`
	Address2 pgtype.Text `json:"address_2"`
	Ward     pgtype.Text `json:"ward"`
	District pgtype.Text `json:"district"`
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
	)
	return i, err
}

const deleteAddress = `-- name: DeleteAddress :exec
DELETE FROM
    user_addresses
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
    id, user_id, phone, address_1, address_2, ward, district, city, is_primary
FROM
    user_addresses
WHERE
    id = $1 AND user_id = $2 AND is_primary = COALESCE($3, is_primary)
LIMIT 1
`

type GetAddressParams struct {
	ID        int64       `json:"id"`
	UserID    int64       `json:"user_id"`
	IsPrimary pgtype.Bool `json:"is_primary"`
}

func (q *Queries) GetAddress(ctx context.Context, arg GetAddressParams) (UserAddress, error) {
	row := q.db.QueryRow(ctx, getAddress, arg.ID, arg.UserID, arg.IsPrimary)
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
	)
	return i, err
}

const listAddresses = `-- name: ListAddresses :many
SELECT
    id, user_id, phone, address_1, address_2, ward, district, city, is_primary
FROM
    user_addresses
WHERE
    user_id = $1
ORDER BY
    id
LIMIT $2
OFFSET $3
`

type ListAddressesParams struct {
	UserID int64 `json:"user_id"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListAddresses(ctx context.Context, arg ListAddressesParams) ([]UserAddress, error) {
	rows, err := q.db.Query(ctx, listAddresses, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []UserAddress{}
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

const setPrimaryAddress = `-- name: SetPrimaryAddress :exec
UPDATE
    user_addresses
SET
    is_primary = $1,
    updated_at = now()
WHERE
    user_id = $2 AND id = $3
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
    is_primary = coalesce($7, is_primary),
    updated_at = now()
WHERE
    id = $8 AND user_id = $9
RETURNING id, user_id, phone, address_1, address_2, ward, district, city, is_primary
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
	)
	return i, err
}