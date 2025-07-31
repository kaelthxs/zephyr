package auth

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
    "zephyr-backend/internal/repository"
)

type jwtService struct {
    secret string
}

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
