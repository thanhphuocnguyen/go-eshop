package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
)

// @Summary Get a display product by ID
// @Schemes http
// @Description get a display product by ID
// @Tags products
// @Accept json
// @Param productId path int true "Product ID"
// @Produce json
// @Success 200 {object} ApiResponse[ManageProductDetailResp]
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /products/{productId} [get]
func (sv *Server) GetDisplayProductByIdHandler(c *gin.Context) {
	var params URISlugParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
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
			c.JSON(http.StatusNotFound, createErr(NotFoundCode, err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	productDetail := mapToProductResponse(productRow)
	prodAttr, err := sv.repo.GetProductAttributesByProductID(c, productRow.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}
	attrs := make([]int32, len(prodAttr))
	for i, attr := range prodAttr {
		attrs[i] = attr.AttributeID
	}
	productDetail.Attributes = attrs
	c.JSON(http.StatusOK, createDataResp(c, productDetail, nil, nil))
}

// @Summary Get list of display products
// @Schemes http
// @Description get list of products
// @Tags products
// @Accept json
// @Param page query int true "Page number"
// @Param pageSize query int true "Page size"
// @Produce json
// @Success 200 {array} ApiResponse[[]ManageProductListModel]
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /products [get]
func (sv *Server) GetDisplayProductsHandler(c *gin.Context) {
	var queries ProductQueries
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	dbParams := repository.GetDisplayProductsParams{
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

	if len(queries.CategoryIDs) > 0 {
		dbParams.CategoryIds = make([]uuid.UUID, len(queries.CategoryIDs))
		for i, id := range queries.CategoryIDs {
			dbParams.CategoryIds[i] = uuid.MustParse(id)
		}
	}

	if queries.CollectionID != nil {
		dbParams.CollectionIds = []uuid.UUID{uuid.MustParse(*queries.CollectionID)}
	}

	if queries.BrandID != nil {
		dbParams.BrandIds = []uuid.UUID{uuid.MustParse(*queries.BrandID)}
	}

	products, err := sv.repo.GetDisplayProducts(c, dbParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	productCnt, err := sv.repo.CountProducts(c, repository.CountProductsParams{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	productResponses := make([]ManageProductListModel, 0)
	for _, product := range products {
		price, _ := product.MinPrice.Float64Value()
		prod := ManageProductListModel{
			ID:          product.ID.String(),
			Name:        product.Name,
			Description: product.Description,
			BasePrice:   price.Float64,
			Sku:         product.BaseSku,
			Slug:        product.Slug,
			ImageUrl:    product.ImageUrl,
		}
		if product.AvgRating.Valid {
			avgRating, _ := product.AvgRating.Float64Value()
			prod.AvgRating = &avgRating.Float64
		}
		productResponses = append(productResponses, prod)
	}

	c.JSON(http.StatusOK, createDataResp(c, productResponses, createPagination(queries.Page, queries.PageSize, productCnt), nil))
}
