# API 文档

本文档详细描述了 Go Demo 项目的所有 API 接口。

## 基础信息

- **Base URL**: `http://localhost:8080`
- **Content-Type**: `application/json`
- **认证方式**: Bearer Token (JWT)

## 响应格式

所有 API 响应都遵循统一的格式：

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {},
  "timestamp": "2023-12-01T10:00:00Z"
}
```

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未认证或认证失败 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 409 | 资源冲突 |
| 500 | 服务器内部错误 |

---

## 认证接口

### 1. 用户注册

**接口地址**: `POST /api/v1/auth/register`

**请求参数**:
```json
{
  "username": "testuser",
  "password": "123456",
  "name": "测试用户",
  "email": "test@example.com",
  "mobile": "13812345678"
}
```

**参数说明**:
- `username` (string, 必填): 用户名，3-20个字符
- `password` (string, 必填): 密码，最少6个字符
- `name` (string, 必填): 真实姓名，1-50个字符
- `email` (string, 可选): 邮箱地址
- `mobile` (string, 必填): 手机号码

**响应示例**:
```json
{
  "code": 200,
  "message": "注册成功",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "phone": "13812345678",
    "status": 1,
    "role": "user",
    "created_at": "2023-12-01T10:00:00Z",
    "updated_at": "2023-12-01T10:00:00Z"
  }
}
```

### 2. 用户登录

**接口地址**: `POST /api/v1/auth/login`

**请求参数**:
```json
{
  "username": "testuser",
  "password": "123456"
}
```

**参数说明**:
- `username` (string, 必填): 用户名
- `password` (string, 必填): 密码

**响应示例**:
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 3600,
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "phone": "13812345678",
      "status": 1,
      "role": "user"
    }
  }
}
```

### 3. 刷新令牌

**接口地址**: `POST /api/v1/auth/refresh`

**请求头**:
```
Authorization: Bearer {refresh_token}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "令牌刷新成功",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer"
  }
}
```

### 4. 用户登出

**接口地址**: `POST /api/v1/auth/logout`

**请求头**:
```
Authorization: Bearer {access_token}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "登出成功",
  "data": null
}
```

### 5. 获取当前用户信息

**接口地址**: `GET /api/v1/auth/profile`

**请求头**:
```
Authorization: Bearer {access_token}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "phone": "13812345678",
    "status": 1,
    "role": "user",
    "created_at": "2023-12-01T10:00:00Z",
    "updated_at": "2023-12-01T10:00:00Z",
    "last_login": "2023-12-01T10:30:00Z"
  }
}
```

---

## 用户管理接口

### 1. 获取用户列表

**接口地址**: `GET /api/v1/users`

**请求头**:
```
Authorization: Bearer {access_token}
```

**查询参数**:
- `page` (int, 可选): 页码，默认为1
- `size` (int, 可选): 每页数量，默认为10，最大100

**请求示例**:
```
GET /api/v1/users?page=1&size=10
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "users": [
      {
        "id": 1,
        "username": "testuser",
        "email": "test@example.com",
        "phone": "13812345678",
        "status": 1,
        "role": "user",
        "created_at": "2023-12-01T10:00:00Z",
        "updated_at": "2023-12-01T10:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "size": 10
  }
}
```

### 2. 获取用户详情

**接口地址**: `GET /api/v1/users/{id}`

**请求头**:
```
Authorization: Bearer {access_token}
```

