package auth

import (
	"log"
	"strconv"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/sqlc"
)

type PasetoTokenGenerator struct {
	symmetricKey paseto.V4SymmetricKey
	implicit     []byte
}

func NewPasetoTokenGenerator() TokenGenerator {
	return &PasetoTokenGenerator{
		symmetricKey: paseto.NewV4SymmetricKey(),
		implicit:     []byte("implicit claim"),
	}
}

func (g *PasetoTokenGenerator) GenerateToken(userID int64, username string, userRole sqlc.UserRole, duration time.Duration) (string, Payload, error) {
	payload := Payload{
		ID:        uuid.New(),
		Username:  username,
		Role:      userRole,
		UserId:    userID,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	token := paseto.NewToken()
	token.Set("username", payload.Username)
	token.Set("id", payload.ID)
	token.Set("role", payload.Role)
	token.Set("user_id", payload.UserId)
	token.SetExpiration(payload.ExpiredAt)
	token.SetNotBefore(payload.IssuedAt)
	token.SetIssuedAt(payload.IssuedAt)

	return token.V4Encrypt(g.symmetricKey, g.implicit), payload, nil
}

func (g *PasetoTokenGenerator) VerifyToken(token string) (*Payload, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())
	parsedToken, err := parser.ParseV4Local(g.symmetricKey, token, g.implicit)

	if err != nil {
		if paseto.RuleError.Is(paseto.RuleError{}, err) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, err := getPayloadFromParsedData(parsedToken)
	if err != nil {
		log.Println(err)
		return nil, ErrInvalidToken
	}
	return payload, nil
}

func getPayloadFromParsedData(t *paseto.Token) (*Payload, error) {
	username, err := t.GetString("username")
	if err != nil {
		return nil, err
	}
	id, err := (t.GetString("id"))
	if err != nil {
		return nil, err
	}
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	issuedAt, err := t.GetIssuedAt()
	if err != nil {
		return nil, err
	}
	expiredAt, err := t.GetExpiration()
	if err != nil {
		return nil, err
	}

	role, err := t.GetString("role")
	if err != nil {
		return nil, err
	}

	userIDStr, err := t.GetString("user_id")
	if err != nil {
		return nil, err
	}

	userId, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return nil, err
	}

	return &Payload{
		ID:        idUUID,
		Username:  username,
		UserId:    userId,
		Role:      sqlc.UserRole(role),
		IssuedAt:  issuedAt,
		ExpiredAt: expiredAt,
	}, nil
}
