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

type getBrandParams struct {
	ID        string  `uri:"id" binding:"required,uuid"`
	ProductID *string `json:"product_id,omitempty"`
}
type getBrandsQueries struct {
	PaginationQueryParams
	Brands []int32 `form:"Brand_ids,omitempty"`
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
// @Success 201 {object} GenericResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /Brands [post]
func (sv *Server) createBrand(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}
	params := repository.CreateBrandParams{
		ID:   uuid.New(),
		Name: req.Name,
		Slug: req.Slug,
	}

	if req.Description != nil {
		params.Description = utils.GetPgTypeText(*req.Description)
	}
	if req.Image != nil {
		fileName, fileID, err := sv.uploadService.UploadFile(c, *req.Image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
			return
		}
		params.ImageUrl = utils.GetPgTypeText(fileName)
		params.ImageID = utils.GetPgTypeText(fileID)
	}

	col, err := sv.repo.CreateBrand(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	c.JSON(http.StatusCreated, createSuccessResponse(c, col, "", nil, nil))
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
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}
	var dbQueries repository.GetBrandsParams = repository.GetBrandsParams{
		Limit:  20,
		Offset: 0,
	}
	dbQueries.Limit = queries.PageSize
	dbQueries.Offset = (queries.Page - 1) * queries.PageSize

	rows, err := sv.repo.GetBrands(c, dbQueries)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	cnt, err := sv.repo.CountBrands(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(
		c,
		rows,
		"",
		&Pagination{
			Page:            queries.Page,
			Total:           cnt,
			PageSize:        queries.PageSize,
			TotalPages:      int(cnt / int64(queries.PageSize)),
			HasNextPage:     cnt > int64((queries.Page-1)*queries.PageSize+queries.PageSize),
			HasPreviousPage: queries.Page > 1,
		}, nil,
	))
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
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	result, err := sv.repo.GetBrandByID(c, uuid.MustParse(param.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse(http.StatusBadRequest, "", fmt.Errorf("brand with ID %s not found", param.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	colResp := CategoryResponse{
		ID:          result.ID,
		Description: &result.Description.String,
		Slug:        result.Slug,
		Published:   result.Published,
		Remarkable:  result.Remarkable.Bool,
		CreatedAt:   result.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   result.UpdatedAt.Format("2006-01-02 15:04:05"),
		ImageUrl:    &result.ImageUrl.String,
		Products:    nil,
		Name:        result.Name,
	}
	c.JSON(http.StatusOK, createSuccessResponse(c, colResp, "", nil, nil))
}

// getProductsByBrand retrieves a list of products in a Brand.
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
func (sv *Server) getProductsByBrand(c *gin.Context) {
	var param getBrandParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}
	var queries PaginationQueryParams
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}
	_, err := sv.repo.GetBrandByID(c, uuid.MustParse(param.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse(http.StatusBadRequest, "", fmt.Errorf("brand with ID %s not found", param.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	getProductsParams := repository.GetProductsByBrandIDParams{
		BrandID: utils.GetPgTypeUUID(uuid.MustParse(param.ID)),
		Limit:   20,
		Offset:  0,
	}

	getProductsParams.Limit = queries.PageSize
	getProductsParams.Offset = (queries.Page - 1) * queries.PageSize

	productRows, err := sv.repo.GetProductsByBrandID(c, getProductsParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}
	products := []ProductListModel{}
	for _, p := range productRows {
		price, _ := p.BasePrice.Float64Value()
		products = append(products, ProductListModel{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description.String,
			Price:       price.Float64,
			Slug:        p.Slug,
			Sku:         p.BaseSku.String,
		})
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, products, "", nil, nil))
}

// updateBrand updates a Brand.
// @Summary Update a Brand
// @Description Update a Brand
// @ID update-Brand
// @Accept json
// @Produce json
// @Param id path int true "Brand ID"
// @Param request body BrandRequest true "Brand request"
// @Success 200 {object} GenericResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /Brands/{id} [put]
func (sv *Server) updateBrand(c *gin.Context) {
	var param getBrandParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}
	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	brand, err := sv.repo.GetBrandByID(c, uuid.MustParse(param.ID))

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse(http.StatusBadRequest, "", fmt.Errorf("brand with ID %s not found", param.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	updateParam := repository.UpdateBrandWithParams{
		ID:   brand.ID,
		Name: utils.GetPgTypeText(*req.Name),
	}
	if req.Image != nil {
		fileName, fileID, err := sv.uploadService.UploadFile(c, *req.Image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
			return
		}
		updateParam.ImageUrl = utils.GetPgTypeText(fileName)
		updateParam.ImageID = utils.GetPgTypeText(fileID)
	}

	if req.Slug != nil {
		updateParam.Slug = utils.GetPgTypeText(*req.Slug)
	}
	if req.Description != nil {
		updateParam.Description = utils.GetPgTypeText(*req.Description)
	}

	if req.Remarkable != nil {
		updateParam.Remarkable = utils.GetPgTypeBool(*req.Remarkable)
	}

	if req.Published != nil {
		updateParam.Published = utils.GetPgTypeBool(*req.Published)
	}

	col, err := sv.repo.UpdateBrandWith(c, updateParam)

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, col, "", nil, nil))
}

// deleteBrand delete a Brand.
// @Summary Delete a Brand
// @Description Delete a Brand
// @ID delete-Brand
// @Accept json
// @Produce json
// @Param id path int true "Brand ID"
// @Success 204 {object} GenericResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /Brands/{id} [delete]
func (sv *Server) deleteBrand(c *gin.Context) {
	var colID getBrandParams
	if err := c.ShouldBindUri(&colID); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	_, err := sv.repo.GetBrandByID(c, uuid.MustParse(colID.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse(http.StatusBadRequest, "", fmt.Errorf("brand with ID %s not found", colID.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	err = sv.repo.DeleteBrand(c, uuid.MustParse(colID.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}
	message := fmt.Sprintf("brand with ID %s deleted", colID.ID)
	c.JSON(http.StatusOK, createSuccessResponse(c, message, "", nil, nil))
}

// ------------------------------------------ Helpers ------------------------------------------
// func groupGetBrandsRows(rows []repository.GetBrandsRow) []BrandResponse {
// 	Brands := []BrandResponse{}
// 	lastID := int32(-1)
// 	for _, r := range rows {
// 		var product ProductListModel
// 		if r.ID.Valid {
// 			priceFrom, _ := r.PriceFrom.(pgtype.Numeric).Float64Value()
// 			priceTo, _ := r.PriceTo.(pgtype.Numeric).Float64Value()
// 			discount := r.Discount.(int16)
// 			product = ProductListModel{
// 				ID:           uuid.UUID(r.ID.Bytes),
// 				Name:         r.ProductName.String,
// 				Description:  r.Description.String,
// 				VariantCount: r.VariantCount,
// 				ImageUrl:     &r.ImageUrl.String,
// 				CreatedAt:    r.CreatedAt.Format("2006-01-02 15:04:05"),
// 				DiscountTo:   discount,
// 				PriceFrom:    priceFrom.Float64,
// 				PriceTo:      priceTo.Float64,
// 			}
// 		}
// 		if r.ID == lastID && r.ID.Valid {
// 			Brands[len(Brands)-1].Products = append(Brands[len(Brands)-1].Products, product)
// 		} else {
// 			productList := []ProductListModel{}
// 			if product.ID.String() != "" {
// 				productList = append(productList, product)
// 			}
// 			Brands = append(Brands, BrandResponse{
// 				ID:          r.ID,
// 				Name:        r.Name,
// 				Description: r.Description.String,
// 				Products:    productList,
// 			})
// 			lastID = r.ID
// 		}
// 	}
// 	return Brands
// }

// func groupGetBrandByIDsRows(rows []repository.GetBrandsByIDsRow) []BrandResponse {
// 	Brands := []BrandResponse{}
// 	lastID := int32(-1)
// 	for _, r := range rows {
// 		var product ProductListModel
// 		if r.ID.Valid {
// 			priceFrom, _ := r.PriceFrom.(pgtype.Numeric).Float64Value()
// 			priceTo, _ := r.PriceTo.(pgtype.Numeric).Float64Value()
// 			discount := r.Discount.(int16)
// 			product = ProductListModel{
// 				ID:           uuid.UUID(r.ID.Bytes),
// 				Name:         r.ProductName.String,
// 				Description:  r.Description.String,
// 				VariantCount: r.VariantCount,
// 				ImageUrl:     &r.ImageUrl.String,
// 				CreatedAt:    r.CreatedAt.Format("2006-01-02 15:04:05"),
// 				DiscountTo:   discount,
// 				PriceFrom:    priceFrom.Float64,
// 				PriceTo:      priceTo.Float64,
// 			}
// 		}
// 		if r.ID == lastID && r.ID.Valid {
// 			Brands[len(Brands)-1].Products = append(Brands[len(Brands)-1].Products, product)
// 		} else {
// 			productList := []ProductListModel{}
// 			if product.ID.String() != "" {
// 				productList = append(productList, product)
// 			}
// 			Brands = append(Brands, BrandResponse{
// 				ID:          r.ID,
// 				Name:        r.Name,
// 				Description: r.Description.String,
// 				Products:    productList,
// 			})
// 			lastID = r.ID
// 		}
// 	}
// 	return Brands
// }
