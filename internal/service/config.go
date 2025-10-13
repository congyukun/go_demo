package service

import "go_demo/internal/repository"

// ServiceConfig 用于聚合 service 层所需依赖 // service.ServiceConfig
type ServiceConfig struct {
	UserRepo repository.UserRepository // service.ServiceConfig.UserRepo
}

// NewServiceConfig 构造 ServiceConfig // service.NewServiceConfig()
func NewServiceConfig(userRepo repository.UserRepository) *ServiceConfig {
	return &ServiceConfig{UserRepo: userRepo}
}

// NewAuthServiceWithConfig 基于配置创建 AuthService // service.NewAuthServiceWithConfig()
func NewAuthServiceWithConfig(cfg *ServiceConfig) AuthService {
	return NewAuthService(cfg.UserRepo)
}

// NewUserServiceWithConfig 基于配置创建 UserService // service.NewUserServiceWithConfig()
func NewUserServiceWithConfig(cfg *ServiceConfig) UserService {
	return NewUserService(cfg.UserRepo)
}