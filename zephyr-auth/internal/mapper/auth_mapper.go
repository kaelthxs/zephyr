package mapper

import (
	"time"
	"zephyr-auth/internal/domain/entity"
	"zephyr-common/dto"
)

func ToAuthUserFromRegisterDTO(d *dto.UserRegisterRequestDTO) *entity.User {
	return &entity.User{
		Username:        d.Username,
		Email:           d.Email,
		Password:        d.Password,
		BirthDate:       parseBirthDate(d.BirthDate),
		Gender:          d.Gender,
		Pronouns:        d.Pronouns,
		IsEmailVerified: d.IsEmailVerified,
	}
}

func parseBirthDate(birthDate *string) *time.Time {
	if birthDate == nil || *birthDate == "" {
		return nil
	}
	parsed, err := time.Parse("2006-01-02", *birthDate)
	if err != nil {
		return nil
	}
	return &parsed
}
