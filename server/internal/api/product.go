package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/util"
	"golang.org/x/sync/errgroup"
)

type createProductRequest struct {
	Name        string           `json:"name" binding:"required,min=3,max=100"`
	Description string           `json:"description" binding:"required,min=10,max=1000"`
	Variants    []variantRequest `json:"variants" binding:"omitempty,dive"`
	CategoryID  *int32           `json:"category_id,omitempty"`
}

type updateProductRequest struct {
	Name        *string                            `json:"name" binding:"omitempty,min=3,max=100"`
	Description *string                            `json:"description" binding:"omitempty,min=10,max=1000"`
	CategoryID  *int32                             `json:"category_id,omitempty"`
	Variants    []repository.UpdateVariantTxParams `json:"variants" binding:"omitempty,dive"`
}

type productParam struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type productQueries struct {
	QueryParams
	Name *string `form:"name" binding:"omitempty,min=3,max=100"`
	Sku  *string `form:"sku" binding:"omitempty,alphanum"`
}

type productVariantModel struct {
	ID         int64       `json:"id"`
	Price      float64     `json:"price"`
	StockQty   int32       `json:"stock_qty"`
	Attributes []Attribute `json:"attributes,omitempty"`
	Sku        *string     `json:"sku,omitempty"`
	ImageUrl   *string     `json:"image_url,omitempty"`
}

