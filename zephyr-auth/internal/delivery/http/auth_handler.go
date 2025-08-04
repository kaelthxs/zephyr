package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zephyr-auth/internal/mapper"
	"zephyr-auth/internal/usecase"
	"zephyr-common/dto"
)

type AuthHandler struct {
	AuthUseCase *usecase.AuthUseCase
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{AuthUseCase: authUseCase}
}

func (h *AuthHandler) Login(c *gin.Context) {

}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.UserRegisterRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userEntity := mapper.ToAuthUserFromRegisterDTO(&req) // Mapping DTO → Entity

	minimalUser, err := h.AuthUseCase.CreateUser(userEntity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, minimalUser)
}

func (h *AuthHandler) Logout(c *gin.Context) {

}
