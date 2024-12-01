package auth

import (
	"time"

	"github.com/thanhphuocnguyen/go-eshop/internal/db/sqlc"
)

type TokenGenerator interface {
	GenerateToken(userID int64, username string, userRole sqlc.UserRole, duration time.Duration) (string, Payload, error)
	VerifyToken(token string) (*Payload, error)
}
