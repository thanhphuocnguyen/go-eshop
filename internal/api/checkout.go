package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
	"github.com/thanhphuocnguyen/go-eshop/internal/processors"
	"github.com/thanhphuocnguyen/go-eshop/internal/utils"
	"github.com/thanhphuocnguyen/go-eshop/pkg/payment"
)

// @Summary Update product items in the cart
// @Schemes http
// @Description update product items in the cart
// @Tags carts
// @Accept json
// @Param input body models.CheckoutModel true "Checkout input"
// @Produce json
// @Success 200 {object} dto.ApiResponse[repository.CreatePaymentResult]
// @Failure 400 {object} dto.ErrorResp
// @Failure 404 {object} dto.ErrorResp
// @Failure 403 {object} dto.ErrorResp
// @Failure 401 {object} dto.ErrorResp
// @Failure 500 {object} dto.ErrorResp
// @Router /carts/checkout [post]
func (s *Server) checkout(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	c := r.Context()

	userID := claims["userId"].(uuid.UUID)
	if err != nil {
		RespondBadRequest(w, InvalidBodyCode, errors.New("user not found"))
		return
	}
	// verify request body
	var req models.CheckoutModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	// validate request
	if err := s.validator.Struct(req); err != nil {
		RespondBadRequest(w, InvalidBodyCode, err)
		return
	}

	user, err := s.repo.GetUserDetailsByID(c, userID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, errors.New("user not found"))
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	cart, err := s.repo.GetCart(c, repository.GetCartParams{
		UserID: utils.GetPgTypeUUID(userID),
	})

	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, InternalServerErrorCode, errors.New("cart not found"))
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}
	var shippingAddr repository.ShippingAddressSnapshot

	address, err := s.repo.GetDefaultAddress(c, userID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			RespondNotFound(w, NotFoundCode, errors.New("address not found"))
			return
		}
		RespondInternalServerError(w, InternalServerErrorCode, err)
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
			RespondInternalServerError(w, InternalServerErrorCode, err)
			return
		}
		if cartUserId != userID {
			RespondForbidden(w, PermissionDeniedCode, errors.New("you are not allowed to access this cart"))
			return
		}
	}

	itemRows, err := s.repo.GetCartItems(c, cart.ID)
	if err != nil {
		log.Error().Err(err).Msg("GetCartItems")
		return
	}

	// Process discounts
	discountResult, err := s.discountProcessor.ProcessDiscounts(c, processors.DiscountContext{User: user, CartItems: itemRows}, req.DiscountCodes)
	if err != nil {
		log.Error().Err(err).Msg("ProcessDiscounts")
		RespondBadRequest(w, InvalidBodyCode, err)
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
			RespondInternalServerError(w, InternalServerErrorCode, err)
			return
		}
	}

	params := repository.CheckoutCartTxArgs{
		CartID:          cart.ID,
		TotalPrice:      totalPrice,
		ShippingAddress: shippingAddr,
		UserID:          userID,
		CustomerInfo: repository.CustomerInfoTxArgs{
			FullName: user.FirstName + " " + user.LastName,
			Email:    user.Email,
			Phone:    user.PhoneNumber,
		},
		CreateOrderItemParams: createOrderItemParams,
		DiscountPrice:         discountResult.TotalDiscount,
		DiscountIDs:           discountResult.AppliedDiscounts,
		PaymentMethodID:       uuid.MustParse(req.PaymentMethodId),
	}

	params.CreatePaymentFn = func(ctx context.Context, orderID uuid.UUID, method string) (paymentIntentID string, clientSecretID *string, err error) {
		// create payment intent
		intent, err := s.paymentSrv.CreatePaymentIntent(ctx, method, payment.PaymentRequest{
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

	rs, err := s.repo.CheckoutCartTx(c, params)
	if err != nil {
		RespondInternalServerError(w, InternalServerErrorCode, err)
		return
	}

	RespondSuccess(w, r, rs)
}
