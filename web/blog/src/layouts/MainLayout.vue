<template>
  <div class="main-layout">
    <header class="header">
      <div class="container-blog">
        <div class="header-content">
          <div class="logo">
            <router-link to="/">InkSpace</router-link>
          </div>
          <nav class="nav">
            <router-link to="/">首页</router-link>
            <router-link to="/blog">博客</router-link>
            <router-link to="/works">作品</router-link>
            <router-link to="/links">友链</router-link>
            <router-link to="/about">关于</router-link>
          </nav>
          <div class="header-actions">
            <template v-if="userStore.isLoggedIn">
              <NotificationDropdown />
              <DropdownMenu>
                <DropdownMenuTrigger as-child>
                  <span class="user-info cursor-pointer">
                    <Avatar :size="32">
                      <AvatarImage :src="userStore.user?.avatar" :alt="userStore.user?.nickname" />
                      <AvatarFallback>{{ (userStore.user?.nickname || userStore.user?.username || '?')[0] }}</AvatarFallback>
                    </Avatar>
                    <span class="username">{{ userStore.user?.nickname || userStore.user?.username }}</span>
                  </span>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" :side-offset="8">
                  <DropdownMenuItem class="cursor-pointer" @click="handleCommand('profile')">
                    <User class="mr-2 h-4 w-4" /> 个人主页
                  </DropdownMenuItem>
                  <DropdownMenuItem class="cursor-pointer" @click="handleCommand('dashboard')">
                    <Gauge class="mr-2 h-4 w-4" /> 用户中心
                  </DropdownMenuItem>
                  <DropdownMenuItem class="cursor-pointer" @click="handleCommand('edit')">
                    <Pencil class="mr-2 h-4 w-4" /> 编辑资料
                  </DropdownMenuItem>
                  <DropdownMenuItem class="cursor-pointer" @click="handleCommand('favorites')">
                    <Bookmark class="mr-2 h-4 w-4" /> 我的收藏
                  </DropdownMenuItem>
                  <DropdownMenuItem class="cursor-pointer" @click="handleCommand('searchUser')">
                    <Search class="mr-2 h-4 w-4" /> 搜索用户
                  </DropdownMenuItem>
                  <Separator class="my-1" />
                  <DropdownMenuItem class="cursor-pointer" @click="handleCommand('logout')">
                    <LogOut class="mr-2 h-4 w-4" /> 退出登录
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </template>
            <template v-else>
              <Button @click="$router.push('/login')">登录</Button>
            </template>
          </div>
        </div>
      </div>
    </header>

    <main class="main-content">
      <RouterView />
    </main>

    <Toaster richColors theme="system" />

    <footer class="footer">
      <div class="container-blog">
        <p>{{ siteSettings.site_copyright || '© 2024 InkSpace. All rights reserved.' }}</p>
        <p>{{ siteSettings.site_description || 'Powered by Go + Gin + Vue' }}</p>
        <p v-if="siteSettings.site_icp">
          <a :href="`https://beian.miit.gov.cn/`" target="_blank" rel="noopener">
            {{ siteSettings.site_icp }}
          </a>
        </p>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { useRouter } from 'vue-router'
import { User, Pencil, Bookmark, LogOut, Gauge, Search } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'
import { DropdownMenu, DropdownMenuTrigger, DropdownMenuContent, DropdownMenuItem } from '@/components/ui/dropdown-menu'
import { Separator } from '@/components/ui/separator'
import { Toaster } from 'vue-sonner'
import NotificationDropdown from '@/components/NotificationDropdown.vue'
import api from '@/utils/api'

const userStore = useUserStore()
const router = useRouter()
const siteSettings = ref({})

// Load user profile if logged in
if (userStore.isLoggedIn && !userStore.user) {
  userStore.fetchProfile()
}

// Load public settings
const loadSiteSettings = async () => {
  try {
    const response = await api.get('/settings/public')
    siteSettings.value = response.data || {}
    
    // 更新页面标题
    if (siteSettings.value.site_name) {
      document.title = siteSettings.value.site_name
    }
    
    // 更新favicon
    if (siteSettings.value.site_logo) {
      updateFavicon(siteSettings.value.site_logo)
    }
  } catch (error) {
    console.error('Failed to load site settings:', error)
  }
}

// 更新favicon
const updateFavicon = (logoUrl) => {
  // 移除旧的favicon
  const oldFavicon = document.querySelector('link[rel="icon"]')
  if (oldFavicon) {
    oldFavicon.remove()
  }
  
  // 创建新的favicon
  const link = document.createElement('link')
  link.rel = 'icon'
  link.type = 'image/png'
  link.href = logoUrl
  document.head.appendChild(link)
}

