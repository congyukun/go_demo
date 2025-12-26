# Go Demo 项目 Docker 部署完整指南

## 📋 项目概述

这是一个标准的 Go Web 应用项目，采用分层架构设计，包含用户管理和认证功能。

**技术栈**：
- **后端框架**: Gin Web Framework
- **数据库**: MySQL 8.0
- **缓存**: Redis 7
- **ORM**: GORM
- **认证**: JWT
- **日志**: Zap
- **依赖注入**: Google Wire
- **反向代理**: Nginx

## 🚀 快速开始（一键部署）

### 前置要求

- Docker 20.x 或更高版本
- Docker Compose 2.x 或更高版本

```bash
# 检查 Docker 版本
docker --version
docker-compose --version
```

### 一键启动所有服务

```bash
# 1. 克隆项目（如果还没有）
git clone <repository-url>
cd go_demo

# 2. 进入部署目录
cd deployments

# 3. 启动所有服务（包括 MySQL、Redis、应用、Nginx）
docker-compose up -d

# 4. 查看服务状态
docker-compose ps

# 5. 查看应用日志
docker-compose logs -f app
```

**访问地址**：
- 应用 API: http://localhost:8080
- Nginx 代理: http://localhost
- Swagger 文档: http://localhost:8080/swagger/index.html
- 健康检查: http://localhost:8080/health

### 验证部署

```bash
# 健康检查
curl http://localhost:8080/health

# 测试用户注册
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Test123456",
    "email": "test@example.com",
    "name": "测试用户"
  }'

# 测试用户登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Test123456"
  }'
```

## 📦 部署架构

### 服务组件

```
┌─────────────────────────────────────────────────────┐
│                    Nginx (80/443)                    │
│              反向代理 + 负载均衡 + SSL                │
└──────────────────────┬──────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────┐
│                Go Application (8080)                 │
│         Gin + GORM + JWT + Wire + Zap               │
└──────────────┬──────────────────┬───────────────────┘
               │                  │
       ┌───────▼────────┐  ┌─────▼──────┐
       │  MySQL (3306)  │  │ Redis (6379)│
       │   数据持久化    │  │  缓存+限流   │
       └────────────────┘  └─────────────┘
```

### Docker Compose 服务说明

| 服务名 | 容器名 | 端口映射 | 说明 |
|--------|--------|----------|------|
| app | go-demo-app | 8080:8080 | Go 应用主服务 |
| mysql | go-demo-mysql | 3306:3306 | MySQL 数据库 |
| redis | go-demo-redis | 6379:6379 | Redis 缓存 |
| nginx | go-demo-nginx | 80:80, 443:443 | Nginx 反向代理 |

## 🔧 详细部署步骤

### 方式一：完整部署（推荐生产环境）

包含所有服务：应用 + MySQL + Redis + Nginx

```bash
cd deployments

# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看所有日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f app
docker-compose logs -f mysql
docker-compose logs -f redis
docker-compose logs -f nginx
```

### 方式二：简化部署（开发/测试环境）

只启动核心服务：应用 + MySQL + Redis

```bash
cd deployments

# 使用简化配置启动
docker-compose -f docker-compose.simple.yml up -d

# 查看服务状态
docker-compose -f docker-compose.simple.yml ps

# 查看日志
docker-compose -f docker-compose.simple.yml logs -f app
```

### 方式三：仅启动依赖服务

只启动 MySQL 和 Redis，应用在本地运行（适合开发调试）

```bash
cd deployments

# 只启动数据库服务
docker-compose up -d mysql redis

# 在项目根目录运行应用
cd ..
go run main.go server --config=./configs/config.dev.yaml
```

## 📝 配置说明

### 环境配置文件

项目提供多个环境配置：

| 配置文件 | 用途 | 数据库地址 | 日志级别 |
|---------|------|-----------|---------|
| `config.docker.yaml` | Docker 环境 | mysql:3306 | debug |
| `config.dev.yaml` | 本地开发 | localhost:3306 | debug |
| `config.yaml` | 生产环境 | 环境变量 | warn |

