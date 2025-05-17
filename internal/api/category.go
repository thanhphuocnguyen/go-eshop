package api

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

// ------------------------------------------ Request and Response ------------------------------------------
type UpdateCategoryRequest struct {
	Name         *string               `form:"name" binding:"omitempty,min=3,max=255"`
	Description  *string               `form:"description" binding:"omitempty,max=1000"`
	Slug         *string               `form:"slug" binding:"omitempty,min=3,max=255"`
	Published    *bool                 `form:"published" binding:"omitempty"`
	Remarkable   *bool                 `form:"remarkable" binding:"omitempty"`
	DisplayOrder *int16                `form:"display_order" binding:"omitempty"`
	Image        *multipart.FileHeader `form:"image" binding:"omitempty"`
}

type CreateCategoryRequest struct {
	Name         string                `form:"name" binding:"required,min=3,max=255"`
	Slug         string                `form:"slug" binding:"required,min=3,max=255"`
	Description  *string               `form:"description" binding:"omitempty,max=1000"`
	DisplayOrder *int16                `form:"display_order" binding:"omitempty"`
	Remarkable   *bool                 `form:"remarkable" binding:"omitempty"`
	Image        *multipart.FileHeader `form:"image" binding:"omitempty"`
}

type CategoryLinkedProduct struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	VariantCount int32   `json:"variant_count"`
	ImageUrl     *string `json:"image_url,omitempty"`
}

type CategoryListResponse struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description *string            `json:"description,omitempty"`
	Slug        string             `json:"slug"`
	Published   bool               `json:"published,omitempty"`
	Remarkable  bool               `json:"remarkable,omitempty"`
	CreatedAt   string             `json:"created_at,omitempty"`
	UpdatedAt   string             `json:"updated_at,omitempty"`
	ImageUrl    *string            `json:"image_url,omitempty"`
	Products    []ProductListModel `json:"products,omitempty"`
}

type CategoryResponse struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	Description *string                 `json:"description,omitempty"`
	Slug        string                  `json:"slug"`
	Published   bool                    `json:"published,omitempty"`
	Remarkable  bool                    `json:"remarkable,omitempty"`
	CreatedAt   string                  `json:"created_at,omitempty"`
	UpdatedAt   string                  `json:"updated_at,omitempty"`
	ImageUrl    *string                 `json:"image_url,omitempty"`
	Products    []CategoryLinkedProduct `json:"products"`
}

type SlugParam struct {
	Slug string `uri:"slug" binding:"required"`
}

type CategoryProductRequest struct {
	SortOrder int16 `json:"sort_order,omitempty"`
}

// ------------------------------------------ API Handlers ------------------------------------------

// --- Public API ---

// getCategoriesHandler retrieves a list of Categories.
// @Summary Get a list of Categories
// @Description Get a list of Categories
// @ID get-Categories
// @Accept json
// @Tags Categories
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} ApiResponse[CategoryResponse]
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /categories [get]
func (sv *Server) getCategoriesHandler(c *gin.Context) {
	var query PaginationQueryParams
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[CategoryResponse](InvalidBodyCode, "", err))
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
		c.JSON(http.StatusInternalServerError, createErrorResponse[CategoryResponse](InternalServerErrorCode, "", err))
		return
	}

	count, err := sv.repo.CountCategories(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[CategoryResponse](InternalServerErrorCode, "", err))
		return
	}

	categoriesResp := make([]CategoryListResponse, len(categories))
	catIds := make([]uuid.UUID, len(categories))
	for i, category := range categories {
		categoriesResp[i] = CategoryListResponse{
			ID:          category.ID.String(),
			Name:        category.Name,
			Slug:        category.Slug,
			Published:   category.Published,
			CreatedAt:   category.CreatedAt.String(),
			UpdatedAt:   category.UpdatedAt.String(),
			Description: category.Description,
			Remarkable:  *category.Remarkable,
			ImageUrl:    category.ImageUrl,
		}

		catIds[i] = category.ID
	}

	c.JSON(http.StatusOK, createSuccessResponse(
		c,
		categoriesResp,
		fmt.Sprintf("Total %d categories", count),
		nil,
		nil,
	))
}

// getCategoryBySlugHandler retrieves a list of Products by Category Slug.
// @Summary Get a list of Products by Category Slug
// @Description Get a list of Products by Category Slug
// @ID get-Products-by-Category-Slug
// @Accept json
// @Tags Categories
// @Produce json
// @Param slug path string true "Category Slug"
// @Param page_size query int false "Page size"
// @Success 200 {object} ApiResponse[CategoryResponse]
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /categories/slug/{slug} [get]
func (sv *Server) getCategoryBySlugHandler(c *gin.Context) {
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

	category, err := sv.repo.GetCategoryBySlug(c, param.Slug)
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
			ID:          category.ID.String(),
			Name:        category.Name,
			Description: category.Description,
			Slug:        category.Slug,
			ImageUrl:    category.ImageUrl,
			CreatedAt:   category.CreatedAt.String(),
		},
		"",
		nil,
		nil,
	))
}

