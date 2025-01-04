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
	}

	var upCmd = &cobra.Command{
		Use:   "up",
		Short: "Apply all up migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := migrate.New(cfg.MigrationPath, cfg.DbUrl)
			if err != nil {
				return fmt.Errorf("cannot create migration: %w", err)
			}
			log.Info().Msg("starting migration")
			err = m.Up()
			log.Info().Msg("migration finished")
			return err
		},
	}

	var downCmd = &cobra.Command{
		Use:   "down",
		Short: "Apply all down migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := migrate.New(cfg.MigrationPath, cfg.DbUrl)
			if err != nil {
				return fmt.Errorf("cannot create migration: %w", err)
			}
			log.Info().Msg("starting migration")
			err = m.Down()
			log.Info().Msg("migration finished")
			return err
		},
	}

	var redoCmd = &cobra.Command{
		Use:   "redo",
		Short: "Redo the last migration",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := migrate.New(cfg.MigrationPath, cfg.DbUrl)
			if err != nil {
				return fmt.Errorf("cannot create migration: %w", err)
			}
			log.Info().Msg("starting redo migration")
			err = m.Steps(-1)
			log.Info().Msg("migration finished")
			return err
		},
	}

	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Print the current migration status",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := migrate.New(cfg.MigrationPath, cfg.DbUrl)
			if err != nil {
				return fmt.Errorf("cannot create migration: %w", err)
			}
			status, dirty, sErr := m.Version()
			if sErr != nil {
				return sErr
			}
			fmt.Println(status, dirty)
			return nil
		},
	}

	migrateCmd.AddCommand(upCmd, downCmd, redoCmd, statusCmd)

	if err := migrateCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("failed to execute migrate command")
		return 1
	}
	return 0
}
