package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/constants"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
	"github.com/thanhphuocnguyen/go-eshop/internal/processors"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
)

// CheckDiscountApplicabilityHandler godoc
// @Summary Check discount applicability
// @Description Check if a discount code is applicable to the current cart
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param input body CheckDiscountApplicabilityRequest true "Discount applicability info"
// @Success 200 {object} ApiResponse[CheckDiscountApplicabilityResponseModel]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /discounts/check-applicability [post]
func (sv *Server) CheckDiscountsApplicabilityHandler(c *gin.Context) {
	// Check discount applicability
	authPayload := c.MustGet(constants.AuthPayLoad).(*auth.TokenPayload)
	var req models.CheckDiscountApplicabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	user, err := sv.repo.GetUserDetailsByID(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, errors.New("user not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	itemRows, err := sv.repo.GetCartItems(c, uuid.MustParse(req.CartID))
	if err != nil {
		log.Error().Err(err).Msg("GetCartItems")
		return
	}

	discountResult, err := sv.discountProcessor.ProcessDiscounts(c, processors.DiscountContext{User: user, CartItems: itemRows}, req.DiscountCodes)
	if err != nil {
		log.Error().Err(err).Msg("ProcessDiscounts")
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, discountResult, nil, nil))
}

// CreateDiscountHandler godoc
// @Summary Create a new discount
// @Description Create a new discount
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param input body CreateDiscountRequest true "Discount info"
// @Success 201 {object} ApiResponse[DiscountDetailResponseModel]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/discounts [post]
func (sv *Server) CreateDiscountHandler(c *gin.Context) {
	// Create a new discount
	var req models.AddDiscountModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
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

	discount, err := sv.repo.InsertDiscount(c, sqlParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusCreated, dto.CreateDataResp(c, discount.String(), nil, nil))
}

// GetDiscountsHandler godoc
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
// @Success 200 {object} ApiResponse[[]DiscountListItemResponseModel]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/discounts [get]
func (sv *Server) GetDiscountsHandler(c *gin.Context) {
	var queries models.DiscountListQuery
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	// Get all discounts
	sqlParams := repository.GetDiscountsParams{
		Limit:  queries.PageSize,
		Offset: (queries.Page - 1) * queries.PageSize,
		// Search:       queries.Search,
		// DiscountType: queries.DiscountType,
		// IsActive:     queries.IsActive,
	}

	if queries.FromDate != nil {
		sqlParams.FromDate = utils.GetPgTypeTimestamp(*queries.FromDate)
	}
	if queries.ToDate != nil {
		sqlParams.ToDate = utils.GetPgTypeTimestamp(*queries.ToDate)
	}

	discounts, err := sv.repo.GetDiscounts(c, sqlParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	total, err := sv.repo.CountDiscounts(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
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

	c.JSON(http.StatusOK, dto.CreateDataResp(c, listData, pagination, nil))
}

// GetDiscountByIDHandler godoc
// @Summary Get discount by ID
// @Description Get discount by ID
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param id path string true "Discount ID"
// @Success 200 {object} ApiResponse[DiscountDetailResponseModel]
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /discounts/{id} [get]
func (sv *Server) GetDiscountByIDHandler(c *gin.Context) {
	// Get discount by ID
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, nil))
		return
	}

	discount, err := sv.repo.GetDiscountByID(c, uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	discountUsageRows, err := sv.repo.GetDiscountUsages(c, discount.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	discountUsages := make([]dto.DiscountUsageHistory, len(discountUsageRows))
	for i, usage := range discountUsageRows {
		discountAmount, _ := usage.DiscountAmount.Float64Value()
		amount, _ := usage.TotalPrice.Float64Value()
		discountUsages[i] = dto.DiscountUsageHistory{
			ID:             discount.ID.String(),
			CustomerName:   usage.CustomerName,
			Amount:         amount.Float64,
			DiscountAmount: discountAmount.Float64,
			Date:           usage.CreatedAt,
			OrderID:        usage.OrderID.String(),
		}
	}

	discountRuleRows, err := sv.repo.GetDiscountRules(c, discount.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	discountRules := make([]dto.DiscountRuleDetail, len(discountRuleRows))
	for i, rule := range discountRuleRows {
		ruleDetail, err := dto.MapToDiscountRuleDetail(rule)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
			return
		}
		discountRules[i] = ruleDetail
	}

	discountValue, _ := discount.DiscountValue.Float64Value()

	resp := dto.DiscountDetail{
		ID:            discount.ID.String(),
		Code:          discount.Code,
		DiscountType:  string(discount.DiscountType),
		DiscountValue: discountValue.Float64,
		IsActive:      discount.IsActive,
		TimesUsed:     discount.TimesUsed,
		UsageLimit:    discount.UsageLimit,
		Description:   discount.Description,
		ValidFrom:     discount.ValidFrom.String(),
		UsageHistory:  discountUsages,
		DiscountRules: discountRules,
		CreatedAt:     discount.CreatedAt.String(),
		UpdatedAt:     discount.UpdatedAt.String(),
	}
	if discount.ValidUntil.Valid {
		resp.ValidUntil = discount.ValidUntil.Time.String()
	}

	if discount.MinOrderValue.Valid {
		minPurchaseAmount, _ := discount.MinOrderValue.Float64Value()
		resp.MinPurchase = minPurchaseAmount.Float64
	}

	if discount.MaxDiscountAmount.Valid {
		maxDiscountAmount, _ := discount.MaxDiscountAmount.Float64Value()
		resp.MaxDiscount = maxDiscountAmount.Float64
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, resp, nil, nil))
}

// UpdateDiscountHandler godoc
// @Summary Update discount by ID
// @Description Update discount by ID
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param id path string true "Discount ID"
// @Param input body UpdateDiscountRequest true "Discount info"
// @Success 200 {object} ApiResponse[string]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /discounts/{id} [put]
func (sv *Server) UpdateDiscountHandler(c *gin.Context) {
	// Update discount by ID
	var param models.UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	var req models.UpdateDiscountModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	discount, err := sv.repo.GetDiscountByID(c, uuid.MustParse(param.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
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

	updated, err := sv.repo.UpdateDiscount(c, sqlParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, updated, nil, nil))
}

// AdminDeleteDiscountHandler godoc
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
func (sv *Server) AdminDeleteDiscountHandler(c *gin.Context) {
	// Delete discount by ID
	var param models.UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	err := sv.repo.DeleteDiscount(c, uuid.MustParse(param.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.Status(http.StatusNoContent)
}

// AdminAddDiscountRuleHandler godoc
// @Summary Add a discount rule to a discount
// @Description Add a discount rule to a discount
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param id path string true "Discount ID"
// @Param input body AddDiscountRuleRequest true "Discount rule info"
// @Success 201 {object} ApiResponse[string]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/discounts/{id}/rules [post]
func (sv *Server) AdminAddDiscountRuleHandler(c *gin.Context) {
	// Add a discount rule to a discount
	var param models.UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	var req models.AddDiscountRuleModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	var ruleVal []byte
	switch req.RuleType {
	case "first_time_buyer":
		var ruleValue models.FirstTimeBuyerRule
		if err := mapstructure.Decode(req.RuleValue, &ruleValue); err != nil {
			c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
			return
		}
		bs, err := json.Marshal(ruleValue)

		if err != nil {
			c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
			return
		}
		ruleVal = bs
	case "product":
		var ruleValue models.ProductRule
		if err := mapstructure.Decode(req.RuleValue, &ruleValue); err != nil {
			c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
			return
		}
		bs, err := json.Marshal(ruleValue)

		if err != nil {
			c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
			return
		}
		ruleVal = bs
	case "category":
		var ruleValue models.CategoryRule
		if err := mapstructure.Decode(req.RuleValue, &ruleValue); err != nil {
			c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
			return
		}
		bs, err := json.Marshal(ruleValue)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
			return
		}
		ruleVal = bs
	case "customer_segment":
		var ruleValue models.CustomerSegmentRule
		if err := mapstructure.Decode(req.RuleValue, &ruleValue); err != nil {
			c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
			return
		}
		bs, err := json.Marshal(ruleValue)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
			return
		}
		ruleVal = bs
	}

	sqlParams := repository.InsertDiscountRuleParams{
		DiscountID: uuid.MustParse(param.ID),
		RuleType:   req.RuleType,
		RuleValue:  ruleVal,
	}

	rule, err := sv.repo.InsertDiscountRule(c, sqlParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusCreated, dto.CreateDataResp(c, rule, nil, nil))
}

// GetDiscountRulesHandler godoc
// @Summary Get all discount rules for a discount
// @Description Get all discount rules for a specific discount
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param id path string true "Discount ID"
// @Success 200 {object} ApiResponse[[]DiscountRule]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/discounts/{id}/rules [get]
func (sv *Server) GetDiscountRulesHandler(c *gin.Context) {
	var param models.UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	rules, err := sv.repo.GetDiscountRules(c, uuid.MustParse(param.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	var ruleDetails []dto.DiscountRuleDetail
	for _, rule := range rules {
		ruleDetail, err := dto.MapToDiscountRuleDetail(rule)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
			return
		}
		ruleDetails = append(ruleDetails, ruleDetail)
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, ruleDetails, nil, nil))
}

// GetDiscountRuleByIDHandler godoc
// @Summary Get a specific discount rule by ID
// @Description Get a specific discount rule by ID
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param id path string true "Discount ID"
// @Param ruleId path string true "Rule ID"
// @Success 200 {object} ApiResponse[DiscountRule]
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/discounts/{id}/rules/{ruleId} [get]
func (sv *Server) GetDiscountRuleByIDHandler(c *gin.Context) {
	var param models.UriRuleIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	rule, err := sv.repo.GetDiscountRuleByID(c, uuid.MustParse(param.RuleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	ruleDetail, err := dto.MapToDiscountRuleDetail(rule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, ruleDetail, nil, nil))
}

// UpdateDiscountRuleHandler godoc
// @Summary Update a discount rule
// @Description Update a discount rule
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param id path string true "Discount ID"
// @Param ruleId path string true "Rule ID"
// @Param input body UpdateDiscountRuleModel true "Updated discount rule info"
// @Success 200 {object} ApiResponse[string]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/discounts/{id}/rules/{ruleId} [put]
func (sv *Server) UpdateDiscountRuleHandler(c *gin.Context) {
	var param models.UriRuleIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	var req models.UpdateDiscountRuleModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	sqlParams := repository.UpdateDiscountRuleParams{
		ID: uuid.MustParse(param.RuleID),
	}

	if req.RuleType != nil {
		sqlParams.RuleType = req.RuleType
	}
	if req.RuleValue != nil {
		ruleValueBytes, err := json.Marshal(req.RuleValue)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
			return
		}
		sqlParams.RuleValue = ruleValueBytes
	}

	rule, err := sv.repo.UpdateDiscountRule(c, sqlParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, rule, nil, nil))
}

// DeleteDiscountRuleHandler godoc
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
func (sv *Server) DeleteDiscountRuleHandler(c *gin.Context) {
	var param models.UriRuleIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	err := sv.repo.DeleteDiscountRule(c, uuid.MustParse(param.RuleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (sv *Server) addDiscountRoutes(rg *gin.RouterGroup) {
	discountsGroup := rg.Group("discounts")
	{
		discountsGroup.GET("", sv.GetDiscountsHandler)
		discountsGroup.GET("/:id", sv.GetDiscountByIDHandler)
	}
}
