package cmd

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/thanhphuocnguyen/go-eshop/config"
	"github.com/thanhphuocnguyen/go-eshop/elasticsearch"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
)

func ExecuteIndexes(ctx context.Context) int {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
		return 1
	}

	var rootCmd = &cobra.Command{
		Use:   "seed",
		Short: "Seed data to database",
		Long:  "Seed data to database from seed files",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := repository.GetPostgresInstance(ctx, cfg)
			if err != nil {
				log.Error().Err(err).Msg("failed to connect to postgres")
				return err
			}
			client, err := elasticsearch.NewClient(cfg.EsUrl)
			if err != nil {
				log.Error().Err(err).Msg("failed to create elasticsearch client")
				return err
			}

			if len(args) == 0 {
				log.Info().Msg("Starting comprehensive database seeding with error handling")
				// Execute all seeding operations with proper error handling
				var seedingError error
				// Seed base data first (sequential order matters)
				if seedingError = indexProducts(ctx, client, store); seedingError != nil {
					return seedingError
				}

				log.Info().Msg("All seeding operations completed successfully")
			} else {
				// Individual seeding commands with error handling
				switch args[0] {
				case "products":
					if err := indexProducts(ctx, client, store); err != nil {
						return err
					}
				default:
					log.Error().Msg("invalid seed command")
				}
			}
			return nil
		},
	}

	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("failed to execute seed command")
		return 1
	}

	return 0
}

func indexProducts(ctx context.Context, client *elasticsearch.ESStore, pg repository.Store) error {
	products, err := pg.GetAdminProductList(ctx, repository.GetAdminProductListParams{
		Limit:  1000,
		Offset: 0,
	})
	for len(products) > 0 {
		productIndexer := elasticsearch.NewProductIndexer(client)
		if _, err := productIndexer.BulkIndexProducts(products); err != nil {
			log.Error().Err(err).Msg("failed to bulk index products")
			return err
		}

		offset := len(products)
		products, err = pg.GetAdminProductList(ctx, repository.GetAdminProductListParams{
			Limit:  1000,
			Offset: int64(offset),
		})
	}

	if err != nil {
		log.Error().Err(err).Msg("failed to get products from postgres")
		return err
	}
	productIndexer := elasticsearch.NewProductIndexer(client)
	productIndexer.BulkDeleteProducts()
	log.Info().Msg("successfully indexed products")
	return nil
}
