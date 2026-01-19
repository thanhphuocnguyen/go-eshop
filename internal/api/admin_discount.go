package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

// adminGetDiscounts godoc
// @Summary Get all discounts
// @Description Get all discounts
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Param search query string false "Search by code"
// @Param discountType query string false "Discount type" default(percentage)
// @Param isActive query bool false "Is active" default(true)
// @Success 200 {object} dto.ApiResponse[[]dto.DiscountListItem]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/discounts [get]
func (s *Server) adminGetDiscounts(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	var queries models.PaginationQuery = ParsePaginationQuery(r)
	discountType := r.URL.Query().Get("discountType")
	isActive := r.URL.Query().Get("isActive")
	fromDateQ := r.URL.Query().Get("fromDate")
	toDateQ := r.URL.Query().Get("toDate")

	// Get all discounts
	sqlParams := repository.GetDiscountsParams{
		Limit:  queries.PageSize,
		Offset: (queries.Page - 1) * queries.PageSize,
		Search: queries.Search,
	}
	var discountTypeEnum repository.NullDiscountType

	if err := discountTypeEnum.Scan(discountType); discountType != "" && err == nil {
		sqlParams.DiscountType = discountTypeEnum
	}

	if isActive != "" {
		if isActive == "true" || isActive == "false" {
			isActiveBool := isActive == "true"
			sqlParams.IsActive = &isActiveBool
		} else {
			RespondBadRequest(w, InvalidBodyCode, errors.New("invalid isActive value"))
			return
		}
	}

	fromDate, err := time.Parse("2006-01-02T15:04:05Z07:00", fromDateQ)
	if fromDateQ != "" && err == nil {
		sqlParams.FromDate = utils.GetPgTypeTimestamp(fromDate)
	}
	toDate, err := time.Parse("2006-01-02T15:04:05Z07:00", toDateQ)
	if toDateQ != "" && err == nil {
		sqlParams.ToDate = utils.GetPgTypeTimestamp(toDate)
	}

	discounts, err := s.repo.GetDiscounts(c, sqlParams)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	total, err := s.repo.CountDiscounts(c)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	listData := make([]dto.DiscountListItem, len(discounts))
	for i, discount := range discounts {
		discountValue, _ := discount.DiscountValue.Float64Value()

		listData[i] = dto.DiscountListItem{
			ID:            discount.ID.String(),
			Code:          discount.Code,
			DiscountType:  string(discount.DiscountType),
			DiscountValue: discountValue.Float64,
			IsActive:      discount.IsActive,
			TimeUsed:      discount.TimesUsed,
			UsageLimit:    discount.UsageLimit,
			Description:   discount.Description,
			ValidFrom:     discount.ValidFrom.String(),
			CreatedAt:     discount.CreatedAt.String(),
			UpdatedAt:     discount.UpdatedAt.String(),
		}
		if discount.ValidUntil.Valid {
			listData[i].ValidUntil = discount.ValidUntil.Time.String()
		}
		if discount.MinOrderValue.Valid {
			minPurchaseAmount, _ := discount.MinOrderValue.Float64Value()
			listData[i].MinPurchase = &minPurchaseAmount.Float64
		}

		if discount.MaxDiscountAmount.Valid {
			maxDiscountAmount, _ := discount.MaxDiscountAmount.Float64Value()
			listData[i].MaxDiscount = &maxDiscountAmount.Float64
		}
	}
	pagination := dto.CreatePagination(queries.Page, queries.PageSize, total)

	RespondSuccessWithPagination(w, listData, pagination)
}

// adminCreateDiscount godoc
// @Summary Create a new discount
// @Description Create a new discount
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param input body models.AddDiscount true "Discount info"
// @Success 201 {object} dto.ApiResponse[uuid.UUID]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/discounts [post]
func (s *Server) adminCreateDiscount(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	// Create a new discount
	var req models.AddDiscount
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	sqlParams := repository.InsertDiscountParams{
		Code:          req.Code,
		DiscountType:  repository.DiscountType(req.DiscountType),
		DiscountValue: utils.GetPgNumericFromFloat(req.DiscountValue),
		IsActive:      req.IsActive,
		UsageLimit:    req.UsageLimit,
		Description:   req.Description,
		ValidFrom:     utils.GetPgTypeTimestamp(req.ValidFrom),
		ValidUntil:    utils.GetPgTypeTimestamp(req.ValidUntil),
		Name:          req.Name,
		UsagePerUser:  req.UsagePerUser,
		IsStackable:   req.IsStackable,
		Priority:      req.Priority,
	}

	if req.MinOrderValue != nil {
		sqlParams.MinOrderValue = utils.GetPgNumericFromFloat(*req.MinOrderValue)
	}
	if req.MaxDiscountAmount != nil {
		sqlParams.MaxDiscountAmount = utils.GetPgNumericFromFloat(*req.MaxDiscountAmount)
	}

	discount, err := s.repo.InsertDiscount(c, sqlParams)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondSuccess(w, discount.String())
}

