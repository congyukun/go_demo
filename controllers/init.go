package controllers

import "github.com/gin-gonic/gin"

// response 统一响应格式
func response(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	})
}
