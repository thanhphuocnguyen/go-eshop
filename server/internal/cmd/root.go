/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
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

	rootCmd.PersistentFlags().StringVarP(&cfg.HttpAddr, "http-addr", "a", cfg.HttpAddr, "HTTP address")
	rootCmd.PersistentFlags().BoolVarP(&profile, "profile", "p", false, "enable profiling")

	rootCmd.AddCommand(apiCmd(ctx, cfg))
	rootCmd.AddCommand(workerCmd(ctx, cfg))

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
			repo, err := repository.GetPostgresInstance(ctx, cfg)
			if err != nil {
				return err
			}
			log.Info().Msg("Postgres instance created")
			redisCfg := asynq.RedisClientOpt{
				Addr: cfg.RedisUrl,
			}
			taskDistributor := worker.NewRedisTaskDistributor(redisCfg)
			uploadService := upload.NewCloudinaryUploadService(cfg)
			paymentCtx := &payment.PaymentContext{}

			api, err := api.NewAPI(cfg, repo, taskDistributor, uploadService, paymentCtx)
			if err != nil {
				return err
			}
			server := api.Server(cfg.HttpAddr)

			go func() {
				log.Info().Str("addr", cfg.HttpAddr).Msg("API server started")
				_ = server.ListenAndServe()
			}()

			<-ctx.Done()

			_ = server.Shutdown(ctx)
			log.Info().Msg("API server shutdown")
			_ = taskDistributor.Shutdown()
			log.Info().Msg("Task distributor shutdown")
			repo.Close()
			log.Info().Msg("Postgres instance closed")
			return nil
		},
	}

	return cmd
}

func workerCmd(ctx context.Context, cfg config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "worker",
		Args:  cobra.ExactArgs(0),
		Short: "Run the asynq worker to process tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			pgDB, err := repository.GetPostgresInstance(ctx, cfg)
			if err != nil {
				log.Error().Err(err).Msg("failed to connect to postgres")
				return err
			}
			mailer := mailer.NewEmailSender(cfg.SmtpUsername, cfg.SmtpPassword, cfg.Env)
			redisCfg := asynq.RedisClientOpt{
				Addr: cfg.RedisUrl,
			}

			distributor := worker.NewRedisTaskProcessor(redisCfg, pgDB, mailer, cfg)

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
