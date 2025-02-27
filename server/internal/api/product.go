package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"golang.org/x/sync/errgroup"
)

type CreateProductRequest struct {
	Name         string           `json:"name" binding:"required,min=3,max=100"`
	Description  string           `json:"description" binding:"required,min=6,max=1000"`
	CategoryID   *int32           `json:"category_id,omitempty"`
	BrandID      *int32           `json:"brand_id,omitempty"`
	CollectionID *int32           `json:"collection_id,omitempty"`
	Variants     []VariantRequest `json:"variants" binding:"omitempty,dive"`
}

type UpdateProductRequest struct {
	Name        *string                            `json:"name" binding:"omitempty,min=3,max=100"`
	Description *string                            `json:"description" binding:"omitempty,min=10,max=1000"`
	CategoryID  *int32                             `json:"category_id,omitempty"`
	Variants    []repository.UpdateVariantTxParams `json:"variants" binding:"omitempty,dive"`
}

type ProductParam struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type ProductQueries struct {
	PaginationQueryParams
	Name *string `form:"name" binding:"omitempty,min=3,max=100"`
	Sku  *string `form:"sku" binding:"omitempty,alphanum"`
}

type ProductVariantModel struct {
	ID         string                   `json:"id"`
	Price      float64                  `json:"price"`
	StockQty   int32                    `json:"stock_qty"`
	Discount   int16                    `json:"discount"`
	Sku        *string                  `json:"sku,omitempty"`
	Attributes []ProductAttributeDetail `json:"attributes,omitempty"`
	ImageUrl   *string                  `json:"image_url,omitempty"`
}

type ProductDetailModel struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	UpdatedAt   string                `json:"updated_at"`
	CreatedAt   string                `json:"created_at"`
	ImageUrl    string                `json:"image_url,omitempty"`
	Variants    []ProductVariantModel `json:"variants,omitempty"`
}

type ProductListModel struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	VariantCount int64   `json:"variant_count"`
	ImageUrl     *string `json:"image_url,omitempty"`
	PriceFrom    float64 `json:"price_from,omitempty"`
	PriceTo      float64 `json:"price_to,omitempty"`
	DiscountTo   int16   `json:"discount_to"`
	CreatedAt    string  `json:"created_at,omitempty"`
}

// ------------------------------ Mapper ------------------------------
func getVariantFromRow(row repository.GetProductDetailRow) ProductVariantModel {
	if !row.VariantID.Valid {
		return ProductVariantModel{}
	}
	price, _ := row.Price.Float64Value()
	variantID := uuid.UUID(row.VariantID.Bytes).String()
	attributes := []ProductAttributeDetail{}
	variantModel := ProductVariantModel{
		ID:       variantID,
		Price:    price.Float64,
		StockQty: row.StockQuantity.Int32,
		Discount: row.Discount.Int16,
	}
	if row.VariantAttributeID.Valid {
		attributes = append(attributes, ProductAttributeDetail{
			ID:    row.VariantAttributeID.Int32,
			Name:  row.AttributeName.String,
			Value: []string{row.VariantAttributeValue.String},
		})
		variantModel.Attributes = attributes
	}

	if row.Sku.Valid {
		variantModel.Sku = &row.Sku.String
	}
	return variantModel
}

func mapToProductResponse(productRows []repository.GetProductDetailRow) ProductDetailModel {
	if len(productRows) == 0 {
		return ProductDetailModel{}
	}
	product := productRows[0].Product
	resp := ProductDetailModel{
		ID:          product.ProductID.String(),
		Name:        product.Name,
		Description: product.Description,
		UpdatedAt:   product.UpdatedAt.String(),
		CreatedAt:   product.CreatedAt.String(),
		Variants:    make([]ProductVariantModel, 0),
	}
	if productRows[0].ImgProductID.Valid {
		resp.ImageUrl = productRows[0].ImageUrl.String
	}
	firstVariant := getVariantFromRow(productRows[0])
	if firstVariant.ID != "" {
		resp.Variants = append(resp.Variants, getVariantFromRow(productRows[0]))
	}

	for i := 1; i < len(productRows); i++ {
		row := productRows[i]
		if !row.VariantID.Valid {
			continue
		}
		// add variants
		// get variant price
		variantID := uuid.UUID(row.VariantID.Bytes).String()
		prevVariantID := uuid.UUID(productRows[i-1].VariantID.Bytes).String()
		// add variant if it's the first variant or different from the previous one
		if variantID != prevVariantID {
			resp.Variants = append(
				resp.Variants,
				getVariantFromRow(row),
			)
		} else if row.AttributeID.Valid && productRows[i-1].AttributeID.Valid {
			// add attribute value to existing variant
			lastVariant := resp.Variants[len(resp.Variants)-1]
			attributeID := row.AttributeID.Int32
			lastAttributeID := productRows[i-1].AttributeID.Int32
			if attributeID != lastAttributeID {
				if row.VariantAttributeID.Valid {
					lastVariant.Attributes = append(lastVariant.Attributes, ProductAttributeDetail{
						ID:   row.VariantAttributeID.Int32,
						Name: row.AttributeName.String,
						Value: []string{
							row.VariantAttributeValue.String,
						},
					})
				}
			} else {
				lastVariant.Attributes[len(lastVariant.Attributes)-1].Value = append(lastVariant.Attributes[len(lastVariant.Attributes)-1].Value, row.VariantAttributeValue.String)
			}
		}
		if row.ImgVariantID.Valid && uuid.UUID(row.ImgVariantID.Bytes).String() == variantID && resp.Variants[len(resp.Variants)-1].ID == variantID {
			resp.Variants[len(resp.Variants)-1].ImageUrl = &row.ImageUrl.String
		}
	}
	return resp
}

