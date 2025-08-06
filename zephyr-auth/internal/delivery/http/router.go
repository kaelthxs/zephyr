package handler

import (
	"github.com/gin-gonic/gin"
	"zephyr-auth/internal/usecase"
)

func RegisterAuthRoutes(router *gin.Engine, authUseCase *usecase.AuthUseCase) {
	handler := NewAuthHandler(authUseCase)
	router.POST("/auth/register", handler.Register)
	router.POST("/auth/login", handler.Login)
	router.POST("/auth/refresh", handler.RefreshToken)
	router.POST("/auth/logout", handler.Logout)
}
