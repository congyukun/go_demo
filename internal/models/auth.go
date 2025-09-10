package models

import "time"

// LoginRequest 登录请求结构体
type LoginRequest struct {
	Username string `json:"username" binding:"required" validate:"required,min=3,max=20"`
	Password string `json:"password" binding:"required" validate:"required,min=6"`
}

// RegisterRequest 注册请求结构体
type RegisterRequest struct {
	Username string `json:"username" binding:"required" validate:"required,min=3,max=20"`
	Password string `json:"password" binding:"required" validate:"required,min=6"`
	Email    string `json:"email" binding:"required" validate:"required,email"`
	Name     string `json:"name" binding:"required" validate:"required,min=1,max=50"`
}

// LoginResponse 登录响应结构体
type LoginResponse struct {
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expires_at"`
	User      UserResponse `json:"user"`
}

// TokenClaims JWT token claims
type TokenClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
}

// UpdateUserRequest 更新用户请求结构体
type UpdateUserRequest struct {
	Email  string `json:"email" validate:"omitempty,email"`
	Name   string `json:"name" validate:"omitempty,min=1,max=50"`
	Status *int   `json:"status" validate:"omitempty,oneof=0 1"`
}
