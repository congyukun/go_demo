<template>
  <div class="dashboard-container">
    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stat-cards">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background-color: #409EFF;">
              <el-icon><User /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.total_users || 0 }}</div>
              <div class="stat-label">总用户数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background-color: #67C23A;">
              <el-icon><CircleCheck /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.active_users || 0 }}</div>
              <div class="stat-label">活跃用户</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background-color: #E6A23C;">
              <el-icon><UserFilled /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.admin_users || 0 }}</div>
              <div class="stat-label">管理员</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background-color: #F56C6C;">
              <el-icon><TrendCharts /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.normal_users || 0 }}</div>
              <div class="stat-label">普通用户</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 欢迎信息 -->
    <el-card class="welcome-card">
      <template #header>
        <div class="card-header">
          <span>欢迎回来</span>
        </div>
      </template>
      <div class="welcome-content">
        <el-avatar :size="64" icon="UserFilled" />
        <div class="welcome-info">
          <h2>{{ userStore.userInfo?.name || userStore.username }}</h2>
          <p>{{ greeting }}，祝您工作愉快！</p>
          <p class="login-time">
            <el-icon><Clock /></el-icon>
            上次登录：{{ userStore.userInfo?.last_login || '首次登录' }}
          </p>
        </div>
      </div>
    </el-card>
    
    <!-- 快捷操作 -->
    <el-card class="quick-actions">
      <template #header>
        <div class="card-header">
          <span>快捷操作</span>
        </div>
      </template>
      <el-row :gutter="20">
        <el-col :span="6">
          <div class="action-item" @click="$router.push('/users')">
            <el-icon size="32" color="#409EFF"><User /></el-icon>
            <span>用户管理</span>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="action-item" @click="$router.push('/profile')">
            <el-icon size="32" color="#67C23A"><Setting /></el-icon>
            <span>个人设置</span>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="action-item" @click="handleRefresh">
            <el-icon size="32" color="#E6A23C"><Refresh /></el-icon>
            <span>刷新数据</span>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="action-item" @click="handleLogout">
            <el-icon size="32" color="#F56C6C"><SwitchButton /></el-icon>
            <span>退出登录</span>
          </div>
        </el-col>
      </el-row>
    </el-card>
    
    <!-- 最近用户 -->
    <el-card class="recent-users">
      <template #header>
        <div class="card-header">
          <span>最近注册用户</span>
          <el-button type="primary" link @click="$router.push('/users')">
            查看全部
          </el-button>
        </div>
      </template>
      <el-table :data="recentUsers" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="email" label="邮箱" />
        <el-table-column prop="name" label="姓名" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="注册时间" width="180" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessageBox } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { getUsers, getUserStats } from '@/api/user'

const userStore = useUserStore()

const loading = ref(false)
const stats = ref({})
const recentUsers = ref([])

const greeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 6) return '凌晨好'
  if (hour < 9) return '早上好'
  if (hour < 12) return '上午好'
  if (hour < 14) return '中午好'
  if (hour < 17) return '下午好'
  if (hour < 19) return '傍晚好'
  return '晚上好'
})

const fetchData = async () => {
  loading.value = true
  try {
    // 获取用户统计 - GET /api/v1/users/stats
    const statsRes = await getUserStats()
    stats.value = statsRes.data || {}
    
    // 获取最近用户 - GET /api/v1/users
    const usersRes = await getUsers({ page: 1, size: 5 })
    // 后端返回格式: { users: [], total: number, page: number, size: number }
    recentUsers.value = usersRes.data?.users || []
  } catch (error) {
    console.error('获取数据失败:', error)
  } finally {
    loading.value = false
  }
}

const handleRefresh = () => {
  fetchData()
}

const handleLogout = () => {
  ElMessageBox.confirm('确定要退出登录吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    userStore.logout()
  }).catch(() => {})
}

onMounted(() => {
  fetchData()
})
</script>

<style lang="scss" scoped>
.dashboard-container {
  .stat-cards {
    margin-bottom: 20px;
  }
  
  .stat-card {
    .stat-content {
      display: flex;
      align-items: center;
      
      .stat-icon {
        width: 60px;
        height: 60px;
        border-radius: 10px;
        display: flex;
        align-items: center;
        justify-content: center;
        margin-right: 15px;
        
        .el-icon {
          font-size: 28px;
          color: #fff;
        }
      }
      
      .stat-info {
        .stat-value {
          font-size: 28px;
          font-weight: bold;
          color: #333;
        }
        
        .stat-label {
          font-size: 14px;
          color: #999;
          margin-top: 5px;
        }
      }
    }
  }
  
  .welcome-card {
    margin-bottom: 20px;
    
    .welcome-content {
      display: flex;
      align-items: center;
      
      .welcome-info {
        margin-left: 20px;
        
        h2 {
          margin: 0 0 10px;
          color: #333;
        }
        
        p {
          margin: 0;
          color: #666;
          
          &.login-time {
            margin-top: 10px;
            color: #999;
            font-size: 13px;
            display: flex;
            align-items: center;
            
            .el-icon {
              margin-right: 5px;
            }
          }
        }
      }
    }
  }
  
  .quick-actions {
    margin-bottom: 20px;
    
    .action-item {
      display: flex;
      flex-direction: column;
      align-items: center;
      padding: 20px;
      cursor: pointer;
      border-radius: 8px;
      transition: all 0.3s;
      
      &:hover {
        background-color: #f5f7fa;
      }
      
      span {
        margin-top: 10px;
        color: #666;
      }
    }
  }
  
  .recent-users {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
  }
}
</style>
