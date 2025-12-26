package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MySQLConfig MySQL配置
type MySQLConfig struct {
	Driver          string `mapstructure:"driver" yaml:"driver"`
	DSN             string `mapstructure:"dsn" yaml:"dsn"`
	MaxOpenConns    int    `mapstructure:"max_open_conns" yaml:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns" yaml:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime" yaml:"conn_max_lifetime"`   // 秒
	ConnMaxIdleTime int    `mapstructure:"conn_max_idle_time" yaml:"conn_max_idle_time"` // 秒
	LogMode         bool   `mapstructure:"log_mode" yaml:"log_mode"`
	SlowThreshold   int    `mapstructure:"slow_threshold" yaml:"slow_threshold"` // 毫秒
}

// NewMySQL 创建MySQL连接
func NewMySQL(config MySQLConfig) (*gorm.DB, error) {
	// 设置日志级别
	logLevel := logger.Silent
	if config.LogMode {
		logLevel = logger.Info
	}

	// 配置GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(config.DSN), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 获取底层的sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库实例失败: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(config.ConnMaxIdleTime) * time.Second)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("数据库连接测试失败: %w", err)
	}

	return db, nil
}

// Close 关闭数据库连接
func Close(db *gorm.DB) error {
	if db == nil {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// HealthCheck 数据库健康检查
func HealthCheck(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("数据库连接为空")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}
	return sqlDB.Ping()
}

// GetStats 获取数据库连接池统计信息
func GetStats(db *gorm.DB) (map[string]interface{}, error) {
	if db == nil {
		return nil, fmt.Errorf("数据库连接为空")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration.String(),
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}, nil
}
