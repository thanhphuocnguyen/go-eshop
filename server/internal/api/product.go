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
	Stock       int32            `json:"stock" binding:"required,gt=0"`
	Price       float64          `json:"price" binding:"required,gt=0,lt=10000"`
	Sku         *string          `json:"sku" binding:"omitempty,max=100"`
	Discount    *int32           `json:"discount" binding:"omitempty,gte=0,lt=10000"`
	CategoryID  *int32           `json:"category_id,omitempty"`
	Variants    []variantRequest `json:"variants" binding:"omitempty,dive"`
}

type updateProductRequest struct {
	Name        *string                            `json:"name" binding:"omitempty,min=3,max=100"`
	Description *string                            `json:"description" binding:"omitempty,min=10,max=1000"`
	Sku         *string                            `json:"sku" binding:"omitempty,max=100"`
	Stock       *int32                             `json:"stock" binding:"omitempty,gt=0,lt=10000"`
	Price       *float64                           `json:"price" binding:"omitempty,gt=0"`
	Discount    *int32                             `json:"discount" binding:"omitempty,gte=0,lte=100"`
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

type productImageModel struct {
	ID        int32  `json:"id"`
	IsPrimary bool   `json:"is_primary"`
	ImageUrl  string `json:"image_url"`
}

type productVariantModel struct {
	ID         int64             `json:"id"`
	Name       string            `json:"name"`
	Price      float64           `json:"price"`
	Attributes map[string]string `json:"attributes"`
}

type productDetailsModel struct {
	ID          int64                 `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Sku         string                `json:"sku"`
	Stock       int32                 `json:"stock"`
	Price       float64               `json:"price"`
	UpdatedAt   string                `json:"updated_at"`
	CreatedAt   string                `json:"created_at"`
	Images      []productImageModel   `json:"images"`
	Variants    []productVariantModel `json:"variants"`
}

type productListModel struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	VariantCount int64   `json:"variant_count"`
	Sku          string  `json:"sku"`
	ImageUrl     *string `json:"image_url,omitempty"`
	Stock        int32   `json:"stock,omitempty"`
	Price        float64 `json:"price,omitempty"`
	CreatedAt    string  `json:"created_at,omitempty"`
}

// ------------------------------ Mapper ------------------------------

func mapToProductResponse(productRows []repository.GetProductDetailRow) productDetailsModel {
	if len(productRows) == 0 {
		return productDetailsModel{}
	}
	product := productRows[0].Product
	price, _ := product.Price.Float64Value()
	resp := productDetailsModel{
		ID:          product.ProductID,
		Name:        product.Name,
		Description: product.Description,
		Sku:         product.Sku.String,
		Stock:       product.Stock,
		Price:       price.Float64,
		UpdatedAt:   product.UpdatedAt.String(),
		CreatedAt:   product.CreatedAt.String(),
		Images:      make([]productImageModel, 0),
		Variants:    make([]productVariantModel, 0),
	}
	for i, row := range productRows {
		if row.ImageUrl.Valid {
			if row.ImageUrl.String != resp.Images[len(resp.Images)-1].ImageUrl {
				resp.Images = append(resp.Images, productImageModel{
					ID:        row.ImageID.Int32,
					IsPrimary: row.ImagePrimary.Bool,
					ImageUrl:  row.ImageUrl.String,
				})
			}
		}
		if row.VariantID.Valid {
			price, _ := row.VariantPrice.Float64Value()
			if row.VariantID.Valid && (i == 0 || row.VariantID.Int64 != productRows[i-1].VariantID.Int64) {
				resp.Variants = append(
					resp.Variants,
					productVariantModel{
						ID:         row.VariantID.Int64,
						Name:       row.VariantName.String,
						Price:      price.Float64,
						Attributes: make(map[string]string),
					})
				if row.AttributeName.Valid {
					resp.Variants[len(resp.Variants)-1].Attributes[row.AttributeName.String] = row.AttributeValue.String
				}
			} else {
				if row.AttributeName.Valid {
					resp.Variants[len(resp.Variants)-1].Attributes[row.AttributeName.String] = row.AttributeValue.String
				}
			}
		}
	}

	return resp
}

func mapToListProductResponse(productRow repository.GetProductsRow) productListModel {
	price, _ := productRow.Price.Float64Value()
	product := productListModel{
		ID:           productRow.ProductID,
		Name:         productRow.Name,
		Description:  productRow.Description,
		Sku:          productRow.Sku.String,
		Stock:        productRow.Stock,
		Price:        price.Float64,
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
		Price:       req.Price,
		Discount:    req.Discount,
		Stock:       req.Stock,
		BrandID:     nil,
		CategoryID:  req.CategoryID,
	}
	if len(req.Variants) > 0 {
		createProductTxParams.Variants = make([]repository.CreateVariantTxParam, len(req.Variants))
		for i, variant := range req.Variants {
			createProductTxParams.Variants[i] = repository.CreateVariantTxParam{
				VariantName:  variant.Name,
				VariantPrice: variant.Price,
				VariantStock: variant.Stock,
				Attributes:   variant.Attributes,
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
// @Success 200 {object} GenericResponse[productListResponse]
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
// @Success 200 {array} GenericListResponse[productListResponse]
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

	c.JSON(http.StatusOK, GenericListResponse[productListModel]{productResponses, <-cntChan, nil, nil})
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

	if req.Name == nil && req.Description == nil && req.Sku == nil && req.Price == nil {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("at least one field is required")))
		return
	}

	updated, err := sv.repo.UpdateProductTx(c, repository.UpdateProductTxParam{
		ProductID:   params.ID,
		Name:        req.Name,
		Description: req.Description,
		Sku:         req.Sku,
		Price:       req.Price,
		Stock:       req.Stock,
		CategoryID:  req.CategoryID,
		Discount:    req.Discount,
		Variants:    req.Variants,
		BrandID:     nil,
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
