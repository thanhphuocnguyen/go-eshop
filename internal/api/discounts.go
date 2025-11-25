package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

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
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
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
	}

	if req.MinOrderValue != nil {
		sqlParams.MinOrderValue = utils.GetPgNumericFromFloat(*req.MinOrderValue)
	}
	if req.MaxDiscountAmount != nil {
		sqlParams.MaxDiscountAmount = utils.GetPgNumericFromFloat(*req.MaxDiscountAmount)
	}

	discount, err := sv.repo.InsertDiscount(c, sqlParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusCreated, createDataResp(c, discount.String(), nil, nil))
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
// @Router /discounts [get]
func (sv *Server) GetDiscountsHandler(c *gin.Context) {
	var queries models.DiscountListQuery
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
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
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}
	total, err := sv.repo.CountDiscounts(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
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
	pagination := createPagination(queries.Page, queries.PageSize, total)

	c.JSON(http.StatusOK, createDataResp(c, listData, pagination, nil))
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
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, nil))
		return
	}

	discount, err := sv.repo.GetDiscountByID(c, uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	discountUsageRows, err := sv.repo.GetDiscountUsages(c, discount.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
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

	c.JSON(http.StatusOK, createDataResp(c, resp, nil, nil))
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
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	var req models.UpdateDiscountModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}
	discount, err := sv.repo.GetDiscountByID(c, uuid.MustParse(param.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	sqlParams := repository.UpdateDiscountParams{
		ID:          discount.ID,
		Name:        req.Name,
		Code:        req.Code,
		IsActive:    req.IsActive,
		UsageLimit:  req.UsageLimit,
		IsStackable: req.IsStackable,
		Priority:    req.Priority,
		Description: req.Description,
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
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusOK, createDataResp(c, updated, nil, nil))
}

// DeleteDiscountHandler godoc
// @Summary Delete discount by ID
// @Description Delete discount by ID
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param id path string true "Discount ID"
// @Success 204
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /discounts/{id} [delete]
func (sv *Server) DeleteDiscountHandler(c *gin.Context) {
	// Delete discount by ID
	var param models.UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	err := sv.repo.DeleteDiscount(c, uuid.MustParse(param.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.Status(http.StatusNoContent)
}

// AddDiscountRuleHandler godoc
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
func (sv *Server) AddDiscountRuleHandler(c *gin.Context) {
	// Add a discount rule to a discount
	var param models.UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	var req models.AddDiscountRuleModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	sqlParams := repository.InsertDiscountRuleParams{
		DiscountID: uuid.MustParse(param.ID),
		RuleType:   req.RuleType,
		RuleValue:  []byte(req.ConditionType),
	}

	rule, err := sv.repo.InsertDiscountRule(c, sqlParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusCreated, createDataResp(c, rule.String(), nil, nil))
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
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	rules, err := sv.repo.GetDiscountRules(c, uuid.MustParse(param.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusOK, createDataResp(c, rules, nil, nil))
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
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	rule, err := sv.repo.GetDiscountRuleByID(c, uuid.MustParse(param.RuleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusOK, createDataResp(c, rule, nil, nil))
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
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	var req models.UpdateDiscountRuleModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	sqlParams := repository.UpdateDiscountRuleParams{
		ID: uuid.MustParse(param.RuleID),
	}

	if req.RuleType != nil {
		sqlParams.RuleType = req.RuleType
	}
	if req.ConditionType != nil {
		sqlParams.RuleValue = []byte(*req.ConditionType)
	}

	ruleID, err := sv.repo.UpdateDiscountRule(c, sqlParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusOK, createDataResp(c, ruleID.String(), nil, nil))
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
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	err := sv.repo.DeleteDiscountRule(c, uuid.MustParse(param.RuleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.Status(http.StatusNoContent)
}
