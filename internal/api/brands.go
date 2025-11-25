package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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
func (sv *Server) GetShopBrandsHandler(c *gin.Context) {
	var queries models.PaginationQuery
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
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
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	cnt, err := sv.repo.CountBrands(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
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

	resp := createDataResp(c, data, createPagination(queries.Page, queries.PageSize, cnt), nil)

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
func (sv *Server) GetShopBrandBySlugHandler(c *gin.Context) {
	var param models.URISlugParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}
	var query models.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	var queries models.PaginationQuery
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
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
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	cnt, err := sv.repo.CountBrands(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
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

	resp := createDataResp(c, data, createPagination(queries.Page, queries.PageSize, cnt), nil)
	c.JSON(http.StatusOK, resp)
}

// --- Admin API ---
// @Summary Create a new Brand
// @Description Create a new Brand
// @Tags Admin
// @ID create-Brand
// @Accept json
// @Produce json
// @Param request body CreateCategoryRequest true "Brand request"
// @Success 201 {object} ApiResponse[CategoryDto]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/brands [post]
func (sv *Server) CreateBrandHandler(c *gin.Context) {
	var req models.CreateCategoryModel
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}
	params := repository.CreateBrandParams{
		Name: req.Name,
		Slug: req.Slug,
	}

	if req.Description != nil {
		params.Description = req.Description
	}
	if req.Image != nil {
		publicID, imgUrl, err := sv.uploadService.UploadFile(c, req.Image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErr(UploadFileCode, err))
			return
		}
		params.ImageUrl = &imgUrl
		params.ImageID = &publicID
	}

	col, err := sv.repo.CreateBrand(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusCreated, createDataResp(c, col, nil, nil))
}

// @Summary Get a list of brands
// @Description Get a list of brands
// @Tags Admin
// @ID get-brands
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} ApiResponse[[]CategoryDto]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /brands [get]
func (sv *Server) GetBrandsHandler(c *gin.Context) {
	var queries models.PaginationQuery
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
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
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	cnt, err := sv.repo.CountBrands(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	data := make([]dto.CategoryDetail, 0, len(rows))

	for _, row := range rows {
		data = append(data, dto.CategoryDetail{
			ID:          row.ID.String(),
			Name:        row.Name,
			Description: row.Description,
			Slug:        row.Slug,
			ImageUrl:    row.ImageUrl,
		})
	}

	pagination := createPagination(queries.Page, queries.PageSize, cnt)

	resp := createDataResp(c, data, pagination, nil)
	c.JSON(http.StatusOK, resp)
}

// @Summary Get a Brand by ID
// @Description Get a Brand by ID
// @ID get-Brand-by-id
// @Accept json
// @Tags Admin
// @Produce json
// @Param id path int true "Brand ID"
// @Success 200 {object} ApiResponse[CategoryDto]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/brands/{id} [get]
func (sv *Server) GetBrandByIDHandler(c *gin.Context) {
	var param models.UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	result, err := sv.repo.GetBrandByID(c, uuid.MustParse(param.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErr(NotFoundCode, fmt.Errorf("brand with ID %s not found", param.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	colResp := dto.AdminCategoryDetail{
		ID:          result.ID.String(),
		Name:        result.Name,
		Description: result.Description,
		Slug:        result.Slug,
		Published:   result.Published,
		CreatedAt:   result.CreatedAt.Format("2006-01-02 15:04:05"),
		ImageUrl:    result.ImageUrl,
		UpdatedAt:   result.UpdatedAt.String(),
	}

	c.JSON(http.StatusOK, createDataResp(c, colResp, nil, nil))
}

// @Summary Update a Brand
// @Description Update a Brand
// @ID update-Brand
// @Accept json
// @Produce json
// @Tags Admin
// @Param id path int true "Brand ID"
// @Param request body UpdateCategoryRequest true "Brand request"
// @Success 200 {object} ApiResponse[CategoryDto]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/brands/{id} [put]
func (sv *Server) UpdateBrandHandler(c *gin.Context) {
	var param models.UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}
	var req models.UpdateCategoryModel
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	brand, err := sv.repo.GetBrandByID(c, uuid.MustParse(param.ID))

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErr(NotFoundCode, fmt.Errorf("brand with ID %s not found", param.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	updateParam := repository.UpdateBrandWithParams{
		ID:   brand.ID,
		Name: req.Name,
	}

	if req.Image != nil {

		imgID, imgUrl, err := sv.uploadService.UploadFile(c, req.Image)
		if err != nil {
			log.Error().Err(err).Interface("value", req.Image.Header).Msg("error when upload image")
			c.JSON(http.StatusInternalServerError, createErr(UploadFileCode, err))
			return
		}
		updateParam.ImageUrl = &imgUrl
		updateParam.ImageID = &imgID
		oldImageID := brand.ImageID
		if oldImageID != nil {
			_, err := sv.uploadService.RemoveFile(c, *oldImageID)
			if err != nil {
				log.Error().Err(err).Msg("error when remove old image")
				c.JSON(http.StatusInternalServerError, createErr(UploadFileCode, err))
				return
			}
			log.Info().Msgf("old image %s removed", *oldImageID)
		}
	}

	if req.Slug != nil {
		updateParam.Slug = req.Slug
	}
	if req.Description != nil {
		updateParam.Description = req.Description
	}

	if req.Published != nil {
		updateParam.Published = req.Published
	}

	col, err := sv.repo.UpdateBrandWith(c, updateParam)

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusOK, createDataResp(c, col, nil, nil))
}

// @Summary Delete a Brand
// @Description Delete a Brand
// @ID delete-Brand
// @Accept json
// @Tags Admin
// @Produce json
// @Param id path int true "Brand ID"
// @Success 204 {object} ApiResponse[bool]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/brands/{id} [delete]
func (sv *Server) DeleteBrandHandler(c *gin.Context) {
	var colID models.UriIDParam
	if err := c.ShouldBindUri(&colID); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	_, err := sv.repo.GetBrandByID(c, uuid.MustParse(colID.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErr(NotFoundCode, fmt.Errorf("brand with ID %s not found", colID.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	err = sv.repo.DeleteBrand(c, uuid.MustParse(colID.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}
	c.JSON(http.StatusOK, createDataResp(c, true, nil, nil))
}
