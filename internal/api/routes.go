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

	// Setup global middleware
	s.setEnvModeMiddleware(router)
	router.Use(corsMiddleware())

	// Setup validator (consider moving to server initialization if used elsewhere)
	validate := validator.New()
	validate.RegisterValidation("uuidslice", uuidSlice)

	docs.SwaggerInfo.BasePath = "/api/v1"

	// Assign router to server
	s.router = router

	// Setup static file serving
	s.setupStaticRoutes()

	// Setup main routes
	s.setupMainRoutes()

	// Setup webhook routes (outside API versioning)
	s.addWebhookRoutes(s.router)

	// Setup Swagger
	s.router.Get("/swagger/*", httpSwagger.WrapHandler)
}

// setupStaticRoutes handles static file serving
func (s *Server) setupStaticRoutes() {
	fileServer := http.FileServer(http.Dir("./assets/"))
	s.router.Handle("/assets/*", http.StripPrefix("/assets/", fileServer))
	s.router.Get("/verify-email", s.VerifyEmail)
}

// setupMainRoutes organizes API routes into public and protected groups
func (s *Server) setupMainRoutes() {
	s.router.Route("/api/v1", func(r chi.Router) {
		// Health check endpoint
		r.Get("/health", s.healthCheck)

		// Public routes
		r.Get("/homepage", s.getHomePage)
		s.addAuthRoutes(r)
		s.addPublicRoutes(r)

		// Protected routes (require authentication)
		r.Group(func(protected chi.Router) {
			protected.Use(jwtauth.Verifier(s.tokenAuth))
			protected.Use(jwtauth.Authenticator(s.tokenAuth))

			s.addAdminRoutes(protected)
			s.addUserRoutes(protected)
			s.addOrderRoutes(protected)
			s.addPaymentRoutes(protected)
			s.addRatingRoutes(protected)
			s.addDiscountRoutes(protected)
			protected.Delete("/images/remove-external/{id}", s.removeImageByPublicID)
		})
	})
}

// addPublicRoutes groups all public routes that don't require authentication
func (s *Server) addPublicRoutes(r chi.Router) {
	s.addProductRoutes(r)
	s.addImageRoutes(r)
	s.addCartRoutes(r)
	s.addCategoryRoutes(r)
	s.addCollectionRoutes(r)
	s.addBrandRoutes(r)
}

// healthCheck handles the health check endpoint
func (s *Server) healthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
