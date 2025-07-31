package repository

// AuthService defines the methods required for password hashing,
// password verification and JWT token generation. Implementations
// of this interface live in the infrastructure layer and should not
// depend on any higher‑level packages.
type AuthService interface {
    // HashPassword should take a plain text password and return a hashed
    // representation. An error should be returned if the hashing fails.
    HashPassword(password string) (string, error)
    // CheckPassword should compare a plain text password with its hashed
    // representation and return true when they match.
    CheckPassword(password, hashed string) bool
    // GenerateToken should produce a signed JWT for the given userID.
    GenerateToken(userID string) (string, error)
}
