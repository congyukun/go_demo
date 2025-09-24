<template>
  <div class="register-container">
    <div class="register-card">
      <div class="register-header">
        <h2>用户注册</h2>
        <p>创建您的新账户</p>
      </div>
      <form @submit.prevent="handleRegister" class="register-form">
        <div class="form-grid">
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
            <label for="email">邮箱</label>
            <input
              id="email"
              v-model="form.email"
              type="email"
              placeholder="请输入邮箱地址"
              required
              :disabled="loading"
            />
          </div>
          
          <div class="form-group">
            <label for="name">姓名</label>
            <input
              id="name"
              v-model="form.name"
              type="text"
              placeholder="请输入您的姓名"
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
          
          <div class="form-group form-group-full">
            <label for="confirmpassword">确认密码</label>
            <input
              id="confirmpassword"
              v-model="form.confirmpassword"
              type="password"
              placeholder="请再次输入密码"
              required
              :disabled="loading"
            />
          </div>
        </div>
        
        <div v-if="error" class="error-message">
          {{ error }}
        </div>
        
        <div v-if="success" class="success-message">
          {{ success }}
        </div>
        
        <button type="submit" class="register-btn" :disabled="loading">
          <span v-if="loading">注册中...</span>
          <span v-else>注册</span>
        </button>
      </form>
      
      <div class="register-footer">
        <p>已有账户？ 
          <router-link to="/login" class="login-link">立即登录</router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script>
import { authAPI } from '../services/api.js'

export default {
  name: 'Register',
  data() {
    return {
      form: {
        username: '',
        email: '',
        name: '',
        password: '',
        confirmpassword: ''
      },
      loading: false,
      error: '',
      success: ''
    }
  },
  methods: {
    async handleRegister() {
      this.loading = true
      this.error = ''
      this.success = ''
      
      // 验证密码匹配
      if (this.form.password !== this.form.confirmpassword) {
        this.error = '两次输入的密码不一致'
        this.loading = false
        return
      }
      
      // 验证密码长度
      if (this.form.password.length < 6) {
        this.error = '密码长度至少6位'
        this.loading = false
        return
      }
      
      try {
        const response = await authAPI.register({
          username: this.form.username,
          email: this.form.email,
          name: this.form.name,
          password: this.form.password
        })
        
        if (response.code === 201) {
          this.success = '注册成功！3秒后跳转到登录页面...'
          
          // 3秒后跳转到登录页面
          setTimeout(() => {
            this.$router.push('/login')
          }, 3000)
        } else {
          this.error = response.message || '注册失败'
        }
      } catch (error) {
        console.error('注册错误:', error)
        this.error = error.response?.data?.message || '注册失败，请检查网络连接'
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>
.register-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 20px;
}

.register-card {
  background: white;
  border-radius: 16px;
  box-shadow: 0 15px 50px rgba(0, 0, 0, 0.15);
  padding: 40px;
  width: 100%;
  max-width: 650px;
}

@media (min-width: 1024px) {
  .register-card {
    max-width: 900px;
    padding: 50px 80px;
  }
}

.register-header {
  text-align: center;
  margin-bottom: 30px;
}

.register-header h2 {
  color: #333;
  margin-bottom: 12px;
  font-size: 32px;
  font-weight: 600;
}

.register-header p {
  color: #666;
  font-size: 16px;
}

.register-form {
  margin-bottom: 24px;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 24px;
}

@media (min-width: 768px) {
  .form-grid {
    grid-template-columns: 1fr 1fr;
    gap: 20px 30px;
  }
  
  .form-group-full {
    grid-column: 1 / -1;
  }
}

.form-group {
  margin-bottom: 0;
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
  .register-header h2 {
    font-size: 36px;
    margin-bottom: 16px;
  }
  
  .register-header p {
    font-size: 18px;
  }
  
  .register-form {
    margin-bottom: 28px;
  }
  
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

.success-message {
  background-color: #f0fff4;
  color: #38a169;
  padding: 16px;
  border-radius: 8px;
  margin-bottom: 24px;
  font-size: 16px;
  text-align: center;
}

.register-btn {
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
  
  .success-message {
    padding: 20px;
    font-size: 18px;
    margin-bottom: 28px;
  }
  
  .register-btn {
    padding: 22px;
    font-size: 20px;
    border-radius: 12px;
  }
}

.register-btn:hover:not(:disabled) {
  transform: translateY(-1px);
}

.register-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.register-footer {
  text-align: center;
  margin-top: 20px;
}

.register-footer p {
  color: #666;
  font-size: 16px;
}

.login-link {
  color: #667eea;
  text-decoration: none;
  font-weight: 500;
}

.login-link:hover {
  text-decoration: underline;
}

@media (min-width: 1024px) {
  .register-footer p {
    font-size: 18px;
  }
}
</style>