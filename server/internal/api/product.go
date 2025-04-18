package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"golang.org/x/sync/errgroup"
)

type ProductParam struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type ProductQueries struct {
	PaginationQueryParams
	Name *string `form:"name" binding:"omitempty,min=3,max=100"`
	Sku  *string `form:"sku" binding:"omitempty,alphanum"`
}

type ProductAttributeDetail struct {
	ID     int32            `json:"id"`
	Name   string           `json:"name"`
	Values []AttributeValue `json:"values"`
}

type ProductVariantModel struct {
	ID         string                   `json:"id"`
	Price      float64                  `json:"price"`
	StockQty   int32                    `json:"stock_qty"`
	IsActive   bool                     `json:"is_active"`
	Sku        *string                  `json:"sku,omitempty"`
	Attributes []ProductAttributeDetail `json:"attributes"`
	ImageUrl   *string                  `json:"image_url,omitempty"`
}
type GeneralCategoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProductModel struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	BasePrice   float64                  `json:"price,omitzero"`
	BaseSku     string                   `json:"sku"`
	UpdatedAt   string                   `json:"updated_at"`
	Slug        string                   `json:"slug,omitempty"`
	IsActive    bool                     `json:"is_active"`
	Images      []ImageResponse          `json:"images,omitempty"`
	Variants    []ProductVariantModel    `json:"variants,omitempty"`
	Collection  *GeneralCategoryResponse `json:"collection,omitempty"`
	Brand       *GeneralCategoryResponse `json:"brand,omitempty"`
	Category    *GeneralCategoryResponse `json:"category,omitempty"`
	CreatedAt   string                   `json:"created_at"`
}

type ProductListModel struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	VariantCount int32     `json:"variant_count,omitzero"`
	Price        float64   `json:"price,omitzero"`
	Slug         string    `json:"slug,omitempty"`
	Sku          string    `json:"sku"`
	ImageUrl     string    `json:"image_url,omitempty"`
	ImageID      int32     `json:"image_id,omitempty"`
	CreatedAt    string    `json:"created_at,omitempty"`
	UpdatedAt    string    `json:"updated_at,omitempty"`
}

// ------------------------------ Mapper ------------------------------

