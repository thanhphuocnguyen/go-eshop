package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhphuocnguyen/go-eshop/internal/auth"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/util"
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
type paymentInfo struct {
	ID             int32   `json:"id"`
	PaymentAmount  float64 `json:"payment_amount"`
	PaymentGateWay *string `json:"payment_gateway"`
	PaymentMethod  string  `json:"payment_method"`
	PaymentStatus  string  `json:"payment_status"`
}

type orderDetailResponse struct {
	ID          int64                  `json:"id"`
	Total       float64                `json:"total"`
	Status      repository.OrderStatus `json:"status"`
	PaymentInfo paymentInfo            `json:"payment_info"`
	Products    []orderItemResp        `json:"products"`
}
type orderListResp struct {
	ID            int64                    `json:"id"`
	Total         float64                  `json:"total"`
	TotalItems    int32                    `json:"total_items"`
	Status        repository.OrderStatus   `json:"status"`
	PaymentStatus repository.PaymentStatus `json:"payment_status"`
	CreatedAt     time.Time                `json:"created_at"`
	UpdatedAt     time.Time                `json:"updated_at"`
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

	listOrderRows, err := sv.repo.ListOrders(c, repository.ListOrdersParams{
		UserID: user.UserID,
		Limit:  orderListQuery.PageSize,
		Offset: (orderListQuery.Page - 1) * orderListQuery.PageSize,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	count, err := sv.repo.CountOrders(c, repository.CountOrdersParams{UserID: user.UserID})

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

	getOrderDetailRows, err := sv.repo.GetOrderDetails(c, params.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if len(getOrderDetailRows) == 0 {
		c.JSON(http.StatusNotFound, mapErrResp(errors.New("order not found")))
		return
	}

	orderDetailRow := getOrderDetailRows[0]
	totalGet := util.PgNumericToFloat64(orderDetailRow.TotalPrice)

	paymentInfo := paymentInfo{
		ID:            int32(orderDetailRow.PaymentID.Int32),
		PaymentAmount: util.PgNumericToFloat64(orderDetailRow.PaymentAmount),
		PaymentMethod: string(orderDetailRow.PaymentMethod.PaymentMethod),
		PaymentStatus: string(orderDetailRow.PaymentStatus.PaymentStatus),
	}
	if orderDetailRow.PaymentGateway.Valid {
		paymentInfo.PaymentGateWay = (*string)(&orderDetailRow.PaymentGateway.PaymentGateway)
	}

	orderDetail := &orderDetailResponse{
		ID:          orderDetailRow.OrderID,
		Total:       totalGet,
		Status:      orderDetailRow.Status,
		PaymentInfo: paymentInfo,
		Products:    []orderItemResp{},
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
// @Success 200 {object} GenericResponse[repository.Order]
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
	order, err := sv.repo.GetOrder(c, params.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if order.UserID != user.UserID {
		c.JSON(http.StatusUnauthorized, mapErrResp(errors.New("user does not have permission")))
		return
	}

	// if order
	err = sv.repo.CancelOrderTx(c, repository.CancelOrderTxParams{OrderID: params.ID})

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	c.JSON(http.StatusOK, GenericResponse[repository.Order]{&order, nil, nil})
}

// @Summary Change order status
// @Description Change order status by order ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param status body string true "Status"
// @Security ApiKeyAuth
// @Success 200 {object} GenericResponse[repository.Order]
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
	order, err := sv.repo.GetOrder(c, params.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	if order.UserID != user.UserID {
		c.JSON(http.StatusUnauthorized, mapErrResp(errors.New("user does not have permission")))
		return
	}
	rs, err := sv.repo.UpdateOrder(c, repository.UpdateOrderParams{
		Status: repository.NullOrderStatus{
			OrderStatus: repository.OrderStatus(req.Status),
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusOK, GenericResponse[repository.Order]{&rs, nil, nil})
}
