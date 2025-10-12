package service

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"go_demo/internal/models"
	"go_demo/internal/repository"
	"go_demo/internal/utils"
	"go_demo/pkg/errors"
	"go_demo/pkg/logger"
)

// AuthService 认证服务接口
type AuthService interface {
	Login(c *gin.Context, req models.LoginRequest) (*models.LoginResponse, error)
	Register(c *gin.Context, req models.RegisterRequest) (*models.UserResponse, error)
	ValidateToken(token string) (*models.TokenClaims, error)
	RefreshToken(refreshToken string) (*models.LoginResponse, error)
	Logout(token string) error
	AssignRole(currentUserID, targetUserID int64, roles []string) error
	RevokeRole(currentUserID, targetUserID int64, roles []string) error
	GetUserRoles(userID int64) ([]models.Role, error)
	GetAllRoles() ([]models.Role, error)
	CreateRole(currentUserID int64, req models.CreateRoleRequest) (*models.Role, error)
	UpdateRole(currentUserID int64, roleID int, req models.UpdateRoleRequest) (*models.Role, error)
}

// authService 认证服务实现
type authService struct {
	userRepo repository.UserRepository
}

// NewAuthService 创建认证服务实例
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

// Login 用户登录
func (s *authService) Login(c *gin.Context, req models.LoginRequest) (*models.LoginResponse, error) {
	// 验证参数
	if req.Username == "" || req.Password == "" {
		return nil, errors.NewValidationError("用户名或密码不能为空")
	}

	// 查找用户
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Info("登录失败：用户不存在",
				logger.String("username", req.Username),
				logger.String("client_ip", utils.GetClientIP(c)),
			)
			return nil, errors.ErrInvalidCredentials
		}
		logger.Error("登录失败：查询用户错误",
			logger.String("username", req.Username),
			logger.Err(err),
		)
		return nil, errors.NewInternalServerError("查询用户失败").WithCause(err)
	}

	// 验证密码
	if !s.verifyPassword(req.Password, user.Password) {
		logger.Info("登录失败：密码错误",
			logger.String("username", req.Username),
			logger.Int64("user_id", int64(user.ID)),
			logger.String("client_ip", utils.GetClientIP(c)),
		)
		return nil, errors.ErrInvalidCredentials
	}

	// 检查用户状态
	if user.Status != 1 {
		logger.Info("登录失败：用户已被禁用",
			logger.String("username", req.Username),
			logger.Int64("user_id", int64(user.ID)),
			logger.Int("status", user.Status),
			logger.String("client_ip", utils.GetClientIP(c)),
		)
		return nil, errors.NewForbiddenError("用户已被禁用")
	}

	// 加载用户角色和权限
	if err := s.userRepo.LoadUserRoles(user); err != nil {
		logger.Error("登录失败：加载用户角色错误",
			logger.String("username", req.Username),
			logger.Int64("user_id", int64(user.ID)),
			logger.Err(err),
		)
		return nil, errors.NewInternalServerError("加载用户角色失败").WithCause(err)
	}

	// 生成JWT token
	// 包含用户角色信息
	roles := user.GetRoleCodes()
	token, err := utils.GenerateAccessTokenWithRoles(int64(user.ID), user.Username, roles)
	if err != nil {
		logger.Error("登录失败：生成token错误",
			logger.String("username", req.Username),
			logger.Int64("user_id", int64(user.ID)),
			logger.Err(err),
		)
		return nil, errors.NewInternalServerError("生成token失败").WithCause(err)
	}

	// 生成刷新token
	refreshToken, err := utils.GenerateRefreshToken(int64(user.ID))
	if err != nil {
		logger.Error("登录失败：生成刷新token错误",
			logger.String("username", req.Username),
			logger.Int64("user_id", int64(user.ID)),
			logger.Err(err),
		)
		return nil, errors.NewInternalServerError("生成刷新token失败").WithCause(err)
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour) // 7天

	response := &models.LoginResponse{
		Token:            token,
		RefreshToken:     refreshToken,
		ExpiresAt:        expiresAt,
		RefreshExpiresAt: &refreshExpiresAt,
		User:             *user.ToResponse(),
	}

	logger.Info("用户登录成功",
		logger.String("username", req.Username),
		logger.Int64("user_id", int64(user.ID)),
		logger.Strings("roles", roles),
		logger.String("client_ip", utils.GetClientIP(c)),
	)

	return response, nil
}

