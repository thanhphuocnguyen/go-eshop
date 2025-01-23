package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"golang.org/x/sync/errgroup"
)

type CreateProductRequest struct {
	Name        string           `json:"name" binding:"required,min=3,max=100"`
	Description string           `json:"description" binding:"required,min=10,max=1000"`
	CategoryID  *int32           `json:"category_id,omitempty"`
	Variants    []VariantRequest `json:"variants" binding:"omitempty,dive"`
}

type UpdateProductRequest struct {
	Name        *string                            `json:"name" binding:"omitempty,min=3,max=100"`
	Description *string                            `json:"description" binding:"omitempty,min=10,max=1000"`
	CategoryID  *int32                             `json:"category_id,omitempty"`
	Variants    []repository.UpdateVariantTxParams `json:"variants" binding:"omitempty,dive"`
}

type ProductParam struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type ProductQueries struct {
	QueryParams
	Name *string `form:"name" binding:"omitempty,min=3,max=100"`
	Sku  *string `form:"sku" binding:"omitempty,alphanum"`
}

type ProductVariantModel struct {
	ID         int64                    `json:"id"`
	Price      float64                  `json:"price"`
	StockQty   int32                    `json:"stock_qty"`
	Sku        *string                  `json:"sku,omitempty"`
	Attributes []ProductAttributeDetail `json:"attributes,omitempty"`
	ImageUrl   *string                  `json:"image_url,omitempty"`
}

type ProductDetailModel struct {
	ID          int64                 `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	UpdatedAt   string                `json:"updated_at"`
	CreatedAt   string                `json:"created_at"`
	Images      []ImageModel          `json:"images,omitempty"`
	Variants    []ProductVariantModel `json:"variants,omitempty"`
}

type ProductListModel struct {
	ID           int64   `json:"id"`
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

func mapToProductResponse(productRows []repository.GetProductDetailRow) ProductDetailModel {
	if len(productRows) == 0 {
		return ProductDetailModel{}
	}
	product := productRows[0].Product
	resp := ProductDetailModel{
		ID:          product.ProductID,
		Name:        product.Name,
		Description: product.Description,
		UpdatedAt:   product.UpdatedAt.String(),
		CreatedAt:   product.CreatedAt.String(),
		Images:      make([]ImageModel, 0),
		Variants:    make([]ProductVariantModel, 0),
	}

	for i, row := range productRows {
		// add variants
		// get variant price
		price, _ := row.Price.Float64Value()
		// add variant if it's the first variant or different from the previous one
		if i == 0 || row.VariantID != productRows[i-1].VariantID {
			variantModel := ProductVariantModel{
				ID:       row.VariantID,
				Price:    price.Float64,
				StockQty: row.StockQuantity,
				Attributes: []ProductAttributeDetail{
					{
						ID:    row.VariantAttributeID,
						Name:  row.AttributeName,
						Value: row.VariantAttributeValue,
					},
				},
			}
			if row.Sku.Valid {
				variantModel.Sku = &row.Sku.String
			}

			if row.ImgVariantID.Valid {
				variantModel.ImageUrl = &row.ImageUrl.String
			}
			resp.Variants = append(
				resp.Variants,
				variantModel,
			)
		} else if row.VariantID == productRows[i-1].VariantID && row.AttributeID != productRows[i-1].AttributeID {
			// add attribute value to existing variant
			lastVariant := &resp.Variants[len(resp.Variants)-1]
			lastVariant.Attributes = append(lastVariant.Attributes, ProductAttributeDetail{
				ID:    row.VariantAttributeID,
				Name:  row.AttributeName,
				Value: row.VariantAttributeValue,
			})
		}

		// Add images for product or variants
		if row.ImgProductID.Valid {
			if row.ImageUrl.String != resp.Images[len(resp.Images)-1].ImageUrl {
				resp.Images = append(resp.Images, ImageModel{
					ID:       row.ImageID.Int32,
					ImageUrl: row.ImageUrl.String,
				})
			}
		}
	}
	return resp
}

func mapToListProductResponse(productRow repository.GetProductsRow) ProductListModel {
	minPrice, _ := productRow.MinPrice.Float64Value()
	maxPrice, _ := productRow.MaxPrice.Float64Value()
	product := ProductListModel{
		ID:          productRow.ProductID,
		Name:        productRow.Name,
		Description: productRow.Description,
		PriceFrom:   minPrice.Float64,
		PriceTo:     maxPrice.Float64,

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
// @Router /product [post]
func (sv *Server) createProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	createProductTxParams := repository.CreateProductTxParam{
		Name:        req.Name,
		Description: req.Description,
		CategoryID:  req.CategoryID,
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
		ProductID: params.ID,
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

	productsQueryParams := repository.GetProductsParams{
		Limit:  queries.PageSize,
		Offset: (queries.Page - 1) * queries.PageSize,
	}

	if queries.Name != nil {
		productsQueryParams.Name = utils.GetPgTypeText(*queries.Name)
	}
	if queries.Sku != nil {
		productsQueryParams.Sku = utils.GetPgTypeText(*queries.Sku)
	}

	errGroup, ctx := errgroup.WithContext(c)

	productChan := make(chan []repository.GetProductsRow, 1)
	cntChan := make(chan int64, 1)
	defer close(productChan)
	defer close(cntChan)

	errGroup.Go(func() error {
		products, err := sv.repo.GetProducts(ctx, productsQueryParams)
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

	c.JSON(http.StatusOK, GenericListResponse[ProductListModel]{&productResponses, <-cntChan, nil, nil})
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
		ProductID:   params.ID,
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
		ProductID: params.ID,
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(err))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	err = sv.repo.DeleteProduct(c, params.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	message := "product deleted"
	success := true
	c.JSON(http.StatusOK, GenericResponse[bool]{&success, &message, nil})
}
