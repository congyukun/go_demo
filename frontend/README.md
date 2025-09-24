# Go Demo Frontend

这是Go Demo项目的Vue.js前端应用，提供用户登录注册功能。

## 功能特性

- 用户登录
- 用户注册
- 用户仪表板
- 用户列表展示
- JWT Token认证
- 响应式设计

## 技术栈

- Vue 3
- Vue Router 4
- Axios
- Vite

## 安装和运行

### 1. 安装依赖

```bash
cd frontend
npm install
```

### 2. 启动开发服务器

```bash
npm run dev
```

应用将在 `http://localhost:3000` 启动

### 3. 构建生产版本

```bash
npm run build
```

## 项目结构

```
frontend/
├── src/
│   ├── components/          # Vue组件
│   │   ├── Login.vue       # 登录页面
│   │   ├── Register.vue    # 注册页面
│   │   └── Dashboard.vue   # 仪表板页面
│   ├── services/           # API服务
│   │   └── api.js         # API请求封装
│   ├── App.vue            # 根组件
│   └── main.js            # 应用入口
├── index.html             # HTML模板
├── vite.config.js         # Vite配置
└── package.json           # 项目配置
```

## API接口

前端应用与后端Go API进行通信，主要接口包括：

- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/logout` - 用户登出
- `GET /api/v1/users` - 获取用户列表

## 使用说明

1. 首次访问会跳转到登录页面
2. 可以使用默认账户登录：
   - 用户名: `admin`
   - 密码: `123456`
3. 也可以注册新账户
4. 登录成功后会跳转到仪表板页面
5. 仪表板显示当前用户信息和所有用户列表

## 注意事项

- 确保后端Go服务已启动在 `http://localhost:8080`
- 前端开发服务器配置了代理，会自动转发API请求到后端
- Token会保存在localStorage中，刷新页面不会丢失登录状态