<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <h2>用户登录</h2>
        <p>欢迎回来，请登录您的账户</p>
      </div>
      
      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label for="username">用户名</label>
          <input
            id="username"
            v-model="form.username"
            type="text"
            placeholder="请输入用户名"
            required
            :disabled="loading"
          />
        </div>
        
        <div class="form-group">
          <label for="password">密码</label>
          <input
            id="password"
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            required
            :disabled="loading"
          />
        </div>
        
        <div v-if="error" class="error-message">
          {{ error }}
        </div>
        
        <button type="submit" class="login-btn" :disabled="loading">
          <span v-if="loading">登录中...</span>
          <span v-else>登录</span>
        </button>
      </form>
      
      <div class="login-footer">
        <p>还没有账户？ 
          <router-link to="/register" class="register-link">立即注册</router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script>
import { authAPI } from '../services/api.js'

export default {
  name: 'Login',
  data() {
    return {
      form: {
        username: '',
        password: ''
      },
      loading: false,
      error: ''
    }
  },
  methods: {
    async handleLogin() {
      this.loading = true
      this.error = ''
      
      try {
        const response = await authAPI.login({
          username: this.form.username,
          password: this.form.password
        })
        
        if (response.code === 200) {
          // 保存token和用户信息
          localStorage.setItem('token', response.data.token)
          localStorage.setItem('user', JSON.stringify(response.data.user))
          
          // 跳转到仪表板
          this.$router.push('/dashboard')
        } else {
          this.error = response.message || '登录失败'
        }
      } catch (error) {
        console.error('登录错误:', error)
        this.error = error.response?.data?.message || '登录失败，请检查网络连接'
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 20px;
}

.login-card {
  background: white;
  border-radius: 16px;
  box-shadow: 0 15px 50px rgba(0, 0, 0, 0.15);
  padding: 40px;
  width: 100%;
  max-width: 600px;
}

@media (min-width: 1024px) {
  .login-card {
    max-width: 800px;
    padding: 50px 80px;
  }
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-header h2 {
  color: #333;
  margin-bottom: 12px;
  font-size: 32px;
  font-weight: 600;
}

.login-header p {
  color: #666;
  font-size: 16px;
}

@media (min-width: 1024px) {
  .login-header h2 {
    font-size: 36px;
    margin-bottom: 16px;
  }
  
  .login-header p {
    font-size: 18px;
  }
}

.login-form {
  margin-bottom: 20px;
}

.form-group {
  margin-bottom: 24px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #333;
  font-weight: 500;
  font-size: 16px;
}

.form-group input {
  width: 100%;
  padding: 16px 20px;
  border: 2px solid #e1e5e9;
  border-radius: 10px;
  font-size: 16px;
  transition: border-color 0.3s ease;
}

@media (min-width: 1024px) {
  .form-group {
    margin-bottom: 28px;
  }
  
  .form-group label {
    font-size: 18px;
    margin-bottom: 10px;
  }
  
  .form-group input {
    padding: 18px 24px;
    font-size: 18px;
    border-radius: 12px;
  }
}

.form-group input:focus {
  outline: none;
  border-color: #667eea;
}

.form-group input:disabled {
  background-color: #f5f5f5;
  cursor: not-allowed;
}

.error-message {
  background-color: #fee;
  color: #c53030;
  padding: 16px;
  border-radius: 8px;
  margin-bottom: 24px;
  font-size: 16px;
  text-align: center;
}

.login-btn {
  width: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 18px;
  border-radius: 10px;
  font-size: 18px;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s ease;
}

@media (min-width: 1024px) {
  .error-message {
    padding: 20px;
    font-size: 18px;
    margin-bottom: 28px;
  }
  
  .login-btn {
    padding: 22px;
    font-size: 20px;
    border-radius: 12px;
  }
}

.login-btn:hover:not(:disabled) {
  transform: translateY(-1px);
}

.login-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.login-footer {
  text-align: center;
  margin-top: 20px;
}

.login-footer p {
  color: #666;
  font-size: 16px;
}

.register-link {
  color: #667eea;
  text-decoration: none;
  font-weight: 500;
}

.register-link:hover {
  text-decoration: underline;
}

@media (min-width: 1024px) {
  .login-footer p {
    font-size: 18px;
  }
}
</style>