func mapToProductResponse(productRows []repository.GetProductVariantsRow) ProductModel {
	if len(productRows) == 0 {
		return ProductModel{}
	}
	product := productRows[0]
	basePrice, _ := product.BasePrice.Float64Value()
	resp := ProductModel{
		ID:          product.ProductID.String(),
		Name:        product.Name,
		BasePrice:   basePrice.Float64,
		Description: product.Description.String,
		BaseSku:     product.BaseSku.String,
		Slug:        product.Slug,
		UpdatedAt:   product.UpdatedAt.String(),
		CreatedAt:   product.CreatedAt.String(),
		IsActive:    product.IsActive.Bool,
		Images:      make([]ImageResponse, 0),
		Variants:    make([]ProductVariantModel, 0),
	}
	if product.CollectionID.Valid {
		collectionID, _ := uuid.FromBytes(product.CollectionID.Bytes[:])
		resp.Collection = &GeneralCategoryResponse{
			ID:   collectionID.String(),
			Name: product.CollectionName.String,
		}
	}
	if product.BrandID.Valid {
		brandID, _ := uuid.FromBytes(product.BrandID.Bytes[:])
		resp.Brand = &GeneralCategoryResponse{
			ID:   brandID.String(),
			Name: product.BrandName.String,
		}
	}
	if product.CategoryID.Valid {
		categoryID, _ := uuid.FromBytes(product.CategoryID.Bytes[:])
		resp.Category = &GeneralCategoryResponse{
			ID:   categoryID.String(),
			Name: product.CategoryName.String,
		}
	}

	for rowIdx, row := range productRows {
		if !productRows[rowIdx].VariantID.Valid {
			continue
		}
		variantID, _ := uuid.FromBytes(row.VariantID.Bytes[:])
		// Check if the variant already exists in the response
		variantIdx := -1
		for vid, variant := range resp.Variants {
			if variant.ID == variantID.String() {
				variantIdx = vid
				break
			}
		}
		log.Debug().Msgf("variantIdx: %d", variantIdx)
		if variantIdx < 0 {
			/* If the variant does not exist, create a new variant and add it to the response */
			variantPrice, _ := row.VariantPrice.Float64Value()
			resp.Variants = append(resp.Variants, ProductVariantModel{
				ID:       variantID.String(),
				Price:    variantPrice.Float64,
				StockQty: row.VariantStock.Int32,
				IsActive: row.IsActive.Bool,
				Sku:      &row.VariantSku.String,
				ImageUrl: &row.VariantImageUrl.String,
				Attributes: []ProductAttributeDetail{
					{
						ID:   row.AttributeID.Int32,
						Name: row.AttributeName.String,
						Values: []AttributeValue{
							{
								ID:           row.AttributeValueID.Int32,
								Value:        row.AttributeValue.String,
								DisplayValue: row.AttributeDisplayValue.String,
								IsActive:     row.AttributeValueIsActive.Bool,
								DisplayOrder: row.AttributeDisplayOrder.Int16,
							},
						},
					},
				},
			})
		} else {
			/*
				If the variant already exists, check if the attribute already exists
				in the variant. If it does not exist, add it to the variant.
				Otherwise, do nothing.
			*/
			attrIdx := -1
			for i, attribute := range resp.Variants[variantIdx].Attributes {
				if attribute.ID == row.AttributeID.Int32 {
					attrIdx = i
					break
				}
			}
			log.Debug().Msgf("attributeIdx: %d", attrIdx)
			if attrIdx < 0 {
				resp.Variants[variantIdx].Attributes = append(resp.Variants[variantIdx].Attributes, ProductAttributeDetail{
					ID:   row.AttributeID.Int32,
					Name: row.AttributeName.String,
					Values: []AttributeValue{
						{
							ID:           row.AttributeValueID.Int32,
							Value:        row.AttributeValue.String,
							DisplayValue: row.AttributeDisplayValue.String,
							IsActive:     row.AttributeValueIsActive.Bool,
							DisplayOrder: row.AttributeDisplayOrder.Int16,
						},
					},
				})
			} else {
				attrValIdx := -1
				for i, value := range resp.Variants[variantIdx].Attributes[attrIdx].Values {
					if value.ID == row.AttributeValueID.Int32 {
						attrValIdx = i
						break
					}
				}
				log.Debug().Msgf("attrValIdx: %d", attrValIdx)
				if attrValIdx < 0 {
					// Add new attribute to the existing variant
					resp.Variants[variantIdx].Attributes[attrIdx].Values = append(resp.Variants[variantIdx].Attributes[attrIdx].Values, AttributeValue{
						ID:           row.AttributeValueID.Int32,
						Value:        row.AttributeValue.String,
						DisplayValue: row.AttributeDisplayValue.String,
						IsActive:     row.AttributeValueIsActive.Bool,
						DisplayOrder: row.AttributeDisplayOrder.Int16,
					})
				}
			}
		}
		/* If the image does not exist, add it to the list of images
		Otherwise, do nothing. */
		existedImage := false
		for _, image := range resp.Images {
			if image.ID == row.ImageID.Int32 {
				existedImage = true
				break
			}
		}
		if !existedImage {
			resp.Images = append(resp.Images, ImageResponse{
				ID:           row.ImageID.Int32,
				Url:          row.ImageUrl.String,
				ExternalID:   row.ImageExternalID.String,
				DisplayOrder: row.ImageDisplayOrder.Int16,
				Role:         row.ImageRole.String,
			})
		}
	}
	return resp
}

func mapToListProductResponse(productRow repository.GetProductsRow) ProductListModel {
	price, _ := productRow.BasePrice.Float64Value()
	product := ProductListModel{
		ID:           productRow.ID,
		Name:         productRow.Name,
		Description:  productRow.Description.String,
		Price:        price.Float64,
		Sku:          productRow.BaseSku.String,
		Slug:         productRow.Slug,
		ImageUrl:     productRow.ImageUrl.String,
		ImageID:      productRow.ImageID.Int32,
		VariantCount: int32(productRow.VariantCount),
		CreatedAt:    productRow.CreatedAt.String(),
		UpdatedAt:    productRow.UpdatedAt.String(),
	}

	return product
}

