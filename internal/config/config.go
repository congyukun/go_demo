package config

import (
	"fmt"

	"github.com/spf13/viper"

	"go_demo/internal/utils"
	"go_demo/pkg/database"
	"go_demo/pkg/logger"
)

// Config 应用配置结构
type Config struct {
	Server        ServerConfig         `mapstructure:"server" yaml:"server"`
	Database      database.MySQLConfig `mapstructure:"database" yaml:"database"`
	JWT           utils.JWTConfig      `mapstructure:"jwt" yaml:"jwt"`
	Log           logger.LogConfig     `mapstructure:"log" yaml:"log"`
	Redis         RedisConfig          `mapstructure:"redis" yaml:"redis"`
	RateLimiter   RateLimiterConfig    `mapstructure:"rate_limiter" yaml:"rate_limiter"`
	CircuitBreaker CircuitBreakerConfig `mapstructure:"circuit_breaker" yaml:"circuit_breaker"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int    `mapstructure:"port" yaml:"port"`
	Mode         string `mapstructure:"mode" yaml:"mode"`                   // debug, release, test
	ReadTimeout  int    `mapstructure:"read_timeout" yaml:"read_timeout"`   // 秒
	WriteTimeout int    `mapstructure:"write_timeout" yaml:"write_timeout"` // 秒
	MaxHeaderMB  int    `mapstructure:"max_header_mb" yaml:"max_header_mb"` // MB
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host         string `mapstructure:"host" yaml:"host"`
	Port         int    `mapstructure:"port" yaml:"port"`
	Password     string `mapstructure:"Password" yaml:"Password"`
	DB           int    `mapstructure:"db" yaml:"db"`
	PoolSize     int    `mapstructure:"pool_size" yaml:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns" yaml:"min_idle_conns"`
	MaxRetries   int    `mapstructure:"max_retries" yaml:"max_retries"`
}

// RateLimiterConfig 限流配置
type RateLimiterConfig struct {
	Enabled      bool `mapstructure:"enabled" yaml:"enabled"`
	GlobalLimit  int  `mapstructure:"global_limit" yaml:"global_limit"`    // 全局限流，单位：请求/秒
	UserLimit    int  `mapstructure:"user_limit" yaml:"user_limit"`        // 用户限流，单位：请求/分钟
	IPLimit      int  `mapstructure:"ip_limit" yaml:"ip_limit"`            // IP限流，单位：请求/分钟
	Window       int  `mapstructure:"window" yaml:"window"`                // 滑动窗口大小，单位：秒
	Algorithm    string `mapstructure:"algorithm" yaml:"algorithm"`        // 算法类型：fixed, sliding
}

// CircuitBreakerConfig 熔断器配置
type CircuitBreakerConfig struct {
	Enabled              bool    `mapstructure:"enabled" yaml:"enabled"`                           // 是否启用熔断器
	MaxRequests          int     `mapstructure:"max_requests" yaml:"max_requests"`                  // 触发熔断检查的最小请求数
	HalfOpenMaxRequests  int     `mapstructure:"half_open_max_requests" yaml:"half_open_max_requests"` // 半开状态最大请求数
	Timeout              int     `mapstructure:"timeout" yaml:"timeout"`                             // 熔断器打开后的超时时间，单位：秒
	ErrorThreshold       float64 `mapstructure:"error_threshold" yaml:"error_threshold"`             // 错误率阈值（0-1之间）
}

// 全局配置实例
var GlobalConfig *Config

// Load 加载配置文件
func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 设置默认值
	setDefaults()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}
	
	// 设置环境变量前缀，在读取配置文件之后，确保环境变量优先级更高
	viper.SetEnvPrefix("GO_DEMO")
	viper.AutomaticEnv()

	// 解析配置
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}
	// 验证配置
	if err := validateConfig(&config); err != nil {
		logger.Debug("配置内容", logger.Any("config", config))
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	// 设置全局配置
	GlobalConfig = &config

	return &config, nil
}

