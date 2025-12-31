package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID          uint   `gorm:"primarykey"`
	Username    string `gorm:"uniqueIndex;size:50;not null"`
	Email       string `gorm:"uniqueIndex;size:100;not null"`
	Password    string `gorm:"size:255;not null"` // 存储哈希后的密码
	Mobile      string `gorm:"size:20"`           // 手机号
	Name        string `gorm:"size:100"`
	Avatar      string `gorm:"size:255"`
	Status      int    `gorm:"default:1"`    // 状态：0=禁用，1=启用
	IsActivated bool   `gorm:"default:true"` // 是否激活
	LastLogin   *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"index"`
}

// ToResponse 转换为响应格式
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Mobile:    u.Mobile,
		Name:      u.Name,
		Avatar:    u.Avatar,
		Status:    u.Status,
		LastLogin: u.LastLogin,
		CreatedAt: u.CreatedAt,
	}
}

// IsActive 检查用户是否激活
func (u *User) IsActive() bool {
	return u.IsActivated
}

// UserResponse 用户响应格式
type UserResponse struct {
	ID        uint       `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Mobile    string     `json:"mobile"`
	Name      string     `json:"name"`
	Avatar    string     `json:"avatar"`
	Status    int        `json:"status"`
	LastLogin *time.Time `json:"last_login"`
	CreatedAt time.Time  `json:"created_at"`
}

// UserQuery 用户查询参数
type UserQuery struct {
	Page     int    `json:"page" form:"page"`
	Size     int    `json:"size" form:"size"`
	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
	Status   *int   `json:"status" form:"status"`
}

// GetOffset 获取偏移量
func (q *UserQuery) GetOffset() int {
	if q.Page <= 0 {
		q.Page = 1
	}
	return (q.Page - 1) * q.GetSize()
}

// GetSize 获取每页大小
func (q *UserQuery) GetSize() int {
	if q.Size <= 0 {
		q.Size = 10
	}
	if q.Size > 100 {
		q.Size = 100
	}
	return q.Size
}

// GetPage 获取页码
func (q *UserQuery) GetPage() int {
	if q.Page <= 0 {
		q.Page = 1
	}
	return q.Page
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required" validate:"required,min=6,max=50"`
	NewPassword string `json:"new_password" binding:"required" validate:"required,min=6,max=50"`
}

// UpdateProfileRequest 更新用户资料请求
type UpdateProfileRequest struct {
	Name   string `json:"name" validate:"omitempty,max=100"`
	Mobile string `json:"mobile" validate:"omitempty,len=11"`
	Avatar string `json:"avatar" validate:"omitempty,url"`
}

// UserCreateRequest 创建用户请求
type UserCreateRequest struct {
	Username string `json:"username" binding:"required" validate:"required,min=3,max=50"`
	Password string `json:"password" binding:"required" validate:"required,min=6,max=50"`
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Mobile   string `json:"mobile" validate:"omitempty,len=11"`
	Name     string `json:"name" validate:"omitempty,max=100"`
}

// UserProfileUpdateRequest 更新用户资料请求
type UserProfileUpdateRequest struct {
	Email  string `json:"email" validate:"omitempty,email"`
	Mobile string `json:"mobile" validate:"omitempty,len=11"`
	Name   string `json:"name" validate:"omitempty,max=100"`
	Avatar string `json:"avatar" validate:"omitempty,url"`
}