### Docker Compose 环境变量

在 [`docker-compose.yml`](deployments/docker-compose.yml) 中配置：

```yaml
environment:
  # 数据库配置
  - MYSQL_HOST=mysql
  - MYSQL_PORT=3306
  - MYSQL_USER=root
  - MYSQL_PASSWORD=123456
  - MYSQL_DATABASE=go_demo
  
  # Redis 配置
  - REDIS_HOST=redis
  - REDIS_PORT=6379
  
  # 等待服务配置
  - MAX_RETRIES=30
  - RETRY_INTERVAL=2
```

### 修改默认配置

#### 修改数据库密码

编辑 [`docker-compose.yml`](deployments/docker-compose.yml)：

```yaml
mysql:
  environment:
    MYSQL_ROOT_PASSWORD: your_new_password  # 修改这里
```

同时修改 [`configs/config.docker.yaml`](configs/config.docker.yaml)：

```yaml
database:
  dsn: "root:your_new_password@tcp(mysql:3306)/go_demo?charset=utf8mb4&parseTime=True&loc=Local"
```

#### 修改应用端口

编辑 [`docker-compose.yml`](deployments/docker-compose.yml)：

```yaml
app:
  ports:
    - "9090:8080"  # 将主机端口改为 9090
```

#### 修改 JWT 密钥

编辑 [`configs/config.docker.yaml`](configs/config.docker.yaml)：

```yaml
jwt:
  secret_key: "your-secret-jwt-key-at-least-32-characters"
```

生成安全密钥：
```bash
openssl rand -base64 32
```

## 🐳 Docker 镜像构建

### 查看 Dockerfile

项目使用多阶段构建优化镜像大小：

```dockerfile
# 构建阶段 - 使用 golang:1.24-alpine
FROM golang:1.24-alpine AS builder
WORKDIR /app
# ... 编译应用

# 运行阶段 - 使用 alpine:3.19
FROM alpine:3.19
# ... 复制二进制文件和配置
```

### 手动构建镜像

```bash
# 在项目根目录执行
docker build -f deployments/Dockerfile -t go-demo:latest .

# 查看镜像
docker images | grep go-demo

# 运行容器
docker run -d \
  --name go-demo-app \
  -p 8080:8080 \
  -e MYSQL_HOST=host.docker.internal \
  -e REDIS_HOST=host.docker.internal \
  go-demo:latest
```

### 构建优化

项目已配置 [`.dockerignore`](.dockerignore) 排除不必要的文件：

```
# 排除的内容
.git/
*.md
tests/
logs/
.vscode/
.idea/
```

## 🔍 服务管理

### 启动服务

```bash
# 启动所有服务
docker-compose up -d

# 启动特定服务
docker-compose up -d app
docker-compose up -d mysql redis

# 前台运行（查看实时日志）
docker-compose up
```

### 停止服务

```bash
# 停止所有服务
docker-compose down

# 停止并删除数据卷（⚠️ 会删除数据库数据）
docker-compose down -v

# 停止特定服务
docker-compose stop app
```

### 重启服务

```bash
# 重启所有服务
docker-compose restart

# 重启特定服务
docker-compose restart app
docker-compose restart mysql
```

### 查看服务状态

```bash
# 查看所有服务状态
docker-compose ps

# 查看详细信息
docker-compose ps -a

# 查看资源使用
docker stats
```

### 查看日志

```bash
# 查看所有服务日志
docker-compose logs

# 实时跟踪日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f app

# 查看最近 100 行日志
docker-compose logs --tail=100 app

# 查看带时间戳的日志
docker-compose logs -f -t app
```

### 进入容器

```bash
# 进入应用容器
docker-compose exec app sh

# 进入 MySQL 容器
docker-compose exec mysql bash

# 进入 Redis 容器
docker-compose exec redis sh

# 以 root 用户进入
docker-compose exec -u root app sh
```

