package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config 应用配置结构
type Config struct {
	App      AppConfig      `yaml:"app"`
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Log      LogConfig      `yaml:"log"`
}

// AppConfig 应用基础配置
type AppConfig struct {
	Name    string `yaml:"name"`
	Env     string `yaml:"env"`
	Version string `yaml:"version"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int `yaml:"port"`
	Timeout      int `yaml:"timeout"`
	ReadTimeout  int `yaml:"read_timeout"`
	WriteTimeout int `yaml:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	MySQL    MySQLConfig    `yaml:"mysql"`
	Postgres PostgresConfig `yaml:"postgres"`
	SQLite   SQLiteConfig   `yaml:"sqlite"`
	MongoDB  MongoDBConfig  `yaml:"mongodb"`
}

// MySQLConfig MySQL配置
type MySQLConfig struct {
	Driver          string `yaml:"driver"`
	DSN             string `yaml:"dsn"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime int    `yaml:"conn_max_idle_time"`
	LogMode         bool   `yaml:"log_mode"`
	SlowThreshold   int    `yaml:"slow_threshold"`
}

// PostgresConfig PostgreSQL配置
type PostgresConfig struct {
	Driver          string `yaml:"driver"`
	DSN             string `yaml:"dsn"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime int    `yaml:"conn_max_idle_time"`
}

// SQLiteConfig SQLite配置
type SQLiteConfig struct {
	Driver          string `yaml:"driver"`
	DSN             string `yaml:"dsn"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}

// MongoDBConfig MongoDB配置
type MongoDBConfig struct {
	URI         string `yaml:"uri"`
	Database    string `yaml:"database"`
	MaxPoolSize int    `yaml:"max_pool_size"`
	MinPoolSize int    `yaml:"min_pool_size"`
	Timeout     int    `yaml:"timeout"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr         string `yaml:"addr"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"db"`
	PoolSize     int    `yaml:"pool_size"`
	MinIdleConns int    `yaml:"min_idle_conns"`
	IdleTimeout  int    `yaml:"idle_timeout"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level       string `yaml:"level"`
	Format      string `yaml:"format"`
	OutputPath  string `yaml:"output_path"`
	ReqLogPath  string `yaml:"req_log_path"`
	MaxSize     int    `yaml:"max_size"`
	MaxBackup   int    `yaml:"max_backup"`
	MaxAge      int    `yaml:"max_age"`
	Compress    bool   `yaml:"compress"`
}

var GlobalConfig *Config

// Load 加载配置文件
func Load() (*Config, error) {
	configPath := getConfigPath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 设置默认值
	setDefaults(&config)

	GlobalConfig = &config
	return &config, nil
}

// getConfigPath 获取配置文件路径
func getConfigPath() string {
	// 优先使用环境变量指定的配置文件
	if configPath := os.Getenv("CONFIG_PATH"); configPath != "" {
		return configPath
	}

	// 默认配置文件路径
	return filepath.Join("configs", "config.yaml")
}

// setDefaults 设置默认值
func setDefaults(config *Config) {
	if config.App.Name == "" {
		config.App.Name = "go-demo"
	}
	if config.App.Env == "" {
		config.App.Env = "dev"
	}
	if config.App.Version == "" {
		config.App.Version = "1.0.0"
	}

	if config.Server.Port == 0 {
		config.Server.Port = 8080
	}
	if config.Server.ReadTimeout == 0 {
		config.Server.ReadTimeout = 15
	}
	if config.Server.WriteTimeout == 0 {
		config.Server.WriteTimeout = 15
	}

	if config.Log.Level == "" {
		config.Log.Level = "info"
	}
	if config.Log.Format == "" {
		config.Log.Format = "console"
	}
	if config.Log.OutputPath == "" {
		config.Log.OutputPath = "logs/app.log"
	}
}
