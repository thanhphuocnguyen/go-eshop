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

type getCollectionParams struct {
	ID        string  `uri:"id" binding:"required,uuid"`
	ProductID *string `json:"product_id,omitempty"`
}
type getCollectionsQueries struct {
	PaginationQueryParams
	Collections *[]int32 `form:"collection_ids,omitempty"`
}

type CollectionProductRequest struct {
	SortOrder int16 `json:"sort_order,omitempty"`
}
type CollectionResponse struct {
	ID          string  `json:"id"`
	Slug        string  `json:"slug"`
	Description *string `json:"description,omitempty"`
	Published   bool    `json:"published"`
	Name        string  `json:"name"`
	Remarkable  bool    `json:"remarkable"`
	ImageUrl    *string `json:"image_url,omitempty"`
	CreatedAt   string  `json:"created_at"`
}

// ------------------------------------------ API Handlers ------------------------------------------

// --- Public API ---
// getShopCollectionsHandler retrieves a list of Collections.
// @Summary Get a list of Collections
// @Description Get a list of Collections
// @ID get-Shop-Collections
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} ApiResponse[CollectionResponse]
// @Failure 400 {object} ApiResponse[CollectionResponse]
// @Failure 500 {object} ApiResponse[CollectionResponse]
// @Router /shop/collections [get]
func (sv *Server) getShopCollectionsHandler(c *gin.Context) {
}

// --- Admin API ---

// createCollection creates a new Collection.
// @Summary Create a new Collection
// @Description Create a new Collection
// @ID create-Collection
// @Accept json
// @Produce json
// @Param request body CreateCategoryRequest true "Collection info"
// @Success 201 {object} ApiResponse[CollectionResponse]
// @Failure 400 {object} ApiResponse[CollectionResponse]
// @Failure 500 {object} ApiResponse[CollectionResponse]
// @Router /collections [post]
func (sv *Server) createCollection(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[CollectionResponse](InvalidBodyCode, "", err))
		return
	}
	createParams := repository.CreateCollectionParams{
		ID:   uuid.New(),
		Name: req.Name,
		Slug: req.Slug,
	}

	if req.Description != nil {
		createParams.Description = req.Description
	}

	if req.Image != nil {
		ID, url, err := sv.uploadService.UploadFile(c, req.Image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse[CollectionResponse](UploadFileCode, "", err))
			return
		}

		createParams.ImageID = &ID
		createParams.ImageUrl = &url
	}

	col, err := sv.repo.CreateCollection(c, createParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[CollectionResponse](InternalServerErrorCode, "", err))
		return
	}

	c.JSON(http.StatusCreated, createSuccessResponse(c, col, "collection", nil, nil))
}

// getCollections retrieves a list of Collections.
// @Summary Get a list of Collections
// @Description Get a list of Collections
// @ID get-Collections
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} ApiResponse[CollectionResponse]
// @Failure 400 {object} ApiResponse[CollectionResponse]
// @Failure 500 {object} ApiResponse[CollectionResponse]
// @Router /collections [get]
func (sv *Server) getCollections(c *gin.Context) {
	var queries getCollectionsQueries
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[CollectionResponse](InvalidBodyCode, "", err))
		return
	}

	dbQueries := repository.GetCollectionsParams{
		Limit:  20,
		Offset: 0,
	}

	dbQueries.Offset = int64(queries.Page-1) * int64(queries.PageSize)
	dbQueries.Limit = int64(queries.PageSize)
	rows, err := sv.repo.GetCollections(c, dbQueries)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[CollectionResponse](InternalServerErrorCode, "", err))
		return
	}

	cnt, err := sv.repo.CountCollections(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[CollectionResponse](InternalServerErrorCode, "", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(
		c,
		rows,
		"collections",
		&Pagination{
			Page:            queries.Page,
			PageSize:        queries.PageSize,
			Total:           cnt,
			TotalPages:      int(cnt / int64(queries.PageSize)),
			HasNextPage:     cnt > int64(queries.Page*queries.PageSize),
			HasPreviousPage: queries.Page > 1,
		}, nil,
	))

}

