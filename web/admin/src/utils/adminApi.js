import axios from 'axios'
import { toast } from 'vue-sonner'

const adminApi = axios.create({
  baseURL: '/api',  // 通过vite proxy转发到8083
  timeout: 30000
})

// Request interceptor
adminApi.interceptors.request.use(
  (config) => {
    const adminToken = localStorage.getItem('admin_token')
    if (adminToken) {
      config.headers.Authorization = `Bearer ${adminToken}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor
adminApi.interceptors.response.use(
  (response) => {
    const data = response.data
    if (data.code === 0 || response.status === 200) {
      return data
    } else {
      toast.error(data.message || '请求失败')
      return Promise.reject(new Error(data.message || '请求失败'))
    }
  },
  (error) => {
    if (error.response) {
      const status = error.response.status
      switch (status) {
        case 401:
          toast.error('未登录或登录已过期')
          localStorage.removeItem('admin_token')
          window.location.href = '/admin/login'
          break
        case 403:
          toast.error('没有权限')
          break
        case 404:
          toast.error('请求的资源不存在')
          break
        case 500:
          toast.error('服务器错误')
          break
        default:
          toast.error(error.response.data?.message || '请求失败')
      }
    } else {
      toast.error('网络错误')
    }
    return Promise.reject(error)
  }
)

export default adminApi

