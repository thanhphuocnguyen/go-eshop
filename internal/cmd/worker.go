package cmd

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/thanhphuocnguyen/go-eshop/config"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/postgres"
	"github.com/thanhphuocnguyen/go-eshop/internal/worker"
)

func WorkerCmd(ctx context.Context, cfg config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "worker",
		Args:  cobra.ExactArgs(0),
		Short: "Run the asynq worker to process tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			pgDB, err := postgres.GetPostgresInstance(ctx, cfg)
			if err != nil {
				log.Error().Err(err).Msg("failed to connect to postgres")
				return err
			}
			redisCfg := asynq.RedisClientOpt{
				Addr: cfg.RedisUrl,
			}
			distributor := worker.NewRedisTaskProcessor(
				redisCfg,
				pgDB,
			)

			err = distributor.Start()
			if err != nil {
				log.Error().Err(err).Msg("failed to start task processor")
				return err
			}

			log.Info().Msg("task processor started")
			<-ctx.Done()
			log.Info().Msg("shutting down task processor")
			distributor.Shutdown()

			return nil
		},
	}
	return cmd
}
