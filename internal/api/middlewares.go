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

// Setup environment mode based on configuration
func (s *Server) registerMiddlewares(r *chi.Mux) {
	// Add server state validation middleware
	r.Use(s.serverStateMiddleware)
	// Add custom panic recovery middleware with better logging
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5, "text/html", "text/css"))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE", "HEAD", "OPTION"},
		AllowedHeaders:   []string{"User-Agent", "Content-Type", "Accept", "Accept-Encoding", "Accept-Language", "Cache-Control", "Connection", "DNT", "Host", "Origin", "Pragma", "Referer"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
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
