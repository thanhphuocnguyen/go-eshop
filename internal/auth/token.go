package auth

import (
	"time"

	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
)

type TokenGenerator interface {
	GenerateToken(userID int64, username string, userRole repository.UserRole, duration time.Duration) (string, Payload, error)
	VerifyToken(token string) (*Payload, error)
}
