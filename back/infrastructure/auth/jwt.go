package auth

import (
    "time"
    "fmt"
    "github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
    secret string
}

func NewService(secret string) *jwtService {
    return &jwtService{secret: secret}
}

func (s *jwtService) GenerateToken(userID string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    
    signedToken, err := token.SignedString([]byte(s.secret))
    if err != nil {
        fmt.Println("Ошибка подписи токена:", err)
    } else {
        fmt.Println("Токен:", signedToken)
    }

    return signedToken, err
}
