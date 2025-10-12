package repository

import (
	"go_demo/internal/models"

	"gorm.io/gorm"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// 基础CRUD操作
	Create(user *models.User) error
	CreateWithTx(tx *gorm.DB, user *models.User) error
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

	// 事务操作
	BeginTransaction() *gorm.DB

	// RBAC相关操作
	LoadUserRoles(user *models.User) error
	CreateUserRole(userRole *models.UserRole) error
	CreateUserRoleWithTx(tx *gorm.DB, userRole *models.UserRole) error
	DeleteUserRoles(userID int) error
	DeleteUserRoleWithTx(tx *gorm.DB, userID, roleID int) error
	ClearUserRolesWithTx(tx *gorm.DB, userID int) error
	GetRoleByCode(code string) (*models.Role, error)
	GetRoleByID(id int) (*models.Role, error)
	GetPermissionByCode(code string) (*models.Permission, error)
	CreateRoleWithTx(tx *gorm.DB, role *models.Role) error
	UpdateRoleWithTx(tx *gorm.DB, role *models.Role) error
	ClearRolePermissionsWithTx(tx *gorm.DB, roleID int) error
	CreateRolePermissionWithTx(tx *gorm.DB, rolePermission *models.RolePermission) error
	GetRolesByIDs(roleIDs []int) ([]models.Role, error)
	GetPermissionsByIDs(permissionIDs []int) ([]models.Permission, error)
	GetUserRoles(userID int) ([]models.Role, error)
	GetUserPermissions(userID int) ([]models.Permission, error)
	AssignRoleToUser(userID, roleID int) error
	RemoveRoleFromUser(userID, roleID int) error
	AssignPermissionToRole(roleID, permissionID int) error
	RemovePermissionFromRole(roleID, permissionID int) error
	CreateRole(role *models.Role) error
	UpdateRole(role *models.Role) error
	DeleteRole(id int) error
	GetRoles() ([]models.Role, error)
	GetRolePermissions(roleID int) ([]models.Permission, error)
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

// CreateWithTx 在事务中创建用户
func (r *userRepository) CreateWithTx(tx *gorm.DB, user *models.User) error {
	return tx.Create(user).Error
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

// BeginTransaction 开始事务
func (r *userRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

// RBAC相关操作

// LoadUserRoles 加载用户角色和权限
func (r *userRepository) LoadUserRoles(user *models.User) error {
	// 加载用户角色
	err := r.db.Preload("Role.Permissions").Where("user_id = ?", user.ID).Find(&user.UserRoles).Error
	if err != nil {
		return err
	}

	// 提取角色
	roles := make([]models.Role, 0, len(user.UserRoles))
	for _, ur := range user.UserRoles {
		roles = append(roles, ur.Role)
	}
	user.Roles = roles

	// 提取权限
	permissionMap := make(map[int]models.Permission)
	for _, ur := range user.UserRoles {
		for _, p := range ur.Role.Permissions {
			permissionMap[int(p.ID)] = p
		}
	}

	permissions := make([]models.Permission, 0, len(permissionMap))
	for _, p := range permissionMap {
		permissions = append(permissions, p)
	}
	user.Permissions = permissions

	return nil
}

// CreateUserRole 创建用户角色关联
func (r *userRepository) CreateUserRole(userRole *models.UserRole) error {
	return r.db.Create(userRole).Error
}

// CreateUserRoleWithTx 在事务中创建用户角色关联
func (r *userRepository) CreateUserRoleWithTx(tx *gorm.DB, userRole *models.UserRole) error {
	return tx.Create(userRole).Error
}

// DeleteUserRoles 删除用户所有角色
func (r *userRepository) DeleteUserRoles(userID int) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.UserRole{}).Error
}

// ClearUserRolesWithTx 在事务中删除用户所有角色
func (r *userRepository) ClearUserRolesWithTx(tx *gorm.DB, userID int) error {
	return tx.Where("user_id = ?", userID).Delete(&models.UserRole{}).Error
}

// DeleteUserRoleWithTx 在事务中删除用户特定角色
func (r *userRepository) DeleteUserRoleWithTx(tx *gorm.DB, userID, roleID int) error {
	return tx.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&models.UserRole{}).Error
}

