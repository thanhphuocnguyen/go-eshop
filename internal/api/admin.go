package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/constants"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
	"github.com/thanhphuocnguyen/go-eshop/pkg/payment"
)

// parseQuery is a utility function to parse URL query parameters into a struct
func parseQuery(r *http.Request, dest interface{}) error {
	validate := validator.New()
	return validate.Struct(dest)
}

// Setup admin-related routes
func (sv *Server) addAdminRoutes(r chi.Router) {
	r.Route("/admin", func(r chi.Router) {
		// Apply authentication and authorization middleware
		r.Use(authenticateMiddleware(nil, sv.tokenGenerator))
		r.Use(authorizeMiddleware(nil, "admin"))

		// User routes
		r.Route("/users", func(r chi.Router) {
			r.Get("/", sv.adminGetUsers)
			r.Get("/{id}", sv.adminGetUser)
		})

		// Product routes
		r.Route("/products", func(r chi.Router) {
			r.Get("/", sv.adminGetProducts)
			r.Post("/", sv.adminAddProduct)
			r.Put("/{id}", sv.adminUpdateProduct)
			r.Delete("/{id}", sv.adminDeleteProduct)

			r.Route("/{id}", func(r chi.Router) {
				r.Post("/images", sv.adminUploadProductImage)

				r.Route("/variants", func(r chi.Router) {
					r.Post("/", sv.adminAddVariant)
					r.Get("/", sv.adminGetVariants)
					r.Get("/{variantId}", sv.adminGetVariant)
					r.Put("/{variantId}", sv.adminUpdateVariant)
					r.Post("/{variantId}/images", sv.adminUploadVariantImage)
					r.Delete("/{variantId}", sv.adminDeleteVariant)
				})
			})
		})

		// Attribute routes
		r.Route("/attributes", func(r chi.Router) {
			r.Post("/", sv.adminCreateAttribute)
			r.Get("/", sv.adminGetAttributes)
			r.Get("/{id}", sv.adminGetAttributeByID)
			r.Put("/{id}", sv.adminUpdateAttribute)
			r.Delete("/{id}", sv.adminRemoveAttribute)

			r.Get("/product/{id}", sv.adminGetAttributeValuesForProduct)

			r.Route("/{id}", func(r chi.Router) {
				r.Post("/create", sv.adminAddAttributeValue)
				r.Put("/update/{valueId}", sv.adminUpdateAttrValue)
				r.Delete("/remove/{valueId}", sv.adminRemoveAttrValue)
			})
		})

		// Order routes
		r.Route("/orders", func(r chi.Router) {
			r.Get("/", sv.adminGetOrders)
			r.Get("/{id}", sv.adminGetOrderDetail)
			r.Put("/{id}/status", sv.adminChangeOrderStatus)
			r.Post("/{id}/cancel", sv.adminCancelOrder)
			r.Post("/{id}/refund", sv.adminRefundOrder)
		})

		// Category routes
		r.Route("/categories", func(r chi.Router) {
			r.Get("/", sv.adminGetCategories)
			r.Get("/{id}", sv.adminGetCategoryByID)
			r.Post("/", sv.adminCreateCategory)
			r.Put("/{id}", sv.adminUpdateCategory)
			r.Delete("/{id}", sv.adminDeleteCategory)
		})

		// Brand routes
		r.Route("/brands", func(r chi.Router) {
			r.Get("/", sv.adminGetBrands)
			r.Get("/{id}", sv.adminGetBrandByID)
			r.Post("/", sv.adminCreateBrand)
			r.Put("/{id}", sv.adminUpdateBrand)
			r.Delete("/{id}", sv.adminDeleteBrand)
		})

		// Collection routes
		r.Route("/collections", func(r chi.Router) {
			r.Get("/", sv.adminGetCollections)
			r.Post("/", sv.adminCreateCollection)
			r.Get("/{id}", sv.adminGetCollectionByID)
			r.Put("/{id}", sv.adminUpdateCollection)
			r.Delete("/{id}", sv.adminDeleteCollection)
		})

		// Rating routes
		r.Route("/ratings", func(r chi.Router) {
			r.Get("/", sv.adminGetRatings)
			r.Delete("/{id}", sv.adminDeleteRating)
			r.Put("/{id}/approve", sv.adminApproveRating)
			r.Put("/{id}/ban", sv.adminBanUserRating)
		})

		// Discount routes
		r.Route("/discounts", func(r chi.Router) {
			r.Post("/", sv.adminCreateDiscount)
			r.Get("/", sv.adminGetDiscounts)
			r.Get("/{id}", sv.getDiscountByID)
			r.Put("/{id}", sv.adminUpdateDiscount)
			r.Delete("/{id}", sv.adminDeleteDiscount)

			r.Route("/{id}/rules", func(r chi.Router) {
				r.Post("/", sv.adminAddDiscountRule)
				r.Get("/", sv.adminGetDiscountRules)
				r.Get("/{ruleId}", sv.adminGetDiscountRuleByID)
				r.Put("/{ruleId}", sv.adminUpdateDiscountRule)
				r.Delete("/{ruleId}", sv.adminDeleteDiscountRule)
			})
		})
	})
}

// adminGetUsers godoc
// @Summary List users
// @Description List users
// @Tags users
// @Accept  json
// @Produce  json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} ApiResponse[[]UserDetail]
// @Failure 500 {object} ErrorResp
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Router /admin/users [get]
func (sv *Server) adminGetUsers(w http.ResponseWriter, r *http.Request) {
	authPayload, ok := r.Context().Value("auth").(*auth.TokenPayload)
	if !ok {
		RespondInternalServerError(w, InternalServerErrorCode, errors.New("authorization payload is not provided"))
		return
	}

	// Parse query parameters
	queries := ParsePaginationQuery(r)

	users, err := sv.repo.GetUsers(r.Context(), repository.GetUsersParams{
		Limit:  queries.PageSize,
		Offset: (queries.Page - 1) * queries.PageSize,
	})

	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	total, err := sv.repo.CountUsers(r.Context())
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	userResp := make([]dto.UserDetail, 0)
	for _, user := range users {
		userResp = append(userResp, dto.MapToUserResponse(user, authPayload.RoleCode))
	}

	pagination := dto.CreatePagination(queries.Page, queries.PageSize, total)
	RespondSuccessWithPagination(w, r, userResp, pagination)
}

// adminGetUser godoc
// @Summary Get user info
// @Description Get user info
// @Tags admin
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} ApiResponse[UserDetail]
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/users/{id} [get]
func (sv *Server) adminGetUser(w http.ResponseWriter, r *http.Request) {
	authPayload, ok := r.Context().Value("auth").(*auth.TokenPayload)
	if !ok {
		RespondInternalServerError(w, InternalServerErrorCode, errors.New("authorization payload is not provided"))
		return
	}

	id, err := GetURLParam(r, "id")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	user, err := sv.repo.GetUserByID(r.Context(), uuid.MustParse(id))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, err)
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	userResp := dto.MapToUserResponse(user, authPayload.RoleCode)
	RespondSuccess(w, r, userResp)
}

// adminGetProducts godoc
// @Summary Get admin list of products
// @Schemes http
// @Description get admin list of products
// @Tags products
// @Accept json
// @Param page query int true "Page number"
// @Param pageSize query int true "Page size"
// @Produce json
// @Success 200 {array} ApiResponse[[]ProductSummary]
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/products [get]
func (sv *Server) adminGetProducts(w http.ResponseWriter, r *http.Request) {
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

	products, err := sv.repo.GetAdminProductList(r.Context(), dbParams)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	productCnt, err := sv.repo.CountProducts(r.Context(), repository.CountProductsParams{})
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	productResponses := make([]dto.ProductListItem, 0)
	for _, product := range products {
		productResponses = append(productResponses, dto.MapToAdminProductResponse(product))
	}

	pagination := dto.CreatePagination(queries.Page, queries.PageSize, productCnt)
	RespondSuccessWithPagination(w, r, productResponses, pagination)
}

// @Summary Create a new product
// @Schemes http
// @Description create a new product with the input payload
// @Tags products
// @Accept json
// @Param input body CreateProductReq true "Product input"
// @Produce json
// @Success 200 {object} ApiResponse[repository.Product]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /products [post]
func (sv *Server) adminAddProduct(w http.ResponseWriter, r *http.Request) {
	var req models.CreateProductModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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

	// Use transaction to ensure all operations succeed or fail together
	txArgs := repository.CreateProductTxArgs{
		Product:       createParams,
		Attributes:    req.Attributes,
		CategoryIDs:   req.CategoryIDs,
		CollectionIDs: req.CollectionIDs,
	}

	product, err := sv.repo.CreateProductTx(r.Context(), txArgs)
	if err != nil {
		log.Error().Err(err).Timestamp().Msg("CreateProductTx")
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondCreated(w, r, product)
}

// @Summary Update a product by ID
// @Schemes http
// @Description update a product by ID
// @Tags products
// @Accept json
// @Param productId path int true "Product ID"
// @Param input body UpdateProductReq true "Product update input"
// @Produce json
// @Success 200 {object} ApiResponse[repository.Product]
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/products/{productId} [put]
func (sv *Server) adminUpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := GetURLParam(r, "id")
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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

	updated, err := sv.repo.UpdateProductTx(r.Context(), txArgs)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, err)
			return
		}
		log.Error().Err(err).Timestamp().Msg("UpdateProductTx")
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondSuccess(w, r, updated)
}

