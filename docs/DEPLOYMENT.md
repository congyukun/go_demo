# Go Demo 项目部署启动指南

## 📋 目录

- [环境要求](#环境要求)
- [快速开始](#快速开始)
- [配置说明](#配置说明)
- [部署方式](#部署方式)
- [常见问题](#常见问题)

## 🔧 环境要求

### 基础环境
- **Go**: 1.24 或更高版本
- **MySQL**: 5.7 或更高版本
- **Redis**: 5.0 或更高版本
- **Docker**: 20.x 或更高版本（可选）
- **Docker Compose**: 2.x 或更高版本（可选）

### 检查环境
```bash
# 检查 Go 版本
go version

# 检查 Docker 版本
docker --version
docker-compose --version
```

## 🚀 快速开始

### 方式一：本地开发环境

#### 1. 克隆项目
```bash
git clone <repository-url>
cd go_demo
```

#### 2. 安装依赖
```bash
go mod tidy
```

#### 3. 启动数据库服务
```bash
# 使用 Docker 启动 MySQL 和 Redis
docker-compose -f deployments/docker-compose.simple.yml up -d mysql redis

# 或手动启动本地 MySQL 和 Redis
```

#### 4. 配置环境
```bash
# 直接使用开发配置启动
go run main.go server --config=./configs/config.dev.yaml
```

#### 5. 初始化数据库
```bash
# 数据库会自动创建表结构（GORM AutoMigrate）
# 或手动执行 SQL
mysql -h localhost -u root -p123456 go_demo < deployments/mysql/init.sql
```

#### 6. 启动应用
```bash
# 使用开发配置启动
go run main.go server --config=./configs/config.dev.yaml

# 或使用 Air 热重载（推荐开发使用）
air
```

#### 7. 验证服务
```bash
# 健康检查
curl http://localhost:8080/health

# 访问 API 文档
open http://localhost:8080/swagger/index.html
```

### 方式二：Docker 部署（推荐）

#### 1. 使用 Docker Compose 一键启动
```bash
# 进入部署目录
cd deployments

# 启动所有服务（包括应用、MySQL、Redis）
docker-compose -f docker-compose.simple.yml up -d

# 查看服务状态
docker-compose -f docker-compose.simple.yml ps

# 查看日志
docker-compose -f docker-compose.simple.yml logs -f app
```

#### 2. 验证服务
```bash
# 健康检查
curl http://localhost:8080/health

# 测试 API
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Test123456",
    "email": "test@example.com",
    "name": "测试用户",
    "mobile": "13800138000"
  }'
```

#### 3. 停止服务
```bash
# 停止服务
docker-compose -f docker-compose.simple.yml down

# 停止并删除数据卷
docker-compose -f docker-compose.simple.yml down -v
```

## ⚙️ 配置说明

### 配置文件

项目提供两个环境配置：

| 文件 | 环境 | 说明 |
|------|------|------|
| [`config.dev.yaml`](../configs/config.dev.yaml) | 开发环境 | 本地开发使用，包含调试配置 |
| [`config.yaml`](../configs/config.yaml) | 生产环境（默认） | 生产部署使用，使用环境变量 |

### 开发环境配置

**特点**：
- 数据库：`localhost:3306`
- 日志级别：`debug`
- SQL 日志：开启
- 连接池：50

**启动方式**：
```bash
go run main.go server --config=./configs/config.dev.yaml
```

### 生产环境配置（默认）

**特点**：
- 使用环境变量注入敏感信息
- 日志级别：`warn`
- SQL 日志：关闭
- 连接池：200

**启动方式**：
```bash
# 1. 设置环境变量
export DATABASE_DSN="user:pass@tcp(host:3306)/go_demo?charset=utf8mb4&parseTime=True&loc=Local"
export JWT_secret_KEY="your-production-secret-key"
export REDIS_HOST="redis-host"
export REDIS_Password="redis-password"

# 2. 启动应用（使用默认配置）
./go_demo server

# 或明确指定配置文件
./go_demo server --config=./configs/config.yaml
```

### 环境变量说明

| 环境变量 | 说明 | 示例 |
|---------|------|------|
| `DATABASE_DSN` | 数据库连接字符串 | `root:pass@tcp(localhost:3306)/go_demo` |
| `JWT_secret_KEY` | JWT 签名密钥 | 至少32位随机字符串 |
| `REDIS_HOST` | Redis 主机地址 | `localhost` 或 `redis` |
| `REDIS_Password` | Redis 密码 | 可选 |

### 生成安全密钥
```bash
# 生成 JWT 密钥
openssl rand -base64 32
```

## 🐳 部署方式

### 1. 本地开发部署

```bash
# 启动依赖服务
docker-compose -f deployments/docker-compose.simple.yml up -d mysql redis

# 启动应用
go run main.go server --config=./configs/config.dev.yaml
```

**访问地址**：
- 应用：http://localhost:8080
- API 文档：http://localhost:8080/swagger/index.html
- MySQL：localhost:3306
- Redis：localhost:6379

### 2. Docker 完整部署

```bash
# 构建并启动所有服务
cd deployments
docker-compose -f docker-compose.simple.yml up -d

# 查看日志
docker-compose -f docker-compose.simple.yml logs -f

# 进入应用容器
docker-compose -f docker-compose.simple.yml exec app sh
```

### 3. 生产环境部署

#### 方式 A：直接部署

```bash
# 1. 构建应用
go build -o go_demo main.go

# 2. 设置环境变量
export DATABASE_DSN="..."
export JWT_secret_KEY="..."
export REDIS_HOST="..."
export REDIS_Password="..."

# 3. 启动应用
./go_demo server --config=./configs/config.prod.yaml
```

#### 方式 B：Docker 部署

```bash
# 1. 构建镜像
docker build -f deployments/Dockerfile -t go-demo:latest .

# 2. 运行容器
docker run -d \
  --name go-demo \
  -p 8080:8080 \
  -e DATABASE_DSN="..." \
  -e JWT_secret_KEY="..." \
  -e REDIS_HOST="..." \
  -e REDIS_Password="..." \
  go-demo:latest
```

#### 方式 C：使用 Nginx 反向代理

```bash
# 1. 启动完整服务栈（包含 Nginx）
cd deployments
docker-compose up -d

# 2. 访问服务
# HTTP: http://localhost
# HTTPS: https://localhost
```

## 🧪 测试部署

### 1. 健康检查
```bash
# 基础健康检查
curl http://localhost:8080/health

# 详细健康检查
curl http://localhost:8080/health/check
```

### 2. 测试 API

#### 用户注册
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Test123456",
    "email": "test@example.com",
    "name": "测试用户",
    "mobile": "13800138000"
  }'
```

#### 用户登录
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Test123456"
  }'
```

#### 获取用户列表（需要认证）
```bash
# 先登录获取 token，然后：
curl -X GET "http://localhost:8080/api/v1/users?page=1&size=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 3. 测试数据库连接
```bash
# 进入 MySQL 容器
docker-compose -f deployments/docker-compose.simple.yml exec mysql mysql -uroot -p123456

# 在 MySQL 中执行
USE go_demo;
SHOW TABLES;
SELECT * FROM users;
```

### 4. 测试 Redis 连接
```bash
# 进入 Redis 容器
docker-compose -f deployments/docker-compose.simple.yml exec redis redis-cli

# 在 Redis 中执行
PING
KEYS *
```

## 🔍 常见问题

### Q1: 端口被占用怎么办？

**A**: 修改配置文件或 docker-compose.yml 中的端口映射

```yaml
# docker-compose.simple.yml
services:
  app:
    ports:
      - "9090:8080"  # 改为 9090
```

### Q2: 数据库连接失败？

**A**: 检查以下几点：
1. MySQL 是否已启动
2. 数据库连接字符串是否正确
3. 数据库是否已创建
4. 用户权限是否正确

```bash
# 检查 MySQL 状态
docker-compose -f deployments/docker-compose.simple.yml ps mysql

# 查看 MySQL 日志
docker-compose -f deployments/docker-compose.simple.yml logs mysql

# 手动创建数据库
docker-compose -f deployments/docker-compose.simple.yml exec mysql \
  mysql -uroot -p123456 -e "CREATE DATABASE IF NOT EXISTS go_demo"
```

### Q3: Redis 连接失败？

**A**: 检查 Redis 服务状态

```bash
# 检查 Redis 状态
docker-compose -f deployments/docker-compose.simple.yml ps redis

# 测试 Redis 连接
docker-compose -f deployments/docker-compose.simple.yml exec redis redis-cli ping
```

### Q4: Docker 构建失败？

**A**: 常见原因和解决方案：

```bash
# 1. 清理 Docker 缓存
docker system prune -a

# 2. 重新构建（不使用缓存）
docker-compose -f deployments/docker-compose.simple.yml build --no-cache

# 3. 检查 Dockerfile 中的 Go 版本是否匹配
```

### Q5: 如何查看应用日志？

**A**: 多种方式查看日志

```bash
# Docker 日志
docker-compose -f deployments/docker-compose.simple.yml logs -f app

# 本地日志文件
tail -f logs/app.log
tail -f logs/request.log

# 进入容器查看
docker-compose -f deployments/docker-compose.simple.yml exec app sh
cat /app/logs/app.log
```

### Q6: 如何重置数据库？

**A**: 删除并重新创建

```bash
# 停止服务
docker-compose -f deployments/docker-compose.simple.yml down

# 删除数据卷
docker-compose -f deployments/docker-compose.simple.yml down -v

# 重新启动
docker-compose -f deployments/docker-compose.simple.yml up -d
```

### Q7: 生产环境如何配置 HTTPS？

**A**: 使用 Nginx 配置 SSL

```bash
# 1. 准备 SSL 证书
# 将证书放到 deployments/nginx/ssl/ 目录

# 2. 修改 Nginx 配置
# 编辑 deployments/nginx/conf.d/default.conf

# 3. 重启 Nginx
docker-compose restart nginx
```

### Q8: 如何备份数据？

**A**: 备份数据库和 Redis

```bash
# 备份 MySQL
docker-compose -f deployments/docker-compose.simple.yml exec mysql \
  mysqldump -uroot -p123456 go_demo > backup_$(date +%Y%m%d).sql

# 备份 Redis
docker-compose -f deployments/docker-compose.simple.yml exec redis \
  redis-cli SAVE
docker cp go-demo-redis:/data/dump.rdb ./redis_backup_$(date +%Y%m%d).rdb
```

## 📊 监控和维护

### 查看资源使用
```bash
# 查看容器资源使用
docker stats

# 查看特定容器
docker stats go-demo-app
```

### 清理资源
```bash
# 清理未使用的镜像
docker image prune

# 清理未使用的容器
docker container prune

# 清理未使用的数据卷
docker volume prune

# 清理所有未使用的资源
docker system prune -a
```

### 更新应用
```bash
# 1. 拉取最新代码
git pull

# 2. 重新构建
docker-compose -f deployments/docker-compose.simple.yml build

# 3. 重启服务
docker-compose -f deployments/docker-compose.simple.yml up -d
```

## 🎯 推荐工作流

### 开发环境
```bash
# 1. 启动依赖服务
docker-compose -f deployments/docker-compose.simple.yml up -d mysql redis

# 2. 使用 Air 热重载开发
air

# 3. 访问 API 文档进行测试
open http://localhost:8080/swagger/index.html
```

### 生产环境
```bash
# 1. 设置环境变量
export DATABASE_DSN="..."
export JWT_secret_KEY="..."

# 2. 构建应用
go build -o go_demo main.go

# 3. 启动应用
./go_demo server --config=./configs/config.prod.yaml

# 4. 配置进程管理（systemd/supervisor）
# 5. 配置 Nginx 反向代理
# 6. 配置监控和日志收集
```

## 📚 相关资源

- **API 文档**: http://localhost:8080/swagger/index.html
- **健康检查**: http://localhost:8080/health
- **项目仓库**: [GitHub Repository]
- **Docker Hub**: [Docker Hub Repository]

## 💡 最佳实践

1. **开发环境**：使用 [`config.dev.yaml`](../configs/config.dev.yaml) + Air 热重载
2. **生产环境**：使用 [`config.prod.yaml`](../configs/config.prod.yaml) + 环境变量
3. **定期备份**：每天备份数据库和 Redis
4. **监控日志**：使用日志收集工具（ELK/Loki）
5. **性能监控**：使用 Prometheus + Grafana
6. **安全加固**：定期更新依赖，使用强密码

---

**最后更新**: 2025-12-26  
**文档版本**: v1.0.0
