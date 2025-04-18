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

// ------------------------------------------ API Handlers ------------------------------------------
// createCollection creates a new Collection.
// @Summary Create a new Collection
// @Description Create a new Collection
// @ID create-Collection
// @Accept json
// @Produce json
// @Param request body CollectionRequest true "Collection request"
// @Success 201 {object} GenericResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /collections [post]
func (sv *Server) createCollection(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}
	createParams := repository.CreateCollectionParams{
		ID:   uuid.New(),
		Name: req.Name,
		Slug: req.Slug,
	}

	if req.Description != nil {
		createParams.Description = utils.GetPgTypeText(*req.Description)
	}

	if req.Image != nil {
		ID, url, err := sv.uploadService.UploadFile(c, req.Image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
			return
		}

		createParams.ImageID = utils.GetPgTypeText(ID)
		createParams.ImageUrl = utils.GetPgTypeText(url)
	}

	col, err := sv.repo.CreateCollection(c, createParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
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
// @Success 200 {object} []CollectionResp
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /collections [get]
func (sv *Server) getCollections(c *gin.Context) {
	var queries getCollectionsQueries
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	dbQueries := repository.GetCollectionsParams{
		Limit:  20,
		Offset: 0,
	}

	dbQueries.Limit = queries.PageSize
	dbQueries.Offset = (queries.Page - 1) * queries.PageSize
	rows, err := sv.repo.GetCollections(c, dbQueries)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	cnt, err := sv.repo.CountCollections(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
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
// @Success 200 {object} CollectionResp
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /collections/{id} [get]
func (sv *Server) getCollectionByID(c *gin.Context) {
	var param getCollectionParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	rows, err := sv.repo.GetCollectionByIDWithProducts(c, uuid.MustParse(param.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse(http.StatusBadRequest, "", fmt.Errorf("collection with ID %s not found", param.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	if len(rows) == 0 {
		c.JSON(http.StatusNotFound, createErrorResponse(http.StatusBadRequest, "", fmt.Errorf("collection with ID %s not found", param.ID)))
		return
	}

	firstRow := rows[0]
	products := make([]ProductListModel, 0)
	colResp := CategoryResponse{
		ID:          firstRow.ID,
		Slug:        firstRow.Slug,
		Description: &firstRow.Description.String,
		Published:   firstRow.Published,
		Name:        firstRow.Name,
		Remarkable:  firstRow.Remarkable.Bool,
		ImageUrl:    &firstRow.ImageUrl.String,
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
				ID:          productID,
				Name:        row.ProductName.String,
				Price:       price.Float64,
				Sku:         row.ProductSku.String,
				Description: row.Description.String,
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
// @Success 200 {object} []ProductListModel
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /collections/{id}/products [get]
func (sv *Server) getProductsByCollection(c *gin.Context) {
	var param getCollectionParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}
	arg := repository.GetProductsByCollectionIDParams{
		CollectionID: utils.GetPgTypeUUID(uuid.MustParse(param.ID)),
		Limit:        20,
		Offset:       0,
	}

	rows, err := sv.repo.GetProductsByCollectionID(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
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
// @Param request body CollectionRequest true "Collection request"
// @Success 200 {object} GenericResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /collections/{id} [put]
func (sv *Server) updateCollection(c *gin.Context) {
	var param getCollectionParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}
	var req UpdateCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	collection, err := sv.repo.GetCollectionByID(c, uuid.MustParse(param.ID))

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse(http.StatusNotFound, "", fmt.Errorf("collection with ID %s not found", param.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}

	updateParam := repository.UpdateCollectionWithParams{
		ID: collection.ID,
	}
	if req.Name != nil {
		updateParam.Name = utils.GetPgTypeText(*req.Name)
	}
	if req.Description != nil {
		updateParam.Description = utils.GetPgTypeText(*req.Description)
	}

	if req.Image != nil {
		ID, url, err := sv.uploadService.UploadFile(c, req.Image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
			return
		}

		updateParam.ImageUrl = utils.GetPgTypeText(url)
		updateParam.ImageID = utils.GetPgTypeText(ID)
	}

	if req.Published != nil {
		updateParam.Published = utils.GetPgTypeBool(*req.Published)
	}

	if req.Remarkable != nil {
		updateParam.Remarkable = utils.GetPgTypeBool(*req.Remarkable)
	}

	col, err := sv.repo.UpdateCollectionWith(c, updateParam)

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
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
// @Success 204 {object} GenericResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /collections/{id} [delete]
func (sv *Server) deleteCollection(c *gin.Context) {
	var colID getCollectionParams
	if err := c.ShouldBindUri(&colID); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	_, err := sv.repo.GetCollectionByID(c, uuid.MustParse(colID.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse(http.StatusBadRequest, "", fmt.Errorf("collection with ID %s not found", colID.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	err = sv.repo.DeleteCollection(c, uuid.MustParse(colID.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}
	message := fmt.Sprintf("Collection with ID %s deleted", colID.ID)
	c.JSON(http.StatusOK, createSuccessResponse(c, nil, message, nil, nil))
}

// ------------------------------------------ Helpers ------------------------------------------
// func groupGetCollectionsRows(rows []repository.GetCollectionsRow) []CollectionResponse {
// 	collections := []CollectionResponse{}
// 	lastID := int32(-1)
// 	for _, r := range rows {
// 		var product ProductListModel
// 		if row.ProductID.Valid {
// 			priceFrom, _ := row.PriceFrom.(pgtype.Numeric).Float64Value()
// 			priceTo, _ := row.PriceTo.(pgtype.Numeric).Float64Value()
// 			discount := row.Discount.(int16)
// 			product = ProductListModel{
// 				ID:           uuid.UUID(row.ProductID.Bytes),
// 				Name:         row.ProductName.String,
// 				Description:  row.Description.String,
// 				VariantCount: row.VariantCount,
// 				ImageUrl:     &row.ImageUrl.String,
// 				CreatedAt:    row.CreatedAt.Format("2006-01-02 15:04:05"),
// 				DiscountTo:   discount,
// 				PriceFrom:    priceFrom.Float64,
// 				PriceTo:      priceTo.Float64,
// 			}
// 		}
// 		if row.ID == lastID && row.ProductID.Valid {
// 			collections[len(collections)-1].Products = append(collections[len(collections)-1].Products, product)
// 		} else {
// 			productList := []ProductListModel{}
// 			if product.ID.String() != "" {
// 				productList = append(productList, product)
// 			}
// 			collections = append(collections, CollectionResponse{
// 				ID:          row.ID,
// 				Name:        row.Name,
// 				Description: row.Description.String,
// 				Products:    productList,
// 			})
// 			lastID = row.ID
// 		}
// 	}
// 	return collections
// }

// func groupGetCollectionByIDsRows(rows []repository.GetCollectionsByIDsRow) []CollectionResponse {
// 	collections := []CollectionResponse{}
// 	lastID := int32(-1)
// 	for _, r := range rows {
// 		var product ProductListModel
// 		if row.ProductID.Valid {
// 			priceFrom, _ := row.PriceFrom.(pgtype.Numeric).Float64Value()
// 			priceTo, _ := row.PriceTo.(pgtype.Numeric).Float64Value()
// 			discount := row.Discount.(int16)
// 			product = ProductListModel{
// 				ID:           uuid.UUID(row.ProductID.Bytes),
// 				Name:         row.ProductName.String,
// 				Description:  row.Description.String,
// 				VariantCount: row.VariantCount,
// 				ImageUrl:     &row.ImageUrl.String,
// 				CreatedAt:    row.CreatedAt.Format("2006-01-02 15:04:05"),
// 				DiscountTo:   discount,
// 				PriceFrom:    priceFrom.Float64,
// 				PriceTo:      priceTo.Float64,
// 			}
// 		}
// 		if row.ID == lastID && row.ProductID.Valid {
// 			collections[len(collections)-1].Products = append(collections[len(collections)-1].Products, product)
// 		} else {
// 			productList := []ProductListModel{}
// 			productList = append(productList, product)
// 			collections = append(collections, CollectionResponse{
// 				ID:          row.ID,
// 				Name:        row.Name,
// 				Description: row.Description.String,
// 				Products:    productList,
// 			})
// 			lastID = row.ID
// 		}
// 	}
// 	return collections
// }
