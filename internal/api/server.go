package api

import (
	"encoding/gob"
	"net/http"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v81"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/thanhphuocnguyen/go-eshop/config"
	docs "github.com/thanhphuocnguyen/go-eshop/docs"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/worker"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
	"github.com/thanhphuocnguyen/go-eshop/pkg/cache"
	"github.com/thanhphuocnguyen/go-eshop/pkg/payment"
	"github.com/thanhphuocnguyen/go-eshop/pkg/upload"
)

// gin-swagger middleware
// swagger embed files

// @BasePath /api/v1
// @title           E-Commerce API
// @description     This is a sample server for a simple e-commerce API.
// @BasePath /api/v1
// @host      localhost:4000
type Server struct {
	config          config.Config
	router          *gin.Engine
	repo            repository.Repository
	tokenGenerator  auth.TokenGenerator
	uploadService   upload.UploadService
	paymentCtx      *payment.PaymentContext
	cacheService    cache.Cache
	taskDistributor worker.TaskDistributor
}

func NewAPI(
	cfg config.Config,
	repo repository.Repository,
	cacheService cache.Cache,
	taskDistributor worker.TaskDistributor,
	uploadService upload.UploadService,
	paymentCtx *payment.PaymentContext,
) (*Server, error) {
	tokenGenerator, err := auth.NewJwtGenerator(cfg.SymmetricKey)
	if err != nil {
		return nil, err
	}
	server := &Server{
		tokenGenerator:  tokenGenerator,
		repo:            repo,
		config:          cfg,
		taskDistributor: taskDistributor,
		uploadService:   uploadService,
		cacheService:    cacheService,
		paymentCtx:      paymentCtx,
	}
	server.initializeRouter()
	return server, nil
}

func (sv *Server) initializeRouter() {
	router := gin.Default()
	gob.Register(&stripe.PaymentIntent{})

	// Setup environment mode
	sv.setupEnvironmentMode(router)

	// Load HTML templates
	router.LoadHTMLGlob("static/templates/*")

	// Setup CORS
	router.Use(sv.setupCORS())

	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	docs.SwaggerInfo.BasePath = "/api/v1"

	router.Static("/assets", "./assets")

	router.GET("verify-email", sv.verifyEmailHandler)
	// Setup API routes
	v1 := router.Group("/api/v1")
	{
		// Health check endpoint
		v1.GET("health", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"status ": "ok"})
		})

		v1.GET("homepage", sv.getHomePageHandler)

		// Register API route groups
		sv.setupAuthRoutes(v1)
		sv.setupAdminRoutes(v1)
		sv.setupUserRoutes(v1)
		sv.setupProductRoutes(v1)
		sv.setupImageRoutes(v1)
		sv.setupCartRoutes(v1)
		sv.setupOrderRoutes(v1)
		sv.setupPaymentRoutes(v1)
		sv.setupCategoryRoutes(v1)
		sv.setupCollectionRoutes(v1)
		sv.setupBrandRoutes(v1)
		sv.setupRatingRoutes(v1)
	}

	// Setup webhook routes
	sv.setupWebhookRoutes(router)

	// Setup Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	sv.router = router
}

// Setup environment mode based on configuration
func (sv *Server) setupEnvironmentMode(router *gin.Engine) {
	if sv.config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
		router.Use(gin.Recovery())
	}
	if sv.config.Env == "development" {
		router.Use(gin.Logger())
	}
}

// Setup CORS configuration
func (sv *Server) setupCORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001", "http://localhost:8080"},
		AllowHeaders:     []string{"Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"},
		AllowFiles:       true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowCredentials: true,
		// AllowAllOrigins:  sv.config.Env == "development",
	})
}