// @Summary Remove a product by ID
// @Schemes http
// @Description remove a product by ID
// @Tags products
// @Accept json
// @Param productId path int true "Product ID"
// @Produce json
// @Success 200 {object} ApiResponse[string]
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/products/{productId} [delete]
func (sv *Server) adminDeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := GetURLParam(r, "id")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	product, err := sv.repo.GetProductByID(r.Context(), repository.GetProductByIDParams{ID: uuid.MustParse(id)})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if errors.Is(err, repository.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.CreateErr(NotFoundCode, err))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	err = sv.repo.DeleteProduct(r.Context(), product.ID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Upload product image
// @Schemes http
// @Description upload product image
// @Tags products
// @Accept multipart/form-data
// @Param id path string true "Product ID"
// @Produce json
// @Success 200 {object} ApiResponse[string]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/products/{id}/image [post]
func (sv *Server) adminUploadProductImage(w http.ResponseWriter, r *http.Request) {
	id, err := GetURLParam(r, "id")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	// Parse multipart form
	err = r.ParseMultipartForm(10 << 20) // 10MB max
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.CreateErr(InvalidBodyCode, err))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, errors.New("image file is required"))
		return
	}
	defer file.Close()

	prod, err := sv.repo.GetProductByID(r.Context(), repository.GetProductByIDParams{ID: uuid.MustParse(id)})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dto.CreateErr(NotFoundCode, err))
		return
	}
	if prod.ImageID != nil {
		msg, err := sv.uploadService.Remove(r.Context(), *prod.ImageID)
		if err != nil {
			log.Error().Err(err).Timestamp().Msg(msg)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.CreateErr(InternalServerErrorCode, err))
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

	uploadID, url, err := sv.uploadService.Upload(r.Context(), fileHeader)
	if err != nil {
		log.Error().Err(err).Timestamp().Msg("UploadFile")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	updated, err := sv.repo.UpdateProduct(r.Context(), repository.UpdateProductParams{
		ImageUrl: &url,
		ImageID:  &uploadID,
		ID:       uuid.MustParse(id),
	})
	if err != nil {
		log.Error().Err(err).Timestamp().Msg("CreateProduct")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.CreateDataResp(r, updated, nil, nil))
}

// @Summary Create a new product variant
// @Schemes http
// @Description create a new product with the input payload
// @Tags products
// @Accept json
// @Param input body CreateProdVariantReq true "Product variant input"
// @Produce json
// @Success 200 {object} ApiResponse[string]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/products/{id}/variants [post]
func (sv *Server) adminAddVariant(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.CreateErr(InvalidBodyCode, errors.New("id parameter is required")))
		return
	}

	var req models.CreateProdVariantModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.CreateErr(InvalidBodyCode, err))
		return
	}
	if len(req.AttributeValues) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.CreateErr(InvalidBodyCode, errors.New("attribute values are required")))
		return
	}

	prod, err := sv.repo.GetProductByID(r.Context(), repository.GetProductByIDParams{ID: uuid.MustParse(id)})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dto.CreateErr(NotFoundCode, err))
		return
	}

	prodAttrs, err := sv.repo.GetProductAttributesByProductID(r.Context(), prod.ID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	if len(prodAttrs) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.CreateErr(InvalidBodyCode, errors.New("product has no attributes")))
		return
	}

	prodAttrIds := make([]int32, len(prodAttrs))
	for i, attr := range prodAttrs {
		prodAttrIds[i] = attr.AttributeID
	}

	attributeValues, err := sv.repo.GetAttributeValuesByIDs(r.Context(), req.AttributeValues)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	if len(attributeValues) != len(prodAttrIds) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.CreateErr(InvalidBodyCode, errors.New("attribute values do not match product attributes")))
		return
	}

	for _, attrVal := range attributeValues {
		if !slices.Contains(prodAttrIds, attrVal.AttributeID) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.CreateErr(InvalidBodyCode, errors.New("attribute value does not belong to product attributes")))
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

	variant, err := sv.repo.CreateProductVariant(r.Context(), createParams)
	if err != nil {
		log.Error().Err(err).Timestamp().Msg("CreateProduct")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	variantAttrParams := make([]repository.CreateBulkProductVariantAttributeParams, len(req.AttributeValues))

	for i, attrValID := range req.AttributeValues {
		variantAttrParams[i] = repository.CreateBulkProductVariantAttributeParams{
			VariantID:        variant.ID,
			AttributeValueID: attrValID,
		}
	}
	_, err = sv.repo.CreateBulkProductVariantAttribute(r.Context(), variantAttrParams)
	if err != nil {
		log.Error().Err(err).Msg("CreateBulkProductVariantAttribute")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.CreateDataResp(r, variant.ID.String(), nil, nil))
}

// @Summary Get product variants
// @Schemes http
// @Description get product variants
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} ApiResponse[[]VariantModelDto]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/products/{id}/variants [get]
func (sv *Server) adminGetVariants(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.CreateErr(InvalidBodyCode, errors.New("id parameter is required")))
		return
	}

	prod, err := sv.repo.GetProductByID(r.Context(), repository.GetProductByIDParams{ID: uuid.MustParse(id)})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dto.CreateErr(NotFoundCode, err))
		return
	}
	variantRows, err := sv.repo.GetProductVariantList(r.Context(), repository.GetProductVariantListParams{ProductID: prod.ID})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	resp := make([]dto.VariantDetail, len(variantRows))
	for i, row := range variantRows {
		resp[i] = dto.MapToVariantListModelDto(row)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.CreateDataResp(r, resp, nil, nil))
}

// @Summary Get product variant
// @Schemes http
// @Description get product variant
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param variantID path string true "Product Variant ID"
// @Success 200 {object} ApiResponse[VariantModelDto]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/products/{id}/variants/{variantID} [get]
func (sv *Server) adminGetVariant(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	variantId := chi.URLParam(r, "variantId")
	if id == "" || variantId == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.CreateErr(InvalidBodyCode, errors.New("id and variantId parameters are required")))
		return
	}

	prod, err := sv.repo.GetProductByID(r.Context(), repository.GetProductByIDParams{ID: uuid.MustParse(id)})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dto.CreateErr(NotFoundCode, err))
		return
	}

	rows, err := sv.repo.GetVariantDetailByID(r.Context(), repository.GetVariantDetailByIDParams{
		ID:        uuid.MustParse(variantId),
		ProductID: prod.ID,
	})

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dto.CreateErr(NotFoundCode, err))
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.CreateDataResp(r, resp, nil, nil))
}

// @Summary Update a product variant
// @Schemes http
// @Description update a product with the input payload
// @Tags products
// @Accept json
// @Param input body UpdateProdVariantReq true "Product variant input"
// @Produce json
// @Success 200 {object} ApiResponse[repository.ProductVariant]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/products/{id}/variants/{variantId} [put]
func (sv *Server) adminUpdateVariant(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	variantId := chi.URLParam(r, "variantId")
	if id == "" || variantId == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.CreateErr(InvalidBodyCode, errors.New("id and variantId parameters are required")))
		return
	}

	var req models.UpdateProdVariantModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.CreateErr(InvalidBodyCode, err))
		return
	}

	prod, err := sv.repo.GetProductByID(r.Context(), repository.GetProductByIDParams{ID: uuid.MustParse(id)})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dto.CreateErr(NotFoundCode, err))
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

	updatedVariant, err := sv.repo.UpdateProductVariant(r.Context(), updateParams)
	if err != nil {
		log.Error().Err(err).Timestamp().Msg("UpdateProductVariant")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.CreateDataResp(r, updatedVariant, nil, nil))
}

