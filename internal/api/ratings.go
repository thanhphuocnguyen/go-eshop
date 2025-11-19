package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
)

// @Summary Post a rating
// @Description Post a product rating
// @Tags ratings
// @Accept json
// @Produce json
// @Param orderItemId formData string true "Order Item ID"
// @Param rating formData float64 true "Rating (1-5)"
// @Param title formData string true "Review Title"
// @Param content formData string true "Review Content"
// @Param files formData file false "Images"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[ProductRatingModel]
// @Failure 400 {object} ErrorResp
// @Failure 403 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /ratings [post]
func (s *Server) postRatingHandler(c *gin.Context) {
	auth, _ := c.MustGet(AuthPayLoad).(*auth.Payload)

	var req PostRatingFormData
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, createErr(InvalidBodyCode, err))
		return
	}

	orderItem, err := s.repo.GetOrderItemByID(c, uuid.MustParse(req.OrderItemID))
	if err != nil {
		c.JSON(400, createErr(InvalidBodyCode, err))
		return
	}
	if orderItem.CustomerID != auth.UserID {
		c.JSON(403, createErr(InvalidBodyCode, nil))
		return
	}

	rating, err := s.repo.InsertProductRating(c, repository.InsertProductRatingParams{
		ProductID:        orderItem.ProductID,
		UserID:           auth.UserID,
		OrderItemID:      utils.GetPgTypeUUID(orderItem.OrderItemID),
		Rating:           utils.GetPgNumericFromFloat(req.Rating),
		ReviewTitle:      &req.Title,
		ReviewContent:    &req.Content,
		VerifiedPurchase: true,
	})

	if err != nil {
		c.JSON(500, createErr(InternalServerErrorCode, err))
		return
	}

	resp := ProductRatingModel{
		ID:               rating.ID,
		UserID:           rating.UserID,
		Rating:           req.Rating,
		ReviewTitle:      *rating.ReviewTitle,
		ReviewContent:    *rating.ReviewContent,
		VerifiedPurchase: rating.VerifiedPurchase,
		// Images:           images,
	}

	c.JSON(200, createDataResp(c, resp, nil, nil))
}

// @Summary Post a helpful rating
// @Description Post a helpful rating
// @Tags ratings
// @Accept json
// @Produce json
// @Param ratingId body string true "Rating ID"
// @Param helpful body bool true "Helpful"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[string]
// @Failure 400 {object} ErrorResp
// @Failure 403 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /ratings/{ratingId}/helpful [post]
func (s *Server) postRatingHelpfulHandler(c *gin.Context) {
	var param UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(400, createErr(InvalidBodyCode, err))
		return
	}
	auth, _ := c.MustGet(AuthPayLoad).(*auth.Payload)
	var req PostHelpfulRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, createErr(InvalidBodyCode, err))
		return
	}

	rating, err := s.repo.GetProductRating(c, uuid.MustParse(param.ID))
	if err != nil {
		c.JSON(400, createErr(InvalidBodyCode, err))
		return
	}
	if rating.UserID == auth.UserID {
		c.JSON(403, createErr(InvalidBodyCode, nil))
		return
	}

	id, err := s.repo.VoteHelpfulRatingTx(c, repository.VoteHelpfulRatingTxArgs{
		UserID:   auth.UserID,
		RatingID: rating.ID,
		Helpful:  req.Helpful,
	})

	if err != nil {
		c.JSON(500, createErr(InternalServerErrorCode, err))
	}

	c.JSON(200, createDataResp(c, id, nil, nil))
}

