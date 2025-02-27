package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

// ------------------------------------------ Request and Response ------------------------------------------
type CollectionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Published   *bool  `json:"published"`
	SortOrder   *int16 `json:"sort_order"`
}

type CollectionResponse struct {
	ID          int32              `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Products    []ProductListModel `json:"products,omitempty"`
}

type getCollectionParams struct {
	CollectionID int32   `uri:"id"`
	ProductID    *string `json:"product_id,omitempty"`
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
// @Success 201 {object} GenericResponse[repository.Collection]
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /collections [post]
func (sv *Server) createCollection(c *gin.Context) {
	var req CollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	params := repository.CreateCollectionParams{
		Name:        req.Name,
		Description: utils.GetPgTypeText(req.Description),
	}

	col, err := sv.repo.CreateCollection(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusCreated, GenericResponse[repository.Collection]{&col, nil, nil})
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
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	var resp []CollectionResponse = []CollectionResponse{}
	cnt := 0
	if queries.Collections != nil {
		dbParams := repository.GetCollectionsByIDsParams{
			CollectionIds: *queries.Collections,
			Limit:         20,
			Offset:        0,
		}
		if queries.Page != nil {
			dbParams.Limit = *queries.PageSize
			dbParams.Offset = (*queries.Page - 1) * *queries.PageSize
		}
		rows, err := sv.repo.GetCollectionsByIDs(c, dbParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, mapErrResp(err))
			return
		}
		resp = append(resp, groupGetCollectionByIDsRows(rows)...)
	} else {
		var dbParams repository.GetCollectionsParams = repository.GetCollectionsParams{
			Limit:  20,
			Offset: 0,
		}
		if queries.Page != nil {
			dbParams.Limit = *queries.PageSize
			dbParams.Offset = (*queries.Page - 1) * *queries.PageSize
		}
		rows, err := sv.repo.GetCollections(c, dbParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, mapErrResp(err))
			return
		}
		resp = append(resp, groupGetCollectionsRows(rows)...)
	}

	c.JSON(http.StatusOK, GenericListResponse[CollectionResponse]{resp, int64(cnt), nil, nil})
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
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	result, err := sv.repo.GetCollectionByID(c, param.CollectionID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("collection with ID %d not found", param.CollectionID)))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	colResp := CollectionResponse{
		ID:   result.CollectionID,
		Name: result.Name,
	}
	c.JSON(http.StatusOK, GenericResponse[CollectionResponse]{&colResp, nil, nil})
}

// getCollectionProducts retrieves a list of products in a Collection.
// @Summary Get a list of products in a Collection
// @Description Get a list of products in a Collection
// @ID get-Collection-products
// @Accept json
// @Produce json
// @Param id path int true "Collection ID"
// @Success 200 {object} []ProductListModel
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /Collection/{id}/products [get]
func (sv *Server) getCollectionProducts(c *gin.Context) {
	var param getCollectionParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	var queries PaginationQueryParams
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	Collection, err := sv.repo.GetCollectionByID(c, param.CollectionID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("collection with ID %d not found", param.CollectionID)))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	getProductsParams := repository.GetProductsByCollectionParams{
		CollectionID: utils.GetPgTypeInt4(Collection.CollectionID),
		Limit:        20,
		Offset:       0,
	}

	if queries.PageSize != nil {
		getProductsParams.Limit = *queries.PageSize
		if queries.Page != nil {
			getProductsParams.Offset = (*queries.Page - 1) * *queries.PageSize
		}
	}
	productRows, err := sv.repo.GetProductsByCollection(c, getProductsParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	products := []ProductListModel{}
	for _, p := range productRows {
		priceFrom, _ := p.MinPrice.Float64Value()
		priceTo, _ := p.MaxPrice.Float64Value()
		products = append(products, ProductListModel{
			ID:           p.ProductID.String(),
			Name:         p.Name,
			PriceFrom:    priceFrom.Float64,
			VariantCount: p.VariantCount,
			PriceTo:      priceTo.Float64,
			Description:  p.Description,
			DiscountTo:   p.Discount,
			ImageUrl:     &p.ImageUrl.String,
			CreatedAt:    p.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, GenericResponse[[]ProductListModel]{&products, nil, nil})
}

// updateCollection updates a Collection.
// @Summary Update a Collection
// @Description Update a Collection
// @ID update-Collection
// @Accept json
// @Produce json
// @Param id path int true "Collection ID"
// @Param request body CollectionRequest true "Collection request"
// @Success 200 {object} GenericResponse[repository.Collection]
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /collections/{id} [put]
func (sv *Server) updateCollection(c *gin.Context) {
	var param getCollectionParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	var req CollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	Collection, err := sv.repo.GetCollectionByID(c, param.CollectionID)

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("collection with ID %d not found", param.CollectionID)))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	updateParam := repository.UpdateCollectionWithParams{
		CollectionID: Collection.CollectionID,
		Name:         utils.GetPgTypeText(req.Name),
	}

	col, err := sv.repo.UpdateCollectionWith(c, updateParam)

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusOK, GenericResponse[repository.Collection]{&col, nil, nil})
}

// deleteCollection delete a Collection.
// @Summary Delete a Collection
// @Description Delete a Collection
// @ID delete-Collection
// @Accept json
// @Produce json
// @Param id path int true "Collection ID"
// @Success 204 {object} GenericResponse[repository.Collection]
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /collections/{id} [delete]
func (sv *Server) deleteCollection(c *gin.Context) {
	var colID getCollectionParams
	if err := c.ShouldBindUri(&colID); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	_, err := sv.repo.GetCollectionByID(c, colID.CollectionID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("collection with ID %d not found", colID.CollectionID)))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	err = sv.repo.DeleteCollection(c, colID.CollectionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	message := fmt.Sprintf("Collection with ID %d deleted", colID.CollectionID)
	success := true
	c.JSON(http.StatusOK, GenericResponse[bool]{&success, &message, nil})
}

// ------------------------------------------ Helpers ------------------------------------------
func groupGetCollectionsRows(rows []repository.GetCollectionsRow) []CollectionResponse {
	collections := []CollectionResponse{}
	lastCollectionID := int32(-1)
	for _, r := range rows {
		var product ProductListModel
		if r.ProductID.Valid {
			priceFrom, _ := r.PriceFrom.(pgtype.Numeric).Float64Value()
			priceTo, _ := r.PriceTo.(pgtype.Numeric).Float64Value()
			discount := r.Discount.(int16)
			product = ProductListModel{
				ID:           uuid.UUID(r.ProductID.Bytes).String(),
				Name:         r.ProductName.String,
				Description:  r.Description.String,
				VariantCount: r.VariantCount,
				ImageUrl:     &r.ImageUrl.String,
				CreatedAt:    r.CreatedAt.Format("2006-01-02 15:04:05"),
				DiscountTo:   discount,
				PriceFrom:    priceFrom.Float64,
				PriceTo:      priceTo.Float64,
			}
		}
		if r.CollectionID == lastCollectionID && r.ProductID.Valid {
			collections[len(collections)-1].Products = append(collections[len(collections)-1].Products, product)
		} else {
			productList := []ProductListModel{}
			if product.ID != "" {
				productList = append(productList, product)
			}
			collections = append(collections, CollectionResponse{
				ID:          r.CollectionID,
				Name:        r.Name,
				Description: r.Description.String,
				Products:    productList,
			})
			lastCollectionID = r.CollectionID
		}
	}
	return collections
}

func groupGetCollectionByIDsRows(rows []repository.GetCollectionsByIDsRow) []CollectionResponse {
	collections := []CollectionResponse{}
	lastCollectionID := int32(-1)
	for _, r := range rows {
		var product ProductListModel
		if r.ProductID.Valid {
			priceFrom, _ := r.PriceFrom.(pgtype.Numeric).Float64Value()
			priceTo, _ := r.PriceTo.(pgtype.Numeric).Float64Value()
			discount := r.Discount.(int16)
			product = ProductListModel{
				ID:           uuid.UUID(r.ProductID.Bytes).String(),
				Name:         r.ProductName.String,
				Description:  r.Description.String,
				VariantCount: r.VariantCount,
				ImageUrl:     &r.ImageUrl.String,
				CreatedAt:    r.CreatedAt.Format("2006-01-02 15:04:05"),
				DiscountTo:   discount,
				PriceFrom:    priceFrom.Float64,
				PriceTo:      priceTo.Float64,
			}
		}
		if r.CollectionID == lastCollectionID && r.ProductID.Valid {
			collections[len(collections)-1].Products = append(collections[len(collections)-1].Products, product)
		} else {
			productList := []ProductListModel{}
			if product.ID != "" {
				productList = append(productList, product)
			}
			collections = append(collections, CollectionResponse{
				ID:          r.CollectionID,
				Name:        r.Name,
				Description: r.Description.String,
				Products:    productList,
			})
			lastCollectionID = r.CollectionID
		}
	}
	return collections
}
