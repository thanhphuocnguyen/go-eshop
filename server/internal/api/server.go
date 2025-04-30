package api

import (
	"encoding/gob"
	"net/http"
	"time"

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
	repo            repository.Repository
	taskDistributor worker.TaskDistributor
	tokenGenerator  auth.TokenGenerator
	uploadService   upload.UploadService
	router          *gin.Engine
	paymentCtx      *payment.PaymentContext
}

func NewAPI(
	cfg config.Config,
	repo repository.Repository,
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

	// Setup CORS
	router.Use(sv.setupCORS())

	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	docs.SwaggerInfo.BasePath = "/api/v1"

	router.Static("/assets", "./assets")

	// Setup API routes
	v1 := router.Group("/api/v1")
	{
		// Health check endpoint
		v1.GET("health", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"status ": "ok"})
		})

		v1.GET("verify-email", sv.verifyEmailHandler)

		// Register API route groups
		sv.setupAuthRoutes(v1)
		sv.setupUserRoutes(v1)
		sv.setupModeratorRoutes(v1)
		sv.setupProductRoutes(v1)
		sv.setupImageRoutes(v1)
		sv.setupAttributeRoutes(v1)
		sv.setupCartRoutes(v1)
		sv.setupOrderRoutes(v1)
		sv.setupPaymentRoutes(v1)
		sv.setupCategoryRoutes(v1)
		sv.setupCollectionRoutes(v1)
		sv.setupBrandRoutes(v1)
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

// Setup authentication routes
func (sv *Server) setupAuthRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("register", sv.registerHandler)
		auth.POST("login", sv.loginHandler)
		auth.POST("refresh-token", sv.refreshTokenHandler)
		authRoutes := auth.Use(authMiddleware(sv.tokenGenerator))
		authRoutes.POST("send-verify-email", sv.sendVerifyEmailHandler)
	}
}