// setDefaults 设置默认配置值
func setDefaults() {
	// 服务器默认配置
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", 60)
	viper.SetDefault("server.write_timeout", 60)
	viper.SetDefault("server.max_header_mb", 1)

	// 数据库默认配置
	viper.SetDefault("database.driver", "mysql")
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.conn_max_lifetime", 3600)
	viper.SetDefault("database.conn_max_idle_time", 1800)
	viper.SetDefault("database.log_mode", true)
	viper.SetDefault("database.slow_threshold", 200)

	// JWT默认配置
	viper.SetDefault("jwt.secret_key", "wVvHx4e0ZvJ8a9d6r5t3yu0i9o7p6y4t3r2e1w2q3s4d5f6")
	viper.SetDefault("jwt.access_expire", 3600)    // 1小时
	viper.SetDefault("jwt.refresh_expire", 604800) // 7天
	viper.SetDefault("jwt.issuer", "go_demo")

	// 日志默认配置
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.output_path", "./logs/app.log")
	viper.SetDefault("log.req_log_path", "./logs/request.log")
	viper.SetDefault("log.max_size", 100)
	viper.SetDefault("log.max_backup", 10)
	viper.SetDefault("log.max_age", 30)
	viper.SetDefault("log.compress", true)

	// Redis默认配置
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)
	viper.SetDefault("redis.min_idle_conns", 5)
	viper.SetDefault("redis.max_retries", 3)

	// 限流默认配置
	viper.SetDefault("rate_limiter.enabled", true)
	viper.SetDefault("rate_limiter.global_limit", 1000)    // 全局限流1000请求/秒
	viper.SetDefault("rate_limiter.user_limit", 100)       // 用户限流100请求/分钟
	viper.SetDefault("rate_limiter.ip_limit", 200)         // IP限流200请求/分钟
	viper.SetDefault("rate_limiter.window", 60)            // 滑动窗口60秒
	viper.SetDefault("rate_limiter.algorithm", "sliding")  // 默认使用滑动窗口算法

	// 熔断器默认配置
	viper.SetDefault("circuit_breaker.enabled", true)
	viper.SetDefault("circuit_breaker.max_requests", 100)            // 触发熔断检查的最小请求数
	viper.SetDefault("circuit_breaker.half_open_max_requests", 10)   // 半开状态最大10个请求
	viper.SetDefault("circuit_breaker.timeout", 30)                  // 超时时间30秒
	viper.SetDefault("circuit_breaker.error_threshold", 0.5)        // 错误率阈值50%
}

// validateConfig 验证配置
func validateConfig(config *Config) error {
	// 验证服务器配置
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("无效的服务器端口: %d", config.Server.Port)
	}

	if config.Server.Mode != "debug" && config.Server.Mode != "release" && config.Server.Mode != "test" {
		return fmt.Errorf("无效的服务器模式: %s", config.Server.Mode)
	}

	// 验证数据库配置
	if config.Database.DSN == "" {
		return fmt.Errorf("数据库DSN不能为空")
	}

	// 验证JWT配置
	fmt.Printf("Debug - JWT config: %+v\n", config.JWT)
	if config.JWT.SecretKey == "" {
		return fmt.Errorf("JWT密钥不能为空")
	}

	if config.JWT.AccessExpire <= 0 {
		return fmt.Errorf("JWT访问token过期时间必须大于0")
	}

	// 验证日志配置
	if config.Log.OutputPath == "" {
		return fmt.Errorf("日志输出路径不能为空")
	}

	return nil
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	return GlobalConfig
}

// GetServerConfig 获取服务器配置
func GetServerConfig() ServerConfig {
	if GlobalConfig == nil {
		return ServerConfig{}
	}
	return GlobalConfig.Server
}

// GetDatabaseConfig 获取数据库配置
func GetDatabaseConfig() database.MySQLConfig {
	if GlobalConfig == nil {
		return database.MySQLConfig{}
	}
	return GlobalConfig.Database
}

// GetJWTConfig 获取JWT配置
func GetJWTConfig() utils.JWTConfig {
	if GlobalConfig == nil {
		return utils.JWTConfig{}
	}
	return GlobalConfig.JWT
}

// GetLogConfig 获取日志配置
func GetLogConfig() logger.LogConfig {
	if GlobalConfig == nil {
		return logger.LogConfig{}
	}
	return GlobalConfig.Log
}

// GetRedisConfig 获取Redis配置
func GetRedisConfig() RedisConfig {
	if GlobalConfig == nil {
		return RedisConfig{}
	}
	return GlobalConfig.Redis
}

// GetRateLimiterConfig 获取限流配置
func GetRateLimiterConfig() RateLimiterConfig {
	if GlobalConfig == nil {
		return RateLimiterConfig{}
	}
	return GlobalConfig.RateLimiter
}

// GetCircuitBreakerConfig 获取熔断器配置
func GetCircuitBreakerConfig() CircuitBreakerConfig {
	if GlobalConfig == nil {
		return CircuitBreakerConfig{}
	}
	return GlobalConfig.CircuitBreaker
}

// IsProduction 判断是否为生产环境
func IsProduction() bool {
	if GlobalConfig == nil {
		return false
	}
	return GlobalConfig.Server.Mode == "release"
}

// IsDevelopment 判断是否为开发环境
func IsDevelopment() bool {
	if GlobalConfig == nil {
		return true
	}
	return GlobalConfig.Server.Mode == "debug"
}

// IsTest 判断是否为测试环境
func IsTest() bool {
	if GlobalConfig == nil {
		return false
	}
	return GlobalConfig.Server.Mode == "test"
}

// LoadFromEnv 从环境变量加载配置（用于容器化部署）
func LoadFromEnv() (*Config, error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("GO_DEMO")

	// 设置默认值
	setDefaults()

	// 解析配置
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("从环境变量解析配置失败: %w", err)
	}

	// 验证配置
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("配置验证失败Env: %w", err)
	}

	// 设置全局配置
	GlobalConfig = &config

	return &config, nil
}

// ReloadConfig 重新加载配置（热更新）
func ReloadConfig(configPath string) error {
	newConfig, err := Load(configPath)
	if err != nil {
		return err
	}

	GlobalConfig = newConfig
	return nil
}
