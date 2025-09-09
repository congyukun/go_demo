package db

import (
	"context"
	"fmt"
	"go_demo/config"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var mySQLConn *gorm.DB

func InitMySQLConn() error {
	cfg := config.GlobalConfig.Mysql
	// 配置log
	var logLevel logger.LogLevel
	if cfg.LogMode {
		logLevel = logger.Info
	} else {
		logLevel = logger.Silent
	}
	// 自定义日志配置
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 输出目标
		logger.Config{
			SlowThreshold:             cfg.SlowThreshold * time.Millisecond, // 慢查询阈值
			LogLevel:                  logLevel,                             // 日志级别
			IgnoreRecordNotFoundError: true,                                 // 忽略记录不存在错误
			Colorful:                  true,                                 // 彩色输出
		},
	)

	// 连接mysql
	// 连接数据库
	db, err := gorm.Open(mysql.Open(cfg.Dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取sqlDB失败: %w", err)
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	if cfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}
	if cfg.ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	}
	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("数据库 ping 失败: %w", err)
	}

	mySQLConn = db
	fmt.Println("MySQL 连接初始化成功")
	return nil
}

// Close 关闭数据库连接
func Close() error {
	if mySQLConn != nil {
		sqlDB, err := mySQLConn.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
