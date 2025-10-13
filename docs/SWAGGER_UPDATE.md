# Swagger 文档更新指南

## 📋 当前状态

项目已经包含完整的Swagger文档，包括：

### ✅ 已生成的文档
- `docs/docs.go` - 自动生成的Swagger文档（1097行）
- `api/openapi.yaml` - OpenAPI 3.0规范文档（469行）

### ✅ 已覆盖的API端点

#### 认证接口
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/logout` - 用户登出
- `GET /api/v1/auth/profile` - 获取当前用户信息
- `POST /api/v1/auth/refresh` - 刷新访问令牌

#### 用户管理接口
- `GET /api/v1/users` - 获取用户列表（分页）
- `POST /api/v1/users` - 创建新用户
- `GET /api/v1/users/{id}` - 获取用户详情
- `PUT /api/v1/users/{id}` - 更新用户信息
- `DELETE /api/v1/users/{id}` - 删除用户
- `PUT /api/v1/users/profile` - 更新个人资料
- `PUT /api/v1/users/password` - 修改密码
- `GET /api/v1/users/stats` - 获取用户统计

#### 系统接口
- `GET /health` - 健康检查

### ✅ 完整的Schema定义
- 用户模型（User）
- 认证请求/响应模型
- 分页响应模型
- 错误响应模型
- JWT认证方案

## 🔄 如何更新Swagger文档

### 方法1：使用swag工具（推荐）
```bash
# 安装swag工具
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文档
swag init -g cmd/server/main.go

# 验证文档
swag fmt
```

### 方法2：手动更新OpenAPI规范
```bash
# 编辑api/openapi.yaml文件
# 确保与代码注释保持一致
```

## 📊 文档验证

### 验证Swagger UI
1. 启动服务：`go run cmd/server/main.go`
2. 访问：`http://localhost:8080/swagger/index.html`
3. 测试所有端点功能

### 验证OpenAPI规范
```bash
# 使用在线验证器
# https://editor.swagger.io/
# 上传api/openapi.yaml文件
```

## 🎯 文档特点

### 1. 完整的API覆盖
- 所有RESTful端点都有详细文档
- 包含请求/响应示例
- 参数验证规则清晰

### 2. 认证集成
- JWT Bearer Token认证
- 安全方案定义完整
- 权限控制说明

### 3. 错误处理
- 标准HTTP状态码
- 统一的错误响应格式
- 详细的错误描述

### 4. 分页支持
- 分页参数说明
- 响应数据结构完整
- 示例数据准确

## 🚀 使用示例

### 通过Swagger UI测试
1. 访问 `http://localhost:8080/swagger/index.html`
2. 点击"Try it out"按钮
3. 填写请求参数
4. 点击"Execute"执行请求

### 通过curl测试
```bash
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
    "username": "testuser",
    "password": "123456"
  }'

# 获取用户列表（需要认证）
curl -X GET http://localhost:8080/api/v1/users?page=1&limit=10 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 📋 维护建议

### 定期更新
- 每次API变更后更新Swagger文档
- 保持代码注释与文档同步
- 定期验证文档准确性

### 版本管理
- 使用Git管理文档变更
- 在CHANGELOG中记录API变更
- 维护API版本历史

### 团队协作
- 代码审查时检查文档更新
- 使用Pull Request模板提醒文档更新
- 定期同步文档与实现

## 🎯 最佳实践

1. **保持同步**：代码变更时同步更新文档
2. **详细描述**：为每个字段提供清晰描述
3. **示例数据**：提供真实的请求/响应示例
4. **错误处理**：详细说明各种错误情况
5. **版本控制**：使用语义化版本管理API变更

当前Swagger文档已经是最新的，包含了所有已实现的功能。如需添加新功能，请按照上述指南更新文档。