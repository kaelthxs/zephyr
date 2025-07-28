package repository

import "zephyr-backend/internal/domain"

type UserRepository interface {
    CreateUser(username, email, hashed, birth_date, phone_number string) error
    GetByEmail(email string) (*domain.User, error)

    SetPhoneVerificationCode(phone, code string) error
    GetByPhone(phone string) (*domain.User, error)
    SetPhoneVerified(phone string) error
    SetEmailVerified(email string) error
}