// GetRoleByCode 根据代码获取角色
func (r *userRepository) GetRoleByCode(code string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("code = ?", code).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetRoleByID 根据ID获取角色
func (r *userRepository) GetRoleByID(id int) (*models.Role, error) {
	var role models.Role
	err := r.db.First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetRolesByIDs 根据ID列表获取角色
func (r *userRepository) GetRolesByIDs(roleIDs []int) ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Where("id IN ?", roleIDs).Find(&roles).Error
	return roles, err
}

// GetPermissionsByIDs 根据ID列表获取权限
func (r *userRepository) GetPermissionsByIDs(permissionIDs []int) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Where("id IN ?", permissionIDs).Find(&permissions).Error
	return permissions, err
}

// GetUserRoles 获取用户角色
func (r *userRepository) GetUserRoles(userID int) ([]models.Role, error) {
	var userRoles []models.UserRole
	err := r.db.Preload("Role").Where("user_id = ?", userID).Find(&userRoles).Error
	if err != nil {
		return nil, err
	}
	
	// 转换为角色列表
	roles := make([]models.Role, 0, len(userRoles))
	for _, ur := range userRoles {
		roles = append(roles, ur.Role)
	}
	
	return roles, nil
}

// GetUserPermissions 获取用户权限
func (r *userRepository) GetUserPermissions(userID int) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Table("permissions").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Distinct("permissions.*").
		Find(&permissions).Error
	return permissions, err
}

// AssignRoleToUser 为用户分配角色
func (r *userRepository) AssignRoleToUser(userID, roleID int) error {
	// 检查是否已存在该角色
	var count int64
	err := r.db.Model(&models.UserRole{}).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // 已存在，无需重复分配
	}

	userRole := &models.UserRole{
		UserID: uint(userID),
		RoleID: uint(roleID),
	}
	return r.db.Create(userRole).Error
}

// RemoveRoleFromUser 移除用户角色
func (r *userRepository) RemoveRoleFromUser(userID, roleID int) error {
	return r.db.Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&models.UserRole{}).Error
}

// AssignPermissionToRole 为角色分配权限
func (r *userRepository) AssignPermissionToRole(roleID, permissionID int) error {
	// 检查是否已存在该权限
	var count int64
	err := r.db.Model(&models.RolePermission{}).
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // 已存在，无需重复分配
	}

	rolePermission := &models.RolePermission{
		RoleID:       uint(roleID),
		PermissionID: uint(permissionID),
	}
	return r.db.Create(rolePermission).Error
}

// RemovePermissionFromRole 移除角色权限
func (r *userRepository) RemovePermissionFromRole(roleID, permissionID int) error {
	return r.db.Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Delete(&models.RolePermission{}).Error
}

// CreateRole 创建角色
func (r *userRepository) CreateRole(role *models.Role) error {
	return r.db.Create(role).Error
}

// CreateRoleWithTx 在事务中创建角色
func (r *userRepository) CreateRoleWithTx(tx *gorm.DB, role *models.Role) error {
	return tx.Create(role).Error
}

// UpdateRole 更新角色
func (r *userRepository) UpdateRole(role *models.Role) error {
	return r.db.Save(role).Error
}

// UpdateRoleWithTx 在事务中更新角色
func (r *userRepository) UpdateRoleWithTx(tx *gorm.DB, role *models.Role) error {
	return tx.Save(role).Error
}

// DeleteRole 删除角色
func (r *userRepository) DeleteRole(id int) error {
	// 先删除角色权限关联
	if err := r.db.Where("role_id = ?", id).Delete(&models.RolePermission{}).Error; err != nil {
		return err
	}

	// 再删除用户角色关联
	if err := r.db.Where("role_id = ?", id).Delete(&models.UserRole{}).Error; err != nil {
		return err
	}

	// 最后删除角色
	return r.db.Delete(&models.Role{}, id).Error
}

// GetRoles 获取所有角色
func (r *userRepository) GetRoles() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Find(&roles).Error
	return roles, err
}

// GetRolePermissions 获取角色权限
func (r *userRepository) GetRolePermissions(roleID int) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Table("permissions").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	return permissions, err
}

// ClearRolePermissionsWithTx 在事务中清除角色所有权限
func (r *userRepository) ClearRolePermissionsWithTx(tx *gorm.DB, roleID int) error {
	return tx.Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error
}

// CreateRolePermissionWithTx 在事务中创建角色权限关联
func (r *userRepository) CreateRolePermissionWithTx(tx *gorm.DB, rolePermission *models.RolePermission) error {
	return tx.Create(rolePermission).Error
}

// GetPermissionByCode 根据代码获取权限
func (r *userRepository) GetPermissionByCode(code string) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.Where("code = ?", code).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}
