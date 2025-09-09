# go_demo

## 项目简介
本项目基于 Gin 实现的简单 RESTful API 服务，支持用户注册、登录与文章的增删改查（CRUD）操作，采用内存存储，适合学习和二次开发。

## 主要特性
- 用户注册与登录接口
- 文章的创建、查询、更新、删除、列表接口
- 统一响应格式，基础参数校验
- 结构清晰，便于扩展

## 目录结构
```
.
├── controllers/         # 控制器层，业务逻辑
│   ├── article.go       # 文章相关接口
│   ├── login.go         # 用户注册/登录接口
│   └── init.go          # 统一响应函数
├── routes/              # 路由注册
│   ├── article.go
│   ├── user.go
│   └── route.go
├── registry/            # 预留注册相关
├── main.go              # 程序入口
├── go.mod/go.sum        # Go 依赖管理
└── README.md
```

## 依赖安装
需先安装 Go 1.18+，并拉取依赖：
```bash
go mod tidy
```

## 启动方式
```bash
go run main.go
```
默认监听 8080 端口。

## 主要接口示例

### 用户注册
- `POST /register`
- 请求体：
  ```json
  {
    "username": "test",
    "Password": "123456"
  }
  ```
- 响应：
  ```json
  {
    "code": 200,
    "message": "注册成功!",
    "data": {
      "username": "test",
      "Password": "123456"
    }
  }
  ```

### 用户登录
- `POST /login`
- 请求体同上

### 创建文章
- `POST /article`
- 请求体：
  ```json
  {
    "title": "标题",
    "content": "内容"
  }
  ```

### 获取文章
- `GET /article/{id}`

### 更新文章
- `PUT /article/{id}`

### 删除文章
- `DELETE /article/{id}`

### 文章列表
- `GET /articles`

## 说明
- 所有接口均返回统一 JSON 格式，包含 code、message、data 字段。
- 本项目为演示用途，数据存储于内存，重启即丢失。
- 可根据实际需求扩展数据库、鉴权、日志等功能。
