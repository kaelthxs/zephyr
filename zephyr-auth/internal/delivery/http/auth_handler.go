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
	var req dto.UserLoginDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		return
	}

	tokens, err := h.AuthUseCase.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.SetCookie("refresh_token", tokens.RefreshToken, 60*60*24*30, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"accessToken":  tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
	})
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

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	tokens, err := h.AuthUseCase.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("refresh_token", tokens.RefreshToken, 60*60*24*30, "/", "localhost", false, true)

	c.JSON(http.StatusOK, tokens)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token missing"})
		return
	}

	err = h.AuthUseCase.Logout(refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "logout failed"})
		return
	}

	c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}
