package constants

type DiscountRule string

const (
	ProductRule          DiscountRule = "product"
	CategoryRule         DiscountRule = "category"
	CollectionRule       DiscountRule = "collection"
	PurchaseQuantityRule DiscountRule = "purchase_quantity"
	FirstTimeBuyerRule   DiscountRule = "first_time_buyer"
	BrandRule            DiscountRule = "brand"
	CustomerSegmentRule  DiscountRule = "customer_segment"
)
