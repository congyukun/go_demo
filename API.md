# Go Demo API 文档

## 项目结构

```
go_demo/
├── config/                 # 配置文件
│   ├── config.go          # 配置结构体和初始化
│   └── config.yaml        # 配置文件
├── controllers/           # 控制器层
│   ├── article_controller.go  # 文章控制器
│   ├── login_controller.go    # 登录控制器
│   └── user_controller.go     # 用户控制器
├── db/                    # 数据库相关
│   ├── gorm_logger.go     # GORM日志配置
│   └── mysql.go           # MySQL连接
├── logger/                # 日志相关
│   └── zap.go            # Zap日志配置
├── models/                # 数据模型
│   ├── auth.go           # 认证相关模型
│   └── user.go           # 用户模型
├── registry/              # 依赖注入
│   └── registry.go       # 控制器注册
├── routes/                # 路由相关
│   ├── middleware.go     # 中间件
│   └── router.go         # 路由配置
├── services/              # 服务层
│   ├── auth_service.go   # 认证服务
│   └── user_service.go   # 用户服务
├── utils/                 # 工具类
│   └── response.go       # 响应工具
├── logs/                  # 日志文件
├── main.go               # 程序入口
├── go.mod                # Go模块文件
└── README.md             # 项目说明
```

## API 接口

### 基础信息
- 基础URL: `http://localhost:8081`
- 内容类型: `application/json`

### 认证接口

#### 1. 用户登录
- **URL**: `POST /api/v1/auth/login`
- **描述**: 用户登录获取token
- **请求体**:
```json
{
  "username": "admin",
  "password": "123456"
}
```
- **响应**:
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "token_c0c297598732595a97196e21bb0a5666",
    "expires_at": "2025-09-11T11:44:35.311421+08:00",
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "name": "管理员",
      "created_at": "0001-01-01T00:00:00Z",
      "updated_at": "0001-01-01T00:00:00Z"
    }
  },
  "request_id": "20250910114435-WWWWWW"
}
```

#### 2. 用户注册
- **URL**: `POST /api/v1/auth/register`
- **描述**: 用户注册
- **请求体**:
```json
{
  "username": "testuser",
  "password": "123456",
  "email": "test@example.com",
  "name": "测试用户"
}
```
- **响应**:
```json
{
  "code": 201,
  "message": "注册成功",
  "data": {
    "id": 5,
    "username": "testuser",
    "email": "test@example.com",
    "name": "测试用户",
    "created_at": "0001-01-01T00:00:00Z",
    "updated_at": "0001-01-01T00:00:00Z"
  },
  "request_id": "20250910114501-kkkkkk"
}
```

#### 3. 用户登出
- **URL**: `POST /api/v1/auth/logout`
- **描述**: 用户登出
- **请求头**: `Authorization: Bearer <token>`
- **响应**:
```json
{
  "code": 200,
  "message": "登出成功",
  "request_id": "xxx"
}
```

### 用户管理接口

#### 1. 获取用户列表
- **URL**: `GET /api/v1/users`
- **描述**: 获取所有用户列表
- **响应**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "name": "管理员",
      "created_at": "0001-01-01T00:00:00Z",
      "updated_at": "0001-01-01T00:00:00Z"
    }
  ],
  "request_id": "xxx"
}
```

#### 2. 获取单个用户
- **URL**: `GET /api/v1/users/{id}`
- **描述**: 根据ID获取用户信息
- **响应**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "name": "管理员",
    "created_at": "0001-01-01T00:00:00Z",
    "updated_at": "0001-01-01T00:00:00Z"
  },
  "request_id": "xxx"
}
```

#### 3. 创建用户
- **URL**: `POST /api/v1/users`
- **描述**: 创建新用户
- **请求体**:
```json
{
  "username": "newuser",
  "email": "newuser@example.com",
  "name": "新用户",
  "password": "123456"
}
```

#### 4. 更新用户
- **URL**: `PUT /api/v1/users/{id}`
- **描述**: 更新用户信息
- **请求体**:
```json
{
  "email": "updated@example.com",
  "name": "更新后的名称"
}
```

#### 5. 删除用户
- **URL**: `DELETE /api/v1/users/{id}`
- **描述**: 删除用户

### 文章管理接口

#### 1. 获取文章列表
- **URL**: `GET /api/v1/articles`

#### 2. 获取单个文章
- **URL**: `GET /api/v1/articles/{id}`

#### 3. 创建文章
- **URL**: `POST /api/v1/articles`

#### 4. 更新文章
- **URL**: `PUT /api/v1/articles/{id}`

#### 5. 删除文章
- **URL**: `DELETE /api/v1/articles/{id}`

### 其他接口

#### 健康检查
- **URL**: `GET /health`
- **描述**: 服务健康检查
- **响应**:
```json
{
  "code": 200,
  "message": "服务正常运行",
  "timestamp": {},
  "request_id": "xxx"
}
```

## 错误码说明

- `200`: 成功
- `201`: 创建成功
- `400`: 请求参数错误
- `401`: 未授权
- `404`: 资源不存在
- `409`: 资源冲突（如用户名已存在）
- `500`: 服务器内部错误

## 架构说明

### 分层架构
1. **Controller层**: 处理HTTP请求，参数验证，调用Service层
2. **Service层**: 业务逻辑处理，调用数据访问层
3. **Model层**: 数据模型定义
4. **Utils层**: 工具类和通用函数

### 中间件
- **ZapLogger**: 请求日志记录
- **ZapRecovery**: Panic恢复
- **RequestID**: 请求ID生成
- **CORS**: 跨域处理

### 特性
- 结构化日志（Zap）
- 请求ID追踪
- 统一错误处理
- 参数验证
- 数据库连接池
- 配置文件管理

## 运行说明

1. 启动服务：
```bash
go run main.go
```

2. 服务将在 `http://localhost:8081` 启动

3. 测试登录：
```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}'