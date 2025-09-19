package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TraceIDKey struct{}

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 尝试从请求头获取 trace_id（如果是上游服务传递的）
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			// 2. 如果请求头中没有，则生成一个新的 trace_id
			traceID = uuid.New().String()
		}
		// 3. 将 trace_id 存储到上下文中
		c.Set("trace_id", traceID)
		// 4. 将 trace_id 添加到响应头，方便客户端调试
		c.Header("X-Trace-ID", traceID)
		//
		ctx := context.WithValue(c.Request.Context(), TraceIDKey{}, traceID)
		c.Request = c.Request.WithContext(ctx)
		// 继续处理请求
		c.Next()
	}
}

// GetTraceID 从 Gin 上下文获取 trace_id
func GetTraceID(ctx *gin.Context) string {
	traceID, _ := ctx.Get("trace_id")
	if id, ok := traceID.(string); ok {
		return id
	}
	return ""
}

// GetTraceIDFromContext 从 context 中获取 trace_id
func GetTraceIDFromContext(ctx context.Context) string {
	traceId := ctx.Value(TraceIDKey{})
	if id, ok := traceId.(string); ok {
		return id
	}
	return ""
}
