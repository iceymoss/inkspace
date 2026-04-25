<template>
  <div class="user-center-layout">
    <aside class="sidebar">
      <div class="logo">
        <router-link to="/">InkSpace</router-link>
      </div>
      <nav class="menu">
        <router-link
          v-for="item in menuItems"
          :key="item.path"
          :to="item.path"
          class="menu-item"
          :class="{ active: activeMenu === item.path }"
        >
          <component :is="item.icon" class="menu-icon" />
          <span>{{ item.label }}</span>
          <Badge v-if="item.badge && unreadCount > 0" variant="destructive" class="ml-auto">{{ unreadCount }}</Badge>
        </router-link>
      </nav>
    </aside>

    <div class="main-wrapper">
      <header class="header">
        <div class="header-content">
          <Breadcrumb>
            <BreadcrumbList>
              <BreadcrumbItem>
                <BreadcrumbLink as-child>
                  <router-link to="/">首页</router-link>
                </BreadcrumbLink>
              </BreadcrumbItem>
              <BreadcrumbSeparator />
              <BreadcrumbItem>
                <BreadcrumbLink>{{ breadcrumbTitle }}</BreadcrumbLink>
              </BreadcrumbItem>
            </BreadcrumbList>
          </Breadcrumb>
          <div class="header-actions">
            <Button variant="ghost" @click="$router.push('/')">返回网站</Button>
            <Button v-if="userStore.isAdmin" variant="accent" @click="goToAdminBackend">
              管理后台
            </Button>
            <DropdownMenu>
              <DropdownMenuTrigger as-child>
                <span class="user-info cursor-pointer">
                  <Avatar :size="32">
                    <AvatarImage :src="userStore.user?.avatar" :alt="userStore.user?.nickname" />
                    <AvatarFallback>{{ (userStore.user?.nickname || userStore.user?.username || '?')[0] }}</AvatarFallback>
                  </Avatar>
                  <span>{{ userStore.user?.nickname || userStore.user?.username }}</span>
                </span>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" :side-offset="8">
                <DropdownMenuItem class="cursor-pointer" @click="handleCommand('logout')">
                  <LogOut class="mr-2 h-4 w-4" /> 退出登录
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </header>

      <main class="main">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, provide, markRaw } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { Gauge, FileText, Image, MessageCircle, Bookmark, Bell, User as UserIcon, Settings, LogOut } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { DropdownMenu, DropdownMenuTrigger, DropdownMenuContent, DropdownMenuItem } from '@/components/ui/dropdown-menu'
import { Breadcrumb, BreadcrumbList, BreadcrumbItem, BreadcrumbLink, BreadcrumbSeparator } from '@/components/ui/breadcrumb'
import api from '@/utils/api'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const unreadCount = ref(0)
const adminBackendUrl = ref('')

const menuItems = [
  { path: '/dashboard', label: '我的主页', icon: markRaw(Gauge) },
  { path: '/dashboard/articles', label: '我的文章', icon: markRaw(FileText) },
  { path: '/dashboard/works', label: '我的作品', icon: markRaw(Image) },
  { path: '/dashboard/comments', label: '我的评论', icon: markRaw(MessageCircle) },
  { path: '/favorites', label: '我的收藏', icon: markRaw(Bookmark) },
  { path: '/dashboard/notifications', label: '我的通知', icon: markRaw(Bell), badge: true },
  { path: '/profile/edit', label: '个人设置', icon: markRaw(Settings) },
]

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
  @apply min-h-screen flex;
  background-color: var(--theme-bg-secondary);
}

.sidebar {
  @apply flex flex-col shrink-0;
  width: 220px;
  background-color: var(--theme-bg-card);
  box-shadow: var(--shadow-md), var(--card-inset-shadow);
  border-right: 1px solid var(--theme-border);
  transition: background-color var(--transition-slow), border-color var(--transition-slow);
}

.logo {
  @apply flex items-center justify-center font-bold;
  height: 60px;
  font-size: var(--font-size-lg);
  border-bottom: 1px solid var(--theme-border);
  background-color: var(--theme-bg-card);
  transition: background-color var(--transition-slow), border-color var(--transition-slow);
}

.logo a {
  @apply font-serif font-bold cursor-pointer inline-block relative;
  font-size: 24px;
  letter-spacing: 1px;
  background: linear-gradient(135deg, var(--theme-primary) 0%, var(--theme-accent) 50%, var(--color-accent) 100%);
  background-size: 200% 200%;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-decoration: none;
  transition: all var(--transition-base);
  animation: gradientShift 3s ease infinite;
}

.logo a:hover {
  letter-spacing: 2px;
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
  @apply flex flex-col py-sm;
}

.menu-item {
  @apply flex items-center gap-sm px-md py-sm text-sm font-medium cursor-pointer relative;
  color: var(--theme-text-secondary);
  transition: background-color var(--transition-fast), color var(--transition-fast);
  text-decoration: none;
  border-left: 3px solid transparent;
}

.menu-item:hover {
  background-color: var(--theme-bg-hover);
  color: var(--theme-primary);
}

.menu-item.active {
  background-color: var(--theme-bg-hover);
  color: var(--theme-primary);
  border-left-color: var(--theme-primary);
  font-weight: 600;
}

.menu-icon {
  @apply w-[18px] h-[18px] shrink-0;
}

.main-wrapper {
  @apply flex flex-col flex-1 min-w-0;
}

.header {
  background-color: var(--theme-bg-card);
  border-bottom: 1px solid var(--theme-border);
  box-shadow: var(--shadow-sm);
  @apply px-md;
  color: var(--theme-text-primary);
  transition: background-color var(--transition-slow), border-color var(--transition-slow);
  height: 56px;
}

.header-content {
  @apply h-full flex justify-between items-center;
}

.header-actions {
  @apply flex items-center gap-md;
}

.user-info {
  @apply flex items-center gap-sm cursor-pointer;
  transition: opacity var(--transition-base);
}

.user-info:hover {
  @apply opacity-80;
}

.main {
  @apply p-md flex-1;
}
</style>
