package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/postgres"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/sqlc"
)

type createCartRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
}

func (sv *Server) createCart(c *gin.Context) {
	var cart createCartRequest
	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	newCart, err := sv.postgres.CreateCart(c, cart.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, newCart)
}

func (sv *Server) getCart(c *gin.Context) {
	userID := c.GetInt64("user_id")

	cart, err := sv.postgres.GetCartByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, mapToCartResponse(cart))
}

type cartItem struct {
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
	ID           int64      `json:"id"`
	CheckedOutAt time.Time  `json:"checked_out"`
	UserID       int64      `json:"user_id"`
	UpdatedAt    time.Time  `json:"updated_at"`
	CreatedAt    time.Time  `json:"created_at"`
	CartItems    []cartItem `json:"cart_items"`
}

func mapToCartResponse(cart []sqlc.GetCartByUserIDRow) cartResponse {
	if len(cart) == 0 {
		return cartResponse{}
	}

	cartItems := make([]cartItem, 0)
	for _, item := range cart {
		prodPrice, _ := item.Product.Price.Float64Value()
		cartItems = append(cartItems, cartItem{
			ID:          item.CartItem.ID,
			ProductID:   item.Product.ID,
			Name:        item.Product.Name,
			Description: item.Product.Description,
			SKU:         item.Product.Sku,
			ImageURL:    item.Product.ImageUrl.String,
			Quantity:    item.CartItem.Quantity,
			Price:       float64(item.CartItem.Quantity) * prodPrice.Float64,
		})
	}

	return cartResponse{
		ID:           cart[0].Cart.ID,
		CheckedOutAt: cart[0].Cart.CheckedOutAt.Time,
		UserID:       cart[0].Cart.UserID,
		UpdatedAt:    cart[0].Cart.UpdatedAt,
		CreatedAt:    cart[0].Cart.CreatedAt,
		CartItems:    cartItems,
	}
}

type addProductToCartRequest struct {
	UserID    int64 `json:"user_id" binding:"required"`
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int16 `json:"quantity" binding:"required"`
}

func (sv *Server) addProductToCart(c *gin.Context) {
	var req addProductToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var cartId int64
	cart, err := sv.postgres.GetCartByUserID(c, req.UserID)
	if err != nil {
		if err == postgres.ErrorRecordNotFound {
			newCart, err := sv.postgres.CreateCart(c, req.UserID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
			cartId = newCart.ID
		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	} else {
		cartId = cart[0].Cart.ID
	}

	_, err = sv.postgres.AddProductToCart(c, sqlc.AddProductToCartParams{
		ProductID: req.ProductID,
		CartID:    cartId,
		Quantity:  req.Quantity,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, "product added to cart")
}

type removeProductFromCartRequest struct {
	CartID    int64 `json:"cart_id" binding:"required"`
	ProductID int64 `json:"product_id" binding:"required"`
}

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
func (sv *Server) checkout(c *gin.Context) {

}

type itemUpdate struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int16 `json:"quantity" binding:"required"`
}
type updateCartRequest struct {
	CartID int64        `json:"cart_id" binding:"required"`
	Items  []itemUpdate `json:"items" binding:"required"`
}

func (sv *Server) updateCartProductItems(c *gin.Context) {
	var req updateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	errChan := make(chan error)
	for _, item := range req.Items {
		go func(item itemUpdate) {
			err := sv.postgres.UpdateProductQuantity(c, sqlc.UpdateProductQuantityParams{
				CartID:    req.CartID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
			})
			if err != nil {
				errChan <- err
				return
			}
		}(item)
	}
	errs := []error{}
	for err := range errChan {
		errs = append(errs, err)
		if len(errs) == len(req.Items) {
			break
		}
	}
	if len(errs) > 0 {
		c.JSON(http.StatusInternalServerError, errorsResponse(errs))
	}

	c.JSON(http.StatusOK, "cart updated")
}
