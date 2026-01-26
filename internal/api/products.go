package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/textquerytype"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
)

// parseFloatParam parses a string parameter into a float64 pointer
func parseFloatParam(param string) (*float64, error) {
	if param == "" {
		return nil, nil
	}

	value, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return nil, err
	}

	return &value, nil
}

// @Summary Get a product detail by ID
// @Schemes http
// @Description get a product detail by ID
// @Tags products
// @Accept json
// @Param id path string true "Product ID"
// @Produce json
// @Success 200 {object} dto.ApiResponse[dto.ProductDetail]
// @Failure 404 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /products/{id} [get]
func (s *Server) getProductById(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		RespondBadRequest(w, InvalidBodyCode, errors.New("id parameter is required"))
		return
	}
	c := r.Context()
	sqlParams := repository.GetProductDetailParams{}
	err := uuid.Validate(idParam)
	if err == nil {
		sqlParams.ID = uuid.MustParse(idParam)
	} else {
		sqlParams.Slug = idParam
	}

	productRow, err := s.repo.GetProductDetail(c, sqlParams)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, err)
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	productDetail := dto.MapToProductDetailResponse(productRow)

	RespondSuccess(w, productDetail)
}

// @Summary Get list of products
// @Schemes http
// @Description get list of products
// @Tags products
// @Accept json
// @Param page query int true "Page number"
// @Param pageSize query int true "Page size"
// @Produce json
// @Success 200 {array} dto.ApiResponse[[]dto.ProductSummary]
// @Failure 404 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /products [get]
func (s *Server) getProducts(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	paginationQuery := ParsePaginationQuery(r)
	queryParams := r.URL.Query()
	var queries models.ProductQuery
	queries.Page = paginationQuery.Page
	queries.PageSize = paginationQuery.PageSize

	// Parse search parameter
	if search := queryParams.Get("search"); search != "" {
		queries.Search = &search
	}
	if brand := queryParams.Get("brand"); brand != "" {
		queries.Brand = &brand
	}
	if priceFrom := queryParams.Get("priceFrom"); priceFrom != "" {

		if pf, err := parseFloatParam(priceFrom); err == nil {
			queries.PriceFrom = pf
		}
	}
	if priceTo := queryParams.Get("priceTo"); priceTo != "" {
		if pt, err := parseFloatParam(priceTo); err == nil {
			queries.PriceTo = pt
		}
	}
	if categories := queryParams["categories"]; len(categories) > 0 {
		queries.Categories = categories
	}
	if collections := queryParams["collections"]; len(collections) > 0 {
		queries.Collections = collections
	}
	if brand := queryParams["brand"]; len(brand) > 0 {
		queries.Brand = &brand[0]
	}
	// Parse attributes parameter
	attributes := make([]string, 0)
	if attrs := queryParams["attributes"]; len(attrs) > 0 {
		for _, attr := range attrs {
			attributes = append(attributes, attr)
		}
		queries.Attributes = attributes
	}
	from := (queries.Page - 1) * queries.PageSize
	rangeQuery := map[string]types.RangeQuery{
		"price": &types.NumberRangeQuery{
			Gte: (*types.Float64)(queries.PriceFrom),
			Lte: (*types.Float64)(queries.PriceTo),
		},
	}
	mustQuery := []types.Query{}
	if queries.Search != nil && len(*queries.Search) > 0 {
		mustQuery = append(mustQuery, types.Query{
			MultiMatch: &types.MultiMatchQuery{
				Query:     *queries.Search,
				Fields:    []string{"name^4", "brand^2", "description", "shortDescription"},
				Type:      &textquerytype.Bestfields,
				Fuzziness: "AUTO",
			},
		})
	}

	filterQuery := []types.Query{
		{
			Term: map[string]types.TermQuery{
				"inStock": {
					Value: true,
				},
			},
		},
		{
			Range: rangeQuery,
		},
	}
	if len(queries.Attributes) > 0 {
		filterQuery = append(filterQuery, types.Query{
			Terms: &types.TermsQuery{
				TermsQuery: map[string]types.TermsQueryField{
					"attributes.keyword": queries.Attributes,
				},
			},
		})
	}
	if len(queries.Categories) > 0 {
		filterQuery = append(filterQuery, types.Query{
			Terms: &types.TermsQuery{
				TermsQuery: map[string]types.TermsQueryField{
					"categories.keyword": queries.Categories,
				},
			},
		})
	}
	if len(queries.Collections) > 0 {
		filterQuery = append(filterQuery, types.Query{
			Terms: &types.TermsQuery{
				TermsQuery: map[string]types.TermsQueryField{
					"collections.keyword": queries.Collections,
				},
			},
		})
	}

	var boostFeature float32 = 1.5
	shouldQuery := []types.Query{{RankFeature: &types.RankFeatureQuery{Boost: &boostFeature, Field: "popularity_score"}}}
	if queries.Brand != nil && len(*queries.Brand) > 0 {
		var boostBrand float32 = 3.0
		shouldQuery = append(shouldQuery, types.Query{
			Term: map[string]types.TermQuery{
				"brand.keyword": {
					Value: queries.Brand,
					Boost: &boostBrand,
				},
			},
		})
	}

	searchSize := int(queries.PageSize)
	searchFrom := int(from)
	q := &search.Request{
		Size: &searchSize,
		From: &searchFrom,
		Query: &types.Query{
			Bool: &types.BoolQuery{
				Must:               mustQuery,
				Filter:             filterQuery,
				Should:             shouldQuery,
				MinimumShouldMatch: 0,
			},
		},
	}
	esProducts, err := s.elasticClient.SearchProducts(c, q) // Parse categoryIds parameter
	if err == nil {
		RespondSuccessWithPagination(w, esProducts, dto.CreatePagination(queries.Page, queries.PageSize, 0))
		return
	}

	// Fallback to DB search
	if err := s.validator.Struct(&queries); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	dbParams := repository.GetProductListParams{
		Limit:  int64(queries.PageSize),
		Offset: int64((queries.Page - 1) * queries.PageSize),
	}

	if queries.Search != nil && len(*queries.Search) > 0 {
		search := *queries.Search
		search = strings.ReplaceAll(search, " ", "%")
		search = strings.ReplaceAll(search, ",", "%")
		search = strings.ReplaceAll(search, ":", "%")
		search = "%" + search + "%"
		dbParams.Search = &search
	}

	products, err := s.repo.GetProductList(c, dbParams)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	productCnt, err := s.repo.CountProducts(c, repository.CountProductsParams{})
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	productResponses := make([]dto.ProductSummary, 0)
	for _, product := range products {
		productResponses = append(productResponses, dto.MapToShopProductResponse(product))
	}

	RespondSuccessWithPagination(w, productResponses, dto.CreatePagination(queries.Page, queries.PageSize, productCnt))
}

// Setup product-related routes
func (s *Server) addProductRoutes(r chi.Router) {
	r.Route("/products", func(r chi.Router) {
		r.Get("/", s.getProducts)
		r.Get("/{id}", s.getProductById)
		r.Get("/{id}/variants", s.getProductVariants)
		r.Get("/{id}/variants/{variantId}", s.getVariantByProductId)
		r.Get("/{id}/ratings", s.getRatingsByProduct)
	})
}
