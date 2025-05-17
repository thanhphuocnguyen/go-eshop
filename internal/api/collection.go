package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

// ------------------------------------------ Request and Response ------------------------------------------

type getCollectionsQueries struct {
	PaginationQueryParams
	Collections *[]int32 `form:"collection_ids,omitempty"`
}

type CollectionProductRequest struct {
	SortOrder int16 `json:"sort_order,omitempty"`
}

// ------------------------------------------ API Handlers ------------------------------------------

// --- Public API ---

// getCollectionHandler retrieves a list of Collections.
// @Summary Get a list of Collections
// @Description Get a list of Collections
// @ID get-Shop-Collection-by-slug
// @Accept json
// @Tags Collections
// @Produce json
// @Param slug path string true "Collection slug"
// @Success 200 {object} ApiResponse[CategoryResponse]
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /shop/collections/{slug} [get]
func (sv *Server) getCollectionHandler(c *gin.Context) {
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
	var category repository.Collection
	var err error
	if id, pErr := uuid.Parse(param.Slug); pErr == nil {
		category, err = sv.repo.GetCollectionByID(c, id)
	} else {
		category, err = sv.repo.GetCollectionBySlug(c, param.Slug)
	}

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[CategoryResponse](NotFoundCode, "", fmt.Errorf("category with slug %s not found", param.Slug)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[CategoryResponse](InternalServerErrorCode, "", err))
		return
	}

	resp := CategoryResponse{
		ID:          category.ID.String(),
		Name:        category.Name,
		Description: category.Description,
		Slug:        category.Slug,
		ImageUrl:    category.ImageUrl,
		CreatedAt:   category.CreatedAt.String(),
	}
	c.JSON(http.StatusOK, createSuccessResponse(c, resp, "", nil, nil))
}

// getCollectionProductsHandler retrieves a list of Products in a Collection.
// @Summary Get a list of Products in a Collection
// @Description Get a list of Products in a Collection
// @ID get-Shop-Collection-Products-by-slug
// @Accept json
// @Tags Collections
// @Produce json
// @Param slug path string true "Collection slug"
// @Param request body CollectionProductRequest true "Collection info"
// @Success 200 {object} ApiResponse[ProductListModel]
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /shop/collections/{slug}/products [get]
func (sv *Server) getCollectionProductsHandler(c *gin.Context) {
	var param SlugParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[ProductListModel](InvalidBodyCode, "", err))
		return
	}
	var query PaginationQueryParams
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[ProductListModel](InvalidBodyCode, "", err))
		return
	}

	var req CollectionProductRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[ProductListModel](InvalidBodyCode, "", err))
		return
	}

	colID, err := sv.repo.GetCollectionBySlug(c, param.Slug)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[ProductListModel](NotFoundCode, "", fmt.Errorf("collection with slug %s not found", param.Slug)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[ProductListModel](InternalServerErrorCode, "", err))
		return
	}

	getProductsParams := repository.GetProductsParams{
		CollectionID: []uuid.UUID{colID.ID},
		IsActive:     utils.BoolPtr(true),
		Limit:        int64(query.PageSize),
		Offset:       int64((query.Page - 1) * query.PageSize),
	}

	productRows, err := sv.repo.GetProducts(c, getProductsParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[ProductListModel](InternalServerErrorCode, "", err))
		return
	}

	cnt, err := sv.repo.CountProducts(c, repository.CountProductsParams{
		IsActive:     utils.BoolPtr(true),
		Name:         nil,
		CollectionID: utils.GetPgTypeUUID(colID.ID),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[ProductListModel](InternalServerErrorCode, "", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(
		c,
		productRows,
		"products",
		&Pagination{
			Page:            query.Page,
			PageSize:        query.PageSize,
			Total:           cnt,
			TotalPages:      cnt / int64(query.PageSize),
			HasNextPage:     cnt > int64(query.Page*query.PageSize),
			HasPreviousPage: query.Page > 1,
		}, nil,
	))
}

// --- Admin API ---

// createCollectionHandler creates a new Collection.
// @Summary Create a new Collection
// @Description Create a new Collection
// @ID create-Collection
// @Accept json
// @Tags Admin
// @Produce json
// @Param request body CreateCategoryRequest true "Collection info"
// @Success 201 {object} ApiResponse[CategoryResponse]
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /admin/collections [post]
func (sv *Server) createCollectionHandler(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[CategoryResponse](InvalidBodyCode, "", err))
		return
	}
	createParams := repository.CreateCollectionParams{
		Name: req.Name,
		Slug: req.Slug,
	}

	if req.Description != nil {
		createParams.Description = req.Description
	}

	if req.Remarkable != nil {
		createParams.Remarkable = req.Remarkable
	}

	if req.Image != nil {
		ID, url, err := sv.uploadService.UploadFile(c, req.Image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse[CategoryResponse](UploadFileCode, "", err))
			return
		}

		createParams.ImageID = &ID
		createParams.ImageUrl = &url
	}

	col, err := sv.repo.CreateCollection(c, createParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[CategoryResponse](InternalServerErrorCode, "", err))
		return
	}

	c.JSON(http.StatusCreated, createSuccessResponse(c, col, "collection", nil, nil))
}

