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
	// ErrorTypeConflict 资源冲突
	ErrorTypeConflict
	// ErrorTypeTooManyRequests 请求过多
	ErrorTypeTooManyRequests
	// ErrorTypeServiceUnavailable 服务不可用
	ErrorTypeServiceUnavailable
)

// AppError 应用错误
type AppError struct {
	Type      ErrorType // 错误类型
	Err       error     // 原始错误
	Message   string    // 错误消息
	Details   string    // 错误详情
	Stack     string    // 错误堆栈
	HTTPCode  int       // HTTP状态码
	ErrorCode string    // 业务错误码
	Cause     error     // 根本原因
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Details)
	}
	return e.Message
}

// Unwrap 返回原始错误
func (e *AppError) Unwrap() error {
	return e.Err
}

// Is 检查错误类型
func (e *AppError) Is(target error) bool {
	if t, ok := target.(*AppError); ok {
		return e.Type == t.Type
	}
	return false
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

// NewWithDetails 创建带详情的应用错误
func NewWithDetails(errType ErrorType, message, details string) *AppError {
	err := New(errType, message)
	err.Details = details
	return err
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

// WrapWithDetails 包装已有错误并添加详情
func WrapWithDetails(err error, errType ErrorType, message, details string) *AppError {
	appErr := Wrap(err, errType, message)
	if appErr != nil {
		appErr.Details = details
	}
	return appErr
}

// WrapAndLog 包装错误并记录日志
func WrapAndLog(err error, errType ErrorType, message string) *AppError {
	appErr := Wrap(err, errType, message)
	if appErr != nil {
		logger.Error(message,
			logger.Err(err),
			logger.String("error_code", appErr.ErrorCode),
			logger.String("error_type", getErrorTypeName(errType)),
			logger.String("stack", appErr.Stack),
		)
	}
	return appErr
}

// WithCause 添加根本原因
func (e *AppError) WithCause(cause error) *AppError {
	e.Cause = cause
	return e
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
	case ErrorTypeConflict:
		return http.StatusConflict
	case ErrorTypeTooManyRequests:
		return http.StatusTooManyRequests
	case ErrorTypeServiceUnavailable:
		return http.StatusServiceUnavailable
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
	case ErrorTypeConflict:
		return "E3002"
	case ErrorTypeTooManyRequests:
		return "E4001"
	case ErrorTypeServiceUnavailable:
		return "E4002"
	case ErrorTypeDatabase:
		return "E5001"
	case ErrorTypeInternal:
		return "E5002"
	default:
		return "E9999"
	}
}

// getErrorTypeName 获取错误类型名称
func getErrorTypeName(errType ErrorType) string {
	switch errType {
	case ErrorTypeValidation:
		return "Validation"
	case ErrorTypeAuthorization:
		return "Authorization"
	case ErrorTypeNotFound:
		return "NotFound"
	case ErrorTypeConflict:
		return "Conflict"
	case ErrorTypeTooManyRequests:
		return "TooManyRequests"
	case ErrorTypeServiceUnavailable:
		return "ServiceUnavailable"
	case ErrorTypeDatabase:
		return "Database"
	case ErrorTypeInternal:
		return "Internal"
	default:
		return "Unknown"
	}
}

// 预定义的常用错误
var (
	ErrInvalidCredentials = New(ErrorTypeAuthorization, "无效的凭据")
	ErrTokenExpired       = New(ErrorTypeAuthorization, "令牌已过期")
	ErrTokenInvalid       = New(ErrorTypeAuthorization, "令牌无效")
	ErrInvalidToken       = New(ErrorTypeAuthorization, "令牌无效")
	ErrUserNotFound       = New(ErrorTypeNotFound, "用户不存在")
	ErrUserExists         = New(ErrorTypeConflict, "用户已存在")
	ErrInvalidRequest     = New(ErrorTypeValidation, "无效的请求")
	ErrPermissionDenied   = New(ErrorTypeAuthorization, "权限不足")
	ErrServiceUnavailable = New(ErrorTypeServiceUnavailable, "服务暂时不可用")
	ErrRateLimitExceeded  = New(ErrorTypeTooManyRequests, "请求频率超限")
)

// NewValidationError 创建验证错误
func NewValidationError(message string) *AppError {
	return New(ErrorTypeValidation, message)
}

// NewInternalServerError 创建内部服务器错误
func NewInternalServerError(message string) *AppError {
	return New(ErrorTypeInternal, message)
}

// NewConflictError 创建冲突错误
func NewConflictError(message string) *AppError {
	return New(ErrorTypeConflict, message)
}

// NewForbiddenError 创建禁止访问错误
func NewForbiddenError(message string) *AppError {
	return New(ErrorTypeAuthorization, message)
}

// NewNotFoundError 创建资源不存在错误
func NewNotFoundError(message string) *AppError {
	return New(ErrorTypeNotFound, message)
}