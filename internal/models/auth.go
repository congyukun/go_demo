package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

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

// Validate 验证注册请求
func (r *RegisterRequest) Validate() error {
	// 这里可以添加自定义验证逻辑
	// 目前使用结构体标签验证，所以返回nil
	return nil
}

// LoginResponse 登录响应结构体
type LoginResponse struct {
	Token            string       `json:"token"`
	RefreshToken     string       `json:"refresh_token"`
	ExpiresAt        time.Time    `json:"expires_at"`
	RefreshExpiresAt *time.Time   `json:"refresh_expires_at"`
	User             UserResponse `json:"user"`
}

// TokenClaims JWT token claims
type TokenClaims struct {
	UserID   int      `json:"user_id"`
	Username string   `json:"username"`
	Role     string   `json:"role"`    // 主角色，用于向后兼容
	Roles    []string `json:"roles"`   // 所有角色
	jwt.RegisteredClaims
}

// UpdateUserRequest 更新用户请求结构体
type UpdateUserRequest struct {
	Email  string `json:"email" validate:"omitempty,email" label:"邮箱"`
	Name   string `json:"name" validate:"omitempty,min=1,max=50" label:"姓名"`
	Status *int   `json:"status" validate:"omitempty,oneof=0 1" label:"状态"`
}

// PermissionResp 权限响应
type PermissionResp struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
}

// RoleRequest 角色请求
type RoleRequest struct {
	Code        string   `json:"code" binding:"required" validate:"required,min=2,max=50"`
	Name        string   `json:"name" binding:"required" validate:"required,min=2,max=100"`
	Description string   `json:"description" validate:"max=255"`
	Level       int      `json:"level" validate:"min=1"`
	Permissions []string `json:"permissions"` // 权限代码列表
}

// RoleResponse 角色响应
type RoleResponse struct {
	ID          uint              `json:"id"`
	Code        string            `json:"code"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Level       int               `json:"level"`
	Status      int               `json:"status"`
	Permissions []PermissionResp  `json:"permissions"`
	CreatedAt   time.Time         `json:"created_at"`
}