// getCollectionsHandler retrieves a list of Collections.
// @Summary Get a list of Collections
// @Description Get a list of Collections
// @ID get-Collections
// @Accept json
// @Tags Admin
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} ApiResponse[CategoryResponse]
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /admin/collections [get]
func (sv *Server) getCollectionsHandler(c *gin.Context) {
	var queries getCollectionsQueries
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[CategoryResponse](InvalidBodyCode, "", err))
		return
	}

	dbQueries := repository.GetCollectionsParams{
		Limit:  20,
		Offset: 0,
	}

	dbQueries.Offset = int64(queries.Page-1) * int64(queries.PageSize)
	dbQueries.Limit = int64(queries.PageSize)
	collectionRows, err := sv.repo.GetCollections(c, dbQueries)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[CategoryResponse](InternalServerErrorCode, "", err))
		return
	}

	cnt, err := sv.repo.CountCollections(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[CategoryResponse](InternalServerErrorCode, "", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(
		c,
		collectionRows,
		"collections",
		&Pagination{
			Page:            queries.Page,
			PageSize:        queries.PageSize,
			Total:           cnt,
			TotalPages:      cnt / int64(queries.PageSize),
			HasNextPage:     cnt > int64(queries.Page*queries.PageSize),
			HasPreviousPage: queries.Page > 1,
		}, nil,
	))
}

// getCollectionByIDHandler retrieves a Collection by its ID.
// @Summary Get a Collection by ID
// @Description Get a Collection by ID
// @ID get-Collection-by-id
// @Accept json
// @Tags Admin
// @Produce json
// @Param id path int true "Collection ID"
// @Success 200 {object} ApiResponse[CategoryResponse]
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /admin/collections/{id} [get]
func (sv *Server) getCollectionByIDHandler(c *gin.Context) {
	var param URIParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[CategoryResponse](InvalidBodyCode, "", err))
		return
	}

	rows, err := sv.repo.GetCollectionByIDWithProducts(c, uuid.MustParse(param.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[CategoryResponse](NotFoundCode, "", fmt.Errorf("collection with ID %s not found", param.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[CategoryResponse](InternalServerErrorCode, "", err))
		return
	}

	firstRow := rows[0]
	colResp := CategoryResponse{
		ID:          firstRow.ID.String(),
		Slug:        firstRow.Slug,
		Description: firstRow.Description,
		Published:   firstRow.Published,
		Name:        firstRow.Name,
		ImageUrl:    firstRow.ImageUrl,
		CreatedAt:   firstRow.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   firstRow.UpdatedAt.Format("2006-01-02 15:04:05"),
		Products:    []CategoryLinkedProduct{},
	}
	if firstRow.Remarkable != nil {
		colResp.Remarkable = *firstRow.Remarkable
	}

	getProductsParams := repository.GetLinkedProductsByCategoryParams{
		CollectionID: utils.GetPgTypeUUID(uuid.MustParse(param.ID)),
		Limit:        200,
		Offset:       0,
	}

	productRows, err := sv.repo.GetLinkedProductsByCategory(c, getProductsParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[[]ProductListModel](InternalServerErrorCode, "", err))
		return
	}
	for _, row := range productRows {
		colResp.Products = append(colResp.Products, CategoryLinkedProduct{
			ID:           row.ID.String(),
			Name:         row.Name,
			VariantCount: int32(row.VariantCount),
			ImageUrl:     &row.ImgUrl,
		})
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, colResp, "collection", nil, nil))
}

