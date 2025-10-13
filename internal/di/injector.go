//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
)

// InitializeContainer 使用 Wire 注入依赖（仅在 wireinject 构建标签下编译）
func InitializeContainer() (*AppContainer, error) {
	wire.Build(
		ProviderSet,
		NewAppContainer,
	)
	return &AppContainer{}, nil
}