package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/thanhphuocnguyen/go-eshop/internal/constants"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
	"github.com/thanhphuocnguyen/go-eshop/internal/worker"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
)

// updateUser godoc
// @Summary Update user info
// @Description Update user info
// @Tags users
// @Accept  json
// @Produce  json
// @Param input body UpdateUserRequest true "User info"
// @Success 200 {object} ApiResponse[repository.UpdateUserRow]
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /users/{id} [patch]
func (sv *Server) updateUser(c *gin.Context) {
	authPayload := c.MustGet(constants.AuthPayLoad).(*auth.TokenPayload)
	var req models.UpdateUserModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidEmailCode, err))
		return
	}
	userId := uuid.MustParse(req.UserID)
	user, err := sv.repo.GetUserByID(c, userId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.CreateErr(UnauthorizedCode, err))
		return
	}

	if authPayload.RoleCode != "admin" && user.ID != userId {
		c.JSON(http.StatusUnauthorized, dto.CreateErr(UnauthorizedCode, err))
		return
	}

	arg := repository.UpdateUserParams{
		ID: userId,
	}

	boolVal := false

	if req.Email != nil {
		arg.Email = req.Email
		if user.Email != *req.Email {
			arg.VerifiedEmail = &boolVal
		}
	}

	if req.FirstName != nil {
		arg.FirstName = req.FirstName
	}

	if req.LastName != nil {
		arg.LastName = req.LastName
	}

	if req.Phone != nil {
		arg.PhoneNumber = req.Phone
		if user.PhoneNumber != *req.Phone {
			arg.VerifiedPhone = &boolVal
		}
	}

	updatedUser, err := sv.repo.UpdateUser(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, updatedUser, nil, nil))
}

// getCurrentUser godoc
// @Summary Get user info
// @Description Get user info
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} ApiResponse[UserDetail]
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /users/me [get]
func (sv *Server) getCurrentUser(c *gin.Context) {
	authPayload, ok := c.MustGet(constants.AuthPayLoad).(*auth.TokenPayload)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, errors.New("authorization payload is not provided")))
		return
	}

	var userResp dto.UserDetail

	user, err := sv.repo.GetUserByID(c, authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.CreateErr(NotFoundCode, err))
		return
	}

	userAddress, err := sv.repo.GetAddresses(c, user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	addressResp := make([]dto.AddressDetail, 0)
	for _, address := range userAddress {
		addressResp = append(addressResp, dto.MapAddressResponse(address))
	}
	userResp = dto.MapToUserResponse(user, authPayload.RoleCode)
	userResp.Addresses = addressResp

	c.JSON(http.StatusOK, dto.CreateDataResp(c, userResp, nil, nil))
}

// sendVerifyEmail godoc
// @Summary Send verify email
// @Description Send verify email
// @Tags users
// @Accept  json
// @Produce  json
// @Success 204 {object} nil
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /users/verify-email [post]
// @Security BearerAuth
func (sv *Server) sendVerifyEmail(c *gin.Context) {
	authPayload, ok := c.MustGet(constants.AuthPayLoad).(*auth.TokenPayload)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, errors.New("authorization payload is not provided")))
		return
	}
	user, err := sv.repo.GetUserByID(c, authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	if user.VerifiedEmail {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidEmailCode, fmt.Errorf("email already verified")))
		return
	}

	err = sv.taskDistributor.SendVerifyAccountEmail(
		c,
		&worker.PayloadVerifyEmail{
			UserID: authPayload.UserID,
		},
		asynq.MaxRetry(3),
		asynq.ProcessIn(5*time.Second),
		asynq.Queue(worker.QueueCritical),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.Status(http.StatusNoContent)
}

// VerifyEmail godoc
// @Summary Verify email
// @Description Verify email
// @Tags users
// @Accept  json
// @Produce  json
// @Param id query int true "ID"
// @Param verify_code query string true "Verify code"
// @Success 200 {object} nil
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /users/verify-email [get]
func (sv *Server) VerifyEmail(c *gin.Context) {
	var query models.VerifyEmailQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidEmailCode, err))
		return
	}
	verifyEmail, err := sv.repo.GetVerifyEmailByVerifyCode(c, query.VerifyCode)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, err))
		return
	}

	// Create a transaction to ensure both operations succeed or fail together
	err = sv.repo.VerifyEmailTx(c, repository.VerifyEmailTxArgs{
		VerifyEmail: verifyEmail,
		VerifyCode:  query.VerifyCode,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	user, err := sv.repo.GetUserByID(c, verifyEmail.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	// Render HTML success page
	c.Header("Content-Type", "text/html")
	c.HTML(http.StatusOK, "verification-success.html", gin.H{
		"username": user.Username,
		"email":    user.Email,
	})
}

// Setup user-related routes
func (sv *Server) addUserRoutes(rg *gin.RouterGroup) {
	users := rg.Group("users", authenticateMiddleware(sv.tokenGenerator))
	{
		users.GET("me", sv.getCurrentUser)
		users.PATCH("me", sv.updateUser)
		users.POST("send-verify-email", sv.sendVerifyEmail)
		userAddresses := users.Group("addresses")
		{
			userAddresses.POST("", sv.createAddress)
			userAddresses.PATCH(":id/default", sv.setDefaultAddress)
			userAddresses.GET("", sv.getAddresses)
			userAddresses.PATCH(":id", sv.updateAddress)
			userAddresses.DELETE(":id", sv.removeAddress)
		}
	}
}
