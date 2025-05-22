package api

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"sync"

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
// @Param order_item_id formData string true "Order Item ID"
// @Param rating formData float64 true "Rating (1-5)"
// @Param title formData string true "Review Title"
// @Param content formData string true "Review Content"
// @Param files formData file false "Images"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[ProductRatingModel]
// @Failure 400 {object} ApiResponse[bool]
// @Failure 403 {object} ApiResponse[bool]
// @Failure 500 {object} ApiResponse[bool]
// @Router /ratings [post]
func (s *Server) postRatingHandler(c *gin.Context) {
	auth, _ := c.MustGet(authorizationPayload).(*auth.Payload)

	var req PostRatingFormData
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, createErrorResponse[gin.H](InvalidBodyCode, "invalid request", err))
		return
	}

	orderItem, err := s.repo.GetOrderItemByID(c, uuid.MustParse(req.OrderItemID))
	if err != nil {
		c.JSON(400, createErrorResponse[gin.H](InvalidBodyCode, "invalid request", err))
		return
	}
	if orderItem.CustomerID != auth.UserID {
		c.JSON(403, createErrorResponse[gin.H](InvalidBodyCode, "forbidden", nil))
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
		c.JSON(500, createErrorResponse[bool](InternalServerErrorCode, "internal server error", err))
		return
	}

	images := make([]RatingImageModel, 0)
	if len(req.Files) > 0 {
		var wg sync.WaitGroup
		wg.Add(len(req.Files))
		uploadResults := make([]repository.InsertImageParams, len(req.Files))
		assignments := make([]repository.InsertImageAssignmentParams, len(req.Files))
		errs := make([]error, len(req.Files))
		for i, file := range req.Files {
			go func(i int, file *multipart.FileHeader) {
				defer wg.Done()
				ID, url, err := s.uploadService.UploadFile(c, file)
				if err != nil {
					errs[i] = err
					return
				}
				uploadResults[i] = repository.InsertImageParams{
					ExternalID: ID,
					Url:        url,
					AltText:    utils.StringPtr(file.Filename),
					Caption:    utils.StringPtr(fmt.Sprintf("Image %d", i+1)),
					MimeType:   utils.StringPtr(file.Header.Get("Content-Type")),
					FileSize:   utils.Int64Ptr(file.Size),
				}
				assignments[i] = repository.InsertImageAssignmentParams{
					EntityID:   rating.ID,
					EntityType: "product_rating",
					Role:       string(repository.ImageRoleGallery),
				}
			}(i, file)
		}
		wg.Wait()
		for _, err := range errs {
			if err != nil {
				c.JSON(500, createErrorResponse[gin.H](InternalServerErrorCode, "internal server error", err))
				return
			}
		}

		// Insert product images
		for i, uploadResult := range uploadResults {
			img, err := s.repo.InsertImage(c, uploadResult)
			if err != nil {
				c.JSON(500, createErrorResponse[gin.H](InternalServerErrorCode, "internal server error", err))
				return
			}
			assignments[i].ImageID = img.ID
			_, err = s.repo.InsertImageAssignment(c, assignments[i])
			if err != nil {
				c.JSON(500, createErrorResponse[gin.H](InternalServerErrorCode, "internal server error", err))
				return
			}
			images = append(images, struct {
				ID  string `json:"id"`
				URL string `json:"url"`
			}{
				ID:  img.ID.String(),
				URL: img.Url,
			})
		}
	}

	resp := ProductRatingModel{
		ID:               rating.ID,
		UserID:           rating.UserID,
		Rating:           req.Rating,
		ReviewTitle:      *rating.ReviewTitle,
		ReviewContent:    *rating.ReviewContent,
		VerifiedPurchase: rating.VerifiedPurchase,
		HelpfulVotes:     rating.HelpfulVotes,
		UnhelpfulVotes:   rating.UnhelpfulVotes,
		Images:           images,
	}

	c.JSON(200, createSuccessResponse(c, resp, "success", nil, nil))
}

