package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/util"
)

// ------------------------------ API Models ------------------------------
type AttributeValue struct {
	ID    int32   `json:"id"`
	Value string  `json:"value"`
	Color *string `json:"color,omitempty"`
}

type CreateAttributeValueParams struct {
	AttributeID int32            `json:"attribute_id" binding:"required"`
	Values      []AttributeValue `json:"values" binding:"required"`
}

type AttributeValueParams struct {
	AttributeParam
	ValueID int32 `uri:"value_id" binding:"required"`
}

// ------------------------------ API Handlers ------------------------------
// @Summary Create attribute values
// @Description Create attribute values
// @Tags attributes
// @Accept json
// @Produce json
// @Param attribute_id path int true "Attribute ID"
// @Param values body CreateAttributeValueParams true "Attribute values"
// @Success 201 {object} GenericResponse[int64]
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /attributes/{attribute_id}/values [post]
func (sv *Server) createAttributeValues(c *gin.Context) {
	var params CreateAttributeValueParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	createParams := make([]repository.CreateBulkAttributeValuesParams, len(params.Values))
	for i, value := range params.Values {
		createParams[i] = repository.CreateBulkAttributeValuesParams{
			AttributeID:    params.AttributeID,
			AttributeValue: value.Value,
		}
		if value.Color != nil {
			createParams[i].Color = util.GetPgTypeText(*value.Color)
		}
	}

	attributeValue, err := sv.repo.CreateBulkAttributeValues(c, createParams)

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	msg := "Attribute values created successfully"
	c.JSON(http.StatusCreated, GenericResponse[int64]{&attributeValue, &msg, nil})
}

// @Summary Remove attribute value
// @Description Remove attribute value
// @Tags attributes
// @Accept json
// @Produce json
// @Param attribute_id path int true "Attribute ID"
// @Param id path int true "Attribute value ID"
// @Success 204 {object} GenericResponse[bool]
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /attributes/{attribute_id}/values/{id} [delete]
func (sv *Server) deleteAttributeValue(c *gin.Context) {
	var attributeParam AttributeValueParams
	if err := c.ShouldBindUri(&attributeParam); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	attributeValue, err := sv.repo.GetAttributeValueByID(c, repository.GetAttributeValueByIDParams{
		AttributeValueID: attributeParam.ValueID,
		AttributeID:      attributeParam.ID,
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(err))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	err = sv.repo.DeleteAttributeValue(c, attributeValue.AttributeValueID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	msg := "Attribute value removed successfully"

	c.JSON(http.StatusNoContent, GenericResponse[bool]{nil, &msg, nil})
}

// @Summary Update attribute value
// @Description Update attribute value
// @Tags attributes
// @Accept json
// @Produce json
// @Param attribute_id path int true "Attribute ID"
// @Param id path int true "Attribute value ID"
// @Param value body AttributeValue true "Attribute value"
// @Success 200 {object} GenericResponse[AttributeValue]
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /attributes/{attribute_id}/values/{id} [put]
func (sv *Server) updateAttributeValue(c *gin.Context) {
	var attributeParam AttributeValueParams
	if err := c.ShouldBindUri(&attributeParam); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	var req AttributeValue
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	attributeValue, err := sv.repo.GetAttributeValueByID(c, repository.GetAttributeValueByIDParams{
		AttributeValueID: attributeParam.ValueID,
		AttributeID:      attributeParam.ID,
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(err))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if attributeValue.AttributeValue == req.Value {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("attribute value is the same as the existing one")))
		return
	}

	updateParams := repository.UpdateAttributeValueParams{
		AttributeValueID: attributeValue.AttributeValueID,
		AttributeValue:   req.Value,
	}
	if req.Color != nil {
		updateParams.Color = util.GetPgTypeText(*req.Color)
	}

	attributeValue, err = sv.repo.UpdateAttributeValue(c, updateParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	msg := "Attribute value updated successfully"
	c.JSON(http.StatusOK, GenericResponse[repository.AttributeValue]{&attributeValue, &msg, nil})
}
