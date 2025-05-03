package api

import (
	"errors"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v81"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
	"github.com/thanhphuocnguyen/go-eshop/pkg/payment"
)

type ProductVariantParam struct {
	ID string `uri:"variant_id" binding:"required,uuid"`
}

type CartItemParam struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type OrderItemAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CartItemResponse struct {
	ID         string                             `json:"id" binding:"required,uuid"`
	ProductID  string                             `json:"product_id" binding:"required,uuid"`
	VariantID  string                             `json:"variant_id" binding:"required,uuid"`
	Name       string                             `json:"name"`
	Quantity   int16                              `json:"quantity"`
	Price      float64                            `json:"price"`
	Discount   int16                              `json:"discount"`
	StockQty   int32                              `json:"stock"`
	Sku        *string                            `json:"sku,omitempty"`
	ImageURL   *string                            `json:"image_url,omitempty"`
	Attributes []repository.AttributeDataSnapshot `json:"attributes"`
}

type CartDetailResponse struct {
	ID         uuid.UUID          `json:"id"`
	TotalPrice float64            `json:"total_price"`
	CartItems  []CartItemResponse `json:"cart_items"`
	UpdatedAt  time.Time          `json:"updated_at,omitempty"`
	CreatedAt  time.Time          `json:"created_at"`
}

type UpdateCartItemQtyReq struct {
	Quantity int16 `json:"quantity" binding:"required,gt=0"`
}

type GetCartItemParam struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type CheckoutRequest struct {
	PaymentMethod      string   `json:"payment_method" binding:"required,oneof=code stripe"`
	PaymentGateway     *string  `json:"payment_gateway" binding:"omitempty,oneof=stripe"`
	AddressID          *int64   `json:"address_id" binding:"omitempty"`
	Email              *string  `json:"email" binding:"omitempty,email"`
	FullName           *string  `json:"full_name" binding:"omitempty"`
	Address            *Address `json:"address" binding:"omitempty"`
	PaymentRecipeEmail *string  `json:"payment_receipt_email" binding:"omitempty,email"`
}

type CheckoutResponse struct {
	OrderID         uuid.UUID `json:"order_id"`
	PaymentID       string    `json:"payment_id"`
	PaymentIntentID *string   `json:"payment_intent_id,omitempty"`
	ClientSecret    *string   `json:"client_secret,omitempty"`
}

// ------------------------------ Handlers ------------------------------

