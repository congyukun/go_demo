package models

import (
	"time"
)

// Permission 权限定义
type Permission struct {
	ID          uint      `gorm:"primarykey"`
	Code        string    `gorm:"uniqueIndex;size:50;not null"`        // 权限代码
	Name        string    `gorm:"size:100;not null"`                   // 权限名称
	Description string    `gorm:"size:255"`                            // 权限描述
	Resource    string    `gorm:"size:50;not null"`                    // 资源
	Action      string    `gorm:"size:50;not null"`                    // 操作
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Role 角色定义
type Role struct {
	ID          uint         `gorm:"primarykey"`
	Code        string       `gorm:"uniqueIndex;size:50;not null"`     // 角色代码
	Name        string       `gorm:"size:100;not null"`                // 角色名称
	Description string       `gorm:"size:255"`                         // 角色描述
	Level       int          `gorm:"default:1"`                        // 角色级别，数字越大权限越高
	Status      int          `gorm:"default:1"`                        // 状态：0=禁用，1=启用
	Permissions []Permission `gorm:"many2many:role_permissions;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// UserRole 用户角色关联
type UserRole struct {
	ID        uint      `gorm:"primarykey"`
	UserID    uint      `gorm:"not null"`
	RoleID    uint      `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Role      Role      `gorm:"foreignKey:RoleID"`                  // 角色信息
}

// RolePermission 角色权限关联
type RolePermission struct {
	ID           uint `gorm:"primarykey"`
	RoleID       uint `gorm:"not null"`
	PermissionID uint `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// User 用户模型
type User struct {
	ID           uint       `gorm:"primarykey"`
	Username     string     `gorm:"uniqueIndex;size:50;not null"`
	Email        string     `gorm:"uniqueIndex;size:100;not null"`
	Password     string     `gorm:"size:255;not null"`              // 存储哈希后的密码
	Phone        string     `gorm:"size:20"`
	Mobile       string     `gorm:"size:20"`                        // 手机号
	Name         string     `gorm:"size:100"`
	Avatar       string     `gorm:"size:255"`
	Status       int        `gorm:"default:1"`                      // 状态：0=禁用，1=启用
	IsActivated  bool       `gorm:"default:true"`                   // 是否激活
	LastLogin    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `gorm:"index"`
	Roles        []Role     `gorm:"many2many:user_roles;"`
	UserRoles    []UserRole `gorm:"foreignKey:UserID"`              // 用户角色关联
	Permissions  []Permission                                       // 用户权限
}

// ToResponse 转换为响应格式
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Phone:     u.Phone,
		Name:      u.Name,
		Avatar:    u.Avatar,
		Status:    u.Status,
		LastLogin: u.LastLogin,
		CreatedAt: u.CreatedAt,
		Roles:     u.GetRoleCodes(),
	}
}

// GetRoleCodes 获取用户角色代码列表
func (u *User) GetRoleCodes() []string {
	var roleCodes []string
	for _, role := range u.Roles {
		roleCodes = append(roleCodes, role.Code)
	}
	return roleCodes
}

// HasRole 检查用户是否拥有指定角色
func (u *User) HasRole(roleCode string) bool {
	for _, role := range u.Roles {
		if role.Code == roleCode {
			return true
		}
	}
	return false
}

// HasPermission 检查用户是否拥有指定权限
func (u *User) HasPermission(resource, action string) bool {
	for _, role := range u.Roles {
		for _, permission := range role.Permissions {
			if permission.Resource == resource && permission.Action == action {
				return true
			}
		}
	}
	return false
}

// GetMaxRoleLevel 获取用户的最大角色级别
func (u *User) GetMaxRoleLevel() int {
	maxLevel := 0
	for _, role := range u.Roles {
		if role.Level > maxLevel {
			maxLevel = role.Level
		}
	}
	return maxLevel
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
	Phone     string     `json:"phone"`
	Name      string     `json:"name"`
	Avatar    string     `json:"avatar"`
	Status    int        `json:"status"`
	LastLogin *time.Time `json:"last_login"`
	CreatedAt time.Time  `json:"created_at"`
	Roles     []string   `json:"roles"`
}

// UserQuery 用户查询参数
type UserQuery struct {
	Page     int   `json:"page" form:"page"`
	Size     int   `json:"size" form:"size"`
	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
	Status   *int  `json:"status" form:"status"`
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
	Phone  string `json:"phone" validate:"omitempty,len=11"`
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
	Email   string `json:"email" validate:"omitempty,email"`
	Mobile  string `json:"mobile" validate:"omitempty,len=11"`
	Name    string `json:"name" validate:"omitempty,max=100"`
	Phone   string `json:"phone" validate:"omitempty,len=11"`
	Avatar  string `json:"avatar" validate:"omitempty,url"`
}

// AssignRoleRequest 分配角色请求
type AssignRoleRequest struct {
	Roles []string `json:"roles" binding:"required" validate:"required,min=1,dive,required"` // 角色代码列表
}

// RevokeRoleRequest 撤销角色请求
type RevokeRoleRequest struct {
	Roles []string `json:"roles" binding:"required" validate:"required,min=1,dive,required"` // 角色代码列表
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Name        string   `json:"name" binding:"required" validate:"required,max=100"` // 角色名称
	Code        string   `json:"code" binding:"required" validate:"required,max=50"`  // 角色代码
	Description string   `json:"description" validate:"max=200"`                      // 角色描述
	Permissions []string `json:"permissions" validate:"dive,required"`                // 权限代码列表
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	Name        *string  `json:"name" validate:"omitempty,max=100"`  // 角色名称
	Description *string  `json:"description" validate:"omitempty,max=200"` // 角色描述
	Permissions []string `json:"permissions" validate:"dive,required"` // 权限代码列表
}
