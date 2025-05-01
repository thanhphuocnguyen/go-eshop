package api

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"sync"

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
	Image        *multipart.FileHeader `form:"image" binding:"omitempty"`
}

type CategoryResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Slug        string  `json:"slug"`
	Published   bool    `json:"published,omitempty"`
	Remarkable  bool    `json:"remarkable,omitempty"`
	CreatedAt   string  `json:"created_at,omitempty"`
	UpdatedAt   string  `json:"updated_at,omitempty"`
	ImageUrl    *string `json:"image_url,omitempty"`
}

type CategoryWithProductsResponse struct {
	CategoryResponse `json:"category,inline"`
	Products         []ProductListModel `json:"products"`
}

type getCategoryParams struct {
	CategoryID string  `uri:"id" binding:"required,uuid"`
	ProductID  *string `json:"product_id,omitempty"`
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
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} ApiResponse[CategoryResponse]
// @Failure 400 {object} ApiResponse[CategoryResponse]
// @Failure 500 {object} ApiResponse[CategoryResponse]
// @Router /categories [get]
func (sv *Server) getCategoriesHandler(c *gin.Context) {
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

	categoriesResp := make([]CategoryWithProductsResponse, len(categories))
	var wg sync.WaitGroup
	productChannel := make(chan []ProductListModel, len(categories))
	for i, category := range categories {
		wg.Add(1)
		categoriesResp[i] = CategoryWithProductsResponse{
			CategoryResponse: CategoryResponse{
				ID:          category.ID.String(),
				Name:        category.Name,
				Slug:        category.Slug,
				Published:   category.Published,
				CreatedAt:   category.CreatedAt.String(),
				UpdatedAt:   category.UpdatedAt.String(),
				Description: category.Description,
				Remarkable:  *category.Remarkable,
				ImageUrl:    category.ImageUrl,
			},
			Products: []ProductListModel{},
		}
		go func() {
			defer wg.Done()
			prodByCategoryRows, err := sv.repo.GetProductsByCategoryID(c, repository.GetProductsByCategoryIDParams{
				CategoryID: utils.GetPgTypeUUID(category.ID),
				Limit:      10,
				Offset:     0,
			})
			if err != nil {
				productChannel <- []ProductListModel{}
				return
			}
			productsResp := make([]ProductListModel, len(prodByCategoryRows))
			for j, product := range prodByCategoryRows {
				price, _ := product.BasePrice.Float64Value()
				productsResp[j] = ProductListModel{
					ID:           product.ID.String(),
					Name:         product.Name,
					Slug:         product.Slug,
					CreatedAt:    product.CreatedAt.String(),
					UpdatedAt:    product.UpdatedAt.String(),
					Description:  *product.Description,
					Price:        price.Float64,
					VariantCount: int32(product.VariantCount),
					Sku:          product.BaseSku,
					ImgUrl:       product.ImgUrl,
					ImgID:        product.ImgID,
				}
			}
			productChannel <- productsResp
		}()
	}

	wg.Wait()
	close(productChannel)
	i := 0

	for products := range productChannel {
		categoriesResp[i].Products = products
		i++
	}

	c.JSON(http.StatusOK, createSuccessResponse(
		c,
		categoriesResp,
		fmt.Sprintf("Total %d categories", count),
		nil,
		nil,
	))
}

// --- Admin API ---
// addCategoryHandler creates a new Category.
// @Summary Create a new Category
// @Description Create a new Category
// @ID create-Category
// @Accept json
// @Produce json
// @Param request body CreateCategoryRequest true "Category request"
// @Success 201 {object} ApiResponse[CategoryResponse]
// @Failure 400 {object} ApiResponse[CategoryResponse]
// @Failure 500 {object} ApiResponse[CategoryResponse]
// @Router /categories [post]
func (sv *Server) addCategoryHandler(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[CategoryResponse](InvalidBodyCode, "", err))
		return
	}
	params := repository.CreateCategoryParams{
		ID:   uuid.New(),
		Name: req.Name,
		Slug: req.Slug,
	}

	if req.Description != nil {
		params.Description = req.Description
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
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} ApiResponse[CategoryResponse]
// @Failure 400 {object} ApiResponse[CategoryResponse]
// @Failure 500 {object} ApiResponse[CategoryResponse]
// @Router /categories [get]
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
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} ApiResponse[CategoryResponse]
// @Failure 400 {object} ApiResponse[CategoryResponse]
// @Failure 404 {object} ApiResponse[CategoryResponse]
// @Failure 500 {object} ApiResponse[CategoryResponse]
// @Router /categories/{id} [get]
func (sv *Server) getCategoryByID(c *gin.Context) {
	var param getCategoryParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[CategoryResponse](InvalidBodyCode, "", err))
		return
	}

	result, err := sv.repo.GetCategoryByID(c, uuid.MustParse(param.CategoryID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[CategoryResponse](InvalidBodyCode, "", fmt.Errorf("category with ID %s not found", param.CategoryID)))
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

	c.JSON(http.StatusOK, createSuccessResponse(c, resp, "", nil, nil))
}

