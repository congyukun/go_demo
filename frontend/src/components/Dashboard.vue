<template>
  <div class="dashboard-container">
    <div class="dashboard-header">
      <h1>欢迎回来，{{ user?.name || '用户' }}！</h1>
      <button @click="handleLogout" class="logout-btn">退出登录</button>
    </div>
    
    <div class="dashboard-content">
      <div class="user-info-card">
        <h2>用户信息</h2>
        <div class="user-details" v-if="user">
          <div class="detail-item">
            <label>用户ID:</label>
            <span>{{ user.id }}</span>
          </div>
          <div class="detail-item">
            <label>用户名:</label>
            <span>{{ user.username }}</span>
          </div>
          <div class="detail-item">
            <label>邮箱:</label>
            <span>{{ user.email }}</span>
          </div>
          <div class="detail-item">
            <label>姓名:</label>
            <span>{{ user.name }}</span>
          </div>
          <div class="detail-item">
            <label>创建时间:</label>
            <span>{{ formatDate(user.created_at) }}</span>
          </div>
        </div>
      </div>
      
      <div class="users-list-card">
        <h2>用户列表</h2>
        <div v-if="loading" class="loading">加载中...</div>
        <div v-else-if="error" class="error">{{ error }}</div>
        <div v-else class="users-grid">
          <div v-for="u in users" :key="u.id" class="user-card">
            <div class="user-avatar">{{ u.name?.charAt(0) || u.username?.charAt(0) }}</div>
            <div class="user-info">
              <h3>{{ u.name }}</h3>
              <p>@{{ u.username }}</p>
              <p class="user-email">{{ u.email }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { authAPI, userAPI } from '../services/api.js'

export default {
  name: 'Dashboard',
  data() {
    return {
      user: null,
      users: [],
      loading: false,
      error: ''
    }
  },
  async mounted() {
    this.loadUserInfo()
    await this.loadUsers()
  },
  methods: {
    loadUserInfo() {
      const userStr = localStorage.getItem('user')
      if (userStr) {
        this.user = JSON.parse(userStr)
      }
    },
    
    async loadUsers() {
      this.loading = true
      this.error = ''
      
      try {
        const response = await userAPI.getUsers()
        if (response.code === 200) {
          this.users = response.data || []
        } else {
          this.error = response.message || '获取用户列表失败'
        }
      } catch (error) {
        console.error('获取用户列表错误:', error)
        this.error = '获取用户列表失败'
      } finally {
        this.loading = false
      }
    },
    
    async handleLogout() {
      try {
        await authAPI.logout()
      } catch (error) {
        console.error('登出错误:', error)
      } finally {
        // 清除本地存储
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        
        // 跳转到登录页面
        this.$router.push('/login')
      }
    },
    
    formatDate(dateStr) {
      if (!dateStr || dateStr === '0001-01-01T00:00:00Z') {
        return '未设置'
      }
      return new Date(dateStr).toLocaleString('zh-CN')
    }
  }
}
</script>

<style scoped>
.dashboard-container {
  min-height: 100vh;
  background: #f5f7fa;
  padding: 20px;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: white;
  padding: 20px 30px;
  border-radius: 12px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  margin-bottom: 30px;
}

.dashboard-header h1 {
  color: #333;
  font-size: 24px;
  font-weight: 600;
  margin: 0;
}

.logout-btn {
  background: #dc3545;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: background-color 0.3s ease;
}

.logout-btn:hover {
  background: #c82333;
}

.dashboard-content {
  display: grid;
  grid-template-columns: 1fr 2fr;
  gap: 30px;
}

@media (max-width: 768px) {
  .dashboard-content {
    grid-template-columns: 1fr;
  }
}

.user-info-card,
.users-list-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  padding: 30px;
}

.user-info-card h2,
.users-list-card h2 {
  color: #333;
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 20px;
  border-bottom: 2px solid #f0f0f0;
  padding-bottom: 10px;
}

.user-details {
  space-y: 15px;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.detail-item:last-child {
  border-bottom: none;
}

.detail-item label {
  font-weight: 500;
  color: #666;
  min-width: 80px;
}

.detail-item span {
  color: #333;
  font-weight: 400;
}

.loading {
  text-align: center;
  color: #666;
  padding: 20px;
}

.error {
  text-align: center;
  color: #dc3545;
  padding: 20px;
  background: #fee;
  border-radius: 6px;
}

.users-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 20px;
}

.user-card {
  display: flex;
  align-items: center;
  padding: 15px;
  border: 1px solid #e1e5e9;
  border-radius: 8px;
  transition: box-shadow 0.3s ease;
}

.user-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.user-avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: 600;
  margin-right: 15px;
  text-transform: uppercase;
}

.user-info h3 {
  margin: 0 0 5px 0;
  color: #333;
  font-size: 16px;
  font-weight: 600;
}

.user-info p {
  margin: 0 0 3px 0;
  color: #666;
  font-size: 14px;
}

.user-email {
  color: #999 !important;
  font-size: 12px !important;
}
</style>