func mapToListProductResponse(productRow repository.GetProductsRow) ProductListModel {
	minPrice, _ := productRow.MinPrice.Float64Value()
	maxPrice, _ := productRow.MaxPrice.Float64Value()
	product := ProductListModel{
		ID:           productRow.ProductID.String(),
		Name:         productRow.Name,
		Description:  productRow.Description,
		PriceFrom:    minPrice.Float64,
		PriceTo:      maxPrice.Float64,
		ImageUrl:     &productRow.ImageUrl.String,
		VariantCount: productRow.VariantCount,
		CreatedAt:    productRow.CreatedAt.String(),
	}
	if productRow.ImageUrl.Valid {
		product.ImageUrl = &productRow.ImageUrl.String
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
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	createProductTxParams := repository.CreateProductTxParam{
		Name:         req.Name,
		Description:  req.Description,
		CategoryID:   req.CategoryID,
		CollectionID: req.CollectionID,
		BrandID:      req.BrandID,
	}

	if len(req.Variants) > 0 {
		createProductTxParams.Variants = make([]repository.CreateVariantTxParam, len(req.Variants))
		for i, variantReq := range req.Variants {
			attributeParams := make([]struct {
				AttributeID int32
				Value       string
			}, len(variantReq.Attributes))
			for j, attribute := range variantReq.Attributes {
				attributeParams[j] = struct {
					AttributeID int32
					Value       string
				}{
					AttributeID: attribute.AttributeID,
					Value:       attribute.Value,
				}
			}

			createProductTxParams.Variants[i] = repository.CreateVariantTxParam{
				VariantPrice: StandardizeDecimal(variantReq.Price),
				VariantStock: variantReq.Stock,
				Attributes:   attributeParams,
				VariantSku:   variantReq.Sku,
				Discount:     variantReq.Discount,
			}
		}
	}
	newProduct, err := sv.repo.CreateProductTx(c, createProductTxParams)

	if err != nil {
		if errors.Is(err, repository.ErrUniqueViolation) {
			c.JSON(http.StatusBadRequest, mapErrResp(errors.New("product has already existed")))
			return
		}

		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusCreated, GenericResponse[repository.CreateProductTxResult]{&newProduct, nil, nil})
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
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	productRow, err := sv.repo.GetProductDetail(c, repository.GetProductDetailParams{
		ProductID: uuid.MustParse(params.ID),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if len(productRow) == 0 {
		c.JSON(http.StatusNotFound, mapErrResp(errors.New("product not found")))
		return
	}

	productDetail := mapToProductResponse(productRow)
	c.JSON(http.StatusOK, GenericResponse[ProductDetailModel]{&productDetail, nil, nil})
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
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	dbParams := repository.GetProductsParams{
		Limit:  20,
		Offset: 0,
	}
	if queries.PageSize != nil {
		dbParams.Limit = *queries.PageSize
		if queries.Page != nil {
			dbParams.Offset = (*queries.Page - 1) * *queries.PageSize
		}
	}

	if queries.Name != nil {
		dbParams.Name = utils.GetPgTypeText(*queries.Name)
	}
	if queries.Sku != nil {
		dbParams.Sku = utils.GetPgTypeText(*queries.Sku)
	}

	errGroup, ctx := errgroup.WithContext(c)

	productChan := make(chan []repository.GetProductsRow, 1)
	cntChan := make(chan int64, 1)
	defer close(productChan)
	defer close(cntChan)

	errGroup.Go(func() error {
		products, err := sv.repo.GetProducts(ctx, dbParams)
		if err != nil {
			return err
		}

		productChan <- products
		return nil
	})

	errGroup.Go(func() error {
		productCnt, err := sv.repo.CountProducts(ctx, repository.CountProductsParams{})
		if err != nil {
			return err
		}

		cntChan <- productCnt
		return nil
	})

	if err := errGroup.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	productResponses := make([]ProductListModel, 0)
	for _, product := range <-productChan {
		productResponses = append(productResponses, mapToListProductResponse(product))
	}

	c.JSON(http.StatusOK, GenericListResponse[ProductListModel]{productResponses, <-cntChan, nil, nil})
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
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	if req.Name == nil && req.Description == nil && len(req.Variants) == 0 {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("at least one field is required")))
		return
	}

	updated, err := sv.repo.UpdateProductTx(c, repository.UpdateProductTxParam{
		ProductID:   uuid.MustParse(params.ID),
		Name:        req.Name,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		Variants:    req.Variants,
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(err))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusOK, GenericResponse[repository.UpdateProductTxResult]{&updated, nil, nil})
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
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	_, err := sv.repo.GetProductByID(c, repository.GetProductByIDParams{
		ProductID: uuid.MustParse(params.ID),
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(err))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	err = sv.repo.DeleteProduct(c, uuid.MustParse(params.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	message := "product deleted"
	success := true
	c.JSON(http.StatusOK, GenericResponse[bool]{&success, &message, nil})
}