// Register 用户注册
func (s *authService) Register(c *gin.Context, req models.RegisterRequest) (*models.UserResponse, error) {
	// 验证参数
	if err := req.Validate(); err != nil {
		return nil, errors.NewValidationError(err.Error())
	}

	logger.Info("用户注册请求",
		logger.String("request_id", utils.GetRequestID(c)),
		logger.String("mobile", req.Mobile),
		logger.String("username", req.Username),
		logger.String("email", req.Email),
		logger.String("client_ip", utils.GetClientIP(c)),
	)

	// 检查用户名是否已存在
	if _, err := s.userRepo.GetByUsername(req.Username); err == nil {
		return nil, errors.NewConflictError("用户名已存在")
	} else if err != gorm.ErrRecordNotFound {
		logger.Error("注册失败：检查用户名错误",
			logger.String("username", req.Username),
			logger.Err(err),
		)
		return nil, errors.NewInternalServerError("检查用户名失败").WithCause(err)
	}

	// 检查电话号是否已存在
	if _, err := s.userRepo.GetByMobile(req.Mobile); err == nil {
		return nil, errors.NewConflictError("手机号已存在")
	} else if err != gorm.ErrRecordNotFound {
		logger.Error("注册失败：检查手机号错误",
			logger.String("mobile", req.Mobile),
			logger.Err(err),
		)
		return nil, errors.NewInternalServerError("检查手机号失败").WithCause(err)
	}

	// 创建用户
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Name:     req.Name,
		Password: s.hashPassword(req.Password),
		Status:   1,
		Mobile:   req.Mobile,
	}

	// 开始事务
	tx := s.userRepo.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建用户
	if err := s.userRepo.CreateWithTx(tx, user); err != nil {
		tx.Rollback()
		logger.Error("注册失败：创建用户错误",
			logger.String("username", req.Username),
			logger.Err(err),
		)
		return nil, errors.NewInternalServerError("创建用户失败").WithCause(err)
	}

	// 为新用户分配默认角色
	defaultRole := &models.Role{
		Code: "user",
	}
	role, err := s.userRepo.GetRoleByCode(defaultRole.Code)
	if err != nil {
		tx.Rollback()
		logger.Error("注册失败：获取默认角色错误",
			logger.String("username", req.Username),
			logger.String("role_code", defaultRole.Code),
			logger.Err(err),
		)
		return nil, errors.NewInternalServerError("获取默认角色失败").WithCause(err)
	}

	// 分配角色
	userRole := &models.UserRole{
		UserID: user.ID,
		RoleID: role.ID,
	}
	if err := s.userRepo.CreateUserRoleWithTx(tx, userRole); err != nil {
		tx.Rollback()
		logger.Error("注册失败：分配用户角色错误",
			logger.String("username", req.Username),
			logger.Int64("user_id", int64(user.ID)),
			logger.Int64("role_id", int64(role.ID)),
			logger.Err(err),
		)
		return nil, errors.NewInternalServerError("分配用户角色失败").WithCause(err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		logger.Error("注册失败：提交事务错误",
			logger.String("username", req.Username),
			logger.Err(err),
		)
		return nil, errors.NewInternalServerError("注册失败").WithCause(err)
	}

	// 重新加载用户角色
	if err := s.userRepo.LoadUserRoles(user); err != nil {
		logger.Error("注册成功但加载角色失败",
			logger.String("username", req.Username),
			logger.Int64("user_id", int64(user.ID)),
			logger.Err(err),
		)
		// 不返回错误，因为注册已成功
	}

	logger.Info("用户注册成功",
		logger.String("username", req.Username),
		logger.Int64("user_id", int64(user.ID)),
		logger.String("role_code", role.Code),
		logger.String("client_ip", utils.GetClientIP(c)),
	)

	return user.ToResponse(), nil
}

