package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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
func (sv *Server) getAvailableDiscounts(w http.ResponseWriter, r *http.Request) {
	authPayload, ok := r.Context().Value("auth").(*auth.TokenPayload)
	if !ok {
		RespondInternalServerError(w, UnauthorizedCode, errors.New("authorization payload is not provided"))
		return
	}
	// Get available discounts
	discountRows, err := sv.repo.GetAvailableDiscountsForUser(r.Context(), authPayload.UserID)
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

	RespondSuccess(w, r, discounts)
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
func (sv *Server) checkDiscountsApplicability(w http.ResponseWriter, r *http.Request) {
	// Check discount applicability
	authPayload, ok := r.Context().Value("auth").(*auth.TokenPayload)
	if !ok {
		RespondInternalServerError(w, UnauthorizedCode, errors.New("authorization payload is not provided"))
		return
	}
	var req models.CheckDiscountApplicabilityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	validate := validator.New()
	if err := validate.Struct(&req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	user, err := sv.repo.GetUserDetailsByID(r.Context(), authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, errors.New("user not found"))
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	itemRows, err := sv.repo.GetCartItems(r.Context(), uuid.MustParse(req.CartID))
	if err != nil {
		log.Error().Err(err).Msg("GetCartItems")
		return
	}

	discountResult, err := sv.discountProcessor.ProcessDiscounts(r.Context(), processors.DiscountContext{User: user, CartItems: itemRows}, req.DiscountCodes)
	if err != nil {
		log.Error().Err(err).Msg("ProcessDiscounts")
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	RespondSuccess(w, r, discountResult)
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
func (sv *Server) getDiscountByID(w http.ResponseWriter, r *http.Request) {
	// Get discount by ID
	id := chi.URLParam(r, "id")
	if id == "" {
		RespondBadRequest(w, InvalidBodyCode, errors.New("id parameter is required"))
		return
	}

	discount, err := sv.repo.GetDiscountByID(r.Context(), uuid.MustParse(id))
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	discountUsageRows, err := sv.repo.GetDiscountUsages(r.Context(), discount.ID)
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

	discountRuleRows, err := sv.repo.GetDiscountRules(r.Context(), discount.ID)
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

	RespondSuccess(w, r, resp)
}

func (sv *Server) addDiscountRoutes(r chi.Router) {
	r.Route("/discounts", func(r chi.Router) {
		// Apply authentication middleware
		r.Use(func(h http.Handler) http.Handler {
			return authenticateMiddleware(h, sv.tokenGenerator)
		})

		r.Get("/available", sv.getAvailableDiscounts)
		r.Post("/check-applicability", sv.checkDiscountsApplicability)
		r.Get("/{id}", sv.getDiscountByID)
	})
}
