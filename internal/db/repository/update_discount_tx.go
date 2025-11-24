package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

type UpdateDiscountTxArgs struct {
	Description       *string          `json:"description" binding:"omitempty"`
	DiscountType      NullDiscountType `json:"discount_type" binding:"omitempty,oneof=percentage fixed_amount"`
	DiscountValue     *float64         `json:"discount_value" binding:"omitempty,gt=0"`
	MinPurchaseAmount *float64         `json:"min_purchase_amount" binding:"omitempty,gt=0"`
	MaxDiscountAmount *float64         `json:"max_discount_amount" binding:"omitempty,gt=0"`
	UsageLimit        *int32           `json:"usage_limit" binding:"omitempty,gte=0"`
	IsActive          *bool            `json:"is_active" binding:"omitempty"`
	ValidFrom         *time.Time       `json:"starts_at" binding:"omitempty"`
	ValidUntil        *time.Time       `json:"expires_at" binding:"omitempty"`
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
		if arg.DiscountType.Valid {
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
		if arg.ValidFrom != nil {
			sqlParams.ValidFrom = utils.GetPgTypeTimestamp(*arg.ValidFrom)
		}
		if arg.ValidUntil != nil {
			sqlParams.ValidUntil = utils.GetPgTypeTimestamp(*arg.ValidUntil)
		}

		if arg.MaxDiscountAmount != nil {
			sqlParams.MaxDiscountAmount = utils.GetPgNumericFromFloat(*arg.MaxDiscountAmount)
		}
		_, err = q.UpdateDiscount(ctx, sqlParams)

		if err != nil {
			return err
		}

		return nil
	})
	return err
}
