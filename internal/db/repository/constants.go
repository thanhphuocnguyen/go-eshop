package repository

import (
	"time"

	"github.com/google/uuid"
)

const (
	ThumbnailRole = "thumbnail"
	SmallRole     = "small"
	GalleryRole   = "gallery"
	SliderRole    = "slider"
	IconRole      = "icon"
	LogoRole      = "logo"
	AvatarRole    = "avatar"
)

// User role constants for backward compatibility
type Role string

const (
	UserRoleCodeAdmin     Role = "admin"
	UserRoleCodeUser      Role = "user"
	UserRoleCodeModerator Role = "moderator"
)

// EntityType reference table struct
type EntityType struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	IsActive    bool      `json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Entity type constants for backward compatibility
const (
	EntityTypeCodeProduct       = "product"
	EntityTypeCodeVariant       = "variant"
	EntityTypeCodeProductRating = "product_rating"
	EntityTypeCodeCategory      = "category"
	EntityTypeCodeCollection    = "collection"
	EntityTypeCodeBrand         = "brand"
	EntityTypeCodeUser          = "user"
)

type NullUserRole struct {
	UserRole UserRole
	Valid    bool
}

type NullPaymentMethod struct {
	PaymentMethod PaymentMethod
	Valid         bool
}

type NullEntityType struct {
	EntityType EntityType
	Valid      bool
}

// VerifyEmail type for email verification
type VerifyEmail struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Email      string
	VerifyCode string
	IsUsed     bool
	CreatedAt  time.Time
	ExpiredAt  time.Time
}
