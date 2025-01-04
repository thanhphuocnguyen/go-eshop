package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/util"
)

// ------------------------------------------ Request and Response ------------------------------------------
type collectionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Published   *bool  `json:"published"`
	SortOrder   *int16 `json:"sort_order"`
}

type collectionResp struct {
	ID          int32                 `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Products    []productListResponse `json:"products"`
}

type getCollectionParams struct {
	ID int32 `uri:"id"`
}

type getCollectionsQueries struct {
	Categories []int32 `form:"categories"`
}

type collectionProductRequest struct {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	params := repository.CreateCollectionParams{
		Name:        req.Name,
		Description: util.GetPgTypeText(req.Description),
		SortOrder:   0,
	}
	if req.Published != nil {
		params.Published = *req.Published
	}
	col, err := sv.repo.CreateCollection(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	colRows, err := sv.repo.GetCollections(c, queries.Categories)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(colRows) == 0 {
		c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("no collections found")))
		return
	}

	cols := []collectionResp{}
	for _, col := range colRows {
		colResp := collectionResp{
			ID:          col.CategoryID,
			Name:        col.Name,
			Description: col.Description.String,
			Products:    []productListResponse{},
		}
		for _, p := range colRows {
			price, _ := p.Price.Float64Value()
			colResp.Products = append(colResp.Products, productListResponse{
				ID:          p.ProductID,
				Name:        p.Name,
				Price:       price.Float64,
				Description: p.Description.String,
				ImageUrl:    &p.ImageUrl.String,
			})
		}
		cols = append(cols, colResp)
	}
	c.JSON(http.StatusOK, GenericResponse[[]collectionResp]{&cols, nil, nil})
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
	var colID getCollectionParams
	if err := c.ShouldBindUri(&colID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	colRows, err := sv.repo.GetCollection(c, colID.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(colRows) == 0 {
		c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("collection with ID %d not found", colID.ID)))
		return
	}

	col := colRows[0]
	colResp := collectionResp{
		ID:          col.CategoryID,
		Name:        col.Name,
		Description: col.Description.String,
		Products:    []productListResponse{},
	}
	for _, p := range colRows {
		price, _ := p.Price.Float64Value()
		colResp.Products = append(colResp.Products, productListResponse{
			ID:          p.ProductID,
			Name:        p.Name,
			Price:       price.Float64,
			Description: p.Description.String,
			ImageUrl:    &p.ImageUrl.String,
		})
	}
	c.JSON(http.StatusOK, GenericResponse[collectionResp]{&colResp, nil, nil})
}

func (sv *Server) updateCollection(c *gin.Context) {
	var param getCollectionParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req collectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := sv.repo.GetCollection(c, param.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("collection with ID %d not found", param.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	updateParam := repository.UpdateCollectionParams{
		CategoryID:  param.ID,
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, GenericResponse[repository.Category]{&col, nil, nil})
}

func (sv *Server) addProductToCollection(c *gin.Context) {
	var req collectionProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := sv.repo.GetProduct(c, repository.GetProductParams{
		ProductID: req.ProductID,
		Archived:  util.GetPgTypeBool(false),
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("product with ID %d not found", req.ProductID)))
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var colID getCollectionParams
	if err := c.ShouldBindUri(&colID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err = sv.repo.GetCollectionProduct(c, repository.GetCollectionProductParams{
		CategoryID: colID.ID,
		ProductID:  req.ProductID,
	})
	if err == nil {
		c.JSON(http.StatusConflict, mapErrResp(fmt.Errorf("product with ID %d already exists in collection with ID %d", req.ProductID, colID.ID)))
		return
	}

	cp, err := sv.repo.AddProductToCollection(c, repository.AddProductToCollectionParams{
		CategoryID: colID.ID,
		ProductID:  req.ProductID,
		SortOrder:  0},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := sv.repo.GetCollection(c, colID.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("collection with ID %d not found", colID.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = sv.repo.RemoveCollection(c, colID.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	message := fmt.Sprintf("collection with ID %d deleted", colID.ID)
	success := true
	c.JSON(http.StatusOK, GenericResponse[bool]{&success, &message, nil})
}

func (sv *Server) deleteProductFromCollection(c *gin.Context) {
	var req collectionProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var param getCollectionParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := sv.repo.GetCollectionProduct(c, repository.GetCollectionProductParams{
		CategoryID: param.ID,
		ProductID:  req.ProductID,
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("product with ID %d not found in collection with ID %d", req.ProductID, param.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = sv.repo.RemoveProductFromCollection(c, repository.RemoveProductFromCollectionParams{
		CategoryID: param.ID,
		ProductID:  req.ProductID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
