package validator

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

// Init 初始化验证器
func Init() error {
	// 这里可以添加任何需要初始化的逻辑
	// 目前为空，但保持此函数以满足测试文件的需求
	return nil
}

// Validator 自定义验证器
var (
	// 错误消息映射
	errorMessages = map[string]string{
		"required": "%s不能为空",
		"min":      "%s长度不能小于%d",
		"max":      "%s长度不能大于%d",
		"email":    "%s必须是有效的邮箱地址",
		"mobile":    "%s必须是有效的手机号",
		"url":      "%s必须是有效的URL",
	}
)

// ValidateStruct 验证结构体
func ValidateStruct(obj interface{}) error {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("验证对象必须是结构体")
	}

	typ := val.Type()
	var errors []string

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// 获取json标签
		jsonTag := fieldType.Tag.Get("json")
		// 获取label标签
		labelTag := fieldType.Tag.Get("label")
		fieldName := getFieldName(jsonTag, fieldType.Name, labelTag)

		// 获取验证标签
		validateTag := fieldType.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		// 分割验证规则
		rules := strings.Split(validateTag, ",")
		
		// 检查是否有omitempty规则
		hasOmitEmpty := false
		otherRules := []string{}
		for _, rule := range rules {
			if rule == "omitempty" {
				hasOmitEmpty = true
			} else {
				otherRules = append(otherRules, rule)
			}
		}
		
		// 如果有omitempty规则且字段为空，则跳过其他验证
		if hasOmitEmpty && isEmpty(field) {
			continue
		}
		
		// 验证其他规则
		for _, rule := range otherRules {
			err := validateField(field, fieldName, rule)
			if err != nil {
				errors = append(errors, err.Error())
			}
		}
	}

	if len(errors) > 0 {
		const errorFormat = "%s"
		return fmt.Errorf(errorFormat, strings.Join(errors, "; "))
	}

	return nil
}

// getFieldName 获取字段名称
func getFieldName(jsonTag, defaultName string, labelTag string) string {
	// 优先使用label标签
	if labelTag != "" {
		return labelTag
	}
	
	// 其次使用json标签
	if jsonTag != "" {
		parts := strings.Split(jsonTag, ",")
		if len(parts) > 0 && parts[0] != "-" {
			return parts[0]
		}
	}
	
	// 最后使用字段名
	return defaultName
}

// isEmpty 检查字段是否为空
func isEmpty(field reflect.Value) bool {
	switch field.Kind() {
	case reflect.String:
		return field.String() == ""
	case reflect.Ptr, reflect.Interface:
		return field.IsNil()
	case reflect.Slice, reflect.Array, reflect.Map:
		return field.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return field.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return field.Float() == 0
	case reflect.Bool:
		return !field.Bool()
	}
	return false
}

// validateField 验证单个字段
func validateField(field reflect.Value, fieldName, rule string) error {
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			if rule == "required" {
				return fmt.Errorf("%s不能为空", fieldName)
			}
			return nil
		}
		field = field.Elem()
	}

	switch rule {
	case "required":
		return validateRequired(field, fieldName)
	default:
		if strings.HasPrefix(rule, "min=") {
			return validateMin(field, fieldName, rule[4:])
		}
		if strings.HasPrefix(rule, "max=") {
			return validateMax(field, fieldName, rule[4:])
		}
		if rule == "email" {
			return validateEmail(field, fieldName)
		}
		if rule == "mobile" {
			return validateMobile(field, fieldName)
		}
		if rule == "url" {
			return validateURL(field, fieldName)
		}
		if rule == "mobile" {
			return validateMobile(field, fieldName)
		}
		if strings.HasPrefix(rule, "oneof=") {
			return validateOneOf(field, fieldName, rule[6:])
		}
	}

	return nil
}

// validateRequired 验证必填
func validateRequired(field reflect.Value, fieldName string) error {
	switch field.Kind() {
	case reflect.String:
		if field.String() == "" {
			return fmt.Errorf("%s不能为空", fieldName)
		}
	case reflect.Slice, reflect.Array, reflect.Map:
		if field.Len() == 0 {
			return fmt.Errorf("%s不能为空", fieldName)
		}
	}

	return nil
}

