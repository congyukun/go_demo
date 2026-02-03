# Go Demo 项目

一个标准的 Go Web 应用项目，采用分层架构设计，包含用户管理和认证功能，并配套 Vue 3 后台管理系统。

## 🚀 项目特性

### 后端特性
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

### 前端特性
- ✅ **Vue 3**: 使用 Composition API
- ✅ **Vite**: 下一代前端构建工具
- ✅ **Element Plus**: Vue 3 UI 组件库
- ✅ **Pinia**: Vue 状态管理
- ✅ **Vue Router**: 路由管理 + 权限守卫
- ✅ **Axios**: HTTP 请求封装
- ✅ **响应式布局**: 可折叠侧边栏

## 📁 项目结构

```
go_demo/
├── cmd/                    # 应用程序入口
│   └── server/
│       └── init.go        # 服务初始化
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
├── web/                   # 🌐 Vue 3 前端项目
│   ├── src/
│   │   ├── api/          # API 接口封装
│   │   ├── assets/       # 静态资源
│   │   ├── layout/       # 布局组件
│   │   ├── router/       # 路由配置
│   │   ├── stores/       # Pinia 状态管理
│   │   ├── styles/       # 全局样式
│   │   ├── utils/        # 工具函数
│   │   └── views/        # 页面组件
│   │       ├── login/    # 登录页
│   │       ├── register/ # 注册页
│   │       ├── dashboard/# 仪表盘
│   │       ├── users/    # 用户管理
│   │       └── profile/  # 个人中心
│   ├── package.json      # 前端依赖
│   └── vite.config.js    # Vite 配置
├── configs/              # 配置文件
├── api/                  # API 文档（OpenAPI规范）
├── docs/                 # 项目文档
│   ├── DEPLOYMENT.md     # 部署文档
│   ├── DEPLOYMENT_OPTIMIZATION.md  # 部署优化
│   └── DOCKER_GUIDE.md   # Docker 指南
├── scripts/              # 脚本文件（构建、部署、迁移）
├── tests/                # 测试文件
├── deployments/          # 部署配置（Docker、Nginx）
├── logs/                 # 日志文件（运行时生成）
├── go.mod                # Go模块定义
├── go.sum               # Go依赖校验
├── Makefile             # 构建脚本
├── main.go              # 主程序入口
└── API.md               # API使用文档
```

## 🛠️ 技术栈

### 后端技术栈
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

### 前端技术栈
- **Vue 3**: 渐进式 JavaScript 框架
- **Vite**: 下一代前端构建工具
- **Vue Router**: 官方路由管理
- **Pinia**: Vue 状态管理库
- **Element Plus**: Vue 3 UI 组件库
- **Axios**: HTTP 请求库
- **Sass**: CSS 预处理器
- **Swagger**: API文档生成

## 🚀 快速开始

### 环境要求

- Go 1.24 或更高版本
- Node.js 16+ (前端)
- MySQL 5.7 或更高版本
- Redis 5.0+ (必需)

### 1. 克隆项目

```bash
git clone <repository-url>
cd go_demo
```

### 2. 安装依赖

```bash
# 安装所有依赖（后端 + 前端）
make install-all

# 或者分别安装
go mod tidy           # 后端依赖
cd web && npm install # 前端依赖
```
### 3. 配置环境

项目支持多环境配置，通过 `--config` 参数指定配置文件：

```bash
# 复制环境变量示例文件（可选）
cp .env.example .env
vim .env
```

配置文件说明：
- [`config.yaml`](configs/config.yaml) - 默认配置（Docker/生产环境）
- [`config.dev.yaml`](configs/config.dev.yaml) - 开发环境配置
- [`config.docker.yaml`](configs/config.docker.yaml) - Docker 环境配置

### 4. 配置数据库

创建数据库：
```sql
CREATE DATABASE go_demo CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 5. 运行应用

#### 方式一：使用 Makefile（推荐）

```bash
# 🚀 同时启动前后端开发服务器
make dev-all

# 或者分别启动
make dev      # 启动后端 (http://localhost:8080)
make web-dev  # 启动前端 (http://localhost:3000)
```

#### 方式二：手动启动

```bash
# 启动后端
go run main.go server --config=./configs/config.dev.yaml

