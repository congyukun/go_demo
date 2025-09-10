# 项目架构文档

## 概述

Go Demo 是一个采用标准 Go 项目布局的 Web 应用程序，实现了用户管理和认证功能。项目遵循分层架构设计原则，具有良好的可维护性和可扩展性。

## 架构设计

### 分层架构

```
┌─────────────────┐
│   Handler 层    │  ← HTTP 请求处理、参数验证
├─────────────────┤
│   Service 层    │  ← 业务逻辑处理
├─────────────────┤
│  Repository 层  │  ← 数据访问抽象
├─────────────────┤
│   Database 层   │  ← 数据持久化
└─────────────────┘
```

### 目录结构说明

#### `/cmd`
应用程序入口点。每个应用程序的 main 函数都应该放在这里。

- `cmd/server/main.go`: 主服务器入口

#### `/internal`
私有应用程序和库代码。这是不希望其他应用程序或库导入的代码。

- `internal/config/`: 配置管理
- `internal/handler/`: HTTP 处理器（控制器层）
- `internal/service/`: 业务逻辑层
- `internal/repository/`: 数据访问层
- `internal/models/`: 数据模型定义

#### `/pkg`
可以被外部应用程序使用的库代码。

- `pkg/database/`: 数据库连接工具
- `pkg/logger/`: 日志工具

#### `/configs`
配置文件模板或默认配置。

#### `/api`
OpenAPI/Swagger 规范，JSON 模式文件，协议定义文件。

#### `/docs`
设计和用户文档。

#### `/scripts`
执行各种构建，安装，分析等操作的脚本。

#### `/deployments`
IaaS，PaaS，系统和容器编排部署配置和模板。

#### `/tests`
额外的外部测试应用程序和测试数据。

## 核心组件

### 配置管理 (internal/config)

负责应用程序配置的加载和管理：

- 支持 YAML 格式配置文件
- 环境变量覆盖
- 默认值设置
- 多环境配置支持

### 数据模型 (internal/models)

定义应用程序的数据结构：

- `User`: 用户实体
- `LoginRequest/Response`: 登录相关数据传输对象
- `RegisterRequest`: 注册请求数据传输对象
- `UpdateUserRequest`: 用户更新请求数据传输对象

### 数据访问层 (internal/repository)

提供数据访问的抽象接口：

- `UserRepository`: 用户数据访问接口
- 支持 CRUD 操作
- 分页查询支持
- 错误处理

### 业务逻辑层 (internal/service)

实现核心业务逻辑：

- `AuthService`: 认证服务（登录、注册、token 验证）
- `UserService`: 用户管理服务
- 业务规则验证
- 数据转换

### HTTP 处理层 (internal/handler)

处理 HTTP 请求和响应：

- `AuthHandler`: 认证相关接口
- `UserHandler`: 用户管理接口
- 统一响应格式
- 错误处理
- 中间件支持

### 中间件 (internal/handler/middleware.go)

提供横切关注点：

- 请求 ID 生成
- CORS 支持
- 认证验证
- 日志记录
- 错误恢复

## 数据流

### 用户注册流程

```
1. HTTP Request → AuthHandler.Register()
2. 参数验证 → models.RegisterRequest
3. 调用 AuthService.Register()
4. 检查用户名/邮箱唯一性
5. 创建用户 → UserRepository.Create()
6. 返回用户信息 → models.UserResponse
```

### 用户登录流程

```
1. HTTP Request → AuthHandler.Login()
2. 参数验证 → models.LoginRequest
3. 调用 AuthService.Login()
4. 验证用户凭据 → UserRepository.GetByUsername()
5. 生成 Token
6. 返回登录信息 → models.LoginResponse
```

### 用户查询流程

```
1. HTTP Request → UserHandler.GetUsers()
2. 认证中间件验证
3. 调用 UserService.GetUsers()
4. 分页查询 → UserRepository.List()
5. 数据转换 → models.UserResponse
6. 返回用户列表
```

## 技术选型

### 核心框架
- **Gin**: 高性能 HTTP Web 框架
- **GORM**: ORM 框架，支持多种数据库
- **Zap**: 高性能结构化日志库

### 数据库
- **MySQL**: 主数据库
- **Redis**: 缓存（可选）

### 配置管理
- **Viper**: 配置管理库（通过 YAML）

### 部署
- **Docker**: 容器化部署
- **Docker Compose**: 多服务编排

## 设计原则

### 1. 单一职责原则
每个组件都有明确的职责：
- Handler 只负责 HTTP 请求处理
- Service 只负责业务逻辑
- Repository 只负责数据访问

### 2. 依赖倒置原则
- Service 层依赖 Repository 接口，而不是具体实现
- 便于单元测试和模块替换

### 3. 开闭原则
- 通过接口设计，支持功能扩展
- 新增功能不需要修改现有代码

### 4. 接口隔离原则
- 定义小而专一的接口
- 避免接口污染

## 安全考虑

### 认证机制
- Token 基础认证（可扩展为 JWT）
- 密码哈希存储
- 用户状态管理

### 输入验证
- 请求参数验证
- SQL 注入防护（通过 ORM）
- XSS 防护

### 日志安全
- 敏感信息脱敏
- 请求追踪
- 错误日志记录

## 性能优化

### 数据库优化
- 连接池配置
- 索引优化
- 慢查询监控
- 集成zap日志系统，统一SQL执行日志
- 自定义GORM日志适配器，支持结构化日志

### 缓存策略
- Redis 缓存支持
- 查询结果缓存
- 会话缓存

### 日志优化
- 异步日志写入
- 日志轮转
- 结构化日志

## 监控和运维

### 健康检查
- `/health` 端点
- 数据库连接检查
- 服务状态监控

### 日志管理
- 结构化日志输出
- 日志级别控制
- 日志文件轮转

### 指标监控
- HTTP 请求指标
- 数据库连接指标
- 业务指标

## 扩展性

### 水平扩展
- 无状态设计
- 负载均衡支持
- 数据库读写分离

### 功能扩展
- 插件化架构
- 中间件机制
- 配置驱动

### 集成扩展
- RESTful API 设计
- 标准化响应格式
- 版本控制支持

## 最佳实践

### 代码规范
- Go 官方代码规范
- 统一的错误处理
- 完整的单元测试

### 部署规范
- 容器化部署
- 环境隔离
- 配置外部化

### 运维规范
- 日志标准化
- 监控告警
- 备份策略

## 未来规划

### 短期目标
- 完善单元测试覆盖率
- 添加集成测试
- 性能基准测试

### 中期目标
- 微服务拆分
- 消息队列集成
- 缓存优化

### 长期目标
- 云原生部署
- 服务网格集成
- 自动化运维