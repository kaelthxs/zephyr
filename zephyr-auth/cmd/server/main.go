package main

import (
	"log"
	"zephyr-auth/internal/infrastructure/redis"

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

	log.Printf("fsdlkdsfajjsdkfjoasfdoijadsfoij")
	log.Print(db)

	log.Print(redisClient)

	router := gin.Default()
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
