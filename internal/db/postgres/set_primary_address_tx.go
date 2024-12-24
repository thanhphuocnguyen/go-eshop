package postgres

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/sqlc"
)

type SetPrimaryAddressTxParams struct {
	NewPrimaryID int64 `json:"new_primary_id"`
	UserID       int64 `json:"user_id"`
}

func (s *Postgres) SetPrimaryAddressTx(ctx context.Context, arg SetPrimaryAddressTxParams) error {
	err := s.execTx(ctx, func(q *sqlc.Queries) error {
		var err error
		err = s.ResetPrimaryAddress(ctx, arg.UserID)
		if err != nil {
			log.Error().Err(err).Msg("Cannot reset primary address")
			return err
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
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("execTx")
	}

	return err
}