// ValidateToken 验证JWT token
func (s *authService) ValidateToken(token string) (*models.TokenClaims, error) {
	if token == "" {
		return nil, errors.NewValidationError("token不能为空")
	}

	// 使用JWT验证token
	jwtClaims, err := utils.ValidateToken(token)
	if err != nil {
		logger.Debug("token验证失败",
			logger.Err(err),
		)
		return nil, errors.ErrInvalidToken
	}

	// 转换为TokenClaims格式
	claims := &models.TokenClaims{
		UserID:   int(jwtClaims.UserID),
		Username: jwtClaims.Username,
		Role:     jwtClaims.Role,     // 主角色
		Roles:    jwtClaims.Roles,    // 所有角色
	}

	return claims, nil
}

// RefreshToken 刷新token
func (s *authService) RefreshToken(refreshToken string) (*models.LoginResponse, error) {
	if refreshToken == "" {
		return nil, errors.NewValidationError("刷新token不能为空")
	}

	// 验证刷新token
	jwtClaims, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		logger.Debug("刷新token验证失败",
			logger.Err(err),
		)
		return nil, errors.ErrInvalidToken
	}

	// 获取用户信息
	user, err := s.userRepo.GetByID(int(jwtClaims.UserID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrInvalidToken
		}
		logger.Error("刷新token失败：获取用户信息错误",
			logger.Int64("user_id", jwtClaims.UserID),
			logger.Err(err),
		)
		return nil, errors.NewInternalServerError("获取用户信息失败").WithCause(err)
	}

	// 检查用户状态
	if user.Status != 1 {
		logger.Info("刷新token失败：用户已被禁用",
			logger.Int64("user_id", int64(user.ID)),
			logger.String("username", user.Username),
			logger.Int("status", user.Status),
		)
		return nil, errors.NewForbiddenError("用户已被禁用")
	}

	// 加载用户角色和权限
	if err := s.userRepo.LoadUserRoles(user); err != nil {
		logger.Error("刷新token失败：加载用户角色错误",
			logger.String("username", user.Username),
			logger.Int64("user_id", int64(user.ID)),
			logger.Err(err),
		)
		return nil, errors.NewInternalServerError("加载用户角色失败").WithCause(err)
	}

	// 生成新的JWT token
	roles := user.GetRoleCodes()
	token, err := utils.GenerateAccessTokenWithRoles(int64(user.ID), user.Username, roles)
	if err != nil {
		logger.Error("刷新token失败：生成新token错误",
			logger.String("username", user.Username),
			logger.Int64("user_id", int64(user.ID)),
			logger.Err(err),
		)
		return nil, errors.NewInternalServerError("生成token失败").WithCause(err)
	}

	// 生成新的刷新token
	newRefreshToken, err := utils.GenerateRefreshToken(int64(user.ID))
	if err != nil {
		logger.Error("刷新token失败：生成新刷新token错误",
			logger.String("username", user.Username),
			logger.Int64("user_id", int64(user.ID)),
			logger.Err(err),
		)
		return nil, errors.NewInternalServerError("生成刷新token失败").WithCause(err)
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour) // 7天

	response := &models.LoginResponse{
		Token:            token,
		RefreshToken:     newRefreshToken,
		ExpiresAt:        expiresAt,
		RefreshExpiresAt: &refreshExpiresAt,
		User:             *user.ToResponse(),
	}

	logger.Info("token刷新成功",
		logger.String("username", user.Username),
		logger.Int64("user_id", int64(user.ID)),
		logger.Strings("roles", roles),
	)

	return response, nil
}

// Logout 用户登出
func (s *authService) Logout(token string) error {
	if token == "" {
		return errors.NewValidationError("token不能为空")
	}

	// 验证token并获取用户信息
	claims, err := utils.ValidateToken(token)
	if err != nil {
		logger.Debug("登出时token验证失败",
			logger.Err(err),
		)
		return errors.ErrInvalidToken
	}

	// 在实际应用中，可以将token加入黑名单
	// 这里只是记录日志
	logger.Info("用户登出",
		logger.String("username", claims.Username),
		logger.Int64("user_id", claims.UserID),
	)

	return nil
}

