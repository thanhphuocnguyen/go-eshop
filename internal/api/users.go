package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/worker"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
	"github.com/thanhphuocnguyen/go-eshop/pkg/cachesrv"
)

// updateUserHandler godoc
// @Summary Update user info
// @Description Update user info
// @Tags users
// @Accept  json
// @Produce  json
// @Param input body UpdateUserRequest true "User info"
// @Success 200 {object} ApiResponse[repository.UpdateUserRow]
// @Failure 400 {object} ApiResponse[gin.H]
// @Failure 401 {object} ApiResponse[gin.H]
// @Failure 500 {object} ApiResponse[gin.H]
// @Router /users/{id} [patch]
func (sv *Server) updateUserHandler(c *gin.Context) {
	authPayload, _ := c.MustGet(AuthPayLoad).(*auth.Payload)
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[repository.UpdateUserRow](InvalidEmailCode, "", err))
		return
	}

	user, err := sv.repo.GetUserByID(c, req.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, createErrorResponse[repository.UpdateUserRow](UnauthorizedCode, "", err))
		return
	}

	if authPayload.RoleCode != repository.UserRoleCodeAdmin && user.ID != req.UserID {
		c.JSON(http.StatusUnauthorized, createErrorResponse[repository.UpdateUserRow](UnauthorizedCode, "", err))
		return
	}

	arg := repository.UpdateUserParams{
		ID: req.UserID,
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
		c.JSON(http.StatusInternalServerError, createErrorResponse[repository.UpdateUserRow](InternalServerErrorCode, "", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, updatedUser, "", nil, nil))
}

// getCurrentUserHandler godoc
// @Summary Get user info
// @Description Get user info
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} ApiResponse[UserResponse]
// @Failure 404 {object} ApiResponse[gin.H]
// @Failure 500 {object} ApiResponse[gin.H]
// @Router /users/me [get]
func (sv *Server) getCurrentUserHandler(c *gin.Context) {
	authPayload, ok := c.MustGet(AuthPayLoad).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, createErrorResponse[UserResponse](InternalServerErrorCode, "", errors.New("authorization payload is not provided")))
		return
	}

	var userResp UserResponse
	err := sv.cachesrv.Get(c, cachesrv.USER_KEY_PREFIX+authPayload.UserID.String(), &userResp)
	if err == nil {
		c.JSON(http.StatusOK, createSuccessResponse(c, userResp, "", nil, nil))
		return
	}

	user, err := sv.repo.GetUserByID(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[UserResponse](NotFoundCode, "", err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[UserResponse](InternalServerErrorCode, "", err))
		return
	}

	userAddress, err := sv.repo.GetAddresses(c, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[UserResponse](InternalServerErrorCode, "", err))
		return
	}

	addressResp := make([]AddressResponse, 0)
	for _, address := range userAddress {
		addressResp = append(addressResp, mapAddressToAddressResponse(address))
	}
	userResp = mapToUserResponse(user)
	userResp.Addresses = addressResp

	if err = sv.cachesrv.Set(c, cachesrv.USER_KEY_PREFIX+authPayload.UserID.String(), userResp, &cachesrv.DEFAULT_EXPIRATION); err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[UserResponse](InternalServerErrorCode, "", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, userResp, "", nil, nil))
}

// getUsersHandler godoc
// @Summary List users
// @Description List users
// @Tags users
// @Accept  json
// @Produce  json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} ApiResponse[[]UserResponse]
// @Failure 500 {object} ApiResponse[gin.H]
// @Failure 400 {object} ApiResponse[gin.H]
// @Failure 401 {object} ApiResponse[gin.H]
// @Router /admin/users [get]
func (sv *Server) getUsersHandler(c *gin.Context) {
	var queries PaginationQueryParams
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[[]UserResponse](InvalidBodyCode, "", err))
		return
	}

	users, err := sv.repo.GetUsers(c, repository.GetUsersParams{
		Limit:  queries.PageSize,
		Offset: (queries.Page - 1) * queries.PageSize,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[[]UserResponse](InternalServerErrorCode, "", err))
		return
	}

	total, err := sv.repo.CountUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[[]UserResponse](InternalServerErrorCode, "", err))
		return
	}

	userResp := make([]UserResponse, 0)
	for _, user := range users {
		userResp = append(userResp, mapToUserResponse(user))
	}

	pagination := createPagination(queries.Page, queries.PageSize, total)
	c.JSON(http.StatusOK, createSuccessResponse(c, userResp, "", pagination, nil))
}

