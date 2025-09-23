package router

import (
	"go_demo/internal/handler"
	"go_demo/internal/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Router 路由管理器
type Router struct {
	engine      *gin.Engine
	authHandler *handler.AuthHandler
	userHandler *handler.UserHandler
}

// NewRouter 创建新的路由管理器
func NewRouter(authHandler *handler.AuthHandler, userHandler *handler.UserHandler) *Router {
	return &Router{
		authHandler: authHandler,
		userHandler: userHandler,
	}
}

// Setup 设置路由
func (r *Router) Setup() *gin.Engine {
	r.engine = gin.New()

	// 注册中间件
	r.setupMiddleware()

	// 注册路由
	r.setupRoutes()

	return r.engine
}

// setupMiddleware 设置中间件
func (r *Router) setupMiddleware() {
	// 基础中间件
	r.engine.Use(gin.Logger())
	r.engine.Use(gin.Recovery())
	
	// 自定义中间件
	r.engine.Use(middleware.RequestID())
	r.engine.Use(middleware.CORS())
	r.engine.Use(middleware.Trace()) // 添加链路追踪中间件
	r.engine.Use(middleware.RequestLog()) // 添加请求日志中间件
}

// setupRoutes 设置路由
func (r *Router) setupRoutes() {
	// 健康检查
	r.setupHealthRoutes()
	
	// API 路由
	r.setupAPIRoutes()
}

// setupHealthRoutes 设置健康检查路由
func (r *Router) setupHealthRoutes() {
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})
}

// setupAPIRoutes 设置 API 路由
func (r *Router) setupAPIRoutes() {
	// API v1 路由组
	v1 := r.engine.Group("/api/v1")
	
	// 认证路由
	r.setupAuthRoutes(v1)
	
	// 用户路由
	r.setupUserRoutes(v1)
}

// setupAuthRoutes 设置认证路由
func (r *Router) setupAuthRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", r.authHandler.Login)
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/logout", middleware.Auth(), r.authHandler.Logout)
	}
}

// setupUserRoutes 设置用户路由
func (r *Router) setupUserRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	users.Use(middleware.Auth()) // 用户相关接口需要JWT认证
	{
		users.GET("", r.userHandler.GetUsers)
		users.GET("/:id", r.userHandler.GetUserByID)
		users.PUT("/:id", r.userHandler.UpdateUser)
		users.DELETE("/:id", r.userHandler.DeleteUser)
	}
}

// GetEngine 获取 Gin 引擎
func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}