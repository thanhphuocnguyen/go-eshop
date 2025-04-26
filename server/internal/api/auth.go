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
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"github.com/thanhphuocnguyen/go-eshop/internal/worker"
)

type RegisterRequestBody struct {
	Username string            `json:"username" binding:"required,min=3,max=32,lowercase"`
	Password string            `json:"password" binding:"required,min=6,max=32"`
	FullName string            `json:"fullname" binding:"required,min=3,max=32"`
	Phone    string            `json:"phone" binding:"required,min=10,max=15"`
	Email    string            `json:"email" binding:"required,email,max=255,min=6"`
	Address  *CreateAddressReq `json:"address" binding:"omitempty,required"`
}

type LoginRequest struct {
	Username *string `json:"username" binding:"omitempty,min=3,max=32"`
	Email    *string `json:"email" binding:"omitempty,email,max=255,min=6"`
	Password string  `json:"password" binding:"required,min=6,max=32"`
}

type LoginResponse struct {
	ID                    uuid.UUID `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_in"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

type RefreshTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

// ------------------------------ Handlers ------------------------------

// registerHandler godoc
// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept  json
// @Produce  json
// @Param input body RegisterRequestBody true "User info"
// @Success 200 {object} ApiResponse[UserResponse]
// @Failure 400 {object} ApiResponse[UserResponse]
// @Failure 500 {object} ApiResponse[UserResponse]
// @Router /users [post]
func (sv *Server) registerHandler(c *gin.Context) {
	var req RegisterRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[UserResponse](http.StatusBadRequest, "", err))
		return
	}

	_, err := sv.repo.GetUserByUsername(c, req.Username)
	if err != nil && !errors.Is(err, repository.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, createErrorResponse[UserResponse](http.StatusBadRequest, "", err))
		return
	}

	if err == nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[UserResponse](http.StatusBadRequest, "", fmt.Errorf("username %s is already taken", req.Username)))
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[UserResponse](http.StatusBadRequest, "", err))
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
		c.JSON(http.StatusInternalServerError, createErrorResponse[UserResponse](http.StatusBadRequest, "", err))
		return
	}

	createAddressArgs := repository.CreateAddressParams{
		UserID:   user.ID,
		Phone:    req.Address.Phone,
		Street:   req.Address.Street,
		City:     req.Address.City,
		District: req.Address.District,
		Default:  true,
	}
	if req.Address.Ward != nil {
		createAddressArgs.Ward = utils.GetPgTypeText(*req.Address.Ward)
	}
	createdAddress, err := sv.repo.CreateAddress(c, createAddressArgs)

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[UserResponse](http.StatusBadRequest, "", err))
		return
	}

	err = sv.taskDistributor.SendVerifyEmail(c, &worker.PayloadVerifyEmail{UserID: user.ID}, asynq.MaxRetry(3), asynq.ProcessIn(5*time.Second), asynq.Queue(worker.QueueCritical))

	if err != nil {
		createErrorResponse[UserResponse](http.StatusInternalServerError, "Please verify your email address to activate your account", err)
		return
	}

	ward := ""
	if createdAddress.Ward.Valid {
		ward = createdAddress.Ward.String
	}

	userResp := &UserResponse{
		ID:            user.ID,
		Username:      user.Username,
		Role:          user.Role,
		Email:         user.Email,
		CreatedAt:     user.CreatedAt.String(),
		VerifiedEmail: user.VerifiedEmail,
		VerifiedPhone: user.VerifiedPhone,
		UpdatedAt:     user.UpdatedAt.String(),
		FullName:      user.Fullname,
		Addresses: []AddressResponse{{
			ID: createdAddress.ID,
			Address: Address{
				Phone:    createdAddress.Phone,
				Street:   createdAddress.Street,
				Ward:     &ward,
				District: createdAddress.District,
				City:     createdAddress.City,
			},
			Default:   createdAddress.Default,
			CreatedAt: createdAddress.CreatedAt,
		}},
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, userResp, "Created user and address successfully", nil, nil))
}

// loginHandler godoc
// @Summary Login to the system
// @Description Login to the system
// @Tags users
// @Accept  json
// @Produce  json
// @Param input body LoginRequest true "User info"
// @Success 200 {object} ApiResponse[LoginResponse]
// @Failure 401 {object} ApiResponse[LoginResponse]
// @Failure 500 {object} ApiResponse[LoginResponse]
// @Router /users/loginHandler [post]
func (sv *Server) loginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[LoginResponse](http.StatusBadRequest, "", err))
		return
	}
	if req.Username == nil && req.Email == nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[LoginResponse](http.StatusBadRequest, "", fmt.Errorf("username or email is required")))
		return
	}

	var user repository.User
	var err error
	if req.Username != nil {
		user, err = sv.repo.GetUserByUsername(c, *req.Username)
	} else {
		user, err = sv.repo.GetUserByEmail(c, *req.Email)
	}

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, createErrorResponse[LoginResponse](http.StatusUnauthorized, "User not existed", err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[LoginResponse](http.StatusInternalServerError, "Internal Server Error", err))
		return
	}

	if err := auth.CheckPassword(req.Password, user.HashedPassword); err != nil {
		c.JSON(http.StatusUnauthorized, createErrorResponse[LoginResponse](http.StatusUnauthorized, "Invalid credentials", err))
		return
	}

	accessToken, payload, err := sv.tokenGenerator.GenerateToken(user.ID, user.Username, user.Role, sv.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[LoginResponse](http.StatusInternalServerError, "", err))
		return
	}

	refreshToken, rfPayload, err := sv.tokenGenerator.GenerateToken(user.ID, user.Username, user.Role, sv.config.RefreshTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[LoginResponse](http.StatusInternalServerError, "", err))
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
		c.JSON(http.StatusInternalServerError, createErrorResponse[LoginResponse](http.StatusInternalServerError, "", err))
		return
	}

	loginResp := LoginResponse{
		ID:                    session.ID,
		AccessTokenExpiresAt:  payload.Expires,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: rfPayload.Expires,
	}
	c.JSON(http.StatusOK, createSuccessResponse(c, loginResp, "success", nil, nil))
}

// refreshTokenHandler godoc
// @Summary Refresh token
// @Description Refresh token
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} ApiResponse[RefreshTokenResponse]
// @Failure 401 {object} ApiResponse[RefreshTokenResponse]
// @Failure 500 {object} ApiResponse[RefreshTokenResponse]
// @Router /users/refresh-token [post]
func (sv *Server) refreshTokenHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, createErrorResponse[RefreshTokenResponse](http.StatusUnauthorized, "", fmt.Errorf("refresh token is required")))
		return
	}

	refreshToken := authHeader[len("Bearer "):]
	refreshTokenPayload, err := sv.tokenGenerator.VerifyToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, createErrorResponse[RefreshTokenResponse](http.StatusUnauthorized, "", err))
		return
	}

	session, err := sv.repo.GetSession(c, refreshTokenPayload.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound,
				createErrorResponse[RefreshTokenResponse](http.StatusNotFound, "Not found", fmt.Errorf("session not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[RefreshTokenResponse](http.StatusInternalServerError, "", err))
		return
	}

	if session.ID != refreshTokenPayload.ID {
		err := errors.New("refresh token is not valid")
		c.JSON(http.StatusUnauthorized, createErrorResponse[RefreshTokenResponse](http.StatusUnauthorized, "", err))
		return
	}

	if session.RefreshToken != refreshToken {
		err := errors.New("refresh token is not valid")
		c.JSON(http.StatusUnauthorized, createErrorResponse[RefreshTokenResponse](http.StatusUnauthorized, "", err))
		return
	}

	if session.Blocked {
		err := errors.New("session is blocked")
		c.JSON(http.StatusUnauthorized, createErrorResponse[RefreshTokenResponse](http.StatusUnauthorized, "", err))
		return
	}

	if time.Now().After(session.ExpiredAt) {
		err := errors.New("refresh token was expired")
		c.JSON(http.StatusUnauthorized, createErrorResponse[RefreshTokenResponse](http.StatusUnauthorized, "", err))
		return
	}
	accessToken, _, err := sv.tokenGenerator.GenerateToken(session.UserID, refreshTokenPayload.Username, refreshTokenPayload.Role, sv.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[RefreshTokenResponse](http.StatusInternalServerError, "", err))
		return
	}

	c.JSON(http.StatusOK,
		createSuccessResponse(c, RefreshTokenResponse{AccessToken: accessToken, AccessTokenExpiresAt: time.Now().Add(sv.config.AccessTokenDuration)}, "success", nil, nil))
}