# 启动前端（新终端）
cd web && npm run dev
```

## 📋 Makefile 命令参考

> 💡 **提示**: 运行 `make help` 可查看所有可用命令的完整帮助信息

### 基础命令

| 命令 | 描述 |
|------|------|
| `make help` | 📖 显示所有命令的帮助信息 |
| `make deps` | 安装 Go 依赖 |
| `make fmt` | 格式化代码 |
| `make vet` | 代码检查 |
| `make test` | 运行测试 |
| `make test-coverage` | 生成测试覆盖率报告 |
| `make build` | 构建应用 |
| `make build-all` | 构建多平台版本 |
| `make run` | 运行应用 |
| `make dev` | 开发模式运行（带热重载） |
| `make clean` | 清理构建文件 |
| `make lint` | 代码质量检查 |
| `make docs` | 生成 API 文档 |
| `make migrate` | 数据库迁移 |
| `make health` | 健康检查 |
| `make install-tools` | 安装开发工具（air, golangci-lint） |

### 前端命令

| 命令 | 描述 |
|------|------|
| `make web-install` | 安装前端依赖 |
| `make web-dev` | 启动前端开发服务器 |
| `make web-build` | 构建前端 |
| `make web-preview` | 预览前端构建 |
| `make web-lint` | 前端代码检查 |
| `make web-clean` | 清理前端构建 |

### 全栈开发命令

| 命令 | 描述 |
|------|------|
| `make dev-all` | 🚀 同时启动前后端开发服务器 |
| `make install-all` | 安装所有依赖（后端 + 前端） |
| `make build-all-stack` | 构建前后端 |

### Docker 命令

| 命令 | 描述 |
|------|------|
| `make docker-build` | 构建 Docker 镜像 |
| `make docker-run` | 运行 Docker 容器 |
| `make docker-deploy` | 完整部署（应用 + MySQL + Redis + Nginx） |
| `make docker-deploy-simple` | 简化部署（应用 + MySQL + Redis） |
| `make docker-deps` | 仅启动依赖服务（MySQL + Redis） |
| `make docker-up` | 启动 Docker Compose |
| `make docker-down` | 停止 Docker Compose |
| `make docker-stop` | 停止所有服务 |
| `make docker-status` | 查看服务状态 |
| `make docker-logs` | 查看应用日志 |
| `make docker-logs-all` | 查看所有服务日志 |
| `make docker-restart` | 重启应用 |
| `make docker-restart-all` | 重启所有服务 |
| `make docker-clean` | 清理所有数据（危险操作） |
| `make docker-info` | 显示服务信息 |

### Podman 命令

| 命令 | 描述 |
|------|------|
| `make podman-build` | 构建 Podman 镜像 |
| `make podman-run` | 运行 Podman 容器 |
| `make podman-deploy` | Podman 完整部署 |
| `make podman-deploy-simple` | Podman 简化部署 |
| `make podman-deps` | Podman 启动依赖服务 |
| `make podman-up` | 启动 Podman Compose |
| `make podman-down` | 停止 Podman Compose |
| `make podman-stop` | 停止所有服务 |
| `make podman-status` | 查看服务状态 |
| `make podman-logs` | 查看应用日志 |
| `make podman-logs-all` | 查看所有服务日志 |
| `make podman-restart` | 重启应用 |
| `make podman-restart-all` | 重启所有服务 |
| `make podman-clean` | 清理所有数据 |
| `make podman-info` | 显示服务信息 |

服务地址：
- **后端 API**: http://localhost:8080
- **前端界面**: http://localhost:3000
- **Swagger 文档**: http://localhost:8080/swagger/index.html

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

## 🔧 配置管理

### 多环境配置

项目支持多环境配置管理，通过 `--config` 参数指定配置文件：

```bash
# 默认配置（不指定参数）
go run main.go server

# 开发环境
go run main.go server --config=./configs/config.dev.yaml

# 测试环境
go run main.go server --config=./configs/config.test.yaml

# 生产环境
./go_demo server --config=./configs/config.prod.yaml
```

### 配置文件结构

```
configs/
├── config.yaml           # 默认配置
├── config.dev.yaml       # 开发环境配置 ✅ 提交
├── config.docker.yaml    # Docker 环境配置 ✅ 提交
```

### 环境变量支持

所有配置项都支持通过环境变量覆盖，命名规则：`GO_DEMO_<SECTION>_<KEY>`

```bash
# 覆盖服务器端口
export GO_DEMO_SERVER_PORT=9090

# 覆盖数据库连接
export GO_DEMO_DATABASE_DSN="root:pass@tcp(localhost:3306)/go_demo"

# 覆盖 JWT 密钥
export GO_DEMO_JWT_secret_KEY="your-secret-key"
```

### 配置优先级

1. **环境变量** - 最高优先级
2. **配置文件** - 中等优先级
3. **默认值** - 最低优先级

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
- [🐳 Docker部署指南](docs/DOCKER_GUIDE.md) - Docker部署说明
- [🚀 部署文档](docs/DEPLOYMENT.md) - 部署指南
- [⚡ 部署优化](docs/DEPLOYMENT_OPTIMIZATION.md) - 部署优化建议

### 快速导航
- [API文档](api/openapi.yaml) - OpenAPI 3.0规范
- [Swagger文档](docs/swagger.yaml) - Swagger API 文档
- [前端项目](web/) - Vue 3 后台管理系统
- [部署配置](deployments/) - Docker 部署配置
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
- 发送邮件至维护者 cunliwakun@163.com
- 参与讨论

---

**注意**: 这是一个演示项目，生产环境使用前请进行适当的安全配置和性能优化。

**最近更新**:
- 2026-02-03 - 新增 Vue 3 后台管理系统前端，集成 Makefile 前端命令
- 2025-12-26 - 完善多环境配置管理，支持环境变量覆盖
- 2025-10-13 - 新增分布式限流系统和Redis缓存支持
