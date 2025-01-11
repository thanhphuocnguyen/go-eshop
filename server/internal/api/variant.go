package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/util"
)

// ----------------------------------------------------------------------------- STRUCTS ----------------------------------------------------------------------------- //
type variantRequest struct {
	Name       string   `json:"name" binding:"required"`
	SKU        *string  `json:"sku,omitempty"`
	Price      float64  `json:"price"`
	Stock      int32    `json:"stock"`
	Attributes *[]int32 `json:"attributes,omitempty"`
}

type updateVariantRequest struct {
	Name       *string  `json:"name,omitempty"`
	SKU        *string  `json:"sku,omitempty"`
	Price      *float64 `json:"price,omitempty"`
	Stock      *int32   `json:"stock,omitempty"`
	Attributes *[]int32 `json:"attributes,omitempty"`
}

type getVariantParams struct {
	ID int64 `uri:"variant_id" binding:"required"`
}

type variantResponse struct {
	VariantID  int64             `json:"variant_id"`
	ProductID  int64             `json:"product_id"`
	Name       string            `json:"name"`
	SKU        *string           `json:"sku,omitempty"`
	Price      float64           `json:"price"`
	Stock      int32             `json:"stock"`
	Attributes map[string]string `json:"attributes,omitempty"`
	CreatedAt  string            `json:"created_at"`
	UpdatedAt  string            `json:"updated_at"`
}

// ----------------------------------------------------------------------------- API HANDLERS ----------------------------------------------------------------------------- //
func (sv *Server) createVariant(c *gin.Context) {
	var req variantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var productID getProductParams
	if err := c.ShouldBindUri(&productID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product, err := sv.repo.GetProduct(c, repository.GetProductParams{
		ProductID: productID.ID,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, mapErrResp(err))
		return
	}
	createVariantParam := repository.CreateVariantTxParam{
		ProductID:    product.ProductID,
		VariantName:  req.Name,
		VariantPrice: req.Price,
		VariantStock: req.Stock,
	}

	if req.Attributes != nil {
		createVariantParam.Attributes = *req.Attributes
	}

	variant, err := sv.repo.CreateVariantTx(c, createVariantParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusCreated, GenericResponse[repository.ProductVariant]{&variant.Variant, nil, nil})
}

func (sv *Server) updateVariant(c *gin.Context) {
	var req updateVariantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name == nil && req.SKU == nil && req.Price == nil && req.Stock == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "at least one field is required"})
		return
	}

	var params getVariantParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	variant, err := sv.repo.GetVariantByID(c, params.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(err))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	updateVariantParam := repository.UpdateVariantParams{
		VariantID: variant.VariantID,
	}

	if req.Name != nil {
		updateVariantParam.VariantName = util.GetPgTypeText(*req.Name)
	}
	if req.SKU != nil {
		updateVariantParam.VariantSku = util.GetPgTypeText(*req.SKU)
	}
	if req.Price != nil {
		updateVariantParam.VariantPrice = util.GetPgNumericFromFloat(*req.Price)
	}
	if req.Stock != nil {
		updateVariantParam.VariantStock = util.GetPgTypeInt4(*req.Stock)
	}

	updated, err := sv.repo.UpdateVariant(c, updateVariantParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	c.JSON(http.StatusOK, GenericResponse[repository.ProductVariant]{&updated, nil, nil})
}

func (sv *Server) deleteVariant(c *gin.Context) {
	var params getVariantParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	variant, err := sv.repo.GetVariantByID(c, params.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(err))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	err = sv.repo.DeleteVariant(c, variant.VariantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (sv *Server) getVariant(c *gin.Context) {
	var params getVariantParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	variantRows, err := sv.repo.GetVariantDetails(c, params.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	if len(variantRows) == 0 {
		c.JSON(http.StatusNotFound, mapErrResp(repository.ErrRecordNotFound))
		return
	}

	variant := variantRows[0]
	price, _ := variant.VariantPrice.Float64Value()
	resp := variantResponse{
		VariantID:  variant.VariantID,
		ProductID:  variant.ProductID,
		Name:       variant.VariantName,
		Price:      price.Float64,
		Stock:      variant.VariantStock,
		Attributes: make(map[string]string),
		CreatedAt:  variant.CreatedAt.String(),
		UpdatedAt:  variant.UpdatedAt.String(),
	}
	if variant.VariantSku.Valid {
		resp.SKU = &variant.VariantSku.String
	}
	for _, attr := range variantRows {
		resp.Attributes[attr.AttributeName] = attr.AttributeValue
	}

	c.JSON(http.StatusOK, GenericResponse[variantResponse]{&resp, nil, nil})
}

func (sv *Server) getVariants(c *gin.Context) {

}
