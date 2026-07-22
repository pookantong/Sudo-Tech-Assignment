package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/idtoken"

	"cinema-booking-backend/internal/model"
)

var ErrInvalidGoogleToken = errors.New("invalid google id token")

type GoogleVerifier struct {
	clientID string
}

func NewGoogleVerifier(clientID string) *GoogleVerifier {
	return &GoogleVerifier{clientID: clientID}
}

type GoogleIdentity struct {
	Sub   string
	Email string
	Name  string
}

func (v *GoogleVerifier) Verify(ctx context.Context, rawIDToken string) (*GoogleIdentity, error) {
	payload, err := idtoken.Validate(ctx, rawIDToken, v.clientID)
	if err != nil {
		return nil, ErrInvalidGoogleToken
	}

	sub, _ := payload.Claims["sub"].(string)
	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)
	if sub == "" {
		return nil, ErrInvalidGoogleToken
	}

	return &GoogleIdentity{Sub: sub, Email: email, Name: name}, nil
}

type TokenIssuer struct {
	secret []byte
	ttl    time.Duration
}

func NewTokenIssuer(secret string, ttl time.Duration) *TokenIssuer {
	return &TokenIssuer{secret: []byte(secret), ttl: ttl}
}

func (t *TokenIssuer) Issue(userID string, role model.Role) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"iat":     now.Unix(),
		"exp":     now.Add(t.ttl).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(t.secret)
}