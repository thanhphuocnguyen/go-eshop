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
	Price      float64 `json:"price" binding:"required,gt=0,lt=1000000"`
	Stock      int32   `json:"stock" binding:"required,gt=0,lt=1000000"`
	Sku        *string `json:"sku,omitempty" binding:"omitempty"`
	Discount   *int16  `json:"discount" binding:"omitempty,gte=0,lt=10000"`
	Attributes []struct {
		AttributeID int32  `json:"attribute_id" binding:"required"`
		Value       string `json:"value" binding:"required,gt=0,lt=255"`
	} `json:"attributes,omitempty" binding:"omitempty"`
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
	VariantID  int64                    `json:"variant_id"`
	Price      float64                  `json:"price"`
	Stock      int32                    `json:"stock"`
	Discount   int16                    `json:"discount,omitempty"`
	Sku        *string                  `json:"sku,omitempty"`
	Attributes []productAttributeDetail `json:"attributes,omitempty"`
	CreatedAt  string                   `json:"created_at"`
	UpdatedAt  string                   `json:"updated_at"`
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

	createAttributeParams := make([]struct {
		AttributeID int32
		Value       string
	}, len(req.Attributes))
	for i, attr := range req.Attributes {
		createAttributeParams[i] = struct {
			AttributeID int32
			Value       string
		}{
			AttributeID: attr.AttributeID,
			Value:       attr.Value,
		}
	}

	createVariantParam := repository.CreateVariantTxParam{
		ProductID:    param.ID,
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

	if len(req.Attributes) > 0 {
		for _, attr := range req.Attributes {
			_, err := sv.repo.UpdateVariantAttribute(c, repository.UpdateVariantAttributeParams{
				VariantAttributeID: attr.ID,
				Value:              util.GetPgTypeText(attr.Value),
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, mapErrResp(err))
				return
			}

		}
	}

	_, err = sv.repo.UpdateVariant(c, updateParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	msg := "Variant updated successfully"

	c.JSON(http.StatusOK, GenericResponse[bool]{nil, &msg, nil})
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
		Price:      price.Float64,
		Stock:      variant.StockQuantity,
		Discount:   variant.Discount,
		Attributes: make([]productAttributeDetail, 0),
		CreatedAt:  variant.CreatedAt.String(),
		UpdatedAt:  variant.UpdatedAt.String(),
	}

	if variant.Sku.Valid {
		resp.Sku = &variant.Sku.String
	}
	for _, attr := range variantRows {
		resp.Attributes = append(resp.Attributes, productAttributeDetail{
			Name:  attr.AttributeName,
			Value: attr.Value,
		})
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
				Price:      price.Float64,
				Stock:      variant.StockQuantity,
				Attributes: make([]productAttributeDetail, 0),
				CreatedAt:  variant.CreatedAt.String(),
				UpdatedAt:  variant.UpdatedAt.String(),
			}
			if variant.Sku.Valid {
				resp.Sku = &variant.Sku.String
			}
			variantResponses = append(variantResponses, resp)
		} else {
			latest := &variantResponses[len(variantResponses)-1]
			latest.Attributes = append(latest.Attributes, productAttributeDetail{
				Name:  variant.AttributeName,
				Value: variant.AttributeValue,
			})
		}
	}

	c.JSON(http.StatusOK, GenericListResponse[variantResponse]{&variantResponses, int64(len(variantResponses)), nil, nil})
}
