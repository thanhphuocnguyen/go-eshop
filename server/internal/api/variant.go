package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

// ----------------------------------------------------------------------------- STRUCTS ----------------------------------------------------------------------------- //
type VariantAttributeRequest struct {
	AttributeID int32  `json:"attribute_id" binding:"required"`
	Value       string `json:"value" binding:"required,gt=0,lt=255"`
}
type VariantRequest struct {
	Price      float64                   `json:"price" binding:"required,gt=0,lt=1000000"`
	Stock      int32                     `json:"stock" binding:"required,gt=0,lt=1000000"`
	Discount   int16                     `json:"discount" binding:"omitempty,gte=0,lt=10000"`
	Sku        *string                   `json:"sku,omitempty" binding:"omitempty"`
	Attributes []VariantAttributeRequest `json:"attributes" binding:"dive,omitempty"`
}

type VariantUpdateAttributeRequest struct {
	ID    int32  `json:"id" binding:"required"`
	Value string `json:"value" binding:"required,gt=0,lt=255"`
}

type UpdateVariantRequest struct {
	Name       *string                         `json:"name,omitempty" binding:"omitempty,gt=0,lt=255"`
	Sku        *string                         `json:"sku,omitempty" binding:"omitempty"`
	Price      *float64                        `json:"price,omitempty" binding:"omitempty,gt=0,lt=1000000"`
	Stock      *int32                          `json:"stock,omitempty" binding:"omitempty,gt=0,lt=1000000"`
	Attributes []VariantUpdateAttributeRequest `json:"attributes,omitempty" binding:"omitempty"`
}

type VariantParams struct {
	ProductParam
	VariantID string `uri:"variant_id" binding:"required,uuid"`
}

type VariantResponse struct {
	VariantID  string                   `json:"variant_id"`
	Price      float64                  `json:"price"`
	Stock      int32                    `json:"stock"`
	Discount   int16                    `json:"discount,omitempty"`
	Sku        *string                  `json:"sku,omitempty"`
	Attributes []ProductAttributeDetail `json:"attributes,omitempty"`
	CreatedAt  string                   `json:"created_at"`
	UpdatedAt  string                   `json:"updated_at"`
}

// -------------------------------------------------------------------- API HANDLERS --------------------------------------------------------------------- //

// @Summary Create a variant
// @Description Create a variant
// @Tags variant
// @Accept json
// @Produce json
// @Param product_id path int true "Product ID"
// @Param VariantRequest body VariantRequest true "Create Variant Request"
// @Success 201 {object} repository.ProductVariant
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /variant/{product_id} [post]
func (sv *Server) createVariant(c *gin.Context) {
	var req VariantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	var param ProductParam
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
		ProductID:    uuid.MustParse(param.ID),
		VariantPrice: req.Price,
		VariantStock: req.Stock,
		Discount:     req.Discount,
		VariantSku:   req.Sku,
		Attributes:   createAttributeParams,
	}

	variant, err := sv.repo.CreateVariantTx(c, createVariantParam)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, errors.New("product not found"))
			return
		}
		if errors.Is(err, repository.ErrForeignKeyViolation) {
			c.JSON(http.StatusBadRequest, errors.New("variant is already exist"))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusCreated, GenericResponse[repository.ProductVariant]{&variant, nil, nil})
}

