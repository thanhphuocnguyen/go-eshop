package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/netip"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	repository "github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"github.com/thanhphuocnguyen/go-eshop/internal/worker"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
)

// ------------------------------ Handlers ------------------------------

// RegisterHandler godoc
// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept  json
// @Produce  json
// @Param input body RegisterRequestBody true "User info"
// @Success 200 {object} ApiResponse[UserDetail]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /users [post]
func (sv *Server) RegisterHandler(c *gin.Context) {
	var req RegisterRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, "", err))
		return
	}

	_, err := sv.repo.GetUserByUsername(c, req.Username)
	if err != nil && !errors.Is(err, repository.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, "", err))
		return
	}

	if err == nil {
		c.JSON(http.StatusBadRequest, createErr(UsernameExistedCode, "", fmt.Errorf("username %s is already taken", req.Username)))
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(HashPasswordCode, "", err))
		return
	}

	userRole, err := sv.repo.GetRoleByCode(c, string(repository.UserRoleCodeUser))
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, "", err))
		return
	}
	arg := repository.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		PhoneNumber:    req.Phone,
		RoleID:         userRole.ID,
	}

	if req.Username == "admin" {
		adminRole, err := sv.repo.GetRoleByCode(c, string(repository.UserRoleCodeAdmin))
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, "", err))
			return
		}
		arg.RoleID = adminRole.ID
	}
	user, err := sv.repo.CreateUser(c, arg)

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, "", err))
		return
	}

	createAddressArgs := repository.CreateAddressParams{
		UserID:      user.ID,
		PhoneNumber: req.Address.Phone,
		Street:      req.Address.Street,
		City:        req.Address.City,
		District:    req.Address.District,
		IsDefault:   true,
	}

	if req.Address.Ward != nil {
		createAddressArgs.Ward = req.Address.Ward
	}
	createdAddress, err := sv.repo.CreateAddress(c, createAddressArgs)

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(AddressCodeCode, "", err))
		return
	}

	emailPayload := &worker.PayloadVerifyEmail{UserID: user.ID}
	err = sv.taskDistributor.SendVerifyAccountEmail(
		c,
		emailPayload,
		asynq.MaxRetry(3),
		asynq.ProcessIn(5*time.Second),
		asynq.Queue(worker.QueueDefault),
	)

	if err != nil {
		createErr(ActivateUserCode, "Please verify your email address to activate your account", err)
		return
	}

	ward := ""
	if createdAddress.Ward != nil {
		ward = *createdAddress.Ward
	}

	userResp := &UserDetail{
		ID:            user.ID,
		Username:      user.Username,
		RoleID:        user.RoleID.String(),
		RoleCode:      repository.Role(userRole.Code),
		Email:         user.Email,
		CreatedAt:     user.CreatedAt.String(),
		VerifiedEmail: user.VerifiedEmail,
		VerifiedPhone: user.VerifiedPhone,
		UpdatedAt:     user.UpdatedAt.String(),
		FirstName:     user.FirstName,
		Locked:        user.Locked,
		Phone:         user.PhoneNumber,
		AvatarURL:     user.AvatarUrl,
		AvatarID:      user.AvatarImageID,
		LastName:      user.LastName,
		Addresses: []AddressResponse{{
			ID:        createdAddress.ID.String(),
			Phone:     createdAddress.PhoneNumber,
			Street:    createdAddress.Street,
			Ward:      &ward,
			District:  createdAddress.District,
			City:      createdAddress.City,
			Default:   createdAddress.IsDefault,
			CreatedAt: createdAddress.CreatedAt,
		}},
	}

	c.JSON(http.StatusOK, createDataResp(c, userResp, "Created user and address successfully", nil, nil))
}

