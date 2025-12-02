package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/dto"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"github.com/thanhphuocnguyen/go-eshop/pkg/auth"
	"github.com/thanhphuocnguyen/go-eshop/pkg/payment"
)

// DiscountProcessor handles discount validation and calculation
type DiscountProcessor struct {
	repo repository.Repository
}

// DiscountContext contains all necessary data for discount processing
type DiscountContext struct {
	User       repository.GetUserDetailsByIDRow
	CartItems  []repository.GetCartItemsRow
	Discounts  []repository.GetDiscountByCodesRow
	TotalPrice float64
}

// ItemDiscount represents discount applied to a specific item
type ItemDiscount struct {
	ItemIndex      int
	DiscountAmount float64
	DiscountID     uuid.UUID
}

// DiscountResult contains the final discount calculation results
type DiscountResult struct {
	ItemDiscounts    []ItemDiscount
	TotalDiscount    float64
	AppliedDiscounts []uuid.UUID
}

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
	authPayload, ok := c.MustGet(AuthPayLoad).(*auth.TokenPayload)
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
	authPayload, ok := c.MustGet(AuthPayLoad).(*auth.TokenPayload)
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
	authPayload, _ := c.MustGet(AuthPayLoad).(*auth.TokenPayload)
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

	authPayload, ok := c.MustGet(AuthPayLoad).(*auth.TokenPayload)
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
				ID:        newCart.ID,
				UserID:    newCart.UserID,
				SessionID: newCart.SessionID,
				ItemCount: 0,
				CreatedAt: newCart.CreatedAt,
				UpdatedAt: newCart.UpdatedAt,
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
				ID:        uuid.New(),
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
	authPayload, ok := c.MustGet(AuthPayLoad).(*auth.TokenPayload)
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
	authPayload, ok := c.MustGet(AuthPayLoad).(*auth.TokenPayload)
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
	discountProcessor := &DiscountProcessor{repo: sv.repo}
	discountResult, err := discountProcessor.ProcessDiscounts(c, DiscountContext{
		User:      user,
		CartItems: itemRows,
		Discounts: nil, // will be populated if discount codes are provided
	}, req.DiscountCodes)
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
	authPayload, ok := c.MustGet(AuthPayLoad).(*auth.TokenPayload)
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

// ------------------------------ Discount Processing Methods ------------------------------

// ProcessDiscounts processes all discount codes and calculates the final discount amounts
func (dp *DiscountProcessor) ProcessDiscounts(c *gin.Context, ctx DiscountContext, discountCodes []string) (*DiscountResult, error) {
	result := &DiscountResult{
		ItemDiscounts:    []ItemDiscount{},
		TotalDiscount:    0,
		AppliedDiscounts: []uuid.UUID{},
	}

	if len(discountCodes) == 0 {
		return result, nil
	}

	// Get discount records
	discountRows, err := dp.repo.GetDiscountByCodes(c, discountCodes)
	if err != nil {
		return nil, fmt.Errorf("failed to get discount codes: %w", err)
	}

	// Validate discount applicability
	if err := dp.validateDiscountApplicability(c, discountRows); err != nil {
		return nil, err
	}

	// Process each discount
	for _, discountRow := range discountRows {
		if discountRow.DiscountType == repository.DiscountTypeFreeShipping {
			result.AppliedDiscounts = append(result.AppliedDiscounts, discountRow.ID)
			continue
		}

		itemDiscounts, err := dp.processDiscountForItems(ctx, discountRow)
		if err != nil {
			return nil, fmt.Errorf("failed to process discount %s: %w", discountRow.Code, err)
		}

		result.ItemDiscounts = append(result.ItemDiscounts, itemDiscounts...)
		result.AppliedDiscounts = append(result.AppliedDiscounts, discountRow.ID)
	}

	// Calculate total discount
	for _, itemDiscount := range result.ItemDiscounts {
		result.TotalDiscount += itemDiscount.DiscountAmount
	}

	return result, nil
}

// validateDiscountApplicability validates basic discount rules
func (dp *DiscountProcessor) validateDiscountApplicability(c *gin.Context, discounts []repository.GetDiscountByCodesRow) error {
	authPayload, _ := c.MustGet(AuthPayLoad).(*auth.TokenPayload)
	stackCnt := 0

	for _, discount := range discounts {
		// Check validity period
		if discount.ValidFrom.After(time.Now().UTC()) {
			return fmt.Errorf("discount code %s is not valid yet", discount.Code)
		}
		if discount.ValidUntil.Valid && discount.ValidUntil.Time.Before(time.Now().UTC()) {
			return fmt.Errorf("discount code %s has expired", discount.Code)
		}

		// Check stacking rules
		if discount.IsStackable {
			stackCnt++
			if stackCnt > 1 {
				return fmt.Errorf("only one stackable discount code is allowed")
			}
		}

		// Check usage limits
		if err := dp.validateUsageLimits(c, discount, authPayload.UserID); err != nil {
			return err
		}
	}

	return nil
}

