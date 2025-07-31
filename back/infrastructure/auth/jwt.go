package auth

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
    "zephyr-backend/internal/repository"
)

// jwtService is a concrete implementation of repository.AuthService that
// generates and verifies JWT tokens. It also delegates password hashing
// and comparison to the bcrypt implementation defined in bcrypt.go.
type jwtService struct {
    secret string
}

// NewService constructs a new JWT service that can be used to hash
// passwords, check password hashes and generate JWT tokens. It returns
// the interface type to keep callers decoupled from the concrete type.
func NewService(secret string) repository.AuthService {
    return &jwtService{secret: secret}
}

func (s *jwtService) GenerateToken(userID string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.secret))
}
