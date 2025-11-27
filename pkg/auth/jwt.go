package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
)

const MinSecretKeySize = 32

type JwtGenerator struct {
	secretKey string
}

func NewJwtGenerator(secretKey string) (TokenGenerator, error) {
	if len(secretKey) < MinSecretKeySize {
		return nil, fmt.Errorf("secret key must be at least %d characters", MinSecretKeySize)
	}
	return &JwtGenerator{secretKey}, nil
}

func (g *JwtGenerator) GenerateToken(userID uuid.UUID, username string, role repository.UserRole, duration time.Duration) (string, *TokenPayload, error) {
	payload, err := NewPayload(userID, username, role, duration)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	token, err := jwtToken.SignedString([]byte(g.secretKey))
	return token, payload, err
}

func (g *JwtGenerator) VerifyToken(token string) (*TokenPayload, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}

		return []byte(g.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &TokenPayload{}, keyFunc)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	// .(*Payload) type assertion
	payload, ok := jwtToken.Claims.(*TokenPayload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