// @Summary Post a helpful rating
// @Description Post a helpful rating
// @Tags ratings
// @Accept json
// @Produce json
// @Param rating_id body string true "Rating ID"
// @Param helpful body bool true "Helpful"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[bool]
// @Failure 400 {object} ApiResponse[bool]
// @Failure 403 {object} ApiResponse[bool]
// @Failure 500 {object} ApiResponse[bool]
// @Router /ratings/{rating_id}/helpful [post]
func (s *Server) postRatingHelpfulHandler(c *gin.Context) {
	var param UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(400, createErrorResponse[gin.H](InvalidBodyCode, "invalid request", err))
		return
	}
	auth, _ := c.MustGet(authorizationPayload).(*auth.Payload)
	var req PostHelpfulRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, createErrorResponse[bool](InvalidBodyCode, "invalid request", err))
		return
	}

	rating, err := s.repo.GetProductRating(c, uuid.MustParse(param.ID))
	if err != nil {
		c.JSON(400, createErrorResponse[bool](InvalidBodyCode, "invalid request", err))
		return
	}
	if rating.UserID == auth.UserID {
		c.JSON(403, createErrorResponse[bool](InvalidBodyCode, "forbidden", nil))
		return
	}

	id, err := s.repo.VoteHelpfulRatingTx(c, repository.VoteHelpfulRatingTxArgs{
		UserID:         auth.UserID,
		RatingID:       rating.ID,
		Helpful:        req.Helpful,
		HelpfulVotes:   rating.HelpfulVotes,
		UnhelpfulVotes: rating.UnhelpfulVotes,
	})

	if err != nil {
		c.JSON(500, createErrorResponse[bool](InternalServerErrorCode, "internal server error", err))
	}

	c.JSON(200, createSuccessResponse(c, id, "success", nil, nil))
}

// @Summary Post a reply to a rating
// @Description Post a reply to a product rating
// @Tags ratings
// @Accept json
// @Produce json
// @Param rating_id path string true "Rating ID"
// @Param content body string true "Reply Content"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[bool]
// @Failure 400 {object} ApiResponse[bool]
// @Failure 403 {object} ApiResponse[bool]
// @Failure 500 {object} ApiResponse[bool]
// @Router /ratings/{rating_id}/reply [post]
func (s *Server) postReplyRatingHandler(c *gin.Context) {
	var param UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(400, createErrorResponse[gin.H](InvalidBodyCode, "invalid request", err))
		return
	}
	auth, _ := c.MustGet(authorizationPayload).(*auth.Payload)

	var req PostReplyRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, createErrorResponse[bool](InvalidBodyCode, "invalid request", err))
		return
	}
	rating, err := s.repo.GetProductRating(c, uuid.MustParse(param.ID))
	if err != nil {
		c.JSON(400, createErrorResponse[bool](InvalidBodyCode, "invalid request", err))
		return
	}

	reply, err := s.repo.InsertRatingReply(c, repository.InsertRatingReplyParams{
		RatingID: rating.ID,
		ReplyBy:  auth.UserID,
		Content:  req.Content,
	})
	if err != nil {
		c.JSON(500, createErrorResponse[bool](InternalServerErrorCode, "internal server error", err))
		return
	}

	c.JSON(200, createSuccessResponse(c, reply.ID, "success", nil, nil))
}

// @Summary Get product ratings
// @Description Get ratings for a specific product
// @Tags ratings
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Success 200 {object} ApiResponse[[]ProductRatingModel]
// @Failure 400 {object} ApiResponse[bool]
// @Failure 404 {object} ApiResponse[bool]
// @Failure 500 {object} ApiResponse[bool]
// @Router /admin/ratings [get]
func (s *Server) getRatingsHandler(c *gin.Context) {
	var queries RatingsQueryParams
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(400, createErrorResponse[bool](InvalidBodyCode, "invalid request", err))
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
		c.JSON(400, createErrorResponse[bool](InvalidBodyCode, "invalid request", err))
		return
	}

	ratingsCount, err := s.repo.CountProductRatings(c, pgtype.UUID{
		Bytes: uuid.Nil,
		Valid: false,
	})

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
		if prIdx != -1 && rating.ImageID.Valid {
			id, _ := uuid.FromBytes(rating.ImageID.Bytes[:])
			productRatings[prIdx].Images = append(productRatings[prIdx].Images, RatingImageModel{
				ID:  id.String(),
				URL: *rating.ImageUrl,
			})
			continue
		}
		model := ProductRatingModel{
			ID:               rating.ID,
			UserID:           rating.UserID,
			Name:             rating.Fullname,
			ProductName:      rating.ProductName,
			Rating:           ratingPoint.Float64,
			IsVisible:        rating.IsVisible,
			IsApproved:       rating.IsApproved,
			ReviewTitle:      *rating.ReviewTitle,
			ReviewContent:    *rating.ReviewContent,
			VerifiedPurchase: rating.VerifiedPurchase,
			HelpfulVotes:     rating.HelpfulVotes,
			UnhelpfulVotes:   rating.UnhelpfulVotes,
		}
		if rating.ImageID.Valid {
			id, _ := uuid.FromBytes(rating.ImageID.Bytes[:])
			model.Images = append(model.Images, RatingImageModel{
				ID:  id.String(),
				URL: *rating.ImageUrl,
			})
		}
		productRatings = append(productRatings, model)
	}
	c.JSON(200, createSuccessResponse(
		c,
		productRatings,
		"success",
		&Pagination{
			Total:           ratingsCount,
			Page:            queries.Page,
			PageSize:        queries.PageSize,
			TotalPages:      utils.CalculateTotalPages(ratingsCount, queries.PageSize),
			HasNextPage:     queries.Page*queries.PageSize < ratingsCount,
			HasPreviousPage: queries.Page > 1,
		}, nil),
	)
}

