// Package server 提供HTTP服务器的启动和管理功能
//
// 主要功能：
//   - 通过 di.InitializeServerApp() 使用 Wire 注入构建服务器应用
//   - 支持命令行参数和配置文件的灵活配置
//   - 实现优雅关闭，确保资源正确释放
//   - 使用 Cobra 框架构建命令行界面
package server

import (
	"context"
	"fmt"
	"go_demo/internal/config"
	"go_demo/internal/di"
	"go_demo/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

// 命令行参数
var (
	configFile string // 配置文件路径
	port       string // 服务器端口号（覆盖配置文件）
)

// Execute 执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// rootCmd 根命令
var rootCmd = &cobra.Command{
	Use:   "go_demo",
	Short: "Go Demo 项目",
	Long:  `一个使用 Go 语言开发的示例项目，集成了 Gin、GORM、Redis 等常见组件。`,
}

// serverCmd 服务器子命令
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动 HTTP 服务器",
	Long:  `启动 HTTP 服务器，提供 RESTful API 服务。`,
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// 全局参数
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "./configs/config.yaml", "配置文件路径")

	// server 命令参数
	serverCmd.Flags().StringVar(&port, "port", "", "服务器端口号（覆盖配置文件）")
}

// startServer 启动HTTP服务器
func startServer() {
	// 初始化服务器应用（带重试机制）
	app := initServerApp()
	defer cleanupResources(app)

	// 创建并启动HTTP服务器
	srv := createHTTPServer(app)
	startHTTPServer(srv)

	// 等待关闭信号
	waitForShutdown(srv)
}

// initServerApp 初始化服务器应用，带重试机制
func initServerApp() *di.ServerApp {
	const (
		maxRetries = 3
		baseDelay  = 2 * time.Second
	)

	var app *di.ServerApp
	var err error
	retryDelay := baseDelay

	for i := 0; i < maxRetries; i++ {
		app, err = di.InitializeServerApp(configFile)
		if err == nil {
			return app
		}

		if i < maxRetries-1 {
			logger.Error("服务器初始化失败，正在重试...",
				logger.Err(err),
				logger.Int("retry", i+1),
				logger.Int("max_retries", maxRetries))
			time.Sleep(retryDelay)
			retryDelay *= 2 // 指数退避
		}
	}

	logger.Fatal("服务器初始化失败，已达最大重试次数", logger.Err(err))
	os.Exit(1)
	return nil
}

// cleanupResources 清理资源
func cleanupResources(app *di.ServerApp) {
	logger.Info("开始清理资源...")
	app.Cleanup()
	logger.Info("资源清理完成")
}

// createHTTPServer 创建HTTP服务器
func createHTTPServer(app *di.ServerApp) *http.Server {
	cfg := config.GetServerConfig()

	// 确定服务器地址：命令行参数优先于配置文件
	addr := fmt.Sprintf(":%d", cfg.Port)
	if port != "" {
		addr = ":" + port
	}

	// 获取超时配置，使用默认值兜底
	readTimeout := getTimeout(cfg.ReadTimeout, 30)
	writeTimeout := getTimeout(cfg.WriteTimeout, 30)
	maxHeaderBytes := getMaxHeaderBytes(cfg.MaxHeaderMB, 1)

	return &http.Server{
		Addr:           addr,
		Handler:        app.Engine,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}
}

// getTimeout 获取超时时间，如果配置值无效则使用默认值
func getTimeout(seconds, defaultSeconds int) time.Duration {
	if seconds > 0 {
		return time.Duration(seconds) * time.Second
	}
	return time.Duration(defaultSeconds) * time.Second
}

// getMaxHeaderBytes 获取最大请求头大小（字节），如果配置值无效则使用默认值
func getMaxHeaderBytes(mb, defaultMB int) int {
	if mb > 0 {
		return mb << 20
	}
	return defaultMB << 20
}

// startHTTPServer 启动HTTP服务器（非阻塞）
func startHTTPServer(srv *http.Server) {
	go func() {
		logger.Info("HTTP服务器启动",
			logger.String("addr", srv.Addr),
			logger.String("config", configFile))

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("服务器运行失败", logger.Err(err))
		}
	}()
}

// waitForShutdown 等待关闭信号并优雅关闭服务器
func waitForShutdown(srv *http.Server) {
	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	logger.Info("收到关闭信号", logger.String("signal", sig.String()))

	// 优雅关闭
	gracefulShutdown(srv)
}

// gracefulShutdown 优雅关闭HTTP服务器
func gracefulShutdown(srv *http.Server) {
	const shutdownTimeout = 30 * time.Second

	logger.Info("服务器正在关闭...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("服务器优雅关闭失败，尝试强制关闭", logger.Err(err))
		if err := srv.Close(); err != nil {
			logger.Error("服务器强制关闭失败", logger.Err(err))
		}
	}

	logger.Info("服务器已关闭")
}
