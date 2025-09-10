package db

import (
	"context"
	"fmt"
	"go_demo/config"
	"go_demo/logger"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var mySQLConn *gorm.DB

func InitMySQLConn() error {
	cfg := config.GlobalConfig.Database.Mysql
	
	logger.Info("开始初始化MySQL连接",
		zap.String("dsn", maskDSN(cfg.Dsn)),
		zap.Int("max_open_conns", cfg.MaxOpenConns),
		zap.Int("max_idle_conns", cfg.MaxIdleConns))

	// 配置GORM日志
	var logLevel gormLogger.LogLevel
	if cfg.LogMode {
		logLevel = gormLogger.Info
	} else {
		logLevel = gormLogger.Silent
	}

	// 创建自定义的GORM日志适配器
	gormZapLogger := NewGormZapLogger(
		logger.Logger,
		gormLogger.Config{
			SlowThreshold:             cfg.SlowThreshold * time.Millisecond,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false, // zap已经处理颜色
		},
	)

	// 连接数据库
	db, err := gorm.Open(mysql.Open(cfg.Dsn), &gorm.Config{
		Logger: gormZapLogger,
	})

	if err != nil {
		logger.Error("数据库连接失败", zap.Error(err))
		return fmt.Errorf("数据库连接失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("获取sqlDB失败", zap.Error(err))
		return fmt.Errorf("获取sqlDB失败: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	if cfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
		logger.Debug("设置连接最大存活时间", zap.Duration("conn_max_lifetime", cfg.ConnMaxLifetime))
	}
	if cfg.ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
		logger.Debug("设置连接最大空闲时间", zap.Duration("conn_max_idle_time", cfg.ConnMaxIdleTime))
	}

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		logger.Error("数据库ping失败", zap.Error(err))
		return fmt.Errorf("数据库 ping 失败: %w", err)
	}

	mySQLConn = db
	logger.Info("MySQL连接初始化成功",
		zap.Int("max_open_conns", cfg.MaxOpenConns),
		zap.Int("max_idle_conns", cfg.MaxIdleConns),
		zap.Bool("log_mode", cfg.LogMode))
	return nil
}

// Close 关闭数据库连接
func Close() error {
	if mySQLConn != nil {
		logger.Info("正在关闭MySQL连接...")
		sqlDB, err := mySQLConn.DB()
		if err != nil {
			logger.Error("获取sqlDB失败", zap.Error(err))
			return err
		}
		if err := sqlDB.Close(); err != nil {
			logger.Error("关闭MySQL连接失败", zap.Error(err))
			return err
		}
		logger.Info("MySQL连接已关闭")
	}
	return nil
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return mySQLConn
}

// maskDSN 隐藏DSN中的敏感信息用于日志记录
func maskDSN(dsn string) string {
	// 简单的掩码处理，隐藏密码部分
	// 例如: root:password@tcp... -> root:***@tcp...
	if len(dsn) > 10 {
		return dsn[:10] + "***" + dsn[len(dsn)-20:]
	}
	return "***"
}
