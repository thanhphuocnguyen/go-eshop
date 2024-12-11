package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhphuocnguyen/go-eshop/internal/auth"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/sqlc"
)

type orderListQuery struct {
	Limit  int32 `form:"limit,default=10" binding:"required,min=5,max=100"`
	Offset int32 `form:"offset,default=0" binding:"required,min=0"`
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

// @Summary List orders
// @Description List orders of the current user
// @Tags orders
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Security ApiKeyAuth
// @Success 200 {object} ListOrdersResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
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
		Limit:  orderListQuery.Limit,
		Offset: orderListQuery.Offset,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, orders)
}

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

func (sv *Server) changeOrderStatus(c *gin.Context) {

}
func (sv *Server) changeOrderPaymentStatus(c *gin.Context) {

}
