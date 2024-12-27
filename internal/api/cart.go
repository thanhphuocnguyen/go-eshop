package api

import (
	"errors"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/thanhphuocnguyen/go-eshop/internal/auth"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/postgres"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/sqlc"
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
	ID           int32              `json:"id"`
	CheckedOutAt time.Time          `json:"checked_out"`
	UserID       int64              `json:"user_id"`
	UpdatedAt    time.Time          `json:"updated_at"`
	CreatedAt    time.Time          `json:"created_at"`
	CartItems    []cartItemResponse `json:"cart_items"`
	TotalPrice   float64            `json:"total_price"`
}

type addProductToCartRequest struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int16 `json:"quantity" binding:"required"`
}

type getCartItemParam struct {
	ID int32 `uri:"id" binding:"required,gt=0"`
}

type checkoutRequest struct {
	PaymentMethod string `json:"payment_method" binding:"required,oneof=credit_card paypal cod debit_card apple_pay wallet postpaid"`
	AddressID     *int64 `json:"address_id" binding:"omitempty,required"`
}

// ------------------------------ Mappers ------------------------------

func mapToCartResponse(cart sqlc.Cart, cartItems []sqlc.GetCartItemsRow) cartResponse {
	var totalPrice float64
	products := make([]cartItemResponse, len(cartItems))
	for i, item := range cartItems {
		price, _ := item.ProductPrice.Float64Value()
		totalPrice += price.Float64 * float64(item.Quantity)
		products[i] = cartItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Name:      item.ProductName,
			ImageURL:  item.ImageUrl.String,
			Quantity:  item.Quantity,
			Price:     price.Float64,
		}
	}

	return cartResponse{
		ID:         cart.ID,
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
// @Success 200 {object} cartResponse
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /carts [post]
func (sv *Server) createCart(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("user not found")))
		return
	}
	user, err := sv.postgres.GetUserByID(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(errors.New("user not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	_, err = sv.postgres.GetCart(c, authPayload.UserID)
	if err == nil {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("cart already existed")))
		return
	}

	newCart, err := sv.postgres.CreateCart(c, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusOK, mapDefaultResp(newCart, nil, nil))
}

