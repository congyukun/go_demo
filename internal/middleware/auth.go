package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"go_demo/internal/utils"
	"go_demo/pkg/logger"
)

// JWTAuthMiddleware JWT认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := utils.GetRequestID(c)

		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("JWT认证失败：未提供Authorization头",
				logger.String("request_id", requestID),
				logger.String("path", c.Request.URL.Path),
				logger.String("client_ip", c.ClientIP()),
			)
			utils.ResponseError(c, http.StatusUnauthorized, "未认证")
			c.Abort()
			return
		}

		// 验证token格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			logger.Warn("JWT认证失败：Authorization格式错误",
				logger.String("request_id", requestID),
				logger.String("path", c.Request.URL.Path),
				logger.String("client_ip", c.ClientIP()),
			)
			utils.ResponseError(c, http.StatusUnauthorized, "认证格式错误")
			c.Abort()
			return
		}

		// 解析token
		tokenString := parts[1]
		// 使用正确的token验证函数名称
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			logger.Warn("JWT认证失败：token解析错误",
				logger.String("request_id", requestID),
				logger.String("path", c.Request.URL.Path),
				logger.String("client_ip", c.ClientIP()),
				logger.Err(err),
			)
			utils.ResponseError(c, http.StatusUnauthorized, "认证过期或无效")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		userID := claims.UserID
		username := claims.Username

		c.Set("user_id", userID)
		c.Set("username", username)

		logger.Debug("JWT认证通过",
			logger.String("request_id", requestID),
			logger.Int64("user_id", userID),
			logger.String("username", username),
			logger.String("path", c.Request.URL.Path),
		)

		c.Next()
	}
}
