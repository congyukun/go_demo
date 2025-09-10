package service

import (
	"fmt"
	"go_demo/internal/models"
	"go_demo/internal/repository"

	"gorm.io/gorm"
)

// UserService 用户服务接口
type UserService interface {
	GetUsers(page, pageSize int) ([]*models.UserResponse, int64, error)
	GetUserByID(id int) (*models.UserResponse, error)
	UpdateUser(id int, req models.UpdateUserRequest) (*models.UserResponse, error)
	DeleteUser(id int) error
}

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// GetUsers 获取用户列表
func (s *userService) GetUsers(page, pageSize int) ([]*models.UserResponse, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	users, total, err := s.userRepo.List(offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("获取用户列表失败: %w", err)
	}

	// 转换为响应结构体
	responses := make([]*models.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return responses, total, nil
}

// GetUserByID 根据ID获取用户
func (s *userService) GetUserByID(id int) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}

	return user.ToResponse(), nil
}

// UpdateUser 更新用户
func (s *userService) UpdateUser(id int, req models.UpdateUserRequest) (*models.UserResponse, error) {
	// 获取用户
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}

	// 更新字段
	if req.Email != "" {
		// 检查邮箱是否已被其他用户使用
		existingUser, err := s.userRepo.GetByEmail(req.Email)
		if err == nil && existingUser.ID != id {
			return nil, fmt.Errorf("邮箱已被使用")
		} else if err != nil && err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("检查邮箱失败: %w", err)
		}
		user.Email = req.Email
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Status != nil {
		user.Status = *req.Status
	}

	// 保存更新
	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("更新用户失败: %w", err)
	}

	return user.ToResponse(), nil
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(id int) error {
	// 检查用户是否存在
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("用户不存在")
		}
		return fmt.Errorf("获取用户失败: %w", err)
	}

	// 删除用户
	if err := s.userRepo.Delete(id); err != nil {
		return fmt.Errorf("删除用户失败: %w", err)
	}

	return nil
}
