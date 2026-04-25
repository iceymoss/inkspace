<template>
  <div class="admin-layout">
    <aside class="sidebar">
      <div class="logo">
        <h3>后台管理</h3>
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
                  <router-link to="/">管理后台</router-link>
                </BreadcrumbLink>
              </BreadcrumbItem>
              <BreadcrumbSeparator />
              <BreadcrumbItem>
                <BreadcrumbLink>{{ breadcrumbTitle }}</BreadcrumbLink>
              </BreadcrumbItem>
            </BreadcrumbList>
          </Breadcrumb>
          <div class="header-actions">
            <DropdownMenu>
              <DropdownMenuTrigger as-child>
                <span class="user-info cursor-pointer">
                  <Avatar :size="32">
                    <AvatarImage :src="adminStore.admin?.avatar" />
                    <AvatarFallback>{{ (adminStore.admin?.nickname || adminStore.admin?.username || '?')[0] }}</AvatarFallback>
                  </Avatar>
                  <span>{{ adminStore.admin?.nickname || adminStore.admin?.username }}</span>
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

      <main class="main-content">
        <RouterView />
      </main>
    </div>

    <Toaster richColors />
  </div>
</template>

<script setup>
import { computed, onMounted, markRaw } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAdminStore } from '@/stores/admin'
import adminApi from '@/utils/adminApi'
import {
  Gauge,
  FileText,
  Image as ImageIcon,
  Folder,
  Tag,
  MessageCircle,
  Link as LinkIcon,
  Settings,
  User as UserIcon,
  Megaphone,
  LogOut
} from 'lucide-vue-next'
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'
import { DropdownMenu, DropdownMenuTrigger, DropdownMenuContent, DropdownMenuItem } from '@/components/ui/dropdown-menu'
import { Breadcrumb, BreadcrumbList, BreadcrumbItem, BreadcrumbLink, BreadcrumbSeparator } from '@/components/ui/breadcrumb'
import { Toaster } from '@/components/ui/toast'

const route = useRoute()
const router = useRouter()
const adminStore = useAdminStore()

const menuItems = [
  { path: '/', label: '控制台', icon: markRaw(Gauge) },
  { path: '/articles', label: '文章管理', icon: markRaw(FileText) },
  { path: '/works', label: '作品管理', icon: markRaw(ImageIcon) },
  { path: '/categories', label: '分类管理', icon: markRaw(Folder) },
  { path: '/tags', label: '标签管理', icon: markRaw(Tag) },
  { path: '/comments', label: '评论管理', icon: markRaw(MessageCircle) },
  { path: '/links', label: '友链管理', icon: markRaw(LinkIcon) },
  { path: '/settings', label: '系统配置', icon: markRaw(Settings) },
  { path: '/users', label: '用户管理', icon: markRaw(UserIcon) },
  { path: '/ads', label: '广告管理', icon: markRaw(Megaphone) },
]

const activeMenu = computed(() => route.path)

const breadcrumbTitle = computed(() => {
  const titles = {
    '/': '控制台',
    '/articles': '文章管理',
    '/works': '作品管理',
    '/categories': '分类管理',
    '/tags': '标签管理',
    '/comments': '评论管理',
    '/links': '友链管理',
    '/settings': '系统配置',
    '/users': '用户管理',
    '/ads': '广告管理'
  }
  return titles[route.path] || '管理'
})

const handleCommand = (command) => {
  if (command === 'logout') {
    adminStore.logout()
    router.push('/login')
  }
}

const loadSiteSettings = async () => {
  try {
    const response = await adminApi.get('/admin/settings')
    const settings = response.data || []
    const siteLogoSetting = settings.find(s => s.key === 'site_logo')
    
    if (siteLogoSetting && siteLogoSetting.value) {
      updateFavicon(siteLogoSetting.value)
    }
    
    const siteNameSetting = settings.find(s => s.key === 'site_name')
    if (siteNameSetting && siteNameSetting.value) {
      document.title = `${siteNameSetting.value} - 管理后台`
    }
  } catch (error) {
    console.error('Failed to load site settings:', error)
  }
}

const updateFavicon = (logoUrl) => {
  const oldFavicon = document.querySelector('link[rel="icon"]')
  if (oldFavicon) {
    oldFavicon.remove()
  }
  
  const link = document.createElement('link')
  link.rel = 'icon'
  link.type = 'image/png'
  link.href = logoUrl
  document.head.appendChild(link)
}

onMounted(() => {
  loadSiteSettings()
})
</script>

<style scoped>
.admin-layout {
  @apply min-h-screen flex;
}

.sidebar {
  @apply flex flex-col shrink-0;
  width: 200px;
  background-color: #1E1B4B;
}

.logo {
  @apply flex items-center justify-center;
  height: 60px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.logo h3 {
  @apply m-0 font-semibold tracking-wider;
  font-size: var(--font-size-lg);
  color: rgba(255, 255, 255, 0.95);
}

.menu {
  @apply flex flex-col py-2;
}

.menu-item {
  @apply flex items-center gap-2 px-4 py-2.5 text-sm font-medium cursor-pointer relative no-underline;
  color: rgba(255, 255, 255, 0.65);
  transition: all 150ms ease;
  border-left: 3px solid transparent;
}

.menu-item:hover {
  background-color: rgba(255, 255, 255, 0.08);
  color: rgba(255, 255, 255, 0.95);
}

.menu-item.active {
  background-color: rgba(255, 255, 255, 0.12);
  color: #fff;
  border-left-color: rgba(255, 255, 255, 0.8);
  font-weight: 600;
}

.menu-icon {
  @apply w-[18px] h-[18px] shrink-0;
}

.main-wrapper {
  @apply flex flex-col flex-1 min-w-0;
}

.header {
  @apply px-4;
  background-color: var(--color-bg-card);
  border-bottom: 1px solid var(--color-border-lighter);
  box-shadow: var(--shadow-sm);
  transition: background-color var(--transition-slow), border-color var(--transition-slow);
  height: 60px;
}

.header-content {
  @apply h-full flex justify-between items-center;
}

.header-actions {
  @apply flex items-center gap-4;
}

.user-info {
  @apply flex items-center gap-2 cursor-pointer;
  transition: opacity var(--transition-base);
}

.user-info:hover {
  @apply opacity-80;
}

.main-content {
  @apply p-4 flex-1;
  background: var(--color-bg-secondary);
  transition: background-color var(--transition-slow);
}
</style>
