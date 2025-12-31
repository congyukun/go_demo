package service

import (
	"fmt"
	"go_demo/internal/models"
	"go_demo/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService 用户服务接口
type UserService interface {
	// 基础CRUD
	GetUsers(page, pageSize int) ([]*models.UserResponse, int64, error)
	GetUserByID(id int) (*models.UserResponse, error)
	UpdateUser(id int, req models.UpdateUserRequest) (*models.UserResponse, error)
	DeleteUser(id int) error

	// 用户管理
	CreateUser(req models.UserCreateRequest) (*models.UserResponse, error)
	UpdateUserProfile(id int, req models.UserProfileUpdateRequest) (*models.UserResponse, error)
	ChangePassword(id int, req models.ChangePasswordRequest) error
	UpdateUserStatus(id int, status int) error

	// 查询方法
	SearchUsers(keyword string, limit int) ([]*models.UserResponse, error)
	GetActiveUsers() ([]*models.UserResponse, error)
	GetRecentUsers(limit int) ([]*models.UserResponse, error)
	GetUserlist(page, pageSize int) ([]*models.UserResponse, int64, error)

	// 统计方法
	GetUserCount() (int64, error)
	GetUserStats() (map[string]interface{}, error)
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

func (s *userService) GetUserlist(page, pageSize int) ([]*models.UserResponse, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	query := &models.UserQuery{
		Page: page,
		Size: pageSize,
	}

	users, total, err := s.userRepo.GetUserList(query)

	if err != nil {
		return nil, 0, fmt.Errorf("获取用户列表失败: %w", err)
	}

	res := make([]*models.UserResponse, len(users))

	for k, user := range users {
		res[k] = user.ToResponse()
	}

	return res, total, nil

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

	query := &models.UserQuery{
		Page: page,
		Size: pageSize,
	}
	users, total, err := s.userRepo.List(query)
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
		if err == nil && int(existingUser.ID) != id {
			return nil, fmt.Errorf("邮箱已被使用")
		} else if err != nil && err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("检查邮箱失败: %w", err)
		}
		user.Email = req.Email
	}

	// 更新姓名字段
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

// CreateUser 创建用户
func (s *userService) CreateUser(req models.UserCreateRequest) (*models.UserResponse, error) {
	// 检查用户名是否已存在
	if _, err := s.userRepo.GetByUsername(req.Username); err == nil {
		return nil, fmt.Errorf("用户名已存在")
	} else if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("检查用户名失败: %w", err)
	}

	// 检查邮箱是否已存在
	if _, err := s.userRepo.GetByEmail(req.Email); err == nil {
		return nil, fmt.Errorf("邮箱已存在")
	} else if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("检查邮箱失败: %w", err)
	}

	// 检查手机号是否已存在（如果提供）
	if req.Mobile != "" {
		if _, err := s.userRepo.GetByMobile(req.Mobile); err == nil {
			return nil, fmt.Errorf("手机号已存在")
		} else if err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("检查手机号失败: %w", err)
		}
	}

	// 哈希密码
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, fmt.Errorf("密码哈希失败: %w", err)
	}

	// 创建用户
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Name:     req.Name,
		Password: string(hashedBytes),
		Mobile:   req.Mobile,
		Status:   1,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	return user.ToResponse(), nil
}

// UpdateUserProfile 更新用户资料
func (s *userService) UpdateUserProfile(id int, req models.UserProfileUpdateRequest) (*models.UserResponse, error) {
	// 获取用户
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}

	// 更新字段
	if req.Email != "" && req.Email != user.Email {
		// 检查邮箱是否已被其他用户使用
		existingUser, err := s.userRepo.GetByEmail(req.Email)
		if err == nil && int(existingUser.ID) != id {
			return nil, fmt.Errorf("邮箱已被使用")
		} else if err != nil && err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("检查邮箱失败: %w", err)
		}
		user.Email = req.Email
	}

	if req.Mobile != "" && req.Mobile != user.Mobile {
		// 检查手机号是否已被其他用户使用
		existingUser, err := s.userRepo.GetByMobile(req.Mobile)
		if err == nil && int(existingUser.ID) != id {
			return nil, fmt.Errorf("手机号已被使用")
		} else if err != nil && err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("检查手机号失败: %w", err)
		}
		user.Mobile = req.Mobile
	}

	// 更新姓名
	if req.Name != "" {
		user.Name = req.Name
	}

	// 更新手机号
	if req.Mobile != "" {
		user.Mobile = req.Mobile
	}

	// 更新头像
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	// 保存更新
	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("更新用户失败: %w", err)
	}

	return user.ToResponse(), nil
}

