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
	UserID    int64  `json:"user_id" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	Address1  string `json:"address_1" binding:"required"`
	Address2  string `json:"address_2"`
	Ward      string `json:"ward" binding:"required"`
	District  string `json:"district" binding:"required"`
	City      string `json:"city" binding:"required"`
	IsDefault bool   `json:"is_default"`
}

type UpdateAddressParams struct {
	UserID    int64  `json:"user_id" binding:"required"`
	Phone     string `json:"phone"`
	Address1  string `json:"address_1"`
	Address2  string `json:"address_2"`
	Ward      string `json:"ward"`
	District  string `json:"district"`
	City      string `json:"city"`
	IsDefault bool   `json:"is_default"`
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
	var input CreateAddressParams
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	payload := sqlc.CreateAddressParams{
		UserID:   input.UserID,
		Phone:    input.Phone,
		Address1: input.Address1,
		City:     input.City,
	}
	if input.Address2 != "" {
		payload.Address2 = util.GetPgTypeText(input.Address2)
	}
	if input.Ward != "" {
		payload.Ward = util.GetPgTypeText(input.Ward)
	}
	if input.District != "" {
		payload.District = util.GetPgTypeText(input.District)
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

	payload := sqlc.UpdateAddressParams{
		ID:     param.ID,
		UserID: input.UserID,
	}
	if input.Phone != "" {
		payload.Phone = util.GetPgTypeText(input.Phone)
	}
	if input.City != "" {
		payload.City = util.GetPgTypeText(input.City)
	}
	if input.Address1 != "" {
		payload.Address1 = util.GetPgTypeText(input.Address1)
	}
	if input.Address2 != "" {
		payload.Address2 = util.GetPgTypeText(input.Address2)
	}
	if input.Ward != "" {
		payload.Ward = util.GetPgTypeText(input.Ward)
	}
	if input.District != "" {
		payload.District = util.GetPgTypeText(input.District)
	}
	if input.IsDefault {
		payload.IsPrimary = util.GetPgTypeBool(true)
	}

	address, err = sv.postgres.UpdateAddress(c, payload)
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

	err := sv.postgres.DeleteAddress(c, sqlc.DeleteAddressParams{
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
	c.JSON(http.StatusOK, nil)
}
