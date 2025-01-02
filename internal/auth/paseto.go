package auth

import (
	"strconv"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/thanhphuocnguyen/go-eshop/internal/db/repository"
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

func (g *PasetoTokenGenerator) GenerateToken(userID int64, username string, userRole repository.UserRole, duration time.Duration) (string, Payload, error) {
	payload := Payload{
		ID:        uuid.New(),
		Username:  username,
		Role:      userRole,
		UserID:    userID,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	token := paseto.NewToken()
	token.SetString("username", payload.Username)
	token.SetString("id", payload.ID.String())
	token.SetString("role", string(payload.Role))
	token.SetString("user_id", strconv.Itoa(int(payload.UserID)))
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
		log.Error().Err(err).Msg("failed to parse token")
		if paseto.RuleError.Is(paseto.RuleError{}, err) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, err := getPayloadFromParsedData(parsedToken)
	if err != nil {
		log.Error().Err(err).Msg("failed to get payload from parsed data")
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
		UserID:    userId,
		Role:      repository.UserRole(role),
		IssuedAt:  issuedAt,
		ExpiredAt: expiredAt,
	}, nil
}
