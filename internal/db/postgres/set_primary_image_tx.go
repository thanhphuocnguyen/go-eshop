package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/sqlc"
)

type SetPrimaryImageTxParams struct {
	NewPrimaryID int32
	ProductID    int64
}

func (s *Postgres) SetPrimaryImageTx(ctx context.Context, arg SetPrimaryImageTxParams) error {
	err := s.execTx(ctx, func(q *sqlc.Queries) error {
		var err error
		err = s.UnsetPrimaryImage(ctx, pgtype.Int8{
			Int64: arg.ProductID,
			Valid: true,
		})
		if err != nil {
			log.Error().Err(err).Msg("Cannot reset primary image")
			return err
		}

		err = s.SetPrimaryImage(ctx, arg.NewPrimaryID)

		if err != nil {
			log.Error().Err(err).Msg("Cannot toggle new primary image")
			return err
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("execTx")
	}

	return err
}