// validateUsageLimits checks user and global usage limits
func (dp *DiscountProcessor) validateUsageLimits(c *gin.Context, discount repository.GetDiscountByCodesRow, userID uuid.UUID) error {
	if discount.UsagePerUser != nil {
		usageCount, err := dp.repo.CountDiscountUsageByDiscountAndUser(c, repository.CountDiscountUsageByDiscountAndUserParams{
			DiscountID: discount.ID,
			UserID:     userID,
		})
		if err != nil {
			return fmt.Errorf("failed to check user usage limit: %w", err)
		}
		if int32(usageCount) >= *discount.UsagePerUser {
			return fmt.Errorf("you have reached the maximum usage for discount code %s", discount.Code)
		}
	}

	if discount.UsageLimit != nil && discount.TimesUsed >= *discount.UsageLimit {
		return fmt.Errorf("discount code %s has reached its usage limit", discount.Code)
	}

	return nil
}

// processDiscountForItems applies discount rules to cart items and calculates discount amounts
func (dp *DiscountProcessor) processDiscountForItems(ctx DiscountContext, discount repository.GetDiscountByCodesRow) ([]ItemDiscount, error) {
	var discountRules []repository.DiscountRule
	if err := json.Unmarshal(discount.Rules, &discountRules); err != nil {
		return nil, fmt.Errorf("failed to unmarshal discount rules: %w", err)
	}

	var itemDiscounts []ItemDiscount
	discountVal, _ := discount.DiscountValue.Float64Value()

	for i, item := range ctx.CartItems {
		if dp.isDiscountApplicableToItem(item, ctx.User, discountRules) {
			discountAmount := dp.calculateItemDiscount(item, discount.DiscountType, discountVal.Float64)
			if discountAmount > 0 {
				itemDiscounts = append(itemDiscounts, ItemDiscount{
					ItemIndex:      i,
					DiscountAmount: discountAmount,
					DiscountID:     discount.ID,
				})
			}
		}
	}

	return itemDiscounts, nil
}

// isDiscountApplicableToItem checks if a discount is applicable to a specific cart item
func (dp *DiscountProcessor) isDiscountApplicableToItem(item repository.GetCartItemsRow, user repository.GetUserDetailsByIDRow, rules []repository.DiscountRule) bool {
	for _, rule := range rules {
		switch DiscountRule(rule.RuleType) {
		case ProductRule:
			if !dp.validateProductRule(item, rule.RuleValue) {
				return false
			}
		case CategoryRule:
			if !dp.validateCategoryRule(item, rule.RuleValue) {
				return false
			}
		case BrandRule:
			if !dp.validateBrandRule(item, rule.RuleValue) {
				return false
			}
		case PurchaseQuantityRule:
			if !dp.validateQuantityRule(item, rule.RuleValue) {
				return false
			}
		case FirstTimeBuyerRule:
			if !dp.validateFirstTimeBuyerRule(user, rule.RuleValue) {
				return false
			}
		case CustomerSegmentRule:
			if !dp.validateCustomerSegmentRule(user, rule.RuleValue) {
				return false
			}
		}
	}
	return true
}

// Rule validation methods
func (dp *DiscountProcessor) validateProductRule(item repository.GetCartItemsRow, ruleValueBytes json.RawMessage) bool {
	var ruleValue models.ProductRule
	if err := json.Unmarshal(ruleValueBytes, &ruleValue); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal ProductRule")
		return false
	}

	for _, pID := range ruleValue.ProductIDs {
		if pID == item.ProductID {
			return true
		}
	}
	return false
}

func (dp *DiscountProcessor) validateCategoryRule(item repository.GetCartItemsRow, ruleValueBytes json.RawMessage) bool {
	var ruleValue models.CategoryRule
	if err := json.Unmarshal(ruleValueBytes, &ruleValue); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal CategoryRule")
		return false
	}

	for _, cID := range ruleValue.CategoryIDs {
		for _, itemCId := range item.CategoryIds {
			if cID == itemCId {
				return true
			}
		}
	}
	return false
}

