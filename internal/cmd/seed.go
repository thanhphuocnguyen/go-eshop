package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/thanhphuocnguyen/go-eshop/config"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
)

var AttributesSeed = []string{"Color",
	"Size",
	"Material",
	"Warranty",
	"Storage",
	"Style",
	"Processor",
}

type VariantSeedModel struct {
	Price      float64           `json:"price"`
	Stock      int32             `json:"stock"`
	Sku        string            `json:"sku"`
	Discount   int16             `json:"discount,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}
type Product struct {
	Name               string             `json:"name"`
	Description        string             `json:"description"`
	ShortDescription   *string            `json:"short_description,omitempty"`
	DiscountPercentage *int16             `json:"discount_percentage,omitempty"`
	Variants           []VariantSeedModel `json:"variants"`
	Price              float64            `json:"base_price"`
	Stock              int32              `json:"stock"`
	Sku                string             `json:"base_sku"`
	Brand              string             `json:"brand,omitempty"`
	Attributes         map[string]string  `json:"attributes,omitempty"`
	Categories         []string           `json:"categories,omitempty"`
	Collections        []string           `json:"collections,omitempty"`
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

type Brand struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
	Published   bool   `json:"published"`
}

type Collection struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
	Published   bool   `json:"published"`
	Slug        string `json:"slug"`
}

type AttributeValues struct {
	Attribute string   `json:"attribute"`
	Values    []string `json:"values"`
}

type Discount struct {
	Code           string   `json:"code"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	DiscountType   string   `json:"discount_type"`
	DiscountValue  float64  `json:"discount_value"`
	MinOrderValue  *float64 `json:"min_order_value"`
	MaxDiscountAmt *float64 `json:"max_discount_amount"`
	UsageLimit     *int32   `json:"usage_limit"`
	UsagePerUser   *int32   `json:"usage_per_user"`
	IsActive       bool     `json:"is_active"`
	IsStackable    bool     `json:"is_stackable"`
	Priority       int32    `json:"priority"`
	ValidFrom      string   `json:"valid_from"`
	ValidUntil     *string  `json:"valid_until"`
}

type ShippingMethod struct {
	Name                  string `json:"name"`
	Description           string `json:"description"`
	IsActive              bool   `json:"is_active"`
	RequiresAddress       bool   `json:"requires_address"`
	EstimatedDeliveryTime string `json:"estimated_delivery_time"`
}

type ShippingZone struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Countries   []string `json:"countries"`
	States      []string `json:"states"`
	ZipCodes    []string `json:"zip_codes"`
	IsActive    bool     `json:"is_active"`
}

