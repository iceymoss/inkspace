<template>
  <div class="admin-login">
    <div class="login-container">
      <el-card class="login-card">
        <div class="login-header">
          <h1>管理后台登录</h1>
          <p>Management System</p>
        </div>

        <el-form
          ref="formRef"
          :model="form"
          :rules="rules"
          @keyup.enter="handleLogin"
        >
          <el-form-item prop="username">
            <el-input
              v-model="form.username"
              placeholder="管理员账号"
              size="large"
              prefix-icon="User"
            />
          </el-form-item>

          <el-form-item prop="password">
            <el-input
              v-model="form.password"
              type="password"
              placeholder="密码"
              size="large"
              prefix-icon="Lock"
              show-password
            />
          </el-form-item>

          <el-form-item>
            <el-button
              type="primary"
              size="large"
              :loading="loading"
              @click="handleLogin"
              style="width: 100%"
            >
              登录
            </el-button>
          </el-form-item>
        </el-form>

        <div class="login-footer">
          <el-link type="primary" @click="goToHome">返回首页</el-link>
        </div>
      </el-card>

      <div class="login-tips">
        <el-alert
          title="安全提示"
          type="warning"
          :closable="false"
          description="这是系统管理后台，仅限管理员访问。请妥善保管您的账号密码。"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAdminStore } from '@/stores/admin'

const router = useRouter()
const adminStore = useAdminStore()
const formRef = ref(null)
const loading = ref(false)

const form = reactive({
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

const handleLogin = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      await adminStore.login(form)
      ElMessage.success('登录成功')
      router.push('/admin')
    } catch (error) {
      ElMessage.error(error.message || '登录失败，请检查账号密码')
    } finally {
      loading.value = false
    }
  })
}

const goToHome = () => {
  router.push('/')
}
</script>

<style scoped>
.admin-login {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.login-container {
  width: 100%;
  max-width: 400px;
}

.login-card {
  margin-bottom: 20px;
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-header h1 {
  margin: 0 0 10px 0;
  font-size: 28px;
  color: #303133;
}

.login-header p {
  margin: 0;
  color: #909399;
  font-size: 14px;
}

.login-footer {
  text-align: center;
  margin-top: 20px;
}

.login-tips {
  text-align: center;
}

:deep(.el-card__body) {
  padding: 40px;
}

:deep(.el-form-item) {
  margin-bottom: 24px;
}

:deep(.el-alert) {
  background-color: rgba(255, 255, 255, 0.9);
}
</style>

