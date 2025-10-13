package di

import (
	"go_demo/internal/config"
	"go_demo/internal/handler"
	"go_demo/internal/repository"
	"go_demo/internal/router"
	"go_demo/internal/service"
	"go_demo/internal/utils"
	"go_demo/pkg/database"

	"github.com/google/wire"
)

// ProviderSet 定义所有依赖提供者集合（仅声明依赖关系，不包含注入器）
var ProviderSet = wire.NewSet(
	// 配置
	config.NewConfig,
	wire.FieldsOf(new(*config.Config), "JWT", "Database"), // 提供 utils.JWTConfig 和 database.MySQLConfig

	// 工具
	utils.NewJWTManager,

	// 基础设施
	database.NewMySQL, // *gorm.DB

	// 仓储
	repository.NewUserRepository,

	// 服务
	service.NewAuthService,
	service.NewUserService,

	// 处理器
	handler.NewAuthHandler,
	handler.NewUserHandler,

	// 路由
	router.NewRouter,
)

// AppContainer 应用容器结构
type AppContainer struct {
	Config   *config.Config
	JWT      *utils.JWTManager
	UserRepo repository.UserRepository
	AuthSvc  service.AuthService
	UserSvc  service.UserService
	AuthHdl  *handler.AuthHandler
	UserHdl  *handler.UserHandler
	Router   *router.Router
}

// NewAppContainer 创建应用容器（由 Wire 组装调用）
func NewAppContainer(
	cfg *config.Config,
	jwt *utils.JWTManager,
	userRepo repository.UserRepository,
	authSvc service.AuthService,
	userSvc service.UserService,
	authHdl *handler.AuthHandler,
	userHdl *handler.UserHandler,
	rt *router.Router,
) *AppContainer {
	return &AppContainer{
		Config:   cfg,
		JWT:      jwt,
		UserRepo: userRepo,
		AuthSvc:  authSvc,
		UserSvc:  userSvc,
		AuthHdl:  authHdl,
		UserHdl:  userHdl,
		Router:   rt,
	}
}