// validateMin 验证最小值
func validateMin(field reflect.Value, fieldName, minStr string) error {
	var min int
	if _, err := fmt.Sscanf(minStr, "%d", &min); err != nil {
		return fmt.Errorf("无效的min规则")
	}

	switch field.Kind() {
	case reflect.String:
		if len(field.String()) < min {
			return fmt.Errorf("%s长度不能小于%d", fieldName, min)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Int() < int64(min) {
			return fmt.Errorf("%s不能小于%d", fieldName, min)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if field.Uint() < uint64(min) {
			return fmt.Errorf("%s不能小于%d", fieldName, min)
		}
	}

	return nil
}

// validateMax 验证最大值
func validateMax(field reflect.Value, fieldName, maxStr string) error {
	var max int
	if _, err := fmt.Sscanf(maxStr, "%d", &max); err != nil {
		return fmt.Errorf("无效的max规则")
	}

	switch field.Kind() {
	case reflect.String:
		if len(field.String()) > max {
			return fmt.Errorf("%s长度不能大于%d", fieldName, max)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Int() > int64(max) {
			return fmt.Errorf("%s不能大于%d", fieldName, max)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if field.Uint() > uint64(max) {
			return fmt.Errorf("%s不能大于%d", fieldName, max)
		}
	}

	return nil
}

// validateEmail 验证邮箱
func validateEmail(field reflect.Value, fieldName string) error {
	if field.Kind() != reflect.String {
		return fmt.Errorf("%s必须是字符串类型", fieldName)
	}

	email := field.String()
	if email == "" {
		return nil // 允许为空，必填由required规则控制
	}

	// 简单的邮箱格式验证
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return fmt.Errorf("%s必须是有效的邮箱地址", fieldName)
	}

	return nil
}


// validateURL 验证URL
func validateURL(field reflect.Value, fieldName string) error {
	if field.Kind() != reflect.String {
		return fmt.Errorf("%s必须是字符串类型", fieldName)
	}

	url := field.String()
	if url == "" {
		return nil // 允许为空，必填由required规则控制
	}

	// 简单的URL格式验证
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return fmt.Errorf("%s必须是有效的URL", fieldName)
	}

	return nil
}

// validateMobile 验证手机号
func validateMobile(field reflect.Value, fieldName string) error {
	if field.Kind() != reflect.String {
		return fmt.Errorf("%s必须是字符串类型", fieldName)
	}

	mobile := field.String()
	if mobile == "" {
		return nil // 允许为空，必填由required规则控制
	}

	// 手机号格式验证（11位数字，以1开头）
	if len(mobile) != 11 || mobile[0] != '1' {
		return fmt.Errorf("%s必须是有效的手机号", fieldName)
	}

	// 检查第二位是否符合要求（3, 4, 5, 7, 8中的一个）
	secondDigit := mobile[1]
	if secondDigit != '3' && secondDigit != '4' && secondDigit != '5' && secondDigit != '7' && secondDigit != '8' {
		return fmt.Errorf("%s必须是有效的手机号", fieldName)
	}

	// 检查是否全为数字
	for _, r := range mobile {
		if !unicode.IsDigit(r) {
			return fmt.Errorf("%s必须是有效的手机号", fieldName)
		}
	}

	return nil
}

// validateOneOf 验证字段值是否在指定范围内
func validateOneOf(field reflect.Value, fieldName, allowedStr string) error {
	// 如果字段为空，则跳过验证
	if isEmpty(field) {
		return nil
	}

	allowedValues := strings.Split(allowedStr, " ")
	valueStr := fmt.Sprintf("%v", field.Interface())

	for _, allowed := range allowedValues {
		if valueStr == allowed {
			return nil
		}
	}

	return fmt.Errorf("%s必须是以下值之一: %s", fieldName, allowedStr)
}
