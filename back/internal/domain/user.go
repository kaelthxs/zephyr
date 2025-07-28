package domain

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username string    `gorm:"column:username"`
	Email    string    `gorm:"column:email"`
	Password string    `gorm:"column:password_hash"`

	BirthDate   string `gorm:"column:birth_date"` // формат "YYYY-MM-DD"
	PhoneNumber string `gorm:"column:phone_number"`

	FirstName string `gorm:"column:first_name"`
	LastName  string `gorm:"column:last_name"`
	Gender    string `gorm:"column:gender"` // 'male' | 'female' | 'unspecified'

	YandexID      string `gorm:"column:yandex_id"`
	OauthProvider string `gorm:"column:oauth_provider"`

	IsPhoneVerified  bool `gorm:"column:is_phone_verified"`
	IsEmailVerified  bool `gorm:"column:is_email_verified"`

	// опционально:
	// CreatedAt time.Time `gorm:"column:created_at"`
	// UpdatedAt time.Time `gorm:"column:updated_at"`
}
