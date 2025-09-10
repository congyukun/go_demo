package models

import "time"

// LoginRequest 登录请求结构体
type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20" label:"用户名"`
	Password string `json:"password" validate:"required,min=6" label:"密码"`
}

// RegisterRequest 注册请求结构体
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20" label:"用户名"`
	Password string `json:"password" validate:"required,min=6" label:"密码"`
	Email    string `json:"email" validate:"omitempty,email" label:"邮箱"`
	Name     string `json:"name" validate:"required,min=1,max=50" label:"姓名"`
	Mobile   string `json:"mobile" validate:"required,mobile" label:"手机号"`
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
	Email  string `json:"email" validate:"omitempty,email" label:"邮箱"`
	Name   string `json:"name" validate:"omitempty,min=1,max=50" label:"姓名"`
	Status *int   `json:"status" validate:"omitempty,oneof=0 1" label:"状态"`
}
