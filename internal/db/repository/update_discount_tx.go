package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

type UpdateDiscountTxArgs struct {
	Description       *string    `json:"description" binding:"omitempty"`
	DiscountType      *string    `json:"discount_type" binding:"omitempty,oneof=percentage fixed_amount"`
	DiscountValue     *float64   `json:"discount_value" binding:"omitempty,gt=0"`
	MinPurchaseAmount *float64   `json:"min_purchase_amount" binding:"omitempty,gt=0"`
	MaxDiscountAmount *float64   `json:"max_discount_amount" binding:"omitempty,gt=0"`
	UsageLimit        *int32     `json:"usage_limit" binding:"omitempty,gte=0"`
	IsActive          *bool      `json:"is_active" binding:"omitempty"`
	StartsAt          *time.Time `json:"starts_at" binding:"omitempty"`
	ExpiresAt         *time.Time `json:"expires_at" binding:"omitempty"`
	// Related entities
	Products   []string `json:"products,omitempty" binding:"omitempty,uuidslice"`
	Categories []string `json:"categories,omitempty" binding:"omitempty,uuidslice"`
	Users      []string `json:"users,omitempty" binding:"omitempty,uuidslice"`
}

func (s *pgRepo) UpdateDiscountTx(ctx context.Context, id uuid.UUID, arg UpdateDiscountTxArgs) error {
	err := s.execTx(ctx, func(q *Queries) (err error) {
		sqlParams := UpdateDiscountParams{
			ID: id,
		}

		if arg.Description != nil {
			sqlParams.Description = arg.Description
		}
		if arg.DiscountType != nil {
			sqlParams.DiscountType = arg.DiscountType
		}
		if arg.DiscountValue != nil {
			sqlParams.DiscountValue = utils.GetPgNumericFromFloat(*arg.DiscountValue)
		}
		if arg.UsageLimit != nil {
			sqlParams.UsageLimit = arg.UsageLimit
		}
		if arg.IsActive != nil {
			sqlParams.IsActive = arg.IsActive
		}
		if arg.StartsAt != nil {
			sqlParams.StartsAt = utils.GetPgTypeTimestamp(*arg.StartsAt)
		}
		if arg.ExpiresAt != nil {
			sqlParams.ExpiresAt = utils.GetPgTypeTimestamp(*arg.ExpiresAt)
		}
		if arg.MinPurchaseAmount != nil {
			sqlParams.MinPurchaseAmount = utils.GetPgNumericFromFloat(*arg.MinPurchaseAmount)
		}
		if arg.MaxDiscountAmount != nil {
			sqlParams.MaxDiscountAmount = utils.GetPgNumericFromFloat(*arg.MaxDiscountAmount)
		}
		discountID, err := q.UpdateDiscount(ctx, sqlParams)

		if len(arg.Products) > 0 {
			if err := q.DeleteProductDiscountsByDiscountID(ctx, discountID); err != nil {
				return err
			}

			sqlParams := make([]InsertBulkProductDiscountsParams, len(arg.Products))
			for i, id := range arg.Products {
				sqlParams[i] = InsertBulkProductDiscountsParams{
					DiscountID: discountID,
					ProductID:  uuid.MustParse(id),
				}
			}
			if _, err = q.InsertBulkProductDiscounts(ctx, sqlParams); err != nil {
				return err
			}
		}
		if len(arg.Categories) > 0 {
			err := q.DeleteCategoryDiscountsByDiscountID(ctx, discountID)
			sqlParams := make([]InsertBulkCategoryDiscountsParams, len(arg.Categories))
			for i, id := range arg.Categories {
				sqlParams[i] = InsertBulkCategoryDiscountsParams{
					DiscountID: discountID,
					CategoryID: uuid.MustParse(id),
				}
			}
			if _, err = q.InsertBulkCategoryDiscounts(ctx, sqlParams); err != nil {
				return err
			}
		}

		if len(arg.Users) > 0 {
			if err := q.DeleteUserDiscountsByDiscountID(ctx, discountID); err != nil {
				return err
			}
			sqlParams := make([]InsertBulkUserDiscountsParams, len(arg.Users))
			for i, id := range arg.Users {
				sqlParams[i] = InsertBulkUserDiscountsParams{
					DiscountID: discountID,
					UserID:     uuid.MustParse(id),
				}
			}
			if _, err = q.InsertBulkUserDiscounts(ctx, sqlParams); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