// adminUpdateDiscount godoc
// @Summary Update discount by ID
// @Description Update discount by ID
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param id path string true "Discount ID"
// @Param input body models.UpdateDiscount true "Discount info"
// @Success 200 {object} dto.ApiResponse[string]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /discounts/{id} [put]
func (s *Server) adminUpdateDiscount(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	// Update discount by ID
	var param models.UriIDParam
	if err := s.GetRequestBody(r, &param); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	var req models.UpdateDiscount
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	discount, err := s.repo.GetDiscountByID(c, uuid.MustParse(param.ID))
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	sqlParams := repository.UpdateDiscountParams{
		ID:           discount.ID,
		Name:         req.Name,
		Code:         req.Code,
		IsActive:     req.IsActive,
		UsageLimit:   req.UsageLimit,
		IsStackable:  req.IsStackable,
		Priority:     req.Priority,
		Description:  req.Description,
		UsagePerUser: req.UsagePerUser,
	}

	if req.DiscountType != nil {
		sqlParams.DiscountType.Scan(req.DiscountType)
	}
	if req.DiscountValue != nil {
		sqlParams.DiscountValue = utils.GetPgNumericFromFloat(*req.DiscountValue)
	}
	if req.ValidFrom != nil {
		sqlParams.ValidFrom = utils.GetPgTypeTimestamp(*req.ValidFrom)
	}
	if req.ValidUntil != nil {
		sqlParams.ValidUntil = utils.GetPgTypeTimestamp(*req.ValidUntil)
	}
	if req.MinOrderValue != nil {
		sqlParams.MinOrderValue = utils.GetPgNumericFromFloat(*req.MinOrderValue)
	}
	if req.MaxDiscountAmount != nil {
		sqlParams.MaxDiscountAmount = utils.GetPgNumericFromFloat(*req.MaxDiscountAmount)
	}

	updated, err := s.repo.UpdateDiscount(c, sqlParams)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondSuccess(w, updated)
}

// adminDeleteDiscount godoc
// @Summary Delete discount by ID
// @Description Delete discount by ID
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param id path string true "Discount ID"
// @Success 204
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/discounts/{id} [delete]
func (s *Server) adminDeleteDiscount(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	// Delete discount by ID
	var param models.UriIDParam
	if err := s.GetRequestBody(r, &param); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	err := s.repo.DeleteDiscount(c, uuid.MustParse(param.ID))
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondNoContent(w)
}

// adminAddDiscountRule godoc
// @Summary Add a discount rule to a discount
// @Description Add a discount rule to a discount
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param id path string true "Discount ID"
// @Param input body models.AddDiscountRule true "Discount rule info"
// @Success 201 {object} dto.ApiResponse[string]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/discounts/{id}/rules [post]
func (s *Server) adminAddDiscountRule(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	// Add a discount rule to a discount
	var param models.UriIDParam
	if err := s.GetRequestBody(r, &param); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	var req models.AddDiscountRule
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	var ruleVal []byte
	switch req.RuleType {
	case "first_time_buyer":
		var ruleValue models.FirstTimeBuyerRule
		if err := mapstructure.Decode(req.RuleValue, &ruleValue); err != nil {
			RespondBadRequest(w, InvalidBodyCode, err)
			return
		}
		bs, err := json.Marshal(ruleValue)

		if err != nil {
			RespondBadRequest(w, InvalidBodyCode, err)
			return
		}
		ruleVal = bs
	case "product":
		var ruleValue models.ProductRule
		if err := mapstructure.Decode(req.RuleValue, &ruleValue); err != nil {
			RespondBadRequest(w, InvalidBodyCode, err)
			return
		}
		bs, err := json.Marshal(ruleValue)

		if err != nil {
			RespondBadRequest(w, InvalidBodyCode, err)
			return
		}
		ruleVal = bs
	case "category":
		var ruleValue models.CategoryRule
		if err := mapstructure.Decode(req.RuleValue, &ruleValue); err != nil {
			RespondBadRequest(w, InvalidBodyCode, err)
			return
		}
		bs, err := json.Marshal(ruleValue)
		if err != nil {
			RespondBadRequest(w, InvalidBodyCode, err)
			return
		}
		ruleVal = bs
	case "customer_segment":
		var ruleValue models.CustomerSegmentRule
		if err := mapstructure.Decode(req.RuleValue, &ruleValue); err != nil {
			RespondBadRequest(w, InvalidBodyCode, err)
			return
		}
		bs, err := json.Marshal(ruleValue)
		if err != nil {
			RespondBadRequest(w, InvalidBodyCode, err)
			return
		}
		ruleVal = bs
	}

	sqlParams := repository.InsertDiscountRuleParams{
		DiscountID: uuid.MustParse(param.ID),
		RuleType:   req.RuleType,
		RuleValue:  ruleVal,
	}

	rule, err := s.repo.InsertDiscountRule(c, sqlParams)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondCreated(w, dto.CreateDataResp(rule, nil, nil))
}

