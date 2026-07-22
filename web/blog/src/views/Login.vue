<template>
  <div class="login-page">
    <el-card class="login-card">
      <div class="login-kicker">
        INKSPACE · MEMBER
      </div>
      <h2>{{ isRegister ? '注册账号' : '用户登录' }}</h2>
      
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        @submit.prevent="handleSubmit"
      >
        <el-form-item prop="username">
          <el-input
            v-model="form.username"
            placeholder="用户名"
            :prefix-icon="User"
          />
        </el-form-item>

        <el-form-item
          v-if="isRegister"
          prop="email"
        >
          <el-input
            v-model="form.email"
            placeholder="邮箱"
            :prefix-icon="Message"
          />
        </el-form-item>

        <el-form-item
          v-if="isRegister"
          prop="nickname"
        >
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

        <el-form-item
          v-if="isRegister"
          prop="confirmPassword"
        >
          <el-input
            v-model="form.confirmPassword"
            type="password"
            placeholder="确认密码"
            :prefix-icon="Lock"
            show-password
          />
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            style="width: 100%"
            :loading="loading"
            @click="handleSubmit"
          >
            {{ isRegister ? '注册' : '登录' }}
          </el-button>
        </el-form-item>

        <div class="form-footer">
          <el-link
            type="primary"
            @click="isRegister = !isRegister"
          >
            {{ isRegister ? '已有账号？去登录' : '没有账号？去注册' }}
          </el-link>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock, Message } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'

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
  margin: 20px;
}

.login-card h2 {
  text-align: center;
  margin-bottom: 30px;
  color: var(--theme-text-primary);
}

.form-footer {
  text-align: center;
  margin-top: 10px;
}

/* Magazine adaptation */
.login-page { min-height: calc(100vh - 64px); padding: 64px 24px; background: var(--theme-bg-primary); }
.login-card { max-width: 430px; margin: 0; border: 1px solid var(--theme-border); border-radius: 0; box-shadow: none; background: transparent; }
.login-card :deep(.el-card__body) { padding: 48px 46px; }
.login-kicker { margin-bottom: 12px; color: var(--theme-primary); font-family: Georgia, 'Songti SC', serif; font-size: 11px; letter-spacing: .24em; text-align: center; }
.login-card h2 { margin-bottom: 38px; font-family: Georgia, 'Songti SC', 'Noto Serif SC', SimSun, serif; font-size: 34px; font-weight: 500; letter-spacing: .05em; }
.login-card :deep(.el-input__wrapper), .login-card :deep(.el-button) { border-radius: 1px; box-shadow: none; }
.login-card :deep(.el-input__wrapper) { background: transparent; border-bottom: 1px solid var(--theme-border); }
.form-footer :deep(.el-link) { color: var(--theme-text-secondary); }
.form-footer :deep(.el-link:hover) { color: var(--theme-primary); }
@media (max-width: 560px) { .login-page { align-items: flex-start; padding: 32px 18px; } .login-card { border-right: 0; border-left: 0; } .login-card :deep(.el-card__body) { padding: 38px 18px; } }
</style>

