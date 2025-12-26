# Makefile 使用指南

## 概述

已将 [`quick-start.sh`](quick-start.sh) 的功能整合到 [`Makefile`](Makefile) 中，提供更标准化的项目管理方式。

## 为什么选择 Makefile？

1. **更标准化** - Makefile 是 Go 项目的标准实践
2. **更简洁** - 命令更短，易于记忆（如 `make docker-deploy`）
3. **更灵活** - 可以轻松组合和扩展命令
4. **IDE 集成** - 大多数 IDE 都支持 Makefile 命令提示

## 快速开始

### 查看所有可用命令
```bash
make help
```

### 常用部署命令

#### 1. 完整部署（推荐用于生产环境）
```bash
make docker-deploy
```
包含：应用 + MySQL + Redis + Nginx

#### 2. 简化部署（推荐用于开发环境）
```bash
make docker-deploy-simple
```
包含：应用 + MySQL + Redis

#### 3. 仅启动依赖服务（本地开发）
```bash
make docker-deps
```
然后在本地运行应用：
```bash
go run main.go server --config=./configs/config.dev.yaml
```

### 常用管理命令

```bash
# 查看服务状态
make docker-status

# 查看应用日志
make docker-logs

# 查看所有服务日志
make docker-logs-all

# 重启应用
make docker-restart

# 停止所有服务
make docker-stop

# 健康检查
make health

# 显示服务信息
make docker-info
```

### 开发命令

```bash
# 安装依赖
make deps

# 格式化代码
make fmt

# 运行测试
make test

# 代码质量检查
make lint

# 开发模式运行（热重载）
make dev

# 生成 API 文档
make docs
```

## 命令对比

| quick-start.sh 功能 | Makefile 命令 | 说明 |
|-------------------|--------------|------|
| 选项 1: 完整部署 | `make docker-deploy` | 一键完整部署 |
| 选项 2: 简化部署 | `make docker-deploy-simple` | 简化部署 |
| 选项 3: 仅依赖服务 | `make docker-deps` | 仅启动 MySQL + Redis |
| 选项 4: 停止服务 | `make docker-stop` | 停止所有服务 |
| 选项 5: 查看状态 | `make docker-status` | 查看服务状态 |
| 选项 6: 查看日志 | `make docker-logs` | 查看日志 |
| 选项 7: 重启服务 | `make docker-restart` | 重启服务 |
| 选项 8: 清理数据 | `make docker-clean` | 清理所有数据 |

## 优势

### Makefile 的优势
- ✅ 命令简短：`make docker-deploy` vs `./quick-start.sh` 然后选择选项
- ✅ 可组合：可以在 CI/CD 中轻松使用
- ✅ 自动补全：支持 shell 自动补全
- ✅ 并行执行：可以使用 `-j` 参数并行执行
- ✅ 依赖管理：可以定义命令之间的依赖关系

### quick-start.sh 的优势
- ✅ 交互式菜单：更友好的用户界面
- ✅ 颜色输出：更好的视觉效果
- ✅ 详细提示：更多的帮助信息

## 建议

1. **日常开发**：使用 Makefile 命令（更快捷）
2. **新手入门**：可以保留 quick-start.sh（更友好）
3. **CI/CD**：使用 Makefile 命令（更标准）

## 保留 quick-start.sh 的场景

如果你希望保留 quick-start.sh，可以用于：
- 团队新成员的快速上手
- 不熟悉 Makefile 的用户
- 需要更详细的交互式引导

## 删除 quick-start.sh

如果决定只使用 Makefile，可以删除：
```bash
rm quick-start.sh
```

然后更新相关文档中的引用。

## 总结

**推荐方案**：保留 [`Makefile`](Makefile) 作为主要工具，可选保留 [`quick-start.sh`](quick-start.sh) 作为新手友好的备选方案。

所有 quick-start.sh 的核心功能都已集成到 Makefile 中，并且提供了更多的灵活性和扩展性。