## 🧪 测试和验证

### 健康检查

```bash
# 应用健康检查
curl http://localhost:8080/health

# 预期响应
{
  "status": "ok",
  "timestamp": "2025-12-26T09:45:00}
```

### 测试数据库连接

```bash
# 进入 MySQL 容器
docker-compose exec mysql mysql -uroot -p123456

# 在 MySQL 中执行
USE go_demo;
SHOW TABLES;
SELECT * FROM users;
EXIT;
```

### 测试 Redis 连接

```bash
# 进入 Redis 容器
docker-compose exec redis redis-cli

# 在 Redis 中执行
PING
KEYS *
INFO
EXIT
```

### API 测试脚本

创建测试脚本 `test-api.sh`：

```bash
#!/bin/bash

BASE_URL="http://localhost:8080"

echo "=== 1. 健康检查 ==="
curl -s $BASE_URL/health | jq

echo -e "\n=== 2. 用户注册 ==="
curl -s -X POST $BASE_URL/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Test123456",
    "email": "test@example.com",
    "name": "测试用户"
  }' | jq

echo -e "\n=== 3. 用户登录 ==="
TOKEN=$(curl -s -X POST $BASE_URL/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Test123456"
  }' | jq -r '.data.access_token')

echo "Token: $TOKEN"

echo -e "\n=== 4. 获取用户列表 ==="
curl -s -X GET "$BASE_URL/api/v1/users?page=1&size=10" \
  -H "Authorization: Bearer $TOKEN" | jq
```

运行测试：
```bash
chmod +x test-api.sh
./test-api.sh
```

## 🔧 故障排查

### 常见问题

#### 1. 端口被占用

**错误信息**：
```
Error starting userland proxy: listen tcp4 0.0.0.0:8080: bind: address already in use
```

**解决方案**：
```bash
# 查看端口占用
lsof -i :8080
netstat -tuln | grep 8080

# 修改端口映射
# 编辑 docker-compose.yml
ports:
  - "9090:8080"  # 改为其他端口
```

#### 2. 数据库连接失败

**错误信息**：
```
Error 2003: Can't connect to MySQL server on 'mysql'
```

**解决方案**：
```bash
# 1. 检查 MySQL 容器状态
docker-compose ps mysql

# 2. 查看 MySQL 日志
docker-compose logs mysql

# 3. 检查健康状态
docker-compose exec mysql mysqladmin ping -h localhost -u root -p123456

# 4. 等待 MySQL 完全启动（通常需要 30 秒）
docker-compose logs -f mysql
# 看到 "ready for connections" 表示启动完成

# 5. 重启应用服务
docker-compose restart app
```

#### 3. Redis 连接失败

**解决方案**：
```bash
# 检查 Redis 状态
docker-compose ps redis

# 测试 Redis 连接
docker-compose exec redis redis-cli ping

# 查看 Redis 日志
docker-compose logs redis
```

#### 4. 应用启动失败

**解决方案**：
```bash
# 查看应用日志
docker-compose logs app

# 检查配置文件
docker-compose exec app cat /app/configs/config.docker.yaml

# 检查环境变量
docker-compose exec app env | grep -E "MYSQL|REDIS"

# 重新构建镜像
docker-compose build --no-cache app
docker-compose up -d app
```

#### 5. 数据卷权限问题

**错误信息**：
```
Permission denied: '/app/logs/app.log'
```

**解决方案**：
```bash
# 检查日志目录权限
ls -la logs/

# 修复权限
chmod -R 755 logs/
chown -R 1001:1001 logs/

# 或删除数据卷重新创建
docker-compose down -v
docker-compose up -d
```

### 调试技巧

#### 查看容器详细信息

```bash
# 查看容器配置
docker inspect go-demo-app

# 查看容器网络
docker network inspect deployments_go-demo-network

# 查看数据卷
docker volume ls
docker volume inspect deployments_mysql_data
```

#### 实时监控

