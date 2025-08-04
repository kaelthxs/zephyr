package auth

import (
	"gorm.io/gorm"
	"zephyr-auth/internal/domain/entity"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (repo *AuthRepository) CreateUser(user *entity.User) error {
	return repo.db.Create(&user).Error
}

func (repo *AuthRepository) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	return &user, repo.db.Where("email = ?", email).First(&user).Error
}
