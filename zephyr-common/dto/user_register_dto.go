package dto

import "time"

type UserRegisterRequestDTO struct {
	Username        string  `json:"username" binding:"required"`
	Email           string  `json:"email" binding:"required,email"`
	Password        string  `json:"password" binding:"required"`
	BirthDate       *string `json:"birth_date"`
	Gender          string  `json:"gender"`
	Pronouns        string  `json:"pronouns"`
	IsEmailVerified *bool   `json:"is_email_verified"` // <-- ADD THIS
}

func ParseBirthDate(birthDate *string) *time.Time {
	if birthDate == nil || *birthDate == "" {
		return nil
	}

	parsed, err := time.Parse("2006-01-02", *birthDate) // Assuming frontend sends "YYYY-MM-DD"
	if err != nil {
		// Handle invalid date format (could log or return nil silently)
		return nil
	}

	return &parsed
}
