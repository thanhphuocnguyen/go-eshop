package cmd

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/thanhphuocnguyen/go-eshop/config"
	"github.com/thanhphuocnguyen/go-eshop/internal/api"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/postgres"
	"github.com/thanhphuocnguyen/go-eshop/internal/uploadsrv"
	"github.com/thanhphuocnguyen/go-eshop/internal/worker"
)

func APICmd(ctx context.Context, cfg config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Args:  cobra.ExactArgs(0),
		Short: "Run Gin Gonic API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			pg, err := postgres.GetPostgresInstance(ctx, cfg)
			if err != nil {
				return err
			}
			log.Info().Msg("Postgres instance created")
			redisCfg := asynq.RedisClientOpt{
				Addr: cfg.RedisUrl,
			}
			taskDistributor := worker.NewRedisTaskDistributor(redisCfg)
			uploadService := uploadsrv.NewCloudinaryUploadService(cfg)

			api, err := api.NewAPI(cfg, pg, taskDistributor, uploadService)
			if err != nil {
				return err
			}
			server := api.Server(cfg.HttpAddr)

			go func() {
				log.Info().Msgf("api server started at :%s", cfg.HttpAddr)
				_ = server.ListenAndServe()
			}()

			<-ctx.Done()

			_ = server.Shutdown(ctx)
			log.Info().Msg("API server shutdown")
			_ = taskDistributor.Shutdown()
			log.Info().Msg("Task distributor shutdown")
			pg.Close()
			log.Info().Msg("Postgres instance closed")
			return nil
		},
	}

	return cmd
}