// getCollectionByID retrieves a Collection by its ID.
// @Summary Get a Collection by ID
// @Description Get a Collection by ID
// @ID get-Collection-by-id
// @Accept json
// @Produce json
// @Param id path int true "Collection ID"
// @Success 200 {object} ApiResponse[CategoryResponse]
// @Failure 400 {object} ApiResponse[CategoryResponse]
// @Failure 404 {object} ApiResponse[CategoryResponse]
// @Failure 500 {object} ApiResponse[CategoryResponse]
// @Router /collections/{id} [get]
func (sv *Server) getCollectionByID(c *gin.Context) {
	var param getCollectionParams
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
	products := make([]ProductListModel, 0)
	colResp := CategoryResponse{
		ID:          firstRow.ID.String(),
		Slug:        firstRow.Slug,
		Description: firstRow.Description,
		Published:   firstRow.Published,
		Name:        firstRow.Name,
		Remarkable:  *firstRow.Remarkable,
		ImageUrl:    firstRow.ImageUrl,
		CreatedAt:   firstRow.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   firstRow.UpdatedAt.Format("2006-01-02 15:04:05"),
		Products:    products,
	}

	for i, row := range rows {
		if !row.ProductID.Valid {
			continue
		}
		if i == 0 || row.ProductID != rows[i-1].ProductID {
			productID, _ := uuid.FromBytes(row.ProductID.Bytes[:])
			price, _ := row.ProductPrice.Float64Value()
			productMode := ProductListModel{
				ID:          productID.String(),
				Name:        *row.ProductName,
				Price:       price.Float64,
				Sku:         *row.ProductSku,
				Description: *row.Description,
			}
			colResp.Products = append(colResp.Products, productMode)
		}
	}
	c.JSON(http.StatusOK, createSuccessResponse(c, colResp, "collection", nil, nil))
}

// getProductsByCollection retrieves a list of Products by Collection ID.
// @Summary Get a list of Products by Collection ID
// @Description Get a list of Products by Collection ID
// @ID get-Products-by-Collection
// @Accept json
// @Produce json
// @Param id path int true "Collection ID"
// @Success 200 {object} ApiResponse[[]ProductListModel]
// @Failure 400 {object} ApiResponse[[]ProductListModel]
// @Failure 404 {object} ApiResponse[[]ProductListModel]
// @Failure 500 {object} ApiResponse[[]ProductListModel]
// @Router /collections/{id}/products [get]
func (sv *Server) getProductsByCollection(c *gin.Context) {
	var param getCollectionParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[[]ProductListModel](InvalidBodyCode, "", err))
		return
	}
	arg := repository.GetProductsByCollectionIDParams{
		CollectionID: utils.GetPgTypeUUID(uuid.MustParse(param.ID)),
		Limit:        20,
		Offset:       0,
	}

	rows, err := sv.repo.GetProductsByCollectionID(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[[]ProductListModel](InternalServerErrorCode, "", err))
		return
	}

	resp := make([]ProductListModel, 0)
	for _, row := range rows {
		price, _ := row.BasePrice.Float64Value()
		productMode := ProductListModel{
			ID:          row.ID.String(),
			Name:        row.Name,
			Price:       price.Float64,
			Sku:         row.BaseSku,
			Slug:        row.Slug,
			ImgUrl:      row.ImgUrl,
			ImgID:       row.ImgID,
			CreatedAt:   row.CreatedAt.Format("2006-01-02 15:04:05"),
			Description: *row.Description,
		}
		resp = append(resp, productMode)
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, rows, "products", nil, nil))
}

// updateCollection updates a Collection.
// @Summary Update a Collection
// @Description Update a Collection
// @ID update-Collection
// @Accept json
// @Produce json
// @Param id path int true "Collection ID"
// @Param request body CreateCategoryRequest true "Collection info"
// @Success 200 {object} ApiResponse[CollectionResponse]
// @Failure 400 {object} ApiResponse[CollectionResponse]
// @Failure 500 {object} ApiResponse[CollectionResponse]
// @Router /collections/{id} [put]
func (sv *Server) updateCollection(c *gin.Context) {
	var param getCollectionParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[CollectionResponse](InvalidBodyCode, "", err))
		return
	}
	var req UpdateCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[CollectionResponse](InvalidBodyCode, "", err))
		return
	}

	collection, err := sv.repo.GetCollectionByID(c, uuid.MustParse(param.ID))

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[CollectionResponse](NotFoundCode, "", fmt.Errorf("collection with ID %s not found", param.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[CollectionResponse](InternalServerErrorCode, "", err))
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
		ID, url, err := sv.uploadService.UploadFile(c, req.Image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse[CollectionResponse](InternalServerErrorCode, "", err))
			return
		}

		updateParam.ImageUrl = &url
		updateParam.ImageID = &ID
	}

	if req.Published != nil {
		updateParam.Published = req.Published
	}

	if req.Remarkable != nil {
		updateParam.Remarkable = req.Remarkable
	}

	col, err := sv.repo.UpdateCollectionWith(c, updateParam)

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[CollectionResponse](InternalServerErrorCode, "", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, col, "", nil, nil))
}

// deleteCollection delete a Collection.
// @Summary Delete a Collection
// @Description Delete a Collection
// @ID delete-Collection
// @Accept json
// @Produce json
// @Param id path int true "Collection ID"
// @Success 204 {object} ApiResponse[bool]
// @Failure 400 {object} ApiResponse[bool]
// @Failure 500 {object} ApiResponse[bool]
// @Router /collections/{id} [delete]
func (sv *Server) deleteCollection(c *gin.Context) {
	var colID getCollectionParams
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
	c.JSON(http.StatusOK, createSuccessResponse[bool](c, true, message, nil, nil))
}
