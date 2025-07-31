package repository

type AuthService interface {
    HashPassword(password string) (string, error)
    CheckPassword(password, hashed string) bool
    GenerateToken(userID string) (string, error)
}
