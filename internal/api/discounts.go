package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

// createDiscountHandler godoc
// @Summary Create a new discount
// @Description Create a new discount
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param input body CreateDiscountRequest true "Discount info"
// @Success 201 {object} ApiResponse[DiscountDetailResponseModel]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /discounts [post]
func (sv *Server) createDiscountHandler(c *gin.Context) {
	// Create a new discount
	var req CreateDiscountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}
	sqlParams := repository.InsertDiscountParams{
		Code:          req.Code,
		DiscountType:  req.DiscountType,
		DiscountValue: utils.GetPgNumericFromFloat(req.DiscountValue),
		IsActive:      req.IsActive,
		UsageLimit:    req.UsageLimit,
		Description:   req.Description,
		StartsAt:      utils.GetPgTypeTimestamp(req.StartsAt),
		ExpiresAt:     utils.GetPgTypeTimestamp(req.ExpiresAt),
	}

	if req.MinPurchaseAmount != nil {
		sqlParams.MinPurchaseAmount = utils.GetPgNumericFromFloat(*req.MinPurchaseAmount)
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

// getDiscountsHandler godoc
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
func (sv *Server) getDiscountsHandler(c *gin.Context) {
	var queries DiscountListQuery
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
	listData := make([]DiscountListItemResponseModel, len(discounts))
	for i, discount := range discounts {
		discountValue, _ := discount.DiscountValue.Float64Value()

		listData[i] = DiscountListItemResponseModel{
			ID:            discount.ID.String(),
			Code:          discount.Code,
			DiscountType:  discount.DiscountType,
			DiscountValue: discountValue.Float64,
			IsActive:      discount.IsActive,
			UsedCount:     discount.UsedCount,
			UsageLimit:    discount.UsageLimit,
			Description:   discount.Description,
			StartsAt:      discount.StartsAt.String(),
			ExpiresAt:     discount.ExpiresAt.String(),
		}
		if discount.MinPurchaseAmount.Valid {
			minPurchaseAmount, _ := discount.MinPurchaseAmount.Float64Value()
			listData[i].MinPurchase = minPurchaseAmount.Float64
		}
		if discount.MaxDiscountAmount.Valid {
			maxDiscountAmount, _ := discount.MaxDiscountAmount.Float64Value()
			listData[i].MaxDiscount = maxDiscountAmount.Float64
		}
	}
	pagination := createPagination(queries.Page, queries.PageSize, total)

	c.JSON(http.StatusOK, createDataResp(c, listData, pagination, nil))
}

// getDiscountByIDHandler godoc
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
func (sv *Server) getDiscountByIDHandler(c *gin.Context) {
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

	discountUsages := make([]DiscountUsageHistory, len(discountUsageRows))
	for i, usage := range discountUsageRows {
		discountAmount, _ := usage.DiscountAmount.Float64Value()
		amount, _ := usage.TotalPrice.Float64Value()
		discountUsages[i] = DiscountUsageHistory{
			OrderID:        usage.OrderID.String(),
			ID:             discount.ID.String(),
			CustomerName:   usage.CustomerName,
			Amount:         amount.Float64,
			DiscountAmount: discountAmount.Float64,
			Date:           usage.CreatedAt.Time,
		}
	}

	discountValue, _ := discount.DiscountValue.Float64Value()

	resp := DiscountDetailResponseModel{
		ID:            discount.ID.String(),
		Code:          discount.Code,
		DiscountType:  discount.DiscountType,
		DiscountValue: discountValue.Float64,
		IsActive:      discount.IsActive,
		UsedCount:     discount.UsedCount,
		UsageLimit:    discount.UsageLimit,
		Description:   discount.Description,
		StartsAt:      discount.StartsAt.String(),
		ExpiresAt:     discount.ExpiresAt.String(),
		CreatedAt:     discount.CreatedAt.String(),
		UpdatedAt:     discount.UpdatedAt.String(),
		UsageHistory:  discountUsages,
	}

	if discount.MinPurchaseAmount.Valid {
		minPurchaseAmount, _ := discount.MinPurchaseAmount.Float64Value()
		resp.MinPurchase = minPurchaseAmount.Float64
	}

	if discount.MaxDiscountAmount.Valid {
		maxDiscountAmount, _ := discount.MaxDiscountAmount.Float64Value()
		resp.MaxDiscount = maxDiscountAmount.Float64
	}

	c.JSON(http.StatusOK, createDataResp(c, resp, nil, nil))
}

// getDiscountProductsByIDHandler godoc
// @Summary Get discount products by ID
// @Description Get discount products by ID
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param id path string true "Discount ID"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Success 200 {object} ApiResponse[DiscountLinkObject]
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /discounts/{id}/products [get]
func (sv *Server) getDiscountProductsByIDHandler(c *gin.Context) {
	var param UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}
	var queries PaginationQueryParams
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	// Get discount products by ID
	discountProductRows, err := sv.repo.GetDiscountProducts(c,
		repository.GetDiscountProductsParams{
			DiscountID: uuid.MustParse(param.ID),
			Limit:      queries.PageSize,
			Offset:     (queries.Page - 1) * queries.PageSize,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	resp := make([]DiscountLinkObject, len(discountProductRows))

	total, err := sv.repo.CountDiscountProducts(c, uuid.MustParse(param.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}
	for i, discountProduct := range discountProductRows {
		basePrice, _ := discountProduct.BasePrice.Float64Value()

		resp[i] = DiscountLinkObject{
			ID:    discountProduct.ProductID.String(),
			Name:  discountProduct.Name,
			Price: &basePrice.Float64,
		}
	}

	pagination := createPagination(queries.Page, queries.PageSize, total)

	c.JSON(http.StatusOK, createDataResp(c, resp, pagination, nil))
}

// getDiscountCategoriesByIDHandler godoc
// @Summary Get discount categories by ID
// @Description Get discount categories by ID
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param id path string true "Discount ID"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Success 200 {object} ApiResponse[DiscountLinkObject]
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /discounts/{id}/categories [get]
func (sv *Server) getDiscountCategoriesByIDHandler(c *gin.Context) {
	var param UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}
	var queries PaginationQueryParams
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	// Get discount categories by ID
	discountCategoryRows, err := sv.repo.GetDiscountCategories(c,
		repository.GetDiscountCategoriesParams{
			DiscountID: uuid.MustParse(param.ID),
			Limit:      queries.PageSize,
			Offset:     (queries.Page - 1) * queries.PageSize,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	resp := make([]DiscountLinkObject, len(discountCategoryRows))

	total, err := sv.repo.CountDiscountCategories(c, uuid.MustParse(param.ID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	for i, discountCategory := range discountCategoryRows {
		resp[i] = DiscountLinkObject{
			ID:   discountCategory.CategoryID.String(),
			Name: discountCategory.Name,
		}
	}

	pagination := createPagination(queries.Page, queries.PageSize, total)

	c.JSON(http.StatusOK, createDataResp(c, resp, pagination, nil))
}

// getDiscountUsersByIDHandler godoc
// @Summary Get discount users by ID
// @Description Get discount users by ID
// @Tags discounts
// @Accept  json
// @Produce  json
// @Param id path string true "Discount ID"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Success 200 {object} ApiResponse[DiscountLinkObject]
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /discounts/{id}/users [get]
func (sv *Server) getDiscountUsersByIDHandler(c *gin.Context) {
	var param UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}
	var queries PaginationQueryParams
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	// Get discount Users by ID
	discountUserRows, err := sv.repo.GetDiscountUsers(c,
		repository.GetDiscountUsersParams{
			DiscountID: uuid.MustParse(param.ID),
			Limit:      queries.PageSize,
			Offset:     (queries.Page - 1) * queries.PageSize,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	resp := make([]DiscountLinkObject, len(discountUserRows))

	total, err := sv.repo.CountDiscountUsers(c, uuid.MustParse(param.ID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	for i, discountUser := range discountUserRows {
		name := ""
		if discountUser.Fullname != nil {
			name = *discountUser.Fullname
		}
		resp[i] = DiscountLinkObject{
			ID:   discountUser.ID.String(),
			Name: name,
		}
	}

	pagination := createPagination(queries.Page, queries.PageSize, total)

	c.JSON(http.StatusOK, createDataResp(c, resp, pagination, nil))
}

func (sv *Server) updateDiscountHandler(c *gin.Context) {
	var param UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	var req repository.UpdateDiscountTxArgs
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}
	// Update discount
	err := sv.repo.UpdateDiscountTx(c, uuid.MustParse(param.ID), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}
	c.JSON(http.StatusOK, createDataResp(c, struct{}{}, nil, nil))
}

func (sv *Server) deleteDiscountHandler(c *gin.Context) {
	// Delete discount
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, nil))
		return
	}

	err := sv.repo.DeleteDiscount(c, uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.Status(http.StatusNoContent)
}
