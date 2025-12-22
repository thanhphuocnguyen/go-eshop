package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
)

func authorizeMiddleware(next http.Handler, roles ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			err := dto.CreateErr(UnauthorizedCode, fmt.Errorf("unauthorized"))
			w.WriteHeader(http.StatusUnauthorized)
			jsoResp, _ := json.Marshal(err)
			w.Write(jsoResp)
			return
		}

		hasRole := false
		for _, role := range roles {
			if claims["roleCode"] == role {
				hasRole = true
				break
			}
		}
		if !hasRole {
			err := dto.CreateErr(UnauthorizedCode, fmt.Errorf("forbidden"))
			w.WriteHeader(http.StatusForbidden)
			jsoResp, _ := json.Marshal(err)
			w.Write(jsoResp)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// Setup CORS configuration
func corsMiddleware() func(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3001", "http://localhost:8080"},
		AllowedHeaders: []string{"Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		MaxAge:         300,
		// AllowAllOrigins:  sv.config.Env == "development",
	})
}

// Setup environment mode based on configuration
func (s *Server) setEnvModeMiddleware(r *chi.Mux) {
	if s.config.Env == "development" {
		r.Use(middleware.Logger)
	}
	// r.Use(middleware.CleanPath)

	// Add server state validation middleware
	r.Use(s.serverStateMiddleware)

	// Add custom panic recovery middleware with better logging
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5, "text/html", "text/css"))
}

// Middleware to validate server state before processing requests
func (s *Server) serverStateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for nil server components
		if s == nil {
			http.Error(w, "Server instance is nil", http.StatusInternalServerError)
			return
		}
		if s.repo == nil {
			http.Error(w, "Database repository is not initialized", http.StatusInternalServerError)
			return
		}
		if s.tokenAuth == nil {
			http.Error(w, "Token authentication is not initialized", http.StatusInternalServerError)
			return
		}
		if s.validator == nil {
			http.Error(w, "Validator is not initialized", http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}
