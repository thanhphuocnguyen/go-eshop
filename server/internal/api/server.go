package api

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	tokenGenerator := auth.NewPasetoTokenGenerator()
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
	if sv.config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
		router.Use(gin.Recovery())
	}

	cors := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
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

		v1.GET("verify-email", sv.verifyEmail)
		user := v1.Group("/user")
		{
			user.POST("register", sv.createUser)
			user.POST("login", sv.loginUser)
			user.POST("refresh-token", sv.refreshToken)
			userAuthRoutes := user.Group("").Use(authMiddleware(sv.tokenGenerator))
			userAuthRoutes.GET("", sv.getUser)
			userAuthRoutes.POST("send-verify-email", sv.sendVerifyEmail)
			userAuthRoutes.PATCH("", sv.updateUser)
			userModeratorRoutes := user.Group("").Use(
				authMiddleware(sv.tokenGenerator),
				roleMiddleware(
					sv.repo,
					repository.UserRoleAdmin,
					repository.UserRoleModerator),
			)
			userModeratorRoutes.GET("list", sv.listUsers)
		}

		userAddress := v1.Group("/address").Use(authMiddleware(sv.tokenGenerator))
		{
			userAddress.POST("", sv.createAddress)
			userAddress.GET("list", sv.listAddresses)
			userAddress.PATCH(":id", sv.updateAddress)
			userAddress.DELETE(":id", sv.removeAddress)
		}

		product := v1.Group("product")
		{
			product.GET("list", sv.getProducts)
			product.GET(":id", sv.getProductDetail)

			productAdmin := product.Group("").
				Use(authMiddleware(sv.tokenGenerator), roleMiddleware(sv.repo, repository.UserRoleAdmin))
			{
				productAdmin.POST("", sv.createProduct)
				productAdmin.PUT(":id", sv.updateProduct)
				productAdmin.DELETE(":id", sv.removeProduct)
			}

			variant := product.Group(":id/variant").
				Use(authMiddleware(sv.tokenGenerator), roleMiddleware(sv.repo, repository.UserRoleAdmin, repository.UserRoleModerator))
			{
				variant.POST("", sv.createVariant)
				variant.GET("list", sv.getVariants)
				variant.GET(":variant_id", sv.getVariant)
				variant.PUT(":variant_id", sv.updateVariant)
				variant.DELETE(":variant_id", sv.deleteVariant)
			}
		}

		images := v1.Group("images")
		{
			// product images
			productImage := images.Group("product").
				Use(authMiddleware(sv.tokenGenerator), roleMiddleware(sv.repo, repository.UserRoleAdmin, repository.UserRoleModerator))
			productImage.POST("", sv.uploadProductImage)
			productImage.GET(":product_id", sv.getProductImage)
			productImage.PUT(":product_id", sv.uploadProductImage)
			productImage.DELETE(":product_id/remove", sv.removeProductImage)

			// TODO: implement variant images
			// variant images
			// variantImage := images.Group("variant").
			// 	Use(authMiddleware(sv.tokenGenerator), roleMiddleware(sv.repo, repository.UserRoleAdmin, repository.UserRoleModerator))
			// variantImage.POST("", sv.uploadVariantImage)
			// variantImage.GET(":variant_id", sv.getVariantImages)
			// variantImage.PUT(":variant_id", sv.uploadVariantImage)
			// variantImage.DELETE(":id/image/remove", sv.removeVariantImage)
		}

		attribute := v1.Group("attribute").
			Use(authMiddleware(sv.tokenGenerator), roleMiddleware(sv.repo, repository.UserRoleAdmin))
		{
			attribute.POST("", sv.createAttribute)
			attribute.GET("list", sv.getAttributes)
			attribute.GET(":id", sv.getAttributeByID)
			attribute.PUT(":id", sv.updateAttribute)
			attribute.DELETE(":id", sv.deleteAttribute)
		}

		cart := v1.Group("/cart")
		{
			cart.Use(authMiddleware(sv.tokenGenerator))
			cart.POST("", sv.createCart)
			cart.GET("", sv.getCartDetail)
			cart.POST("checkout", sv.checkout)
			cart.PUT("clear", sv.clearCart)

			cartItem := cart.Group("/item")
			cartItem.POST("", sv.addCartItem)
			cartItem.DELETE(":id", sv.removeCartItem)
			cartItem.PUT(":id/quantity", sv.updateCartItemQuantity)
		}

		order := v1.Group("/order").Use(authMiddleware(sv.tokenGenerator))
		{
			order.GET("list", sv.orderList)
			order.GET(":id", sv.orderDetail)
			order.PUT(":id/cancel", sv.cancelOrder)

			adminOrder := order.Use(roleMiddleware(sv.repo, repository.UserRoleAdmin, repository.UserRoleModerator))
			adminOrder.PUT(":id/refund", sv.refundOrder)
			adminOrder.PUT(":id/status", sv.changeOrderStatus)
		}

		payment := v1.Group("/payment").Use(authMiddleware(sv.tokenGenerator))
		{
			payment.GET("stripe-config", sv.getStripeConfig)
			payment.POST("initiate", sv.initiatePayment)
			payment.GET(":order_id", sv.getPayment)
			payment.PUT(":order_id", sv.changePaymentStatus)
		}

		collection := v1.Group("/collection")
		{
			collection.GET("list", sv.getCollections)
			collection.GET(":id", sv.getCollectionByID)

			collectionAuthRoutes := collection.Group("").Use(
				authMiddleware(sv.tokenGenerator),
				roleMiddleware(sv.repo, repository.UserRoleAdmin),
			)
			collectionAuthRoutes.POST("", sv.createCollection)
			collectionAuthRoutes.PUT(":id", sv.updateCollection)
			collectionAuthRoutes.DELETE(":id", sv.removeCollection)
			collectionAuthRoutes.POST(":id/product", sv.addProductToCollection)
			collectionAuthRoutes.PUT(":id/product/sort_order", sv.updateProductSortOrder)
			collectionAuthRoutes.DELETE(":id/product/:product_id", sv.deleteProductFromCollection)
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

type errorResponse struct {
	Error string `json:"error"`
}

func mapErrResp(err error) errorResponse {
	return errorResponse{Error: err.Error()}
}

type GenericResponse[T any] struct {
	Data    *T      `json:"data,omitempty"`
	Message *string `json:"message,omitempty"`
	Error   *string `json:"error,omitempty"`
}

type GenericListResponse[T any] struct {
	Data    *[]T    `json:"data,omitempty"`
	Total   int64   `json:"total,omitempty"`
	Message *string `json:"message,omitempty"`
	Error   *string `json:"error,omitempty"`
}

type QueryParams struct {
	Page     int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=100"`
}
