package repository

import "zephyr-backend/internal/domain"

type UserRepository interface {
    CreateUser(
        username, email, passwordHash, birthDate, phoneNumber,
        firstName, lastName, gender, yandexID, oauthProvider string,
    ) error
    GetByEmail(email string) (*domain.User, error)
    GetByPhone(phone string) (*domain.User, error)
    SetPhoneVerified(phone string) error
    SetEmailVerified(email string) error
}
