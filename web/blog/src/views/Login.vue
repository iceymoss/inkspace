<template>
  <div class="login-page">
    <Card class="login-card">
      <CardHeader class="text-center">
        <CardTitle>{{ isRegister ? '注册账号' : '用户登录' }}</CardTitle>
      </CardHeader>
      <CardContent>
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <div class="space-y-2">
            <div class="input-wrapper">
              <User class="input-icon" />
              <Input v-model="form.username" placeholder="用户名" class="pl-9" />
            </div>
            <p v-if="errors.username" class="text-sm text-destructive">{{ errors.username }}</p>
          </div>

          <div v-if="isRegister" class="space-y-2">
            <div class="input-wrapper">
              <Mail class="input-icon" />
              <Input v-model="form.email" placeholder="邮箱" class="pl-9" />
            </div>
            <p v-if="errors.email" class="text-sm text-destructive">{{ errors.email }}</p>
          </div>

          <div v-if="isRegister" class="space-y-2">
            <div class="input-wrapper">
              <User class="input-icon" />
              <Input v-model="form.nickname" placeholder="昵称（选填）" class="pl-9" />
            </div>
          </div>

          <div class="space-y-2">
            <div class="input-wrapper">
              <Lock class="input-icon" />
              <Input v-model="form.password" type="password" placeholder="密码" class="pl-9" />
            </div>
            <p v-if="errors.password" class="text-sm text-destructive">{{ errors.password }}</p>
          </div>

          <div v-if="isRegister" class="space-y-2">
            <div class="input-wrapper">
              <Lock class="input-icon" />
              <Input v-model="form.confirmPassword" type="password" placeholder="确认密码" class="pl-9" />
            </div>
            <p v-if="errors.confirmPassword" class="text-sm text-destructive">{{ errors.confirmPassword }}</p>
          </div>

          <Button type="submit" class="w-full" :disabled="loading">
            <span v-if="loading" class="mr-2 animate-spin">⟳</span>
            {{ isRegister ? '注册' : '登录' }}
          </Button>

          <div class="form-footer">
            <button type="button" class="switch-link" @click="isRegister = !isRegister">
              {{ isRegister ? '已有账号？去登录' : '没有账号？去注册' }}
            </button>
          </div>
        </form>
      </CardContent>
    </Card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { User, Lock, Mail } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardHeader, CardContent, CardTitle } from '@/components/ui/card'
import { useUserStore } from '@/stores/user'
import { loadTheme } from '@/utils/theme'

const router = useRouter()
const userStore = useUserStore()

const isRegister = ref(false)
const loading = ref(false)
const errors = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const form = reactive({
  username: '',
  email: '',
  nickname: '',
  password: '',
  confirmPassword: ''
})

const validateForm = () => {
  let valid = true
  errors.username = ''
  errors.email = ''
  errors.password = ''
  errors.confirmPassword = ''

  if (!form.username || form.username.length < 3 || form.username.length > 50) {
    errors.username = '用户名长度在 3 到 50 个字符'
    valid = false
  }

  if (isRegister.value) {
    if (!form.email) {
      errors.email = '请输入邮箱'
      valid = false
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
      errors.email = '请输入正确的邮箱地址'
      valid = false
    }
  }

  if (!form.password || form.password.length < 6 || form.password.length > 50) {
    errors.password = '密码长度在 6 到 50 个字符'
    valid = false
  }

  if (isRegister.value) {
    if (!form.confirmPassword) {
      errors.confirmPassword = '请确认密码'
      valid = false
    } else if (form.confirmPassword !== form.password) {
      errors.confirmPassword = '两次输入的密码不一致'
      valid = false
    }
  }

  return valid
}

const handleSubmit = async () => {
  if (!validateForm()) return

  loading.value = true
  try {
    if (isRegister.value) {
      await userStore.register({
        username: form.username,
        email: form.email,
        nickname: form.nickname,
        password: form.password
      })
      toast.success('注册成功！请登录')
      isRegister.value = false
      form.password = ''
      form.confirmPassword = ''
    } else {
      await userStore.login({
        username: form.username,
        password: form.password
      })
      toast.success('登录成功')
      router.push('/')
    }
  } catch (error) {
    console.error('Login/Register error:', error)
  } finally {
    loading.value = false
  }
}

// 加载主题配置
onMounted(() => {
  loadTheme()
})
</script>

<style scoped>
.login-page {
  @apply min-h-screen flex items-center justify-center;
  background: var(--theme-bg-secondary);
}

.login-card {
  @apply w-full mx-lg;
  max-width: 400px;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md), var(--card-inset-shadow);
  background: var(--theme-bg-card);
  transition: box-shadow var(--transition-base);
}

.login-card:hover {
  box-shadow: var(--shadow-lg), var(--card-inset-shadow);
}

.input-wrapper {
  @apply relative;
}

.input-icon {
  @apply absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground pointer-events-none;
}

.form-footer {
  @apply text-center mt-sm;
}

.switch-link {
  @apply text-sm font-medium cursor-pointer bg-transparent border-none;
  color: var(--theme-primary);
  transition: color var(--transition-fast);
}

.switch-link:hover {
  color: var(--theme-primary-hover);
  text-decoration: underline;
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
