package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
	"github.com/thanhphuocnguyen/go-eshop/internal/worker"
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
func (s *Server) updateUser(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
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
	user, err := s.repo.GetUserByID(r.Context(), userId)
	if err != nil {
		RespondUnauthorized(w, UnauthorizedCode, err)
		return
	}

	if claims["roleCode"] != "admin" && user.ID != userId {
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

	updatedUser, err := s.repo.UpdateUser(r.Context(), arg)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondSuccess(w, r, updatedUser)
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
func (s *Server) getCurrentUser(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	userId := claims["userId"].(uuid.UUID)
	roleCode := claims["roleCode"].(string)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, errors.New("authorization payload is not provided"))
		return
	}

	var userResp dto.UserDetail

	user, err := s.repo.GetUserByID(r.Context(), userId)
	if err != nil {
		RespondNotFound(w, NotFoundCode, err)
		return
	}

	userAddress, err := s.repo.GetAddresses(r.Context(), user.ID)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	addressResp := make([]dto.AddressDetail, 0)
	for _, address := range userAddress {
		addressResp = append(addressResp, dto.MapAddressResponse(address))
	}
	userResp = dto.MapToUserResponse(user, roleCode)
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
func (s *Server) sendVerifyEmail(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, errors.New("authorization payload is not provided"))
		return
	}
	userId := claims["userId"].(uuid.UUID)
	user, err := s.repo.GetUserByID(r.Context(), userId)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	if user.VerifiedEmail {
		RespondBadRequest(w, InvalidEmailCode, fmt.Errorf("email already verified"))
		return
	}

	err = s.taskDistributor.SendVerifyAccountEmail(
		r.Context(),
		&worker.PayloadVerifyEmail{
			UserID: userId,
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
func (s *Server) VerifyEmail(w http.ResponseWriter, r *http.Request) {
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

	verifyEmail, err := s.repo.GetVerifyEmailByVerifyCode(r.Context(), query.VerifyCode)
	if err != nil {
		RespondNotFound(w, NotFoundCode, err)
		return
	}

	// Create a transaction to ensure both operations succeed or fail together
	err = s.repo.VerifyEmailTx(r.Context(), repository.VerifyEmailTxArgs{
		VerifyEmail: verifyEmail,
		VerifyCode:  query.VerifyCode,
	})

	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	user, err := s.repo.GetUserByID(r.Context(), verifyEmail.UserID)
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
func (s *Server) addUserRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/me", s.getCurrentUser)
		r.Patch("/me", s.updateUser)
		r.Post("/send-verify-email", s.sendVerifyEmail)

		// Address routes
		r.Route("/addresses", func(r chi.Router) {
			r.Post("/", s.createAddress)
			r.Patch("/{id}/default", s.setDefaultAddress)
			r.Get("/", s.getAddresses)
			r.Patch("/{id}", s.updateAddress)
			r.Delete("/{id}", s.removeAddress)
		})
	})
}
