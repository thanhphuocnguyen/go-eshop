package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/thanhphuocnguyen/go-eshop/internal/auth"
	repository "github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"github.com/thanhphuocnguyen/go-eshop/internal/worker"
)

type UserResponse struct {
	ID                uuid.UUID           `json:"id"`
	Username          string              `json:"username"`
	Role              repository.UserRole `json:"role"`
	FullName          string              `json:"fullname"`
	Email             string              `json:"email,omitempty"`
	CreatedAt         string              `json:"created_at,omitempty"`
	VerifiedEmail     bool                `json:"verified_email,omitempty"`
	VerifiedPhone     bool                `json:"verified_phone,omitempty"`
	UpdatedAt         string              `json:"updated_at,omitempty"`
	PasswordChangedAt string              `json:"password_changed_at,omitempty"`
	Addresses         []addressResponse   `json:"addresses,omitempty"`
}
type addressResponse struct {
	Address  string `json:"address"`
	Address2 string `json:"address_2"`
	City     string `json:"city"`
	District string `json:"district"`
	Ward     string `json:"ward"`
	Phone    string `json:"phone"`
}

type ListUserParams struct {
	Page     int `form:"page" binding:"required,min=1"`
	PageSize int `form:"page_size" binding:"required,min=1,max=100"`
}

type UpdateUserRequest struct {
	UserID   uuid.UUID           `json:"user_id" binding:"required,uuid"`
	FullName *string             `json:"fullname,omitempty" binding:"omitempty,min=3,max=32,alphanum"`
	Email    string              `json:"email" binding:"email,max=255,min=6"`
	Role     repository.UserRole `json:"role"`
}

type VerifyEmailQuery struct {
	ID         int32  `form:"id" binding:"required,min=1"`
	VerifyCode string `form:"verify_code" binding:"required,min=1"`
}

// ------------------------------ Mappers ------------------------------

func mapToUserResponse(user repository.User) UserResponse {
	return UserResponse{
		Email:             user.Email,
		FullName:          user.Fullname,
		Role:              user.Role,
		Username:          user.Username,
		VerifiedEmail:     user.VerifiedEmail,
		VerifiedPhone:     user.VerifiedPhone,
		CreatedAt:         user.CreatedAt.String(),
		UpdatedAt:         user.UpdatedAt.String(),
		PasswordChangedAt: user.PasswordChangedAt.String(),
	}
}

func mapAddressToAddressResponse(address repository.UserAddress) addressResponse {
	return addressResponse{
		Address:  address.Street,
		City:     address.City,
		District: address.District,
		Ward:     address.Ward.String,
		Phone:    address.Phone,
	}
}

// ------------------------------ Handlers ------------------------------

// updateUser godoc
// @Summary Update user info
// @Description Update user info
// @Tags users
// @Accept  json
// @Produce  json
// @Param input body updateUserRequest true "User info"
// @Success 200 {object} ApiResponse
// @Failure 400 {object} ApiResponse
// @Failure 401 {object} ApiResponse
// @Failure 500 {object} ApiResponse
// @Router /users/{id} [patch]
func (sv *Server) updateUser(c *gin.Context) {
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	user, err := sv.repo.GetUserByID(c, req.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, createErrorResponse(http.StatusUnauthorized, "", err))
		return
	}

	if user.Role != repository.UserRoleAdmin && user.ID != req.UserID {
		c.JSON(http.StatusUnauthorized, createErrorResponse(http.StatusUnauthorized, "", err))
		return
	}

	arg := repository.UpdateUserParams{
		ID: req.UserID,
		Email: pgtype.Text{
			String: req.Email,
			Valid:  true,
		},
	}

	if req.FullName != nil {
		arg.Fullname = utils.GetPgTypeText(*req.FullName)
	}

	if user.Role == repository.UserRoleAdmin {
		arg.Role = repository.NullUserRole{
			UserRole: req.Role,
		}
	}
	updatedUser, err := sv.repo.UpdateUser(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, updatedUser, "", nil, nil))
}

// getUser godoc
// @Summary Get user info
// @Description Get user info
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} ApiResponse
// @Failure 404 {object} ApiResponse
// @Failure 500 {object} ApiResponse
// @Router /users [get]
func (sv *Server) getUser(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", errors.New("authorization payload is not provided")))
		return
	}

	user, err := sv.repo.GetUserByID(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse(http.StatusBadRequest, "", err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	userAddress, err := sv.repo.GetAddresses(c, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	addressResp := make([]addressResponse, 0)
	for _, address := range userAddress {
		addressResp = append(addressResp, mapAddressToAddressResponse(address))
	}
	userResp := mapToUserResponse(user)
	userResp.Addresses = addressResp

	c.JSON(http.StatusOK, createSuccessResponse(c, userResp, "", nil, nil))
}

// listUsers godoc
// @Summary List users
// @Description List users
// @Tags users
// @Accept  json
// @Produce  json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} ApiResponse
// @Failure 500 {object} ApiResponse
// @Router /users/list [get]
func (sv *Server) listUsers(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", errors.New("authorization payload is not provided")))
		return
	}
	user, err := sv.repo.GetUserByID(c, authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	if user.Role != repository.UserRoleAdmin {
		c.JSON(http.StatusUnauthorized, createErrorResponse(http.StatusUnauthorized, "", errors.New("user is not admin")))
		return
	}

	var queries ListUserParams
	if err := c.ShouldBindUri(&queries); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	users, err := sv.repo.ListUsers(c, repository.ListUsersParams{
		Limit:  int32(queries.PageSize),
		Offset: int32((queries.Page - 1) * queries.PageSize),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	userResp := make([]UserResponse, 0)
	for _, user := range users {
		userAddress, err := sv.repo.GetAddresses(c, user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
			return
		}

		addressResp := make([]addressResponse, 0)
		for _, address := range userAddress {
			addressResp = append(addressResp, mapAddressToAddressResponse(address))
		}
		userResp = append(userResp, mapToUserResponse(user))
		userResp[len(userResp)-1].Addresses = addressResp
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, userResp, "", nil, nil))
}

func (sv *Server) sendVerifyEmail(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", errors.New("authorization payload is not provided")))
		return
	}

	err := sv.taskDistributor.SendVerifyEmail(
		c,
		&worker.PayloadVerifyEmail{
			UserID: authPayload.UserID,
		},
		asynq.MaxRetry(3),
		asynq.ProcessIn(5*time.Second),
		asynq.Queue(worker.QueueCritical),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (sv *Server) verifyEmail(c *gin.Context) {
	var query VerifyEmailQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}
	verifyEmail, err := sv.repo.GetVerifyEmailByID(c, query.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, createErrorResponse(http.StatusNotFound, "", err))
		return
	}

	if verifyEmail.ExpiredAt.Before(time.Now()) {
		c.JSON(http.StatusNotFound, createErrorResponse(http.StatusNotFound, "", errors.New("verify code expired")))
		return
	}

	if verifyEmail.VerifyCode != query.VerifyCode {
		c.JSON(http.StatusNotFound, createErrorResponse(http.StatusNotFound, "", err))
		return
	}
	_, err = sv.repo.UpdateVerifyEmail(c, repository.UpdateVerifyEmailParams{
		ID:         verifyEmail.ID,
		VerifyCode: query.VerifyCode,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	_, err = sv.repo.UpdateUser(c, repository.UpdateUserParams{
		ID: verifyEmail.UserID,
		VerifiedEmail: pgtype.Bool{
			Bool:  true,
			Valid: true,
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}
	c.Status(http.StatusNoContent)
}
