import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'

const api = axios.create({
  baseURL: '/api',
  timeout: 30000
})

// Request interceptor
api.interceptors.request.use(
  (config) => {
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor
api.interceptors.response.use(
  (response) => {
    const data = response.data
    // 检查响应体中的 code 字段，只有 code === 0 才是成功
    if (data.code === 0) {
      // 返回整个响应对象，包含code, message, data
      return data
    } else {
      // 即使 HTTP 状态码是 200，如果 code 不是 0，也是错误
      ElMessage.error(data.message || '请求失败')
      return Promise.reject(new Error(data.message || '请求失败'))
    }
  },
  (error) => {
    if (error.response) {
      const status = error.response.status
      const data = error.response.data
      
      // 优先显示响应体中的 message，如果没有再显示默认消息
      const errorMessage = data?.message
      
      switch (status) {
        case 401:
          ElMessage.error(errorMessage || '未登录或登录已过期')
          const userStore = useUserStore()
          userStore.logout()
          window.location.href = '/login'
          break
        case 403:
          ElMessage.error(errorMessage || '没有权限')
          break
        case 404:
          ElMessage.error(errorMessage || '请求的资源不存在')
          break
        case 500:
          ElMessage.error(errorMessage || '服务器错误')
          break
        default:
          ElMessage.error(errorMessage || '请求失败')
      }
    } else {
      ElMessage.error('网络错误')
    }
    return Promise.reject(error)
  }
)

export default api

