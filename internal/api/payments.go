package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"github.com/thanhphuocnguyen/go-eshop/pkg/payment"
)

// Setup payment-related routes
func (s *Server) addPaymentRoutes(r chi.Router) {
	r.Route("/payments", func(r chi.Router) {
		r.Get("/{id}", s.getPayment)
		if s.config.Env == DevEnv {
			r.Post("/{id}/confirm", s.confirmPayment)
		}
		r.Get("/methods", s.getPaymentMethods)
		r.Get("/stripe-config", s.getStripeConfig)
		r.Post("/", s.createPaymentIntent)
		r.Put("/{orderId}", s.changePaymentStatus)
	})
}

func (s *Server) getStripeConfig(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{
		"public_key": s.config.StripePublishableKey,
	}
	RespondSuccess(w, resp)
}

// @Summary Initiate payment
// @Description Initiate payment
// @Tags payment
// @Accept json
// @Produce json
// @Param request body models.PaymentModel true "Payment request"
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse[dto.PaymentIntentSecret]
// @Failure 400 {object} dto.ErrorResp
// @Failure 401 {object} dto.ErrorResp
// @Failure 403 {object} dto.ErrorResp
// @Failure 404 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /payments [post]
func (s *Server) createPaymentIntent(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	_, claims, err := jwtauth.FromContext(c)
	if err != nil {
		RespondUnauthorized(w, UnauthorizedCode, err)
	}
	userID := uuid.MustParse(claims["userId"].(string))
	user, err := s.repo.GetUserByID(c, userID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, err)
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	var req models.PaymentModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	ord, err := s.repo.GetOrder(c, uuid.MustParse(req.OrderID))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, err)
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, errors.New("order not found"))
		return
	}

	if ord.UserID != user.ID {
		RespondForbidden(w, PermissionDeniedCode, errors.New("permission denied"))
		return
	}

	pmRow, err := s.repo.GetPaymentByOrderID(c, ord.ID)
	if err != nil && !errors.Is(err, repository.ErrRecordNotFound) {
		RespondInternalServerError(w, InternalServerErrorCode, errors.New("order not found"))
		return
	}

	if pmRow.ID != uuid.Nil && pmRow.Status != repository.PaymentStatusCancelled {
		RespondBadRequest(w, InvalidPaymentCode, errors.New("payment already exists"))
		return
	}

	total, _ := ord.TotalPrice.Float64Value()
	paymentMethodId := uuid.MustParse(req.PaymentMethodID)
	// create new payment
	createPaymentParams := repository.CreatePaymentParams{
		OrderID:         ord.ID,
		Amount:          utils.GetPgNumericFromFloat(total.Float64),
		PaymentMethodID: paymentMethodId,
	}
	var resp dto.PaymentIntentSecret
	paymentMethod, err := s.repo.GetPaymentMethodByID(c, paymentMethodId)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, err)
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	intent, err := s.paymentSrv.CreatePaymentIntent(c, paymentMethod.Code, payment.PaymentRequest{
		Amount:      int64(total.Float64 * 100), // convert to smallest currency unit
		Currency:    payment.USD,
		Email:       user.Email,
		Description: "Payment for order " + ord.ID.String(),
		Metadata:    map[string]string{"order_id": ord.ID.String()},
	})
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	createPaymentParams.PaymentIntentID = &intent.ID
	resp.ClientSecret = &intent.ClientSecret

	pmRow, err = s.repo.CreatePayment(c, createPaymentParams)

	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	resp.PaymentID = pmRow.ID.String()
	RespondSuccess(w, resp)
}

// @Summary Get payment  by order ID
// @Description Get payment  by order ID
// @Tags payment
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse[dto.PaymentDetail]
// @Failure 400 {object} dto.ErrorResp
// @Failure 401 {object} dto.ErrorResp
// @Failure 403 {object} dto.ErrorResp
// @Failure 404 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /payments/{id} [get]
func (s *Server) getPayment(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id := chi.URLParam(r, "id")
	if id == "" {
		RespondBadRequest(w, InvalidBodyCode, errors.New("missing id parameter"))
		return
	}

	payment, err := s.repo.GetPaymentByID(c, uuid.MustParse(id))
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	var details interface{}
	paymentMethod, err := s.repo.GetPaymentMethodByID(c, payment.PaymentMethodID)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	details, err = s.paymentSrv.GetPayment(c, *payment.PaymentIntentID, paymentMethod.Code)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	resp := dto.PaymentDetail{
		ID:      payment.ID.String(),
		Gateway: payment.Gateway,
		Status:  payment.Status,
		Details: details,
	}

	RespondSuccess(w, resp)
}