// updateCollectionHandler updates a Collection.
// @Summary Update a Collection
// @Description Update a Collection
// @ID update-Collection
// @Accept json
// @Tags Admin
// @Produce json
// @Param id path int true "Collection ID"
// @Param request body CreateCategoryRequest true "Collection info"
// @Success 200 {object} ApiResponse[CategoryResponse]
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /admin/collections/{id} [put]
func (sv *Server) updateCollectionHandler(c *gin.Context) {
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

	collection, err := sv.repo.GetCollectionByID(c, uuid.MustParse(param.ID))

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[CategoryResponse](NotFoundCode, "", fmt.Errorf("collection with ID %s not found", param.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[CategoryResponse](InternalServerErrorCode, "", err))
		return
	}

	updateParam := repository.UpdateCollectionWithParams{
		ID: collection.ID,
	}
	if req.Name != nil {
		updateParam.Name = req.Name
	}
	if req.Description != nil {
		updateParam.Description = req.Description
	}

	if req.Image != nil {
		oldImageID := collection.ImageID
		oldImageUrl := collection.ImageUrl
		ID, url, err := sv.uploadService.UploadFile(c, req.Image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse[CategoryResponse](InternalServerErrorCode, "", err))
			return
		}

		updateParam.ImageUrl = &url
		updateParam.ImageID = &ID

		// Delete old image
		if oldImageID != nil && oldImageUrl != nil {
			if msg, err := sv.uploadService.RemoveFile(c, *oldImageID); err != nil {
				c.JSON(http.StatusInternalServerError, createErrorResponse[CategoryResponse](InternalServerErrorCode, msg, err))
				return
			}
		}
	}

	if req.Published != nil {
		updateParam.Published = req.Published
	}

	if req.Remarkable != nil {
		updateParam.Remarkable = req.Remarkable
	}

	col, err := sv.repo.UpdateCollectionWith(c, updateParam)

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[CategoryResponse](InternalServerErrorCode, "", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, col, "", nil, nil))
}

// deleteCollection delete a Collection.
// @Summary Delete a Collection
// @Description Delete a Collection
// @ID delete-Collection
// @Accept json
// @Tags Admin
// @Produce json
// @Param id path int true "Collection ID"
// @Success 204 {object} ApiResponse[bool]
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /admin/collections/{id} [delete]
func (sv *Server) deleteCollection(c *gin.Context) {
	var colID URIParam
	if err := c.ShouldBindUri(&colID); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[bool](InvalidBodyCode, "", err))
		return
	}

	_, err := sv.repo.GetCollectionByID(c, uuid.MustParse(colID.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[bool](NotFoundCode, "", fmt.Errorf("collection with ID %s not found", colID.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
		return
	}

	err = sv.repo.DeleteCollection(c, uuid.MustParse(colID.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
		return
	}
	message := fmt.Sprintf("Collection with ID %s deleted", colID.ID)
	c.JSON(http.StatusOK, createSuccessResponse(c, true, message, nil, nil))
}
