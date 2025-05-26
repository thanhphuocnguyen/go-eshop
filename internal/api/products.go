package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"golang.org/x/sync/errgroup"
)

// @Summary Create a new product
// @Schemes http
// @Description create a new product with the input payload
// @Tags products
// @Accept json
// @Param input body repository.CreateProductTxArgs true "Product input"
// @Produce json
// @Success 200 {object} ApiResponse[ProductListModel]
// @Failure 400 {object} ApiResponse[gin.H]
// @Failure 500 {object} ApiResponse[gin.H]
// @Router /products [post]
func (sv *Server) addProductHandler(c *gin.Context) {
	var req repository.CreateProductTxArgs
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[ProductListModel](InvalidBodyCode, "", err))
		return
	}

	productID, err := sv.repo.CreateProductTx(c, req)

	if err != nil {
		log.Error().Err(err).Msg("CreateProduct")
		c.JSON(http.StatusInternalServerError, createErrorResponse[ProductListModel](InternalServerErrorCode, "", err))
		return
	}

	c.JSON(http.StatusCreated, createSuccessResponse(c, productID.String(), "", nil, nil))
}

// @Summary Get a product detail by ID
// @Schemes http
// @Description get a product detail by ID
// @Tags products
// @Accept json
// @Param product_id path int true "Product ID"
// @Produce json
// @Success 200 {object} ApiResponse[gin.H]
// @Failure 404 {object} ApiResponse[gin.H]
// @Failure 500 {object} ApiResponse[gin.H]
// @Router /products/{product_id} [get]
func (sv *Server) getProductDetailHandler(c *gin.Context) {
	var params URISlugParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[gin.H](InvalidBodyCode, "", err))
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
			c.JSON(http.StatusNotFound, createErrorResponse[gin.H](NotFoundCode, "", err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "", err))
		return
	}

	prodID := productRow.ProductID
	productDetail := mapToProductResponse(productRow)

	variantRows, err := sv.repo.GetProductVariants(c, repository.GetProductVariantsParams{
		ProductID: prodID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "", err))
		return
	}

	entityIds := make([]uuid.UUID, 0)
	idMap := make(map[uuid.UUID]bool)

	for _, row := range variantRows {
		if _, ok := idMap[row.ID]; !ok {
			idMap[row.ID] = true
			entityIds = append(entityIds, row.ID)
		}
	}
	// Add the product ID to the entityIds slice
	// This ensures that the product ID is included in the list of entity IDs
	// when fetching product images
	entityIds = append(entityIds, prodID)

	images, err := sv.repo.GetProductImagesAssigned(c, entityIds)
	imageResp := mapToProductImages(prodID, images)
	variants := mapToVariantResp(variantRows)
	productDetail.Variants = variants
	productDetail.ProductImages = imageResp
	c.JSON(http.StatusOK, createSuccessResponse(c, productDetail, "product retrieved", nil, nil))
}

// @Summary Get list of products
// @Schemes http
// @Description get list of products
// @Tags products
// @Accept json
// @Param page query int true "Page number"
// @Param pageSize query int true "Page size"
// @Produce json
// @Success 200 {array} ApiResponse[ProductListModel]
// @Failure 404 {object} ApiResponse[ProductListModel]
// @Failure 500 {object} ApiResponse[ProductListModel]
// @Router /products [get]
func (sv *Server) getProductsHandler(c *gin.Context) {
	var queries ProductQueries
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[[]ProductListModel](InvalidBodyCode, "", err))
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
		c.JSON(http.StatusInternalServerError, createErrorResponse[[]ProductListModel](InternalServerErrorCode, "Server error", err))
		return
	}

	productCnt, err := sv.repo.CountProducts(c, repository.CountProductsParams{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[[]ProductListModel](InternalServerErrorCode, "Server error", err))
		return
	}

	productResponses := make([]ProductListModel, 0)
	for _, product := range products {
		productResponses = append(productResponses, mapToListProductResponse(product))
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, productResponses, "products retrieved", &Pagination{
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
// @Param product_id path int true "Product ID"
// @Param input body repository.UpdateProductTxParams true "Product input"
// @Produce json
// @Success 200 {object} ApiResponse[ProductListModel]
// @Failure 404 {object} ApiResponse[ProductListModel]
// @Failure 500 {object} ApiResponse[ProductListModel]
// @Router /products/{product_id} [put]
func (sv *Server) updateProductHandler(c *gin.Context) {
	var param UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[ProductListModel](InvalidBodyCode, "", err))
		return
	}
	uuid, err := uuid.Parse(param.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[ProductListModel](InvalidBodyCode, "", err))
		return
	}

	var req repository.UpdateProductTxParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[ProductListModel](InvalidBodyCode, "", err))
		return
	}
	err = sv.repo.UpdateProductTx(c, uuid, req)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[ProductListModel](NotFoundCode, "", err))
			return
		}
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, struct{}{}, "product updated", nil, nil))
}

// @Summary Remove a product by ID
// @Schemes http
// @Description remove a product by ID
// @Tags products
// @Accept json
// @Param product_id path int true "Product ID"
// @Produce json
// @Success 200 {object} ApiResponse[bool]
// @Failure 404 {object} ApiResponse[gin.H]
// @Failure 500 {object} ApiResponse[gin.H]
// @Router /products/{product_id} [delete]
func (sv *Server) deleteProductHandler(c *gin.Context) {
	var params UriIDParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[bool](InvalidBodyCode, "", err))
		return
	}

	product, err := sv.repo.GetProductByID(c, repository.GetProductByIDParams{ID: uuid.MustParse(params.ID)})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[bool](NotFoundCode, "", err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
		return
	}

	// Remove the product image
	images, err := sv.repo.GetProductImagesAssigned(c, []uuid.UUID{product.ID})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[bool](NotFoundCode, "", err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
		return
	}
	errGroup, _ := errgroup.WithContext(c)
	for _, image := range images {
		errGroup.Go(func() (err error) {
			img, err := sv.repo.GetImageFromID(c, repository.GetImageFromIDParams{
				ID:         image.ID,
				EntityType: string(repository.EntityTypeProductVariant),
			})

			if err != nil {
				if errors.Is(err, repository.ErrRecordNotFound) {
					c.JSON(http.StatusNotFound, createErrorResponse[bool](NotFoundCode, "", err))
					return
				}
				c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
				return
			}
			// Remove image from storage
			msg, err := sv.removeImageUtil(c, img.ExternalID)
			if err != nil {
				return fmt.Errorf("failed to remove image: %w, reason: %s", err, msg)
			}

			// Remove image from product
			err = sv.repo.DeleteProductImage(c, image.ID)
			if err != nil {
				if errors.Is(err, repository.ErrRecordNotFound) {
					c.JSON(http.StatusNotFound, createErrorResponse[bool](NotFoundCode, "", err))
					return
				}
				c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
				return
			}
			return
		})
	}

	err = errGroup.Wait()

	err = sv.repo.DeleteProduct(c, uuid.MustParse(params.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, true, "Removed!", nil, nil))
}
