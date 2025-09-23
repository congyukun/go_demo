package middleware

import (
	"crypto/rand"
	"fmt"
	"go_demo/internal/utils"
	"go_demo/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequestID 请求ID中间件
func RequestID() gin.HandlerFunc {
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

// CORS CORS中间件
func CORS() gin.HandlerFunc {
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

// Logger 日志中间件
func Logger() gin.HandlerFunc {
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

// Recovery 恢复中间件
func Recovery() gin.HandlerFunc {
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

// generateRequestID 生成请求ID
func generateRequestID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}