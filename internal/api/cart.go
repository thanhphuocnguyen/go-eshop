package api

import (
	"errors"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v81"
	"github.com/thanhphuocnguyen/go-eshop/internal/auth"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/pkg/payment"
)

type updateCartItemRequest struct {
	Quantity int16 `json:"quantity" binding:"required,gt=0"`
}

type cartItemResponse struct {
	ID        int32   `json:"id"`
	ProductID int64   `json:"product_id"`
	Name      string  `json:"name"`
	ImageURL  string  `json:"image_url"`
	Quantity  int16   `json:"quantity"`
	Price     float64 `json:"price"`
}

type cartResponse struct {
	ID         int32              `json:"id"`
	UserID     int64              `json:"user_id"`
	UpdatedAt  time.Time          `json:"updated_at,omitempty"`
	CreatedAt  time.Time          `json:"created_at"`
	CartItems  []cartItemResponse `json:"cart_items,omitempty"`
	TotalPrice float64            `json:"total_price"`
}

type addProductToCartRequest struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int16 `json:"quantity" binding:"required"`
}

type getCartItemParam struct {
	ID int32 `uri:"id" binding:"required,gt=0"`
}

type checkoutRequest struct {
	PaymentGateway *string `json:"payment_gateway" binding:"required,oneof=stripe paypal visa mastercard apple_pay google_pay postpaid momo zalo_pay vn_pay"`
	AddressID      *int64  `json:"address_id" binding:"omitempty,required"`
}

type checkoutResponse struct {
	OrderID   int64  `json:"order_id"`
	PaymentID string `json:"payment_id"`
}

// ------------------------------ Mappers ------------------------------

func mapToCartResponse(cart repository.Cart, cartItems []repository.GetCartItemsRow) cartResponse {
	var totalPrice float64
	products := make([]cartItemResponse, len(cartItems))
	for i, item := range cartItems {
		price, _ := item.ProductPrice.Float64Value()
		totalPrice += price.Float64 * float64(item.Quantity)
		products[i] = cartItemResponse{
			ID:        item.CartItemID,
			ProductID: item.ProductID,
			Name:      item.ProductName,
			Quantity:  item.Quantity,
			Price:     price.Float64,
		}

		if item.ImageUrl.Valid {
			products[i].ImageURL = item.ImageUrl.String
		}
	}

	return cartResponse{
		ID:         cart.CartID,
		UserID:     cart.UserID,
		UpdatedAt:  cart.UpdatedAt,
		CreatedAt:  cart.CreatedAt,
		CartItems:  products,
		TotalPrice: math.Round(totalPrice*100) / 100,
	}
}

// ------------------------------ Handlers ------------------------------

// createCart godoc
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
		if errors.Is(err, repository.ErrorRecordNotFound) {
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

	newCart, err := sv.repo.CreateCart(c, user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusOK, GenericResponse[repository.Cart]{&newCart, nil, nil})
}

