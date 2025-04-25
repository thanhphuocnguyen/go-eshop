package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

// ------------------------------ API Models ------------------------------
type AttributeValue struct {
	ID           int32   `json:"id"`
	Value        string  `json:"value"`
	DisplayValue *string `json:"display_value"`
	IsActive     *bool   `json:"is_active"`
	DisplayOrder *int16  `json:"display_order"`
}

type AttributeValueRequest struct {
	Value        string  `json:"value" binding:"required"`
	DisplayValue *string `json:"display_value" binding:"omitempty"`
	DisplayOrder *int16  `json:"display_order" binding:"omitempty,min=0"`
	IsActive     *bool   `json:"is_active" binding:"omitempty"`
}

type UpdateAttributeValueRequest struct {
	ID                    *int32 `json:"id" binding:"omitempty"`
	AttributeValueRequest `json:",inline"`
}

type AttributeResponse struct {
	ID        int32            `json:"id"`
	Name      string           `json:"name"`
	Values    []AttributeValue `json:"values,omitempty"`
	CreatedAt string           `json:"created_at"`
	UpdatedAt string           `json:"updated_at"`
}

type CreateAttributeRequest struct {
	Name   string                  `json:"name" binding:"required"`
	Values []AttributeValueRequest `json:"values,omitempty"`
}

type UpdateAttributeRequest struct {
	Name   string                        `json:"name" binding:"required"`
	Values []UpdateAttributeValueRequest `json:"values,omitempty"`
}

type AttributeParam struct {
	ID int32 `uri:"id" binding:"required"`
}

type GetAttributesQuery struct {
	IDs []int32 `form:"ids" binding:"omitempty"`
}

// ------------------------------ API Handlers ------------------------------

// @Summary Create an attribute
// @Description Create an attribute
// @Tags attributes
// @Accept json
// @Produce json
// @Param params body CreateAttributeRequest true "Attribute name"
// @Success 201 {object} ApiResponse[AttributeResponse]
// @Failure 400 {object} ApiResponse[AttributeResponse]
// @Failure 500 {object} ApiResponse[AttributeResponse]
// @Router /attributes [post]
func (sv *Server) createAttribute(c *gin.Context) {
	var params CreateAttributeRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[AttributeResponse](http.StatusBadRequest, "", err))
		return
	}

	attribute, err := sv.repo.CreateAttribute(c, params.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[AttributeResponse](http.StatusBadRequest, "", err))
		return
	}
	attributeValues := []AttributeValue{}
	if len(params.Values) > 0 {
		for _, value := range params.Values {
			createParams := repository.CreateAttributeValueParams{
				AttributeID: attribute.ID,
				Value:       value.Value,
			}
			if value.DisplayValue != nil {
				createParams.DisplayValue = utils.GetPgTypeText(*value.DisplayValue)
			}
			if value.DisplayOrder != nil {
				createParams.DisplayOrder = *value.DisplayOrder
			}
			value, err := sv.repo.CreateAttributeValue(c, createParams)
			if err != nil {
				c.JSON(http.StatusInternalServerError, createErrorResponse[AttributeResponse](http.StatusBadRequest, "", err))
				return
			}

			attributeValues = append(attributeValues, AttributeValue{
				ID:    value.ID,
				Value: value.Value,
			})

		}
	}

	attributeResp := AttributeResponse{
		ID:        attribute.ID,
		Name:      attribute.Name,
		Values:    attributeValues,
		CreatedAt: attribute.CreatedAt.String(),
	}

	c.JSON(http.StatusCreated, createSuccessResponse(c, attributeResp, "", nil, nil))
}

