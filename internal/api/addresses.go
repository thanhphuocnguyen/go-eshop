package api

import (
	"encoding/json"
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
// @Param input body CreateAddressRequest true "Create Address"
// @Success 200 {object} ApiResponse[AddressResponse]
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Router /users/addresses [post]
func (sv *Server) createAddress(w http.ResponseWriter, r *http.Request) {
	token, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {

		RespondInternalServerError(w, UnauthorizedCode, fmt.Errorf("authorization payload is not provided"))
		return
	}

	var req models.CreateAddress
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	validate := validator.New()
	if err := validate.Struct(&req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	addresses, err := sv.repo.GetAddresses(r.Context(), claims["user_id"].(uuid.UUID))
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	if len(addresses) >= 10 {
		RespondBadRequest(w, InvalidBodyCode, fmt.Errorf("maximum number of addresses reached"))
		return
	}

	payload := repository.CreateAddressParams{
		PhoneNumber: req.Phone,
		UserID:      claims["userId"].(uuid.UUID),
		Street:      req.Street,
		IsDefault:   req.IsDefault,
		City:        req.City,
		District:    req.District,
	}

	if req.Ward != nil {
		payload.Ward = req.Ward
	}

	created, err := sv.repo.CreateAddress(r.Context(), payload)

	if req.IsDefault {
		err := sv.repo.SetPrimaryAddressTx(r.Context(), repository.SetPrimaryAddressTxArgs{
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
	RespondSuccess(w, r, addressDetail)
}

// getAddresses godoc
// @Summary Get list of addresses
// @Description Get list of addresses
// @Tags addresses
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} ApiResponse[[]AddressResponse]
// @Failure 401 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /users/addresses [get]
func (sv *Server) getAddresses(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if !ok {
		RespondInternalServerError(w, UnauthorizedCode, fmt.Errorf("authorization payload is not provided"))
		return
	}

	addresses, err := sv.repo.GetAddresses(r.Context(), userID)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	addressesResponse := make([]dto.AddressDetail, len(addresses))
	for i, addresses := range addresses {
		addressesResponse[i] = dto.MapAddressResponse(addresses)
	}

	RespondSuccess(w, r, addressesResponse)
}

// updateAddress godoc
// @Summary Update an addresses
// @Description Update an addresses
// @Tags addresses
// @Accept json
// @Produce json
// @Param id path int true "Address ID"
// @Param input body UpdateAddressRequest true "Update Address"
// @Success 200 {object} ApiResponse[AddressResponse]
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /users/addresses/{id} [put]
func (sv *Server) updateAddress(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if !ok {
		RespondUnauthorized(w, UnauthorizedCode, fmt.Errorf("authorization payload is not provided"))
		return
	}
	var input models.UpdateAddress
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	validate := validator.New()
	if err := validate.Struct(&input); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		RespondBadRequest(w, InvalidBodyCode, fmt.Errorf("id parameter is required"))
		return
	}

	_, err := sv.repo.GetAddress(r.Context(), repository.GetAddressParams{
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

	if input.Phone != nil {
		payload.PhoneNumber = input.Phone
	}

	if input.City != nil {
		payload.City = input.City
	}

	if input.Address != nil {
		payload.Street = input.Address
	}

	if input.Ward != nil {
		payload.Ward = input.Ward
	}

	if input.District != nil {
		payload.District = input.District
	}

	if input.IsDefault != nil {
		if *input.IsDefault {
			err := sv.repo.SetPrimaryAddressTx(r.Context(), repository.SetPrimaryAddressTxArgs{
				NewPrimaryID: uuid.MustParse(idParam),
				UserID:       userID,
			})
			if err != nil {
				RespondInternalServerError(w, InternalServerErrorCode, fmt.Errorf("failed to set primary addresses: %w", err))
				return
			}
		}
	}

	addresses, err := sv.repo.UpdateAddress(r.Context(), payload)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	addressDetail := dto.MapAddressResponse(addresses)
	RespondSuccess(w, r, addressDetail)
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
func (sv *Server) removeAddress(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if !ok {
		RespondInternalServerError(w, UnauthorizedCode, fmt.Errorf("authorization payload is not provided"))
		return
	}

	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		RespondBadRequest(w, InvalidBodyCode, fmt.Errorf("id parameter is required"))
		return
	}

	addresses, err := sv.repo.GetAddress(r.Context(), repository.GetAddressParams{
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

	err = sv.repo.DeleteAddress(r.Context(), repository.DeleteAddressParams{
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
func (sv *Server) setDefaultAddress(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if !ok {
		RespondInternalServerError(w, UnauthorizedCode, fmt.Errorf("authorization payload is not provided"))
		return
	}

	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		RespondBadRequest(w, InvalidBodyCode, fmt.Errorf("id parameter is required"))
		return
	}

	_, err := sv.repo.GetAddress(r.Context(), repository.GetAddressParams{
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

	err = sv.repo.SetPrimaryAddressTx(r.Context(), repository.SetPrimaryAddressTxArgs{
		NewPrimaryID: uuid.MustParse(idParam),
		UserID:       userID,
	})
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, fmt.Errorf("failed to set primary addresses: %w", err))
		return
	}
	RespondNoContent(w)
}
