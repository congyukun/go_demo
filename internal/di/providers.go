// Package di 简洁的依赖注入实现
// 设计原则：
// 1. 使用聚合器模式减少 Provider 数量和参数复杂度
// 2. 保持清晰的分层结构：Config -> Infrastructure -> Business -> Presentation
// 3. 统一初始化顺序，确保依赖关系正确
package di

import (
	"fmt"
	"go_demo/internal/config"
	"go_demo/internal/router"
	"go_demo/internal/utils"
	"go_demo/pkg/cache"
	"go_demo/pkg/database"
	"go_demo/pkg/logger"
	"go_demo/pkg/validator"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ===== 基础设施层 =====

// ProvideConfig 加载配置 // di.ProvideConfig()
func ProvideConfig(configPath string) (*config.Config, error) {
	// 如果未指定，默认使用 ./configs/config.yaml
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}
	return cfg, nil
}

// AppInit 应用初始化标记 // di.AppInit
type AppInit struct{}

// ProvideAppInit 初始化应用基础组件 // di.ProvideAppInit()
func ProvideAppInit(cfg *config.Config) (AppInit, error) {
	// 初始化日志
	if err := logger.Init(cfg.Log); err != nil {
		return AppInit{}, fmt.Errorf("日志初始化失败: %w", err)
	}

	// 初始化JWT
	utils.InitJWT(cfg.JWT)

	// 初始化验证器
	if err := validator.Init(); err != nil {
		return AppInit{}, fmt.Errorf("验证器初始化失败: %w", err)
	}

	// 设置Gin模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	return AppInit{}, nil
}

// ProvideDB 初始化数据库 // di.ProvideDB()
func ProvideDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := database.NewMySQL(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("数据库初始化失败: %w", err)
	}
	logger.Info("MySQL数据库初始化成功", logger.String("addr", cfg.Database.DSN))
	return db, nil
}

// ProvideCache 初始化缓存 // di.ProvideCache()
func ProvideCache(cfg *config.Config) (cache.CacheInterface, error) {
	redisCfg := cache.RedisConfig{
		Host:         cfg.Redis.Host,
		Port:         cfg.Redis.Port,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
		MaxRetries:   cfg.Redis.MaxRetries,
	}

	redisCache, err := cache.NewRedisCache(redisCfg)
	if err != nil {
		return nil, fmt.Errorf("redis缓存初始化失败: %w", err)
	}

	logger.Info("Redis缓存初始化成功", logger.String("addr", fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)))
	return redisCache, nil
}

// ===== 业务层聚合 =====

// ProvideRepository 初始化仓储层 // di.ProvideRepository()
func ProvideRepository(db *gorm.DB) *Repository {
	return NewRepository(db)
}

// ProvideServices 初始化服务层聚合器 // di.ProvideServices()
func ProvideServices(repo *Repository) *Services {
	return NewServices(repo)
}

// ProvideHandlers 初始化处理器层聚合器 // di.ProvideHandlers()
func ProvideHandlers(services *Services) *Handlers {
	return NewHandlers(services)
}

// ProvideAppDependencies 初始化应用依赖聚合器 // di.ProvideAppDependencies()
func ProvideAppDependencies(
	cfg *config.Config,
	db *gorm.DB,
	cache cache.CacheInterface,
	repo *Repository,
	services *Services,
	handlers *Handlers,
) *AppDependencies {
	return &AppDependencies{
		Config:     cfg,
		DB:         db,
		Cache:      cache,
		Repository: repo,
		Services:   services,
		Handlers:   handlers,
	}
}

// ===== 路由层 =====

// ProvideRouter 初始化路由器 // di.ProvideRouter()
func ProvideRouter(handlers *Handlers) *router.Router {
	return router.NewRouter(handlers.Auth, handlers.User)
}

// ProvideGinEngine 初始化Gin引擎 // di.ProvideGinEngine()
func ProvideGinEngine(_ AppInit, r *router.Router) *gin.Engine {
	return r.Setup()
}

// ===== 资源清理 =====

// ProvideCleanup 提供资源清理函数 // di.ProvideCleanup()
func ProvideCleanup(deps *AppDependencies) func() {
	return func() {
		// 关闭缓存连接
		if deps.Cache != nil {
			if closer, ok := deps.Cache.(interface{ Close() error }); ok {
				if err := closer.Close(); err != nil {
					logger.Error("关闭缓存连接失败", logger.Err(err))
				}
			}
		}

		// 关闭数据库连接
		if deps.DB != nil {
			if err := database.Close(deps.DB); err != nil {
				logger.Error("关闭数据库连接失败", logger.Err(err))
			}
		}

		// 同步日志
		logger.Sync()
	}
}
