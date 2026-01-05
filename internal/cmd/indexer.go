package cmd

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/thanhphuocnguyen/go-eshop/config"
	"github.com/thanhphuocnguyen/go-eshop/elasticsearch"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
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
			client, err := elasticsearch.NewClient(cfg.EsUrl)
			if err != nil {
				log.Error().Err(err).Msg("failed to create elasticsearch client")
				return err
			}

			if err := client.Ping(); err != nil {
				log.Error().Err(err).Msg("failed to ping elasticsearch")
				return err
			}

			if len(args) == 0 {
				log.Info().Msg("Starting comprehensive database indexing with error handling")
				// Execute all indexing operations with proper error handling
				var indexingError error
				// Index base data first (sequential order matters)
				if indexingError = indexProducts(ctx, client, store); indexingError != nil {
					return indexingError
				}

				log.Info().Msg("All indexing operations completed successfully")
			} else {
				// Individual indexing commands with error handling
				switch args[0] {
				case "products":
					if err := indexProducts(ctx, client, store); err != nil {
						return err
					}
				default:
					log.Error().Msg("invalid indexer command")
				}
			}
			return nil
		},
	}

	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("failed to execute indexer command")
		return 1
	}

	return 0
}

func indexProducts(ctx context.Context, client *elasticsearch.ESStore, pg repository.Store) error {
	// Create the products index with mapping first
	if err := createProductIndex(client); err != nil {
		log.Error().Err(err).Msg("failed to create product index")
		return err
	}

	offset := int64(0)
	limit := int64(1000)

	for {
		products, err := pg.GetAdminProductList(ctx, repository.GetAdminProductListParams{
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

		productIndexer := elasticsearch.NewProductIndexer(client)
		if _, err := productIndexer.BulkIndexProducts(ctx, products); err != nil {
			log.Error().Err(err).Msg("failed to bulk index products")
			return err
		}

		log.Info().Int("count", len(products)).Int64("offset", offset).Msg("indexed products batch")
		offset += int64(len(products))
	}

	log.Info().Msg("successfully indexed all products")
	return nil
}

func createProductIndex(client *elasticsearch.ESStore) error {
	// Check if index already exists
	if exists, err := client.IndexExists("products"); err != nil {
		return err
	} else if exists {
		log.Info().Msg("products index already exists, skipping creation")
		return nil
	}

	// Read the product mapping from file
	mappingFile := filepath.Join("elasticsearch", "mappings", "product_mapping.json")
	mappingData, err := os.ReadFile(mappingFile)
	if err != nil {
		return err
	}

	var mapping map[string]interface{}
	if err := json.Unmarshal(mappingData, &mapping); err != nil {
		return err
	}

	log.Info().Msg("creating products index")
	return client.CreateIndex("products", mapping)
}
