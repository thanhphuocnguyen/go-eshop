package api

import (
	"errors"
	"net/http"
	"slices"
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
// @Param input body repository.CreateProductReq true "Product input"
// @Produce json
// @Success 200 {object} ApiResponse[string]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /products [post]
func (sv *Server) AddProductHandler(c *gin.Context) {
	var req CreateProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	createParams := repository.CreateProductParams{
		Name:        req.Name,
		Description: req.Description,
	}

	createParams.BasePrice = utils.GetPgNumericFromFloat(req.BasePrice)
	createParams.ShortDescription = req.ShortDescription
	createParams.Slug = req.Slug
	createParams.BaseSku = req.BaseSku

	createParams.CategoryID = utils.GetPgTypeUUIDFromString(req.CategoryID)

	createParams.BrandID = utils.GetPgTypeUUIDFromString(req.BrandID)
	if req.CollectionID != nil {
		createParams.CollectionID = utils.GetPgTypeUUIDFromString(*req.CollectionID)
	}

	product, err := sv.repo.CreateProduct(c, createParams)

	for _, attrID := range req.Attributes {
		_, err := sv.repo.CreateProductAttribute(c, repository.CreateProductAttributeParams{
			ProductID:   product.ID,
			AttributeID: attrID,
		})
		if err != nil {
			log.Error().Err(err).Msg("CreateProductAttribute")
			c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
			return
		}
	}

	if err != nil {
		log.Error().Err(err).Timestamp().Msg("CreateProduct")

		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusCreated, createDataResp(c, product.ID.String(), nil, nil))
}

// @Summary Get a product detail by ID
// @Schemes http
// @Description get a product detail by ID
// @Tags products
// @Accept json
// @Param productId path int true "Product ID"
// @Produce json
// @Success 200 {object} ApiResponse[ManageProductDetailResp]
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /products/{productId} [get]
func (sv *Server) GetProductByIdHandler(c *gin.Context) {
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
func (sv *Server) GetManageProductsHandler(c *gin.Context) {
	var queries ProductQueries
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	dbParams := repository.GetAdminProductsParams{
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

	products, err := sv.repo.GetAdminProducts(c, dbParams)
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
		productResponses = append(productResponses, mapToListProductResponse(product))
	}

	c.JSON(http.StatusOK, createDataResp(c, productResponses, createPagination(queries.Page, queries.PageSize, productCnt), nil))
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
func (sv *Server) UpdateProductHandler(c *gin.Context) {
	var param UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}
	uuid, err := uuid.Parse(param.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	var req UpdateProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}
	if isStructEmpty(req) {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, errors.New("at least one field must be provided to update")))
		return
	}

	updateParams := repository.UpdateProductParams{
		ID: uuid,
	}
	if req.Name != nil {
		updateParams.Name = req.Name
	}
	if req.Description != nil {
		updateParams.Description = req.Description
	}
	if req.ShortDescription != nil {
		updateParams.ShortDescription = req.ShortDescription
	}
	if req.Slug != nil {
		updateParams.Slug = req.Slug
	}
	if req.BasePrice != nil {
		updateParams.BasePrice = utils.GetPgNumericFromFloat(*req.BasePrice)
	}
	if req.IsActive != nil {
		updateParams.IsActive = req.IsActive
	}
	if req.CategoryID != nil {
		updateParams.CategoryID = utils.GetPgTypeUUIDFromString(*req.CategoryID)
	}
	if req.BrandID != nil {
		updateParams.BrandID = utils.GetPgTypeUUIDFromString(*req.BrandID)
	}
	if req.CollectionID != nil {
		updateParams.CollectionID = utils.GetPgTypeUUIDFromString(*req.CollectionID)
	}

	updated, err := sv.repo.UpdateProduct(c, updateParams)

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErr(NotFoundCode, err))
			return
		}
	}

	if req.Attributes != nil {
		err = sv.repo.DeleteProductAttributesByProductID(c, uuid)
		if err != nil {
			log.Error().Err(err).Msg("DeleteProductAttributesByProductID")
			c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
			return
		}
		prodAttrParams := make([]repository.CreateBulkProductAttributesParams, len(*req.Attributes))
		for i, attrID := range *req.Attributes {
			prodAttrParams[i] = repository.CreateBulkProductAttributesParams{
				ProductID:   uuid,
				AttributeID: attrID,
			}
		}
		_, err = sv.repo.CreateBulkProductAttributes(c, prodAttrParams)
		if err != nil {
			log.Error().Err(err).Msg("CreateBulkProductAttributes")
			c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
			return
		}
	}

	c.JSON(http.StatusOK, createDataResp(c, updated, nil, nil))
}

