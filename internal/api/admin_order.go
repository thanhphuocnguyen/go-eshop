package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
	"github.com/thanhphuocnguyen/go-eshop/pkg/payment"
)

// @Summary Cancel order
// @Description Cancel order by order ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse[uuid.UUID]
// @Failure 400 {object} dto.ErrorResp
// @Failure 401 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /admin/orders/{orderId}/cancel [put]
func (s *Server) adminCancelOrder(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	_, claims, err := jwtauth.FromContext(c)
	if err != nil {
		RespondInternalServerError(w, UnauthorizedCode, fmt.Errorf("authorization payload is not provided"))
		return
	}

	id, err := GetUrlParam(r, "id")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	var req models.CancelOrderModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	order, err := s.repo.GetOrder(c, uuid.MustParse(id))
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	if order.Status == repository.OrderStatusCancelled || order.Status == repository.OrderStatusRefunded {
		RespondBadRequest(w, InvalidPaymentCode, errors.New("order is already cancelled or refunded"))
		return
	}

	userRole := claims["role"].(string)
	userID := uuid.MustParse(claims["userId"].(string))

	if order.UserID != userID && userRole != "admin" {
		RespondForbidden(w, PermissionDeniedCode, errors.New("you do not have permission to access this order"))
		return
	}

	paymentRow, err := s.repo.GetPaymentByOrderID(c, order.ID)
	if err != nil && !errors.Is(err, repository.ErrRecordNotFound) {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	// if order status is not pending or user is not admin
	if order.Status != repository.OrderStatusPending || (!errors.Is(err, repository.ErrRecordNotFound) && paymentRow.Status != repository.PaymentStatusPending) {
		RespondBadRequest(w, PermissionDeniedCode, errors.New("order cannot be cancelled"))
		return
	}

	// if order
	cancelOrderTxParams := repository.CancelOrderTxArgs{
		OrderID: uuid.MustParse(id),
		CancelPaymentFromMethod: func(paymentID string, method string) error {
			req := payment.RefundRequest{
				TransactionID: paymentID,
				Amount:        paymentRow.Amount.Int.Int64(),
			}
			_, err = s.paymentSrv.RefundPayment(c, req, *paymentRow.Gateway)
			return err
		},
	}
	ordId, err := s.repo.CancelOrderTx(c, cancelOrderTxParams)

	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	s.cacheSrv.Delete(c, "order_detail:"+id)
	RespondSuccess(w, ordId)
}

// @Summary Refund order
// @Description Refund order by order ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse[repository.GetOrderRow]
// @Failure 400 {object} dto.ErrorResp
// @Failure 401 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /admin/order/{orderId}/refund [put]
func (s *Server) adminRefundOrder(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id, err := GetUrlParam(r, "id")
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	var req models.RefundOrderModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	order, err := s.repo.GetOrder(c, uuid.MustParse(id))
	if err != nil {
		if err == repository.ErrRecordNotFound {
			RespondNotFound(w, NotFoundCode, fmt.Errorf("order with ID %s not found", id))
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	if order.Status != repository.OrderStatusDelivered {
		RespondBadRequest(w, InvalidPaymentCode, errors.New("order cannot be refunded"))
		return
	}

	err = s.repo.RefundOrderTx(c, repository.RefundOrderTxArgs{
		OrderID: uuid.MustParse(id),
		RefundPaymentFromMethod: func(paymentID string, method string) (string, error) {
			req := payment.RefundRequest{
				TransactionID: paymentID,
				Amount:        order.TotalPrice.Int.Int64(),
				Reason:        req.Reason,
			}
			rs, err := s.paymentSrv.RefundPayment(c, req, method)
			return rs.Reason, err
		},
	})

	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	s.cacheSrv.Delete(c, "order_detail:"+id)

	RespondSuccess(w, order)
}

// @Summary Get all orders (admin endpoint)
// @Description Get all orders with pagination and filtering
// @Tags admin
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Param status query string false "Filter by status"
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse[[]dto.OrderListItem]
// @Failure 401 {object} ErrorResp
// @Failure 403 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/orders [get]
func (s *Server) adminGetOrders(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	var orderListQuery models.PaginationQuery = ParsePaginationQuery(r)
	status := r.URL.Query().Get("status")
	paymentStatus := r.URL.Query().Get("paymentStatus")

	dbParams := repository.GetOrdersParams{
		Limit:  orderListQuery.PageSize,
		Offset: (orderListQuery.Page - 1) * orderListQuery.PageSize,
	}
	countParams := repository.CountOrdersParams{}

	if status != "" {
		dbParams.Status = repository.NullOrderStatus{
			OrderStatus: repository.OrderStatus(status),
			Valid:       true,
		}
		countParams.Status = repository.NullOrderStatus{
			OrderStatus: repository.OrderStatus(status),
			Valid:       true,
		}
	}

	if paymentStatus != "" {
		dbParams.PaymentStatus = repository.NullPaymentStatus{
			PaymentStatus: repository.PaymentStatus(paymentStatus),
			Valid:         true,
		}
		countParams.PaymentStatus = repository.NullPaymentStatus{
			PaymentStatus: repository.PaymentStatus(paymentStatus),
			Valid:         true,
		}
	}

	fetchedOrderRows, err := s.repo.GetOrders(c, dbParams)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	count, err := s.repo.CountOrders(c, countParams)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	var orderResponses []dto.OrderListItem
	for _, aggregated := range fetchedOrderRows {
		// Convert PaymentStatus interface{} to PaymentStatus type
		paymentStatus := repository.PaymentStatusPending
		if aggregated.PaymentStatus.Valid {
			paymentStatus = aggregated.PaymentStatus.PaymentStatus
		}

		total, _ := aggregated.TotalPrice.Float64Value()
		orderResponses = append(orderResponses, dto.OrderListItem{
			ID:            aggregated.ID,
			Total:         total.Float64,
			TotalItems:    int32(aggregated.TotalItems),
			Status:        aggregated.Status,
			CustomerName:  aggregated.CustomerName,
			CustomerEmail: aggregated.CustomerEmail,
			PaymentStatus: paymentStatus,
			CreatedAt:     aggregated.CreatedAt.UTC(),
			UpdatedAt:     aggregated.UpdatedAt.UTC(),
		})
	}

	RespondSuccessWithPagination(w, orderResponses, dto.CreatePagination(orderListQuery.Page, orderListQuery.PageSize, count))
}

// @Summary Get order details by ID (admin endpoint)
// @Description Get detailed information about an order by its ID
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse[dto.OrderDetail]
// @Failure 401 {object} ErrorResp
// @Failure 403 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/orders/{id} [get]
func (s *Server) adminGetOrderDetail(w http.ResponseWriter, r *http.Request) {
	// Reuse the existing order detail handler since admin has access to all orders
	s.getOrderDetail(w, r)
}

// @Summary Change order status
// @Description Change order status by order ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param status body string true "Status"
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse[uuid.UUID]
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /admin/orders/{orderId}/status [put]
func (s *Server) adminChangeOrderStatus(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	id, err := GetUrlParam(r, "id")
	var req models.OrderStatusModel
	if err := s.GetRequestBody(r, &req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}
	order, err := s.repo.GetOrder(c, uuid.MustParse(id))
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, err)
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	if order.Status == repository.OrderStatusDelivered || order.Status == repository.OrderStatusCancelled || order.Status == repository.OrderStatusRefunded {
		RespondBadRequest(w, InvalidPaymentCode, errors.New("order cannot be changed"))
		return
	}

	status := repository.OrderStatus(req.Status)

	updateParams := repository.UpdateOrderParams{
		ID: order.ID,
		Status: repository.NullOrderStatus{
			OrderStatus: status,
			Valid:       true,
		},
	}
	if status == repository.OrderStatusConfirmed {
		updateParams.ConfirmedAt = pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		}
	}
	if status == repository.OrderStatusDelivering {
		updateParams.DeliveredAt = pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		}
	}

	rs, err := s.repo.UpdateOrder(c, updateParams)

	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	if err := s.cacheSrv.Delete(c, "order_detail:"+id); err != nil {
		log.Err(err).Msg("failed to delete order detail cache")
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	s.cacheSrv.Delete(c, "order_detail:"+id)

	RespondSuccess(w, rs)
}
