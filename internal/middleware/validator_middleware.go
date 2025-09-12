package middleware

import (
	"encoding/json"
	"go_demo/internal/common"
	"go_demo/pkg/validator"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ValidateJSON 验证JSON请求体的中间件
func ValidateJSON(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 绑定JSON到结构体
		if err := c.ShouldBindJSON(obj); err != nil {
			common.ResponseError(c, http.StatusBadRequest, "请求参数格式错误: "+err.Error())
			c.Abort()
			return
		}

		// 使用自定义验证器进行验证
		if err := validator.ValidateStruct(obj); err != nil {
			common.ResponseError(c, http.StatusBadRequest, err.Error())
			c.Abort()
			return
		}

		// 将验证后的对象存储到上下文中
		c.Set("validated_data", obj)
		c.Next()
	}
}

// ValidateQuery 验证查询参数的中间件
func ValidateQuery(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 绑定查询参数到结构体
		if err := c.ShouldBindQuery(obj); err != nil {
			common.ResponseError(c, http.StatusBadRequest, "查询参数格式错误: "+err.Error())
			c.Abort()
			return
		}

		// 使用自定义验证器进行验证
		if err := validator.ValidateStruct(obj); err != nil {
			common.ResponseError(c, http.StatusBadRequest, err.Error())
			c.Abort()
			return
		}

		// 将验证后的对象存储到上下文中
		c.Set("validated_query", obj)
		c.Next()
	}
}

// ValidateURI 验证URI参数的中间件
func ValidateURI(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 绑定URI参数到结构体
		if err := c.ShouldBindUri(obj); err != nil {
			common.ResponseError(c, http.StatusBadRequest, "URI参数格式错误: "+err.Error())
			c.Abort()
			return
		}

		// 使用自定义验证器进行验证
		if err := validator.ValidateStruct(obj); err != nil {
			common.ResponseError(c, http.StatusBadRequest, err.Error())
			c.Abort()
			return
		}

		// 将验证后的对象存储到上下文中
		c.Set("validated_uri", obj)
		c.Next()
	}
}

// GetValidatedData 从上下文中获取验证后的JSON数据
func GetValidatedData(c *gin.Context) (interface{}, bool) {
	return c.Get("validated_data")
}

// GetValidatedQuery 从上下文中获取验证后的查询参数
func GetValidatedQuery(c *gin.Context) (interface{}, bool) {
	return c.Get("validated_query")
}

// GetValidatedURI 从上下文中获取验证后的URI参数
func GetValidatedURI(c *gin.Context) (interface{}, bool) {
	return c.Get("validated_uri")
}

// ValidateStruct 直接验证结构体的辅助函数
func ValidateStruct(c *gin.Context, obj interface{}) bool {
	if err := validator.ValidateStruct(obj); err != nil {
		common.ResponseError(c, http.StatusBadRequest, err.Error())
		return false
	}
	return true
}

// ValidateAndBind 绑定并验证JSON请求的辅助函数
func ValidateAndBind(c *gin.Context, obj interface{}) bool {
	// 直接读取请求体并解析JSON，完全绕过Gin的验证器
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		common.ResponseError(c, http.StatusBadRequest, "读取请求体失败: "+err.Error())
		return false
	}

	// 解析JSON
	if err := json.Unmarshal(body, obj); err != nil {
		common.ResponseError(c, http.StatusBadRequest, "请求参数格式错误: "+err.Error())
		return false
	}

	// 使用自定义验证器进行验证
	if err := validator.ValidateStruct(obj); err != nil {
		common.ResponseError(c, http.StatusBadRequest, err.Error())
		return false
	}

	return true
}