func (sv *Server) setupAdminRoutes(rg *gin.RouterGroup) {
	admin := rg.Group("/admin", authMiddleware(sv.tokenGenerator), roleMiddleware(sv.repo, repository.UserRoleAdmin))
	{
		users := admin.Group("users")
		{
			users.GET("", sv.getUsersHandler)
		}

		productGroup := admin.Group("products")
		{
			productGroup.POST("", sv.createProduct)
			productGroup.PUT(":id", sv.updateProduct)
			productGroup.PUT(":id/variants", sv.updateProductVariants)
			productGroup.DELETE(":id", sv.removeProduct)
		}

		attributeGroup := admin.Group("attributes")
		{
			attributeGroup.POST("", sv.createAttributeHandler)
			attributeGroup.GET("", sv.getAttributesHandler)
			attributeGroup.GET(":id", sv.getAttributeByIDHandler)
			attributeGroup.PUT(":id", sv.updateAttributeHandler)
			attributeGroup.DELETE(":id", sv.deleteAttribute)
		}

		adminOrder := admin.Group("orders")
		{
			adminOrder.GET("", sv.getAdminOrdersHandler)
			adminOrder.GET(":id", sv.getAdminOrderDetailHandler)
			adminOrder.PUT(":id/refund", sv.refundOrder)
			adminOrder.PUT(":id/status", sv.changeOrderStatus)
		}

		categories := admin.Group("categories")
		{
			categories.GET("", sv.getAdminCategoriesHandler)
			categories.GET(":id", sv.getCategoryByID)
			categories.POST("", sv.addCategoryHandler)
			categories.PUT(":id", sv.updateCategoryHandler)
			categories.DELETE(":id", sv.deleteCategory)
		}

		brands := admin.Group("brands")
		{

			brands.GET("", sv.getBrandsHandler)
			brands.POST("", sv.createBrandHandler)
			brands.GET(":id", sv.getBrandByIDHandler)
			brands.PUT(":id", sv.updateBrand)
			brands.DELETE(":id", sv.deleteBrand)
		}

		collections := admin.Group("collections")
		{
			collections.GET("", sv.getCollectionsHandler)
			collections.POST("", sv.addCollectionHandler)
			collections.GET(":id", sv.getCollectionByIDHandler)
			collections.PUT(":id", sv.updateCollectionHandler)
			collections.DELETE(":id", sv.deleteCollection)
		}

		images := admin.Group("images")
		{
			productImages := images.Group("products")
			productImages.POST(":entity_id", sv.uploadProductImagesHandler)
			productImages.DELETE(":entity_id", sv.removeImage)
		}

		ratings := admin.Group("ratings")
		{
			ratings.DELETE(":id", sv.deleteRatingHandler)
			ratings.POST(":id/approve", sv.approveRatingHandler)
			ratings.POST(":id/ban", sv.banUserRatingHandler)
		}

	}
}

// Setup authentication routes
func (sv *Server) setupAuthRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("register", sv.registerHandler)
		auth.POST("login", sv.loginHandler)
		auth.POST("refresh-token", sv.refreshTokenHandler)
	}
}

// Setup user-related routes
func (sv *Server) setupUserRoutes(rg *gin.RouterGroup) {
	user := rg.Group("/user", authMiddleware(sv.tokenGenerator))
	{
		user.GET("", sv.getUserHandler)
		user.PATCH("", sv.updateUserHandler)
		user.POST("send-verify-email", sv.sendVerifyEmailHandler)

		userAddress := user.Group("addresses")
		{
			userAddress.POST("", sv.createAddressHandler)
			userAddress.PATCH(":id/default", sv.setDefaultAddressHandler)
			userAddress.GET("", sv.getAddressesHandlers)
			userAddress.PATCH(":id", sv.updateAddressHandlers)
			userAddress.DELETE(":id", sv.removeAddressHandlers)
		}
	}
}

// Setup product-related routes
func (sv *Server) setupProductRoutes(rg *gin.RouterGroup) {
	product := rg.Group("products")
	{
		product.GET("", sv.getProductsHandler)
		product.GET(":id", sv.getProductDetailHandler)
		product.GET(":id/ratings", sv.getProductRatingsHandler)
	}
}

// Setup image-related routes
func (sv *Server) setupImageRoutes(rg *gin.RouterGroup) {
	images := rg.Group("images", authMiddleware(sv.tokenGenerator))
	{
		images.DELETE("remove-external/:public_id", sv.removeImageByPublicID)
		images.GET("", sv.getProductImagesHandler)
	}
}

// Setup cart-related routes
func (sv *Server) setupCartRoutes(rg *gin.RouterGroup) {
	cart := rg.Group("/cart", authMiddleware(sv.tokenGenerator))
	{
		cart.POST("", sv.createCart)
		cart.GET("", sv.getCartHandler)
		cart.POST("checkout", sv.checkoutHandler)
		cart.PUT("clear", sv.clearCart)
		cartItem := cart.Group("item")
		cartItem.DELETE(":id", sv.removeCartItem)
		cartItem.PUT(":id/quantity", sv.updateCartItemQtyHandler)
	}
}

