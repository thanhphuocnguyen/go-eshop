package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type SetPrimaryAddressTxArgs struct {
	NewPrimaryID uuid.UUID `json:"new_primary_id"`
	UserID       uuid.UUID `json:"user_id"`
}

func (s *pgRepo) SetPrimaryAddressTx(ctx context.Context, arg SetPrimaryAddressTxArgs) error {
	err := s.execTx(ctx, func(q *Queries) error {
		var err error
		err = s.ResetPrimaryAddress(ctx, arg.UserID)
		if err != nil {
			log.Error().Err(err).Msg("Cannot reset primary address")
			return err
		}

		err = s.SetPrimaryAddress(ctx, SetPrimaryAddressParams{
			Default: true,
			UserID:  arg.UserID,
			ID:      arg.NewPrimaryID,
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
