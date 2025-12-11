package processors

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/constants"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
	"github.com/thanhphuocnguyen/go-eshop/internal/models"
)

// DiscountProcessor handles discount validation and calculation
type DiscountProcessor struct {
	repo repository.Store
}

func NewDiscountProcessor(repo repository.Store) *DiscountProcessor {
	return &DiscountProcessor{
		repo: repo,
	}
}

// DiscountContext contains all necessary data for discount processing
type DiscountContext struct {
	User      repository.GetUserDetailsByIDRow
	CartItems []repository.GetCartItemsRow
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
	if err := dp.validateDiscountApplicability(c, ctx, discountRows); err != nil {
		return nil, err
	}

	// Process each discount
	for _, discountRow := range discountRows {
		if discountRow.DiscountType == repository.DiscountTypeFreeShipping {
			result.AppliedDiscounts = append(result.AppliedDiscounts, discountRow.ID)
			continue
		}

		itemDiscounts, err := dp.processDiscountForItems(c, ctx, discountRow)
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
func (dp *DiscountProcessor) validateDiscountApplicability(c *gin.Context, ctx DiscountContext, discounts []repository.Discount) error {
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
		if !discount.IsStackable {
			if stackCnt > 1 {
				return fmt.Errorf("only one stackable discount code is allowed")
			}
			stackCnt++
		}

		// Check usage limits
		if err := dp.validateUsageLimits(c, discount, ctx.User.ID); err != nil {
			return err
		}
	}

	return nil
}

// validateUsageLimits checks user and global usage limits
func (dp *DiscountProcessor) validateUsageLimits(c *gin.Context, discount repository.Discount, userID uuid.UUID) error {
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
func (dp *DiscountProcessor) processDiscountForItems(c *gin.Context, ctx DiscountContext, discount repository.Discount) ([]ItemDiscount, error) {
	ruleRows, err := dp.repo.GetDiscountRules(c, discount.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get discount rules: %w", err)
	}

	var itemDiscounts []ItemDiscount
	discountVal, _ := discount.DiscountValue.Float64Value()

	for i, item := range ctx.CartItems {
		if dp.isDiscountApplicableToItem(item, ctx.User, ruleRows) {
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
		switch constants.DiscountRule(rule.RuleType) {
		case constants.ProductRule:
			if !dp.validateProductRule(item, rule.RuleValue) {
				return false
			}
		case constants.CategoryRule:
			if !dp.validateCategoryRule(item, rule.RuleValue) {
				return false
			}
		case constants.BrandRule:
			if !dp.validateBrandRule(item, rule.RuleValue) {
				return false
			}
		case constants.PurchaseQuantityRule:
			if !dp.validateQuantityRule(item, rule.RuleValue) {
				return false
			}
		case constants.FirstTimeBuyerRule:
			if !dp.validateFirstTimeBuyerRule(user, rule.RuleValue) {
				return false
			}
		case constants.CustomerSegmentRule:
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
			if cID == itemCId.String() {
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
