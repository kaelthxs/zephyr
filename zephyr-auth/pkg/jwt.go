package pkg

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTServiceImpl implements the JWTService interface defined in
// internal/domain/service/auth.go. It encapsulates the secret used for
// signing and verifying tokens and provides methods to generate and parse
// JSON Web Tokens that carry a user identifier and an expiration time.
type JWTServiceImpl struct {
	Secret string
}

// GenerateToken creates a signed JWT containing the provided userID in its
// claims. The token expires 24 hours from the time of issuance. It
// returns the signed token string or an error if signing fails.
func (s *JWTServiceImpl) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.Secret))
}

// ParseToken validates a JWT and extracts the userID from its claims. If
// the token is invalid or the claims are missing, an error is returned.
func (s *JWTServiceImpl) ParseToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.Secret), nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}
	userIDVal, ok := claims["user_id"]
	if !ok {
		return "", errors.New("user_id not found in token")
	}
	userID, ok := userIDVal.(string)
	if !ok {
		return "", errors.New("user_id claim is not a string")
	}
	return userID, nil
}
