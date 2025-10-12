package errors

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"go_demo/pkg/logger"
)

// ErrorType 是错误类型
type ErrorType uint

const (
	// ErrorTypeUnknown 未知错误
	ErrorTypeUnknown ErrorType = iota
	// ErrorTypeInternal 内部错误
	ErrorTypeInternal
	// ErrorTypeValidation 验证错误
	ErrorTypeValidation
	// ErrorTypeDatabase 数据库错误
	ErrorTypeDatabase
	// ErrorTypeAuthorization 授权错误
	ErrorTypeAuthorization
	// ErrorTypeNotFound 资源不存在
	ErrorTypeNotFound
)

// AppError 应用错误
type AppError struct {
	Type      ErrorType // 错误类型
	Err       error     // 原始错误
	Message   string    // 错误消息
	Stack     string    // 错误堆栈
	HTTPCode  int       // HTTP状态码
	ErrorCode string    // 业务错误码
}

// Error 实现error接口
func (e *AppError) Error() string {
	return e.Message
}

// Unwrap 返回原始错误
func (e *AppError) Unwrap() error {
	return e.Err
}

// New 创建新的应用错误
func New(errType ErrorType, message string) *AppError {
	stack := getStack(2)
	httpCode := getHTTPCode(errType)
	errorCode := getErrorCode(errType)

	return &AppError{
		Type:      errType,
		Message:   message,
		Stack:     stack,
		HTTPCode:  httpCode,
		ErrorCode: errorCode,
	}
}

// Wrap 包装已有错误
func Wrap(err error, errType ErrorType, message string) *AppError {
	if err == nil {
		return nil
	}

	stack := getStack(2)
	httpCode := getHTTPCode(errType)
	errorCode := getErrorCode(errType)

	return &AppError{
		Type:      errType,
		Err:       err,
		Message:   message,
		Stack:     stack,
		HTTPCode:  httpCode,
		ErrorCode: errorCode,
	}
}

// WrapAndLog 包装错误并记录日志
func WrapAndLog(err error, errType ErrorType, message string) *AppError {
	appErr := Wrap(err, errType, message)
	if appErr != nil {
		logger.Error(message,
			logger.Err(err),
			logger.String("stack", appErr.Stack),
			logger.String("error_code", appErr.ErrorCode),
		)
	}
	return appErr
}

// getStack 获取调用堆栈
func getStack(skip int) string {
	var sb strings.Builder
	for i := skip; i < skip+5; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		sb.WriteString(fmt.Sprintf("%s:%d %s\n", file, line, fn.Name()))
	}
	return sb.String()
}

// getHTTPCode 根据错误类型获取HTTP状态码
func getHTTPCode(errType ErrorType) int {
	switch errType {
	case ErrorTypeValidation:
		return http.StatusBadRequest
	case ErrorTypeAuthorization:
		return http.StatusUnauthorized
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeDatabase, ErrorTypeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// getErrorCode 根据错误类型获取业务错误码
func getErrorCode(errType ErrorType) string {
	switch errType {
	case ErrorTypeValidation:
		return "E1001"
	case ErrorTypeAuthorization:
		return "E2001"
	case ErrorTypeNotFound:
		return "E3001"
	case ErrorTypeDatabase:
		return "E4001"
	case ErrorTypeInternal:
		return "E5001"
	default:
		return "E9999"
	}
}