// LoginHandler godoc
// @Summary Login to the system
// @Description Login to the system
// @Tags users
// @Accept  json
// @Produce  json
// @Param input body LoginRequest true "User info"
// @Success 200 {object} ApiResponse[LoginResponse]
// @Failure 401 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /users/LoginHandler [post]
func (sv *Server) LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, "", err))
		return
	}
	if req.Username == nil && req.Email == nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidEmailCode, "", fmt.Errorf("username or email is required")))
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
			c.JSON(http.StatusUnauthorized, createErr(NotFoundCode, "User not existed", err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, "Internal Server Error", err))
		return
	}

	if err := auth.CheckPassword(req.Password, user.HashedPassword); err != nil {
		c.JSON(http.StatusUnauthorized, createErr(UnauthorizedCode, "Invalid credentials", err))
		return
	}

	role, err := sv.repo.GetRoleByID(c, user.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, "", err))
		return
	}
	accessToken, payload, err := sv.tokenGenerator.GenerateToken(user.ID, user.Username, role, sv.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InvalidTokenCode, "", err))
		return
	}

	refreshToken, rfPayload, err := sv.tokenGenerator.GenerateToken(user.ID, user.Username, role, sv.config.RefreshTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InvalidTokenCode, "", err))
		return
	}

	clientIP, err := netip.ParseAddr(c.ClientIP())
	if err != nil {
		// Fallback to localhost if parsing fails
		clientIP = netip.MustParseAddr("127.0.0.1")
	}

	session, err := sv.repo.InsertSession(c, repository.InsertSessionParams{
		ID:           rfPayload.ID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    c.GetHeader("User-Agent"),
		ClientIp:     clientIP,
		Blocked:      false,
		ExpiredAt:    utils.GetPgTypeTimestamp(time.Now().Add(sv.config.RefreshTokenDuration)),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, "", err))
		return
	}

	loginResp := LoginResponse{
		ID:                    session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  payload.Expires,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: rfPayload.Expires,
	}
	c.JSON(http.StatusOK, createDataResp(c, loginResp, "success", nil, nil))
}

// refreshTokenHandler godoc
// @Summary Refresh token
// @Description Refresh token
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} ApiResponse[RefreshTokenResponse]
// @Failure 401 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /users/refresh-token [post]
func (sv *Server) refreshTokenHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, createErr(UnauthorizedCode, "", fmt.Errorf("refresh token is required")))
		return
	}

	refreshToken := authHeader[len("Bearer "):]
	refreshTokenPayload, err := sv.tokenGenerator.VerifyToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, createErr(UnauthorizedCode, "", err))
		return
	}

	session, err := sv.repo.GetSession(c, refreshTokenPayload.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound,
				createErr(NotFoundCode, "Not found", fmt.Errorf("session not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, "", err))
		return
	}

	if session.ID != refreshTokenPayload.ID {
		err := errors.New("refresh token is not valid")
		c.JSON(http.StatusUnauthorized, createErr(InvalidTokenCode, "", err))
		return
	}

	if session.RefreshToken != refreshToken {
		err := errors.New("refresh token is not valid")
		c.JSON(http.StatusUnauthorized, createErr(InvalidTokenCode, "", err))
		return
	}

	if session.Blocked {
		err := errors.New("session is blocked")
		c.JSON(http.StatusUnauthorized, createErr(InvalidSessionCode, "", err))
		return
	}

	if time.Now().After(session.ExpiredAt) {
		err := errors.New("refresh token was expired")
		c.JSON(http.StatusUnauthorized, createErr(InvalidSessionCode, "", err))
		return
	}

	role, err := sv.repo.GetRoleByID(c, refreshTokenPayload.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, "", err))
		return
	}

	accessToken, _, err := sv.tokenGenerator.GenerateToken(session.UserID, refreshTokenPayload.Username, role, sv.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, "", err))
		return
	}

	c.JSON(http.StatusOK,
		createDataResp(c, RefreshTokenResponse{AccessToken: accessToken, AccessTokenExpiresAt: time.Now().Add(sv.config.AccessTokenDuration)}, "success", nil, nil))
}
