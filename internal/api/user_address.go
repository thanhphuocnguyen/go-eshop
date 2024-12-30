package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhphuocnguyen/go-eshop/internal/auth"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/postgres"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/sqlc"
	"github.com/thanhphuocnguyen/go-eshop/internal/util"
)

// ------------------------------ Params ------------------------------
type GetAddressParams struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type CreateAddressParams struct {
	Phone     string `json:"phone" binding:"required,min=10,max=15"`
	Street    string `json:"street" binding:"required"`
	Ward      string `json:"ward,omitempty" binding:"omitempty"`
	District  string `json:"district" binding:"required"`
	City      string `json:"city" binding:"required"`
	IsDefault bool   `json:"is_default,omitempty" binding:"omitempty"`
}

type UpdateAddressParams struct {
	Phone     *string `json:"phone" binding:"omitempty"`
	Address1  *string `json:"address_1" binding:"omitempty"`
	Address2  *string `json:"address_2" binding:"omitempty"`
	Ward      *string `json:"ward" binding:"omitempty"`
	District  *string `json:"district" binding:"omitempty"`
	City      *string `json:"city" binding:"omitempty"`
	IsDefault *bool   `json:"is_default" binding:"omitempty"`
}

// ------------------------------ API Models ------------------------------
type AddressResponse struct {
	ID        int64   `json:"id"`
	Phone     string  `json:"phone"`
	Address1  string  `json:"address_1"`
	Address2  *string `json:"address_2,omitempty"`
	Ward      *string `json:"ward"`
	District  string  `json:"district"`
	City      string  `json:"city"`
	IsDefault bool    `json:"is_default"`
}

// ------------------------------ Mapper ------------------------------
func mapAddressResponse(address sqlc.UserAddress) AddressResponse {
	return AddressResponse{
		ID:        address.ID,
		Phone:     address.Phone,
		Address1:  address.Street,
		Ward:      &address.Ward.String,
		District:  address.District,
		City:      address.City,
		IsDefault: address.IsPrimary,
	}
}

// ------------------------------ Handlers ------------------------------

// createAddress godoc
// @Summary Create a new address
// @Description Create a new address
// @Tags address
// @Accept json
// @Produce json
// @Param input body CreateAddressParams true "Create Address"
// @Success 200 {object} GenericResponse[AddressResponse]
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Router /address [post]
func (sv *Server) createAddress(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, mapErrResp(fmt.Errorf("authorization payload is not provided")))
		return
	}
	addresses, err := sv.postgres.GetAddresses(c, authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if len(addresses) >= 10 {
		c.JSON(http.StatusBadRequest, mapErrResp(fmt.Errorf("maximum number of addresses reached")))
		return
	}
	var req CreateAddressParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	payload := sqlc.CreateAddressParams{
		UserID:   authPayload.UserID,
		Phone:    req.Phone,
		Street:   req.Street,
		City:     req.City,
		District: req.District,
	}

	if req.Ward != "" {
		payload.Ward = util.GetPgTypeText(req.Ward)
	}

	address, err := sv.postgres.CreateAddress(c, payload)

	if req.IsDefault {
		err := sv.postgres.SetPrimaryAddressTx(c, postgres.SetPrimaryAddressTxParams{
			NewPrimaryID: address.ID,
			UserID:       authPayload.UserID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, mapErrResp(fmt.Errorf("failed to set primary address: %w", err)))
			return
		}
		address.IsPrimary = true
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	addressDetail := mapAddressResponse(address)
	c.JSON(http.StatusOK, GenericResponse[AddressResponse]{&addressDetail, nil, nil})
}

// getAddresses godoc
// @Summary Get list of addresses
// @Description Get list of addresses
// @Tags address
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} GenericResponse[[]AddressResponse]
// @Router /address [get]
func (sv *Server) listAddresses(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, mapErrResp(fmt.Errorf("authorization payload is not provided")))
		return
	}

	addresses, err := sv.postgres.GetAddresses(c, authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	addressesResponse := make([]AddressResponse, len(addresses))
	for i, address := range addresses {
		addressesResponse[i] = mapAddressResponse(address)
	}

	c.JSON(http.StatusOK, GenericResponse[[]AddressResponse]{&addressesResponse, nil, nil})
}

// updateAddress godoc
// @Summary Update an address
// @Description Update an address
// @Tags address
// @Accept json
// @Produce json
// @Param id path int true "Address ID"
// @Param input body UpdateAddressParams true "Update Address"
// @Success 200 {object} GenericResponse[AddressResponse]
// @Router /address/{id} [put]
func (sv *Server) updateAddress(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, mapErrResp(fmt.Errorf("authorization payload is not provided")))
		return
	}
	var input UpdateAddressParams
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	var param GetAddressParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	_, err := sv.postgres.GetAddress(c, sqlc.GetAddressParams{
		ID:     param.ID,
		UserID: authPayload.UserID,
	})
	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(err))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	payload := sqlc.UpdateAddressParams{
		ID:     param.ID,
		UserID: authPayload.UserID,
	}

	if input.Phone != nil {
		payload.Phone = util.GetPgTypeText(*input.Phone)
	}

	if input.City != nil {
		payload.City = util.GetPgTypeText(*input.City)
	}

	if input.Address1 != nil {
		payload.Street = util.GetPgTypeText(*input.Address1)
	}

	if input.Ward != nil {
		payload.Ward = util.GetPgTypeText(*input.Ward)
	}

	if input.District != nil {
		payload.District = util.GetPgTypeText(*input.District)
	}

	if input.IsDefault != nil {
		if *input.IsDefault {
			err := sv.postgres.SetPrimaryAddressTx(c, postgres.SetPrimaryAddressTxParams{
				NewPrimaryID: param.ID,
				UserID:       authPayload.UserID,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, mapErrResp(fmt.Errorf("failed to set primary address: %w", err)))
				return
			}
		}
	}

	address, err := sv.postgres.UpdateAddress(c, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	addressDetail := mapAddressResponse(address)
	c.JSON(http.StatusOK, GenericResponse[AddressResponse]{&addressDetail, nil, nil})
}

// removeAddress godoc
// @Summary Remove an address
// @Description Remove an address
// @Tags address
// @Accept json
// @Produce json
// @Param id path int true "Address ID"
// @Success 200 {object} GenericResponse[bool]
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /address/{id} [delete]
func (sv *Server) removeAddress(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, mapErrResp(fmt.Errorf("authorization payload is not provided")))
		return
	}
	var param GetAddressParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	address, err := sv.postgres.GetAddress(c, sqlc.GetAddressParams{
		ID:     param.ID,
		UserID: authPayload.UserID,
	})

	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(err))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if address.IsDeleted {
		c.JSON(http.StatusNotFound, mapErrResp(fmt.Errorf("address has been removed")))
		return
	}

	err = sv.postgres.DeleteAddress(c, sqlc.DeleteAddressParams{
		ID:     param.ID,
		UserID: authPayload.UserID,
	})
	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(err))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	status := true
	c.JSON(http.StatusOK, GenericResponse[bool]{&status, nil, nil})
}