func ExecuteSeed(ctx context.Context) int {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
		return 1
	}
	pg, err := repository.GetPostgresInstance(ctx, cfg)

	if err != nil {
		log.Error().Err(err).Msg("failed to get postgres instance")
		return 1
	}

	var rootCmd = &cobra.Command{
		Use:   "seed",
		Short: "Seed data to database",
		Long:  "Seed data to database from seed files",
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) == 0 {
				log.Info().Msg("Starting comprehensive database seeding with error handling")

				// Execute all seeding operations with proper error handling
				var seedingError error
				// Seed base data first (sequential order matters)
				if seedingError = seedAttributes(ctx, pg); seedingError != nil {
					return seedingError
				}
				if seedingError = seedAttributeValues(ctx, pg); seedingError != nil {
					return seedingError
				}
				if seedingError = seedCategories(ctx, pg); seedingError != nil {
					return seedingError
				}
				if seedingError = seedBrands(ctx, pg); seedingError != nil {
					return seedingError
				}
				if seedingError = seedCollections(ctx, pg); seedingError != nil {
					return seedingError
				}
				if seedingError = seedShippingMethods(ctx, pg); seedingError != nil {
					return seedingError
				}
				if seedingError = seedShippingZones(ctx, pg); seedingError != nil {
					return seedingError
				}
				if seedingError = seedDiscounts(ctx, pg); seedingError != nil {
					return seedingError
				}
				if seedingError = seedUsers(ctx, pg); seedingError != nil {
					return seedingError
				}
				if seedingError = seedUserAddresses(ctx, pg); seedingError != nil {
					return seedingError
				}
				if seedingError = seedProducts(ctx, pg); seedingError != nil {
					return seedingError
				}

				log.Info().Msg("All seeding operations completed successfully")

			} else {
				// Individual seeding commands with error handling
				switch args[0] {
				case "products":
					if err := seedProducts(ctx, pg); err != nil {
						return err
					}
				case "users":
					if err := seedUsers(ctx, pg); err != nil {
						return err
					}
				case "categories":
					if err := seedCategories(ctx, pg); err != nil {
						return err
					}
				case "attributes":
					if err := seedAttributes(ctx, pg); err != nil {
						return err
					}
				case "attribute-values":
					if err := seedAttributeValues(ctx, pg); err != nil {
						return err
					}
				case "brands":
					if err := seedBrands(ctx, pg); err != nil {
						return err
					}
				case "collections":
					if err := seedCollections(ctx, pg); err != nil {
						return err
					}
				case "discounts":
					if err := seedDiscounts(ctx, pg); err != nil {
						return err
					}
				case "shipping":
					if err := seedShippingMethods(ctx, pg); err != nil {
						return err
					}
					if err := seedShippingZones(ctx, pg); err != nil {
						return err
					}
				case "variants":
					if err := seedExistingProductVariants(ctx, pg); err != nil {
						return err
					}
				case "user-addresses":
					if err := seedUserAddresses(ctx, pg); err != nil {
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

func seedProducts(ctx context.Context, repo repository.Store) error {
	// remove all products first
	cnt, err := repo.CountProducts(ctx, repository.CountProductsParams{})
	if err != nil {
		log.Error().Err(err).Msg("failed to count products")
		return err
	}
	if cnt > 0 {
		_, err := repo.QueryRaw(ctx, "TRUNCATE TABLE products CASCADE;")
		if err != nil {
			log.Error().Err(err).Msg("failed to truncate products table")
			return err
		}
	}
	log.Info().Msg("seeding products")
	productData, err := os.ReadFile("seeds/products.json")
	if err != nil {
		log.Error().Err(err).Msg("failed to read product data")
		return err
	}

	var products []Product
	log.Info().Msg("parsing product data")
	err = json.Unmarshal(productData, &products)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse product data")
		return err
	}

	// Get all reference data
	brands, err := repo.GetBrands(ctx, repository.GetBrandsParams{
		Limit:  1000,
		Offset: 0,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to get brands")
	}

	categories, err := repo.GetCategories(ctx, repository.GetCategoriesParams{
		Limit:  1000,
		Offset: 0,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to get categories")
	}

	collections, err := repo.GetCollections(ctx, repository.GetCollectionsParams{
		Limit:  1000,
		Offset: 0,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to get collections")
	}

	log.Info().Int("product count", len(products)).Msg("creating products")
	for _, product := range products {

		params := repository.CreateProductParams{
			Name:               product.Name,
			Description:        product.Description,
			ShortDescription:   product.ShortDescription,
			BasePrice:          utils.GetPgNumericFromFloat(product.Price),
			BaseSku:            product.Sku,
			Slug:               utils.Slugify(product.Name),
			DiscountPercentage: product.DiscountPercentage,
		}

		// Find and set brand
		if product.Brand != "" {
			brand := slices.IndexFunc(brands, func(b repository.Brand) bool {
				return strings.EqualFold(b.Name, product.Brand)
			})
			if brand != -1 {
				params.BrandID = utils.GetPgTypeUUID(brands[brand].ID)
			} else {
				log.Warn().Str("brand", product.Brand).Str("product", product.Name).Msg("brand not found")
			}
		}

		createdProduct, err := repo.CreateProduct(ctx, params)
		if err != nil {
			log.Error().Err(err).Msg("failed to create product")
			continue
		}

		// Create category relationships
		if len(product.Categories) > 0 {
			err := createProductCategoryRelationships(ctx, repo, createdProduct.ID, product.Categories, categories)
			if err != nil {
				log.Error().Err(err).Str("product", product.Name).Msg("failed to create category relationships")
			}
		}

		// Create collection relationships
		if len(product.Collections) > 0 {
			err := createProductCollectionRelationships(ctx, repo, createdProduct.ID, product.Collections, collections)
			if err != nil {
				log.Error().Err(err).Str("product", product.Name).Msg("failed to create collection relationships")
			}
		}

		// Create attribute relationships
		if len(product.Attributes) > 0 {
			err := createProductAttributeRelationships(ctx, repo, createdProduct.ID, product.Attributes)
			if err != nil {
				log.Error().Err(err).Str("product", product.Name).Msg("failed to create attribute relationships")
			}
		}

		// Create variants for this product
		if len(product.Variants) > 0 {
			err := seedProductVariants(ctx, repo, createdProduct.ID, product.Variants)
			if err != nil {
				log.Error().Err(err).Str("product", product.Name).Msg("failed to create variants")
			}
		}
	}

	log.Info().Msg("products and variants created")
	return nil
}

func seedProductVariants(ctx context.Context, repo repository.Store, productID uuid.UUID, variants []VariantSeedModel) error {
	log.Info().Int("variant count", len(variants)).Str("product_id", productID.String()).Msg("creating product variants")

	for _, variant := range variants {
		// Create the variant
		variantParams := repository.CreateProductVariantParams{
			ProductID: productID,
			Sku:       variant.Sku,
			Price:     utils.GetPgNumericFromFloat(variant.Price),
			Stock:     variant.Stock,
		}

		// Add description if we have attributes
		if len(variant.Attributes) > 0 {
			descParts := []string{}
			for key, value := range variant.Attributes {
				descParts = append(descParts, fmt.Sprintf("%s: %s", key, value))
			}
			description := strings.Join(descParts, ", ")
			variantParams.Description = &description
		}

		createdVariant, err := repo.CreateProductVariant(ctx, variantParams)
		if err != nil {
			log.Error().Err(err).Str("sku", variant.Sku).Msg("failed to create variant")
			continue
		}

		// Create variant attributes
		if len(variant.Attributes) > 0 {
			err := createVariantAttributes(ctx, repo, createdVariant.ID, variant.Attributes)
			if err != nil {
				log.Error().Err(err).Str("variant_id", createdVariant.ID.String()).Msg("failed to create variant attributes")
			}
		}
	}

	log.Info().Str("product_id", productID.String()).Msg("variants created successfully")
	return nil
}

func createVariantAttributes(ctx context.Context, repo repository.Store, variantID uuid.UUID, attributes map[string]string) error {
	log.Debug().Str("variant_id", variantID.String()).Interface("attributes", attributes).Msg("creating variant attributes")

	var attributeValueParams []repository.CreateBulkProductVariantAttributeParams

	for attrName, attrValue := range attributes {
		// Get the attribute by name
		attr, err := repo.GetAttributeByName(ctx, attrName)
		if err != nil {
			log.Error().Err(err).Str("attribute", attrName).Msg("failed to get attribute")
			continue
		}

		// Get the attribute value
		attributeValues, err := repo.GetAttributeValues(ctx, attr.ID)
		if err != nil {
			log.Error().Err(err).Str("attribute", attrName).Msg("failed to get attribute values")
			continue
		}

		// Find the matching attribute value
		var matchingValueID *int64
		for _, av := range attributeValues {
			if av.Value == attrValue {
				matchingValueID = &av.ID
				break
			}
		}

		if matchingValueID == nil {
			log.Warn().Str("attribute", attrName).Str("value", attrValue).Msg("attribute value not found, creating it")
			// Create the missing attribute value
			createdValue, err := repo.CreateAttributeValue(ctx, repository.CreateAttributeValueParams{
				AttributeID: attr.ID,
				Value:       attrValue,
			})
			if err != nil {
				log.Error().Err(err).Str("attribute", attrName).Str("value", attrValue).Msg("failed to create attribute value")
				continue
			}
			matchingValueID = &createdValue.ID
		}

		// Add to bulk parameters
		attributeValueParams = append(attributeValueParams, repository.CreateBulkProductVariantAttributeParams{
			VariantID:        variantID,
			AttributeValueID: *matchingValueID,
		})
	}

	// Bulk create variant attributes
	if len(attributeValueParams) > 0 {
		_, err := repo.CreateBulkProductVariantAttribute(ctx, attributeValueParams)
		if err != nil {
			return fmt.Errorf("failed to create variant attributes: %w", err)
		}
		log.Debug().Str("variant_id", variantID.String()).Int("count", len(attributeValueParams)).Msg("variant attributes created")
	}

	return nil
}

func seedCategories(ctx context.Context, repo repository.Store) error {
	_, err := repo.QueryRaw(ctx, "TRUNCATE TABLE categories CASCADE;")
	if err != nil {
		log.Error().Err(err).Msg("failed to truncate categories table")
		return err
	}

	log.Info().Msg("seeding categories")
	categoryData, err := os.ReadFile("seeds/categories.json")
	if err != nil {
		log.Error().Err(err).Msg("failed to read category data")
		return err
	}

	var categories []Categories
	log.Info().Msg("parsing category data")
	err = json.Unmarshal(categoryData, &categories)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse category data")
		return err
	}

	log.Info().Msgf("category count: %d", len(categories))
	params := make([]repository.SeedCategoriesParams, len(categories))
	for i, category := range categories {
		params[i] = repository.SeedCategoriesParams{
			Name:        category.Name,
			Description: &category.Description,
			Slug:        utils.Slugify(category.Name),
		}
	}
	_, err = repo.SeedCategories(ctx, params)
	if err != nil {
		log.Error().Err(err).Msgf("failed to create categories: %v", err)
		return err
	}
	log.Info().Msg("categories created")
	return nil
}

func seedBrands(ctx context.Context, repo repository.Store) error {
	_, err := repo.QueryRaw(ctx, "TRUNCATE TABLE brands CASCADE;")
	if err != nil {
		log.Error().Err(err).Msg("failed to truncate brands table")
		return err
	}

	log.Info().Msg("seeding brands")
	brandData, err := os.ReadFile("seeds/brands.json")
	if err != nil {
		log.Error().Err(err).Msg("failed to read brand data")
		return err
	}

	var brands []Brand
	err = json.Unmarshal(brandData, &brands)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse brand data")
		return err
	}

	log.Info().Msgf("brand count: %d", len(brands))
	params := make([]repository.SeedBrandsParams, len(brands))
	for i, brand := range brands {
		params[i] = repository.SeedBrandsParams{
			Name:         brand.Name,
			Slug:         utils.Slugify(brand.Name),
			Description:  &brand.Description,
			DisplayOrder: utils.Int32Ptr(int32(brand.SortOrder)),
			Published:    brand.Published,
		}
	}
	_, err = repo.SeedBrands(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("failed to create brands")
		return err
	}
	log.Info().Msg("brands created")
	return nil
}

func seedCollections(ctx context.Context, repo repository.Store) error {
	_, err := repo.QueryRaw(ctx, "TRUNCATE TABLE collections CASCADE;")
	if err != nil {
		log.Error().Err(err).Msg("failed to truncate collections table")
		return err
	}

	log.Info().Msg("seeding collections")
	collectionData, err := os.ReadFile("seeds/collections.json")
	if err != nil {
		log.Error().Err(err).Msg("failed to read collection data")
		return err
	}

	var collections []Collection
	err = json.Unmarshal(collectionData, &collections)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse collection data")
		return err
	}

	log.Info().Msgf("collection count: %d", len(collections))
	params := make([]repository.SeedCollectionsParams, len(collections))
	for i, collection := range collections {
		params[i] = repository.SeedCollectionsParams{
			Name:        collection.Name,
			Description: &collection.Description,
			ImageUrl:    nil,
			Slug:        collection.Slug,
		}
	}
	_, err = repo.SeedCollections(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("failed to create collections")
		return err
	}
	log.Info().Msg("collections created")
	return nil
}

func seedUsers(ctx context.Context, pg repository.Store) error {
	_, err := pg.QueryRaw(ctx, "TRUNCATE TABLE users CASCADE;")
	if err != nil {
		log.Error().Err(err).Msg("failed to truncate users table")
		return err
	}

	log.Info().Str("from file", "users.json").Msg("parsing user data")
	userData, err := os.ReadFile("seeds/users.json")
	if err != nil {
		log.Error().Err(err).Msg("failed to read user data")
		return err
	}

	var users []User
	log.Info().Msg("parsing users from data")
	err = json.Unmarshal(userData, &users)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse user data")
		return err
	}

	log.Info().Int("user count", len(users)).Msg("creating users")
	params := make([]repository.SeedUsersParams, len(users))
	role, err := pg.GetRoleByCode(ctx, "user")
	if err != nil {
		log.Error().Err(err).Msg("failed to get user role")
		return err
	}
	adminRole, err := pg.GetRoleByCode(ctx, "admin")
	if err != nil {
		log.Error().Err(err).Msg("failed to get admin role")
		return err
	}
	for i, user := range users {
		hashed, _ := auth.HashPwd(user.Password)
		fields := strings.Fields(user.FullName)
		params[i] = repository.SeedUsersParams{
			Email:          user.Email,
			Username:       user.Username,
			PhoneNumber:    user.Phone,
			HashedPassword: hashed,
			FirstName:      user.FullName,
			LastName:       "",
			RoleID:         role.ID,
		}
		if len(fields) > 1 {
			params[i].FirstName = fields[0]
			params[i].LastName = strings.Join(fields[1:], " ")
		}
		if params[i].Username == "admin" {
			params[i].RoleID = adminRole.ID
		}
	}

	_, err = pg.SeedUsers(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("failed to create users")
	}

	log.Info().Msg("users created")
	return nil
}

func seedAttributes(ctx context.Context, repo repository.Store) error {
	_, err := repo.QueryRaw(ctx, "TRUNCATE TABLE attributes CASCADE;")
	if err != nil {
		log.Error().Err(err).Msg("failed to truncate attributes table")
		return err
	}

	log.Info().Msg("seeding attributes")
	for _, attribute := range AttributesSeed {
		_, err := repo.CreateAttribute(ctx, attribute)
		if err != nil {
			log.Error().Err(err).Msg("failed to create attribute")
			return err
		}
	}

	log.Info().Msg("attributes created")
	return nil
}

func seedAttributeValues(ctx context.Context, repo repository.Store) error {
	_, err := repo.QueryRaw(ctx, "TRUNCATE TABLE attribute_values CASCADE;")
	if err != nil {
		log.Error().Err(err).Msg("failed to truncate attribute values table")
		return err
	}
	log.Info().Msg("seeding attribute values")
	attributeData, err := os.ReadFile("seeds/attribute_values.json")
	if err != nil {
		log.Error().Err(err).Msg("failed to read attribute values data")
		return err
	}

	var attributeValues []AttributeValues
	err = json.Unmarshal(attributeData, &attributeValues)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse attribute values data")
		return err
	}

	for _, attrVal := range attributeValues {
		attr, err := repo.GetAttributeByName(ctx, attrVal.Attribute)
		if err != nil {
			log.Error().Err(err).Str("attribute", attrVal.Attribute).Msg("failed to get attribute")
			continue
		}

		var params []repository.CreateAttributeValuesParams
		for _, value := range attrVal.Values {
			params = append(params, repository.CreateAttributeValuesParams{
				AttributeID: attr.ID,
				Value:       value,
			})
		}

		_, err = repo.CreateAttributeValues(ctx, params)
		if err != nil {
			log.Error().Err(err).Str("attribute", attrVal.Attribute).Msg("failed to create attribute values")
		}
	}

	log.Info().Msg("attribute values created")
	return nil
}