// @Summary Change payment status
// @Description Change payment status
// @Tags payment
// @Accept json
// @Produce json
// @Param paymentId path string true "Payment ID"
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse[dto.PaymentDetail]
// @Failure 400 {object} dto.ErrorResp
// @Failure 401 {object} dto.ErrorResp
// @Failure 403 {object} dto.ErrorResp
// @Failure 404 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /payments/{paymentId} [get]
func (s *Server) changePaymentStatus(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	orderId := chi.URLParam(r, "orderId")
	if orderId == "" {
		RespondBadRequest(w, InvalidBodyCode, errors.New("missing orderId parameter"))
		return
	}
	var req models.UpdatePaymentStatusModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, errors.New("invalid request body"))
		return
	}
	payment, err := s.repo.GetPaymentByID(c, uuid.MustParse(orderId))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, errors.New("order not found"))
		}
		RespondInternalServerError(w, InternalServerErrorCode, errors.New("order not found"))
		return
	}
	order, err := s.repo.GetOrder(c, payment.OrderID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, err)
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	if order.Status != repository.OrderStatusPending {
		RespondBadRequest(w, InvalidPaymentCode, errors.New("order is not pending"))
		return
	}

	if payment.Gateway != nil {
		RespondBadRequest(w, InvalidPaymentCode, errors.New("cannot change payment status for stripe payment"))
		return
	}

	if req.Status == repository.PaymentStatusCancelled {
		if payment.Status == repository.PaymentStatusSuccess {
			RespondBadRequest(w, InvalidPaymentCode, errors.New("cannot cancel payment that is already success"))
			return
		}
		if payment.Status == repository.PaymentStatusCancelled {
			RespondBadRequest(w, InvalidPaymentCode, errors.New("payment is already cancelled"))
			return
		}
	}

	if req.Status == repository.PaymentStatusSuccess {
		_, err := s.repo.UpdateOrder(c, repository.UpdateOrderParams{
			ID: order.ID,
			Status: repository.NullOrderStatus{
				OrderStatus: repository.OrderStatusDelivered,
				Valid:       true,
			},
			DeliveredAt: utils.GetPgTypeTimestamp(time.Now()),
		})
		if err != nil {
			RespondInternalServerError(w, InternalServerErrorCode, err)
			return
		}

	}

	err = s.repo.UpdatePayment(c, repository.UpdatePaymentParams{
		ID: payment.ID,
		Status: repository.NullPaymentStatus{
			PaymentStatus: req.Status,
			Valid:         true,
		},
	})
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	resp := dto.PaymentDetail{
		ID:     payment.ID.String(),
		Status: req.Status,
	}
	RespondSuccess(w, resp)
}

// @Summary Get payment methods
// @Description Get payment methods
// @Tags payment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse[[]dto.PaymentMethodResponse]
// @Failure 400 {object} dto.ErrorResp
// @Failure 401 {object} dto.ErrorResp
// @Failure 403 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /payments/methods [get]
func (s *Server) getPaymentMethods(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	paymentMethods, err := s.repo.ListPaymentMethods(c)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	var resp []dto.PaymentMethodResponse
	for _, pm := range paymentMethods {
		resp = append(resp, dto.PaymentMethodResponse{
			ID:   pm.ID.String(),
			Name: pm.Name,
			Code: pm.Code,
		})
	}
	RespondSuccess(w, resp)
}

// @Summary Confirm Payment
// @Description Confirm Payment
// @Tags payment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse[payment.PaymentResult]
// @Failure 400 {object} dto.ErrorResp
// @Failure 401 {object} dto.ErrorResp
// @Failure 403 {object} dto.ErrorResp
// @Failure 404 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /payments/{id}/confirm [post]
func (s *Server) confirmPayment(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id := chi.URLParam(r, "id")
	if id == "" {
		RespondBadRequest(w, InvalidBodyCode, errors.New("missing id parameter"))
		return
	}

	payment, err := s.repo.GetPaymentByID(c, uuid.MustParse(id))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, err)
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	paymentMethod, err := s.repo.GetPaymentMethodByID(c, payment.PaymentMethodID)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	rs, err := s.paymentSrv.ConfirmPayment(c, *payment.PaymentIntentID, paymentMethod.Code)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondSuccess(w, rs)
}
