package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

// --- Public API ---

// @Summary Get a list of brands for the shop
// @Description Get a list of brands for the shop
// @ID get-shop-brands
// @Accept json
// @Tags Brands
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} ApiResponse[[]CategoryDto]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /shop/brands [get]
func (s *Server) getShopBrands(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	var queries models.PaginationQuery = GetPaginationQuery(r)

	var dbQueries repository.GetBrandsParams = repository.GetBrandsParams{
		Limit:     20,
		Offset:    0,
		Published: utils.BoolPtr(true),
	}

	dbQueries.Limit = queries.PageSize
	dbQueries.Offset = (queries.Page - 1) * queries.PageSize

	rows, err := s.repo.GetBrands(c, dbQueries)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	cnt, err := s.repo.CountBrands(c)

	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	data := make([]dto.CategoryDetail, len(rows))

	for i, row := range rows {
		model := dto.CategoryDetail{
			ID:          row.ID.String(),
			Name:        row.Name,
			Description: row.Description,
			Slug:        row.Slug,
			Published:   row.Published,
			CreatedAt:   row.CreatedAt.String(),
			ImageUrl:    row.ImageUrl,
		}

		data[i] = model
	}

	resp := dto.CreateDataResp(c, data, dto.CreatePagination(queries.Page, queries.PageSize, cnt), nil)

	RespondSuccess(w, r, resp)
}

// @Summary Get a list of brands for the shop
// @Description Get a list of brands for the shop
// @ID get-shop-brand-by-slug
// @Accept json
// @Tags Brands
// @Produce json
// @Param slug path string true "Brand slug"
// @Success 200 {object} ApiResponse[[]CategoryDto]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /shop/brands/{slug} [get]
func (s *Server) getShopBrandBySlug(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	slug, err := GetUrlParam(r, "slug")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	brandRow, err := s.repo.GetBrandBySlug(c, slug)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	resp := dto.CreateDataResp(c, brandRow, nil, nil)
	RespondSuccess(w, r, resp)
}

// Setup brand-related routes
func (s *Server) addBrandRoutes(r chi.Router) {
	r.Route("brands", func(r chi.Router) {
		r.Get("", s.getShopBrands)
		r.Get(":slug", s.getShopBrandBySlug)
	})
}
