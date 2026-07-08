package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/solomonolatunji/vessel/internal/types"
)

// TokenService handles generation and validation of cryptographic JWT access tokens.
type TokenService struct {
	secretKey []byte
}

// NewTokenService initializes a TokenService using the vessel JWT secret key.
func NewTokenService() *TokenService {
	return &TokenService{
		secretKey: []byte("vessel-super-secret-jwt-signing-key-change-in-prod"),
	}
}

// GenerateToken creates an HMAC-SHA256 JWT access token for a validated user identity.
func (ts *TokenService) GenerateToken(user *types.User) (string, error) {
	if user == nil {
		return "", errors.New("user cannot be nil when generating token")
	}

	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(72 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
		"iss":   "vessel-auth",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(ts.secretKey)
}

// ValidateToken parses and verifies the HMAC signature and expiration timestamp of a JWT token string.
func (ts *TokenService) ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return ts.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims or signature")
	}

	return claims, nil
}
