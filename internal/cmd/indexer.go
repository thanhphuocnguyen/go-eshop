package cmd

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/thanhphuocnguyen/go-eshop/config"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/pkg/elasticsearch"
)

func ExecuteIndexer(ctx context.Context) int {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
		return 1
	}

	var rootCmd = &cobra.Command{
		Use:   "indexer",
		Short: "Index data to Elasticsearch from database",
		Long:  "Index data to Elasticsearch from database",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := repository.GetPostgresInstance(ctx, cfg)
			if err != nil {
				log.Error().Err(err).Msg("failed to connect to postgres")
				return err
			}
			client, err := elasticsearch.NewClient(cfg.EsUrl, cfg.EsUserName, cfg.EsPassword)
			if err != nil {
				log.Error().Err(err).Msg("failed to create elasticsearch client")
				return err
			}

			// if err := client.Ping(ctx); err != nil {
			// 	log.Error().Err(err).Msg("failed to ping elasticsearch")
			// 	return err
			// }

			productIndexer := elasticsearch.NewProductIndexer(client)

			if len(args) == 0 {
				log.Info().Msg("Starting comprehensive database indexing with error handling")
				// Execute all indexing operations with proper error handling
				var indexingError error
				// Index base data first (sequential order matters)
				if indexingError = indexProducts(ctx, productIndexer, store); indexingError != nil {
					return indexingError
				}

				log.Info().Msg("All indexing operations completed successfully")
			} else {
				// Individual indexing commands with error handling
				switch args[0] {
				case "index-products":
					if err := indexProducts(ctx, productIndexer, store); err != nil {
						return err
					}
				case "create-product-index":
					client.DeleteIndex(ctx, elasticsearch.PRODUCT_INDEX)
					log.Info().Msg("Existing product index deleted")
					if err := createProductIndex(ctx, productIndexer); err != nil {
						return err
					}
					log.Info().Msg("Product index created successfully")
				default:
					log.Error().Msg("invalid indexer command")
				}
			}
			client.Close(ctx)
			return nil
		},
	}

	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("failed to execute indexer command")
		return 1
	}

	return 0
}

func indexProducts(ctx context.Context, client *elasticsearch.ProductIndexer, pg repository.Store) error {
	offset := int64(0)
	limit := int64(1000)

	for {
		products, err := pg.GetProductsForIndexing(ctx, repository.GetProductsForIndexingParams{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to get products from postgres")
			return err
		}

		// Break if no more products
		if len(products) == 0 {
			break
		}

		if _, err := client.BulkIndexProducts(ctx, products); err != nil {
			log.Error().Err(err).Msg("failed to bulk index products")
			return err
		}

		log.Info().Int("count", len(products)).Int64("offset", offset).Msg("indexed products batch")
		offset += int64(len(products))
	}

	log.Info().Msg("successfully indexed all products")
	return nil
}

func createProductIndex(ctx context.Context, client *elasticsearch.ProductIndexer) error {
	return client.CreateIndex(ctx)
}
