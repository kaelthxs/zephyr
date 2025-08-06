package dto

type UserRegisterRequestDTO struct {
	Username        string  `json:"username" binding:"required"`
	Email           string  `json:"email" binding:"required,email"`
	Password        string  `json:"password" binding:"required"`
	BirthDate       *string `json:"birth_date"`
	Gender          string  `json:"gender"`
	Pronouns        string  `json:"pronouns"`
	IsEmailVerified *bool   `json:"is_email_verified"`
}
