package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/netip"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hibiken/asynq"
	repository "github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"github.com/thanhphuocnguyen/go-eshop/internal/worker"
)

// ------------------------------ s ------------------------------

// register godoc
// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept  json
// @Produce  json
// @Param input body RegisterRequestBody true "User info"
// @Success 200 {object} ApiResponse[UserDetail]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /auth/register [post]
func (sv *Server) register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	_, err := sv.repo.GetUserByUsername(r.Context(), req.Username)
	if err != nil && !errors.Is(err, repository.ErrRecordNotFound) {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	if err == nil {
		RespondBadRequest(w, UsernameExistedCode, fmt.Errorf("username %s is already taken", req.Username))
		return
	}

	hashedPassword, err := auth.HashPwd(req.Password)
	if err != nil {
		RespondInternalServerError(w, HashPasswordCode, err)
		return
	}

	userRole, err := sv.repo.GetRoleByCode(r.Context(), "user")
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
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
		adminRole, err := sv.repo.GetRoleByCode(r.Context(), "admin")
		if err != nil {
			RespondInternalServerError(w, InternalServerErrorCode, err)
			return
		}
		arg.RoleID = adminRole.ID
	}
	user, err := sv.repo.CreateUser(r.Context(), arg)

	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	var address dto.AddressDetail
	if req.Address != nil {

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

		createdAddress, err := sv.repo.CreateAddress(r.Context(), createAddressArgs)

		if err != nil {
			RespondInternalServerError(w, AddressCodeCode, err)
			return
		}
		ward := ""
		if createdAddress.Ward != nil {
			ward = *createdAddress.Ward
		}
		address = dto.AddressDetail{
			ID:        createdAddress.ID.String(),
			Phone:     createdAddress.PhoneNumber,
			Street:    createdAddress.Street,
			Ward:      &ward,
			District:  createdAddress.District,
			City:      createdAddress.City,
			Default:   createdAddress.IsDefault,
			CreatedAt: createdAddress.CreatedAt,
		}
	}

	emailPayload := &worker.PayloadVerifyEmail{UserID: user.ID}
	err = sv.taskDistributor.SendVerifyAccountEmail(
		r.Context(),
		emailPayload,
		asynq.MaxRetry(3),
		asynq.ProcessIn(5*time.Second),
		asynq.Queue(worker.QueueDefault),
	)

	if err != nil {
		RespondInternalServerError(w, ActivateUserCode, err)
		return
	}

	userResp := &dto.UserDetail{
		ID:            user.ID,
		Username:      user.Username,
		RoleID:        user.RoleID.String(),
		RoleCode:      userRole.Code,
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
	}
	if !dto.IsStructEmpty(address) {
		userResp.Addresses = []dto.AddressDetail{address}
	}

	RespondSuccess(w, r, userResp)
}

// login godoc
// @Summary login to the system
// @Description login to the system
// @Tags users
// @Accept  json
// @Produce  json
// @Param input body LoginRequest true "User info"
// @Success 200 {object} ApiResponse[LoginResponse]
// @Failure 401 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /auth/login [post]
func (sv *Server) login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	if req.Username == nil && req.Email == nil {
		RespondBadRequest(w, InvalidEmailCode, fmt.Errorf("username or email is required"))
		return
	}

	var user repository.User
	var err error = nil
	if req.Username != nil {
		user, err = sv.repo.GetUserByUsername(r.Context(), *req.Username)
	} else {
		user, err = sv.repo.GetUserByEmail(r.Context(), *req.Email)
	}

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondUnauthorized(w, NotFoundCode, err)
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	if err = auth.ComparePwd(req.Password, user.HashedPassword); err != nil {
		RespondUnauthorized(w, UnauthorizedCode, err)
		return
	}

	role, err := sv.repo.GetRoleByID(r.Context(), user.RoleID)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	accessToken, payload, err := sv.tokenGenerator.GenerateToken(user.ID, user.Username, role, sv.config.AccessTokenDuration)
	if err != nil {
		RespondInternalServerError(w, InvalidTokenCode, err)
		return
	}

	refreshToken, rfPayload, err := sv.tokenGenerator.GenerateToken(user.ID, user.Username, role, sv.config.RefreshTokenDuration)
	if err != nil {
		RespondInternalServerError(w, InvalidTokenCode, err)
		return
	}

	clientIP, err := netip.ParseAddr(r.RemoteAddr)
	if err != nil {
		// Fallback to localhost if parsing fails
		clientIP = netip.MustParseAddr("127.0.0.1")
	}

	session, err := sv.repo.InsertSession(r.Context(), repository.InsertSessionParams{
		ID:           rfPayload.ID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    r.Header.Get("User-Agent"),
		ClientIp:     clientIP,
		Blocked:      false,
		ExpiredAt:    utils.GetPgTypeTimestamp(time.Now().Add(sv.config.RefreshTokenDuration)),
	})

	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	loginResp := dto.LoginResponse{
		ID:                    session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  payload.ExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: rfPayload.ExpiresAt,
	}
	RespondSuccess(w, r, loginResp)
}

// refreshToken godoc
// @Summary Refresh token
// @Description Refresh token
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} ApiResponse[RefreshTokenResponse]
// @Failure 401 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /auth/refresh-token [post]
func (sv *Server) refreshToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		RespondUnauthorized(w, UnauthorizedCode, fmt.Errorf("refresh token is required"))
		return
	}

	refreshToken := authHeader[len("Bearer "):]
	refreshTokenPayload, err := sv.tokenGenerator.VerifyToken(refreshToken)
	if err != nil {
		RespondUnauthorized(w, UnauthorizedCode, err)
		return
	}

	session, err := sv.repo.GetSession(r.Context(), refreshTokenPayload.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, fmt.Errorf("session not found"))
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	if session.ID != refreshTokenPayload.ID {
		err := errors.New("refresh token is not valid")
		RespondUnauthorized(w, InvalidTokenCode, err)
		return
	}

	if session.RefreshToken != refreshToken {
		err := errors.New("refresh token is not valid")
		RespondUnauthorized(w, InvalidTokenCode, err)
		return
	}

	if session.Blocked {
		err := errors.New("session is blocked")
		RespondUnauthorized(w, InvalidSessionCode, err)
		return
	}

	if time.Now().After(session.ExpiredAt) {
		err := errors.New("refresh token was expired")
		RespondUnauthorized(w, InvalidSessionCode, err)
		return
	}

	role, err := sv.repo.GetRoleByID(r.Context(), refreshTokenPayload.RoleID)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	accessToken, _, err := sv.tokenGenerator.GenerateToken(session.UserID, refreshTokenPayload.Username, role, sv.config.AccessTokenDuration)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	resp := dto.RefreshToken{AccessToken: accessToken, AccessTokenExpiresAt: time.Now().Add(sv.config.AccessTokenDuration)}
	RespondSuccess(w, r, resp)
}

// Setup authentication routes
func (sv *Server) addAuthRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", sv.register)
		r.Post("/login", sv.login)
		r.Post("/refresh-token", sv.refreshToken)
	})
}