// RevokeRole 撤销用户角色
func (s *authService) RevokeRole(currentUserID, targetUserID int64, roles []string) error {
	// 验证参数
	if len(roles) == 0 {
		return errors.NewValidationError("角色列表不能为空")
	}

	// 获取当前用户信息
	currentUser, err := s.userRepo.GetByID(int(currentUserID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrInvalidCredentials
		}
		logger.Error("撤销角色失败：获取当前用户信息错误",
			logger.Int64("current_user_id", currentUserID),
			logger.Err(err),
		)
		return errors.NewInternalServerError("获取当前用户信息失败").WithCause(err)
	}

	// 检查当前用户是否有管理员权限
	if !currentUser.HasRole("admin") {
		logger.Warn("撤销角色失败：当前用户没有管理员权限",
			logger.Int64("current_user_id", currentUserID),
			logger.String("username", currentUser.Username),
			logger.Int64("target_user_id", targetUserID),
		)
		return errors.NewForbiddenError("没有权限撤销角色")
	}

	// 获取目标用户信息
	targetUser, err := s.userRepo.GetByID(int(targetUserID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError("目标用户不存在")
		}
		logger.Error("撤销角色失败：获取目标用户信息错误",
			logger.Int64("target_user_id", targetUserID),
			logger.Err(err),
		)
		return errors.NewInternalServerError("获取目标用户信息失败").WithCause(err)
	}

	// 开始事务
	tx := s.userRepo.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 撤销指定角色
	for _, roleCode := range roles {
		// 获取角色信息
		role, err := s.userRepo.GetRoleByCode(roleCode)
		if err != nil {
			tx.Rollback()
			logger.Error("撤销角色失败：获取角色信息错误",
				logger.Int64("target_user_id", targetUserID),
				logger.String("role_code", roleCode),
				logger.Err(err),
			)
			return errors.NewNotFoundError("角色不存在: " + roleCode).WithCause(err)
		}

		// 删除用户角色关联
		if err := s.userRepo.DeleteUserRoleWithTx(tx, int(targetUser.ID), int(role.ID)); err != nil {
			tx.Rollback()
			logger.Error("撤销角色失败：删除用户角色关联错误",
				logger.Int64("target_user_id", targetUserID),
				logger.String("role_code", roleCode),
				logger.Err(err),
			)
			return errors.NewInternalServerError("撤销角色失败").WithCause(err)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		logger.Error("撤销角色失败：提交事务错误",
			logger.Int64("target_user_id", targetUserID),
			logger.Err(err),
		)
		return errors.NewInternalServerError("撤销角色失败").WithCause(err)
	}

	logger.Info("撤销角色成功",
		logger.Int64("current_user_id", currentUserID),
		logger.String("current_username", currentUser.Username),
		logger.Int64("target_user_id", targetUserID),
		logger.String("target_username", targetUser.Username),
		logger.Strings("roles", roles),
	)

	return nil
}

// GetUserRoles 获取用户角色
func (s *authService) GetUserRoles(userID int64) ([]models.Role, error) {
	// 获取用户信息
	user, err := s.userRepo.GetByID(int(userID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("用户不存在")
		}
		logger.Error("获取用户角色失败：获取用户信息错误",
			logger.Int64("user_id", userID),
			logger.Err(err),
		)
		return nil, errors.NewInternalServerError("获取用户信息失败").WithCause(err)
	}

	// 获取用户角色
	roles, err := s.userRepo.GetUserRoles(int(user.ID))
	if err != nil {
		logger.Error("获取用户角色失败",
			logger.Int64("user_id", userID),
			logger.Err(err),
		)
		return nil, errors.NewInternalServerError("获取用户角色失败").WithCause(err)
	}

	return roles, nil
}

// GetAllRoles 获取所有角色
func (s *authService) GetAllRoles() ([]models.Role, error) {
	roles, err := s.userRepo.GetRoles()
	if err != nil {
		logger.Error("获取所有角色失败",
			logger.Err(err),
		)
		return nil, errors.NewInternalServerError("获取所有角色失败").WithCause(err)
	}

	return roles, nil
}

// CreateRole 创建角色
func (s *authService) CreateRole(currentUserID int64, req models.CreateRoleRequest) (*models.Role, error) {
	// 获取当前用户信息
	currentUser, err := s.userRepo.GetByID(int(currentUserID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrInvalidCredentials
		}
		logger.Error("创建角色失败：获取当前用户信息错误",
			logger.Int64("current_user_id", currentUserID),
			logger.Err(err),
		)
		return nil, errors.NewInternalServerError("获取当前用户信息失败").WithCause(err)
	}

	// 检查当前用户是否有管理员权限
	if !currentUser.HasRole("admin") {
		logger.Warn("创建角色失败：当前用户没有管理员权限",
			logger.Int64("current_user_id", currentUserID),
			logger.String("username", currentUser.Username),
		)
		return nil, errors.NewForbiddenError("没有权限创建角色")
	}

	// 检查角色代码是否已存在
	existingRole, err := s.userRepo.GetRoleByCode(req.Code)
	if err == nil && existingRole != nil {
		return nil, errors.NewValidationError("角色代码已存在")
	}

	// 开始事务
	tx := s.userRepo.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建角色
	role := &models.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
	}
	if err := s.userRepo.CreateRoleWithTx(tx, role); err != nil {
		tx.Rollback()
		logger.Error("创建角色失败",
			logger.String("role_name", req.Name),
			logger.String("role_code", req.Code),
			logger.Err(err),
		)
		return nil, errors.NewInternalServerError("创建角色失败").WithCause(err)
	}

	// 分配权限
	for _, permissionCode := range req.Permissions {
		// 获取权限信息
		permission, err := s.userRepo.GetPermissionByCode(permissionCode)
		if err != nil {
			tx.Rollback()
			logger.Error("创建角色失败：获取权限信息错误",
				logger.String("role_code", req.Code),
				logger.String("permission_code", permissionCode),
				logger.Err(err),
			)
			return nil, errors.NewNotFoundError("权限不存在: " + permissionCode).WithCause(err)
		}

		// 创建角色权限关联
		rolePermission := &models.RolePermission{
			RoleID:       role.ID,
			PermissionID: permission.ID,
		}
		if err := s.userRepo.CreateRolePermissionWithTx(tx, rolePermission); err != nil {
			tx.Rollback()
			logger.Error("创建角色失败：创建角色权限关联错误",
				logger.String("role_code", req.Code),
				logger.String("permission_code", permissionCode),
				logger.Err(err),
			)
			return nil, errors.NewInternalServerError("创建角色失败").WithCause(err)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		logger.Error("创建角色失败：提交事务错误",
			logger.String("role_code", req.Code),
			logger.Err(err),
		)
		return nil, errors.NewInternalServerError("创建角色失败").WithCause(err)
	}

	// 重新加载角色信息
	role, err = s.userRepo.GetRoleByCode(req.Code)
	if err != nil {
		logger.Error("创建角色失败：重新加载角色信息错误",
			logger.String("role_code", req.Code),
			logger.Err(err),
		)
		return nil, errors.NewInternalServerError("创建角色失败").WithCause(err)
	}

	logger.Info("创建角色成功",
		logger.Int64("current_user_id", currentUserID),
		logger.String("current_username", currentUser.Username),
		logger.String("role_name", req.Name),
		logger.String("role_code", req.Code),
	)

	return role, nil
}

// UpdateRole 更新角色
func (s *authService) UpdateRole(currentUserID int64, roleID int, req models.UpdateRoleRequest) (*models.Role, error) {
	// 获取当前用户信息
	currentUser, err := s.userRepo.GetByID(int(currentUserID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrInvalidCredentials
		}
		logger.Error("更新角色失败：获取当前用户信息错误",
			logger.Int64("current_user_id", currentUserID),
			logger.Err(err),
		)
		return nil, errors.NewInternalServerError("获取当前用户信息失败").WithCause(err)
	}

	// 检查当前用户是否有管理员权限
	if !currentUser.HasRole("admin") {
		logger.Warn("更新角色失败：当前用户没有管理员权限",
			logger.Int64("current_user_id", currentUserID),
			logger.String("username", currentUser.Username),
		)
		return nil, errors.NewForbiddenError("没有权限更新角色")
	}

	// 获取角色信息
	role, err := s.userRepo.GetRoleByID(roleID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("角色不存在")
		}
		logger.Error("更新角色失败：获取角色信息错误",
			logger.Int("role_id", roleID),
			logger.Err(err),
		)
		return nil, errors.NewInternalServerError("获取角色信息失败").WithCause(err)
	}

	// 开始事务
	tx := s.userRepo.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新角色信息
	if req.Name != nil {
		role.Name = *req.Name
	}
	if req.Description != nil {
		role.Description = *req.Description
	}
	if err := s.userRepo.UpdateRoleWithTx(tx, role); err != nil {
		tx.Rollback()
		logger.Error("更新角色失败",
			logger.Int("role_id", roleID),
			logger.Err(err),
		)
		return nil, errors.NewInternalServerError("更新角色失败").WithCause(err)
	}

	// 如果提供了权限列表，则更新权限
	if len(req.Permissions) > 0 {
		// 清除现有权限
		if err := s.userRepo.ClearRolePermissionsWithTx(tx, roleID); err != nil {
			tx.Rollback()
			logger.Error("更新角色失败：清除角色现有权限错误",
				logger.Int("role_id", roleID),
				logger.Err(err),
			)
			return nil, errors.NewInternalServerError("更新角色失败").WithCause(err)
		}

		// 分配新权限
		for _, permissionCode := range req.Permissions {
			// 获取权限信息
			permission, err := s.userRepo.GetPermissionByCode(permissionCode)
			if err != nil {
				tx.Rollback()
				logger.Error("更新角色失败：获取权限信息错误",
					logger.Int("role_id", roleID),
					logger.String("permission_code", permissionCode),
					logger.Err(err),
				)
				return nil, errors.NewNotFoundError("权限不存在: " + permissionCode).WithCause(err)
			}

			// 创建角色权限关联
			rolePermission := &models.RolePermission{
				RoleID:       role.ID,
				PermissionID: permission.ID,
			}
			if err := s.userRepo.CreateRolePermissionWithTx(tx, rolePermission); err != nil {
				tx.Rollback()
				logger.Error("更新角色失败：创建角色权限关联错误",
					logger.Int("role_id", roleID),
					logger.String("permission_code", permissionCode),
					logger.Err(err),
				)
				return nil, errors.NewInternalServerError("更新角色失败").WithCause(err)
			}
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		logger.Error("更新角色失败：提交事务错误",
			logger.Int("role_id", roleID),
			logger.Err(err),
		)
		return nil, errors.NewInternalServerError("更新角色失败").WithCause(err)
	}

	// 重新加载角色信息
	role, err = s.userRepo.GetRoleByID(roleID)
	if err != nil {
		logger.Error("更新角色失败：重新加载角色信息错误",
			logger.Int("role_id", roleID),
			logger.Err(err),
		)
		return nil, errors.NewInternalServerError("更新角色失败").WithCause(err)
	}

	logger.Info("更新角色成功",
		logger.Int64("current_user_id", currentUserID),
		logger.String("current_username", currentUser.Username),
		logger.Int("role_id", roleID),
		logger.String("role_name", role.Name),
	)

	return role, nil
}

// hashPassword 密码哈希 - 使用bcrypt替代MD5
func (s *authService) hashPassword(password string) string {
	// 使用bcrypt进行密码哈希，成本为10
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		// 如果bcrypt失败，回退到MD5（不推荐用于生产环境）
		logger.Error("密码哈希失败，回退到MD5",
			logger.Err(err),
		)
		hash := md5.Sum([]byte(password))
		return fmt.Sprintf("%x", hash)
	}
	return string(hashedBytes)
}

