package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// User 结构体（仅演示，实际应有持久化与加密）
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginController 用户相关控制器
type LoginController struct{}

// 响应辅助函数

// 用户注册
// @Summary 用户注册
// @Description 注册新用户
// @Accept json
// @Produce json
// @Param user body User true "用户信息"
// @Success 200 {object} User
// @Failure 400 {string} string "参数错误"
// @Router /register [post]
func (lc *LoginController) Register(c *gin.Context) {
	var req User
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		response(c, http.StatusBadRequest, "参数错误，用户名和密码必填", nil)
		return
	}
	// 实际应保存用户到数据库，此处仅演示
	response(c, http.StatusOK, "注册成功!", req)
}

// 用户登录
// @Summary 用户登录
// @Description 用户登录接口
// @Accept json
// @Produce json
// @Param user body User true "用户信息"
// @Success 200 {object} User
// @Failure 400 {string} string "参数错误"
// @Router /login [post]
func (lc *LoginController) Login(c *gin.Context) {
	var req User
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		response(c, http.StatusBadRequest, "参数错误，用户名和密码必填", nil)
		return
	}
	// 实际应校验用户名密码，此处仅演示
	response(c, http.StatusOK, "登录成功!", req)
}
