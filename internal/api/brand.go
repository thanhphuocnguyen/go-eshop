package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

// ------------------------------------------ Request and Response ------------------------------------------

type BrandsQueries struct {
	PaginationQueryParams
	brands []int32 `form:"Brand_ids,omitempty"`
}

type BrandProductRequest struct {
	SortOrder int16 `json:"sort_order,omitempty"`
}

// ------------------------------------------ API Handlers ------------------------------------------

// --- Public API ---
// getShopBrandsHandler retrieves a list of brands for the shop.
// @Summary Get a list of brands for the shop
// @Description Get a list of brands for the shop
// @ID get-shop-brands
// @Accept json
// @Tags Brands
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} ApiResponse[CategoryResponse]
// @Failure 400 {object} ApiResponse[gin.H]
// @Failure 500 {object} ApiResponse[gin.H]
// @Router /shop/brands [get]
func (sv *Server) getShopBrandsHandler(c *gin.Context) {
	var queries BrandsQueries
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[gin.H](InvalidBodyCode, "", err))
		return
	}
	var dbQueries repository.GetBrandsParams = repository.GetBrandsParams{
		Limit:     20,
		Offset:    0,
		Published: utils.BoolPtr(true),
	}

	dbQueries.Limit = queries.PageSize
	dbQueries.Offset = (queries.Page - 1) * queries.PageSize
	var cached *struct {
		Data       []CategoryResponse `json:"data"`
		Pagination *Pagination        `json:"pagination"`
	}

	if err := sv.cacheService.Get(c, fmt.Sprintf("brands-%d-%d", queries.Page, queries.PageSize), &cached); err == nil {
		if cached != nil {
			resp := createSuccessResponse(c, cached.Data, "Cached", cached.Pagination, nil)
			c.JSON(http.StatusOK, resp)
			return
		}
	}

	rows, err := sv.repo.GetBrands(c, dbQueries)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "", err))
		return
	}

	cnt, err := sv.repo.CountBrands(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "", err))
		return
	}

	data := make([]CategoryResponse, len(rows))

	for i, row := range rows {
		model := CategoryResponse{
			ID:          row.ID.String(),
			Name:        row.Name,
			Description: row.Description,
			Slug:        row.Slug,
			Published:   row.Published,
			CreatedAt:   row.CreatedAt.String(),
			UpdatedAt:   row.UpdatedAt.String(),
			ImageUrl:    row.ImageUrl,
		}

		if row.Remarkable != nil {
			model.Remarkable = *row.Remarkable
		}

		data[i] = model
	}
	cached = &struct {
		Data       []CategoryResponse `json:"data"`
		Pagination *Pagination        `json:"pagination"`
	}{
		Data: data,
		Pagination: &Pagination{
			Page:            queries.Page,
			Total:           cnt,
			PageSize:        queries.PageSize,
			TotalPages:      cnt / int64(queries.PageSize),
			HasNextPage:     cnt > int64((queries.Page-1)*queries.PageSize+queries.PageSize),
			HasPreviousPage: queries.Page > 1,
		},
	}

	resp := createSuccessResponse(c, cached.Data, "", cached.Pagination, nil)
	if err = sv.cacheService.Set(c, fmt.Sprintf("brands-%d-%d", queries.Page, queries.PageSize), resp, nil); err != nil {
		log.Error().Err(err).Msg("error when set brands to cache")
	}
	c.JSON(http.StatusOK, resp)
}

