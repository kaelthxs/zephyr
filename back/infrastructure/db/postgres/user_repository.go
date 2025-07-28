package postgres

import (
	"gorm.io/gorm"
	"zephyr-backend/internal/domain"
)

type userRepo struct {
	db *gorm.DB
}
func NewUserRepository(db *gorm.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(
	username, email, password, birthDate, phoneNumber, firstName, lastName, gender, yandexID, oauthProvider string,
) error {
	return r.db.Create(&domain.User{
		Username:       username,
		Email:          email,
		Password:       password,
		BirthDate:      birthDate,
		PhoneNumber:    phoneNumber,
		FirstName:      firstName,
		LastName:       lastName,
		Gender:         gender,
		YandexID:       yandexID,
		OauthProvider:  oauthProvider,
		IsEmailVerified: true,
	}).Error
}



func (r *userRepo) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepo) SetPhoneVerificationCode(phone, code string) error {
	return r.db.Model(&domain.User{}).
		Where("phone_number = ?", phone).
		Update("phone_verification_code", code).Error
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