// getUserHandler godoc
// @Summary Get user info
// @Description Get user info
// @Tags Admin
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} ApiResponse[UserResponse]
// @Failure 400 {object} ApiResponse[UserResponse]
// @Failure 404 {object} ApiResponse[UserResponse]
// @Failure 500 {object} ApiResponse[UserResponse]
// @Router /admin/users/{id} [get]
func (sv *Server) getUserHandler(c *gin.Context) {
	var param UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[UserResponse](InvalidBodyCode, "", err))
		return
	}

	user, err := sv.repo.GetUserByID(c, uuid.MustParse(param.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[UserResponse](NotFoundCode, "", err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[UserResponse](InternalServerErrorCode, "", err))
		return
	}

	userResp := mapToUserResponse(user)
	c.JSON(http.StatusOK, createSuccessResponse(c, userResp, "", nil, nil))
}

// sendVerifyEmailHandler godoc
// @Summary Send verify email
// @Description Send verify email
// @Tags users
// @Accept  json
// @Produce  json
// @Success 204 {object} ApiResponse[gin.H]
// @Failure 400 {object} ApiResponse[gin.H]
// @Failure 401 {object} ApiResponse[gin.H]
// @Failure 500 {object} ApiResponse[gin.H]
// @Router /users/verify-email [post]
// @Security BearerAuth
func (sv *Server) sendVerifyEmailHandler(c *gin.Context) {
	authPayload, ok := c.MustGet(AuthPayLoad).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "", errors.New("authorization payload is not provided")))
		return
	}
	user, err := sv.repo.GetUserByID(c, authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "", err))
		return
	}
	if user.VerifiedEmail {
		c.JSON(http.StatusBadRequest, createErrorResponse[gin.H](InvalidEmailCode, "email already verified", fmt.Errorf("email already verified")))
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
		c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "", err))
		return
	}

	c.Status(http.StatusNoContent)
}

// verifyEmailHandler godoc
// @Summary Verify email
// @Description Verify email
// @Tags users
// @Accept  json
// @Produce  json
// @Param id query int true "ID"
// @Param verify_code query string true "Verify code"
// @Success 200 {object} ApiResponse[gin.H]
// @Failure 400 {object} ApiResponse[gin.H]
// @Failure 401 {object} ApiResponse[gin.H]
// @Failure 404 {object} ApiResponse[gin.H]
// @Failure 500 {object} ApiResponse[gin.H]
// @Router /users/verify-email [get]
// @Security BearerAuth
func (sv *Server) verifyEmailHandler(c *gin.Context) {
	var query VerifyEmailQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[bool](InvalidEmailCode, "", err))
		return
	}
	verifyEmail, err := sv.repo.GetVerifyEmailByVerifyCode(c, query.VerifyCode)
	if err != nil {
		c.JSON(http.StatusNotFound, createErrorResponse[bool](NotFoundCode, "", err))
		return
	}

	// Create a transaction to ensure both operations succeed or fail together
	err = sv.repo.VerifyEmailTx(c, repository.VerifyEmailTxArgs{
		VerifyEmail: verifyEmail,
		VerifyCode:  query.VerifyCode,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
		return
	}

	user, err := sv.repo.GetUserByID(c, verifyEmail.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
		return
	}
	if err := sv.cachesrv.Set(c, cachesrv.USER_KEY_PREFIX+user.ID.String(), mapToUserResponse(user), &cachesrv.DEFAULT_EXPIRATION); err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
		return
	}

	// Render HTML success page
	c.Header("Content-Type", "text/html")
	c.HTML(http.StatusOK, "verification-success.html", gin.H{
		"username": user.Username,
		"email":    user.Email,
	})
}
