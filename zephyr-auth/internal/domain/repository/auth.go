package repository

import (
	"zephyr-auth/internal/domain/entity"
)

type AuthRepository interface {
	CreateUser(user *entity.User) error
	GetUserByEmail(email string) (*entity.User, error)
}
