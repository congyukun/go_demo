// Package di 负责依赖注入各项 Provider 的定义与初始化。
// 设计要点：
// 1) ProvideConfig(string) (*config.Config, error) —— 加载配置，所有 Provider 的根依赖。
// 2) ProvideAppInit(cfg) (AppInit, error) —— 统一完成副作用初始化：logger.Init(cfg.Log)、utils.InitJWT(cfg.JWT)、validator.Init()，并根据 cfg.Server.Mode 设置 Gin 模式。
// 3) ProvideDB(cfg) (*gorm.DB, error) —— 创建数据库连接。
// 4) ProvideCache(cfg) (cache.CacheInterface, error) —— 优先 Redis，失败回退内存缓存。
// 5) ProvideRepositories(db) repository.UserRepository —— 组合仓储层（当前为 UserRepository）。
// 6) ProvideAuthService(userRepo) / ProvideUserService(userRepo) —— 构建业务服务。
// 7) ProvideAuthHandler(authSvc, userSvc) / ProvideUserHandler(userSvc) —— 构建 Handler。
// 8) ProvideRateLimiterFactory(cfg, cache) / ProvideCircuitBreakerFactory(cfg) —— 中间件工厂。
// 9) ProvideRouter(authHandler, userHandler, rl, cb) *router.Router —— 构建 Router，集中注册中间件与路由分组。
// 10) ProvideGinEngine(r *router.Router) *gin.Engine —— 调用 r.Setup() 返回最终可用的引擎。
// Wire 在 internal/di/wire.go 中声明 Sets，并在 wire_gen.go 中生成 InitializeServer 注入器调用这些 Provider。
package di

import (
	"fmt"
	"go_demo/internal/config"
	"go_demo/internal/handler"
	"go_demo/internal/middleware"
	"go_demo/internal/repository"
	"go_demo/internal/router"
	"go_demo/internal/service"
	"go_demo/internal/utils"
	"go_demo/pkg/cache"
	"go_demo/pkg/database"
	"go_demo/pkg/logger"
	"go_demo/pkg/validator"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 统一缓存策略说明：
// - 项目层面统一采用 Redis 作为缓存实现，启动阶段若 Redis 不可用则启动失败（不再回退至内存）。
// - 限流中间件统一使用分布式策略（依赖 Redis），全局采用固定窗口算法，用户级采用滑动窗口算法。
// - 熔断器仍为进程内实现，但其日志与配置通过 DI 注入统一管理。

// ProvideConfig 加载配置
func ProvideConfig(configPath string) (*config.Config, error) { // di.ProvideConfig()
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}
	return cfg, nil
}

// AppInit 为需要副作用初始化的组件提供统一入口
type AppInit struct{}

// ProvideAppInit 初始化日志、JWT、验证器等副作用组件
func ProvideAppInit(cfg *config.Config) (AppInit, error) { // di.ProvideAppInit()
	logConfig := cfg.Log
	if err := logger.Init(logConfig); err != nil {
		return AppInit{}, err
	}
	// 初始化JWT
	utils.InitJWT(cfg.JWT)
	// 初始化验证器
	if err := validator.Init(); err != nil {
		return AppInit{}, err
	}
	// 设置Gin模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	return AppInit{}, nil
}

// ProvideDB 初始化数据库
func ProvideDB(cfg *config.Config) (*gorm.DB, error) { // di.ProvideDB()
	db, err := database.NewMySQL(cfg.Database)
	if err != nil {
		return nil, err
	}
	logger.Info("MySQL数据库初始化成功", logger.String("addr", cfg.Database.DSN))
	return db, nil
}
// ProvideCache 初始化缓存（统一采用 Redis，不再使用内存回退）
func ProvideCache(cfg *config.Config) (cache.CacheInterface, error) { // di.ProvideCache()
	redisCfg := cache.RedisConfig{
		Addr:         fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
		MaxRetries:   cfg.Redis.MaxRetries,
	}
	redisCache, err := cache.NewRedisCache(redisCfg)
	if err != nil {
		// 统一使用 Redis，直接返回错误以便启动阶段显式失败，避免静默降级
		return nil, fmt.Errorf("redis缓存初始化失败: %w", err)
	}
	logger.Info("Redis缓存初始化成功", logger.String("addr", redisCfg.Addr))
	return redisCache, nil
}

// ProvideCacheWithInit 包装缓存 Provider，引入对 appReady 的依赖，保证在日志/JWT/校验器初始化后再创建缓存
func ProvideCacheWithInit(_ appReady, cfg *config.Config) (cache.CacheInterface, error) { // di.ProvideCacheWithInit()
	return ProvideCache(cfg)
}

// ProvideRepositories 初始化仓储
func ProvideRepositories(db *gorm.DB) repository.UserRepository { // di.ProvideRepositories()
	return repository.NewUserRepository(db)
}