// @Summary Upload product image
// @Schemes http
// @Description upload product image
// @Tags products
// @Accept json
// @Param input "
// @Produce json
// @Success 200 {object} ApiResponse[string]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /products/{id}/image [post]
func (sv *Server) UploadProductImageHandler(c *gin.Context) {
	var param UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}
	file, _ := c.FormFile("file")
	if file == nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, errors.New("image file is required")))
		return
	}

	id, url, err := sv.uploadService.UploadFile(c, file)
	if err != nil {
		log.Error().Err(err).Timestamp().Msg("UploadFile")

		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	updated, err := sv.repo.UpdateProduct(c, repository.UpdateProductParams{
		ImageUrl: &url,
		ImageID:  &id,
		ID:       uuid.MustParse(param.ID),
	})
	if err != nil {
		log.Error().Err(err).Timestamp().Msg("CreateProduct")

		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusCreated, createDataResp(c, updated, nil, nil))
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
func (sv *Server) DeleteProductHandler(c *gin.Context) {
	var params UriIDParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	product, err := sv.repo.GetProductByID(c, repository.GetProductByIDParams{ID: uuid.MustParse(params.ID)})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErr(NotFoundCode, err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	err = sv.repo.DeleteProduct(c, product.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Create a new product variant
// @Schemes http
// @Description create a new product with the input payload
// @Tags products
// @Accept json
// @Param input body repository.CreateProductReq true "Product input"
// @Produce json
// @Success 200 {object} ApiResponse[string]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /products/{id}/variants [post]
func (sv *Server) AddProductVariantHandler(c *gin.Context) {
	var prodId URISlugParam
	if err := c.ShouldBindUri(&prodId); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}
	var req CreateProdVariantReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}
	if len(req.AttributeValues) == 0 {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, errors.New("attribute values are required")))
		return
	}

	prod, err := sv.repo.GetProductByID(c, repository.GetProductByIDParams{ID: uuid.MustParse(prodId.ID)})
	if err != nil {
		c.JSON(http.StatusNotFound, createErr(NotFoundCode, err))
		return
	}
	prodAttrs, err := sv.repo.GetProductAttributesByProductID(c, prod.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}
	if len(prodAttrs) == 0 {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, errors.New("product has no attributes")))
		return
	}
	prodAttrIds := make([]int32, len(prodAttrs))
	for i, attr := range prodAttrs {
		prodAttrIds[i] = attr.AttributeID
	}

	attributeValues, err := sv.repo.GetAttributeValuesByIDs(c, req.AttributeValues)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}
	for _, attrVal := range attributeValues {
		if !slices.Contains(prodAttrIds, attrVal.AttributeID) {
			c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, errors.New("attribute value does not belong to product attributes")))
			return
		}
	}

	variantSku := repository.GetVariantSKUWithAttributeNames(prod.BaseSku, attributeValues)

	createParams := repository.CreateProductVariantParams{
		ProductID:   prod.ID,
		Description: req.Description,
		Sku:         variantSku,
		Price:       utils.GetPgNumericFromFloat(req.Price),
		Stock:       req.StockQty,
	}
	if req.Weight != nil {
		createParams.Weight = utils.GetPgNumericFromFloat(*req.Weight)
	}

	variant, err := sv.repo.CreateProductVariant(c, createParams)
	if err != nil {
		log.Error().Err(err).Timestamp().Msg("CreateProduct")

		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}
	variantAttrParams := make([]repository.CreateBulkProductVariantAttributeParams, len(req.AttributeValues))

	for i, attrValID := range req.AttributeValues {
		variantAttrParams[i] = repository.CreateBulkProductVariantAttributeParams{
			VariantID:        variant.ID,
			AttributeValueID: attrValID,
		}
	}
	_, err = sv.repo.CreateBulkProductVariantAttribute(c, variantAttrParams)
	if err != nil {
		log.Error().Err(err).Msg("CreateBulkProductVariantAttribute")
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusCreated, createDataResp(c, variant.ID.String(), nil, nil))
}

// @Summary Upload product image
// @Schemes http
// @Description upload product image
// @Tags products
// @Accept json
// @Param input "
// @Produce json
// @Success 200 {object} ApiResponse[string]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /products/{id}/variants/{variantId}/image [post]
func (sv *Server) UploadProductVariantImageHandler(c *gin.Context) {
	var param UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}
	var variantParam ProductVariantParam
	if err := c.ShouldBindUri(&variantParam); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}
	file, _ := c.FormFile("file")
	if file == nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, errors.New("image file is required")))
		return
	}

	product, err := sv.repo.GetProductByID(c, repository.GetProductByIDParams{ID: uuid.MustParse(param.ID)})
	if err != nil {
		c.JSON(http.StatusNotFound, createErr(NotFoundCode, err))
		return
	}
	variant, err := sv.repo.GetProductVariantByID(c, repository.GetProductVariantByIDParams{
		ID:        uuid.MustParse(variantParam.ID),
		ProductID: product.ID,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, createErr(NotFoundCode, err))
		return
	}

	id, url, err := sv.uploadService.UploadFile(c, file)
	if err != nil {
		log.Error().Err(err).Timestamp().Msg("UploadFile")

		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	updated, err := sv.repo.UpdateProductVariant(c, repository.UpdateProductVariantParams{
		ImageUrl:  &url,
		ImageID:   &id,
		ID:        variant.ID,
		ProductID: product.ID,
	})

	if err != nil {
		log.Error().Err(err).Timestamp().Msg("UpdateProductVariant")
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusCreated, createDataResp(c, updated, nil, nil))
}

