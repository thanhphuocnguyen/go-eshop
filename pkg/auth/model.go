package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
)

var (
	ErrExpiredToken = errors.New("token is expired")
	ErrInvalidToken = errors.New("token is invalid")
)

type TokenPayload struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"userId"`
	Username string    `json:"username"`
	RoleID   uuid.UUID `json:"roleId"`
	RoleCode string    `json:"roleCode"`
	IssuedAt time.Time `json:"issuedAt"`
	Expires  time.Time `json:"expires"`
}

func NewPayload(userID uuid.UUID, username string, rol repository.UserRole, duration time.Duration) (*TokenPayload, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &TokenPayload{
		ID:       id,
		UserID:   userID,
		RoleCode: rol.Code,
		RoleID:   rol.ID,
		Username: username,
		IssuedAt: time.Now(),
		Expires:  time.Now().Add(duration),
	}

	return payload, nil
}
