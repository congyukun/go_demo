// Package di 的 Wire 声明文件（仅在生成器阶段编译）
// 简化的依赖注入配置，使用聚合器模式减少复杂度
//go:build wireinject
// +build wireinject

package di

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// 基础设施集合
var infrastructureSet = wire.NewSet(
	ProvideConfig,
	ProvideAppInit,
	ProvideDB,
	ProvideCache,
)

// 业务逻辑集合
var businessSet = wire.NewSet(
	ProvideRepository,
	ProvideServices,
	ProvideHandlers,
	ProvideAppDependencies,
)

// 应用层集合
var applicationSet = wire.NewSet(
	ProvideRouter,
	ProvideGinEngine,
)

// InitializeServer 使用 Wire 构建 Gin Engine
func InitializeServer(configPath string) (*gin.Engine, error) { // di.InitializeServer()
	wire.Build(
		infrastructureSet,
		businessSet,
		applicationSet,
	)

	return nil, fmt.Errorf("wire build failed") // 实际由 wire 生成替换
}

// InitializeApp 初始化完整应用依赖（可选，用于测试或其他场景）
func InitializeApp(configPath string) (*AppDependencies, error) { // di.InitializeApp()
	wire.Build(
		infrastructureSet,
		businessSet,
	)

	return nil, fmt.Errorf("wire build failed") // 实际由 wire 生成替换
}
