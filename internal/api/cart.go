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

type updateCartItemRequest struct {
	ID       int64 `json:"id" binding:"required,gt=0"`
	Quantity int16 `json:"quantity" binding:"required,gt=0"`
}

type cartItemResponse struct {
	ID          int64   `json:"id"`
	ProductID   int64   `json:"product_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	SKU         string  `json:"sku"`
	ImageURL    string  `json:"image_url"`
	Quantity    int16   `json:"quantity"`
	Price       float64 `json:"price"`
}

type cartResponse struct {
	ID           int64              `json:"id"`
	CheckedOutAt time.Time          `json:"checked_out"`
	UserID       int64              `json:"user_id"`
	UpdatedAt    time.Time          `json:"updated_at"`
	CreatedAt    time.Time          `json:"created_at"`
	CartItems    []cartItemResponse `json:"cart_items"`
}

type addProductToCartRequest struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int16 `json:"quantity" binding:"required"`
}

type removeProductFromCartRequest struct {
	CartID    int64 `json:"cart_id" binding:"required"`
	ProductID int64 `json:"product_id" binding:"required"`
}

type checkoutRequest struct {
	CartID      int64  `json:"cart_id" binding:"required"`
	IsCod       bool   `json:"is_cod" binding:"required"`
	PaymentType string `json:"payment_type" binding:"required"`
}

// ------------------------------ Mappers ------------------------------

func mapToCartResponse(cartItems []sqlc.GetCartDetailRow) cartResponse {
	if len(cartItems) == 0 {
		return cartResponse{}
	}

	products := make([]cartItemResponse, 0)
	for _, item := range cartItems {
		product := item.Product
		cartItem := item.CartItem
		prodPrice, _ := product.Price.Float64Value()
		products = append(products, cartItemResponse{
			ID:          cartItem.ID,
			ProductID:   product.ID,
			Name:        product.Name,
			Description: product.Description,
			SKU:         product.Sku,
			ImageURL:    product.ImageUrl.String,
			Quantity:    cartItem.Quantity,
			Price:       float64(cartItem.Quantity) * prodPrice.Float64,
		})
	}
	cart := cartItems[0].Cart

	return cartResponse{
		ID:           cart.ID,
		CheckedOutAt: cart.CheckedOutAt.Time,
		UserID:       cart.UserID,
		UpdatedAt:    cart.UpdatedAt,
		CreatedAt:    cart.CreatedAt,
		CartItems:    products,
	}
}

// ------------------------------ Handlers ------------------------------

// CreateCart godoc
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}
	user, err := sv.postgres.GetUserByID(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user is not existed"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	_, err = sv.postgres.GetCart(c, authPayload.UserID)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cart already existed"})
		return
	}

	newCart, err := sv.postgres.CreateCart(c, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, newCart)
}

// GetCart godoc
// @Summary Get cart details by user ID
// @Schemes http
// @Description get cart details by user ID
// @Tags carts
// @Accept json
// @Produce json
// @Success 200 {object} cartResponse
// @Failure 500 {object} gin.H
// @Router /carts [get]
func (sv *Server) getCart(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	cartDetails, err := sv.postgres.GetCartDetail(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "cart not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if len(cartDetails) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "cart not found"})
		return
	}

	c.JSON(http.StatusOK, mapToCartResponse(cartDetails))
}

// AddProductToCart godoc
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
func (sv *Server) addProductToCart(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	var req addProductToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	product, err := sv.postgres.GetProduct(c, req.ProductID)
	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	cartDetail, err := sv.postgres.GetCartDetail(c, authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if len(cartDetail) == 0 {
		newCart, err := sv.postgres.CreateCart(c, authPayload.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		cartItem, err := sv.postgres.AddProductToCart(c, sqlc.AddProductToCartParams{
			ProductID: req.ProductID,
			CartID:    cartDetail[0].Cart.ID,
			Quantity:  req.Quantity,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		cartDetail = append(cartDetail, sqlc.GetCartDetailRow{
			Cart:     newCart,
			Product:  product,
			CartItem: cartItem,
		})
	} else {
		cartItem, err := sv.postgres.AddProductToCart(c, sqlc.AddProductToCartParams{
			ProductID: req.ProductID,
			CartID:    cartDetail[0].Cart.ID,
			Quantity:  req.Quantity,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		cartDetail = append(cartDetail, sqlc.GetCartDetailRow{
			Cart:     cartDetail[0].Cart,
			CartItem: cartItem,
			Product:  product,
		})
	}

	err = sv.postgres.UpdateCart(c, cartDetail[0].Cart.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, mapToCartResponse(cartDetail))
}

// removeProductFromCart godoc
// @Summary Remove a product from the cart
// @Schemes http
// @Description remove a product from the cart
// @Tags carts
// @Accept json
// @Param input body removeProductFromCartRequest true "Remove product from cart input"
// @Produce json
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /carts/products [delete]
func (sv *Server) removeProductFromCart(c *gin.Context) {
	var req removeProductFromCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := sv.postgres.RemoveProductFromCart(c, sqlc.RemoveProductFromCartParams{
		CartID:    req.CartID,
		ProductID: req.ProductID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, "product removed from cart")
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
// @Router /carts/products [post]
func (sv *Server) checkout(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	var req checkoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	cart, err := sv.postgres.GetCart(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "cart not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if cart.UserID != user.UserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cart does not belong to the user"})
		return
	}

	order, err := sv.postgres.CheckoutCartTx(c, postgres.CheckoutCartParams{
		UserID: user.UserID,
		CartID: cart.ID,
		CreateOrderParams: sqlc.CreateOrderParams{
			UserID:      user.UserID,
			PaymentType: sqlc.PaymentType(req.PaymentType),
			IsCod:       req.IsCod,
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, order.Order)
}

// updateCartItemQuantity godoc
// @Summary Update product items in the cart
// @Schemes http
// @Description update product items in the cart
// @Tags carts
// @Accept json
// @Param input body updateCartRequest true "Update cart items input"
// @Produce json
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /carts/products [put]
func (sv *Server) updateCartItemQuantity(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}
	var req updateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cartDetail, err := sv.postgres.GetCartDetail(c, authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if len(cartDetail) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "cart not found"})
		return
	}

	if cartDetail[0].Cart.UserID != authPayload.UserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cart does not belong to the user"})
		return
	}

	for i, item := range cartDetail {
		if item.CartItem.ID == req.ID {
			err := sv.postgres.UpdateCartItemQuantity(c, sqlc.UpdateCartItemQuantityParams{
				ID:       req.ID,
				Quantity: req.Quantity,
			})
			cartDetail[i].CartItem.Quantity = req.Quantity
			if err != nil {
				c.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		}
	}

	c.JSON(http.StatusOK, mapToCartResponse(cartDetail))
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	cart, err := sv.postgres.GetCart(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, postgres.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "cart not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if cart.UserID != authPayload.UserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cart does not belong to the user"})
		return
	}

	err = sv.postgres.ClearCart(c, cart.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, "cart cleared")
}
