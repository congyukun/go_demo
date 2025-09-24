# Go Demo 项目

一个标准的 Go Web 应用项目，采用分层架构设计，包含用户管理和认证功能。

## 🚀 项目特性

- ✅ **标准项目结构**: 遵循 Go 项目布局标准
- ✅ **分层架构**: Handler -> Service -> Repository 清晰分层
- ✅ **用户管理**: 完整的用户 CRUD 操作
- ✅ **认证系统**: 登录/注册/登出功能
- ✅ **配置管理**: 支持多环境配置
- ✅ **日志系统**: 结构化日志记录，MySQL集成zap日志
- ✅ **数据库支持**: MySQL/PostgreSQL/SQLite/MongoDB
- ✅ **缓存支持**: Redis 集成
- ✅ **API 文档**: OpenAPI 3.0 规范
- ✅ **容器化**: Docker 和 Docker Compose 支持
- ✅ **健康检查**: 服务健康状态监控

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
│   └── models/           # 数据模型
├── pkg/                   # 可重用的库代码
│   ├── database/         # 数据库连接
│   └── logger/           # 日志工具
├── configs/              # 配置文件
├── api/                  # API 文档
├── docs/                 # 项目文档
├── scripts/              # 脚本文件
├── tests/                # 测试文件
├── deployments/          # 部署配置
├── logs/                 # 日志文件（运行时生成）
└── data/                 # 数据文件（运行时生成）
```

## 🛠️ 技术栈

- **Go 1.19+**: 编程语言
- **Gin**: Web 框架
- **GORM**: ORM 框架
- **MySQL**: 主数据库
- **Redis**: 缓存数据库
- **Zap**: 结构化日志
- **YAML**: 配置文件格式
- **Docker**: 容器化部署

## 🚀 快速开始

### 环境要求

- Go 1.19 或更高版本
- MySQL 5.7 或更高版本
- Redis（可选）

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

### 4. 运行应用

```bash
# 开发环境运行
go run cmd/server/main.go

# 或者构建后运行
go build -o bin/server cmd/server/main.go
./bin/server
```

应用将在 `http://localhost:8080` 启动。

### 5. 验证服务

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

### 认证接口

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/v1/auth/register` | 用户注册 |
| POST | `/api/v1/auth/login` | 用户登录 |
| POST | `/api/v1/auth/logout` | 用户登出 |

### 用户管理接口

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | `/api/v1/users` | 获取用户列表 | ✅ |
| GET | `/api/v1/users/:id` | 获取用户详情 | ✅ |
| PUT | `/api/v1/users/:id` | 更新用户信息 | ✅ |
| DELETE | `/api/v1/users/:id` | 删除用户 | ✅ |

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

### 代码规范

- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 添加必要的注释和文档
- 完整的错误处理
- 详细的日志记录

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

可以通过环境变量 `CONFIG_PATH` 指定配置文件路径。

## 📊 监控和日志

### 日志

- 使用 Zap 结构化日志
- 支持控制台和文件输出
- 自动日志轮转
- 多级别日志记录

### 健康检查

访问 `/health` 端点获取服务健康状态。

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
