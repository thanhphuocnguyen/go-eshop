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

// UpdateUserHandler godoc
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
func (sv *Server) updateUserHandler(c *gin.Context) {
	authPayload, _ := c.MustGet(AuthPayLoad).(*auth.Payload)
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidEmailCode, err))
		return
	}

	user, err := sv.repo.GetUserByID(c, req.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, createErr(UnauthorizedCode, err))
		return
	}

	if authPayload.RoleCode != repository.UserRoleCodeAdmin && user.ID != req.UserID {
		c.JSON(http.StatusUnauthorized, createErr(UnauthorizedCode, err))
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
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusOK, createDataResp(c, updatedUser, nil, nil))
}

// GetCurrentUserHandler godoc
// @Summary Get user info
// @Description Get user info
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} ApiResponse[UserDetail]
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /users/me [get]
func (sv *Server) GetCurrentUserHandler(c *gin.Context) {
	authPayload, ok := c.MustGet(AuthPayLoad).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, errors.New("authorization payload is not provided")))
		return
	}

	var userResp UserDetail
	err := sv.cachesrv.Get(c, cachesrv.USER_KEY_PREFIX+authPayload.UserID.String(), &userResp)
	if err == nil {
		c.JSON(http.StatusOK, createDataResp(c, userResp, nil, nil))
		return
	}

	user, err := sv.repo.GetUserByID(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErr(NotFoundCode, err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	userAddress, err := sv.repo.GetAddresses(c, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	addressResp := make([]AddressResponse, 0)
	for _, address := range userAddress {
		addressResp = append(addressResp, mapAddressToAddressResponse(address))
	}
	userResp = mapToUserResponse(user, authPayload.RoleCode)
	userResp.Addresses = addressResp

	if err = sv.cachesrv.Set(c, cachesrv.USER_KEY_PREFIX+authPayload.UserID.String(), userResp, &cachesrv.DEFAULT_EXPIRATION); err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusOK, createDataResp(c, userResp, nil, nil))
}

// GetUsersHandler godoc
// @Summary List users
// @Description List users
// @Tags users
// @Accept  json
// @Produce  json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} ApiResponse[[]UserResponse]
// @Failure 500 {object} ErrorResp
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Router /admin/users [get]
func (sv *Server) GetUsersHandler(c *gin.Context) {
	authPayload, ok := c.MustGet(AuthPayLoad).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, errors.New("authorization payload is not provided")))
		return
	}
	var queries PaginationQueryParams
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	users, err := sv.repo.GetUsers(c, repository.GetUsersParams{
		Limit:  queries.PageSize,
		Offset: (queries.Page - 1) * queries.PageSize,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	total, err := sv.repo.CountUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	userResp := make([]UserDetail, 0)
	for _, user := range users {
		userResp = append(userResp, mapToUserResponse(user, authPayload.RoleCode))
	}

	pagination := createPagination(queries.Page, queries.PageSize, total)
	c.JSON(http.StatusOK, createDataResp(c, userResp, pagination, nil))
}

// GetUserHandler godoc
// @Summary Get user info
// @Description Get user info
// @Tags Admin
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} ApiResponse[UserDetail]
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/users/{id} [get]
func (sv *Server) GetUserHandler(c *gin.Context) {
	authPayload, ok := c.MustGet(AuthPayLoad).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, errors.New("authorization payload is not provided")))
		return
	}
	var param UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	user, err := sv.repo.GetUserByID(c, uuid.MustParse(param.ID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErr(NotFoundCode, err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	userResp := mapToUserResponse(user, authPayload.RoleCode)
	c.JSON(http.StatusOK, createDataResp(c, userResp, nil, nil))
}

// SendVerifyEmailHandler godoc
// @Summary Send verify email
// @Description Send verify email
// @Tags users
// @Accept  json
// @Produce  json
// @Success 204 {object}
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /users/verify-email [post]
// @Security BearerAuth
func (sv *Server) SendVerifyEmailHandler(c *gin.Context) {
	authPayload, ok := c.MustGet(AuthPayLoad).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, errors.New("authorization payload is not provided")))
		return
	}
	user, err := sv.repo.GetUserByID(c, authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}
	if user.VerifiedEmail {
		c.JSON(http.StatusBadRequest, createErr(InvalidEmailCode, fmt.Errorf("email already verified")))
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
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	c.Status(http.StatusNoContent)
}

// VerifyEmailHandler godoc
// @Summary Verify email
// @Description Verify email
// @Tags users
// @Accept  json
// @Produce  json
// @Param id query int true "ID"
// @Param verify_code query string true "Verify code"
// @Success 200 {object} html
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /users/verify-email [get]
// @Security BearerAuth
func (sv *Server) VerifyEmailHandler(c *gin.Context) {
	authPayload, ok := c.MustGet(AuthPayLoad).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, errors.New("authorization payload is not provided")))
		return
	}
	var query VerifyEmailQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidEmailCode, err))
		return
	}
	verifyEmail, err := sv.repo.GetVerifyEmailByVerifyCode(c, query.VerifyCode)
	if err != nil {
		c.JSON(http.StatusNotFound, createErr(NotFoundCode, err))
		return
	}

	// Create a transaction to ensure both operations succeed or fail together
	err = sv.repo.VerifyEmailTx(c, repository.VerifyEmailTxArgs{
		VerifyEmail: verifyEmail,
		VerifyCode:  query.VerifyCode,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	user, err := sv.repo.GetUserByID(c, verifyEmail.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}
	if err := sv.cachesrv.Set(c, cachesrv.USER_KEY_PREFIX+user.ID.String(), mapToUserResponse(user, authPayload.RoleCode), &cachesrv.DEFAULT_EXPIRATION); err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	// Render HTML success page
	c.Header("Content-Type", "text/html")
	c.HTML(http.StatusOK, "verification-success.html", gin.H{
		"username": user.Username,
		"email":    user.Email,
	})
}
