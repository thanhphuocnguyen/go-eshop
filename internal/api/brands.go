package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
func (sv *Server) getShopBrands(c *gin.Context) {
	var queries models.PaginationQuery
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	var dbQueries repository.GetBrandsParams = repository.GetBrandsParams{
		Limit:     20,
		Offset:    0,
		Published: utils.BoolPtr(true),
	}

	dbQueries.Limit = queries.PageSize
	dbQueries.Offset = (queries.Page - 1) * queries.PageSize

	rows, err := sv.repo.GetBrands(c, dbQueries)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	cnt, err := sv.repo.CountBrands(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
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

	c.JSON(http.StatusOK, resp)
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
func (sv *Server) getShopBrandBySlug(c *gin.Context) {
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

	var queries models.PaginationQuery
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	var dbQueries repository.GetBrandsParams = repository.GetBrandsParams{
		Limit:  20,
		Offset: 0,
	}
	dbQueries.Limit = queries.PageSize
	dbQueries.Offset = (queries.Page - 1) * queries.PageSize

	rows, err := sv.repo.GetBrands(c, dbQueries)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	cnt, err := sv.repo.CountBrands(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	data := make([]dto.CategoryDetail, len(rows))

	for i, row := range rows {
		data[i] = dto.CategoryDetail{
			ID:          row.ID.String(),
			Name:        row.Name,
			Description: row.Description,
			Slug:        row.Slug,
			ImageUrl:    row.ImageUrl,
		}
	}

	resp := dto.CreateDataResp(c, data, dto.CreatePagination(queries.Page, queries.PageSize, cnt), nil)
	c.JSON(http.StatusOK, resp)
}

// Setup brand-related routes
func (sv *Server) addBrandRoutes(r chi.Router) {
	brands := rg.Group("brands")
	{
		brands.GET("", sv.getShopBrands)
		brands.GET(":slug", sv.getShopBrandBySlug)
	}
}
