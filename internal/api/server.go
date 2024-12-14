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
}

const (
	imageAssetsDir = "assets/images/"
)

func NewAPI(cfg config.Config, postgres postgres.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenGenerator := auth.NewPasetoTokenGenerator()
	server := &Server{
		tokenGenerator:  tokenGenerator,
		postgres:        postgres,
		config:          cfg,
		taskDistributor: taskDistributor,
	}
	server.initializeRouter()
	return server, nil
}

func (sv *Server) initializeRouter() {
	router := gin.Default()
	router.Use(gin.Recovery())
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
		user := v1.Group("/users")
		{
			user.POST("/register", sv.createUser)
			user.POST("/login", sv.loginUser)
			user.Use(authMiddleware(sv.tokenGenerator)).PATCH("", sv.updateUser)
		}

		product := v1.Group("/products")
		{
			product.GET(":product_id", sv.getProduct)
			product.GET("", sv.listProducts)
			productAuthRoutes := product.Group("").Use(authMiddleware(sv.tokenGenerator), roleMiddleware(sqlc.UserRoleAdmin))
			productAuthRoutes.POST("", sv.createProduct)
			productAuthRoutes.PUT(":product_id", sv.updateProduct)
			productAuthRoutes.DELETE(":product_id", sv.removeProduct)
			productAuthRoutes.POST(":product_id/upload-image", sv.uploadProductImage)
			productAuthRoutes.DELETE(":product_id/remove-image", sv.removeProductImage)
		}

		cart := v1.Group("/carts").Use(authMiddleware(sv.tokenGenerator))
		{
			cart.POST("", sv.createCart)
			cart.GET("", sv.getCart)
			cart.POST("/add-product", sv.addProductToCart)
			cart.POST("/remove", sv.removeProductFromCart)
			cart.POST("/checkout", sv.checkout)
			cart.PUT("/item-quantity", sv.updateCartItemQuantity)
			cart.PUT("/clear", sv.clearCart)
		}

		order := v1.Group("/orders").Use(authMiddleware(sv.tokenGenerator))
		{
			order.GET("", sv.orderList)
			order.GET(":id", sv.orderDetail)
			order.PUT(":id/cancel", sv.cancelOrder)
			adminOrder := order.Use(roleMiddleware(sqlc.UserRoleAdmin))
			adminOrder.PUT(":id/change-status", sv.changeOrderStatus)
			adminOrder.PUT(":id/change-payment-status", sv.changeOrderPaymentStatus)
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
