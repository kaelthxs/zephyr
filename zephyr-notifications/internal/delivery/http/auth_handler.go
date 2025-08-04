package http

import "zephyr-auth/config"

type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
	cfg         *config.Config
}

func NewAuthHandler(authUC *usecase.AuthUseCase, cfg *config.Config) *AuthHandler {
	return &AuthHandler{authUseCase: authUC, cfg: cfg}
}
