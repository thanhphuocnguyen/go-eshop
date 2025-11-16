package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
)

// createAddressHandler godoc
// @Summary Create a new address
// @Description Create a new address
// @Tags address
// @Accept json
// @Produce json
// @Param input body CreateAddressRequest true "Create Address"
// @Success 200 {object} ApiResponse[AddressResponse]
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Router /address [post]
func (sv *Server) createAddressHandler(c *gin.Context) {
	authPayload, ok := c.MustGet(AuthPayLoad).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, createErr(UnauthorizedCode, fmt.Errorf("authorization payload is not provided")))
		return
	}

	var req CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	addresses, err := sv.repo.GetAddresses(c, authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	if len(addresses) >= 10 {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, fmt.Errorf("maximum number of addresses reached")))
		return
	}

	payload := repository.CreateAddressParams{
		PhoneNumber: req.Phone,
		UserID:      authPayload.UserID,
		Street:      req.Street,
		IsDefault:   req.IsDefault,
		City:        req.City,
		District:    req.District,
	}

	if req.Ward != nil {
		payload.Ward = req.Ward
	}

	address, err := sv.repo.CreateAddress(c, payload)

	if req.IsDefault {
		err := sv.repo.SetPrimaryAddressTx(c, repository.SetPrimaryAddressTxArgs{
			NewPrimaryID: address.ID,
			UserID:       authPayload.UserID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, fmt.Errorf("failed to set primary address: %w", err)))
			return
		}
		address.IsDefault = true
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	addressDetail := mapAddressResponse(address)
	c.JSON(http.StatusOK, createDataResp(c, addressDetail, nil, nil))
}

// getAddresses godoc
// @Summary Get list of addresses
// @Description Get list of addresses
// @Tags address
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} ApiResponse[[]AddressResponse]
// @Failure 401 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /address [get]
func (sv *Server) getAddressesHandlers(c *gin.Context) {
	authPayload, ok := c.MustGet(AuthPayLoad).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, createErr(UnauthorizedCode, fmt.Errorf("authorization payload is not provided")))
		return
	}

	addresses, err := sv.repo.GetAddresses(c, authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}
	addressesResponse := make([]AddressResponse, len(addresses))
	for i, address := range addresses {
		addressesResponse[i] = mapAddressResponse(address)
	}

	c.JSON(http.StatusOK, createDataResp(c, addressesResponse, nil, nil))
}

// updateAddressHandlers godoc
// @Summary Update an address
// @Description Update an address
// @Tags address
// @Accept json
// @Produce json
// @Param id path int true "Address ID"
// @Param input body UpdateAddressRequest true "Update Address"
// @Success 200 {object} ApiResponse[AddressResponse]
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /address/{id} [put]
func (sv *Server) updateAddressHandlers(c *gin.Context) {
	authPayload, ok := c.MustGet(AuthPayLoad).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusUnauthorized, createErr(UnauthorizedCode, fmt.Errorf("authorization payload is not provided")))
		return
	}
	var input UpdateAddressRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}
	var param UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	_, err := sv.repo.GetAddress(c, repository.GetAddressParams{
		ID:     uuid.MustParse(param.ID),
		UserID: authPayload.UserID,
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErr(NotFoundCode, err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	payload := repository.UpdateAddressParams{
		ID:     uuid.MustParse(param.ID),
		UserID: authPayload.UserID,
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
			err := sv.repo.SetPrimaryAddressTx(c, repository.SetPrimaryAddressTxArgs{
				NewPrimaryID: uuid.MustParse(param.ID),
				UserID:       authPayload.UserID,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, fmt.Errorf("failed to set primary address: %w", err)))
				return
			}
		}
	}

	address, err := sv.repo.UpdateAddress(c, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}
	addressDetail := mapAddressResponse(address)
	c.JSON(http.StatusOK, createDataResp(c, addressDetail, nil, nil))
}

// removeAddressHandlers godoc
// @Summary Remove an address
// @Description Remove an address
// @Tags address
// @Accept json
// @Produce json
// @Param id path int true "Address ID"
// @Success 204 {object}
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Router /address/{id} [delete]
func (sv *Server) removeAddressHandlers(c *gin.Context) {
	authPayload, ok := c.MustGet(AuthPayLoad).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, createErr(UnauthorizedCode, fmt.Errorf("authorization payload is not provided")))
		return
	}
	var param UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	address, err := sv.repo.GetAddress(c, repository.GetAddressParams{
		ID:     uuid.MustParse(param.ID),
		UserID: authPayload.UserID,
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErr(NotFoundCode, err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	err = sv.repo.DeleteAddress(c, repository.DeleteAddressParams{
		ID:     address.ID,
		UserID: address.UserID,
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErr(NotFoundCode, err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}
	c.Status(http.StatusNoContent)
}

// setDefaultAddressHandler godoc
// @Summary Set default address
// @Description Set default address
// @Tags address
// @Accept json
// @Produce json
// @Param id path int true "Address ID"
// @Success 204 {object}
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Router /address/{id}/default [put]
func (sv *Server) setDefaultAddressHandler(c *gin.Context) {
	authPayload, ok := c.MustGet(AuthPayLoad).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, createErr(UnauthorizedCode, fmt.Errorf("authorization payload is not provided")))
		return
	}
	var param UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErr(InvalidBodyCode, err))
		return
	}

	_, err := sv.repo.GetAddress(c, repository.GetAddressParams{
		ID:     uuid.MustParse(param.ID),
		UserID: authPayload.UserID,
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErr(NotFoundCode, err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, err))
		return
	}

	err = sv.repo.SetPrimaryAddressTx(c, repository.SetPrimaryAddressTxArgs{
		NewPrimaryID: uuid.MustParse(param.ID),
		UserID:       authPayload.UserID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErr(InternalServerErrorCode, fmt.Errorf("failed to set primary address: %w", err)))
		return
	}
	c.Status(http.StatusNoContent)
}
