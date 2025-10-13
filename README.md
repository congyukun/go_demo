# Go Demo 项目

一个标准的 Go Web 应用项目，采用分层架构设计，包含用户管理和认证功能。

## 🚀 项目特性

- ✅ **标准项目结构**: 遵循 Go 项目布局标准
- ✅ **分层架构**: Handler -> Service -> Repository 清晰分层
- ✅ **用户管理**: 完整的用户 CRUD 操作
- ✅ **认证系统**: 登录/注册/登出功能
- ✅ **配置管理**: 支持多环境配置
- ✅ **日志系统**: 结构化日志记录，支持文件和控制台输出
- ✅ **数据库支持**: MySQL (已配置，可扩展支持PostgreSQL)
- ✅ **缓存支持**: Redis 集成
- ✅ **限流系统**: 分布式限流中间件
- ✅ **API 文档**: OpenAPI 3.0 规范
- ✅ **容器化**: Docker 和 Docker Compose 支持
- ✅ **健康检查**: 服务健康状态监控
- ✅ **依赖注入**: Google Wire 编译期依赖注入
- ✅ **熔断保护**: 服务熔断和降级机制

## 📁 项目结构

```
go_demo/
├── cmd/                    # 应用程序入口
│   └── server/
│       └── main.go        # 主程序入口
├── internal/              # 内部应用代码
│   ├── config/           # 配置管理
│   ├── handler/          # HTTP 处理器（控制器层）
│   ├── service/          # 业务逻辑层
│   ├── repository/       # 数据访问层
│   ├── models/           # 数据模型
│   ├── middleware/       # 中间件（限流、认证等）
│   └── di/               # 依赖注入（Wire）
├── pkg/                   # 可重用的库代码
│   ├── cache/            # Redis缓存封装
│   ├── database/         # 数据库连接
│   ├── errors/           # 错误处理
│   ├── logger/           # 日志工具
│   └── validator/        # 参数验证
├── configs/              # 配置文件
├── api/                  # API 文档（OpenAPI规范）
├── docs/                 # 项目文档
│   ├── ARCHITECTURE.md   # 架构文档
│   └── TECH_SUMMARY.md   # 技术特性总结
├── scripts/              # 脚本文件（构建、部署、迁移）
├── tests/                # 测试文件
├── deployments/          # 部署配置（Docker、Nginx）
├── logs/                 # 日志文件（运行时生成）
├── go.mod                # Go模块定义
├── go.sum               # Go依赖校验
├── Makefile             # 构建脚本
└── API.md               # API使用文档
```

## 🛠️ 技术栈

- **Go 1.24**: 编程语言
- **Gin**: Web 框架
- **GORM**: ORM 框架
- **MySQL**: 主数据库
- **Redis**: 缓存数据库 + 分布式限流
- **Zap**: 结构化日志
- **Viper**: 配置管理
- **JWT**: 无状态认证
- **Wire**: 依赖注入
- **Docker**: 容器化部署
- **Swagger**: API文档生成

## 🚀 快速开始

### 环境要求

- Go 1.24 或更高版本
- MySQL 5.7 或更高版本
- Redis 5.0+ (必需)

### 1. 克隆项目

```bash
git clone <repository-url>
cd go_demo
```

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 配置数据库

创建数据库：
```sql
CREATE DATABASE go_test CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

修改配置文件 `configs/config.yaml` 中的数据库连接信息。

### 4. 配置Redis

确保Redis服务已启动，修改配置文件：
```yaml
redis:
  addr: "localhost:6379"
  password: ""
  db: 0
  pool_size: 10
```

### 5. 运行应用

```bash
# 开发环境运行
go run cmd/server/main.go

# 或者构建后运行
go build -o bin/server cmd/server/main.go
./bin/server
```

应用将在 `http://localhost:8080` 启动。

### 6. 验证服务

```bash
# 健康检查
curl http://localhost:8080/health

# 用户注册
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456",
    "email": "test@example.com",
    "name": "Test User"
  }'

# 用户登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "123456"
  }'
```

## 📚 API 文档

### Swagger UI 文档

项目已集成 Swagger 文档，启动服务后可通过以下方式访问：

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Swagger JSON**: http://localhost:8080/swagger/doc.json

