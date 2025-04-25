package api

import (
	"errors"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/auth"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"github.com/thanhphuocnguyen/go-eshop/pkg/payment"
)

type updateCartItemRequest struct {
	Quantity int16 `json:"quantity" binding:"required,gt=0"`
}

type ProductVariantParam struct {
	ID string `uri:"variant_id" binding:"required,uuid"`
}

type CartItemAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CartItemResponse struct {
	ID         string              `json:"id" binding:"required,uuid"`
	ProductID  string              `json:"product_id" binding:"required,uuid"`
	VariantID  string              `json:"variant_id" binding:"required,uuid"`
	Name       string              `json:"name"`
	Quantity   int16               `json:"quantity"`
	Price      float64             `json:"price"`
	Discount   int16               `json:"discount"`
	StockQty   int32               `json:"stock"`
	Sku        *string             `json:"sku,omitempty"`
	ImageURL   *string             `json:"image_url,omitempty"`
	Attributes []CartItemAttribute `json:"attributes"`
}

type CartDetailResponse struct {
	ID         uuid.UUID          `json:"id"`
	TotalPrice float64            `json:"total_price"`
	CartItems  []CartItemResponse `json:"cart_items"`
	UpdatedAt  time.Time          `json:"updated_at,omitempty"`
	CreatedAt  time.Time          `json:"created_at"`
}

type PutCartItemReq struct {
	VariantID string `json:"variant_id" binding:"required,uuid"`
	Quantity  int16  `json:"quantity" binding:"required,gt=0"`
}

type GetCartItemParam struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type CheckoutRequest struct {
	PaymentGateway *string `json:"payment_gateway" binding:"required,oneof=stripe paypal visa mastercard apple_pay google_pay postpaid momo zalo_pay vn_pay"`
	AddressID      *int64  `json:"address_id" binding:"omitempty,required"`
}

type CheckoutResponse struct {
	OrderID   uuid.UUID `json:"order_id"`
	PaymentID string    `json:"payment_id"`
}

// ------------------------------ Handlers ------------------------------

// @Summary Create a new cart
// @Schemes http
// @Description create a new cart for a user
// @Tags carts
// @Accept json
// @Produce json
// @Success 200 {object} ApiResponse[CartDetailResponse]
// @Failure 400 {object} ApiResponse[CartDetailResponse]
// @Failure 500 {object} ApiResponse[CartDetailResponse]
// @Router /cart [post]
func (sv *Server) createCart(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, createErrorResponse[CartDetailResponse](http.StatusBadRequest, "", errors.New("user not found")))
		return
	}
	user, err := sv.repo.GetUserByID(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[CartDetailResponse](http.StatusNotFound, "", errors.New("user not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[CartDetailResponse](http.StatusInternalServerError, "", err))
		return
	}
	_, err = sv.repo.GetCart(c, repository.GetCartParams{
		UserID: authPayload.UserID,
	})
	if err == nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[CartDetailResponse](http.StatusBadRequest, "", errors.New("cart already exists")))
		return
	}

	newCart, err := sv.repo.CreateCart(c, repository.CreateCartParams{
		ID:     uuid.New(),
		UserID: user.ID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[CartDetailResponse](http.StatusInternalServerError, "", err))
		return
	}
	resp := &CartDetailResponse{
		ID:         newCart.ID,
		TotalPrice: 0,
		CartItems:  []CartItemResponse{},
		UpdatedAt:  newCart.UpdatedAt,
		CreatedAt:  newCart.CreatedAt,
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, resp, "", nil, nil))
}