// @Summary Update a variant
// @Description Update a variant
// @Tags variant
// @Accept json
// @Produce json
// @Param variant_id path int true "Variant ID"
// @Param UpdateVariantRequest body UpdateVariantRequest true "Update Variant Request"
// @Success 200 {object} repository.ProductVariant
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /product/{product_id}/variant/{variant_id} [put]
func (sv *Server) updateVariant(c *gin.Context) {
	var params VariantParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	var req UpdateVariantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	if req.Name == nil && req.Sku == nil && req.Price == nil && req.Stock == nil && len(req.Attributes) == 0 {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("at least one field is required")))
		return
	}

	variant, err := sv.repo.GetVariantByID(c, repository.GetVariantByIDParams{
		VariantID: uuid.MustParse(params.VariantID),
		ProductID: uuid.MustParse(params.ID),
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, errors.New("variant not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	updateParams := repository.UpdateVariantParams{
		VariantID: variant.VariantID,
	}

	if req.Sku != nil {
		updateParams.Sku = utils.GetPgTypeText(*req.Sku)
	}
	if req.Price != nil {
		updateParams.Price = utils.GetPgNumericFromFloat(*req.Price)
	}
	if req.Stock != nil {
		updateParams.StockQuantity = utils.GetPgTypeInt4(*req.Stock)
	}

	if len(req.Attributes) > 0 {
		for _, attr := range req.Attributes {
			_, err := sv.repo.UpdateVariantAttribute(c, repository.UpdateVariantAttributeParams{
				VariantAttributeID: attr.ID,
				Value:              utils.GetPgTypeText(attr.Value),
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
	var params VariantParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	variant, err := sv.repo.GetVariantByID(c, repository.GetVariantByIDParams{
		VariantID: uuid.MustParse(params.VariantID),
		ProductID: uuid.MustParse(params.ID),
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, errors.New("variant not found"))
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

// @Summary Get a variant
// @Description Get a variant
// @Tags variant
// @Accept json
// @Produce json
// @Param variant_id path int true "Variant ID"
// @Success 200 {object} VariantResponse
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /variant/{variant_id} [get]
func (sv *Server) getVariant(c *gin.Context) {
	var params VariantParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	variantRows, err := sv.repo.GetVariantDetails(c, uuid.MustParse(params.ID))
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
	resp := VariantResponse{
		VariantID:  variant.VariantID.String(),
		Price:      price.Float64,
		Stock:      variant.StockQuantity,
		Discount:   variant.Discount,
		Attributes: make([]ProductAttributeDetail, 0),
		CreatedAt:  variant.CreatedAt.String(),
		UpdatedAt:  variant.UpdatedAt.String(),
	}

	if variant.Sku.Valid {
		resp.Sku = &variant.Sku.String
	}
	for _, attr := range variantRows {
		resp.Attributes = append(resp.Attributes, ProductAttributeDetail{
			Name:  attr.AttributeName,
			Value: []string{attr.Value},
		})
	}

	c.JSON(http.StatusOK, GenericResponse[VariantResponse]{&resp, nil, nil})
}

func (sv *Server) getVariants(c *gin.Context) {
	var params ProductParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	variants, err := sv.repo.GetVariantByProductID(c, uuid.MustParse(params.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if len(variants) == 0 {
		c.JSON(http.StatusNotFound, mapErrResp(repository.ErrRecordNotFound))
		return
	}
	variantResponses := make([]VariantResponse, 0)
	for _, variant := range variants {
		if len(variantResponses) == 0 || variantResponses[len(variantResponses)-1].VariantID != variant.VariantID.String() {
			price, _ := variant.Price.Float64Value()
			resp := VariantResponse{
				VariantID:  variant.VariantID.String(),
				Price:      price.Float64,
				Stock:      variant.StockQuantity,
				Attributes: make([]ProductAttributeDetail, 0),
				CreatedAt:  variant.CreatedAt.String(),
				UpdatedAt:  variant.UpdatedAt.String(),
			}
			if variant.Sku.Valid {
				resp.Sku = &variant.Sku.String
			}
			variantResponses = append(variantResponses, resp)
		} else {
			latest := &variantResponses[len(variantResponses)-1]
			latest.Attributes = append(latest.Attributes, ProductAttributeDetail{
				ID:    variant.VariantAttributeID,
				Name:  variant.AttributeName,
				Value: []string{variant.AttributeValue},
			})
		}
	}

	c.JSON(http.StatusOK, GenericListResponse[VariantResponse]{variantResponses, int64(len(variantResponses)), nil, nil})
}
