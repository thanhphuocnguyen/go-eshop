package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
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

type CreateProductReq struct {
	Name             string                                    `json:"name" binding:"required,min=3,max=100"`
	Description      string                                    `json:"description" binding:"omitempty,min=6,max=5000"`
	ShortDescription *string                                   `json:"short_description" binding:"omitempty,max=2000"`
	Price            float64                                   `json:"price" binding:"required,gt=0"`
	Sku              string                                    `json:"sku" binding:"required"`
	Slug             string                                    `json:"slug" binding:"omitempty"`
	CategoryID       string                                    `json:"category_id,omitempty" binding:"omitempty,uuid"`
	BrandID          string                                    `json:"brand_id,omitempty" binding:"omitempty,uuid"`
	CollectionID     *string                                   `json:"collection_id,omitempty" binding:"omitempty,uuid"`
	Attributes       []string                                  `json:"attributes" binding:"min=1"`
	Variants         []repository.CreateProductVariantTxParams `json:"variants,omitempty"`
}
type UpdateProductImageAssignments struct {
	ID           int32  `json:"id"`
	EntityID     string `json:"entity_id"`
	DisplayOrder int16  `json:"display_order"`
}

type UpdateProductImages struct {
	ID          string   `json:"id"`
	Role        *string  `json:"role"`
	IsRemoved   *bool    `json:"omitempty,is_removed"`
	Assignments []string `json:"assignments,omitempty"`
}

type UpdateProductReq struct {
	Name             *string               `json:"name" binding:"omitempty,min=3,max=100"`
	Description      *string               `json:"description" binding:"omitempty,min=6,max=5000"`
	ShortDescription *string               `json:"short_description" binding:"omitempty,max=2000"`
	Price            *float64              `json:"price" binding:"omitempty,gt=0"`
	Sku              *string               `json:"sku" binding:"omitempty"`
	Slug             *string               `json:"slug" binding:"omitempty"`
	Stock            *int32                `json:"stock" binding:"omitempty,gt=0"`
	CategoryID       *string               `json:"category_id,omitempty" binding:"omitempty,uuid"`
	BrandID          *string               `json:"brand_id,omitempty" binding:"omitempty,uuid"`
	CollectionID     *string               `json:"collection_id,omitempty" binding:"omitempty,uuid"`
	Attributes       []string              `json:"attributes" binding:"omitempty"`
	Images           []UpdateProductImages `json:"images" binding:"omitempty,dive"`
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
	ImgUrl       string  `json:"image_url,omitempty"`
	ImgID        string  `json:"image_id,omitempty"`
	CreatedAt    string  `json:"created_at,omitempty"`
	UpdatedAt    string  `json:"updated_at,omitempty"`
}

// ------------------------------ Handlers ------------------------------

