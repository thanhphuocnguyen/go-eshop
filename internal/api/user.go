package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	repository "github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/worker"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
	"github.com/thanhphuocnguyen/go-eshop/pkg/cache"
)

// ------------------------------ Structs ------------------------------
type UserResponse struct {
	ID                uuid.UUID           `json:"id"`
	Role              repository.UserRole `json:"role"`
	Username          string              `json:"username"`
	FullName          string              `json:"fullname"`
	Email             string              `json:"email,omitempty"`
	Phone             string              `json:"phone,omitempty"`
	VerifiedEmail     bool                `json:"verified_email,omitempty"`
	VerifiedPhone     bool                `json:"verified_phone,omitempty"`
	PasswordChangedAt string              `json:"password_changed_at,omitempty"`
	Addresses         []AddressResponse   `json:"addresses"`
	CreatedAt         string              `json:"created_at,omitempty"`
	UpdatedAt         string              `json:"updated_at,omitempty"`
}

type UpdateUserRequest struct {
	UserID   uuid.UUID `json:"user_id" binding:"required,uuid"`
	FullName *string   `json:"fullname,omitempty" binding:"omitempty,min=3,max=32"`
	Email    *string   `json:"email" binding:"email,max=255,min=6"`
	Phone    *string   `json:"phone" binding:"omitempty,min=8,max=15"`
}

type VerifyEmailQuery struct {
	VerifyCode string `form:"verify_code" binding:"required,min=1"`
}

// ------------------------------ Mappers ------------------------------

func mapToUserResponse(user repository.User) UserResponse {
	return UserResponse{
		ID:                user.ID,
		Addresses:         []AddressResponse{},
		Email:             user.Email,
		FullName:          user.Fullname,
		Role:              user.Role,
		Phone:             user.Phone,
		Username:          user.Username,
		VerifiedEmail:     user.VerifiedEmail,
		VerifiedPhone:     user.VerifiedPhone,
		CreatedAt:         user.CreatedAt.String(),
		UpdatedAt:         user.UpdatedAt.String(),
		PasswordChangedAt: user.PasswordChangedAt.String(),
	}
}

func mapAddressToAddressResponse(address repository.UserAddress) AddressResponse {
	return AddressResponse{
		Phone:    address.Phone,
		Street:   address.Street,
		Ward:     address.Ward,
		District: address.District,
		City:     address.City,
		Default:  address.Default,
		ID:       address.ID,
	}
}

// ------------------------------ Handlers ------------------------------

// updateUserHandler godoc
// @Summary Update user info
// @Description Update user info
// @Tags users
// @Accept  json
// @Produce  json
// @Param input body UpdateUserRequest true "User info"
// @Success 200 {object} ApiResponse[repository.UpdateUserRow]
// @Failure 400 {object} ApiResponse[repository.UpdateUserRow]
// @Failure 401 {object} ApiResponse[repository.UpdateUserRow]
// @Failure 500 {object} ApiResponse[repository.UpdateUserRow]
// @Router /users/{id} [patch]
func (sv *Server) updateUserHandler(c *gin.Context) {
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

	if user.Role != repository.UserRoleAdmin && user.ID != req.UserID {
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

	if req.FullName != nil {
		arg.Fullname = req.FullName
	}

	if req.Phone != nil {
		arg.Phone = req.Phone
		if user.Phone != *req.Phone {
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

// getUserHandler godoc
// @Summary Get user info
// @Description Get user info
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} ApiResponse[UserResponse]
// @Failure 404 {object} ApiResponse[UserResponse]
// @Failure 500 {object} ApiResponse[UserResponse]
// @Router /users [get]
func (sv *Server) getUserHandler(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, createErrorResponse[UserResponse](InternalServerErrorCode, "", errors.New("authorization payload is not provided")))
		return
	}
	var userResp UserResponse
	err := sv.cacheService.Get(c, cache.USER_KEY_PREFIX+authPayload.UserID.String(), &userResp)
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
	sv.cacheService.Set(c, cache.USER_KEY_PREFIX+authPayload.UserID.String(), userResp, &cache.DEFAULT_EXPIRATION)

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
// @Failure 500 {object} ApiResponse[[]UserResponse]
// @Router /users/list [get]
func (sv *Server) getUsersHandler(c *gin.Context) {
	var queries PaginationQueryParams
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[[]UserResponse](InvalidBodyCode, "", err))
		return
	}

	users, err := sv.repo.ListUsers(c, repository.ListUsersParams{
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

	c.JSON(http.StatusOK, createSuccessResponse(c, userResp, "", &Pagination{
		Total:           int64(len(userResp)),
		Page:            queries.Page,
		PageSize:        queries.PageSize,
		TotalPages:      total / queries.PageSize,
		HasNextPage:     len(userResp) > int(queries.PageSize),
		HasPreviousPage: queries.Page > 1,
	}, nil))
}

func (sv *Server) sendVerifyEmailHandler(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", errors.New("authorization payload is not provided")))
		return
	}
	user, err := sv.repo.GetUserByID(c, authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
		return
	}
	if user.VerifiedEmail {
		c.JSON(http.StatusBadRequest, createErrorResponse[bool](InvalidEmailCode, "email already verified", nil))
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
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
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
// @Success 200 {object} ApiResponse[bool]
// @Failure 400 {object} ApiResponse[bool]
// @Failure 401 {object} ApiResponse[bool]
// @Failure 404 {object} ApiResponse[bool]
// @Failure 500 {object} ApiResponse[bool]
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

	// Render HTML success page
	c.Header("Content-Type", "text/html")
	c.HTML(http.StatusOK, "verification-success.html", gin.H{
		"username": user.Username,
		"email":    user.Email,
	})
}
