<template>
  <div class="profile-edit">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>编辑个人资料</span>
          <el-button text @click="$router.back()">返回</el-button>
        </div>
      </template>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
        style="max-width: 600px"
      >
        <el-form-item label="用户名">
          <el-input v-model="currentUser.username" disabled />
          <div class="form-tip">用户名不可修改</div>
        </el-form-item>

        <el-form-item label="昵称" prop="nickname">
          <el-input
            v-model="form.nickname"
            placeholder="请输入昵称"
            maxlength="50"
            show-word-limit
          />
        </el-form-item>

        <el-form-item label="邮箱" prop="email">
          <el-input
            v-model="form.email"
            placeholder="请输入邮箱"
            type="email"
          />
        </el-form-item>

        <el-form-item label="头像" prop="avatar">
          <div class="avatar-upload">
            <el-upload
              class="avatar-uploader"
              :action="uploadUrl"
              :headers="uploadHeaders"
              :show-file-list="false"
              :on-success="handleAvatarSuccess"
              :before-upload="beforeAvatarUpload"
              accept="image/jpeg,image/jpg,image/png,image/webp"
            >
              <el-image
                v-if="form.avatar"
                :src="form.avatar"
                fit="cover"
                class="avatar-image"
              >
                <template #error>
                  <div class="image-error">
                    <el-icon><icon-picture /></el-icon>
                  </div>
                </template>
              </el-image>
              <el-icon v-else class="avatar-uploader-icon">
                <Plus />
              </el-icon>
            </el-upload>
            <div class="avatar-tips">
              <div>点击上传头像</div>
              <div>支持 jpg, png, webp 格式，文件不超过2MB</div>
            </div>
          </div>
        </el-form-item>

        <el-form-item label="个人简介" prop="bio">
          <el-input
            v-model="form.bio"
            type="textarea"
            :rows="4"
            placeholder="请输入个人简介"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="loading">
            保存修改
          </el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 修改密码区域 -->
    <el-card style="margin-top: 20px">
      <template #header>
        <span>修改密码</span>
      </template>

      <el-form
        ref="passwordFormRef"
        :model="passwordForm"
        :rules="passwordRules"
        label-width="100px"
        style="max-width: 600px"
      >
        <el-form-item label="当前密码" prop="oldPassword">
          <el-input
            v-model="passwordForm.oldPassword"
            type="password"
            placeholder="请输入当前密码"
            show-password
          />
        </el-form-item>

        <el-form-item label="新密码" prop="newPassword">
          <el-input
            v-model="passwordForm.newPassword"
            type="password"
            placeholder="请输入新密码（6-50位）"
            show-password
          />
        </el-form-item>

        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input
            v-model="passwordForm.confirmPassword"
            type="password"
            placeholder="请再次输入新密码"
            show-password
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handlePasswordChange" :loading="passwordLoading">
            修改密码
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import { Picture as IconPicture, Plus } from '@element-plus/icons-vue'
import api from '@/utils/api'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref(null)
const passwordFormRef = ref(null)
const loading = ref(false)
const passwordLoading = ref(false)

// 上传配置
const uploadUrl = computed(() => {
  // 使用vite proxy
  return '/api/upload/avatar'
})

const uploadHeaders = computed(() => {
  return {
    Authorization: `Bearer ${userStore.token}`
  }
})

const currentUser = reactive({
  username: '',
  nickname: '',
  email: '',
  avatar: '',
  bio: ''
})

const form = reactive({
  nickname: '',
  email: '',
  avatar: '',
  bio: ''
})

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const rules = {
  nickname: [
    { max: 50, message: '昵称最多50个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  bio: [
    { max: 500, message: '个人简介最多500个字符', trigger: 'blur' }
  ]
}

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== passwordForm.newPassword) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const passwordRules = {
  oldPassword: [
    { required: true, message: '请输入当前密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 50, message: '密码长度为6-50位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

// 获取用户信息
const fetchUserProfile = async () => {
  try {
    const data = await api.get('/profile')
    Object.assign(currentUser, data.data)
    Object.assign(form, {
      nickname: data.data.nickname || '',
      email: data.data.email || '',
      avatar: data.data.avatar || '',
      bio: data.data.bio || ''
    })
  } catch (error) {
    ElMessage.error('获取用户信息失败')
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      const data = await api.put('/profile', form)
      ElMessage.success('保存成功')
      
      // 更新store中的用户信息
      userStore.setUser(data.data)
      
      // 刷新当前数据
      Object.assign(currentUser, data.data)
    } catch (error) {
      ElMessage.error(error.message || '保存失败')
    } finally {
      loading.value = false
    }
  })
}

// 重置表单
const handleReset = () => {
  Object.assign(form, {
    nickname: currentUser.nickname || '',
    email: currentUser.email || '',
    avatar: currentUser.avatar || '',
    bio: currentUser.bio || ''
  })
}

// 修改密码
const handlePasswordChange = async () => {
  if (!passwordFormRef.value) return
  
  await passwordFormRef.value.validate(async (valid) => {
    if (!valid) return

    passwordLoading.value = true
    try {
      await api.put('/profile/password', {
        old_password: passwordForm.oldPassword,
        new_password: passwordForm.newPassword
      })
      
      ElMessage.success('密码修改成功，请重新登录')
      
      // 清空表单
      passwordForm.oldPassword = ''
      passwordForm.newPassword = ''
      passwordForm.confirmPassword = ''
      
      // 退出登录
      setTimeout(() => {
        userStore.logout()
        router.push('/login')
      }, 1500)
    } catch (error) {
      ElMessage.error(error.message || '密码修改失败')
    } finally {
      passwordLoading.value = false
    }
  })
}

// 头像上传前验证
const beforeAvatarUpload = (file) => {
  const isImage = ['image/jpeg', 'image/jpg', 'image/png', 'image/webp'].includes(file.type)
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isImage) {
    ElMessage.error('只支持上传 jpg, png, webp 格式的图片')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('图片大小不能超过 2MB')
    return false
  }
  return true
}

// 头像上传成功
const handleAvatarSuccess = (response) => {
  if (response.code === 0 && response.data) {
    // 后端返回的是相对路径，需要拼接完整URL
    const avatarUrl = `http://localhost:8081${response.data.url}`
    form.avatar = avatarUrl
    ElMessage.success('头像上传成功')
  } else {
    ElMessage.error(response.message || '头像上传失败')
  }
}

onMounted(() => {
  fetchUserProfile()
})
</script>

<style scoped>
.profile-edit {
  max-width: 800px;
  margin: 0 auto;
  padding: 40px 20px;
  background-color: #f5f7fa;
  min-height: 100vh;
}

.profile-edit .el-card {
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
}

.avatar-upload {
  display: flex;
  align-items: flex-start;
  gap: 20px;
}

.avatar-uploader {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: border-color 0.3s;
}

.avatar-uploader:hover {
  border-color: #409eff;
}

.avatar-image {
  width: 120px;
  height: 120px;
  display: block;
}

.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 120px;
  height: 120px;
  text-align: center;
  line-height: 120px;
}

.avatar-tips {
  font-size: 12px;
  color: #909399;
  line-height: 1.8;
}

.image-error {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  background: #f5f7fa;
  color: #909399;
  font-size: 30px;
}

:deep(.el-textarea__inner) {
  font-family: inherit;
}
</style>