// @Summary Create a new product
// @Schemes http
// @Description create a new product with the input payload
// @Tags products
// @Accept json
// @Param input body CreateProductReq true "Product input"
// @Produce json
// @Success 200 {object} ApiResponse[ProductListModel]
// @Failure 400 {object} ApiResponse[ProductListModel]
// @Failure 500 {object} ApiResponse[ProductListModel]
// @Router /products [post]
func (sv *Server) createProduct(c *gin.Context) {
	var req CreateProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[ProductListModel](InvalidBodyCode, "", err))
		return
	}
	attributes := make([]uuid.UUID, len(req.Attributes))
	for i, attr := range req.Attributes {
		attributes[i] = uuid.MustParse(attr)
	}
	createParams := repository.CreateProductParams{
		Name:        req.Name,
		Attributes:  attributes,
		Description: req.Description,
	}

	createParams.BasePrice = utils.GetPgNumericFromFloat(req.Price)
	createParams.ShortDescription = req.ShortDescription
	createParams.Slug = req.Slug
	createParams.BaseSku = req.Sku

	createParams.CategoryID = utils.GetPgTypeUUIDFromString(req.CategoryID)

	createParams.BrandID = utils.GetPgTypeUUIDFromString(req.BrandID)
	if req.CollectionID != nil {
		createParams.CollectionID = utils.GetPgTypeUUIDFromString(*req.CollectionID)
	}

	product, err := sv.repo.CreateProduct(c, createParams)
	if err != nil {
		log.Error().Err(err).Timestamp()

		if errors.Is(err, repository.ErrUniqueViolation) {
			c.JSON(http.StatusBadRequest, createErrorResponse[ProductListModel](ConflictCode, "", err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[ProductListModel](InternalServerErrorCode, "", err))
		return
	}

	resp := ProductListModel{
		ID:          product.ID.String(),
		Name:        product.Name,
		Description: product.Description,
		Slug:        product.Slug,
		Sku:         product.BaseSku,
		CreatedAt:   product.CreatedAt.String(),
		UpdatedAt:   product.UpdatedAt.String(),
	}

	c.JSON(http.StatusCreated, createSuccessResponse(c, resp, "product created", nil, nil))
}

// @Summary Get a product detail by ID
// @Schemes http
// @Description get a product detail by ID
// @Tags product detail
// @Accept json
// @Param product_id path int true "Product ID"
// @Produce json
// @Success 200 {object} ApiResponse[ProductListModel]
// @Failure 404 {object} ApiResponse[ProductListModel]
// @Failure 500 {object} ApiResponse[ProductListModel]
// @Router /products/{product_id} [get]
func (sv *Server) getProductDetailHandler(c *gin.Context) {
	var params ProductParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[ProductModel](InvalidBodyCode, "", err))
		return
	}
	productID := uuid.MustParse(params.ID)
	productRows, err := sv.repo.GetProductDetail(c, repository.GetProductDetailParams{
		ID: productID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[ProductModel](InternalServerErrorCode, "", err))
		return
	}

	if len(productRows) == 0 {
		c.JSON(http.StatusNotFound, createErrorResponse[ProductModel](NotFoundCode, "", errors.New("product not found")))
		return
	}

	productDetail := mapToProductResponse(productRows)

	variantRows, err := sv.repo.GetProductVariants(c, repository.GetProductVariantsParams{
		ProductID: uuid.MustParse(params.ID),
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
	entityIds = append(entityIds, productID)

	images, err := sv.repo.GetProductImagesAssigned(c, entityIds)
	imageResp := mapToProductImages(productID, images)
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

	if queries.Search != nil {
		search := *queries.Search
		search = strings.ReplaceAll(search, " ", "%")
		search = strings.ReplaceAll(search, ",", "%")
		search = strings.ReplaceAll(search, ":", "%")
		search = "%" + search + "%"
		dbParams.Search = &search
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
// @Param input body UpdateProductReq true "Product input"
// @Produce json
// @Success 200 {object} ApiResponse[ProductListModel]
// @Failure 404 {object} ApiResponse[ProductListModel]
// @Failure 500 {object} ApiResponse[ProductListModel]
// @Router /products/{product_id} [put]
func (sv *Server) updateProduct(c *gin.Context) {
	var param ProductParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[ProductListModel](InvalidBodyCode, "", err))
		return
	}

	var req UpdateProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[ProductListModel](InvalidBodyCode, "", err))
		return
	}
	productID := uuid.MustParse(param.ID)
	updateProductParam := repository.UpdateProductParams{
		ID: productID,
	}

	if req.Name != nil {
		updateProductParam.Name = req.Name
	}
	if req.Description != nil {
		updateProductParam.Description = req.Description
	}
	if req.ShortDescription != nil {
		updateProductParam.ShortDescription = req.ShortDescription
	}
	if req.Slug != nil {
		updateProductParam.Slug = req.Slug
	}
	if req.Sku != nil {
		updateProductParam.BaseSku = req.Sku
	}
	if req.CategoryID != nil {
		updateProductParam.CategoryID = utils.GetPgTypeUUIDFromString(*req.CategoryID)
	}
	if req.CollectionID != nil {
		updateProductParam.CollectionID = utils.GetPgTypeUUIDFromString(*req.CollectionID)
	}
	if req.BrandID != nil {
		updateProductParam.BrandID = utils.GetPgTypeUUIDFromString(*req.BrandID)
	}
	if req.Price != nil {
		updateProductParam.BasePrice = utils.GetPgNumericFromFloat(*req.Price)
	}
	if len(req.Attributes) > 0 {
		attributes := make([]uuid.UUID, len(req.Attributes))
		for i, attr := range req.Attributes {
			attributes[i] = uuid.MustParse(attr)
		}
		updateProductParam.Attributes = attributes
	}

	product, err := sv.repo.UpdateProduct(c, updateProductParam)

	if err != nil {
		log.Error().Err(err).Msg("UpdateProduct")
		c.JSON(http.StatusInternalServerError, createErrorResponse[ProductListModel](InternalServerErrorCode, "", err))
		return
	}

	var removeImgErr *ApiError
	if len(req.Images) > 0 {
		errGroup, _ := errgroup.WithContext(c)
		imgAssignmentArgs := make([]repository.UpdateProdImagesTxArgs, 0)
		for _, image := range req.Images {
			if image.IsRemoved != nil && *image.IsRemoved {
				errGroup.Go(func() (err error) {
					img, err := sv.repo.GetImageFromID(c, repository.GetImageFromIDParams{
						ID:         uuid.MustParse(image.ID),
						EntityType: repository.ProductEntityType,
					})

					if err != nil {
						if errors.Is(err, repository.ErrRecordNotFound) {
							c.JSON(http.StatusNotFound, createErrorResponse[ProductListModel](NotFoundCode, "", err))
							return
						}
						c.JSON(http.StatusInternalServerError, createErrorResponse[ProductListModel](InternalServerErrorCode, "", err))
						return
					}

					msg, err := sv.removeImageUtil(c, img.ExternalID)
					if err != nil {
						return fmt.Errorf("failed to remove image: %w, reason: %s", err, msg)
					}

					// Remove image from product
					err = sv.repo.DeleteProductImage(c, uuid.MustParse(image.ID))
					if err != nil {
						if errors.Is(err, repository.ErrRecordNotFound) {
							c.JSON(http.StatusNotFound, createErrorResponse[ProductListModel](NotFoundCode, "", err))
							return
						}
						c.JSON(http.StatusInternalServerError, createErrorResponse[ProductListModel](InternalServerErrorCode, "", err))
						return
					}
					return err
				})
			} else {
				// parse the assignments to UUIDs
				assignmentIds := make([]uuid.UUID, 0)
				for _, assignment := range image.Assignments {
					if assignment == productID.String() {
						continue
					}
					assignmentIds = append(assignmentIds, uuid.MustParse(assignment))
				}
				// append the image assignment to the list
				args := repository.UpdateProdImagesTxArgs{
					ImageID:    image.ID,
					Role:       image.Role,
					EntityID:   productID,
					EntityType: repository.ProductEntityType,
					VariantIDs: assignmentIds,
				}

				imgAssignmentArgs = append(imgAssignmentArgs, args)
			}

		}

		err = errGroup.Wait()
		if err != nil {
			removeImgErr = &ApiError{
				Code:    strconv.Itoa(http.StatusBadRequest),
				Details: "Some images are not removed",
				Stack:   err.Error(),
			}
		}
		if len(imgAssignmentArgs) != 0 {
			// update the image assignments
			err := sv.repo.UpdateProductImagesTx(c, imgAssignmentArgs)
			if err != nil {
				if errors.Is(err, repository.ErrRecordNotFound) {
					c.JSON(http.StatusNotFound, createErrorResponse[ProductListModel](InternalServerErrorCode, "", err))
					return
				}
				c.JSON(http.StatusInternalServerError, createErrorResponse[ProductListModel](InternalServerErrorCode, "", err))
				return
			}
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse[ProductListModel](InternalServerErrorCode, "", err))
			return
		}
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, product, "product updated", nil, removeImgErr))
}