// getCategoryProductsHandler retrieves a list of Products by Category ID.
// @Summary Get a list of Products by Category ID
// @Description Get a list of Products by Category ID
// @ID get-Products-by-Category-ID
// @Accept json
// @Tags Categories
// @Produce json
// @Param id path string true "Category ID"
// @Param page_size query int false "Page size"
// @Success 200 {object} ApiResponse[CategoryResponse]
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /categories/{id}/products [get]
func (sv *Server) getCategoryProductsHandler(c *gin.Context) {
	var param URIParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[CategoryResponse](InvalidBodyCode, "", err))
		return
	}
	var query PaginationQueryParams
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[CategoryResponse](InvalidBodyCode, "", err))
		return
	}

	category, err := sv.repo.GetCategoryByID(c, uuid.MustParse(param.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[CategoryResponse](NotFoundCode, "", fmt.Errorf("category with ID %s not found", param.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[CategoryResponse](InternalServerErrorCode, "", err))
		return
	}

	productRows, err := sv.repo.GetProducts(c, repository.GetProductsParams{
		Limit:        query.PageSize,
		Offset:       (query.Page - 1) * query.PageSize,
		IsActive:     utils.BoolPtr(true),
		Search:       nil,
		CategoryIds:  []uuid.UUID{category.ID},
		CollectionID: nil,
		Slug:         nil,
	})

	c.JSON(http.StatusOK, createSuccessResponse(
		c,
		productRows,
		fmt.Sprintf("Total %d products", len(productRows)),
		nil,
		nil,
	))
}

// --- Admin API ---
// createCategoryHandler creates a new Category.
// @Summary Create a new Category
// @Description Create a new Category
// @ID create-Category
// @Accept json
// @Tags Categories
// @Produce json
// @Param request body CreateCategoryRequest true "Category request"
// @Success 201 {object} ApiResponse[CategoryResponse]
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /admin/categories [post]
func (sv *Server) createCategoryHandler(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[CategoryResponse](InvalidBodyCode, "", err))
		return
	}
	params := repository.CreateCategoryParams{
		Name: req.Name,
		Slug: req.Slug,
	}

	if req.Description != nil {
		params.Description = req.Description
	}

	if req.Remarkable != nil {
		params.Remarkable = req.Remarkable
	}

	if req.Image != nil {
		imageID, imageURL, err := sv.uploadService.UploadFile(c, req.Image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse[CategoryResponse](UploadFileCode, "", err))
			return
		}
		params.ImageID = &imageID
		params.ImageUrl = &imageURL
	}

	col, err := sv.repo.CreateCategory(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[CategoryResponse](InternalServerErrorCode, "", err))
		return
	}
	resp := CategoryResponse{
		ID:          col.ID.String(),
		Name:        col.Name,
		Slug:        col.Slug,
		Published:   col.Published,
		CreatedAt:   col.CreatedAt.String(),
		UpdatedAt:   col.UpdatedAt.String(),
		Description: col.Description,
		Remarkable:  *col.Remarkable,
		ImageUrl:    col.ImageUrl,
	}

	c.JSON(http.StatusCreated, createSuccessResponse(c, resp, "", nil, nil))
}

// getAdminCategoriesHandler retrieves a list of Categories.
// @Summary Get a list of Categories
// @Description Get a list of Categories
// @ID get-Categories
// @Accept json
// @Tags Categories
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} ApiResponse[CategoryResponse]
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router/admin/categories [get]
func (sv *Server) getAdminCategoriesHandler(c *gin.Context) {
	var query PaginationQueryParams
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[CategoryResponse](InvalidBodyCode, "", err))
		return
	}
	params := repository.GetCategoriesParams{
		Limit:  10,
		Offset: 0,
	}
	params.Offset = (params.Limit) * int64(query.Page-1)
	params.Limit = int64(query.PageSize)

	categories, err := sv.repo.GetCategories(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[CategoryResponse](InternalServerErrorCode, "", err))
		return
	}

	count, err := sv.repo.CountCategories(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[CategoryResponse](InternalServerErrorCode, "", err))
		return
	}
	categoriesResp := make([]CategoryResponse, len(categories))

	for i, category := range categories {
		categoriesResp[i] = CategoryResponse{
			ID:          category.ID.String(),
			Name:        category.Name,
			Slug:        category.Slug,
			Published:   category.Published,
			CreatedAt:   category.CreatedAt.String(),
			UpdatedAt:   category.UpdatedAt.String(),
			Description: category.Description,
			Remarkable:  *category.Remarkable,
			ImageUrl:    category.ImageUrl,
		}
	}
	c.JSON(http.StatusOK, createSuccessResponse(
		c,
		categoriesResp,
		"",
		&Pagination{
			Page:            query.Page,
			PageSize:        query.PageSize,
			Total:           count,
			TotalPages:      count / int64(query.PageSize),
			HasNextPage:     count > int64(query.Page*query.PageSize),
			HasPreviousPage: query.Page > 1,
		}, nil,
	))
}