// ProvideAuthService 初始化认证服务
func ProvideAuthService(userRepo repository.UserRepository) service.AuthService { // di.ProvideAuthService()
	return service.NewAuthService(userRepo)
}

// ProvideUserService 初始化用户服务
func ProvideUserService(userRepo repository.UserRepository) service.UserService { // di.ProvideUserService()
	return service.NewUserService(userRepo)
}

// ProvideRateLimiterFactory 初始化限流器工厂（优化：统一分布式依赖Redis，全局采用固定窗口，用户采用滑动窗口）
func ProvideRateLimiterFactory(cfg *config.Config, c cache.CacheInterface) *middleware.RateLimiterFactory { // di.ProvideRateLimiterFactory()
	global := middleware.DefaultRateLimiterConfig()
	global.MaxRequests = cfg.RateLimiter.GlobalLimit
	global.Window = time.Duration(cfg.RateLimiter.Window) * time.Second
	// 全局限流更适合固定窗口，便于统一控制峰值
	global.Algorithm = "fixed_window"
	global.Distributed = true
	global.Cache = c

	userCfg := middleware.DefaultRateLimiterConfig()
	userCfg.MaxRequests = cfg.RateLimiter.UserLimit
	userCfg.Window = time.Duration(cfg.RateLimiter.Window) * time.Second
	// 用户限流采用滑动窗口，提升平滑度与用户体验
	userCfg.Algorithm = "sliding_window"
	userCfg.Distributed = true
	userCfg.Cache = c

	return middleware.NewRateLimiterFactory(global, userCfg, c)
}

// ProvideCircuitBreakerFactory 初始化熔断器工厂
func ProvideCircuitBreakerFactory(cfg *config.Config) *middleware.CircuitBreakerFactory { // di.ProvideCircuitBreakerFactory()
	cbCfg := middleware.DefaultCircuitBreakerConfig("global")
	cbCfg.MaxRequests = uint32(cfg.CircuitBreaker.MaxRequests)
	cbCfg.MaxHalfOpenRequests = uint32(cfg.CircuitBreaker.HalfOpenMaxRequests)
	cbCfg.ErrorThreshold = cfg.CircuitBreaker.ErrorThreshold
	cbCfg.Timeout = time.Duration(cfg.CircuitBreaker.Timeout) * time.Second
	cbCfg.Enabled = cfg.CircuitBreaker.Enabled
	return middleware.NewCircuitBreakerFactory(cbCfg)
}

// ProvideAuthHandler 初始化 AuthHandler
func ProvideAuthHandler(authSvc service.AuthService, userSvc service.UserService) *handler.AuthHandler { // di.ProvideAuthHandler()
	return handler.NewAuthHandler(authSvc, userSvc)
}

// ProvideUserHandler 初始化 UserHandler
func ProvideUserHandler(userSvc service.UserService) *handler.UserHandler { // di.ProvideUserHandler()
	return handler.NewUserHandler(userSvc)
}

// ProvideRouter 构建 Router 并注册路由
func ProvideRouter(authHandler *handler.AuthHandler, userHandler *handler.UserHandler, rl *middleware.RateLimiterFactory, cb *middleware.CircuitBreakerFactory) *router.Router { // di.ProvideRouter()
	return router.NewRouterWithMiddleware(authHandler, userHandler, rl, cb)
}

// ProvideGinEngine 创建 Gin Engine（调用 Router.Setup）
func ProvideGinEngine(r *router.Router) *gin.Engine { // di.ProvideGinEngine()
	return r.Setup()
}

// ProvideAppReady 作为中间依赖，强制在构造链路中执行 ProvideAppInit
// 不改变对外签名，通过在 wire.go 中将其纳入依赖图，保证 AppInit 在 Cache/DB 使用前完成
type appReady struct{}

func ProvideAppReady(_ AppInit) appReady { // di.ProvideAppReady()
	return appReady{}
}

// ProvideGinEngineWithInit 在不改变原有签名的前提下，引入对 appReady 的依赖，确保初始化顺序
func ProvideGinEngineWithInit(_ appReady, r *router.Router) *gin.Engine { // di.ProvideGinEngineWithInit()
	return r.Setup()
}

// ProvideCleanup 组合资源清理函数
func ProvideCleanup(db *gorm.DB, c cache.CacheInterface) func() { // di.ProvideCleanup()
	return func() {
		// 关闭缓存连接
		if c != nil {
			if closer, ok := c.(interface{ Close() error }); ok {
				if err := closer.Close(); err != nil {
					logger.Error("关闭缓存连接失败", logger.Err(err))
				}
			}
		}

		// 关闭数据库连接
		if db != nil {
			if err := database.Close(db); err != nil {
				logger.Error("关闭数据库连接失败", logger.Err(err))
			}
		}

		// 同步日志
		logger.Sync()
	}
}
