<template>
  <div class="user-center-layout">
    <el-container>
      <el-aside width="200px" class="sidebar">
        <div class="logo">
          <router-link to="/">InkSpace</router-link>
        </div>
        <el-menu
          :default-active="activeMenu"
          :router="true"
          class="menu"
        >
          <el-menu-item index="/dashboard">
            <el-icon><Odometer /></el-icon>
            <span>我的主页</span>
          </el-menu-item>
          <el-menu-item index="/dashboard/articles">
            <el-icon><Document /></el-icon>
            <span>我的文章</span>
          </el-menu-item>
          <el-menu-item index="/dashboard/works">
            <el-icon><Picture /></el-icon>
            <span>我的作品</span>
          </el-menu-item>
          <el-menu-item index="/dashboard/comments">
            <el-icon><ChatDotRound /></el-icon>
            <span>我的评论</span>
          </el-menu-item>
          <el-menu-item index="/favorites">
            <el-icon><Collection /></el-icon>
            <span>我的收藏</span>
          </el-menu-item>
          <el-menu-item index="/dashboard/notifications">
            <el-icon><Bell /></el-icon>
            <span>我的通知</span>
            <el-badge :value="unreadCount" :hidden="unreadCount === 0" class="notification-badge" />
          </el-menu-item>
          <el-menu-item index="/profile/edit">
            <el-icon><User /></el-icon>
            <span>个人设置</span>
          </el-menu-item>
        </el-menu>
      </el-aside>

      <el-container>
        <el-header class="header">
          <div class="header-content">
            <el-breadcrumb separator="/">
              <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
              <el-breadcrumb-item>{{ breadcrumbTitle }}</el-breadcrumb-item>
            </el-breadcrumb>
            <div class="header-actions">
              <el-button text @click="$router.push('/')">返回网站</el-button>
              <el-button v-if="userStore.isAdmin" type="primary" @click="goToAdminBackend">
                管理后台
              </el-button>
              <el-dropdown @command="handleCommand">
                <span class="user-info">
                  <el-avatar :size="32" :src="userStore.user?.avatar" />
                  <span>{{ userStore.user?.nickname || userStore.user?.username }}</span>
                </span>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="logout">
                      <el-icon><SwitchButton /></el-icon> 退出登录
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </div>
        </el-header>

        <el-main class="main">
          <RouterView />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, provide } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import api from '@/utils/api'
import {
  Odometer,
  Document,
  Picture,
  Collection,
  Bell,
  User,
  SwitchButton,
  ChatDotRound
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const unreadCount = ref(0)
const adminBackendUrl = ref('')

const loadUnreadCount = async () => {
  try {
    const response = await api.get('/notifications/unread-count')
    unreadCount.value = response.data.count || 0
  } catch (error) {
    console.error('Failed to load unread count:', error)
  }
}

// 提供更新未读数量的方法给子组件
const refreshUnreadCount = () => {
  loadUnreadCount()
}
provide('refreshUnreadCount', refreshUnreadCount)

// 每30秒刷新一次未读数量
let intervalId = null

onMounted(() => {
  loadUnreadCount()
  loadAdminBackendUrl()
  intervalId = setInterval(loadUnreadCount, 30000)
  
  // 监听路由变化，在通知页面时更频繁地刷新
  router.afterEach((to) => {
    if (to.path === '/dashboard/notifications') {
      loadUnreadCount()
    }
  })
})

onUnmounted(() => {
  if (intervalId) {
    clearInterval(intervalId)
  }
})

const activeMenu = computed(() => route.path)

const breadcrumbTitle = computed(() => {
  const titles = {
    '/dashboard': '我的主页',
    '/dashboard/articles': '我的文章',
    '/dashboard/articles/create': '写文章',
    '/dashboard/works': '我的作品',
    '/dashboard/works/create': '创建作品',
    '/dashboard/comments': '我的评论',
    '/favorites': '我的收藏',
    '/dashboard/notifications': '我的通知',
    '/profile/edit': '个人设置'
  }
  // 检查文章编辑路径
  if (route.path.includes('/dashboard/articles/') && route.path.includes('/edit')) {
    return '编辑文章'
  }
  // 检查作品编辑路径
  if (route.path.includes('/dashboard/works/') && route.path.includes('/edit')) {
    return '编辑作品'
  }
  return titles[route.path] || '用户中心'
})

const loadAdminBackendUrl = async () => {
  try {
    const response = await api.get('/settings/public')
    const settings = response.data || {}
    adminBackendUrl.value = settings.admin_backend_url || '/admin'
  } catch (error) {
    console.error('Failed to load admin backend URL:', error)
    adminBackendUrl.value = '/admin' // 默认值
  }
}

const goToAdminBackend = () => {
  const url = adminBackendUrl.value || '/admin'
  // 如果是绝对URL，使用 window.open，否则使用 router.push
  if (url.startsWith('http://') || url.startsWith('https://')) {
    window.open(url, '_blank')
  } else {
    router.push(url)
  }
}

const handleCommand = (command) => {
  if (command === 'logout') {
    userStore.logout()
    router.push('/')
  }
}
</script>

<style scoped>
.user-center-layout {
  min-height: 100vh;
  background-color: var(--theme-bg-secondary);
}

.sidebar {
  background-color: var(--theme-bg-card);
  box-shadow: 2px 0 8px var(--theme-shadow);
  border-right: 1px solid var(--theme-border);
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  font-weight: bold;
  border-bottom: 1px solid var(--theme-border);
  background-color: var(--theme-bg-card);
}

.logo a {
  font-size: 24px;
  font-weight: 600;
  background: linear-gradient(135deg, #667eea 0%, var(--theme-primary) 30%, #764ba2 70%, #f093fb 100%);
  background-size: 200% 200%;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-decoration: none;
  letter-spacing: 2px;
  position: relative;
  display: inline-block;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  font-style: italic;
  text-shadow: 0 0 30px rgba(102, 126, 234, 0.3);
  animation: gradientShift 3s ease infinite;
}

.logo a:hover {
  transform: translateY(-2px) scale(1.03);
  filter: brightness(1.15);
  letter-spacing: 3px;
}

@keyframes gradientShift {
  0%, 100% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
}

.menu {
  border-right: none;
  background-color: var(--theme-bg-card);
}

.header {
  background-color: var(--theme-bg-card);
  border-bottom: 1px solid var(--theme-border);
  padding: 0 20px;
  color: var(--theme-text-primary);
}

.header-content {
  height: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 15px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
}

.main {
  padding: 20px;
}

/* 修复通知气泡位置 */
:deep(.el-menu-item) {
  position: relative;
}

:deep(.el-menu-item .notification-badge) {
  position: absolute;
  right: 10px;
  top: 50%;
  transform: translateY(-50%);
  margin-top: 0;
}

:deep(.el-menu-item .notification-badge .el-badge__content) {
  position: static;
  transform: none;
}
</style>

