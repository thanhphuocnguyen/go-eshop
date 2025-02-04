package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"golang.org/x/sync/errgroup"
)

// ------------------------------------------ Request and Response ------------------------------------------
type CategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Published   *bool  `json:"published"`
	SortOrder   *int16 `json:"sort_order"`
}

type CategoryResp struct {
	ID          int32              `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Products    []ProductListModel `json:"products"`
}

type getCategoryParams struct {
	CategoryID int32  `uri:"id"`
	ProductID  *int64 `json:"product_id,omitempty"`
}
type getCategoriesQueries struct {
	Categories *[]int32 `form:"categories,omitempty"`
}

type CategoryProductRequest struct {
	SortOrder int16 `json:"sort_order,omitempty"`
}

type addProductToCategoryRequest struct {
	ProductID int64 `json:"product_id"`
}

// ------------------------------------------ API Handlers ------------------------------------------
// createCategory creates a new Category.
// @Summary Create a new Category
// @Description Create a new Category
// @ID create-Category
// @Accept json
// @Produce json
// @Param request body CategoryRequest true "Category request"
// @Success 201 {object} GenericResponse[repository.Category]
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /Categories [post]
func (sv *Server) createCategory(c *gin.Context) {
	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	params := repository.CreateCategoryParams{
		Name:        req.Name,
		Description: utils.GetPgTypeText(req.Description),
		SortOrder:   0,
	}

	col, err := sv.repo.CreateCategory(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusCreated, GenericResponse[repository.Category]{&col, nil, nil})
}

// getCategories retrieves a list of Categories.
// @Summary Get a list of Categories
// @Description Get a list of Categories
// @ID get-Categories
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} []CategoryResp
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /Categories [get]
func (sv *Server) getCategories(c *gin.Context) {
	var queries getCategoriesQueries
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	errGroup, ctx := errgroup.WithContext(c)

	total := make(chan int64, 1)
	colRows := make(chan []repository.GetCategoriesRow, 1)
	defer close(total)
	defer close(colRows)
	errGroup.Go(func() error {
		rows, err := sv.repo.GetCategories(ctx, utils.GetPgTypeBool(true))
		if err != nil {
			return err
		}
		colRows <- rows
		return nil
	})

	errGroup.Go(func() error {
		cnt, err := sv.repo.CountCategories(ctx, pgtype.Int4{
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
		c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("no Categories found")))
		return
	}
	resp := <-colRows

	c.JSON(http.StatusOK, GenericListResponse[repository.GetCategoriesRow]{&resp, <-total, nil, nil})
}

// getCategoryByID retrieves a Category by its ID.
// @Summary Get a Category by ID
// @Description Get a Category by ID
// @ID get-Category-by-id
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} CategoryResp
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /Categories/{id} [get]
func (sv *Server) getCategoryByID(c *gin.Context) {
	var param getCategoryParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	colRows, err := sv.repo.GetCategoryWithProduct(c, param.CategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	if len(colRows) == 0 {
		c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("category with ID %d not found", param.CategoryID)))
		return
	}

	col := colRows[0]
	colResp := CategoryResp{
		ID:          col.CategoryID,
		Name:        col.Name,
		Description: col.Description.String,
		Products:    []ProductListModel{},
	}
	for _, p := range colRows {
		priceFrom, _ := p.PriceFrom.Float64Value()
		priceTo, _ := p.PriceTo.Float64Value()
		colResp.Products = append(colResp.Products, ProductListModel{
			ID:           p.ProductID,
			Name:         p.ProductName,
			PriceFrom:    priceFrom.Float64,
			VariantCount: p.VariantCount,
			PriceTo:      priceTo.Float64,
			Description:  p.Description.String,
			DiscountTo:   p.Discount,
			ImageUrl:     &p.ImageUrl.String,
			CreatedAt:    p.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	c.JSON(http.StatusOK, GenericResponse[CategoryResp]{&colResp, nil, nil})
}

// updateCategory updates a Category.
// @Summary Update a Category
// @Description Update a Category
// @ID update-Category
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param request body CategoryRequest true "Category request"
// @Success 200 {object} GenericResponse[repository.Category]
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /Categories/{id} [put]
func (sv *Server) updateCategory(c *gin.Context) {
	var param getCategoryParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	_, err := sv.repo.GetCategoryWithProduct(c, param.CategoryID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("category with ID %d not found", param.CategoryID)))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	updateParam := repository.UpdateCategoryWithParams{
		CategoryID:  param.CategoryID,
		Name:        utils.GetPgTypeText(req.Name),
		Description: utils.GetPgTypeText(req.Description),
	}
	if req.Published != nil {
		updateParam.Published = utils.GetPgTypeBool(*req.Published)
	}
	if req.SortOrder != nil {
		updateParam.SortOrder = utils.GetPgTypeInt2(*req.SortOrder)
	}
	col, err := sv.repo.UpdateCategoryWith(c, updateParam)

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusOK, GenericResponse[repository.Category]{&col, nil, nil})
}

// addProductToCategory adds a product to a Category.
// @Summary Add a product to a Category
// @Description Add a product to a Category
// @ID add-product-to-Category
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param request body addProductToCategoryRequest true "Product ID"
// @Success 201 {object} GenericResponse[repository.CategoryProduct]
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /Categories/{id}/product [post]
func (sv *Server) addProductToCategory(c *gin.Context) {
	var req addProductToCategoryRequest
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

	var param getCategoryParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	cnt, err := sv.repo.CountCategories(c, utils.GetPgTypeInt4(param.CategoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	_, err = sv.repo.GetCategoryProduct(c, repository.GetCategoryProductParams{
		CategoryID: param.CategoryID,
		ProductID:  req.ProductID,
	})

	if err == nil {
		c.JSON(http.StatusConflict, mapErrResp(fmt.Errorf("product with ID %d already exists in Category with ID %d", req.ProductID, param.CategoryID)))
		return
	}

	var maxSortOrder int16
	if cnt == 0 {
		maxSortOrder = 0
	} else {
		maxSortOrder, err := sv.repo.GetMaxSortOrderInCategory(c, param.CategoryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, mapErrResp(err))
			return
		}
		if maxSortOrder == nil {
			maxSortOrder = 0
		}
	}

	cp, err := sv.repo.AddProductToCategory(c, repository.AddProductToCategoryParams{
		CategoryID: param.CategoryID,
		ProductID:  req.ProductID,
		SortOrder:  maxSortOrder + 1,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusCreated, GenericResponse[repository.CategoryProduct]{&cp, nil, nil})
}

// removeCategory deletes a Category.
// @Summary Delete a Category
// @Description Delete a Category
// @ID delete-Category
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 204 {object} GenericResponse[repository.Category]
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /Categories/{id} [delete]
func (sv *Server) removeCategory(c *gin.Context) {
	var colID getCategoryParams
	if err := c.ShouldBindUri(&colID); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	_, err := sv.repo.GetCategoryByID(c, colID.CategoryID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("category with ID %d not found", colID.CategoryID)))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	err = sv.repo.RemoveCategory(c, colID.CategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	message := fmt.Sprintf("Category with ID %d deleted", colID.CategoryID)
	success := true
	c.JSON(http.StatusOK, GenericResponse[bool]{&success, &message, nil})
}

// deleteProductFromCategory deletes a product from a Category.
// @Summary Delete a product from a Category
// @Description Delete a product from a Category
// @ID delete-product-from-Category
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param product_id body CategoryProductRequest true "Product ID"
// @Success 204
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /Categories/{id}/product [delete]
func (sv *Server) deleteProductFromCategory(c *gin.Context) {
	var param getCategoryParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	if param.ProductID == nil {
		c.JSON(http.StatusBadRequest, mapErrResp(fmt.Errorf("product_id is required")))
		return
	}

	_, err := sv.repo.GetCategoryProduct(c, repository.GetCategoryProductParams{
		CategoryID: param.CategoryID,
		ProductID:  *param.ProductID,
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("product with ID %d not found in Category with ID %d", *param.ProductID, param.CategoryID)))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	err = sv.repo.RemoveProductFromCategory(c, repository.RemoveProductFromCategoryParams{
		CategoryID: param.CategoryID,
		ProductID:  *param.ProductID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// updateProductSortOrder updates the sort order of a product in a Category.
// @Summary Update the sort order of a product in a Category
// @Description Update the sort order of a product in a Category
// @ID update-product-sort-order
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param product_id body CategoryProductRequest true "Product ID"
// @Param sort_order body int16 true "Sort order"
// @Success 204
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /Categories/{id}/product/sort-order [put]
func (sv *Server) updateProductSortOrder(c *gin.Context) {
	var req CategoryProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	var param getCategoryParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	if param.ProductID == nil {
		c.JSON(http.StatusBadRequest, mapErrResp(fmt.Errorf("product_id is required")))
		return
	}

	_, err := sv.repo.GetCategoryProduct(c, repository.GetCategoryProductParams{
		CategoryID: param.CategoryID,
		ProductID:  *param.ProductID,
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("product with ID %d not found in Category with ID %d", *param.ProductID, param.CategoryID)))
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

	err = sv.repo.UpdateProductSortOrderInCategory(c, repository.UpdateProductSortOrderInCategoryParams{
		CategoryID: param.CategoryID,
		ProductID:  *param.ProductID,
		SortOrder:  sortReq.SortOrder,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
