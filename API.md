# API 接口文档

## 1. 认证接口

### 1.1 用户注册

**请求**
```
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "password123",
  "name": "New User"
}
```

**响应**
```
{
  "code": 0,
  "message": "注册成功",
  "data": {
    "id": 1001,
    "username": "newuser",
    "email": "newuser@example.com",
    "name": "New User",
    "mobile": "",
    "status": 1,
    "last_login": null,
    "created_at": "2023-12-01T10:00:00Z",
    "updated_at": "2023-12-01T10:00:00Z"
  }
}
```

### 1.2 用户登录

**请求**
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}
```

**响应**
```
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2023-12-01T11:00:00Z",
    "refresh_expires_at": "2023-12-08T10:00:00Z",
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "name": "Test User",
      "mobile": "13800138000",
      "status": 1,
      "last_login": "2023-12-01T10:00:00Z",
      "created_at": "2023-12-01T10:00:00Z",
      "updated_at": "2023-12-01T10:00:00Z"
    }
  }
}
```

### 1.3 刷新访问令牌

**请求**
```
POST /api/v1/auth/refresh
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应**
```
{
  "code": 0,
  "message": "令牌刷新成功",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 3600
  }
}
```

### 1.4 用户登出

**请求**
```
POST /api/v1/auth/logout
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应**
```
{
  "code": 0,
  "message": "登出成功",
  "data": null
}
```

### 1.5 获取当前用户信息

**请求**
```
GET /api/v1/auth/profile
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应**
```
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "name": "Test User",
    "mobile": "13800138000",
    "status": 1,
    "last_login": "2023-12-01T10:00:00Z",
    "created_at": "2023-12-01T10:00:00Z",
    "updated_at": "2023-12-01T10:00:00Z"
  }
}
```

## 2. 用户管理接口

### 2.1 获取用户列表

**请求**
```
GET /api/v1/users?page=1&limit=10&keyword=
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应**
```
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "users": [
      {
        "id": 1,
        "username": "testuser",
        "email": "test@example.com",
        "name": "Test User",
        "mobile": "13800138000",
        "status": 1,
        "last_login": "2023-12-01T10:00:00Z",
        "created_at": "2023-12-01T10:00:00Z",
        "updated_at": "2023-12-01T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 100,
      "pages": 10
    }
  }
}
```

### 2.2 获取用户详情

**请求**
```
GET /api/v1/users/1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应**
```
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "name": "Test User",
    "mobile": "13800138000",
    "status": 1,
    "last_login": "2023-12-01T10:00:00Z",
    "created_at": "2023-12-01T10:00:00Z",
    "updated_at": "2023-12-01T10:00:00Z"
  }
}
```

### 2.3 创建用户

**请求**
```
POST /api/v1/users
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "password123",
  "name": "New User",
  "mobile": "13800138001"
}
```

**响应**
```
{
  "code": 0,
  "message": "创建成功",
  "data": {
    "id": 1001,
    "username": "newuser",
    "email": "newuser@example.com",
    "name": "New User",
    "mobile": "13800138001",
    "status": 1,
    "last_login": null,
    "created_at": "2023-12-01T10:00:00Z",
    "updated_at": "2023-12-01T10:00:00Z"
  }
}
```

### 2.4 更新用户信息

**请求**
```
PUT /api/v1/users/1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "name": "Updated User",
  "email": "updated@example.com",
  "mobile": "13800138002"
}
```

**响应**
```
{
  "code": 0,
  "message": "更新成功",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "updated@example.com",
    "name": "Updated User",
    "mobile": "13800138002",
    "status": 1,
    "last_login": "2023-12-01T10:00:00Z",
    "created_at": "2023-12-01T10:00:00Z",
    "updated_at": "2023-12-01T11:00:00Z"
  }
}
```

### 2.5 删除用户

**请求**
```
DELETE /api/v1/users/1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应**
```
{
  "code": 0,
  "message": "删除成功",
  "data": null
}
```

### 2.6 更新个人资料

**请求**
```
PUT /api/v1/users/profile
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "name": "Updated Profile",
  "email": "profile@example.com",
  "mobile": "13800138003"
}
```

**响应**
```
{
  "code": 0,
  "message": "更新成功",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "profile@example.com",
    "name": "Updated Profile",
    "mobile": "13800138003",
    "status": 1,
    "last_login": "2023-12-01T10:00:00Z",
    "created_at": "2023-12-01T10:00:00Z",
    "updated_at": "2023-12-01T11:00:00Z"
  }
}
```

### 2.7 修改密码

**请求**
```
PUT /api/v1/users/password
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "old_password": "oldpassword123",
  "new_password": "newpassword456"
}
```

**响应**
```
{
  "code": 0,
  "message": "密码修改成功",
  "data": null
}
```

### 2.8 获取用户统计信息

**请求**
```
GET /api/v1/users/stats
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应**
```
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "total_users": 100,
    "active_users": 85,
    "admin_users": 0,
    "normal_users": 0
  }
}
```

## 注意事项

1. 所有接口请求都需要通过HTTPS协议进行传输
2. 请求和响应的数据格式均为JSON
3. 认证接口不需要Authorization头，其他接口都需要
4. 时间格式为ISO 8601标准格式
5. 分页参数page从1开始，limit默认为10
6. 某些管理接口可能需要管理员权限