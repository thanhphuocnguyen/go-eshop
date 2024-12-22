package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhphuocnguyen/go-eshop/internal/auth"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/sqlc"
)

// ---------------------------------------------- API Models ----------------------------------------------
type orderListQuery struct {
	Page     int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=100"`
}

type orderIDParams struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type orderDetailResponse struct {
	ID            int64              `json:"id"`
	Total         float64            `json:"total"`
	Status        sqlc.OrderStatus   `json:"status"`
	PaymentType   sqlc.PaymentType   `json:"payment_type"`
	PaymentStatus sqlc.PaymentStatus `json:"payment_status"`
	Products      []productResponse  `json:"products"`
}
type orderResponse struct {
	ID            int64              `json:"id"`
	Total         float64            `json:"total"`
	TotalItems    int32              `json:"total_items"`
	Status        sqlc.OrderStatus   `json:"status"`
	PaymentType   sqlc.PaymentType   `json:"payment_type"`
	PaymentStatus sqlc.PaymentStatus `json:"payment_status"`
	CreatedAt     string             `json:"created_at"`
	UpdatedAt     string             `json:"updated_at"`
}
type listOrderResponse struct {
	Orders []orderResponse `json:"orders"`
}

//---------------------------------------------- API Handlers ----------------------------------------------

// @Summary List orders
// @Description List orders of the current user
// @Tags orders
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Security ApiKeyAuth
// @Success 200 {object} listOrderResponse
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /orders [get]
func (sv *Server) orderList(c *gin.Context) {
	user, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("user not found")))
		return
	}
	var orderListQuery orderListQuery
	if err := c.ShouldBindQuery(&orderListQuery); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	orders, err := sv.postgres.ListOrders(c, sqlc.ListOrdersParams{
		UserID: user.UserID,
		Limit:  orderListQuery.PageSize,
		Offset: (orderListQuery.Page - 1) * orderListQuery.PageSize,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	var orderResponses []orderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, orderResponse{
			ID:            order.Order.ID,
			Total:         float64(order.TotalPrice),
			TotalItems:    int32(order.TotalItems),
			Status:        order.Order.Status,
			PaymentType:   order.Order.PaymentType,
			PaymentStatus: order.Order.PaymentStatus,
			CreatedAt:     order.Order.CreatedAt.String(),
			UpdatedAt:     order.Order.UpdatedAt.String(),
		})
	}

	c.JSON(http.StatusOK, listOrderResponse{orderResponses})
}

// @Summary Get order detail
// @Description Get order detail by order ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Security ApiKeyAuth
// @Success 200 {object} orderDetailResponse
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /orders/{id} [get]
func (sv *Server) orderDetail(c *gin.Context) {
	var params orderIDParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	orderDetails, err := sv.postgres.GetOrderDetails(c, params.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	var orderDetailResponse orderDetailResponse
	orderDetailResponse.ID = orderDetails[0].Order.ID

	orderDetailResponse.PaymentStatus = orderDetails[0].Order.PaymentStatus
	orderDetailResponse.PaymentType = orderDetails[0].Order.PaymentType
	orderDetailResponse.Status = orderDetails[0].Order.Status
	var total float64
	for _, order := range orderDetails {
		price, _ := order.Product.Price.Float64Value()
		total += price.Float64 * float64(order.OrderItem.Quantity)
		orderDetailResponse.Products = append(orderDetailResponse.Products, productResponse{
			ID:          order.Product.ID,
			Name:        order.Product.Name,
			Price:       price.Float64,
			Description: order.Product.Description,
			Sku:         order.Product.Sku,
			ImageURL:    order.Product.ImageUrl.String,
			Stock:       order.Product.Stock,
			UpdatedAt:   order.Product.UpdatedAt.String(),
			CreatedAt:   order.Product.CreatedAt.String(),
		})
	}
	orderDetailResponse.Total = total
	c.JSON(http.StatusOK, orderDetailResponse)
}

// @Summary Cancel order
// @Description Cancel order by order ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Security ApiKeyAuth
// @Success 200 {object} sqlc.Order
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /orders/{id}/cancel [put]
func (sv *Server) cancelOrder(c *gin.Context) {
	user, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("user not found")))
		return
	}
	var params orderIDParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	order, err := sv.postgres.GetOrder(c, params.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if order.UserID != user.UserID {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("user does not have permission")))
		return
	}

	// if order
	order, err = sv.postgres.UpdateOrder(c, sqlc.UpdateOrderParams{
		ID: params.ID,
		Status: sqlc.NullOrderStatus{
			OrderStatus: sqlc.OrderStatusCancelled,
			Valid:       true,
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, order)
}

// @Summary Change order status
// @Description Change order status by order ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param status body string true "Status"
// @Security ApiKeyAuth
// @Success 200 {object} sqlc.Order
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /orders/{id}/status [put]
func (sv *Server) changeOrderStatus(c *gin.Context) {

}

// @Summary Change order payment status
// @Description Change order payment status by order ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param status body string true "Payment Status"
// @Security ApiKeyAuth
// @Success 200 {object} sqlc.Order
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /orders/{id}/payment [put]
func (sv *Server) changeOrderPaymentStatus(c *gin.Context) {

}
