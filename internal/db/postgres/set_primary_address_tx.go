package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/sqlc"
)

type SetPrimaryAddressTxParams struct {
	NewPrimaryID int64 `json:"new_primary_id"`
	UserID       int64 `json:"user_id"`
}

type SetPrimaryAddressTxResult struct {
	Success bool `json:"success"`
}

func (s *Postgres) SetPrimaryAddressTx(ctx context.Context, arg SetPrimaryAddressTxParams) (SetPrimaryAddressTxResult, error) {
	var result SetPrimaryAddressTxResult
	err := s.execTx(ctx, func(q *sqlc.Queries) error {
		var err error
		// get all user addresses
		oldPrimary, err := s.GetAddress(ctx, sqlc.GetAddressParams{
			UserID: arg.UserID,
			IsPrimary: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("ListUserAddresses")
			return err
		}
		// set primary address
		err = s.SetPrimaryAddress(ctx, sqlc.SetPrimaryAddressParams{
			IsPrimary: false,
			UserID:    arg.UserID,
			ID:        oldPrimary.ID,
		})

		if err != nil {
			log.Error().Err(err).Msg("Cannot toggle old primary address")
		}
		err = s.SetPrimaryAddress(ctx, sqlc.SetPrimaryAddressParams{
			IsPrimary: true,
			UserID:    arg.UserID,
			ID:        arg.NewPrimaryID,
		})

		if err != nil {
			log.Error().Err(err).Msg("Cannot toggle new primary address")
			return err
		}
		result.Success = true
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("execTx")
	}

	return result, err
}
