package api

import (
	"errors"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
)

// @Summary Create an attribute
// @Description Create an attribute
// @Tags attributes
// @Accept json
// @Produce json
// @Param params body AttributeValuesReq true "Attribute name"
// @Success 201 {object} ApiResponse[AttributeRespModel]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /attributes [post]
func (sv *Server) CreateAttributeHandler(c *gin.Context) {
	var req models.AttributeModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	attribute, err := sv.repo.CreateAttribute(c, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	attributeResp := dto.AttributeDetail{
		ID:   attribute.ID,
		Name: attribute.Name,
	}

	c.JSON(http.StatusCreated, dto.CreateDataResp(c, attributeResp, nil, nil))
}

// @Summary Get an attribute
// @Description Get an attribute
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Success 200 {object} ApiResponse[AttributeRespModel]
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /attributes/{id} [get]
func (sv *Server) GetAttributeByIDHandler(c *gin.Context) {
	var attributeParam models.AttributeParam
	if err := c.ShouldBindUri(&attributeParam); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	attr, err := sv.repo.GetAttributeByID(c, attributeParam.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	attributeResp := dto.AttributeDetail{
		Name: attr.Name,
		ID:   attr.ID,
	}

	values, err := sv.repo.GetAttributeValues(c, attr.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	attributeResp.Values = make([]dto.AttributeValueDetail, len(values))

	for i, val := range values {
		attributeResp.Values[i] = dto.AttributeValueDetail{
			ID:    val.ID,
			Value: val.Value,
		}
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, attributeResp, nil, nil))
}

// @Summary Get all attributes
// @Description Get all attributes
// @Tags attributes
// @Accept json
// @Produce json
// @Success 200 {object} ApiResponse[[]AttributeRespModel]
// @Failure 500 {object} ErrorResp
// @Router /attributes [get]
func (sv *Server) GetAttributesHandler(c *gin.Context) {
	var queries models.AttributesQuery
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	attributeRows, err := sv.repo.GetAttributes(c, queries.IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	var attributeResp = []dto.AttributeDetail{}
	for i := range attributeRows {
		attrVal := attributeRows[i]
		if i == 0 || attributeRows[i].ID != attributeRows[i-1].ID {
			attributeResp = append(attributeResp, dto.AttributeDetail{
				ID:     attrVal.ID,
				Name:   attrVal.Name,
				Values: []dto.AttributeValueDetail{},
			})
			if attrVal.AttrValueID != nil {
				id := *attrVal.AttrValueID
				attributeResp[len(attributeResp)-1].Values = append(attributeResp[len(attributeResp)-1].Values, dto.AttributeValueDetail{
					ID:    id,
					Value: *attrVal.AttrValue,
				})
			}
		} else if attrVal.AttrValueID != nil {
			id := *attrVal.AttrValueID
			attributeResp[len(attributeResp)-1].Values = append(attributeResp[len(attributeResp)-1].Values, dto.AttributeValueDetail{
				ID:    id,
				Value: *attrVal.AttrValue,
			})
		}
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, attributeResp, nil, nil))
}

// @Summary Get attributes and their values by for a product
// @Description Get attributes and their values for a product
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} ApiResponse[[]AttributeRespModel]
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /attributes/product/{id} [get]
func (sv *Server) GetAttributeValuesForProductHandler(c *gin.Context) {
	var uri models.UriIDParam
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	attrs, err := sv.repo.GetProductAttributeValuesByProductID(c, uuid.MustParse(uri.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	resp := make([]dto.AttributeDetail, 0)
	for _, attr := range attrs {

		if slices.ContainsFunc(resp, func(a dto.AttributeDetail) bool {
			return *attr.AttributeID == a.ID
		}) {
			// push value to existing attribute
			for i, r := range resp {
				if r.ID == *attr.AttributeID {
					resp[i].Values = append(resp[i].Values, dto.AttributeValueDetail{
						ID:    *attr.AttributeValueID,
						Value: *attr.AttributeValue,
					})
					break
				}
			}
		} else {
			// create new attribute
			attrResp := dto.AttributeDetail{
				ID:   *attr.AttributeID,
				Name: *attr.AttributeName,
				Values: []dto.AttributeValueDetail{
					{
						ID:    *attr.AttributeValueID,
						Value: *attr.AttributeValue,
					},
				},
			}
			resp = append(resp, attrResp)
		}
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, resp, nil, nil))
}

// @Summary Update an attribute
// @Description Update an attribute
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Param params body AttributeRequest true "Attribute name"
// @Success 200 {object} ApiResponse[repository.Attribute]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /attributes/{id} [put]
func (sv *Server) UpdateAttributeHandler(c *gin.Context) {
	var param models.AttributeParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	var req models.AttributeModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	attr, err := sv.repo.GetAttributeByID(c, param.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	attribute, err := sv.repo.UpdateAttribute(c, repository.UpdateAttributeParams{
		ID:   attr.ID,
		Name: req.Name,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, attribute, nil, nil))
}

// @Summary Add new attribute value
// @Description Add new attribute value
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Param params body AttributeValuesReq true "Attribute value"
// @Success 200 {object} ApiResponse[bool]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /attributes/{id}/create [post]
func (sv *Server) AddAttributeValueHandler(c *gin.Context) {
	var param models.AttributeParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	var req models.AttributeValueModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	attr, err := sv.repo.GetAttributeByID(c, param.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	obj, err := sv.repo.CreateAttributeValue(c, repository.CreateAttributeValueParams{
		AttributeID: attr.ID,
		Value:       req.Value,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusCreated, dto.CreateDataResp(c, obj, nil, nil))
}

// @Summary update attribute value
// @Description update attribute value
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Param params body AttributeValuesReq true "Attribute value"
// @Success 200 {object} ApiResponse[bool]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /attributes/{id}/update/{valueId} [put]
func (sv *Server) UpdateAttrValueHandler(c *gin.Context) {
	var param models.AttributeParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	if param.ValueID == nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, nil))
		return
	}
	var req models.AttributeValueModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	attr, err := sv.repo.GetAttributeByID(c, param.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, err))
		return
	}

	res, err := sv.repo.UpdateAttributeValue(c, repository.UpdateAttributeValueParams{
		AttributeID: attr.ID,
		ID:          *param.ValueID,
		Value:       req.Value,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, res, nil, nil))
}

// @Summary remove an attribute value
// @Description remove an attribute value
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Success 204 {object} nil
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /attributes/{id}/remove/{valueId} [delete]
func (sv *Server) RemoveAttrValueHandler(c *gin.Context) {
	var param models.AttributeParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	if param.ValueID == nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, nil))
		return
	}

	attr, err := sv.repo.GetAttributeByID(c, param.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	err = sv.repo.DeleteAttributeValueByValueID(c, repository.DeleteAttributeValueByValueIDParams{
		AttributeID: attr.ID,
		ID:          *param.ValueID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusNoContent, dto.CreateDataResp(c, struct{}{}, nil, nil))
}

// @Summary Delete an attribute
// @Description Delete an attribute
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Success 204 {object} nil
// @Failure 500 {object} ErrorResp
// @Router /attributes/{id} [delete]
func (sv *Server) RemoveAttributeHandler(c *gin.Context) {
	var params models.AttributeParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	attr, err := sv.repo.GetAttributeByID(c, params.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, err))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	err = sv.repo.DeleteAttribute(c, attr.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.Status(http.StatusNoContent)
}
