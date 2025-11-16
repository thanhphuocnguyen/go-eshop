package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
)

// @Summary Create an attribute
// @Description Create an attribute
// @Tags attributes
// @Accept json
// @Produce json
// @Param params body CreateAttributeRequest true "Attribute name"
// @Success 201 {object} ApiResponse[AttributeResponse]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /attributes [post]
func (sv *Server) CreateAttributeHandler(c *gin.Context) {
	var params CreateAttributeRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, "", err))
		return
	}

	attribute, err := sv.repo.CreateAttribute(c, params.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InvalidBodyCode, "", err))
		return
	}

	attrValues := []AttributeValue{}
	if len(params.Values) > 0 {
		for _, value := range params.Values {
			createParams := repository.CreateAttributeValueParams{
				AttributeID: attribute.ID,
				Value:       value,
			}

			attrVal, err := sv.repo.CreateAttributeValue(c, createParams)
			if err != nil {
				c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, "", err))
				return
			}

			attrValues = append(attrValues, AttributeValue{
				Value: attrVal.Value,
			})

		}
	}

	attributeResp := AttributeResponse{
		ID:     attribute.ID,
		Name:   attribute.Name,
		Values: attrValues,
	}

	c.JSON(http.StatusCreated, createDataResp(c, attributeResp, "", nil, nil))
}

// @Summary Get an attribute
// @Description Get an attribute
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Success 200 {object} ApiResponse[AttributeResponse]
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /attributes/{id} [get]
func (sv *Server) GetAttributeByIDHandler(c *gin.Context) {
	var attributeParam AttributeParam
	if err := c.ShouldBindUri(&attributeParam); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, "", err))
		return
	}

	attr, err := sv.repo.GetAttributeByID(c, attributeParam.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, "", err))
		return
	}

	attributeResp := AttributeResponse{
		Name: attr.Name,
		ID:   attr.ID,
	}

	values, err := sv.repo.GetAttributeValues(c, attr.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, "", err))
		return
	}
	attributeResp.Values = make([]AttributeValue, len(values))

	for i, val := range values {
		attributeResp.Values[i] = AttributeValue{
			ID:    val.ID,
			Value: val.Value,
		}
	}

	c.JSON(http.StatusOK, createDataResp(c, attributeResp, "", nil, nil))
}

// @Summary Get all attributes
// @Description Get all attributes
// @Tags attributes
// @Accept json
// @Produce json
// @Success 200 {object} ApiResponse[[]AttributeResponse]
// @Failure 500 {object} ErrorResp
// @Router /attributes [get]
func (sv *Server) getAttributesHandler(c *gin.Context) {
	var queries GetAttributesQuery
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, "", err))
		return
	}
	// var cached *[]AttributeResponse
	// err := sv.cachesrv.Get(c, fmt.Sprintf("attributes-%s", queries.IDs), &cached)
	// if err != nil && !errors.Is(err, cachesrv.ErrCacheMiss) {
	// 	log.Error().Err(err).Msg("failed to get attributes from cache")
	// }
	// if cached != nil {
	// 	c.JSON(http.StatusOK, createSuccessResponse(c, &cached, "", nil, nil))
	// 	return
	// }

	attributeRows, err := sv.repo.GetAttributes(c, queries.IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, "", err))
		return
	}

	cnt, err := sv.repo.CountAttributes(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, "", err))
		return
	}

	var attributeResp = []AttributeResponse{}
	for i := range attributeRows {
		attrVal := attributeRows[i]
		if i == 0 || attributeRows[i].ID != attributeRows[i-1].ID {
			attributeResp = append(attributeResp, AttributeResponse{
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
	// if err := sv.cachesrv.Set(c, fmt.Sprintf("attributes-%s", queries.IDs), attributeResp, nil); err != nil {
	// 	log.Error().Err(err).Msg("failed to cache attributes")
	// }

	c.JSON(http.StatusOK, createDataResp(c, attributeResp, "", &Pagination{Total: cnt}, nil))
}

// @Summary Update an attribute
// @Description Update an attribute
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Param params body UpdateAttributeRequest true "Attribute name"
// @Success 200 {object} ApiResponse[repository.Attribute]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /attributes/{id} [put]
func (sv *Server) updateAttributeHandler(c *gin.Context) {
	var param AttributeParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, "", err))
		return
	}
	var req UpdateAttributeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, "", err))
		return
	}

	attr, err := sv.repo.GetAttributeByID(c, param.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, "", err))
		return
	}

	attribute, err := sv.repo.UpdateAttribute(c, repository.UpdateAttributeParams{
		ID:   attr.ID,
		Name: req.Name,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, "", err))
		return
	}

	c.JSON(http.StatusOK, createDataResp(c, attribute, "", nil, nil))
}

// @Summary Add new attribute values
// @Description Add new attribute values
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Param params body AddAttributeValuesRequest true "Attribute values"
// @Success 200 {object} ApiResponse[bool]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /attributes/{id}/create [put]
func (sv *Server) AddAttrValuesHandler(c *gin.Context) {
	var param AttributeParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, "", err))
		return
	}
	var req []string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, "", err))
		return
	}

	attr, err := sv.repo.GetAttributeByID(c, param.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, "", err))
		return
	}

	createParams := make([]repository.CreateAttributeValuesParams, len(req))
	for i, val := range req {
		createParams[i] = repository.CreateAttributeValuesParams{
			AttributeID: attr.ID,
			Value:       val,
		}
	}
	_, err = sv.repo.CreateAttributeValues(c, createParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, "", err))
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Add new attribute values
// @Description Add new attribute values
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Param params body AddAttributeValuesRequest true "Attribute values"
// @Success 200 {object} ApiResponse[bool]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /attributes/{id}/remove/{valueId} [put]
func (sv *Server) RemoveAttrValueHandler(c *gin.Context) {
	var param AttributeParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, "", err))
		return
	}

	if param.ValueID == nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, "valueId is required", nil))
		return
	}

	attr, err := sv.repo.GetAttributeByID(c, param.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, "", err))
		return
	}

	err = sv.repo.DeleteAttributeValueByValueID(c, repository.DeleteAttributeValueByValueIDParams{
		AttributeID: attr.ID,
		ID:          *param.ValueID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, "", err))
		return
	}

	c.Status(http.StatusNoContent)
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
func (sv *Server) deleteAttributeHandler(c *gin.Context) {
	var params AttributeParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, "", err))
		return
	}

	attr, err := sv.repo.GetAttributeByID(c, params.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErr(NotFoundCode, "", err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, "", err))
		return
	}

	err = sv.repo.DeleteAttribute(c, attr.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, "", err))
		return
	}

	c.Status(http.StatusNoContent)
}