func seedDiscounts(ctx context.Context, repo repository.Store) error {
	_, err := repo.QueryRaw(ctx, "TRUNCATE TABLE discounts CASCADE;")
	if err != nil {
		log.Error().Err(err).Msg("failed to truncate discounts table")
		return err
	}

	log.Info().Msg("seeding discounts")
	discountData, err := os.ReadFile("seeds/discounts.json")
	if err != nil {
		log.Error().Err(err).Msg("failed to read discount data")
		return err
	}

	var discounts []Discount
	err = json.Unmarshal(discountData, &discounts)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse discount data")
		return err
	}

	log.Info().Msgf("discount count: %d", len(discounts))
	var params []repository.SeedDiscountsParams
	for _, discount := range discounts {
		validFrom := pgtype.Timestamptz{}
		t, _ := time.Parse(time.RFC3339, discount.ValidFrom)
		validFrom.Scan(t)

		var validUntil pgtype.Timestamptz
		if discount.ValidUntil != nil {
			t, _ := time.Parse(time.RFC3339, *discount.ValidUntil)
			validUntil.Scan(t)
		}

		var minOrderValue, maxDiscountAmt pgtype.Numeric
		if discount.MinOrderValue != nil {
			minOrderValue = utils.GetPgNumericFromFloat(*discount.MinOrderValue)
		}
		if discount.MaxDiscountAmt != nil {
			maxDiscountAmt = utils.GetPgNumericFromFloat(*discount.MaxDiscountAmt)
		}

		params = append(params, repository.SeedDiscountsParams{
			Code:              discount.Code,
			Name:              discount.Name,
			Description:       &discount.Description,
			DiscountType:      repository.DiscountType(discount.DiscountType),
			DiscountValue:     utils.GetPgNumericFromFloat(discount.DiscountValue),
			MinOrderValue:     minOrderValue,
			MaxDiscountAmount: maxDiscountAmt,
			UsageLimit:        discount.UsageLimit,
			UsagePerUser:      discount.UsagePerUser,
			IsActive:          discount.IsActive,
			IsStackable:       discount.IsStackable,
			Priority:          &discount.Priority,
			ValidFrom:         validFrom,
			ValidUntil:        validUntil,
		})
	}

	_, err = repo.SeedDiscounts(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("failed to create discounts")
		return err
	}
	log.Info().Msg("discounts created")
	return nil
}