// @Summary Upload a product variant image
// @Schemes http
// @Description upload a product variant image with the input payload
// @Tags products
// @Accept multipart/form-data
// @Param id path string true "Product ID"
// @Param variantId path string true "Product Variant ID"
// @Produce json
// @Success 200 {object} ApiResponse[repository.ProductVariant]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/products/{id}/variants/{variantId}/images [post]
func (sv *Server) adminUploadVariantImage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	variantId := chi.URLParam(r, "variantId")
	if id == "" || variantId == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.CreateErr(InvalidBodyCode, errors.New("id and variantId parameters are required")))
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10MB max
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.CreateErr(InvalidBodyCode, err))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.CreateErr(InvalidBodyCode, err))
		return
	}
	defer file.Close()

	prod, err := sv.repo.GetProductByID(r.Context(), repository.GetProductByIDParams{ID: uuid.MustParse(id)})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dto.CreateErr(NotFoundCode, err))
		return
	}

	variant, err := sv.repo.GetProductVariantByID(r.Context(), repository.GetProductVariantByIDParams{
		ID:        uuid.MustParse(variantId),
		ProductID: prod.ID,
	})

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dto.CreateErr(NotFoundCode, err))
		return
	}

	if variant.ImageID != nil {
		msg, err := sv.uploadService.Remove(r.Context(), *variant.ImageID)
		if err != nil {
			log.Error().Err(err).Timestamp().Msg(msg)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.CreateErr(InternalServerErrorCode, err))
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

	uploadID, url, err := sv.uploadService.Upload(r.Context(), fileHeader)
	if err != nil {
		log.Error().Err(err).Timestamp().Msg("UploadFile")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	updateParam := repository.UpdateProductVariantParams{
		ProductID: prod.ID,
		ID:        variant.ID,
		ImageID:   &uploadID,
		ImageUrl:  &url,
	}

	updatedVariant, err := sv.repo.UpdateProductVariant(r.Context(), updateParam)
	if err != nil {
		log.Error().Err(err).Timestamp().Msg("UpdateProductVariant")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.CreateDataResp(r, updatedVariant, nil, nil))
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
// @Router /admin/products/{id}/variant/{variantID} [delete]
func (sv *Server) adminDeleteVariant(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	variantId := chi.URLParam(r, "variantId")
	if id == "" || variantId == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.CreateErr(InvalidBodyCode, errors.New("id and variantId parameters are required")))
		return
	}

	prod, err := sv.repo.GetProductByID(r.Context(), repository.GetProductByIDParams{ID: uuid.MustParse(id)})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dto.CreateErr(NotFoundCode, err))
		return
	}

	err = sv.repo.DeleteProductVariant(r.Context(), repository.DeleteProductVariantParams{
		ProductID: prod.ID,
		ID:        uuid.MustParse(variantId),
	})

	if err != nil {
		log.Error().Err(err).Timestamp().Msg("DeleteVariant")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(dto.CreateDataResp(r, struct{}{}, nil, nil))
}

// @Summary Create an attribute
// @Description Create an attribute
// @Tags attributes
// @Accept json
// @Produce json
// @Param params body AttributeValuesReq true "Attribute name"
// @Success 201 {object} ApiResponse[AttributeRespModel]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/attributes [post]
func (sv *Server) adminCreateAttribute(w http.ResponseWriter, r *http.Request) {
	var req models.AttributeModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.CreateErr(InvalidBodyCode, err))
		return
	}

	attribute, err := sv.repo.CreateAttribute(r.Context(), req.Name)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.CreateErr(InvalidBodyCode, err))
		return
	}

	attributeResp := dto.AttributeDetail{
		ID:   attribute.ID,
		Name: attribute.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.CreateDataResp(r, attributeResp, nil, nil))
}

