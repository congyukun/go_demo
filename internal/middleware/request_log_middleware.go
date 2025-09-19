package middleware

import (
	"bytes"
	"encoding/json"
	"go_demo/internal/utils"
	"go_demo/pkg/logger"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// responseWriter 包装gin.ResponseWriter以捕获响应内容
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// RequestLogMiddleware 请求日志中间件 - 记录请求参数和返回值
func RequestLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		requestID := utils.GetRequestID(c)
		traceId := c.GetString("trace_id")

		// 读取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// 重新设置请求体，以便后续处理器可以读取
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 获取查询参数
		queryParams := c.Request.URL.Query()

		// 获取路径参数
		pathParams := make(map[string]string)
		for _, param := range c.Params {
			pathParams[param.Key] = param.Value
		}

		// 包装ResponseWriter以捕获响应内容
		blw := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 计算处理时间
		duration := time.Since(startTime)

		// 获取响应内容
		responseBody := blw.body.String()

		// 尝试解析响应为JSON以便更好的日志记录
		var responseData interface{}
		if err := json.Unmarshal([]byte(responseBody), &responseData); err != nil {
			// 如果不是JSON，直接记录字符串
			responseData = responseBody
		}

		// 记录完整的请求和响应信息
		logger.ReqInfo("req",
			logger.String("trace_id", traceId),
			logger.String("request_id", requestID),
			logger.String("method", c.Request.Method),
			logger.String("path", c.Request.URL.Path),
			logger.String("client_ip", c.ClientIP()),
			logger.String("user_agent", c.Request.UserAgent()),
			logger.Any("query_params", queryParams),
			logger.Any("path_params", pathParams),
			logger.String("request_body", string(requestBody)),
			logger.Int("status_code", c.Writer.Status()),
			logger.String("duration", duration.String()),
			logger.Any("response_body", responseData),
		)
	}
}

// RequestLogMiddlewareWithConfig 带配置的请求日志中间件
func RequestLogMiddlewareWithConfig(config RequestLogConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		requestID := utils.GetRequestID(c)

		// 读取请求体
		var requestBody []byte
		if c.Request.Body != nil && config.LogRequestBody {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// 重新设置请求体，以便后续处理器可以读取
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 获取查询参数
		var queryParams map[string][]string
		if config.LogQueryParams {
			queryParams = c.Request.URL.Query()
		}

		// 获取路径参数
		var pathParams map[string]string
		if config.LogPathParams {
			pathParams = make(map[string]string)
			for _, param := range c.Params {
				pathParams[param.Key] = param.Value
			}
		}

		// 获取请求头
		var headers map[string][]string
		if config.LogHeaders {
			headers = c.Request.Header
		}

		// 包装ResponseWriter以捕获响应内容
		var blw *responseWriter
		if config.LogResponseBody {
			blw = &responseWriter{
				ResponseWriter: c.Writer,
				body:           bytes.NewBufferString(""),
			}
			c.Writer = blw
		}

		// 处理请求
		c.Next()

		// 计算处理时间
		duration := time.Since(startTime)

		// 构建完整的请求和响应日志字段
		logFields := []zap.Field{
			logger.String("trace_id", c.GetString("trace_id")),
			logger.String("request_id", requestID),
			logger.String("method", c.Request.Method),
			logger.String("path", c.Request.URL.Path),
			logger.String("client_ip", c.ClientIP()),
			logger.Int("status_code", c.Writer.Status()),
			logger.String("duration", duration.String()),
		}

		// 添加可选字段
		if config.LogUserAgent {
			logFields = append(logFields, logger.String("user_agent", c.Request.UserAgent()))
		}
		if config.LogQueryParams && queryParams != nil {
			logFields = append(logFields, logger.Any("query_params", queryParams))
		}
		if config.LogPathParams && pathParams != nil {
			logFields = append(logFields, logger.Any("path_params", pathParams))
		}
		if config.LogHeaders && headers != nil {
			logFields = append(logFields, logger.Any("headers", headers))
		}
		if config.LogRequestBody && len(requestBody) > 0 {
			logFields = append(logFields, logger.String("request_body", string(requestBody)))
		}
		if config.LogResponseBody && blw != nil {
			responseBody := blw.body.String()
			var responseData interface{}
			if err := json.Unmarshal([]byte(responseBody), &responseData); err != nil {
				// 如果不是JSON，直接记录字符串
				responseData = responseBody
			}
			logFields = append(logFields, logger.Any("response_body", responseData))
		}

		// 记录完整的请求和响应信息
		logger.ReqInfo("req", logFields...)
	}
}

// RequestLogConfig 请求日志配置
type RequestLogConfig struct {
	LogRequestBody  bool // 是否记录请求体
	LogResponseBody bool // 是否记录响应体
	LogQueryParams  bool // 是否记录查询参数
	LogPathParams   bool // 是否记录路径参数
	LogHeaders      bool // 是否记录请求头
	LogUserAgent    bool // 是否记录User-Agent
}

// DefaultRequestLogConfig 默认请求日志配置
func DefaultRequestLogConfig() RequestLogConfig {
	return RequestLogConfig{
		LogRequestBody:  true,
		LogResponseBody: true,
		LogQueryParams:  true,
		LogPathParams:   true,
		LogHeaders:      false, // 默认不记录请求头，可能包含敏感信息
		LogUserAgent:    true,
	}
}