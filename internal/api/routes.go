package api

import (
	"encoding/gob"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stripe/stripe-go/v81"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "github.com/thanhphuocnguyen/go-eshop/docs"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
)

func (sv *Server) setupAdminRoutes(rg *gin.RouterGroup) {
	admin := rg.Group("/admin", authenticateMiddleware(sv.tokenGenerator), authorizeMiddleware(sv.repo, repository.UserRoleAdmin))
	{
		users := admin.Group("users")
		{
			users.GET("", sv.getUsersHandler)
			users.GET(":id", sv.getUserHandler)
		}

		productGroup := admin.Group("products")
		{
			productGroup.POST("", sv.addProductHandler)
			productGroup.PUT(":id", sv.updateProductHandler)
			productGroup.DELETE(":id", sv.deleteProductHandler)
		}

		attributeGroup := admin.Group("attributes")
		{
			attributeGroup.POST("", sv.createAttributeHandler)
			attributeGroup.GET("", sv.getAttributesHandler)
			attributeGroup.GET(":id", sv.getAttributeByIDHandler)
			attributeGroup.PUT(":id", sv.updateAttributeHandler)
			attributeGroup.DELETE(":id", sv.deleteAttributeHandler)
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
			categories.POST("", sv.createCategoryHandler)
			categories.PUT(":id", sv.updateCategoryHandler)
			categories.DELETE(":id", sv.deleteCategoryHandler)
		}

		brands := admin.Group("brands")
		{

			brands.GET("", sv.getBrandsHandler)
			brands.GET(":id", sv.getBrandByIDHandler)
			brands.POST("", sv.createBrandHandler)
			brands.PUT(":id", sv.updateBrandHandler)
			brands.DELETE(":id", sv.deleteBrand)
		}

		collections := admin.Group("collections")
		{
			collections.GET("", sv.getCollectionsHandler)
			collections.POST("", sv.createCollectionHandler)
			collections.GET(":id", sv.getCollectionByIDHandler)
			collections.PUT(":id", sv.updateCollectionHandler)
			collections.DELETE(":id", sv.deleteCollectionHandler)
		}

		images := admin.Group("images")
		{
			productImages := images.Group("products")
			productImages.POST(":id", sv.uploadProductImagesHandler)
			productImages.DELETE(":id", sv.removeImageHandler)
		}

		ratings := admin.Group("ratings")
		{
			ratings.GET("", sv.getRatingsHandler)
			ratings.DELETE(":id", sv.deleteRatingHandler)
			ratings.PUT(":id/approve", sv.approveRatingHandler)
			ratings.PUT(":id/ban", sv.banUserRatingHandler)
		}

		discounts := admin.Group("discounts")
		{
			discounts.POST("", sv.createDiscountHandler)
			discounts.GET("", sv.getDiscountsHandler)
			discounts.GET(":id", sv.getDiscountByIDHandler)
			discounts.PUT(":id", sv.updateDiscountHandler)
			discounts.DELETE(":id", sv.deleteDiscountHandler)
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
	users := rg.Group("user", authenticateMiddleware(sv.tokenGenerator))
	{
		users.GET("me", sv.getCurrentUserHandler)
		users.PATCH("me", sv.updateUserHandler)
		users.POST("send-verify-email", sv.sendVerifyEmailHandler)
		userAddresses := users.Group("addresses")
		{
			userAddresses.POST("", sv.createAddressHandler)
			userAddresses.PATCH(":id/default", sv.setDefaultAddressHandler)
			userAddresses.GET("", sv.getAddressesHandlers)
			userAddresses.PATCH(":id", sv.updateAddressHandlers)
			userAddresses.DELETE(":id", sv.removeAddressHandlers)
		}
	}
}

// Setup product-related routes
func (sv *Server) setupProductRoutes(rg *gin.RouterGroup) {
	products := rg.Group("products")
	{
		products.GET("", sv.getProductsHandler)
		products.GET(":id", sv.getProductDetailHandler)
		products.GET(":id/ratings", sv.getRatingsByProductHandler)
	}
}

// Setup image-related routes
func (sv *Server) setupImageRoutes(rg *gin.RouterGroup) {
	images := rg.Group("images", authenticateMiddleware(sv.tokenGenerator))
	{
		images.DELETE(
			"remove-external/:public_id",
			authorizeMiddleware(sv.repo, repository.UserRoleAdmin),
			sv.removeImageByPublicIDHandler)
		images.GET("", sv.getProductImagesHandler)
	}
}

// Setup cart-related routes
func (sv *Server) setupCartRoutes(rg *gin.RouterGroup) {
	cart := rg.Group("/cart", authenticateMiddleware(sv.tokenGenerator))
	{
		cart.POST("", sv.createCart)
		cart.GET("", sv.getCartHandler)
		cart.POST("checkout", sv.checkoutHandler)
		cart.PUT("clear", sv.clearCart)
		cartItems := cart.Group("item")
		cartItems.DELETE(":id", sv.removeCartItem)
		cartItems.PUT(":id/quantity", sv.updateCartItemQtyHandler)
	}
}

// Setup order-related routes
func (sv *Server) setupOrderRoutes(rg *gin.RouterGroup) {
	orders := rg.Group("/orders", authenticateMiddleware(sv.tokenGenerator))
	{
		orders.GET("", sv.getOrdersHandler)
		orders.GET(":id", sv.getOrderDetailHandler)
		orders.PUT(":id/confirm-received", sv.confirmOrderPayment)
		orders.PUT(":id/cancel", sv.cancelOrder)
	}
}

// Setup payment-related routes
func (sv *Server) setupPaymentRoutes(rg *gin.RouterGroup) {
	payments := rg.Group("/payments").Use(authenticateMiddleware(sv.tokenGenerator))
	{
		payments.GET(":id", sv.getPaymentHandler)
		payments.GET("stripe-config", sv.getStripeConfig)
		payments.POST("", sv.createPaymentIntentHandler)
		payments.PUT(":order_id", sv.changePaymentStatusHandler)
	}
}

// Setup category-related routes
func (sv *Server) setupCategoryRoutes(rg *gin.RouterGroup) {
	categories := rg.Group("categories")
	{
		categories.GET("", sv.getCategoriesHandler)
		categories.GET(":slug", sv.getCategoryBySlugHandler)
		categories.GET(":slug/products", sv.getCategoryBySlugHandler)
	}
}

// Setup collection-related routes
func (sv *Server) setupCollectionRoutes(rg *gin.RouterGroup) {
	collections := rg.Group("collections")
	{
		collections.GET("", sv.getCollectionsHandler)
		collections.GET(":slug", sv.getCollectionBySlugHandler)
	}
}

// Setup brand-related routes
func (sv *Server) setupBrandRoutes(rg *gin.RouterGroup) {
	brands := rg.Group("brands")
	{
		brands.GET("", sv.getShopBrandsHandler)
		brands.GET(":slug", sv.getShopBrandBySlugHandler)
	}
}

// Setup brand-related routes
func (sv *Server) setupRatingRoutes(rg *gin.RouterGroup) {
	ratings := rg.Group("ratings", authenticateMiddleware(sv.tokenGenerator))
	{
		ratings.POST("", sv.postRatingHandler)
		ratings.GET(":order_id", sv.getOrderRatingsHandler)
		ratings.POST(":id/helpful", sv.postRatingHelpfulHandler)
		ratings.POST(":id/reply", sv.postReplyRatingHandler)
	}
}

// Setup discount-related routes

// Setup webhook routes
func (sv *Server) setupWebhookRoutes(router *gin.Engine) {
	webhooks := router.Group("/webhook/v1")
	{
		webhooks.POST("stripe", sv.stripeEventHandler)
	}
}

func (sv *Server) initializeRouter() {
	router := gin.Default()
	gob.Register(&stripe.PaymentIntent{})

	// Setup environment mode
	sv.setEnvModeMiddleware(router)

	// Load HTML templates
	router.LoadHTMLGlob("static/templates/*")

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("uuidslice", uuidSlice)
	}

	// Setup CORS
	router.Use(corsMiddleware())

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

		v1.GET("homepage", sv.getHomePageHandler)
		v1.GET("verify-email", sv.verifyEmailHandler)

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
