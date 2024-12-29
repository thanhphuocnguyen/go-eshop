package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/thanhphuocnguyen/go-eshop/internal/auth"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/postgres"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/sqlc"
	"github.com/thanhphuocnguyen/go-eshop/internal/util"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32,alphanum,lowercase"`
	Password string `json:"password" binding:"required,min=6,max=32,alphanum"`
	FullName string `json:"full_name" binding:"required,min=3,max=32"`
	Phone    string `json:"phone" binding:"required,min=10,max=15"`
	Email    string `json:"email" binding:"required,email,max=255,min=6"`
}

type userResponse struct {
	Email             string            `json:"email"`
	FullName          string            `json:"full_name"`
	Username          string            `json:"username"`
	CreatedAt         string            `json:"created_at"`
	UpdatedAt         string            `json:"updated_at"`
	PasswordChangedAt string            `json:"password_changed_at"`
	Addresses         []addressResponse `json:"addresses"`
}
type addressResponse struct {
	Address  string `json:"address"`
	Address2 string `json:"address_2"`
	City     string `json:"city"`
	District string `json:"district"`
	Ward     string `json:"ward"`
	Phone    string `json:"phone"`
}
type loginUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32,alphanum"`
	Password string `json:"password" binding:"required,min=6,max=32,alphanum"`
}

type updateUserRequest struct {
	UserID   int64         `json:"user_id" binding:"required,min=1"`
	FullName *string       `json:"full_name,omitempty" binding:"omitempty,min=3,max=32,alphanum"`
	Email    string        `json:"email" binding:"email,max=255,min=6"`
	Role     sqlc.UserRole `json:"role"`
}
type loginResponse struct {
	SessionID            uuid.UUID    `json:"session_id"`
	Token                string       `json:"token"`
	TokenExpireAt        time.Time    `json:"token_expire_at"`
	RefreshToken         string       `json:"refresh_token"`
	RefreshTokenExpireAt time.Time    `json:"refresh_token_expire_at"`
	User                 userResponse `json:"user"`
}

// ------------------------------ Mappers ------------------------------

func mapToUserResponse(user sqlc.User) userResponse {
	return userResponse{
		Email:             user.Email,
		FullName:          user.FullName,
		Username:          user.Username,
		CreatedAt:         user.CreatedAt.String(),
		UpdatedAt:         user.UpdatedAt.String(),
		PasswordChangedAt: user.PasswordChangedAt.String(),
	}
}

func mapAddressToAddressResponse(address sqlc.UserAddress) addressResponse {
	return addressResponse{
		Address:  address.Address1,
		Address2: address.Address2.String,
		City:     address.City,
		District: address.District,
		Ward:     address.Ward.String,
		Phone:    address.Phone,
	}
}

// ------------------------------ Handlers ------------------------------

// createUser godoc
// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept  json
// @Produce  json
// @Param input body createUserRequest true "User info"
// @Success 200 {object} GenericResponse[sqlc.CreateUserRow]
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /users [post]
func (sv *Server) createUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	arg := sqlc.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
		Phone:          req.Phone,
	}
	if req.Username == "admin" {
		arg.Role = sqlc.UserRoleAdmin
	}
	user, err := sv.postgres.CreateUser(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusOK, GenericResponse[sqlc.CreateUserRow]{&user, nil, nil})
}

// loginUser godoc
// @Summary Login to the system
// @Description Login to the system
// @Tags users
// @Accept  json
// @Produce  json
// @Param input body loginUserRequest true "User info"
// @Success 200 {object} GenericResponse[loginResponse]
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /users/login [post]
func (sv *Server) loginUser(c *gin.Context) {
	var req loginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	user, err := sv.postgres.GetUserByUsername(c, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if err := auth.CheckPassword(req.Password, user.HashedPassword); err != nil {
		c.JSON(http.StatusUnauthorized, mapErrResp(err))
		return
	}

	token, payload, err := sv.tokenGenerator.GenerateToken(user.ID, user.Username, user.Role, sv.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	refreshToken, rfPayload, err := sv.tokenGenerator.GenerateToken(user.ID, user.Username, user.Role, sv.config.RefreshTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	sessionID := uuid.New()
	session, err := sv.postgres.CreateSession(c, sqlc.CreateSessionParams{
		ID:           sessionID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    c.GetHeader("User-Agent"),
		ClientIp:     c.ClientIP(),
		IsBlocked:    false,
		ExpiredAt:    time.Now().Add(sv.config.RefreshTokenDuration),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	loginResp := loginResponse{
		SessionID:            session.ID,
		TokenExpireAt:        payload.ExpiredAt,
		Token:                token,
		RefreshToken:         refreshToken,
		RefreshTokenExpireAt: rfPayload.ExpiredAt,
		User:                 mapToUserResponse(user),
	}
	c.JSON(http.StatusOK, GenericResponse[loginResponse]{&loginResp, nil, nil})
}

// updateUser godoc
// @Summary Update user info
// @Description Update user info
// @Tags users
// @Accept  json
// @Produce  json
// @Param input body updateUserRequest true "User info"
// @Success 200 {object} GenericResponse[sqlc.UpdateUserRow]
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /users/{id} [patch]
func (sv *Server) updateUser(c *gin.Context) {
	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	user, err := sv.postgres.GetUserByID(c, 1)
	if err != nil {
		c.JSON(http.StatusUnauthorized, mapErrResp(err))
		return
	}

	if user.Role != sqlc.UserRoleAdmin && user.ID != req.UserID {
		c.JSON(http.StatusUnauthorized, mapErrResp(err))
		return
	}

	arg := sqlc.UpdateUserParams{
		ID: req.UserID,
		Email: pgtype.Text{
			String: req.Email,
			Valid:  true,
		},
	}

	if req.FullName != nil {
		arg.FullName = util.GetPgTypeText(*req.FullName)
	}

	if user.Role == sqlc.UserRoleAdmin {
		arg.Role = sqlc.NullUserRole{
			UserRole: req.Role,
		}
	}
	updatedUser, err := sv.postgres.UpdateUser(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusOK, GenericResponse[sqlc.UpdateUserRow]{&updatedUser, nil, nil})
}

// getUser godoc
// @Summary Get user info
// @Description Get user info
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} GenericResponse[userResponse]
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /users [get]
func (sv *Server) getUser(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, mapErrResp(fmt.Errorf("authorization payload is not provided")))
		return
	}

	user, err := sv.postgres.GetUserByID(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(err))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	userAddress, err := sv.postgres.GetAddresses(c, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	addressResp := make([]addressResponse, 0)
	for _, address := range userAddress {
		addressResp = append(addressResp, mapAddressToAddressResponse(address))
	}
	userResp := mapToUserResponse(user)
	userResp.Addresses = addressResp

	c.JSON(http.StatusOK, GenericResponse[userResponse]{&userResp, nil, nil})
}
