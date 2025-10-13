package main

import (
	"context"
	"errors"
	"fmt"
	"go_demo/internal/di"
	"go_demo/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// @title Go Demo API
// @version 1.0.0
// @description Go Demo 项目的API文档
// @description 包含用户认证、用户管理等功能的RESTful API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@godemo.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// 初始化容器
	container, err := di.InitializeContainer()
	if err != nil {
		panic(fmt.Sprintf("初始化容器失败: %v", err))
	}

	// 初始化日志
	defer logger.Sync()

	logger.Info("服务启动中...",
		logger.String("app_name", "go_demo"),
		logger.String("version", "1.0.0"),
		logger.String("mode", container.Config.Server.Mode),
	)

	// 设置Gin模式
	if container.Config.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由
	ginEngine := container.Router.Setup()

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", container.Config.Server.Port),
		Handler:      ginEngine,
		ReadTimeout:  time.Duration(container.Config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(container.Config.Server.WriteTimeout) * time.Second,
	}

	// 启动服务器
	go func() {
		logger.Info("HTTP服务器启动", logger.Int("port", container.Config.Server.Port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
