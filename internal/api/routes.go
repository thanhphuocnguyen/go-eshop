package api

import (
	"encoding/gob"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/stripe/stripe-go/v84"
	httpSwagger "github.com/swaggo/http-swagger"
	docs "github.com/thanhphuocnguyen/go-eshop/docs"
)

// Setup image-related routes
func (sv *Server) addImageRoutes(r chi.Router) {
	r.Route("/images", func(r chi.Router) {
		r.Delete("/remove-external/{public_id}", sv.removeImageByPublicID)
		r.Get("/", sv.getProductImages)
	})
}

// Setup discount-related routes

// Setup webhook routes
func (sv *Server) addWebhookRoutes(r chi.Router) {
	r.Route("/webhook/v1", func(r chi.Router) {
		r.Post("/stripe", sv.sendStripeEvent)
	})
}

func (sv *Server) initializeRouter() {
	router := chi.NewRouter()
	gob.Register(&stripe.PaymentIntent{})

	// Setup environment mode
	sv.setEnvModeMiddleware(router)

	// Setup validator
	validate := validator.New()
	validate.RegisterValidation("uuidslice", uuidSlice)

	// Setup CORS
	router.Use(corsMiddleware())

	docs.SwaggerInfo.BasePath = "/api/v1"

	// Serve static files
	fileServer := http.FileServer(http.Dir("./assets/"))
	router.Handle("/assets/*", http.StripPrefix("/assets/", fileServer))

	router.Get("/verify-email", sv.VerifyEmail)
	// Setup API routes
	router.Route("/api/v1", func(r chi.Router) {
		// Health check endpoint
		r.Get("/health", func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		})

		r.Get("/homepage", sv.getHomePage)

		// Register API route groups
		sv.addAuthRoutes(r)
		sv.router.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(sv.tokenAuth))
			r.Use(jwtauth.Authenticator(sv.tokenAuth))
			sv.addAdminRoutes(r)
			sv.addUserRoutes(r)
			sv.addOrderRoutes(r)
			sv.addPaymentRoutes(r)
			sv.addRatingRoutes(r)
		})
		sv.addProductRoutes(r)
		sv.addImageRoutes(r)
		sv.addCartRoutes(r)
		sv.addCategoryRoutes(r)
		sv.addCollectionRoutes(r)
		sv.addBrandRoutes(r)
		sv.addDiscountRoutes(r)
	})

	// Setup webhook routes
	sv.addWebhookRoutes(router)

	// Setup Swagger
	router.Get("/swagger/*", httpSwagger.WrapHandler)

	sv.router = router
}
