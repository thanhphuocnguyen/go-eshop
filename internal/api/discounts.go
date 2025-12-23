package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
	"github.com/thanhphuocnguyen/go-eshop/internal/processors"
)

// getAvailableDiscounts godoc
// @Summary Get available discounts
// @Description Get a list of available discounts
// @Tags discounts
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.ApiResponse[[]dto.DiscountListItem]
// @Failure 500 {object} dto.ErrorResp
// @Router /discounts/available [get]
func (s *Server) getAvailableDiscounts(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	_, claims, err := jwtauth.FromContext(c)
	if err != nil {
		RespondInternalServerError(w, UnauthorizedCode, errors.New("authorization payload is not provided"))
		return
	}
	// Get available discounts
	userID := uuid.MustParse(claims["userId"].(string))

	discountRows, err := s.repo.GetAvailableDiscountsForUser(c, userID)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
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

	RespondSuccess(w, discounts)
}

// CheckDiscountApplicability godoc
// @Summary Check discount applicability
// @Descript dion Check if a discount code is applicable to the current cart
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param input body models.CheckDiscountApplicabilityRequest true "Discount applicability info"
// @Success 200 {object} dto.ApiResponse[processors.DiscountResult]
// @Failure 400 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /discounts/check-applicability [post]
func (s *Server) checkDiscountsApplicability(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	// Check discount applicability
	_, claims, err := jwtauth.FromContext(c)
	if err != nil {
		RespondInternalServerError(w, UnauthorizedCode, errors.New("authorization payload is not provided"))
		return
	}
	userID := uuid.MustParse(claims["userId"].(string))
	var req models.CheckDiscountApplicabilityRequest
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	validate := validator.New()
	if err := validate.Struct(&req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	user, err := s.repo.GetUserDetailsByID(c, userID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, errors.New("user not found"))
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	itemRows, err := s.repo.GetCartItems(c, uuid.MustParse(req.CartID))
	if err != nil {
		log.Error().Err(err).Msg("GetCartItems")
		return
	}

	discountResult, err := s.discountProcessor.ProcessDiscounts(c, processors.DiscountContext{User: user, CartItems: itemRows}, req.DiscountCodes)
	if err != nil {
		log.Error().Err(err).Msg("ProcessDiscounts")
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	RespondSuccess(w, discountResult)
}

// getDiscountByID godoc
// @Summary Get discount by ID
// @Description Get discount by ID
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param id path string true "Discount ID"
// @Success 200 {object} dto.ApiResponse[dto.DiscountDetail]
// @Failure 400 {object} dto.ErrorResp
// @Failure 404 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /discounts/{id} [get]
func (s *Server) getDiscountByID(w http.ResponseWriter, r *http.Request) {
	// Get discount by ID
	c := r.Context()
	id := chi.URLParam(r, "id")
	if id == "" {
		RespondBadRequest(w, InvalidBodyCode, errors.New("id parameter is required"))
		return
	}

	discount, err := s.repo.GetDiscountByID(c, uuid.MustParse(id))
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	discountUsageRows, err := s.repo.GetDiscountUsages(c, discount.ID)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
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

	discountRuleRows, err := s.repo.GetDiscountRules(c, discount.ID)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	discountRules := make([]dto.DiscountRuleDetail, len(discountRuleRows))
	for i, rule := range discountRuleRows {
		ruleDetail, err := dto.MapToDiscountRuleDetail(rule)
		if err != nil {
			RespondInternalServerError(w, InternalServerErrorCode, err)
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

	RespondSuccess(w, resp)
}

func (s *Server) addDiscountRoutes(r chi.Router) {
	r.Route("/discounts", func(r chi.Router) {
		r.Get("/available", s.getAvailableDiscounts)
		r.Post("/check-applicability", s.checkDiscountsApplicability)
		r.Get("/{id}", s.getDiscountByID)
	})
}
