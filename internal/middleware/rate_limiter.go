package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"go_demo/internal/utils"
	"go_demo/pkg/cache"
	"go_demo/pkg/logger"
)

// RateLimiterConfig 限流配置
type RateLimiterConfig struct {
	// 限流键生成函数，默认使用客户端IP
	KeyGenerator func(*gin.Context) string
	// 限流算法，支持"fixed_window"（固定窗口）和"sliding_window"（滑动窗口）
	Algorithm string
	// 时间窗口大小（秒）
	Window time.Duration
	// 窗口内允许的最大请求数
	MaxRequests int
	// 是否启用分布式限流
	Distributed bool
	// 缓存接口，用于分布式限流
	Cache cache.CacheInterface
	// 限流失败时的处理函数
	OnLimitReached func(*gin.Context)
}

// DefaultRateLimiterConfig 返回默认的限流配置
func DefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		KeyGenerator: func(c *gin.Context) string {
			return "rate_limit:" + c.ClientIP()
		},
		Algorithm:   "sliding_window",
		Window:      time.Minute,
		MaxRequests: 100,
		Distributed: true,
		OnLimitReached: func(c *gin.Context) {
			utils.ResponseError(c, http.StatusTooManyRequests, "请求过于频繁，请稍后再试")
		},
	}
}

// RateLimiter 限流中间件
func RateLimiter(config RateLimiterConfig) gin.HandlerFunc {
	// 设置默认值
	if config.KeyGenerator == nil {
		config.KeyGenerator = DefaultRateLimiterConfig().KeyGenerator
	}
	if config.Window == 0 {
		config.Window = DefaultRateLimiterConfig().Window
	}
	if config.MaxRequests == 0 {
		config.MaxRequests = DefaultRateLimiterConfig().MaxRequests
	}
	if config.OnLimitReached == nil {
		config.OnLimitReached = DefaultRateLimiterConfig().OnLimitReached
	}

	return func(c *gin.Context) {
		requestID := utils.GetRequestID(c)
		key := config.KeyGenerator(c)

		// 检查是否超过限流
		allowed, err := isAllowed(c.Request.Context(), key, config)
		if err != nil {
			logger.Error("限流检查失败",
				logger.String("request_id", requestID),
				logger.String("key", key),
				logger.Err(err),
			)
			// 限流检查失败时，允许请求通过
			c.Next()
			return
		}

		if !allowed {
			logger.Warn("请求被限流",
				logger.String("request_id", requestID),
				logger.String("key", key),
				logger.String("path", c.Request.URL.Path),
				logger.String("client_ip", c.ClientIP()),
			)
			config.OnLimitReached(c)
			c.Abort()
			return
		}

		c.Next()
	}
}

// isAllowed 检查是否允许请求
func isAllowed(ctx context.Context, key string, config RateLimiterConfig) (bool, error) {
	if config.Distributed && config.Cache != nil {
		return isAllowedDistributed(ctx, key, config)
	}
	return isAllowedLocal(key, config)
}

// IsAllowed 导出的isAllowed函数，供其他包使用
func IsAllowed(ctx context.Context, key string, config RateLimiterConfig) (bool, error) {
	return isAllowed(ctx, key, config)
}

// isAllowedDistributed 分布式限流检查
func isAllowedDistributed(ctx context.Context, key string, config RateLimiterConfig) (bool, error) {
	switch config.Algorithm {
	case "fixed_window":
		return isAllowedFixedWindowDistributed(ctx, key, config)
	case "sliding_window":
		return isAllowedSlidingWindowDistributed(ctx, key, config)
	default:
		return isAllowedSlidingWindowDistributed(ctx, key, config)
	}
}

// isAllowedFixedWindowDistributed 固定窗口分布式限流
func isAllowedFixedWindowDistributed(ctx context.Context, key string, config RateLimiterConfig) (bool, error) {
	// 使用Redis的INCR和EXPIRE实现固定窗口限流
	windowKey := fmt.Sprintf("%s:fixed:%d", key, time.Now().Unix()/int64(config.Window.Seconds()))
	
	// 尝试增加计数
	count, err := config.Cache.Increment(windowKey)
	if err != nil {
		return false, err
	}
	
	// 如果是第一次设置，设置过期时间
	if count == 1 {
		err = config.Cache.SetExpire(ctx, windowKey, config.Window)
		if err != nil {
			return false, err
		}
	}
	
	return count <= int64(config.MaxRequests), nil
}