// adminGetDiscountRules godoc
// @Summary Get all discount rules for a discount
// @Description Get all discount rules for a specific discount
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param id path string true "Discount ID"
// @Success 200 {object} dto.ApiResponse[[]dto.DiscountRuleDetail]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/discounts/{id}/rules [get]
func (s *Server) adminGetDiscountRules(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id, err := GetUrlParam(r, "id")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	rules, err := s.repo.GetDiscountRules(c, uuid.MustParse(id))
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	var ruleDetails []dto.DiscountRuleDetail
	for _, rule := range rules {
		ruleDetail, err := dto.MapToDiscountRuleDetail(rule)
		if err != nil {
			RespondInternalServerError(w, InternalServerErrorCode, err)
			return
		}
		ruleDetails = append(ruleDetails, ruleDetail)
	}

	RespondSuccess(w, dto.CreateDataResp(ruleDetails, nil, nil))
}

// adminGetDiscountRuleByID godoc
// @Summary Get a specific discount rule by ID
// @Description Get a specific discount rule by ID
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param id path string true "Discount ID"
// @Param ruleId path string true "Rule ID"
// @Success 200 {object} dto.ApiResponse[dto.DiscountRuleDetail]
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/discounts/{id}/rules/{ruleId} [get]
func (s *Server) adminGetDiscountRuleByID(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	ruleId, err := GetUrlParam(r, "ruleId")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	rule, err := s.repo.GetDiscountRuleByID(c, uuid.MustParse(ruleId))
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	ruleDetail, err := dto.MapToDiscountRuleDetail(rule)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondSuccess(w, dto.CreateDataResp(ruleDetail, nil, nil))
}

// adminUpdateDiscountRule godoc
// @Summary Update a discount rule
// @Description Update a discount rule
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param id path string true "Discount ID"
// @Param ruleId path string true "Rule ID"
// @Param input body models.UpdateDiscountRule true "Updated discount rule info"
// @Success 200 {object} dto.ApiResponse[string]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/discounts/{id}/rules/{ruleId} [put]
func (s *Server) adminUpdateDiscountRule(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	ruleId, err := GetUrlParam(r, "ruleId")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	var req models.UpdateDiscountRule
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	sqlParams := repository.UpdateDiscountRuleParams{
		ID: uuid.MustParse(ruleId),
	}

	if req.RuleType != nil {
		sqlParams.RuleType = req.RuleType
	}
	if req.RuleValue != nil {
		ruleValueBytes, err := json.Marshal(req.RuleValue)
		if err != nil {
			RespondBadRequest(w, InvalidBodyCode, err)
			return
		}
		sqlParams.RuleValue = ruleValueBytes
	}

	rule, err := s.repo.UpdateDiscountRule(c, sqlParams)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondSuccess(w, dto.CreateDataResp(rule, nil, nil))
}

// adminDeleteDiscountRule godoc
// @Summary Delete a discount rule
// @Description Delete a discount rule
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param id path string true "Discount ID"
// @Param ruleId path string true "Rule ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/discounts/{id}/rules/{ruleId} [delete]
func (s *Server) adminDeleteDiscountRule(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	ruleId, err := GetUrlParam(r, "ruleId")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	err = s.repo.DeleteDiscountRule(c, uuid.MustParse(ruleId))
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondNoContent(w)
}
