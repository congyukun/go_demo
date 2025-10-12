package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"go_demo/internal/middleware"
	"go_demo/internal/models"
	"go_demo/internal/service"
	"go_demo/internal/utils"
	"go_demo/pkg/logger"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Users []models.UserResponse `json:"users"`
	Total int64                 `json:"total"`
	Page  int                   `json:"page"`
	Size  int                   `json:"size"`
}

// GetUsers 获取用户列表
// @Summary 获取用户列表
// @Description 分页获取用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} utils.Response{data=models.UserListResponse} "获取成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 401 {object} utils.Response "未认证"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	requestID := middleware.GetTraceID(c)

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	logger.Info("获取用户列表请求",
		logger.String("request_id", requestID),
		logger.Int("page", page),
		logger.Int("size", size),
		logger.String("client_ip", c.ClientIP()),
	)

	// 调用服务层获取用户列表
	users, total, err := h.userService.GetUsers(page, size)
	if err != nil {
		logger.Error("获取用户列表失败",
			logger.String("request_id", requestID),
			logger.Err(err),
		)
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 构造响应
	response := UserListResponse{
		Users: make([]models.UserResponse, len(users)),
		Total: total,
		Page:  page,
		Size:  size,
	}

	for i, user := range users {
		response.Users[i] = *user
	}

	logger.Info("获取用户列表成功",
		logger.String("request_id", requestID),
		logger.Int64("total", total),
		logger.Int("count", len(users)),
	)

	utils.ResponseSuccess(c, "获取成功", response)
}

// GetUser 获取用户详情
// @Summary 获取用户详情
// @Description 根据用户ID获取用户详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} utils.Response{data=models.UserResponse} "获取成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 401 {object} utils.Response "未认证"
// @Failure 404 {object} utils.Response "用户不存在"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	requestID := middleware.GetTraceID(c)

	// 获取用户ID参数
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	logger.Info("获取用户详情请求",
		logger.String("request_id", requestID),
		logger.Int("user_id", id),
		logger.String("client_ip", c.ClientIP()),
	)

	// 调用服务层获取用户信息
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		logger.Warn("获取用户详情失败",
			logger.String("request_id", requestID),
			logger.Int("user_id", id),
			logger.Err(err),
		)

		if err.Error() == "用户不存在" {
			utils.ResponseError(c, http.StatusNotFound, err.Error())
		} else {
			utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	logger.Info("获取用户详情成功",
		logger.String("request_id", requestID),
		logger.Int("user_id", id),
	)

	utils.ResponseSuccess(c, "获取成功", user)
}

// CreateUser 创建用户
// @Summary 创建用户
// @Description 创建新用户（管理员功能）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.UserCreateRequest true "创建用户请求"
// @Success 200 {object} utils.Response{data=models.UserResponse} "创建成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 401 {object} utils.Response "未认证"
// @Failure 403 {object} utils.Response "权限不足"
// @Failure 409 {object} utils.Response "用户已存在"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	requestID := middleware.GetTraceID(c)

	// 绑定并验证请求参数
	var req models.UserCreateRequest
	if !middleware.ValidateAndBind(c, &req) {
		return
	}

	logger.Info("创建用户请求",
		logger.String("request_id", requestID),
		logger.String("username", req.Username),
		logger.String("email", req.Email),
		logger.String("client_ip", c.ClientIP()),
	)

	// 调用服务层创建用户
	user, err := h.userService.CreateUser(req)
	if err != nil {
		logger.Warn("创建用户失败",
			logger.String("request_id", requestID),
			logger.String("username", req.Username),
			logger.Err(err),
		)

		// 根据错误类型返回不同的HTTP状态码
		if err.Error() == "用户名已存在" || err.Error() == "邮箱已存在" || err.Error() == "手机号已存在" {
			utils.ResponseError(c, http.StatusConflict, err.Error())
		} else {
			utils.ResponseError(c, http.StatusBadRequest, err.Error())
		}
		return
	}

	logger.Info("创建用户成功",
		logger.String("request_id", requestID),
		logger.String("username", req.Username),
		logger.Int("user_id", int(user.ID)),
	)

	utils.ResponseSuccess(c, "创建成功", user)
}