// getCategoryByID retrieves a Category by its ID.
// @Summary Get a Category by ID
// @Description Get a Category by ID
// @ID get-Category-by-id
// @Accept json
// @Tags Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} ApiResponse[CategoryResponse]
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /categories/{id} [get]
func (sv *Server) getCategoryByID(c *gin.Context) {
	var param URIParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[CategoryResponse](InvalidBodyCode, "", err))
		return
	}

	result, err := sv.repo.GetCategoryByID(c, uuid.MustParse(param.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[CategoryResponse](InvalidBodyCode, "", fmt.Errorf("category with ID %s not found", param.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[CategoryResponse](InternalServerErrorCode, "", err))
		return
	}

	resp := CategoryResponse{
		ID:          result.ID.String(),
		Name:        result.Name,
		Slug:        result.Slug,
		Published:   result.Published,
		CreatedAt:   result.CreatedAt.String(),
		UpdatedAt:   result.UpdatedAt.String(),
		Description: result.Description,
		Remarkable:  *result.Remarkable,
		ImageUrl:    result.ImageUrl,
	}

	getProductsParams := repository.GetLinkedProductsByCategoryParams{
		CategoryID: utils.GetPgTypeUUID(uuid.MustParse(param.ID)),
		Limit:      200,
		Offset:     0,
	}

	productRows, err := sv.repo.GetLinkedProductsByCategory(c, getProductsParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[[]ProductListModel](InternalServerErrorCode, "", err))
		return
	}

	for _, row := range productRows {
		resp.Products = append(resp.Products, CategoryLinkedProduct{
			ID:           row.ID.String(),
			Name:         row.Name,
			VariantCount: int32(row.VariantCount),
			ImageUrl:     &row.ImgUrl,
		})
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, resp, "", nil, nil))
}

// updateCategoryHandler updates a Category.
// @Summary Update a Category
// @Description Update a Category
// @ID update-Category
// @Accept json
// @Tags Admin
// @Produce json
// @Param id path int true "Category ID"
// @Param request body UpdateCategoryRequest true "Category request"
// @Success 200 {object} ApiResponse[CategoryResponse]
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /admin/categories/{id} [put]
func (sv *Server) updateCategoryHandler(c *gin.Context) {
	var param URIParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[CategoryResponse](InvalidBodyCode, "", err))
		return
	}
	var req UpdateCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[CategoryResponse](InvalidBodyCode, "", err))
		return
	}

	category, err := sv.repo.GetCategoryByID(c, uuid.MustParse(param.ID))

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[CategoryResponse](NotFoundCode, "", fmt.Errorf("category with ID %s not found", param.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[CategoryResponse](InternalServerErrorCode, "", err))
		return
	}

	updateParam := repository.UpdateCategoryParams{
		ID: category.ID,
	}

	if req.Name != nil {
		updateParam.Name = req.Name
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
	msg := ""
	var apiErr *ApiError

	imageID, imageURL := "", ""
	if req.Image != nil {
		oldImageID := category.ImageID
		oldImageURL := category.ImageUrl
		// remove old image
		if oldImageID != nil && oldImageURL != nil {
			msg, err = sv.uploadService.RemoveFile(c, *oldImageID)
			if err != nil {
				apiErr = &ApiError{
					Code:    UploadFileCode,
					Details: "Failed to remove old image",
					Stack:   err.Error()}
			}
		}
		imageID, imageURL, err = sv.uploadService.UploadFile(c, req.Image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse[CategoryResponse](UploadFileCode, "", err))
			return
		}
		updateParam.ImageID = &imageID
		updateParam.ImageUrl = &imageURL
	}
	col, err := sv.repo.UpdateCategory(c, updateParam)

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[CategoryResponse](InternalServerErrorCode, "", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, col, msg, nil, apiErr))
}

// deleteCategoryHandler delete a Category.
// @Summary Delete a Category
// @Description Delete a Category
// @ID delete-Category
// @Accept json
// @Tags Admin
// @Produce json
// @Param id path int true "Category ID"
// @Success 204 {object} ApiResponse[bool]
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /admin/categories/{id} [delete]
func (sv *Server) deleteCategoryHandler(c *gin.Context) {
	var param URIParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[bool](InvalidBodyCode, "", err))
		return
	}

	_, err := sv.repo.GetCategoryByID(c, uuid.MustParse(param.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[bool](NotFoundCode, "", fmt.Errorf("category with ID %s not found", param.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
		return
	}

	err = sv.repo.DeleteCategory(c, uuid.MustParse(param.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
		return
	}
	c.JSON(http.StatusOK, createSuccessResponse(c, true, fmt.Sprintf("Category with ID %s deleted", param.ID), nil, nil))
}
