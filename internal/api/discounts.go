package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/constants"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
	"github.com/thanhphuocnguyen/go-eshop/internal/processors"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
)

// getAvailableDiscounts godoc
// @Summary Get available discounts
// @Description Get a list of available discounts
// @Tags discounts
// @Accept  json
// @Produce  json
// @Success 200 {object} ApiResponse[[]DiscountListItemResponseModel]
// @Failure 500 {object} ErrorResp
// @Router /discounts/available [get]
func (sv *Server) getAvailableDiscounts(c *gin.Context) {
	authPayload := c.MustGet(constants.AuthPayLoad).(*auth.TokenPayload)
	// Get available discounts
	discountRows, err := sv.repo.GetAvailableDiscountsForUser(c, authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	discounts := make([]dto.DiscountListItem, len(discountRows))
	for i, discount := range discountRows {
		discountValue, _ := discount.DiscountValue.Float64Value()
		discounts[i] = dto.DiscountListItem{
			ID:            discount.ID.String(),
			Code:          discount.Code,
			DiscountType:  string(discount.DiscountType),
			DiscountValue: discountValue.Float64,
			IsActive:      discount.IsActive,
			TimeUsed:      discount.TimesUsed,
			UsageLimit:    discount.UsageLimit,
			ValidFrom:     discount.ValidFrom.String(),
			CreatedAt:     discount.CreatedAt.String(),
		}
		if discount.ValidUntil.Valid {
			discounts[i].ValidUntil = discount.ValidUntil.Time.String()
		}
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, discounts, nil, nil))
}

// CheckDiscountApplicability godoc
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
func (sv *Server) checkDiscountsApplicability(c *gin.Context) {
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

// getDiscountByID godoc
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
func (sv *Server) getDiscountByID(c *gin.Context) {
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

func (sv *Server) addDiscountRoutes(rg *gin.RouterGroup) {
	discountsGroup := rg.Group("discounts", authenticateMiddleware(sv.tokenGenerator))
	{
		discountsGroup.GET("/available", sv.getAvailableDiscounts)
		discountsGroup.POST("/check-applicability", sv.checkDiscountsApplicability)
		discountsGroup.GET("/:id", sv.getDiscountByID)
	}
}
