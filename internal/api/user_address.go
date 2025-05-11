package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
)

// ------------------------------ Params ------------------------------
type GetAddressParams struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type CreateAddressReq struct {
	Phone     string  `json:"phone" binding:"required,min=10,max=15"`
	Street    string  `json:"street" binding:"required"`
	District  string  `json:"district" binding:"required"`
	City      string  `json:"city" binding:"required"`
	Ward      *string `json:"ward,omitempty" binding:"omitempty,max=100"`
	IsDefault bool    `json:"is_default,omitempty" binding:"omitempty"`
}

type UpdateAddressReq struct {
	Phone     *string `json:"phone" binding:"omitempty"`
	Address   *string `json:"address_1" binding:"omitempty"`
	Ward      *string `json:"ward" binding:"omitempty"`
	District  *string `json:"district" binding:"omitempty"`
	City      *string `json:"city" binding:"omitempty"`
	IsDefault *bool   `json:"is_default" binding:"omitempty"`
}

// ------------------------------ API Models ------------------------------
type Address struct {
	Phone    string  `json:"phone"`
	Street   string  `json:"street"`
	Ward     *string `json:"ward,omitempty"`
	District string  `json:"district"`
	City     string  `json:"city"`
}
type AddressResponse struct {
	ID        string    `json:"id"`
	Default   bool      `json:"default"`
	CreatedAt time.Time `json:"created_at"`
	Phone     string    `json:"phone"`
	Street    string    `json:"street"`
	Ward      *string   `json:"ward,omitempty"`
	District  string    `json:"district"`
	City      string    `json:"city"`
}

// ------------------------------ Mapper ------------------------------
func mapAddressResponse(address repository.UserAddress) AddressResponse {
	return AddressResponse{
		ID:        address.ID.String(),
		Default:   address.Default,
		CreatedAt: address.CreatedAt,
		Phone:     address.Phone,
		Street:    address.Street,
		Ward:      address.Ward,
		District:  address.District,
		City:      address.City,
	}
}

// ------------------------------ Handlers ------------------------------

// createAddressHandler godoc
// @Summary Create a new address
// @Description Create a new address
// @Tags address
// @Accept json
// @Produce json
// @Param input body CreateAddressReq true "Create Address"
// @Success 200 {object} ApiResponse[AddressResponse]
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Router /address [post]
func (sv *Server) createAddressHandler(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, createErrorResponse[AddressResponse](UnauthorizedCode, "", fmt.Errorf("authorization payload is not provided")))
		return
	}

	var req CreateAddressReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[AddressResponse](InvalidBodyCode, "", err))
		return
	}

	addresses, err := sv.repo.GetAddresses(c, authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[AddressResponse](InternalServerErrorCode, "", err))
		return
	}

	if len(addresses) >= 10 {
		c.JSON(http.StatusBadRequest, createErrorResponse[AddressResponse](InvalidBodyCode, "", fmt.Errorf("maximum number of addresses reached")))
		return
	}

	payload := repository.CreateAddressParams{
		Phone:    req.Phone,
		UserID:   authPayload.UserID,
		Street:   req.Street,
		Default:  req.IsDefault,
		City:     req.City,
		District: req.District,
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
			c.JSON(http.StatusInternalServerError, createErrorResponse[AddressResponse](InternalServerErrorCode, "", fmt.Errorf("failed to set primary address: %w", err)))
			return
		}
		address.Default = true
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[AddressResponse](InternalServerErrorCode, "", err))
		return
	}

	addressDetail := mapAddressResponse(address)
	c.JSON(http.StatusOK, createSuccessResponse(c, addressDetail, "", nil, nil))
}

// getAddresses godoc
// @Summary Get list of addresses
// @Description Get list of addresses
// @Tags address
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} ApiResponse[[]AddressResponse]
// @Router /address [get]
func (sv *Server) getAddressesHandlers(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, createErrorResponse[[]AddressResponse](UnauthorizedCode, "", fmt.Errorf("authorization payload is not provided")))
		return
	}

	addresses, err := sv.repo.GetAddresses(c, authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[[]AddressResponse](InternalServerErrorCode, "", err))
		return
	}
	addressesResponse := make([]AddressResponse, len(addresses))
	for i, address := range addresses {
		addressesResponse[i] = mapAddressResponse(address)
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, addressesResponse, "", nil, nil))
}

