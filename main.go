package main

import (
	"go_demo/config"
	"go_demo/db"
	"go_demo/logger"
	"go_demo/routes"
	"strconv"

	"go.uber.org/zap"
)

func main() {

	// 初始化配置
	if err := config.InitConfig(""); err != nil {
		panic("config init error:" + err.Error())
	}
	cfg := config.GlobalConfig
	serverConfig := cfg.Server
	if serverConfig.Port == 0 {
		panic("port is invalid")
	}

	// 初始化zap日志
	if err := logger.InitZapLogger(); err != nil {
		panic("logger init error:" + err.Error())
	}
	defer logger.Close()

	logger.Info("应用启动中...",
		zap.String("app_name", cfg.App.Name),
		zap.String("version", cfg.App.Version),
		zap.String("env", cfg.App.Env),
		zap.Int("port", serverConfig.Port))

	// 初始化数据库
	if err := db.InitMySQLConn(); err != nil {
		logger.Fatal("mysql连接失败", zap.Error(err))
	}
	defer db.Close()

	// 初始化路由
	r := routes.SetupRouter()

	// 启动服务
	addr := ":" + strconv.Itoa(serverConfig.Port)
	logger.Infof("服务启动在端口: %s", addr)
	if err := r.Run(addr); err != nil {
		logger.Fatal("服务启动失败", zap.Error(err))
	}
}