// @Summary Get an attribute
// @Description Get an attribute
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Success 200 {object} ApiResponse[AttributeResponse]
// @Failure 404 {object} ApiResponse[AttributeResponse]
// @Failure 500 {object} ApiResponse[AttributeResponse]
// @Router /attributes/{id} [get]
func (sv *Server) getAttributeByIDHandler(c *gin.Context) {
	var attributeParam AttributeParam
	if err := c.ShouldBindUri(&attributeParam); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[AttributeResponse](http.StatusBadRequest, "", err))
		return
	}

	attributeRows, err := sv.repo.GetAttributeByID(c, attributeParam.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[AttributeResponse](http.StatusBadRequest, "", err))
		return
	}

	if attributeRows == nil {
		c.JSON(http.StatusNotFound, createErrorResponse[AttributeResponse](http.StatusBadRequest, "", errors.New("attribute not found")))
		return
	}

	attributeResp := AttributeResponse{
		Name:      attributeRows[0].Name,
		ID:        attributeRows[0].ID,
		CreatedAt: attributeRows[0].CreatedAt.String(),
		UpdatedAt: attributeRows[0].CreatedAt.String(),
		Values:    []AttributeValue{},
	}

	for i := 0; i < len(attributeRows); i++ {
		if !attributeRows[i].AttributeValueID.Valid {
			continue
		}
		attributeResp.Values = append(attributeResp.Values, AttributeValue{
			ID:           attributeRows[i].AttributeValueID.Int32,
			Value:        attributeRows[i].Value.String,
			DisplayValue: &attributeRows[i].DisplayValue.String,
			IsActive:     &attributeRows[i].AttributeValueIsActive.Bool,
			DisplayOrder: &attributeRows[i].DisplayOrder.Int16,
		})
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, attributeResp, "", nil, nil))
}

// @Summary Get all attributes
// @Description Get all attributes
// @Tags attributes
// @Accept json
// @Produce json
// @Success 200 {object} ApiResponse[[]AttributeResponse]
// @Failure 500 {object} ApiResponse[[]AttributeResponse]
// @Router /attributes [get]
func (sv *Server) getAttributesHandler(c *gin.Context) {
	var queries GetAttributesQuery

	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[[]AttributeResponse](http.StatusBadRequest, "", err))
		return
	}

	attributeRows, err := sv.repo.GetAttributes(c, queries.IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[[]AttributeResponse](http.StatusBadRequest, "", err))
		return
	}

	cnt, err := sv.repo.CountAttributes(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[[]AttributeResponse](http.StatusBadRequest, "", err))
		return
	}

	var attributeResp = []AttributeResponse{}
	for i := 0; i < len(attributeRows); i++ {
		attrVal := attributeRows[i]
		if i == 0 || attributeRows[i].ID != attributeRows[i-1].ID {
			attributeResp = append(attributeResp, AttributeResponse{
				ID:        attrVal.ID,
				Name:      attrVal.Name,
				CreatedAt: attrVal.CreatedAt.String(),
				Values:    []AttributeValue{},
			})
			if attrVal.AttributeValueID.Valid {
				attributeResp[len(attributeResp)-1].Values = append(attributeResp[len(attributeResp)-1].Values, AttributeValue{
					ID:           attrVal.AttributeValueID.Int32,
					Value:        attrVal.Value.String,
					DisplayValue: &attrVal.DisplayValue.String,
					DisplayOrder: &attrVal.DisplayOrder.Int16,
					IsActive:     &attrVal.AttributeValueIsActive.Bool,
				})
			}
		} else if attrVal.AttributeValueID.Valid {
			attributeResp[len(attributeResp)-1].Values = append(attributeResp[len(attributeResp)-1].Values, AttributeValue{
				ID:           attrVal.AttributeValueID.Int32,
				Value:        attrVal.Value.String,
				DisplayValue: &attrVal.DisplayValue.String,
				IsActive:     &attrVal.AttributeValueIsActive.Bool,
				DisplayOrder: &attrVal.DisplayOrder.Int16,
			})
		}
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, attributeResp, "", &Pagination{Total: cnt}, nil))
}