// isAllowedSlidingWindowDistributed 滑动窗口分布式限流
func isAllowedSlidingWindowDistributed(ctx context.Context, key string, config RateLimiterConfig) (bool, error) {
	now := time.Now()
	windowStart := now.Add(-config.Window).Unix()
	
	// 使用Redis的ZSET实现滑动窗口限流
	// 1. 移除窗口外的记录
	_, err := config.Cache.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart, 10))
	// 如果有序集合不存在，忽略错误
	if err != nil && err.Error() != "有序集合不存在" {
		return false, err
	}
	
	// 2. 获取当前窗口内的请求数
	count, err := config.Cache.ZCard(key)
	if err != nil && err.Error() != "有序集合不存在" {
		return false, err
	}
	
	// 3. 如果超过限制，拒绝请求
	if count >= int64(config.MaxRequests) {
		return false, nil
	}
	
	// 4. 添加当前请求到窗口
	member := &cache.ZMember{
		Score:  float64(now.Unix()),
		Member: fmt.Sprintf("%d", now.UnixNano()),
	}
	err = config.Cache.ZAdd(key, member)
	if err != nil {
		return false, err
	}
	
	// 5. 设置过期时间
	err = config.Cache.SetExpire(ctx, key, config.Window)
	if err != nil {
		return false, err
	}
	
	return true, nil
}

// isAllowedLocal 本地限流检查
func isAllowedLocal(key string, config RateLimiterConfig) (bool, error) {
	// 这里可以使用内存中的限流器，例如golang.org/x/time/rate
	// 为了简单起见，这里只实现一个基本的本地限流
	// 在生产环境中，建议使用更成熟的限流库
	
	// 这里只是一个示例，实际应用中应该使用更精确的限流算法
	// 例如令牌桶或滑动窗口算法
	
	// 由于我们已经有Redis缓存，这里简化实现，直接使用分布式限流
	return false, fmt.Errorf("本地限流未实现，请使用分布式限流")
}

// APIRateLimiter API级别的限流中间件
func APIRateLimiter(path string, maxRequests int, window time.Duration) gin.HandlerFunc {
	config := DefaultRateLimiterConfig()
	config.KeyGenerator = func(c *gin.Context) string {
		return fmt.Sprintf("api_rate_limit:%s:%s", path, c.ClientIP())
	}
	config.MaxRequests = maxRequests
	config.Window = window
	
	return RateLimiter(config)
}

// UserRateLimiter 用户级别的限流中间件
func UserRateLimiter(maxRequests int, window time.Duration) gin.HandlerFunc {
	config := DefaultRateLimiterConfig()
	config.KeyGenerator = func(c *gin.Context) string {
		// 尝试从上下文获取用户ID
		if userID, exists := c.Get("user_id"); exists {
			return fmt.Sprintf("user_rate_limit:%v", userID)
		}
		// 如果没有用户ID，使用IP
		return fmt.Sprintf("user_rate_limit:ip:%s", c.ClientIP())
	}
	config.MaxRequests = maxRequests
	config.Window = window
	
	return RateLimiter(config)
}

// RateLimiterFactory 限流器工厂
type RateLimiterFactory struct {
	globalConfig RateLimiterConfig
	userConfig   RateLimiterConfig
	Cache        cache.CacheInterface // 导出cache字段
}

// NewRateLimiterFactory 创建限流器工厂
func NewRateLimiterFactory(globalConfig, userConfig RateLimiterConfig, cache cache.CacheInterface) *RateLimiterFactory {
	// 设置默认值
	if globalConfig.Cache == nil {
		globalConfig.Cache = cache
	}
	if userConfig.Cache == nil {
		userConfig.Cache = cache
	}
	
	return &RateLimiterFactory{
		globalConfig: globalConfig,
		userConfig:   userConfig,
		Cache:        cache,
	}
}

// GlobalMiddleware 返回全局限流中间件
func (f *RateLimiterFactory) GlobalMiddleware() gin.HandlerFunc {
	return RateLimiter(f.globalConfig)
}

// UserMiddleware 返回用户级限流中间件
func (f *RateLimiterFactory) UserMiddleware() gin.HandlerFunc {
	return RateLimiter(f.userConfig)
}