// getCartDetail godoc
// @Summary Get cart details by user ID
// @Schemes http
// @Description get cart details by user ID
// @Tags carts
// @Accept json
// @Produce json
// @Success 200 {object} cartResponse
// @Failure 500 {object} gin.H
// @Router /carts [get]
func (sv *Server) getCartDetail(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}
	cart, err := sv.postgres.GetCart(c, authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	cartItems, err := sv.postgres.GetCartItems(c, cart.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	c.JSON(http.StatusOK, mapToCartResponse(cart, cartItems))
}

// addCartItem godoc
// @Summary Add a product to the cart
// @Schemes http
// @Description add a product to the cart
// @Tags carts
// @Accept json
// @Param input body addProductToCartRequest true "Add product to cart input"
// @Produce json
// @Success 200 {object} cartResponse
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /carts/products [post]
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

	product, err := sv.postgres.GetProduct(c, sqlc.GetProductParams{
		ID: req.ProductID,
		Archived: pgtype.Bool{
			Bool:  false,
			Valid: true,
		},
	})
	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(errors.New("product not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if product.Stock < int32(req.Quantity) {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("insufficient stock")))
		return
	}

	cart, err := sv.postgres.GetCart(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if cart.UserID != authPayload.UserID {
		c.JSON(http.StatusUnauthorized, mapErrResp(errors.New("cart does not belong to the user")))
		return
	}

	// check if the product is already in the cart
	cartItem, err := sv.postgres.GetCartItemByProductID(c, req.ProductID)
	if err != nil && errors.Is(err, postgres.ErrorRecordNotFound) {
		cartItem, err = sv.postgres.AddProductToCart(c, sqlc.AddProductToCartParams{
			ProductID: req.ProductID,
			CartID:    cart.ID,
			Quantity:  req.Quantity,
		})
	} else if err == nil {
		if int32(cartItem.Quantity+req.Quantity) > product.Stock {
			c.JSON(http.StatusBadRequest, mapErrResp(errors.New("insufficient stock")))
			return
		}
		err = sv.postgres.UpdateCartItemQuantity(c, sqlc.UpdateCartItemQuantityParams{
			Quantity: cartItem.Quantity + req.Quantity,
			ID:       cartItem.ID,
		})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	err = sv.postgres.UpdateCart(c, cart.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusOK, mapDefaultResp(cartItem, nil, nil))
}

// removeCartItem godoc
// @Summary Remove a product from the cart
// @Schemes http
// @Description remove a product from the cart
// @Tags carts
// @Accept json
// @Param id path int true "Product ID"
// @Produce json
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /carts/products [delete]
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

	cart, err := sv.postgres.GetCart(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if cart.UserID != authPayload.UserID {
		c.JSON(http.StatusUnauthorized, mapErrResp(errors.New("cart does not belong to the user")))
		return
	}

	err = sv.postgres.RemoveProductFromCart(c, sqlc.RemoveProductFromCartParams{
		CartID: cart.ID,
		ID:     param.ID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusOK, mapDefaultResp("product removed", nil, nil))
}

// checkout godoc
// @Summary Update product items in the cart
// @Schemes http
// @Description update product items in the cart
// @Tags carts
// @Accept json
// @Param input body checkoutRequest true "Update cart items input"
// @Produce json
// @Success 200 {object} sqlc.Order
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /carts/checkout [post]
func (sv *Server) checkout(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("user not found")))
		return
	}

	var req checkoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, mapErrResp(err))
		return
	}

	addresses, err := sv.postgres.GetAddresses(c, authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}
	if len(addresses) == 0 {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("no address found")))
		return
	}

	primaryAddressID := int64(addresses[0].ID)
	if req.AddressID == nil {
		isAddressExist := false
		for _, address := range addresses {
			if address.IsPrimary {
				primaryAddressID = address.ID
			}
			if address.ID == *req.AddressID {
				isAddressExist = true
				break
			}
		}
		if !isAddressExist {
			c.JSON(http.StatusBadRequest, mapErrResp(errors.New("address not found")))
			return
		}
	}

	cart, err := sv.postgres.GetCart(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	itemCnt, err := sv.postgres.CountCartItem(c, cart.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if itemCnt == 0 {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("cart is empty")))
		return
	}

	if cart.UserID != authPayload.UserID {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("cart does not belong to the user")))
		return
	}

	params := postgres.CheckoutCartTxParams{
		UserID:        authPayload.UserID,
		CartID:        cart.ID,
		PaymentMethod: sqlc.PaymentMethod(req.PaymentMethod),
	}

	if req.AddressID != nil {
		params.AddressID = *req.AddressID
	} else {
		params.AddressID = primaryAddressID
	}

	order, err := sv.postgres.CheckoutCartTx(c, params)

	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusOK, mapDefaultResp(order, nil, nil))
}

// updateCartItemQuantity godoc
// @Summary Update product items in the cart
// @Schemes http
// @Description update product items in the cart
// @Tags carts
// @Accept json
// @Param input body updateCartItemRequest true "Update cart items input"
// @Produce json
// @Success 200 {object} cartResponse
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /carts/products [put]
func (sv *Server) updateCartItemQuantity(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
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

	cart, err := sv.postgres.GetCart(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if cart.UserID != authPayload.UserID {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("cart does not belong to the user")))
		return
	}

	err = sv.postgres.UpdateCartItemQuantity(c, sqlc.UpdateCartItemQuantityParams{
		Quantity: req.Quantity,
		ID:       param.ID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	cartItem, err := sv.postgres.GetCartItem(c, param.ID)
	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(errors.New("cart item not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	product, err := sv.postgres.GetProduct(c, sqlc.GetProductParams{
		ID: cartItem.ProductID,
	})

	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
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

	c.JSON(http.StatusOK, mapDefaultResp(cartItem, nil, nil))
}

// clearCart godoc
// @Summary  Clear the cart
// @Schemes http
// @Description  clear the cart
// @Tags carts
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /carts/clear [delete]
func (sv *Server) clearCart(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("user not found")))
		return
	}

	cart, err := sv.postgres.GetCart(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, mapErrResp(errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	if cart.UserID != authPayload.UserID {
		c.JSON(http.StatusBadRequest, mapErrResp(errors.New("cart does not belong to the user")))
		return
	}

	err = sv.postgres.ClearCart(c, cart.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, mapErrResp(err))
		return
	}

	c.JSON(http.StatusOK, "cart cleared")
}
