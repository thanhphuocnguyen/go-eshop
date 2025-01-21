package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/util"
	"golang.org/x/sync/errgroup"
)

// ------------------------------------------ Request and Response ------------------------------------------
type collectionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Published   *bool  `json:"published"`
	SortOrder   *int16 `json:"sort_order"`
}

type collectionResp struct {
	ID          int32              `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Products    []productListModel `json:"products"`
}

type getCollectionParams struct {
	CollectionID int32  `uri:"id"`
	ProductID    *int64 `json:"product_id,omitempty"`
}
type getCollectionsQueries struct {
	Categories *[]int32 `form:"categories,omitempty"`
}

type collectionProductRequest struct {
	SortOrder int16 `json:"sort_order,omitempty"`
}

type addProductToCollectionRequest struct {
	ProductID int64 `json:"product_id"`
}

// ------------------------------------------ API Handlers ------------------------------------------
// createCollection creates a new collection.
// @Summary Create a new collection
// @Description Create a new collection
// @ID create-collection
// @Accept json
// @Produce json
// @Param request body collectionRequest true "Collection request"
// @Success 201 {object} GenericResponse[repository.Category]
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /collections [post]
func (sv *Server) createCollection(c *gin.Context) {
	var req collectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	params := repository.CreateCollectionParams{
		Name:        req.Name,
		Description: util.GetPgTypeText(req.Description),
		SortOrder:   0,
	}

	col, err := sv.repo.CreateCollection(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusCreated, GenericResponse[repository.Category]{&col, nil, nil})
}

// getCollections retrieves a list of collections.
// @Summary Get a list of collections
// @Description Get a list of collections
// @ID get-collections
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} []collectionResp
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /collections [get]
func (sv *Server) getCollections(c *gin.Context) {
	var queries getCollectionsQueries
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	errGroup, ctx := errgroup.WithContext(c)

	total := make(chan int64, 1)
	colRows := make(chan []repository.Category, 1)
	defer close(total)
	defer close(colRows)
	errGroup.Go(func() error {
		rows, err := sv.repo.GetCollections(ctx)
		if err != nil {
			return err
		}
		colRows <- rows
		return nil
	})

	errGroup.Go(func() error {
		cnt, err := sv.repo.CountCollections(ctx, pgtype.Int4{
			Valid: false,
		})

		if err != nil {
			return err
		}
		total <- cnt
		return nil
	})

	if err := errGroup.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if len(colRows) == 0 {
		c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("no collections found")))
		return
	}
	resp := <-colRows

	c.JSON(http.StatusOK, GenericListResponse[repository.Category]{&resp, <-total, nil, nil})
}

// getCollectionByID retrieves a collection by its ID.
// @Summary Get a collection by ID
// @Description Get a collection by ID
// @ID get-collection-by-id
// @Accept json
// @Produce json
// @Param id path int true "Collection ID"
// @Success 200 {object} collectionResp
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /collections/{id} [get]
func (sv *Server) getCollectionByID(c *gin.Context) {
	var param getCollectionParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	colRows, err := sv.repo.GetCollection(c, param.CollectionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	if len(colRows) == 0 {
		c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("collection with ID %d not found", param.CollectionID)))
		return
	}

	col := colRows[0]
	colResp := collectionResp{
		ID:          col.CategoryID,
		Name:        col.Name,
		Description: col.Description.String,
		Products:    []productListModel{},
	}
	for _, p := range colRows {
		priceFrom, _ := p.PriceFrom.Float64Value()
		priceTo, _ := p.PriceTo.Float64Value()
		colResp.Products = append(colResp.Products, productListModel{
			ID:           p.ProductID,
			Name:         p.Name,
			PriceFrom:    priceFrom.Float64,
			VariantCount: p.VariantCount,
			PriceTo:      priceTo.Float64,
			Description:  p.Description.String,
			ImageUrl:     &p.ImageUrl.String,
		})
	}
	c.JSON(http.StatusOK, GenericResponse[collectionResp]{&colResp, nil, nil})
}

func (sv *Server) updateCollection(c *gin.Context) {
	var param getCollectionParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	var req collectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	_, err := sv.repo.GetCollection(c, param.CollectionID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("collection with ID %d not found", param.CollectionID)))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	updateParam := repository.UpdateCollectionParams{
		CategoryID:  param.CollectionID,
		Name:        util.GetPgTypeText(req.Name),
		Description: util.GetPgTypeText(req.Description),
	}
	if req.Published != nil {
		updateParam.Published = util.GetPgTypeBool(*req.Published)
	}
	if req.SortOrder != nil {
		updateParam.SortOrder = util.GetPgTypeInt2(*req.SortOrder)
	}
	col, err := sv.repo.UpdateCollection(c, updateParam)

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusOK, GenericResponse[repository.Category]{&col, nil, nil})
}

func (sv *Server) addProductToCollection(c *gin.Context) {
	var req addProductToCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	_, err := sv.repo.GetProductByID(c, repository.GetProductByIDParams{
		ProductID: req.ProductID,
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("product with ID %d not found", req.ProductID)))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	var param getCollectionParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	cnt, err := sv.repo.CountCollections(c, util.GetPgTypeInt4(param.CollectionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	_, err = sv.repo.GetCollectionProduct(c, repository.GetCollectionProductParams{
		CategoryID: param.CollectionID,
		ProductID:  *param.ProductID,
	})

	if err == nil {
		c.JSON(http.StatusConflict, mapErrResp(fmt.Errorf("product with ID %d already exists in collection with ID %d", req.ProductID, param.CollectionID)))
		return
	}

	var maxSortOrder int16
	if cnt == 0 {
		maxSortOrder = 0
	} else {
		maxSortOrder, err = sv.repo.GetMaxSortOrderInCollection(c, param.CollectionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, mapErrResp(err))
			return
		}
	}

	cp, err := sv.repo.AddProductToCollection(c, repository.AddProductToCollectionParams{
		CategoryID: param.CollectionID,
		ProductID:  *param.ProductID,
		SortOrder:  maxSortOrder + 1,
	},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusCreated, GenericResponse[repository.CategoryProduct]{&cp, nil, nil})
}

// removeCollection deletes a collection.
// @Summary Delete a collection
// @Description Delete a collection
// @ID delete-collection
// @Accept json
// @Produce json
// @Param id path int true "Collection ID"
// @Success 204 {object} GenericResponse[repository.Category]
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /collections/{id} [delete]
func (sv *Server) removeCollection(c *gin.Context) {
	var colID getCollectionParams
	if err := c.ShouldBindUri(&colID); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	_, err := sv.repo.GetCollection(c, colID.CollectionID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("collection with ID %d not found", colID.CollectionID)))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	err = sv.repo.RemoveCollection(c, colID.CollectionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	message := fmt.Sprintf("collection with ID %d deleted", colID.CollectionID)
	success := true
	c.JSON(http.StatusOK, GenericResponse[bool]{&success, &message, nil})
}

// deleteProductFromCollection deletes a product from a collection.
// @Summary Delete a product from a collection
// @Description Delete a product from a collection
// @ID delete-product-from-collection
// @Accept json
// @Produce json
// @Param id path int true "Collection ID"
// @Param product_id body collectionProductRequest true "Product ID"
// @Success 204
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /collections/{id}/product [delete]
func (sv *Server) deleteProductFromCollection(c *gin.Context) {
	var param getCollectionParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	if param.ProductID == nil {
		c.JSON(http.StatusBadRequest, mapErrResp(fmt.Errorf("product_id is required")))
		return
	}

	_, err := sv.repo.GetCollectionProduct(c, repository.GetCollectionProductParams{
		CategoryID: param.CollectionID,
		ProductID:  *param.ProductID,
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("product with ID %d not found in collection with ID %d", *param.ProductID, param.CollectionID)))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	err = sv.repo.RemoveProductFromCollection(c, repository.RemoveProductFromCollectionParams{
		CategoryID: param.CollectionID,
		ProductID:  *param.ProductID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// updateProductSortOrder updates the sort order of a product in a collection.
// @Summary Update the sort order of a product in a collection
// @Description Update the sort order of a product in a collection
// @ID update-product-sort-order
// @Accept json
// @Produce json
// @Param id path int true "Collection ID"
// @Param product_id body collectionProductRequest true "Product ID"
// @Param sort_order body int16 true "Sort order"
// @Success 204
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /collections/{id}/product/sort-order [put]
func (sv *Server) updateProductSortOrder(c *gin.Context) {
	var req collectionProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	var param getCollectionParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	if param.ProductID == nil {
		c.JSON(http.StatusBadRequest, mapErrResp(fmt.Errorf("product_id is required")))
		return
	}

	_, err := sv.repo.GetCollectionProduct(c, repository.GetCollectionProductParams{
		CategoryID: param.CollectionID,
		ProductID:  *param.ProductID,
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("product with ID %d not found in collection with ID %d", *param.ProductID, param.CollectionID)))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	var sortReq struct {
		SortOrder int16 `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&sortReq); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	err = sv.repo.UpdateProductSortOrderInCollection(c, repository.UpdateProductSortOrderInCollectionParams{
		CategoryID: param.CollectionID,
		ProductID:  *param.ProductID,
		SortOrder:  sortReq.SortOrder,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
