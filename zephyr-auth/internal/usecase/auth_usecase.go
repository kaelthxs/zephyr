package usecase

import (
	"github.com/google/uuid"
	"zephyr-auth/internal/domain/entity"
	"zephyr-auth/internal/domain/repository"
	"zephyr-auth/internal/domain/service"
	"zephyr-auth/internal/infrastructure/redis"
	"zephyr-common/dto"
)

type AuthUseCase struct {
	AuthRepository repository.AuthRepository
	HashService    service.HashService
	JWTService     service.JWTService
	RedisClient    *redis.RedisClient // Add this
}

func NewAuthUseCase(authRepository repository.AuthRepository, hashService service.HashService, jwtService service.JWTService, redisClient *redis.RedisClient) *AuthUseCase {
	return &AuthUseCase{
		AuthRepository: authRepository,
		HashService:    hashService,
		JWTService:     jwtService,
		RedisClient:    redisClient,
	}
}

func (a *AuthUseCase) CreateUser(user *entity.User) (*dto.UserMinimalDTO, error) {
	hashedPassword, err := a.HashService.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword
	user.ID = uuid.New()

	if user.IsEmailVerified == nil {
		defaultVerified := false
		user.IsEmailVerified = &defaultVerified
	}

	if user.Gender == "" {
		defaultGender := "не выбран"
		user.Gender = defaultGender
	}

	if user.Pronouns == "" {
		defaultPronouns := "не выбрано"
		user.Pronouns = defaultPronouns
	}

	if err := a.AuthRepository.CreateUser(user); err != nil {
		return nil, err // Handle DB insertion error
	}

	return &dto.UserMinimalDTO{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
