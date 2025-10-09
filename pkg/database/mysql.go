package database

import (
	"context"
	"fmt"
	"go_demo/pkg/logger"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// MySQLConfig MySQL配置
type MySQLConfig struct {
	Driver          string `mapstructure:"driver" yaml:"driver"`
	DSN             string `mapstructure:"dsn" yaml:"dsn"`
	MaxOpenConns    int    `mapstructure:"max_open_conns" yaml:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns" yaml:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime" yaml:"conn_max_lifetime"`
	ConnMaxIdleTime int    `mapstructure:"conn_max_idle_time" yaml:"conn_max_idle_time"`
	LogMode         bool   `mapstructure:"log_mode" yaml:"log_mode"`
	SlowThreshold   int    `mapstructure:"slow_threshold" yaml:"slow_threshold"`
}

// NewMySQL 创建MySQL数据库连接
func NewMySQL(cfg MySQLConfig) (*gorm.DB, error) {
	// 创建自定义的GORM日志适配器
	gormConfig := &gorm.Config{
		Logger: NewZapGormLogger(cfg),
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(cfg.DSN), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("连接MySQL数据库失败: %w", err)
	}

	// 获取底层的sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接池失败: %w", err)
	}

	// 设置连接池参数
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
	}
	if cfg.ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleTime) * time.Second)
	}

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

// ZapGormLogger 基于zap的GORM日志适配器
type ZapGormLogger struct {
	LogLevel                  gormLogger.LogLevel
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
}

// NewZapGormLogger 创建新的zap GORM日志适配器
func NewZapGormLogger(cfg MySQLConfig) *ZapGormLogger {
	var logLevel gormLogger.LogLevel
	if cfg.LogMode {
		logLevel = gormLogger.Info
	} else {
		logLevel = gormLogger.Silent
	}

	return &ZapGormLogger{
		LogLevel:                  logLevel,
		SlowThreshold:             time.Duration(cfg.SlowThreshold) * time.Millisecond,
		IgnoreRecordNotFoundError: true,
	}
}

// LogMode 设置日志级别
func (l *ZapGormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info 记录信息日志
func (l *ZapGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Info {
		logger.Info("GORM Info", logger.String("message", fmt.Sprintf(msg, data...)))
	}
}

// Warn 记录警告日志
func (l *ZapGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Warn {
		logger.Warn("GORM Warning", logger.String("message", fmt.Sprintf(msg, data...)))
	}
}

// Error 记录错误日志
func (l *ZapGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Error {
		logger.Error("GORM Error", logger.String("message", fmt.Sprintf(msg, data...)))
	}
}

// Trace 记录SQL执行日志
func (l *ZapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= gormLogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []zap.Field{
		logger.String("sql", sql),
		logger.String("duration", elapsed.String()),
		logger.Int64("rows", rows),
	}

	switch {
	case err != nil && l.LogLevel >= gormLogger.Error && (!l.IgnoreRecordNotFoundError || err != gormLogger.ErrRecordNotFound):
		logger.Error("SQL执行错误", append(fields, logger.Err(err))...)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormLogger.Warn:
		logger.Warn("慢查询检测", append(fields, logger.String("threshold", l.SlowThreshold.String()))...)
	case l.LogLevel == gormLogger.Info:
		logger.Debug("SQL执行", fields...)
	}
}
