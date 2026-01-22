package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"github.com/thanhphuocnguyen/go-eshop/internal/worker"
)

// adminGetProducts godoc
// @Summary Get admin list of products
// @Schemes http
// @Description get admin list of products
// @Tags products
// @Accept json
// @Param page query int true "Page number"
// @Param pageSize query int true "Page size"
// @Produce json
// @Success 200 {array} dto.ApiResponse[[]dto.ProductSummary]
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/products [get]
func (s *Server) adminGetProducts(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	// Parse query parameters
	var queries models.ProductQuery
	queries.Page = 1
	queries.PageSize = 10

	if page := r.URL.Query().Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			queries.Page = int64(p)
		}
	}
	if pageSize := r.URL.Query().Get("pageSize"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil {
			queries.PageSize = int64(ps)
		}
	}
	if search := r.URL.Query().Get("search"); search != "" {
		queries.Search = &search
	}

	dbParams := repository.GetAdminProductListParams{
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

	products, err := s.repo.GetAdminProductList(c, dbParams)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	productCnt, err := s.repo.CountProducts(c, repository.CountProductsParams{})
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	productResponses := make([]dto.ProductListItem, 0)
	for _, product := range products {
		productResponses = append(productResponses, dto.MapToAdminProductResponse(product))
	}

	pagination := dto.CreatePagination(queries.Page, queries.PageSize, productCnt)
	RespondSuccessWithPagination(w, productResponses, pagination)
}

// @Summary Create a new product
// @Schemes http
// @Description create a new product with the input payload
// @Tags products
// @Accept json
// @Param input body models.CreateProductModel true "Product input"
// @Produce json
// @Success 200 {object} dto.ApiResponse[repository.Product]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /products [post]
func (s *Server) createProduct(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	var req models.CreateProductModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	createParams := repository.CreateProductParams{
		Name:               req.Name,
		Description:        req.Description,
		DiscountPercentage: req.DiscountPercentage,
	}

	createParams.BasePrice = utils.GetPgNumericFromFloat(req.BasePrice)
	createParams.ShortDescription = req.ShortDescription
	createParams.Slug = req.Slug
	createParams.BaseSku = req.BaseSku
	createParams.BrandID = utils.GetPgTypeUUIDFromString(req.BrandID)

	// Convert variants if provided
	variants := make([]repository.CreateProductVariantTxParams, 0, len(req.Variants))
	for _, variantReq := range req.Variants {
		// Convert int64 to int32 for AttributeValues
		attributeValues := make([]int64, len(variantReq.AttributeValues))
		copy(attributeValues, variantReq.AttributeValues)

		variant := repository.CreateProductVariantTxParams{
			Description:     variantReq.Description,
			Price:           variantReq.Price,
			Stock:           variantReq.StockQty,
			AttributeValues: attributeValues,
		}

		if variantReq.Description != "" {
			variant.Description = variantReq.Description
		}
		if variantReq.Weight != nil {
			variant.Weight = variantReq.Weight
		}

		variants = append(variants, variant)
	}

	// Use transaction to ensure all operations succeed or fail together
	txArgs := repository.CreateProductTxArgs{
		Product:       createParams,
		Attributes:    req.Attributes,
		CategoryIDs:   req.CategoryIDs,
		CollectionIDs: req.CollectionIDs,
		Variants:      variants,
	}

	product, err := s.repo.CreateProductTx(c, txArgs)
	if err != nil {
		log.Error().Err(err).Timestamp().Msg("CreateProductTx")
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	err = s.taskDistributor.SendIndexProductTask(c, worker.PayloadIndexProduct{
		ProductID: product.ID,
	})
	if err != nil {
		log.Error().Err(err).Msg("SendIndexProductTask")
	}

	RespondCreated(w, product)
}

// @Summary Update a product by ID
// @Schemes http
// @Description update a product by ID
// @Tags products
// @Accept json
// @Param productId path int true "Product ID"
// @Param input body models.UpdateProductModel true "Product update input"
// @Produce json
// @Success 200 {object} dto.ApiResponse[repository.Product]
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/products/{productId} [put]
func (s *Server) updateProduct(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id, err := GetUrlParam(r, "id")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	productID, err := uuid.Parse(id)

	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	var req models.UpdateProductModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	if dto.IsStructEmpty(req) {
		RespondBadRequest(w, InvalidBodyCode, errors.New("at least one field must be provided to update"))
		return
	}

	updateParams := repository.UpdateProductParams{
		ID:                 productID,
		DiscountPercentage: req.DiscountPercentage,
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

	if req.BrandID != nil {
		updateParams.BrandID = utils.GetPgTypeUUIDFromString(*req.BrandID)
	}

	// Use transaction to ensure all operations succeed or fail together
	txArgs := repository.UpdateProductTxArgs{
		Product:       updateParams,
		Attributes:    req.Attributes,
		CategoryIDs:   req.CategoryIDs,
		CollectionIDs: req.CollectionIDs,
	}

	updated, err := s.repo.UpdateProductTx(c, txArgs)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, err)
			return
		}
		log.Error().Err(err).Timestamp().Msg("UpdateProductTx")
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	err = s.taskDistributor.SendUpdateIndexProductTask(c, worker.PayloadUpdateProductIndex{
		ProductID: updated.ID,
		Data:      req,
	})
	if err != nil {
		log.Error().Err(err).Msg("SendUpdateIndexProductTask")
	}

	RespondSuccess(w, updated)
}

