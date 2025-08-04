package service

type HashService interface {
	HashPassword(password string) (string, error)
	CheckPassword(password, hashed string) bool
}

type JWTService interface {
	GenerateToken(userID string) (string, error)
	ParseToken(token string) (string, error)
}
