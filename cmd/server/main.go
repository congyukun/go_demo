// 应用入口 main.go
// 说明：
// - 通过 di.InitializeServer("./configs/config.yaml") 使用 Wire 注入构建 *gin.Engine。
//   Wire 的注入流程见 internal/di/wire.go 与生成文件 internal/di/wire_gen.go。
// - JWT 初始化在 ProvideAppInit 中统一完成（utils.InitJWT(cfg.JWT)），确保 middleware.JWTAuthMiddleware 与 utils.*Token 使用同一配置。
// - 路由由 Router.Setup() 进行集中注册：健康检查、Swagger、API v1 分组（auth 与 users），并在 Router 层挂载全局/分组级中间件（日志、恢复、Trace、CORS、请求日志、限流、熔断）。
// - 优雅关闭：保持 HTTP Server 的 Shutdown 逻辑，资源清理由各组件自行负责（数据库/缓存在各自 Provider 的 Close 实现中处理）。
package main

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
	 // 使用 Wire 初始化 Gin 引擎
		engine, err := di.InitializeServer("./configs/config.yaml")
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize server with wire: %v", err))
	}

	// 读取端口与超时配置（从已初始化的日志中输出，但 Gin 引擎内部路由与中间件已通过 Wire 配置好）
	// 由于 Wire 内部完成了 config.Load，这里通过 expose 的 engine 配置 HTTP 服务器
	// 注意：Gin 模式在 ProvideLogger/ProvideJWT/ProvideValidator 后由 Router.Setup 统一设置

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    ":8080", // 端口由路由/配置决定；若需严格从配置读取，可在 di.InitializeServer 返回时一并返回 cfg.Server.Port
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
