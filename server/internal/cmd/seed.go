package cmd

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/thanhphuocnguyen/go-eshop/config"
	"github.com/thanhphuocnguyen/go-eshop/internal/auth"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/util"
)

type Product struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int32   `json:"stock"`
	Sku         string  `json:"sku"`
}

type Collection struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
	Published   bool   `json:"published"`
}

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"fullname"`
	Phone    string `json:"phone"`
}

type Address struct {
	Street   string `json:"street"`
	City     string `json:"city"`
	District string `json:"district"`
	Ward     string `json:"ward"`
	Phone    string `json:"phone"`
	UserID   int64  `json:"user_id"`
}

func ExecuteSeed(ctx context.Context) int {
	var rootCmd = &cobra.Command{
		Use:   "seed",
		Short: "Seed data to database",
		Long:  "Seed data to database from seed files",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadConfig(".")
			if err != nil {
				log.Fatal().Err(err).Msg("failed to load config")
				return err
			}
			pg, err := repository.GetPostgresInstance(ctx, cfg)
			if err != nil {
				log.Error().Err(err).Msg("failed to get postgres instance")
				return err
			}
			defer pg.Close()
			if len(args) == 0 {
				var waitGroup sync.WaitGroup
				defer waitGroup.Wait()
				waitGroup.Add(3)
				go func() {
					defer waitGroup.Done()
					seedProducts(ctx, pg)
				}()
				go func() {
					defer waitGroup.Done()
					seedUsers(ctx, pg)
					seedUserAddresses(ctx, pg)
				}()
				go func() {
					defer waitGroup.Done()
					seedCollections(ctx, pg)
				}()
			} else {
				switch args[0] {
				case "products":
					seedProducts(ctx, pg)
				case "users":
					seedUsers(ctx, pg)
				case "addresses":
					seedUserAddresses(ctx, pg)
				case "collections":
					seedCollections(ctx, pg)
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

func seedProducts(ctx context.Context, repo repository.Repository) {
	countProducts, err := repo.CountProducts(ctx, repository.CountProductsParams{})
	if err != nil {
		log.Error().Err(err).Msg("failed to count products")
		return
	}

	if countProducts > 0 {
		log.Info().Msg("products already seeded")
		return
	}

	log.Info().Msg("seeding products")
	productData, err := os.ReadFile("seeds/products.json")
	if err != nil {
		log.Error().Err(err).Msg("failed to read product data")
		return
	}

	var products []Product
	log.Info().Msg("parsing product data")
	err = json.Unmarshal(productData, &products)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse product data")
		return
	}

	log.Info().Msgf("product count: %d", len(products))
	params := make([]repository.SeedProductsParams, len(products))
	for i, product := range products {
		price := util.GetPgNumericFromFloat(product.Price)
		params[i] = repository.SeedProductsParams{
			Name:        product.Name,
			Description: product.Description,
			Sku:         product.Sku,
			Stock:       product.Stock,
			Price:       price,
		}
	}
	_, err = repo.SeedProducts(ctx, params)
	if err != nil {
		log.Error().Err(err).Msgf("failed to create product: %v", err)
	}
	log.Info().Msg("products created")
}

func seedCollections(ctx context.Context, repo repository.Repository) {
	countCollections, err := repo.CountCollections(ctx, pgtype.Int4{})
	if err != nil {
		log.Error().Err(err).Msg("failed to count collections")
		return
	}

	if countCollections > 0 {
		log.Info().Msg("collections already seeded")
		return
	}

	log.Info().Msg("seeding collections")
	collectionData, err := os.ReadFile("seeds/collections.json")
	if err != nil {
		log.Error().Err(err).Msg("failed to read collection data")
		return
	}

	var collections []Collection
	log.Info().Msg("parsing collection data")
	err = json.Unmarshal(collectionData, &collections)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse collection data")
		return
	}

	log.Info().Msgf("collection count: %d", len(collections))
	params := make([]repository.SeedCollectionsParams, len(collections))
	for i, collection := range collections {
		params[i] = repository.SeedCollectionsParams{
			Name: collection.Name,
			Description: pgtype.Text{
				String: collection.Description,
				Valid:  true,
			},
			SortOrder: int16(collection.SortOrder),
			Published: collection.Published,
		}
	}
	_, err = repo.SeedCollections(ctx, params)
	if err != nil {
		log.Error().Err(err).Msgf("failed to create collection: %v", err)
	}
	log.Info().Msg("collections created")
}

func seedUsers(ctx context.Context, pg repository.Repository) {
	countUsers, err := pg.CountUsers(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to count users")
		return
	}

	if countUsers > 0 {
		log.Info().Msg("users already seeded")
		return
	}

	log.Info().Str("from file: ", "users.json").Msg("parsing user data")
	userData, err := os.ReadFile("seeds/users.json")
	if err != nil {
		return
	}

	var users []User
	log.Info().Msg("parsing user from data")
	err = json.Unmarshal(userData, &users)
	if err != nil {
		return
	}

	log.Info().Int("with user count: ", len(users)).Msg("creating users")
	params := make([]repository.SeedUsersParams, len(users))
	for i, user := range users {
		hashed, _ := auth.HashPassword(user.Password)

		params[i] = repository.SeedUsersParams{
			Email:          user.Email,
			Username:       user.Username,
			Phone:          user.Phone,
			HashedPassword: hashed,
			Fullname:       user.FullName,
			Role:           repository.UserRoleUser,
		}
		if params[i].Username == "admin" {
			params[i].Role = repository.UserRoleAdmin
		}

	}
	_, err = pg.SeedUsers(ctx, params)
	if err != nil {
		log.Error().Err(err).Msgf("failed to create user: %v", err)
	}
	log.Info().Msg("users created")
}

func seedUserAddresses(ctx context.Context, pg repository.Repository) {
	countAddresses, err := pg.CountAddresses(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to count addresses")
		return
	}

	if countAddresses > 0 {
		log.Info().Msg("addresses already seeded")
		return
	}

	log.Info().Msg("parsing addresses json")
	addressData, err := os.ReadFile("seeds/addresses.json")
	if err != nil {
		return
	}

	var addresses []Address
	log.Info().Msg("parsing addresses from data")
	err = json.Unmarshal(addressData, &addresses)
	if err != nil {
		return
	}

	log.Info().Msg("creating addresses")
	log.Info().Msgf("addresses count: %d", len(addresses))
	params := make([]repository.SeedAddressesParams, len(addresses))
	for i, user := range addresses {
		params[i] = repository.SeedAddressesParams{
			UserID:   user.UserID,
			Phone:    user.Phone,
			Street:   user.Street,
			Ward:     util.GetPgTypeText(user.Ward),
			District: user.District,
			City:     user.City,
		}
	}
	_, err = pg.SeedAddresses(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("failed to seed addresses")
		return
	}
	log.Info().Msg("addresses created")
}
