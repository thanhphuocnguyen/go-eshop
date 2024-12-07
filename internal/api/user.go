package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/thanhphuocnguyen/go-eshop/internal/auth"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/sqlc"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	FullName string `json:"full_name" binding:"required,min=3,max=32"`
	Email    string `json:"email" binding:"required,email,max=255"`
}

// createUser godoc
// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept  json
// @Produce  json
// @Param input body createUserRequest true "User info"
// @Success 200 {object} createUserResponse
// @Router /users [post]
func (sv *Server) createUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := sqlc.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}
	user, err := sv.postgres.CreateUser(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, user)
}

type createUserResponse struct {
	Email             string `json:"email"`
	FullName          string `json:"full_name"`
	Username          string `json:"username"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
	PasswordChangedAt string `json:"password_changed_at"`
}
type loginUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
}

type loginResponse struct {
	SessionID            uuid.UUID          `json:"session_id"`
	Token                string             `json:"token"`
	TokenExpireAt        time.Time          `json:"token_expire_at"`
	RefreshToken         string             `json:"refresh_token"`
	RefreshTokenExpireAt time.Time          `json:"refresh_token_expire_at"`
	User                 createUserResponse `json:"user"`
}

// loginUser godoc
// @Summary Login to the system
// @Description Login to the system
// @Tags users
// @Accept  json
// @Produce  json
// @Param input body loginUserRequest true "User info"
// @Success 200 {object} loginResponse
// @Router /users/login [post]
func (sv *Server) loginUser(c *gin.Context) {
	var req loginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := sv.postgres.GetUserByUsername(c, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err := auth.CheckPassword(req.Password, user.HashedPassword); err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	token, payload, err := sv.tokenGenerator.GenerateToken(user.ID, user.Username, user.Role, sv.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, rfPayload, err := sv.tokenGenerator.GenerateToken(user.ID, user.Username, user.Role, sv.config.RefreshTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
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
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, loginResponse{
		SessionID:            session.ID,
		TokenExpireAt:        payload.ExpiredAt,
		Token:                token,
		RefreshToken:         refreshToken,
		RefreshTokenExpireAt: rfPayload.ExpiredAt,
		User:                 mapToUserResponse(user),
	})
}

type updateUserRequest struct {
	UserID   int64         `json:"user_id" binding:"required"`
	FullName string        `json:"full_name"`
	Email    string        `json:"email" binding:"email"`
	Role     sqlc.UserRole `json:"role"`
}

// updateUser godoc
// @Summary Update user info
// @Description Update user info
// @Tags users
// @Accept  json
// @Produce  json
// @Param input body updateUserRequest true "User info"
// @Success 200 {object} createUserResponse
// @Router /users/{id} [patch]
func (sv *Server) updateUser(c *gin.Context) {
	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := sv.postgres.GetUserByID(c, 1)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if user.Role != sqlc.UserRoleAdmin && user.ID != req.UserID {
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := sqlc.UpdateUserParams{
		ID: req.UserID,
		Email: pgtype.Text{
			String: req.Email,
			Valid:  true,
		},
		Role: sqlc.NullUserRole{
			UserRole: req.Role,
		},
		FullName: pgtype.Text{
			String: req.FullName,
			Valid:  true,
		},
	}
	updatedUser, err := sv.postgres.UpdateUser(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func mapToUserResponse(user sqlc.User) createUserResponse {
	return createUserResponse{
		Email:             user.Email,
		FullName:          user.FullName,
		Username:          user.Username,
		CreatedAt:         user.CreatedAt.String(),
		UpdatedAt:         user.UpdatedAt.String(),
		PasswordChangedAt: user.PasswordChangedAt.String(),
	}
}
