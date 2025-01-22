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
	ID            int32                    `json:"id"`
	ProductID     int64                    `json:"product_id"`
	VariantID     int64                    `json:"variant_id"`
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
	ProductID int64 `json:"product_id" binding:"required,min=1"`
	VariantID int64 `json:"variant_id" binding:"required,min=1"`
	Quantity  int16 `json:"quantity" binding:"required,gt=0"`
}

type getCartItemParam struct {
	ID int32 `uri:"id" binding:"required,gt=0"`
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

		if i == 0 || cartItems[lastIdx].VariantID != row.VariantID {
			priceParsed, _ := row.Price.Float64Value()
			totalPrice += priceParsed.Float64 * float64(row.Quantity)
			productVariant := cartItemResponse{
				ID:            row.CartItemID,
				Name:          row.ProductName,
				Quantity:      row.Quantity,
				Price:         priceParsed.Float64,
				Discount:      row.Discount,
				StockQuantity: row.StockQuantity,
				VariantID:     row.VariantID,
				ProductID:     row.ProductID,
				Attributes: []cartItemAttributeModel{
					{
						Name:  row.AttributeName,
						Value: row.AttributeValue,
					},
				},
			}
			if row.Sku.Valid {
				productVariant.Sku = &row.Sku.String
			}
			if row.ImageUrl.Valid {
				productVariant.ImageURL = &row.ImageUrl.String
			}
			cartItems = append(cartItems, productVariant)
		} else {
			cartItems[lastIdx].Attributes = append(cartItems[lastIdx].Attributes, cartItemAttributeModel{
				Name:  row.AttributeName,
				Value: row.AttributeValue,
			})
		}
	}

	return cartResponse{
		ID:         cart.CartID,
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
// @Success 200 {object} GenericResponse[repository.Cart]
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /cart [post]
func (sv *Server) createCart(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("user not found")))
		return
	}
	user, err := sv.repo.GetUserByID(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(errors.New("user not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	_, err = sv.repo.GetCart(c, authPayload.UserID)
	if err == nil {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("cart already existed")))
		return
	}

	newCart, err := sv.repo.CreateCart(c, repository.CreateCartParams{
		CartID: uuid.New(),
		UserID: user.UserID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusOK, GenericResponse[repository.Cart]{&newCart, nil, nil})
}

// @Summary Get cart details by user ID
// @Schemes http
// @Description get cart details by user ID
// @Tags cart
// @Accept json
// @Produce json
// @Success 200 {object} GenericResponse[cartResponse]
// @Failure 500 {object} errorResponse
// @Router /cart [get]
func (sv *Server) getCartDetail(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("user not found")))
		return
	}
	cart, err := sv.repo.GetCart(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	cartItems, err := sv.repo.GetCartItemsByID(c, cart.CartID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	cartDetail := mapToCartResponse(cart, cartItems)

	c.JSON(http.StatusOK, GenericResponse[cartResponse]{&cartDetail, nil, nil})
}

// @Summary Add a product to the cart
// @Schemes http
// @Description add a product to the cart
// @Tags carts
// @Accept json
// @Param input body addProductToCartRequest true "Add product to cart input"
// @Produce json
// @Success 200 {object} GenericResponse[cartItemResponse]
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /cart/item [post]
func (sv *Server) addCartItem(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("user not found")))
		return
	}

	var req addProductToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	cart, err := sv.repo.GetCart(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			// create a new cart if not found
			cart, err = sv.repo.CreateCart(c, repository.CreateCartParams{
				CartID: uuid.New(),
				UserID: authPayload.UserID,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, mapErrResp(err))
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, mapErrResp(err))
			return
		}
	}

	if cart.UserID != authPayload.UserID {
		c.JSON(http.StatusForbidden, mapErrResp(errors.New("cart does not belong to the user")))
		return
	}

	productVariant, err := sv.repo.GetVariantByID(c, req.VariantID)

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(errors.New("product not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if productVariant.StockQuantity < int32(req.Quantity) {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("insufficient stock")))
		return
	}

	cartItem, err := sv.repo.GetCartItemByVariantID(c, req.VariantID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			cartItem, err = sv.repo.CreateCartItem(c, repository.CreateCartItemParams{
				CartID:    cart.CartID,
				ProductID: req.ProductID,
				VariantID: req.VariantID,
				Quantity:  req.Quantity,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, mapErrResp(err))
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, mapErrResp(err))
			return
		}
	} else {
		if cartItem.ProductID != req.ProductID {
			c.JSON(http.StatusBadRequest, mapErrResp(errors.New("product not found")))
			return
		}
		if int32(cartItem.Quantity+req.Quantity) > productVariant.StockQuantity {
			c.JSON(http.StatusBadRequest, mapErrResp(errors.New("insufficient stock")))
			return
		}

		err = sv.repo.UpdateCartItemQuantity(c, repository.UpdateCartItemQuantityParams{
			Quantity:   cartItem.Quantity + req.Quantity,
			CartItemID: cartItem.CartItemID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, mapErrResp(err))
			return
		}
	}

	err = sv.repo.UpdateCart(c, cart.CartID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	createdID := cartItem.CartItemID

	c.JSON(http.StatusOK, GenericResponse[int32]{&createdID, nil, nil})
}