// @Summary Get all attributes
// @Description Get all attributes
// @Tags attributes
// @Accept json
// @Produce json
// @Success 200 {object} ApiResponse[[]AttributeRespModel]
// @Failure 500 {object} ErrorResp
// @Router /admin/attributes [get]
func (sv *Server) adminGetAttributes(w http.ResponseWriter, r *http.Request) {
	var queries models.AttributesQuery
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	attributeRows, err := sv.repo.GetAttributes(c, queries.IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	var attributeResp = []dto.AttributeDetail{}
	for i := range attributeRows {
		attrVal := attributeRows[i]
		if i == 0 || attributeRows[i].ID != attributeRows[i-1].ID {
			attributeResp = append(attributeResp, dto.AttributeDetail{
				ID:     attrVal.ID,
				Name:   attrVal.Name,
				Values: []dto.AttributeValueDetail{},
			})
			if attrVal.AttrValueID != nil {
				id := *attrVal.AttrValueID
				attributeResp[len(attributeResp)-1].Values = append(attributeResp[len(attributeResp)-1].Values, dto.AttributeValueDetail{
					ID:    id,
					Value: *attrVal.AttrValue,
				})
			}
		} else if attrVal.AttrValueID != nil {
			id := *attrVal.AttrValueID
			attributeResp[len(attributeResp)-1].Values = append(attributeResp[len(attributeResp)-1].Values, dto.AttributeValueDetail{
				ID:    id,
				Value: *attrVal.AttrValue,
			})
		}
	}

	RespondSuccess(w, r, attributeResp)
}

// @Summary Get an attribute
// @Description Get an attribute
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Success 200 {object} ApiResponse[AttributeRespModel]
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/attributes/{id} [get]
func (sv *Server) adminGetAttributeByID(w http.ResponseWriter, r *http.Request) {
	var attributeParam models.AttributeParam
	if err := c.ShouldBindUri(&attributeParam); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	attr, err := sv.repo.GetAttributeByID(c, attributeParam.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	attributeResp := dto.AttributeDetail{
		Name: attr.Name,
		ID:   attr.ID,
	}

	values, err := sv.repo.GetAttributeValues(c, attr.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	attributeResp.Values = make([]dto.AttributeValueDetail, len(values))

	for i, val := range values {
		attributeResp.Values[i] = dto.AttributeValueDetail{
			ID:    val.ID,
			Value: val.Value,
		}
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, attributeResp, nil, nil))
}

// @Summary Update an attribute
// @Description Update an attribute
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Param params body AttributeRequest true "Attribute name"
// @Success 200 {object} ApiResponse[repository.Attribute]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/attributes/{id} [put]
func (sv *Server) adminUpdateAttribute(c *gin.Context) {
	var param models.AttributeParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	var req models.AttributeModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	attr, err := sv.repo.GetAttributeByID(c, param.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	attribute, err := sv.repo.UpdateAttribute(c, repository.UpdateAttributeParams{
		ID:   attr.ID,
		Name: req.Name,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, attribute, nil, nil))
}

// @Summary Delete an attribute
// @Description Delete an attribute
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Success 204 {object} nil
// @Failure 500 {object} ErrorResp
// @Router /admin/attributes/{id} [delete]
func (sv *Server) adminRemoveAttribute(c *gin.Context) {
	var params models.AttributeParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	attr, err := sv.repo.GetAttributeByID(c, params.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, err))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	err = sv.repo.DeleteAttribute(c, attr.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Get attributes and their values by for a product
// @Description Get attributes and their values for a product
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} ApiResponse[[]AttributeRespModel]
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/attributes/product/{id} [get]
func (sv *Server) adminGetAttributeValuesForProduct(c *gin.Context) {
	var uri models.UriIDParam
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	attrs, err := sv.repo.GetProductAttributeValuesByProductID(c, uuid.MustParse(uri.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	resp := make([]dto.AttributeDetail, 0)
	for _, attr := range attrs {

		if slices.ContainsFunc(resp, func(a dto.AttributeDetail) bool {
			return *attr.AttributeID == a.ID
		}) {
			// push value to existing attribute
			for i, r := range resp {
				if r.ID == *attr.AttributeID {
					resp[i].Values = append(resp[i].Values, dto.AttributeValueDetail{
						ID:    *attr.AttributeValueID,
						Value: *attr.AttributeValue,
					})
					break
				}
			}
		} else {
			// create new attribute
			attrResp := dto.AttributeDetail{
				ID:   *attr.AttributeID,
				Name: *attr.AttributeName,
				Values: []dto.AttributeValueDetail{
					{
						ID:    *attr.AttributeValueID,
						Value: *attr.AttributeValue,
					},
				},
			}
			resp = append(resp, attrResp)
		}
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, resp, nil, nil))
}

// @Summary Add new attribute value
// @Description Add new attribute value
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Param params body AttributeValuesReq true "Attribute value"
// @Success 200 {object} ApiResponse[bool]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/attributes/{id}/create [post]
func (sv *Server) adminAddAttributeValue(c *gin.Context) {
	var param models.AttributeParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	var req models.AttributeValueModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	attr, err := sv.repo.GetAttributeByID(c, param.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	obj, err := sv.repo.CreateAttributeValue(c, repository.CreateAttributeValueParams{
		AttributeID: attr.ID,
		Value:       req.Value,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusCreated, dto.CreateDataResp(c, obj, nil, nil))
}

// @Summary update attribute value
// @Description update attribute value
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Param params body AttributeValuesReq true "Attribute value"
// @Success 200 {object} ApiResponse[bool]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/attributes/{id}/update/{valueId} [put]
func (sv *Server) adminUpdateAttrValue(c *gin.Context) {
	var param models.AttributeParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	if param.ValueID == nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, nil))
		return
	}
	var req models.AttributeValueModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	attr, err := sv.repo.GetAttributeByID(c, param.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, err))
		return
	}

	res, err := sv.repo.UpdateAttributeValue(c, repository.UpdateAttributeValueParams{
		AttributeID: attr.ID,
		ID:          *param.ValueID,
		Value:       req.Value,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, res, nil, nil))
}

// @Summary remove an attribute value
// @Description remove an attribute value
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Success 204 {object} nil
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/attributes/{id}/remove/{valueId} [delete]
func (sv *Server) adminRemoveAttrValue(c *gin.Context) {
	var param models.AttributeParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	if param.ValueID == nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, nil))
		return
	}

	attr, err := sv.repo.GetAttributeByID(c, param.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	err = sv.repo.DeleteAttributeValueByValueID(c, repository.DeleteAttributeValueByValueIDParams{
		AttributeID: attr.ID,
		ID:          *param.ValueID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusNoContent, dto.CreateDataResp(c, struct{}{}, nil, nil))
}

// @Summary Get all orders (admin endpoint)
// @Description Get all orders with pagination and filtering
// @Tags admin
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Param status query string false "Filter by status"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[[]OrderListResponse]
// @Failure 401 {object} ErrorResp
// @Failure 403 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/orders [get]
func (sv *Server) adminGetOrders(c *gin.Context) {
	var orderListQuery models.OrderListQuery
	if err := c.ShouldBindQuery(&orderListQuery); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	dbParams := repository.GetOrdersParams{
		Limit:  orderListQuery.PageSize,
		Offset: (orderListQuery.Page - 1) * orderListQuery.PageSize,
	}

	if orderListQuery.Status != nil {
		dbParams.Status = repository.NullOrderStatus{
			OrderStatus: repository.OrderStatus(*orderListQuery.Status),
			Valid:       true,
		}
	}

	if orderListQuery.PaymentStatus != nil {
		dbParams.PaymentStatus = repository.NullPaymentStatus{
			PaymentStatus: repository.PaymentStatus(*orderListQuery.PaymentStatus),
			Valid:         true,
		}
	}

	fetchedOrderRows, err := sv.repo.GetOrders(c, dbParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	countParams := repository.CountOrdersParams{}
	if orderListQuery.Status != nil {
		countParams.Status = repository.NullOrderStatus{
			OrderStatus: repository.OrderStatus(*orderListQuery.Status),
			Valid:       true,
		}
	}

	if orderListQuery.PaymentStatus != nil {
		countParams.PaymentStatus = repository.NullPaymentStatus{
			PaymentStatus: repository.PaymentStatus(*orderListQuery.PaymentStatus),
			Valid:         true,
		}
	}

	count, err := sv.repo.CountOrders(c, countParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	var orderResponses []dto.OrderListItem
	for _, aggregated := range fetchedOrderRows {
		// Convert PaymentStatus interface{} to PaymentStatus type
		paymentStatus := repository.PaymentStatusPending
		if aggregated.PaymentStatus.Valid {
			paymentStatus = aggregated.PaymentStatus.PaymentStatus
		}

		total, _ := aggregated.TotalPrice.Float64Value()
		orderResponses = append(orderResponses, dto.OrderListItem{
			ID:            aggregated.ID,
			Total:         total.Float64,
			TotalItems:    int32(aggregated.TotalItems),
			Status:        aggregated.Status,
			CustomerName:  aggregated.CustomerName,
			CustomerEmail: aggregated.CustomerEmail,
			PaymentStatus: paymentStatus,
			CreatedAt:     aggregated.CreatedAt.UTC(),
			UpdatedAt:     aggregated.UpdatedAt.UTC(),
		})
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, orderResponses, dto.CreatePagination(orderListQuery.Page, orderListQuery.PageSize, count), nil))
}

// @Summary Get order details by ID (admin endpoint)
// @Description Get detailed information about an order by its ID
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[OrderDetailResponse]
// @Failure 401 {object} ErrorResp
// @Failure 403 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/orders/{id} [get]
func (sv *Server) adminGetOrderDetail(c *gin.Context) {
	// Reuse the existing order detail handler since admin has access to all orders
	sv.getOrderDetail(c)
}

// @Summary Change order status
// @Description Change order status by order ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param status body string true "Status"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[OrderListResponse]
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/orders/{orderId}/status [put]
func (sv *Server) adminChangeOrderStatus(c *gin.Context) {
	var params models.UriIDParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	var req models.OrderStatusModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	order, err := sv.repo.GetOrder(c, uuid.MustParse(params.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, err))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	if order.Status == repository.OrderStatusDelivered || order.Status == repository.OrderStatusCancelled || order.Status == repository.OrderStatusRefunded {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidPaymentCode, errors.New("order cannot be changed")))
		return
	}

	status := repository.OrderStatus(req.Status)

	updateParams := repository.UpdateOrderParams{
		ID: order.ID,
		Status: repository.NullOrderStatus{
			OrderStatus: status,
			Valid:       true,
		},
	}
	if status == repository.OrderStatusConfirmed {
		updateParams.ConfirmedAt = pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		}
	}
	if status == repository.OrderStatusDelivering {
		updateParams.DeliveredAt = pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		}
	}

	rs, err := sv.repo.UpdateOrder(c, updateParams)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	if err := sv.cacheSrv.Delete(c, "order_detail:"+params.ID); err != nil {
		log.Err(err).Msg("failed to delete order detail cache")
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	sv.cacheSrv.Delete(c, "order_detail:"+params.ID)

	c.JSON(http.StatusOK, dto.CreateDataResp(c, rs, nil, nil))
}

// @Summary Cancel order
// @Description Cancel order by order ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[OrderListResponse]
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/orders/{orderId}/cancel [put]
func (sv *Server) adminCancelOrder(c *gin.Context) {
	tokenPayload, _ := c.MustGet(constants.AuthPayLoad).(*auth.TokenPayload)

	var params models.UriIDParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	var req models.CancelOrderModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	order, err := sv.repo.GetOrder(c, uuid.MustParse(params.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	if order.Status == repository.OrderStatusCancelled || order.Status == repository.OrderStatusRefunded {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidPaymentCode, errors.New("order is already cancelled or refunded")))
		return
	}

	userRole := c.GetString(constants.UserRole)
	if order.UserID != tokenPayload.UserID && userRole != "admin" {
		c.JSON(http.StatusForbidden, dto.CreateErr(PermissionDeniedCode, errors.New("you do not have permission to access this order")))
		return
	}

	paymentRow, err := sv.repo.GetPaymentByOrderID(c, order.ID)
	if err != nil && !errors.Is(err, repository.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	// if order status is not pending or user is not admin
	if order.Status != repository.OrderStatusPending || (!errors.Is(err, repository.ErrRecordNotFound) && paymentRow.Status != repository.PaymentStatusPending) {
		c.JSON(http.StatusBadRequest, dto.CreateErr(PermissionDeniedCode, errors.New("order cannot be cancelled")))
		return
	}

	// if order
	cancelOrderTxParams := repository.CancelOrderTxArgs{
		OrderID: uuid.MustParse(params.ID),
		CancelPaymentFromMethod: func(paymentID string, method string) error {
			req := payment.RefundRequest{
				TransactionID: paymentID,
				Amount:        paymentRow.Amount.Int.Int64(),
			}
			_, err = sv.paymentSrv.RefundPayment(c, req, *paymentRow.Gateway)
			return err
		},
	}
	ordId, err := sv.repo.CancelOrderTx(c, cancelOrderTxParams)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(repository.ErrDeadlockDetected.InternalQuery, err))
		return
	}
	sv.cacheSrv.Delete(c, "order_detail:"+params.ID)
	c.JSON(http.StatusOK, dto.CreateDataResp(c, ordId, nil, nil))
}

