import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import adminApi from '@/utils/adminApi'

export const useAdminStore = defineStore('admin', () => {
  const token = ref(localStorage.getItem('admin_token') || '')
  const admin = ref(null)

  const isLoggedIn = computed(() => !!token.value)

  function setToken(newToken) {
    token.value = newToken
    localStorage.setItem('admin_token', newToken)
  }

  function setAdmin(adminData) {
    admin.value = adminData
  }

  async function login(credentials) {
    try {
      // 调用管理员登录API
      const response = await adminApi.post('/admin/auth/login', credentials)
      
      // adminApi的响应拦截器已经返回了response.data
      // 所以response实际上就是后端返回的数据 { code: 0, data: { token, user } }
      const { token: authToken, user } = response.data
      
      // 检查是否是管理员
      if (user.role !== 'admin') {
        throw new Error('该账号不是管理员')
      }
      
      setToken(authToken)
      setAdmin(user)
      return response.data
    } catch (error) {
      throw error
    }
  }

  async function fetchProfile() {
    if (!token.value) return
    
    try {
      const response = await adminApi.get('/admin/auth/profile')
      setAdmin(response.data)
    } catch (error) {
      logout()
      throw error
    }
  }

  function logout() {
    token.value = ''
    admin.value = null
    localStorage.removeItem('admin_token')
  }

  return {
    token,
    admin,
    isLoggedIn,
    setToken,
    setAdmin,
    login,
    fetchProfile,
    logout
  }
})