// Setup user-related routes
func (sv *Server) setupUserRoutes(rg *gin.RouterGroup) {
	user := rg.Group("/user", authMiddleware(sv.tokenGenerator))
	{
		user.GET("", sv.getUser)
		user.PATCH("", sv.updateUser)

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

// Setup moderator-specific routes
func (sv *Server) setupModeratorRoutes(rg *gin.RouterGroup) {
	moderatorRoutes := rg.Group("/moderator",
		authMiddleware(sv.tokenGenerator),
		roleMiddleware(
			sv.repo,
			repository.UserRoleAdmin,
			repository.UserRoleModerator))
	{
		moderatorRoutes.GET("list", sv.listUsers)
	}
}

// Setup product-related routes
func (sv *Server) setupProductRoutes(rg *gin.RouterGroup) {
	product := rg.Group("products")
	{
		product.GET("", sv.getProducts)
		product.GET(":id", sv.getProductDetail)
		productAdmin := product.Use(authMiddleware(sv.tokenGenerator), roleMiddleware(sv.repo, repository.UserRoleAdmin))
		{
			productAdmin.POST("", sv.createProduct)
			productAdmin.PUT(":id", sv.updateProduct)
			productAdmin.PUT(":id/variants", sv.updateProductVariants)
			productAdmin.DELETE(":id", sv.removeProduct)
		}
	}
}

// Setup image-related routes
func (sv *Server) setupImageRoutes(rg *gin.RouterGroup) {
	images := rg.Group(
		"images",
		authMiddleware(sv.tokenGenerator),
		roleMiddleware(sv.repo, repository.UserRoleAdmin, repository.UserRoleModerator))
	{
		images.DELETE("remove-external/:public_id", sv.removeImageByPublicID)
		images.GET("", sv.getImages)
		images.POST("product/:entity_id", sv.uploadProductImages)
		images.DELETE(":entity_id", sv.removeImage)
	}
}

// Setup attribute-related routes
func (sv *Server) setupAttributeRoutes(rg *gin.RouterGroup) {
	attribute := rg.Group(
		"attributes",
		authMiddleware(sv.tokenGenerator),
		roleMiddleware(sv.repo, repository.UserRoleAdmin))
	{
		// C
		attribute.POST("", sv.createAttribute)
		// R
		attribute.GET("", sv.getAttributesHandler)
		attribute.GET(":id", sv.getAttributeByIDHandler)
		// U
		attribute.PUT(":id", sv.updateAttributeHandler)
		// D
		attribute.DELETE(":id", sv.deleteAttribute)
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
	order := rg.Group("/order", authMiddleware(sv.tokenGenerator))
	{
		order.GET("list", sv.getOrdersHandler)
		order.GET(":id", sv.getOrderDetailHandler)
		order.PUT(":id/cancel", sv.cancelOrder)

		adminOrder := order.Group(
			"",
			roleMiddleware(sv.repo, repository.UserRoleAdmin, repository.UserRoleModerator))
		adminOrder.PUT(":id/refund", sv.refundOrder)
		adminOrder.PUT(":id/status", sv.changeOrderStatus)
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
		category.GET("", sv.getCategories)
		//TODO: category.GET("dashboard", sv.getCategoriesForDashboard)
		category.GET(":id", sv.getCategoryByID)
		category.GET(":id/products", sv.getProductsByCategory)

		// Routes that require authentication
		categoryAuthRoutes := category.Group("").Use(
			authMiddleware(sv.tokenGenerator),
			roleMiddleware(sv.repo, repository.UserRoleAdmin),
		)
		categoryAuthRoutes.POST("", sv.addCategoryHandler)
		categoryAuthRoutes.PUT(":id", sv.updateCategory)
		categoryAuthRoutes.DELETE(":id", sv.deleteCategory)
	}
}

// Setup collection-related routes
func (sv *Server) setupCollectionRoutes(rg *gin.RouterGroup) {
	collection := rg.Group("collections")
	{
		collection.GET("", sv.getCollections)
		collection.GET(":id", sv.getCollectionByID)
		collection.GET(":id/products", sv.getProductsByCollection)

		// Routes that require authentication
		collectionAuthRoutes := collection.Group("").Use(
			authMiddleware(sv.tokenGenerator),
			roleMiddleware(sv.repo, repository.UserRoleAdmin),
		)
		collectionAuthRoutes.POST("", sv.createCollection)
		collectionAuthRoutes.PUT(":id", sv.updateCollection)
		collectionAuthRoutes.DELETE(":id", sv.deleteCollection)
	}
}

// Setup brand-related routes
func (sv *Server) setupBrandRoutes(rg *gin.RouterGroup) {
	brand := rg.Group("brands")
	{
		brand.GET("", sv.getBrandsHandler)
		brand.GET(":id", sv.getBrandByIDHandler)
		brand.GET(":id/products", sv.getProductsByBrandHandler)

		// Routes that require authentication
		brandAuthRoutes := brand.Group("").Use(
			authMiddleware(sv.tokenGenerator),
			roleMiddleware(sv.repo, repository.UserRoleAdmin),
		)
		brandAuthRoutes.POST("", sv.createBrandHandler)
		brandAuthRoutes.PUT(":id", sv.updateBrand)
		brandAuthRoutes.DELETE(":id", sv.deleteBrand)
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

// Response types - unchanged
type ApiResponse[T any] struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       *T          `json:"data,omitempty"`
	Error      *ApiError   `json:"error,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Meta       *MetaInfo   `json:"meta"`
}

// Error structure for detailed errors
type ApiError struct {
	Code    string `json:"code"`
	Details string `json:"details"`
	Stack   string `json:"stack,omitempty"` // Hide in production
}

// Pagination info (for paginated endpoints)
type Pagination struct {
	Total           int64 `json:"total"`
	Page            int64 `json:"page"`
	PageSize        int64 `json:"pageSize"`
	TotalPages      int   `json:"totalPages"`
	HasNextPage     bool  `json:"hasNextPage"`
	HasPreviousPage bool  `json:"hasPreviousPage"`
}

// Meta info about the request
type MetaInfo struct {
	Timestamp string `json:"timestamp"`
	RequestID string `json:"requestId"`
	Path      string `json:"path"`
	Method    string `json:"method"`
}

type PaginationQueryParams struct {
	Page     int64 `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize int64 `form:"page_size,default=20" binding:"omitempty,min=1,max=100"`
}

func createErrorResponse[T any](code string, msg string, err error) ApiResponse[T] {
	return ApiResponse[T]{
		Success: false,
		Data:    nil,
		Error: &ApiError{
			Code:    code,
			Details: msg,
			Stack:   err.Error(),
		},
	}
}

func createSuccessResponse[T any](c *gin.Context, data T, message string, pagination *Pagination, err *ApiError) ApiResponse[T] {
	resp := ApiResponse[T]{
		Success:    true,
		Message:    message,
		Data:       &data,
		Pagination: pagination,

		Meta: &MetaInfo{
			Timestamp: time.Now().Format(time.RFC3339),
			RequestID: c.GetString("RequestID"),
			Path:      c.FullPath(),
			Method:    c.Request.Method,
		},
	}
	if err != nil {
		resp.Error = err
	}
	return resp
}