// @Summary Refund order
// @Description Refund order by order ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[OrderListResponse]
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/order/{orderId}/refund [put]
func (sv *Server) adminRefundOrder(c *gin.Context) {
	var params models.UriIDParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	var req models.RefundOrderModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	order, err := sv.repo.GetOrder(c, uuid.MustParse(params.ID))
	if err != nil {
		if err == repository.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, err))
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	if order.Status != repository.OrderStatusDelivered {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidPaymentCode, errors.New("order cannot be refunded")))
		return
	}

	err = sv.repo.RefundOrderTx(c, repository.RefundOrderTxArgs{
		OrderID: uuid.MustParse(params.ID),
		RefundPaymentFromMethod: func(paymentID string, method string) (string, error) {
			req := payment.RefundRequest{
				TransactionID: paymentID,
				Amount:        order.TotalPrice.Int.Int64(),
				Reason:        req.Reason,
			}
			rs, err := sv.paymentSrv.RefundPayment(c, req, method)
			return rs.Reason, err
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	sv.cacheSrv.Delete(c, "order_detail:"+params.ID)

	c.JSON(http.StatusOK, dto.CreateDataResp(c, order, nil, nil))
}

// adminGetCategories retrieves a list of Categories.
// @Summary Get a list of Categories
// @Description Get a list of Categories
// @ID get-admin-Categories
// @Accept json
// @Tags Categories
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} ApiResponse[[]dto.CategoryDetail]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/categories [get]
func (sv *Server) adminGetCategories(w http.ResponseWriter, r *http.Request) {
	var query models.PaginationQuery
	query.Page = 1
	query.PageSize = 10
	if page := r.URL.Query().Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			query.Page = int64(p)
		}
	}
	if pageSize := r.URL.Query().Get("pageSize"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil {
			query.PageSize = int64(ps)
		}
	}
	params := repository.GetCategoriesParams{
		Limit:  10,
		Offset: 0,
	}
	params.Offset = (params.Limit) * int64(query.Page-1)
	params.Limit = int64(query.PageSize)

	categories, err := sv.repo.GetCategories(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	count, err := sv.repo.CountCategories(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	categoriesResp := make([]dto.AdminCategoryDetail, len(categories))

	for i, category := range categories {
		categoriesResp[i] = dto.AdminCategoryDetail{
			ID:          category.ID.String(),
			Name:        category.Name,
			Slug:        category.Slug,
			Published:   category.Published,
			CreatedAt:   category.CreatedAt.String(),
			Description: category.Description,
			ImageUrl:    category.ImageUrl,
			UpdatedAt:   category.UpdatedAt.String(),
		}
	}
	c.JSON(http.StatusOK, dto.CreateDataResp(c, categoriesResp, dto.CreatePagination(query.Page, query.PageSize, count), nil))
}

// adminGetCategoryByID retrieves a Category by its ID.
// @Summary Get a Category by ID
// @Description Get a Category by ID
// @ID get-Category-by-id
// @Accept json
// @Tags Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} ApiResponse[dto.CategoryDetail]
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/categories/{id} [get]
func (sv *Server) adminGetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.CreateErr(InvalidBodyCode, errors.New("id parameter is required")))
		return
	}

	category, err := sv.repo.GetCategoryByID(c, uuid.MustParse(param.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(InvalidBodyCode, fmt.Errorf("category with ID %s not found", param.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	resp := dto.CategoryDetail{
		ID:          category.ID.String(),
		Name:        category.Name,
		Slug:        category.Slug,
		Published:   category.Published,
		CreatedAt:   category.CreatedAt.String(),
		Description: category.Description,
		ImageUrl:    category.ImageUrl,
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, resp, nil, nil))
}

// adminCreateCategory creates a new Category.
// @Summary Create a new Category
// @Description Create a new Category
// @ID create-Category
// @Accept json
// @Tags Categories
// @Produce json
// @Param request body CreateCategoryRequest true "Category request"
// @Success 201 {object} ApiResponse[dto.CategoryDetail]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/categories [post]
func (sv *Server) adminCreateCategory(c *gin.Context) {
	var req models.CreateCategoryModel
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	params := repository.CreateCategoryParams{
		Name: req.Name,
		Slug: req.Slug,
	}

	if req.Description != nil {
		params.Description = req.Description
	}

	if req.Image != nil {
		imageID, imageURL, err := sv.uploadService.Upload(c, req.Image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.CreateErr(UploadFileCode, err))
			return
		}
		params.ImageID = &imageID
		params.ImageUrl = &imageURL
	}

	col, err := sv.repo.CreateCategory(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	resp := dto.CategoryDetail{
		ID:          col.ID.String(),
		Name:        col.Name,
		Slug:        col.Slug,
		Published:   col.Published,
		CreatedAt:   col.CreatedAt.String(),
		Description: col.Description,
		ImageUrl:    col.ImageUrl,
	}

	c.JSON(http.StatusCreated, dto.CreateDataResp(c, resp, nil, nil))
}

// adminUpdateCategory updates a Category.
// @Summary Update a Category
// @Description Update a Category
// @ID update-Category
// @Accept json
// @Tags admin
// @Produce json
// @Param id path int true "Category ID"
// @Param request body models.UpdateCategoryModel true "Category request"
// @Success 200 {object} ApiResponse[repository.Category]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/categories/{id} [put]
func (sv *Server) adminUpdateCategory(c *gin.Context) {
	var param models.UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	var req models.UpdateCategoryModel
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	category, err := sv.repo.GetCategoryByID(c, uuid.MustParse(param.ID))

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, fmt.Errorf("category with ID %s not found", param.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	updateParam := repository.UpdateCategoryParams{
		ID: category.ID,
	}

	if req.Name != nil {
		updateParam.Name = req.Name
	}

	if req.Slug != nil {
		updateParam.Slug = req.Slug
	}

	if req.Description != nil {
		updateParam.Description = req.Description
	}

	if req.Published != nil {
		updateParam.Published = req.Published
	}
	var apiErr *dto.ApiError

	imageID, imageURL := "", ""
	if req.Image != nil {
		oldImageID := category.ImageID
		oldImageURL := category.ImageUrl
		// remove old image
		if oldImageID != nil && oldImageURL != nil {
			_, err = sv.uploadService.Remove(c, *oldImageID)
			if err != nil {
				apiErr = &dto.ApiError{
					Code:    UploadFileCode,
					Details: err.Error(),
					Stack:   err}
			}
		}
		imageID, imageURL, err = sv.uploadService.Upload(c, req.Image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.CreateErr(UploadFileCode, err))
			return
		}
		updateParam.ImageID = &imageID
		updateParam.ImageUrl = &imageURL
	}
	col, err := sv.repo.UpdateCategory(c, updateParam)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, col, nil, apiErr))
}

