<template>
  <main class="bootstrap-page">
    <el-card class="bootstrap-card">
      <h1>初始化管理后台</h1>
      <div v-if="loading">正在检查系统状态…</div>
      <div v-else-if="statusError" role="alert">系统暂不可用。无法确认系统初始化状态，请稍后联系维护人员。</div>
      <div v-else-if="status?.state === 'unavailable'" role="alert">系统暂不可用。系统状态异常，无法通过此页面恢复。请联系维护人员。</div>
      <div v-else-if="status?.state === 'initialized'" role="status">系统已完成初始化。请使用管理员账号登录。</div>
      <div v-else-if="status?.state === 'pending' && !status.secret_configured" role="alert">初始化暂不可用。服务器尚未配置初始化密钥，请联系维护人员。</div>
      <el-form v-else-if="status?.state === 'pending'" ref="formRef" :model="form" :rules="rules" label-position="top" @submit.prevent="submit">
        <el-form-item label="初始化密钥" prop="secret"><el-input v-model="form.secret" type="password" autocomplete="off" show-password /></el-form-item>
        <el-form-item label="管理员用户名" prop="username"><el-input v-model="form.username" autocomplete="username" /></el-form-item>
        <el-form-item label="邮箱" prop="email"><el-input v-model="form.email" autocomplete="email" /></el-form-item>
        <el-form-item label="密码" prop="password"><el-input v-model="form.password" type="password" autocomplete="new-password" show-password /><span v-if="passwordMessage" class="validation-error">{{ passwordMessage }}</span></el-form-item>
        <el-form-item label="确认密码" prop="confirm_password"><el-input v-model="form.confirm_password" type="password" autocomplete="new-password" show-password /><span v-if="confirmationMessage" class="validation-error">{{ confirmationMessage }}</span></el-form-item>
        <div v-if="submitError" class="form-error" role="alert">{{ submitError }}</div>
        <el-button native-type="submit" type="primary" :loading="submitting">创建管理员</el-button>
      </el-form>
      <el-button v-if="status?.state === 'initialized'" type="primary" @click="router.push('/login')">前往登录</el-button>
    </el-card>
  </main>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import adminApi from '@/utils/adminApi'

const router = useRouter()
const status = ref(null)
const loading = ref(true)
const statusError = ref(false)
const submitting = ref(false)
const submitError = ref('')
const formRef = ref(null)
const form = reactive({ secret: '', username: '', email: '', password: '', confirm_password: '' })
const passwordMessage = ref('')
const confirmationMessage = ref('')
const validatePasswordBytes = (_rule, value, callback) => {
  if (new TextEncoder().encode(value).length > 72) {
    passwordMessage.value = '密码不能超过72个字节'
    callback(new Error(passwordMessage.value))
  } else {
    passwordMessage.value = ''
    callback()
  }
}
const validateConfirmation = (_rule, value, callback) => {
  if (value !== form.password) {
    confirmationMessage.value = '两次输入的密码不一致'
    callback(new Error(confirmationMessage.value))
  } else {
    confirmationMessage.value = ''
    callback()
  }
}
const rules = {
  secret: [{ required: true, message: '请输入初始化密钥', trigger: 'blur' }],
  username: [{ required: true, message: '请输入管理员用户名', trigger: 'blur' }, { min: 3, max: 50, message: '用户名长度为3到50位', trigger: 'blur' }],
  email: [{ required: true, type: 'email', message: '请输入有效邮箱', trigger: 'blur' }, { max: 100, message: '邮箱长度不能超过100位', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }, { min: 6, max: 50, message: '密码长度为6到50位', trigger: 'blur' }, { validator: validatePasswordBytes, trigger: 'blur' }],
  confirm_password: [{ required: true, message: '请确认密码', trigger: 'blur' }, { validator: validateConfirmation, trigger: 'blur' }]
}

onMounted(async () => {
  try {
    const response = await adminApi.get('/admin/auth/bootstrap/status', { skipAuth: true, suppressGlobalError: true })
    if (['pending', 'initialized', 'unavailable'].includes(response.data?.state)) status.value = response.data
    else statusError.value = true
  } catch {
    statusError.value = true
  } finally {
    loading.value = false
  }
})

async function submit() {
  if (!formRef.value || submitting.value) return
  try {
    await formRef.value.validate()
    submitting.value = true
    submitError.value = ''
    passwordMessage.value = ''
    confirmationMessage.value = ''
    await adminApi.post('/admin/auth/bootstrap', { ...form }, { skipAuth: true, suppressGlobalError: true })
    await router.push('/login')
  } catch (error) {
    if (error?.message === '系统已初始化，请登录') {
      await router.push('/login')
      return
    }
    if (error?.message && error.message !== 'invalid') {
      const safeMessages = ['用户名或邮箱已被占用，请更换', '初始化密钥错误或未配置', '系统已初始化，请登录', '系统初始化不可用，请联系维护者', '密码不能超过72个字节', '请检查密钥、用户名、邮箱、密码及确认密码', '请检查初始化信息']
      submitError.value = safeMessages.includes(error.message) ? error.message : '提交失败，请稍后重试或联系维护人员。'
    }
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.bootstrap-page { min-height: 100vh; display: grid; place-items: center; padding: 24px; background: linear-gradient(135deg, #667eea, #764ba2); }
.bootstrap-card { width: min(100%, 480px); }
h1 { margin: 0 0 24px; text-align: center; color: #303133; }
.bootstrap-card > div { margin: 12px 0; color: #606266; }
.validation-error { display: block; color: #f56c6c; font-size: 12px; }
.form-error { margin: 0 0 18px; color: #f56c6c; }
.el-button { width: 100%; }
</style>
