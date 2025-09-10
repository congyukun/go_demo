package main

import (
	"context"
	"fmt"
	"go_demo/internal/config"
	"go_demo/internal/handler"
	"go_demo/internal/repository"
	"go_demo/internal/router"
	"go_demo/internal/service"
	"go_demo/pkg/database"
	"go_demo/pkg/logger"
	"go_demo/pkg/validator"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("加载配置失败: %v", err))
	}

	// 初始化日志
	logConfig := logger.LogConfig{
		Level:       cfg.Log.Level,
		Format:      cfg.Log.Format,
		OutputPath:  cfg.Log.OutputPath,
		ReqLogPath:  cfg.Log.ReqLogPath,
		MaxSize:     cfg.Log.MaxSize,
		MaxBackup:   cfg.Log.MaxBackup,
		MaxAge:      cfg.Log.MaxAge,
		Compress:    cfg.Log.Compress,
	}
	if err := logger.Init(logConfig); err != nil {
		panic(fmt.Sprintf("初始化日志失败: %v", err))
	}
	defer logger.Sync()

	logger.Info("服务启动中...",
		logger.String("app_name", cfg.App.Name),
		logger.String("version", cfg.App.Version),
		logger.String("env", cfg.App.Env),
	)

	// 初始化数据库
	mysqlConfig := database.MySQLConfig{
		Driver:          cfg.Database.MySQL.Driver,
		DSN:             cfg.Database.MySQL.DSN,
		MaxOpenConns:    cfg.Database.MySQL.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MySQL.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.MySQL.ConnMaxLifetime,
		ConnMaxIdleTime: cfg.Database.MySQL.ConnMaxIdleTime,
		LogMode:         cfg.Database.MySQL.LogMode,
		SlowThreshold:   cfg.Database.MySQL.SlowThreshold,
	}
	db, err := database.NewMySQL(mysqlConfig)
	if err != nil {
		logger.Fatal("数据库连接失败", logger.Err(err))
	}
	defer database.Close(db)

	// 初始化验证器
	if err := validator.Init(); err != nil {
		logger.Fatal("验证器初始化失败", logger.Err(err))
	}
	logger.Info("验证器初始化成功")

	// 初始化仓储层
	userRepo := repository.NewUserRepository(db)

	// 初始化服务层
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)

	// 初始化处理器
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	// 设置Gin模式
	if cfg.App.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由
	appRouter := router.NewRouter(authHandler, userHandler)
	ginEngine := appRouter.Setup()

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      ginEngine,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// 启动服务器
	go func() {
		logger.Info("HTTP服务器启动", logger.Int("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("服务器启动失败", logger.Err(err))
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("服务器正在关闭...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("服务器强制关闭", logger.Err(err))
	}

	logger.Info("服务器已关闭")
}
