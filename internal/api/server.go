package api

import (
	"encoding/gob"
	"net/http"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("uuidslice", uuidSlice)
	}

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
	router.Use(gin.Recovery())
	if sv.config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
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

func (s *Server) Server(addr string) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: s.router.Handler(),
	}
}

type DashboardData struct {
	Categories  []CategoryResponse `json:"categories"`
	Collections []CategoryResponse `json:"collections"`
}

func (s *Server) getHomePageHandler(ctx *gin.Context) {
	var wg sync.WaitGroup
	wg.Add(2)

	// Use channels to collect results
	categoriesChan := make(chan []CategoryResponse, 1)
	collectionsChan := make(chan []CategoryResponse, 1)

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
			collectionsChan <- []CategoryResponse{}
			return
		}

		collectionModel := make([]CategoryResponse, len(collectionRows))
		for i, collection := range collectionRows {
			collectionModel[i] = CategoryResponse{
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
