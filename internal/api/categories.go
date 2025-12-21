package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

// getCategories retrieves a list of Categories.
// @Summary Get a list of Categories
// @Description Get a list of Categories
// @ID get-Categories
// @Accept json
// @Tags Categories
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} dto.ApiResponse[[]dto.AdminCategoryDetail]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /categories [get]
func (s *Server) getCategories(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	var query models.PaginationQuery = GetPaginationQuery(r)

	dbParams := repository.GetCategoriesParams{
		Limit:     10,
		Offset:    0,
		Published: utils.BoolPtr(true),
	}
	dbParams.Offset = (dbParams.Limit) * int64(query.Page-1)
	dbParams.Limit = int64(query.PageSize)

	categories, err := s.repo.GetCategories(c, dbParams)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	cnt, err := s.repo.CountCategories(c)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	categoriesResp := make([]dto.AdminCategoryDetail, len(categories))
	catIds := make([]uuid.UUID, len(categories))
	for i, category := range categories {
		categoriesResp[i] = dto.AdminCategoryDetail{
			ID:          category.ID.String(),
			Name:        category.Name,
			Slug:        category.Slug,
			Published:   category.Published,
			CreatedAt:   category.CreatedAt.String(),
			UpdatedAt:   category.UpdatedAt.String(),
			Description: category.Description,
			ImageUrl:    category.ImageUrl,
		}

		catIds[i] = category.ID
	}

	RespondSuccess(w, r, dto.CreateDataResp(c, categoriesResp, dto.CreatePagination(cnt, query.Page, query.PageSize), nil))
}

// getCategoryBySlug retrieves a list of Products by Category Slug.
// @Summary Get a list of Products by Category Slug
// @Description Get a list of Products by Category Slug
// @ID get-Products-by-Category-Slug
// @Accept json
// @Tags Categories
// @Produce json
// @Param slug path string true "Category Slug"
// @Param pageSize query int false "Page size"
// @Success 200 {object} dto.ApiResponse[dto.CategoryDetail]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /categories/slug/{slug} [get]
func (s *Server) getCategoryBySlug(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	slug, err := GetUrlParam(r, "slug")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	var query models.PaginationQuery = GetPaginationQuery(r)

	category, err := s.repo.GetCategoryBySlug(c, slug)

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, fmt.Errorf("category with Slug %s not found", slug))
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	resp := dto.CategoryDetail{
		ID:          category.ID.String(),
		Name:        category.Name,
		Slug:        category.Slug,
		Published:   category.Published,
		CreatedAt:   category.CreatedAt.String(),
		Description: category.Description,
		ImageUrl:    category.ImageUrl,
	}

	products, err := s.repo.GetProductList(c, repository.GetProductListParams{
		CategoryIds: []uuid.UUID{category.ID},
		Limit:       query.PageSize,
		Offset:      (query.PageSize) * int64(query.Page-1),
	})
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	productResponses := make([]dto.ProductSummary, len(products))
	for i, product := range products {
		productResponses[i] = dto.MapToShopProductResponse(product)
	}
	resp.Products = productResponses

	RespondSuccess(w, r, resp)
}

// Setup category-related routes
func (s *Server) addCategoryRoutes(r chi.Router) {
	r.Route("/categories", func(r chi.Router) {
		r.Get("/", s.getCategories)
		r.Get("/:slug", s.getCategoryBySlug)
		r.Get("/:slug/products", s.getCategoryBySlug)
	})

}