// Setup order-related routes
func (sv *Server) setupOrderRoutes(rg *gin.RouterGroup) {
	order := rg.Group("/orders", authMiddleware(sv.tokenGenerator))
	{
		order.GET("", sv.getOrdersHandler)
		order.GET(":id", sv.getOrderDetailHandler)
		order.PUT(":id/confirm-received", sv.confirmOrderPayment)
		order.PUT(":id/cancel", sv.cancelOrder)
	}
}

// Setup payment-related routes
func (sv *Server) setupPaymentRoutes(rg *gin.RouterGroup) {
	payment := rg.Group("/payments").Use(authMiddleware(sv.tokenGenerator))
	{
		payment.GET(":id", sv.getPaymentHandler)
		payment.GET("stripe-config", sv.getStripeConfig)
		payment.POST("", sv.createPaymentIntentHandler)
		payment.PUT(":order_id", sv.changePaymentStatus)
	}
}

// Setup category-related routes
func (sv *Server) setupCategoryRoutes(rg *gin.RouterGroup) {
	category := rg.Group("categories")
	{
		category.GET("", sv.getCategoriesHandler)
		category.GET(":slug", sv.getCategoryBySlug)
	}
}

// Setup collection-related routes
func (sv *Server) setupCollectionRoutes(rg *gin.RouterGroup) {
	collections := rg.Group("collections")
	{
		collections.GET("", sv.getShopCollectionsHandler)
	}
}

// Setup brand-related routes
func (sv *Server) setupBrandRoutes(rg *gin.RouterGroup) {
	brands := rg.Group("brands")
	{
		brands.GET("", sv.getShopBrandsHandler)
	}
}

// Setup brand-related routes
func (sv *Server) setupRatingRoutes(rg *gin.RouterGroup) {
	ratings := rg.Group("ratings", authMiddleware(sv.tokenGenerator))
	{
		ratings.POST("", sv.postRatingHandler)
		ratings.POST("helpful", sv.postRatingHelpfulHandler)
		ratings.POST(":id/reply", sv.postReplyRatingHandler)
	}
}

// Setup webhook routes
func (sv *Server) setupWebhookRoutes(router *gin.Engine) {
	webhook := router.Group("/webhook/v1")
	{
		webhook.POST("stripe", sv.stripeWebhook)
	}
}

func (s *Server) Server(addr string) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: s.router.Handler(),
	}
}

type DashboardData struct {
	Categories  []CategoryResponse   `json:"categories"`
	Collections []CollectionResponse `json:"collections"`
}

func (s *Server) getHomePageHandler(ctx *gin.Context) {
	var wg sync.WaitGroup
	wg.Add(2)

	// Use channels to collect results
	categoriesChan := make(chan []CategoryResponse, 1)
	collectionsChan := make(chan []CollectionResponse, 1)

	// Fetch categories
	go func() {
		defer wg.Done()
		categoryRows, err := s.repo.GetCategories(ctx, repository.GetCategoriesParams{
			Limit:  5,
			Offset: 0,
		})

		if err != nil {
			categoriesChan <- []CategoryResponse{}
			return
		}

		categoryModel := make([]CategoryResponse, len(categoryRows))
		for i, category := range categoryRows {
			categoryModel[i] = CategoryResponse{
				ID:          category.ID.String(),
				Name:        category.Name,
				Description: category.Description,
				ImageUrl:    category.ImageUrl,
				Slug:        category.Slug,
			}
		}

		categoriesChan <- categoryModel
	}()

	// Fetch collections
	go func() {
		defer wg.Done()
		collectionRows, err := s.repo.GetCollections(ctx, repository.GetCollectionsParams{
			Limit:  5,
			Offset: 0,
		})

		if err != nil {
			collectionsChan <- []CollectionResponse{}
			return
		}

		collectionModel := make([]CollectionResponse, len(collectionRows))
		for i, collection := range collectionRows {
			collectionModel[i] = CollectionResponse{
				ID:          collection.ID.String(),
				Name:        collection.Name,
				Description: collection.Description,
				ImageUrl:    collection.ImageUrl,
				Slug:        collection.Slug,
			}
		}

		collectionsChan <- collectionModel
	}()

	// Wait for all goroutines to finish
	wg.Wait()

	// Read the results from channels
	categories := <-categoriesChan
	collections := <-collectionsChan

	// Create the response
	response := DashboardData{
		Categories:  categories,
		Collections: collections,
	}

	ctx.JSON(http.StatusOK, createSuccessResponse(ctx, response, "Get homepage data successfully", nil, nil))
}