// @Summary Get product ratings
// @Description Get ratings for a specific product
// @Tags ratings
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Success 200 {object} ApiResponse[[]ProductRatingModel]
// @Failure 400 {object} ApiResponse[bool]
// @Failure 404 {object} ApiResponse[bool]
// @Failure 500 {object} ApiResponse[bool]
// @Router /ratings/products/{product_id} [get]
func (s *Server) getRatingsByProductHandler(c *gin.Context) {
	var param UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(400, createErrorResponse[bool](InvalidBodyCode, "invalid request", err))
		return
	}
	var queries PaginationQueryParams
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(400, createErrorResponse[bool](InvalidBodyCode, "invalid request", err))
		return
	}

	ratings, err := s.repo.GetProductRatings(c, repository.GetProductRatingsParams{
		ProductID: utils.GetPgTypeUUID(uuid.MustParse(param.ID)),
		Limit:     queries.PageSize,
		Offset:    (queries.Page - 1) * queries.PageSize,
	})
	if err != nil {
		c.JSON(400, createErrorResponse[bool](InvalidBodyCode, "invalid request", err))
		return
	}

	ratingsCount, err := s.repo.CountProductRatings(c, utils.GetPgTypeUUID(uuid.MustParse(param.ID)))

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
		if prIdx != -1 && rating.ImageID.Valid {
			id, _ := uuid.FromBytes(rating.ImageID.Bytes[:])
			productRatings[prIdx].Images = append(productRatings[prIdx].Images, RatingImageModel{
				ID:  id.String(),
				URL: *rating.ImageUrl,
			})
			continue
		}
		model := ProductRatingModel{
			ID:               rating.ID,
			UserID:           rating.UserID,
			Name:             rating.Fullname,
			Rating:           ratingPoint.Float64,
			ProductName:      rating.ProductName,
			IsVisible:        rating.IsVisible,
			ReviewTitle:      *rating.ReviewTitle,
			ReviewContent:    *rating.ReviewContent,
			VerifiedPurchase: rating.VerifiedPurchase,
			HelpfulVotes:     rating.HelpfulVotes,
			UnhelpfulVotes:   rating.UnhelpfulVotes,
		}
		if rating.ImageID.Valid {
			id, _ := uuid.FromBytes(rating.ImageID.Bytes[:])
			model.Images = append(model.Images, RatingImageModel{
				ID:  id.String(),
				URL: *rating.ImageUrl,
			})
		}
		productRatings = append(productRatings, model)
	}

	c.JSON(200, createSuccessResponse(
		c,
		productRatings,
		"success",
		&Pagination{
			Total:           ratingsCount,
			Page:            queries.Page,
			PageSize:        queries.PageSize,
			TotalPages:      utils.CalculateTotalPages(ratingsCount, queries.PageSize),
			HasNextPage:     queries.Page*queries.PageSize < ratingsCount,
			HasPreviousPage: queries.Page > 1,
		}, nil),
	)
}

