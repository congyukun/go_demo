package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        int            `json:"id" gorm:"primaryKey;autoIncrement"`
	Username  string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;size:100"`
	Name      string         `json:"name" gorm:"size:100;not null"`
	Password  string         `json:"-" gorm:"size:255;not null"`
	Status    int            `json:"status" gorm:"default:1;comment:状态 1:正常 0:禁用"`
	Mobile    string         `json:"mobile" gorm:"size:20;not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// UserResponse 用户响应结构体
type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Status   int    `json:"status"`
	Mobile   string `json:"mobile"`
}

// ToResponse 转换为响应结构体
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Name:     u.Name,
		Status:   u.Status,
		Mobile:   u.Mobile,
	}
}
