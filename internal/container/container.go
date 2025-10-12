package container

import (
	"fmt"
	"time"

	"gorm.io/gorm"

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
)

// Container 依赖注入容器
// 负责管理应用程序所有组件的生命周期和依赖关系
type Container struct {
	Config         *config.Config
	DB             *gorm.DB
	Cache          cache.CacheInterface
	UserRepo       repository.UserRepository
	AuthService    service.AuthService
	UserService    service.UserService
	AuthHandler    *handler.AuthHandler
	UserHandler    *handler.UserHandler
	Router         *router.Router
	RateLimiterFactory *middleware.RateLimiterFactory
	CircuitBreakerFactory *middleware.CircuitBreakerFactory
}

// NewContainer 创建依赖注入容器
// 按照依赖关系顺序初始化所有组件
func NewContainer(cfg *config.Config) (*Container, error) {
	container := &Container{
		Config: cfg,
	}

	// 按照依赖关系顺序初始化组件
	if err := container.initComponents(); err != nil {
		return nil, err
	}

	return container, nil
}

// initComponents 初始化所有组件
// 初始化顺序：基础组件 -> 数据层 -> 业务层 -> 表现层
func (c *Container) initComponents() error {
	// 1. 初始化基础组件
	if err := c.initBaseComponents(); err != nil {
		return err
	}

	// 2. 初始化数据层组件
	if err := c.initDataComponents(); err != nil {
		return err
	}

	// 3. 初始化业务层组件
	c.initServiceComponents()

	// 4. 初始化表现层组件
	c.initHandlerComponents()

	return nil
}

// initBaseComponents 初始化基础组件（日志、JWT、验证器等）
func (c *Container) initBaseComponents() error {
	// 初始化日志
	logConfig := logger.LogConfig{
		Level:      c.Config.Log.Level,
		Format:     c.Config.Log.Format,
		OutputPath: c.Config.Log.OutputPath,
		ReqLogPath: c.Config.Log.ReqLogPath,
		MaxSize:    c.Config.Log.MaxSize,
		MaxBackup:  c.Config.Log.MaxBackup,
		MaxAge:     c.Config.Log.MaxAge,
		Compress:   c.Config.Log.Compress,
	}
	if err := logger.Init(logConfig); err != nil {
		return err
	}

	// 初始化JWT
	jwtConfig := c.Config.JWT
	utils.InitJWT(jwtConfig)

	// 初始化验证器
	if err := validator.Init(); err != nil {
		return err
	}

	return nil
}

// initDataComponents 初始化数据层组件（数据库、缓存、仓储等）
func (c *Container) initDataComponents() error {
	// 初始化数据库
	mysqlConfig := c.Config.Database
	db, err := database.NewMySQL(mysqlConfig)
	if err != nil {
		return err
	}
	c.DB = db

	// 初始化Redis缓存
	redisConfig := cache.RedisConfig{
		Addr:         c.getRedisAddr(),
		Password:     c.Config.Redis.Password,
		DB:           c.Config.Redis.DB,
		PoolSize:     c.Config.Redis.PoolSize,
		MinIdleConns: c.Config.Redis.MinIdleConns,
		MaxRetries:   c.Config.Redis.MaxRetries,
	}

	// 尝试初始化Redis，但不强制要求成功
	redisCache, err := cache.NewRedisCache(redisConfig)
	if err != nil {
		logger.Warn("Redis初始化失败，将使用内存缓存", logger.Err(err))
		// 使用内存缓存作为后备
		c.Cache = cache.NewMemoryCache()
	} else {
		c.Cache = redisCache
		logger.Info("Redis缓存初始化成功")
	}

	// 初始化仓储层
	c.initRepositories()

	return nil
}

// initServiceComponents 初始化业务层组件
func (c *Container) initServiceComponents() {
	c.AuthService = service.NewAuthService(c.UserRepo)
	c.UserService = service.NewUserService(c.UserRepo)
	
	// 初始化限流器工厂
	globalRateLimiterConfig := middleware.DefaultRateLimiterConfig()
	globalRateLimiterConfig.MaxRequests = c.Config.RateLimiter.GlobalLimit
	globalRateLimiterConfig.Window = time.Duration(c.Config.RateLimiter.Window) * time.Second
	globalRateLimiterConfig.Algorithm = c.Config.RateLimiter.Algorithm
	
	userRateLimiterConfig := middleware.DefaultRateLimiterConfig()
	userRateLimiterConfig.MaxRequests = c.Config.RateLimiter.UserLimit
	userRateLimiterConfig.Window = time.Duration(c.Config.RateLimiter.Window) * time.Second
	userRateLimiterConfig.Algorithm = c.Config.RateLimiter.Algorithm
	
	c.RateLimiterFactory = middleware.NewRateLimiterFactory(globalRateLimiterConfig, userRateLimiterConfig, c.Cache)
	
	// 初始化熔断器工厂
	circuitBreakerConfig := middleware.DefaultCircuitBreakerConfig("global")
	circuitBreakerConfig.MaxRequests = uint32(c.Config.CircuitBreaker.MaxRequests)
	circuitBreakerConfig.MaxHalfOpenRequests = uint32(c.Config.CircuitBreaker.HalfOpenMaxRequests)
	circuitBreakerConfig.ErrorThreshold = c.Config.CircuitBreaker.ErrorThreshold
	circuitBreakerConfig.Timeout = time.Duration(c.Config.CircuitBreaker.Timeout) * time.Second
	circuitBreakerConfig.Enabled = c.Config.CircuitBreaker.Enabled
	
	c.CircuitBreakerFactory = middleware.NewCircuitBreakerFactory(circuitBreakerConfig)
}

// initHandlerComponents 初始化表现层组件
func (c *Container) initHandlerComponents() {
	c.AuthHandler = handler.NewAuthHandler(c.AuthService, c.UserService)
	c.UserHandler = handler.NewUserHandler(c.UserService)
	
	// 使用工厂模式创建中间件
	c.Router = router.NewRouterWithMiddleware(
		c.AuthHandler,
		c.UserHandler,
		c.RateLimiterFactory,
		c.CircuitBreakerFactory,
	)
	
	c.Router.Setup()
}

// initRepositories 初始化仓储层
func (c *Container) initRepositories() {
	c.UserRepo = repository.NewUserRepository(c.DB)
}

// getRedisAddr 获取Redis连接地址
func (c *Container) getRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Config.Redis.Host, c.Config.Redis.Port)
}

// Cleanup 清理资源
// 按照与初始化相反的顺序释放资源
func (c *Container) Cleanup() error {
	var lastErr error

	// 关闭缓存连接
	if c.Cache != nil {
		if closer, ok := c.Cache.(interface{ Close() error }); ok {
			if err := closer.Close(); err != nil {
				logger.Error("关闭缓存连接失败", logger.Err(err))
				lastErr = err
			}
		}
	}

	// 关闭数据库连接
	if c.DB != nil {
		if err := database.Close(c.DB); err != nil {
			logger.Error("关闭数据库连接失败", logger.Err(err))
			lastErr = err
		}
	}

	// 同步日志
	logger.Sync()

	return lastErr
}
