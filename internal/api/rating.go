package api

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
)

// ----- Types -----
type PostRatingFormData struct {
	OrderItemID string                  `form:"order_item_id" binding:"required"`
	Rating      float64                 `form:"rating" binding:"required,min=1,max=5"`
	Title       string                  `form:"title" binding:"required"`
	Content     string                  `form:"content" binding:"required"`
	Files       []*multipart.FileHeader `form:"files" binding:"omitempty"`
}

type PostHelpfulRatingRequest struct {
	RatingID string `json:"rating_id" binding:"required"`
	Helpful  bool   `json:"helpful" binding:"required"`
}

type PostReplyRatingRequest struct {
	RatingID string `json:"rating_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

type ProductRatingModel struct {
	ID               uuid.UUID `json:"id"`
	UserID           uuid.UUID `json:"user_id"`
	Rating           float64   `json:"rating"`
	ReviewTitle      string    `json:"review_title"`
	ReviewContent    string    `json:"review_content"`
	VerifiedPurchase bool      `json:"verified_purchase"`
	HelpfulVotes     int32     `json:"helpful_votes"`
	UnhelpfulVotes   int32     `json:"unhelpful_votes"`
}

// ----- Rating Handlers -----
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
		ID:               uuid.New(),
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
					AltText:    StringPtr(file.Filename),
					Caption:    StringPtr(fmt.Sprintf("Image %d", i+1)),
					MimeType:   StringPtr(file.Header.Get("Content-Type")),
					FileSize:   Int64Ptr(file.Size),
				}
				assignments[i] = repository.InsertImageAssignmentParams{
					EntityID:   rating.ID,
					EntityType: string(repository.EntityTypeRating),
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
		}
	}

	c.JSON(200, createSuccessResponse(c, rating, "success", nil, nil))
}

func (s *Server) postRatingHelpfulHandler(c *gin.Context) {
	auth, _ := c.MustGet(authorizationPayload).(*auth.Payload)

	var req PostHelpfulRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, createErrorResponse[bool](InvalidBodyCode, "invalid request", err))
		return
	}
	rating, err := s.repo.GetProductRating(c, uuid.MustParse(req.RatingID))
	if err != nil {
		c.JSON(400, createErrorResponse[bool](InvalidBodyCode, "invalid request", err))
		return
	}
	if rating.UserID != auth.UserID {
		c.JSON(403, createErrorResponse[bool](InvalidBodyCode, "forbidden", nil))
		return
	}

	// add rating helpful record
	helpfulVote, err := s.repo.InsertRatingVotes(c, repository.InsertRatingVotesParams{
		RatingID:  uuid.MustParse(req.RatingID),
		UserID:    auth.UserID,
		IsHelpful: req.Helpful,
	})
	if err != nil {
		c.JSON(500, createErrorResponse[bool](InternalServerErrorCode, "internal server error", err))
		return
	}

	updateRatingParams := repository.UpdateProductRatingParams{
		ID: uuid.MustParse(req.RatingID),
	}
	if req.Helpful {
		updateRatingParams.HelpfulVotes = Int32Ptr(rating.HelpfulVotes + 1)
	} else {
		updateRatingParams.UnhelpfulVotes = Int32Ptr(rating.UnhelpfulVotes + 1)
	}
	_, err = s.repo.UpdateProductRating(c, updateRatingParams)
	if err != nil {
		c.JSON(500, createErrorResponse[bool](InternalServerErrorCode, "internal server error", err))
		return
	}

	c.JSON(200, createSuccessResponse(c, helpfulVote.ID, "success", nil, nil))
}

func (s *Server) postReplyRatingHandler(c *gin.Context) {
	auth, _ := c.MustGet(authorizationPayload).(*auth.Payload)

	var req PostReplyRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, createErrorResponse[bool](InvalidBodyCode, "invalid request", err))
		return
	}
	rating, err := s.repo.GetProductRating(c, uuid.MustParse(req.RatingID))
	if err != nil {
		c.JSON(400, createErrorResponse[bool](InvalidBodyCode, "invalid request", err))
		return
	}

	reply, err := s.repo.InsertRatingReply(c, repository.InsertRatingReplyParams{
		ID:       uuid.New(),
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

func (s *Server) getProductRatingsHandler(c *gin.Context) {
	var param ProductParam
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
		ProductID: uuid.MustParse(param.ID),
		Limit:     queries.PageSize,
		Offset:    (queries.Page - 1) * queries.PageSize,
	})
	if err != nil {
		c.JSON(400, createErrorResponse[bool](InvalidBodyCode, "invalid request", err))
		return
	}

	ratingsCount, err := s.repo.CountProductRatings(c, uuid.MustParse(param.ID))

	productRatings := make([]ProductRatingModel, 0)
	for _, rating := range ratings {
		ratingPoint, _ := rating.ProductRating.Rating.Float64Value()
		productRatings = append(productRatings, ProductRatingModel{
			ID:               rating.ProductRating.ID,
			UserID:           rating.UserID,
			Rating:           ratingPoint.Float64,
			ReviewTitle:      *rating.ProductRating.ReviewTitle,
			ReviewContent:    *rating.ProductRating.ReviewContent,
			VerifiedPurchase: rating.ProductRating.VerifiedPurchase,
			HelpfulVotes:     rating.ProductRating.HelpfulVotes,
			UnhelpfulVotes:   rating.ProductRating.UnhelpfulVotes,
		})
	}

	c.JSON(200, createSuccessResponse(
		c,
		productRatings,
		"success",
		&Pagination{
			Total:           ratingsCount,
			Page:            queries.Page,
			PageSize:        queries.PageSize,
			TotalPages:      CalculateTotalPages(ratingsCount, queries.PageSize),
			HasNextPage:     queries.Page*queries.PageSize < ratingsCount,
			HasPreviousPage: queries.Page > 1,
		}, nil),
	)
}

// @Summary Delete a rating
// @Description Delete a product rating by ID
// @Tags admin, ratings
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
// @Tags admin, ratings
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
// @Tags admin, ratings
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
