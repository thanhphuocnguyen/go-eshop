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

// getCollectionBySlugHandler retrieves a list of Collections.
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
// @Router /collections/{slug} [get]
func (sv *Server) getCollectionBySlugHandler(c *gin.Context) {
	var param SlugParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[CategoryResponse](InvalidBodyCode, "", err))
		return
	}

	collection, err := sv.repo.GetCollectionBySlug(c, param.Slug)

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[CategoryResponse](NotFoundCode, "", fmt.Errorf("category with slug %s not found", param.Slug)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[CategoryResponse](InternalServerErrorCode, "", err))
		return
	}

	filters, err := sv.repo.GetFilterListForCollectionID(c, collection.ID)

	collectionResp := CollectionDetailResponse{
		Collection: CategoryResponse{
			ID:          collection.ID.String(),
			Name:        collection.Name,
			Description: collection.Description,
			Slug:        collection.Slug,
			ImageUrl:    collection.ImageUrl,
			CreatedAt:   collection.CreatedAt.String(),
			Published:   collection.Published,
			UpdatedAt:   collection.UpdatedAt.String(),
		},
		Categories: make([]FiltersModel, 0),
		Brands:     make([]FiltersModel, 0),
		Attributes: make(map[string][]FiltersModel),
	}
	listAttrs := make([]uuid.UUID, 0)
	for _, row := range filters {
		if row.CategoryID.Valid {
			idx := -1
			id, _ := uuid.FromBytes(row.CategoryID.Bytes[:])
			for i, c := range collectionResp.Categories {
				if c.ID == id.String() {
					idx = i
					break
				}
			}
			if idx == -1 {
				collectionResp.Categories = append(collectionResp.Categories, FiltersModel{
					ID:   id.String(),
					Name: *row.CategoryName,
				})
			}
		}
		if row.BrandID.Valid {
			id, _ := uuid.FromBytes(row.BrandID.Bytes[:])
			idx := -1
			for i, b := range collectionResp.Brands {
				if b.ID == id.String() {
					idx = i
					break
				}
			}
			if idx == -1 {
				collectionResp.Brands = append(collectionResp.Brands, FiltersModel{
					ID:   id.String(),
					Name: *row.BrandName,
				})
			}
		}

		if len(row.Attributes) > 0 {
			for _, attr := range row.Attributes {
				idx := -1
				for i, a := range listAttrs {
					if a == attr {
						idx = i
						break
					}
				}
				if idx == -1 {
					listAttrs = append(listAttrs, attr)
				}
			}
		}
	}
	attributes, err := sv.repo.GetAttributeWithValuesByIDs(c, listAttrs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[CategoryResponse](InternalServerErrorCode, "", err))
		return
	}

	for _, attr := range attributes {
		if !attr.AttributeValueID.Valid {
			continue
		}

		id, _ := uuid.FromBytes(attr.AttributeValueID.Bytes[:])
		if _, ok := collectionResp.Attributes[attr.AttributeName]; !ok {
			collectionResp.Attributes[attr.AttributeName] = []FiltersModel{{
				ID:   id.String(),
				Name: *attr.AttributeValueName,
			}}

		}
		idx := -1
		for i, a := range collectionResp.Attributes[attr.AttributeName] {
			if a.ID == id.String() {
				idx = i
				break
			}
		}
		if idx == -1 {
			collectionResp.Attributes[attr.AttributeName] = append(collectionResp.Attributes[attr.AttributeName], FiltersModel{
				ID:   id.String(),
				Name: *attr.AttributeValueName,
			})
		}
	}

	if collection.Remarkable != nil {
		collectionResp.Collection.Remarkable = *collection.Remarkable
	}
	c.JSON(http.StatusOK, createSuccessResponse(c, collectionResp, "", nil, nil))
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
	var param UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
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

	colResp := CategoryResponse{
		ID:          collection.ID.String(),
		Slug:        collection.Slug,
		Description: collection.Description,
		Published:   collection.Published,
		Name:        collection.Name,
		ImageUrl:    collection.ImageUrl,
		CreatedAt:   collection.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   collection.UpdatedAt.Format("2006-01-02 15:04:05"),
		Products:    []CategoryLinkedProduct{},
	}
	if collection.Remarkable != nil {
		colResp.Remarkable = *collection.Remarkable
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
	var param UriIDParam
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

// deleteCollectionHandler delete a Collection.
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
func (sv *Server) deleteCollectionHandler(c *gin.Context) {
	var colID UriIDParam
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
