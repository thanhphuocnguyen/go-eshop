package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Setup admin-related routes
func (s *Server) addAdminRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(func(h http.Handler) http.Handler {
			return authorizeMiddleware(h, "admin")
		})
		r.Route("/admin", func(r chi.Router) {
			// Apply authentication and authorization middleware
			// User routes
			r.Route("/users", func(r chi.Router) {
				r.Get("/", s.adminGetUsers)
				r.Get("/{id}", s.adminGetUser)
			})

			// Product routes
			r.Route("/products", func(r chi.Router) {
				r.Get("/", s.adminGetProducts)
				r.Post("/", s.createProduct)

				r.Route("/{id}", func(r chi.Router) {
					r.Put("/", s.updateProduct)
					r.Delete("/", s.deleteProduct)
					r.Post("/images", s.uploadProductImage)

					r.Route("/variants", func(r chi.Router) {
						r.Post("/", s.createVariant)
						r.Get("/", s.getProductVariants)
						r.Get("/{variantId}", s.getVariantByProductId)
						r.Put("/{variantId}", s.updateVariant)
						r.Post("/{variantId}/images", s.uploadVariantImage)
						r.Delete("/{variantId}", s.deleteVariant)
					})
				})
			})

			// Attribute routes
			r.Route("/attributes", func(r chi.Router) {
				r.Post("/", s.createAttribute)
				r.Get("/", s.adminGetAttributes)
				r.Get("/{id}", s.adminGetAttributeByID)
				r.Put("/{id}", s.updateAttribute)
				r.Delete("/{id}", s.removeAttribute)

				r.Get("/product/{id}", s.adminGetAttributeValuesForProduct)

				r.Route("/{id}", func(r chi.Router) {
					r.Post("/create", s.adminAddAttributeValue)
					r.Put("/update/{valueId}", s.adminUpdateAttrValue)
					r.Delete("/remove/{valueId}", s.adminRemoveAttrValue)
				})
			})

			// Order routes
			r.Route("/orders", func(r chi.Router) {
				r.Get("/", s.adminGetOrders)
				r.Get("/{id}", s.adminGetOrderDetail)
				r.Put("/{id}/status", s.adminChangeOrderStatus)
				r.Post("/{id}/cancel", s.adminCancelOrder)
				r.Post("/{id}/refund", s.adminRefundOrder)
				r.Delete("/{id}", s.adminDeleteOrder)
			})

			// Category routes
			r.Route("/categories", func(r chi.Router) {
				r.Get("/", s.adminGetCategories)
				r.Get("/{id}", s.adminGetCategoryByID)
				r.Post("/", s.createCategory)
				r.Put("/{id}", s.updateCategory)
				r.Delete("/{id}", s.deleteCategory)
			})

			// Brand routes
			r.Route("/brands", func(r chi.Router) {
				r.Get("/", s.adminGetBrands)
				r.Get("/{id}", s.adminGetBrandByID)
				r.Post("/", s.createBrand)
				r.Put("/{id}", s.updateBrand)
				r.Delete("/{id}", s.deleteBrand)
			})

			// Collection routes
			r.Route("/collections", func(r chi.Router) {
				r.Get("/", s.getCollections)
				r.Get("/{id}", s.adminGetCollectionByID)
				r.Post("/", s.adminCreateCollection)
				r.Put("/{id}", s.adminUpdateCollection)
				r.Delete("/{id}", s.adminDeleteCollection)
			})

			// Rating routes
			r.Route("/ratings", func(r chi.Router) {
				r.Get("/", s.adminGetRatings)
				r.Delete("/{id}", s.adminDeleteRating)
				r.Put("/{id}/approve", s.adminApproveRating)
				r.Put("/{id}/ban", s.adminBanUserRating)
			})

			// Discount routes
			r.Route("/discounts", func(r chi.Router) {
				r.Post("/", s.createDiscount)
				r.Get("/", s.adminGetDiscounts)
				r.Get("/{id}", s.getDiscountByID)
				r.Put("/{id}", s.updateDiscount)
				r.Delete("/{id}", s.adminDeleteDiscount)

				r.Route("/{id}/rules", func(r chi.Router) {
					r.Post("/", s.createDiscountRule)
					r.Get("/", s.adminGetDiscountRules)
					r.Get("/{ruleId}", s.adminGetDiscountRuleByID)
					r.Put("/{ruleId}", s.updateDiscountRule)
					r.Delete("/{ruleId}", s.adminDeleteDiscountRule)
				})
			})
		})
	})
}