```bash
# 监控资源使用
docker stats

# 监控特定容器
docker stats go-demo-app go-demo-mysql go-demo-redis

# 查看容器进程
docker-compose top
```

#### 网络诊断

```bash
# 进入应用容器测试网络
docker-compose exec app sh

# 测试 MySQL 连接
nc -zv mysql 3306
ping mysql

# 测试 Redis 连接
nc -zv redis 6379
ping redis

# 测试 DNS 解析
nslookup mysql
nslookup redis
```

## 📊 数据管理

### 数据备份

#### 备份 MySQL 数据

```bash
# 备份整个数据库
docker-compose exec mysql mysqldump -uroot -p123456 go_demo > backup_$(date +%Y%m%d_%H%M%S).sql

# 备份特定表
docker-compose exec mysql mysqldump -uroot -p123456 go_demo users > users_backup.sql

# 备份所有数据库
docker-compose exec mysql mysqldump -uroot -p123456 --all-databases > all_backup.sql
```

#### 备份 Redis 数据

```bash
# 触发 Redis 保存
docker-compose exec redis redis-cli SAVE

# 复制 RDB 文件
docker cp go-demo-redis:/data/dump.rdb ./redis_backup_$(date +%Y%m%d_%H%M%S).rdb
```

#### 备份数据卷

```bash
# 备份 MySQL 数据卷
docker run --rm \
  -v deployments_mysql_data:/data \
  -v $(pwd):/backup \
  alpine tar czf /backup/mysql_data_backup.tar.gz -C /data .

# 备份 Redis 数据卷
docker run --rm \
  -v deployments_redis_data:/data \
  -v $(pwd):/backup \
  alpine tar czf /backup/redis_data_backup.tar.gz -C /data .
```

### 数据恢复

#### 恢复 MySQL 数据

```bash
# 从 SQL 文件恢复
docker-compose exec -T mysql mysql -uroot -p123456 go_demo < backup.sql

# 或进入容器恢复
docker cp backup.sql go-demo-mysql:/tmp/
docker-compose exec mysql mysql -uroot -p123456 go_demo -e "source /tmp/backup.sql"
```

#### 恢复 Redis 数据

```bash
# 停止 Redis
docker-compose stop redis

# 复制 RDB 文件
docker cp redis_backup.rdb go-demo-redis:/data/dump.rdb

# 启动 Redis
docker-compose start redis
```

#### 恢复数据卷

```bash
# 恢复 MySQL 数据卷
docker run --rm \
  -v deployments_mysql_data:/data \
  -v $(pwd):/backup \
  alpine sh -c "cd /data && tar xzf /backup/mysql_data_backup.tar.gz"
```

### 数据库初始化

项目包含初始化 SQL 脚本：

```bash
# 查看初始化脚本
cat deployments/mysql/init.sql

# 手动执行初始化
docker-compose exec mysql mysql -uroot -p123456 go_demo < deployments/mysql/init.sql
```

## 🚀 生产环境部署

### 生产环境检查清单

- [ ] 修改默认密码（MySQL、Redis、JWT）
- [ ] 配置 HTTPS/SSL 证书
- [ ] 设置合适的资源限制
- [ ] 配置日志轮转
- [ ] 设置数据备份策略
- [ ] 配置监控告警
- [ ] 启用防火墙规则
- [ ] 配置域名和 DNS

### 安全加固

#### 1. 修改默认密码

```yaml
# docker-compose.yml
mysql:
  environment:
    MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}  # 使用环境变量

redis:
  command:
    - redis-server
    - --requirepass ${REDIS_PASSWORD}  # 设置密码
```

创建 `.env` 文件：
```bash
MYSQL_ROOT_PASSWORD=your_strong_password_here
REDIS_PASSWORD=your_redis_password_here
JWT_secret_KEY=your_jwt_secret_key_at_least_32_chars
```

#### 2. 配置 HTTPS

