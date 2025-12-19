package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
)

func authorizeMiddleware(next http.Handler, roles ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authPayload, ok := r.Context().Value("auth").(*auth.TokenPayload)
		if !ok {
			err := dto.CreateErr(UnauthorizedCode, fmt.Errorf("unauthorized"))
			w.WriteHeader(http.StatusUnauthorized)
			jsoResp, _ := json.Marshal(err)
			w.Write(jsoResp)
			return
		}

		hasRole := false
		for _, role := range roles {
			if authPayload.RoleCode == role {
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
func (sv *Server) setEnvModeMiddleware(r *chi.Mux) {
	if sv.config.Env == "development" {
		r.Use(middleware.Logger)
	}
	r.Use(middleware.CleanPath)

	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5, "text/html", "text/css"))
}
