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
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"go_demo/internal/di"
	"go_demo/pkg/logger"
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
	engine, err := di.InitializeServer(configFile)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize server with wire: %v", err))
	}

	// 读取端口与超时配置（从已初始化的日志中输出，但 Gin 引擎内部路由与中间件已通过 Wire 配置好）
	// 由于 Wire 内部完成了 config.Load，这里通过 expose 的 engine 配置 HTTP 服务器
	// 注意：Gin 模式在 ProvideLogger/ProvideJWT/ProvideValidator 后由 Router.Setup 统一设置

	// 创建HTTP服务器
	serverAddr := ":8080"
	if port != "" {
		serverAddr = fmt.Sprintf(":%s", port)
	}

	srv := &http.Server{
		Addr:    serverAddr, // 端口由路由/配置决定；若需严格从配置读取，可在 di.InitializeServer 返回时一并返回 cfg.Server.Port
		Handler: engine,
		// 如需严格控制超时，可在 di 中增加 ProvideServerConfig 返回具体值，这里保持与原实现一致的超时策略
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// 启动服务器
	go func() {
		logger.Info("HTTP服务器启动", logger.String("addr", srv.Addr))
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

	// 关闭HTTP服务器
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("服务器强制关闭", logger.Err(err))
	}

	// 资源清理由进程退出时的defer或各自组件负责；此处无全局cleanup

	logger.Info("服务器已关闭")
}
