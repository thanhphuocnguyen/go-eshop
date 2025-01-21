package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/util"
)

// ----------------------------------------------------------------------------- STRUCTS ----------------------------------------------------------------------------- //
type variantAttributeReq struct {
	AttributeID int32  `json:"attribute_id" binding:"required"`
	VariantID   int64  `json:"variant_id" binding:"omitempty"`
	Value       string `json:"value" binding:"required,gt=0,lt=255"`
}
type variantRequest struct {
	Name       string                `json:"name" binding:"required,gt=0,lt=255"`
	Price      float64               `json:"price" binding:"required,gt=0,lt=1000000"`
	Stock      int32                 `json:"stock" binding:"required,gt=0,lt=1000000"`
	Sku        *string               `json:"sku,omitempty" binding:"omitempty"`
	Discount   *int32                `json:"discount" binding:"omitempty,gte=0,lt=10000"`
	Attributes []variantAttributeReq `json:"attributes,omitempty" binding:"omitempty"`
}

type updateAttributeRequest struct {
	ID    int32  `json:"id" binding:"required"`
	Value string `json:"value" binding:"required,gt=0,lt=255"`
}

type updateVariantRequest struct {
	Name       *string                  `json:"name,omitempty" binding:"omitempty,gt=0,lt=255"`
	Sku        *string                  `json:"sku,omitempty" binding:"omitempty"`
	Price      *float64                 `json:"price,omitempty" binding:"omitempty,gt=0,lt=1000000"`
	Stock      *int32                   `json:"stock,omitempty" binding:"omitempty,gt=0,lt=1000000"`
	Attributes []updateAttributeRequest `json:"attributes,omitempty" binding:"omitempty"`
}

type variantParams struct {
	productParam
	ID int64 `uri:"variant_id" binding:"required"`
}

