package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
)

// ------------------------------ API Models ------------------------------
type AttributeValue struct {
	ID           uuid.UUID `json:"id"`
	Code         string    `json:"code"`
	Name         *string   `json:"name"`
	IsActive     *bool     `json:"is_active"`
	DisplayOrder *int16    `json:"display_order"`
}

type AttributeValueRequest struct {
	Code         string  `json:"code" binding:"required"`
	Name         *string `json:"name" binding:"omitempty"`
	DisplayOrder *int16  `json:"display_order" binding:"omitempty,min=0"`
	IsActive     *bool   `json:"is_active" binding:"omitempty"`
}

type UpdateAttributeValueRequest struct {
	ID                    *string `json:"id" binding:"omitempty,uuid"`
	AttributeValueRequest `json:",inline"`
}

type AttributeResponse struct {
	ID        uuid.UUID        `json:"id"`
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
	ID string `uri:"id" binding:"required,uuid"`
}

type GetAttributesQuery struct {
	IDs []uuid.UUID `form:"ids" binding:"omitempty"`
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
func (sv *Server) createAttributeHandler(c *gin.Context) {
	var params CreateAttributeRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[AttributeResponse](InvalidBodyCode, "", err))
		return
	}

	attribute, err := sv.repo.CreateAttribute(c, params.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[AttributeResponse](InvalidBodyCode, "", err))
		return
	}
	attributeValues := []AttributeValue{}
	if len(params.Values) > 0 {
		for _, value := range params.Values {
			createParams := repository.CreateAttributeValueParams{
				AttributeID: attribute.ID,
				Code:        value.Code,
				Name:        *value.Name,
			}

			if value.DisplayOrder != nil {
				createParams.DisplayOrder = *value.DisplayOrder
			}
			value, err := sv.repo.CreateAttributeValue(c, createParams)
			if err != nil {
				c.JSON(http.StatusInternalServerError, createErrorResponse[AttributeResponse](InternalServerErrorCode, "", err))
				return
			}

			attributeValues = append(attributeValues, AttributeValue{
				Code: value.Code,
				Name: &value.Name,
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
		c.JSON(http.StatusBadRequest, createErrorResponse[AttributeResponse](InvalidBodyCode, "", err))
		return
	}

	attributeRows, err := sv.repo.GetAttributeByID(c, uuid.MustParse(attributeParam.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[AttributeResponse](InternalServerErrorCode, "", err))
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
		id, _ := uuid.FromBytes(attributeRows[i].AttributeValueID.Bytes[:])
		attributeResp.Values = append(attributeResp.Values, AttributeValue{
			ID:           id,
			Name:         attributeRows[i].AttrValName,
			Code:         *attributeRows[i].AttrValCode,
			IsActive:     attributeRows[i].AttributeValueIsActive,
			DisplayOrder: attributeRows[i].DisplayOrder,
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
		c.JSON(http.StatusBadRequest, createErrorResponse[[]AttributeResponse](InvalidBodyCode, "", err))
		return
	}
	var cached *[]AttributeResponse
	err := sv.cacheService.Get(c, fmt.Sprintf("attributes-%s", queries.IDs), &cached)
	if cached != nil {
		c.JSON(http.StatusOK, createSuccessResponse(c, &cached, "", nil, nil))
		return
	}

	attributeRows, err := sv.repo.GetAttributes(c, queries.IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[[]AttributeResponse](InternalServerErrorCode, "", err))
		return
	}

	cnt, err := sv.repo.CountAttributes(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[[]AttributeResponse](InternalServerErrorCode, "", err))
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
				id, _ := uuid.FromBytes(attrVal.AttributeValueID.Bytes[:])
				attributeResp[len(attributeResp)-1].Values = append(attributeResp[len(attributeResp)-1].Values, AttributeValue{
					ID:           id,
					Code:         *attrVal.AttrValCode,
					Name:         attrVal.AttrValName,
					DisplayOrder: attrVal.DisplayOrder,
					IsActive:     attrVal.AttributeValueIsActive,
				})
			}
		} else if attrVal.AttributeValueID.Valid {
			id, _ := uuid.FromBytes(attrVal.AttributeValueID.Bytes[:])
			attributeResp[len(attributeResp)-1].Values = append(attributeResp[len(attributeResp)-1].Values, AttributeValue{
				ID:           id,
				Code:         *attrVal.AttrValCode,
				Name:         attrVal.AttrValName,
				IsActive:     attrVal.AttributeValueIsActive,
				DisplayOrder: attrVal.DisplayOrder,
			})
		}
	}
	if err := sv.cacheService.Set(c, fmt.Sprintf("attributes-%s", queries.IDs), attributeResp, nil); err != nil {
		log.Error().Err(err).Msg("failed to cache attributes")
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
func (sv *Server) updateAttributeHandler(c *gin.Context) {
	var param AttributeParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[AttributeResponse](InvalidBodyCode, "", err))
		return
	}

	var req UpdateAttributeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[AttributeResponse](InvalidBodyCode, "", err))
		return
	}

	existingAttributeRows, err := sv.repo.GetAttributeByID(c, uuid.MustParse(param.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[AttributeResponse](InternalServerErrorCode, "", err))
		return
	}

	existed := existingAttributeRows[0]

	currentAttributeValues, err := sv.repo.GetAttributeValues(c, existed.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[AttributeResponse](InternalServerErrorCode, "", err))
		return
	}

	attribute, err := sv.repo.UpdateAttribute(c, repository.UpdateAttributeParams{
		ID:   existed.ID,
		Name: req.Name,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[AttributeResponse](InternalServerErrorCode, "", err))
		return
	}

	currentAttributeValuesMap := make(map[string]AttributeValue)
	for _, value := range currentAttributeValues {
		currentAttributeValuesMap[value.ID.String()] = AttributeValue{
			ID:           value.ID,
			Code:         value.Code,
			Name:         &value.Name,
			IsActive:     value.IsActive,
			DisplayOrder: &value.DisplayOrder,
		}
	}

	attrVals := []AttributeValue{}

	for _, value := range req.Values {
		if value.ID != nil {
			if _, ok := currentAttributeValuesMap[*value.ID]; ok {
				params := repository.UpdateAttributeValueParams{
					ID: uuid.MustParse(*value.ID),
				}
				if value.Code != "" {
					params.Code = &value.Code
				}
				if value.Name != nil {
					params.Name = value.Name
				}
				if value.DisplayOrder != nil {
					params.DisplayOrder = value.DisplayOrder
				}
				updated, err := sv.repo.UpdateAttributeValue(c, params)
				if err != nil {
					c.JSON(http.StatusInternalServerError, createErrorResponse[AttributeResponse](InternalServerErrorCode, "", err))
					return
				}
				delete(currentAttributeValuesMap, *value.ID)
				attrVals = append(attrVals, AttributeValue{
					ID:           updated.ID,
					Code:         updated.Code,
					Name:         &updated.Name,
					IsActive:     updated.IsActive,
					DisplayOrder: &updated.DisplayOrder,
				})
			}
		} else {
			createParams := repository.CreateAttributeValueParams{
				AttributeID: existed.ID,
				Code:        value.Code,
			}
			if value.Name != nil {
				createParams.Name = *value.Name
			}
			if value.DisplayOrder != nil {
				createParams.DisplayOrder = *value.DisplayOrder
			}
			newAttr, err := sv.repo.CreateAttributeValue(c, createParams)
			if err != nil {
				c.JSON(http.StatusInternalServerError, createErrorResponse[AttributeResponse](InternalServerErrorCode, "", err))
				return
			}
			attrVals = append(attrVals, AttributeValue{
				ID:           newAttr.ID,
				Code:         newAttr.Code,
				Name:         &newAttr.Name,
				IsActive:     newAttr.IsActive,
				DisplayOrder: &newAttr.DisplayOrder,
			})
		}
	}

	for id := range currentAttributeValuesMap {
		err := sv.repo.DeleteAttributeValue(c, uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse[AttributeResponse](InternalServerErrorCode, "", err))
			return
		}
	}

	attributeResp := AttributeResponse{
		ID:        attribute.ID,
		Name:      attribute.Name,
		Values:    attrVals,
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
func (sv *Server) deleteAttributeHandler(c *gin.Context) {
	var params AttributeParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[bool](InvalidBodyCode, "", err))
		return
	}

	attribute, err := sv.repo.GetAttributeByID(c, uuid.MustParse(params.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[bool](NotFoundCode, "", err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
		return
	}

	err = sv.repo.DeleteAttribute(c, attribute[0].ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
		return
	}
	c.Status(http.StatusNoContent)
}
