package api

import (
	"errors"
	"net/http"
	"slices"
	"strconv"
	"strings"

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
// @Param params body models.AttributeModel true "Attribute name"
// @Success 201 {object} dto.ApiResponse[dto.AttributeDetail]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/attributes [post]
func (s *Server) createAttribute(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	var req models.AttributeModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	attribute, err := s.repo.CreateAttribute(c, req.Name)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	if len(req.Values) > 0 {
		params := make([]repository.CreateAttributeValuesParams, len(req.Values))
		for i, val := range req.Values {
			params[i] = repository.CreateAttributeValuesParams{
				AttributeID: attribute.ID,
				Value:       val,
			}
		}

		_, err = s.repo.CreateAttributeValues(c, params)
		if err != nil {
			RespondInternalServerError(w, InternalServerErrorCode, err)
			return
		}
	}

	attributeResp := dto.AttributeDetail{
		ID:   attribute.ID,
		Name: attribute.Name,
	}

	RespondSuccess(w, attributeResp)
}

// @Summary Get all attributes
// @Description Get all attributes
// @Tags attributes
// @Accept json
// @Produce json
// @Success 200 {object} dto.ApiResponse[[]dto.AttributeDetail]
// @Failure 500 {object} ErrorResp
// @Router /admin/attributes [get]
func (s *Server) adminGetAttributes(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	query := r.URL.Query().Get("ids")

	ids := []int32{}
	if query != "" {
		for _, idStr := range strings.Split(query, ",") {
			id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 32)
			if err != nil {
				RespondBadRequest(w, InvalidBodyCode, err)
				return
			}
			ids = append(ids, int32(id))
		}
	}
	attributeRows, err := s.repo.GetAttributes(c, ids)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
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

	RespondSuccess(w, attributeResp)
}

// @Summary Get an attribute
// @Description Get an attribute
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Success 200 {object} dto.ApiResponse[dto.AttributeDetail]
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/attributes/{id} [get]
func (s *Server) adminGetAttributeByID(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	idParam, err := GetUrlParam(r, "id")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	id, _ := strconv.Atoi(idParam)

	attr, err := s.repo.GetAttributeByID(c, int32(id))
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	attributeResp := dto.AttributeDetail{
		Name: attr.Name,
		ID:   attr.ID,
	}

	values, err := s.repo.GetAttributeValues(c, attr.ID)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	attributeResp.Values = make([]dto.AttributeValueDetail, len(values))

	for i, val := range values {
		attributeResp.Values[i] = dto.AttributeValueDetail{
			ID:    val.ID,
			Value: val.Value,
		}
	}

	RespondSuccess(w, attributeResp)
}

// @Summary Update an attribute
// @Description Update an attribute
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Param params body models.AttributeModel true "Attribute name"
// @Success 200 {object} dto.ApiResponse[repository.Attribute]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/attributes/{id} [put]
func (s *Server) updateAttribute(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	idParam, err := GetUrlParam(r, "id")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	id, _ := strconv.Atoi(idParam)

	var req models.AttributeModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	attr, err := s.repo.GetAttributeByID(c, int32(id))
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	attribute, err := s.repo.UpdateAttribute(c, repository.UpdateAttributeParams{
		ID:   attr.ID,
		Name: req.Name,
	})

	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondSuccess(w, attribute)
}

// @Summary Delete an attribute
// @Description Delete an attribute
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Success 204 {object} nil
// @Failure 500 {object} ErrorResp
// @Router /admin/attributes/{id} [delete]
func (s *Server) removeAttribute(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	idParam, err := GetUrlParam(r, "id")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	id, _ := strconv.Atoi(idParam)

	attr, err := s.repo.GetAttributeByID(c, int32(id))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, err)
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	err = s.repo.DeleteAttribute(c, attr.ID)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondNoContent(w)
}

// @Summary Get attributes and their values by for a product
// @Description Get attributes and their values for a product
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.ApiResponse[[]dto.AttributeDetail]
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/attributes/product/{id} [get]
func (s *Server) adminGetAttributeValuesForProduct(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	idParam, err := GetUrlParam(r, "id")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	id, err := uuid.Parse(idParam)
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	attrs, err := s.repo.GetProductAttributeValuesByProductID(c, id)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
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

	RespondSuccess(w, resp)
}

// @Summary Add new attribute value
// @Description Add new attribute value
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Param params body models.AttributeValueModel true "Attribute value"
// @Success 200 {object} dto.ApiResponse[bool]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/attributes/{id}/create [post]
func (s *Server) adminAddAttributeValue(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	idParam, err := GetUrlParam(r, "id")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	var req models.AttributeValueModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	id, _ := strconv.Atoi(idParam)

	attr, err := s.repo.GetAttributeByID(c, int32(id))
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	obj, err := s.repo.CreateAttributeValue(c, repository.CreateAttributeValueParams{
		AttributeID: attr.ID,
		Value:       req.Value,
	})
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondCreated(w, obj)
}

// @Summary update attribute value
// @Description update attribute value
// @Tags attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Param params body models.AttributeValueModel true "Attribute value"
// @Success 200 {object} dto.ApiResponse[bool]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/attributes/{id}/update/{valueId} [put]
func (s *Server) adminUpdateAttrValue(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	idParam, err := GetUrlParam(r, "id")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	valueIdParam, err := GetUrlParam(r, "valueId")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	if valueIdParam == "" {
		RespondBadRequest(w, InvalidBodyCode, nil)
		return
	}
	var req models.AttributeValueModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	id, _ := strconv.Atoi(idParam)
	valueId, _ := strconv.Atoi(valueIdParam)

	attr, err := s.repo.GetAttributeByID(c, int32(id))
	if err != nil {
		RespondNotFound(w, NotFoundCode, err)
		return
	}

	res, err := s.repo.UpdateAttributeValue(c, repository.UpdateAttributeValueParams{
		AttributeID: attr.ID,
		ID:          int64(valueId),
		Value:       req.Value,
	})
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondSuccess(w, res)
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
// @Router /admin/attributes/{id}/remove/{valueId} [delete]
func (s *Server) adminRemoveAttrValue(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	idParam, err := GetUrlParam(r, "id")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	valueIdParam, err := GetUrlParam(r, "valueId")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	if valueIdParam == "" {
		RespondBadRequest(w, InvalidBodyCode, nil)
		return
	}

	id, _ := strconv.Atoi(idParam)
	valueId, _ := strconv.Atoi(valueIdParam)

	attr, err := s.repo.GetAttributeByID(c, int32(id))
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	err = s.repo.DeleteAttributeValueByValueID(c, repository.DeleteAttributeValueByValueIDParams{
		AttributeID: attr.ID,
		ID:          int64(valueId),
	})

	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondNoContent(w)
}
