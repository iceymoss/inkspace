<template>
  <div class="admin-login">
    <div class="login-container">
      <Card class="login-card">
        <CardContent class="p-6">
          <div class="login-header">
            <h1>管理后台登录</h1>
            <p>Management System</p>
          </div>

          <form @submit.prevent="handleLogin">
            <div class="space-y-4">
              <div class="space-y-2">
                <div class="relative">
                  <UserIcon class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                  <Input
                    v-model="form.username"
                    placeholder="管理员账号"
                    class="pl-9 h-11"
                  />
                </div>
                <p v-if="errors.username" class="text-sm text-destructive">{{ errors.username }}</p>
              </div>

              <div class="space-y-2">
                <div class="relative">
                  <Lock class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                  <Input
                    v-model="form.password"
                    :type="showPassword ? 'text' : 'password'"
                    placeholder="密码"
                    class="pl-9 pr-9 h-11"
                  />
                  <button
                    type="button"
                    class="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
                    @click="showPassword = !showPassword"
                  >
                    <Eye v-if="!showPassword" class="h-4 w-4" />
                    <EyeOff v-else class="h-4 w-4" />
                  </button>
                </div>
                <p v-if="errors.password" class="text-sm text-destructive">{{ errors.password }}</p>
              </div>

              <Button
                type="submit"
                class="w-full h-11"
                :disabled="loading"
              >
                <Loader2 v-if="loading" class="mr-2 h-4 w-4 animate-spin" />
                登录
              </Button>
            </div>
          </form>

          <div class="login-footer">
            <router-link to="/" class="text-primary hover:underline">返回首页</router-link>
          </div>
        </CardContent>
      </Card>

      <div class="login-tips">
        <Alert>
          <AlertTitle>安全提示</AlertTitle>
          <AlertDescription>这是系统管理后台，仅限管理员访问。请妥善保管您的账号密码。</AlertDescription>
        </Alert>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { useAdminStore } from '@/stores/admin'
import { User as UserIcon, Lock, Eye, EyeOff, Loader2 } from 'lucide-vue-next'
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Alert, AlertTitle, AlertDescription } from '@/components/ui/alert'

const router = useRouter()
const adminStore = useAdminStore()
const loading = ref(false)
const showPassword = ref(false)

const form = reactive({
  username: '',
  password: ''
})

const errors = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [
    { required: true, message: '请输入管理员账号', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少6位', trigger: 'blur' }
  ]
}

const validate = () => {
  errors.username = ''
  errors.password = ''

  for (const rule of rules.username) {
    if (rule.required && !form.username) {
      errors.username = rule.message
      return false
    }
  }

  for (const rule of rules.password) {
    if (rule.required && !form.password) {
      errors.password = rule.message
      return false
    }
    if (rule.min && form.password.length < rule.min) {
      errors.password = rule.message
      return false
    }
  }

  return true
}

const handleLogin = async () => {
  if (!validate()) return

  loading.value = true
  try {
    await adminStore.login(form)
    toast.success('登录成功')
    router.push('/')
  } catch (error) {
    toast.error(error.message || '登录失败，请检查账号密码')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.admin-login {
  @apply min-h-screen flex items-center justify-center p-4;
  background: linear-gradient(135deg, var(--color-primary) 0%, var(--color-accent) 100%);
}

.login-container {
  @apply w-full max-w-[400px];
}

.login-card {
  @apply mb-4 rounded-lg;
  box-shadow: var(--shadow-lg);
}

.login-header {
  @apply text-center mb-6;
}

.login-header h1 {
  @apply m-0 mb-2 text-2xl font-bold;
  color: var(--color-text-primary);
}

.login-header p {
  @apply m-0 text-sm;
  color: var(--color-text-tertiary);
}

.login-footer {
  @apply text-center mt-4;
}

.login-tips {
  @apply text-center;
}
</style>