编辑 [`deployments/nginx/conf.d/default.conf`](deployments/nginx/conf.d/default.conf)：

```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    
    # SSL 配置
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    
    location / {
        proxy_pass http://app:8080;
        # ... 其他配置
    }
}

# HTTP 重定向到 HTTPS
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}
```

挂载 SSL 证书：
```yaml
nginx:
  volumes:
    - ./nginx/ssl:/etc/nginx/ssl:ro
```

#### 3. 资源限制

```yaml
# docker-compose.yml
app:
  deploy:
    resources:
      limits:
        cpus: '2'
        memory: 1G
      reservations:
        cpus: '0.5'
        memory: 512M

mysql:
  deploy:
    resources:
      limits:
        cpus: '2'
        memory: 2G
      reservations:
        cpus: '1'
        memory: 1G
```

### 性能优化

#### MySQL 优化

```yaml
mysql:
  command:
    - --max_connections=500
    - --innodb_buffer_pool_size=1G
    - --innodb_log_file_size=256M
    - --query_cache_size=64M
```

#### Redis 优化

```yaml
redis:
  command:
    - redis-server
    - --maxmemory 512mb
    - --maxmemory-policy allkeys-lru
    - --save 900 1
    - --save 300 10
```

### 监控配置

#### 添加 Prometheus 监控

创建 `docker-compose.monitoring.yml`：

```yaml
services:
  prometheus:
    image: prom/prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    networks:
      - go-demo-network

  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"
    volumes:
      - grafana_data:/var/lib/grafana
    networks:
      - go-demo-network

volumes:
  prometheus_data:
  grafana_data:
```

启动监控：
```bash
docker-compose -f docker-compose.yml -f docker-compose.monitoring.yml up -d
```

## 📚 常用命令速查

### Docker Compose 命令

```bash
# 启动服务
docker-compose up -d                    # 后台启动
docker-compose up                       # 前台启动
docker-compose up -d --build            # 重新构建并启动

# 停止服务
docker-compose stop                     # 停止服务
docker-compose down                     # 停止并删除容器
docker-compose down -v                  # 停止并删除容器和数据卷

# 查看状态
docker-compose ps                       # 查看服务状态
docker-compose logs -f                  # 查看日志
docker-compose top                      # 查看进程

# 重启服务
docker-compose restart                  # 重启所有服务
docker-compose restart app              # 重启特定服务

# 执行命令
docker-compose exec app sh              # 进入容器
docker-compose exec mysql mysql -uroot -p123456  # 执行命令

# 构建镜像
docker-compose build                    # 构建所有镜像
docker-compose build --no-cache app     # 不使用缓存构建
```

### Docker 命令

```bash
# 镜像管理
docker images                           # 查看镜像
docker rmi <image-id>                   # 删除镜像
docker image prune                      # 清理未使用镜像

# 容器管理
docker ps                               # 查看运行中容器
docker ps -a                            # 查看所有容器
docker rm <container-id>                # 删除容器
docker container prune                  # 清理停止的容器

# 日志查看
docker logs -f <container-name>         # 查看容器日志
docker logs --tail 100 <container-name> # 查看最近100行

# 资源管理
docker stats                            # 查看资源使用
docker system df                        # 查看磁盘使用
docker system prune -a                  # 清理所有未使用资源
```

## 🔗 相关文档

- [项目 README](README.md) - 项目概述和快速开始
- [API 文档](API.md) - API 接口说明
- [部署文档](docs/DEPLOYMENT.md) - 详细部署指南
- [架构文档](docs/ARCHITECTURE.md) - 系统架构说明
- [Docker 指南](docs/DOCKER_GUIDE.md) - Docker 使用指南

## 📞 获取帮助

如遇到问题：

1. 查看日志：`docker-compose logs -f app`
2. 检查服务状态：`docker-compose ps`
3. 查看本文档的故障排查章节
4. 提交 Issue 或联系维护者

---

**最后更新**: 2025-12-26
**维护者**: cunliwakun@163.com
