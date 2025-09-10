package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig  `mapstructure:"redis"`
	Log      LogConfig    `mapstructure:"log"`
	App      AppConfig    `mapstructure:"app"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Mysql MySQLConfig `mapstructure:"mysql"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Env string `mapstructure:"env"`
	Version string `mapstructure:"version"`
}
// ServerConfig 服务器配置
type ServerConfig struct {
	Port     int    `mapstructure:"port"`
	Mode     string `mapstructure:"mode"`
	Timeout  int    `mapstructure:"timeout"`
	LogLevel string `mapstructure:"log_level"`
}

type DbConfig struct {
	Driver       string `mapstructure:"driver"`
	Dsn          string `mapstructure:"dsn"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxLifeTime  int    `mapstructure:"max_life_time"`
}

// MySQLConfig MySQL配置详情
type MySQLConfig struct {
	Driver         string        `mapstructure:"driver"`
	Dsn            string        `mapstructure:"dsn"`
	MaxOpenConns   int           `mapstructure:"max_open_conns"`
	MaxIdleConns   int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
	LogMode        bool          `mapstructure:"log_mode"`
	SlowThreshold  time.Duration `mapstructure:"slow_threshold"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	Db       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	OutputPath string `mapstructure:"output_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackup  int    `mapstructure:"max_backup"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

var GlobalConfig Config

func InitConfig(configPath string) error {
	// 如果 configPath 为空，则使用默认配置文件路径
	if configPath == "" {
		// 获取当前文件所在目录
		workDir, _ := os.Getwd()
		// 构造配置文件路径
		configPath = filepath.Join(workDir, "config")
	}
	// 初始化 viper
	viper.SetConfigName("config")   // 配置文件名称（不含扩展名）
	viper.SetConfigType("yaml")     // 配置文件类型
	viper.AddConfigPath(configPath) // 配置文件路径

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}
	// 解析配置文件到全局变量
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}
	// 监测配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件变化, 重新加载配置...")
		if err := viper.Unmarshal(&GlobalConfig); err != nil {
			fmt.Printf("解析配置文件失败: %v\n", err)
		} else {
			fmt.Println("配置文件重新加载成功!")
		}
	})
	// 支持环境变量覆盖配置（可选）
	viper.AutomaticEnv()
	// viper.SetEnvPrefix("app")                              // 环境变量前缀
	// viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // 替换分隔符

	return nil

}
