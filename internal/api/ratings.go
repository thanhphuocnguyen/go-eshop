package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/constants"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
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
func (s *Server) postRating(c *gin.Context) {
	auth, _ := c.MustGet(constants.AuthPayLoad).(*auth.TokenPayload)

	var req models.PostRatingFormData
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	orderItem, err := s.repo.GetOrderItemByID(c, uuid.MustParse(req.OrderItemID))
	if err != nil {
		c.JSON(400, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	if orderItem.UserID != auth.UserID {
		c.JSON(403, dto.CreateErr(InvalidBodyCode, nil))
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
func (s *Server) postRatingHelpful(c *gin.Context) {
	var param models.UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(400, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	auth, _ := c.MustGet(constants.AuthPayLoad).(*auth.TokenPayload)
	var req models.PostHelpfulRatingModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	rating, err := s.repo.GetProductRating(c, uuid.MustParse(param.ID))
	if err != nil {
		c.JSON(400, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	if rating.UserID == auth.UserID {
		c.JSON(403, dto.CreateErr(InvalidBodyCode, nil))
		return
	}

	id, err := s.repo.VoteHelpfulRatingTx(c, repository.VoteHelpfulRatingTxArgs{
		UserID:   auth.UserID,
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
func (s *Server) postReplyRating(c *gin.Context) {
	var param models.UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(400, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	auth, _ := c.MustGet(constants.AuthPayLoad).(*auth.TokenPayload)

	var req models.PostReplyRatingModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	rating, err := s.repo.GetProductRating(c, uuid.MustParse(param.ID))
	if err != nil {
		c.JSON(400, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	reply, err := s.repo.InsertRatingReply(c, repository.InsertRatingReplyParams{
		RatingID: rating.ID,
		ReplyBy:  auth.UserID,
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
func (s *Server) getRatingsByProduct(c *gin.Context) {
	var param models.UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(400, dto.CreateErr(InvalidBodyCode, err))
		return
	}
	var queries models.PaginationQuery
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(400, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	ratings, err := s.repo.GetProductRatings(c, repository.GetProductRatingsParams{
		ProductID: utils.GetPgTypeUUID(uuid.MustParse(param.ID)),
		Limit:     queries.PageSize,
		Offset:    (queries.Page - 1) * queries.PageSize,
	})
	if err != nil {
		c.JSON(400, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	ratingsCount, err := s.repo.CountProductRatings(c, utils.GetPgTypeUUID(uuid.MustParse(param.ID)))
	if err != nil {
		c.JSON(400, dto.CreateErr(InvalidBodyCode, err))
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

	c.JSON(200, dto.CreateDataResp(c, productRatings, dto.CreatePagination(queries.Page, queries.PageSize, ratingsCount), nil))
}

// Setup brand-related routes
func (sv *Server) addRatingRoutes(rg *gin.RouterGroup) {
	ratings := rg.Group("ratings", authenticateMiddleware(sv.tokenGenerator))
	{
		ratings.POST("", sv.postRating)
		ratings.GET(":orderId", sv.AdminGetOrderRatings)
		ratings.POST(":id/helpful", sv.postRatingHelpful)
		ratings.POST(":id/reply", sv.postReplyRating)
	}
}
