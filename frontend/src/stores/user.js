import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/utils/api'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(null)

  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  function setToken(newToken) {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  function setUser(userData) {
    user.value = userData
  }

  async function login(credentials) {
    const response = await api.post('/login', credentials)
    setToken(response.data.token)
    setUser(response.data.user)
    return response.data
  }

  async function register(userData) {
    const response = await api.post('/register', userData)
    return response.data
  }

  async function fetchProfile() {
    if (!token.value) return
    try {
      const response = await api.get('/profile')
      setUser(response.data)
    } catch (error) {
      logout()
    }
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
  }

  return {
    token,
    user,
    isLoggedIn,
    isAdmin,
    setToken,
    setUser,
    login,
    register,
    fetchProfile,
    logout
  }
})