// getProductsByCategory retrieves a list of Products by Category ID.
// @Summary Get a list of Products by Category ID
// @Description Get a list of Products by Category ID
// @ID get-Products-by-Category
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} ApiResponse[[]ProductListModel]
// @Failure 400 {object} ApiResponse[ProductListModel]
// @Failure 500 {object} ApiResponse[ProductListModel]
// @Router /categories/{id}/products [get]
func (sv *Server) getProductsByCategory(c *gin.Context) {
	var param getCategoryParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[ProductListModel](InvalidBodyCode, "", err))
		return
	}
	var query PaginationQueryParams
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[ProductListModel](InvalidBodyCode, "", err))
		return
	}

	params := repository.GetProductsByCategoryIDParams{
		CategoryID: utils.GetPgTypeUUID(uuid.MustParse(param.CategoryID)),
		Limit:      10,
		Offset:     0,
	}
	params.Offset = (params.Limit) * int64(query.Page-1)
	params.Limit = int64(query.PageSize)

	products, err := sv.repo.GetProductsByCategoryID(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[ProductListModel](InternalServerErrorCode, "", err))
		return
	}

	productsResp := make([]ProductListModel, len(products))
	for i, product := range products {
		price, _ := product.BasePrice.Float64Value()
		productsResp[i] = ProductListModel{
			ID:           product.ID.String(),
			Name:         product.Name,
			Slug:         product.Slug,
			CreatedAt:    product.CreatedAt.String(),
			UpdatedAt:    product.UpdatedAt.String(),
			Description:  *product.Description,
			Price:        price.Float64,
			VariantCount: int32(product.VariantCount),
			Sku:          product.BaseSku,
			ImgUrl:       product.ImgUrl,
			ImgID:        product.ImgID,
		}
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, productsResp, "", nil, nil))
}

// updateCategory updates a Category.
// @Summary Update a Category
// @Description Update a Category
// @ID update-Category
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param request body UpdateCategoryRequest true "Category request"
// @Success 200 {object} ApiResponse[CategoryResponse]
// @Failure 400 {object} ApiResponse[CategoryResponse]
// @Failure 500 {object} ApiResponse[CategoryResponse]
// @Router /categories/{id} [put]
func (sv *Server) updateCategory(c *gin.Context) {
	var param getCategoryParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[CategoryResponse](InvalidBodyCode, "", err))
		return
	}
	var req UpdateCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[CategoryResponse](InvalidBodyCode, "", err))
		return
	}

	category, err := sv.repo.GetCategoryByID(c, uuid.MustParse(param.CategoryID))

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[CategoryResponse](NotFoundCode, "", fmt.Errorf("category with ID %s not found", param.CategoryID)))
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

	imageID, imageURL := "", ""
	if req.Image != nil {
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

	c.JSON(http.StatusOK, createSuccessResponse(c, col, "", nil, nil))
}

// deleteCategory delete a Category.
// @Summary Delete a Category
// @Description Delete a Category
// @ID delete-Category
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 204 {object} ApiResponse[bool]
// @Failure 400 {object} ApiResponse[bool]
// @Failure 500 {object} ApiResponse[bool]
// @Router /categories/{id} [delete]
func (sv *Server) deleteCategory(c *gin.Context) {
	var colID getCategoryParams
	if err := c.ShouldBindUri(&colID); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[bool](InvalidBodyCode, "", err))
		return
	}

	_, err := sv.repo.GetCategoryByID(c, uuid.MustParse(colID.CategoryID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[bool](NotFoundCode, "", fmt.Errorf("category with ID %s not found", colID.CategoryID)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
		return
	}

	err = sv.repo.DeleteCategory(c, uuid.MustParse(colID.CategoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
		return
	}
	c.JSON(http.StatusOK, createSuccessResponse(c, true, fmt.Sprintf("Category with ID %s deleted", colID.CategoryID), nil, nil))
}