type productDetailsModel struct {
	ID          int64                 `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	UpdatedAt   string                `json:"updated_at"`
	CreatedAt   string                `json:"created_at"`
	Images      []productImageModel   `json:"images,omitempty"`
	Variants    []productVariantModel `json:"variants,omitempty"`
}

type productListModel struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	VariantCount int64   `json:"variant_count"`
	ImageUrl     *string `json:"image_url,omitempty"`
	PriceFrom    float64 `json:"price_from,omitempty"`
	PriceTo      float64 `json:"price_to,omitempty"`
	CreatedAt    string  `json:"created_at,omitempty"`
}

// ------------------------------ Mapper ------------------------------

func mapToProductResponse(productRows []repository.GetProductDetailRow) productDetailsModel {
	if len(productRows) == 0 {
		return productDetailsModel{}
	}
	product := productRows[0].Product
	resp := productDetailsModel{
		ID:          product.ProductID,
		Name:        product.Name,
		Description: product.Description,
		UpdatedAt:   product.UpdatedAt.String(),
		CreatedAt:   product.CreatedAt.String(),
		Images:      make([]productImageModel, 0),
		Variants:    make([]productVariantModel, 0),
	}

	for i, row := range productRows {
		// add variants
		// get variant price
		price, _ := row.Price.Float64Value()
		// add variant if it's the first variant or different from the previous one
		if i == 0 || row.VariantID != productRows[i-1].VariantID {
			variantModel := productVariantModel{
				ID:    row.VariantID,
				Price: price.Float64,
				Attributes: []Attribute{{
					ID:              row.AttributeID,
					Name:            row.AttributeName,
					AttributeValues: []string{row.VariantAttributeValue},
				}},
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
			lastVariant.Attributes = append(lastVariant.Attributes, Attribute{
				ID:              row.AttributeID,
				Name:            row.AttributeName,
				AttributeValues: []string{row.VariantAttributeValue},
			})
		} else {
			// add attribute value to existing attribute
			lastVariant := &resp.Variants[len(resp.Variants)-1]
			lastAttribute := &lastVariant.Attributes[len(lastVariant.Attributes)-1]
			lastAttribute.AttributeValues = append(lastAttribute.AttributeValues, row.VariantAttributeValue)
		}

		// Add images for product or variants
		if row.ImgProductID.Valid {
			if row.ImageUrl.String != resp.Images[len(resp.Images)-1].ImageUrl {
				resp.Images = append(resp.Images, productImageModel{
					ID:        row.ImageID.Int32,
					IsPrimary: row.ImagePrimary.Bool,
					ImageUrl:  row.ImageUrl.String,
				})
			}
		}
	}
	return resp
}

func mapToListProductResponse(productRow repository.GetProductsRow) productListModel {
	minPrice, _ := productRow.MinPrice.Float64Value()
	maxPrice, _ := productRow.MaxPrice.Float64Value()
	product := productListModel{
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

// createProduct godoc
// @Summary Create a new product
// @Schemes http
// @Description create a new product with the input payload
// @Tags products
// @Accept json
// @Param input body createProductRequest true "Product input"
// @Produce json
// @Success 200 {object} GenericResponse[repository.Product]
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /products [post]
func (sv *Server) createProduct(c *gin.Context) {
	var req createProductRequest
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
		for i, variant := range req.Variants {
			attributeParams := make([]repository.CreateVariantAttributeParams, len(variant.Attributes))
			for j, attribute := range variant.Attributes {
				attributeParams[j] = repository.CreateVariantAttributeParams{
					AttributeID: attribute.AttributeID,
					Value:       attribute.Value,
				}
			}

			createProductTxParams.Variants[i] = repository.CreateVariantTxParam{
				VariantName:  variant.Name,
				VariantPrice: variant.Price,
				VariantStock: variant.Stock,
				Attributes:   attributeParams,
				VariantSku:   variant.Sku,
			}
		}
	}
	newProduct, err := sv.repo.CreateProductTx(c, createProductTxParams)

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusCreated, GenericResponse[repository.CreateProductTxResult]{&newProduct, nil, nil})
}

// getProductDetail godoc
// @Summary Get a product detail by ID
// @Schemes http
// @Description get a product detail by ID
// @Tags product detail
// @Accept json
// @Param product_id path int true "Product ID"
// @Produce json
// @Success 200 {object} GenericResponse[productListModel]
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /products/{product_id} [get]
func (sv *Server) getProductDetail(c *gin.Context) {
	var params productParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	productRow, err := sv.repo.GetProductDetail(c, repository.GetProductDetailParams{
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

	if len(productRow) == 0 {
		c.JSON(http.StatusNotFound, mapErrResp(errors.New("product not found")))
		return
	}

	productDetail := mapToProductResponse(productRow)
	c.JSON(http.StatusOK, GenericResponse[productDetailsModel]{&productDetail, nil, nil})
}

// getProducts godoc
// @Summary Get list of products
// @Schemes http
// @Description get list of products
// @Tags products
// @Accept json
// @Param page query int true "Page number"
// @Param page_size query int true "Page size"
// @Produce json
// @Success 200 {array} GenericListResponse[productListModel]
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /products [get]
func (sv *Server) getProducts(c *gin.Context) {
	var queries productQueries
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	productsQueryParams := repository.GetProductsParams{
		Limit:  queries.PageSize,
		Offset: (queries.Page - 1) * queries.PageSize,
	}

	if queries.Name != nil {
		productsQueryParams.Name = util.GetPgTypeText(*queries.Name)
	}
	if queries.Sku != nil {
		productsQueryParams.Sku = util.GetPgTypeText(*queries.Sku)
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

	productResponses := make([]productListModel, 0)
	for _, product := range <-productChan {
		productResponses = append(productResponses, mapToListProductResponse(product))
	}

	c.JSON(http.StatusOK, GenericListResponse[productListModel]{&productResponses, <-cntChan, nil, nil})
}

// updateProduct godoc
// @Summary Update a product by ID
// @Schemes http
// @Description update a product by ID
// @Tags products
// @Accept json
// @Param product_id path int true "Product ID"
// @Param input body updateProductRequest true "Product input"
// @Produce json
// @Success 200 {object} GenericResponse[repository.Product]
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /products/{product_id} [put]
func (sv *Server) updateProduct(c *gin.Context) {
	var params productParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	var req updateProductRequest
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

// removeProduct godoc
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
	var params productParam
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
