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
type BrandRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Published   *bool  `json:"published"`
	SortOrder   *int16 `json:"sort_order"`
}

type BrandResponse struct {
	ID          int32              `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Products    []ProductListModel `json:"products,omitempty"`
}

type getBrandParams struct {
	BrandID   int32   `uri:"id"`
	ProductID *string `json:"product_id,omitempty"`
}
type getBrandsQueries struct {
	PaginationQueryParams
	Brands *[]int32 `form:"Brand_ids,omitempty"`
}

type BrandProductRequest struct {
	SortOrder int16 `json:"sort_order,omitempty"`
}

// ------------------------------------------ API Handlers ------------------------------------------
// createBrand creates a new Brand.
// @Summary Create a new Brand
// @Description Create a new Brand
// @ID create-Brand
// @Accept json
// @Produce json
// @Param request body BrandRequest true "Brand request"
// @Success 201 {object} GenericResponse[repository.Brand]
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /Brands [post]
func (sv *Server) createBrand(c *gin.Context) {
	var req BrandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	params := repository.CreateBrandParams{
		Name:        req.Name,
		Description: utils.GetPgTypeText(req.Description),
	}

	col, err := sv.repo.CreateBrand(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusCreated, GenericResponse[repository.Brand]{&col, nil, nil})
}

// getBrands retrieves a list of Brands.
// @Summary Get a list of Brands
// @Description Get a list of Brands
// @ID get-Brands
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} []BrandResp
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /Brands [get]
func (sv *Server) getBrands(c *gin.Context) {
	var queries getBrandsQueries
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	var resp []BrandResponse
	if queries.Brands != nil {
		dbParams := repository.GetBrandsByIDsParams{
			BrandIds: *queries.Brands,
			Limit:    20,
			Offset:   0,
		}
		if queries.Page != nil {
			dbParams.Limit = *queries.PageSize
			dbParams.Offset = (*queries.Page - 1) * *queries.PageSize
		}
		rows, err := sv.repo.GetBrandsByIDs(c, dbParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, mapErrResp(err))
			return
		}
		resp = groupGetBrandByIDsRows(rows)
	} else {
		var dbParams repository.GetBrandsParams = repository.GetBrandsParams{
			Limit:  20,
			Offset: 0,
		}
		if queries.Page != nil {
			dbParams.Limit = *queries.PageSize
			dbParams.Offset = (*queries.Page - 1) * *queries.PageSize
		}
		rows, err := sv.repo.GetBrands(c, dbParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, mapErrResp(err))
			return
		}
		resp = groupGetBrandsRows(rows)
	}

	cnt, err := sv.repo.CountBrands(c, pgtype.Int4{
		Valid: false,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusOK, GenericListResponse[BrandResponse]{resp, cnt, nil, nil})
}

// getBrandByID retrieves a Brand by its ID.
// @Summary Get a Brand by ID
// @Description Get a Brand by ID
// @ID get-Brand-by-id
// @Accept json
// @Produce json
// @Param id path int true "Brand ID"
// @Success 200 {object} BrandResp
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /Brands/{id} [get]
func (sv *Server) getBrandByID(c *gin.Context) {
	var param getBrandParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	result, err := sv.repo.GetBrandByID(c, param.BrandID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("brand with ID %d not found", param.BrandID)))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	colResp := BrandResponse{
		ID:   result.BrandID,
		Name: result.Name,
	}
	c.JSON(http.StatusOK, GenericResponse[BrandResponse]{&colResp, nil, nil})
}

// getBrandProducts retrieves a list of products in a Brand.
// @Summary Get a list of products in a Brand
// @Description Get a list of products in a Brand
// @ID get-Brand-products
// @Accept json
// @Produce json
// @Param id path int true "Brand ID"
// @Success 200 {object} []ProductListModel
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /Brand/{id}/products [get]
func (sv *Server) getBrandProducts(c *gin.Context) {
	var param getBrandParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	var queries PaginationQueryParams
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	Brand, err := sv.repo.GetBrandByID(c, param.BrandID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("brand with ID %d not found", param.BrandID)))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	getProductsParams := repository.GetProductsByBrandParams{
		BrandID: utils.GetPgTypeInt4(Brand.BrandID),
		Limit:   20,
		Offset:  0,
	}

	if queries.PageSize != nil {
		getProductsParams.Limit = *queries.PageSize
		if queries.Page != nil {
			getProductsParams.Offset = (*queries.Page - 1) * *queries.PageSize
		}
	}
	productRows, err := sv.repo.GetProductsByBrand(c, getProductsParams)
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

// updateBrand updates a Brand.
// @Summary Update a Brand
// @Description Update a Brand
// @ID update-Brand
// @Accept json
// @Produce json
// @Param id path int true "Brand ID"
// @Param request body BrandRequest true "Brand request"
// @Success 200 {object} GenericResponse[repository.Brand]
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /Brands/{id} [put]
func (sv *Server) updateBrand(c *gin.Context) {
	var param getBrandParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	var req BrandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	Brand, err := sv.repo.GetBrandByID(c, param.BrandID)

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("brand with ID %d not found", param.BrandID)))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	updateParam := repository.UpdateBrandWithParams{
		BrandID: Brand.BrandID,
		Name:    utils.GetPgTypeText(req.Name),
	}

	col, err := sv.repo.UpdateBrandWith(c, updateParam)

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusOK, GenericResponse[repository.Brand]{&col, nil, nil})
}

// deleteBrand delete a Brand.
// @Summary Delete a Brand
// @Description Delete a Brand
// @ID delete-Brand
// @Accept json
// @Produce json
// @Param id path int true "Brand ID"
// @Success 204 {object} GenericResponse[repository.Brand]
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /Brands/{id} [delete]
func (sv *Server) deleteBrand(c *gin.Context) {
	var colID getBrandParams
	if err := c.ShouldBindUri(&colID); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	_, err := sv.repo.GetBrandByID(c, colID.BrandID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("brand with ID %d not found", colID.BrandID)))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	err = sv.repo.DeleteBrand(c, colID.BrandID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	message := fmt.Sprintf("brand with ID %d deleted", colID.BrandID)
	success := true
	c.JSON(http.StatusOK, GenericResponse[bool]{&success, &message, nil})
}

// ------------------------------------------ Helpers ------------------------------------------
func groupGetBrandsRows(rows []repository.GetBrandsRow) []BrandResponse {
	Brands := []BrandResponse{}
	lastBrandID := int32(-1)
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
		if r.BrandID == lastBrandID && r.ProductID.Valid {
			Brands[len(Brands)-1].Products = append(Brands[len(Brands)-1].Products, product)
		} else {
			productList := []ProductListModel{}
			if product.ID != "" {
				productList = append(productList, product)
			}
			Brands = append(Brands, BrandResponse{
				ID:          r.BrandID,
				Name:        r.Name,
				Description: r.Description.String,
				Products:    productList,
			})
			lastBrandID = r.BrandID
		}
	}
	return Brands
}

func groupGetBrandByIDsRows(rows []repository.GetBrandsByIDsRow) []BrandResponse {
	Brands := []BrandResponse{}
	lastBrandID := int32(-1)
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
		if r.BrandID == lastBrandID && r.ProductID.Valid {
			Brands[len(Brands)-1].Products = append(Brands[len(Brands)-1].Products, product)
		} else {
			productList := []ProductListModel{}
			if product.ID != "" {
				productList = append(productList, product)
			}
			Brands = append(Brands, BrandResponse{
				ID:          r.BrandID,
				Name:        r.Name,
				Description: r.Description.String,
				Products:    productList,
			})
			lastBrandID = r.BrandID
		}
	}
	return Brands
}
