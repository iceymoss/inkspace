<template>
  <div class="login-page">
    <el-card class="login-card">
      <h2>{{ isRegister ? '注册账号' : '用户登录' }}</h2>
      
      <el-form :model="form" :rules="rules" ref="formRef" @submit.prevent="handleSubmit">
        <el-form-item prop="username">
          <el-input
            v-model="form.username"
            placeholder="用户名"
            :prefix-icon="User"
          />
        </el-form-item>

        <el-form-item v-if="isRegister" prop="email">
          <el-input
            v-model="form.email"
            placeholder="邮箱"
            :prefix-icon="Message"
          />
        </el-form-item>

        <el-form-item v-if="isRegister" prop="nickname">
          <el-input
            v-model="form.nickname"
            placeholder="昵称（选填）"
            :prefix-icon="User"
          />
        </el-form-item>

        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="密码"
            :prefix-icon="Lock"
            show-password
          />
        </el-form-item>

        <el-form-item v-if="isRegister" prop="confirmPassword">
          <el-input
            v-model="form.confirmPassword"
            type="password"
            placeholder="确认密码"
            :prefix-icon="Lock"
            show-password
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" class="submit-btn" @click="handleSubmit" :loading="loading">
            {{ isRegister ? '注册' : '登录' }}
          </el-button>
        </el-form-item>

        <div class="form-footer">
          <el-link type="primary" @click="isRegister = !isRegister">
            {{ isRegister ? '已有账号？去登录' : '没有账号？去注册' }}
          </el-link>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock, Message } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { loadTheme } from '@/utils/theme'

const router = useRouter()
const userStore = useUserStore()

const formRef = ref()
const isRegister = ref(false)
const loading = ref(false)

const form = reactive({
  username: '',
  email: '',
  nickname: '',
  password: '',
  confirmPassword: ''
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== form.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules = reactive({
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 50, message: '用户名长度在 3 到 50 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 50, message: '密码长度在 6 到 50 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
})

const handleSubmit = async () => {
  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      if (isRegister.value) {
        await userStore.register({
          username: form.username,
          email: form.email,
          nickname: form.nickname,
          password: form.password
        })
        ElMessage.success('注册成功！请登录')
        isRegister.value = false
        form.password = ''
        form.confirmPassword = ''
      } else {
        await userStore.login({
          username: form.username,
          password: form.password
        })
        ElMessage.success('登录成功')
        router.push('/')
      }
    } catch (error) {
      // API 拦截器已经处理了错误消息的显示，这里不需要再次显示
      // 只需要处理其他逻辑（如果需要的话）
      console.error('Login/Register error:', error)
    } finally {
      loading.value = false
    }
  })
}

// 加载主题配置
onMounted(() => {
  loadTheme()
})
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--theme-bg-secondary);
}

.login-card {
  width: 100%;
  max-width: 400px;
  margin: var(--spacing-lg);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  background: var(--theme-bg-card);
  transition: box-shadow var(--transition-base);
}

.login-card h2 {
  text-align: center;
  margin-bottom: var(--spacing-xl);
  color: var(--theme-text-primary);
  font-size: var(--font-size-2xl);
  font-weight: 600;
  line-height: var(--line-height-tight);
  font-family: var(--font-sans);
}

.form-footer {
  text-align: center;
  margin-top: var(--spacing-sm);
}

.form-footer :deep(.el-link) {
  cursor: pointer;
  transition: color var(--transition-fast);
}

.submit-btn {
  width: 100%;
  cursor: pointer;
  transition: background-color var(--transition-fast), border-color var(--transition-fast);
}
</style>

<style>
body.theme-day .login-page {
  background: linear-gradient(135deg, var(--theme-bg-secondary) 0%, var(--theme-bg-hover) 50%, var(--theme-bg-secondary) 100%);
}

body.theme-night .login-page {
  background: var(--theme-hero-gradient);
}

body.theme-holiday .login-page {
  background: var(--theme-bg-secondary);
}

body.theme-mourning .login-page {
  background: var(--theme-hero-gradient);
}
</style>

