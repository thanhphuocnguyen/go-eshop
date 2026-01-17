package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
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
			r.Get("/indexing-products", s.getIndexingProducts)
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

func (s *Server) getIndexingProducts(w http.ResponseWriter, r *http.Request) {
	pageQ := r.URL.Query().Get("page")
	perPageQ := r.URL.Query().Get("per_page")
	var page int = 1
	var perPage int = 100
	if pageQ != "" {
		page, _ = strconv.Atoi(pageQ)
	}
	if perPageQ != "" {
		perPage, _ = strconv.Atoi(perPageQ)
	}
	products, err := s.repo.GetProductsForIndexing(r.Context(), repository.GetProductsForIndexingParams{
		Limit:  int64(perPage),
		Offset: int64((page - 1) * perPage),
	})
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondJSON(w, http.StatusOK, products)
}
