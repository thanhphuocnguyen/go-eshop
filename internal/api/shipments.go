package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

func (s *Server) registerShipmentRoutes(r chi.Router) {
	r.Route("shipping-methods", func(r chi.Router) {
		r.Post("/", s.handleCreateShippingMethod)
		r.Get("/{shippingMethodID}", s.handleGetShippingMethod)
		r.Put("/{shippingMethodID}", s.handleUpdateShippingMethod)
		r.Delete("/{shippingMethodID}", s.handleDeleteShippingMethod)
	})

	r.Route("shipping-zones", func(r chi.Router) {
		r.Post("/", s.handleCreateShippingZone)
		r.Get("/{shippingZoneID}", s.handleGetShippingZone)
		r.Put("/{shippingZoneID}", s.handleUpdateShippingZone)
		r.Delete("/{shippingZoneID}", s.handleDeleteShippingZone)
	})

	r.Route("shipping-rates", func(r chi.Router) {
		r.Post("/", s.handleCreateShippingRate)
		r.Get("/{shippingRateID}", s.handleGetShippingRate)
		r.Put("/{shippingRateID}", s.handleUpdateShippingRate)
		r.Delete("/{shippingRateID}", s.handleDeleteShippingRate)
	})

	r.Route("shipments", func(r chi.Router) {
		r.Post("/", s.handleCreateShipment)
		r.Get("/{shipmentID}", s.handleGetShipment)
		r.Put("/{shipmentID}", s.handleUpdateShipment)
		r.Delete("/{shipmentID}", s.handleDeleteShipment)
	})
}

func (s *Server) handleCreateShippingMethod(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	var req models.ShippingMethodModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	estDays := fmt.Sprintf("%d", *req.EstimatedDays)
	createdRow, err := s.repo.CreateShippingMethod(c, repository.CreateShippingMethodParams{
		Name:                  req.Name,
		Description:           req.Description,
		IsActive:              req.Active,
		RequiresAddress:       true,
		EstimatedDeliveryTime: &estDays,
	})
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	RespondSuccess(w, createdRow)
}

func (s *Server) handleGetShippingMethod(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id, err := GetUrlParam(r, "shippingMethodID")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	shippingMethod, err := s.repo.GetShippingMethodByID(c, uuid.MustParse(id))
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	RespondSuccess(w, shippingMethod)
}

func (s *Server) handleUpdateShippingMethod(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id, err := GetUrlParam(r, "shippingMethodID")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	var req models.ShippingMethodModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	estDays := fmt.Sprintf("%d", *req.EstimatedDays)
	updatedRow, err := s.repo.UpdateShippingMethod(c, repository.UpdateShippingMethodParams{
		ID:                    uuid.MustParse(id),
		Name:                  &req.Name,
		Description:           req.Description,
		IsActive:              &req.Active,
		RequiresAddress:       &req.RequiresAddress,
		EstimatedDeliveryTime: &estDays,
	})
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	RespondSuccess(w, updatedRow)
}

func (s *Server) handleDeleteShippingMethod(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id, err := GetUrlParam(r, "shippingMethodID")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	err = s.repo.DeleteShippingMethod(c, uuid.MustParse(id))
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	RespondSuccess(w, map[string]string{"message": "Shipping method deleted successfully"})
}

func (s *Server) handleCreateShippingZone(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	var req models.ShippingZoneModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	createdRow, err := s.repo.CreateShippingZone(c, repository.CreateShippingZoneParams{
		Name:        req.Name,
		Description: req.Description,
		Countries:   req.Countries,
	})
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	RespondSuccess(w, createdRow)
}

func (s *Server) handleGetShippingZone(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id, err := GetUrlParam(r, "shippingZoneID")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	shippingZone, err := s.repo.GetShippingZoneByID(c, uuid.MustParse(id))
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	RespondSuccess(w, shippingZone)
}

func (s *Server) handleUpdateShippingZone(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id, err := GetUrlParam(r, "shippingZoneID")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	var req models.ShippingZoneModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	updatedRow, err := s.repo.UpdateShippingZone(c, repository.UpdateShippingZoneParams{
		ID:          uuid.MustParse(id),
		Name:        &req.Name,
		Description: req.Description,
		Countries:   req.Countries,
	})
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	RespondSuccess(w, updatedRow)
}

func (s *Server) handleDeleteShippingZone(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id, err := GetUrlParam(r, "shippingZoneID")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	err = s.repo.DeleteShippingZone(c, uuid.MustParse(id))
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	RespondNoContent(w)
}

func (s *Server) handleCreateShippingRate(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	var req models.ShippingRateModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	createdRow, err := s.repo.CreateShippingRate(c, repository.CreateShippingRateParams{
		ShippingMethodID:      uuid.MustParse(req.ShippingMethodID),
		ShippingZoneID:        uuid.MustParse(req.ShippingZoneID),
		Name:                  req.Name,
		BaseRate:              utils.GetPgNumericFromFloat(req.Price),
		MinOrderAmount:        utils.GetPgNumericFromFloat(*req.MinOrderAmount),
		MaxOrderAmount:        utils.GetPgNumericFromFloat(*req.MaxOrderAmount),
		FreeShippingThreshold: utils.GetPgNumericFromFloat(*req.FreeShippingThreshold),
		IsActive:              req.IsActive,
	})
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	RespondSuccess(w, createdRow)
}

func (s *Server) handleGetShippingRate(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id, err := GetUrlParam(r, "shippingRateID")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	shippingRate, err := s.repo.GetShippingRateByID(c, uuid.MustParse(id))
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	RespondSuccess(w, shippingRate)
}

func (s *Server) handleUpdateShippingRate(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id, err := GetUrlParam(r, "shippingRateID")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	var req models.UpdateShippingRateModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	var price *float64
	if req.Price != nil {
		price = req.Price
	}
	updatedRow, err := s.repo.UpdateShippingRate(c, repository.UpdateShippingRateParams{
		ID:    uuid.MustParse(id),
		Price: price,
	})
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	RespondSuccess(w, updatedRow)
}

func (s *Server) handleDeleteShippingRate(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id, err := GetUrlParam(r, "shippingRateID")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	err = s.repo.DeleteShippingRate(c, uuid.MustParse(id))
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	RespondNoContent(w)
}

func (s *Server) handleCreateShipment(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	var req models.ShipmentModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	createdRow, err := s.repo.CreateShipment(c, repository.CreateShipmentParams{
		OrderID
Status
ShippedAt
DeliveredAt
TrackingNumber
TrackingUrl
ShippingProvider
ShippingNotes
	})
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	RespondSuccess(w, createdRow)
}

func (s *Server) handleGetShipment(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id, err := GetUrlParam(r, "shipmentID")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	shipment, err := s.repo.GetShipmentByID(c, uuid.MustParse(id))
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	RespondSuccess(w, shipment)
}

func (s *Server) handleUpdateShipment(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id, err := GetUrlParam(r, "shipmentID")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	var req models.UpdateShipmentModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	updatedRow, err := s.repo.UpdateShipment(c, repository.UpdateShipmentParams{
		ID:           uuid.MustParse(id),
		Name:         req.Name,
		Phone:        req.Phone,
		Address:      req.Address,
		City:         req.City,
		State:        req.State,
		Country:      req.Country,
		PostalCode:   req.PostalCode,
		Instructions: req.Instructions,
	})
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	RespondSuccess(w, updatedRow)
}

func (s *Server) handleDeleteShipment(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id, err := GetUrlParam(r, "shipmentID")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	err = s.repo.DeleteShipment(c, uuid.MustParse(id))
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	RespondNoContent(w)
}