// ------------------------------ Handlers ------------------------------

// @Summary Create a new product
// @Schemes http
// @Description create a new product with the input payload
// @Tags products
// @Accept json
// @Param input body CreateProductRequest true "Product input"
// @Produce json
// @Success 200 {object} GenericResponse[repository.Product]
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /products [post]
func (sv *Server) createProduct(c *gin.Context) {
	var req repository.CreateProductTxParam
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	result, err := sv.repo.CreateProductTx(c, req)

	if err != nil {
		if errors.Is(err, repository.ErrUniqueViolation) {
			c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
			return
		}

		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))

		return
	}

	c.JSON(http.StatusCreated, createSuccessResponse(c, result, "product created", nil, nil))
}

// @Summary Get a product detail by ID
// @Schemes http
// @Description get a product detail by ID
// @Tags product detail
// @Accept json
// @Param product_id path int true "Product ID"
// @Produce json
// @Success 200 {object} GenericResponse[ProductListModel]
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /products/{product_id} [get]
func (sv *Server) getProductDetail(c *gin.Context) {
	var params ProductParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	productRow, err := sv.repo.GetProductVariants(c, repository.GetProductVariantsParams{
		ID: uuid.MustParse(params.ID),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}

	if len(productRow) == 0 {
		c.JSON(http.StatusNotFound, createErrorResponse(http.StatusNotFound, "", errors.New("product not found")))
		return
	}

	productDetail := mapToProductResponse(productRow)
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
// @Success 200 {array} GenericListResponse[ProductListModel]
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /products [get]
func (sv *Server) getProducts(c *gin.Context) {
	var queries ProductQueries
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	dbParams := repository.GetProductsParams{
		Limit:  20,
		Offset: 0,
	}

	dbParams.Limit = queries.PageSize
	dbParams.Offset = (queries.Page - 1) * queries.PageSize

	if queries.Name != nil {
		dbParams.Name = utils.GetPgTypeText(*queries.Name)
	}

	if queries.Sku != nil {
		dbParams.BaseSku = utils.GetPgTypeText(*queries.Sku)
	}

	products, err := sv.repo.GetProducts(c, dbParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "Server error", err))
		return
	}

	productCnt, err := sv.repo.CountProducts(c, repository.CountProductsParams{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "Server error", err))
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
		TotalPages:      int((productCnt + int64(queries.PageSize) - 1) / int64(queries.PageSize)),
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
// @Param input body UpdateProductRequest true "Product input"
// @Produce json
// @Success 200 {object} GenericResponse[repository.Product]
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /products/{product_id} [put]
func (sv *Server) updateProduct(c *gin.Context) {
	var params ProductParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	var req repository.UpdateProductTxParam
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	updated, err := sv.repo.UpdateProductTx(c, uuid.MustParse(params.ID), req)

	var removeImgErr *ApiError
	if len(req.RemovedImages) > 0 {
		errGroup, _ := errgroup.WithContext(c)
		for _, imgID := range req.RemovedImages {
			imgID := imgID
			errGroup.Go(func() error {
				return sv.removeImageUtil(c, imgID)
			})
		}
		err = errGroup.Wait()
		if err != nil {
			removeImgErr = &ApiError{
				Code:    strconv.Itoa(http.StatusBadRequest),
				Details: "Some images are not removed",
				Stack:   err.Error(),
			}
		}
	}
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse(http.StatusNotFound, "", err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, updated, "product updated", nil, removeImgErr))
}

// @Summary Remove a product by ID
// @Schemes http
// @Description remove a product by ID
// @Tags products
// @Accept json
// @Param product_id path int true "Product ID"
// @Produce json
// @Success 200 {object} GenericResponse[bool]
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /products/{product_id} [delete]
func (sv *Server) removeProduct(c *gin.Context) {
	var params ProductParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	_, err := sv.repo.GetProductByID(c, repository.GetProductByIDParams{
		ID: uuid.MustParse(params.ID),
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse(http.StatusNotFound, "", err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}

	err = sv.repo.DeleteProduct(c, uuid.MustParse(params.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}

	message := "product deleted"
	success := true
	c.JSON(http.StatusOK, createSuccessResponse(c, success, message, nil, nil))
}