type variantResponse struct {
	VariantID  int64    `json:"variant_id"`
	ProductID  int64    `json:"product_id"`
	Price      float64  `json:"price"`
	Stock      int32    `json:"stock"`
	Attributes []string `json:"attributes,omitempty"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
	Sku        *string  `json:"sku,omitempty"`
}

// -------------------------------------------------------------------- API HANDLERS --------------------------------------------------------------------- //

// godoc
// @Summary Create a variant
// @Description Create a variant
// @Tags variant
// @Accept json
// @Produce json
// @Param product_id path int true "Product ID"
// @Param variantRequest body variantRequest true "Create Variant Request"
// @Success 201 {object} repository.ProductVariant
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /variant/{product_id} [post]
func (sv *Server) createVariant(c *gin.Context) {
	var req variantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	var param productParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	product, err := sv.repo.GetProductByID(c, repository.GetProductByIDParams{
		ProductID: param.ID,
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(err))
			return
		}

		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	createAttributeParams := make([]repository.CreateVariantAttributeParams, len(req.Attributes))
	for i, attr := range req.Attributes {
		createAttributeParams[i] = repository.CreateVariantAttributeParams{
			AttributeID: attr.AttributeID,
			VariantID:   attr.VariantID,
			Value:       attr.Value,
		}
	}

	createVariantParam := repository.CreateVariantTxParam{
		ProductID:    product.ProductID,
		VariantName:  req.Name,
		VariantPrice: req.Price,
		VariantStock: req.Stock,
		VariantSku:   req.Sku,
		Attributes:   createAttributeParams,
	}

	variant, err := sv.repo.CreateVariantTx(c, createVariantParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusCreated, GenericResponse[repository.ProductVariant]{&variant, nil, nil})
}

// godoc
// @Summary Update a variant
// @Description Update a variant
// @Tags variant
// @Accept json
// @Produce json
// @Param variant_id path int true "Variant ID"
// @Param updateVariantRequest body updateVariantRequest true "Update Variant Request"
// @Success 200 {object} repository.ProductVariant
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /variant/{variant_id} [put]
func (sv *Server) updateVariant(c *gin.Context) {
	var params variantParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	var req updateVariantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	if req.Name == nil && req.Sku == nil && req.Price == nil && req.Stock == nil && len(req.Attributes) == 0 {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("at least one field is required")))
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

	updateParams := repository.UpdateVariantParams{
		VariantID: variant.VariantID,
	}

	if req.Sku != nil {
		updateParams.Sku = util.GetPgTypeText(*req.Sku)
	}
	if req.Price != nil {
		updateParams.Price = util.GetPgNumericFromFloat(*req.Price)
	}
	if req.Stock != nil {
		updateParams.StockQuantity = util.GetPgTypeInt4(*req.Stock)
	}
	resp := variantResponse{
		Attributes: make([]string, 0),
	}
	if len(req.Attributes) > 0 {
		for _, attr := range req.Attributes {
			attributeUpdated, err := sv.repo.UpdateVariantAttribute(c, repository.UpdateVariantAttributeParams{
				VariantAttributeID: attr.ID,
				Value:              util.GetPgTypeText(attr.Value),
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, mapErrResp(err))
				return
			}
			resp.Attributes = append(resp.Attributes, attributeUpdated.Value)
		}
	}

	updated, err := sv.repo.UpdateVariant(c, updateParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	resp.VariantID = updated.VariantID
	resp.ProductID = updated.ProductID
	price, _ := updated.Price.Float64Value()
	resp.Price = price.Float64
	resp.Stock = updated.StockQuantity
	resp.CreatedAt = updated.CreatedAt.String()
	resp.UpdatedAt = updated.UpdatedAt.String()

	c.JSON(http.StatusOK, GenericResponse[variantResponse]{&resp, nil, nil})
}

// godoc
// @Summary Delete a variant
// @Description Delete a variant
// @Tags variant
// @Accept json
// @Produce json
// @Param variant_id path int true "Variant ID"
// @Success 204
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /variant/{variant_id} [delete]
func (sv *Server) deleteVariant(c *gin.Context) {
	var params variantParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
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

// godoc
// @Summary Get a variant
// @Description Get a variant
// @Tags variant
// @Accept json
// @Produce json
// @Param variant_id path int true "Variant ID"
// @Success 200 {object} variantResponse
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /variant/{variant_id} [get]
func (sv *Server) getVariant(c *gin.Context) {
	var params variantParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
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
	price, _ := variant.Price.Float64Value()
	resp := variantResponse{
		VariantID:  variant.VariantID,
		ProductID:  variant.ProductID,
		Price:      price.Float64,
		Stock:      variant.StockQuantity,
		Attributes: make([]string, 0),
		CreatedAt:  variant.CreatedAt.String(),
		UpdatedAt:  variant.UpdatedAt.String(),
	}
	if variant.Sku.Valid {
		resp.Sku = &variant.Sku.String
	}
	for _, attr := range variantRows {

		resp.Attributes = append(resp.Attributes, attr.AttributeName)
	}

	c.JSON(http.StatusOK, GenericResponse[variantResponse]{&resp, nil, nil})
}

func (sv *Server) getVariants(c *gin.Context) {
	var params productParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	variants, err := sv.repo.GetVariantByProductID(c, params.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if len(variants) == 0 {
		c.JSON(http.StatusNotFound, mapErrResp(repository.ErrRecordNotFound))
		return
	}
	variantResponses := make([]variantResponse, 0)
	for _, variant := range variants {
		if len(variantResponses) == 0 || variantResponses[len(variantResponses)-1].VariantID != variant.VariantID {
			price, _ := variant.Price.Float64Value()
			resp := variantResponse{
				VariantID:  variant.VariantID,
				ProductID:  variant.ProductID,
				Price:      price.Float64,
				Stock:      variant.StockQuantity,
				Attributes: make([]string, 0),
				CreatedAt:  variant.CreatedAt.String(),
				UpdatedAt:  variant.UpdatedAt.String(),
			}
			if variant.Sku.Valid {
				resp.Sku = &variant.Sku.String
			}
			variantResponses = append(variantResponses, resp)
		} else {
			latest := &variantResponses[len(variantResponses)-1]
			latest.Attributes = append(latest.Attributes, variant.AttributeName)
		}
	}

	c.JSON(http.StatusOK, GenericListResponse[variantResponse]{&variantResponses, int64(len(variantResponses)), nil, nil})
}
