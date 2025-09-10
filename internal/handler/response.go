package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

// ResponseSuccess 成功响应
func ResponseSuccess(c *gin.Context, message string, data interface{}) {
	requestID := GetRequestID(c)
	c.JSON(http.StatusOK, Response{
		Code:      200,
		Message:   message,
		Data:      data,
		RequestID: requestID,
	})
}

// ResponseError 错误响应
func ResponseError(c *gin.Context, httpCode int, message string) {
	requestID := GetRequestID(c)
	c.JSON(httpCode, Response{
		Code:      httpCode,
		Message:   message,
		RequestID: requestID,
	})
}

// ResponseErrorWithCode 带自定义错误码的错误响应
func ResponseErrorWithCode(c *gin.Context, httpCode, errCode int, message string) {
	requestID := GetRequestID(c)
	c.JSON(httpCode, Response{
		Code:      errCode,
		Message:   message,
		RequestID: requestID,
	})
}

// GetRequestID 获取请求ID
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}