// getCartDetail godoc
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
		if errors.Is(err, repository.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	cartItems, err := sv.repo.GetCartItems(c, cart.CartID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	cartDetail := mapToCartResponse(cart, cartItems)

	c.JSON(http.StatusOK, GenericResponse[cartResponse]{&cartDetail, nil, nil})
}

// addCartItem godoc
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

	product, err := sv.repo.GetProductWithImage(c, repository.GetProductWithImageParams{
		ProductID: req.ProductID,
		Archived: pgtype.Bool{
			Bool:  false,
			Valid: true,
		},
	})

	if err != nil {
		if errors.Is(err, repository.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(errors.New("product not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if product.Archived {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("product is archived")))
		return
	}

	if product.Stock < int32(req.Quantity) {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("insufficient stock")))
		return
	}

	cart, err := sv.repo.GetCart(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrorRecordNotFound) {
			cart, err = sv.repo.CreateCart(c, authPayload.UserID)
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

	// check if the product is already in the cart
	cartItem, err := sv.repo.GetCartItemByProductID(c, req.ProductID)
	if err != nil && errors.Is(err, repository.ErrorRecordNotFound) {
		cartItem, err = sv.repo.AddProductToCart(c, repository.AddProductToCartParams{
			ProductID: req.ProductID,
			CartID:    cart.CartID,
			Quantity:  req.Quantity,
		})
	} else if err == nil {
		if int32(cartItem.Quantity+req.Quantity) > product.Stock {
			c.JSON(http.StatusBadRequest, mapErrResp(errors.New("insufficient stock")))
			return
		}
		err = sv.repo.UpdateCartItemQuantity(c, repository.UpdateCartItemQuantityParams{
			Quantity:   cartItem.Quantity + req.Quantity,
			CartItemID: cartItem.CartItemID,
		})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	err = sv.repo.UpdateCart(c, cart.CartID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	productPrice, _ := product.Price.Float64Value()
	itemResp := cartItemResponse{
		ID:        cartItem.CartItemID,
		ProductID: cartItem.ProductID,
		Name:      product.Name,
		ImageURL:  product.ImageUrl.String,
		Quantity:  cartItem.Quantity,
		Price:     productPrice.Float64,
	}

	c.JSON(http.StatusOK, GenericResponse[cartItemResponse]{&itemResp, nil, nil})
}

// removeCartItem godoc
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
		if errors.Is(err, repository.ErrorRecordNotFound) {
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

// checkout godoc
// @Summary Update product items in the cart
// @Schemes http
// @Description update product items in the cart
// @Tags carts
// @Accept json
// @Param input body checkoutRequest true "Update cart items input"
// @Produce json
// @Success 200 {object} GenericResponse[repository.CheckoutCartTxResult]
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
		if errors.Is(err, repository.ErrorRecordNotFound) {
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
		if errors.Is(err, repository.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	addresses, err := sv.repo.GetAddresses(c, authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if len(addresses) == 0 {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("no address found")))
		return
	}

	addressID := int64(addresses[0].UserAddressID)
	if req.AddressID == nil {
		isAddressExist := false
		for _, address := range addresses {
			if address.Default {
				addressID = address.UserAddressID
			}
			if address.UserAddressID == *req.AddressID {
				isAddressExist = true
				break
			}
		}
		if !isAddressExist {
			c.JSON(http.StatusBadRequest, mapErrResp(errors.New("address not found")))
			return
		}
	} else {
		addressID = *req.AddressID
	}

	if cart.UserID != authPayload.UserID {
		c.JSON(http.StatusForbidden, mapErrResp(errors.New("cart does not belong to the user")))
		return
	}

	cartItems, err := sv.repo.GetCartItems(c, cart.CartID)
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
		if item.ProductStock < int32(item.Quantity) {
			log.Error().Msg("Product out of stock")
			return
		}
		price, _ := item.ProductPrice.Float64Value()
		createOrderItemParams[i] = repository.CreateOrderItemParams{
			ProductID: item.ProductID,
			Quantity:  int32(item.Quantity),
			Price:     item.ProductPrice,
		}
		totalPrice += price.Float64 * float64(item.Quantity)
	}

	log.Info().Msgf("Total price: %f", totalPrice)

	paymentMethod := repository.PaymentMethodCod
	if req.PaymentGateway != nil {
		paymentGateway := repository.PaymentGateway(*req.PaymentGateway)
		switch paymentGateway {

		case repository.PaymentGatewayStripe:
			stripeInstance, err := payment.NewStripePayment(sv.config.StripeSecretKey)
			if err != nil {
				c.JSON(http.StatusInternalServerError, mapErrResp(err))
				return
			}
			stripe.Key = sv.config.StripeSecretKey
			paymentMethod = repository.PaymentMethodCard
			stripePayment := stripeInstance
			sv.paymentCtx.SetStrategy(stripePayment)
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

// updateCartItemQuantity godoc
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
		if errors.Is(err, repository.ErrorRecordNotFound) {
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
		if errors.Is(err, repository.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(errors.New("cart item not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	product, err := sv.repo.GetProduct(c, repository.GetProductParams{
		ProductID: cartItem.ProductID,
	})

	if err != nil {
		if errors.Is(err, repository.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(errors.New("product not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
	}

	if product.Stock < int32(req.Quantity) {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("insufficient stock")))
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	productPrice, _ := product.Price.Float64Value()
	itemResp := cartItemResponse{
		ID:        cartItem.CartID,
		ProductID: cartItem.ProductID,
		Name:      product.Name,
		ImageURL:  cartItem.ImageUrl.String,
		Quantity:  req.Quantity,
		Price:     productPrice.Float64,
	}

	c.JSON(http.StatusOK, GenericResponse[cartItemResponse]{&itemResp, nil, nil})
}

// clearCart godoc
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
		if errors.Is(err, repository.ErrorRecordNotFound) {
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
	success := true
	c.JSON(http.StatusOK, GenericResponse[bool]{&success, &msg, nil})
}