// @Summary Get cart details by user ID
// @Schemes http
// @Description get cart details by user ID
// @Tags cart
// @Accept json
// @Produce json
// @Success 200 {object} ApiResponse[CartDetailResponse]
// @Failure 500 {object} ApiResponse[CartDetailResponse]
// @Router /cart [get]
func (sv *Server) getCart(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, createErrorResponse[CartDetailResponse](http.StatusBadRequest, "", errors.New("user not found")))
		return
	}

	cart, err := sv.repo.GetCart(c, repository.GetCartParams{
		UserID: authPayload.UserID,
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			cart, err := sv.repo.CreateCart(c, repository.CreateCartParams{
				ID:     uuid.New(),
				UserID: authPayload.UserID,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, createErrorResponse[CartDetailResponse](http.StatusInternalServerError, "", err))
				return
			}
			c.JSON(http.StatusOK, createSuccessResponse(c, CartDetailResponse{
				ID:         cart.ID,
				TotalPrice: 0,
				CartItems:  []CartItemResponse{},
				UpdatedAt:  cart.UpdatedAt,
				CreatedAt:  cart.CreatedAt,
			}, "", nil, nil))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[CartDetailResponse](http.StatusInternalServerError, "", err))
		return
	}

	cartDetail := CartDetailResponse{
		ID:         cart.ID,
		TotalPrice: 0,
		CartItems:  []CartItemResponse{},
		UpdatedAt:  cart.UpdatedAt,
		CreatedAt:  cart.CreatedAt,
	}
	cartItemRows, err := sv.repo.GetCartItems(c, cart.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[CartDetailResponse](http.StatusInternalServerError, "", err))
		return
	}

	cartDetail.CartItems, cartDetail.TotalPrice = mapToCartItemsResp(cartItemRows)

	c.JSON(http.StatusOK, createSuccessResponse(c, cartDetail, "", nil, nil))
}

// @Summary Add a product to the cart
// @Schemes http
// @Description add a product to the cart
// @Tags carts
// @Accept json
// @Param input body PutCartItemReq true "Add product to cart input"
// @Produce json
// @Success 200 {object} ApiResponse[uuid.UUID]
// @Failure 400 {object} ApiResponse[uuid.UUID]
// @Failure 500 {object} ApiResponse[uuid.UUID]
// @Router /cart/item/{variant_id} [post]
func (sv *Server) putCartItemHandler(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, createErrorResponse[uuid.UUID](http.StatusBadRequest, "", errors.New("user not found")))
		return
	}

	var req PutCartItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[uuid.UUID](http.StatusInternalServerError, "", err))
		return
	}

	cart, err := sv.repo.GetCart(c, repository.GetCartParams{
		UserID: authPayload.UserID,
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[uuid.UUID](http.StatusNotFound, "", errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[uuid.UUID](http.StatusInternalServerError, "", err))
		return
	}

	cartItem, err := sv.repo.GetCartItemByProductVariantID(c, repository.GetCartItemByProductVariantIDParams{
		VariantID: uuid.MustParse(req.VariantID),
		CartID:    cart.ID,
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			cartItem, err = sv.repo.CreateCartItem(c, repository.CreateCartItemParams{
				ID:        uuid.New(),
				CartID:    cart.ID,
				VariantID: uuid.MustParse(req.VariantID),
				Quantity:  req.Quantity,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, createErrorResponse[uuid.UUID](http.StatusInternalServerError, "", err))
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, createErrorResponse[uuid.UUID](http.StatusInternalServerError, "", err))
			return
		}
	} else {
		err = sv.repo.UpdateCartItemQuantity(c, repository.UpdateCartItemQuantityParams{
			Quantity: cartItem.Quantity + req.Quantity,
			ID:       cartItem.ID,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse[uuid.UUID](http.StatusInternalServerError, "", err))
			return
		}
	}

	err = sv.repo.UpdateCartTimestamp(c, cart.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[uuid.UUID](http.StatusInternalServerError, "", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, cartItem.ID, "Success!", nil, nil))
}

