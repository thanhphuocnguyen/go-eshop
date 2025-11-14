package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
)

var (
	ErrExpiredToken = errors.New("token is expired")
	ErrInvalidToken = errors.New("token is invalid")
)

type Payload struct {
	ID       uuid.UUID       `json:"id"`
	UserID   uuid.UUID       `json:"userId"`
	Username string          `json:"username"`
	RoleID   uuid.UUID       `json:"roleId"`
	RoleCode repository.Role `json:"roleCode"`
	IssuedAt time.Time       `json:"issuedAt"`
	Expires  time.Time       `json:"expires"`
}

func (payload *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(payload.Expires), nil
}
func (payload *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(payload.IssuedAt), nil
}
func (payload *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(payload.IssuedAt), nil
}
func (payload *Payload) GetIssuer() (string, error) {
	return "simple-life-shop", nil
}
func (payload *Payload) GetSubject() (string, error) {
	return payload.Username, nil
}
func (payload *Payload) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{payload.Username}, nil
}

func NewPayload(userID uuid.UUID, username string, rol repository.UserRole, duration time.Duration) (*Payload, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID:       id,
		UserID:   userID,
		RoleCode: repository.Role(rol.Code),
		RoleID:   rol.ID,
		Username: username,
		IssuedAt: time.Now(),
		Expires:  time.Now().Add(duration),
	}

	return payload, nil
}
