package tests

import (
	"go_demo/pkg/database"
	"go_demo/pkg/logger"
	"testing"
)

func TestMySQLZapLogger(t *testing.T) {
	// 初始化日志
	logConfig := logger.LogConfig{
		Level:      "debug",
		Format:     "console",
		OutputPath: "logs/test.log",
		MaxSize:    10,
		MaxBackup:  3,
		MaxAge:     7,
		Compress:   false,
	}
	
	err := logger.Init(logConfig)
	if err != nil {
		t.Fatalf("初始化日志失败: %v", err)
	}
	defer logger.Sync()

	// 测试MySQL配置
	mysqlConfig := database.MySQLConfig{
		Driver:            "mysql",
		DSN:               "test:test@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local",
		MaxOpenConns:      10,
		MaxIdleConns:      5,
		ConnMaxLifetime:   3600,
		ConnMaxIdleTime:   1800,
		LogMode:           true,
		SlowThreshold:     100, // 100ms慢查询阈值
	}

	// 创建zap GORM日志适配器
	zapLogger := database.NewZapGormLogger(mysqlConfig)
	
	if zapLogger == nil {
		t.Error("创建zap GORM日志适配器失败")
	}

	// 测试日志级别设置
	if zapLogger.LogLevel == 0 {
		t.Error("日志级别设置错误")
	}

	// 测试慢查询阈值
	if zapLogger.SlowThreshold.Milliseconds() != 100 {
		t.Errorf("慢查询阈值设置错误，期望100ms，实际%dms", zapLogger.SlowThreshold.Milliseconds())
	}

	t.Log("MySQL zap日志集成测试通过")
}

func TestMySQLConnection(t *testing.T) {
	// 这个测试需要实际的MySQL连接，在CI/CD环境中可能会跳过
	t.Skip("跳过MySQL连接测试，需要实际的数据库环境")

	// 如果有测试数据库，可以取消注释以下代码
	/*
	logConfig := logger.LogConfig{
		Level:      "info",
		Format:     "console",
		OutputPath: "logs/test.log",
		MaxSize:    10,
		MaxBackup:  3,
		MaxAge:     7,
		Compress:   false,
	}
	
	err := logger.Init(logConfig)
	if err != nil {
		t.Fatalf("初始化日志失败: %v", err)
	}
	defer logger.Sync()

	mysqlConfig := database.MySQLConfig{
		Driver:            "mysql",
		DSN:               "root:123456@tcp(127.0.0.1:3306)/go_test?charset=utf8mb4&parseTime=True&loc=Local",
		MaxOpenConns:      10,
		MaxIdleConns:      5,
		ConnMaxLifetime:   3600,
		ConnMaxIdleTime:   1800,
		LogMode:           true,
		SlowThreshold:     100,
	}

	db, err := database.NewMySQL(mysqlConfig)
	if err != nil {
		t.Fatalf("连接MySQL失败: %v", err)
	}
	defer database.Close(db)

	// 执行一个简单的查询来测试日志
	var result int
	err = db.Raw("SELECT 1").Scan(&result).Error
	if err != nil {
		t.Fatalf("执行测试查询失败: %v", err)
	}

	if result != 1 {
		t.Errorf("查询结果错误，期望1，实际%d", result)
	}

	t.Log("MySQL连接和zap日志测试通过")
	*/
}