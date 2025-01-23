package api

import (
	"errors"
	"fmt"
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

type createUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32,alphanum,lowercase"`
	Password string `json:"password" binding:"required,min=6,max=32,alphanum"`
	FullName string `json:"fullname" binding:"required,min=3,max=32"`
	Phone    string `json:"phone" binding:"required,min=10,max=15"`
	Email    string `json:"email" binding:"required,email,max=255,min=6"`
}

type userResponse struct {
	Email             string              `json:"email"`
	FullName          string              `json:"fullname"`
	Username          string              `json:"username"`
	CreatedAt         string              `json:"created_at"`
	VerifiedEmail     bool                `json:"verified_email"`
	VerifiedPhone     bool                `json:"verified_phone"`
	Role              repository.UserRole `json:"role"`
	UpdatedAt         string              `json:"updated_at"`
	PasswordChangedAt string              `json:"password_changed_at"`
	Addresses         []addressResponse   `json:"addresses"`
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

type listUserParams struct {
	Page     int `form:"page" binding:"required,min=1"`
	PageSize int `form:"page_size" binding:"required,min=1,max=100"`
}

type updateUserRequest struct {
	UserID   uuid.UUID           `json:"user_id" binding:"required,uuid"`
	FullName *string             `json:"fullname,omitempty" binding:"omitempty,min=3,max=32,alphanum"`
	Email    string              `json:"email" binding:"email,max=255,min=6"`
	Role     repository.UserRole `json:"role"`
}
type loginResponse struct {
	SessionID            uuid.UUID    `json:"session_id"`
	Token                string       `json:"token"`
	TokenExpireAt        time.Time    `json:"token_expire_at"`
	RefreshToken         string       `json:"refresh_token"`
	RefreshTokenExpireAt time.Time    `json:"refresh_token_expire_at"`
	User                 userResponse `json:"user"`
}

type renewAccessTokenResp struct {
	AccessToken          string        `json:"access_token"`
	AccessTokenExpiresAt time.Duration `json:"access_token_expires_at"`
}

type verifyEmailQuery struct {
	ID         int32  `form:"id" binding:"required,min=1"`
	VerifyCode string `form:"verify_code" binding:"required,min=1"`
}

// ------------------------------ Mappers ------------------------------

func mapToUserResponse(user repository.User) userResponse {
	return userResponse{
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

// createUser godoc
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
func (sv *Server) createUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	_, err := sv.repo.GetUserByUsername(c, req.Username)
	if err != nil && !errors.Is(err, repository.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if err == nil {
		c.JSON(http.StatusBadRequest, mapErrResp(fmt.Errorf("username %s is already taken", req.Username)))
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	arg := repository.CreateUserParams{
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
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	err = sv.taskDistributor.SendVerifyEmail(
		c,
		&worker.PayloadVerifyEmail{
			UserID: user.UserID,
		},
		asynq.MaxRetry(3),
		asynq.ProcessIn(5*time.Second),
		asynq.Queue(worker.QueueCritical),
	)
	message := "Please verify your email address to activate your account"
	errMsg := ""
	if err != nil {
		message = "Failed to send email verification"
		errMsg = err.Error()
	}

	c.JSON(http.StatusOK, GenericResponse[repository.CreateUserRow]{&user, &message, &errMsg})
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

	user, err := sv.repo.GetUserByUsername(c, req.Username)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, mapErrResp(err))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	// TODO: check if user has confirmed email or phone number

	if err := auth.CheckPassword(req.Password, user.HashedPassword); err != nil {
		c.JSON(http.StatusUnauthorized, mapErrResp(err))
		return
	}

	token, payload, err := sv.tokenGenerator.GenerateToken(user.UserID, user.Username, user.Email, sv.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	refreshToken, rfPayload, err := sv.tokenGenerator.GenerateToken(user.UserID, user.Username, user.Email, sv.config.RefreshTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	sessionID := uuid.New()
	session, err := sv.repo.CreateSession(c, repository.CreateSessionParams{
		SessionID:    sessionID,
		UserID:       user.UserID,
		RefreshToken: refreshToken,
		UserAgent:    c.GetHeader("User-Agent"),
		ClientIp:     c.ClientIP(),
		Blocked:      false,
		ExpiredAt:    time.Now().Add(sv.config.RefreshTokenDuration),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	loginResp := loginResponse{
		SessionID:            session.SessionID,
		TokenExpireAt:        payload.ExpiredAt,
		Token:                token,
		RefreshToken:         refreshToken,
		RefreshTokenExpireAt: rfPayload.ExpiredAt,
		User:                 mapToUserResponse(user),
	}
	c.JSON(http.StatusOK, GenericResponse[loginResponse]{&loginResp, nil, nil})
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
	refreshToken := c.GetHeader("Authorization")
	if refreshToken == "" {
		c.JSON(http.StatusUnauthorized, mapErrResp(fmt.Errorf("refresh token is required")))
		return
	}
	refreshTokenPayload, err := sv.tokenGenerator.VerifyToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, mapErrResp(err))
		return
	}

	session, err := sv.repo.GetSession(c, refreshTokenPayload.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, mapErrResp(err))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if session.UserID != refreshTokenPayload.UserID {
		err := errors.New("refresh token is not valid")
		c.JSON(http.StatusUnauthorized, mapErrResp(err))
		return
	}

	if session.RefreshToken != refreshToken {
		err := errors.New("refresh token is not valid")
		c.JSON(http.StatusUnauthorized, mapErrResp(err))
		return
	}

	if session.Blocked {
		err := errors.New("session is blocked")
		c.JSON(http.StatusUnauthorized, mapErrResp(err))
		return
	}

	if time.Now().After(session.ExpiredAt) {
		err := errors.New("refresh token was expired")
		c.JSON(http.StatusUnauthorized, mapErrResp(err))
		return
	}
	accessToken, _, err := sv.tokenGenerator.GenerateToken(session.UserID, refreshTokenPayload.Username, refreshTokenPayload.Email, sv.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusOK, &renewAccessTokenResp{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: sv.config.AccessTokenDuration,
	})
}

// updateUser godoc
// @Summary Update user info
// @Description Update user info
// @Tags users
// @Accept  json
// @Produce  json
// @Param input body updateUserRequest true "User info"
// @Success 200 {object} GenericResponse[repository.UpdateUserRow]
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

	user, err := sv.repo.GetUserByID(c, req.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, mapErrResp(err))
		return
	}

	if user.Role != repository.UserRoleAdmin && user.UserID != req.UserID {
		c.JSON(http.StatusUnauthorized, mapErrResp(err))
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
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusOK, GenericResponse[repository.UpdateUserRow]{&updatedUser, nil, nil})
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

	user, err := sv.repo.GetUserByID(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(err))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	userAddress, err := sv.repo.GetAddresses(c, user.UserID)
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

// listUsers godoc
// @Summary List users
// @Description List users
// @Tags users
// @Accept  json
// @Produce  json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} GenericResponse[[]userResponse]
// @Failure 500 {object} gin.H
// @Router /users/list [get]
func (sv *Server) listUsers(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, mapErrResp(fmt.Errorf("authorization payload is not provided")))
		return
	}
	user, err := sv.repo.GetUserByID(c, authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if user.Role != repository.UserRoleAdmin {
		c.JSON(http.StatusUnauthorized, mapErrResp(fmt.Errorf("user does not have permission")))
		return
	}

	var queries listUserParams
	if err := c.ShouldBindUri(&queries); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	users, err := sv.repo.ListUsers(c, repository.ListUsersParams{
		Limit:  int32(queries.PageSize),
		Offset: int32((queries.Page - 1) * queries.PageSize),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	userResp := make([]userResponse, 0)
	for _, user := range users {
		userAddress, err := sv.repo.GetAddresses(c, user.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, mapErrResp(err))
			return
		}

		addressResp := make([]addressResponse, 0)
		for _, address := range userAddress {
			addressResp = append(addressResp, mapAddressToAddressResponse(address))
		}
		userResp = append(userResp, mapToUserResponse(user))
		userResp[len(userResp)-1].Addresses = addressResp
	}

	c.JSON(http.StatusOK, GenericResponse[[]userResponse]{&userResp, nil, nil})
}

func (sv *Server) sendVerifyEmail(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, mapErrResp(fmt.Errorf("authorization payload is not provided")))
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
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusOK, GenericResponse[interface{}]{nil, nil, nil})
}

func (sv *Server) verifyEmail(c *gin.Context) {
	var query verifyEmailQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	verifyEmail, err := sv.repo.GetVerifyEmailByID(c, query.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("not found")))
		return
	}

	if verifyEmail.ExpiredAt.Before(time.Now()) {
		c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("expired")))
		return
	}

	if verifyEmail.VerifyCode != query.VerifyCode {
		c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("not found")))
		return
	}
	_, err = sv.repo.UpdateVerifyEmail(c, repository.UpdateVerifyEmailParams{
		ID:         verifyEmail.ID,
		VerifyCode: query.VerifyCode,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
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
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

}
