package service

import (
	"crypto/md5"
	"fmt"
	"go_demo/internal/models"
	"go_demo/internal/repository"
	"go_demo/internal/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService 认证服务接口
type AuthService interface {
	Login(req models.LoginRequest) (*models.LoginResponse, error)
	Register(req models.RegisterRequest) (*models.UserResponse, error)
	ValidateToken(token string) (*models.TokenClaims, error)
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
func (s *authService) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	// 验证参数
	if req.Username == "" || req.Password == "" {
		return nil, fmt.Errorf("用户名或密码不能为空")
	}

	// 查找用户
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户名或密码错误")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 验证密码
	if !s.verifyPassword(req.Password, user.Password) {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, fmt.Errorf("用户已被禁用")
	}

	// 生成JWT token
	token, err := utils.GenerateAccessToken(int64(user.ID), user.Username, "user")
	if err != nil {
		return nil, fmt.Errorf("生成token失败: %w", err)
	}
	expiresAt := time.Now().Add(24 * time.Hour)

	response := &models.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      *user.ToResponse(),
	}

	return response, nil
}

// Register 用户注册
func (s *authService) Register(req models.RegisterRequest) (*models.UserResponse, error) {

	// 检查用户名是否已存在
	if _, err := s.userRepo.GetByUsername(req.Username); err == nil {
		return nil, fmt.Errorf("用户名已存在")
	} else if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("检查用户名失败: %w", err)
	}

	// 检查电话号是否已存在
	if _, err := s.userRepo.GetByMobile(req.Mobile); err == nil {
		return nil, fmt.Errorf("手机号已存在")
	} else if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("检查手机号失败: %w", err)
	}

	// 创建用户
	// 创建用户
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Name:     req.Name,
		Password:     s.hashPassword(req.Password),
		Status:   1,
		Mobile:   req.Mobile,
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	return user.ToResponse(), nil
}

// ValidateToken 验证JWT token
func (s *authService) ValidateToken(token string) (*models.TokenClaims, error) {
	if token == "" {
		return nil, fmt.Errorf("token不能为空")
	}

	// 使用JWT验证token
	jwtClaims, err := utils.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("token验证失败: %w", err)
	}
	
	// 转换为TokenClaims格式
	claims := &models.TokenClaims{
		UserID:   int(jwtClaims.UserID),
		Username: jwtClaims.Username,
	}
	if err != nil {
		return nil, fmt.Errorf("token验证失败: %w", err)
	}

	return claims, nil
}


// hashPassword 密码哈希 - 使用bcrypt替代MD5
func (s *authService) hashPassword(Password string) string {
	// 使用bcrypt进行密码哈希，成本为10
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(Password), 10)
	if err != nil {
		// 如果bcrypt失败，回退到MD5（不推荐用于生产环境）
		hash := md5.Sum([]byte(Password))
		return fmt.Sprintf("%x", hash)
	}
	return string(hashedBytes)
}

// verifyPassword 验证密码 - 支持bcrypt和MD5
func (s *authService) verifyPassword(Password, hashedPassword string) bool {
	// 首先尝试bcrypt验证
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(Password))
	if err == nil {
		return true
	}
	
	// 如果bcrypt失败，尝试MD5验证（向后兼容）
	hash := md5.Sum([]byte(Password))
	md5Hash := fmt.Sprintf("%x", hash)
	return md5Hash == hashedPassword
}
