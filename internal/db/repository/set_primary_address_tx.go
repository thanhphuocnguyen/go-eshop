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

func (repo *pgRepo) SetPrimaryAddressTx(ctx context.Context, arg SetPrimaryAddressTxArgs) error {
	err := repo.execTx(ctx, func(q *Queries) error {
		var err error
		err = q.ResetPrimaryAddress(ctx, arg.UserID)
		if err != nil {
			log.Error().Err(err).Msg("Cannot reset primary address")
			return err
		}

		err = q.SetPrimaryAddress(ctx, SetPrimaryAddressParams{
			UserID:    arg.UserID,
			ID:        arg.NewPrimaryID,
			IsDefault: true,
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
