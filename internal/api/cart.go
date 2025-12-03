package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/constants"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
	"github.com/thanhphuocnguyen/go-eshop/internal/processors"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
	"github.com/thanhphuocnguyen/go-eshop/pkg/payment"
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
func (sv *Server) CreateCart(c *gin.Context) {
	authPayload, ok := c.MustGet(constants.AuthPayLoad).(*auth.TokenPayload)
	if !ok {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, errors.New("user not found")))
		return
	}
	user, err := sv.repo.GetUserByID(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, errors.New("user not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	_, err = sv.repo.GetCart(c, repository.GetCartParams{
		UserID: utils.GetPgTypeUUID(authPayload.UserID),
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
func (sv *Server) GetCartHandler(c *gin.Context) {
	authPayload, ok := c.MustGet(constants.AuthPayLoad).(*auth.TokenPayload)
	if !ok {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, errors.New("user not found")))
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
func (sv *Server) GetCartAvailableDiscountsHandler(c *gin.Context) {
	authPayload, _ := c.MustGet(constants.AuthPayLoad).(*auth.TokenPayload)
	_, err := sv.repo.GetCart(c, repository.GetCartParams{
		UserID: utils.GetPgTypeUUID(authPayload.UserID),
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
func (sv *Server) UpdateCartItemQtyHandler(c *gin.Context) {
	var param models.UriIDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, errors.New("invalid variant id")))
		return
	}

	authPayload, ok := c.MustGet(constants.AuthPayLoad).(*auth.TokenPayload)
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
		UserID: utils.GetPgTypeUUID(authPayload.UserID),
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			newCart, createCartErr := sv.repo.CreateCart(c, repository.CreateCartParams{
				UserID: utils.GetPgTypeUUID(authPayload.UserID),
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
func (sv *Server) RemoveCartItem(c *gin.Context) {
	authPayload, ok := c.MustGet(constants.AuthPayLoad).(*auth.TokenPayload)
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
		UserID: utils.GetPgTypeUUID(authPayload.UserID),
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
		if cartUserID != authPayload.UserID {
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

// @Summary Update product items in the cart
// @Schemes http
// @Description update product items in the cart
// @Tags carts
// @Accept json
// @Param input body CheckoutRequest true "Checkout input"
// @Produce json
// @Success 200 {object} ApiResponse[repository.CreatePaymentResult]
// @Failure 400 {object} ErrorResp
// @Failure 404 {object} ErrorResp
// @Failure 403 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /carts/checkout [post]
func (sv *Server) CheckoutHandler(c *gin.Context) {
	authPayload, ok := c.MustGet(constants.AuthPayLoad).(*auth.TokenPayload)
	if !ok {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, errors.New("user not found")))
		return
	}

	var req models.CheckoutModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	user, err := sv.repo.GetUserDetailsByID(c, authPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, errors.New("user not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	cart, err := sv.repo.GetCart(c, repository.GetCartParams{
		UserID: utils.GetPgTypeUUID(authPayload.UserID),
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(InternalServerErrorCode, errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}
	var shippingAddr repository.ShippingAddressSnapshot

	address, err := sv.repo.GetAddress(c, repository.GetAddressParams{
		ID:     uuid.MustParse(req.AddressId),
		UserID: user.ID,
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, errors.New("address not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	shippingAddr = repository.ShippingAddressSnapshot{
		Street:   address.Street,
		Ward:     *address.Ward,
		District: address.District,
		City:     address.City,
		Phone:    address.PhoneNumber,
	}

	if cart.UserID.Valid {
		cartUserId, err := uuid.FromBytes(cart.UserID.Bytes[:])
		if err != nil {
			log.Error().Err(err).Msg("GetCart")
			c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
			return
		}
		if cartUserId != authPayload.UserID {
			c.JSON(http.StatusForbidden, dto.CreateErr(PermissionDeniedCode, errors.New("you are not allowed to access this cart")))
			return
		}
	}

	itemRows, err := sv.repo.GetCartItems(c, cart.ID)
	if err != nil {
		log.Error().Err(err).Msg("GetCartItems")
		return
	}

	// Process discounts
	discountResult, err := sv.discountProcessor.ProcessDiscounts(c, processors.DiscountContext{User: user, CartItems: itemRows}, req.DiscountCodes)
	if err != nil {
		log.Error().Err(err).Msg("ProcessDiscounts")
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, err))
		return
	}

	// Calculate totals and create order items
	var totalPrice float64
	createOrderItemParams := make([]repository.CreateBulkOrderItemsParams, len(itemRows))

	for i, item := range itemRows {
		itemPrice, _ := item.VariantPrice.Float64Value()
		lineTotal := float64(item.CartItem.Quantity) * itemPrice.Float64
		totalPrice += lineTotal

		// Find discount for this item
		discountAmount := 0.0
		if item.ProductDiscountPercentage != nil {
			discountAmount = lineTotal * (float64(*item.ProductDiscountPercentage) / 100)
		}

		// Add applied discounts
		for _, itemDiscount := range discountResult.ItemDiscounts {
			if itemDiscount.ItemIndex == i {
				discountAmount += itemDiscount.DiscountAmount
			}
		}

		createOrderItemParams[i] = repository.CreateBulkOrderItemsParams{
			VariantID:            item.CartItem.VariantID,
			Quantity:             item.CartItem.Quantity,
			PricePerUnitSnapshot: item.VariantPrice,
			VariantSkuSnapshot:   item.VariantSku,
			ProductNameSnapshot:  item.ProductName,
			LineTotalSnapshot:    utils.GetPgNumericFromFloat(lineTotal),
			DiscountedPrice:      utils.GetPgNumericFromFloat(discountAmount),
		}

		err := json.Unmarshal(item.Attributes, &createOrderItemParams[i].AttributesSnapshot)
		if err != nil {
			log.Error().Err(err).Msg("Unmarshal attributes")
			c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
			return
		}
	}

	params := repository.CheckoutCartTxArgs{
		CartID:          cart.ID,
		TotalPrice:      totalPrice,
		ShippingAddress: shippingAddr,
		UserID:          authPayload.UserID,
		CustomerInfo: repository.CustomerInfoTxArgs{
			FullName: user.FirstName + " " + user.LastName,
			Email:    user.Email,
			Phone:    user.PhoneNumber,
		},
		PaymentGateway:        &req.PaymentMethodId,
		CreateOrderItemParams: createOrderItemParams,
		DiscountPrice:         discountResult.TotalDiscount,
		DiscountIDs:           discountResult.AppliedDiscounts,
		PaymentMethodID:       uuid.MustParse(req.PaymentMethodId),
	}

	params.CreatePaymentFn = func(orderID uuid.UUID, method string) (paymentIntentID string, clientSecretID *string, err error) {
		// create payment intent
		intent, err := sv.paymentSrv.CreatePaymentIntent(c, method, payment.PaymentRequest{
			Amount:   int64((totalPrice - discountResult.TotalDiscount) * 100), // convert to cents
			Currency: "usd",
			Email:    user.Email,
			Metadata: map[string]string{
				"OrderID": orderID.String(),
			},
		})

		if err != nil {
			log.Error().Err(err).Msg("CreatePaymentIntent")
			return "", nil, err
		}
		return intent.ID, &intent.ClientSecret, nil
	}

	rs, err := sv.repo.CheckoutCartTx(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	c.JSON(http.StatusOK, dto.CreateDataResp(c, rs, nil, nil))
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
func (sv *Server) ClearCart(c *gin.Context) {
	authPayload, ok := c.MustGet(constants.AuthPayLoad).(*auth.TokenPayload)
	if !ok {
		c.JSON(http.StatusBadRequest, dto.CreateErr(InvalidBodyCode, errors.New("user not found")))
		return
	}

	cart, err := sv.repo.GetCart(c, repository.GetCartParams{
		UserID: utils.GetPgTypeUUID(authPayload.UserID),
	})
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.CreateErr(NotFoundCode, errors.New("cart not found")))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.CreateErr(InternalServerErrorCode, err))
		return
	}

	if string(cart.UserID.Bytes[:]) != authPayload.UserID.String() {
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
func (sv *Server) addCartRoutes(rg *gin.RouterGroup) {
	cart := rg.Group("/carts", authenticateMiddleware(sv.tokenGenerator))
	{
		cart.POST("", sv.CreateCart)
		cart.GET("", sv.GetCartHandler)
		cart.POST("checkout", sv.CheckoutHandler)
		cart.PUT("clear", sv.ClearCart)

		cart.GET("available-discounts", sv.GetCartAvailableDiscountsHandler)
		cartItems := cart.Group("items")
		{
			cartItems.PUT(":id/quantity", sv.UpdateCartItemQtyHandler)
			cartItems.DELETE(":id", sv.RemoveCartItem)
		}
	}
}
