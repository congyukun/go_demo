# 🚀 Go Demo 项目 - 快速启动指南

## 一键启动（推荐）

使用交互式启动脚本，最简单的方式：

```bash
./quick-start.sh
```

脚本会自动：
- ✅ 检查 Docker 环境
- ✅ 检查端口占用
- ✅ 提供多种部署选项
- ✅ 自动健康检查
- ✅ 显示访问地址

## 手动启动

### 方式一：完整部署（生产环境）

包含：应用 + MySQL + Redis + Nginx

```bash
cd deployments
docker-compose up -d
```

**访问地址**：
- 应用: http://localhost:8080
- Nginx: http://localhost
- Swagger: http://localhost:8080/swagger/index.html

### 方式二：简化部署（开发环境）

包含：应用 + MySQL + Redis

```bash
cd deployments
docker-compose -f docker-compose.simple.yml up -d
```

### 方式三：本地开发

只启动数据库，应用在本地运行：

```bash
# 1. 启动数据库
cd deployments
docker-compose up -d mysql redis

# 2. 运行应用
cd ..
go run main.go server --config=./configs/config.dev.yaml
```

## 验证部署

```bash
# 健康检查
curl http://localhost:8080/health

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

## 常用命令

```bash
# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f app

# 停止服务
docker-compose down

# 重启服务
docker-compose restart app

# 进入容器
docker-compose exec app sh
```

## 服务端口

| 服务 | 端口 | 说明 |
|------|------|------|
| 应用 | 8080 | Go 应用主服务 |
| MySQL | 3306 | 数据库 |
| Redis | 6379 | 缓存 |
| Nginx | 80/443 | 反向代理 |

## 默认账号

**MySQL**:
- 用户: root
- 密码: 123456
- 数据库: go_demo

**Redis**:
- 无密码

## 故障排查

### 端口被占用

```bash
# 查看端口占用
lsof -i :8080

# 修改端口（编辑 docker-compose.yml）
ports:
  - "9090:8080"
```

### 服务启动失败

```bash
# 查看日志
docker-compose logs app

# 重新构建
docker-compose build --no-cache
docker-compose up -d
```

### 数据库连接失败

```bash
# 检查 MySQL 状态
docker-compose ps mysql
docker-compose logs mysql

# 等待 MySQL 完全启动（约30秒）
docker-compose logs -f mysql
# 看到 "ready for connections" 后重启应用
docker-compose restart app
```

## 更多文档

- 📖 [完整部署指南](DOCKER_DEPLOYMENT_GUIDE.md) - 详细的 Docker 部署文档
- 📖 [项目文档](README.md) - 项目概述和功能说明
- 📖 [API 文档](API.md) - API 接口说明
- 📖 [部署优化](docs/DEPLOYMENT_OPTIMI 性能优化建议

## 获取帮助

遇到问题？

1. 查看日志：`docker-compose logs -f app`
2. 检查状态：`docker-compose ps`
3. 查看文档：[DOCKER_DEPLOYMENT_GUIDE.md](DOCKER_DEPLOYMENT_GUIDE.md)
4. 提交 Issue 或联系维护者

---

**快速启动脚本功能**：
- ✅ 完整部署
- ✅ 简化部署
- ✅ 仅启动依赖
- ✅ 停止服务
- ✅ 查看状态
- ✅ 查看日志
- ✅ 重启服务
- ✅ 清理数据

**推荐使用**: `./quick-start.sh` 获得最佳体验！