// @Summary Remove a product from the cart
// @Schemes http
// @Description remove a product from the cart
// @Tags carts
// @Accept json
// @Param id path int true "Product ID"
// @Produce json
// @Success 200 {object} ApiResponse[string]
// @Failure 400 {object} ApiResponse[string]
// @Failure 500 {object} ApiResponse[string]
// @Router /cart/item/{id} [delete]
func (sv *Server) removeCartItem(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, createErrorResponse[string](http.StatusBadRequest, "", errors.New("user not found")))
		return
	}

	var param GetCartItemParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[string](http.StatusInternalServerError, "", err))
		return
	}

	cart, err := sv.repo.GetCart(c, repository.GetCartParams{
		UserID: authPayload.UserID,
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[string](http.StatusNotFound, "", errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[string](http.StatusInternalServerError, "", err))
		return
	}

	if cart.UserID != authPayload.UserID {
		c.JSON(http.StatusForbidden, createErrorResponse[string](http.StatusForbidden, "", errors.New("user not found")))
		return
	}

	err = sv.repo.RemoveProductFromCart(c, repository.RemoveProductFromCartParams{
		CartID: cart.ID,
		ID:     uuid.MustParse(param.ID),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[string](http.StatusInternalServerError, "", err))
		return
	}

	message := "item removed"
	c.JSON(http.StatusOK, createSuccessResponse(c, message, "", nil, nil))
}

