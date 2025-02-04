// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: copyfrom.go

package repository

import (
	"context"
)

// iteratorForAddBulkAttributes implements pgx.CopyFromSource.
type iteratorForAddBulkAttributes struct {
	rows                 []string
	skippedFirstNextCall bool
}

func (r *iteratorForAddBulkAttributes) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForAddBulkAttributes) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0],
	}, nil
}

func (r iteratorForAddBulkAttributes) Err() error {
	return nil
}

func (q *Queries) AddBulkAttributes(ctx context.Context, name []string) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"attributes"}, []string{"name"}, &iteratorForAddBulkAttributes{rows: name})
}

// iteratorForAddBulkProducts implements pgx.CopyFromSource.
type iteratorForAddBulkProducts struct {
	rows                 []AddBulkProductsParams
	skippedFirstNextCall bool
}

func (r *iteratorForAddBulkProducts) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForAddBulkProducts) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].Name,
		r.rows[0].Description,
	}, nil
}

func (r iteratorForAddBulkProducts) Err() error {
	return nil
}

func (q *Queries) AddBulkProducts(ctx context.Context, arg []AddBulkProductsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"products"}, []string{"name", "description"}, &iteratorForAddBulkProducts{rows: arg})
}

// iteratorForCreateBulkVariantAttributes implements pgx.CopyFromSource.
type iteratorForCreateBulkVariantAttributes struct {
	rows                 []CreateBulkVariantAttributesParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateBulkVariantAttributes) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateBulkVariantAttributes) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].VariantID,
		r.rows[0].AttributeID,
		r.rows[0].Value,
	}, nil
}

func (r iteratorForCreateBulkVariantAttributes) Err() error {
	return nil
}

func (q *Queries) CreateBulkVariantAttributes(ctx context.Context, arg []CreateBulkVariantAttributesParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"variant_attributes"}, []string{"variant_id", "attribute_id", "value"}, &iteratorForCreateBulkVariantAttributes{rows: arg})
}

// iteratorForSeedAddresses implements pgx.CopyFromSource.
type iteratorForSeedAddresses struct {
	rows                 []SeedAddressesParams
	skippedFirstNextCall bool
}

func (r *iteratorForSeedAddresses) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForSeedAddresses) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].UserID,
		r.rows[0].Phone,
		r.rows[0].Street,
		r.rows[0].Ward,
		r.rows[0].District,
		r.rows[0].City,
		r.rows[0].Default,
	}, nil
}

func (r iteratorForSeedAddresses) Err() error {
	return nil
}

func (q *Queries) SeedAddresses(ctx context.Context, arg []SeedAddressesParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"user_addresses"}, []string{"user_id", "phone", "street", "ward", "district", "city", "default"}, &iteratorForSeedAddresses{rows: arg})
}

// iteratorForSeedCategories implements pgx.CopyFromSource.
type iteratorForSeedCategories struct {
	rows                 []SeedCategoriesParams
	skippedFirstNextCall bool
}

func (r *iteratorForSeedCategories) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForSeedCategories) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].Name,
		r.rows[0].Description,
		r.rows[0].SortOrder,
		r.rows[0].Published,
	}, nil
}

func (r iteratorForSeedCategories) Err() error {
	return nil
}

func (q *Queries) SeedCategories(ctx context.Context, arg []SeedCategoriesParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"categories"}, []string{"name", "description", "sort_order", "published"}, &iteratorForSeedCategories{rows: arg})
}

// iteratorForSeedUsers implements pgx.CopyFromSource.
type iteratorForSeedUsers struct {
	rows                 []SeedUsersParams
	skippedFirstNextCall bool
}

func (r *iteratorForSeedUsers) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForSeedUsers) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].UserID,
		r.rows[0].Email,
		r.rows[0].Username,
		r.rows[0].Phone,
		r.rows[0].Fullname,
		r.rows[0].HashedPassword,
		r.rows[0].Role,
	}, nil
}

func (r iteratorForSeedUsers) Err() error {
	return nil
}

func (q *Queries) SeedUsers(ctx context.Context, arg []SeedUsersParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"users"}, []string{"user_id", "email", "username", "phone", "fullname", "hashed_password", "role"}, &iteratorForSeedUsers{rows: arg})
}
