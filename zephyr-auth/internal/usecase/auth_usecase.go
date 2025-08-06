package usecase

import (
	"errors"
	"github.com/google/uuid"
	redisLib "github.com/redis/go-redis/v9"
	"log"
	"time"
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
		return nil, err
	}

	return &dto.UserMinimalDTO{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (a *AuthUseCase) Login(loginDTO dto.UserLoginDTO) (*dto.AuthResponse, error) {
	user, err := a.AuthRepository.GetUserByEmail(loginDTO.Email)
	if err != nil {
		return nil, errors.New("invalid email")
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	if !a.HashService.CheckPassword(loginDTO.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := a.JWTService.GenerateToken(user.ID.String())
	if err != nil {
		return nil, errors.New("invalid token")
	}

	refreshToken := uuid.NewString()

	err = a.RedisClient.Client.Set(a.RedisClient.Ctx, "refresh:"+refreshToken, user.ID.String(), time.Hour*24*30).Err()
	if err != nil {
		return nil, err
	}

	log.Printf("User %s logged in. AccessToken issued. RefreshToken: %s", user.ID.String(), refreshToken)

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *AuthUseCase) RefreshToken(providedRefreshToken string) (*dto.AuthResponse, error) {
	key := "refresh:" + providedRefreshToken
	userID, err := a.RedisClient.Client.Get(a.RedisClient.Ctx, key).Result()
	if errors.Is(redisLib.Nil, err) {
		return nil, errors.New("invalid refresh token")
	} else if err != nil {
		return nil, err
	}

	_ = a.RedisClient.Client.Del(a.RedisClient.Ctx, key).Err()

	accessToken, err := a.JWTService.GenerateToken(userID)
	if err != nil {
		return nil, err
	}

	newRefreshToken := uuid.NewString()
	err = a.RedisClient.Client.Set(a.RedisClient.Ctx, "refresh:"+newRefreshToken, userID, time.Hour*24*30).Err()
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil

}
