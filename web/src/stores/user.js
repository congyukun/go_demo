import { defineStore } from 'pinia'
import { login as loginApi, logout as logoutApi, getProfile } from '@/api/auth'
import router from '@/router'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    refreshToken: localStorage.getItem('refreshToken') || '',
    userInfo: JSON.parse(localStorage.getItem('userInfo') || 'null')
  }),
  
  getters: {
    isLoggedIn: (state) => !!state.token,
    username: (state) => state.userInfo?.username || '',
    userId: (state) => state.userInfo?.id || null
  },
  
  actions: {
    // 登录
    async login(loginForm) {
      try {
        const res = await loginApi(loginForm)
        // 后端返回格式: { token, refresh_token, expires_at, refresh_expires_at, user }
        const { token, refresh_token, user } = res.data
        
        this.token = token
        this.refreshToken = refresh_token
        this.userInfo = user
        
        // 持久化存储
        localStorage.setItem('token', token)
        localStorage.setItem('refreshToken', refresh_token)
        localStorage.setItem('userInfo', JSON.stringify(user))
        
        return res
      } catch (error) {
        throw error
      }
    },
    
    // 获取用户信息
    async fetchUserInfo() {
      try {
        const res = await getProfile()
        this.userInfo = res.data
        localStorage.setItem('userInfo', JSON.stringify(res.data))
        return res.data
      } catch (error) {
        throw error
      }
    },
    
    // 登出
    async logout() {
      try {
        if (this.token) {
          await logoutApi()
        }
      } catch (error) {
        console.error('登出请求失败:', error)
      } finally {
        this.resetState()
        router.push('/login')
      }
    },
    
    // 重置状态
    resetState() {
      this.token = ''
      this.refreshToken = ''
      this.userInfo = null
      
      localStorage.removeItem('token')
      localStorage.removeItem('refreshToken')
      localStorage.removeItem('userInfo')
    },
    
    // 更新 Token
    setToken(token) {
      this.token = token
      localStorage.setItem('token', token)
    }
  }
})