func seedShippingMethods(ctx context.Context, repo repository.Store) error {
	_, err := repo.QueryRaw(ctx, "TRUNCATE TABLE shipping_methods CASCADE;")
	if err != nil {
		log.Error().Err(err).Msg("failed to truncate shipping methods table")
		return err
	}

	log.Info().Msg("seeding shipping methods")
	methodData, err := os.ReadFile("seeds/shipping_methods.json")
	if err != nil {
		log.Error().Err(err).Msg("failed to read shipping method data")
		return err
	}

	var methods []ShippingMethod
	err = json.Unmarshal(methodData, &methods)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse shipping method data")
		return err
	}

	log.Info().Msgf("shipping method count: %d", len(methods))
	var params []repository.SeedShippingMethodsParams
	for _, method := range methods {
		params = append(params, repository.SeedShippingMethodsParams{
			Name:                  method.Name,
			Description:           &method.Description,
			IsActive:              method.IsActive,
			RequiresAddress:       method.RequiresAddress,
			EstimatedDeliveryTime: &method.EstimatedDeliveryTime,
		})
	}

	_, err = repo.SeedShippingMethods(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("failed to create shipping methods")
		return err
	}
	log.Info().Msg("shipping methods created")
	return nil
}

func seedShippingZones(ctx context.Context, repo repository.Store) error {
	_, err := repo.QueryRaw(ctx, "TRUNCATE TABLE shipping_zones CASCADE;")
	if err != nil {
		log.Error().Err(err).Msg("failed to truncate shipping zones table")
		return err
	}

	log.Info().Msg("seeding shipping zones")
	zoneData, err := os.ReadFile("seeds/shipping_zones.json")
	if err != nil {
		log.Error().Err(err).Msg("failed to read shipping zone data")
		return err
	}

	var zones []ShippingZone
	err = json.Unmarshal(zoneData, &zones)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse shipping zone data")
		return err
	}

	log.Info().Msgf("shipping zone count: %d", len(zones))
	var params []repository.SeedShippingZonesParams
	for _, zone := range zones {
		params = append(params, repository.SeedShippingZonesParams{
			Name:        zone.Name,
			Description: &zone.Description,
			Countries:   zone.Countries,
			States:      zone.States,
			ZipCodes:    zone.ZipCodes,
			IsActive:    zone.IsActive,
		})
	}

	_, err = repo.SeedShippingZones(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("failed to create shipping zones")
		return err
	}
	log.Info().Msg("shipping zones created")
	return nil
}