**路径参数**:
- `id` (int, 必填): 用户ID

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "phone": "13812345678",
    "status": 1,
    "role": "user",
    "created_at": "2023-12-01T10:00:00Z",
    "updated_at": "2023-12-01T10:00:00Z",
    "last_login": "2023-12-01T10:30:00Z"
  }
}
```

### 3. 创建用户

**接口地址**: `POST /api/v1/users`

**请求头**:
```
Authorization: Bearer {access_token}
```

**请求参数**:
```json
{
  "username": "newuser",
  "password": "123456",
  "email": "newuser@example.com"
}
```

**参数说明**:
- `username` (string, 必填): 用户名，3-20个字符
- `password` (string, 必填): 密码，最少6个字符
- `email` (string, 可选): 邮箱地址

**响应示例**:
```json
{
  "code": 200,
  "message": "创建成功",
  "data": {
    "id": 2,
    "username": "newuser",
    "email": "newuser@example.com",
    "phone": "",
    "status": 1,
    "role": "user",
    "created_at": "2023-12-01T11:00:00Z",
    "updated_at": "2023-12-01T11:00:00Z"
  }
}
```

### 4. 更新用户信息

**接口地址**: `PUT /api/v1/users/{id}`

**请求头**:
```
Authorization: Bearer {access_token}
```

**路径参数**:
- `id` (int, 必填): 用户ID

**请求参数**:
```json
{
  "email": "updated@example.com",
  "status": 0
}
```

**参数说明**:
- `email` (string, 可选): 邮箱地址
- `status` (int, 可选): 用户状态，0=禁用，1=启用

**响应示例**:
```json
{
  "code": 200,
  "message": "更新成功",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "updated@example.com",
    "phone": "13812345678",
    "status": 0,
    "role": "user",
    "created_at": "2023-12-01T10:00:00Z",
    "updated_at": "2023-12-01T11:30:00Z"
  }
}
```

### 5. 删除用户

**接口地址**: `DELETE /api/v1/users/{id}`

**请求头**:
```
Authorization: Bearer {access_token}
```

**路径参数**:
- `id` (int, 必填): 用户ID

**响应示例**:
```json
{
  "code": 200,
  "message": "删除成功",
  "data": null
}
```

### 6. 更新个人资料

**接口地址**: `PUT /api/v1/users/profile`

**请求头**:
```
Authorization: Bearer {access_token}
```

**请求参数**:
```json
{
  "email": "newemail@example.com"
}
```

**参数说明**:
- `email` (string, 可选): 邮箱地址

**响应示例**:
```json
{
  "code": 200,
  "message": "更新成功",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "newemail@example.com",
    "phone": "13812345678",
    "status": 1,
    "role": "user",
    "created_at": "2023-12-01T10:00:00Z",
    "updated_at": "2023-12-01T12:00:00Z"
  }
}
```

### 7. 修改密码

**接口地址**: `PUT /api/v1/users/Password`

**请求头**:
```
Authorization: Bearer {access_token}
```

**请求参数**:
```json
{
  "old_password": "123456",
  "new_password": "newpass123"
}
```

**参数说明**:
- `old_password` (string, 必填): 原密码
- `new_password` (string, 必填): 新密码，最少6个字符

**响应示例**:
```json
{
  "code": 200,
  "message": "修改成功",
  "data": null
}
```

### 8. 获取用户统计信息

**接口地址**: `GET /api/v1/users/stats`

**请求头**:
```
Authorization: Bearer {access_token}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "total": 100,
    "active": 85,
    "inactive": 15,
    "today_registered": 5,
    "this_month_registered": 20
  }
}
```

---

## 系统接口

### 健康检查

**接口地址**: `GET /health`

**响应示例**:
```json
{
  "status": "ok",
  "time": "2023-12-01T10:00:00Z"
}
```

---

## 使用示例

### 完整的用户注册登录流程

```bash
# 1. 用户注册
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456",
    "name": "测试用户",
    "email": "test@example.com",
    "mobile": "13812345678"
  }'

# 2. 用户登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456"
  }'

# 3. 使用返回的token访问受保护的接口
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# 4. 获取当前用户信息
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

## 注意事项

1. **认证**: 除了注册、登录和健康检查接口外，其他接口都需要在请求头中携带有效的JWT token
2. **权限**: 某些管理接口可能需要管理员权限
3. **限流**: 生产环境建议对API接口进行限流保护
4. **HTTPS**: 生产环境必须使用HTTPS协议
5. **参数验证**: 所有接口都会进行严格的参数验证
6. **错误处理**: 接口会返回详细的错误信息，便于调试和处理