// @Summary Update product items in the cart
// @Schemes http
// @Description update product items in the cart
// @Tags carts
// @Accept json
// @Param input body CheckoutRequest true "Checkout input"
// @Produce json
// @Success 200 {object} ApiResponse[CheckoutResponse]
// @Failure 400 {object} ApiResponse[CheckoutResponse]
// @Failure 500 {object} ApiResponse[CheckoutResponse]
// @Router /cart/checkout [post]
func (sv *Server) checkout(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, createErrorResponse[CheckoutResponse](http.StatusBadRequest, "", errors.New("user not found")))
		return
	}

	user, err := sv.repo.GetUserByID(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[CheckoutResponse](http.StatusBadRequest, "", errors.New("user not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[CheckoutResponse](http.StatusInternalServerError, "", err))
		return
	}

	var req CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[CheckoutResponse](http.StatusInternalServerError, "", err))
		return
	}

	cart, err := sv.repo.GetCart(c, repository.GetCartParams{
		UserID: authPayload.UserID,
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[CheckoutResponse](http.StatusNotFound, "", errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[CheckoutResponse](http.StatusInternalServerError, "", err))
		return
	}
	var addressID int64

	if req.AddressID == nil {
		defaultAddress, err := sv.repo.GetDefaultAddress(c, user.ID)
		if err != nil {
			if errors.Is(err, repository.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, createErrorResponse[CheckoutResponse](http.StatusNotFound, "", errors.New("address not found")))
				return
			}
			c.JSON(http.StatusInternalServerError, createErrorResponse[CheckoutResponse](http.StatusInternalServerError, "", err))
			return
		}
		addressID = defaultAddress.ID
	} else {
		address, err := sv.repo.GetAddress(c, repository.GetAddressParams{
			ID:     *req.AddressID,
			UserID: user.ID,
		})
		if err != nil {
			if errors.Is(err, repository.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, createErrorResponse[CheckoutResponse](http.StatusNotFound, "", errors.New("address not found")))
				return
			}
			c.JSON(http.StatusInternalServerError, createErrorResponse[CheckoutResponse](http.StatusInternalServerError, "", err))
			return
		}
		addressID = address.ID
	}

	if cart.UserID != authPayload.UserID {
		c.JSON(http.StatusForbidden, createErrorResponse[CheckoutResponse](http.StatusForbidden, "", errors.New("user not found")))
		return
	}

	cartItems, err := sv.repo.GetCartItems(c, cart.ID)
	if err != nil {
		log.Error().Err(err).Msg("GetCartItems")
		return
	}

	// create order
	createOrderItemParams := make([]repository.CreateOrderItemParams, 0)
	attributeList := make([]AttributeValue, 0)
	totalPrice := float64(0)
	for _, item := range cartItems {
		price, _ := item.Price.Float64Value()
		paramIdx := -1
		for j, param := range createOrderItemParams {
			if param.VariantID == item.CartItem.VariantID {
				paramIdx = j
				break
			}
		}
		if paramIdx != -1 {
			itemParam := repository.CreateOrderItemParams{
				VariantID:            item.CartItem.VariantID,
				Quantity:             item.CartItem.Quantity,
				PricePerUnitSnapshot: item.Price,
				ID:                   uuid.New(),
				OrderID:              uuid.New(),
				VariantSkuSnapshot:   item.Sku,
				ProductNameSnapshot:  item.ProductName,
				LineTotalSnapshot:    utils.GetPgNumericFromFloat(float64(item.Quantity) * price.Float64),
			}
			createOrderItemParams = append(createOrderItemParams, itemParam)
			attributeList = append(attributeList, AttributeValue{
				ID:           item.AttrID,
				Value:        item.AttrValText,
				DisplayValue: &item.AttrDisplayVal.String,
			})
		} else {
			attrIdx := -1
			for j, attr := range attributeList {
				if attr.ID == item.AttrID {
					attrIdx = j
					break
				}
			}
			if attrIdx == -1 {
				attributeList = append(attributeList, AttributeValue{
					ID:           item.AttrValID,
					Value:        item.AttrValText,
					DisplayValue: &item.AttrDisplayVal.String,
				})
			} else {

			}
		}
		totalPrice += price.Float64 * float64(item.CartItem.Quantity)
	}

	paymentMethod := repository.PaymentMethodCod
	if req.PaymentGateway != nil {
		paymentGateway := repository.PaymentGateway(*req.PaymentGateway)
		switch paymentGateway {
		case repository.PaymentGatewayStripe:
			paymentMethod = repository.PaymentMethodCard
			stripeInstance, err := payment.NewStripePayment(sv.config.StripeSecretKey)
			if err != nil {
				c.JSON(http.StatusInternalServerError, createErrorResponse[CheckoutResponse](http.StatusInternalServerError, "", err))
				return
			}
			sv.paymentCtx.SetStrategy(stripeInstance)
		default:
			c.JSON(http.StatusBadRequest, createErrorResponse[CheckoutResponse](http.StatusBadRequest, "", errors.New("payment gateway not supported")))
			return
		}
	}

	var paymentID string = uuid.New().String()
	if req.PaymentGateway != nil {
		paymentID, err = sv.paymentCtx.InitiatePayment(totalPrice, user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse[CheckoutResponse](http.StatusInternalServerError, "", err))
			return
		}
	}

	params := repository.CreateOrderTxParams{
		CartID:                cart.ID,
		PaymentMethod:         paymentMethod,
		AddressID:             addressID,
		TotalPrice:            totalPrice,
		UserID:                authPayload.UserID,
		CreateOrderItemParams: createOrderItemParams,
		PaymentID:             paymentID,
	}

	if req.PaymentGateway != nil {
		params.PaymentGateway = repository.PaymentGatewayStripe
	}

	checkoutResult, err := sv.repo.CreateOrderTx(c, params)

	if err != nil {
		sv.paymentCtx.CancelPayment(paymentID, "order creation failed")
		c.JSON(http.StatusInternalServerError, createErrorResponse[CheckoutResponse](http.StatusInternalServerError, "", err))
		return
	}
	resp := CheckoutResponse{
		OrderID:   checkoutResult,
		PaymentID: paymentID,
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, resp, "", nil, nil))
}

