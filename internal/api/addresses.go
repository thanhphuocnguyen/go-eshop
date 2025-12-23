package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
)

// createAddress godoc
// @Summary Create a new addresses
// @Description Create a new addresses
// @Tags addresses
// @Accept json
// @Produce json
// @Param input body models.CreateAddress true "Create Address"
// @Success 200 {object} dto.ApiResponse[dto.AddressDetail]
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Router /users/addresses [post]
func (s *Server) createAddress(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	_, claims, err := jwtauth.FromContext(c)
	if err != nil {
		RespondInternalServerError(w, UnauthorizedCode, fmt.Errorf("authorization payload is not provided"))
		return
	}

	var req models.CreateAddress
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	validate := validator.New()
	if err := validate.Struct(&req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	addresses, err := s.repo.GetAddresses(c, claims["user_id"].(uuid.UUID))
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	if len(addresses) >= 10 {
		RespondBadRequest(w, InvalidBodyCode, fmt.Errorf("maximum number of addresses reached"))
		return
	}
	userID := uuid.MustParse(claims["userId"].(string))

	payload := repository.CreateAddressParams{
		PhoneNumber: req.Phone,
		UserID:      userID,
		Street:      req.Street,
		IsDefault:   req.IsDefault,
		City:        req.City,
		District:    req.District,
	}

	if req.Ward != nil {
		payload.Ward = req.Ward
	}

	created, err := s.repo.CreateAddress(c, payload)

	if req.IsDefault {
		err := s.repo.SetPrimaryAddressTx(c, repository.SetPrimaryAddressTxArgs{
			NewPrimaryID: created.ID,
			UserID:       claims["user_id"].(uuid.UUID),
		})
		if err != nil {
			RespondInternalServerError(w, InternalServerErrorCode, fmt.Errorf("failed to set primary addresses: %w", err))
			return
		}
		created.IsDefault = true
	}

	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	addressDetail := dto.MapAddressResponse(created)
	RespondSuccess(w, addressDetail)
}

// getAddresses godoc
// @Summary Get list of addresses
// @Description Get list of addresses
// @Tags addresses
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} dto.ApiResponse[[]dto.AddressDetail]
// @Failure 401 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /users/addresses [get]
func (s *Server) getAddresses(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	_, claims, err := jwtauth.FromContext(c)
	if err != nil {
		RespondInternalServerError(w, UnauthorizedCode, fmt.Errorf("authorization payload is not provided"))
		return
	}
	userID := uuid.MustParse(claims["userId"].(string))

	addresses, err := s.repo.GetAddresses(c, userID)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	addressesResponse := make([]dto.AddressDetail, len(addresses))
	for i, addresses := range addresses {
		addressesResponse[i] = dto.MapAddressResponse(addresses)
	}

	RespondSuccess(w, addressesResponse)
}

// updateAddress godoc
// @Summary Update an addresses
// @Description Update an addresses
// @Tags addresses
// @Accept json
// @Produce json
// @Param id path int true "Address ID"
// @Param input body models.UpdateAddress true "Update Address"
// @Success 200 {object} dto.ApiResponse[dto.AddressDetail]
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /users/addresses/{id} [put]
func (s *Server) updateAddress(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	_, claims, err := jwtauth.FromContext(c)
	if err != nil {
		RespondUnauthorized(w, UnauthorizedCode, fmt.Errorf("authorization payload is not provided"))
		return
	}
	userID := uuid.MustParse(claims["userId"].(string))

	var req models.UpdateAddress
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	validate := validator.New()
	if err := validate.Struct(&req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		RespondBadRequest(w, InvalidBodyCode, fmt.Errorf("id parameter is required"))
		return
	}

	_, err = s.repo.GetAddress(c, repository.GetAddressParams{
		ID:     uuid.MustParse(idParam),
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, err)
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	payload := repository.UpdateAddressParams{
		ID:     uuid.MustParse(idParam),
		UserID: userID,
	}

	if req.Phone != nil {
		payload.PhoneNumber = req.Phone
	}

	if req.City != nil {
		payload.City = req.City
	}

	if req.Address != nil {
		payload.Street = req.Address
	}

	if req.Ward != nil {
		payload.Ward = req.Ward
	}

	if req.District != nil {
		payload.District = req.District
	}

	if req.IsDefault != nil {
		if *req.IsDefault {
			err := s.repo.SetPrimaryAddressTx(c, repository.SetPrimaryAddressTxArgs{
				NewPrimaryID: uuid.MustParse(idParam),
				UserID:       userID,
			})
			if err != nil {
				RespondInternalServerError(w, InternalServerErrorCode, fmt.Errorf("failed to set primary addresses: %w", err))
				return
			}
		}
	}

	addresses, err := s.repo.UpdateAddress(c, payload)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	addressDetail := dto.MapAddressResponse(addresses)
	RespondSuccess(w, addressDetail)
}

// removeAddress godoc
// @Summary Remove an addresses
// @Description Remove an addresses
// @Tags addresses
// @Accept json
// @Produce json
// @Param id path int true "Address ID"
// @Success 204 {object} nil
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Router /users/addresses/{id} [delete]
func (s *Server) removeAddress(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	_, claims, err := jwtauth.FromContext(c)
	if err != nil {
		RespondInternalServerError(w, UnauthorizedCode, fmt.Errorf("authorization payload is not provided"))
		return
	}
	userID := uuid.MustParse(claims["userId"].(string))

	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		RespondBadRequest(w, InvalidBodyCode, fmt.Errorf("id parameter is required"))
		return
	}

	addresses, err := s.repo.GetAddress(c, repository.GetAddressParams{
		ID:     uuid.MustParse(idParam),
		UserID: userID,
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, err)
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	err = s.repo.DeleteAddress(c, repository.DeleteAddressParams{
		ID:     addresses.ID,
		UserID: addresses.UserID,
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, err)
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	RespondNoContent(w)
}

// setDefaultAddress godoc
// @Summary Set default addresses
// @Description Set default addresses
// @Tags addresses
// @Accept json
// @Produce json
// @Param id path int true "Address ID"
// @Success 204 {object} nil
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Router /users/addresses/{id}/default [put]
func (s *Server) setDefaultAddress(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	_, claims, err := jwtauth.FromContext(c)
	if err != nil {
		RespondInternalServerError(w, UnauthorizedCode, fmt.Errorf("authorization payload is not provided"))
		return
	}
	userID := uuid.MustParse(claims["userId"].(string))

	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		RespondBadRequest(w, InvalidBodyCode, fmt.Errorf("id parameter is required"))
		return
	}

	_, err = s.repo.GetAddress(c, repository.GetAddressParams{
		ID:     uuid.MustParse(idParam),
		UserID: userID,
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, err)
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	err = s.repo.SetPrimaryAddressTx(c, repository.SetPrimaryAddressTxArgs{
		NewPrimaryID: uuid.MustParse(idParam),
		UserID:       userID,
	})
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, fmt.Errorf("failed to set primary addresses: %w", err))
		return
	}
	RespondNoContent(w)
}