func seedUserAddresses(ctx context.Context, repo repository.Store) error {
	// First, get users before truncating to avoid lock conflicts
	cnt, err := repo.CountAddresses(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to count addresses")
		return err
	}
	if cnt > 0 {
		// Now truncate user_addresses table
		_, err := repo.QueryRaw(ctx, "TRUNCATE TABLE user_addresses;")
		if err != nil {
			log.Error().Err(err).Msg("failed to truncate user_addresses table")
			return err
		}
	}
	log.Info().Msg("fetching users for address seeding")
	userList, err := repo.GetUsers(ctx, repository.GetUsersParams{
		Limit:  1000,
		Offset: 0,
	})

	if err != nil {
		log.Error().Err(err).Msg("failed to get users for address seeding")
		return err
	}

	if len(userList) == 0 {
		log.Warn().Msg("no users found, skipping address seeding")
		return err
	}

	userIDs := make([]uuid.UUID, len(userList))
	for i, user := range userList {
		userIDs[i] = user.ID
	}

	log.Info().Msg("parsing addresses json")
	addressData, err := os.ReadFile("seeds/addresses.json")
	if err != nil {
		log.Error().Err(err).Msg("failed to read address data")
		return err
	}

	var addresses []Address
	log.Info().Msg("parsing addresses from data")
	err = json.Unmarshal(addressData, &addresses)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse address data")
		return err
	}

	log.Info().Msg("creating addresses")
	log.Info().Int("addresses count", len(addresses))
	params := make([]repository.SeedAddressesParams, len(addresses))
	for i, address := range addresses {
		idx := rand.Intn(len(userIDs))
		params[i] = repository.SeedAddressesParams{
			UserID:      userIDs[idx],
			PhoneNumber: address.Phone,
			Street:      address.Street,
			Ward:        &address.Ward,
			District:    address.District,
			City:        address.City,
		}
	}
	_, err = repo.SeedAddresses(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("failed to seed addresses")
		return err
	}
	log.Info().Msg("addresses created")
	return nil
}

