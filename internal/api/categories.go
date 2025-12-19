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
// @Success 200 {object} ApiResponse[[]dto.AdminCategoryDetail]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /categories [get]
func (sv *Server) getCategories(w http.ResponseWriter, r *http.Request) {
	var query models.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	params := repository.GetCategoriesParams{
		Limit:     10,
		Offset:    0,
		Published: utils.BoolPtr(true),
	}
	params.Offset = (params.Limit) * int64(query.Page-1)
	params.Limit = int64(query.PageSize)

	categories, err := sv.repo.GetCategories(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	cnt, err := sv.repo.CountCategories(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
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

	c.JSON(http.StatusOK, dto.CreateDataResp(c, categoriesResp, dto.CreatePagination(cnt, query.Page, query.PageSize), nil))
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
// @Success 200 {object} ApiResponse[dto.CategoryDetail]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /categories/slug/{slug} [get]
func (sv *Server) getCategoryBySlug(w http.ResponseWriter, r *http.Request) {
	var param models.URISlugParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	var query models.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	category, err := sv.repo.GetCategoryBySlug(c, param.Slug)

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(InvalidBodyCode, fmt.Errorf("category with Slug %s not found", param.Slug)))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
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

	products, err := sv.repo.GetProductList(c, repository.GetProductListParams{
		CategoryIds: []uuid.UUID{category.ID},
		Limit:       query.PageSize,
		Offset:      (query.PageSize) * int64(query.Page-1),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	productResponses := make([]dto.ProductSummary, len(products))
	for i, product := range products {
		productResponses[i] = dto.MapToShopProductResponse(product)
	}
	resp.Products = productResponses

	c.JSON(http.StatusOK, dto.CreateDataResp(c, resp, nil, nil))
}

// Setup category-related routes
func (sv *Server) addCategoryRoutes(r chi.Router) {
	r.Route("categories", func(r chi.Router) {
		r.Get("", sv.getCategories)
		r.Get(":slug", sv.getCategoryBySlug)
		r.Get(":slug/products", sv.getCategoryBySlug)
	})

}
