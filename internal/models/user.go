package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型 - 匹配现有数据库表结构
type User struct {
	ID        int        `json:"id" gorm:"primarykey;column:id;type:int(11);AUTO_INCREMENT"`
	Username  string     `json:"username" gorm:"uniqueIndex:idx_username;size:50;not null;column:username"`
	Email     string     `json:"email" gorm:"uniqueIndex:idx_email;size:100;column:email"`
	Name      string     `json:"name" gorm:"size:100;not null;column:name"`
	Password  string     `json:"-" gorm:"size:255;not null;column:password"`
	Status    int        `json:"status" gorm:"default:1;not null;column:status;comment:状态 1:正常 0:禁用"`
	Mobile    string     `json:"mobile" gorm:"uniqueIndex:idx_mobile;size:20;not null;column:mobile"`
	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"column:updated_at;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `json:"-" gorm:"index:idx_deleted_at;column:deleted_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// UserResponse 用户响应结构
type UserResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Mobile    string    `json:"mobile"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse 转换为响应结构
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Name:      u.Name,
		Mobile:    u.Mobile,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Users []UserResponse `json:"users"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
}

// UserCreateRequest 创建用户请求
type UserCreateRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50" label:"用户名"`
	Email    string `json:"email" validate:"required,email" label:"邮箱"`
	Password string `json:"password" validate:"required,min=6,max=50" label:"密码"`
	Name     string `json:"name" validate:"required,min=1,max=100" label:"姓名"`
	Mobile   string `json:"mobile" validate:"required,mobile" label:"手机号"`
}

// UserUpdateRequest 更新用户请求
type UserUpdateRequest struct {
	Email  string `json:"email" validate:"omitempty,email" label:"邮箱"`
	Name   string `json:"name" validate:"omitempty,min=1,max=100" label:"姓名"`
	Mobile string `json:"mobile" validate:"omitempty,mobile" label:"手机号"`
	Status *int   `json:"status" validate:"omitempty,oneof=0 1" label:"状态"`
}

// UserQuery 用户查询参数
type UserQuery struct {
	Page     int    `form:"page" validate:"omitempty,min=1" label:"页码"`
	Size     int    `form:"size" validate:"omitempty,min=1,max=100" label:"每页数量"`
	Username string `form:"username" validate:"omitempty,max=50" label:"用户名"`
	Email    string `form:"email" validate:"omitempty,email" label:"邮箱"`
	Status   *int   `form:"status" validate:"omitempty,oneof=0 1" label:"状态"`
}

// GetPage 获取页码，默认为1
func (q *UserQuery) GetPage() int {
	if q.Page <= 0 {
		return 1
	}
	return q.Page
}

// GetSize 获取每页数量，默认为10
func (q *UserQuery) GetSize() int {
	if q.Size <= 0 {
		return 10
	}
	if q.Size > 100 {
		return 100
	}
	return q.Size
}

// GetOffset 获取偏移量
func (q *UserQuery) GetOffset() int {
	return (q.GetPage() - 1) * q.GetSize()
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required" label:"原密码"`
	NewPassword string `json:"new_password" validate:"required,min=6,max=50" label:"新密码"`
}

// UserProfileUpdateRequest 用户资料更新请求
type UserProfileUpdateRequest struct {
	Email  string `json:"email" validate:"omitempty,email" label:"邮箱"`
	Name   string `json:"name" validate:"omitempty,min=1,max=100" label:"姓名"`
	Mobile string `json:"mobile" validate:"omitempty,mobile" label:"手机号"`
}

// IsActive 检查用户是否激活
func (u *User) IsActive() bool {
	return u.Status == 1
}

// BeforeCreate GORM钩子：创建前
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// 设置默认值
	if u.Status == 0 {
		u.Status = 1
	}
	return nil
}

// BeforeUpdate GORM钩子：更新前
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// 可以在这里添加更新前的逻辑
	return nil
}
