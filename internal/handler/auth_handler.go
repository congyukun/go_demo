package handler

import (
	"go_demo/internal/middleware"
	"go_demo/internal/models"
	"go_demo/internal/service"
	"go_demo/internal/utils"
	"go_demo/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

	logger.Info("用户登录请求",
		logger.String("request_id", requestID),
		logger.String("username", req.Username),
		logger.String("client_ip", c.ClientIP()),
	)

	// 调用服务层进行登录
	response, err := h.authService.Login(req)
	if err != nil {
		logger.Warn("用户登录失败",
			logger.String("request_id", requestID),
			logger.String("username", req.Username),
			logger.String("client_ip", c.ClientIP()),
			logger.Err(err),
		)
		utils.ResponseError(c, http.StatusUnauthorized, err.Error())
		return
	}

	logger.Info("用户登录成功",
		logger.String("request_id", requestID),
		logger.String("username", req.Username),
		logger.Int("user_id", int(response.User.ID)),
		logger.String("client_ip", c.ClientIP()),
	)

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
		logger.String("mobile", req.Mobile),
		logger.String("username", req.Username),
		logger.String("email", req.Email),
		logger.String("client_ip", c.ClientIP()),
	)

	// 调用服务层进行注册
	user, err := h.authService.Register(req)
	if err != nil {
		logger.Warn("用户注册失败",
			logger.String("request_id", requestID),
			logger.String("username", req.Username),
			logger.String("email", req.Email),
			logger.String("client_ip", c.ClientIP()),
			logger.Err(err),
		)

		// 根据错误类型返回不同的HTTP状态码
		if err.Error() == "用户名已存在" || err.Error() == "手机号已存在" {
			utils.ResponseError(c, http.StatusConflict, err.Error())
		} else {
			utils.ResponseError(c, http.StatusBadRequest, err.Error())
		}
		return
	}

	logger.Info("用户注册成功",
		logger.String("request_id", requestID),
		logger.String("username", req.Username),
		logger.String("mobile", req.Mobile),
		logger.Int("user_id", int(user.ID)),
		logger.String("client_ip", c.ClientIP()),
	)

	utils.ResponseSuccess(c, "注册成功", user)
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
	// 这里需要从刷新令牌中获取用户信息，然后生成新的访问令牌
	claims, err := utils.ValidateToken(refreshToken)
	if err != nil {
		logger.Warn("刷新令牌验证失败",
			logger.String("request_id", requestID),
			logger.String("client_ip", c.ClientIP()),
			logger.Err(err),
		)
		utils.ResponseError(c, http.StatusUnauthorized, "刷新令牌无效")
		return
	}

	// 验证这是一个刷新令牌（刷新令牌的role为空）
	if claims.Role != "" {
		utils.ResponseError(c, http.StatusUnauthorized, "无效的刷新令牌")
		return
	}

	// 生成新的访问令牌（这里需要从数据库获取用户当前角色）
	// 为了简化，这里假设用户角色为"user"，实际应该查询数据库
	newAccessToken, err := utils.RefreshAccessToken(refreshToken, "user")
	if err != nil {
		logger.Error("生成新访问令牌失败",
			logger.String("request_id", requestID),
			logger.String("client_ip", c.ClientIP()),
			logger.Err(err),
		)
		utils.ResponseError(c, http.StatusInternalServerError, "生成新令牌失败")
		return
	}

	logger.Info("令牌刷新成功",
		logger.String("request_id", requestID),
		logger.Int64("user_id", claims.UserID),
		logger.String("client_ip", c.ClientIP()),
	)

	response := map[string]string{
		"access_token": newAccessToken,
		"token_type":   "Bearer",
	}

	utils.ResponseSuccess(c, "令牌刷新成功", response)
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户登出接口（可选实现，主要用于记录日志）
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response "登出成功"
// @Failure 401 {object} utils.Response "未认证"
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	requestID := middleware.GetTraceID(c)

	// 从上下文获取用户信息（由认证中间件设置）
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ResponseError(c, http.StatusUnauthorized, "未认证")
		return
	}

	username, _ := c.Get("username")

	logger.Info("用户登出",
		logger.String("request_id", requestID),
		logger.Any("user_id", userID),
		logger.Any("username", username),
		logger.String("client_ip", c.ClientIP()),
	)

	// 在实际应用中，可以在这里：
	// 1. 将令牌加入黑名单
	// 2. 清除服务器端会话
	// 3. 记录登出日志等

	utils.ResponseSuccess(c, "登出成功", nil)
}

// GetProfile 获取当前用户信息
// @Summary 获取当前用户信息
// @Description 获取当前登录用户的详细信息
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

	// 从上下文获取用户信息（由认证中间件设置）
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ResponseError(c, http.StatusUnauthorized, "未认证")
		return
	}

	userID, ok := userIDInterface.(int64)
	if !ok {
		utils.ResponseError(c, http.StatusInternalServerError, "用户ID类型错误")
		return
	}

	logger.Info("获取用户信息请求",
		logger.String("request_id", requestID),
		logger.Int64("user_id", userID),
		logger.String("client_ip", c.ClientIP()),
	)

	// 调用用户服务获取用户详细信息
	user, err := h.userService.GetUserByID(int(userID))
	if err != nil {
		logger.Error("获取用户信息失败",
			logger.String("request_id", requestID),
			logger.Int64("user_id", userID),
			logger.Err(err),
		)
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Info("获取用户信息成功",
		logger.String("request_id", requestID),
		logger.Int64("user_id", userID),
	)

	utils.ResponseSuccess(c, "获取成功", user)
}
