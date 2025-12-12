package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
)

// @Summary Get a product detail by ID
// @Schemes http
// @Description get a product detail by ID
// @Tags products
// @Accept json
// @Param productId path int true "Product ID"
// @Produce json
// @Success 200 {object} ApiResponse[ProductDetailDto]
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /products/{productId} [get]
func (sv *Server) getProductById(c *gin.Context) {
	var params models.UriIDParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	sqlParams := repository.GetProductDetailParams{}
	err := uuid.Validate(params.ID)
	if err == nil {
		sqlParams.ID = uuid.MustParse(params.ID)
	} else {
		sqlParams.Slug = params.ID
	}

	productRow, err := sv.repo.GetProductDetail(c, sqlParams)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, err))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	productDetail := dto.MapToProductDetailResponse(productRow)

	c.JSON(http.StatusOK, dto.CreateDataResp(c, productDetail, nil, nil))
}

// @Summary Get list of products
// @Schemes http
// @Description get list of products
// @Tags products
// @Accept json
// @Param page query int true "Page number"
// @Param pageSize query int true "Page size"
// @Produce json
// @Success 200 {array} ApiResponse[[]ProductSummary]
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /products [get]
func (sv *Server) getProducts(c *gin.Context) {
	var queries models.ProductQuery
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	dbParams := repository.GetProductListParams{
		Limit:  queries.PageSize,
		Offset: (queries.Page - 1) * queries.PageSize,
	}

	if queries.Search != nil && len(*queries.Search) > 0 {
		search := *queries.Search
		search = strings.ReplaceAll(search, " ", "%")
		search = strings.ReplaceAll(search, ",", "%")
		search = strings.ReplaceAll(search, ":", "%")
		search = "%" + search + "%"
		dbParams.Search = &search
	}

	if queries.BrandIDs != nil {
		dbParams.BrandIds = make([]uuid.UUID, 0)
		for _, id := range *queries.BrandIDs {
			dbParams.BrandIds = append(dbParams.BrandIds, uuid.MustParse(id))
		}
	}

	if queries.CategoryIDs != nil {
		dbParams.CategoryIds = make([]uuid.UUID, len(*queries.CategoryIDs))
		for i, id := range *queries.CategoryIDs {
			dbParams.CategoryIds[i] = uuid.MustParse(id)
		}
	}

	products, err := sv.repo.GetProductList(c, dbParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	productCnt, err := sv.repo.CountProducts(c, repository.CountProductsParams{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	productResponses := make([]dto.ProductSummary, 0)
	for _, product := range products {
		productResponses = append(productResponses, dto.MapToShopProductResponse(product))
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, productResponses, dto.CreatePagination(queries.Page, queries.PageSize, productCnt), nil))
}

// Setup product-related routes
func (sv *Server) addProductRoutes(rg *gin.RouterGroup) {
	products := rg.Group("products")
	{
		products.GET("", sv.getProducts)
		products.GET(":id", sv.getProductById)
		products.GET(":id/ratings", sv.getRatingsByProduct)
	}
}
