package pkg

import "golang.org/x/crypto/bcrypt"

type HashServiceImpl struct{}

func (s *HashServiceImpl) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *HashServiceImpl) CheckPassword(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}