// @Summary Post a reply to a rating
// @Description Post a reply to a product rating
// @Tags ratings
// @Accept json
// @Produce json
// @Param ratingId path string true "Rating ID"
// @Param content body string true "Reply Content"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[string]
// @Failure 400 {object} ErrorResp
// @Failure 403 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /ratings/{ratingId}/reply [post]
func (s *Server) postReplyRatingHandler(c *gin.Context) {
	var param UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(400, createErr(InvalidBodyCode, err))
		return
	}
	auth, _ := c.MustGet(AuthPayLoad).(*auth.Payload)

	var req PostReplyRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, createErr(InvalidBodyCode, err))
		return
	}
	rating, err := s.repo.GetProductRating(c, uuid.MustParse(param.ID))
	if err != nil {
		c.JSON(400, createErr(InvalidBodyCode, err))
		return
	}

	reply, err := s.repo.InsertRatingReply(c, repository.InsertRatingReplyParams{
		RatingID: rating.ID,
		ReplyBy:  auth.UserID,
		Content:  req.Content,
	})
	if err != nil {
		c.JSON(500, createErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(200, createDataResp(c, reply.ID, nil, nil))
}

// @Summary Get product ratings
// @Description Get ratings for a specific product
// @Tags ratings
// @Accept json
// @Produce json
// @Param productId path string true "Product ID"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Success 200 {object} ApiResponse[[]ProductRatingModel]
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/ratings [get]
func (s *Server) getRatingsHandler(c *gin.Context) {
	var queries RatingsQueryParams
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(400, createErr(InvalidBodyCode, err))
		return
	}
	sqlParams := repository.GetProductRatingsParams{
		Limit:  queries.PageSize,
		Offset: (queries.Page - 1) * queries.PageSize,
	}

	if queries.Status != nil {
		switch *queries.Status {
		case "approved":
			sqlParams.IsApproved = utils.BoolPtr(true)
		case "rejected":
			sqlParams.IsApproved = utils.BoolPtr(false)
			sqlParams.IsVisible = utils.BoolPtr(false)
		case "pending":
			sqlParams.IsApproved = nil
		default:
		}
	}
	ratings, err := s.repo.GetProductRatings(c, sqlParams)
	if err != nil {
		c.JSON(400, createErr(InvalidBodyCode, err))
		return
	}

	ratingsCount, err := s.repo.CountProductRatings(c, pgtype.UUID{
		Bytes: uuid.Nil,
		Valid: false,
	})
	if err != nil {
		c.JSON(400, createErr(InvalidBodyCode, err))
		return
	}

	productRatings := make([]ProductRatingModel, 0)
	for _, rating := range ratings {
		ratingPoint, _ := rating.Rating.Float64Value()
		prIdx := -1
		for i, pr := range productRatings {
			if pr.ID == rating.ID {
				prIdx = i
				break
			}
		}
		if prIdx != -1 && rating.ImageID != nil {
			productRatings[prIdx].Images = append(productRatings[prIdx].Images, RatingImageModel{
				ID:  *rating.ImageID,
				URL: *rating.ImageUrl,
			})
			continue
		}
		model := ProductRatingModel{
			ID:               rating.ID,
			UserID:           rating.UserID,
			FirstName:        rating.FirstName,
			LastName:         rating.LastName,
			ProductName:      rating.ProductName,
			Rating:           ratingPoint.Float64,
			IsVisible:        rating.IsVisible,
			IsApproved:       rating.IsApproved,
			ReviewTitle:      *rating.ReviewTitle,
			ReviewContent:    *rating.ReviewContent,
			VerifiedPurchase: rating.VerifiedPurchase,
			Count:            ratingsCount,
		}
		if rating.ImageID != nil {
			model.Images = append(model.Images, RatingImageModel{
				ID:  *rating.ImageID,
				URL: *rating.ImageUrl,
			})
		}
		productRatings = append(productRatings, model)
	}
	c.JSON(200, createDataResp(c, productRatings, createPagination(queries.Page, queries.PageSize, ratingsCount), nil))
}

// @Summary Get product ratings
// @Description Get ratings for a specific product
// @Tags ratings
// @Accept json
// @Produce json
// @Param productId path string true "Product ID"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Success 200 {object} ApiResponse[[]ProductRatingModel]
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /ratings/products/{productId} [get]
func (s *Server) getRatingsByProductHandler(c *gin.Context) {
	var param UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(400, createErr(InvalidBodyCode, err))
		return
	}
	var queries PaginationQueryParams
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(400, createErr(InvalidBodyCode, err))
		return
	}

	ratings, err := s.repo.GetProductRatings(c, repository.GetProductRatingsParams{
		ProductID: utils.GetPgTypeUUID(uuid.MustParse(param.ID)),
		Limit:     queries.PageSize,
		Offset:    (queries.Page - 1) * queries.PageSize,
	})
	if err != nil {
		c.JSON(400, createErr(InvalidBodyCode, err))
		return
	}

	ratingsCount, err := s.repo.CountProductRatings(c, utils.GetPgTypeUUID(uuid.MustParse(param.ID)))
	if err != nil {
		c.JSON(400, createErr(InvalidBodyCode, err))
		return
	}
	productRatings := make([]ProductRatingModel, 0)
	for _, rating := range ratings {
		ratingPoint, _ := rating.Rating.Float64Value()
		prIdx := -1
		for i, pr := range productRatings {
			if pr.ID == rating.ID {
				prIdx = i
				break
			}
		}
		if prIdx != -1 && rating.ImageID != nil {
			productRatings[prIdx].Images = append(productRatings[prIdx].Images, RatingImageModel{
				ID:  *rating.ImageID,
				URL: *rating.ImageUrl,
			})
			continue
		}
		model := ProductRatingModel{
			ID:               rating.ID,
			UserID:           rating.UserID,
			FirstName:        rating.FirstName,
			LastName:         rating.LastName,
			Rating:           ratingPoint.Float64,
			ProductName:      rating.ProductName,
			IsVisible:        rating.IsVisible,
			ReviewTitle:      *rating.ReviewTitle,
			ReviewContent:    *rating.ReviewContent,
			VerifiedPurchase: rating.VerifiedPurchase,
		}
		if rating.ImageID != nil {
			model.Images = append(model.Images, RatingImageModel{
				ID:  *rating.ImageID,
				URL: *rating.ImageUrl,
			})
		}
		productRatings = append(productRatings, model)
	}

	c.JSON(200, createDataResp(c, productRatings, createPagination(queries.Page, queries.PageSize, ratingsCount), nil))
}

