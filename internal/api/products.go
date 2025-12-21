package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
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
func (s *Server) getProductById(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		RespondBadRequest(w, InvalidBodyCode, errors.New("id parameter is required"))
		return
	}

	sqlParams := repository.GetProductDetailParams{}
	err := uuid.Validate(idParam)
	if err == nil {
		sqlParams.ID = uuid.MustParse(idParam)
	} else {
		sqlParams.Slug = idParam
	}

	productRow, err := s.repo.GetProductDetail(r.Context(), sqlParams)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, err)
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	productDetail := dto.MapToProductDetailResponse(productRow)

	RespondSuccess(w, r, productDetail)
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
func (s *Server) getProducts(w http.ResponseWriter, r *http.Request) {
	paginationQuery := ParsePaginationQuery(r)

	queryParams := r.URL.Query()
	var queries models.ProductQuery
	queries.Page = paginationQuery.Page
	queries.PageSize = paginationQuery.PageSize

	// Parse search parameter
	if search := queryParams.Get("search"); search != "" {
		queries.Search = &search
	}

	// Parse brandIds parameter
	if brandIDs := queryParams["brandIds"]; len(brandIDs) > 0 {
		queries.BrandIDs = &brandIDs
	}

	// Parse categoryIds parameter
	if categoryIDs := queryParams["categoryIds"]; len(categoryIDs) > 0 {
		queries.CategoryIDs = &categoryIDs
	}

	validate := validator.New()
	if err := validate.Struct(&queries); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
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

	products, err := s.repo.GetProductList(r.Context(), dbParams)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	productCnt, err := s.repo.CountProducts(r.Context(), repository.CountProductsParams{})
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	productResponses := make([]dto.ProductSummary, 0)
	for _, product := range products {
		productResponses = append(productResponses, dto.MapToShopProductResponse(product))
	}

	RespondSuccessWithPagination(w, r, productResponses, dto.CreatePagination(queries.Page, queries.PageSize, productCnt))
}

// Setup product-related routes
func (s *Server) addProductRoutes(rg chi.Router) {
	rg.Route("/products", func(r chi.Router) {
		r.Get("/", s.getProducts)
		r.Get("/{id}", s.getProductById)
		r.Get("/{id}/ratings", s.getRatingsByProduct)
	})
}