// ChangePassword 修改密码
func (s *userService) ChangePassword(id int, req models.ChangePasswordRequest) error {
	// 获取用户
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("用户不存在")
		}
		return fmt.Errorf("获取用户失败: %w", err)
	}

	// 验证原密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
	if err != nil {
		return fmt.Errorf("原密码错误")
	}

	// 哈希新密码
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 10)
	if err != nil {
		return fmt.Errorf("密码哈希失败: %w", err)
	}

	// 更新密码
	user.Password = string(hashedBytes)
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	return nil
}

// UpdateUserStatus 更新用户状态
func (s *userService) UpdateUserStatus(id int, status int) error {
	// 检查用户是否存在
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("用户不存在")
		}
		return fmt.Errorf("获取用户失败: %w", err)
	}

	// 更新状态
	if err := s.userRepo.UpdateStatus(id, status); err != nil {
		return fmt.Errorf("更新用户状态失败: %w", err)
	}

	return nil
}

// SearchUsers 搜索用户
func (s *userService) SearchUsers(keyword string, limit int) ([]*models.UserResponse, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	users, err := s.userRepo.SearchUsers(keyword, limit)
	if err != nil {
		return nil, fmt.Errorf("搜索用户失败: %w", err)
	}

	responses := make([]*models.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return responses, nil
}

// GetActiveUsers 获取活跃用户
func (s *userService) GetActiveUsers() ([]*models.UserResponse, error) {
	users, err := s.userRepo.GetActiveUsers()
	if err != nil {
		return nil, fmt.Errorf("获取活跃用户失败: %w", err)
	}

	responses := make([]*models.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return responses, nil
}

// GetRecentUsers 获取最近注册的用户
func (s *userService) GetRecentUsers(limit int) ([]*models.UserResponse, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	users, err := s.userRepo.GetRecentUsers(limit)
	if err != nil {
		return nil, fmt.Errorf("获取最近用户失败: %w", err)
	}

	responses := make([]*models.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return responses, nil
}

// GetUserCount 获取用户总数
func (s *userService) GetUserCount() (int64, error) {
	count, err := s.userRepo.Count()
	if err != nil {
		return 0, fmt.Errorf("获取用户总数失败: %w", err)
	}
	return count, nil
}

// GetUserStats 获取用户统计信息
func (s *userService) GetUserStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 总用户数
	totalCount, err := s.userRepo.Count()
	if err != nil {
		return nil, fmt.Errorf("获取总用户数失败: %w", err)
	}
	stats["total_users"] = totalCount

	// 活跃用户数
	activeUsers, err := s.userRepo.GetActiveUsers()
	if err != nil {
		return nil, fmt.Errorf("获取活跃用户失败: %w", err)
	}
	stats["active_users"] = len(activeUsers)

	// 由于数据库表结构中没有角色字段，暂时设置为0
	stats["admin_users"] = 0
	stats["normal_users"] = 0

	return stats, nil
}

// GetUserRoles 获取用户角色
/*
func (s *userService) GetUserRoles(userID int64) ([]*models.Role, error) {
	user, err := s.userRepo.GetByID(int(userID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}

	// 转换 []models.Role 为 []*models.Role
	roles := make([]*models.Role, len(user.Roles))
	for i := range user.Roles {
		roles[i] = &user.Roles[i]
	}

	return roles, nil
}
*/

// AssignRole 为用户分配角色
/*
func (s *userService) AssignRole(userID int64, roleName string) error {
	// 检查用户是否存在
	user, err := s.userRepo.GetByID(int(userID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("用户不存在")
		}
		return fmt.Errorf("获取用户失败: %w", err)
	}

	// 检查角色是否存在
	// 注意：这里需要实现获取角色的方法，目前假设角色存在
	// 在实际应用中，应该从角色服务或仓库中获取角色信息

	// 检查用户是否已有该角色
	for _, role := range user.Roles {
		if role.Code == roleName {
			return fmt.Errorf("用户已拥有该角色")
		}
	}

	// 添加角色
	// 注意：这里需要实现添加角色的方法，目前只是示例
	// 在实际应用中，应该通过用户仓库或角色仓库来处理关联关系

	return nil
}
*/

// RevokeRole 撤销用户角色
/*
func (s *userService) RevokeRole(userID int64, roleName string) error {
	// 检查用户是否存在
	user, err := s.userRepo.GetByID(int(userID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("用户不存在")
		}
		return fmt.Errorf("获取用户失败: %w", err)
	}

	// 检查用户是否拥有该角色
	hasRole := false
	for _, role := range user.Roles {
		if role.Code == roleName {
			hasRole = true
			break
		}
	}

	if !hasRole {
		return fmt.Errorf("用户不拥有该角色")
	}

	// 移除角色
	// 注意：这里需要实现移除角色的方法，目前只是示例
	// 在实际应用中，应该通过用户仓库或角色仓库来处理关联关系

	return nil
}
*/
