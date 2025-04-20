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
	ID    int32          `json:"id"`
	Name  string         `json:"name"`
	Value AttributeValue `json:"value"`
}

type ImageAssignment struct {
	ID           int32  `json:"id"`
	EntityID     string `json:"entity_id"`
	EntityType   string `json:"entity_type"`
	DisplayOrder int16  `json:"display_order"`
	Role         string `json:"role"`
}

type VariantImageModel struct {
	ID          int32             `json:"id"`
	Url         string            `json:"url"`
	ExternalID  string            `json:"external_id"`
	Assignments []ImageAssignment `json:"assignments"`
}

type ProductVariantModel struct {
	ID         string                   `json:"id"`
	Price      float64                  `json:"price"`
	StockQty   int32                    `json:"stock_qty"`
	IsActive   bool                     `json:"is_active"`
	Sku        *string                  `json:"sku,omitempty"`
	Attributes []ProductAttributeDetail `json:"attributes"`
}
type GeneralCategoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProductModel struct {
	ID            string                   `json:"id"`
	Name          string                   `json:"name"`
	Description   string                   `json:"description"`
	BasePrice     float64                  `json:"price,omitzero"`
	BaseSku       string                   `json:"sku"`
	UpdatedAt     string                   `json:"updated_at"`
	IsActive      bool                     `json:"is_active"`
	Slug          string                   `json:"slug"`
	CreatedAt     string                   `json:"created_at"`
	Variants      []ProductVariantModel    `json:"variants"`
	VariantImages []VariantImageModel      `json:"variant_images"`
	Images        []ImageResponse          `json:"images"`
	Collection    *GeneralCategoryResponse `json:"collection,omitempty"`
	Brand         *GeneralCategoryResponse `json:"brand,omitempty"`
	Category      *GeneralCategoryResponse `json:"category,omitempty"`
}

