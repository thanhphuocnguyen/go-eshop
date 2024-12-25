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
	"github.com/thanhphuocnguyen/go-eshop/internal/upload"
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
	uploadService   upload.UploadService
}

const (
	imageAssetsDir = "assets/images/"
)

func NewAPI(cfg config.Config, postgres postgres.Store, taskDistributor worker.TaskDistributor, uploadService upload.UploadService) (*Server, error) {
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
			//TODO: user.POST("refresh-token", sv.refreshToken)
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
			product.GET(":id", sv.getProduct)
			product.GET("", sv.listProducts)

			productAuthRoutes := product.Group("").Use(authMiddleware(sv.tokenGenerator), roleMiddleware(sqlc.UserRoleAdmin))
			productAuthRoutes.POST("", sv.createProduct)
			productAuthRoutes.PUT(":id", sv.updateProduct)
			productAuthRoutes.DELETE(":id", sv.removeProduct)
			productAuthRoutes.POST(":id/upload-image", sv.uploadProductImage)
			productAuthRoutes.DELETE(":id/remove-image", sv.removeProductImage)
		}

		cart := v1.Group("/cart")
		{
			cart.Use(authMiddleware(sv.tokenGenerator))
			cart.POST("", sv.createCart)
			cart.GET("", sv.getCart)
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
			order.PUT(":id/cancel", sv.cancelOrder)
			adminOrder := order.Use(roleMiddleware(sqlc.UserRoleAdmin))
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

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func responseMapper(data interface{}, message *string, err *error) gin.H {
	response := gin.H{"data": data}
	if message != nil {
		response["message"] = *message
	}
	if err != nil {
		response["error"] = (*err).Error()
	}
	return response
}
