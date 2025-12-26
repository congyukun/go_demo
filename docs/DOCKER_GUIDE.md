# Docker 启动指南

## 🚀 快速启动

### 方式一：简化版（推荐开发使用）

**不使用Nginx，直接访问Go应用**

```bash
# 1. 进入部署目录
cd deployments

# 2. 启动服务
docker-compose -f docker-compose.simple.yml up -d

# 3. 查看日志
docker-compose -f docker-compose.simple.yml logs -f app

# 4. 测试服务
curl http://localhost:8080/health

# 5. 停止服务
docker-compose -f docker-compose.simple.yml down
```

**访问地址**：
- Go应用：http://localhost:8080
- MySQL：localhost:3306
- Redis：localhost:6379

### 方式二：完整版（推荐生产使用）

**使用Nginx作为反向代理**

```bash
# 1. 进入部署目录
cd deployments

# 2. 启动服务
docker-compose up -d

# 3. 查看日志
docker-compose logs -f

# 4. 测试服务
curl http://localhost/health

# 5. 停止服务
docker-compose down
```

**访问地址**：
- Nginx（HTTP）：http://localhost
- Nginx（HTTPS）：https://localhost
- Go应用（内部）：http://localhost:8080
- MySQL：localhost:3306
- Redis：localhost:6379

## 📋 启动前准备

### 1. 确保配置文件存在

```bash
# 检查配置文件
ls -la configs/config.yaml

# 如果不存在，从模板复制
cp configs/config.example.yaml configs/config.yaml
```

### 2. 确保Docker和Docker Compose已安装

```bash
# 检查Docker版本
docker --version
# 应该显示：Docker version 20.x.x 或更高

# 检查Docker Compose版本
docker-compose --version
# 应该显示：Docker Compose version 2.x.x 或更高
```

### 3. 确保端口未被占用

```bash
# 检查端口占用（macOS/Linux）
lsof -i :8080  # Go应用
lsof -i :3306  # MySQL
lsof -i :6379  # Redis
lsof -i :80    # Nginx HTTP
lsof -i :443   # Nginx HTTPS

# 如果有进程占用，可以停止或修改docker-compose中的端口映射
```

## 🔧 常用命令

### 启动和停止

```bash
# 启动所有服务（后台运行）
docker-compose up -d

# 启动所有服务（前台运行，查看日志）
docker-compose up

# 停止所有服务
docker-compose stop

# 停止并删除容器
docker-compose down

# 停止并删除容器、网络、数据卷
docker-compose down -v
```

### 查看状态

```bash
# 查看所有服务状态
docker-compose ps

# 查看服务日志
docker-compose logs

# 实时查看日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f app
docker-compose logs -f mysql
docker-compose logs -f redis
```

### 重启服务

```bash
# 重启所有服务
docker-compose restart

# 重启特定服务
docker-compose restart app
docker-compose restart mysql
```

### 进入容器

```bash
# 进入Go应用容器
docker-compose exec app sh

# 进入MySQL容器
docker-compose exec mysql bash

# 进入Redis容器
docker-compose exec redis sh
```

### 重新构建

```bash
# 重新构建并启动
docker-compose up -d --build

# 仅重新构建
docker-compose build

# 强制重新构建（不使用缓存）
docker-compose build --no-cache
```

## 🧪 测试服务

### 1. 健康检查

```bash
# 基础健康检查
curl http://localhost:8080/health

# 详细健康检查
curl http://localhost:8080/health/check

# 就绪检查
curl http://localhost:8080/health/ready

# 存活检查
curl http://localhost:8080/health/live
```

### 2. 测试API

```bash
# 用户注册
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Test123456",
    "email": "test@example.com",
    "name": "测试用户"
  }'

# 用户登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Test123456"
  }'
```

### 3. 测试数据库连接

```bash
# 进入MySQL容器
docker-compose exec mysql mysql -uroot -p123456

# 在MySQL中执行
USE go_demo;
SHOW TABLES;
SELECT * FROM users;
```

### 4. 测试Redis连接

```bash
# 进入Redis容器
docker-compose exec redis redis-cli

# 在Redis中执行
PING
KEYS *
```

## 🐛 故障排查

### 问题1：容器启动失败

```bash
# 查看详细日志
docker-compose logs app

# 常见原因：
# 1. 端口被占用 - 修改docker-compose.yml中的端口映射
# 2. 配置文件错误 - 检查configs/config.yaml
# 3. 数据库未就绪 - 等待MySQL完全启动
```

### 问题2：无法连接数据库

```bash
# 检查MySQL是否启动
docker-compose ps mysql

# 查看MySQL日志
docker-compose logs mysql

# 检查网络连接
docker-compose exec app ping mysql

# 解决方案：
# 1. 确保config.yaml中数据库地址为 "mysql:3306"
# 2. 等待MySQL健康检查通过
# 3. 检查数据库密码是否正确
```

### 问题3：无法连接Redis

```bash
# 检查Redis是否启动
docker-compose ps redis

# 查看Redis日志
docker-compose logs redis

# 测试连接
docker-compose exec app ping redis

# 解决方案：
# 1. 确保config.yaml中Redis地址为 "redis:6379"
# 2. 检查Redis是否正常运行
```

### 问题4：Go应用编译失败

```bash
# 查看构建日志
docker-compose build app

# 常见原因：
# 1. Go版本不匹配 - 检查Dockerfile中的Go版本
# 2. 依赖下载失败 - 检查网络或使用代理
# 3. 代码语法错误 - 本地先运行 go build 测试
```

## 📊 监控和维护

### 查看资源使用

```bash
# 查看容器资源使用情况
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

### 备份数据

```bash
# 备份MySQL数据
docker-compose exec mysql mysqldump -uroot -p123456 go_demo > backup.sql

# 备份Redis数据
docker-compose exec redis redis-cli SAVE
docker cp go-demo-redis:/data/dump.rdb ./redis-backup.rdb
```

## 🎯 推荐工作流

### 开发环境

```bash
# 1. 使用简化版启动
cd deployments
docker-compose -f docker-compose.simple.yml up -d

# 2. 查看日志确认启动成功
docker-compose -f docker-compose.simple.yml logs -f app

# 3. 开发和测试
# ... 你的开发工作 ...

# 4. 修改代码后重新构建
docker-compose -f docker-compose.simple.yml up -d --build

# 5. 完成后停止
docker-compose -f docker-compose.simple.yml down
```

### 生产环境

```bash
# 1. 使用完整版启动
cd deployments
docker-compose up -d

# 2. 检查所有服务状态
docker-compose ps

# 3. 查看日志
docker-compose logs -f

# 4. 监控运行状态
docker stats

# 5. 定期备份数据
# ... 执行备份脚本 ...
```

## 📚 相关文档

- [Docker官方文档](https://docs.docker.com/)
- [Docker Compose文档](https://docs.docker.com/compose/)
- [项目改进文档](../docs/IMPROVEMENTS.md)
- [Nginx使用指南](../docs/NGINX_GUIDE.md)

## 💡 提示

1. **开发时使用简化版**：`docker-compose.simple.yml`
2. **生产时使用完整版**：`docker-compose.yml`
3. **修改代码后记得重新构建**：`--build`
4. **定期清理Docker资源**：避免磁盘空间不足
5. **查看日志排查问题**：`docker-compose logs -f`

---

**快速启动命令**：
```bash
cd deployments && docker-compose -f docker-compose.simple.yml up -d
```