// @Summary Remove a product by ID
// @Schemes http
// @Description remove a product by ID
// @Tags products
// @Accept json
// @Param productId path int true "Product ID"
// @Produce json
// @Success 200 {object} dto.ApiResponse[string]
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/products/{productId} [delete]
func (s *Server) deleteProduct(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id, err := GetUrlParam(r, "id")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	product, err := s.repo.GetProductByID(c, repository.GetProductByIDParams{ID: uuid.MustParse(id)})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, err)
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	err = s.repo.DeleteProduct(c, product.ID)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	err = s.taskDistributor.SendDeleteIndexProductTask(c, worker.PayloadIndexProduct{
		ProductID: product.ID,
	})
	if err != nil {
		log.Error().Err(err).Msg("SendDeleteIndexProductTask")
	}
	RespondNoContent(w)
}

// @Summary Upload product image
// @Schemes http
// @Description upload product image
// @Tags products
// @Accept multipart/form-data
// @Param id path string true "Product ID"
// @Produce json
// @Success 200 {object} dto.ApiResponse[string]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/products/{id}/image [post]
func (s *Server) uploadProductImage(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id, err := GetUrlParam(r, "id")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, errors.New("image file is required"))
		return
	}
	defer file.Close()

	prod, err := s.repo.GetProductByID(c, repository.GetProductByIDParams{ID: uuid.MustParse(id)})
	if err != nil {
		RespondNotFound(w, NotFoundCode, err)
		return
	}
	if prod.ImageID != nil {
		msg, err := s.uploadService.Remove(c, *prod.ImageID)
		if err != nil {
			log.Error().Err(err).Timestamp().Msg(msg)
			RespondInternalServerError(w, InternalServerErrorCode, err)
			return
		}
		return
	}

	fileHeader := &struct {
		Filename string
		Header   map[string][]string
		Size     int64
	}{
		Filename: header.Filename,
		Header:   header.Header,
		Size:     header.Size,
	}

	uploadID, url, err := s.uploadService.Upload(c, fileHeader)
	if err != nil {
		log.Error().Err(err).Timestamp().Msg("UploadFile")
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	updated, err := s.repo.UpdateProduct(c, repository.UpdateProductParams{
		ImageUrl: &url,
		ImageID:  &uploadID,
		ID:       uuid.MustParse(id),
	})
	if err != nil {
		log.Error().Err(err).Timestamp().Msg("CreateProduct")
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	RespondSuccess(w, updated)
}

// @Summary Create a new product variant
// @Schemes http
// @Description create a new product with the input payload
// @Tags products
// @Accept json
// @Param input body models.CreateProdVariantModel true "Product variant input"
// @Produce json
// @Success 200 {object} dto.ApiResponse[string]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/products/{id}/variants [post]
func (s *Server) createVariant(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id := chi.URLParam(r, "id")
	if id == "" {
		RespondBadRequest(w, InvalidBodyCode, errors.New("id parameter is required"))
		return
	}

	var req models.CreateProdVariantModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	if len(req.AttributeValues) == 0 {
		RespondBadRequest(w, InvalidBodyCode, errors.New("attribute values are required"))
		return
	}

	prod, err := s.repo.GetProductByID(c, repository.GetProductByIDParams{ID: uuid.MustParse(id)})
	if err != nil {
		RespondNotFound(w, NotFoundCode, err)
		return
	}

	prodAttrs, err := s.repo.GetProductAttributesByProductID(c, prod.ID)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	if len(prodAttrs) == 0 {
		RespondBadRequest(w, InvalidBodyCode, errors.New("product has no attributes"))
		return
	}

	prodAttrIds := make([]int32, len(prodAttrs))
	for i, attr := range prodAttrs {
		prodAttrIds[i] = attr.AttributeID
	}

	attributeValues, err := s.repo.GetAttributeValuesByIDs(c, req.AttributeValues)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	if len(attributeValues) != len(prodAttrIds) {
		RespondBadRequest(w, InvalidBodyCode, errors.New("attribute values do not match product attributes"))
		return
	}

	for _, attrVal := range attributeValues {
		if !slices.Contains(prodAttrIds, attrVal.AttributeID) {
			RespondBadRequest(w, InvalidBodyCode, errors.New("attribute value does not belong to product attributes"))
			return
		}
	}

	variantSku := repository.GetVariantSKUWithAttributeNames(prod.BaseSku, attributeValues)

	createParams := repository.CreateProductVariantParams{
		ProductID:   prod.ID,
		Description: &req.Description,
		Sku:         variantSku,
		Price:       utils.GetPgNumericFromFloat(req.Price),
		Stock:       req.StockQty,
	}
	if req.Weight != nil {
		createParams.Weight = utils.GetPgNumericFromFloat(*req.Weight)
	}

	// Use transaction to ensure all operations succeed or fail together
	txArgs := repository.CreateProductVariantTxArgs{
		Variant:         createParams,
		AttributeValues: req.AttributeValues,
	}

	variant, err := s.repo.CreateProductVariantTx(c, txArgs)
	if err != nil {
		log.Error().Err(err).Timestamp().Msg("CreateProductVariantTx")
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	err = s.taskDistributor.SendIndexProductTask(c, worker.PayloadIndexProduct{
		ProductID: prod.ID,
	})
	if err != nil {
		log.Error().Err(err).Msg("SendIndexProductTask")
	}

	RespondCreated(w, variant)
}

// @Summary Get product variants
// @Schemes http
// @Description get product variants
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} dto.ApiResponse[[]dto.VariantDetail]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/products/{id}/variants [get]
func (s *Server) getProductVariants(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id := chi.URLParam(r, "id")
	if id == "" {
		RespondBadRequest(w, InvalidBodyCode, errors.New("id parameter is required"))
		return
	}

	prod, err := s.repo.GetProductByID(c, repository.GetProductByIDParams{ID: uuid.MustParse(id)})
	if err != nil {
		RespondNotFound(w, NotFoundCode, err)
		return
	}
	variantRows, err := s.repo.GetProductVariantList(c, repository.GetProductVariantListParams{ProductID: prod.ID})
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	resp := make([]dto.VariantDetail, len(variantRows))
	for i, row := range variantRows {
		resp[i] = dto.MapToVariantListModelDto(row)
	}

	RespondSuccess(w, resp)
}

// @Summary Get product variant
// @Schemes http
// @Description get product variant
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param variantID path string true "Product Variant ID"
// @Success 200 {object} dto.ApiResponse[dto.VariantDetail]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/products/{id}/variants/{variantID} [get]
func (s *Server) getVariantByProductId(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id := chi.URLParam(r, "id")
	variantId := chi.URLParam(r, "variantId")
	if id == "" || variantId == "" {
		RespondBadRequest(w, InvalidBodyCode, errors.New("id and variantId parameters are required"))
		return
	}

	prod, err := s.repo.GetProductByID(c, repository.GetProductByIDParams{ID: uuid.MustParse(id)})
	if err != nil {
		RespondNotFound(w, NotFoundCode, err)
		return
	}

	rows, err := s.repo.GetVariantDetailByID(c, repository.GetVariantDetailByIDParams{
		ID:        uuid.MustParse(variantId),
		ProductID: prod.ID,
	})

	if err != nil {
		RespondNotFound(w, NotFoundCode, err)
		return
	}
	first := rows[0]
	price, _ := first.Price.Float64Value()
	resp := dto.VariantDetail{
		ID:         first.ID.String(),
		Price:      price.Float64,
		Stock:      first.Stock,
		Sku:        first.Sku,
		Attributes: make([]dto.AttributeValueDetail, len(rows)),
		CreatedAt:  first.CreatedAt.String(),
		UpdatedAt:  first.UpdatedAt.String(),
		IsActive:   *first.IsActive,
	}
	if first.Weight.Valid {
		weightFloat, _ := first.Weight.Float64Value()
		resp.Weight = &weightFloat.Float64
	}
	if first.ImageUrl != nil {
		resp.ImageUrl = first.ImageUrl
	}
	if first.ImageID != nil {
		resp.ImageID = first.ImageID
	}
	for i, row := range rows {
		attr := dto.AttributeValueDetail{
			ID:    *row.AttributeValueID,
			Value: *row.AttributeValue,
		}
		resp.Attributes[i] = attr
	}

	RespondSuccess(w, resp)
}

// @Summary Update a product variant
// @Schemes http
// @Description update a product with the input payload
// @Tags products
// @Accept json
// @Param input body models.UpdateProdVariantModel true "Product variant input"
// @Produce json
// @Success 200 {object} dto.ApiResponse[repository.ProductVariant]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/products/{id}/variants/{variantId} [put]
func (s *Server) updateVariant(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id := chi.URLParam(r, "id")
	variantId := chi.URLParam(r, "variantId")
	if id == "" || variantId == "" {
		RespondBadRequest(w, InvalidBodyCode, errors.New("id and variantId parameters are required"))
		return
	}

	var req models.UpdateProdVariantModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	prod, err := s.repo.GetProductByID(c, repository.GetProductByIDParams{ID: uuid.MustParse(id)})
	if err != nil {
		RespondNotFound(w, NotFoundCode, err)
		return
	}

	updateParams := repository.UpdateProductVariantParams{
		ProductID: prod.ID,
		ID:        uuid.MustParse(variantId),
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

	updatedVariant, err := s.repo.UpdateProductVariant(c, updateParams)
	if err != nil {
		log.Error().Err(err).Timestamp().Msg("UpdateProductVariant")
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondSuccess(w, updatedVariant)
}

// @Summary Upload a product variant image
// @Schemes http
// @Description upload a product variant image with the input payload
// @Tags products
// @Accept multipart/form-data
// @Param id path string true "Product ID"
// @Param variantId path string true "Product Variant ID"
// @Produce json
// @Success 200 {object} dto.ApiResponse[repository.ProductVariant]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/products/{id}/variants/{variantId}/images [post]
func (s *Server) uploadVariantImage(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id := chi.URLParam(r, "id")
	variantId := chi.URLParam(r, "variantId")
	if id == "" || variantId == "" {
		RespondBadRequest(w, InvalidBodyCode, errors.New("id and variantId parameters are required"))
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10MB max
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	defer file.Close()

	prod, err := s.repo.GetProductByID(c, repository.GetProductByIDParams{ID: uuid.MustParse(id)})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dto.CreateErr(NotFoundCode, err))
		return
	}

	variant, err := s.repo.GetProductVariantByID(c, repository.GetProductVariantByIDParams{
		ID:        uuid.MustParse(variantId),
		ProductID: prod.ID,
	})

	if err != nil {
		RespondNotFound(w, NotFoundCode, err)
		return
	}

	if variant.ImageID != nil {
		msg, err := s.uploadService.Remove(c, *variant.ImageID)
		if err != nil {
			log.Error().Err(err).Timestamp().Msg(msg)
			RespondInternalServerError(w, InternalServerErrorCode, err)
			return
		}
	}

	fileHeader := &struct {
		Filename string
		Header   map[string][]string
		Size     int64
	}{
		Filename: header.Filename,
		Header:   header.Header,
		Size:     header.Size,
	}

	uploadID, url, err := s.uploadService.Upload(c, fileHeader)
	if err != nil {
		log.Error().Err(err).Timestamp().Msg("UploadFile")
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	updateParam := repository.UpdateProductVariantParams{
		ProductID: prod.ID,
		ID:        variant.ID,
		ImageID:   &uploadID,
		ImageUrl:  &url,
	}

	updatedVariant, err := s.repo.UpdateProductVariant(c, updateParam)
	if err != nil {
		log.Error().Err(err).Timestamp().Msg("UpdateProductVariant")
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondSuccess(w, updatedVariant)
}

// @Summary Delete a product variant
// @Schemes http
// @Description delete a product variant with the input payload
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} dto.ApiResponse[string]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/products/{id}/variant/{variantID} [delete]
func (s *Server) deleteVariant(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id := chi.URLParam(r, "id")
	variantId := chi.URLParam(r, "variantId")
	if id == "" || variantId == "" {
		RespondBadRequest(w, InvalidBodyCode, errors.New("id and variantId parameters are required"))
		return
	}

	prod, err := s.repo.GetProductByID(c, repository.GetProductByIDParams{ID: uuid.MustParse(id)})
	if err != nil {
		RespondNotFound(w, NotFoundCode, err)
		return
	}

	err = s.repo.DeleteProductVariant(c, repository.DeleteProductVariantParams{
		ProductID: prod.ID,
		ID:        uuid.MustParse(variantId),
	})

	if err != nil {
		log.Error().Err(err).Timestamp().Msg("DeleteVariant")
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	RespondNoContent(w)
}
