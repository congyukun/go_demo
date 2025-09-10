package handler

import (
	"go_demo/internal/models"
	"go_demo/internal/service"
	"go_demo/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	requestID := GetRequestID(c)

	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("登录参数错误",
			logger.String("request_id", requestID),
			logger.Err(err),
			logger.String("client_ip", c.ClientIP()),
		)
		ResponseError(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	logger.Info("用户登录请求",
		logger.String("request_id", requestID),
		logger.String("username", req.Username),
		logger.String("client_ip", c.ClientIP()),
		logger.String("user_agent", c.Request.UserAgent()),
	)

	// 调用服务层处理登录
	response, err := h.authService.Login(req)
	if err != nil {
		logger.Warn("登录失败",
			logger.String("request_id", requestID),
			logger.String("username", req.Username),
			logger.String("client_ip", c.ClientIP()),
			logger.Err(err),
		)
		ResponseError(c, http.StatusUnauthorized, err.Error())
		return
	}

	logger.Info("用户登录成功",
		logger.String("request_id", requestID),
		logger.String("username", req.Username),
		logger.Int("user_id", response.User.ID),
		logger.String("client_ip", c.ClientIP()),
	)

	ResponseSuccess(c, "登录成功", response)
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	requestID := GetRequestID(c)

	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("注册参数错误",
			logger.String("request_id", requestID),
			logger.Err(err),
			logger.String("client_ip", c.ClientIP()),
		)
		ResponseError(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	logger.Info("用户注册请求",
		logger.String("request_id", requestID),
		logger.String("username", req.Username),
		logger.String("email", req.Email),
		logger.String("name", req.Name),
		logger.String("client_ip", c.ClientIP()),
	)

	// 调用服务层处理注册
	user, err := h.authService.Register(req)
	if err != nil {
		logger.Warn("注册失败",
			logger.String("request_id", requestID),
			logger.String("username", req.Username),
			logger.Err(err),
		)
		ResponseError(c, http.StatusConflict, err.Error())
		return
	}

	logger.Info("用户注册成功",
		logger.String("request_id", requestID),
		logger.String("username", req.Username),
		logger.String("email", req.Email),
		logger.Int("user_id", user.ID),
		logger.String("client_ip", c.ClientIP()),
	)

	ResponseSuccess(c, "注册成功", user)
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *gin.Context) {
	requestID := GetRequestID(c)

	// 获取Authorization头
	token := c.GetHeader("Authorization")

	logger.Info("用户登出请求",
		logger.String("request_id", requestID),
		logger.String("token", maskToken(token)),
		logger.String("client_ip", c.ClientIP()),
	)

	if token == "" {
		logger.Warn("登出失败：未提供token",
			logger.String("request_id", requestID),
			logger.String("client_ip", c.ClientIP()),
		)
		ResponseError(c, http.StatusUnauthorized, "未提供认证token")
		return
	}

	// 验证token
	_, err := h.authService.ValidateToken(token)
	if err != nil {
		logger.Warn("登出失败：token无效",
			logger.String("request_id", requestID),
			logger.String("client_ip", c.ClientIP()),
			logger.Err(err),
		)
		ResponseError(c, http.StatusUnauthorized, "token无效")
		return
	}

	logger.Info("用户登出成功",
		logger.String("request_id", requestID),
		logger.String("token", maskToken(token)),
		logger.String("client_ip", c.ClientIP()),
	)

	ResponseSuccess(c, "登出成功", nil)
}

// maskToken 掩码token用于日志记录
func maskToken(token string) string {
	if len(token) <= 10 {
		return "***"
	}
	return token[:10] + "***"
}
