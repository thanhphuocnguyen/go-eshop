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

type Variant struct {
	Price    float64 `json:"price"`
	Stock    int32   `json:"stock"`
	Sku      string  `json:"sku"`
	Discount int32   `json:"discount,omitempty"`
}
type Product struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Variants    []Variant `json:"variants"`
	Price       float64   `json:"price"`
	Stock       int32     `json:"stock"`
	Sku         string    `json:"sku"`
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

type AttributeValue struct {
	Value string  `json:"value"`
	Color *string `json:"color,omitempty"`
}
type Attribute struct {
	Name   string           `json:"name"`
	Values []AttributeValue `json:"values"`
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
				waitGroup.Add(4)
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
				go func() {
					defer waitGroup.Done()
					seedAttributes(ctx, pg)
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
			Name:        product.Name,
			Description: product.Description,
		}
		createdProduct, err := repo.CreateProduct(ctx, params)
		if err != nil {
			log.Error().Err(err).Msg("failed to create product")
			continue
		}
		for _, variant := range product.Variants {

			_, err := repo.CreateVariant(ctx, repository.CreateVariantParams{
				ProductID:     createdProduct.ProductID,
				Price:         util.GetPgNumericFromFloat(variant.Price),
				StockQuantity: variant.Stock,
				Sku:           util.GetPgTypeText(variant.Sku),
				Discount:      variant.Discount,
			})
			if err != nil {
				log.Error().Err(err).Msg("failed to create variant")
				continue
			}
		}
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
	log.Info().Msg("parsing users from data")
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
		log.Error().Err(err).Msg("failed to create users")
	}
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

	log.Info().Msg("parsing attributes json")
	attributeData, err := os.ReadFile("seeds/attributes.json")
	if err != nil {
		log.Error().Err(err).Msg("failed to read attribute data")
		return
	}

	var attributes []Attribute
	log.Info().Msg("parsing attributes from data")
	err = json.Unmarshal(attributeData, &attributes)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse attribute data")
		return
	}

	log.Info().Msg("creating attributes")
	log.Info().Int("attributes count: %d", len(attributes))
	for _, attribute := range attributes {
		_, err := repo.CreateAttribute(ctx, attribute.Name)
		if err != nil {
			log.Error().Err(err).Msg("failed to create attribute")
			return
		}
	}

	log.Info().Msg("attributes and values created")
}

func seedUserAddresses(ctx context.Context, repo repository.Repository) {
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
		return
	}

	log.Info().Msg("creating addresses")
	log.Info().Int("addresses count: %d", len(addresses))
	params := make([]repository.SeedAddressesParams, len(addresses))
	for i, address := range addresses {
		params[i] = repository.SeedAddressesParams{
			UserID:   address.UserID,
			Phone:    address.Phone,
			Street:   address.Street,
			Ward:     util.GetPgTypeText(address.Ward),
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
