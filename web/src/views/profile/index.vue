<template>
  <div class="profile-container">
    <el-row :gutter="20">
      <!-- 用户信息卡片 -->
      <el-col :span="8">
        <el-card class="user-card">
          <div class="user-avatar">
            <el-avatar :size="100" icon="UserFilled" />
            <el-upload
              class="avatar-uploader"
              action="#"
              :show-file-list="false"
              :auto-upload="false"
              @change="handleAvatarChange"
            >
              <el-button size="small" type="primary" class="upload-btn">
                更换头像
              </el-button>
            </el-upload>
          </div>
          <div class="user-info">
            <h2>{{ userInfo?.name || userInfo?.username }}</h2>
            <p class="username">@{{ userInfo?.username }}</p>
            <el-divider />
            <div class="info-item">
              <el-icon><Message /></el-icon>
              <span>{{ userInfo?.email }}</span>
            </div>
            <div class="info-item" v-if="userInfo?.mobile">
              <el-icon><Phone /></el-icon>
              <span>{{ userInfo?.mobile }}</span>
            </div>
            <div class="info-item">
              <el-icon><Calendar /></el-icon>
              <span>注册于 {{ userInfo?.created_at }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <!-- 编辑表单 -->
      <el-col :span="16">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>个人资料</span>
            </div>
          </template>
          
          <el-tabs v-model="activeTab">
            <el-tab-pane label="基本信息" name="basic">
              <el-form
                ref="basicFormRef"
                :model="basicForm"
                :rules="basicRules"
                label-width="100px"
                style="max-width: 500px"
              >
                <el-form-item label="用户名">
                  <el-input v-model="basicForm.username" disabled />
                </el-form-item>
                <el-form-item label="姓名" prop="name">
                  <el-input v-model="basicForm.name" placeholder="请输入姓名" />
                </el-form-item>
                <el-form-item label="邮箱" prop="email">
                  <el-input v-model="basicForm.email" placeholder="请输入邮箱" />
                </el-form-item>
                <el-form-item label="手机号" prop="mobile">
                  <el-input v-model="basicForm.mobile" placeholder="请输入手机号" />
                </el-form-item>
                <el-form-item>
                  <el-button type="primary" :loading="basicLoading" @click="handleUpdateBasic">
                    保存修改
                  </el-button>
                </el-form-item>
              </el-form>
            </el-tab-pane>
            
            <el-tab-pane label="修改密码" name="password">
              <el-form
                ref="passwordFormRef"
                :model="passwordForm"
                :rules="passwordRules"
                label-width="100px"
                style="max-width: 500px"
              >
                <el-form-item label="当前密码" prop="old_password">
                  <el-input
                    v-model="passwordForm.old_password"
                    type="password"
                    placeholder="请输入当前密码"
                    show-password
                  />
                </el-form-item>
                <el-form-item label="新密码" prop="new_password">
                  <el-input
                    v-model="passwordForm.new_password"
                    type="password"
                    placeholder="请输入新密码"
                    show-password
                  />
                </el-form-item>
                <el-form-item label="确认密码" prop="confirm_password">
                  <el-input
                    v-model="passwordForm.confirm_password"
                    type="password"
                    placeholder="请再次输入新密码"
                    show-password
                  />
                </el-form-item>
                <el-form-item>
                  <el-button type="primary" :loading="passwordLoading" @click="handleUpdatepassword">
                    修改密码
                  </el-button>
                </el-form-item>
              </el-form>
            </el-tab-pane>
          </el-tabs>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { getProfile } from '@/api/auth'
import { updateProfile, changepassword } from '@/api/user'

const userStore = useUserStore()

const activeTab = ref('basic')
const basicLoading = ref(false)
const passwordLoading = ref(false)
const userInfo = ref(null)

const basicFormRef = ref(null)
const passwordFormRef = ref(null)

const basicForm = reactive({
  username: '',
  name: '',
  email: '',
  mobile: ''
})

const passwordForm = reactive({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

const basicRules = {
  name: [
    { max: 50, message: '姓名长度不能超过 50 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  mobile: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ]
}

const validateConfirmpassword = (rule, value, callback) => {
  if (value !== passwordForm.new_password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const passwordRules = {
  old_password: [
    { required: true, message: '请输入当前密码', trigger: 'blur' }
  ],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 个字符', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    { validator: validateConfirmpassword, trigger: 'blur' }
  ]
}

// 获取用户信息 - GET /api/v1/auth/profile
const fetchUserInfo = async () => {
  try {
    const res = await getProfile()
    userInfo.value = res.data
    
    // 填充表单
    basicForm.username = res.data.username
    basicForm.name = res.data.name
    basicForm.email = res.data.email
    basicForm.mobile = res.data.mobile
  } catch (error) {
    console.error('获取用户信息失败:', error)
  }
}

// 更新基本信息 - PUT /api/v1/users/profile
const handleUpdateBasic = async () => {
  if (!basicFormRef.value) return
  
  await basicFormRef.value.validate(async (valid) => {
    if (!valid) return
    
    basicLoading.value = true
    try {
      const { username, ...data } = basicForm
      await updateProfile(data)
      ElMessage.success('更新成功')
      
      // 刷新用户信息
      await userStore.fetchUserInfo()
      fetchUserInfo()
    } catch (error) {
      console.error('更新失败:', error)
    } finally {
      basicLoading.value = false
    }
  })
}

// 修改密码 - PUT /api/v1/users/password
const handleUpdatepassword = async () => {
  if (!passwordFormRef.value) return
  
  await passwordFormRef.value.validate(async (valid) => {
    if (!valid) return
    
    passwordLoading.value = true
    try {
      await changepassword({
        old_password: passwordForm.old_password,
        new_password: passwordForm.new_password
      })
      ElMessage.success('密码修改成功，请重新登录')
      
      // 清空表单
      passwordForm.old_password = ''
      passwordForm.new_password = ''
      passwordForm.confirm_password = ''
      
      // 登出
      userStore.logout()
    } catch (error) {
      console.error('修改密码失败:', error)
    } finally {
      passwordLoading.value = false
    }
  })
}

// 头像更换
const handleAvatarChange = (file) => {
  ElMessage.info('头像上传功能待实现')
}

onMounted(() => {
  fetchUserInfo()
})
</script>

<style lang="scss" scoped>
.profile-container {
  .user-card {
    text-align: center;
    
    .user-avatar {
      margin-bottom: 20px;
      
      .upload-btn {
        margin-top: 15px;
      }
    }
    
    .user-info {
      h2 {
        margin: 0 0 5px;
        color: #333;
      }
      
      .username {
        color: #999;
        margin: 0 0 15px;
      }
      
      .info-item {
        display: flex;
        align-items: center;
        justify-content: center;
        margin-bottom: 10px;
        color: #666;
        
        .el-icon {
          margin-right: 8px;
          color: #409EFF;
        }
      }
    }
  }
  
  .card-header {
    font-weight: bold;
  }
}
</style>
