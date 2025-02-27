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
type CategoryRequest struct {
	Name      string `json:"name"`
	Published *bool  `json:"published"`
	SortOrder *int16 `json:"sort_order"`
}

type CategoryResponse struct {
	ID        int32              `json:"id"`
	Name      string             `json:"name"`
	SortOrder int16              `json:"sort_order,omitempty"`
	Published bool               `json:"published,omitempty"`
	Products  []ProductListModel `json:"products,omitempty"`
}

type getCategoryParams struct {
	CategoryID int32   `uri:"id"`
	ProductID  *string `json:"product_id,omitempty"`
}
type getCategoriesQueries struct {
	IncludePublished *bool    `form:"include_published,omitempty"`
	Categories       *[]int32 `form:"categories,omitempty"`
}

type CategoryProductRequest struct {
	SortOrder int16 `json:"sort_order,omitempty"`
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
// @Router /categories [post]
func (sv *Server) createCategory(c *gin.Context) {
	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	params := repository.CreateCategoryParams{
		Name:      req.Name,
		SortOrder: 0,
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
// @Router /categories [get]
func (sv *Server) getCategories(c *gin.Context) {
	var queries getCategoriesQueries
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	var resp []CategoryResponse
	if queries.Categories != nil {
		createParams := repository.GetCategoriesByIDsParams{
			CategoryIds: *queries.Categories,
		}
		if queries.IncludePublished != nil {
			createParams.Published = utils.GetPgTypeBool(queries.IncludePublished != nil)

		}
		rows, err := sv.repo.GetCategoriesByIDs(c, createParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, mapErrResp(err))
			return
		}
		resp = groupGetCategoryByIDsRows(rows)
	} else {
		rows, err := sv.repo.GetCategories(c,
			pgtype.Bool{
				Bool:  true,
				Valid: true,
			})
		if err != nil {
			c.JSON(http.StatusInternalServerError, mapErrResp(err))
			return
		}
		resp = groupGetCategoriesRows(rows)
	}

	cnt, err := sv.repo.CountCategories(c, pgtype.Int4{
		Valid: false,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusOK, GenericListResponse[CategoryResponse]{resp, cnt, nil, nil})
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
// @Router /categories/{id} [get]
func (sv *Server) getCategoryByID(c *gin.Context) {
	var param getCategoryParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	result, err := sv.repo.GetCategoryByID(c, param.CategoryID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("category with ID %d not found", param.CategoryID)))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	colResp := CategoryResponse{
		ID:        result.CategoryID,
		Name:      result.Name,
		SortOrder: result.SortOrder,
		Published: result.Published,
		Products:  []ProductListModel{},
	}
	c.JSON(http.StatusOK, GenericResponse[CategoryResponse]{&colResp, nil, nil})
}

// getCategoryProducts retrieves a list of products in a Category.
// @Summary Get a list of products in a Category
// @Description Get a list of products in a Category
// @ID get-Category-products
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} []ProductListModel
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /category/{id}/products [get]
func (sv *Server) getCategoryProducts(c *gin.Context) {
	var param getCategoryParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	var queries PaginationQueryParams
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	category, err := sv.repo.GetCategoryByID(c, param.CategoryID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("category with ID %d not found", param.CategoryID)))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	getProductsParams := repository.GetProductsByCategoryParams{
		CategoryID: utils.GetPgTypeInt4(category.CategoryID),
		Limit:      20,
		Offset:     0,
	}

	if queries.PageSize != nil {
		getProductsParams.Limit = *queries.PageSize
		if queries.Page != nil {
			getProductsParams.Offset = (*queries.Page - 1) * *queries.PageSize
		}
	}
	productRows, err := sv.repo.GetProductsByCategory(c, getProductsParams)
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
// @Router /categories/{id} [put]
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

	category, err := sv.repo.GetCategoryByID(c, param.CategoryID)

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("category with ID %d not found", param.CategoryID)))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	updateParam := repository.UpdateCategoryWithParams{
		CategoryID: category.CategoryID,
		Name:       utils.GetPgTypeText(req.Name),
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

// deleteCategory delete a Category.
// @Summary Delete a Category
// @Description Delete a Category
// @ID delete-Category
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 204 {object} GenericResponse[repository.Category]
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /categories/{id} [delete]
func (sv *Server) deleteCategory(c *gin.Context) {
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

	err = sv.repo.DeleteCategory(c, colID.CategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	message := fmt.Sprintf("Category with ID %d deleted", colID.CategoryID)
	success := true
	c.JSON(http.StatusOK, GenericResponse[bool]{&success, &message, nil})
}

// ------------------------------------------ Helpers ------------------------------------------
func groupGetCategoriesRows(rows []repository.GetCategoriesRow) []CategoryResponse {
	categories := []CategoryResponse{}
	lastCategoryID := int32(-1)
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
		if r.CategoryID == lastCategoryID && r.ProductID.Valid {
			categories[len(categories)-1].Products = append(categories[len(categories)-1].Products, product)
		} else {
			productList := []ProductListModel{}
			if product.ID != "" {
				productList = append(productList, product)
			}
			categories = append(categories, CategoryResponse{
				ID:        r.CategoryID,
				Name:      r.Name,
				SortOrder: r.SortOrder,
				Published: r.Published,
				Products:  productList,
			})
			lastCategoryID = r.CategoryID
		}
	}
	return categories
}

func groupGetCategoryByIDsRows(rows []repository.GetCategoriesByIDsRow) []CategoryResponse {
	categories := []CategoryResponse{}
	lastCategoryID := int32(-1)
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
		if r.CategoryID == lastCategoryID && r.ProductID.Valid {
			categories[len(categories)-1].Products = append(categories[len(categories)-1].Products, product)
		} else {
			productList := []ProductListModel{}
			if product.ID != "" {
				productList = append(productList, product)
			}
			categories = append(categories, CategoryResponse{
				ID:        r.CategoryID,
				Name:      r.Name,
				SortOrder: r.SortOrder,
				Published: r.Published,
				Products:  productList,
			})
			lastCategoryID = r.CategoryID
		}
	}
	return categories
}