// @Summary Create a new cart
// @Schemes http
// @Description create a new cart for a user
// @Tags carts
// @Accept json
// @Produce json
// @Success 200 {object} ApiResponse[CartDetailResponse]
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 401 {object} gin.H
// @Router /cart [post]
func (sv *Server) createCart(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, createErrorResponse[CartDetailResponse](InvalidBodyCode, "", errors.New("user not found")))
		return
	}
	user, err := sv.repo.GetUserByID(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[CartDetailResponse](NotFoundCode, "", errors.New("user not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[CartDetailResponse](InternalServerErrorCode, "", err))
		return
	}
	_, err = sv.repo.GetCart(c, repository.GetCartParams{
		UserID: utils.GetPgTypeUUID(authPayload.UserID),
	})
	if err == nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[CartDetailResponse](InvalidBodyCode, "", errors.New("cart already exists")))
		return
	}

	newCart, err := sv.repo.CreateCart(c, repository.CreateCartParams{
		ID:     uuid.New(),
		UserID: utils.GetPgTypeUUID(user.ID),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[CartDetailResponse](InternalServerErrorCode, "", err))
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
// @Failure 500 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 401 {object} gin.H
// @Router /cart [get]
func (sv *Server) getCartHandler(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, createErrorResponse[CartDetailResponse](InvalidBodyCode, "", errors.New("user not found")))
		return
	}
	cart, err := sv.repo.GetCart(c, repository.GetCartParams{
		UserID: utils.GetPgTypeUUID(authPayload.UserID),
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			cart, err := sv.repo.CreateCart(c, repository.CreateCartParams{
				ID:     uuid.New(),
				UserID: utils.GetPgTypeUUID(authPayload.UserID),
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, createErrorResponse[CartDetailResponse](InternalServerErrorCode, "", err))
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
		c.JSON(http.StatusInternalServerError, createErrorResponse[CartDetailResponse](InternalServerErrorCode, "", err))
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
		c.JSON(http.StatusInternalServerError, createErrorResponse[CartDetailResponse](InternalServerErrorCode, "", err))
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
// @Param input body UpdateCartItemQtyReq true "Add product to cart input"
// @Produce json
// @Success 200 {object} ApiResponse[uuid.UUID]
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /cart/item/{variant_id} [post]
func (sv *Server) updateCartItemQtyHandler(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, createErrorResponse[uuid.UUID](InvalidBodyCode, "", errors.New("user not found")))
		return
	}

	var param CartItemParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[uuid.UUID](InvalidBodyCode, "", errors.New("invalid variant id")))
		return
	}

	var req UpdateCartItemQtyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[uuid.UUID](InternalServerErrorCode, "", err))
		return
	}

	cart, err := sv.repo.GetCart(c, repository.GetCartParams{
		UserID: utils.GetPgTypeUUID(authPayload.UserID),
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[uuid.UUID](NotFoundCode, "", errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[uuid.UUID](InternalServerErrorCode, "", err))
		return
	}

	cartItem, err := sv.repo.GetCartItem(c, repository.GetCartItemParams{
		ID:     uuid.MustParse(param.ID),
		CartID: cart.ID,
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			cartItem, err = sv.repo.CreateCartItem(c, repository.CreateCartItemParams{
				ID:        uuid.New(),
				CartID:    cart.ID,
				VariantID: uuid.MustParse(param.ID),
				Quantity:  req.Quantity,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, createErrorResponse[uuid.UUID](InternalServerErrorCode, "", err))
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, createErrorResponse[uuid.UUID](InternalServerErrorCode, "", err))
			return
		}
	} else {
		err = sv.repo.UpdateCartItemQuantity(c, repository.UpdateCartItemQuantityParams{
			Quantity: req.Quantity,
			ID:       cartItem.ID,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse[uuid.UUID](InternalServerErrorCode, "", err))
			return
		}
	}

	err = sv.repo.UpdateCartTimestamp(c, cart.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[uuid.UUID](InternalServerErrorCode, "", err))
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
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /cart/item/{id} [delete]
func (sv *Server) removeCartItem(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, createErrorResponse[string](InvalidBodyCode, "", errors.New("user not found")))
		return
	}

	var param GetCartItemParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[string](InternalServerErrorCode, "", err))
		return
	}

	cart, err := sv.repo.GetCart(c, repository.GetCartParams{
		UserID: utils.GetPgTypeUUID(authPayload.UserID),
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[string](NotFoundCode, "", errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[string](InternalServerErrorCode, "", err))
		return
	}

	if cart.UserID.Valid && string(cart.UserID.Bytes[:]) != authPayload.UserID.String() {
		c.JSON(http.StatusForbidden, createErrorResponse[string]("forbidden", "", errors.New("user not found")))
		return
	}

	err = sv.repo.RemoveProductFromCart(c, repository.RemoveProductFromCartParams{
		CartID: cart.ID,
		ID:     uuid.MustParse(param.ID),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[string](InternalServerErrorCode, "", err))
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
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /cart/checkoutHandler [post]
func (sv *Server) checkoutHandler(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, createErrorResponse[CheckoutResponse](InvalidBodyCode, "", errors.New("user not found")))
		return
	}

	var req CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[CheckoutResponse](InvalidBodyCode, "", err))
		return
	}

	user, err := sv.repo.GetUserByID(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[CheckoutResponse](NotFoundCode, "", errors.New("user not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[CheckoutResponse](InternalServerErrorCode, "", err))
		return
	}

	cart, err := sv.repo.GetCart(c, repository.GetCartParams{
		UserID: utils.GetPgTypeUUID(authPayload.UserID),
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[CheckoutResponse](InternalServerErrorCode, "", errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[CheckoutResponse](InternalServerErrorCode, "", err))
		return
	}
	var shippingAddr repository.ShippingAddressSnapshot
	if req.AddressID == nil {
		// create new address
		if req.Address == nil {
			c.JSON(http.StatusBadRequest, createErrorResponse[CheckoutResponse](InternalServerErrorCode, "", errors.New("address not found")))
			return
		}
		address, err := sv.repo.CreateAddress(c, repository.CreateAddressParams{
			UserID:   user.ID,
			Phone:    req.Address.Phone,
			Street:   req.Address.Street,
			Ward:     req.Address.Ward,
			District: req.Address.District,
			City:     req.Address.City,
			Default:  false,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse[CheckoutResponse](InternalServerErrorCode, "", err))
			return
		}
		shippingAddr = repository.ShippingAddressSnapshot{
			Street:   address.Street,
			Ward:     *address.Ward,
			District: address.District,
			City:     address.City,
			Phone:    address.Phone,
		}
	} else {
		address, err := sv.repo.GetAddress(c, repository.GetAddressParams{
			ID:     *req.AddressID,
			UserID: user.ID,
		})
		if err != nil {
			if errors.Is(err, repository.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, createErrorResponse[CheckoutResponse](NotFoundCode, "", errors.New("address not found")))
				return
			}
			c.JSON(http.StatusInternalServerError, createErrorResponse[CheckoutResponse](InternalServerErrorCode, "", err))
			return
		}
		shippingAddr = repository.ShippingAddressSnapshot{
			Street:   address.Street,
			Ward:     *address.Ward,
			District: address.District,
			City:     address.City,
			Phone:    address.Phone,
		}
	}
	if cart.UserID.Valid {
		cartUserId, err := uuid.FromBytes(cart.UserID.Bytes[:])
		if err != nil {
			log.Error().Err(err).Msg("GetCart")
			c.JSON(http.StatusInternalServerError, createErrorResponse[CheckoutResponse](InternalServerErrorCode, "", err))
			return
		}
		if cartUserId.String() != authPayload.UserID.String() {
			c.JSON(http.StatusForbidden, createErrorResponse[CheckoutResponse](PermissionDeniedCode, "", errors.New("you are not allowed to access this cart")))
			return
		}
	}

	cartItemRows, err := sv.repo.GetCartItemsForOrder(c, cart.ID)
	if err != nil {
		log.Error().Err(err).Msg("GetCartItems")
		return
	}

	// create order
	createOrderItemParams := make([]repository.CreateBulkOrderItemsParams, 0)
	attributeList := make(map[string][]OrderItemAttribute)
	var totalPrice float64

	for _, item := range cartItemRows {
		price, _ := item.Price.Float64Value()
		paramIdx := -1
		for j, param := range createOrderItemParams {
			if param.VariantID.String() == item.CartItem.VariantID.String() {
				paramIdx = j
				break
			}
		}
		if paramIdx == -1 {
			itemParam := repository.CreateBulkOrderItemsParams{
				ID:                   uuid.New(),
				VariantID:            item.CartItem.VariantID,
				Quantity:             item.CartItem.Quantity,
				PricePerUnitSnapshot: item.Price,
				VariantSkuSnapshot:   item.Sku,
				ProductNameSnapshot:  item.ProductName,
				LineTotalSnapshot:    utils.GetPgNumericFromFloat(float64(item.CartItem.Quantity) * price.Float64),
			}
			createOrderItemParams = append(createOrderItemParams, itemParam)
			attributeList[item.VariantID.String()] = []OrderItemAttribute{{
				Name:  item.AttrName,
				Value: item.AttrValCode,
			}}
			totalPrice += price.Float64 * float64(item.CartItem.Quantity)
		} else {
			attributeList[item.VariantID.String()] = append(attributeList[item.VariantID.String()], OrderItemAttribute{
				Name:  item.AttrName,
				Value: item.AttrValCode,
			})
		}
	}

	params := repository.CreateOrderTxArgs{
		CartID:          cart.ID,
		TotalPrice:      totalPrice,
		ShippingAddress: shippingAddr,
		UserID:          authPayload.UserID,
		CustomerInfo: repository.CustomerInfoTxArgs{
			FullName: user.Fullname,
			Email:    user.Email,
			Phone:    user.Phone,
		},
		CreateOrderItemParams: createOrderItemParams,
	}

	if req.FullName != nil {
		params.CustomerInfo.FullName = *req.FullName
	}

	if req.Email != nil {
		params.CustomerInfo.Email = *req.Email
	}

	orderID, err := sv.repo.CreateOrderTx(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[CheckoutResponse](InternalServerErrorCode, "", err))
		return
	}
	createPaymentArgs := repository.CreatePaymentParams{
		ID:      uuid.New(),
		OrderID: orderID,
		Amount:  utils.GetPgNumericFromFloat(totalPrice),
	}
	switch req.PaymentMethod {
	case string(repository.PaymentMethodStripe):
		createPaymentArgs.PaymentMethod = repository.PaymentMethodStripe
		createPaymentArgs.PaymentGateway = repository.NullPaymentGateway{
			PaymentGateway: repository.PaymentGatewayStripe,
			Valid:          true,
		}

		stripeInstance, err := payment.NewStripePayment(sv.config.StripeSecretKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse[CheckoutResponse](InternalServerErrorCode, "", err))
			return
		}
		sv.paymentCtx.SetStrategy(stripeInstance)
	default:
		c.JSON(http.StatusBadRequest, createErrorResponse[CheckoutResponse](InvalidBodyCode, "", errors.New("payment gateway not supported")))
		return
	}

	// init payment for custom payment

	receiptEmail := ""
	if req.PaymentRecipeEmail != nil {
		receiptEmail = *req.PaymentRecipeEmail
	} else {
		receiptEmail = user.Email
	}

	checkoutResp := CheckoutResponse{
		OrderID: orderID,
	}
	// create payment intent
	checkoutResult, checkoutErr := sv.paymentCtx.CreatePaymentIntent(totalPrice, receiptEmail)
	if checkoutErr != nil {
		c.JSON(http.StatusOK, createSuccessResponse(c, checkoutResp, "", nil, &ApiError{
			Code:    "payment_gateway_error",
			Details: checkoutErr.Error(),
			Stack:   checkoutErr.Error(),
		}))
		return
	}

	paymentIntent := checkoutResult.(*stripe.PaymentIntent)
	createPaymentArgs.GatewayPaymentIntentID = &paymentIntent.ID
	checkoutResp.ClientSecret = &paymentIntent.ClientSecret

	// create payment transaction
	payment, err := sv.repo.CreatePayment(c, createPaymentArgs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[CheckoutResponse](InternalServerErrorCode, "", err))
		return
	}

	checkoutResp.PaymentIntentID = &paymentIntent.ID
	checkoutResp.PaymentID = payment.ID.String()
	c.JSON(http.StatusOK, createSuccessResponse(c, checkoutResp, "", nil, nil))
}

