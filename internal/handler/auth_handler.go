package handler

import (
	"go_demo/internal/middleware"
	"go_demo/internal/models"
	"go_demo/internal/service"
	"go_demo/internal/utils"
	"go_demo/pkg/errors"
	"go_demo/pkg/logger"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// handleServiceError 统一处理服务层返回的错误
func handleServiceError(c *gin.Context, err error, requestID string) {
	// 记录错误日志
	logger.Warn("服务调用失败",
		logger.String("request_id", requestID),
		logger.Err(err),
	)

	// 根据错误类型返回不同的HTTP状态码
	appErr, ok := err.(*errors.AppError)
	if ok {
		utils.ResponseError(c, appErr.HTTPCode, appErr.Error())
		return
	}

	// 默认返回500错误
	utils.ResponseError(c, http.StatusInternalServerError, "服务器内部错误")
}

// AuthHandler 认证处理器
type AuthHandler struct {
	authService service.AuthService
	userService service.UserService
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(authService service.AuthService, userService service.UserService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
	}
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "登录请求"
// @Success 200 {object} utils.Response{data=models.LoginResponse} "登录成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 401 {object} utils.Response "认证失败"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	requestID := middleware.GetTraceID(c)

	// 绑定并验证请求参数
	var req models.LoginRequest
	if !middleware.ValidateAndBind(c, &req) {
		return
	}

	// logger.Info("用户登录请求",
	// 	logger.String("request_id", requestID),
	// 	logger.String("username", req.Username),
	// 	logger.String("client_ip", c.ClientIP()),
	// )

	// 调用服务层进行登录
	response, err := h.authService.Login(c, req)
	if err != nil {
		handleServiceError(c, err, requestID)
		return
	}

	// logger.Info("用户登录成功",
	// 	logger.String("request_id", requestID),
	// 	logger.String("username", req.Username),
	// 	logger.Int("user_id", int(response.User.ID)),
	// 	logger.String("client_ip", c.ClientIP()),
	// )

	utils.ResponseSuccess(c, "登录成功", response)
}

// Register 用户注册
// @Summary 用户注册
// @Description 用户注册接口
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "注册请求"
// @Success 200 {object} utils.Response{data=models.UserResponse} "注册成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 409 {object} utils.Response "用户已存在"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	requestID := middleware.GetTraceID(c)

	// 绑定并验证请求参数
	var req models.RegisterRequest
	if !middleware.ValidateAndBind(c, &req) {
		return
	}

	logger.Info("用户注册请求",
		logger.String("request_id", requestID),
		logger.String("username", req.Username),
		logger.String("email", req.Email),
		logger.String("client_ip", c.ClientIP()),
	)

	// 调用服务层进行注册
	user, err := h.authService.Register(c, req)
	if err != nil {
		handleServiceError(c, err, requestID)
		return
	}

	logger.Info("用户注册成功",
		logger.String("request_id", requestID),
		logger.String("username", req.Username),
		logger.Int("user_id", int(user.ID)),
		logger.String("client_ip", c.ClientIP()),
	)

	utils.ResponseSuccess(c, "注册成功", user)
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户登出接口
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response "登出成功"
// @Failure 401 {object} utils.Response "未认证"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	requestID := middleware.GetTraceID(c)

	// 从Authorization头获取token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.ResponseError(c, http.StatusUnauthorized, "未提供认证令牌")
		return
	}

	// 提取Bearer token
	tokenString := ""
	if len(authHeader) > 7 && strings.ToUpper(authHeader[0:6]) == "BEARER" {
		tokenString = authHeader[7:]
	}

	if tokenString == "" {
		utils.ResponseError(c, http.StatusUnauthorized, "无效的认证令牌")
		return
	}

	logger.Info("用户登出请求",
		logger.String("request_id", requestID),
		logger.String("client_ip", c.ClientIP()),
	)

	// 调用服务层进行登出
	err := h.authService.Logout(tokenString)
	if err != nil {
		handleServiceError(c, err, requestID)
		return
	}

	logger.Info("用户登出成功",
		logger.String("request_id", requestID),
		logger.String("client_ip", c.ClientIP()),
	)

	utils.ResponseSuccess(c, "登出成功", nil)
}

// GetProfile 获取用户资料
// @Summary 获取用户资料
// @Description 获取当前登录用户的资料
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=models.UserResponse} "获取成功"
// @Failure 401 {object} utils.Response "未认证"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	requestID := middleware.GetTraceID(c)

	// 从上下文获取用户ID（由认证中间件设置）
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ResponseError(c, http.StatusUnauthorized, "未认证")
		return
	}

	userID := userIDInterface.(int64)

	logger.Info("获取用户资料请求",
		logger.String("request_id", requestID),
		logger.Int64("user_id", userID),
		logger.String("client_ip", c.ClientIP()),
	)

	// 调用用户服务获取用户信息
	user, err := h.userService.GetUserByID(int(userID))
	if err != nil {
		handleServiceError(c, err, requestID)
		return
	}

	logger.Info("获取用户资料成功",
		logger.String("request_id", requestID),
		logger.Int64("user_id", userID),
		logger.String("client_ip", c.ClientIP()),
	)

	utils.ResponseSuccess(c, "获取成功", user)
}

// RefreshToken 刷新访问令牌
// @Summary 刷新访问令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {refresh_token}"
// @Success 200 {object} utils.Response{data=map[string]string} "刷新成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 401 {object} utils.Response "令牌无效"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	requestID := middleware.GetTraceID(c)

	// 获取Authorization头
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.ResponseError(c, http.StatusBadRequest, "未提供刷新令牌")
		return
	}

	// 提取token（去掉Bearer前缀）
	refreshToken := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		refreshToken = authHeader[7:]
	}

	logger.Info("刷新令牌请求",
		logger.String("request_id", requestID),
		logger.String("client_ip", c.ClientIP()),
	)

	// 验证刷新令牌并生成新的访问令牌
	claims, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		handleServiceError(c, err, requestID)
		return
	}

	// 生成新的访问令牌
	newAccessToken, err := utils.GenerateAccessToken(claims.UserID, claims.Username)
	if err != nil {
		handleServiceError(c, err, requestID)
		return
	}

	logger.Info("令牌刷新成功",
		logger.String("request_id", requestID),
		logger.Int64("user_id", claims.UserID),
		logger.String("client_ip", c.ClientIP()),
	)

	response := map[string]interface{}{
		"access_token": newAccessToken,
		"token_type":   "Bearer",
		"expires_in":   3600, // 1小时，单位秒
	}

	utils.ResponseSuccess(c, "令牌刷新成功", response)
}
