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
)

// Setup image-related routes
func (sv *Server) addImageRoutes(rg *gin.RouterGroup) {
	images := rg.Group("images", authenticateMiddleware(sv.tokenGenerator))
	{
		images.DELETE(
			"remove-external/:public_id",
			authorizeMiddleware("admin"),
			sv.RemoveImageByPublicID)
		images.GET("", sv.GetProductImages)
	}
}

// Setup discount-related routes

// Setup webhook routes
func (sv *Server) addWebhookRoutes(router *gin.Engine) {
	webhooks := router.Group("/webhook/v1")
	{
		webhooks.POST("stripe", sv.stripeEvent)
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

	router.GET("verify-email", sv.VerifyEmail)
	// Setup API routes
	v1 := router.Group("/api/v1")
	{
		// Health check endpoint
		v1.GET("health", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"status ": "ok"})
		})

		v1.GET("homepage", sv.getHomePage)

		// Register API route groups
		sv.addAuthRoutes(v1)
		sv.addAdminRoutes(v1)
		sv.addUserRoutes(v1)
		sv.addProductRoutes(v1)
		sv.addImageRoutes(v1)
		sv.addCartRoutes(v1)
		sv.router.POST("checkout", sv.checkout)
		sv.addOrderRoutes(v1)
		sv.addPaymentRoutes(v1)
		sv.addCategoryRoutes(v1)
		sv.addCollectionRoutes(v1)
		sv.addBrandRoutes(v1)
		sv.addRatingRoutes(v1)
		sv.addDiscountRoutes(v1)
	}

	// Setup webhook routes
	sv.addWebhookRoutes(router)

	// Setup Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	sv.router = router
}
