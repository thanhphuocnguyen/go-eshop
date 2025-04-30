package cmd

import (
	"context"
	"encoding/json"
	"math/rand"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/thanhphuocnguyen/go-eshop/config"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
)

var Attributes []string = []string{"Color",
	"Size",
	"Material",
	"Brand",
	"Warranty",
	"Storage",
	"Style",
}

type Variant struct {
	Price      float64           `json:"price"`
	Stock      int32             `json:"stock"`
	Sku        string            `json:"sku"`
	Discount   int16             `json:"discount,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}
type Product struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Variants    []Variant `json:"variants"`
	Price       float64   `json:"price"`
	Stock       int32     `json:"stock"`
	Sku         string    `json:"sku"`
}

type Categories struct {
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
					seedAttributes(ctx, pg)
					seedProducts(ctx, pg)
				}()
				go func() {
					defer waitGroup.Done()
					seedUsers(ctx, pg)
				}()
				go func() {
					defer waitGroup.Done()
					seedCategories(ctx, pg)
				}()
			} else {
				switch args[0] {
				case "products":
					seedProducts(ctx, pg)
				case "users":
					seedUsers(ctx, pg)
				case "Categories":
					seedCategories(ctx, pg)
				case "attributes":
					seedAttributes(ctx, pg)
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

	log.Info().Int("product count: %d", len(products))
	for _, product := range products {
		params := repository.CreateProductParams{
			ID:          uuid.New(),
			Name:        product.Name,
			Description: &product.Description,
			BasePrice:   utils.GetPgNumericFromFloat(product.Price),
			BaseSku:     product.Sku,
			Slug:        product.Name,
		}
		_, err := repo.CreateProduct(ctx, params)
		if err != nil {
			log.Error().Err(err).Msg("failed to create product")
			continue
		}

	}

	log.Info().Msg("products created")
}

func seedCategories(ctx context.Context, repo repository.Repository) {
	countCategories, err := repo.CountCategories(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to count Categories")
		return
	}

	if countCategories > 0 {
		log.Info().Msg("Categories already seeded")
		return
	}

	log.Info().Msg("seeding Categories")
	CategorieData, err := os.ReadFile("seeds/Categories.json")
	if err != nil {
		log.Error().Err(err).Msg("failed to read Categories data")
		return
	}

	var Categories []Categories
	log.Info().Msg("parsing Categories data")
	err = json.Unmarshal(CategorieData, &Categories)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse Categories data")
		return
	}

	log.Info().Msgf("Categories count: %d", len(Categories))
	params := make([]repository.SeedCategoriesParams, len(Categories))
	for i, Categories := range Categories {
		params[i] = repository.SeedCategoriesParams{
			Name:        Categories.Name,
			Description: &Categories.Description,
			Slug:        Categories.Name,
		}
	}
	_, err = repo.SeedCategories(ctx, params)
	if err != nil {
		log.Error().Err(err).Msgf("failed to create Categories: %v", err)
	}
	log.Info().Msg("Categories created")
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
	log.Info().Msg("parsing users from data")
	err = json.Unmarshal(userData, &users)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse user data")
		return
	}

	log.Info().Int("with user count: ", len(users)).Msg("creating users")
	params := make([]repository.SeedUsersParams, len(users))
	userIDs := make([]uuid.UUID, len(users))
	for i, user := range users {
		hashed, _ := auth.HashPassword(user.Password)
		usrID := uuid.New()
		userIDs[i] = usrID
		params[i] = repository.SeedUsersParams{
			ID:             usrID,
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
		log.Error().Err(err).Msg("failed to create users")
	}
	seedUserAddresses(ctx, pg, userIDs)
	log.Info().Msg("users created")
}

func seedAttributes(ctx context.Context, repo repository.Repository) {
	countAttributes, err := repo.CountAttributes(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to count attributes")
		return
	}

	if countAttributes > 0 {
		log.Info().Msg("attributes already seeded")
		return
	}

	for _, attribute := range Attributes {
		_, err := repo.CreateAttribute(ctx, attribute)
		if err != nil {
			log.Error().Err(err).Msg("failed to create attribute")
			return
		}
	}

	log.Info().Msg("attributes and values created")
}

func seedUserAddresses(ctx context.Context, repo repository.Repository, userIDs []uuid.UUID) {
	countAddresses, err := repo.CountAddresses(ctx)
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
		log.Error().Err(err).Msg("failed to parse address data")
		return
	}

	log.Info().Msg("creating addresses")
	log.Info().Int("addresses count: %d", len(addresses))
	params := make([]repository.SeedAddressesParams, len(addresses))
	for i, address := range addresses {
		idx := rand.Intn(len(userIDs))
		params[i] = repository.SeedAddressesParams{
			UserID:   userIDs[idx],
			Phone:    address.Phone,
			Street:   address.Street,
			Ward:     &address.Ward,
			District: address.District,
			City:     address.City,
		}
	}
	_, err = repo.SeedAddresses(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("failed to seed addresses")
		return
	}
	log.Info().Msg("addresses created")
}
