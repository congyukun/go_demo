# Go Demo 项目

一个标准的 Go Web 应用项目，采用分层架构设计，包含用户管理和认证功能。

## 项目结构

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
├── logs/                 # 日志文件
└── data/                 # 数据文件
```

## 技术栈

- **Go 1.19+**: 编程语言
- **Gin**: Web 框架
- **GORM**: ORM 框架
- **MySQL**: 数据库
- **Zap**: 日志库
- **YAML**: 配置文件格式

## 功能特性

- ✅ 用户注册/登录/登出
- ✅ 用户信息管理（CRUD）
- ✅ JWT 认证（简化版）
- ✅ 请求日志记录
- ✅ 统一错误处理
- ✅ 配置管理
- ✅ 数据库连接池
- ✅ API 文档
- ✅ 健康检查

## 快速开始

### 环境要求

- Go 1.19 或更高版本
- MySQL 5.7 或更高版本

### 安装依赖

```bash
go mod tidy
```

### 配置数据库

1. 创建数据库：
```sql
CREATE DATABASE go_demo CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

2. 修改配置文件 `configs/config.yaml` 中的数据库连接信息

### 运行应用

```bash
# 开发环境
go run cmd/server/main.go

# 或者构建后运行
go build -o bin/server cmd/server/main.go
./bin/server
```

应用将在 `http://localhost:8080` 启动

### 健康检查

```bash
curl http://localhost:8080/health
```

## API 文档

详细的 API 文档请查看：
- [OpenAPI 规范](../api/openapi.yaml)
- 启动应用后访问：`http://localhost:8080/api/v1`

### 主要接口

#### 认证相关
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/logout` - 用户登出

#### 用户管理
- `GET /api/v1/users` - 获取用户列表
- `GET /api/v1/users/:id` - 获取用户详情
- `PUT /api/v1/users/:id` - 更新用户信息
- `DELETE /api/v1/users/:id` - 删除用户

## 开发指南

### 项目架构

项目采用经典的分层架构：

1. **Handler 层**: 处理 HTTP 请求，参数验证，调用 Service 层
2. **Service 层**: 业务逻辑处理，调用 Repository 层
3. **Repository 层**: 数据访问，与数据库交互
4. **Model 层**: 数据模型定义

### 代码规范

- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 添加必要的注释
- 错误处理要完整
- 日志记录要详细

### 配置管理

配置文件位于 `configs/config.yaml`，支持：
- 应用基础配置
- 服务器配置
- 数据库配置
- 日志配置
- Redis 配置

### 日志管理

使用 Zap 日志库，支持：
- 结构化日志
- 日志轮转
- 多级别日志
- 控制台和文件输出

### 数据库迁移

使用 GORM 的 AutoMigrate 功能：

```go
db.AutoMigrate(&models.User{})
```

## 部署

### Docker 部署

```bash
# 构建镜像
docker build -t go-demo .

# 运行容器
docker run -p 8080:8080 go-demo
```

### 生产环境配置

1. 修改配置文件中的环境为 `prod`
2. 配置生产数据库连接
3. 设置适当的日志级别
4. 配置反向代理（如 Nginx）

## 测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/service

# 运行测试并显示覆盖率
go test -cover ./...
```

## 贡献指南

1. Fork 项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

MIT License

## 联系方式

如有问题，请提交 Issue 或联系维护者。