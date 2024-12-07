package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/postgres"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/sqlc"
)

type createProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Sku         string  `json:"sku" binding:"required"`
	Stock       int32   `json:"stock" binding:"required"`
	Price       float32 `json:"price" binding:"required"`
}

// CreateProduct godoc
// @Summary Create a new product
// @Schemes http
// @Description create a new product with the input payload
// @Tags products
// @Accept json
// @Param input body createProductRequest true "Product input"
// @Produce json
// @Success 200 {object} productResponse
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /products [post]
func (sv *Server) createProduct(c *gin.Context) {
	var product createProductRequest
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	price, err := parsePgNumeric(product.Price)

	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
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
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, newProduct)
}

type getProductParams struct {
	ProductID int64 `uri:"product_id" binding:"required"`
}

// GetProductDetail godoc
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
func (sv *Server) getProduct(c *gin.Context) {
	var params getProductParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := sv.postgres.GetProduct(c, params.ProductID)
	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, convertToProductResponse(product))
}

type listProductsParams struct {
	Page     int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=100"`
}

// GetProducts godoc
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
func (sv *Server) listProducts(c *gin.Context) {
	var queries listProductsParams
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	products, err := sv.postgres.ListProducts(c, sqlc.ListProductsParams{
		Limit:  queries.PageSize,
		Offset: (queries.Page - 1) * queries.PageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	productResponses := make([]productResponse, len(products))
	for i, product := range products {
		productResponses[i] = convertToProductResponse(product)
	}

	c.JSON(http.StatusOK, productResponses)
}

type productResponse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Sku         string  `json:"sku"`
	ImageURL    string  `json:"image_url"`
	Stock       int32   `json:"stock"`
	Price       float64 `json:"price"`
	UpdatedAt   string  `json:"updated_at"`
	CreatedAt   string  `json:"created_at"`
}

func convertToProductResponse(product sqlc.Product) productResponse {
	price, _ := product.Price.Float64Value()
	return productResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Sku:         product.Sku,
		ImageURL:    product.ImageUrl.String,
		Stock:       product.Stock,
		Price:       price.Float64,
		UpdatedAt:   product.UpdatedAt.String(),
		CreatedAt:   product.CreatedAt.String(),
	}
}

// Remove Product godoc
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := sv.postgres.DeleteProduct(c, params.ProductID)
	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

// Update Product godoc
// @Summary Update a product by ID
// @Schemes http
// @Description update a product by ID
// @Tags products
// @Accept json
// @Param product_id path int true "Product ID"
// @Param input body createProductRequest true "Product input"
// @Produce json
// @Success 200 {object} productResponse
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /products/{product_id} [put]
type updateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Sku         string  `json:"sku" binding:"required"`
	Stock       int32   `json:"stock" binding:"required"`
	Price       float32 `json:"price" binding:"required"`
}

func (sv *Server) updateProduct(c *gin.Context) {
	var params getProductParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var product updateProductRequest
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	price, err := parsePgNumeric(product.Price)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	updated, err := sv.postgres.UpdateProduct(c, sqlc.UpdateProductParams{
		ID:          params.ProductID,
		Name:        getPgTypeText(product.Name),
		Description: getPgTypeText(product.Description),
		Sku:         getPgTypeText(product.Sku),
		Stock:       pgtype.Int4{Int32: product.Stock, Valid: true},
		Price:       price,
	})
	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, convertToProductResponse(updated))
}

func getPgTypeText(value string) pgtype.Text {
	return pgtype.Text{
		String: value,
		Valid:  true,
	}
}

func parsePgNumeric(value float32) (pgtype.Numeric, error) {

	price := pgtype.Numeric{}
	if err := price.Scan(fmt.Sprintf("%.2f", value)); err != nil {
		return price, err
	}
	return price, nil
}
