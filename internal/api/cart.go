package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
)

// @Summary Create a new cart
// @Schemes http
// @Description create a new cart for a user
// @Tags carts
// @Accept json
// @Produce json
// @Success 200 {object} dto.ApiResponse[dto.CartDetail]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 403 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Router /carts [post]
func (s *Server) createCart(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	_, claims, err := jwtauth.FromContext(c)

	userID := uuid.MustParse(claims["userId"].(string))

	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, errors.New("user not found"))
		return
	}
	user, err := s.repo.GetUserByID(c, userID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, errors.New("user not found"))
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	_, err = s.repo.GetCart(c, repository.GetCartParams{
		UserID: utils.GetPgTypeUUID(userID),
	})
	if err == nil {
		RespondBadRequest(w, InvalidBodyCode, errors.New("cart already exists"))
		return
	}

	newCart, err := s.repo.CreateCart(c, repository.CreateCartParams{
		UserID: utils.GetPgTypeUUID(user.ID),
	})
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	resp := &dto.CartDetail{
		ID:         newCart.ID,
		TotalPrice: 0,
		CartItems:  []dto.CartItemDetail{},
		CreatedAt:  newCart.CreatedAt,
	}

	RespondSuccess(w, r, dto.CreateDataResp(c, resp, nil, nil))
}

// @Summary Get cart details by user ID
// @Schemes http
// @Description get cart details by user ID
// @Tags cart
// @Accept json
// @Produce json
// @Success 200 {object} dto.ApiResponse[dto.CartDetail]
// @Failure 500 {object} ErrorResp
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 403 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Router /carts [get]
func (s *Server) getCart(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	_, claims, err := jwtauth.FromContext(c)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, errors.New("user not found"))
		return
	}

	userID := uuid.MustParse(claims["userId"].(string))
	cart, err := s.repo.GetCart(c, repository.GetCartParams{
		UserID: utils.GetPgTypeUUID(userID),
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			cart, err := s.repo.CreateCart(c, repository.CreateCartParams{
				UserID: utils.GetPgTypeUUID(userID),
			})
			if err != nil {
				RespondInternalServerError(w, InternalServerErrorCode, err)
				return
			}
			RespondSuccess(w, r, dto.CreateDataResp(c, dto.CartDetail{
				ID:         cart.ID,
				TotalPrice: 0,
				CartItems:  []dto.CartItemDetail{},
				UpdatedAt:  &cart.UpdatedAt,
				CreatedAt:  cart.CreatedAt,
			}, nil, nil))
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	cartItemRows, err := s.repo.GetCartItems(c, cart.ID)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
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

	RespondSuccess(w, r, dto.CreateDataResp(c, cartDetail, nil, nil))
}

// @Summary Get cart discounts
// @Schemes http
// @Description get cart discounts
// @Tags carts
// @Accept json
// @Produce json
// @Success 200 {object} dto.ApiResponse[[]repository.GetCartRow]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /carts/available-discounts [get]
func (s *Server) getCartAvailableDiscounts(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	_, claims, err := jwtauth.FromContext(c)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, errors.New("authorization payload is not provided"))
		return
	}
	userID := uuid.MustParse(claims["userId"].(string))

	cart, err := s.repo.GetCart(c, repository.GetCartParams{
		UserID: utils.GetPgTypeUUID(userID),
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, errors.New("cart not found"))
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondSuccess(w, r, dto.CreateDataResp(c, cart, nil, nil))

}

