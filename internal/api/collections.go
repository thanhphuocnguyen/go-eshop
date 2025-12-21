package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
)

// --- Public API ---

// @Summary Get a list of Collections
// @Description Get a list of Collections
// @ID get-Shop-Collection-by-slug
// @Accept json
// @Tags Collections
// @Produce json
// @Param slug path string true "Collection slug"
// @Success 200 {object} dto.ApiResponse[dto.CategoryDetail]
// @Failure 400 {object} dto.ErrorResp
// @Failure 404 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /collections/{slug} [get]
func (s *Server) getCollectionBySlug(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	slug, err := GetUrlParam(r, "slug")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	query := GetPaginationQuery(r)

	collection, err := s.repo.GetCollectionBySlug(c, slug)

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, err)
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	rows, err := s.repo.GetProductList(c, repository.GetProductListParams{
		CollectionIds: []uuid.UUID{collection.ID},
		Limit:         query.PageSize,
		Offset:        (query.PageSize) * int64(query.Page-1),
	})

	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	collectionResp := dto.CategoryDetail{
		ID:          collection.ID.String(),
		Name:        collection.Name,
		Description: collection.Description,
		Slug:        collection.Slug,
		ImageUrl:    collection.ImageUrl,
		CreatedAt:   collection.CreatedAt.String(),
		Products:    make([]dto.ProductSummary, len(rows)),
	}
	for i, row := range rows {
		collectionResp.Products[i] = dto.MapToShopProductResponse(row)
	}

	RespondSuccess(w, r, collectionResp)
}

// Setup collection-related routes
func (s *Server) addCollectionRoutes(r chi.Router) {
	r.Route("/collections", func(r chi.Router) {
		r.Get("/", s.adminGetCollections)
		r.Get("/{slug}", s.getCollectionBySlug)
	})
}
