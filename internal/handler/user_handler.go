package handler

import (
	"go_demo/internal/models"
	"go_demo/internal/service"
	"go_demo/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

// GetUsers 获取用户列表
func (h *UserHandler) GetUsers(c *gin.Context) {
	requestID := GetRequestID(c)

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	logger.Info("获取用户列表请求",
		logger.String("request_id", requestID),
		logger.Int("page", page),
		logger.Int("page_size", pageSize),
		logger.String("client_ip", c.ClientIP()),
	)

	// 调用服务层
	users, total, err := h.userService.GetUsers(page, pageSize)
	if err != nil {
		logger.Error("获取用户列表失败",
			logger.String("request_id", requestID),
			logger.Err(err),
		)
		ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 构造响应数据
	data := gin.H{
		"users": users,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	}

	logger.Info("获取用户列表成功",
		logger.String("request_id", requestID),
		logger.Int64("total", total),
	)

	ResponseSuccess(c, "获取用户列表成功", data)
}

// GetUserByID 根据ID获取用户
func (h *UserHandler) GetUserByID(c *gin.Context) {
	requestID := GetRequestID(c)

	// 获取用户ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Warn("用户ID参数错误",
			logger.String("request_id", requestID),
			logger.String("id", idStr),
			logger.Err(err),
		)
		ResponseError(c, http.StatusBadRequest, "用户ID参数错误")
		return
	}

	logger.Info("获取用户详情请求",
		logger.String("request_id", requestID),
		logger.Int("user_id", id),
		logger.String("client_ip", c.ClientIP()),
	)

	// 调用服务层
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		logger.Warn("获取用户详情失败",
			logger.String("request_id", requestID),
			logger.Int("user_id", id),
			logger.Err(err),
		)
		if err.Error() == "用户不存在" {
			ResponseError(c, http.StatusNotFound, err.Error())
		} else {
			ResponseError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	logger.Info("获取用户详情成功",
		logger.String("request_id", requestID),
		logger.Int("user_id", id),
	)

	ResponseSuccess(c, "获取用户详情成功", user)
}

// UpdateUser 更新用户
func (h *UserHandler) UpdateUser(c *gin.Context) {
	requestID := GetRequestID(c)

	// 获取用户ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Warn("用户ID参数错误",
			logger.String("request_id", requestID),
			logger.String("id", idStr),
			logger.Err(err),
		)
		ResponseError(c, http.StatusBadRequest, "用户ID参数错误")
		return
	}

	// 绑定并验证请求参数
	var req models.UpdateUserRequest
	if !ValidateAndBind(c, &req) {
		logger.Warn("更新用户参数验证失败",
			logger.String("request_id", requestID),
			logger.Int("user_id", id),
		)
		return
	}

	logger.Info("更新用户请求",
		logger.String("request_id", requestID),
		logger.Int("user_id", id),
		logger.String("client_ip", c.ClientIP()),
	)

	// 调用服务层
	user, err := h.userService.UpdateUser(id, req)
	if err != nil {
		logger.Warn("更新用户失败",
			logger.String("request_id", requestID),
			logger.Int("user_id", id),
			logger.Err(err),
		)
		if err.Error() == "用户不存在" {
			ResponseError(c, http.StatusNotFound, err.Error())
		} else {
			ResponseError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	logger.Info("更新用户成功",
		logger.String("request_id", requestID),
		logger.Int("user_id", id),
	)

	ResponseSuccess(c, "更新用户成功", user)
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(c *gin.Context) {
	requestID := GetRequestID(c)

	// 获取用户ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Warn("用户ID参数错误",
			logger.String("request_id", requestID),
			logger.String("id", idStr),
			logger.Err(err),
		)
		ResponseError(c, http.StatusBadRequest, "用户ID参数错误")
		return
	}

	logger.Info("删除用户请求",
		logger.String("request_id", requestID),
		logger.Int("user_id", id),
		logger.String("client_ip", c.ClientIP()),
	)

	// 调用服务层
	err = h.userService.DeleteUser(id)
	if err != nil {
		logger.Warn("删除用户失败",
			logger.String("request_id", requestID),
			logger.Int("user_id", id),
			logger.Err(err),
		)
		if err.Error() == "用户不存在" {
			ResponseError(c, http.StatusNotFound, err.Error())
		} else {
			ResponseError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	logger.Info("删除用户成功",
		logger.String("request_id", requestID),
		logger.Int("user_id", id),
	)

	ResponseSuccess(c, "删除用户成功", nil)
}