// @Summary update product quantity in the cart
// @Schemes http
// @Description add a product to the cart
// @Tags carts
// @Accept json
// @Param input body models.UpdateCartItemQtyModel true "Add product to cart input"
// @Produce json
// @Success 200 {object} dto.ApiResponse[string]
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /carts/items/{variant_id} [post]
func (s *Server) updateCartItemQty(w http.ResponseWriter, r *http.Request) {
	id, err := GetUrlParam(r, "variant_id")
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, errors.New("user not found"))
		return
	}
	userID := uuid.MustParse(claims["userId"].(string))

	c := r.Context()

	var req models.UpdateCartItemQtyModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	cart, err := s.repo.GetCart(c, repository.GetCartParams{
		UserID: utils.GetPgTypeUUID(userID),
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			newCart, createCartErr := s.repo.CreateCart(c, repository.CreateCartParams{
				UserID: utils.GetPgTypeUUID(userID),
			})
			if createCartErr != nil {
				RespondInternalServerError(w, InternalServerErrorCode, createCartErr)
				return
			}
			cart = repository.GetCartRow{
				ID: newCart.ID,
			}
		} else {
			RespondInternalServerError(w, InternalServerErrorCode, err)
			return
		}
	}
	uuid := uuid.MustParse(id)
	cartItem, err := s.repo.GetCartItem(c, repository.GetCartItemParams{
		ID:     uuid,
		CartID: cart.ID,
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			cartItem, err = s.repo.AddCartItem(c, repository.AddCartItemParams{
				CartID:    cart.ID,
				VariantID: uuid,
				Quantity:  req.Quantity,
			})
			if err != nil {
				RespondInternalServerError(w, InternalServerErrorCode, err)
				return
			}
		} else {
			RespondInternalServerError(w, InternalServerErrorCode, err)
			return
		}
	} else {
		err = s.repo.UpdateCartItemQuantity(c, repository.UpdateCartItemQuantityParams{
			Quantity: req.Quantity,
			ID:       cartItem.ID,
		})

		if err != nil {
			RespondInternalServerError(w, InternalServerErrorCode, err)
			return
		}
	}

	err = s.repo.UpdateCartTimestamp(c, cart.ID)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondSuccess(w, r, dto.CreateDataResp(c, cartItem.ID, nil, nil))
}

// @Summary Remove a product from the cart
// @Schemes http
// @Description remove a product from the cart
// @Tags carts
// @Accept json
// @Param id path int true "Product ID"
// @Produce json
// @Success 200 {object} dto.ApiResponse[string]
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 403 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /carts/items/{id} [delete]
func (s *Server) removeCartItem(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, errors.New("authorization payload is not provided"))
		return
	}
	userID := uuid.MustParse(claims["userId"].(string))

	c := r.Context()
	id, err := GetUrlParam(r, "id")
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	cart, err := s.repo.GetCart(c, repository.GetCartParams{
		UserID: utils.GetPgTypeUUID(userID),
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, errors.New("cart not found"))
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	if cart.UserID.Valid {
		cartUserID, _ := uuid.FromBytes(cart.UserID.Bytes[:])
		if cartUserID != userID {
			RespondForbidden(w, UnauthorizedCode, errors.New("forbidden"))
			return
		}
	}

	err = s.repo.RemoveProductFromCart(c, repository.RemoveProductFromCartParams{
		CartID: cart.ID,
		ID:     uuid.MustParse(id),
	})
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondNoContent(w)
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
func (s *Server) clearCart(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	c := r.Context()
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, errors.New("authorization payload is not provided"))
		return
	}
	userID := uuid.MustParse(claims["userId"].(string))

	cart, err := s.repo.GetCart(c, repository.GetCartParams{
		UserID: utils.GetPgTypeUUID(userID),
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, errors.New("cart not found"))
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	if string(cart.UserID.Bytes[:]) != userID.String() {
		RespondForbidden(w, UnauthorizedCode, errors.New("forbidden"))
		return
	}

	err = s.repo.ClearCart(c, cart.ID)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
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
func (s *Server) addCartRoutes(r chi.Router) {
	r.Route("/carts", func(r chi.Router) {
		r.Post("/", s.createCart)
		r.Post("/checkout", s.checkout)
		r.Get("/", s.getCart)
		r.Put("/clear", s.clearCart)

		r.Get("/available-discounts", s.getCartAvailableDiscounts)
		r.Route("/items", func(r chi.Router) {
			r.Put("/{id}/quantity", s.updateCartItemQty)
			r.Delete("/{id}", s.removeCartItem)
		})
	})
}
