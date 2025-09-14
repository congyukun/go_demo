package middleware

import (
	"crypto/rand"
	"fmt"
	"go_demo/internal/service"
	"go_demo/internal/utils"
	"go_demo/pkg/logger"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RequestIDMiddleware 请求ID中间件
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// CORSMiddleware CORS中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Request-ID")
		c.Header("Access-Control-Expose-Headers", "Content-Length, X-Request-ID")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := utils.GetRequestID(c)

		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("认证失败：未提供Authorization头",
				logger.String("request_id", requestID),
				logger.String("path", c.Request.URL.Path),
				logger.String("client_ip", c.ClientIP()),
			)
			utils.ResponseError(c, http.StatusUnauthorized, "未提供认证token")
			c.Abort()
			return
		}

		// 解析token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			// 如果没有Bearer前缀，直接使用原值
			token = authHeader
		}

		// 验证token（这里使用一个简单的验证，实际项目中应该注入AuthService）
		authService := service.NewAuthService(nil) // 临时创建，实际应该通过依赖注入
		claims, err := authService.ValidateToken(token)
		if err != nil {
			logger.Warn("认证失败：token无效",
				logger.String("request_id", requestID),
				logger.String("path", c.Request.URL.Path),
				logger.String("client_ip", c.ClientIP()),
				logger.Err(err),
			)
			utils.ResponseError(c, http.StatusUnauthorized, "token无效")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		logger.Debug("认证成功",
			logger.String("request_id", requestID),
			logger.Int("user_id", claims.UserID),
			logger.String("username", claims.Username),
			logger.String("path", c.Request.URL.Path),
		)

		c.Next()
	}
}

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.Info("HTTP请求",
			logger.String("method", param.Method),
			logger.String("path", param.Path),
			logger.Int("status", param.StatusCode),
			logger.String("client_ip", param.ClientIP),
			logger.String("user_agent", param.Request.UserAgent()),
			logger.String("latency", param.Latency.String()),
		)
		return ""
	})
}

// RecoveryMiddleware 恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		requestID := utils.GetRequestID(c)
		logger.Error("服务器内部错误",
			logger.String("request_id", requestID),
			logger.String("path", c.Request.URL.Path),
			logger.String("method", c.Request.Method),
			logger.Any("error", recovered),
		)
		utils.ResponseError(c, http.StatusInternalServerError, "服务器内部错误")
	})
}

// JWTMiddleware JWT认证中间件
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := utils.GetRequestID(c)

		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("JWT认证失败：未提供Authorization头",
				logger.String("request_id", requestID),
				logger.String("path", c.Request.URL.Path),
				logger.String("client_ip", c.ClientIP()),
			)
			utils.ResponseError(c, http.StatusUnauthorized, "未提供认证token")
			c.Abort()
			return
		}

		// 提取token（去掉Bearer前缀）
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			logger.Warn("JWT认证失败：token格式错误",
				logger.String("request_id", requestID),
				logger.String("path", c.Request.URL.Path),
				logger.String("client_ip", c.ClientIP()),
			)
			utils.ResponseError(c, http.StatusUnauthorized, "token格式错误")
			c.Abort()
			return
		}

		// 验证JWT token
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			logger.Warn("JWT认证失败：token无效",
				logger.String("request_id", requestID),
				logger.String("path", c.Request.URL.Path),
				logger.String("client_ip", c.ClientIP()),
				logger.Err(err),
			)
			utils.ResponseError(c, http.StatusUnauthorized, "token无效")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		logger.Debug("JWT认证成功",
			logger.String("request_id", requestID),
			logger.Int("user_id", claims.UserID),
			logger.String("username", claims.Username),
			logger.String("path", c.Request.URL.Path),
		)

		c.Next()
	}
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
