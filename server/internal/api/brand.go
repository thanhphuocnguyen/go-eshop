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

type BrandParams struct {
	ID        string  `uri:"id" binding:"required,uuid"`
	ProductID *string `json:"product_id,omitempty"`
}
type BrandsQueries struct {
	PaginationQueryParams
	Brands []int32 `form:"Brand_ids,omitempty"`
}

type BrandProductRequest struct {
	SortOrder int16 `json:"sort_order,omitempty"`
}

type BrandResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Slug        string    `json:"slug"`
	ImageUrl    *string   `json:"image_url,omitempty"`
}

// ------------------------------------------ API Handlers ------------------------------------------

// --- Public API ---
// getShopBrandsHandler retrieves a list of Brands for the shop.
// @Summary Get a list of Brands for the shop
// @Description Get a list of Brands for the shop
// @ID get-shop-Brands
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} ApiResponse[BrandResponse]
// @Failure 400 {object} ApiResponse[BrandResponse]
// @Failure 500 {object} ApiResponse[BrandResponse]
// @Router /shop/Brands [get]
func (sv *Server) getShopBrandsHandler(c *gin.Context) {
	var queries BrandsQueries
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[BrandResponse](InvalidBodyCode, "", err))
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
		c.JSON(http.StatusInternalServerError, createErrorResponse[BrandResponse](InternalServerErrorCode, "", err))
		return
	}

	cnt, err := sv.repo.CountBrands(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[BrandResponse](InternalServerErrorCode, "", err))
		return
	}

	resp := make([]BrandResponse, 0, len(rows))

	for _, row := range rows {
		resp = append(resp, BrandResponse{
			ID:          row.ID,
			Name:        row.Name,
			Description: row.Description,
			Slug:        row.Slug,
			ImageUrl:    row.ImageUrl,
		})
	}

	c.JSON(http.StatusOK, createSuccessResponse(
		c,
		resp,
		"",
		&Pagination{
			Page:            queries.Page,
			Total:           cnt,
			PageSize:        queries.PageSize,
			TotalPages:      cnt / int64(queries.PageSize),
			HasNextPage:     cnt > int64((queries.Page-1)*queries.PageSize+queries.PageSize),
			HasPreviousPage: queries.Page > 1,
		}, nil,
	))
}

// --- Admin API ---
// createBrandHandler creates a new Brand.
// @Summary Create a new Brand
// @Description Create a new Brand
// @ID create-Brand
// @Accept json
// @Produce json
// @Param request body CreateCategoryRequest true "Brand request"
// @Success 201 {object} ApiResponse[BrandResponse]
// @Failure 400 {object} ApiResponse[BrandResponse]
// @Failure 500 {object} ApiResponse[BrandResponse]
// @Router /Brands [post]
func (sv *Server) createBrandHandler(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[BrandResponse](InvalidBodyCode, "", err))
		return
	}
	params := repository.CreateBrandParams{
		ID:   uuid.New(),
		Name: req.Name,
		Slug: req.Slug,
	}

	if req.Description != nil {
		params.Description = req.Description
	}
	if req.Image != nil {
		fileName, fileID, err := sv.uploadService.UploadFile(c, req.Image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse[BrandResponse](UploadFileCode, "error when upload image", err))
			return
		}
		params.ImageUrl = &fileName
		params.ImageID = &fileID
	}

	col, err := sv.repo.CreateBrand(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[BrandResponse](InternalServerErrorCode, "", err))
		return
	}

	c.JSON(http.StatusCreated, createSuccessResponse(c, col, "", nil, nil))
}

// getBrandsHandler retrieves a list of Brands.
// @Summary Get a list of Brands
// @Description Get a list of Brands
// @ID get-Brands
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} ApiResponse[BrandResponse]
// @Failure 400 {object} ApiResponse[BrandResponse]
// @Failure 500 {object} ApiResponse[BrandResponse]
// @Router /Brands [get]
func (sv *Server) getBrandsHandler(c *gin.Context) {
	var queries BrandsQueries
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[BrandResponse](InvalidBodyCode, "", err))
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
		c.JSON(http.StatusInternalServerError, createErrorResponse[BrandResponse](InternalServerErrorCode, "", err))
		return
	}

	cnt, err := sv.repo.CountBrands(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[BrandResponse](InternalServerErrorCode, "", err))
		return
	}

	resp := make([]BrandResponse, 0, len(rows))

	for _, row := range rows {
		resp = append(resp, BrandResponse{
			ID:          row.ID,
			Name:        row.Name,
			Description: row.Description,
			Slug:        row.Slug,
			ImageUrl:    row.ImageUrl,
		})
	}

	c.JSON(http.StatusOK, createSuccessResponse(
		c,
		resp,
		"",
		&Pagination{
			Page:            queries.Page,
			Total:           cnt,
			PageSize:        queries.PageSize,
			TotalPages:      cnt / int64(queries.PageSize),
			HasNextPage:     cnt > int64((queries.Page-1)*queries.PageSize+queries.PageSize),
			HasPreviousPage: queries.Page > 1,
		}, nil,
	))
}

