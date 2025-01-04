package auth

import (
	"time"
)

type TokenGenerator interface {
	GenerateToken(userID int64, username, email string, duration time.Duration) (string, Payload, error)
	VerifyToken(token string) (*Payload, error)
}
