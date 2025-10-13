// Package handler provides configuration aggregation for handlers.
package handler

import "go_demo/internal/service"

// HandlerConfig 聚合 handler 层依赖 // handler.HandlerConfig
type HandlerConfig struct {
	AuthService service.AuthService // handler.HandlerConfig.AuthService
	UserService service.UserService // handler.HandlerConfig.UserService
}

// NewHandlerConfig 构造 HandlerConfig // handler.NewHandlerConfig()
func NewHandlerConfig(authSvc service.AuthService, userSvc service.UserService) *HandlerConfig {
	return &HandlerConfig{
		AuthService: authSvc,
		UserService: userSvc,
	}
}

// NewAuthHandlerWithConfig 基于配置创建 AuthHandler // handler.NewAuthHandlerWithConfig()
func NewAuthHandlerWithConfig(cfg *HandlerConfig) *AuthHandler {
	return NewAuthHandler(cfg.AuthService, cfg.UserService)
}

// NewUserHandlerWithConfig 基于配置创建 UserHandler // handler.NewUserHandlerWithConfig()
func NewUserHandlerWithConfig(cfg *HandlerConfig) *UserHandler {
	return NewUserHandler(cfg.UserService)
}