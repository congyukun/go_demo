package middleware

import (
	"go_demo/internal/models"
	"go_demo/internal/utils"
	"go_demo/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PermissionMiddleware 权限检查中间件
func PermissionMiddleware(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			utils.ResponseError(c, http.StatusUnauthorized, "用户未认证")
			c.Abort()
			return
		}

		// 检查是否是管理员角色
		isAdmin := false
		if userResp, ok := user.(models.UserResponse); ok {
			for _, role := range userResp.Roles {
				if role == "admin" {
					isAdmin = true
					break
				}
			}
		}

		// 管理员拥有所有权限
		if isAdmin {
			c.Next()
			return
		}

		// 如果不是管理员，检查是否有特定权限
		// 这里简化处理，实际应该检查用户的具体权限
		utils.ResponseError(c, http.StatusForbidden, "权限不足")
		c.Abort()
	}
}

// RoleMiddleware 角色检查中间件
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			utils.ResponseError(c, http.StatusUnauthorized, "用户未认证")
			c.Abort()
			return
		}

		// 检查用户角色
		hasRole := false
		if userResp, ok := user.(models.UserResponse); ok {
			for _, userRole := range userResp.Roles {
				if userRole == requiredRole {
					hasRole = true
					break
				}
			}
		}

		if hasRole {
			c.Next()
			return
		}

		utils.ResponseError(c, http.StatusForbidden, "角色不足")
		c.Abort()
	}
}

// SelfOrAdminMiddleware 检查是否是本人或管理员的中间件
func SelfOrAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			utils.ResponseError(c, http.StatusUnauthorized, "用户未认证")
			c.Abort()
			return
		}

		// 获取路径中的用户ID
		userIDStr := c.Param("id")
		if userIDStr == "" {
			userIDStr = c.Param("user_id")
		}

		if userIDStr == "" {
			utils.ResponseError(c, http.StatusBadRequest, "未指定用户ID")
			c.Abort()
			return
		}

		targetUserID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			utils.ResponseError(c, http.StatusBadRequest, "用户ID格式错误")
			c.Abort()
			return
		}

		// 检查是否是管理员或本人
		isAdmin := false
		isSelf := false

		if userResp, ok := user.(models.UserResponse); ok {
			// 检查是否是本人
			if int64(userResp.ID) == targetUserID {
				isSelf = true
			}

			// 检查是否是管理员
			for _, role := range userResp.Roles {
				if role == "admin" {
					isAdmin = true
					break
				}
			}
		}

		if isAdmin || isSelf {
			c.Next()
			return
		}

		utils.ResponseError(c, http.StatusForbidden, "无权限操作此资源")
		c.Abort()
	}
}

// MultiResourcePermissionMiddleware 多资源权限检查中间件
func MultiResourcePermissionMiddleware(resourceType string, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			utils.ResponseError(c, http.StatusUnauthorized, "用户未认证")
			c.Abort()
			return
		}

		// 检查是否是管理员角色
		isAdmin := false
		if userResp, ok := user.(models.UserResponse); ok {
			for _, role := range userResp.Roles {
				if role == "admin" {
					isAdmin = true
					break
				}
			}
		}

		// 管理员拥有所有权限
		if isAdmin {
			c.Next()
			return
		}

		// 记录权限检查日志
		requestID := utils.GetRequestID(c)
		logger.Warn("权限检查失败",
			logger.String("request_id", requestID),
			logger.String("resource_type", resourceType),
			logger.String("action", action),
		)

		utils.ResponseError(c, http.StatusForbidden, "权限不足")
		c.Abort()
	}
}

// AdminRequiredMiddleware 管理员权限中间件
func AdminRequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			utils.ResponseError(c, http.StatusUnauthorized, "用户未认证")
			c.Abort()
			return
		}

		// 检查是否是管理员
		isAdmin := false
		if userResp, ok := user.(models.UserResponse); ok {
			for _, role := range userResp.Roles {
				if role == "admin" {
					isAdmin = true
					break
				}
			}
		}

		if isAdmin {
			c.Next()
			return
		}

		utils.ResponseError(c, http.StatusForbidden, "需要管理员权限")
		c.Abort()
	}
}

// GetUserFromContext 从上下文中获取用户信息
func GetUserFromContext(c *gin.Context) (models.UserResponse, bool) {
	user, exists := c.Get("user")
	if !exists {
		return models.UserResponse{}, false
	}

	if userResp, ok := user.(models.UserResponse); ok {
		return userResp, true
	}

	return models.UserResponse{}, false
}

// CheckUserPermission 检查用户是否有指定权限的辅助函数
func CheckUserPermission(user models.UserResponse, permission string) bool {
	// 检查是否是管理员
	for _, role := range user.Roles {
		if role == "admin" {
			return true
		}
	}

	// 简化的权限检查逻辑
	return false
}

// GetUserRoles 获取用户角色的辅助函数
func GetUserRoles(user models.UserResponse) []string {
	return user.Roles
}

// IsUserInRole 检查用户是否在指定角色中的辅助函数
func IsUserInRole(user models.UserResponse, roleName string) bool {
	for _, role := range user.Roles {
		if role == roleName {
			return true
		}
	}
	return false
}

// IsUserInAnyRole 检查用户是否在任意指定角色中的辅助函数
func IsUserInAnyRole(user models.UserResponse, roleNames ...string) bool {
	for _, userRole := range user.Roles {
		for _, roleName := range roleNames {
			if userRole == roleName {
				return true
			}
		}
	}
	return false
}