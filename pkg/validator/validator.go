package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	// Validator 全局验证器实例
	Validator *validator.Validate
	// Trans 全局翻译器实例
	Trans ut.Translator
)

// defaultValidator 实现 Gin 的验证器接口
type defaultValidator struct {
	validate *validator.Validate
}

// ValidateStruct 实现 binding.StructValidator 接口
func (v *defaultValidator) ValidateStruct(obj interface{}) error {
	return v.validate.Struct(obj)
}

// Engine 返回底层的验证器引擎
func (v *defaultValidator) Engine() interface{} {
	return v.validate
}

// Init 初始化验证器
func Init() error {
	// 创建新的验证器实例
	Validator = validator.New()

	// 设置gin使用我们的验证器
	binding.Validator = &defaultValidator{validate: Validator}

	// 注册自定义验证规则
	if err := registerCustomValidators(); err != nil {
		return fmt.Errorf("注册自定义验证规则失败: %v", err)
	}

	// 初始化中文翻译器
	if err := initTranslator(); err != nil {
		return fmt.Errorf("初始化翻译器失败: %v", err)
	}

	// 注册自定义字段名称
	registerFieldNames()

	return nil
}

// registerCustomValidators 注册自定义验证规则
func registerCustomValidators() error {
	// 注册手机号验证规则
	if err := Validator.RegisterValidation("mobile", validateMobile); err != nil {
		return err
	}

	return nil
}

// validateMobile 手机号验证函数
func validateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	// 中国大陆手机号正则表达式
	matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, mobile)
	return matched
}

// initTranslator 初始化中文翻译器
func initTranslator() error {
	// 创建中文翻译器
	zhLocale := zh.New()
	uni := ut.New(zhLocale, zhLocale)

	var ok bool
	Trans, ok = uni.GetTranslator("zh")
	if !ok {
		return fmt.Errorf("获取中文翻译器失败")
	}

	// 注册默认翻译
	if err := zh_translations.RegisterDefaultTranslations(Validator, Trans); err != nil {
		return err
	}

	// 注册自定义验证规则的翻译
	if err := registerCustomTranslations(); err != nil {
		return err
	}

	return nil
}

// registerCustomTranslations 注册自定义翻译
func registerCustomTranslations() error {
	// 注册手机号验证的中文翻译
	if err := Validator.RegisterTranslation("mobile", Trans, func(ut ut.Translator) error {
		return ut.Add("mobile", "{0}必须是有效的手机号码", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("mobile", fe.Field())
		return t
	}); err != nil {
		return err
	}

	// 覆盖一些默认翻译
	translations := map[string]string{
		"required": "{0}为必填字段",
		"min":      "{0}长度不能少于{1}个字符",
		"max":      "{0}长度不能超过{1}个字符",
		"email":    "{0}必须是有效的邮箱地址",
		"oneof":    "{0}必须是[{1}]中的一个",
	}

	for tag, translation := range translations {
		if err := Validator.RegisterTranslation(tag, Trans, func(ut ut.Translator) error {
			return ut.Add(tag, translation, true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(tag, fe.Field(), fe.Param())
			return t
		}); err != nil {
			return err
		}
	}

	return nil
}

// registerFieldNames 注册字段名称映射
// registerFieldNames 注册字段名称映射
func registerFieldNames() {
	Validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		// 优先使用label标签作为字段名
		name := fld.Tag.Get("label")
		if name != "" {
			return name
		}

		// 如果没有label标签，使用json标签
		jsonName := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if jsonName != "" && jsonName != "-" {
			return jsonName
		}

		// 如果都没有，使用字段名
		return fld.Name
	})
}

// ValidateStruct 验证结构体并返回中文错误信息
func ValidateStruct(s interface{}) error {
	if err := Validator.Struct(s); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errorMessages []string
			for _, e := range validationErrors {
				errorMessages = append(errorMessages, e.Translate(Trans))
			}
			return fmt.Errorf(strings.Join(errorMessages, "; "))
		}
		return err
	}
	return nil
}

// GetFieldError 获取单个字段的错误信息
func GetFieldError(err error, field string) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			if e.Field() == field {
				return e.Translate(Trans)
			}
		}
	}
	return ""
}

// GetAllFieldErrors 获取所有字段的错误信息映射
func GetAllFieldErrors(err error) map[string]string {
	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors[e.Field()] = e.Translate(Trans)
		}
	}
	return errors
}