func (dp *DiscountProcessor) validateBrandRule(item repository.GetCartItemsRow, ruleValueBytes json.RawMessage) bool {
	var ruleValue models.BrandRule
	if err := json.Unmarshal(ruleValueBytes, &ruleValue); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal BrandRule")
		return false
	}

	if !item.ProductBrandID.Valid {
		return false
	}

	for _, bID := range ruleValue.BrandIDs {
		if bID == item.ProductBrandID.Bytes {
			return true
		}
	}
	return false
}

func (dp *DiscountProcessor) validateQuantityRule(item repository.GetCartItemsRow, ruleValueBytes json.RawMessage) bool {
	var ruleValue models.PurchaseQuantityRule
	if err := json.Unmarshal(ruleValueBytes, &ruleValue); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal PurchaseQuantityRule")
		return false
	}

	qty := int32(item.CartItem.Quantity)
	return qty >= int32(ruleValue.MinQuantity) && qty <= int32(ruleValue.MaxQuantity)
}

func (dp *DiscountProcessor) validateFirstTimeBuyerRule(user repository.GetUserDetailsByIDRow, ruleValueBytes json.RawMessage) bool {
	var ruleValue models.FirstTimeBuyerRule
	if err := json.Unmarshal(ruleValueBytes, &ruleValue); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal FirstTimeBuyerRule")
		return false
	}

	if !ruleValue.IsFirstTimeBuyer {
		return true // Rule doesn't apply
	}

	return user.TotalOrders == 0
}

func (dp *DiscountProcessor) validateCustomerSegmentRule(user repository.GetUserDetailsByIDRow, ruleValueBytes json.RawMessage) bool {
	var ruleValue models.CustomerSegmentRule
	if err := json.Unmarshal(ruleValueBytes, &ruleValue); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal CustomerSegmentRule")
		return false
	}

	// Check if customer is new customer
	if ruleValue.IsNewCustomer {
		if user.TotalOrders != 0 {
			return false
		}
	}

	// Check maximum previous orders constraint
	if ruleValue.MaxPreviousOrders != nil && *ruleValue.MaxPreviousOrders > 0 {
		if user.TotalOrders > int64(*ruleValue.MaxPreviousOrders) {
			return false
		}
	}

	// Check customer type if specified
	if ruleValue.CustomerType != nil {
		// Map user role or classification to customer type
		// This could be based on user role, order history, or other criteria
		customerType := dp.determineCustomerType(user)
		if customerType != *ruleValue.CustomerType {
			return false
		}
	}

	// Check minimum total spent constraint
	if ruleValue.MinTotalSpent != nil && *ruleValue.MinTotalSpent > 0 {
		// Note: This requires implementing a method to get user's total spending
		// For now, we'll use a placeholder that needs to be implemented
		totalSpent, _ := user.TotalSpent.Float64Value()
		userTotalSpent := totalSpent.Float64
		if userTotalSpent < *ruleValue.MinTotalSpent {
			return false
		}
	}

	return true
}

// determineCustomerType maps user data to customer type categories
func (dp *DiscountProcessor) determineCustomerType(user repository.GetUserDetailsByIDRow) string {
	// Implement customer type logic based on business rules
	// Examples of customer types: "new", "regular", "vip", "premium", "loyal", etc.

	// Primary classification based on order history
	if user.TotalOrders == 0 {
		return "new"
	} else if user.TotalOrders >= 1 && user.TotalOrders <= 3 {
		return "regular"
	} else if user.TotalOrders > 3 && user.TotalOrders <= 10 {
		return "frequent"
	} else if user.TotalOrders > 10 {
		return "loyal"
	}

	// Secondary classification based on user role
	// This can override the order-based classification for special users
	switch user.RoleCode {
	case "premium", "vip":
		return "premium"
	case "admin", "moderator":
		return "staff"
	default:
		// Fall back to order-based classification above
		if user.TotalOrders == 0 {
			return "new"
		} else if user.TotalOrders <= 3 {
			return "regular"
		} else if user.TotalOrders <= 10 {
			return "frequent"
		} else {
			return "loyal"
		}
	}
}

// calculateItemDiscount calculates the discount amount for a specific item
func (dp *DiscountProcessor) calculateItemDiscount(item repository.GetCartItemsRow, discountType repository.DiscountType, discountValue float64) float64 {
	itemPrice, _ := item.VariantPrice.Float64Value()
	lineTotal := float64(item.CartItem.Quantity) * itemPrice.Float64

	switch discountType {
	case repository.DiscountTypeFixedAmount:
		return discountValue
	case repository.DiscountTypePercentage:
		return lineTotal * (discountValue / 100)
	default:
		return 0
	}
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
