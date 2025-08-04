package main

import (
	"log"
	handler "zephyr-auth/internal/delivery/http"
	"zephyr-auth/internal/infrastructure/auth"
	"zephyr-auth/internal/infrastructure/redis"
	"zephyr-auth/internal/usecase"
	"zephyr-auth/pkg"

	"github.com/gin-gonic/gin"
	"zephyr-auth/config"
	"zephyr-auth/internal/infrastructure/postgres"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфига: %v", err)
	}

	db, err := postgres.NewPostgres(cfg.DBUrl)
	if err != nil {
		log.Fatalf("Error connecting to postgres: %v", err)
	}

	redisClient, err := redis.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("Error connecting to redis: %v", err)
	}

	authRepo := auth.NewAuthRepository(db)
	hashService := &pkg.HashServiceImpl{}
	jwtService := &pkg.JWTServiceImpl{}
	authUseCase := usecase.NewAuthUseCase(authRepo, hashService, jwtService, redisClient)

	router := gin.Default()

	handler.RegisterAuthRoutes(router, authUseCase)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
