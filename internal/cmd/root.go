/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"

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
	"github.com/thanhphuocnguyen/go-eshop/pkg/gateways"
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
		// PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 	if !profile {
		// 		return nil
		// 	}

		// 	f, pErr := os.Create(("cpu.pprof"))
		// 	if pErr != nil {
		// 		return pErr
		// 	}
		// 	_ = pprof.StartCPUProfile(f)

		// 	return nil
		// },
		// PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		// 	if !profile {
		// 		return nil
		// 	}
		// 	pprof.StopCPUProfile()

		// 	f, pErr := os.Create(("mem.pprof"))
		// 	if pErr != nil {
		// 		return pErr
		// 	}

		// 	defer f.Close()

		// 	runtime.GC()
		// 	err := pprof.WriteHeapProfile(f)
		// 	return err
		// },
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
		Short: "Run Eshop API server",
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
			if taskDistributor == nil {
				return fmt.Errorf("failed to create task distributor")
			}

			uploadService := upload.NewCloudinaryUploader(cfg)
			if uploadService == nil {
				return fmt.Errorf("failed to create upload service")
			}

			mailer := mailer.NewEmailSender(cfg.SmtpUsername, cfg.SmtpPassword, cfg.Env)
			if mailer == nil {
				return fmt.Errorf("failed to create mailer service")
			}

			service := payment.NewPaymentService()
			if service == nil {
				return fmt.Errorf("failed to create payment service")
			}

			// Register gateways
			paymentGateways, err := pgRepo.GetPaymentMethods(ctx)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to get payment methods from db")
			}
			for _, pg := range paymentGateways {
				switch pg.Code {
				case "stripe":
					service.RegisterGateway(pg.Code, gateways.NewStripeGateway)
				case "paypal":
					service.RegisterGateway(pg.Code, gateways.NewPaypalGateway)
				}
			}

			stripeConfig := payment.GatewayConfig{
				Name:          "stripe",
				APIKey:        cfg.StripePublishableKey,
				SecretKey:     cfg.StripeSecretKey,
				WebhookSecret: cfg.StripeWebhookSecret,
				Environment:   "sandbox",
			}

			if err := service.AddGateway(stripeConfig); err != nil {
				log.Fatal().Err(err).Msg("failed to add stripe gateway")
			}

			taskProcessor := worker.NewRedisTaskProcessor(redisCfg, pgRepo, mailer, cfg)
			if taskProcessor == nil {
				return fmt.Errorf("failed to create task processor")
			}

			api, err := api.NewAPI(cfg, pgRepo, taskDistributor, uploadService, service)
			if err != nil {
				return fmt.Errorf("failed to create API server: %w", err)
			}
			if api == nil {
				return fmt.Errorf("API server is nil")
			}

			server := api.Server(fmt.Sprintf(":%s", cfg.Port))
			if server == nil {
				return fmt.Errorf("HTTP server is nil")
			}

			// Validate all server dependencies before starting
			if err := api.ValidateServerDependencies(ctx); err != nil {
				return fmt.Errorf("server dependency validation failed: %w", err)
			}
			log.Info().Msg("Server dependency validation completed successfully")

			go func() {
				log.Info().Str("addr", fmt.Sprintf("http://%s:%s", cfg.Domain, cfg.Port)).Msg("API server started")
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Error().Err(err).Msg("HTTP server error")
				}
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
