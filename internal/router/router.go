package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go_demo/docs"
	"go_demo/internal/handler"
	"go_demo/internal/middleware"
// 移除未使用的导入 "go_demo/internal/utils"
)

// RateLimiter 接口定义
type RateLimiter interface {
	GlobalMiddleware() gin.HandlerFunc
	UserMiddleware() gin.HandlerFunc
}

// CircuitBreaker 接口定义
type CircuitBreaker interface {
	Middleware() gin.HandlerFunc
}

// Router 路由管理器
type Router struct {
	engine         *gin.Engine
	authHandler    *handler.AuthHandler
	userHandler    *handler.UserHandler
	rateLimiter    RateLimiter
	circuitBreaker CircuitBreaker
}

// NewRouter 创建新的路由管理器
func NewRouter(authHandler *handler.AuthHandler, userHandler *handler.UserHandler) *Router {
	return &Router{
		authHandler: authHandler,
		userHandler: userHandler,
	}
}

// NewRouterWithMiddleware 使用中间件创建新的路由管理器
func NewRouterWithMiddleware(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	rateLimiter RateLimiter,
	circuitBreaker CircuitBreaker,
) *Router {
	return &Router{
		authHandler:    authHandler,
		userHandler:    userHandler,
		rateLimiter:    rateLimiter,
		circuitBreaker: circuitBreaker,
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
	r.engine.Use(
		middleware.Logger(),   // 自定义日志中间件
		middleware.Recovery(), // 自定义恢复中间件
	)

	// 自定义中间件
	r.engine.Use(
		middleware.RequestID(),  // 请求ID，用于追踪请求
		middleware.CORS(),       // 跨域资源共享支持
		middleware.Trace(),      // 链路追踪
		middleware.RequestLog(), // 请求日志记录
	)

	// 添加限流和熔断中间件
	// 全局限流中间件
	if r.rateLimiter != nil {
		r.engine.Use(r.rateLimiter.GlobalMiddleware())
	}

	// 全局熔断中间件
	if r.circuitBreaker != nil {
		r.engine.Use(r.circuitBreaker.Middleware())
	}
}

// setupRoutes 设置路由
func (r *Router) setupRoutes() {
	// 健康检查
	r.setupHealthRoutes()

	// Swagger文档路由
	r.setupSwaggerRoutes()

	// API 路由
	r.setupAPIRoutes()
}

// setupHealthRoutes 设置健康检查路由
func (r *Router) setupHealthRoutes() {
	// 健康检查路由
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})
}

// RouteGroup 定义路由组接口
type RouteGroup interface {
	Group(string, ...gin.HandlerFunc) *gin.RouterGroup
}

// setupAPIRoutes 设置 API 路由
func (r *Router) setupAPIRoutes() {
	// API v1 路由组
	v1 := r.engine.Group("/api/v1")

	// 认证路由
	r.setupAuthRoutes(v1)

	// 用户路由
	r.setupUserRoutes(v1)
	
	// 可以在这里添加更多的路由组
	// 例如：r.setupArticleRoutes(v1) 等
}

// setupAuthRoutes 设置认证路由
func (r *Router) setupAuthRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		// 为认证路由添加用户级限流
		if r.rateLimiter != nil {
			auth.Use(r.rateLimiter.UserMiddleware())
		}

		// 公开路由（不需要认证）
		auth.POST("/login", r.authHandler.Login)
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/refresh", r.authHandler.RefreshToken)
		
		// 需要认证的路由
		auth.POST("/logout", middleware.JWTAuthMiddleware(), r.authHandler.Logout)
		auth.GET("/profile", middleware.JWTAuthMiddleware(), r.authHandler.GetProfile)

	}
}

// setupUserRoutes 设置用户路由
func (r *Router) setupUserRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	users.Use(middleware.JWTAuthMiddleware()) // 用户相关接口需要JWT认证

	// 为用户路由添加用户级限流
	if r.rateLimiter != nil {
		users.Use(r.rateLimiter.UserMiddleware())
	}

	{
		// 用户管理（需要认证）
		users.GET("", r.userHandler.GetUsers)
		users.POST("", r.userHandler.CreateUser)
		users.GET("/stats", r.userHandler.GetUserStats)

		// 用户详情和操作（需要认证）
		users.GET("/:id", r.userHandler.GetUser)
		users.PUT("/:id", r.userHandler.UpdateUser)
		users.DELETE("/:id", r.userHandler.DeleteUser)

		// 用户自己的操作
		users.PUT("/profile", r.userHandler.UpdateProfile)
		users.PUT("/password", r.userHandler.ChangePassword)
	}
}

// setupSwaggerRoutes 设置Swagger文档路由
func (r *Router) setupSwaggerRoutes() {
	// 导入docs包以确保它被使用
	_ = docs.SwaggerInfo
	// Swagger文档路由
	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// RegisterRoutes 注册新的路由
// 该方法允许在运行时动态添加路由
func (r *Router) RegisterRoutes(registerFunc func(*gin.Engine)) {
	if r.engine != nil && registerFunc != nil {
		registerFunc(r.engine)
	}
}

// GetEngine 获取 Gin 引擎
func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}
