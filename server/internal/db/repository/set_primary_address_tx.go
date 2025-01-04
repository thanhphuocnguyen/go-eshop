package repository

import (
	"context"

	"github.com/rs/zerolog/log"
)

type SetPrimaryAddressTxParams struct {
	NewPrimaryID int64 `json:"new_primary_id"`
	UserID       int64 `json:"user_id"`
}

func (s *pgRepo) SetPrimaryAddressTx(ctx context.Context, arg SetPrimaryAddressTxParams) error {
	err := s.execTx(ctx, func(q *Queries) error {
		var err error
		err = s.ResetPrimaryAddress(ctx, arg.UserID)
		if err != nil {
			log.Error().Err(err).Msg("Cannot reset primary address")
			return err
		}

		err = s.SetPrimaryAddress(ctx, SetPrimaryAddressParams{
			Default:       true,
			UserID:        arg.UserID,
			UserAddressID: arg.NewPrimaryID,
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