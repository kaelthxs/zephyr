package postgres

import (
    "gorm.io/gorm"
    "zephyr-backend/internal/domain"
    "zephyr-backend/internal/repository"
)

// userRepo is a Postgres implementation of the repository.UserRepository interface.
// It delegates all persistence operations to GORM. The struct and its methods
// live in the infrastructure layer to keep storage details away from the
// domain and use cases.
type userRepo struct {
    db *gorm.DB
}

// NewUserRepository instantiates a new Postgres backed implementation of
// repository.UserRepository. The returned repository is ready for use.
func NewUserRepository(db *gorm.DB) repository.UserRepository {
    return &userRepo{db: db}
}

// CreateUser persists a new user into the database. Fields that are not
// explicitly set are persisted as zero values. The email verification flag
// defaults to false on creation.
func (r *userRepo) CreateUser(
    username, email, passwordHash, birthDate, phoneNumber, firstName, lastName, gender, yandexID, oauthProvider string,
) error {
    return r.db.Create(&domain.User{
        Username:       username,
        Email:          email,
        Password:       passwordHash,
        BirthDate:      birthDate,
        PhoneNumber:    phoneNumber,
        FirstName:      firstName,
        LastName:       lastName,
        Gender:         gender,
        YandexID:       yandexID,
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