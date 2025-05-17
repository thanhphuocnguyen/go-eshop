package api

import (
	"github.com/gin-gonic/gin"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
)

func (sv *Server) setupAdminRoutes(rg *gin.RouterGroup) {
	admin := rg.Group("/admin", authMiddleware(sv.tokenGenerator), roleMiddleware(sv.repo, repository.UserRoleAdmin))
	{
		users := admin.Group("users")
		{
			users.GET("", sv.getUsersHandler)
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
			brands.POST("", sv.createBrandHandler)
			brands.GET(":id", sv.getBrandByIDHandler)
			brands.PUT(":id", sv.updateBrandHandler)
			brands.DELETE(":id", sv.deleteBrand)
		}

		collections := admin.Group("collections")
		{
			collections.GET("", sv.getCollectionsHandler)
			collections.POST("", sv.createCollectionHandler)
			collections.GET(":id", sv.getCollectionByIDHandler)
			collections.PUT(":id", sv.updateCollectionHandler)
			collections.DELETE(":id", sv.deleteCollection)
		}

		images := admin.Group("images")
		{
			productImages := images.Group("products")
			productImages.POST(":entity_id", sv.uploadProductImagesHandler)
			productImages.DELETE(":entity_id", sv.removeImageHandler)
		}

		ratings := admin.Group("ratings")
		{
			ratings.GET("", sv.getRatingsHandler)
			ratings.DELETE(":id", sv.deleteRatingHandler)
			ratings.PUT(":id/approve", sv.approveRatingHandler)
			ratings.PUT(":id/ban", sv.banUserRatingHandler)
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
	user := rg.Group("user", authMiddleware(sv.tokenGenerator))
	{
		user.GET("me", sv.getUserHandler)
		user.PATCH("me", sv.updateUserHandler)
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
		product.GET(":id/ratings", sv.getRatingsByProductHandler)
	}
}

// Setup image-related routes
func (sv *Server) setupImageRoutes(rg *gin.RouterGroup) {
	images := rg.Group("images", authMiddleware(sv.tokenGenerator))
	{
		images.DELETE(
			"remove-external/:public_id",
			roleMiddleware(sv.repo, repository.UserRoleAdmin),
			sv.removeImageByPublicIDHandler)
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
		category.GET(":slug", sv.getCategoryBySlugHandler)
		category.GET(":slug/products", sv.getCategoryBySlugHandler)
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
	ratings := rg.Group("ratings", authMiddleware(sv.tokenGenerator))
	{
		ratings.POST("", sv.postRatingHandler)
		ratings.GET(":order_id", sv.getOrderRatingsHandler)
		ratings.POST(":id/helpful", sv.postRatingHelpfulHandler)
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
