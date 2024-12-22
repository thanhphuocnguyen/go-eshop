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
	Phone     string  `json:"phone" binding:"required,min=10,max=15"`
	Address1  string  `json:"address_1" binding:"required"`
	Address2  *string `json:"address_2,omitempty" binding:"omitempty"`
	Ward      string  `json:"ward,omitempty" binding:"omitempty"`
	District  string  `json:"district" binding:"required"`
	City      string  `json:"city" binding:"required"`
	IsDefault bool    `json:"is_default,omitempty" binding:"omitempty"`
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

type GetAddressListQuery struct {
	Page     int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=100"`
}

// ------------------------------ Handlers ------------------------------

// createAddress godoc
// @Summary Create a new address
// @Description Create a new address
// @Tags address
// @Accept json
// @Produce json
// @Param input body CreateAddressParams true "Create Address"
// @Success 200 {object} sqlc.UserAddress
// @Router /address [post]
func (sv *Server) createAddress(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("authorization payload is not provided")))
		return
	}

	var input CreateAddressParams
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	payload := sqlc.CreateAddressParams{
		UserID:   authPayload.UserID,
		Phone:    input.Phone,
		Address1: input.Address1,
		City:     input.City,
		District: input.District,
	}
	if input.Address2 != nil {
		payload.Address2 = util.GetPgTypeText(*input.Address2)
	}
	if input.Ward != "" {
		payload.Ward = util.GetPgTypeText(input.Ward)
	}
	address, err := sv.postgres.CreateAddress(c, payload)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, address)
}

// getAddresses godoc
// @Summary Get list of addresses
// @Description Get list of addresses
// @Tags address
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} []sqlc.UserAddress
// @Router /address [get]
func (sv *Server) listAddresses(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("authorization payload is not provided")))
		return
	}
	var query GetAddressListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	addresses, err := sv.postgres.ListAddresses(c, sqlc.ListAddressesParams{
		UserID: authPayload.UserID,
		Limit:  query.PageSize,
		Offset: (query.Page - 1) * query.PageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, addresses)
}

// updateAddress godoc
// @Summary Update an address
// @Description Update an address
// @Tags address
// @Accept json
// @Produce json
// @Param id path int true "Address ID"
// @Param input body UpdateAddressParams true "Update Address"
// @Success 200 {object} sqlc.UserAddress
// @Router /address/{id} [put]
func (sv *Server) updateAddress(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("authorization payload is not provided")))
		return
	}
	var input UpdateAddressParams
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var param GetAddressParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := sv.postgres.GetAddress(c, sqlc.GetAddressParams{
		ID:     param.ID,
		UserID: authPayload.UserID,
	})
	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
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
		payload.Address1 = util.GetPgTypeText(*input.Address1)
	}
	if input.Address2 != nil {
		payload.Address2 = util.GetPgTypeText(*input.Address2)
	}
	if input.Ward != nil {
		payload.Ward = util.GetPgTypeText(*input.Ward)
	}
	if input.District != nil {
		payload.District = util.GetPgTypeText(*input.District)
	}
	if input.IsDefault != nil {
		payload.IsPrimary = util.GetPgTypeBool(*input.IsDefault)
	}

	address, err := sv.postgres.UpdateAddress(c, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, address)
}

// removeAddress godoc
// @Summary Remove an address
// @Description Remove an address
// @Tags address
// @Accept json
// @Produce json
// @Param id path int true "Address ID"
// @Success 200
// @Router /address/{id} [delete]
func (sv *Server) removeAddress(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("authorization payload is not provided")))
		return
	}
	var param GetAddressParams
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	address, err := sv.postgres.GetAddress(c, sqlc.GetAddressParams{
		ID:     param.ID,
		UserID: authPayload.UserID,
	})

	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if address.IsDeleted {
		c.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("address has been removed")))
		return
	}

	err = sv.postgres.DeleteAddress(c, sqlc.DeleteAddressParams{
		ID:     param.ID,
		UserID: authPayload.UserID,
	})
	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, responseMapper(struct {
		Message string `json:"message"`
		Success bool   `json:"success"`
	}{
		Message: "Address has been removed",
		Success: true,
	}, nil))
}
