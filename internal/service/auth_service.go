package service

import (
	"crypto/md5"
	"fmt"
	"go_demo/internal/models"
	"go_demo/internal/repository"
	"go_demo/internal/utils"
	"go_demo/pkg/errors"
	"go_demo/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService 认证服务接口
type AuthService interface {
	Login(c *gin.Context, req models.LoginRequest) (*models.LoginResponse, error)
	Register(c *gin.Context, req models.RegisterRequest) (*models.UserResponse, error)
	ValidateToken(token string) (*models.TokenClaims, error)
	RefreshToken(refreshToken string) (*models.LoginResponse, error)
	Logout(token string) error
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

	// 生成JWT token
	token, err := utils.GenerateAccessToken(int64(user.ID), user.Username)
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

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		logger.Error("注册失败：提交事务错误",
			logger.String("username", req.Username),
			logger.Err(err),
		)
		return nil, errors.NewInternalServerError("注册失败").WithCause(err)
	}

	logger.Info("用户注册成功",
		logger.String("username", req.Username),
		logger.Int64("user_id", int64(user.ID)),
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

	// 生成新的JWT token
	token, err := utils.GenerateAccessToken(int64(user.ID), user.Username)
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
