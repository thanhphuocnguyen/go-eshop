package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/util"
)

// ----------------------------------------------------------------------------- STRUCTS ----------------------------------------------------------------------------- //
type variantRequest struct {
	Name       string  `json:"name" binding:"required,gt=0,lt=255"`
	Price      float64 `json:"price" binding:"required,gt=0,lt=1000000"`
	Stock      int32   `json:"stock" binding:"required,gt=0,lt=1000000"`
	Sku        *string `json:"sku,omitempty" binding:"omitempty"`
	Attributes []int32 `json:"attributes,omitempty" binding:"omitempty"`
}

type updateVariantRequest struct {
	Name       *string  `json:"name,omitempty" binding:"omitempty,gt=0,lt=255"`
	Sku        *string  `json:"sku,omitempty" binding:"omitempty"`
	Price      *float64 `json:"price,omitempty" binding:"omitempty,gt=0,lt=1000000"`
	Stock      *int32   `json:"stock,omitempty" binding:"omitempty,gt=0,lt=1000000"`
	Attributes []int32  `json:"attributes,omitempty" binding:"omitempty"`
}

type variantParams struct {
	productParam
	ID int64 `uri:"variant_id" binding:"required"`
}

type variantResponse struct {
	VariantID  int64            `json:"variant_id"`
	ProductID  int64            `json:"product_id"`
	Name       string           `json:"name"`
	SKU        *string          `json:"sku,omitempty"`
	Price      float64          `json:"price"`
	Stock      int32            `json:"stock"`
	Attributes []AttributeValue `json:"attributes,omitempty"`
	CreatedAt  string           `json:"created_at"`
	UpdatedAt  string           `json:"updated_at"`
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

	createVariantParam := repository.CreateVariantTxParam{
		ProductID:    product.ProductID,
		VariantName:  req.Name,
		VariantPrice: req.Price,
		VariantStock: req.Stock,
		VariantSku:   req.Sku,
		Attributes:   req.Attributes,
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

	if req.Name != nil {
		updateParams.VariantName = util.GetPgTypeText(*req.Name)
	}
	if req.Sku != nil {
		updateParams.VariantSku = util.GetPgTypeText(*req.Sku)
	}
	if req.Price != nil {
		updateParams.VariantPrice = util.GetPgNumericFromFloat(*req.Price)
	}
	if req.Stock != nil {
		updateParams.VariantStock = util.GetPgTypeInt4(*req.Stock)
	}
	if len(req.Attributes) > 0 {
		attributes, err := sv.repo.GetVariantAttributes(c, variant.VariantID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, mapErrResp(err))
			return
		}
		existingAttributes := make(map[int32]bool)
		newAttributes := make([]int32, 0)
		for _, attr := range attributes {
			existingAttributes[attr.AttributeValueID] = false
		}
		for _, attr := range req.Attributes {
			if _, ok := existingAttributes[attr]; ok {
				existingAttributes[attr] = true
			} else {
				newAttributes = append(newAttributes, attr)
			}
		}

		log.Info().Interface("existingAttributes", existingAttributes).Msg("existingAttributes")
		for attr, exists := range existingAttributes {
			if !exists {
				log.Info().Int32("attr", attr).Int32("variantID", int32(variant.VariantID)).Msg("deleting")
				err = sv.repo.DeleteVariantAttribute(c, repository.DeleteVariantAttributeParams{
					AttributeValueID: attr,
					VariantID:        variant.VariantID,
				})
				if err != nil {
					c.JSON(http.StatusInternalServerError, mapErrResp(err))
					return
				}
			}
		}

		for _, attr := range newAttributes {
			_, err := sv.repo.CreateVariantAttribute(c, repository.CreateVariantAttributeParams{
				VariantID:        variant.VariantID,
				AttributeValueID: attr,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, mapErrResp(err))
				return
			}
		}
	}

	updated, err := sv.repo.UpdateVariant(c, updateParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusOK, GenericResponse[repository.ProductVariant]{&updated, nil, nil})
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
	price, _ := variant.VariantPrice.Float64Value()
	resp := variantResponse{
		VariantID:  variant.VariantID,
		ProductID:  variant.ProductID,
		Name:       variant.VariantName,
		Price:      price.Float64,
		Stock:      variant.VariantStock,
		Attributes: make([]AttributeValue, 0),
		CreatedAt:  variant.CreatedAt.String(),
		UpdatedAt:  variant.UpdatedAt.String(),
	}
	if variant.VariantSku.Valid {
		resp.SKU = &variant.VariantSku.String
	}
	for _, attr := range variantRows {
		value := AttributeValue{
			ID:    attr.AttributeValueID,
			Value: attr.AttributeValue,
		}
		if attr.Color.Valid {
			color := attr.Color.String
			value.Color = &color
		}
		resp.Attributes = append(resp.Attributes, value)
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
			price, _ := variant.VariantPrice.Float64Value()
			resp := variantResponse{
				VariantID:  variant.VariantID,
				ProductID:  variant.ProductID,
				Name:       variant.VariantName,
				Price:      price.Float64,
				Stock:      variant.VariantStock,
				Attributes: make([]AttributeValue, 0),
				CreatedAt:  variant.CreatedAt.String(),
				UpdatedAt:  variant.UpdatedAt.String(),
			}
			if variant.VariantSku.Valid {
				resp.SKU = &variant.VariantSku.String
			}
			variantResponses = append(variantResponses, resp)
		} else {
			latest := &variantResponses[len(variantResponses)-1]
			if variant.Color.Valid {
				color := variant.Color.String
				latest.Attributes = append(latest.Attributes, AttributeValue{
					ID:    variant.AttributeValueID,
					Value: variant.AttributeValue,
					Color: &color,
				})
			} else {
				latest.Attributes = append(latest.Attributes, AttributeValue{
					ID:    variant.AttributeValueID,
					Value: variant.AttributeValue,
				})
			}
		}
	}

	c.JSON(http.StatusOK, GenericListResponse[variantResponse]{variantResponses, int64(len(variantResponses)), nil, nil})
}