// adminDeleteCategory delete a Category.
// @Summary Delete a Category
// @Description Delete a Category
// @ID delete-Category
// @Accept json
// @Tags admin
// @Produce json
// @Param id path int true "Category ID"
// @Success 204 {object} nil
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/categories/{id} [delete]
func (sv *Server) adminDeleteCategory(c *gin.Context) {
	var param models.UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	_, err := sv.repo.GetCategoryByID(c, uuid.MustParse(param.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, fmt.Errorf("category with ID %s not found", param.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	err = sv.repo.DeleteCategory(c, uuid.MustParse(param.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary Create a new Brand
// @Description Create a new Brand
// @Tags admin
// @ID create-Brand
// @Accept json
// @Produce json
// @Param request body CreateCategoryRequest true "Brand request"
// @Success 201 {object} ApiResponse[CategoryDto]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/brands [post]
func (sv *Server) adminCreateBrand(c *gin.Context) {
	var req models.CreateCategoryModel
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	params := repository.CreateBrandParams{
		Name: req.Name,
		Slug: req.Slug,
	}

	if req.Description != nil {
		params.Description = req.Description
	}
	if req.Image != nil {
		publicID, imgUrl, err := sv.uploadService.Upload(c, req.Image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.CreateErr(UploadFileCode, err))
			return
		}
		params.ImageUrl = &imgUrl
		params.ImageID = &publicID
	}

	col, err := sv.repo.CreateBrand(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusCreated, dto.CreateDataResp(c, col, nil, nil))
}

// @Summary Get a list of brands
// @Description Get a list of brands
// @Tags admin
// @ID get-brands
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} ApiResponse[[]CategoryDto]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/brands [get]
func (sv *Server) adminGetBrands(c *gin.Context) {
	var queries models.PaginationQuery
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
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
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	cnt, err := sv.repo.CountBrands(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	data := make([]dto.CategoryDetail, 0, len(rows))

	for _, row := range rows {
		data = append(data, dto.CategoryDetail{
			ID:          row.ID.String(),
			Name:        row.Name,
			Description: row.Description,
			Slug:        row.Slug,
			ImageUrl:    row.ImageUrl,
		})
	}

	pagination := dto.CreatePagination(queries.Page, queries.PageSize, cnt)

	resp := dto.CreateDataResp(c, data, pagination, nil)
	c.JSON(http.StatusOK, resp)
}

// @Summary Get a Brand by ID
// @Description Get a Brand by ID
// @ID get-Brand-by-id
// @Accept json
// @Tags admin
// @Produce json
// @Param id path int true "Brand ID"
// @Success 200 {object} ApiResponse[CategoryDto]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/brands/{id} [get]
func (sv *Server) adminGetBrandByID(c *gin.Context) {
	var param models.UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	result, err := sv.repo.GetBrandByID(c, uuid.MustParse(param.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, fmt.Errorf("brand with ID %s not found", param.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	colResp := dto.AdminCategoryDetail{
		ID:          result.ID.String(),
		Name:        result.Name,
		Description: result.Description,
		Slug:        result.Slug,
		Published:   result.Published,
		CreatedAt:   result.CreatedAt.Format("2006-01-02 15:04:05"),
		ImageUrl:    result.ImageUrl,
		UpdatedAt:   result.UpdatedAt.String(),
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, colResp, nil, nil))
}

// @Summary Update a Brand
// @Description Update a Brand
// @ID update-Brand
// @Accept json
// @Produce json
// @Tags admin
// @Param id path int true "Brand ID"
// @Param request body UpdateCategoryRequest true "Brand request"
// @Success 200 {object} ApiResponse[CategoryDto]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/brands/{id} [put]
func (sv *Server) adminUpdateBrand(c *gin.Context) {
	var param models.UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	var req models.UpdateCategoryModel
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	brand, err := sv.repo.GetBrandByID(c, uuid.MustParse(param.ID))

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, fmt.Errorf("brand with ID %s not found", param.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	updateParam := repository.UpdateBrandWithParams{
		ID:   brand.ID,
		Name: req.Name,
	}

	if req.Image != nil {

		imgID, imgUrl, err := sv.uploadService.Upload(c, req.Image)
		if err != nil {
			log.Error().Err(err).Interface("value", req.Image.Header).Msg("error when upload image")
			c.JSON(http.StatusInternalServerError, dto.CreateErr(UploadFileCode, err))
			return
		}
		updateParam.ImageUrl = &imgUrl
		updateParam.ImageID = &imgID
		oldImageID := brand.ImageID
		if oldImageID != nil {
			_, err := sv.uploadService.Remove(c, *oldImageID)
			if err != nil {
				log.Error().Err(err).Msg("error when remove old image")
				c.JSON(http.StatusInternalServerError, dto.CreateErr(UploadFileCode, err))
				return
			}
			log.Info().Msgf("old image %s removed", *oldImageID)
		}
	}

	if req.Slug != nil {
		updateParam.Slug = req.Slug
	}
	if req.Description != nil {
		updateParam.Description = req.Description
	}

	if req.Published != nil {
		updateParam.Published = req.Published
	}

	col, err := sv.repo.UpdateBrandWith(c, updateParam)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, col, nil, nil))
}

// @Summary Delete a Brand
// @Description Delete a Brand
// @ID delete-Brand
// @Accept json
// @Tags admin
// @Produce json
// @Param id path int true "Brand ID"
// @Success 204 {object} ApiResponse[bool]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/brands/{id} [delete]
func (sv *Server) adminDeleteBrand(c *gin.Context) {
	var colID models.UriIDParam
	if err := c.ShouldBindUri(&colID); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	_, err := sv.repo.GetBrandByID(c, uuid.MustParse(colID.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, fmt.Errorf("brand with ID %s not found", colID.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	err = sv.repo.DeleteBrand(c, uuid.MustParse(colID.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	c.JSON(http.StatusOK, dto.CreateDataResp(c, true, nil, nil))
}

// @Summary Create a new Collection
// @Description Create a new Collection
// @ID create-Collection
// @Accept json
// @Tags admin
// @Produce json
// @Param request body models.CreateCategoryModel true "Collection info"
// @Success 201 {object} ApiResponse[dto.CategoryDetail]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/collections [post]
func (sv *Server) adminCreateCollection(c *gin.Context) {
	var req models.CreateCategoryModel
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode,
			err))
		return
	}
	createParams := repository.CreateCollectionParams{
		Name: req.Name,
		Slug: req.Slug,
	}

	if req.Description != nil {
		createParams.Description = req.Description
	}

	if req.Image != nil {
		ID, url, err := sv.uploadService.Upload(c, req.Image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.CreateErr(UploadFileCode,
				err))
			return
		}

		createParams.ImageID = &ID
		createParams.ImageUrl = &url
	}

	col, err := sv.repo.CreateCollection(c, createParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode,
			err))
		return
	}
	sv.cacheSrv.Delete(c, "collections-*")

	c.JSON(http.StatusCreated, dto.CreateDataResp(c, col, nil, nil))
}

// @Summary Get a list of Collections
// @Description Get a list of Collections
// @ID get-Collections
// @Accept json
// @Tags admin
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} ApiResponse[dto.CategoryDetail]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/collections [get]
func (sv *Server) adminGetCollections(c *gin.Context) {
	var queries models.PaginationQuery
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode,
			err))
		return
	}

	dbQueries := repository.GetCollectionsParams{
		Limit:  20,
		Offset: 0,
	}

	dbQueries.Offset = int64(queries.Page-1) * int64(queries.PageSize)
	dbQueries.Limit = int64(queries.PageSize)
	collectionRows, err := sv.repo.GetCollections(c, dbQueries)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode,
			err))
		return
	}

	cnt, err := sv.repo.CountCollections(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode,
			err))
		return
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, collectionRows, dto.CreatePagination(cnt, queries.Page, queries.PageSize), nil))
}

// @Summary Get a Collection by ID
// @Description Get a Collection by ID
// @ID get-Collection-by-id
// @Accept json
// @Tags admin
// @Produce json
// @Param id path int true "Collection ID"
// @Success 200 {object} ApiResponse[dto.CategoryDetail]
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/collections/{id} [get]
func (sv *Server) adminGetCollectionByID(c *gin.Context) {
	var param models.UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode,
			err))
		return
	}

	collection, err := sv.repo.GetCollectionByID(c, uuid.MustParse(param.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode,
				fmt.Errorf("collection with ID %s not found", param.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode,
			err))
		return
	}

	colResp := dto.CategoryDetail{
		ID:          collection.ID.String(),
		Slug:        collection.Slug,
		Description: collection.Description,
		Published:   collection.Published,
		Name:        collection.Name,
		ImageUrl:    collection.ImageUrl,
		CreatedAt:   collection.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, colResp, nil, nil))
}

// @Summary Update a Collection
// @Description Update a Collection
// @ID update-Collection
// @Accept json
// @Tags admin
// @Produce json
// @Param id path int true "Collection ID"
// @Param request body models.CreateCategoryModel true "Collection info"
// @Success 200 {object} ApiResponse[dto.CategoryDetail]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/collections/{id} [put]
func (sv *Server) adminUpdateCollection(c *gin.Context) {
	var param models.UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode,
			err))
		return
	}
	var req models.UpdateCategoryModel
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode,
			err))
		return
	}

	collection, err := sv.repo.GetCollectionByID(c, uuid.MustParse(param.ID))

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode,
				fmt.Errorf("collection with ID %s not found", param.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode,
			err))
		return
	}

	updateParam := repository.UpdateCollectionWithParams{
		ID: collection.ID,
	}
	if req.Name != nil {
		updateParam.Name = req.Name
	}
	if req.Description != nil {
		updateParam.Description = req.Description
	}

	if req.Image != nil {
		oldImageID := collection.ImageID
		oldImageUrl := collection.ImageUrl
		ID, url, err := sv.uploadService.Upload(c, req.Image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode,
				err))
			return
		}

		updateParam.ImageUrl = &url
		updateParam.ImageID = &ID

		// Delete old image
		if oldImageID != nil && oldImageUrl != nil {
			if _, err := sv.uploadService.Remove(c, *oldImageID); err != nil {
				c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
				return
			}
		}
	}

	if req.Published != nil {
		updateParam.Published = req.Published
	}

	col, err := sv.repo.UpdateCollectionWith(c, updateParam)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode,
			err))
		return
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, col,
		nil, nil))
}