// getShopBrandBySlugHandler retrieves a list of brands for the shop.
// @Summary Get a list of brands for the shop
// @Description Get a list of brands for the shop
// @ID get-shop-brand-by-slug
// @Accept json
// @Tags Brands
// @Produce json
// @Param slug path string true "Brand slug"
// @Success 200 {object} ApiResponse[CategoryResponse]
// @Failure 400 {object} ApiResponse[gin.H]
// @Failure 500 {object} ApiResponse[gin.H]
// @Router /shop/brands/{slug} [get]
func (sv *Server) getShopBrandBySlugHandler(c *gin.Context) {
	var param SlugParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[CategoryResponse](InvalidBodyCode, "", err))
		return
	}
	var query PaginationQueryParams
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[CategoryResponse](InvalidBodyCode, "", err))
		return
	}

	brand, err := sv.repo.GetBrandBySlug(c, param.Slug)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[CategoryResponse](NotFoundCode, "", fmt.Errorf("category with slug %s not found", param.Slug)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[CategoryResponse](InternalServerErrorCode, "", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(
		c,
		CategoryResponse{
			ID:          brand.ID.String(),
			Name:        brand.Name,
			Description: brand.Description,
			Slug:        brand.Slug,
			ImageUrl:    brand.ImageUrl,
			CreatedAt:   brand.CreatedAt.String(),
		},
		"",
		nil,
		nil,
	))
	var queries BrandsQueries
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[gin.H](InvalidBodyCode, "", err))
		return
	}
	var dbQueries repository.GetBrandsParams = repository.GetBrandsParams{
		Limit:  20,
		Offset: 0,
	}
	dbQueries.Limit = queries.PageSize
	dbQueries.Offset = (queries.Page - 1) * queries.PageSize
	var cached *struct {
		Data       []CategoryResponse `json:"data"`
		Pagination *Pagination        `json:"pagination"`
	}

	if err := sv.cacheService.Get(c, fmt.Sprintf("brands-%d-%d", queries.Page, queries.PageSize), &cached); err == nil {
		if cached != nil {
			resp := createSuccessResponse(c, cached.Data, "Cached", cached.Pagination, nil)
			c.JSON(http.StatusOK, resp)
			return
		}
	}

	rows, err := sv.repo.GetBrands(c, dbQueries)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "", err))
		return
	}

	cnt, err := sv.repo.CountBrands(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "", err))
		return
	}

	data := make([]CategoryResponse, len(rows))

	for i, row := range rows {
		data[i] = CategoryResponse{
			ID:          row.ID.String(),
			Name:        row.Name,
			Description: row.Description,
			Slug:        row.Slug,
			ImageUrl:    row.ImageUrl,
		}
	}
	cached = &struct {
		Data       []CategoryResponse `json:"data"`
		Pagination *Pagination        `json:"pagination"`
	}{
		Data: data,
		Pagination: &Pagination{
			Page:            queries.Page,
			Total:           cnt,
			PageSize:        queries.PageSize,
			TotalPages:      cnt / int64(queries.PageSize),
			HasNextPage:     cnt > int64((queries.Page-1)*queries.PageSize+queries.PageSize),
			HasPreviousPage: queries.Page > 1,
		},
	}

	resp := createSuccessResponse(c, cached.Data, "", cached.Pagination, nil)
	if err = sv.cacheService.Set(c, fmt.Sprintf("brands-%d-%d", queries.Page, queries.PageSize), resp, nil); err != nil {
		log.Error().Err(err).Msg("error when set brands to cache")
	}
	c.JSON(http.StatusOK, resp)
}

// --- Admin API ---
// createBrandHandler creates a new Brand.
// @Summary Create a new Brand
// @Description Create a new Brand
// @Tags Admin
// @ID create-Brand
// @Accept json
// @Produce json
// @Param request body CreateCategoryRequest true "Brand request"
// @Success 201 {object} ApiResponse[CategoryResponse]
// @Failure 400 {object} ApiResponse[gin.H]
// @Failure 500 {object} ApiResponse[gin.H]
// @Router /brands [post]
func (sv *Server) createBrandHandler(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[gin.H](InvalidBodyCode, "", err))
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
			c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](UploadFileCode, "error when upload image", err))
			return
		}
		params.ImageUrl = &imgUrl
		params.ImageID = &publicID
	}

	col, err := sv.repo.CreateBrand(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "", err))
		return
	}

	c.JSON(http.StatusCreated, createSuccessResponse(c, col, "", nil, nil))
}

// getBrandsHandler retrieves a list of brands.
// @Summary Get a list of brands
// @Description Get a list of brands
// @Tags Admin
// @ID get-brands
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} ApiResponse[gin.H]
// @Failure 400 {object} ApiResponse[gin.H]
// @Failure 500 {object} ApiResponse[CategoryResponse]
// @Router /brands [get]
func (sv *Server) getBrandsHandler(c *gin.Context) {
	var queries BrandsQueries
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[gin.H](InvalidBodyCode, "", err))
		return
	}
	var dbQueries repository.GetBrandsParams = repository.GetBrandsParams{
		Limit:  20,
		Offset: 0,
	}
	dbQueries.Limit = queries.PageSize
	dbQueries.Offset = (queries.Page - 1) * queries.PageSize

	var cached *struct {
		Data       []CategoryResponse `json:"data"`
		Pagination *Pagination        `json:"pagination"`
	}
	if err := sv.cacheService.Get(c, fmt.Sprintf("brands-%d-%d", queries.Page, queries.PageSize), &cached); err == nil {
		if cached != nil {
			resp := createSuccessResponse(c, cached.Data, "Cached", cached.Pagination, nil)
			c.JSON(http.StatusOK, resp)
			return
		}
	}

	rows, err := sv.repo.GetBrands(c, dbQueries)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "", err))
		return
	}

	cnt, err := sv.repo.CountBrands(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "", err))
		return
	}

	data := make([]CategoryResponse, 0, len(rows))

	for _, row := range rows {
		data = append(data, CategoryResponse{
			ID:          row.ID.String(),
			Name:        row.Name,
			Description: row.Description,
			Slug:        row.Slug,
			ImageUrl:    row.ImageUrl,
		})
	}

	pagination := &Pagination{
		Page:            queries.Page,
		Total:           cnt,
		PageSize:        queries.PageSize,
		TotalPages:      cnt / int64(queries.PageSize),
		HasNextPage:     cnt > int64((queries.Page-1)*queries.PageSize+queries.PageSize),
		HasPreviousPage: queries.Page > 1,
	}

	cached = &struct {
		Data       []CategoryResponse "json:\"data\""
		Pagination *Pagination        "json:\"pagination\""
	}{
		Data:       data,
		Pagination: pagination,
	}

	if err = sv.cacheService.Set(c, fmt.Sprintf("brands-%d-%d", queries.Page, queries.PageSize), cached, nil); err != nil {
		log.Error().Err(err).Msg("error when set brands to cache")
	}

	resp := createSuccessResponse(c, data, "", pagination, nil)
	c.JSON(http.StatusOK, resp)
}