### 认证接口

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/v1/auth/register` | 用户注册 |
| POST | `/api/v1/auth/login` | 用户登录 |
| POST | `/api/v1/auth/refresh` | 刷新访问令牌 |
| POST | `/api/v1/auth/logout` | 用户登出 |
| GET  | `/api/v1/auth/profile` | 获取当前用户信息 |

### 用户管理接口

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | `/api/v1/users` | 获取用户列表 | ✅ |
| POST | `/api/v1/users` | 创建新用户 | ✅ |
| GET | `/api/v1/users/:id` | 获取用户详情 | ✅ |
| PUT | `/api/v1/users/:id` | 更新用户信息 | ✅ |
| DELETE | `/api/v1/users/:id` | 删除用户 | ✅ |
| PUT | `/api/v1/users/profile` | 更新当前用户资料 | ✅ |
| PUT | `/api/v1/users/Password` | 修改当前用户密码 | ✅ |
| GET | `/api/v1/users/stats` | 获取用户统计信息 | ✅ |

### 限流配置

系统支持多级限流配置：

- **全局限流**: 100请求/分钟/IP
- **API限流**: 可针对特定API配置
- **用户限流**: 基于用户ID的个性化限流

### 使用 Swagger 文档

1. **启动服务**:
   ```bash
   go run cmd/server/main.go
   ```

2. **访问 Swagger UI**:
   打开浏览器访问: http://localhost:8080/swagger/index.html

3. **认证测试**:
   - 使用 `/api/v1/auth/register` 注册新用户
   - 使用 `/api/v1/auth/login` 登录获取 JWT token
   - 点击 Swagger UI 右上角的 "Authorize" 按钮
   - 输入格式: `Bearer <your_jwt_token>`

4. **生成/更新文档**:
   ```bash
   # 安装 swag 工具
   go install github.com/swaggo/swag/cmd/swag@latest
   
   # 生成文档
   swag init -g cmd/server/main.go
   ```

详细的 API 文档请查看 [OpenAPI 规范](api/openapi.yaml)。

## 🐳 Docker 部署

### 使用 Docker Compose（推荐）

```bash
# 启动所有服务
cd deployments
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f app

# 停止服务
docker-compose down
```

### 单独使用 Docker

```bash
# 构建镜像
docker build -f deployments/Dockerfile -t go-demo .

# 运行容器
docker run -p 8080:8080 go-demo
```

## 🧪 测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/service

# 运行测试并显示覆盖率
go test -cover ./...

# 运行单元测试
go test ./tests -v
```

## 📝 开发指南

### 项目架构

项目采用经典的分层架构：

1. **Handler 层**: 处理 HTTP 请求，参数验证，调用 Service 层
2. **Service 层**: 业务逻辑处理，调用 Repository 层
3. **Repository 层**: 数据访问，与数据库交互
4. **Model 层**: 数据模型定义

### 新增功能开发指南

#### 1. 限流中间件使用

```go
// 使用默认限流配置
router.Use(middleware.RateLimiter(middleware.DefaultRateLimiterConfig()))

// 自定义限流配置
config := middleware.RateLimiterConfig{
    Window:      time.Minute,
    MaxRequests: 100,
    KeyGenerator: func(c *gin.Context) string {
        return "custom:" + c.ClientIP()
    },
}
router.Use(middleware.RateLimiter(config))
```

#### 2. 缓存操作示例

```go
// 获取缓存实例
cache := pkgcache.NewRedisCache(redisConfig)

// 设置缓存
err := cache.Set("user:1", userData, time.Hour)
if err != nil {
    log.Printf("设置缓存失败: %v", err)
}

// 获取缓存
var user User
err = cache.GetObject("user:1", &user)
if err != nil {
    log.Printf("获取缓存失败: %v", err)
}

// 删除缓存
err = cache.Delete("user:1")
if err != nil {
    log.Printf("删除缓存失败: %v", err)
}
```

### 代码规范

- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 添加必要的注释和文档
- 完整的错误处理
- 详细的日志记录
- 使用 Wire 进行依赖注入
- 遵循 Clean Architecture 原则

### 添加新功能

1. 在 `internal/models` 中定义数据模型
2. 在 `internal/repository` 中实现数据访问
3. 在 `internal/service` 中实现业务逻辑
4. 在 `internal/handler` 中实现 HTTP 处理
5. 在 `cmd/server/main.go` 中注册路由

## 🔧 配置说明

配置文件位于 `configs/config.yaml`，支持以下配置：

- **app**: 应用基础配置
- **server**: 服务器配置
- **database**: 数据库配置（支持多种数据库）
- **redis**: Redis 配置
- **log**: 日志配置
- **rate_limit**: 限流配置

可以通过环境变量 `CONFIG_PATH` 指定配置文件路径。

### 限流配置示例

```yaml
rate_limiter:
  enabled: true
  global_limit: 1000
  user_limit: 100
  ip_limit: 200
  window: 60
  algorithm: "sliding"
```

## 📊 监控和日志

### 日志系统

- 使用 Zap 结构化日志
- 支持控制台和文件输出
- 自动日志轮转
- 多级别日志记录

### 健康检查

访问 `/health` 端点获取服务健康状态。

### 限流监控

系统提供以下监控指标：
- `rate_limiter_requests_total`: 总请求数
- `rate_limiter_rejected_total`: 被拒绝请求数
- `rate_limiter_allowed_total`: 被允许请求数
- `redis_connections_active`: Redis活跃连接数

## 📖 项目文档

### 文档索引
- [🏗️ 架构文档](docs/ARCHITECTURE.md) - 系统架构详细说明
- [📊 技术特性](docs/TECH_SUMMARY.md) - 核心特性总结
- [📖 Swagger指南](docs/SWAGGER_UPDATE.md) - API文档使用指南

### 快速导航
- [API文档](api/openapi.yaml) - OpenAPI 3.0规范
- [部署指南](deployments/) - Docker和K8s部署配置
- [测试用例](tests/) - 测试代码和用例

## 🤝 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 📞 联系方式

如有问题或建议，请：
- 提交 Issue
- 发送邮件至维护者
- 参与讨论

---

**注意**: 这是一个演示项目，生产环境使用前请进行适当的安全配置和性能优化。

**最近更新**: 2025-10-13 - 新增分布式限流系统和Redis缓存支持
