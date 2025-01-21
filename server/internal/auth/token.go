package auth

import (
	"time"

	"github.com/google/uuid"
)

type TokenGenerator interface {
	GenerateToken(userID uuid.UUID, username, email string, duration time.Duration) (string, Payload, error)
	VerifyToken(token string) (*Payload, error)
}
