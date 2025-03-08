package controllers

import "github.com/gin-gonic/gin"

type LoginController struct{}

func (lc *LoginController) Login(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "登录成功!",
	})
	// 处理登录逻辑
}

func (lc *LoginController) Register(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "注册成功!",
	})
	// 处理注册逻辑
}
