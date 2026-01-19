package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

// @Summary Get product ratings
// @Description Get ratings for a specific product
// @Tags ratings
// @Accept json
// @Produce json
// @Param productId path string true "Product ID"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Success 200 {object} dto.ApiResponse[[]dto.ProductRatingDetail]
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/ratings [get]
func (s *Server) adminGetRatings(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	var queries models.PaginationQuery = ParsePaginationQuery(r)
	status := r.URL.Query().Get("status")
	sqlParams := repository.GetProductRatingsParams{
		Limit:  queries.PageSize,
		Offset: (queries.Page - 1) * queries.PageSize,
	}

	if status != "" {
		switch status {
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
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	ratingsCount, err := s.repo.CountProductRatings(c, pgtype.UUID{
		Bytes: uuid.Nil,
		Valid: false,
	})
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	productRatings := make([]dto.ProductRatingDetail, 0)
	for _, rating := range ratings {
		ratingPoint, _ := rating.Rating.Float64Value()
		prIdx := -1
		for i, pr := range productRatings {
			if pr.ID == rating.ID.String() {
				prIdx = i
				break
			}
		}
		if prIdx != -1 && rating.ImageID != nil {
			productRatings[prIdx].Images = append(productRatings[prIdx].Images, dto.RatingImage{
				ID:  *rating.ImageID,
				URL: *rating.ImageUrl,
			})
			continue
		}
		model := dto.ProductRatingDetail{
			ID:               rating.ID.String(),
			UserID:           rating.UserID.String(),
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
			model.Images = append(model.Images, dto.RatingImage{
				ID:  *rating.ImageID,
				URL: *rating.ImageUrl,
			})
		}
		productRatings = append(productRatings, model)
	}
	RespondSuccessWithPagination(w, productRatings, dto.CreatePagination(queries.Page, queries.PageSize, ratingsCount))
}

// @Summary Get order ratings
// @Description Get ratings for a specific order
// @Tags ratings
// @Accept json
// @Produce json
// @Param orderId path string true "Order ID"
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse[[]repository.GetProductRatingsByOrderItemIDsRow]
// @Failure 400 {object} ErrorResp
// @Failure 403 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /ratings/orders/{orderId} [get]
func (s *Server) adminGetOrderRatings(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	_, claims, err := jwtauth.FromContext(c)

	orderId, err := GetUrlParam(r, "orderId")

	orderItems, err := s.repo.GetOrderItemsByOrderID(c, uuid.MustParse(orderId))
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	if len(orderItems) == 0 {
		RespondNotFound(w, NotFoundCode, nil)
		return
	}
	userID := uuid.MustParse(claims["userId"].(string))

	if orderItems[0].UserID != userID {
		RespondForbidden(w, PermissionDeniedCode, nil)
		return
	}
	orderItemIds := make([]uuid.UUID, len(orderItems))
	for i, orderItem := range orderItems {
		orderItemIds[i] = orderItem.OrderItemID
	}
	ratings, err := s.repo.GetProductRatingsByOrderItemIDs(c, orderItemIds)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondSuccess(w, ratings)
}

// @Summary Delete a rating
// @Description Delete a product rating by ID
// @Tags admin, ratings
// @Accept json
// @Produce json
// @Param id path string true "Rating ID"
// @Security BearerAuth
// @Success 204 {object} nil
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/ratings/{id} [delete]
func (s *Server) adminDeleteRating(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	// Parse the rating ID from the URL
	id, err := GetUrlParam(r, "id")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	// Convert the ID to UUID
	ratingID, err := uuid.Parse(id)
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	// Check if rating exists first
	_, err = s.repo.GetProductRating(c, ratingID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, err)
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	// Delete the rating
	err = s.repo.DeleteProductRating(c, ratingID)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondNoContent(w)
}

// @Summary Approve a rating
// @Description Approve a product rating by ID
// @Tags admin, ratings
// @Accept json
// @Produce json
// @Param id path string true "Rating ID"
// @Security BearerAuth
// @Success 204 {object} nil
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/ratings/{id}/approve [post]
func (s *Server) adminApproveRating(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	// Parse the rating ID from the URL
	id, err := GetUrlParam(r, "id")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	// Convert the ID to UUID
	ratingID, err := uuid.Parse(id)
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	// Check if rating exists first
	rating, err := s.repo.GetProductRating(c, ratingID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, err)
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	// Set IsApproved to true
	isApproved := true

	// Update the rating
	_, err = s.repo.UpdateProductRating(c, repository.UpdateProductRatingParams{
		ID:         rating.ID,
		IsApproved: &isApproved,
	})

	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondNoContent(w)
}

// @Summary Ban a user from rating
// @Description Ban a user from rating by setting their rating to invisible
// @Tags admin, ratings
// @Accept json
// @Produce json
// @Param id path string true "Rating ID"
// @Security BearerAuth
// @Success 204 {object} nil
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/ratings/{id}/ban [post]
func (s *Server) adminBanUserRating(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	// Parse the rating ID from the URL
	id, err := GetUrlParam(r, "id")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	// Convert the ID to UUID
	ratingID, err := uuid.Parse(id)
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	// Check if rating exists first
	rating, err := s.repo.GetProductRating(c, ratingID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, err)
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	// Set IsVisible to false
	isVisible := false

	// Update the rating
	_, err = s.repo.UpdateProductRating(c, repository.UpdateProductRatingParams{
		ID:        rating.ID,
		IsVisible: &isVisible,
	})

	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondNoContent(w)
}
