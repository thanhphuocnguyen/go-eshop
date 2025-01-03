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

type Payload struct {
	ID        uuid.UUID           `json:"id"`
	Username  string              `json:"username"`
	Role      repository.UserRole `json:"role"`
	UserID    int64               `json:"user_id"`
	IssuedAt  time.Time           `json:"issued_at"`
	ExpiredAt time.Time           `json:"expired_at"`
}
