//   - 通过 di.InitializeServer("./configs/config.yaml") 使用 Wire 注入构建 *gin.Engine。
//     Wire 的注入流程见 internal/di/wire.go 与生成文件 internal/di/wire_gen.go。
//   - JWT 初始化在 ProvideAppInit 中统一完成（utils.InitJWT(cfg.JWT)），确保 middleware.JWTAuthMiddleware 与 utils.*Token 使用同一配置。
//   - 路由由 Router.Setup() 进行集中注册：健康检查、Swagger、API v1 分组（auth 与 users），并在 Router 层挂载全局/分组级中间件（日志、恢复、Trace、CORS、请求日志、限流、熔断）。
//   - 优雅关闭：保持 HTTP Server 的 Shutdown 逻辑，资源清理由各组件自行负责（数据库/缓存在各自 Provider 的 Close 实现中处理）。
//   - 使用 Cobra 框架构建命令行界面
package server

import (
	"context"
	"fmt"
	"go_demo/internal/di"
	"go_demo/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

// 全局配置参数
var (
	configFile string
	port       string
)

// Execute 执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// rootCmd 代表没有调用子命令时的基础命令
var rootCmd = &cobra.Command{
	Use:   "go_demo",
	Short: "一个Go语言示例项目",
	Long:  `这是一个使用Go语言开发的示例项目，包含了Gin、GORM、Redis等常见组件的集成。`,
}

// serverCmd 代表server子命令
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动HTTP服务器",
	Long:  `启动HTTP服务器，提供RESTful API服务。`,
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func init() {
	// 注册子命令
	rootCmd.AddCommand(serverCmd)

	// 配置全局参数
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "./configs/config.yaml", "配置文件路径")

	// 配置server命令的参数
	serverCmd.Flags().StringVar(&port, "port", "", "服务器端口号（覆盖配置文件）")
}

// startServer 启动HTTP服务器
func startServer() {
	// 初始化服务器，带重试机制
	var app *di.ServerApp
	var err error

	maxRetries := 3
	retryDelay := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		app, err = di.InitializeServerApp(configFile)
		if err == nil {
			break
		}

		if i < maxRetries-1 {
			logger.Error("服务器初始化失败，正在重试...",
				logger.Err(err),
				logger.Int("retry", i+1),
				logger.Int("max_retries", maxRetries))
			time.Sleep(retryDelay)
			retryDelay *= 2 // 指数退避
		} else {
			logger.Fatal("服务器初始化失败，已达最大重试次数", logger.Err(err))
			os.Exit(1)
		}
	}

	// 注册资源清理函数
	defer func() {
		logger.Info("开始清理资源...")
		app.Cleanup()
		logger.Info("资源清理完成")
	}()

	// 创建HTTP服务器
	serverAddr := ":8080"
	if port != "" {
		serverAddr = fmt.Sprintf(":%s", port)
	}

	srv := &http.Server{
		Addr:           serverAddr,
		Handler:        app.Engine,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// 启动服务器
	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("HTTP服务器启动",
			logger.String("addr", srv.Addr),
			logger.String("config", configFile))
		serverErrors <- srv.ListenAndServe()
	}()

	// 等待中断信号或服务器错误
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal("服务器运行失败", logger.Err(err))
		}
	case sig := <-quit:
		logger.Info("收到关闭信号", logger.String("signal", sig.String()))
	}

	// 优雅关闭
	logger.Info("服务器正在关闭...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 关闭HTTP服务器
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("服务器关闭失败", logger.Err(err))
		// 强制关闭
		if err := srv.Close(); err != nil {
			logger.Error("服务器强制关闭失败", logger.Err(err))
		}
	}

	logger.Info("服务器已关闭")
}