// getBrandByIDHandler retrieves a Brand by its ID.
// @Summary Get a Brand by ID
// @Description Get a Brand by ID
// @ID get-Brand-by-id
// @Accept json
// @Produce json
// @Param id path int true "Brand ID"
// @Success 200 {object} ApiResponse[CategoryResponse]
// @Failure 400 {object} ApiResponse[CategoryResponse]
// @Failure 500 {object} ApiResponse[CategoryResponse]
// @Router /Brands/{id} [get]
func (sv *Server) getBrandByIDHandler(c *gin.Context) {
	var param BrandParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[CategoryResponse](InvalidBodyCode, "", err))
		return
	}

	result, err := sv.repo.GetBrandByID(c, uuid.MustParse(param.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[CategoryResponse](NotFoundCode, "", fmt.Errorf("brand with ID %s not found", param.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[CategoryResponse](InternalServerErrorCode, "", err))
		return
	}

	colResp := CategoryWithProductsResponse{
		CategoryResponse: CategoryResponse{
			ID:          result.ID.String(),
			Name:        result.Name,
			Description: result.Description,
			Slug:        result.Slug,
			Published:   result.Published,
			Remarkable:  *result.Remarkable,
			CreatedAt:   result.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   result.UpdatedAt.Format("2006-01-02 15:04:05"),
			ImageUrl:    result.ImageUrl,
		},
		Products: nil,
	}
	c.JSON(http.StatusOK, createSuccessResponse(c, colResp, "", nil, nil))
}

// getProductsByBrandHandler retrieves a list of products in a Brand.
// @Summary Get a list of products in a Brand
// @Description Get a list of products in a Brand
// @ID get-Brand-products
// @Accept json
// @Produce json
// @Param id path int true "Brand ID"
// @Success 200 {object} ApiResponse[[]ProductListModel]
// @Failure 400 {object} ApiResponse[[]ProductListModel]
// @Failure 500 {object} ApiResponse[[]ProductListModel]
// @Router /Brand/{id}/products [get]
func (sv *Server) getProductsByBrandHandler(c *gin.Context) {
	var param BrandParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[[]ProductListModel](InvalidBodyCode, "", err))
		return
	}
	var queries PaginationQueryParams
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[[]ProductListModel](InvalidBodyCode, "", err))
		return
	}
	_, err := sv.repo.GetBrandByID(c, uuid.MustParse(param.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[[]ProductListModel](NotFoundCode, "", fmt.Errorf("brand with ID %s not found", param.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[[]ProductListModel](InternalServerErrorCode, "", err))
		return
	}

	getProductsParams := repository.GetProductsByBrandIDParams{
		BrandID: utils.GetPgTypeUUID(uuid.MustParse(param.ID)),
		Limit:   20,
		Offset:  0,
	}

	getProductsParams.Limit = queries.PageSize
	getProductsParams.Offset = (queries.Page - 1) * queries.PageSize

	productRows, err := sv.repo.GetProductsByBrandID(c, getProductsParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[[]ProductListModel](InternalServerErrorCode, "", err))
		return
	}
	products := []ProductListModel{}
	for _, p := range productRows {
		price, _ := p.BasePrice.Float64Value()
		products = append(products, ProductListModel{
			ID:          p.ID.String(),
			Name:        p.Name,
			Description: *p.Description,
			Price:       price.Float64,
			Slug:        p.Slug,
			Sku:         p.BaseSku,
		})
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, products, "", nil, nil))
}

// updateBrand updates a Brand.
// @Summary Update a Brand
// @Description Update a Brand
// @ID update-Brand
// @Accept json
// @Produce json
// @Param id path int true "Brand ID"
// @Param request body UpdateCategoryRequest true "Brand request"
// @Success 200 {object} ApiResponse[BrandResponse]
// @Failure 400 {object} ApiResponse[BrandResponse]
// @Failure 500 {object} ApiResponse[BrandResponse]
// @Router /Brands/{id} [put]
func (sv *Server) updateBrand(c *gin.Context) {
	var param BrandParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[BrandResponse](InvalidBodyCode, "", err))
		return
	}
	var req UpdateCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[BrandResponse](InvalidBodyCode, "", err))
		return
	}

	brand, err := sv.repo.GetBrandByID(c, uuid.MustParse(param.ID))

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[BrandResponse](NotFoundCode, "", fmt.Errorf("brand with ID %s not found", param.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[BrandResponse](InternalServerErrorCode, "", err))
		return
	}

	updateParam := repository.UpdateBrandWithParams{
		ID:   brand.ID,
		Name: req.Name,
	}
	if req.Image != nil {
		fileName, fileID, err := sv.uploadService.UploadFile(c, req.Image)
		if err != nil {
			log.Error().Err(err).Interface("value", req.Image.Header).Msg("error when upload image")
			c.JSON(http.StatusInternalServerError, createErrorResponse[BrandResponse](UploadFileCode, "", err))
			return
		}
		updateParam.ImageUrl = &fileName
		updateParam.ImageID = &fileID
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
		c.JSON(http.StatusInternalServerError, createErrorResponse[BrandResponse](InternalServerErrorCode, "", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, col, "", nil, nil))
}

// deleteBrand delete a Brand.
// @Summary Delete a Brand
// @Description Delete a Brand
// @ID delete-Brand
// @Accept json
// @Produce json
// @Param id path int true "Brand ID"
// @Success 204 {object} ApiResponse[bool]
// @Failure 400 {object} ApiResponse[bool]
// @Failure 500 {object} ApiResponse[bool]
// @Router /Brands/{id} [delete]
func (sv *Server) deleteBrand(c *gin.Context) {
	var colID BrandParams
	if err := c.ShouldBindUri(&colID); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[bool](InvalidBodyCode, "", err))
		return
	}

	_, err := sv.repo.GetBrandByID(c, uuid.MustParse(colID.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[bool](NotFoundCode, "", fmt.Errorf("brand with ID %s not found", colID.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
		return
	}

	err = sv.repo.DeleteBrand(c, uuid.MustParse(colID.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
		return
	}
	c.JSON(http.StatusOK, createSuccessResponse(c, true, fmt.Sprintf("brand with ID %s deleted", colID.ID), nil, nil))
}
