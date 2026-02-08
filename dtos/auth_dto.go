package dtos

import "github.com/cesarbmathec/medical-exams-backend/models"

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	RoleID   uint   `json:"role_id" binding:"required"`
}

type LoginResponse struct {
	User  models.UserResponse `json:"user"`
	Token string              `json:"token"`
}

type RegisterResponse struct {
	User  models.UserResponse `json:"user"`
	Token string              `json:"token,omitempty"`
}