// @Summary Get order ratings
// @Description Get ratings for a specific order
// @Tags ratings
// @Accept json
// @Produce json
// @Param order_id path string true "Order ID"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[[]ProductRatingModel]
// @Failure 400 {object} ApiResponse[bool]
// @Failure 403 {object} ApiResponse[bool]
// @Failure 404 {object} ApiResponse[bool]
// @Failure 500 {object} ApiResponse[bool]
// @Router /ratings/orders/{order_id} [get]
func (s *Server) getOrderRatingsHandler(c *gin.Context) {
	auth, _ := c.MustGet(authorizationPayload).(*auth.Payload)

	var param struct {
		OrderID string `uri:"order_id" binding:"required,uuid"`
	}
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(400, createErrorResponse[bool](InvalidBodyCode, "invalid request", err))
		return
	}

	orderItems, err := s.repo.GetOrderItemsByOrderID(c, uuid.MustParse(param.OrderID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "internal server error", err))
		return
	}
	if len(orderItems) == 0 {
		c.JSON(http.StatusNotFound, createErrorResponse[bool](NotFoundCode, "order items not found", nil))
		return
	}
	if orderItems[0].CustomerID != auth.UserID {
		c.JSON(http.StatusForbidden, createErrorResponse[bool](PermissionDeniedCode, "forbidden", nil))
		return
	}
	orderItemIds := make([]uuid.UUID, len(orderItems))
	for i, orderItem := range orderItems {
		orderItemIds[i] = orderItem.OrderItemID
	}
	ratings, err := s.repo.GetProductRatingsByOrderItemIDs(c, orderItemIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "internal server error", err))
		return
	}

	c.JSON(200, createSuccessResponse(c, ratings, "success", nil, nil))
}

// @Summary Delete a rating
// @Description Delete a product rating by ID
// @Tags Admin, ratings
// @Accept json
// @Produce json
// @Param id path string true "Rating ID"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[bool]
// @Failure 400 {object} ApiResponse[bool]
// @Failure 404 {object} ApiResponse[bool]
// @Failure 500 {object} ApiResponse[bool]
// @Router /admin/ratings/{id} [delete]
func (sv *Server) deleteRatingHandler(c *gin.Context) {
	// Parse the rating ID from the URL
	var param struct {
		ID string `uri:"id" binding:"required,uuid"`
	}

	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[bool](InvalidBodyCode, "Invalid rating ID", err))
		return
	}

	// Convert the ID to UUID
	ratingID, err := uuid.Parse(param.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[bool](InvalidBodyCode, "Invalid rating ID format", err))
		return
	}

	// Check if rating exists first
	_, err = sv.repo.GetProductRating(c, ratingID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[bool](NotFoundCode, "Rating not found", err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
		return
	}

	// Delete the rating
	err = sv.repo.DeleteProductRating(c, ratingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "Failed to delete rating", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, true, "Rating deleted successfully", nil, nil))
}

// @Summary Approve a rating
// @Description Approve a product rating by ID
// @Tags Admin, ratings
// @Accept json
// @Produce json
// @Param id path string true "Rating ID"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[bool]
// @Failure 400 {object} ApiResponse[bool]
// @Failure 404 {object} ApiResponse[bool]
// @Failure 500 {object} ApiResponse[bool]
// @Router /admin/ratings/{id}/approve [post]
func (sv *Server) approveRatingHandler(c *gin.Context) {
	// Parse the rating ID from the URL
	var param struct {
		ID string `uri:"id" binding:"required,uuid"`
	}

	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[bool](InvalidBodyCode, "Invalid rating ID", err))
		return
	}

	// Convert the ID to UUID
	ratingID, err := uuid.Parse(param.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[bool](InvalidBodyCode, "Invalid rating ID format", err))
		return
	}

	// Check if rating exists first
	rating, err := sv.repo.GetProductRating(c, ratingID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[bool](NotFoundCode, "Rating not found", err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
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
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "Failed to approve rating", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, true, "Rating approved successfully", nil, nil))
}

// @Summary Ban a user from rating
// @Description Ban a user from rating by setting their rating to invisible
// @Tags Admin, ratings
// @Accept json
// @Produce json
// @Param id path string true "Rating ID"
// @Security BearerAuth
// @Success 200 {object} ApiResponse[bool]
// @Failure 400 {object} ApiResponse[bool]
// @Failure 404 {object} ApiResponse[bool]
// @Failure 500 {object} ApiResponse[bool]
// @Router /admin/ratings/{id}/ban [post]
func (sv *Server) banUserRatingHandler(c *gin.Context) {
	// Parse the rating ID from the URL
	var param struct {
		ID string `uri:"id" binding:"required,uuid"`
	}

	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[bool](InvalidBodyCode, "Invalid rating ID", err))
		return
	}

	// Convert the ID to UUID
	ratingID, err := uuid.Parse(param.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[bool](InvalidBodyCode, "Invalid rating ID format", err))
		return
	}

	// Check if rating exists first
	rating, err := sv.repo.GetProductRating(c, ratingID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[bool](NotFoundCode, "Rating not found", err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
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
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "Failed to ban user rating", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, true, "User rating banned successfully", nil, nil))
}
