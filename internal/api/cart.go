package api

import (
	"errors"
	"fmt"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v81"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
	"github.com/thanhphuocnguyen/go-eshop/pkg/paymentsrv"
)

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

// @Summary Get cart discounts
// @Schemes http
// @Description get cart discounts
// @Tags carts
// @Accept json
// @Produce json
// @Success 200 {object} ApiResponse[gin.H]
// @Failure 400 {object} ApiResponse[gin.H]
// @Failure 500 {object} ApiResponse[gin.H]
// @Router /cart/discounts [get]
func (sv *Server) getCartDiscountsHandler(c *gin.Context) {
	authPayload, _ := c.MustGet(authorizationPayload).(*auth.Payload)
	cart, err := sv.repo.GetCart(c, repository.GetCartParams{
		UserID: utils.GetPgTypeUUID(authPayload.UserID),
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[gin.H](NotFoundCode, "", errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "", err))
		return
	}

	cartDiscounts, err := sv.repo.GetAvailableDiscountsForCart(c, cart.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, cartDiscounts, "Success!", nil, nil))

}

// @Summary Add a product to the cart
// @Schemes http
// @Description add a product to the cart
// @Tags carts
// @Accept json
// @Param input body UpdateCartItemQtyRequest true "Add product to cart input"
// @Produce json
// @Success 200 {object} ApiResponse[uuid.UUID]
// @Failure 400 {object} ApiResponse[gin.H]
// @Failure 500 {object} ApiResponse[gin.H]
// @Router /cart/item/{variant_id} [post]
func (sv *Server) updateCartItemQtyHandler(c *gin.Context) {
	var param UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[uuid.UUID](InvalidBodyCode, "", errors.New("invalid variant id")))
		return
	}

	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, createErrorResponse[uuid.UUID](InvalidBodyCode, "", errors.New("user not found")))
		return
	}

	var req UpdateCartItemQtyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[uuid.UUID](InternalServerErrorCode, "", err))
		return
	}

	cart, err := sv.repo.GetCart(c, repository.GetCartParams{
		UserID: utils.GetPgTypeUUID(authPayload.UserID),
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			newCart, createCartErr := sv.repo.CreateCart(c, repository.CreateCartParams{
				UserID: utils.GetPgTypeUUID(authPayload.UserID),
			})
			if createCartErr != nil {
				c.JSON(http.StatusInternalServerError, createErrorResponse[uuid.UUID](InternalServerErrorCode, "", createCartErr))
				return
			}
			cart = newCart
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

	var param UriIDParam
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

	if cart.UserID.Valid {
		cartUserID, _ := uuid.FromBytes(cart.UserID.Bytes[:])
		if cartUserID != authPayload.UserID {
			c.JSON(http.StatusForbidden, createErrorResponse[string]("forbidden", "", errors.New("user not found")))
			return
		}
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
// @Success 200 {object} ApiResponse[gin.H]
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /cart/checkoutHandler [post]
func (sv *Server) checkoutHandler(c *gin.Context) {
	authPayload, ok := c.MustGet(authorizationPayload).(*auth.Payload)
	if !ok {
		c.JSON(http.StatusBadRequest, createErrorResponse[gin.H](InvalidBodyCode, "", errors.New("user not found")))
		return
	}

	var req CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse[gin.H](InvalidBodyCode, "", err))
		return
	}

	user, err := sv.repo.GetUserByID(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[gin.H](NotFoundCode, "", errors.New("user not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "", err))
		return
	}

	cart, err := sv.repo.GetCart(c, repository.GetCartParams{
		UserID: utils.GetPgTypeUUID(authPayload.UserID),
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, createErrorResponse[gin.H](InternalServerErrorCode, "", errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "", err))
		return
	}
	var shippingAddr repository.ShippingAddressSnapshot
	if req.AddressID == nil {
		// create new address
		if req.Address == nil {
			c.JSON(http.StatusBadRequest, createErrorResponse[gin.H](InternalServerErrorCode, "", errors.New("must provide address or address ID")))
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
			c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "", err))
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
			ID:     uuid.MustParse(*req.AddressID),
			UserID: user.ID,
		})
		if err != nil {
			if errors.Is(err, repository.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, createErrorResponse[gin.H](NotFoundCode, "", errors.New("address not found")))
				return
			}
			c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "", err))
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
			c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "", err))
			return
		}
		if cartUserId != authPayload.UserID {
			c.JSON(http.StatusForbidden, createErrorResponse[gin.H](PermissionDeniedCode, "", errors.New("you are not allowed to access this cart")))
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
	var totalPrice float64

	for _, item := range cartItemRows {
		variantIdx := -1
		for j, param := range createOrderItemParams {
			if param.VariantID == item.CartItem.VariantID {
				variantIdx = j
				break
			}
		}

		if variantIdx == -1 {
			price, _ := item.Price.Float64Value()
			itemParam := repository.CreateBulkOrderItemsParams{
				VariantID:            item.CartItem.VariantID,
				Quantity:             item.CartItem.Quantity,
				PricePerUnitSnapshot: item.Price,
				VariantSkuSnapshot:   item.Sku,
				ProductNameSnapshot:  item.ProductName,
				LineTotalSnapshot:    utils.GetPgNumericFromFloat(float64(item.CartItem.Quantity) * price.Float64),
				AttributesSnapshot: []repository.AttributeDataSnapshot{{
					Name:  item.AttrName,
					Value: item.AttrValCode,
				}},
			}

			createOrderItemParams = append(createOrderItemParams, itemParam)
			totalPrice += price.Float64 * float64(item.CartItem.Quantity)
		} else {
			createOrderItemParams[variantIdx].AttributesSnapshot = append(
				createOrderItemParams[variantIdx].AttributesSnapshot,
				repository.AttributeDataSnapshot{
					Name:  item.AttrName,
					Value: item.AttrValCode,
				})
		}
	}

	params := repository.CheckoutCartTxArgs{
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
	var discountPrice float64
	var discountID *uuid.UUID

	// check if there is a discount code
	if req.DiscountCode != nil && *req.DiscountCode != "" {
		discount, err := sv.repo.GetDiscountByCode(c, *req.DiscountCode)
		discountID = &discount.ID
		discountProductsAndCategories, err := sv.repo.GetDiscountProductsAndCategories(c, discount.ID)
		if err != nil {
			log.Error().Err(err).Msg("GetDiscountByCode")
			c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "", err))
			return
		}
		discountProductIDs := make(map[uuid.UUID]bool)
		discountCategoryIDs := make(map[uuid.UUID]bool)
		productIDs := make(map[uuid.UUID]bool)
		for _, item := range discountProductsAndCategories {
			if item.ProductID.Valid {
				id, _ := uuid.FromBytes(item.ProductID.Bytes[:])
				discountProductIDs[id] = true
			}
			if item.CategoryID.Valid {
				id, _ := uuid.FromBytes(item.CategoryID.Bytes[:])
				discountCategoryIDs[id] = true
			}
		}

		discountValue, _ := discount.DiscountValue.Float64Value()
		for _, item := range cartItemRows {
			if productIDs[item.ProductID] {
				// skip if the product is already counted
				continue
			}
			productIDs[item.ProductID] = true
			price, _ := item.Price.Float64Value()
			lineTotal := price.Float64 * float64(item.CartItem.Quantity)
			if _, ok := discountProductIDs[item.ProductID]; ok {
				// if the item has a discount price, use it
				if item.Price.Valid {
					if discount.DiscountType == string(repository.PercentageDiscount) {
						discountPrice += lineTotal * (discountValue.Float64 / 100)
					} else {
						discountPrice += lineTotal - discountValue.Float64
					}
				}
			} else if item.CategoryID.Valid {
				catId, _ := uuid.FromBytes(item.CategoryID.Bytes[:])
				if _, ok := discountCategoryIDs[catId]; ok {
					if discount.DiscountType == string(repository.PercentageDiscount) {
						discountPrice += lineTotal * (discountValue.Float64 / 100)
					} else {
						discountPrice += lineTotal - discountValue.Float64
					}
				} else {
					log.Warn().Msgf("Category %s not found in discount categories", catId.String())
				}
			}
		}

	}

	params.DiscountPrice = discountPrice
	params.DiscountID = discountID
	params.PaymentMethod = repository.PaymentMethod(req.PaymentMethod)

	params.CreatePaymentFn = func(orderID uuid.UUID, paymentMethod repository.PaymentMethod) (paymentIntentID string, clientSecretID *string, err error) {
		switch paymentMethod {
		case repository.PaymentMethodStripe:
			stripeInstance, err := paymentsrv.NewStripePayment(sv.config.StripeSecretKey)
			if err != nil {
				return "", utils.StringPtr(""), err
			}
			sv.paymentCtx.SetStrategy(stripeInstance)
		default:
			return "", utils.StringPtr(""), fmt.Errorf("unsupported payment method: %s", paymentMethod)
		}
		receiptEmail := user.Email
		// create payment intent
		checkoutResult, err := sv.paymentCtx.CreatePaymentIntent(totalPrice, receiptEmail)
		if err != nil {
			return "", utils.StringPtr(""), err
		}

		paymentIntent, ok := checkoutResult.(*stripe.PaymentIntent)
		if !ok {
			return "", utils.StringPtr(""), fmt.Errorf("unexpected payment intent type: %T", checkoutResult)
		}
		return paymentIntent.ID, &paymentIntent.ClientSecret, nil
	}

	rs, err := sv.repo.CheckoutCartTx(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse[gin.H](InternalServerErrorCode, "", err))
		return
	}

	c.JSON(http.StatusOK, createSuccessResponse(c, rs, "", nil, nil))
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
				StockQty:  row.StockQty,
				Sku:       &row.Sku,
				Attributes: []repository.AttributeDataSnapshot{
					attr,
				},
				ImageURL: row.ImageUrl,
			}

			// Populate CategoryID if available
			if row.CategoryID.Valid {
				categoryID := row.CategoryID.Bytes
				categoryUUID, err := uuid.FromBytes(categoryID[:])
				if err == nil {
					categoryIDStr := categoryUUID.String()
					cartItem.CategoryID = &categoryIDStr
				}
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
