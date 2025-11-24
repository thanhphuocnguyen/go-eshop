package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
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
	"Brand",
	"Warranty",
	"Storage",
	"Style",
}

type VariantSeedModel struct {
	Price      float64           `json:"price"`
	Stock      int32             `json:"stock"`
	Sku        string            `json:"sku"`
	Discount   int16             `json:"discount,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}
type Product struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Variants    []VariantSeedModel `json:"variants"`
	Price       float64            `json:"price"`
	Stock       int32              `json:"stock"`
	Sku         string             `json:"sku"`
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
				log.Info().Msg("Starting comprehensive database seeding")

				// Seed base data first (sequential order matters)
				seedAttributes(ctx, pg)
				seedAttributeValues(ctx, pg)
				seedCategories(ctx, pg)
				seedBrands(ctx, pg)
				seedCollections(ctx, pg)
				seedShippingMethods(ctx, pg)
				seedShippingZones(ctx, pg)
				seedDiscounts(ctx, pg)

				// Seed users and products in parallel
				var waitGroup sync.WaitGroup
				defer waitGroup.Wait()
				waitGroup.Add(2)
				go func() {
					defer waitGroup.Done()
					seedUsers(ctx, pg)
				}()
				go func() {
					defer waitGroup.Done()
					seedProducts(ctx, pg)
				}()
			} else {
				switch args[0] {
				case "products":
					seedProducts(ctx, pg)
				case "users":
					seedUsers(ctx, pg)
				case "categories":
					seedCategories(ctx, pg)
				case "attributes":
					seedAttributes(ctx, pg)
				case "attribute-values":
					seedAttributeValues(ctx, pg)
				case "brands":
					seedBrands(ctx, pg)
				case "collections":
					seedCollections(ctx, pg)
				case "discounts":
					seedDiscounts(ctx, pg)
				case "shipping":
					seedShippingMethods(ctx, pg)
					seedShippingZones(ctx, pg)
				case "variants":
					seedExistingProductVariants(ctx, pg)
				case "user-addresses":
					seedUserAddresses(ctx, pg)
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

	log.Info().Int("product count", len(products)).Msg("creating products")
	for _, product := range products {
		params := repository.CreateProductParams{
			Name:        product.Name,
			Description: product.Description,
			BasePrice:   utils.GetPgNumericFromFloat(product.Price),
			BaseSku:     product.Sku,
			Slug:        utils.Slugify(product.Name),
		}
		createdProduct, err := repo.CreateProduct(ctx, params)
		if err != nil {
			log.Error().Err(err).Msg("failed to create product")
			continue
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
}

func seedProductVariants(ctx context.Context, repo repository.Repository, productID uuid.UUID, variants []VariantSeedModel) error {
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

func createVariantAttributes(ctx context.Context, repo repository.Repository, variantID uuid.UUID, attributes map[string]string) error {
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

func seedCategories(ctx context.Context, repo repository.Repository) {
	countCategories, err := repo.CountCategories(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to count categories")
		return
	}

	if countCategories > 0 {
		log.Info().Msg("categories already seeded")
		return
	}

	log.Info().Msg("seeding categories")
	categoryData, err := os.ReadFile("seeds/categories.json")
	if err != nil {
		log.Error().Err(err).Msg("failed to read category data")
		return
	}

	var categories []Categories
	log.Info().Msg("parsing category data")
	err = json.Unmarshal(categoryData, &categories)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse category data")
		return
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
		return
	}
	log.Info().Msg("categories created")
}

func seedBrands(ctx context.Context, repo repository.Repository) {
	countBrands, err := repo.CountBrands(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to count brands")
		return
	}

	if countBrands > 0 {
		log.Info().Msg("brands already seeded")
		return
	}

	log.Info().Msg("seeding brands")
	brandData, err := os.ReadFile("seeds/brands.json")
	if err != nil {
		log.Error().Err(err).Msg("failed to read brand data")
		return
	}

	var brands []Brand
	err = json.Unmarshal(brandData, &brands)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse brand data")
		return
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
		return
	}
	log.Info().Msg("brands created")
}

func seedCollections(ctx context.Context, repo repository.Repository) {
	countCollections, err := repo.CountCollections(ctx)
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
	err = json.Unmarshal(collectionData, &collections)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse collection data")
		return
	}

	log.Info().Msgf("collection count: %d", len(collections))
	params := make([]repository.SeedCollectionsParams, len(collections))
	for i, collection := range collections {
		params[i] = repository.SeedCollectionsParams{
			Name:        collection.Name,
			Description: &collection.Description,
			ImageUrl:    nil,
		}
	}
	_, err = repo.SeedCollections(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("failed to create collections")
		return
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

	log.Info().Str("from file", "users.json").Msg("parsing user data")
	userData, err := os.ReadFile("seeds/users.json")
	if err != nil {
		log.Error().Err(err).Msg("failed to read user data")
		return
	}

	var users []User
	log.Info().Msg("parsing users from data")
	err = json.Unmarshal(userData, &users)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse user data")
		return
	}

	log.Info().Int("user count", len(users)).Msg("creating users")
	params := make([]repository.SeedUsersParams, len(users))
	role, err := pg.GetRoleByCode(ctx, "user")
	adminRole, err := pg.GetRoleByCode(ctx, "admin")
	if err != nil {
		log.Error().Err(err).Msg("failed to get user role")
		return
	}
	for i, user := range users {
		hashed, _ := auth.HashPassword(user.Password)
		params[i] = repository.SeedUsersParams{
			Email:          user.Email,
			Username:       user.Username,
			PhoneNumber:    user.Phone,
			HashedPassword: hashed,
			FirstName:      user.FullName,
			LastName:       "",
			RoleID:         role.ID,
		}
		if params[i].Username == "admin" {
			params[i].RoleID = adminRole.ID
		}
	}

	_, err = pg.SeedUsers(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("failed to create users")
	}

	seedUserAddresses(ctx, pg)
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

	log.Info().Msg("seeding attributes")
	for _, attribute := range AttributesSeed {
		_, err := repo.CreateAttribute(ctx, attribute)
		if err != nil {
			log.Error().Err(err).Msg("failed to create attribute")
			return
		}
	}

	log.Info().Msg("attributes created")
}

func seedAttributeValues(ctx context.Context, repo repository.Repository) {
	log.Info().Msg("seeding attribute values")
	attributeData, err := os.ReadFile("seeds/attribute_values.json")
	if err != nil {
		log.Error().Err(err).Msg("failed to read attribute values data")
		return
	}

	var attributeValues []AttributeValues
	err = json.Unmarshal(attributeData, &attributeValues)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse attribute values data")
		return
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
}

func seedDiscounts(ctx context.Context, repo repository.Repository) {
	countDiscounts, err := repo.CountDiscounts(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to count discounts")
		return
	}

	if countDiscounts > 0 {
		log.Info().Msg("discounts already seeded")
		return
	}

	log.Info().Msg("seeding discounts")
	discountData, err := os.ReadFile("seeds/discounts.json")
	if err != nil {
		log.Error().Err(err).Msg("failed to read discount data")
		return
	}

	var discounts []Discount
	err = json.Unmarshal(discountData, &discounts)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse discount data")
		return
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
			minOrderValue.Scan(*discount.MinOrderValue)
		}
		if discount.MaxDiscountAmt != nil {
			maxDiscountAmt.Scan(*discount.MaxDiscountAmt)
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
			IsActive:          &discount.IsActive,
			IsStackable:       &discount.IsStackable,
			Priority:          &discount.Priority,
			ValidFrom:         validFrom,
			ValidUntil:        validUntil,
		})
	}

	_, err = repo.SeedDiscounts(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("failed to create discounts")
		return
	}
	log.Info().Msg("discounts created")
}

func seedShippingMethods(ctx context.Context, repo repository.Repository) {
	countMethods, err := repo.CountShippingMethods(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to count shipping methods")
		return
	}

	if countMethods > 0 {
		log.Info().Msg("shipping methods already seeded")
		return
	}

	log.Info().Msg("seeding shipping methods")
	methodData, err := os.ReadFile("seeds/shipping_methods.json")
	if err != nil {
		log.Error().Err(err).Msg("failed to read shipping method data")
		return
	}

	var methods []ShippingMethod
	err = json.Unmarshal(methodData, &methods)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse shipping method data")
		return
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
		return
	}
	log.Info().Msg("shipping methods created")
}

func seedShippingZones(ctx context.Context, repo repository.Repository) {
	countZones, err := repo.CountShippingZones(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to count shipping zones")
		return
	}

	if countZones > 0 {
		log.Info().Msg("shipping zones already seeded")
		return
	}

	log.Info().Msg("seeding shipping zones")
	zoneData, err := os.ReadFile("seeds/shipping_zones.json")
	if err != nil {
		log.Error().Err(err).Msg("failed to read shipping zone data")
		return
	}

	var zones []ShippingZone
	err = json.Unmarshal(zoneData, &zones)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse shipping zone data")
		return
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
		return
	}
	log.Info().Msg("shipping zones created")
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
		log.Error().Err(err).Msg("failed to read address data")
		return
	}

	var addresses []Address
	log.Info().Msg("parsing addresses from data")
	err = json.Unmarshal(addressData, &addresses)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse address data")
		return
	}
	userList, err := repo.GetUsers(ctx, repository.GetUsersParams{
		Limit:  1000,
		Offset: 0,
	})

	userIDs := make([]uuid.UUID, len(userList))
	for i, user := range userList {
		userIDs[i] = user.ID
	}

	if err != nil {
		log.Error().Err(err).Msg("failed to get created users")
		return
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
		return
	}
	log.Info().Msg("addresses created")
}

// seedExistingProductVariants reads products.json and creates variants for existing products
func seedExistingProductVariants(ctx context.Context, repo repository.Repository) {
	log.Info().Msg("seeding variants for existing products")

	// Read products data
	productData, err := os.ReadFile("seeds/products.json")
	if err != nil {
		log.Error().Err(err).Msg("failed to read product data")
		return
	}

	var products []Product
	err = json.Unmarshal(productData, &products)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse product data")
		return
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
}
