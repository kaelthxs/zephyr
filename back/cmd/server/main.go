package main

import (
    "log"

    "github.com/gin-gonic/gin"
    "zephyr-backend/config"
    "zephyr-backend/infrastructure/auth"
    "zephyr-backend/infrastructure/cache"
    "zephyr-backend/infrastructure/db/postgres"
    "zephyr-backend/infrastructure/mail"
    "zephyr-backend/infrastructure/sms"
    handler "zephyr-backend/infrastructure/web/gin"
    "zephyr-backend/internal/usecase"
)

func main() {
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal("ошибка загрузки конфигов:", err)
    }

    db, err := postgres.NewPostgres(cfg)
    if err != nil {
        log.Fatal("не удалось подключиться к БД:", err)
    }

    cacheClient := cache.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
    userRepo := postgres.NewUserRepository(db)
    authService := auth.NewService(cfg.JWTSecret)
    smsClient := sms.NewSmsClient(cfg.SmsRuApiID)
    mailer := mail.NewSMTPMailer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPFrom)
    userUC := usecase.NewUserUseCase(
        userRepo,
        authService,
        cacheClient,
        smsClient,
        mailer,
        cfg.AppBaseURL,
    )    

    router := gin.Default()
    authMiddleware := auth.AuthMiddleware(cfg.JWTSecret)
    handler.RegisterRoutes(router, userUC, authMiddleware, cfg)

    err = router.Run(":" + cfg.Port)
    if err != nil {
        log.Fatal("сервер не запустился:", err)
    }
}