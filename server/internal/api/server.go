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
	"github.com/thanhphuocnguyen/go-eshop/internal/auth"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/worker"
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
	if sv.config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
		router.Use(gin.Recovery())
	}
	if sv.config.Env == "development" {
		router.Use(gin.Logger())
	}

	cors := cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001", "http://localhost:8080"},
		AllowHeaders:     []string{"Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"},
		AllowFiles:       true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowCredentials: true,
		// AllowAllOrigins:  sv.config.Env == "development",
	})

	router.Use(cors)

	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	docs.SwaggerInfo.BasePath = "/api/v1"

	router.Static("/assets", "./assets")

	v1 := router.Group("/api/v1")
	{
		v1.GET("health", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"status ": "ok"})
		})

		v1.GET("verify-email", sv.verifyEmailHandler)
		auth := v1.Group("/auth")
		{
			auth.POST("register", sv.registerHandler)
			auth.POST("login", sv.loginHandler)
			auth.POST("refresh-token", sv.refreshTokenHandler)
			authRoutes := auth.Use(authMiddleware(sv.tokenGenerator))
			authRoutes.POST("send-verify-email", sv.sendVerifyEmailHandler)
		}

		user := v1.Group("/user", authMiddleware(sv.tokenGenerator))
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

		moderatorRoutes := v1.Group("/moderator",
			authMiddleware(sv.tokenGenerator),
			roleMiddleware(
				sv.repo,
				repository.UserRoleAdmin,
				repository.UserRoleModerator))
		{
			moderatorRoutes.GET("list", sv.listUsers)
		}

		product := v1.Group("products")
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

		images := v1.Group(
			"images",
			authMiddleware(sv.tokenGenerator),
			roleMiddleware(sv.repo, repository.UserRoleAdmin, repository.UserRoleModerator))
		{
			// images
			images.DELETE("remove-external/:public_id", sv.removeImageByPublicID)
			images.GET("", sv.getImages)
			images.POST("product/:entity_id", sv.uploadProductImages)
			images.DELETE(":entity_id", sv.removeImage)

		}

		attribute := v1.Group(
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

		cart := v1.Group("/cart", authMiddleware(sv.tokenGenerator))
		{
			cart.POST("", sv.createCart)
			cart.GET("", sv.getCartHandler)
			cart.POST("checkout", sv.checkoutHandler)
			cart.PUT("clear", sv.clearCart)

			cartItem := cart.Group("item")
			cartItem.DELETE(":id", sv.removeCartItem)
			cartItem.PUT(":id/quantity", sv.updateCartItemQtyHandler)
		}

		order := v1.Group("/order", authMiddleware(sv.tokenGenerator))
		{
			order.GET("list", sv.orderList)
			order.GET(":id", sv.orderDetail)
			order.PUT(":id/cancel", sv.cancelOrder)

			adminOrder := order.Group(
				"",
				roleMiddleware(sv.repo, repository.UserRoleAdmin, repository.UserRoleModerator))
			adminOrder.PUT(":id/refund", sv.refundOrder)
			adminOrder.PUT(":id/status", sv.changeOrderStatus)
		}

		payment := v1.Group("/payments").Use(authMiddleware(sv.tokenGenerator))
		{
			payment.GET("stripe-config", sv.getStripeConfig)
			payment.POST("", sv.createPaymentIntentHandler)
			payment.GET(":order_id", sv.getPaymentHandler)
			payment.PUT(":order_id", sv.changePaymentStatus)
		}

		category := v1.Group("categories")
		{
			category.GET("", sv.getCategories)
			//TODO: category.GET("dashboard", sv.getCategoriesForDashboard)
			category.GET(":id", sv.getCategoryByID)
			category.GET(":id/products", sv.getProductsByCategory)
			categoryAuthRoutes := category.Group("").Use(
				authMiddleware(sv.tokenGenerator),
				roleMiddleware(sv.repo, repository.UserRoleAdmin),
			)
			categoryAuthRoutes.POST("", sv.addCategoryHandler)
			categoryAuthRoutes.PUT(":id", sv.updateCategory)
			categoryAuthRoutes.DELETE(":id", sv.deleteCategory)
		}

		collection := v1.Group("collections")
		{
			// CRUD
			collection.GET("", sv.getCollections)
			collection.GET(":id", sv.getCollectionByID)
			collection.GET(":id/products", sv.getProductsByCollection)
			collectionAuthRoutes := collection.Group("").Use(
				authMiddleware(sv.tokenGenerator),
				roleMiddleware(sv.repo, repository.UserRoleAdmin),
			)
			collectionAuthRoutes.POST("", sv.createCollection)
			collectionAuthRoutes.PUT(":id", sv.updateCollection)
			collectionAuthRoutes.DELETE(":id", sv.deleteCollection)
		}

		brand := v1.Group("brands")
		{
			// CRUD
			brand.GET("", sv.getBrandsHandler)
			brand.GET(":id", sv.getBrandByIDHandler)
			brand.GET(":id/products", sv.getProductsByBrandHandler)
			brandAuthRoutes := brand.Group("").Use(
				authMiddleware(sv.tokenGenerator),
				roleMiddleware(sv.repo, repository.UserRoleAdmin),
			)
			brandAuthRoutes.POST("", sv.createBrandHandler)
			brandAuthRoutes.PUT(":id", sv.updateBrand)
			brandAuthRoutes.DELETE(":id", sv.deleteBrand)
		}
	}

	// webhooks
	webhook := router.Group("/webhook/v1")
	{
		webhook.POST("stripe", sv.stripeWebhook)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	sv.router = router
}

func (s *Server) Server(addr string) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: s.router.Handler(),
	}
}

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
	Page            int32 `json:"page"`
	PageSize        int32 `json:"pageSize"`
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
	Page     int32 `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize int32 `form:"page_size,default=20" binding:"omitempty,min=1,max=100"`
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