// verifyPassword 验证密码 - 支持bcrypt和MD5
func (s *authService) verifyPassword(password, hashedPassword string) bool {
	// 首先尝试bcrypt验证
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == nil {
		return true
	}

	// 如果bcrypt失败，尝试MD5验证（向后兼容）
	hash := md5.Sum([]byte(password))
	md5Hash := fmt.Sprintf("%x", hash)
	if md5Hash == hashedPassword {
		logger.Warn("使用MD5密码验证，建议升级到bcrypt",
			logger.String("password_hash", hashedPassword[:10]+"..."),
		)
		return true
	}

	return false
}

// AssignRole 分配角色给用户
func (s *authService) AssignRole(currentUserID, targetUserID int64, roles []string) error {
	// 验证参数
	if len(roles) == 0 {
		return errors.NewValidationError("角色列表不能为空")
	}

	// 获取当前用户信息
	currentUser, err := s.userRepo.GetByID(int(currentUserID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrInvalidCredentials
		}
		logger.Error("分配角色失败：获取当前用户信息错误",
			logger.Int64("current_user_id", currentUserID),
			logger.Err(err),
		)
		return errors.NewInternalServerError("获取当前用户信息失败").WithCause(err)
	}

	// 检查当前用户是否有管理员权限
	if !currentUser.HasRole("admin") {
		logger.Warn("分配角色失败：当前用户没有管理员权限",
			logger.Int64("current_user_id", currentUserID),
			logger.String("username", currentUser.Username),
			logger.Int64("target_user_id", targetUserID),
		)
		return errors.NewForbiddenError("没有权限分配角色")
	}

	// 获取目标用户信息
	targetUser, err := s.userRepo.GetByID(int(targetUserID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError("目标用户不存在")
		}
		logger.Error("分配角色失败：获取目标用户信息错误",
			logger.Int64("target_user_id", targetUserID),
			logger.Err(err),
		)
		return errors.NewInternalServerError("获取目标用户信息失败").WithCause(err)
	}

	// 开始事务
	tx := s.userRepo.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 清除用户现有角色
	if err := s.userRepo.ClearUserRolesWithTx(tx, int(targetUser.ID)); err != nil {
		tx.Rollback()
		logger.Error("分配角色失败：清除用户现有角色错误",
			logger.Int64("target_user_id", targetUserID),
			logger.Err(err),
		)
		return errors.NewInternalServerError("清除用户现有角色失败").WithCause(err)
	}

	// 分配新角色
	for _, roleCode := range roles {
		// 获取角色信息
		role, err := s.userRepo.GetRoleByCode(roleCode)
		if err != nil {
			tx.Rollback()
			logger.Error("分配角色失败：获取角色信息错误",
				logger.Int64("target_user_id", targetUserID),
				logger.String("role_code", roleCode),
				logger.Err(err),
			)
			return errors.NewNotFoundError("角色不存在: " + roleCode).WithCause(err)
		}

		// 创建用户角色关联
		userRole := &models.UserRole{
			UserID: targetUser.ID,
			RoleID: role.ID,
		}
		if err := s.userRepo.CreateUserRoleWithTx(tx, userRole); err != nil {
			tx.Rollback()
			logger.Error("分配角色失败：创建用户角色关联错误",
				logger.Int64("target_user_id", targetUserID),
				logger.String("role_code", roleCode),
				logger.Err(err),
			)
			return errors.NewInternalServerError("分配角色失败").WithCause(err)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		logger.Error("分配角色失败：提交事务错误",
			logger.Int64("target_user_id", targetUserID),
			logger.Err(err),
		)
		return errors.NewInternalServerError("分配角色失败").WithCause(err)
	}

	logger.Info("分配角色成功",
		logger.Int64("current_user_id", currentUserID),
		logger.String("current_username", currentUser.Username),
		logger.Int64("target_user_id", targetUserID),
		logger.String("target_username", targetUser.Username),
		logger.Strings("roles", roles),
	)

	return nil
}
