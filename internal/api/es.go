package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Elastic search re-index api util
func (s *Server) addEsRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(func(h http.Handler) http.Handler {
			return authorizeMiddleware(h, "admin")
		})
		r.Route("/es", func(r chi.Router) {
			r.Post("/product-mappings", s.createProductIndex)
			r.Get("/product-mappings", s.getProductMappings)
		})
	})
}

func (s *Server) createProductIndex(w http.ResponseWriter, r *http.Request) {
	err := s.elasticClient.CreateIndex()
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondNoContent(w)
}

func (s *Server) getProductMappings(w http.ResponseWriter, r *http.Request) {
	mappings, err := s.elasticClient.GetMapping()
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondJSON(w, http.StatusOK, mappings)
}
