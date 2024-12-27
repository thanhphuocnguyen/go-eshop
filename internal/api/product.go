package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/postgres"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/sqlc"
	"github.com/thanhphuocnguyen/go-eshop/internal/util"
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
	Page     int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=100"`
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

func mapToProductResponse(productRow []sqlc.GetProductDetailRow) productResponse {
	if len(productRow) == 0 {
		return productResponse{}
	}
	product := productRow[0].Product
	price, _ := product.Price.Float64Value()
	resp := productResponse{
		ID:          product.ID,
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
				IsPrimary: img.ImageIsPrimary.Bool,
				ImageUrl:  img.ImageUrl.String,
			})

		}
	}

	return resp
}
func mapToListProductResponse(productRow sqlc.ListProductsRow) productListResponse {
	price, _ := productRow.Price.Float64Value()
	product := productListResponse{
		ID:          productRow.ID,
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
// @Success 200 {object} sqlc.Product
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

	newProduct, err := sv.postgres.CreateProduct(c, sqlc.CreateProductParams{
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

	c.JSON(http.StatusCreated, mapDefaultResp(newProduct, nil, nil))
}

// getProductDetail godoc
// @Summary Get a product detail by ID
// @Schemes http
// @Description get a product detail by ID
// @Tags product detail
// @Accept json
// @Param product_id path int true "Product ID"
// @Produce json
// @Success 200 {object} productResponse
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /products/{product_id} [get]
func (sv *Server) getProductDetail(c *gin.Context) {
	var params getProductParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	productRow, err := sv.postgres.GetProductDetail(c, sqlc.GetProductDetailParams{
		ID: params.ID,
	})

	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
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

	c.JSON(http.StatusOK, mapDefaultResp(mapToProductResponse(productRow), nil, nil))
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
// @Success 200 {array} productResponse
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /products [get]
func (sv *Server) getProducts(c *gin.Context) {
	var queries listProductsParams
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	productsChan := make(chan []sqlc.ListProductsRow)
	countChan := make(chan int64)
	go func() {
		products, err := sv.postgres.ListProducts(c, sqlc.ListProductsParams{
			Limit:  queries.PageSize,
			Offset: (queries.Page - 1) * queries.PageSize,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, mapErrResp(err))
			return
		}
		productsChan <- products
	}()
	go func() {
		count, err := sv.postgres.CountProducts(c, sqlc.CountProductsParams{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, mapErrResp(err))
			return
		}
		countChan <- count
	}()

	productResponses := make([]productListResponse, 0)
	for _, product := range <-productsChan {
		productResponses = append(productResponses, mapToListProductResponse(product))
	}

	c.JSON(http.StatusOK, mapListResp(productResponses, <-countChan))
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
// @Success 200 {object} productResponse
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
	updateBody := sqlc.UpdateProductParams{
		ID: params.ID,
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

	updated, err := sv.postgres.UpdateProduct(c, updateBody)
	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(err))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusOK, mapDefaultResp(updated, nil, nil))
}

// removeProduct godoc
// @Summary Remove a product by ID
// @Schemes http
// @Description remove a product by ID
// @Tags products
// @Accept json
// @Param product_id path int true "Product ID"
// @Produce json
// @Success 200 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /products/{product_id} [delete]
func (sv *Server) removeProduct(c *gin.Context) {
	var params getProductParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	_, err := sv.postgres.GetProduct(c, sqlc.GetProductParams{
		ID: params.ID,
	})
	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(err))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	err = sv.postgres.DeleteProduct(c, params.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	message := "product deleted"
	c.JSON(http.StatusOK, mapDefaultResp(nil, &message, nil))
}