// @Summary Update an attribute
// @Description Update an attribute
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Param params body UpdateAttributeRequest true "Attribute name"
// @Success 200 {object} ApiResponse[AttributeResponse]
// @Failure 400 {object} ApiResponse[AttributeResponse]
// @Failure 500 {object} ApiResponse[[]AttributeResponse]
// @Router /attributes/{id} [put]
func (sv *Server) updateAttribute(c *gin.Context) {
	var param AttributeParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[AttributeResponse](http.StatusBadRequest, "", err))
		return
	}

	var req UpdateAttributeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[AttributeResponse](http.StatusBadRequest, "", err))
		return
	}

	existingAttributeRows, err := sv.repo.GetAttributeByID(c, param.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[AttributeResponse](http.StatusBadRequest, "", err))
		return
	}

	if existingAttributeRows == nil {
		c.JSON(http.StatusNotFound, createErrorResponse[AttributeResponse](http.StatusBadRequest, "", errors.New("attribute not found")))
		return
	}

	existed := existingAttributeRows[0]

	currentAttributeValues, err := sv.repo.GetAttributeValues(c, existed.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[AttributeResponse](http.StatusBadRequest, "", err))
		return
	}

	attribute, err := sv.repo.UpdateAttribute(c, repository.UpdateAttributeParams{
		ID:   existed.ID,
		Name: req.Name,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[AttributeResponse](http.StatusBadRequest, "", err))
		return
	}

	currentAttributeValuesMap := make(map[int32]AttributeValue)
	for _, value := range currentAttributeValues {
		currentAttributeValuesMap[value.ID] = AttributeValue{
			ID:           value.ID,
			Value:        value.Value,
			DisplayValue: &value.DisplayValue.String,
			IsActive:     &value.IsActive.Bool,
			DisplayOrder: &value.DisplayOrder,
		}
	}

	response := []AttributeValue{}

	for _, value := range req.Values {
		if value.ID != nil {
			if _, ok := currentAttributeValuesMap[*value.ID]; ok {
				params := repository.UpdateAttributeValueParams{
					ID:    *value.ID,
					Value: value.Value,
				}
				if value.DisplayValue != nil {
					params.DisplayValue = utils.GetPgTypeText(*value.DisplayValue)
				}
				if value.DisplayOrder != nil {
					params.DisplayOrder = utils.GetPgTypeInt2(*value.DisplayOrder)
				}
				updated, err := sv.repo.UpdateAttributeValue(c, params)
				if err != nil {
					c.JSON(http.StatusInternalServerError, createErrorResponse[AttributeResponse](http.StatusBadRequest, "", err))
					return
				}
				delete(currentAttributeValuesMap, *value.ID)
				response = append(response, AttributeValue{
					ID:           updated.ID,
					Value:        updated.Value,
					DisplayValue: &updated.DisplayValue.String,
					IsActive:     &updated.IsActive.Bool,
					DisplayOrder: &updated.DisplayOrder,
				})
			}
		} else {
			createParams := repository.CreateAttributeValueParams{
				AttributeID: existed.ID,
				Value:       value.Value,
			}
			if value.DisplayValue != nil {
				createParams.DisplayValue = utils.GetPgTypeText(*value.DisplayValue)
			}
			if value.DisplayOrder != nil {
				createParams.DisplayOrder = *value.DisplayOrder
			}
			newAttr, err := sv.repo.CreateAttributeValue(c, createParams)
			if err != nil {
				c.JSON(http.StatusInternalServerError, createErrorResponse[AttributeResponse](http.StatusBadRequest, "", err))
				return
			}
			response = append(response, AttributeValue{
				ID:           newAttr.ID,
				Value:        newAttr.Value,
				DisplayValue: &newAttr.DisplayValue.String,
				IsActive:     &newAttr.IsActive.Bool,
				DisplayOrder: &newAttr.DisplayOrder,
			})
		}
	}

	for id := range currentAttributeValuesMap {
		err := sv.repo.DeleteAttributeValue(c, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse[AttributeResponse](http.StatusBadRequest, "", err))
			return
		}
	}

	attributeResp := AttributeResponse{
		ID:        attribute.ID,
		Name:      attribute.Name,
		Values:    response,
		CreatedAt: attribute.CreatedAt.String(),
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, attributeResp, "", nil, nil))
}

// @Summary Delete an attribute
// @Description Delete an attribute
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Success 204 {object} nil
// @Failure 500 {object} ApiResponse[bool]
// @Router /attributes/{id} [delete]
func (sv *Server) deleteAttribute(c *gin.Context) {
	var params AttributeParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[bool](http.StatusBadRequest, "", err))
		return
	}

	attribute, err := sv.repo.GetAttributeByID(c, params.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](http.StatusInternalServerError, "", err))
		return
	}
	if attribute == nil {
		c.JSON(http.StatusNotFound, createErrorResponse[bool](http.StatusNotFound, "", errors.New("attribute not found")))
		return
	}

	err = sv.repo.DeleteAttribute(c, attribute[0].ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](http.StatusInternalServerError, "", err))
		return
	}
	c.Status(http.StatusNoContent)
}