// @Summary Remove a product from the cart
// @Schemes http
// @Description remove a product from the cart
// @Tags carts
// @Accept json
// @Param id path int true "Product ID"
// @Produce json
// @Success 200 {object} GenericResponse[bool]
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /cart/item/{id} [delete]
func (sv *Server) removeCartItem(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("user not found")))
		return
	}

	var param getCartItemParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	cart, err := sv.repo.GetCart(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if cart.UserID != authPayload.UserID {
		c.JSON(http.StatusForbidden, mapErrResp(errors.New("cart does not belong to the user")))
		return
	}

	err = sv.repo.RemoveProductFromCart(c, repository.RemoveProductFromCartParams{
		CartID:     cart.CartID,
		CartItemID: param.ID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	success := true
	message := "product removed"
	c.JSON(http.StatusOK, GenericResponse[bool]{&success, &message, nil})
}

// @Summary Update product items in the cart
// @Schemes http
// @Description update product items in the cart
// @Tags carts
// @Accept json
// @Param input body checkoutRequest true "Update cart items input"
// @Produce json
// @Success 200 {object} GenericResponse[checkoutResponse]
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /cart/checkout [post]
func (sv *Server) checkout(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("user not found")))
		return
	}

	user, err := sv.repo.GetUserByID(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(errors.New("user not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	var req checkoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	cart, err := sv.repo.GetCart(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	var addressID int64

	if req.AddressID == nil {
		defaultAddress, err := sv.repo.GetDefaultAddress(c, user.UserID)
		if err != nil {
			if errors.Is(err, repository.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, mapErrResp(errors.New("default address not found")))
				return
			}
			c.JSON(http.StatusInternalServerError, mapErrResp(err))
			return
		}
		addressID = defaultAddress.UserAddressID
	} else {
		address, err := sv.repo.GetAddress(c, repository.GetAddressParams{
			UserAddressID: *req.AddressID,
			UserID:        user.UserID,
		})
		if err != nil {
			if errors.Is(err, repository.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, mapErrResp(errors.New("address not found")))
				return
			}
			c.JSON(http.StatusInternalServerError, mapErrResp(err))
			return
		}
		addressID = address.UserAddressID
	}

	if cart.UserID != authPayload.UserID {
		c.JSON(http.StatusForbidden, mapErrResp(errors.New("cart does not belong to the user")))
		return
	}

	cartItems, err := sv.repo.GetCartItemsByID(c, cart.CartID)
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
		price, _ := item.Price.Float64Value()
		createOrderItemParams[i] = repository.CreateOrderItemParams{
			ProductID: item.ProductID,
			VariantID: item.VariantID,
			Quantity:  int32(item.Quantity),
			Price:     item.Price,
		}
		totalPrice += price.Float64 * float64(item.Quantity)
	}

	paymentMethod := repository.PaymentMethodCod
	if req.PaymentGateway != nil {
		paymentGateway := repository.PaymentGateway(*req.PaymentGateway)
		switch paymentGateway {
		case repository.PaymentGatewayStripe:
			paymentMethod = repository.PaymentMethodCard
			stripeInstance, err := payment.NewStripePayment(sv.config.StripeSecretKey)
			if err != nil {
				c.JSON(http.StatusInternalServerError, mapErrResp(err))
				return
			}
			sv.paymentCtx.SetStrategy(stripeInstance)
		default:
			c.JSON(http.StatusBadRequest, mapErrResp(errors.New("currently we only support stripe")))
			return
		}
	}

	var paymentID string = uuid.New().String()
	if req.PaymentGateway != nil {
		paymentID, err = sv.paymentCtx.InitiatePayment(totalPrice, user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, mapErrResp(err))
			return
		}
	}
	params := repository.CreateOrderTxParams{
		CartID:                cart.CartID,
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
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusOK, GenericResponse[checkoutResponse]{&checkoutResponse{
		OrderID:   checkoutResult,
		PaymentID: paymentID,
	}, nil, nil})
}

// @Summary Update product items in the cart
// @Schemes http
// @Description update product items in the cart
// @Tags carts
// @Accept json
// @Param input body updateCartItemRequest true "Update cart items input"
// @Produce json
// @Success 200 {object} GenericResponse[cartItemResponse]
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /cart/item/{id}/quantity [put]
func (sv *Server) updateCartItemQuantity(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("user not found")))
		return
	}

	var param getCartItemParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}
	var req updateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	cart, err := sv.repo.GetCart(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if cart.UserID != authPayload.UserID {
		c.JSON(http.StatusForbidden, mapErrResp(errors.New("cart does not belong to the user")))
		return
	}

	err = sv.repo.UpdateCartItemQuantity(c, repository.UpdateCartItemQuantityParams{
		Quantity:   req.Quantity,
		CartItemID: param.ID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	cartItem, err := sv.repo.GetCartItemWithProduct(c, param.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(errors.New("cart item not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	_, err = sv.repo.GetVariantByID(c, cartItem.VariantID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(errors.New("product not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	msg := "cart item updated"
	c.JSON(http.StatusOK, GenericResponse[cartItemResponse]{nil, &msg, nil})
}

// @Summary  Clear the cart
// @Schemes http
// @Description  clear the cart
// @Tags carts
// @Accept json
// @Produce json
// @Success 200 {object} GenericResponse[bool]
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /cart/clear [put]
func (sv *Server) clearCart(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("user not found")))
		return
	}

	cart, err := sv.repo.GetCart(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if cart.UserID != authPayload.UserID {
		c.JSON(http.StatusForbidden, mapErrResp(errors.New("cart does not belong to the user")))
		return
	}

	err = sv.repo.ClearCart(c, cart.CartID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	msg := "cart cleared"
	c.JSON(http.StatusOK, GenericResponse[string]{nil, &msg, nil})
}