const handleCommand = (command) => {
  if (command === 'logout') {
    userStore.logout()
    router.push('/')
  } else if (command === 'dashboard') {
    router.push('/dashboard')
  } else if (command === 'profile') {
    const currentUserId = userStore.user?.id
    if (!currentUserId) {
      console.error('User ID not found')
      return
    }
    const targetPath = `/users/${currentUserId}`
    const currentPath = router.currentRoute.value.path
    
    // 如果当前路由就是目标路由，使用 replace 强制刷新
    // 否则使用 push 正常跳转
    if (currentPath === targetPath) {
      // 先跳转到首页，然后立即跳转到目标路径，强制组件重新加载
      router.replace('/').then(() => {
        setTimeout(() => {
          router.replace(targetPath)
        }, 50)
      })
    } else {
      // 正常跳转，watch 会监听到路由变化并重新加载数据
      router.push(targetPath)
    }
  } else if (command === 'edit') {
    router.push('/profile/edit')
  } else if (command === 'favorites') {
    router.push('/favorites')
  } else if (command === 'searchUser') {
    router.push('/user-search')
  }
}

onMounted(() => {
  loadSiteSettings()
})
</script>

<style scoped>
.main-layout {
  @apply min-h-screen flex flex-col;
}

.header {
  background: var(--theme-bg-card);
  box-shadow: var(--shadow-sm), var(--card-inset-shadow);
  border-bottom: 1px solid var(--theme-border);
  @apply sticky top-0 z-[1000];
  color: var(--theme-text-primary);
  transition: background-color var(--transition-slow), border-color var(--transition-slow), box-shadow var(--transition-slow);
}

.header-content {
  @apply flex items-center justify-between;
  height: 64px;
}

.logo a {
  @apply font-serif font-bold cursor-pointer inline-block relative;
  font-size: 28px;
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

.nav {
  @apply flex;
  gap: var(--spacing-lg);
}

.nav a {
  color: var(--theme-text-secondary);
  @apply font-medium cursor-pointer relative;
  font-size: var(--font-size-base);
  transition: color var(--transition-base);
  text-decoration: none;
  padding: var(--spacing-xs) 0;
}

.nav a::after {
  content: '';
  @apply absolute left-0 w-0 h-[2px];
  bottom: -2px;
  background: var(--theme-primary);
  border-radius: var(--radius-full);
  transition: width var(--transition-base);
}

.nav a:hover::after,
.nav a.router-link-active::after {
  @apply w-full;
}

.nav a:hover,
.nav a.router-link-active {
  color: var(--theme-primary);
}

.header-actions {
  @apply flex items-center;
  gap: var(--spacing-md);
}

.user-info {
  @apply flex items-center;
  gap: var(--spacing-sm);
  transition: opacity var(--transition-base);
}

.user-info:hover {
  @apply opacity-80;
}

.username {
  color: var(--theme-text-primary);
  @apply font-medium;
  font-size: var(--font-size-sm);
}

.main-content {
  @apply flex-1;
  padding: var(--spacing-sm) 0 var(--spacing-xl);
}

.footer {
  background: var(--color-text-primary);
  border-top: 1px solid var(--theme-border);
  color: rgba(255, 255, 255, 0.8);
  @apply py-lg text-center mt-auto;
  transition: background-color var(--transition-slow), border-color var(--transition-slow);
}

.footer p {
  color: rgba(255, 255, 255, 0.7);
  @apply leading-base;
  font-size: var(--font-size-sm);
  margin-bottom: var(--spacing-xs);
}

.footer a {
  color: rgba(255, 255, 255, 0.8);
  @apply cursor-pointer;
  transition: color var(--transition-base);
}

.footer a:hover {
  @apply text-white;
}

body.theme-mourning .footer {
  background: var(--theme-bg-card);
  border-top: 1px solid var(--theme-border);
  color: var(--theme-text-primary);
}

body.theme-mourning .footer p {
  color: var(--theme-text-secondary);
}

body.theme-mourning .footer a {
  color: var(--theme-text-primary);
}

body.theme-mourning .footer a:hover {
  color: var(--theme-primary);
}

@media (max-width: 768px) {
  .nav {
    gap: var(--spacing-md);
  }

  .nav a {
    font-size: var(--font-size-sm);
  }

  .logo a {
    font-size: 20px;
  }
}
</style>
