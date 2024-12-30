package cmd

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/thanhphuocnguyen/go-eshop/config"
)

func ExecuteMigrate(ctx context.Context) int {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}
	var migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			m, err := migrate.New(cfg.MigrationPath, cfg.DbUrl)
			if err != nil {
				return fmt.Errorf("cannot create migration: %w", err)
			}

			if len(args) == 0 {
				log.Info().Msg("starting migration")
				err = m.Up()
				log.Info().Msg("migration finished")
				return
			}
			switch args[0] {
			case "up":
				log.Info().Msg("starting migration")
				err = m.Up()
				log.Info().Msg("migration finished")
			case "down":
				log.Info().Msg("starting migration")
				err = m.Down()
				log.Info().Msg("migration finished")
			case "redo":
				log.Info().Msg("starting redo migration")
				err = m.Steps(-1)
				log.Info().Msg("migration finished")
			case "status":
				status, dirty, sErr := m.Version()
				if sErr != nil {
					err = sErr
					return
				}
				fmt.Println(status, dirty)
			}
			return
		},
	}
	if err := migrateCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("failed to execute migrate command")
		return 1
	}
	return 0
}
