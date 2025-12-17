package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
	"github.com/thanhphuocnguyen/go-eshop/internal/worker"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
)

// updateUser godoc
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
func (sv *Server) updateUser(w http.ResponseWriter, r *http.Request) {
	authPayload := r.Context().Value("auth").(*auth.TokenPayload)
	var req models.UpdateUserModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondBadRequest(w, InvalidEmailCode, err)
		return
	}

	validate := validator.New()
	if err := validate.Struct(&req); err != nil {
		RespondBadRequest(w, InvalidEmailCode, err)
		return
	}

	userId := uuid.MustParse(req.UserID)
	user, err := sv.repo.GetUserByID(r.Context(), userId)
	if err != nil {
		RespondUnauthorized(w, UnauthorizedCode, err)
		return
	}

	if authPayload.RoleCode != "admin" && user.ID != userId {
		RespondUnauthorized(w, UnauthorizedCode, fmt.Errorf("access denied"))
		return
	}

	arg := repository.UpdateUserParams{
		ID: userId,
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

	updatedUser, err := sv.repo.UpdateUser(r.Context(), arg)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	resp := dto.MapToUserResponse(updatedUser, authPayload.RoleCode)
	RespondSuccess(w, r, resp)
}

// getCurrentUser godoc
// @Summary Get user info
// @Description Get user info
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} ApiResponse[UserDetail]
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /users/me [get]
func (sv *Server) getCurrentUser(w http.ResponseWriter, r *http.Request) {
	authPayload, ok := r.Context().Value("auth").(*auth.TokenPayload)
	if !ok {
		RespondInternalServerError(w, InternalServerErrorCode, errors.New("authorization payload is not provided"))
		return
	}

	var userResp dto.UserDetail

	user, err := sv.repo.GetUserByID(r.Context(), authPayload.UserID)
	if err != nil {
		RespondNotFound(w, NotFoundCode, err)
		return
	}

	userAddress, err := sv.repo.GetAddresses(r.Context(), user.ID)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	addressResp := make([]dto.AddressDetail, 0)
	for _, address := range userAddress {
		addressResp = append(addressResp, dto.MapAddressResponse(address))
	}
	userResp = dto.MapToUserResponse(user, authPayload.RoleCode)
	userResp.Addresses = addressResp

	RespondSuccess(w, r, userResp)
}

// sendVerifyEmail godoc
// @Summary Send verify email
// @Description Send verify email
// @Tags users
// @Accept  json
// @Produce  json
// @Success 204 {object} nil
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /users/verify-email [post]
// @Security BearerAuth
func (sv *Server) sendVerifyEmail(w http.ResponseWriter, r *http.Request) {
	authPayload, ok := r.Context().Value("auth").(*auth.TokenPayload)
	if !ok {
		RespondInternalServerError(w, InternalServerErrorCode, errors.New("authorization payload is not provided"))
		return
	}
	user, err := sv.repo.GetUserByID(r.Context(), authPayload.UserID)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	if user.VerifiedEmail {
		RespondBadRequest(w, InvalidEmailCode, fmt.Errorf("email already verified"))
		return
	}

	err = sv.taskDistributor.SendVerifyAccountEmail(
		r.Context(),
		&worker.PayloadVerifyEmail{
			UserID: authPayload.UserID,
		},
		asynq.MaxRetry(3),
		asynq.ProcessIn(5*time.Second),
		asynq.Queue(worker.QueueCritical),
	)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondNoContent(w)
}

// VerifyEmail godoc
// @Summary Verify email
// @Description Verify email
// @Tags users
// @Accept  json
// @Produce  json
// @Param id query int true "ID"
// @Param verify_code query string true "Verify code"
// @Success 200 {object} nil
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /users/verify-email [get]
func (sv *Server) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	queryParams, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		RespondBadRequest(w, InvalidEmailCode, err)
		return
	}

	var query models.VerifyEmailQuery
	if verifyCode := queryParams.Get("verify_code"); verifyCode != "" {
		query.VerifyCode = verifyCode
	} else {
		RespondBadRequest(w, InvalidEmailCode, fmt.Errorf("verify_code is required"))
		return
	}

	validate := validator.New()
	if err := validate.Struct(&query); err != nil {
		RespondBadRequest(w, InvalidEmailCode, err)
		return
	}

	verifyEmail, err := sv.repo.GetVerifyEmailByVerifyCode(r.Context(), query.VerifyCode)
	if err != nil {
		RespondNotFound(w, NotFoundCode, err)
		return
	}

	// Create a transaction to ensure both operations succeed or fail together
	err = sv.repo.VerifyEmailTx(r.Context(), repository.VerifyEmailTxArgs{
		VerifyEmail: verifyEmail,
		VerifyCode:  query.VerifyCode,
	})

	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	user, err := sv.repo.GetUserByID(r.Context(), verifyEmail.UserID)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	// Render HTML success page
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	htmlContent := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<title>Email Verification Success</title>
	</head>
	<body>
		<h1>Email Verification Successful</h1>
		<p>Hello %s,</p>
		<p>Your email %s has been successfully verified!</p>
	</body>
	</html>
	`, user.Username, user.Email)
	w.Write([]byte(htmlContent))
}

// Setup user-related routes
func (sv *Server) addUserRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		// Apply authentication middleware
		r.Use(func(h http.Handler) http.Handler {
			return authenticateMiddleware(h, sv.tokenGenerator)
		})

		r.Get("/me", sv.getCurrentUser)
		r.Patch("/me", sv.updateUser)
		r.Post("/send-verify-email", sv.sendVerifyEmail)

		// Address routes
		r.Route("/addresses", func(r chi.Router) {
			r.Post("/", sv.createAddress)
			r.Patch("/{id}/default", sv.setDefaultAddress)
			r.Get("/", sv.getAddresses)
			r.Patch("/{id}", sv.updateAddress)
			r.Delete("/{id}", sv.removeAddress)
		})
	})
}