// @Summary Get product variants
// @Schemes http
// @Description get product variants
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} ApiResponse[[]ProductVariantModel]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /products/{id}/variants [get]
func (sv *Server) GetVariantsHandler(c *gin.Context) {
	var prodId URISlugParam
	if err := c.ShouldBindUri(&prodId); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	prod, err := sv.repo.GetProductByID(c, repository.GetProductByIDParams{ID: uuid.MustParse(prodId.ID)})
	if err != nil {
		c.JSON(http.StatusNotFound, createErr(NotFoundCode, err))
		return
	}
	variantRows, err := sv.repo.GetProductVariantList(c, repository.GetProductVariantListParams{ProductID: prod.ID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusOK, createDataResp(c, variantRows, nil, nil))
}

// @Summary Get product variant
// @Schemes http
// @Description get product variant
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} ApiResponse[ProductVariantModel]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /products/{id}/variants/{variantID} [get]
func (sv *Server) GetVariantHandler(c *gin.Context) {
	var prodId URISlugParam
	if err := c.ShouldBindUri(&prodId); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}
	var variantId ProductVariantParam
	if err := c.ShouldBindUri(&variantId); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	prod, err := sv.repo.GetProductByID(c, repository.GetProductByIDParams{ID: uuid.MustParse(prodId.ID)})
	if err != nil {
		c.JSON(http.StatusNotFound, createErr(NotFoundCode, err))
		return
	}

	rows, err := sv.repo.GetVariantDetailByID(c, repository.GetVariantDetailByIDParams{
		ID:        uuid.MustParse(variantId.ID),
		ProductID: prod.ID,
	})

	if err != nil {
		c.JSON(http.StatusNotFound, createErr(NotFoundCode, err))
		return
	}
	first := rows[0]
	price, _ := first.Price.Float64Value()
	resp := VariantModel{
		ID:         first.ID.String(),
		Price:      price.Float64,
		StockQty:   first.Stock,
		Sku:        &first.Sku,
		Attributes: make([]AttributeValue, len(rows)),
		CreatedAt:  first.CreatedAt.String(),
		UpdatedAt:  first.UpdatedAt.String(),
		IsActive:   *first.IsActive,
	}
	if first.Weight.Valid {
		weightFloat, _ := first.Weight.Float64Value()
		resp.Weight = &weightFloat.Float64
	}
	for i, row := range rows {
		attr := AttributeValue{
			ID:    *row.AttributeValueID,
			Value: *row.AttributeValue,
		}
		resp.Attributes[i] = attr
	}

	c.JSON(http.StatusOK, createDataResp(c, resp, nil, nil))
}

// @Summary Update a product variant
// @Schemes http
// @Description update a product with the input payload
// @Tags products
// @Accept json
// @Param input body repository.CreateProductReq true "Product input"
// @Produce json
// @Success 200 {object} ApiResponse[string]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /products/{id}/variants/{variantId} [put]
func (sv *Server) UpdateProductVariantHandler(c *gin.Context) {
	var uris URIVariantParam
	if err := c.ShouldBindUri(&uris); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}
	var req UpdateProdVariantReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	prod, err := sv.repo.GetProductByID(c, repository.GetProductByIDParams{ID: uuid.MustParse(uris.ProductID)})
	if err != nil {
		c.JSON(http.StatusNotFound, createErr(NotFoundCode, err))
		return
	}

	updateParams := repository.UpdateProductVariantParams{
		ProductID: prod.ID,
		ID:        uuid.MustParse(uris.VariantID),
	}

	if req.Price != nil {
		updateParams.Price = utils.GetPgNumericFromFloat(*req.Price)
	}
	if req.StockQty != nil {
		updateParams.Stock = req.StockQty
	}
	if req.Weight != nil {
		updateParams.Weight = utils.GetPgNumericFromFloat(*req.Weight)
	}
	if req.Description != nil {
		updateParams.Description = req.Description
	}

	updatedVariant, err := sv.repo.UpdateProductVariant(c, updateParams)
	if err != nil {
		log.Error().Err(err).Timestamp().Msg("UpdateProductVariant")

		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusOK, createDataResp(c, updatedVariant.ID.String(), nil, nil))
}

// @Summary Delete a product variant
// @Schemes http
// @Description delete a product variant with the input payload
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} ApiResponse[string]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /products/{id}/variant/{variantID} [delete]
func (sv *Server) DeleteProductVariantHandler(c *gin.Context) {
	var uris URIVariantParam
	if err := c.ShouldBindUri(&uris); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	prod, err := sv.repo.GetProductByID(c, repository.GetProductByIDParams{ID: uuid.MustParse(uris.ProductID)})
	if err != nil {
		c.JSON(http.StatusNotFound, createErr(NotFoundCode, err))
		return
	}

	err = sv.repo.DeleteProductVariant(c, repository.DeleteProductVariantParams{
		ProductID: prod.ID,
		ID:        uuid.MustParse(uris.VariantID),
	})
	if err != nil {
		log.Error().Err(err).Timestamp().Msg("DeleteProductVariant")

		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}
	c.JSON(http.StatusNoContent, createDataResp(c, struct{}{}, nil, nil))
}
