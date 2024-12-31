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
	"github.com/thanhphuocnguyen/go-eshop/internal/db/postgres"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/sqlc"
	"github.com/thanhphuocnguyen/go-eshop/internal/uploadsrv"
	"github.com/thanhphuocnguyen/go-eshop/internal/worker"
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
	postgres        postgres.Store
	router          *gin.Engine
	taskDistributor worker.TaskDistributor
	tokenGenerator  auth.TokenGenerator
	uploadService   uploadsrv.UploadService
}

func NewAPI(cfg config.Config, postgres postgres.Store, taskDistributor worker.TaskDistributor, uploadService uploadsrv.UploadService) (*Server, error) {
	tokenGenerator := auth.NewPasetoTokenGenerator()
	server := &Server{
		tokenGenerator:  tokenGenerator,
		postgres:        postgres,
		config:          cfg,
		taskDistributor: taskDistributor,
		uploadService:   uploadService,
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
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	})
	router.Use(cors)

	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	docs.SwaggerInfo.BasePath = "/api/v1"

	router.Static("/assets", "./assets")

	v1 := router.Group("/api/v1")
	{
		user := v1.Group("/user")
		{
			user.POST("register", sv.createUser)
			user.POST("login", sv.loginUser)
			user.POST("refresh-token", sv.refreshToken)
			userAuthRoutes := user.Group("").Use(authMiddleware(sv.tokenGenerator))
			userAuthRoutes.GET("", sv.getUser)
			userAuthRoutes.PATCH("", sv.updateUser)
		}

		userAddress := v1.Group("/address").Use(authMiddleware(sv.tokenGenerator))
		{
			userAddress.POST("", sv.createAddress)
			userAddress.GET("", sv.listAddresses)
			userAddress.PATCH(":id", sv.updateAddress)
			userAddress.DELETE(":id", sv.removeAddress)
		}

		product := v1.Group("/product")
		{
			product.GET(":id", sv.getProductDetail)
			product.GET("", sv.getProducts)

			productAuthRoutes := product.Group("").Use(authMiddleware(sv.tokenGenerator), roleMiddleware(sqlc.UserRoleAdmin))
			productAuthRoutes.POST("", sv.createProduct)
			productAuthRoutes.PUT(":id", sv.updateProduct)
			productAuthRoutes.DELETE(":id", sv.removeProduct)
			productAuthRoutes.POST(":id/image", sv.uploadProductImage)
			productAuthRoutes.GET(":id/image", sv.getProductImages)
			productAuthRoutes.PUT(":id/image/:image_id", sv.uploadProductImage)
			productAuthRoutes.PUT(":id/image/:image_id/primary", sv.setImagesPrimary)
			productAuthRoutes.DELETE(":id/image/:image_id", sv.removeProductImage)
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
			order.GET("", sv.orderList)
			order.GET(":id", sv.orderDetail)
			adminOrder := order.Use(roleMiddleware(sqlc.UserRoleAdmin))
			adminOrder.PUT(":id/cancel", sv.cancelOrder)
			adminOrder.PUT(":id/change-status", sv.changeOrderStatus)
			adminOrder.PUT(":id/change-payment-status", sv.changeOrderPaymentStatus)
		}

		collection := v1.Group("/collection")
		{
			collection.GET("", sv.listCollections)
			collectionAuthRoutes := collection.Group("").Use(authMiddleware(sv.tokenGenerator), roleMiddleware(sqlc.UserRoleAdmin))
			collectionAuthRoutes.POST("", sv.createCollection)
			collectionAuthRoutes.GET(":id", sv.getCollectionByID)
			collectionAuthRoutes.PUT(":id", sv.updateCollection)
			collectionAuthRoutes.DELETE(":id", sv.removeCollection)
			collectionAuthRoutes.POST(":id/product", sv.addProductToCollection)
			collectionAuthRoutes.DELETE(":id/product", sv.deleteProductFromCollection)
		}
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
	Total   *int64  `json:"total,omitempty"`
	Message *string `json:"message,omitempty"`
	Error   *string `json:"error,omitempty"`
}
