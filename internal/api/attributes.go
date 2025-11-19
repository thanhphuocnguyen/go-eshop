package api

import (
	"errors"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
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
	var params AttributeRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	attribute, err := sv.repo.CreateAttribute(c, params.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InvalidBodyCode, err))
		return
	}

	attributeResp := AttributeRespModel{
		ID:   attribute.ID,
		Name: attribute.Name,
	}

	c.JSON(http.StatusCreated, createDataResp(c, attributeResp, nil, nil))
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
	var attributeParam AttributeParam
	if err := c.ShouldBindUri(&attributeParam); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	attr, err := sv.repo.GetAttributeByID(c, attributeParam.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	attributeResp := AttributeRespModel{
		Name: attr.Name,
		ID:   attr.ID,
	}

	values, err := sv.repo.GetAttributeValues(c, attr.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}
	attributeResp.Values = make([]AttributeValue, len(values))

	for i, val := range values {
		attributeResp.Values[i] = AttributeValue{
			ID:    val.ID,
			Value: val.Value,
		}
	}

	c.JSON(http.StatusOK, createDataResp(c, attributeResp, nil, nil))
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
	var queries GetAttributesQuery
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	attributeRows, err := sv.repo.GetAttributes(c, queries.IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	var attributeResp = []AttributeRespModel{}
	for i := range attributeRows {
		attrVal := attributeRows[i]
		if i == 0 || attributeRows[i].ID != attributeRows[i-1].ID {
			attributeResp = append(attributeResp, AttributeRespModel{
				ID:     attrVal.ID,
				Name:   attrVal.Name,
				Values: []AttributeValue{},
			})
			if attrVal.AttrValueID != nil {
				id := *attrVal.AttrValueID
				attributeResp[len(attributeResp)-1].Values = append(attributeResp[len(attributeResp)-1].Values, AttributeValue{
					ID:    id,
					Value: *attrVal.AttrValue,
				})
			}
		} else if attrVal.AttrValueID != nil {
			id := *attrVal.AttrValueID
			attributeResp[len(attributeResp)-1].Values = append(attributeResp[len(attributeResp)-1].Values, AttributeValue{
				ID:    id,
				Value: *attrVal.AttrValue,
			})
		}
	}

	c.JSON(http.StatusOK, createDataResp(c, attributeResp, nil, nil))
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
	var uri UriIDParam
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	attrs, err := sv.repo.GetProductAttributeValuesByProductID(c, uuid.MustParse(uri.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	resp := make([]AttributeRespModel, 0)
	for _, attr := range attrs {

		if slices.ContainsFunc(resp, func(a AttributeRespModel) bool {
			return *attr.AttributeID == a.ID
		}) {
			// push value to existing attribute
			for i, r := range resp {
				if r.ID == *attr.AttributeID {
					resp[i].Values = append(resp[i].Values, AttributeValue{
						ID:    *attr.AttributeValueID,
						Value: *attr.AttributeValue,
					})
					break
				}
			}
		} else {
			// create new attribute
			attrResp := AttributeRespModel{
				ID:   *attr.AttributeID,
				Name: *attr.AttributeName,
				Values: []AttributeValue{
					{
						ID:    *attr.AttributeValueID,
						Value: *attr.AttributeValue,
					},
				},
			}
			resp = append(resp, attrResp)
		}
	}

	c.JSON(http.StatusOK, createDataResp(c, resp, nil, nil))
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
	var param AttributeParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}
	var req AttributeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	attr, err := sv.repo.GetAttributeByID(c, param.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	attribute, err := sv.repo.UpdateAttribute(c, repository.UpdateAttributeParams{
		ID:   attr.ID,
		Name: req.Name,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusOK, createDataResp(c, attribute, nil, nil))
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
	var param AttributeParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}
	var req AttributeValuesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	attr, err := sv.repo.GetAttributeByID(c, param.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	obj, err := sv.repo.CreateAttributeValue(c, repository.CreateAttributeValueParams{
		AttributeID: attr.ID,
		Value:       req.Value,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusCreated, createDataResp(c, obj, nil, nil))
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
	var param AttributeParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}
	if param.ValueID == nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, nil))
		return
	}
	var req AttributeValuesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	attr, err := sv.repo.GetAttributeByID(c, param.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, createErr(NotFoundCode, err))
		return
	}

	res, err := sv.repo.UpdateAttributeValue(c, repository.UpdateAttributeValueParams{
		AttributeID: attr.ID,
		ID:          *param.ValueID,
		Value:       req.Value,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusOK, createDataResp(c, res, nil, nil))
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
	var param AttributeParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	if param.ValueID == nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, nil))
		return
	}

	attr, err := sv.repo.GetAttributeByID(c, param.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	err = sv.repo.DeleteAttributeValueByValueID(c, repository.DeleteAttributeValueByValueIDParams{
		AttributeID: attr.ID,
		ID:          *param.ValueID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusNoContent, createDataResp(c, struct{}{}, nil, nil))
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
	var params AttributeParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	attr, err := sv.repo.GetAttributeByID(c, params.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErr(NotFoundCode, err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	err = sv.repo.DeleteAttribute(c, attr.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.Status(http.StatusNoContent)
}
