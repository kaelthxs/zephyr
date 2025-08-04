package dto

import "github.com/google/uuid"

type UserMinimalDTO struct {
	ID       uuid.UUID `json:"id" binding:"required"`
	Username string    `json:"username" binding:"required"`
	Email    string    `json:"email" binding:"required,email"`
	Password string    `json:"password" binding:"required"`
}
