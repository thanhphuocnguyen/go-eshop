package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/constants"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
)

// @Summary Create a new cart
// @Schemes http
// @Description create a new cart for a user
// @Tags carts
// @Accept json
// @Produce json
// @Success 200 {object} ApiResponse[CartDetailResponse]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 403 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Router /carts [post]
func (sv *Server) createCart(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	userID := claims["userId"].(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, errors.New("user not found")))
		return
	}
	user, err := sv.repo.GetUserByID(c, userID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, errors.New("user not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	_, err = sv.repo.GetCart(c, repository.GetCartParams{
		UserID: utils.GetPgTypeUUID(userID),
	})
	if err == nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, errors.New("cart already exists")))
		return
	}

	newCart, err := sv.repo.CreateCart(c, repository.CreateCartParams{
		UserID: utils.GetPgTypeUUID(user.ID),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	resp := &dto.CartDetail{
		ID:         newCart.ID,
		TotalPrice: 0,
		CartItems:  []dto.CartItemDetail{},
		CreatedAt:  newCart.CreatedAt,
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, resp, nil, nil))
}

// @Summary Get cart details by user ID
// @Schemes http
// @Description get cart details by user ID
// @Tags cart
// @Accept json
// @Produce json
// @Success 200 {object} ApiResponse[CartDetailResponse]
// @Failure 500 {object} ErrorResp
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 403 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Router /carts [get]
func (sv *Server) getCart(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if !ok {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, errors.New("user not found")))
		return
	}
	cart, err := sv.repo.GetCart(c, repository.GetCartParams{
		UserID: utils.GetPgTypeUUID(userID),
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			cart, err := sv.repo.CreateCart(c, repository.CreateCartParams{
				UserID: utils.GetPgTypeUUID(userID),
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
				return
			}
			c.JSON(http.StatusOK, dto.CreateDataResp(c, dto.CartDetail{
				ID:         cart.ID,
				TotalPrice: 0,
				CartItems:  []dto.CartItemDetail{},
				UpdatedAt:  &cart.UpdatedAt,
				CreatedAt:  cart.CreatedAt,
			}, nil, nil))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	cartItemRows, err := sv.repo.GetCartItems(c, cart.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	cartDetail := dto.CartDetail{
		ID:             cart.ID,
		TotalPrice:     0,
		DiscountAmount: 0,
		CartItems:      make([]dto.CartItemDetail, len(cartItemRows)),
		UpdatedAt:      &cart.UpdatedAt,
		CreatedAt:      cart.CreatedAt,
	}

	for i, row := range cartItemRows {
		item := mapToCartItemsResp(row)
		cartDetail.CartItems[i] = item
		cartDetail.TotalPrice += item.Price * float64(item.Quantity)
		cartDetail.DiscountAmount += item.DiscountAmount
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, cartDetail, nil, nil))
}

// @Summary Get cart discounts
// @Schemes http
// @Description get cart discounts
// @Tags carts
// @Accept json
// @Produce json
// @Success 200 {object} ApiResponse[[]repository.GetAvailableDiscountsForCartRow]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /carts/available-discounts [get]
func (sv *Server) getCartAvailableDiscounts(w http.ResponseWriter, r *http.Request) {
	authPayload, _ := c.MustGet(constants.AuthPayLoad).(*auth.TokenPayload)
	_, err := sv.repo.GetCart(c, repository.GetCartParams{
		UserID: utils.GetPgTypeUUID(userID),
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, struct{}{}, nil, nil))

}

// @Summary update product quantity in the cart
// @Schemes http
// @Description add a product to the cart
// @Tags carts
// @Accept json
// @Param input body UpdateCartItemQtyRequest true "Add product to cart input"
// @Produce json
// @Success 200 {object} ApiResponse[string]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /carts/items/{variant_id} [post]
func (sv *Server) updateCartItemQty(w http.ResponseWriter, r *http.Request) {
	var param models.UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, errors.New("invalid variant id")))
		return
	}

	_, claims, err := jwtauth.FromContext(r.Context())
	if !ok {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, errors.New("user not found")))
		return
	}

	var req models.UpdateCartItemQtyModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	cart, err := sv.repo.GetCart(c, repository.GetCartParams{
		UserID: utils.GetPgTypeUUID(userID),
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			newCart, createCartErr := sv.repo.CreateCart(c, repository.CreateCartParams{
				UserID: utils.GetPgTypeUUID(userID),
			})
			if createCartErr != nil {
				c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, createCartErr))
				return
			}
			cart = repository.GetCartRow{
				ID: newCart.ID,
			}
		} else {
			c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
			return
		}
	}

	cartItem, err := sv.repo.GetCartItem(c, repository.GetCartItemParams{
		ID:     uuid.MustParse(param.ID),
		CartID: cart.ID,
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			cartItem, err = sv.repo.AddCartItem(c, repository.AddCartItemParams{
				CartID:    cart.ID,
				VariantID: uuid.MustParse(param.ID),
				Quantity:  req.Quantity,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
			return
		}
	} else {
		err = sv.repo.UpdateCartItemQuantity(c, repository.UpdateCartItemQuantityParams{
			Quantity: req.Quantity,
			ID:       cartItem.ID,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
			return
		}
	}

	err = sv.repo.UpdateCartTimestamp(c, cart.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, cartItem.ID, nil, nil))
}

// @Summary Remove a product from the cart
// @Schemes http
// @Description remove a product from the cart
// @Tags carts
// @Accept json
// @Param id path int true "Product ID"
// @Produce json
// @Success 200 {object} ApiResponse[string]
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 403 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /carts/items/{id} [delete]
func (sv *Server) removeCartItem(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if !ok {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, errors.New("user not found")))
		return
	}

	var param models.UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	cart, err := sv.repo.GetCart(c, repository.GetCartParams{
		UserID: utils.GetPgTypeUUID(userID),
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	if cart.UserID.Valid {
		cartUserID, _ := uuid.FromBytes(cart.UserID.Bytes[:])
		if cartUserID != userID {
			c.JSON(http.StatusForbidden, dto.CreateErr("forbidden", errors.New("user not found")))
			return
		}
	}

	err = sv.repo.RemoveProductFromCart(c, repository.RemoveProductFromCartParams{
		CartID: cart.ID,
		ID:     uuid.MustParse(param.ID),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	message := "item removed"
	c.JSON(http.StatusOK, dto.CreateDataResp(c, message, nil, nil))
}

// @Summary  Clear the cart
// @Schemes http
// @Description  clear the cart
// @Tags carts
// @Accept json
// @Produce json
// @Success 204 {object} nil
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 403 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /carts/clear [put]
func (sv *Server) clearCart(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if !ok {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, errors.New("user not found")))
		return
	}

	cart, err := sv.repo.GetCart(c, repository.GetCartParams{
		UserID: utils.GetPgTypeUUID(userID),
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	if string(cart.UserID.Bytes[:]) != userID.String() {
		c.JSON(http.StatusForbidden, dto.CreateErr("forbidden", errors.New("user not found")))
		return
	}

	err = sv.repo.ClearCart(c, cart.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.Status(http.StatusNoContent)
}

// ------------------------------ Mappers ------------------------------
func mapToCartItemsResp(row repository.GetCartItemsRow) dto.CartItemDetail {

	// if it's the first item or the previous item is different
	var attr []dto.AttributeDetail
	err := json.Unmarshal(row.Attributes, &attr)
	if err != nil {
		log.Error().Err(err).Msg("Unmarshal cart item attributes")
	}
	price, _ := row.VariantPrice.Float64Value()
	qty := row.CartItem.Quantity
	amount := price.Float64 * float64(qty)
	cartItemsResp := dto.CartItemDetail{
		ID:         row.CartItem.ID.String(),
		ProductID:  row.ProductID.String(),
		VariantID:  row.CartItem.VariantID.String(),
		Name:       row.ProductName,
		Quantity:   row.CartItem.Quantity,
		Price:      price.Float64,
		StockQty:   row.VariantStock,
		Sku:        &row.VariantSku,
		ImageURL:   row.VariantImageUrl,
		Attributes: attr,
	}
	discountAmount := 0.0
	if row.ProductDiscountPercentage != nil {
		discountAmount = amount * float64(*row.ProductDiscountPercentage) / 100
		cartItemsResp.DiscountAmount = discountAmount
	}

	return cartItemsResp
}

// Setup cart-related routes
func (sv *Server) addCartRoutes(r chi.Router) {
	cart := rg.Group("/carts", authenticateMiddleware(sv.tokenGenerator))
	{
		cart.POST("", sv.createCart)
		cart.POST("checkout", sv.checkout)
		cart.GET("", sv.getCart)
		cart.PUT("clear", sv.clearCart)

		cart.GET("available-discounts", sv.getCartAvailableDiscounts)
		cartItems := cart.Group("items")
		{
			cartItems.PUT(":id/quantity", sv.updateCartItemQty)
			cartItems.DELETE(":id", sv.removeCartItem)
		}
	}
}
