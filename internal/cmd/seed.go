package cmd

import (
	"context"
	"encoding/json"
	"os"
	"sync"

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
				seedProducts(ctx, pg)
				seedUsers(ctx, pg)
				seedUserAddresses(ctx, pg)
				return nil
			}
			switch args[0] {
			case "products":
				seedProducts(ctx, pg)
			case "users":
				seedUsers(ctx, pg)
			case "addresses":
				seedUserAddresses(ctx, pg)
			default:
				log.Error().Msg("invalid seed command")
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

	var wg sync.WaitGroup
	defer wg.Wait()

	for _, product := range products {
		wg.Add(1)
		go func(product Product) {
			defer wg.Done()
			price, err := util.ParsePgTypeNumber(product.Price)
			if err != nil {
				log.Error().Err(err).Msgf("failed to parse price for product: %s", product.Name)
				return
			}
			_, err = repo.CreateProduct(ctx, repository.CreateProductParams{
				Name:        product.Name,
				Description: product.Description,
				Sku:         product.Sku,
				Stock:       product.Stock,
				Price:       price,
			})
			if err != nil {
				log.Error().Err(err).Msgf("failed to create product: %s", product.Name)
			}
		}(product)
	}
	log.Info().Msg("products created")
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

	log.Info().Msg("parsing user data")
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

	log.Info().Msg("creating users")
	log.Info().Msgf("user count: %d", len(users))

	var wg sync.WaitGroup
	defer wg.Wait()

	for _, user := range users {
		wg.Add(1)
		go func() {
			defer wg.Done()
			hashed, err := auth.HashPassword(user.Password)
			if err != nil {
				return
			}
			params := repository.CreateUserParams{
				Email:          user.Email,
				Username:       user.Username,
				Phone:          user.Phone,
				HashedPassword: hashed,
				Fullname:       user.FullName,
				Role:           repository.UserRoleUser,
			}
			if user.Username == "admin" {
				params.Role = repository.UserRoleAdmin
			}
			_, err = pg.CreateUser(ctx, params)
			if err != nil {
				log.Error().Err(err).Msgf("failed to create user: %s", user.Email)
			}
		}()
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
	var wg sync.WaitGroup
	defer wg.Wait()
	for _, user := range addresses {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err = pg.CreateAddress(ctx, repository.CreateAddressParams{
				UserID:   user.UserID,
				Phone:    user.Phone,
				Street:   user.Street,
				Ward:     util.GetPgTypeText(user.Ward),
				District: user.District,
				City:     user.City,
			})
			if err != nil {
				log.Error().Err(err).Msgf("failed to create address: %s", user.Street)
			}
		}()
	}
	log.Info().Msg("addresses created")
}
