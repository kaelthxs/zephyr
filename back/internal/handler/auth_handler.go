package handler

import (
    "zephyr-backend/internal/usecase"
)

type AuthHandler struct {
    uc           *usecase.UserUseCase
    clientID     string
    clientSecret string
    redirectURI  string
}

func NewAuthHandler(uc *usecase.UserUseCase, clientID, clientSecret, redirectURI string) *AuthHandler {
    return &AuthHandler{
        uc:           uc,
        clientID:     clientID,
        clientSecret: clientSecret,
        redirectURI:  redirectURI,
    }
}
