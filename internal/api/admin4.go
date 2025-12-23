package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

// @Summary Update a product variant
// @Schemes http
// @Description update a product with the input payload
// @Tags products
// @Accept json
// @Param input body models.UpdateProdVariantModel true "Product variant input"
// @Produce json
// @Success 200 {object} dto.ApiResponse[repository.ProductVariant]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/products/{id}/variants/{variantId} [put]
func (s *Server) adminUpdateVariant(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id := chi.URLParam(r, "id")
	variantId := chi.URLParam(r, "variantId")
	if id == "" || variantId == "" {
		RespondBadRequest(w, InvalidBodyCode, errors.New("id and variantId parameters are required"))
		return
	}

	var req models.UpdateProdVariantModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	prod, err := s.repo.GetProductByID(c, repository.GetProductByIDParams{ID: uuid.MustParse(id)})
	if err != nil {
		RespondNotFound(w, NotFoundCode, err)
		return
	}

	updateParams := repository.UpdateProductVariantParams{
		ProductID: prod.ID,
		ID:        uuid.MustParse(variantId),
	}

	if req.Price != nil {
		updateParams.Price = utils.GetPgNumericFromFloat(*req.Price)
	}
	if req.StockQty != nil {
		updateParams.Stock = req.StockQty
	}
	if req.Weight != nil {
		updateParams.Weight = utils.GetPgNumericFromFloat(*req.Weight)
	}
	if req.Description != nil {
		updateParams.Description = req.Description
	}

	updatedVariant, err := s.repo.UpdateProductVariant(c, updateParams)
	if err != nil {
		log.Error().Err(err).Timestamp().Msg("UpdateProductVariant")
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondSuccess(w, updatedVariant)
}

// @Summary Upload a product variant image
// @Schemes http
// @Description upload a product variant image with the input payload
// @Tags products
// @Accept multipart/form-data
// @Param id path string true "Product ID"
// @Param variantId path string true "Product Variant ID"
// @Produce json
// @Success 200 {object} dto.ApiResponse[repository.ProductVariant]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/products/{id}/variants/{variantId}/images [post]
func (s *Server) adminUploadVariantImage(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id := chi.URLParam(r, "id")
	variantId := chi.URLParam(r, "variantId")
	if id == "" || variantId == "" {
		RespondBadRequest(w, InvalidBodyCode, errors.New("id and variantId parameters are required"))
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10MB max
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	defer file.Close()

	prod, err := s.repo.GetProductByID(c, repository.GetProductByIDParams{ID: uuid.MustParse(id)})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dto.CreateErr(NotFoundCode, err))
		return
	}

	variant, err := s.repo.GetProductVariantByID(c, repository.GetProductVariantByIDParams{
		ID:        uuid.MustParse(variantId),
		ProductID: prod.ID,
	})

	if err != nil {
		RespondNotFound(w, NotFoundCode, err)
		return
	}

	if variant.ImageID != nil {
		msg, err := s.uploadService.Remove(c, *variant.ImageID)
		if err != nil {
			log.Error().Err(err).Timestamp().Msg(msg)
			RespondInternalServerError(w, InternalServerErrorCode, err)
			return
		}
	}

	fileHeader := &struct {
		Filename string
		Header   map[string][]string
		Size     int64
	}{
		Filename: header.Filename,
		Header:   header.Header,
		Size:     header.Size,
	}

	uploadID, url, err := s.uploadService.Upload(c, fileHeader)
	if err != nil {
		log.Error().Err(err).Timestamp().Msg("UploadFile")
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	updateParam := repository.UpdateProductVariantParams{
		ProductID: prod.ID,
		ID:        variant.ID,
		ImageID:   &uploadID,
		ImageUrl:  &url,
	}

	updatedVariant, err := s.repo.UpdateProductVariant(c, updateParam)
	if err != nil {
		log.Error().Err(err).Timestamp().Msg("UpdateProductVariant")
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondSuccess(w, updatedVariant)
}

// @Summary Delete a product variant
// @Schemes http
// @Description delete a product variant with the input payload
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} dto.ApiResponse[string]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/products/{id}/variant/{variantID} [delete]
func (s *Server) adminDeleteVariant(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id := chi.URLParam(r, "id")
	variantId := chi.URLParam(r, "variantId")
	if id == "" || variantId == "" {
		RespondBadRequest(w, InvalidBodyCode, errors.New("id and variantId parameters are required"))
		return
	}

	prod, err := s.repo.GetProductByID(c, repository.GetProductByIDParams{ID: uuid.MustParse(id)})
	if err != nil {
		RespondNotFound(w, NotFoundCode, err)
		return
	}

	err = s.repo.DeleteProductVariant(c, repository.DeleteProductVariantParams{
		ProductID: prod.ID,
		ID:        uuid.MustParse(variantId),
	})

	if err != nil {
		log.Error().Err(err).Timestamp().Msg("DeleteVariant")
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	RespondNoContent(w)
}

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
func (s *Server) adminCreateAttribute(w http.ResponseWriter, r *http.Request) {
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
func (s *Server) adminUpdateAttribute(w http.ResponseWriter, r *http.Request) {
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
func (s *Server) adminRemoveAttribute(w http.ResponseWriter, r *http.Request) {
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

	attrs, err := s.repo.GetProductAttributeValuesByProductID(c, uuid.MustParse(idParam))
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

// @Summary Get all orders (admin endpoint)
// @Description Get all orders with pagination and filtering
// @Tags admin
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Param status query string false "Filter by status"
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse[[]dto.OrderListItem]
// @Failure 401 {object} ErrorResp
// @Failure 403 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/orders [get]
func (s *Server) adminGetOrders(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	var orderListQuery models.PaginationQuery = GetPaginationQuery(r)
	status := r.URL.Query().Get("status")
	paymentStatus := r.URL.Query().Get("paymentStatus")

	dbParams := repository.GetOrdersParams{
		Limit:  orderListQuery.PageSize,
		Offset: (orderListQuery.Page - 1) * orderListQuery.PageSize,
	}
	countParams := repository.CountOrdersParams{}

	if status != "" {
		dbParams.Status = repository.NullOrderStatus{
			OrderStatus: repository.OrderStatus(status),
			Valid:       true,
		}
		countParams.Status = repository.NullOrderStatus{
			OrderStatus: repository.OrderStatus(status),
			Valid:       true,
		}
	}

	if paymentStatus != "" {
		dbParams.PaymentStatus = repository.NullPaymentStatus{
			PaymentStatus: repository.PaymentStatus(paymentStatus),
			Valid:         true,
		}
		countParams.PaymentStatus = repository.NullPaymentStatus{
			PaymentStatus: repository.PaymentStatus(paymentStatus),
			Valid:         true,
		}
	}

	fetchedOrderRows, err := s.repo.GetOrders(c, dbParams)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	count, err := s.repo.CountOrders(c, countParams)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	var orderResponses []dto.OrderListItem
	for _, aggregated := range fetchedOrderRows {
		// Convert PaymentStatus interface{} to PaymentStatus type
		paymentStatus := repository.PaymentStatusPending
		if aggregated.PaymentStatus.Valid {
			paymentStatus = aggregated.PaymentStatus.PaymentStatus
		}

		total, _ := aggregated.TotalPrice.Float64Value()
		orderResponses = append(orderResponses, dto.OrderListItem{
			ID:            aggregated.ID,
			Total:         total.Float64,
			TotalItems:    int32(aggregated.TotalItems),
			Status:        aggregated.Status,
			CustomerName:  aggregated.CustomerName,
			CustomerEmail: aggregated.CustomerEmail,
			PaymentStatus: paymentStatus,
			CreatedAt:     aggregated.CreatedAt.UTC(),
			UpdatedAt:     aggregated.UpdatedAt.UTC(),
		})
	}

	RespondSuccessWithPagination(w, orderResponses, dto.CreatePagination(orderListQuery.Page, orderListQuery.PageSize, count))
}

// @Summary Get order details by ID (admin endpoint)
// @Description Get detailed information about an order by its ID
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse[dto.OrderDetail]
// @Failure 401 {object} ErrorResp
// @Failure 403 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/orders/{id} [get]
func (s *Server) adminGetOrderDetail(w http.ResponseWriter, r *http.Request) {
	// Reuse the existing order detail handler since admin has access to all orders
	s.getOrderDetail(w, r)
}

// @Summary Change order status
// @Description Change order status by order ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param status body string true "Status"
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse[uuid.UUID]
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/orders/{orderId}/status [put]
func (s *Server) adminChangeOrderStatus(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id, err := GetUrlParam(r, "id")
	var req models.OrderStatusModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	order, err := s.repo.GetOrder(c, uuid.MustParse(id))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, err)
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	if order.Status == repository.OrderStatusDelivered || order.Status == repository.OrderStatusCancelled || order.Status == repository.OrderStatusRefunded {
		RespondBadRequest(w, InvalidPaymentCode, errors.New("order cannot be changed"))
		return
	}

	status := repository.OrderStatus(req.Status)

	updateParams := repository.UpdateOrderParams{
		ID: order.ID,
		Status: repository.NullOrderStatus{
			OrderStatus: status,
			Valid:       true,
		},
	}
	if status == repository.OrderStatusConfirmed {
		updateParams.ConfirmedAt = pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		}
	}
	if status == repository.OrderStatusDelivering {
		updateParams.DeliveredAt = pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		}
	}

	rs, err := s.repo.UpdateOrder(c, updateParams)

	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	if err := s.cacheSrv.Delete(c, "order_detail:"+id); err != nil {
		log.Err(err).Msg("failed to delete order detail cache")
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	s.cacheSrv.Delete(c, "order_detail:"+id)

	RespondSuccess(w, rs)
}
