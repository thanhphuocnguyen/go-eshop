package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhphuocnguyen/go-eshop/internal/auth"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/postgres"
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

type changeOrderStatusReq struct {
	Status string `json:"status" binding:"required,oneof=pending confirmed delivering delivered completed"`
}

type orderItemResp struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	ImageUrl *string `json:"image_url"`
	Quantity int32   `json:"quantity"`
}

type orderDetailResponse struct {
	ID            int64              `json:"id"`
	Total         float64            `json:"total"`
	Status        sqlc.OrderStatus   `json:"status"`
	PaymentStatus sqlc.PaymentStatus `json:"payment_status"`
	Products      []orderItemResp    `json:"products"`
}
type orderListResp struct {
	ID            int64              `json:"id"`
	Total         float64            `json:"total"`
	TotalItems    int32              `json:"total_items"`
	Status        sqlc.OrderStatus   `json:"status"`
	PaymentStatus sqlc.PaymentStatus `json:"payment_status"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
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
		c.JSON(http.StatusUnauthorized, mapErrResp(errors.New("user not found")))
		return
	}
	var orderListQuery orderListQuery
	if err := c.ShouldBindQuery(&orderListQuery); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	listOrderRows, err := sv.postgres.ListOrders(c, sqlc.ListOrdersParams{
		UserID: user.UserID,
		Limit:  orderListQuery.PageSize,
		Offset: (orderListQuery.Page - 1) * orderListQuery.PageSize,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	count, err := sv.postgres.CountOrders(c, sqlc.CountOrdersParams{UserID: user.UserID})

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	var orderResponses []orderListResp
	for _, aggregated := range listOrderRows {
		total, _ := aggregated.TotalPrice.Float64Value()
		orderResponses = append(orderResponses, orderListResp{
			ID:            aggregated.OrderID,
			Total:         total.Float64,
			TotalItems:    int32(aggregated.TotalItems),
			Status:        aggregated.Status,
			PaymentStatus: aggregated.PaymentStatus.PaymentStatus,
			CreatedAt:     aggregated.CreatedAt.UTC(),
			UpdatedAt:     aggregated.UpdatedAt.UTC(),
		})
	}

	c.JSON(http.StatusOK, GenericListResponse[orderListResp]{&orderResponses, &count, nil, nil})
}

// @Summary Get order detail
// @Description Get order detail by order ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Security ApiKeyAuth
// @Success 200 {object} GenericResponse[orderDetailResponse]
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /orders/{id} [get]
func (sv *Server) orderDetail(c *gin.Context) {
	var params orderIDParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	getOrderDetailRows, err := sv.postgres.GetOrderDetails(c, params.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if len(getOrderDetailRows) == 0 {
		c.JSON(http.StatusNotFound, mapErrResp(errors.New("order not found")))
		return
	}

	orderDetailRow := getOrderDetailRows[0]
	totalGet, _ := orderDetailRow.TotalPrice.Float64Value()
	orderDetail := &orderDetailResponse{
		ID:            orderDetailRow.OrderID,
		Total:         totalGet.Float64,
		Status:        orderDetailRow.Status,
		PaymentStatus: orderDetailRow.PaymentStatus.PaymentStatus,
		Products:      []orderItemResp{},
	}

	for _, item := range getOrderDetailRows {
		orderDetail.Products = append(orderDetail.Products, orderItemResp{
			ID:       item.ProductID.Int64,
			Name:     item.ProductName.String,
			ImageUrl: &item.ImageUrl.String,
			Quantity: item.Quantity.Int32,
		})
	}

	c.JSON(http.StatusOK, GenericResponse[orderDetailResponse]{orderDetail, nil, nil})
}

// @Summary Cancel order
// @Description Cancel order by order ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Security ApiKeyAuth
// @Success 200 {object} GenericResponse[sqlc.Order]
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /orders/{id}/cancel [put]
func (sv *Server) cancelOrder(c *gin.Context) {
	user, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusUnauthorized, mapErrResp(errors.New("user not found")))
		return
	}
	var params orderIDParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	order, err := sv.postgres.GetOrder(c, params.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if order.UserID != user.UserID {
		c.JSON(http.StatusUnauthorized, mapErrResp(errors.New("user does not have permission")))
		return
	}

	// if order
	err = sv.postgres.CancelOrderTx(c, postgres.CancelOrderTxParams{OrderID: params.ID})

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	c.JSON(http.StatusOK, GenericResponse[sqlc.Order]{&order, nil, nil})
}

// @Summary Change order status
// @Description Change order status by order ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param status body string true "Status"
// @Security ApiKeyAuth
// @Success 200 {object} GenericResponse[sqlc.Order]
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /orders/{id}/status [put]
func (sv *Server) changeOrderStatus(c *gin.Context) {
	user, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusUnauthorized, mapErrResp(errors.New("user not found")))
		return
	}
	var params orderIDParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	var req changeOrderStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	order, err := sv.postgres.GetOrder(c, params.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	if order.UserID != user.UserID {
		c.JSON(http.StatusUnauthorized, mapErrResp(errors.New("user does not have permission")))
		return
	}
	rs, err := sv.postgres.UpdateOrder(c, sqlc.UpdateOrderParams{
		Status: sqlc.NullOrderStatus{
			OrderStatus: sqlc.OrderStatus(req.Status),
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusOK, GenericResponse[sqlc.Order]{&rs, nil, nil})
}
