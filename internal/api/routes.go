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
func (s *Server) addImageRoutes(r chi.Router) {
	r.Route("/images", func(r chi.Router) {
		r.Get("/", s.getProductImages)
	})
}

// Setup discount-related routes

// Setup webhook routes
func (s *Server) addWebhookRoutes(r chi.Router) {
	r.Route("/webhook/v1", func(r chi.Router) {
		r.Post("/stripe", s.sendStripeEvent)
	})
}

func (s *Server) initializeRouter() {
	router := chi.NewRouter()
	gob.Register(&stripe.PaymentIntent{})

	// Setup environment mode
	s.setEnvModeMiddleware(router)

	// Setup validator
	validate := validator.New()
	validate.RegisterValidation("uuidslice", uuidSlice)

	// Setup CORS
	router.Use(corsMiddleware())

	docs.SwaggerInfo.BasePath = "/api/v1"

	// Serve static files
	fileServer := http.FileServer(http.Dir("./assets/"))
	router.Handle("/assets/*", http.StripPrefix("/assets/", fileServer))

	router.Get("/verify-email", s.VerifyEmail)
	// Setup API routes
	router.Route("/api/v1", func(r chi.Router) {
		// Health check endpoint
		r.Get("/health", func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		})

		r.Get("/homepage", s.getHomePage)

		// Register API route groups
		s.addAuthRoutes(r)
		s.router.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(s.tokenAuth))
			r.Use(jwtauth.Authenticator(s.tokenAuth))
			s.addAdminRoutes(r)
			s.addUserRoutes(r)
			s.addOrderRoutes(r)
			s.addPaymentRoutes(r)
			s.addRatingRoutes(r)
			s.addDiscountRoutes(r)
			r.Delete("/images/remove-external/{id}", s.removeImageByPublicID)
		})
		s.addProductRoutes(r)
		s.addImageRoutes(r)
		s.addCartRoutes(r)
		s.addCategoryRoutes(r)
		s.addCollectionRoutes(r)
		s.addBrandRoutes(r)
	})

	// Setup webhook routes
	s.addWebhookRoutes(router)

	// Setup Swagger
	router.Get("/swagger/*", httpSwagger.WrapHandler)

	s.router = router
}
