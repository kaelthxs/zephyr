package dto

import (
	"github.com/google/uuid"
	"time"
)

type UserDTO struct {
	ID                 uuid.UUID
	Username           string
	Email              string
	PasswordHash       string
	PhoneNumber        *string
	IsPhoneVerified    bool
	Role               string
	AccountStatus      string
	CreatedAt          time.Time
	UpdatedAt          time.Time
	FirstName          *string
	LastName           *string
	ProfilePicture     *string
	BirthDate          *time.Time
	Gender             *string
	Pronouns           string
	Country            *string
	City               *string
	Profession         *string
	LanguagePreference string
	Theme              string
	PremiumStatus      bool
	PremiumPlan        *string
	TwoFactorEnabled   bool
	TwoFactorMethod    *string
	EmbeddingVector    []float32
	OAuthID            *string
	OAuthProvider      *string
	IsEmailVerified    *bool
}
