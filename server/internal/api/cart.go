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
	"github.com/thanhphuocnguyen/go-eshop/pkg/payment"
)

type updateCartItemRequest struct {
	Quantity int16 `json:"quantity" binding:"required,gt=0"`
}
type cartItemAttributeModel struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type cartItemResponse struct {
	ID            string                   `json:"id" binding:"required,uuid"`
	ProductID     string                   `json:"product_id" binding:"required,uuid"`
	VariantID     string                   `json:"variant_id" binding:"required,uuid"`
	Name          string                   `json:"name"`
	Quantity      int16                    `json:"quantity"`
	Price         float64                  `json:"price"`
	Discount      int16                    `json:"discount"`
	StockQuantity int32                    `json:"stock"`
	Sku           *string                  `json:"sku,omitempty"`
	ImageURL      *string                  `json:"image_url,omitempty"`
	Attributes    []cartItemAttributeModel `json:"attributes,omitempty"`
}

type cartResponse struct {
	ID         uuid.UUID          `json:"id"`
	UserID     uuid.UUID          `json:"user_id"`
	TotalPrice float64            `json:"total_price"`
	CartItems  []cartItemResponse `json:"cart_items,omitempty"`
	UpdatedAt  time.Time          `json:"updated_at,omitempty"`
	CreatedAt  time.Time          `json:"created_at"`
}

type addProductToCartRequest struct {
	VariantID string `json:"variant_id" binding:"required,uuid"`
	Quantity  int16  `json:"quantity" binding:"required,gt=0"`
}

type getCartItemParam struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type checkoutRequest struct {
	PaymentGateway *string `json:"payment_gateway" binding:"required,oneof=stripe paypal visa mastercard apple_pay google_pay postpaid momo zalo_pay vn_pay"`
	AddressID      *int64  `json:"address_id" binding:"omitempty,required"`
}

type checkoutResponse struct {
	OrderID   uuid.UUID `json:"order_id"`
	PaymentID string    `json:"payment_id"`
}

// ------------------------------ Mappers ------------------------------
func mapToCartResponse(cart repository.Cart, dataRows []repository.GetCartItemsByIDRow) cartResponse {
	var totalPrice float64
	cartItems := make([]cartItemResponse, 0)
	for i, row := range dataRows {
		// if it's the first item or the previous item is different
		lastIdx := len(cartItems) - 1

		if i == 0 {
			priceParsed, _ := row.Product.BasePrice.Float64Value()
			totalPrice += priceParsed.Float64 * float64(row.CartItem.Quantity)
			productVariant := cartItemResponse{
				ID:        row.CartItem.CartID.String(),
				Name:      row.Product.Name,
				Quantity:  row.CartItem.Quantity,
				Price:     priceParsed.Float64,
				ProductID: row.Product.ID.String(),
				Attributes: []cartItemAttributeModel{
					{
						Name:  row.AttributeName.String,
						Value: row.AttributeValue.String,
					},
				},
			}
			productVariant.Sku = &row.Product.BaseSku.String
			cartItems = append(cartItems, productVariant)
		} else {
			cartItems[lastIdx].Attributes = append(cartItems[lastIdx].Attributes, cartItemAttributeModel{
				Name:  row.AttributeName.String,
				Value: row.AttributeValue.String,
			})
		}
	}

	return cartResponse{
		ID:         cart.ID,
		UserID:     cart.UserID,
		UpdatedAt:  cart.UpdatedAt,
		CreatedAt:  cart.CreatedAt,
		CartItems:  cartItems,
		TotalPrice: math.Round(totalPrice*100) / 100,
	}
}

// ------------------------------ Handlers ------------------------------

// @Summary Create a new cart
// @Schemes http
// @Description create a new cart for a user
// @Tags carts
// @Accept json
// @Produce json
// @Success 200 {object} ApiResponse
// @Failure 400 {object} ApiResponse
// @Failure 500 {object} ApiResponse
// @Router /cart [post]
func (sv *Server) createCart(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", errors.New("user not found")))
		return
	}
	user, err := sv.repo.GetUserByID(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse(http.StatusNotFound, "", errors.New("user not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}
	_, err = sv.repo.GetCart(c, authPayload.UserID)
	if err == nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", errors.New("cart already exists")))
		return
	}

	newCart, err := sv.repo.CreateCart(c, repository.CreateCartParams{
		ID:     uuid.New(),
		UserID: user.ID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, newCart, "", nil, nil))
}

// @Summary Count the number of items in the cart
// @Schemes http
// @Description count the number of items in the cart
// @Tags carts
// @Accept json
// @Produce json
// @Success 200 {object} ApiResponse
// @Failure 400 {object} ApiResponse
// @Failure 500 {object} ApiResponse
// @Router /cart/items-count [get]
func (sv *Server) countCartItems(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", errors.New("user not found")))
		return
	}

	cart, err := sv.repo.GetCart(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse(http.StatusNotFound, "", errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}

	count, err := sv.repo.CountCartItems(c, cart.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, count, "", nil, nil))
}