// Helper functions to create product relationships

// createProductCategoryRelationships creates relationships between a product and its categories
func createProductCategoryRelationships(ctx context.Context, repo repository.Store, productID uuid.UUID, categoryNames []string, categories []repository.Category) error {
	var categoryProductParams []repository.AddProductsToCategoryParams

	for _, categoryName := range categoryNames {
		category := slices.IndexFunc(categories, func(c repository.Category) bool {
			return strings.EqualFold(c.Name, categoryName)
		})

		if category != -1 {
			categoryProductParams = append(categoryProductParams, repository.AddProductsToCategoryParams{
				CategoryID: categories[category].ID,
				ProductID:  productID,
			})
		} else {
			log.Warn().Str("category", categoryName).Str("product_id", productID.String()).Msg("category not found for product")
		}
	}

	if len(categoryProductParams) > 0 {
		_, err := repo.AddProductsToCategory(ctx, categoryProductParams)
		if err != nil {
			return fmt.Errorf("failed to add product to categories: %w", err)
		}
		log.Debug().Str("product_id", productID.String()).Int("categories", len(categoryProductParams)).Msg("product categories added")
	}

	return nil
}

// createProductCollectionRelationships creates relationships between a product and its collections
func createProductCollectionRelationships(ctx context.Context, repo repository.Store, productID uuid.UUID, collectionNames []string, collections []repository.Collection) error {
	var collectionProductParams []repository.AddProductsToCollectionParams

	for _, collectionName := range collectionNames {
		collection := slices.IndexFunc(collections, func(c repository.Collection) bool {
			return strings.EqualFold(c.Name, collectionName)
		})

		if collection != -1 {
			collectionProductParams = append(collectionProductParams, repository.AddProductsToCollectionParams{
				CollectionID: collections[collection].ID,
				ProductID:    productID,
			})
		} else {
			log.Warn().Str("collection", collectionName).Str("product_id", productID.String()).Msg("collection not found for product")
		}
	}

	if len(collectionProductParams) > 0 {
		_, err := repo.AddProductsToCollection(ctx, collectionProductParams)
		if err != nil {
			return fmt.Errorf("failed to add product to collections: %w", err)
		}
		log.Debug().Str("product_id", productID.String()).Int("collections", len(collectionProductParams)).Msg("product collections added")
	}

	return nil
}

