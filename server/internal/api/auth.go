package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/thanhphuocnguyen/go-eshop/internal/auth"
	repository "github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/worker"
)

type RegisterRequestBody struct {
	Username string `json:"username" binding:"required,min=3,max=32,alphanum,lowercase"`
	Password string `json:"password" binding:"required,min=6,max=32,alphanum"`
	FullName string `json:"fullname" binding:"required,min=3,max=32"`
	Phone    string `json:"phone" binding:"required,min=10,max=15"`
	Email    string `json:"email" binding:"required,email,max=255,min=6"`
}

type LoginRequestBody struct {
	Username string `json:"username" binding:"required,min=3,max=32,alphanum"`
	Password string `json:"password" binding:"required,min=6,max=32,alphanum"`
}

type LoginResponse struct {
	ID                   uuid.UUID    `json:"session_id"`
	Token                string       `json:"token"`
	TokenExpireAt        time.Time    `json:"token_expire_at"`
	RefreshToken         string       `json:"refresh_token"`
	RefreshTokenExpireAt time.Time    `json:"refresh_token_expire_at"`
	User                 UserResponse `json:"user"`
}

type RefreshTokenResponse struct {
	AccessToken          string        `json:"access_token"`
	AccessTokenExpiresAt time.Duration `json:"access_token_expires_at"`
}

// ------------------------------ Handlers ------------------------------

// register godoc
// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept  json
// @Produce  json
// @Param input body createUserRequest true "User info"
// @Success 200 {object} GenericResponse[repository.CreateUserRow]
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /users [post]
func (sv *Server) register(c *gin.Context) {
	var req RegisterRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	_, err := sv.repo.GetUserByUsername(c, req.Username)
	if err != nil && !errors.Is(err, repository.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	if err == nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", fmt.Errorf("username %s is already taken", req.Username)))
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	arg := repository.CreateUserParams{
		ID:             uuid.New(),
		Username:       req.Username,
		HashedPassword: hashedPassword,
		Fullname:       req.FullName,
		Email:          req.Email,
		Phone:          req.Phone,
		Role:           repository.UserRoleUser,
	}

	if req.Username == "admin" {
		arg.Role = repository.UserRoleAdmin
	}
	user, err := sv.repo.CreateUser(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	err = sv.taskDistributor.SendVerifyEmail(
		c,
		&worker.PayloadVerifyEmail{
			UserID: user.ID,
		},
		asynq.MaxRetry(3),
		asynq.ProcessIn(5*time.Second),
		asynq.Queue(worker.QueueCritical),
	)

	message := "Please verify your email address to activate your account"
	if err != nil {
		createErrorResponse(http.StatusInternalServerError, "", err)
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, user, message, nil, nil))
}

// login godoc
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
func (sv *Server) login(c *gin.Context) {
	var req LoginRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", err))
		return
	}

	user, err := sv.repo.GetUserByUsername(c, req.Username)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, createErrorResponse(http.StatusUnauthorized, "", err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}

	// TODO: check if user has confirmed email or phone number

	if err := auth.CheckPassword(req.Password, user.HashedPassword); err != nil {
		c.JSON(http.StatusUnauthorized, createErrorResponse(http.StatusUnauthorized, "", err))
		return
	}

	token, payload, err := sv.tokenGenerator.GenerateToken(user.ID, user.Username, user.Email, sv.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}

	refreshToken, rfPayload, err := sv.tokenGenerator.GenerateToken(user.ID, user.Username, user.Email, sv.config.RefreshTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}

	session, err := sv.repo.CreateSession(c, repository.CreateSessionParams{
		ID:           rfPayload.ID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    c.GetHeader("User-Agent"),
		ClientIp:     c.ClientIP(),
		Blocked:      false,
		ExpiredAt:    time.Now().Add(sv.config.RefreshTokenDuration),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}

	userResp := UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
		FullName: user.Fullname,
	}

	loginResp := LoginResponse{
		ID:                   session.ID,
		TokenExpireAt:        payload.ExpiredAt,
		Token:                token,
		RefreshToken:         refreshToken,
		RefreshTokenExpireAt: rfPayload.ExpiredAt,
	}
	loginResp.User = userResp

	c.JSON(http.StatusOK, createSuccessResponse(c, loginResp, "success", nil, nil))
}

// refreshToken godoc
// @Summary Refresh token
// @Description Refresh token
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} renewAccessTokenResp
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /users/refresh-token [post]
func (sv *Server) refreshToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, createErrorResponse(http.StatusUnauthorized, "", fmt.Errorf("refresh token is required")))
		return
	}
	refreshToken := authHeader[len("Bearer "):]
	refreshTokenPayload, err := sv.tokenGenerator.VerifyToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, createErrorResponse(http.StatusUnauthorized, "", err))
		return
	}

	session, err := sv.repo.GetSession(c, refreshTokenPayload.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound,
				createErrorResponse(http.StatusNotFound, "", fmt.Errorf("session not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}

	if session.ID != refreshTokenPayload.ID {
		err := errors.New("refresh token is not valid")
		c.JSON(http.StatusUnauthorized, createErrorResponse(http.StatusUnauthorized, "", err))
		return
	}

	if session.RefreshToken != refreshToken {
		err := errors.New("refresh token is not valid")
		c.JSON(http.StatusUnauthorized, createErrorResponse(http.StatusUnauthorized, "", err))
		return
	}

	if session.Blocked {
		err := errors.New("session is blocked")
		c.JSON(http.StatusUnauthorized, createErrorResponse(http.StatusUnauthorized, "", err))
		return
	}

	if time.Now().After(session.ExpiredAt) {
		err := errors.New("refresh token was expired")
		c.JSON(http.StatusUnauthorized, createErrorResponse(http.StatusUnauthorized, "", err))
		return
	}
	accessToken, _, err := sv.tokenGenerator.GenerateToken(session.UserID, refreshTokenPayload.Username, refreshTokenPayload.Email, sv.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, RefreshTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: sv.config.AccessTokenDuration,
	}, "success", nil, nil))
}