type ProductListModel struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	VariantCount int32     `json:"variant_count,omitzero"`
	Price        float64   `json:"price,omitzero"`
	Slug         string    `json:"slug,omitempty"`
	Sku          string    `json:"sku"`
	ImgUrl       string    `json:"image_url,omitempty"`
	ImgID        int32     `json:"image_id,omitempty"`
	CreatedAt    string    `json:"created_at,omitempty"`
	UpdatedAt    string    `json:"updated_at,omitempty"`
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

	productRows, err := sv.repo.GetProductDetail(c, repository.GetProductDetailParams{
		ID: uuid.MustParse(params.ID),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}

	if len(productRows) == 0 {
		c.JSON(http.StatusNotFound, createErrorResponse(http.StatusNotFound, "", errors.New("product not found")))
		return
	}

	productDetail := mapToProductResponse(productRows)

	variantRows, err := sv.repo.GetProductVariants(c, repository.GetProductVariantsParams{
		ProductID: uuid.MustParse(params.ID),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}
	variants, variantImages := mapToVariantResp(variantRows)
	productDetail.Variants = variants
	productDetail.VariantImages = variantImages
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
			img, err := sv.repo.GetImageFromID(c, repository.GetImageFromIDParams{
				ID:         imgID,
				EntityType: repository.ProductImageType,
			})
			if err != nil {
				if errors.Is(err, repository.ErrRecordNotFound) {
					c.JSON(http.StatusNotFound, createErrorResponse(http.StatusNotFound, "", err))
					return
				}
				c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
				return
			}

			errGroup.Go(func() error {
				_, err := sv.removeImageUtil(c, img.ExternalID)
				if err != nil {
					log.Error().Msgf("Error removing image %s: %v", img.ExternalID, err)
				}
				return err
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

// ------------------------------ Mapper ------------------------------

func mapToProductResponse(productRows []repository.GetProductDetailRow) ProductModel {
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

	for _, row := range productRows {
		/* If the image does not exist, add it to the list of images
		Otherwise, do nothing. */
		if !row.ImgID.Valid {
			continue
		}
		existedImage := false
		for _, image := range resp.Images {
			if image.ID == row.ImgID.Int32 {
				existedImage = true
				break
			}
		}
		if !existedImage {
			resp.Images = append(resp.Images, ImageResponse{
				ID:           row.ImgID.Int32,
				Url:          row.ImgUrl.String,
				ExternalID:   row.ImgExternalID.String,
				DisplayOrder: row.ImgAssignmentDisplayOrder.Int16,
				Role:         row.ImgAssignmentRole.String,
			})
		}
	}
	return resp
}

func mapToVariantResp(variantRows []repository.GetProductVariantsRow) ([]ProductVariantModel, []VariantImageModel) {
	variants := make([]ProductVariantModel, 0)
	variantImages := make([]VariantImageModel, 0)
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
				if a.ID == row.AttrID.Int32 {
					attrIdx = j
					break
				}
			}

			if attrIdx != -1 {
				// If the attribute already exists, do nothing
				continue
			}
			variants[variantIdx].Attributes = append(variants[variantIdx].Attributes, ProductAttributeDetail{
				ID:   row.AttrID.Int32,
				Name: row.AttrName.String,
				Value: AttributeValue{
					ID:           row.AttrValID.Int32,
					Value:        row.AttrValue.String,
					DisplayValue: row.AttrDisplayValue.String,
					IsActive:     row.IsActive.Bool,
					DisplayOrder: row.AttrDisplayOrder.Int16,
				},
			})
		} else {
			// If the variant does not exist, add it to the list of variants
			price, _ := row.Price.Float64Value()
			variant := ProductVariantModel{
				ID:       row.ID.String(),
				Price:    price.Float64,
				StockQty: row.Stock,
				IsActive: row.IsActive.Bool,
				Sku:      &row.Sku,
				Attributes: []ProductAttributeDetail{
					{
						ID:   row.AttrID.Int32,
						Name: row.AttrName.String,
						Value: AttributeValue{
							ID:           row.AttrValID.Int32,
							Value:        row.AttrValue.String,
							DisplayValue: row.AttrDisplayValue.String,
							IsActive:     row.IsActive.Bool,
							DisplayOrder: row.AttrDisplayOrder.Int16,
						},
					},
				},
			}
			variants = append(variants, variant)
		}
		if row.ImgID.Int32 != 0 {
			// If the image does not exist, add it to the list of images
			imageIdx := -1
			for i, image := range variantImages {
				if image.ID == row.ImgID.Int32 {
					imageIdx = i
					break
				}
			}
			if imageIdx != -1 {
				variantImages = append(variantImages, VariantImageModel{
					ID:         row.ImgID.Int32,
					Url:        row.ImgUrl.String,
					ExternalID: row.ImgExternalID.String,
					Assignments: []ImageAssignment{
						{
							ID:           row.ImgID.Int32,
							EntityID:     row.ID.String(),
							EntityType:   repository.ProductImageType,
							Role:         row.ImgAssignmentRole.String,
							DisplayOrder: row.ImgAssignmentDisplayOrder.Int16,
						},
					},
				})
			} else {
				// If the image already exists, append the assignment to the existing image
				for j, image := range variantImages {
					if image.ID == row.ImgID.Int32 {
						variantImages[j].Assignments = append(variantImages[j].Assignments, ImageAssignment{
							ID:           row.ImgID.Int32,
							EntityID:     row.ID.String(),
							EntityType:   repository.ProductImageType,
							Role:         row.ImgAssignmentRole.String,
							DisplayOrder: row.ImgAssignmentDisplayOrder.Int16,
						})
						break
					}
				}
			}
		}
	}
	return variants, variantImages
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
		ImgUrl:       productRow.ImgUrl,
		ImgID:        productRow.ImgID,
		VariantCount: int32(productRow.VariantCount),
		CreatedAt:    productRow.CreatedAt.String(),
		UpdatedAt:    productRow.UpdatedAt.String(),
	}

	return product
}