// @Summary  Clear the cart
// @Schemes http
// @Description  clear the cart
// @Tags carts
// @Accept json
// @Produce json
// @Success 200 {object} ApiResponse[bool]
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /cart/clear [put]
func (sv *Server) clearCart(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, createErrorResponse[bool](InvalidBodyCode, "", errors.New("user not found")))
		return
	}

	cart, err := sv.repo.GetCart(c, repository.GetCartParams{
		UserID: utils.GetPgTypeUUID(authPayload.UserID),
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[bool](NotFoundCode, "", errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
		return
	}

	if string(cart.UserID.Bytes[:]) != authPayload.UserID.String() {
		c.JSON(http.StatusForbidden, createErrorResponse[bool]("forbidden", "", errors.New("user not found")))
		return
	}

	err = sv.repo.ClearCart(c, cart.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[bool](InternalServerErrorCode, "", err))
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
			attr := repository.AttributeDataSnapshot{
				Name:  row.AttrName,
				Value: row.AttrValName,
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
				Attributes: []repository.AttributeDataSnapshot{
					attr,
				},
				ImageURL: row.ImageUrl,
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
				attr := repository.AttributeDataSnapshot{
					Name:  row.AttrName,
					Value: row.AttrValName,
				}

				cartItemsResp[cartItemIdx].Attributes = append(cartItemsResp[cartItemIdx].Attributes, attr)
			}
		}
	}

	return cartItemsResp, totalPrice
}