// getBrandByIDHandler retrieves a Brand by its ID.
// @Summary Get a Brand by ID
// @Description Get a Brand by ID
// @ID get-Brand-by-id
// @Accept json
// @Tags Admin
// @Produce json
// @Param id path int true "Brand ID"
// @Success 200 {object} ApiResponse[CategoryResponse]
// @Failure 400 {object} ApiResponse[gin.H]
// @Failure 500 {object} ApiResponse[gin.H]
// @Router /admin/brands/{id} [get]
func (sv *Server) getBrandByIDHandler(c *gin.Context) {
	var param URIParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[gin.H](InvalidBodyCode, "", err))
		return
	}

	result, err := sv.repo.GetBrandByID(c, uuid.MustParse(param.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[gin.H](NotFoundCode, "", fmt.Errorf("brand with ID %s not found", param.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "", err))
		return
	}

	colResp := CategoryResponse{
		ID:          result.ID.String(),
		Name:        result.Name,
		Description: result.Description,
		Slug:        result.Slug,
		Published:   result.Published,
		CreatedAt:   result.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   result.UpdatedAt.Format("2006-01-02 15:04:05"),
		ImageUrl:    result.ImageUrl,
		Remarkable:  *result.Remarkable,
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, colResp, "", nil, nil))
}

// updateBrandHandler updates a Brand.
// @Summary Update a Brand
// @Description Update a Brand
// @ID update-Brand
// @Accept json
// @Produce json
// @Tags Admin
// @Param id path int true "Brand ID"
// @Param request body UpdateCategoryRequest true "Brand request"
// @Success 200 {object} ApiResponse[CategoryResponse]
// @Failure 400 {object} ApiResponse[CategoryResponse]
// @Failure 500 {object} ApiResponse[CategoryResponse]
// @Router /brands/{id} [put]
func (sv *Server) updateBrandHandler(c *gin.Context) {
	var param URIParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[gin.H](InvalidBodyCode, "", err))
		return
	}
	var req UpdateCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[gin.H](InvalidBodyCode, "", err))
		return
	}

	brand, err := sv.repo.GetBrandByID(c, uuid.MustParse(param.ID))

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[gin.H](NotFoundCode, "", fmt.Errorf("brand with ID %s not found", param.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "", err))
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
			c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](UploadFileCode, "", err))
			return
		}
		updateParam.ImageUrl = &imgUrl
		updateParam.ImageID = &imgID
		oldImageID := brand.ImageID
		if oldImageID != nil {
			errMsg, err := sv.uploadService.RemoveFile(c, *oldImageID)
			if err != nil {
				log.Error().Err(err).Msg("error when remove old image")
				c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](UploadFileCode, errMsg, err))
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

	if req.Remarkable != nil {
		updateParam.Remarkable = req.Remarkable
	}

	if req.Published != nil {
		updateParam.Published = req.Published
	}

	col, err := sv.repo.UpdateBrandWith(c, updateParam)

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, col, "", nil, nil))
}

// deleteBrand delete a Brand.
// @Summary Delete a Brand
// @Description Delete a Brand
// @ID delete-Brand
// @Accept json
// @Tags Admin
// @Produce json
// @Param id path int true "Brand ID"
// @Success 204 {object} ApiResponse[bool]
// @Failure 400 {object} ApiResponse[gin.H]
// @Failure 500 {object} ApiResponse[gin.H]
// @Router /brands/{id} [delete]
func (sv *Server) deleteBrand(c *gin.Context) {
	var colID URIParam
	if err := c.ShouldBindUri(&colID); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[gin.H](InvalidBodyCode, "", err))
		return
	}

	_, err := sv.repo.GetBrandByID(c, uuid.MustParse(colID.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[gin.H](NotFoundCode, "", fmt.Errorf("brand with ID %s not found", colID.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "", err))
		return
	}

	err = sv.repo.DeleteBrand(c, uuid.MustParse(colID.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "", err))
		return
	}
	c.JSON(http.StatusOK, createSuccessResponse(c, true, fmt.Sprintf("brand with ID %s deleted", colID.ID), nil, nil))
}
