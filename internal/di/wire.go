// Package di 的 Wire 声明文件（仅在生成器阶段编译）。
// 说明：
// - wire.go 使用 //go:build wireinject，使其只在 wire 代码生成时参与编译；实际运行使用 wire_gen.go。
// - baseSet/dataSet/serviceSet/middlewareSet/handlerSet/routerSet 分层组织 Provider，确保依赖清晰。
// - InitializeServer(configPath) 声明注入器，Wire 会在 wire_gen.go 中生成具体实现，返回 *gin.Engine。
// 使用：
//   1) 安装 wire：go install github.com/google/wire/cmd/wire@latest
//   2) 生成代码：cd internal/di && wire
//   3) 启动服务：go build ./cmd/server && ./cmd/server
//go:build wireinject
// +build wireinject

package di

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// Sets
var baseSet = wire.NewSet(
	ProvideConfig,
	ProvideAppInit,
)

var dataSet = wire.NewSet(
	ProvideDB,
	ProvideCache,
	ProvideRepositories,
)

var serviceSet = wire.NewSet(
	ProvideAuthService,
	ProvideUserService,
)

var middlewareSet = wire.NewSet(
	ProvideRateLimiterFactory,
	ProvideCircuitBreakerFactory,
)

var handlerSet = wire.NewSet(
	ProvideAuthHandler,
	ProvideUserHandler,
)

var routerSet = wire.NewSet(
	ProvideRouter,
	ProvideGinEngine,
)

// InitializeServer 使用 Wire 构建 Gin Engine 与清理函数
func InitializeServer(configPath string) (*gin.Engine, error) { // di.InitializeServer()
	wire.Build(
		baseSet,
		dataSet,
		serviceSet,
		middlewareSet,
		handlerSet,
		routerSet,
	)

	return nil, fmt.Errorf("wire build failed") // 实际由 wire 生成替换
}
