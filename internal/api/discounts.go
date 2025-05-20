package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

func (sv *Server) createDiscountHandler(c *gin.Context) {
	// Create a new discount
	var req CreateDiscountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[gin.H](InvalidBodyCode, "", err))
		return
	}
	sqlParams := repository.InsertDiscountParams{
		Code:          req.Code,
		Description:   &req.Description,
		DiscountType:  req.DiscountType,
		DiscountValue: utils.GetPgNumericFromFloat(req.DiscountValue),
		IsActive:      &req.IsActive,
		UsageLimit:    req.UsageLimit,
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
		c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "Failed to create discount", err))
		return
	}

	c.JSON(http.StatusCreated, createSuccessResponse(c, discount.String(), "Discount created successfully", nil, nil))
}

func (sv *Server) getDiscountsHandler(c *gin.Context) {
	var queries DiscountListQuery
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[gin.H](InvalidBodyCode, "", err))
		return
	}

	// Get all discounts
	sqlParams := repository.GetDiscountsParams{
		Limit:        queries.PageSize,
		Offset:       (queries.Page - 1) * queries.PageSize,
		Search:       queries.Search,
		DiscountType: queries.DiscountType,
		IsActive:     queries.IsActive,
	}
	if queries.FromDate != nil {
		sqlParams.FromDate = utils.GetPgTypeTimestamp(*queries.FromDate)
	}
	if queries.ToDate != nil {
		sqlParams.ToDate = utils.GetPgTypeTimestamp(*queries.ToDate)
	}

	discounts, err := sv.repo.GetDiscounts(c, sqlParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "Failed to get discounts", err))
		return
	}
	total, err := sv.repo.CountDiscounts(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "Failed to count discounts", err))
		return
	}
	pagination := createPagination(queries.Page, queries.PageSize, total)

	c.JSON(http.StatusOK, createSuccessResponse(c, discounts, "Discounts retrieved successfully", pagination, nil))
}

func (sv *Server) getDiscountByIDHandler(c *gin.Context) {
	// Get discount by ID
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, createErrorResponse[gin.H](InvalidBodyCode, "Discount ID is required", nil))
		return
	}

	discount, err := sv.repo.GetDiscountByID(c, uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "Failed to get discount", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, discount, "Discount retrieved successfully", nil, nil))
}

func (sv *Server) updateDiscountHandler(c *gin.Context) {
	var param UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[gin.H](InvalidBodyCode, "", err))
		return
	}

	var req repository.UpdateDiscountTxArgs
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[gin.H](InvalidBodyCode, "", err))
		return
	}
	// Update discount
	err := sv.repo.UpdateDiscountTx(c, uuid.MustParse(param.ID), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "Failed to update discount", err))
		return
	}
	c.Status(http.StatusNoContent)
}

func (sv *Server) deleteDiscountHandler(c *gin.Context) {
	// Delete discount
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, createErrorResponse[gin.H](InvalidBodyCode, "Discount ID is required", nil))
		return
	}

	err := sv.repo.DeleteDiscount(c, uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "Failed to delete discount", err))
		return
	}

	c.Status(http.StatusNoContent)
}
