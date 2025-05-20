package repository

const (
	ThumbnailRole = "thumbnail"
	SmallRole     = "small"
	GalleryRole   = "gallery"
	SliderRole    = "slider"
	IconRole      = "icon"
	LogoRole      = "logo"
	AvatarRole    = "avatar"
)

type DiscountType string

const (
	PercentageDiscount  DiscountType = "percentage"
	FixedAmountDiscount DiscountType = "fixed_amount"
)
