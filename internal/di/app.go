// Package di 提供简洁的依赖注入配置
package di

import (
	"go_demo/internal/config"
	"go_demo/internal/handler"
	"go_demo/internal/repository"
	"go_demo/internal/service"
	"go_demo/pkg/cache"

	"gorm.io/gorm"
)

// AppDependencies 应用依赖聚合器 // di.AppDependencies
type AppDependencies struct {
	Config     *config.Config            // di.AppDependencies.Config
	DB         *gorm.DB                  // di.AppDependencies.DB
	Cache      cache.CacheInterface      // di.AppDependencies.Cache
	Repository repository.UserRepository // di.AppDependencies.Repository
	Services   *Services                 // di.AppDependencies.Services
	Handlers   *Handlers                 // di.AppDependencies.Handlers
}

// Services 服务层聚合器 // di.Services
type Services struct {
	Auth service.AuthService // di.Services.Auth
	User service.UserService // di.Services.User
}

// Handlers 处理器层聚合器 // di.Handlers
type Handlers struct {
	Auth *handler.AuthHandler // di.Handlers.Auth
	User *handler.UserHandler // di.Handlers.User
}

// NewServices 创建服务聚合器 // di.NewServices()
func NewServices(repo repository.UserRepository) *Services {
	return &Services{
		Auth: service.NewAuthService(repo),
		User: service.NewUserService(repo),
	}
}

// NewHandlers 创建处理器聚合器 // di.NewHandlers()
func NewHandlers(services *Services) *Handlers {
	return &Handlers{
		Auth: handler.NewAuthHandler(services.Auth, services.User),
		User: handler.NewUserHandler(services.User),
	}
}
