package db

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// GormZapLogger GORM的zap日志适配器
type GormZapLogger struct {
	ZapLogger                 *zap.Logger
	LogLevel                  gormLogger.LogLevel
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
}

// NewGormZapLogger 创建新的GORM zap日志适配器
func NewGormZapLogger(zapLogger *zap.Logger, config gormLogger.Config) gormLogger.Interface {
	return &GormZapLogger{
		ZapLogger:                 zapLogger,
		LogLevel:                  config.LogLevel,
		SlowThreshold:             config.SlowThreshold,
		IgnoreRecordNotFoundError: config.IgnoreRecordNotFoundError,
	}
}

// LogMode 设置日志模式
func (l *GormZapLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info 记录info级别日志
func (l *GormZapLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Info {
		l.ZapLogger.Sugar().Infof(msg, data...)
	}
}

// Warn 记录warn级别日志
func (l *GormZapLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Warn {
		l.ZapLogger.Sugar().Warnf(msg, data...)
	}
}

// Error 记录error级别日志
func (l *GormZapLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Error {
		l.ZapLogger.Sugar().Errorf(msg, data...)
	}
}

// Trace 记录SQL执行日志
func (l *GormZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= gormLogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Duration("elapsed", elapsed),
		zap.Int64("rows", rows),
	}

	switch {
	case err != nil && l.LogLevel >= gormLogger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		l.ZapLogger.Error("SQL执行错误", append(fields, zap.Error(err))...)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormLogger.Warn:
		l.ZapLogger.Warn("慢查询检测", append(fields, zap.Duration("threshold", l.SlowThreshold))...)
	case l.LogLevel == gormLogger.Info:
		l.ZapLogger.Info("SQL执行", fields...)
	}
}

// ParamsFilter 参数过滤器（可选实现）
func (l *GormZapLogger) ParamsFilter(ctx context.Context, sql string, params ...interface{}) (string, []interface{}) {
	if l.LogLevel == gormLogger.Info {
		return sql, params
	}
	return sql, nil
}