// @Summary Update product items in the cart
// @Schemes http
// @Description update product items in the cart
// @Tags carts
// @Accept json
// @Param input body updateCartItemRequest true "Update cart items input"
// @Produce json
// @Success 200 {object} ApiResponse[bool]
// @Failure 400 {object} ApiResponse[bool]
// @Failure 500 {object} ApiResponse[bool]
// @Router /cart/item/{id}/quantity [put]
func (sv *Server) updateCartItemQuantity(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, createErrorResponse[bool](http.StatusBadRequest, "", errors.New("user not found")))
		return
	}

	var param GetCartItemParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[bool](http.StatusInternalServerError, "", err))
		return
	}
	var req updateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[bool](http.StatusInternalServerError, "", err))
		return
	}

	cart, err := sv.repo.GetCart(c, repository.GetCartParams{
		UserID: authPayload.UserID,
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[bool](http.StatusNotFound, "", errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](http.StatusInternalServerError, "", err))
		return
	}

	if cart.UserID != authPayload.UserID {
		c.JSON(http.StatusForbidden, createErrorResponse[bool](http.StatusForbidden, "", errors.New("user not found")))
		return
	}

	err = sv.repo.UpdateCartItemQuantity(c, repository.UpdateCartItemQuantityParams{
		Quantity: req.Quantity,
		ID:       uuid.MustParse(param.ID),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](http.StatusInternalServerError, "", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, true, "cart item updated", nil, nil))
}

// @Summary  Clear the cart
// @Schemes http
// @Description  clear the cart
// @Tags carts
// @Accept json
// @Produce json
// @Success 200 {object} ApiResponse[bool]
// @Failure 400 {object} ApiResponse[bool]
// @Failure 500 {object} ApiResponse[bool]
// @Router /cart/clear [put]
func (sv *Server) clearCart(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, createErrorResponse[bool](http.StatusBadRequest, "", errors.New("user not found")))
		return
	}

	cart, err := sv.repo.GetCart(c, repository.GetCartParams{
		UserID: authPayload.UserID,
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[bool](http.StatusNotFound, "", errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](http.StatusInternalServerError, "", err))
		return
	}

	if cart.UserID != authPayload.UserID {
		c.JSON(http.StatusForbidden, createErrorResponse[bool](http.StatusForbidden, "", errors.New("user not found")))
		return
	}

	err = sv.repo.ClearCart(c, cart.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](http.StatusInternalServerError, "", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, true, "cart cleared", nil, nil))
}

// ------------------------------ Mappers ------------------------------
func mapToCartItemsResp(rows []repository.GetCartItemsRow) ([]CartItemResponse, float64) {
	cartItemsResp := make([]CartItemResponse, 0)

	var totalPrice float64
	for _, row := range rows {
		// if it's the first item or the previous item is different

		cartItemIdx := -1
		for j, cartItem := range cartItemsResp {
			if cartItem.ID == row.CartItem.ID.String() {
				cartItemIdx = j
				break
			}
		}
		if cartItemIdx == -1 {
			price, _ := row.Price.Float64Value()
			attr := CartItemAttribute{
				Name:  row.AttrName,
				Value: row.AttrValText,
			}
			if row.AttrDisplayVal.Valid && row.AttrDisplayVal.String != "" {
				attr.Value = row.AttrDisplayVal.String
			}
			cartItem := CartItemResponse{
				ID:        row.CartItemID.String(),
				ProductID: row.ProductID.String(),
				VariantID: row.VariantID.String(),
				Name:      row.ProductName,
				Quantity:  row.Quantity,
				Price:     math.Round(price.Float64*100) / 100,
				// Discount:      row.Discount.Int16,
				StockQty: row.StockQty,
				Sku:      &row.Sku,
				Attributes: []CartItemAttribute{
					attr,
				},
				ImageURL: &row.ImageUrl.String,
			}
			cartItemsResp = append(cartItemsResp, cartItem)
			totalPrice += cartItem.Price * float64(cartItem.Quantity)
		} else {
			attrIdx := -1
			for i, attr := range cartItemsResp[cartItemIdx].Attributes {
				if attr.Name == row.AttrName {
					attrIdx = i
					break
				}
			}
			if attrIdx == -1 {
				attr := CartItemAttribute{
					Name:  row.AttrName,
					Value: row.AttrValText,
				}
				if row.AttrDisplayVal.Valid && row.AttrDisplayVal.String != "" {
					attr.Value = row.AttrDisplayVal.String
				}
				cartItemsResp[cartItemIdx].Attributes = append(cartItemsResp[cartItemIdx].Attributes, attr)
			}
		}
	}

	return cartItemsResp, totalPrice
}
