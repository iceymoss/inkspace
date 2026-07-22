<template>
  <div
    class="user-center-layout"
    :class="{ 'sidebar-is-open': sidebarOpen }"
  >
    <el-container>
      <button
        v-if="sidebarOpen"
        class="sidebar-scrim"
        type="button"
        aria-label="关闭用户中心导航"
        @click="sidebarOpen = false"
      />
      <el-aside
        width="240px"
        class="sidebar"
        :class="{ 'is-open': sidebarOpen }"
      >
        <div class="logo">
          <router-link to="/">
            <template v-if="isTerminal"><span>$</span>inkspace.log</template>
            <template v-else>Ink<span>Space</span></template>
          </router-link>
          <small :class="{ serif: !isTerminal }">{{ isTerminal ? '~/dashboard' : 'MEMBER · 01' }}</small>
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
          <el-menu-item index="/dashboard/workspaces">
            <el-icon><Reading /></el-icon>
            <span>我的知识库</span>
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
            <el-badge
              :value="unreadCount"
              :hidden="unreadCount === 0"
              class="notification-badge"
            />
          </el-menu-item>
          <el-menu-item index="/profile/edit">
            <el-icon><User /></el-icon>
            <span>个人设置</span>
          </el-menu-item>
          <el-menu-item index="/dashboard/appearance">
            <el-icon><Brush /></el-icon>
            <span>外观与主题</span>
          </el-menu-item>
        </el-menu>
      </el-aside>

      <el-container>
        <el-header class="header">
          <div class="header-content">
            <div class="header-leading">
              <button
                class="sidebar-toggle"
                type="button"
                :aria-expanded="sidebarOpen"
                :aria-label="sidebarOpen ? '关闭用户中心导航' : '打开用户中心导航'"
                @click="sidebarOpen = !sidebarOpen"
              >
                <Menu />
              </button>
              <div class="page-marker">
                <span class="page-label">{{ isTerminal ? 'USER / CONTROL PANEL' : 'PERSONAL EDITION' }}</span>
                <el-breadcrumb separator="/">
                  <el-breadcrumb-item :to="{ path: '/' }">
                    首页
                  </el-breadcrumb-item>
                  <el-breadcrumb-item>{{ breadcrumbTitle }}</el-breadcrumb-item>
                </el-breadcrumb>
              </div>
            </div>
            <div class="header-actions">
              <el-button
                text
                @click="$router.push('/')"
              >
                返回网站
              </el-button>
              <el-button
                v-if="userStore.isAdmin"
                type="primary"
                @click="goToAdminBackend"
              >
                管理后台
              </el-button>
              <el-dropdown @command="handleCommand">
                <span class="user-info">
                  <el-avatar
                    :size="32"
                    :src="userStore.user?.avatar"
                  />
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
import { ref, computed, onMounted, onUnmounted, provide, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useAppearanceStore } from '@/stores/appearance'
import api from '@/utils/api'
import {
  Odometer,
  Document,
  Picture,
  Collection,
  Bell,
  User,
  SwitchButton,
  ChatDotRound,
  Reading,
  Brush,
  Menu
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const appearance = useAppearanceStore()
const isTerminal = computed(() => appearance.activePreference.ui_theme === 'terminal')
const unreadCount = ref(0)
const adminBackendUrl = ref('')
const sidebarOpen = ref(false)

watch(() => route.fullPath, () => {
  sidebarOpen.value = false
})

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

const activeMenu = computed(() => {
  if (route.path.startsWith('/dashboard/workspaces') || route.path.startsWith('/dashboard/docs/')) {
    return '/dashboard/workspaces'
  }
  return route.path
})

const breadcrumbTitle = computed(() => {
  const titles = {
    '/dashboard': '我的主页',
    '/dashboard/articles': '我的文章',
    '/dashboard/articles/create': '写文章',
    '/dashboard/workspaces': '我的知识库',
    '/dashboard/works': '我的作品',
    '/dashboard/works/create': '创建作品',
    '/dashboard/comments': '我的评论',
    '/favorites': '我的收藏',
    '/dashboard/notifications': '我的通知',
    '/dashboard/appearance': '外观与主题',
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
  if (route.path.startsWith('/dashboard/workspaces/')) {
    return '知识库空间'
  }
  if (route.path.startsWith('/dashboard/docs/')) {
    return '编辑知识库文档'
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
  min-height: 100dvh;
  background-color: var(--theme-bg-secondary);
}

.user-center-layout > .el-container {
  min-height: 100vh;
  min-height: 100dvh;
  align-items: stretch;
}

.user-center-layout > .el-container > .el-container {
  min-width: 0;
  min-height: 100vh;
  min-height: 100dvh;
}

.sidebar {
  position: sticky;
  top: 0;
  align-self: flex-start;
  height: 100vh;
  height: 100dvh;
  background-color: var(--theme-bg-card);
  border-right: 1px solid var(--theme-border);
  overflow: hidden;
}

.logo {
  height: 96px;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  flex-direction: column;
  padding: 0 28px;
  font-size: 20px;
  font-weight: bold;
  border-bottom: 1px solid var(--theme-border);
  background-color: var(--theme-bg-card);
}

.logo a {
  color: var(--theme-text-primary);
  font-family: Georgia, 'Songti SC', 'Noto Serif SC', SimSun, serif;
  font-size: 23px;
  font-weight: 600;
  text-decoration: none;
  letter-spacing: .1em;
  display: inline-block;
}

.logo a span {
  color: var(--theme-primary);
}

.logo small {
  margin-top: 3px;
  color: var(--theme-primary);
  font-size: 9px;
  font-weight: 400;
  letter-spacing: .2em;
}

.menu {
  height: calc(100vh - 96px);
  height: calc(100dvh - 96px);
  border-right: none;
  background-color: var(--theme-bg-card);
  overflow-y: auto;
  overflow-x: hidden;
}

.header {
  height: 72px;
  background-color: var(--theme-bg-card);
  border-bottom: 1px solid var(--theme-border);
  padding: 0 32px;
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

.header-leading {
  display: flex;
  align-items: center;
  min-width: 0;
}

.page-marker {
  display: grid;
  gap: 2px;
}

.page-label {
  color: var(--theme-primary);
  font-family: Georgia, serif;
  font-size: 9px;
  letter-spacing: .24em;
}

.sidebar-toggle {
  display: none;
  place-items: center;
  width: 36px;
  height: 36px;
  margin-right: 16px;
  padding: 8px;
  border: 1px solid var(--theme-border);
  border-radius: 50%;
  background: transparent;
  color: var(--theme-text-primary);
  cursor: pointer;
}

.sidebar-toggle svg {
  width: 100%;
  height: 100%;
}

.sidebar-scrim {
  display: none;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
}

.main {
  padding: 34px clamp(20px, 4vw, 56px) 56px;
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

@media (max-width: 900px) {
  .sidebar {
    position: fixed;
    z-index: 1100;
    left: 0;
    width: min(280px, calc(100vw - 48px)) !important;
    height: 100vh;
    height: 100dvh;
    transform: translateX(-100%);
    transition: transform .25s ease;
  }

  .sidebar.is-open {
    transform: none;
  }

  .sidebar-scrim {
    position: fixed;
    z-index: 1090;
    inset: 0;
    display: block;
    border: 0;
    background: color-mix(in srgb, var(--theme-text-primary) 38%, transparent);
    cursor: pointer;
  }

  .sidebar-toggle {
    display: grid;
  }

  .header {
    padding-inline: 20px;
  }

  .main {
    padding: 28px 20px 48px;
  }
}

@media (max-width: 560px) {
  .header {
    height: 62px;
    padding: 0 12px;
  }

  .header-actions {
    gap: 4px;
  }

  .header-actions > .el-button,
  .user-info > span {
    display: none;
  }

  .page-label,
  :deep(.el-breadcrumb__item:first-child) {
    display: none;
  }

  .sidebar-toggle {
    margin-right: 10px;
  }

  .main {
    padding: 20px 12px 40px;
  }
}
</style>
