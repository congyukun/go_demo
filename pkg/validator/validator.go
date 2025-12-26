package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	// 手机号正则
	mobileRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)
	// 邮箱正则
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-]+\.[a-zA-Z]{2,}$`)
	// 用户名正则（字母开头，字母数字下划线，4-20位）
	usernameRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{3,19}$`)
)

// Init 初始化验证器
func Init() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册自定义验证器
		if err := v.RegisterValidation("mobile", validateMobile); err != nil {
			return fmt.Errorf("注册mobile验证器失败: %w", err)
		}
		if err := v.RegisterValidation("username", validateUsername); err != nil {
			return fmt.Errorf("注册username验证器失败: %w", err)
		}
		if err := v.RegisterValidation("strong_password", validateStrongPassword); err != nil {
			return fmt.Errorf("注册strong_password验证器失败: %w", err)
		}

		// 注册自定义标签名称函数
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
	return nil
}

// validateMobile 验证手机号
func validateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	return mobileRegex.MatchString(mobile)
}

// validateUsername 验证用户名
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	return usernameRegex.MatchString(username)
}

// validateStrongPassword 验证强密码（至少8位，包含大小写字母和数字）
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}

	hasUpper := false
	hasLower := false
	hasDigit := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}

// TranslateError 翻译验证错误信息
func TranslateError(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := e.Field()
			tag := e.Tag()

			switch tag {
			case "required":
				errors[field] = fmt.Sprintf("%s是必填字段", field)
			case "email":
				errors[field] = fmt.Sprintf("%s格式不正确", field)
			case "min":
				errors[field] = fmt.Sprintf("%s长度不能小于%s", field, e.Param())
			case "max":
				errors[field] = fmt.Sprintf("%s长度不能大于%s", field, e.Param())
			case "len":
				errors[field] = fmt.Sprintf("%s长度必须为%s", field, e.Param())
			case "mobile":
				errors[field] = fmt.Sprintf("%s格式不正确", field)
			case "username":
				errors[field] = fmt.Sprintf("%s格式不正确（字母开头，4-20位字母数字下划线）", field)
			case "strong_password":
				errors[field] = fmt.Sprintf("%s强度不够（至少8位，包含大小写字母和数字）", field)
			case "eqfield":
				errors[field] = fmt.Sprintf("%s必须等于%s", field, e.Param())
			case "nefield":
				errors[field] = fmt.Sprintf("%s不能等于%s", field, e.Param())
			case "gt":
				errors[field] = fmt.Sprintf("%s必须大于%s", field, e.Param())
			case "gte":
				errors[field] = fmt.Sprintf("%s必须大于等于%s", field, e.Param())
			case "lt":
				errors[field] = fmt.Sprintf("%s必须小于%s", field, e.Param())
			case "lte":
				errors[field] = fmt.Sprintf("%s必须小于等于%s", field, e.Param())
			case "oneof":
				errors[field] = fmt.Sprintf("%s必须是以下值之一: %s", field, e.Param())
			default:
				errors[field] = fmt.Sprintf("%s验证失败", field)
			}
		}
	}

	return errors
}

// ValidateMobile 验证手机号
func ValidateMobile(mobile string) bool {
	return mobileRegex.MatchString(mobile)
}

// ValidateEmail 验证邮箱
func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// ValidateUsername 验证用户名
func ValidateUsername(username string) bool {
	return usernameRegex.MatchString(username)
}

// ValidatePassword 验证密码强度
func ValidatePassword(password string, minLength int) bool {
	if len(password) < minLength {
		return false
	}

	hasUpper := false
	hasLower := false
	hasDigit := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}

// ValidateStruct 验证结构体
func ValidateStruct(obj interface{}) error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.Struct(obj); err != nil {
			// 翻译错误信息
			errors := TranslateError(err)
			// 将错误信息组合成字符串
			var errMsgs []string
			for _, msg := range errors {
				errMsgs = append(errMsgs, msg)
			}
			return fmt.Errorf("%s", strings.Join(errMsgs, "; "))
		}
	}
	return nil
}
