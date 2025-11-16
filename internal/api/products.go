package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

// @Summary Create a new product
// @Schemes http
// @Description create a new product with the input payload
// @Tags products
// @Accept json
// @Param input body repository.CreateProductTxArgs true "Product input"
// @Produce json
// @Success 200 {object} ApiResponse[string]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /products [post]
func (sv *Server) AddProductHandler(c *gin.Context) {
	var req repository.CreateProductTxArgs
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, "", err))
		return
	}

	productID, err := sv.repo.CreateProductTx(c, req)

	if err != nil {
		log.Error().Err(err).Msg("CreateProduct")
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, "", err))
		return
	}

	c.JSON(http.StatusCreated, createDataResp(c, productID.String(), "", nil, nil))
}

// @Summary Get a product detail by ID
// @Schemes http
// @Description get a product detail by ID
// @Tags products
// @Accept json
// @Param productId path int true "Product ID"
// @Produce json
// @Success 200 {object} ApiResponse[ProductDetailItemResponse]
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /products/{productId} [get]
func (sv *Server) GetProductDetailHandler(c *gin.Context) {
	var params URISlugParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, "", err))
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
			c.JSON(http.StatusNotFound, createErr(NotFoundCode, "", err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, "", err))
		return
	}

	prodID := productRow.ProductID
	productDetail := mapToProductResponse(productRow)

	variantRows, err := sv.repo.GetProductVariants(c, repository.GetProductVariantsParams{
		ProductID: prodID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, "", err))
		return
	}

	idMap := make(map[uuid.UUID]bool)

	for _, row := range variantRows {
		if _, ok := idMap[row.ID]; !ok {
			idMap[row.ID] = true
		}
	}

	// imageResp := mapToProductImages(prodID, images)
	variants := mapToVariantResp(variantRows)
	productDetail.Variants = variants
	// productDetail.ProductImages = imageResp
	c.JSON(http.StatusOK, createDataResp(c, productDetail, "product retrieved", nil, nil))
}

// @Summary Get list of products
// @Schemes http
// @Description get list of products
// @Tags products
// @Accept json
// @Param page query int true "Page number"
// @Param pageSize query int true "Page size"
// @Produce json
// @Success 200 {array} ApiResponse
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /products [get]
func (sv *Server) getProductsHandler(c *gin.Context) {
	var queries ProductQueries
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, "", err))
		return
	}

	dbParams := repository.GetProductsParams{
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
		dbParams.CollectionID = utils.GetPgTypeUUID(uuid.MustParse(*queries.CollectionID))
	}

	if queries.BrandID != nil {
		dbParams.BrandID = utils.GetPgTypeUUID(uuid.MustParse(*queries.BrandID))
	}

	products, err := sv.repo.GetProducts(c, dbParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, "Server error", err))
		return
	}

	productCnt, err := sv.repo.CountProducts(c, repository.CountProductsParams{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, "Server error", err))
		return
	}

	productResponses := make([]ProductListModel, 0)
	for _, product := range products {
		productResponses = append(productResponses, mapToListProductResponse(product))
	}

	c.JSON(http.StatusOK, createDataResp(c, productResponses, "products retrieved", &Pagination{
		Total:           productCnt,
		Page:            queries.Page,
		PageSize:        queries.PageSize,
		TotalPages:      (productCnt + queries.PageSize - 1) / queries.PageSize,
		HasNextPage:     int(queries.Page*queries.PageSize) < int(productCnt),
		HasPreviousPage: queries.Page > 1,
	}, nil))
}

// @Summary Update a product by ID
// @Schemes http
// @Description update a product by ID
// @Tags products
// @Accept json
// @Param productId path int true "Product ID"
// @Param input body repository.UpdateProductTxParams true "Product input"
// @Produce json
// @Success 200 {object} ApiResponse
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrResp
// @Router /products/{productId} [put]
func (sv *Server) updateProductHandler(c *gin.Context) {
	var param UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, "", err))
		return
	}
	uuid, err := uuid.Parse(param.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, "", err))
		return
	}

	var req repository.UpdateProductTxParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, "", err))
		return
	}
	err = sv.repo.UpdateProductTx(c, uuid, req)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErr(NotFoundCode, "", err))
			return
		}
	}

	c.JSON(http.StatusOK, createDataResp(c, struct{}{}, "product updated", nil, nil))
}

// @Summary Remove a product by ID
// @Schemes http
// @Description remove a product by ID
// @Tags products
// @Accept json
// @Param productId path int true "Product ID"
// @Produce json
// @Success 200 {object} ApiResponse
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /products/{productId} [delete]
func (sv *Server) deleteProductHandler(c *gin.Context) {
	var params UriIDParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, "", err))
		return
	}

	product, err := sv.repo.GetProductByID(c, repository.GetProductByIDParams{ID: uuid.MustParse(params.ID)})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErr(NotFoundCode, "", err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, "", err))
		return
	}

	err = sv.repo.DeleteProduct(c, product.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, "", err))
		return
	}

	c.Status(http.StatusNoContent)
}