// createProductAttributeRelationships creates relationships between a product and its attributes
func createProductAttributeRelationships(ctx context.Context, repo repository.Store, productID uuid.UUID, attributes map[string]string) error {
	var productAttributeParams []repository.CreateBulkProductAttributesParams

	for attributeName := range attributes {
		// Get the attribute by name
		attr, err := repo.GetAttributeByName(ctx, attributeName)
		if err != nil {
			log.Warn().Err(err).Str("attribute", attributeName).Str("product_id", productID.String()).Msg("attribute not found for product")
			continue
		}

		productAttributeParams = append(productAttributeParams, repository.CreateBulkProductAttributesParams{
			ProductID:   productID,
			AttributeID: attr.ID,
		})
	}

	if len(productAttributeParams) > 0 {
		_, err := repo.CreateBulkProductAttributes(ctx, productAttributeParams)
		if err != nil {
			return fmt.Errorf("failed to add product attributes: %w", err)
		}
		log.Debug().Str("product_id", productID.String()).Int("attributes", len(productAttributeParams)).Msg("product attributes added")
	}

	return nil
}

// seedExistingProductVariants reads products.json and creates variants for existing products
func seedExistingProductVariants(ctx context.Context, repo repository.Store) error {
	_, err := repo.QueryRaw(ctx, "TRUNCATE TABLE product_variants CASCADE;")
	if err != nil {
		log.Error().Err(err).Msg("failed to truncate product_variants table")
		return err
	}

	log.Info().Msg("seeding variants for existing products")

	// Read products data
	productData, err := os.ReadFile("seeds/products.json")
	if err != nil {
		log.Error().Err(err).Msg("failed to read product data")
		return err
	}

	var products []Product
	err = json.Unmarshal(productData, &products)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse product data")
		return err
	}

	log.Info().Int("product count", len(products)).Msg("processing products for variant creation")

	for _, product := range products {
		// Find the existing product by SKU
		existingProduct, err := repo.GetProductBySku(ctx, repository.GetProductBySkuParams{
			BaseSku: product.Sku,
		})
		if err != nil {
			log.Error().Err(err).Str("sku", product.Sku).Msg("product not found")
			continue
		}

		// Check if variants already exist
		existingVariants, err := repo.GetProductVariantsByProductID(ctx, repository.GetProductVariantsByProductIDParams{
			ProductID: existingProduct.ID,
		})
		if err != nil {
			log.Error().Err(err).Str("product_id", existingProduct.ID.String()).Msg("failed to check existing variants")
			continue
		}

		if len(existingVariants) > 0 {
			log.Info().Str("product", product.Name).Int("existing_variants", len(existingVariants)).Msg("variants already exist, skipping")
			continue
		}

		// Create variants
		if len(product.Variants) > 0 {
			err := seedProductVariants(ctx, repo, existingProduct.ID, product.Variants)
			if err != nil {
				log.Error().Err(err).Str("product", product.Name).Msg("failed to create variants")
			} else {
				log.Info().Str("product", product.Name).Int("variant_count", len(product.Variants)).Msg("variants created successfully")
			}
		} else {
			log.Info().Str("product", product.Name).Msg("no variants defined in seed data")
		}
	}

	log.Info().Msg("variant seeding completed")
	return nil
}
