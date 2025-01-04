package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/util"
)

type createProductRequest struct {
	Name        string  `json:"name" binding:"required,min=3,max=100"`
	Description string  `json:"description" binding:"required,min=10,max=1000"`
	Sku         string  `json:"sku" binding:"required,alphanum"`
	Stock       int32   `json:"stock" binding:"required,gt=0"`
	Price       float64 `json:"price" binding:"required,gt=0,lt=10000"`
}

type updateProductRequest struct {
	Name        *string  `json:"name" binding:"omitempty,min=3,max=100"`
	Description *string  `json:"description" binding:"omitempty,min=10,max=1000"`
	Sku         *string  `json:"sku" binding:"omitempty,alphanum"`
	Stock       *int32   `json:"stock" binding:"omitempty,gt=0,lt=10000"`
	Price       *float64 `json:"price" binding:"omitempty,gt=0,lt=10000"`
}

type getProductParams struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type listProductsParams struct {
	Page     int32   `form:"page" binding:"required,min=1"`
	PageSize int32   `form:"page_size" binding:"required,min=5,max=20"`
	Name     *string `form:"name" binding:"omitempty,min=3,max=100"`
	Sku      *string `form:"sku" binding:"omitempty,alphanum"`
}

type productImage struct {
	ID        int32  `json:"id"`
	IsPrimary bool   `json:"is_primary"`
	ImageUrl  string `json:"image_url"`
}

type productResponse struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Sku         string         `json:"sku"`
	Images      []productImage `json:"images"`
	Stock       int32          `json:"stock"`
	Price       float64        `json:"price"`
	UpdatedAt   string         `json:"updated_at"`
	CreatedAt   string         `json:"created_at"`
}

type productListResponse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Sku         string  `json:"sku"`
	ImageUrl    *string `json:"image_url"`
	Stock       int32   `json:"stock"`
	Price       float64 `json:"price"`
	UpdatedAt   string  `json:"updated_at"`
	CreatedAt   string  `json:"created_at"`
}

// ------------------------------ Mapper ------------------------------

func mapToProductResponse(productRow []repository.GetProductDetailRow) productResponse {
	if len(productRow) == 0 {
		return productResponse{}
	}
	product := productRow[0].Product
	price, _ := product.Price.Float64Value()
	resp := productResponse{
		ID:          product.ProductID,
		Name:        product.Name,
		Description: product.Description,
		Sku:         product.Sku,
		Stock:       product.Stock,
		Price:       price.Float64,

		UpdatedAt: product.UpdatedAt.String(),
		CreatedAt: product.CreatedAt.String(),
	}
	for _, img := range productRow {
		if img.ImageID.Valid {
			resp.Images = append(resp.Images, productImage{
				ID:        img.ImageID.Int32,
				IsPrimary: img.ImagePrimary.Bool,
				ImageUrl:  img.ImageUrl.String,
			})

		}
	}

	return resp
}
func mapToListProductResponse(productRow repository.ListProductsRow) productListResponse {
	price, _ := productRow.Price.Float64Value()
	product := productListResponse{
		ID:          productRow.ProductID,
		Name:        productRow.Name,
		Description: productRow.Description,
		Sku:         productRow.Sku,
		Stock:       productRow.Stock,
		Price:       price.Float64,
		UpdatedAt:   productRow.UpdatedAt.String(),
		CreatedAt:   productRow.CreatedAt.String(),
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
	var product createProductRequest
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	price, err := util.ParsePgTypeNumber(product.Price)

	if err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	newProduct, err := sv.repo.CreateProduct(c, repository.CreateProductParams{
		Name:        product.Name,
		Description: product.Description,
		Sku:         product.Sku,
		Stock:       product.Stock,
		Price:       price,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusCreated, GenericResponse[repository.Product]{&newProduct, nil, nil})
}

// getProductDetail godoc
// @Summary Get a product detail by ID
// @Schemes http
// @Description get a product detail by ID
// @Tags product detail
// @Accept json
// @Param product_id path int true "Product ID"
// @Produce json
// @Success 200 {object} GenericResponse[productResponse]
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /products/{product_id} [get]
func (sv *Server) getProductDetail(c *gin.Context) {
	var params getProductParams
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
	c.JSON(http.StatusOK, GenericResponse[productResponse]{&productDetail, nil, nil})
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
	var queries listProductsParams
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	productsQueryParams := repository.ListProductsParams{
		Limit:  queries.PageSize,
		Offset: (queries.Page - 1) * queries.PageSize,
	}
	if queries.Name != nil {
		productsQueryParams.Name = util.GetPgTypeText(*queries.Name)
	}
	if queries.Sku != nil {
		productsQueryParams.Sku = util.GetPgTypeText(*queries.Sku)
	}

	products, err := sv.repo.ListProducts(c, productsQueryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	productCnt, err := sv.repo.CountProducts(c, repository.CountProductsParams{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	productResponses := make([]productListResponse, 0)
	for _, product := range products {
		productResponses = append(productResponses, mapToListProductResponse(product))
	}

	c.JSON(http.StatusOK, GenericListResponse[productListResponse]{&productResponses, &productCnt, nil, nil})
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
	var params getProductParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	var product updateProductRequest
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	updateBody := repository.UpdateProductParams{
		ProductID: params.ID,
	}
	if product.Price != nil {
		price, err := util.ParsePgTypeNumber(*product.Price)
		if err != nil {
			c.JSON(http.StatusBadRequest, mapErrResp(err))
			return
		}
		updateBody.Price = price
	}

	if product.Name != nil {
		updateBody.Name = pgtype.Text{
			String: *product.Name,
			Valid:  true,
		}
	}

	if product.Description != nil {
		updateBody.Description = pgtype.Text{
			String: *product.Description,
			Valid:  true,
		}
	}

	if product.Sku != nil {
		updateBody.Sku = pgtype.Text{
			String: *product.Sku,
			Valid:  true,
		}
	}

	updated, err := sv.repo.UpdateProduct(c, updateBody)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(err))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusOK, GenericResponse[repository.Product]{&updated, nil, nil})
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
	var params getProductParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	_, err := sv.repo.GetProduct(c, repository.GetProductParams{
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
