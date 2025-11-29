package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
)

type UserDetail struct {
	ID                uuid.UUID       `json:"id"`
	RoleID            string          `json:"roleId"`
	RoleCode          string          `json:"roleCode"`
	Username          string          `json:"username"`
	FirstName         string          `json:"firstName"`
	LastName          string          `json:"lastName"`
	Email             string          `json:"email,omitempty"`
	Phone             string          `json:"phone,omitempty"`
	AvatarURL         *string         `json:"avatarUrl,omitempty"`
	AvatarID          *string         `json:"avatarId,omitempty"`
	Locked            bool            `json:"locked,omitempty"`
	VerifiedEmail     bool            `json:"verifiedEmail,omitempty"`
	VerifiedPhone     bool            `json:"verifiedPhone,omitempty"`
	PasswordChangedAt string          `json:"passwordChangedAt,omitempty"`
	Addresses         []AddressDetail `json:"addresses,omitempty"`
	CreatedAt         string          `json:"createdAt,omitempty"`
	LastLoginAt       string          `json:"lastLoginAt,omitempty"`
	UpdatedAt         string          `json:"updatedAt,omitempty"`
}

type LoginResponse struct {
	ID                    string    `json:"sessionId"`
	AccessToken           string    `json:"accessToken"`
	AccessTokenExpiresAt  time.Time `json:"accessTokenExpiresIn"`
	RefreshToken          string    `json:"refreshToken"`
	RefreshTokenExpiresAt time.Time `json:"refreshTokenExpiresAt"`
}

type RefreshToken struct {
	AccessToken          string    `json:"accessToken"`
	AccessTokenExpiresAt time.Time `json:"accessTokenExpiresAt"`
}

func MapToUserResponse(user repository.User, roleCode string) UserDetail {
	return UserDetail{
		ID:                user.ID,
		Addresses:         []AddressDetail{},
		Email:             user.Email,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		RoleID:            user.RoleID.String(),
		RoleCode:          roleCode,
		Phone:             user.PhoneNumber,
		Username:          user.Username,
		VerifiedEmail:     user.VerifiedEmail,
		VerifiedPhone:     user.VerifiedPhone,
		CreatedAt:         user.CreatedAt.String(),
		UpdatedAt:         user.UpdatedAt.String(),
		PasswordChangedAt: user.PasswordChangedAt.String(),
	}
}