// @Summary Get order ratings
// @Description Get ratings for a specific order
// @Tags ratings
// @Accept json
// @Produce json
// @Param orderId path string true "Order ID"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[[]ProductRatingModel]
// @Failure 400 {object} ErrorResp
// @Failure 403 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /ratings/orders/{orderId} [get]
func (s *Server) getOrderRatingsHandler(c *gin.Context) {
	auth, _ := c.MustGet(AuthPayLoad).(*auth.Payload)

	var param struct {
		OrderID string `uri:"orderId" binding:"required,uuid"`
	}
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(400, createErr(InvalidBodyCode, err))
		return
	}

	orderItems, err := s.repo.GetOrderItemsByOrderID(c, uuid.MustParse(param.OrderID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}
	if len(orderItems) == 0 {
		c.JSON(http.StatusNotFound, createErr(NotFoundCode, nil))
		return
	}
	if orderItems[0].CustomerID != auth.UserID {
		c.JSON(http.StatusForbidden, createErr(PermissionDeniedCode, nil))
		return
	}
	orderItemIds := make([]uuid.UUID, len(orderItems))
	for i, orderItem := range orderItems {
		orderItemIds[i] = orderItem.OrderItemID
	}
	ratings, err := s.repo.GetProductRatingsByOrderItemIDs(c, orderItemIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(200, createDataResp(c, ratings, nil, nil))
}

// @Summary Delete a rating
// @Description Delete a product rating by ID
// @Tags Admin, ratings
// @Accept json
// @Produce json
// @Param id path string true "Rating ID"
// @Security BearerAuth
// @Success 204 {object} nil
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/ratings/{id} [delete]
func (sv *Server) deleteRatingHandler(c *gin.Context) {
	// Parse the rating ID from the URL
	var param struct {
		ID string `uri:"id" binding:"required,uuid"`
	}

	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	// Convert the ID to UUID
	ratingID, err := uuid.Parse(param.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	// Check if rating exists first
	_, err = sv.repo.GetProductRating(c, ratingID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErr(NotFoundCode, err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	// Delete the rating
	err = sv.repo.DeleteProductRating(c, ratingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Approve a rating
// @Description Approve a product rating by ID
// @Tags Admin, ratings
// @Accept json
// @Produce json
// @Param id path string true "Rating ID"
// @Security BearerAuth
// @Success 204 {object} nil
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/ratings/{id}/approve [post]
func (sv *Server) approveRatingHandler(c *gin.Context) {
	// Parse the rating ID from the URL
	var param struct {
		ID string `uri:"id" binding:"required,uuid"`
	}

	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	// Convert the ID to UUID
	ratingID, err := uuid.Parse(param.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	// Check if rating exists first
	rating, err := sv.repo.GetProductRating(c, ratingID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErr(NotFoundCode, err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	// Set IsApproved to true
	isApproved := true

	// Update the rating
	_, err = sv.repo.UpdateProductRating(c, repository.UpdateProductRatingParams{
		ID:         rating.ID,
		IsApproved: &isApproved,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Ban a user from rating
// @Description Ban a user from rating by setting their rating to invisible
// @Tags Admin, ratings
// @Accept json
// @Produce json
// @Param id path string true "Rating ID"
// @Security BearerAuth
// @Success 204 {object} nil
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/ratings/{id}/ban [post]
func (sv *Server) banUserRatingHandler(c *gin.Context) {
	// Parse the rating ID from the URL
	var param struct {
		ID string `uri:"id" binding:"required,uuid"`
	}

	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	// Convert the ID to UUID
	ratingID, err := uuid.Parse(param.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	// Check if rating exists first
	rating, err := sv.repo.GetProductRating(c, ratingID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErr(NotFoundCode, err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	// Set IsVisible to false
	isVisible := false

	// Update the rating
	_, err = sv.repo.UpdateProductRating(c, repository.UpdateProductRatingParams{
		ID:        rating.ID,
		IsVisible: &isVisible,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.Status(http.StatusNoContent)
}
