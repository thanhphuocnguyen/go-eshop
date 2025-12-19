package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
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
func (s *Server) postRating(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())

	var req models.PostRatingFormData
	if err := c.ShouldBind(&req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	orderItem, err := s.repo.GetOrderItemByID(c, uuid.MustParse(req.OrderItemID))
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	if orderItem.UserID != claims["userId"].(uuid.UUID) {
		c.JSON(403, dto.CreateErr(InvalidBodyCode, nil))
		return
	}

	rating, err := s.repo.InsertProductRating(c, repository.InsertProductRatingParams{
		ProductID:        orderItem.ProductID,
		UserID:           claims["userId"].(uuid.UUID),
		OrderItemID:      utils.GetPgTypeUUID(orderItem.OrderItemID),
		Rating:           utils.GetPgNumericFromFloat(req.Rating),
		ReviewTitle:      &req.Title,
		ReviewContent:    &req.Content,
		VerifiedPurchase: true,
	})

	if err != nil {
		c.JSON(500, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	resp := dto.ProductRatingDetail{
		ID:               rating.ID.String(),
		UserID:           rating.UserID.String(),
		Rating:           req.Rating,
		ReviewTitle:      *rating.ReviewTitle,
		ReviewContent:    *rating.ReviewContent,
		VerifiedPurchase: rating.VerifiedPurchase,
		// Images:           images,
	}

	c.JSON(200, dto.CreateDataResp(c, resp, nil, nil))
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
func (s *Server) postRatingHelpful(w http.ResponseWriter, r *http.Request) {
	var param models.UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	_, claims, err := jwtauth.FromContext(r.Context())
	var req models.PostHelpfulRatingModel
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	rating, err := s.repo.GetProductRating(c, uuid.MustParse(param.ID))
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	if rating.UserID == claims["userId"].(uuid.UUID) {
		c.JSON(403, dto.CreateErr(InvalidBodyCode, nil))
		return
	}

	id, err := s.repo.VoteHelpfulRatingTx(c, repository.VoteHelpfulRatingTxArgs{
		UserID:   claims["userId"].(uuid.UUID),
		RatingID: rating.ID,
		Helpful:  req.Helpful,
	})

	if err != nil {
		c.JSON(500, dto.CreateErr(InternalServerErrorCode, err))
	}

	c.JSON(200, dto.CreateDataResp(c, id, nil, nil))
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
func (s *Server) postReplyRating(w http.ResponseWriter, r *http.Request) {
	var param models.UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	_, claims, err := jwtauth.FromContext(r.Context())

	var req models.PostReplyRatingModel
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	rating, err := s.repo.GetProductRating(c, uuid.MustParse(param.ID))
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	reply, err := s.repo.InsertRatingReply(c, repository.InsertRatingReplyParams{
		RatingID: rating.ID,
		ReplyBy:  claims["userId"].(uuid.UUID),
		Content:  req.Content,
	})
	if err != nil {
		c.JSON(500, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(200, dto.CreateDataResp(c, reply.ID, nil, nil))
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
func (s *Server) getRatingsByProduct(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		RespondBadRequest(w, InvalidBodyCode, errors.New("id parameter is required"))
		return
	}

	paginationQuery, err := ParsePaginationQuery(r)
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	ratings, err := s.repo.GetProductRatings(r.Context(), repository.GetProductRatingsParams{
		ProductID: utils.GetPgTypeUUID(uuid.MustParse(idParam)),
		Limit:     paginationQuery.PageSize,
		Offset:    (paginationQuery.Page - 1) * paginationQuery.PageSize,
	})
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	ratingsCount, err := s.repo.CountProductRatings(r.Context(), utils.GetPgTypeUUID(uuid.MustParse(idParam)))
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
			Rating:           ratingPoint.Float64,
			ProductName:      rating.ProductName,
			IsVisible:        rating.IsVisible,
			ReviewTitle:      *rating.ReviewTitle,
			ReviewContent:    *rating.ReviewContent,
			VerifiedPurchase: rating.VerifiedPurchase,
		}
		if rating.ImageID != nil {
			model.Images = append(model.Images, dto.RatingImage{
				ID:  *rating.ImageID,
				URL: *rating.ImageUrl,
			})
		}
		productRatings = append(productRatings, model)
	}

	RespondSuccessWithPagination(w, r, productRatings, dto.CreatePagination(paginationQuery.Page, paginationQuery.PageSize, ratingsCount))
}

// Setup brand-related routes
func (sv *Server) addRatingRoutes(r chi.Router) {
	r.Route("ratings", func(r chi.Router) {
		r.Post("", sv.postRating)
		r.Get(":orderId", sv.adminGetOrderRatings)
		r.Post(":id/helpful", sv.postRatingHelpful)
		r.Post(":id/reply", sv.postReplyRating)
	})
}
