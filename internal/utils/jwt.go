package utils

import (
	"fmt"
	"go_demo/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT密钥，实际项目中应该从配置文件读取
var jwtSecret = []byte("your-secret-key-change-this-in-production")

// GenerateJWT 生成JWT token
func GenerateJWT(userID int, username string) (string, error) {
	// 创建claims
	claims := &models.TokenClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24小时过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go_demo",
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名token
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("生成token失败: %w", err)
	}

	return tokenString, nil
}

// ValidateJWT 验证JWT token
func ValidateJWT(tokenString string) (*models.TokenClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("解析token失败: %w", err)
	}

	// 验证token是否有效
	if !token.Valid {
		return nil, fmt.Errorf("token无效")
	}

	// 提取claims
	claims, ok := token.Claims.(*models.TokenClaims)
	if !ok {
		return nil, fmt.Errorf("无法提取token claims")
	}

	return claims, nil
}

// RefreshJWT 刷新JWT token
func RefreshJWT(tokenString string) (string, error) {
	// 验证当前token
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return "", fmt.Errorf("当前token无效: %w", err)
	}

	// 生成新token
	return GenerateJWT(claims.UserID, claims.Username)
}