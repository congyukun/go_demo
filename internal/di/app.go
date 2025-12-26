// Package di 提供简洁的依赖注入配置
package di

import (
	"go_demo/internal/config"
	"go_demo/internal/handler"
	"go_demo/internal/repository"
	"go_demo/internal/service"
	"go_demo/pkg/cache"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ServerApp 服务器应用包装器，包含引擎和清理函数
type ServerApp struct {
	Engine  *gin.Engine
	Cleanup func()
}

// AppDependencies 应用依赖聚合器 // di.AppDependencies
type AppDependencies struct {
	Config     *config.Config       // di.AppDependencies.Config
	DB         *gorm.DB             // di.AppDependencies.DB
	Cache      cache.CacheInterface // di.AppDependencies.Cache
	Repository *Repository          // di.AppDependencies.Repository
	Services   *Services            // di.AppDependencies.Services
	Handlers   *Handlers            // di.AppDependencies.Handlers
}

// Repository 仓储层聚合器 // di.Repository
type Repository struct {
	User repository.UserRepository // di.Repository.User
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

// NewRepository 创建仓储聚合器 // di.NewRepository()
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User: repository.NewUserRepository(db),
	}
}

// NewServices 创建服务聚合器 // di.NewServices()
func NewServices(repo *Repository) *Services {
	return &Services{
		Auth: service.NewAuthService(repo.User),
		User: service.NewUserService(repo.User),
	}
}

// NewHandlers 创建处理器聚合器 // di.NewHandlers()
func NewHandlers(services *Services) *Handlers {
	return &Handlers{
		Auth: handler.NewAuthHandler(services.Auth, services.User),
		User: handler.NewUserHandler(services.User),
	}
}
