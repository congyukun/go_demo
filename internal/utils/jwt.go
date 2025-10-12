package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTConfig JWT配置
type JWTConfig struct {
	SecretKey     string `mapstructure:"secret_key" yaml:"secret_key"`
	AccessExpire  int64  `mapstructure:"access_expire" yaml:"access_expire"`   // 访问token过期时间（秒）
	RefreshExpire int64  `mapstructure:"refresh_expire" yaml:"refresh_expire"` // 刷新token过期时间（秒）
	Issuer        string `mapstructure:"issuer" yaml:"issuer"`                 // 签发者
}

// Claims JWT声明
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWTManager JWT管理器
type JWTManager struct {
	config JWTConfig
}

// NewJWTManager 创建JWT管理器
func NewJWTManager(config JWTConfig) *JWTManager {
	if config.Issuer == "" {
		config.Issuer = "go_demo"
	}
	if config.AccessExpire == 0 {
		config.AccessExpire = 3600 // 默认1小时
	}
	if config.RefreshExpire == 0 {
		config.RefreshExpire = 604800 // 默认7天
	}
	return &JWTManager{
		config: config,
	}
}

// GenerateAccessToken 生成访问token
func (j *JWTManager) GenerateAccessToken(userID int64, username string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.config.Issuer,
			Subject:   username,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(j.config.AccessExpire) * time.Second)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.config.SecretKey))
}

// GenerateRefreshToken 生成刷新token
func (j *JWTManager) GenerateRefreshToken(userID int64) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: "",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.config.Issuer,
			Subject:   "refresh",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(j.config.RefreshExpire) * time.Second)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.config.SecretKey))
}

// GenerateTokenPair 生成token对（访问token和刷新token）
func (j *JWTManager) GenerateTokenPair(userID int64, username string) (accessToken, refreshToken string, err error) {
	accessToken, err = j.GenerateAccessToken(userID, username)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = j.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ParseToken 解析token
func (j *JWTManager) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名方法")
		}
		return []byte(j.config.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的token")
}

// ValidateToken 验证token有效性
func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	// 检查token是否过期
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token已过期")
	}

	// 检查token是否还未生效
	if claims.NotBefore != nil && claims.NotBefore.Time.After(time.Now()) {
		return nil, errors.New("token还未生效")
	}

	return claims, nil
}

// ValidateRefreshToken 验证刷新token有效性
func (j *JWTManager) ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// 验证这是一个刷新token（刷新token的subject为"refresh"）
	if claims.Subject != "refresh" {
		return nil, errors.New("无效的刷新token")
	}

	return claims, nil
}

// GetUserIDFromToken 从token中获取用户ID
func (j *JWTManager) GetUserIDFromToken(tokenString string) (int64, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// GetUsernameFromToken 从token中获取用户名
func (j *JWTManager) GetUsernameFromToken(tokenString string) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Username, nil
}

// IsTokenExpired 检查token是否过期
func (j *JWTManager) IsTokenExpired(tokenString string) bool {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return true
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return true
	}

	return false
}

// GetTokenExpireTime 获取token过期时间
func (j *JWTManager) GetTokenExpireTime(tokenString string) (time.Time, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return time.Time{}, err
	}

	if claims.ExpiresAt != nil {
		return claims.ExpiresAt.Time, nil
	}

	return time.Time{}, errors.New("token没有过期时间")
}

// GetTokenRemainingTime 获取token剩余有效时间
func (j *JWTManager) GetTokenRemainingTime(tokenString string) (time.Duration, error) {
	expireTime, err := j.GetTokenExpireTime(tokenString)
	if err != nil {
		return 0, err
	}

	remaining := expireTime.Sub(time.Now())
	if remaining < 0 {
		return 0, errors.New("token已过期")
	}

	return remaining, nil
}

// 全局JWT管理器实例
var jwtManager *JWTManager

// InitJWT 初始化JWT管理器
func InitJWT(config JWTConfig) {
	jwtManager = NewJWTManager(config)
}

// GetJWTManager 获取JWT管理器实例
func GetJWTManager() *JWTManager {
	return jwtManager
}

// 以下是全局便捷函数，使用默认的JWT管理器

// GenerateAccessToken 生成访问token
func GenerateAccessToken(userID int64, username string) (string, error) {
	if jwtManager == nil {
		return "", errors.New("JWT管理器未初始化")
	}
	return jwtManager.GenerateAccessToken(userID, username)
}

// GenerateRefreshToken 生成刷新token
func GenerateRefreshToken(userID int64) (string, error) {
	if jwtManager == nil {
		return "", errors.New("JWT管理器未初始化")
	}
	return jwtManager.GenerateRefreshToken(userID)
}

// ValidateToken 验证token有效性
func ValidateToken(tokenString string) (*Claims, error) {
	if jwtManager == nil {
		return nil, errors.New("JWT管理器未初始化")
	}
	return jwtManager.ValidateToken(tokenString)
}

// ValidateRefreshToken 验证刷新token有效性
func ValidateRefreshToken(tokenString string) (*Claims, error) {
	if jwtManager == nil {
		return nil, errors.New("JWT管理器未初始化")
	}
	return jwtManager.ValidateRefreshToken(tokenString)
}