// UpdateUser 更新用户
// @Summary 更新用户
// @Description 更新用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Param request body models.UpdateUserRequest true "更新用户请求"
// @Success 200 {object} utils.Response{data=models.UserResponse} "更新成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 401 {object} utils.Response "未认证"
// @Failure 404 {object} utils.Response "用户不存在"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	requestID := middleware.GetTraceID(c)

	// 获取用户ID参数
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	// 绑定并验证请求参数
	var req models.UpdateUserRequest
	if !middleware.ValidateAndBind(c, &req) {
		return
	}

	logger.Info("更新用户请求",
		logger.String("request_id", requestID),
		logger.Int("user_id", id),
		logger.String("client_ip", c.ClientIP()),
	)

	// 调用服务层更新用户
	user, err := h.userService.UpdateUser(id, req)
	if err != nil {
		logger.Warn("更新用户失败",
			logger.String("request_id", requestID),
			logger.Int("user_id", id),
			logger.Err(err),
		)

		if err.Error() == "用户不存在" {
			utils.ResponseError(c, http.StatusNotFound, err.Error())
		} else {
			utils.ResponseError(c, http.StatusBadRequest, err.Error())
		}
		return
	}

	logger.Info("更新用户成功",
		logger.String("request_id", requestID),
		logger.Int("user_id", id),
	)

	utils.ResponseSuccess(c, "更新成功", user)
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 删除用户（软删除）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} utils.Response "删除成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 401 {object} utils.Response "未认证"
// @Failure 404 {object} utils.Response "用户不存在"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	requestID := middleware.GetTraceID(c)

	// 获取用户ID参数
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	logger.Info("删除用户请求",
		logger.String("request_id", requestID),
		logger.Int("user_id", id),
		logger.String("client_ip", c.ClientIP()),
	)

	// 调用服务层删除用户
	err = h.userService.DeleteUser(id)
	if err != nil {
		logger.Warn("删除用户失败",
			logger.String("request_id", requestID),
			logger.Int("user_id", id),
			logger.Err(err),
		)

		if err.Error() == "用户不存在" {
			utils.ResponseError(c, http.StatusNotFound, err.Error())
		} else {
			utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	logger.Info("删除用户成功",
		logger.String("request_id", requestID),
		logger.Int("user_id", id),
	)

	utils.ResponseSuccess(c, "删除成功", nil)
}

// UpdateProfile 更新当前用户资料
// @Summary 更新当前用户资料
// @Description 用户更新自己的资料信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.UserProfileUpdateRequest true "更新资料请求"
// @Success 200 {object} utils.Response{data=models.UserResponse} "更新成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 401 {object} utils.Response "未认证"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	requestID := middleware.GetTraceID(c)

	// 从上下文获取用户ID（由认证中间件设置）
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

	// 绑定并验证请求参数
	var req models.UserProfileUpdateRequest
	if !middleware.ValidateAndBind(c, &req) {
		return
	}

	logger.Info("更新用户资料请求",
		logger.String("request_id", requestID),
		logger.Int64("user_id", userID),
		logger.String("client_ip", c.ClientIP()),
	)

	// 调用服务层更新用户资料
	user, err := h.userService.UpdateUserProfile(int(userID), req)
	if err != nil {
		logger.Warn("更新用户资料失败",
			logger.String("request_id", requestID),
			logger.Int64("user_id", userID),
			logger.Err(err),
		)
		utils.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Info("更新用户资料成功",
		logger.String("request_id", requestID),
		logger.Int64("user_id", userID),
	)

	utils.ResponseSuccess(c, "更新成功", user)
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 用户修改自己的密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ChangePasswordRequest true "修改密码请求"
// @Success 200 {object} utils.Response "修改成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 401 {object} utils.Response "未认证"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/users/Password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	requestID := middleware.GetTraceID(c)

	// 从上下文获取用户ID（由认证中间件设置）
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

	// 绑定并验证请求参数
	var req models.ChangePasswordRequest
	if !middleware.ValidateAndBind(c, &req) {
		return
	}

	logger.Info("修改密码请求",
		logger.String("request_id", requestID),
		logger.Int64("user_id", userID),
		logger.String("client_ip", c.ClientIP()),
	)

	// 调用服务层修改密码
	err := h.userService.ChangePassword(int(userID), req)
	if err != nil {
		logger.Warn("修改密码失败",
			logger.String("request_id", requestID),
			logger.Int64("user_id", userID),
			logger.Err(err),
		)
		utils.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Info("修改密码成功",
		logger.String("request_id", requestID),
		logger.Int64("user_id", userID),
	)

	utils.ResponseSuccess(c, "修改成功", nil)
}

// GetUserStats 获取用户统计信息
// @Summary 获取用户统计信息
// @Description 获取用户统计信息（管理员功能）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=map[string]interface{}} "获取成功"
// @Failure 401 {object} utils.Response "未认证"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/users/stats [get]
func (h *UserHandler) GetUserStats(c *gin.Context) {
	requestID := middleware.GetTraceID(c)

	logger.Info("获取用户统计信息请求",
		logger.String("request_id", requestID),
		logger.String("client_ip", c.ClientIP()),
	)

	// 调用服务层获取统计信息
	stats, err := h.userService.GetUserStats()
	if err != nil {
		logger.Error("获取用户统计信息失败",
			logger.String("request_id", requestID),
			logger.Err(err),
		)
		utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Info("获取用户统计信息成功",
		logger.String("request_id", requestID),
	)

	utils.ResponseSuccess(c, "获取成功", stats)
}

// GetUserRoles 获取用户角色
// @Summary 获取用户角色
// @Description 获取指定用户的角色列表（管理员功能）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id path int true "用户ID"
// @Success 200 {object} utils.Response{data=[]models.Role} "获取成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 401 {object} utils.Response "未认证"
// @Failure 403 {object} utils.Response "权限不足"
// @Failure 404 {object} utils.Response "用户不存在"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/users/{user_id}/roles [get]
func (h *UserHandler) GetUserRoles(c *gin.Context) {
	requestID := middleware.GetTraceID(c)

	// 获取用户ID参数
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	logger.Info("获取用户角色请求",
		logger.String("request_id", requestID),
		logger.Int64("user_id", userID),
		logger.String("client_ip", c.ClientIP()),
	)

	// 调用服务层获取用户角色
	roles, err := h.userService.GetUserRoles(userID)
	if err != nil {
		logger.Warn("获取用户角色失败",
			logger.String("request_id", requestID),
			logger.Int64("user_id", userID),
			logger.Err(err),
		)

		// 检查是否是记录未找到错误
		if strings.Contains(err.Error(), "record not found") || strings.Contains(err.Error(), "not found") {
			utils.ResponseError(c, http.StatusNotFound, "用户不存在")
		} else {
			utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	logger.Info("获取用户角色成功",
		logger.String("request_id", requestID),
		logger.Int64("user_id", userID),
		logger.Int("role_count", len(roles)),
	)

	utils.ResponseSuccess(c, "获取成功", roles)
}

// AssignRole 分配角色
// @Summary 分配角色
// @Description 为用户分配角色（管理员功能）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id path int true "用户ID"
// @Param request body models.AssignRoleRequest true "分配角色请求"
// @Success 200 {object} utils.Response "分配成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 401 {object} utils.Response "未认证"
// @Failure 403 {object} utils.Response "权限不足"
// @Failure 404 {object} utils.Response "用户或角色不存在"
// @Failure 409 {object} utils.Response "用户已拥有该角色"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/users/{user_id}/roles [post]
func (h *UserHandler) AssignRole(c *gin.Context) {
	requestID := middleware.GetTraceID(c)

	// 获取用户ID参数
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	// 绑定并验证请求参数
	var req models.AssignRoleRequest
	if !middleware.ValidateAndBind(c, &req) {
		return
	}

	logger.Info("分配角色请求",
		logger.String("request_id", requestID),
		logger.Int64("user_id", userID),
		logger.Strings("roles", req.Roles),
		logger.String("client_ip", c.ClientIP()),
	)

	// 分配角色，暂时取第一个角色
	roleName := ""
	if len(req.Roles) > 0 {
		roleName = req.Roles[0]
	}
	err = h.userService.AssignRole(userID, roleName)
	if err != nil {
		logger.Warn("分配角色失败",
			logger.String("request_id", requestID),
			logger.Int64("user_id", userID),
			logger.String("role_name", roleName),
			logger.Err(err),
		)

		// 检查错误类型
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "不存在") {
			utils.ResponseError(c, http.StatusNotFound, "角色不存在")
		} else if strings.Contains(err.Error(), "conflict") || strings.Contains(err.Error(), "已拥有") {
			utils.ResponseError(c, http.StatusConflict, "用户已拥有该角色")
		} else {
			utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	logger.Info("分配角色成功",
		logger.String("request_id", requestID),
		logger.Int64("user_id", userID),
		logger.String("role_name", roleName),
	)

	utils.ResponseSuccess(c, "分配成功", nil)
}

// RevokeRole 撤销角色
// @Summary 撤销角色
// @Description 撤销用户的角色（管理员功能）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id path int true "用户ID"
// @Param request body models.RevokeRoleRequest true "撤销角色请求"
// @Success 200 {object} utils.Response "撤销成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 401 {object} utils.Response "未认证"
// @Failure 403 {object} utils.Response "权限不足"
// @Failure 404 {object} utils.Response "用户或角色不存在"
// @Failure 409 {object} utils.Response "用户未拥有该角色"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/users/{user_id}/roles/revoke [post]
func (h *UserHandler) RevokeRole(c *gin.Context) {
	requestID := middleware.GetTraceID(c)

	// 获取用户ID参数
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	// 绑定并验证请求参数
	var req models.RevokeRoleRequest
	if !middleware.ValidateAndBind(c, &req) {
		return
	}

	logger.Info("撤销角色请求",
		logger.String("request_id", requestID),
		logger.Int64("user_id", userID),
		logger.Strings("roles", req.Roles),
		logger.String("client_ip", c.ClientIP()),
	)

	// 撤销角色，暂时取第一个角色
	roleName := ""
	if len(req.Roles) > 0 {
		roleName = req.Roles[0]
	}
	err = h.userService.RevokeRole(userID, roleName)
	if err != nil {
		logger.Warn("撤销角色失败",
			logger.String("request_id", requestID),
			logger.Int64("user_id", userID),
			logger.String("role_name", roleName),
			logger.Err(err),
		)

		// 检查错误类型
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "不存在") {
			utils.ResponseError(c, http.StatusNotFound, "角色不存在")
		} else if strings.Contains(err.Error(), "conflict") || strings.Contains(err.Error(), "未拥有") {
			utils.ResponseError(c, http.StatusConflict, "用户未拥有该角色")
		} else {
			utils.ResponseError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	logger.Info("撤销角色成功",
		logger.String("request_id", requestID),
		logger.Int64("user_id", userID),
		logger.String("role_name", roleName),
	)

	utils.ResponseSuccess(c, "撤销成功", nil)
}
