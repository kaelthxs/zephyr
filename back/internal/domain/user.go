package domain

import "github.com/google/uuid"

type User struct {
    ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    Username string `gorm:"column:username"`
    Email    string `gorm:"column:email"`
    Password string `gorm:"column:password_hash"`
    BirthDate string `gorm:"column:birth_date"`
    PhoneNumber         string `gorm:"column:phone_number"`
	IsPhoneVerified     bool   `gorm:"column:is_phone_verified"`
    IsEmailVerified bool `gorm:"column:is_email_verified"`
}
