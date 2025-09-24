import axios from 'axios'

// 创建axios实例
const api = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
api.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  response => {
    return response.data
  },
  error => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// 认证相关API
export const authAPI = {
  // 登录
  login(credentials) {
    return api.post('/auth/login', credentials)
  },
  
  // 注册
  register(userData) {
    return api.post('/auth/register', userData)
  },
  
  // 登出
  logout() {
    return api.post('/auth/logout')
  }
}

// 用户相关API
export const userAPI = {
  // 获取用户列表
  getUsers() {
    return api.get('/users')
  },
  
  // 获取单个用户
  getUser(id) {
    return api.get(`/users/${id}`)
  }
}

export default api