// @Summary Get cart details by user ID
// @Schemes http
// @Description get cart details by user ID
// @Tags cart
// @Accept json
// @Produce json
// @Success 200 {object} ApiResponse
// @Failure 500 {object} ApiResponse
// @Router /cart [get]
func (sv *Server) getCartDetail(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", errors.New("user not found")))
		return
	}

	cart, err := sv.repo.GetCart(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse(http.StatusNotFound, "", errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}

	cartItems, err := sv.repo.GetCartItemsByID(c, cart.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}

	cartDetail := mapToCartResponse(cart, cartItems)

	c.JSON(http.StatusOK, createSuccessResponse(c, cartDetail, "", nil, nil))
}

// @Summary Add a product to the cart
// @Schemes http
// @Description add a product to the cart
// @Tags carts
// @Accept json
// @Param input body addProductToCartRequest true "Add product to cart input"
// @Produce json
// @Success 200 {object} ApiResponse
// @Failure 400 {object} ApiResponse
// @Failure 500 {object} ApiResponse
// @Router /cart/item [post]
func (sv *Server) addCartItem(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", errors.New("user not found")))
		return
	}

	var req addProductToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}

	cart, err := sv.repo.GetCart(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			// create a new cart if not found
			cart, err = sv.repo.CreateCart(c, repository.CreateCartParams{
				ID:     uuid.New(),
				UserID: authPayload.UserID,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
			return
		}
	}

	if cart.UserID != authPayload.UserID {
		c.JSON(http.StatusForbidden, createErrorResponse(http.StatusForbidden, "", errors.New("user not found")))
		return
	}

	// product, err := sv.repo.GetProductByID(c, repository.GetProductByIDParams{
	// 	ID: uuid.MustParse(req.VariantID),
	// })

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse(http.StatusNotFound, "", errors.New("product not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}

	cartItem, err := sv.repo.GetCartItemByProductID(c, repository.GetCartItemByProductIDParams{
		ID:        cart.ID,
		VariantID: uuid.MustParse(req.VariantID),
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			cartItem, err = sv.repo.CreateCartItem(c, repository.CreateCartItemParams{
				ID:        cart.ID,
				VariantID: uuid.MustParse(req.VariantID),
				Quantity:  req.Quantity,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
			return
		}
	} else {
		if cartItem.VariantID.String() != req.VariantID {
			c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusNotFound, "", errors.New("product not found")))
			return
		}

		err = sv.repo.UpdateCartItemQuantity(c, repository.UpdateCartItemQuantityParams{
			Quantity: cartItem.Quantity + req.Quantity,
			ID:       cartItem.ID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
			return
		}
	}

	err = sv.repo.UpdateCart(c, cart.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}

	createdID := cartItem.ID.String()

	c.JSON(http.StatusOK, createSuccessResponse(c, createdID, "", nil, nil))
}

// @Summary Remove a product from the cart
// @Schemes http
// @Description remove a product from the cart
// @Tags carts
// @Accept json
// @Param id path int true "Product ID"
// @Produce json
// @Success 200 {object} ApiResponse
// @Failure 400 {object} ApiResponse
// @Failure 500 {object} ApiResponse
// @Router /cart/item/{id} [delete]
func (sv *Server) removeCartItem(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", errors.New("user not found")))
		return
	}

	var param getCartItemParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}

	cart, err := sv.repo.GetCart(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse(http.StatusNotFound, "", errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}

	if cart.UserID != authPayload.UserID {
		c.JSON(http.StatusForbidden, createErrorResponse(http.StatusForbidden, "", errors.New("user not found")))
		return
	}

	err = sv.repo.RemoveProductFromCart(c, repository.RemoveProductFromCartParams{
		CartID: cart.ID,
		ID:     uuid.MustParse(param.ID),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
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
// @Param input body checkoutRequest true "Update cart items input"
// @Produce json
// @Success 200 {object} ApiResponse
// @Failure 400 {object} ApiResponse
// @Failure 500 {object} ApiResponse
// @Router /cart/checkout [post]
func (sv *Server) checkout(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", errors.New("user not found")))
		return
	}

	user, err := sv.repo.GetUserByID(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse(http.StatusBadRequest, "", errors.New("user not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}

	var req checkoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}

	cart, err := sv.repo.GetCart(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse(http.StatusNotFound, "", errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}
	var addressID int64

	if req.AddressID == nil {
		defaultAddress, err := sv.repo.GetDefaultAddress(c, user.ID)
		if err != nil {
			if errors.Is(err, repository.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, createErrorResponse(http.StatusNotFound, "", errors.New("address not found")))
				return
			}
			c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
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
				c.JSON(http.StatusNotFound, createErrorResponse(http.StatusNotFound, "", errors.New("address not found")))
				return
			}
			c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
			return
		}
		addressID = address.ID
	}

	if cart.UserID != authPayload.UserID {
		c.JSON(http.StatusForbidden, createErrorResponse(http.StatusForbidden, "", errors.New("user not found")))
		return
	}

	cartItems, err := sv.repo.GetCartItemsByID(c, cart.ID)
	if err != nil {
		log.Error().Err(err).Msg("GetCartItems")
		return
	}
	if len(cartItems) == 0 {
		log.Error().Msg("Cart is empty")
	}
	// create order
	createOrderItemParams := make([]repository.CreateOrderItemParams, len(cartItems))
	totalPrice := float64(0)
	for i, item := range cartItems {
		price, _ := item.Product.BasePrice.Float64Value()
		createOrderItemParams[i] = repository.CreateOrderItemParams{
			VariantID:            item.CartItem.VariantID,
			Quantity:             item.CartItem.Quantity,
			PricePerUnitSnapshot: item.Product.BasePrice,
			ID:                   uuid.New(),
			OrderID:              uuid.New(),
			VariantSkuSnapshot:   item.Product.BaseSku.String,
			ProductNameSnapshot:  item.Product.Name,
			LineTotalSnapshot:    item.Product.BasePrice,
			AttributesSnapshot:   []byte(item.AttributeName.String),
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
				c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
				return
			}
			sv.paymentCtx.SetStrategy(stripeInstance)
		default:
			c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", errors.New("payment gateway not supported")))
			return
		}
	}

	var paymentID string = uuid.New().String()
	if req.PaymentGateway != nil {
		paymentID, err = sv.paymentCtx.InitiatePayment(totalPrice, user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
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
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}
	resp := checkoutResponse{
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
// @Success 200 {object} ApiResponse
// @Failure 400 {object} ApiResponse
// @Failure 500 {object} ApiResponse
// @Router /cart/item/{id}/quantity [put]
func (sv *Server) updateCartItemQuantity(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", errors.New("user not found")))
		return
	}

	var param getCartItemParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}
	var req updateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}

	cart, err := sv.repo.GetCart(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse(http.StatusNotFound, "", errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}

	if cart.UserID != authPayload.UserID {
		c.JSON(http.StatusForbidden, createErrorResponse(http.StatusForbidden, "", errors.New("user not found")))
		return
	}

	err = sv.repo.UpdateCartItemQuantity(c, repository.UpdateCartItemQuantityParams{
		Quantity: req.Quantity,
		ID:       uuid.MustParse(param.ID),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}

	msg := "cart item updated"
	c.JSON(http.StatusOK, createSuccessResponse(c, msg, "", nil, nil))
}

// @Summary  Clear the cart
// @Schemes http
// @Description  clear the cart
// @Tags carts
// @Accept json
// @Produce json
// @Success 200 {object} ApiResponse
// @Failure 400 {object} ApiResponse
// @Failure 500 {object} ApiResponse
// @Router /cart/clear [put]
func (sv *Server) clearCart(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, createErrorResponse(http.StatusBadRequest, "", errors.New("user not found")))
		return
	}

	cart, err := sv.repo.GetCart(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse(http.StatusNotFound, "", errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}

	if cart.UserID != authPayload.UserID {
		c.JSON(http.StatusForbidden, createErrorResponse(http.StatusForbidden, "", errors.New("user not found")))
		return
	}

	err = sv.repo.ClearCart(c, cart.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse(http.StatusInternalServerError, "", err))
		return
	}

	msg := "cart cleared"
	c.JSON(http.StatusOK, createSuccessResponse(c, msg, "", nil, nil))
}
