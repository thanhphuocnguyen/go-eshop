package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
)

// --- Public API ---

// @Summary Get a list of Collections
// @Description Get a list of Collections
// @ID get-Shop-Collection-by-slug
// @Accept json
// @Tags Collections
// @Produce json
// @Param slug path string true "Collection slug"
// @Success 200 {object} ApiResponse[CategoryDto]
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /collections/{slug} [get]
func (sv *Server) GetCollectionBySlugHandler(c *gin.Context) {
	var param SlugParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode,
			err))
		return
	}
	var query PaginationQueryParams
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode,
			err))
		return
	}

	collection, err := sv.repo.GetCollectionBySlug(c, param.Slug)

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErr(NotFoundCode,
				fmt.Errorf("category with slug %s not found", param.Slug)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode,
			err))
		return
	}

	rows, err := sv.repo.GetProductList(c, repository.GetProductListParams{
		CollectionIds: []uuid.UUID{collection.ID},
		Limit:         query.PageSize,
		Offset:        (query.PageSize) * int64(query.Page-1),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode,
			err))
		return
	}
	collectionResp := CategoryDto{
		ID:          collection.ID.String(),
		Name:        collection.Name,
		Description: collection.Description,
		Slug:        collection.Slug,
		ImageUrl:    collection.ImageUrl,
		CreatedAt:   collection.CreatedAt.String(),
		Products:    make([]ProductSummary, len(rows)),
	}
	for i, row := range rows {
		collectionResp.Products[i] = mapToShopProductResponse(row)
	}

	c.JSON(http.StatusOK, createDataResp(c, collectionResp, nil, nil))
}

// --- Admin API ---

// @Summary Create a new Collection
// @Description Create a new Collection
// @ID create-Collection
// @Accept json
// @Tags Admin
// @Produce json
// @Param request body CreateCategoryRequest true "Collection info"
// @Success 201 {object} ApiResponse[CategoryDto]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/collections [post]
func (sv *Server) CreateCollectionHandler(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode,
			err))
		return
	}
	createParams := repository.CreateCollectionParams{
		Name: req.Name,
		Slug: req.Slug,
	}

	if req.Description != nil {
		createParams.Description = req.Description
	}

	if req.Image != nil {
		ID, url, err := sv.uploadService.UploadFile(c, req.Image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErr(UploadFileCode,
				err))
			return
		}

		createParams.ImageID = &ID
		createParams.ImageUrl = &url
	}

	col, err := sv.repo.CreateCollection(c, createParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode,
			err))
		return
	}
	sv.cachesrv.Delete(c, "collections-*")

	c.JSON(http.StatusCreated, createDataResp(c, col, nil, nil))
}

// @Summary Get a list of Collections
// @Description Get a list of Collections
// @ID get-Collections
// @Accept json
// @Tags Admin
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} ApiResponse[CategoryDto]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/collections [get]
func (sv *Server) GetCollectionsHandler(c *gin.Context) {
	var queries CollectionsQueryParams
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode,
			err))
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
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode,
			err))
		return
	}

	cnt, err := sv.repo.CountCollections(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode,
			err))
		return
	}

	c.JSON(http.StatusOK, createDataResp(c, collectionRows, createPagination(cnt, queries.Page, queries.PageSize), nil))
}

// @Summary Get a Collection by ID
// @Description Get a Collection by ID
// @ID get-Collection-by-id
// @Accept json
// @Tags Admin
// @Produce json
// @Param id path int true "Collection ID"
// @Success 200 {object} ApiResponse[CategoryDto]
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/collections/{id} [get]
func (sv *Server) GetCollectionByIDHandler(c *gin.Context) {
	var param UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode,
			err))
		return
	}

	collection, err := sv.repo.GetCollectionByID(c, uuid.MustParse(param.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErr(NotFoundCode,
				fmt.Errorf("collection with ID %s not found", param.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode,
			err))
		return
	}

	colResp := AdminCategoryDto{
		ID:          collection.ID.String(),
		Slug:        collection.Slug,
		Description: collection.Description,
		Published:   collection.Published,
		Name:        collection.Name,
		ImageUrl:    collection.ImageUrl,
		CreatedAt:   collection.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   collection.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, createDataResp(c, colResp, nil, nil))
}

// @Summary Update a Collection
// @Description Update a Collection
// @ID update-Collection
// @Accept json
// @Tags Admin
// @Produce json
// @Param id path int true "Collection ID"
// @Param request body CreateCategoryRequest true "Collection info"
// @Success 200 {object} ApiResponse[CategoryDto]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/collections/{id} [put]
func (sv *Server) UpdateCollectionHandler(c *gin.Context) {
	var param UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode,
			err))
		return
	}
	var req UpdateCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode,
			err))
		return
	}

	collection, err := sv.repo.GetCollectionByID(c, uuid.MustParse(param.ID))

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErr(NotFoundCode,
				fmt.Errorf("collection with ID %s not found", param.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode,
			err))
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
			c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode,
				err))
			return
		}

		updateParam.ImageUrl = &url
		updateParam.ImageID = &ID

		// Delete old image
		if oldImageID != nil && oldImageUrl != nil {
			if _, err := sv.uploadService.RemoveFile(c, *oldImageID); err != nil {
				c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
				return
			}
		}
	}

	if req.Published != nil {
		updateParam.Published = req.Published
	}

	col, err := sv.repo.UpdateCollectionWith(c, updateParam)

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode,
			err))
		return
	}

	c.JSON(http.StatusOK, createDataResp(c, col,
		nil, nil))
}

// @Summary Delete a Collection
// @Description Delete a Collection
// @ID delete-Collection
// @Accept json
// @Tags Admin
// @Produce json
// @Param id path int true "Collection ID"
// @Success 204 {object} nil
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/collections/{id} [delete]
func (sv *Server) DeleteCollectionHandler(c *gin.Context) {
	var colID UriIDParam
	if err := c.ShouldBindUri(&colID); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode,
			err))
		return
	}

	_, err := sv.repo.GetCollectionByID(c, uuid.MustParse(colID.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErr(NotFoundCode,
				fmt.Errorf("collection with ID %s not found", colID.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode,
			err))
		return
	}

	err = sv.repo.DeleteCollection(c, uuid.MustParse(colID.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode,
			err))
		return
	}
	c.Status(http.StatusNoContent)
}