func (sv *Server) updateProductVariants(c *gin.Context) {
	var params ProductParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[repository.UpdateProductVariantsTxResult](InvalidBodyCode, "", err))
		return
	}

	var req repository.UpdateProdVariantsTxArgs
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[repository.UpdateProductVariantsTxResult](InvalidBodyCode, "", err))
		return
	}

	updated, err := sv.repo.UpdateProductVariantsTx(c, uuid.MustParse(params.ID), req)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[repository.UpdateProductVariantsTxResult](NotFoundCode, "", err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[repository.UpdateProductVariantsTxResult](InternalServerErrorCode, "", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, updated, "product variants updated", nil, nil))
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
func (sv *Server) removeProduct(c *gin.Context) {
	var params ProductParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[bool](InvalidBodyCode, "", err))
		return
	}

	product, err := sv.repo.GetProductByID(c, repository.GetProductByIDParams{
		ID: uuid.MustParse(params.ID),
	})
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
		UpdatedAt:        product.UpdatedAt.String(),
		CreatedAt:        product.CreatedAt.String(),
		IsActive:         *product.IsActive,
		Variants:         make([]ProductVariantModel, 0),
		ProductImages:    make([]ProductImageModel, 0),
		Brand: &GeneralCategoryResponse{
			ID:   product.BrandID.String(),
			Name: product.BrandName,
		},
		Category: &GeneralCategoryResponse{
			ID:   product.CategoryID.String(),
			Name: product.CategoryName,
		},
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
			if row.EntityID.String() != productID.String() {
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

			if row.EntityID.String() != productID.String() {
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
	product := ProductListModel{
		ID:           productRow.ID.String(),
		Name:         productRow.Name,
		Description:  productRow.Description,
		MinPrice:     minPrice.Float64,
		MaxPrice:     maxPrice.Float64,
		Sku:          productRow.BaseSku,
		Slug:         productRow.Slug,
		ImgUrl:       productRow.ImgUrl,
		ImgID:        productRow.ImgID.String(),
		VariantCount: productRow.VariantCount,
		CreatedAt:    productRow.CreatedAt.String(),
		UpdatedAt:    productRow.UpdatedAt.String(),
	}

	return product
}
