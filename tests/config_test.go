package tests

import (
	"go_demo/internal/config"
	"os"
	"testing"
)

func TestConfigLoad(t *testing.T) {
	// 创建临时配置文件
	configContent := `
server:
  port: 9090
  mode: test
  read_timeout: 30
  write_timeout: 30

database:
  driver: mysql
  dsn: "test:test@tcp(localhost:3306)/test_db?charset=utf8mb4&parseTime=True&loc=Local"
  max_open_conns: 50
  max_idle_conns: 5

jwt:
  secret_key: "test-secret-key"
  access_expire: 1800
  refresh_expire: 86400
  issuer: "test_app"

log:
  level: debug
  format: console
  output_path: "./test_logs/app.log"
  max_size: 50
  max_backup: 5
  max_age: 7
  compress: false

redis:
  host: localhost
  port: 6379
  db: 1
  pool_size: 5
`

	// 创建临时配置文件
	tmpFile, err := os.CreateTemp("", "config_test_*.yaml")
	if err != nil {
		t.Fatalf("创建临时配置文件失败: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("写入配置文件失败: %v", err)
	}
	tmpFile.Close()

	t.Run("加载配置文件", func(t *testing.T) {
		cfg, err := config.Load(tmpFile.Name())
		if err != nil {
			t.Fatalf("加载配置失败: %v", err)
		}

		// 验证服务器配置
		if cfg.Server.Port != 9090 {
			t.Errorf("服务器端口不匹配: 期望 9090, 实际 %d", cfg.Server.Port)
		}
		if cfg.Server.Mode != "test" {
			t.Errorf("服务器模式不匹配: 期望 test, 实际 %s", cfg.Server.Mode)
		}

		// 验证数据库配置
		if cfg.Database.Driver != "mysql" {
			t.Errorf("数据库驱动不匹配: 期望 mysql, 实际 %s", cfg.Database.Driver)
		}
		if cfg.Database.MaxOpenConns != 50 {
			t.Errorf("最大连接数不匹配: 期望 50, 实际 %d", cfg.Database.MaxOpenConns)
		}

		// 验证JWT配置
		if cfg.JWT.SecretKey != "test-secret-key" {
			t.Errorf("JWT密钥不匹配: 期望 test-secret-key, 实际 %s", cfg.JWT.SecretKey)
		}
		if cfg.JWT.AccessExpire != 1800 {
			t.Errorf("访问token过期时间不匹配: 期望 1800, 实际 %d", cfg.JWT.AccessExpire)
		}

		// 验证日志配置
		if cfg.Log.Level != "debug" {
			t.Errorf("日志级别不匹配: 期望 debug, 实际 %s", cfg.Log.Level)
		}
		if cfg.Log.Format != "console" {
			t.Errorf("日志格式不匹配: 期望 console, 实际 %s", cfg.Log.Format)
		}

		// 验证Redis配置
		if cfg.Redis.Host != "localhost" {
			t.Errorf("Redis主机不匹配: 期望 localhost, 实际 %s", cfg.Redis.Host)
		}
		if cfg.Redis.DB != 1 {
			t.Errorf("Redis数据库不匹配: 期望 1, 实际 %d", cfg.Redis.DB)
		}
	})

	t.Run("全局配置访问", func(t *testing.T) {
		// 加载配置到全局变量
		_, err := config.Load(tmpFile.Name())
		if err != nil {
			t.Fatalf("加载配置失败: %v", err)
		}

		// 测试全局配置访问
		globalCfg := config.GetConfig()
		if globalCfg == nil {
			t.Fatal("全局配置为空")
		}

		serverCfg := config.GetServerConfig()
		if serverCfg.Port != 9090 {
			t.Errorf("全局服务器配置端口不匹配: 期望 9090, 实际 %d", serverCfg.Port)
		}

		dbCfg := config.GetDatabaseConfig()
		if dbCfg.Driver != "mysql" {
			t.Errorf("全局数据库配置驱动不匹配: 期望 mysql, 实际 %s", dbCfg.Driver)
		}

		jwtCfg := config.GetJWTConfig()
		if jwtCfg.SecretKey != "test-secret-key" {
			t.Errorf("全局JWT配置密钥不匹配: 期望 test-secret-key, 实际 %s", jwtCfg.SecretKey)
		}

		logCfg := config.GetLogConfig()
		if logCfg.Level != "debug" {
			t.Errorf("全局日志配置级别不匹配: 期望 debug, 实际 %s", logCfg.Level)
		}

		redisCfg := config.GetRedisConfig()
		if redisCfg.Host != "localhost" {
			t.Errorf("全局Redis配置主机不匹配: 期望 localhost, 实际 %s", redisCfg.Host)
		}
	})

	t.Run("环境判断", func(t *testing.T) {
		// 加载测试配置
		_, err := config.Load(tmpFile.Name())
		if err != nil {
			t.Fatalf("加载配置失败: %v", err)
		}

		if !config.IsTest() {
			t.Error("应该识别为测试环境")
		}
		if config.IsProduction() {
			t.Error("不应该识别为生产环境")
		}
		if config.IsDevelopment() {
			t.Error("不应该识别为开发环境")
		}
	})
}

func TestConfigValidation(t *testing.T) {
	t.Run("无效端口", func(t *testing.T) {
		configContent := `
server:
  port: 70000  # 无效端口
  mode: debug

database:
  dsn: "test:test@tcp(localhost:3306)/test_db"

jwt:
  secret_key: "test-key"
  access_expire: 3600

log:
  output_path: "./logs/app.log"
`
		tmpFile, err := os.CreateTemp("", "invalid_config_*.yaml")
		if err != nil {
			t.Fatalf("创建临时配置文件失败: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		tmpFile.WriteString(configContent)
		tmpFile.Close()

		_, err = config.Load(tmpFile.Name())
		if err == nil {
			t.Error("应该返回端口验证错误")
		}
	})

	t.Run("空JWT密钥", func(t *testing.T) {
		configContent := `
server:
  port: 8080
  mode: debug

database:
  dsn: "test:test@tcp(localhost:3306)/test_db"

jwt:
  secret_key: ""  # 空密钥
  access_expire: 3600

log:
  output_path: "./logs/app.log"
`
		tmpFile, err := os.CreateTemp("", "invalid_jwt_*.yaml")
		if err != nil {
			t.Fatalf("创建临时配置文件失败: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		tmpFile.WriteString(configContent)
		tmpFile.Close()

		_, err = config.Load(tmpFile.Name())
		if err == nil {
			t.Error("应该返回JWT密钥验证错误")
		}
	})
}

func TestConfigDefaults(t *testing.T) {
	// 创建最小配置文件
	configContent := `
database:
  dsn: "test:test@tcp(localhost:3306)/test_db"

jwt:
  secret_key: "test-key"
`

	tmpFile, err := os.CreateTemp("", "minimal_config_*.yaml")
	if err != nil {
		t.Fatalf("创建临时配置文件失败: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString(configContent)
	tmpFile.Close()

	cfg, err := config.Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("加载最小配置失败: %v", err)
	}

	// 验证默认值
	if cfg.Server.Port != 8080 {
		t.Errorf("默认端口不正确: 期望 8080, 实际 %d", cfg.Server.Port)
	}
	if cfg.Server.Mode != "debug" {
		t.Errorf("默认模式不正确: 期望 debug, 实际 %s", cfg.Server.Mode)
	}
	if cfg.JWT.AccessExpire != 3600 {
		t.Errorf("默认访问token过期时间不正确: 期望 3600, 实际 %d", cfg.JWT.AccessExpire)
	}
	if cfg.Log.Level != "info" {
		t.Errorf("默认日志级别不正确: 期望 info, 实际 %s", cfg.Log.Level)
	}
}