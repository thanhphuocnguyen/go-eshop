/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/thanhphuocnguyen/go-eshop/config"
	"github.com/thanhphuocnguyen/go-eshop/internal/api"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/worker"
	"github.com/thanhphuocnguyen/go-eshop/pkg/cache"
	"github.com/thanhphuocnguyen/go-eshop/pkg/mailer"
	"github.com/thanhphuocnguyen/go-eshop/pkg/payment"
	"github.com/thanhphuocnguyen/go-eshop/pkg/upload"
)

// rootCmd represents the base command when called without any subcommands

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ctx context.Context) int {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	if cfg.Env == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	profile := false
	var rootCmd = &cobra.Command{
		Use:   "web",
		Short: "web is a web server",
		Long:  `web is a web server`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if !profile {
				return nil
			}

			f, pErr := os.Create(("cpu.pprof"))
			if pErr != nil {
				return pErr
			}
			_ = pprof.StartCPUProfile(f)

			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if !profile {
				return nil
			}
			pprof.StopCPUProfile()

			f, pErr := os.Create(("mem.pprof"))
			if pErr != nil {
				return pErr
			}

			defer f.Close()

			runtime.GC()
			err := pprof.WriteHeapProfile(f)
			return err
		},
	}

	rootCmd.PersistentFlags().StringVarP(&cfg.Domain, "domain", "d", cfg.Domain, "HTTP domain")
	rootCmd.PersistentFlags().StringVarP(&cfg.Port, "port", "p", cfg.Port, "HTTP port")
	rootCmd.PersistentFlags().BoolVarP(&profile, "profile", "", false, "enable profiling")

	rootCmd.AddCommand(apiCmd(ctx, cfg))

	if err = rootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("failed to execute command")
		return 1
	}
	return 0
}

func apiCmd(ctx context.Context, cfg config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Args:  cobra.ExactArgs(0),
		Short: "Run Gin Gonic API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			pgRepo, err := repository.GetPostgresInstance(ctx, cfg)
			if err != nil {
				return err
			}

			log.Info().Msg("Postgres instance created")

			redisCfg := asynq.RedisClientOpt{
				Addr: cfg.RedisUrl,
			}

			taskDistributor := worker.NewRedisTaskDistributor(redisCfg)
			uploadService := upload.NewCloudinaryUploadService(cfg)
			mailer := mailer.NewEmailSender(cfg.SmtpUsername, cfg.SmtpPassword, cfg.Env)
			paymentCtx := &payment.PaymentContext{}
			taskProcessor := worker.NewRedisTaskProcessor(redisCfg, pgRepo, mailer, cfg)
			cacheService := cache.NewRedisCache(cfg)
			api, err := api.NewAPI(cfg, pgRepo, cacheService, taskDistributor, uploadService, paymentCtx)
			if err != nil {
				return err
			}

			server := api.Server(fmt.Sprintf(":%s", cfg.Port))
			go func() {
				log.Info().Str("addr", fmt.Sprintf("http://%s:%s", cfg.Domain, cfg.Port)).Msg("API server started")
				_ = server.ListenAndServe()
			}()

			go func() {
				log.Info().Msg("Starting task distributor")
				err = taskProcessor.Start()
				if err != nil {
					log.Error().Err(err).Msg("task processor stopped")
				}
			}()

			<-ctx.Done()
			log.Info().Msg("Shutting down API server")
			_ = server.Shutdown(ctx)

			log.Info().Msg("Shutting down task distributor")
			_ = taskDistributor.Shutdown()

			log.Info().Msg("shutting down task processor")
			taskProcessor.Shutdown()

			log.Info().Msg("Shutting down pgRepo")
			pgRepo.Close()

			return nil
		},
	}

	return cmd
}