// @Summary Delete a Collection
// @Description Delete a Collection
// @ID delete-Collection
// @Accept json
// @Tags admin
// @Produce json
// @Param id path int true "Collection ID"
// @Success 204 {object} nil
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/collections/{id} [delete]
func (sv *Server) adminDeleteCollection(c *gin.Context) {
	var colID models.UriIDParam
	if err := c.ShouldBindUri(&colID); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode,
			err))
		return
	}

	_, err := sv.repo.GetCollectionByID(c, uuid.MustParse(colID.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode,
				fmt.Errorf("collection with ID %s not found", colID.ID)))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode,
			err))
		return
	}

	err = sv.repo.DeleteCollection(c, uuid.MustParse(colID.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode,
			err))
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary Get product ratings
// @Description Get ratings for a specific product
// @Tags ratings
// @Accept json
// @Produce json
// @Param productId path string true "Product ID"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Success 200 {object} ApiResponse[[]ProductRatingModel]
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/ratings [get]
func (s *Server) adminGetRatings(c *gin.Context) {
	var queries models.RatingsQueryParams
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(400, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	sqlParams := repository.GetProductRatingsParams{
		Limit:  queries.PageSize,
		Offset: (queries.Page - 1) * queries.PageSize,
	}

	if queries.Status != nil {
		switch *queries.Status {
		case "approved":
			sqlParams.IsApproved = utils.BoolPtr(true)
		case "rejected":
			sqlParams.IsApproved = utils.BoolPtr(false)
			sqlParams.IsVisible = utils.BoolPtr(false)
		case "pending":
			sqlParams.IsApproved = nil
		default:
		}
	}
	ratings, err := s.repo.GetProductRatings(c, sqlParams)
	if err != nil {
		c.JSON(400, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	ratingsCount, err := s.repo.CountProductRatings(c, pgtype.UUID{
		Bytes: uuid.Nil,
		Valid: false,
	})
	if err != nil {
		c.JSON(400, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	productRatings := make([]dto.ProductRatingDetail, 0)
	for _, rating := range ratings {
		ratingPoint, _ := rating.Rating.Float64Value()
		prIdx := -1
		for i, pr := range productRatings {
			if pr.ID == rating.ID.String() {
				prIdx = i
				break
			}
		}
		if prIdx != -1 && rating.ImageID != nil {
			productRatings[prIdx].Images = append(productRatings[prIdx].Images, dto.RatingImage{
				ID:  *rating.ImageID,
				URL: *rating.ImageUrl,
			})
			continue
		}
		model := dto.ProductRatingDetail{
			ID:               rating.ID.String(),
			UserID:           rating.UserID.String(),
			FirstName:        rating.FirstName,
			LastName:         rating.LastName,
			ProductName:      rating.ProductName,
			Rating:           ratingPoint.Float64,
			IsVisible:        rating.IsVisible,
			IsApproved:       rating.IsApproved,
			ReviewTitle:      *rating.ReviewTitle,
			ReviewContent:    *rating.ReviewContent,
			VerifiedPurchase: rating.VerifiedPurchase,
			Count:            ratingsCount,
		}
		if rating.ImageID != nil {
			model.Images = append(model.Images, dto.RatingImage{
				ID:  *rating.ImageID,
				URL: *rating.ImageUrl,
			})
		}
		productRatings = append(productRatings, model)
	}
	c.JSON(200, dto.CreateDataResp(c, productRatings, dto.CreatePagination(queries.Page, queries.PageSize, ratingsCount), nil))
}

// @Summary Get order ratings
// @Description Get ratings for a specific order
// @Tags ratings
// @Accept json
// @Produce json
// @Param orderId path string true "Order ID"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[[]ProductRatingModel]
// @Failure 400 {object} ErrorResp
// @Failure 403 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /ratings/orders/{orderId} [get]
func (s *Server) adminGetOrderRatings(c *gin.Context) {
	auth, _ := c.MustGet(constants.AuthPayLoad).(*auth.TokenPayload)

	var param struct {
		OrderID string `uri:"orderId" binding:"required,uuid"`
	}
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(400, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	orderItems, err := s.repo.GetOrderItemsByOrderID(c, uuid.MustParse(param.OrderID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	if len(orderItems) == 0 {
		c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, nil))
		return
	}
	if orderItems[0].UserID != auth.UserID {
		c.JSON(http.StatusForbidden, dto.CreateErr(PermissionDeniedCode, nil))
		return
	}
	orderItemIds := make([]uuid.UUID, len(orderItems))
	for i, orderItem := range orderItems {
		orderItemIds[i] = orderItem.OrderItemID
	}
	ratings, err := s.repo.GetProductRatingsByOrderItemIDs(c, orderItemIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(200, dto.CreateDataResp(c, ratings, nil, nil))
}

// @Summary Delete a rating
// @Description Delete a product rating by ID
// @Tags admin, ratings
// @Accept json
// @Produce json
// @Param id path string true "Rating ID"
// @Security BearerAuth
// @Success 204 {object} nil
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/ratings/{id} [delete]
func (sv *Server) adminDeleteRating(c *gin.Context) {
	// Parse the rating ID from the URL
	var param struct {
		ID string `uri:"id" binding:"required,uuid"`
	}

	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	// Convert the ID to UUID
	ratingID, err := uuid.Parse(param.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	// Check if rating exists first
	_, err = sv.repo.GetProductRating(c, ratingID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, err))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	// Delete the rating
	err = sv.repo.DeleteProductRating(c, ratingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Approve a rating
// @Description Approve a product rating by ID
// @Tags admin, ratings
// @Accept json
// @Produce json
// @Param id path string true "Rating ID"
// @Security BearerAuth
// @Success 204 {object} nil
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/ratings/{id}/approve [post]
func (sv *Server) adminApproveRating(c *gin.Context) {
	// Parse the rating ID from the URL
	var param struct {
		ID string `uri:"id" binding:"required,uuid"`
	}

	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	// Convert the ID to UUID
	ratingID, err := uuid.Parse(param.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	// Check if rating exists first
	rating, err := sv.repo.GetProductRating(c, ratingID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, err))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	// Set IsApproved to true
	isApproved := true

	// Update the rating
	_, err = sv.repo.UpdateProductRating(c, repository.UpdateProductRatingParams{
		ID:         rating.ID,
		IsApproved: &isApproved,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Ban a user from rating
// @Description Ban a user from rating by setting their rating to invisible
// @Tags admin, ratings
// @Accept json
// @Produce json
// @Param id path string true "Rating ID"
// @Security BearerAuth
// @Success 204 {object} nil
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/ratings/{id}/ban [post]
func (sv *Server) adminBanUserRating(c *gin.Context) {
	// Parse the rating ID from the URL
	var param struct {
		ID string `uri:"id" binding:"required,uuid"`
	}

	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	// Convert the ID to UUID
	ratingID, err := uuid.Parse(param.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	// Check if rating exists first
	rating, err := sv.repo.GetProductRating(c, ratingID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, err))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	// Set IsVisible to false
	isVisible := false

	// Update the rating
	_, err = sv.repo.UpdateProductRating(c, repository.UpdateProductRatingParams{
		ID:        rating.ID,
		IsVisible: &isVisible,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.Status(http.StatusNoContent)
}

// adminGetDiscounts godoc
// @Summary Get all discounts
// @Description Get all discounts
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Param search query string false "Search by code"
// @Param discountType query string false "Discount type" default(percentage)
// @Param isActive query bool false "Is active" default(true)
// @Success 200 {object} ApiResponse[[]DiscountListItemResponseModel]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/discounts [get]
func (sv *Server) adminGetDiscounts(c *gin.Context) {
	var queries models.DiscountListQuery
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	// Get all discounts
	sqlParams := repository.GetDiscountsParams{
		Limit:  queries.PageSize,
		Offset: (queries.Page - 1) * queries.PageSize,
		// Search:       queries.Search,
		// DiscountType: queries.DiscountType,
		// IsActive:     queries.IsActive,
	}

	if queries.FromDate != nil {
		sqlParams.FromDate = utils.GetPgTypeTimestamp(*queries.FromDate)
	}
	if queries.ToDate != nil {
		sqlParams.ToDate = utils.GetPgTypeTimestamp(*queries.ToDate)
	}

	discounts, err := sv.repo.GetDiscounts(c, sqlParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	total, err := sv.repo.CountDiscounts(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	listData := make([]dto.DiscountListItem, len(discounts))
	for i, discount := range discounts {
		discountValue, _ := discount.DiscountValue.Float64Value()

		listData[i] = dto.DiscountListItem{
			ID:            discount.ID.String(),
			Code:          discount.Code,
			DiscountType:  string(discount.DiscountType),
			DiscountValue: discountValue.Float64,
			IsActive:      discount.IsActive,
			TimeUsed:      discount.TimesUsed,
			UsageLimit:    discount.UsageLimit,
			Description:   discount.Description,
			ValidFrom:     discount.ValidFrom.String(),
			CreatedAt:     discount.CreatedAt.String(),
			UpdatedAt:     discount.UpdatedAt.String(),
		}
		if discount.ValidUntil.Valid {
			listData[i].ValidUntil = discount.ValidUntil.Time.String()
		}
		if discount.MinOrderValue.Valid {
			minPurchaseAmount, _ := discount.MinOrderValue.Float64Value()
			listData[i].MinPurchase = &minPurchaseAmount.Float64
		}

		if discount.MaxDiscountAmount.Valid {
			maxDiscountAmount, _ := discount.MaxDiscountAmount.Float64Value()
			listData[i].MaxDiscount = &maxDiscountAmount.Float64
		}
	}
	pagination := dto.CreatePagination(queries.Page, queries.PageSize, total)

	c.JSON(http.StatusOK, dto.CreateDataResp(c, listData, pagination, nil))
}

// adminCreateDiscount godoc
// @Summary Create a new discount
// @Description Create a new discount
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param input body CreateDiscountRequest true "Discount info"
// @Success 201 {object} ApiResponse[DiscountDetailResponseModel]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/discounts [post]
func (sv *Server) adminCreateDiscount(c *gin.Context) {
	// Create a new discount
	var req models.AddDiscount
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	sqlParams := repository.InsertDiscountParams{
		Code:          req.Code,
		DiscountType:  repository.DiscountType(req.DiscountType),
		DiscountValue: utils.GetPgNumericFromFloat(req.DiscountValue),
		IsActive:      req.IsActive,
		UsageLimit:    req.UsageLimit,
		Description:   req.Description,
		ValidFrom:     utils.GetPgTypeTimestamp(req.ValidFrom),
		ValidUntil:    utils.GetPgTypeTimestamp(req.ValidUntil),
		Name:          req.Name,
		UsagePerUser:  req.UsagePerUser,
		IsStackable:   req.IsStackable,
		Priority:      req.Priority,
	}

	if req.MinOrderValue != nil {
		sqlParams.MinOrderValue = utils.GetPgNumericFromFloat(*req.MinOrderValue)
	}
	if req.MaxDiscountAmount != nil {
		sqlParams.MaxDiscountAmount = utils.GetPgNumericFromFloat(*req.MaxDiscountAmount)
	}

	discount, err := sv.repo.InsertDiscount(c, sqlParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusCreated, dto.CreateDataResp(c, discount.String(), nil, nil))
}

// adminUpdateDiscount godoc
// @Summary Update discount by ID
// @Description Update discount by ID
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param id path string true "Discount ID"
// @Param input body UpdateDiscountRequest true "Discount info"
// @Success 200 {object} ApiResponse[string]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /discounts/{id} [put]
func (sv *Server) adminUpdateDiscount(c *gin.Context) {
	// Update discount by ID
	var param models.UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	var req models.UpdateDiscount
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	discount, err := sv.repo.GetDiscountByID(c, uuid.MustParse(param.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	sqlParams := repository.UpdateDiscountParams{
		ID:           discount.ID,
		Name:         req.Name,
		Code:         req.Code,
		IsActive:     req.IsActive,
		UsageLimit:   req.UsageLimit,
		IsStackable:  req.IsStackable,
		Priority:     req.Priority,
		Description:  req.Description,
		UsagePerUser: req.UsagePerUser,
	}

	if req.DiscountType != nil {
		sqlParams.DiscountType.Scan(req.DiscountType)
	}
	if req.DiscountValue != nil {
		sqlParams.DiscountValue = utils.GetPgNumericFromFloat(*req.DiscountValue)
	}
	if req.ValidFrom != nil {
		sqlParams.ValidFrom = utils.GetPgTypeTimestamp(*req.ValidFrom)
	}
	if req.ValidUntil != nil {
		sqlParams.ValidUntil = utils.GetPgTypeTimestamp(*req.ValidUntil)
	}
	if req.MinOrderValue != nil {
		sqlParams.MinOrderValue = utils.GetPgNumericFromFloat(*req.MinOrderValue)
	}
	if req.MaxDiscountAmount != nil {
		sqlParams.MaxDiscountAmount = utils.GetPgNumericFromFloat(*req.MaxDiscountAmount)
	}

	updated, err := sv.repo.UpdateDiscount(c, sqlParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, updated, nil, nil))
}

// adminDeleteDiscount godoc
// @Summary Delete discount by ID
// @Description Delete discount by ID
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param id path string true "Discount ID"
// @Success 204
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/discounts/{id} [delete]
func (sv *Server) adminDeleteDiscount(c *gin.Context) {
	// Delete discount by ID
	var param models.UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	err := sv.repo.DeleteDiscount(c, uuid.MustParse(param.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.Status(http.StatusNoContent)
}

// adminAddDiscountRule godoc
// @Summary Add a discount rule to a discount
// @Description Add a discount rule to a discount
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param id path string true "Discount ID"
// @Param input body AddDiscountRuleRequest true "Discount rule info"
// @Success 201 {object} ApiResponse[string]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/discounts/{id}/rules [post]
func (sv *Server) adminAddDiscountRule(c *gin.Context) {
	// Add a discount rule to a discount
	var param models.UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	var req models.AddDiscountRule
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	var ruleVal []byte
	switch req.RuleType {
	case "first_time_buyer":
		var ruleValue models.FirstTimeBuyerRule
		if err := mapstructure.Decode(req.RuleValue, &ruleValue); err != nil {
			c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
			return
		}
		bs, err := json.Marshal(ruleValue)

		if err != nil {
			c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
			return
		}
		ruleVal = bs
	case "product":
		var ruleValue models.ProductRule
		if err := mapstructure.Decode(req.RuleValue, &ruleValue); err != nil {
			c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
			return
		}
		bs, err := json.Marshal(ruleValue)

		if err != nil {
			c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
			return
		}
		ruleVal = bs
	case "category":
		var ruleValue models.CategoryRule
		if err := mapstructure.Decode(req.RuleValue, &ruleValue); err != nil {
			c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
			return
		}
		bs, err := json.Marshal(ruleValue)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
			return
		}
		ruleVal = bs
	case "customer_segment":
		var ruleValue models.CustomerSegmentRule
		if err := mapstructure.Decode(req.RuleValue, &ruleValue); err != nil {
			c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
			return
		}
		bs, err := json.Marshal(ruleValue)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
			return
		}
		ruleVal = bs
	}

	sqlParams := repository.InsertDiscountRuleParams{
		DiscountID: uuid.MustParse(param.ID),
		RuleType:   req.RuleType,
		RuleValue:  ruleVal,
	}

	rule, err := sv.repo.InsertDiscountRule(c, sqlParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusCreated, dto.CreateDataResp(c, rule, nil, nil))
}

// adminGetDiscountRules godoc
// @Summary Get all discount rules for a discount
// @Description Get all discount rules for a specific discount
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param id path string true "Discount ID"
// @Success 200 {object} ApiResponse[[]DiscountRule]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/discounts/{id}/rules [get]
func (sv *Server) adminGetDiscountRules(c *gin.Context) {
	var param models.UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	rules, err := sv.repo.GetDiscountRules(c, uuid.MustParse(param.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	var ruleDetails []dto.DiscountRuleDetail
	for _, rule := range rules {
		ruleDetail, err := dto.MapToDiscountRuleDetail(rule)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
			return
		}
		ruleDetails = append(ruleDetails, ruleDetail)
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, ruleDetails, nil, nil))
}

// adminGetDiscountRuleByID godoc
// @Summary Get a specific discount rule by ID
// @Description Get a specific discount rule by ID
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param id path string true "Discount ID"
// @Param ruleId path string true "Rule ID"
// @Success 200 {object} ApiResponse[DiscountRule]
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/discounts/{id}/rules/{ruleId} [get]
func (sv *Server) adminGetDiscountRuleByID(c *gin.Context) {
	var param models.UriRuleIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	rule, err := sv.repo.GetDiscountRuleByID(c, uuid.MustParse(param.RuleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	ruleDetail, err := dto.MapToDiscountRuleDetail(rule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, ruleDetail, nil, nil))
}

// adminUpdateDiscountRule godoc
// @Summary Update a discount rule
// @Description Update a discount rule
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param id path string true "Discount ID"
// @Param ruleId path string true "Rule ID"
// @Param input body UpdateDiscountRuleModel true "Updated discount rule info"
// @Success 200 {object} ApiResponse[string]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/discounts/{id}/rules/{ruleId} [put]
func (sv *Server) adminUpdateDiscountRule(c *gin.Context) {
	var param models.UriRuleIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	var req models.UpdateDiscountRule
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	sqlParams := repository.UpdateDiscountRuleParams{
		ID: uuid.MustParse(param.RuleID),
	}

	if req.RuleType != nil {
		sqlParams.RuleType = req.RuleType
	}
	if req.RuleValue != nil {
		ruleValueBytes, err := json.Marshal(req.RuleValue)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
			return
		}
		sqlParams.RuleValue = ruleValueBytes
	}

	rule, err := sv.repo.UpdateDiscountRule(c, sqlParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, rule, nil, nil))
}

// adminDeleteDiscountRule godoc
// @Summary Delete a discount rule
// @Description Delete a discount rule
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param id path string true "Discount ID"
// @Param ruleId path string true "Rule ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/discounts/{id}/rules/{ruleId} [delete]
func (sv *Server) adminDeleteDiscountRule(c *gin.Context) {
	var param models.UriRuleIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	err := sv.repo.DeleteDiscountRule(c, uuid.MustParse(param.RuleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.Status(http.StatusNoContent)
}
