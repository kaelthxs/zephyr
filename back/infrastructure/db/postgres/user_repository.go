package postgres

import (
    "gorm.io/gorm"
    "zephyr-backend/internal/domain"
    "zephyr-backend/internal/repository"
)

type userRepo struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
    return &userRepo{db: db}
}

func (r *userRepo) CreateUser(
    username, email, passwordHash, birthDate, phoneNumber, firstName, lastName, gender, oauthID, oauthProvider string, isVerifiedEmail bool) error {
    return r.db.Create(&domain.User{
        Username:       username,
        Email:          email,
        Password:       passwordHash,
        BirthDate:      birthDate,
        PhoneNumber:    phoneNumber,
        FirstName:      firstName,
        LastName:       lastName,
        Gender:         gender,
        OauthID:       oauthID,
        OauthProvider:  oauthProvider,
        IsEmailVerified: false,
    }).Error
}

func (r *userRepo) GetByEmail(email string) (*domain.User, error) {
    var user domain.User
    err := r.db.Where("email = ?", email).First(&user).Error
    return &user, err
}

func (r *userRepo) GetByPhone(phone string) (*domain.User, error) {
    var user domain.User
    err := r.db.Where("phone_number = ?", phone).First(&user).Error
    return &user, err
}

func (r *userRepo) SetPhoneVerified(phone string) error {
    return r.db.Model(&domain.User{}).
        Where("phone_number = ?", phone).
        Update("is_phone_verified", true).Error
}

func (r *userRepo) SetEmailVerified(email string) error {
    return r.db.Model(&domain.User{}).
        Where("email = ?", email).
        Update("is_email_verified", true).Error
}