package auth

import (
	"time"

	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
)

type TokenGenerator interface {
	GenerateToken(userID uuid.UUID, username string, role repository.UserRole, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
