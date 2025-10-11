package repository

import (
	"go_demo/internal/models"

	"gorm.io/gorm"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// 基础CRUD操作
	Create(user *models.User) error
	GetByID(id int) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByMobile(mobile string) (*models.User, error)
	Update(user *models.User) error
	Delete(id int) error

	// 查询操作
	List(query *models.UserQuery) ([]models.User, int64, error)
	Count() (int64, error)

	// 状态操作
	UpdateStatus(id int, status int) error

	// 扩展查询方法
	SearchUsers(keyword string, limit int) ([]models.User, error)
	GetActiveUsers() ([]models.User, error)
	GetRecentUsers(limit int) ([]models.User, error)

	// 存在性检查
	ExistsByUsername(username string) (bool, error)
	ExistsByEmail(email string) (bool, error)
	ExistsByMobile(mobile string) (bool, error)

	// 批量操作
	BatchUpdateStatus(ids []int, status int) error
}

// userRepository 用户仓储实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create 创建用户
func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(id int) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByPhone 根据手机号获取用户
func (r *userRepository) GetByMobile(mobile string) (*models.User, error) {
	var user models.User
	err := r.db.Where("mobile = ?", mobile).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete 删除用户（软删除）
func (r *userRepository) Delete(id int) error {
	return r.db.Delete(&models.User{}, id).Error
}

// List 获取用户列表
func (r *userRepository) List(query *models.UserQuery) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// 构建查询条件
	db := r.db.Model(&models.User{})

	// 添加过滤条件
	if query.Username != "" {
		db = db.Where("username LIKE ?", "%"+query.Username+"%")
	}
	if query.Email != "" {
		db = db.Where("email LIKE ?", "%"+query.Email+"%")
	}
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := query.GetOffset()
	limit := query.GetSize()

	err := db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Count 获取用户总数
func (r *userRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Count(&count).Error
	return count, err
}

// UpdateStatus 更新用户状态
func (r *userRepository) UpdateStatus(id int, status int) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateLastLogin 更新最后登录时间
func (r *userRepository) UpdateLastLogin(id uint) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("last_login", gorm.Expr("NOW()")).Error
}

// ExistsByUsername 检查用户名是否存在
func (r *userRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// ExistsByEmail 检查邮箱是否存在
func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// ExistsByMobile 检查手机号是否存在
func (r *userRepository) ExistsByMobile(mobile string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("mobile = ?", mobile).Count(&count).Error
	return count > 0, err
}

// GetActiveUsers 获取活跃用户
func (r *userRepository) GetActiveUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Where("status = ?", 1).Find(&users).Error
	return users, err
}

// BatchUpdateStatus 批量更新用户状态
func (r *userRepository) BatchUpdateStatus(ids []int, status int) error {
	return r.db.Model(&models.User{}).Where("id IN ?", ids).Update("status", status).Error
}

// SearchUsers 搜索用户（模糊匹配用户名、邮箱）
func (r *userRepository) SearchUsers(keyword string, limit int) ([]models.User, error) {
	var users []models.User
	err := r.db.Where("username LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Limit(limit).Find(&users).Error
	return users, err
}

// GetRecentUsers 获取最近注册的用户
func (r *userRepository) GetRecentUsers(limit int) ([]models.User, error) {
	var users []models.User
	err := r.db.Order("created_at DESC").Limit(limit).Find(&users).Error
	return users, err
}