// updateAddressHandlers godoc
// @Summary Update an address
// @Description Update an address
// @Tags address
// @Accept json
// @Produce json
// @Param id path int true "Address ID"
// @Param input body UpdateAddressReq true "Update Address"
// @Success 200 {object} ApiResponse[AddressResponse]
// @Router /address/{id} [put]
func (sv *Server) updateAddressHandlers(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusUnauthorized, createErrorResponse[AddressResponse](UnauthorizedCode, "", fmt.Errorf("authorization payload is not provided")))
		return
	}
	var input UpdateAddressReq
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[AddressResponse](InvalidBodyCode, "", err))
		return
	}
	var param GetAddressParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[AddressResponse](InvalidBodyCode, "", err))
		return
	}

	_, err := sv.repo.GetAddress(c, repository.GetAddressParams{
		ID:     uuid.MustParse(param.ID),
		UserID: authPayload.UserID,
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[AddressResponse](NotFoundCode, "", err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[AddressResponse](InternalServerErrorCode, "", err))
		return
	}

	payload := repository.UpdateAddressParams{
		ID:     uuid.MustParse(param.ID),
		UserID: authPayload.UserID,
	}

	if input.Phone != nil {
		payload.Phone = input.Phone
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
				c.JSON(http.StatusInternalServerError, createErrorResponse[AddressResponse](InternalServerErrorCode, "", fmt.Errorf("failed to set primary address: %w", err)))
				return
			}
		}
	}

	address, err := sv.repo.UpdateAddress(c, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[AddressResponse](InternalServerErrorCode, "", err))
		return
	}
	addressDetail := mapAddressResponse(address)
	c.JSON(http.StatusOK, createSuccessResponse(c, addressDetail, "", nil, nil))
}

// removeAddressHandlers godoc
// @Summary Remove an address
// @Description Remove an address
// @Tags address
// @Accept json
// @Produce json
// @Param id path int true "Address ID"
// @Success 200 {object} ApiResponse[bool]
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /address/{id} [delete]
func (sv *Server) removeAddressHandlers(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](UnauthorizedCode, "", fmt.Errorf("authorization payload is not provided")))
		return
	}
	var param GetAddressParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[bool](InvalidBodyCode, "", err))
		return
	}

	address, err := sv.repo.GetAddress(c, repository.GetAddressParams{
		ID:     uuid.MustParse(param.ID),
		UserID: authPayload.UserID,
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[bool](NotFoundCode, "", err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
		return
	}

	if address.Deleted {
		c.JSON(http.StatusNotFound, createErrorResponse[bool](NotFoundCode, "", fmt.Errorf("address has been removed")))
		return
	}

	err = sv.repo.DeleteAddress(c, repository.DeleteAddressParams{
		ID:     uuid.MustParse(param.ID),
		UserID: authPayload.UserID,
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[bool](NotFoundCode, "", err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
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
// @Success 200 {object} ApiResponse[bool]
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /address/{id}/default [put]
func (sv *Server) setDefaultAddressHandler(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](UnauthorizedCode, "", fmt.Errorf("authorization payload is not provided")))
		return
	}
	var param GetAddressParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[bool](InvalidBodyCode, "", err))
		return
	}

	address, err := sv.repo.GetAddress(c, repository.GetAddressParams{
		ID:     uuid.MustParse(param.ID),
		UserID: authPayload.UserID,
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[bool](NotFoundCode, "", err))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
		return
	}

	if address.Deleted {
		c.JSON(http.StatusNotFound, createErrorResponse[bool](NotFoundCode, "", fmt.Errorf("address has been removed")))
		return
	}

	err = sv.repo.SetPrimaryAddressTx(c, repository.SetPrimaryAddressTxArgs{
		NewPrimaryID: uuid.MustParse(param.ID),
		UserID:       authPayload.UserID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", fmt.Errorf("failed to set primary address: %w", err)))
		return
	}
	c.JSON(http.StatusOK, createSuccessResponse(c, true, "", nil, nil))
}
