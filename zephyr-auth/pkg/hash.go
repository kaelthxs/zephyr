package pkg

import "golang.org/x/crypto/bcrypt"

// HashServiceImpl implements the HashService interface defined in
// internal/domain/service/auth.go. It uses the bcrypt algorithm to hash
// and verify passwords. Bcrypt is chosen for its strength and built‑in
// salting mechanism, helping to protect stored passwords from brute-force
// and rainbow table attacks.
type HashServiceImpl struct{}

// HashPassword takes a plaintext password and returns its bcrypt hash. The
// cost parameter is left at the bcrypt package's default, which offers a
// good balance between security and performance. If hashing fails, an
// error is returned.
func (s *HashServiceImpl) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword compares a plaintext password against its bcrypt hash and
// returns true if they match. It returns false for any mismatch or
// underlying error.
func (s *HashServiceImpl) CheckPassword(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}
