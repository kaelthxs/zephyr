package entity

import (
	"github.com/google/uuid"
	"time"
)

type AuthUser struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username string    `gorm:"column:username"`
	Email    string    `gorm:"column:email"`
	Password string    `gorm:"column:password_hash"`

	BirthDate *time.Time `gorm:"column:birth_date"`

	FirstName *string `gorm:"column:first_name"`
	LastName  *string `gorm:"column:last_name"`
	Gender    string  `gorm:"column:gender"`
	Pronouns  string  `gorm:"column:pronouns"`

	OauthID       *string `gorm:"column:oauth_id"`
	OauthProvider *string `gorm:"column:oauth_provider"`

	IsEmailVerified *bool `gorm:"column:is_email_verified"`
}
