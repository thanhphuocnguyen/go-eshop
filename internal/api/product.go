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

type ProductQueries struct {
	Page         int64   `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize     int64   `form:"page_size,default=20" binding:"omitempty,min=1,max=100"`
	Search       *string `form:"search" binding:"omitempty,max=1000"`
	CategoryID   *string `form:"category_id" binding:"omitempty,uuid"`
	BrandID      *string `form:"brand_id" binding:"omitempty,uuid"`
	CollectionID *string `form:"collection_id" binding:"omitempty,uuid"`
}

type ProductAttributeDetail struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	ValueObject AttributeValue `json:"value_object"`
}

type ImageAssignment struct {
	ID           string `json:"id"`
	EntityID     string `json:"entity_id"`
	EntityType   string `json:"entity_type"`
	DisplayOrder int16  `json:"display_order"`
	Role         string `json:"role"`
}

type ProductImageModel struct {
	ID                 string            `json:"id"`
	Url                string            `json:"url"`
	ExternalID         string            `json:"external_id"`
	Role               string            `json:"role"`
	VariantAssignments []ImageAssignment `json:"assignments"`
}

type ProductVariantModel struct {
	ID         string                   `json:"id"`
	Price      float64                  `json:"price"`
	StockQty   int32                    `json:"stock_qty"`
	IsActive   bool                     `json:"is_active"`
	Sku        *string                  `json:"sku,omitempty"`
	Attributes []ProductAttributeDetail `json:"attributes"`
}

type UpdateProductImageAssignments struct {
	ID           int32  `json:"id"`
	EntityID     string `json:"entity_id"`
	DisplayOrder int16  `json:"display_order"`
}

type GeneralCategoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProductModel struct {
	ID               string                   `json:"id"`
	Name             string                   `json:"name"`
	Description      string                   `json:"description"`
	ShortDescription *string                  `json:"short_description"`
	Attributes       []string                 `json:"attributes"`
	BasePrice        float64                  `json:"price,omitzero"`
	BaseSku          string                   `json:"sku"`
	UpdatedAt        string                   `json:"updated_at"`
	IsActive         bool                     `json:"is_active"`
	Slug             string                   `json:"slug"`
	CreatedAt        string                   `json:"created_at"`
	RatingCount      int32                    `json:"rating_count"`
	OneStarCount     int32                    `json:"one_star_count"`
	TwoStarCount     int32                    `json:"two_star_count"`
	ThreeStarCount   int32                    `json:"three_star_count"`
	FourStarCount    int32                    `json:"four_star_count"`
	FiveStarCount    int32                    `json:"five_star_count"`
	Variants         []ProductVariantModel    `json:"variants"`
	ProductImages    []ProductImageModel      `json:"product_images"`
	Collection       *GeneralCategoryResponse `json:"collection,omitempty"`
	Brand            *GeneralCategoryResponse `json:"brand,omitempty"`
	Category         *GeneralCategoryResponse `json:"category,omitempty"`
}

type ProductListModel struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	VariantCount int64   `json:"variant_count,omitzero"`
	MinPrice     float64 `json:"min_price,omitzero"`
	MaxPrice     float64 `json:"max_price,omitzero"`
	Slug         string  `json:"slug,omitempty"`
	Sku          string  `json:"sku"`
	ImgUrl       *string `json:"image_url,omitempty"`
	ImgID        *string `json:"image_id,omitempty"`
	CreatedAt    string  `json:"created_at,omitempty"`
	UpdatedAt    string  `json:"updated_at,omitempty"`
}

type ProductCreateResp struct {
	ID string `json:"id"`
}

// ------------------------------ Handlers ------------------------------

// @Summary Create a new product
// @Schemes http
// @Description create a new product with the input payload
// @Tags products
// @Accept json
// @Param input body repository.CreateProductTxArgs true "Product input"
// @Produce json
// @Success 200 {object} ApiResponse[ProductListModel]
// @Failure 400 {object} ApiResponse[ProductListModel]
// @Failure 500 {object} ApiResponse[ProductListModel]
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

	resp := ProductCreateResp{
		ID: productID.String(),
	}

	c.JSON(http.StatusCreated, createSuccessResponse(c, resp, "", nil, nil))
}

// @Summary Get a product detail by ID
// @Schemes http
// @Description get a product detail by ID
// @Tags products
// @Accept json
// @Param product_id path int true "Product ID"
// @Produce json
// @Success 200 {object} ApiResponse[ProductListModel]
// @Failure 404 {object} ApiResponse[ProductListModel]
// @Failure 500 {object} ApiResponse[ProductListModel]
// @Router /products/{product_id} [get]
func (sv *Server) getProductDetailHandler(c *gin.Context) {
	var params ProductDetailParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[ProductModel](InvalidBodyCode, "", err))
		return
	}
	sqlParams := repository.GetProductDetailParams{}
	err := uuid.Validate(params.ID)
	if err == nil {
		sqlParams.ID = uuid.MustParse(params.ID)
	} else {
		sqlParams.Slug = params.ID
	}

	productRows, err := sv.repo.GetProductDetail(c, sqlParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[ProductModel](InternalServerErrorCode, "", err))
		return
	}

	if len(productRows) == 0 {
		c.JSON(http.StatusNotFound, createErrorResponse[ProductModel](NotFoundCode, "", errors.New("product not found")))
		return
	}
	prodID := productRows[0].ProductID
	productDetail := mapToProductResponse(productRows)

	variantRows, err := sv.repo.GetProductVariants(c, repository.GetProductVariantsParams{
		ProductID: prodID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[ProductModel](InternalServerErrorCode, "", err))
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
// @Param page_size query int true "Page size"
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

	if queries.CategoryID != nil {
		dbParams.CategoryIds = []uuid.UUID{uuid.MustParse(*queries.CategoryID)}
	}

	if queries.CollectionID != nil {
		dbParams.CollectionID = []uuid.UUID{uuid.MustParse(*queries.CollectionID)}
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
	var param URIParam
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

	c.JSON(http.StatusOK, createSuccessResponse(c, ProductCreateResp{ID: uuid.String()}, "product updated", nil, nil))
}

// @Summary Remove a product by ID
// @Schemes http
// @Description remove a product by ID
// @Tags products
// @Accept json
// @Param product_id path int true "Product ID"
// @Produce json
// @Success 200 {object} ApiResponse[bool]
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /products/{product_id} [delete]
func (sv *Server) deleteProductHandler(c *gin.Context) {
	var params URIParam
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
				EntityType: repository.VariantEntityType,
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

// ------------------------------ Mapper ------------------------------

func mapToProductResponse(productRows []repository.GetProductDetailRow) ProductModel {
	if len(productRows) == 0 {
		return ProductModel{}
	}
	product := productRows[0]
	basePrice, _ := product.BasePrice.Float64Value()
	attributes := make([]string, len(product.Attributes))
	for i, attr := range product.Attributes {
		attributes[i] = attr.String()
	}

	resp := ProductModel{
		ID:               product.ProductID.String(),
		Name:             product.Name,
		BasePrice:        basePrice.Float64,
		ShortDescription: product.ShortDescription,
		Attributes:       attributes,
		Description:      product.Description,
		BaseSku:          product.BaseSku,
		Slug:             product.Slug,
		RatingCount:      product.RatingCount,
		OneStarCount:     product.OneStarCount,
		TwoStarCount:     product.TwoStarCount,
		ThreeStarCount:   product.ThreeStarCount,
		FourStarCount:    product.FourStarCount,
		FiveStarCount:    product.FiveStarCount,

		UpdatedAt:     product.UpdatedAt.String(),
		CreatedAt:     product.CreatedAt.String(),
		IsActive:      *product.IsActive,
		Variants:      make([]ProductVariantModel, 0),
		ProductImages: make([]ProductImageModel, 0),
	}

	if product.BrandID.Valid {
		id, _ := uuid.FromBytes(product.BrandID.Bytes[:])
		resp.Brand = &GeneralCategoryResponse{
			ID:   id.String(),
			Name: *product.BrandName,
		}
	}
	if product.CategoryID.Valid {
		id, _ := uuid.FromBytes(product.CategoryID.Bytes[:])
		resp.Category = &GeneralCategoryResponse{
			ID:   id.String(),
			Name: *product.CategoryName,
		}
	}
	if product.CollectionID.Valid {
		collectionID, _ := uuid.FromBytes(product.CollectionID.Bytes[:])
		resp.Collection = &GeneralCategoryResponse{
			ID:   collectionID.String(),
			Name: *product.CollectionName,
		}
	}

	return resp
}

func mapToVariantResp(variantRows []repository.GetProductVariantsRow) []ProductVariantModel {
	variants := make([]ProductVariantModel, 0)
	for _, row := range variantRows {
		variantIdx := -1
		for i, v := range variants {
			if v.ID == row.ID.String() {
				variantIdx = i
				break
			}
		}
		if variantIdx != -1 {
			// If the variant already exists, append the attribute to the existing variant
			attrIdx := -1
			for j, a := range variants[variantIdx].Attributes {
				if a.ID == row.AttrID.String() {
					attrIdx = j
					break
				}
			}

			if attrIdx != -1 {
				// If the attribute already exists, do nothing
				continue
			}

			variants[variantIdx].Attributes = append(variants[variantIdx].Attributes, ProductAttributeDetail{
				ID:   row.AttrID.String(),
				Name: row.AttrName,
				ValueObject: AttributeValue{
					ID:           row.AttrValID,
					Code:         row.AttrValCode,
					Name:         &row.AttrValName,
					IsActive:     row.IsActive,
					DisplayOrder: &row.AttrDisplayOrder,
				},
			})
		} else {
			// If the variant does not exist, add it to the list of variants
			price, _ := row.Price.Float64Value()
			variant := ProductVariantModel{
				ID:       row.ID.String(),
				Price:    price.Float64,
				StockQty: row.Stock,
				IsActive: *row.IsActive,
				Sku:      &row.Sku,
				Attributes: []ProductAttributeDetail{
					{
						ID:   row.AttrID.String(),
						Name: row.AttrName,
						ValueObject: AttributeValue{
							ID:           row.AttrValID,
							Code:         row.AttrValCode,
							Name:         &row.AttrValName,
							IsActive:     row.IsActive,
							DisplayOrder: &row.AttrDisplayOrder,
						},
					},
				},
			}
			variants = append(variants, variant)
		}

	}
	return variants
}

func mapToProductImages(productID uuid.UUID, imageRows []repository.GetProductImagesAssignedRow) []ProductImageModel {
	// log.Debug().Msgf("mapToProductImages: %v", imageRows)
	images := make([]ProductImageModel, 0)
	for _, row := range imageRows {
		existingImageIdx := -1
		for i, image := range images {
			if image.ID == row.ID.String() {
				existingImageIdx = i
				break
			}
		}
		if existingImageIdx != -1 {
			image := ImageAssignment{
				ID:           row.ID.String(),
				EntityID:     row.EntityID.String(),
				EntityType:   row.EntityType,
				Role:         row.Role,
				DisplayOrder: row.DisplayOrder,
			}
			if row.EntityID != productID {
				// If the image already exists, append the assignment to the existing image
				images[existingImageIdx].VariantAssignments = append(images[existingImageIdx].VariantAssignments, image)
			}
		} else {
			// If the image does not exist, add it to the list of images
			image := ProductImageModel{
				ID:                 row.ID.String(),
				Url:                row.Url,
				ExternalID:         row.ExternalID,
				Role:               row.Role,
				VariantAssignments: make([]ImageAssignment, 0),
			}

			if row.EntityID != productID {
				image.VariantAssignments = append(image.VariantAssignments, ImageAssignment{
					ID:           row.ID.String(),
					EntityID:     row.EntityID.String(),
					EntityType:   row.EntityType,
					Role:         row.Role,
					DisplayOrder: row.DisplayOrder,
				})
			}
			images = append(images, image)
		}
	}
	return images
}

func mapToListProductResponse(productRow repository.GetProductsRow) ProductListModel {
	minPrice, _ := productRow.MinPrice.Float64Value()
	maxPrice, _ := productRow.MaxPrice.Float64Value()
	basePrice, _ := productRow.BasePrice.Float64Value()
	if minPrice.Float64 == 0 {
		minPrice = basePrice
	}
	if maxPrice.Float64 == 0 {
		maxPrice = basePrice
	}
	product := ProductListModel{
		ID:           productRow.ID.String(),
		Name:         productRow.Name,
		Description:  productRow.Description,
		MinPrice:     minPrice.Float64,
		MaxPrice:     maxPrice.Float64,
		Sku:          productRow.BaseSku,
		Slug:         productRow.Slug,
		ImgUrl:       productRow.ImgUrl,
		VariantCount: productRow.VariantCount,
		CreatedAt:    productRow.CreatedAt.String(),
		UpdatedAt:    productRow.UpdatedAt.String(),
	}
	if productRow.ImgID.Valid {
		id, _ := uuid.FromBytes(productRow.ImgID.Bytes[:])
		product.ImgID = utils.StringPtr(id.String())
	}

	return product
}
